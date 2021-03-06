package client

import (
	"bufio"
	"codec"
	"context"
	"errors"
	"ext"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/alecthomas/log4go"
	"github.com/gogf/gf/container/gpool"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/os/gmutex"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/robfig/cron"
	"io"
	"net"
	"os"
	"proc"
	"screenshot"
	"sync"
	"syscall"
	"time"
)

const (
	TCP                      = "tcp"
	PACKAGE_TAG       string = "CMD"
	SEND_PACKET       string = "SEND_PACKET"
	RECIVE_PACKET     string = "RECIVE_PACKET"
	CAPTURE_PAGE_SIZE        = 1024
)

type None = struct{}
type Client struct {
	//日志
	logger *log4go.Logger
	//协议池
	cmdPool *gpool.Pool
	//编解码器
	packCodec *codec.Codec
	//截图对象
	capt *screenshot.Screenshot
	//截图id
	captId uint32
	//消息
	pubSub    *ext.ExtGoChanel
	//是否连接
	connected *gtype.Bool
}

func NewClient(logger *log4go.Logger, ps *ext.ExtGoChanel) *Client {
	return &Client{logger: logger, cmdPool: gpool.New(0, func() (interface{}, error) {
		ret := new(proc.Cmd)
		ret.Head = new(proc.Head)
		ret.Content = new(proc.Content)
		return ret, nil
	}, nil), capt: screenshot.NewScreenshotDefault(), pubSub: ps}
}

func (this *Client) Start(ctx context.Context, ret *proc.Broadcast) error {
	var (
		remotAddr *net.TCPAddr
		tcpCon    *net.TCPConn
		tcpErr    error
		wgroup    sync.WaitGroup
	)
	cron_task := cron.New()
	this.connected = gtype.NewBool(false)
	remotAddr, tcpErr = net.ResolveTCPAddr(TCP, fmt.Sprintf("%s:%d", *ret.Body.Server, *ret.Body.Port))
	if nil == tcpErr {
		tcpCon, tcpErr = net.DialTCP(TCP, nil, remotAddr)
		if nil == tcpErr {
			this.logger.Info("服务器连接成功:%s", tcpCon.RemoteAddr().String())
			this.connected.Set(true)
			//编解码器
			this.packCodec = codec.NewCodec(tcpCon)
			//创建接收、发送队列
			writeMessage, writeStatus := this.pubSub.Subscribe(context.Background(), SEND_PACKET)
			readMessage, readStatus := this.pubSub.Subscribe(context.Background(), RECIVE_PACKET)
			//控制消息队列
			connectCtx, connectCancelCtx := context.WithCancel(ctx)
			if nil == writeStatus && nil == readStatus {
				//创建锁保证TCP流顺序发送(保证各个操作之间不相互影响，不然影响服务器收包)
				syncLock := gmutex.New()
				cron_task.AddFunc("*/1 * * * *", func() {
					//每1秒发送心跳协议
					this.ping(syncLock)
				})
				//读取接收到的包并发布到消息队列（读取STREAM到MESSAGE）
				go func() {
					wgroup.Add(1)
					for this.connected.Val() {
						select {
						case <-connectCtx.Done():
							goto EXIT_REC
						case <-ctx.Done():
							goto EXIT_REC
						default:
							//读取应用协议包s
							if rd, e := this.packCodec.Read(); nil == e {
								if len(rd) > 0 {
									if e = this.pubSub.Publish(RECIVE_PACKET, message.NewMessage(watermill.NewUUID(), rd)); nil != e {
										this.logger.Warn(e)
									}
								}
							} else if e == io.EOF || syscall.EINVAL == e {
								this.logger.Warn(e)
								connectCancelCtx()
							} else {
								this.logger.Warn(e)
								//process
								if _, ok := e.(*net.OpError); ok {
									connectCancelCtx()
								}
							}
						}
					}
					//失败信息
				EXIT_REC:
					this.connected.Set(false)
					wgroup.Done()
				}()
				//处理消息队列（写入）并发送至服务器
				go func() {
					wgroup.Add(1)
					for {
						select {
						case <-connectCtx.Done():
							goto EXIT_W
						case <-ctx.Done():
							goto EXIT_W
						case msg := <-writeMessage:
							if nil != msg {
								e := this.packCodec.Write(msg.Payload)
								msg.Ack()
								if nil != e {
									if syscall.EINVAL == e {
										this.logger.Warn(e)
										connectCancelCtx()
									} else {
										this.logger.Warn(e)
										//process
										if _, ok := e.(*net.OpError); ok {
											connectCancelCtx()
										}
									}
								}
							}
						}
					}
				EXIT_W:
					this.connected.Set(false)
					wgroup.Done()
				}()
				//处理接收到的队列消息
				go func() {
					wgroup.Add(1)
					for this.connected.Val() {
						select {
						case <-connectCtx.Done():
							goto EXIT_THIS
						case <-ctx.Done():
							goto EXIT_THIS
						case msg := <-readMessage:
							if nil != msg {
								//处理数据包(接收的应用数据包)
								if cmd, err := this.cmdPool.Get(); nil == err {
									if nil != cmd {
										cd := cmd.(*proc.Cmd)
										if err := proto.Unmarshal(msg.Payload, cd); nil == err {
											switch *cd.Content.Type {
											case proc.ContentType_CLOSE:
												this.logger.Info("服务器已中断连接(%s)", *cd.Content.GetClose().Reason)
												connectCancelCtx()
												break
											case proc.ContentType_CAPTURE:
												//服务器请求马上截图
												this.capture(syncLock)
												break
											}
										}
									}
								}
								msg.Ack()
							}
						}
					}
				EXIT_THIS:
					this.connected.Set(false)
					wgroup.Done()
				}()
				//启动时自动更新一次最新信息
				this.capture(syncLock)
				cron_task.Start()
				select {
				case <-ctx.Done():
					this.logger.Info("客户端已被用户终止(01)")
				case <-connectCtx.Done():
					this.logger.Info("客户端已被服务器终止(02)")
				}
				wgroup.Wait()
				cron_task.Stop()
				return errors.New("客户端关闭,将自动重连,请耐心等候...")
			} else {
				this.logger.Fine(writeStatus)
			}
			if e := this.packCodec.Close(); nil != e {
				this.logger.Error(e)
			}
		} else {
			this.logger.Info(tcpErr)
		}
	} else {
		this.logger.Info(tcpErr)
	}
	return nil
}

