package core

import (
	"fmt"
	"github.com/candbright/go-log/log"
	"os"
	"os/exec"
	"runtime"
)

type Cmd struct {
	*exec.Cmd
}

func (c *Cmd) Run() error {
	log.Infof("Running cmd %v", c.Args)
	return c.Cmd.Run()
}

func Command(name string, arg ...string) *Cmd {
	command := exec.Command(name, arg...)
	return &Cmd{command}

}

func Exists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true // 文件存在
	}
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return false
}

func Download(url string, dstPath string) error {
	switch runtime.GOOS {
	case "linux":
		err := Command("wget", "--no-check-certificate",
			fmt.Sprintf("--user-agent=%s", "\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36\""),
			fmt.Sprintf("--timeout=%d", 600),
			"-P", dstPath, url).Run()
		if err != nil {
			return err
		}
		return nil
	case "windows":
		err := Command("powershell", "Invoke-WebRequest", "-Uri", url, "-OutFile", dstPath).Run()
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func Zip(src, dst string) error {
	switch runtime.GOOS {
	case "linux":
		err := Command("zip", "-r", src, dst).Run()
		if err != nil {
			return err
		}
		return nil
	case "windows":
		err := Command("powershell", "-Command", "Compress-Archive",
			"-Path", src, "-DestinationPath", dst).Run()
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported platform")
	}
}

func Unzip(src, dst string) error {
	switch runtime.GOOS {
	case "linux":
		err := Command("unzip", "-q", src, "-d", dst).Run()
		if err != nil {
			return err
		}
		return nil
	case "windows":
		err := Command("powershell", "-Command", "Expand-Archive",
			"-Path", src, "-DestinationPath", dst,
			"-Force").Run()
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported platform")
	}
}
