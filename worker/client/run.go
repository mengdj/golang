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
	"github.com/robfig/cron"
	"io"
	"net"
	"os"
	"proc"
	"screenshot"
	"syscall"
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
	logger    *log4go.Logger
	cmdPool   *gpool.Pool
	packCodec *codec.Codec
	capt      *screenshot.Screenshot
	captId    uint32
	pubSub    *ext.ExtGoChanel
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
		remotAddr      *net.TCPAddr
		tcpCon         *net.TCPConn
		tcpErr         error
		connectFailure chan None
	)
	cron_task := cron.New()
	//处理连接失败
	connectFailure = make(chan None)
	connectFailureStatus := gtype.NewBool(false)
	defer func() {
		close(connectFailure)
	}()
	remotAddr, tcpErr = net.ResolveTCPAddr(TCP, fmt.Sprintf("%s:%d", *ret.Body.Server, *ret.Body.Port))
	if nil == tcpErr {
		tcpCon, tcpErr = net.DialTCP(TCP, nil, remotAddr)
		if nil == tcpErr {
			this.logger.Info("服务器连接成功:%s", tcpCon.RemoteAddr().String())
			this.packCodec = codec.NewCodec(tcpCon)
			writeMessage, writeStatus := this.pubSub.Subscribe(context.Background(), SEND_PACKET)
			readMessage, readStatus := this.pubSub.Subscribe(context.Background(), RECIVE_PACKET)
			if nil == writeStatus && nil == readStatus {
				//创建锁保证TCP流顺序发送
				syncLock := gmutex.New()
				//每5秒发送心跳协议
				cron_task.AddFunc("*/5 * * * *", func() {
					this.ping(syncLock)
				})
				//每1分钟发送截图数据
				cron_task.AddFunc("0 */1 * * *", func() {
					this.capture(syncLock)
				})
				//读取接收到的包并发布到消息队列（读取）
				go func() {
					for {
						if rd, e := this.packCodec.Read(); nil == e {
							if len(rd) > 0 {
								if e = this.pubSub.Publish(RECIVE_PACKET, message.NewMessage(watermill.NewUUID(), rd)); nil != e {
									this.logger.Warn(e)
								}
							}
						} else if e == io.EOF || syscall.EINVAL == e {
							this.logger.Warn(e)
							break
						} else {
							this.logger.Warn(e)
							//process
							if _, ok := e.(*net.OpError); ok {
								break
							}
						}
					}
					//失败信息
					connectFailure <- None{}
				}()
				//处理消息队列（写入）并发送至服务器
				go func() {
					for msg := range writeMessage {
						if e := this.packCodec.Write(msg.Payload); nil == e {
							msg.Ack()
						} else if syscall.EINVAL == e {
							this.logger.Warn(e)
							break
						} else {
							this.logger.Warn(e)
							//process
							if _, ok := e.(*net.OpError); ok {
								break
							}
						}
					}
					//失败信息
					connectFailure <- None{}
				}()
				//处理接收到的队列消息
				go func() {
					for msg:=range readMessage{
						if connectFailureStatus.Val() {
							break
						}
						//code here
						msg.Ack()
					}
				}()
				cron_task.Start()
				select {
				case <-ctx.Done():
					this.logger.Info("客户端终止")
				case <-connectFailure:
					connectFailureStatus.Set(true)
				}
				//取消队列数据(PING)
				cron_task.Stop()
				return errors.New("客户端关闭")
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
	if lock.TryLock() {
		var (
			hostName    string           = ""
			contentType proc.ContentType = proc.ContentType_PING
		)
		contentPing := &proc.Content_Ping{}
		contentPing.Ping = &proc.Ping{}
		//获取主机名
		if h, e := os.Hostname(); nil == e {
			hostName = h
		}
		contentPing.Ping.Name = &hostName
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
		lock.Unlock()
	}
}

/** 捕捉屏幕截图（可能比较耗时因此加入锁） */
func (this *Client) capture(lock *gmutex.Mutex) {
	if lock.TryLock() {
		var (
			contentType proc.ContentType = proc.ContentType_CAPTURE
			compress    bool             = false
		)
		this.captId += 1
		this.captId %= 65536
		//屏幕截图数据太大，需要多次发送 CAPTURE_PAGE_SIZE
		if nil == this.capt.Capture() {
			if file, err := os.Open(this.capt.GetPath()); nil == err {
				defer func() {
					file.Close()
				}()
				//内部函数实现发送截图数据到队列
				pageSend := func(d []byte, seq uint32, more bool) {
					if ret, err := this.cmdPool.Get(); nil == err {
						if ins, ok := ret.(*proc.Cmd); ok {
							contentCapture := &proc.Content_Capture{}
							contentCapture.Capture = &proc.Capture{}
							contentCapture.Capture.Compress = &compress
							contentCapture.Capture.More = &more
							contentCapture.Capture.Id = &this.captId
							contentCapture.Capture.Data = d
							contentCapture.Capture.Seq = &seq
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
				//
				raw := make([]byte, CAPTURE_PAGE_SIZE)
				buffReader := bufio.NewReader(file)
				var seq uint32 = 1
				for {
					read, err := buffReader.Read(raw)
					if err == io.EOF || read < 0 {
						break
					}
					pageSend(raw[:read], seq, true)
					seq++
				}
				pageSend([]byte{}, seq, false)
			}
		}
		lock.Unlock()
	}
}
