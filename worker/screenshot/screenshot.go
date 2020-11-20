package screenshot

import (
	"fmt"
	"github.com/shirou/gopsutil/host"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
			os := info.OS
			cfg := viper.New()
			cfg.SetConfigType("yaml")
			cfg.ReadConfig(f)
			tmp.cmd = cfg.Get(fmt.Sprintf("%s.cmd", os)).(string)
			tmp.file = cfg.Get(fmt.Sprintf("%s.file", os)).(string)
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

func (this *Screenshot) GetFile() ([]byte,error){
	return ioutil.ReadFile(this.file)
}

func (this *Screenshot) GetPath() string{
	return this.file
}

func (this *Screenshot) GetInfoStat() (ret *host.InfoStat, err error) {
	return host.Info()
}
