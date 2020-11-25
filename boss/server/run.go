package server

import (
	"bufio"
	"codec"
	"context"
	"encoding/json"
	"ext"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/alecthomas/log4go"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gpool"
	"github.com/gogf/gf/container/gtype"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/nfnt/resize"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"proc"
	"strings"
	"time"
	"tool"
	www "web"
	"web/model"
)

const (
	TCP = "tcp"
)

type None = struct{}

type capture struct {
	recv   chan []byte
	file   *os.File
	status *gtype.Bool
}

type connectItem = struct {
	//连接
	conn gnet.Conn
	//客户端机器名称
	name string
	//捕捉数据接受channel
	cap capture
	//客户端断开channel
	close chan None
	//上一次ping的事假
	ping int64
}

type ConnectContainer = gmap.StrAnyMap

type App struct {
	*gnet.EventServer
	port model.Port
	//协程池
	workerPool *goroutine.Pool
	//协议处理
	cmdPool *gpool.Pool
	//chan切片池
	chanBytePool *gpool.Pool
	logger       *log4go.Logger
	//在线连接
	connects *ConnectContainer
	//web服务器
	web *www.Web
	//消息处理
	mill *ext.ExtGoChanel
}

func NewApp(logger *log4go.Logger, p model.Port) *App {
	return &App{logger: logger, connects: gmap.NewStrAnyMap(true), port: p}
}

/** 使用自定义解码器解码 */
func (this *App) Start(ctx context.Context) error {
	return gnet.Serve(this, fmt.Sprintf("%s://0.0.0.0:%d", TCP, this.port.Socket), gnet.WithLogger(Logger{this.logger}), gnet.WithTicker(true), gnet.WithReusePort(true), gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*1), gnet.WithCodec(codec.NewCodec()))
}

func (this *App) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	//协程池
	this.workerPool = goroutine.Default()
	this.logger.Info("启动成功 %s (multi-cores: %t, loops: %d)", srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	this.mill = ext.NewExtGoChannel(gochannel.Config{Persistent: true, BlockPublishUntilSubscriberAck: false}, watermill.NewStdLogger(false, false))
	this.web = www.NewWeb(this.logger, this.mill, this.workerPool)
	//协议池
	this.cmdPool = gpool.New(0, func() (interface{}, error) {
		ret := new(proc.Cmd)
		ret.Head = new(proc.Head)
		ret.Content = new(proc.Content)
		return ret, nil
	}, nil)
	//切片池
	this.chanBytePool = gpool.New(0, func() (interface{}, error) {
		return make(chan []byte), nil
	}, func(i interface{}) {
		if c, ok := i.(chan []byte); ok {
			close(c)
		}
	})
	//启动WEB服务器
	this.workerPool.Submit(func() {
		if err := this.web.Run(fmt.Sprintf(":%d", this.port.Web)); nil != err {
			this.logger.Critical(err)
		}
	})
	//订阅查询所有客户端消息，给其他查询
	this.workerPool.Submit(func() {
		if messages, err := this.mill.Subscribe(context.Background(), tool.QUERY_WORKERS); nil == err {
			for {
				select {
				case <-messages:
					//查询所在的客户信息(传输到消息队列即可,切片拷贝的是地址，不用考虑)
					workers := []model.Worker{}
					this.connects.Iterator(func(k string, v interface{}) bool {
						tar := v.(*connectItem)
						workers = append(workers, model.Worker{Addr: k, Name: tar.name, Ping: tar.ping})
						return true
					})
					if res, err := json.Marshal(workers); nil == err {
						this.mill.Publish(tool.QUERY_WORKERS_RESULT, message.NewMessage(watermill.NewUUID(), res))
					}
				}
			}
		} else {
			this.logger.Error(err)
		}
	})
	//处理接收到的截图并生成缩略图(后台任务)
	this.workerPool.Submit(func() {
		if tasks, err := this.mill.Subscribe(context.Background(), tool.PROCESS_FOR_SLOW); nil == err {
			for {
				select {
				case task := <-tasks:
					if "build_capture" == task.Metadata.Get("name") {
						//生成缩略图(方便给其他接口处理)
						spath := string(task.Payload)
						if sf, errs := os.Open(spath); nil == errs {
							if img, _, erri := image.Decode(sf); nil == erri {
								if nimg := resize.Resize(400, 0, img, resize.InterpolationFunction(resize.Lanczos3)); nil != nimg {
									if tf, errt := os.OpenFile(fmt.Sprintf("%s/capture_thumb%s", filepath.Dir(spath), filepath.Ext(spath)), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm); nil == errt {
										if erre := jpeg.Encode(tf, nimg, &jpeg.Options{Quality: 90}); nil != erre {
											this.logger.Warn(erre)
										}
										tf.Close()
									} else {
										this.logger.Warn(errt)
									}
								}
							} else {
								this.logger.Warn(erri)
							}
							sf.Close()
						} else {
							this.logger.Warn(errs)
						}
					}
					task.Ack()
				}
			}
		} else {
			this.logger.Error(err)
		}
	})
	return
}