var (
	source proc.Source = proc.Source_CLIENT
	cmd    string      = PACKAGE_TAG
)

/** 心跳协议 */
func (this *Client) ping(lock *gmutex.Mutex) {
	if this.connected.Val() {
		lock.TryLockFunc(func() {
			var (
				hostName    string           = ""
				contentType proc.ContentType = proc.ContentType_PING
				unixTime    int64            = time.Now().Unix()
			)
			contentPing := &proc.Content_Ping{}
			contentPing.Ping = &proc.Ping{}
			//获取主机名
			if h, e := os.Hostname(); nil == e {
				hostName = h
			}
			contentPing.Ping.Name = &hostName
			contentPing.Ping.Time = &unixTime
			//每5秒发送一个心跳包
			if ret, err := this.cmdPool.Get(); nil == err {
				if ins, ok := ret.(*proc.Cmd); ok {
					ins.Content.Source = &source
					ins.Head.Cmd = &cmd
					ins.Content.Type = &contentType
					ins.Content.Param = contentPing
					if data, err := proto.Marshal(ins); nil == err {
						if err = this.pubSub.Publish(SEND_PACKET, message.NewMessage(watermill.NewUUID(), data)); nil != err {
							this.logger.Warn(err)
						}
					} else {
						this.logger.Error(err)
					}
					this.cmdPool.Put(ins)
				}
			}
		})
	}
}

/** 发送本机基本信息 */
func (this *Client) info(lock *gmutex.Mutex) {
	if lock.TryLock() {
		//缓存信息
		lock.Unlock()
	}
}

/** 捕捉屏幕截图（可能比较耗时因此加入锁） */
func (this *Client) capture(lock *gmutex.Mutex) {
	if this.connected.Val() {
		var (
			contentType proc.ContentType = proc.ContentType_CAPTURE
			file        *os.File
			err         error
		)
		//截图(裁剪图像)比较耗时，不该放到锁里，否则造成其他包无法投递，尤其是心跳包，服务器直接关掉客户端就麻烦了
		if nil == this.capt.Capture() {
			if err = this.capt.Resize(1920, 0, screenshot.Lanczos3, 90); nil != err {
				this.logger.Warn(err)
			}
			lock.TryLockFunc(func() {
				this.captId += 1
				this.captId %= 65536
				//屏幕截图数据太大，需要多次发送 CAPTURE_PAGE_SIZE
				if file, err = os.Open(this.capt.GetPath()); nil == err {
					defer func() {
						file.Close()
					}()
					//内部函数实现发送截图数据到队列(尽量减少锁的时间)
					pageSend := func(d []byte, seq uint32, more, compress bool) {
						if ret, err := this.cmdPool.Get(); nil == err {
							if ins, ok := ret.(*proc.Cmd); ok {
								contentCapture := &proc.Content_Capture{}
								contentCapture.Capture = &proc.Capture{}
								contentCapture.Capture.Compress = &compress
								contentCapture.Capture.More = &more
								contentCapture.Capture.Id = &this.captId
								contentCapture.Capture.Data = d
								contentCapture.Capture.Seq = &seq
								contentCapture.Capture.Compress = &compress
								ins.Content.Source = &source
								ins.Content.Type = &contentType
								ins.Content.Param = contentCapture
								ins.Head.Cmd = &cmd
								if data, errx := proto.Marshal(ins); nil == errx {
									if errx = this.pubSub.Publish(SEND_PACKET, message.NewMessage(watermill.NewUUID(), data)); nil != errx {
										this.logger.Warn(errx)
									}
								} else {
									this.logger.Error(errx)
								}
								this.cmdPool.Put(ret)
							}
						}
					}
					raw := make([]byte, CAPTURE_PAGE_SIZE)
					buffReader := bufio.NewReader(file)
					var seq uint32 = 0
					pageSend([]byte{}, seq, false, false)
					//空包0+内容+最后一个包（空包）
					seq++
					for {
						read, err := buffReader.Read(raw)
						if err == io.EOF || read < 0 {
							break
						}
						if tmp := snappy.Encode(nil, raw); len(tmp) > 0 {
							pageSend(tmp, seq, true, true)
						} else {
							pageSend(raw[:read], seq, true, false)
						}
						seq++
					}
					pageSend([]byte{}, seq, false, false)
				}
			})
		}
	}
}
