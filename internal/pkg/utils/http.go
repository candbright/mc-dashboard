package utils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	logFunc = func(format string, args ...interface{}) {
		fmt.Printf(format, args...)
	}
	client = defaultHTTPClient()
)

func SetLogFunc(f func(format string, args ...interface{})) {
	if f != nil {
		logFunc = f
	}
}

func SetHTTPClient(c *http.Client) {
	if c != nil {
		client = c
	}
}

func GetHTTPClient() *http.Client {
	return client
}

// CreateHTTPClient 创建一个配置好的 HTTP 客户端
func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second, // 1分钟超时
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConns:          10,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 10 * time.Second,
			DisableCompression:    true,
			ForceAttemptHTTP2:     false,
			MaxConnsPerHost:       1,
			MaxIdleConnsPerHost:   1,
			DisableKeepAlives:     true,
			ResponseHeaderTimeout: 30 * time.Second,
			Proxy:                 http.ProxyFromEnvironment,
		},
	}
}

// CreateWindowsHTTPRequest 创建一个配置好的 HTTP 请求
func CreateWindowsHTTPRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "close")

	return req, nil
}

func Download(url string, dstPath string) error {
	// 创建请求
	req, err := CreateWindowsHTTPRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file, status code: %d", resp.StatusCode)
	}

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 创建目标文件
	file, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer file.Close()

	// 获取文件大小
	contentLength := resp.ContentLength
	var downloaded int64

	// 创建进度写入器
	progressWriter := &progressWriter{
		writer: file,
		onProgress: func(n int64) {
			downloaded += n
			if contentLength > 0 {
				percentage := float64(downloaded) / float64(contentLength) * 100
				if logFunc != nil {
					logFunc("Download progress: %.2f%% (%d/%d bytes)", percentage, downloaded, contentLength)
				}
			}
		},
	}

	// 复制文件内容
	_, err = io.Copy(progressWriter, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	if logFunc != nil {
		logFunc("Download completed: %s", dstPath)
	}

	return nil
}

// progressWriter 用于跟踪下载进度
type progressWriter struct {
	writer     io.Writer
	onProgress func(n int64)
}

func (pw *progressWriter) Write(p []byte) (n int, err error) {
	n, err = pw.writer.Write(p)
	if n > 0 && pw.onProgress != nil {
		pw.onProgress(int64(n))
	}
	return
}