func (this *App) Tick() (delay time.Duration, action gnet.Action) {
	now := time.Now().Unix()
	active := 0
	this.connects.Iterator(func(k string, v interface{}) bool {
		tar := v.(*connectItem)
		if (now - tar.ping) > 10 {
			//触发下线通知
			if !this.mill.IsPause(tool.WORKER_OFFLINE) {
				if d, e := json.Marshal(model.Worker{Addr: strings.Split(tar.conn.RemoteAddr().String(), ":")[0], Name: tar.name, Ping: tar.ping}); nil == e {
					this.mill.Publish(tool.WORKER_OFFLINE_RESULT, message.NewMessage(watermill.NewUUID(), d))
				}
			}
			if err := tar.conn.Close(); nil != err {
				this.logger.Warn(err)
			}
			tar.close <- None{}
			this.chanBytePool.Put(tar.cap.recv)
			this.connects.Remove(k)
		} else {
			active++
		}
		return true
	})
	if active == 0 {
		//适当延迟数据
		return time.Second * 3, gnet.None
	}
	return time.Second * 1, gnet.None
}

func (this *App) OnShutdown(svr gnet.Server) {
	this.workerPool.Release()
	this.chanBytePool.Close()
	this.cmdPool.Close()
}

func (this *App) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if action != gnet.Close && action != gnet.Shutdown {
		if len(frame) > 0 {
			remoteAddr := strings.Split(c.RemoteAddr().String(), ":")[0]
			//read
			if cmd, err := this.cmdPool.Get(); nil == err {
				if nil != cmd {
					msg := cmd.(*proc.Cmd)
					if err := proto.Unmarshal(frame, msg); nil == err {
						switch *msg.Content.Type {
						case proc.ContentType_PING:
							//处理PING请求(用ping请求来辅助检测客户端是否还存活)
							if this.connects.Contains(remoteAddr) {
								tar := this.connects.Get(remoteAddr).(*connectItem)
								tar.ping = time.Now().Unix()
								tar.name = *msg.Content.GetPing().Name
							}
							break
						case proc.ContentType_CAPTURE:
							//发送数据给接收者接收文件数据(状态异常就丢弃数据了，防止都塞)
							if this.connects.Contains(remoteAddr) {
								tar := this.connects.Get(remoteAddr).(*connectItem)
								if *(msg.Content.GetCapture().Compress) {
									if tmp, derr := snappy.Decode(nil, msg.Content.GetCapture().GetData()); nil == derr {
										tar.cap.recv <- tmp
									} else {
										this.logger.Error(derr)
										//恢复一个异常解码数据
									}
								} else {
									tar.cap.recv <- msg.Content.GetCapture().GetData()
								}
							}
							break
						}
					} else {
						this.logger.Info("数据包异常:%s", err.Error())
					}
					this.cmdPool.Put(cmd)
				}
			}
		}
	}
	return
}

