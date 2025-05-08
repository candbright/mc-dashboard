package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// DownloadStatus 下载状态
type DownloadStatus struct {
	TotalBytes    int64   // 文件总大小
	Downloaded    int64   // 已下载大小
	Percentage    float64 // 下载百分比
	Speed         float64 // 下载速度(KB/s)
	IsDownloading bool    // 是否正在下载
	Err           error   // 错误信息
}

// Downloader 下载器
type Downloader struct {
	client       *http.Client
	status       DownloadStatus
	statusMutex  sync.Mutex
	stopChan     chan struct{}
	statusChan   chan DownloadStatus
	lastUpdate   time.Time
	lastDownload int64
}

// NewDownloader 创建新的下载器
func NewDownloader() *Downloader {
	return &Downloader{
		client:     &http.Client{},
		stopChan:   make(chan struct{}),
		statusChan: make(chan DownloadStatus, 10),
	}
}

// Download 异步下载文件
func (d *Downloader) Download(url, filepath string) {
	go func() {
		defer close(d.statusChan)

		// 获取文件信息
		resp, err := d.client.Head(url)
		if err != nil {
			d.updateStatus(DownloadStatus{Err: err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			d.updateStatus(DownloadStatus{Err: fmt.Errorf("HTTP error: %s", resp.Status)})
			return
		}

		totalBytes := resp.ContentLength
		d.updateStatus(DownloadStatus{
			TotalBytes:    totalBytes,
			IsDownloading: true,
		})

		// 创建目标文件
		file, err := os.Create(filepath)
		if err != nil {
			d.updateStatus(DownloadStatus{Err: err})
			return
		}
		defer file.Close()

		// 开始下载
		resp, err = d.client.Get(url)
		if err != nil {
			d.updateStatus(DownloadStatus{Err: err})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			d.updateStatus(DownloadStatus{Err: fmt.Errorf("HTTP error: %s", resp.Status)})
			return
		}

		// 创建带缓冲的reader
		reader := io.TeeReader(resp.Body, &progressWriter{
			total:      totalBytes,
			downloader: d,
		})

		// 复制数据到文件
		_, err = io.Copy(file, reader)
		if err != nil {
			d.updateStatus(DownloadStatus{Err: err})
			return
		}

		// 下载完成
		d.updateStatus(DownloadStatus{
			TotalBytes:    totalBytes,
			Downloaded:    totalBytes,
			Percentage:    100,
			IsDownloading: false,
		})
	}()
}

// Stop 停止下载
func (d *Downloader) Stop() {
	close(d.stopChan)
}

// Status 获取状态通道
func (d *Downloader) Status() <-chan DownloadStatus {
	return d.statusChan
}

func (d *Downloader) GetCurrentStatus() DownloadStatus {
	d.statusMutex.Lock()
	defer d.statusMutex.Unlock()
	return d.status
}

func (d *Downloader) updateStatus(status DownloadStatus) {
	d.statusMutex.Lock()
	defer d.statusMutex.Unlock()

	// 计算下载速度
	now := time.Now()
	if !d.lastUpdate.IsZero() {
		elapsed := now.Sub(d.lastUpdate).Seconds()
		if elapsed > 0 {
			downloaded := status.Downloaded - d.lastDownload
			status.Speed = float64(downloaded) / 1024 / elapsed
		}
	}
	d.lastUpdate = now
	d.lastDownload = status.Downloaded

	// 计算百分比
	if status.TotalBytes > 0 {
		status.Percentage = float64(status.Downloaded) / float64(status.TotalBytes) * 100
	}

	d.status = status
	select {
	case d.statusChan <- status:
	default: // 避免阻塞
	}
}

// progressWriter 用于跟踪下载进度
type progressWriter struct {
	total      int64
	downloaded int64
	downloader *Downloader
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.downloaded += int64(n)

	pw.downloader.updateStatus(DownloadStatus{
		TotalBytes:    pw.total,
		Downloaded:    pw.downloaded,
		IsDownloading: true,
	})

	return n, nil
}

// 使用示例
func main() {
	downloader := NewDownloader()

	// 开始下载
	url := "https://example.com/largefile.zip"
	filepath := "largefile.zip"
	downloader.Download(url, filepath)

	// 实时获取下载状态
	for status := range downloader.Status() {
		if status.Err != nil {
			fmt.Printf("下载错误: %v\n", status.Err)
			break
		}

		fmt.Printf("进度: %.2f%%, 速度: %.2f KB/s, 已下载: %d/%d bytes\n",
			status.Percentage,
			status.Speed,
			status.Downloaded,
			status.TotalBytes)

		if !status.IsDownloading && status.Percentage >= 100 {
			fmt.Println("下载完成!")
			break
		}

		// 模拟中途取消
		// if status.Percentage > 30 {
		// 	fmt.Println("取消下载...")
		// 	downloader.Stop()
		// }
	}
}
