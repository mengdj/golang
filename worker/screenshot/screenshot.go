package screenshot

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/viper"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type InterpolationFunction int

// InterpolationFunction constants
const (
	// Nearest-neighbor interpolation
	NearestNeighbor InterpolationFunction = iota
	// Bilinear interpolation
	Bilinear
	// Bicubic interpolation (with cubic hermite spline)
	Bicubic
	// Mitchell-Netravali interpolation
	MitchellNetravali
	// Lanczos interpolation (a=2)
	Lanczos2
	// Lanczos interpolation (a=3)
	Lanczos3
)

/** 截图功能封装，为CAPTURE准备数据*/
type Screenshot struct {
	cmd  string
	file string
}

//默认读取当前目录下的config/screenshot.yaml文件
func NewScreenshotDefault() *Screenshot {
	if dir, err := os.Getwd(); nil == err {
		return NewScreenshot(fmt.Sprintf("%s/config/screenshot.yaml", dir))
	}
	return nil
}

func NewScreenshot(path string) *Screenshot {
	if f, err := os.Open(path); nil == err {
		tmp := &Screenshot{}
		defer func() {
			f.Close()
		}()
		if info, err := tmp.GetInfoStat(); nil == err {
			//读取配置文件获取对应的命令
			ost := info.OS
			cfg := viper.New()
			cfg.SetConfigType("yaml")
			cfg.ReadConfig(f)
			tmp.cmd = cfg.Get(fmt.Sprintf("%s.cmd", ost)).(string)
			tmp.file = fmt.Sprintf("%s/capture.jpg", os.TempDir())
			log.Println(tmp.file)
			if "" != tmp.cmd && "" != tmp.file {
				tmp.cmd = strings.ReplaceAll(tmp.cmd, "%FILE%", tmp.file)
				return tmp
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
	return nil
}

func (this *Screenshot) Capture() (err error) {
	params := strings.Split(this.cmd, " ")
	cmd := exec.Command(params[0], params[1:]...)
	err = cmd.Start()
	if nil == err {
		err = cmd.Wait()
	}
	return err
}

// Deprecated :效率太低了
func (this *Screenshot) GetFile() ([]byte, error) {
	return ioutil.ReadFile(this.file)
}

func (this *Screenshot) GetPath() string {
	return this.file
}

func (this *Screenshot) GetInfoStat() (ret *host.InfoStat, err error) {
	return host.Info()
}

//设置当前截图的大小，宽和高可独立设置为0，实现等比缩放
func (this *Screenshot) Resize(w, h uint, interp InterpolationFunction, quality int) error {
	return this.ResizePicture(this.file, w, h, interp, quality)
}

//设置制定截图的大小，宽和高可独立设置为0，实现等比缩放
func (this *Screenshot) ResizePicture(path string, w, h uint, interp InterpolationFunction, quality int) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if nil != err {
		panic(err)
	}
	defer file.Close()
	var (
		img, nimg image.Image
		offset    int64
	)
	img, _, err = image.Decode(file)
	if nil == err {
		nimg = resize.Resize(w, h, img, resize.InterpolationFunction(interp))
		if nil != nimg {
			//适配其他格式的代码暂时省略
			buff := bytes.NewBuffer([]byte{})
			//编码到缓冲区避免操作IO
			if err = jpeg.Encode(buff, nimg, &jpeg.Options{Quality: quality}); nil == err {
				//copy(需要清空文件并重置指针,此处代码未做备份处理失败的情况)
				if err = file.Truncate(0); nil == err {
					offset, err = file.Seek(0, io.SeekStart)
					if offset == 0 && nil == err {
						file.Write(buff.Bytes())
					} else {
						err = errors.New("resize store failture")
					}
				}
			}
		} else {
			err = errors.New("resize image encode failture")
		}
	}
	return err
}