func (this *App) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	remoteAddr := strings.Split(c.RemoteAddr().String(), ":")[0]
	nowUnix := time.Now().Unix()
	conitem := &connectItem{conn: c, close: make(chan None), ping: nowUnix}
	conitem.cap = capture{nil, nil, gtype.NewBool(true)}
	if i, e := this.chanBytePool.Get(); nil == e {
		if c, ok := i.(chan []byte); ok {
			conitem.cap.recv = c
		}
	}
	//发布客户端上线消息(保证至少一个管理员上线的情况下才发布消息)
	if !this.mill.IsPause(tool.WORKER_ONLINE) {
		if d, e := json.Marshal(model.Worker{Addr: remoteAddr, Name: "", Ping: nowUnix}); nil == e {
			this.mill.Publish(tool.WORKER_ONLINE_RESULT, message.NewMessage(watermill.NewUUID(), d))
		}
	}
	this.connects.SetIfNotExist(remoteAddr, conitem)

	//接收截图
	this.workerPool.Submit(func() {
		//匿名函数内部使用的只是指针的拷贝，可能对象已被释放(因此需要拷贝)
		item := conitem
		if nil != item {
			if dir, e := os.Getwd(); nil == e {
				dir = fmt.Sprintf("%s/data/%s", dir, strings.Split(remoteAddr, ":")[0])
				if e = os.MkdirAll(dir, os.ModePerm); e != nil {
					this.logger.Error(e)
					conitem.cap.status.Set(false)
				} else {
					//目标文件路径
					save := fmt.Sprintf("%s/capture.jpg", dir)
					//缓冲区写入文件
					var writer *bufio.Writer = nil
					for {
						select {
						case <-item.close:
							goto CLOSE_CON
						case data := <-item.cap.recv:
							//code here
							if item.cap.status.Val() {
								if 0 == len(data) {
									if nil == item.cap.file {
										//start
										if item.cap.file, e = os.OpenFile(save, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, os.ModePerm); nil == e {
											if nil == writer {
												//文件缓冲设置为10KB
												writer = bufio.NewWriterSize(item.cap.file, 1024*10)
											}
											writer.Reset(item.cap.file)
										} else {
											this.logger.Error(e)
											conitem.cap.status.Set(false)
										}
									} else {
										//停止(保存目标文件)
										writer.Flush()
										item.cap.file.Close()
										item.cap.file = nil
										if nm := message.NewMessage(watermill.NewUUID(), []byte(save)); nil != nm {
											nm.Metadata.Set("name", "build_capture")
											this.mill.Publish(tool.PROCESS_FOR_SLOW, nm)
										}
									}
								} else {
									if nil != item.cap.file && nil != writer {
										writer.Write(data)
									}
								}
							}
						}
					}
				CLOSE_CON:
					//
				}
			}
		}
	})
	return
}

func (this *App) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	removeAddr := strings.Split(c.RemoteAddr().String(), ":")[0]
	if this.connects.Contains(removeAddr) {
		tar := this.connects.Get(removeAddr).(*connectItem)
		if nil != tar {
			//发布客户端离线消息
			if !this.mill.IsPause(tool.WORKER_OFFLINE) {
				if d, e := json.Marshal(model.Worker{Addr: removeAddr, Name: tar.name, Ping: tar.ping}); nil == e {
					this.mill.Publish(tool.WORKER_OFFLINE_RESULT, message.NewMessage(watermill.NewUUID(), d))
				}
			}
			if nil != tar.conn {
				tar.close <- None{}
				tool.Try(func() {
					close(tar.close)
				}, func(i interface{}) {
					this.logger.Info(i)
				})
				this.chanBytePool.Put(tar.cap.recv)
				this.connects.Remove(removeAddr)
			}
		}
	}
	return
}