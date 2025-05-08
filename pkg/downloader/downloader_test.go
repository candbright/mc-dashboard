package downloader

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

// 测试辅助函数 - 创建测试HTTP服务器
func createTestServer(content string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Set("Content-Length", string(len(content)))
			return
		}
		w.Write([]byte(content))
	}))
}

// 测试正常下载
func TestDownloader_NormalDownload(t *testing.T) {
	// 创建测试服务器
	testContent := strings.Repeat("This is test content. ", 1000)
	server := createTestServer(testContent)
	defer server.Close()

	// 创建下载器
	downloader := NewDownloader()
	defer downloader.Stop()

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "download_test_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// 开始下载
	downloader.Download(server.URL, tmpFile.Name())

	// 检查下载状态
	var finalStatus DownloadStatus
	for status := range downloader.Status() {
		if status.Err != nil {
			t.Fatalf("下载失败: %v", status.Err)
		}
		finalStatus = status

		if !status.IsDownloading && status.Percentage >= 100 {
			break
		}
	}

	// 验证下载完成状态
	if finalStatus.Percentage != 100 {
		t.Errorf("期望下载完成100%%, 实际得到 %.2f%%", finalStatus.Percentage)
	}

	// 验证文件内容
	fileContent, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("读取下载文件失败: %v", err)
	}

	if string(fileContent) != testContent {
		t.Error("下载文件内容与测试内容不匹配")
	}
}

// 测试下载错误(无效URL)
func TestDownloader_InvalidURL(t *testing.T) {
	downloader := NewDownloader()
	defer downloader.Stop()

	tmpFile, err := os.CreateTemp("", "download_test_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// 使用无效URL
	downloader.Download("http://invalid.url", tmpFile.Name())

	var gotError bool
	for status := range downloader.Status() {
		if status.Err != nil {
			gotError = true
			t.Logf("预期错误: %v", status.Err)
			break
		}
	}

	if !gotError {
		t.Error("期望得到错误，但下载成功完成")
	}
}

// 测试取消下载
func TestDownloader_CancelDownload(t *testing.T) {
	// 创建慢速测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Set("Content-Length", "1000000")
			return
		}
		// 模拟慢速下载
		for i := 0; i < 100; i++ {
			time.Sleep(100 * time.Millisecond)
			w.Write([]byte(strings.Repeat("x", 1000)))
		}
	}))
	defer server.Close()

	downloader := NewDownloader()

	tmpFile, err := os.CreateTemp("", "download_test_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// 开始下载
	downloader.Download(server.URL, tmpFile.Name())

	// 等待部分下载后取消
	time.Sleep(500 * time.Millisecond)
	downloader.Stop()

	// 检查是否收到取消信号
	var gotCancel bool
	for status := range downloader.Status() {
		if status.Err != nil {
			t.Logf("取消后的错误: %v", status.Err)
			gotCancel = true
			break
		}
		if status.Percentage > 0 && !status.IsDownloading {
			gotCancel = true
			break
		}
	}

	if !gotCancel {
		t.Error("未能成功取消下载")
	}
}

// 测试下载进度报告
func TestDownloader_ProgressReporting(t *testing.T) {
	testContent := strings.Repeat("x", 100000) // 100KB内容
	server := createTestServer(testContent)
	defer server.Close()

	downloader := NewDownloader()
	defer downloader.Stop()

	tmpFile, err := os.CreateTemp("", "download_test_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	downloader.Download(server.URL, tmpFile.Name())

	var progressReports []DownloadStatus
	for status := range downloader.Status() {
		progressReports = append(progressReports, status)
		if !status.IsDownloading && status.Percentage >= 100 {
			break
		}
	}

	// 检查是否收到多个进度更新
	if len(progressReports) < 2 {
		t.Errorf("期望至少2个进度报告, 实际得到 %d", len(progressReports))
	}

	// 检查进度是否递增
	var lastPercentage float64
	for i, report := range progressReports {
		if i > 0 && report.Percentage < lastPercentage {
			t.Errorf("进度百分比不应减少: 前一个 %.2f%%, 当前 %.2f%%", lastPercentage, report.Percentage)
		}
		lastPercentage = report.Percentage
	}

	// 检查最终完成状态
	lastReport := progressReports[len(progressReports)-1]
	if lastReport.Percentage != 100 {
		t.Errorf("期望最终进度100%%, 实际得到 %.2f%%", lastReport.Percentage)
	}
}

// 测试下载速度计算
func TestDownloader_DownloadSpeed(t *testing.T) {
	// 创建慢速服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "HEAD" {
			w.Header().Set("Content-Length", "10000")
			return
		}
		// 每100ms发送1KB数据
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			w.Write([]byte(strings.Repeat("x", 1000)))
		}
	}))
	defer server.Close()

	downloader := NewDownloader()
	defer downloader.Stop()

	tmpFile, err := os.CreateTemp("", "download_test_*.txt")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	downloader.Download(server.URL, tmpFile.Name())

	var speedSamples []float64
	for status := range downloader.Status() {
		if status.Speed > 0 {
			speedSamples = append(speedSamples, status.Speed)
		}
		if !status.IsDownloading {
			break
		}
	}

	// 检查速度样本
	if len(speedSamples) == 0 {
		t.Error("未收集到任何下载速度样本")
	}

	// 检查速度是否在合理范围内(预期约10KB/s)
	for _, speed := range speedSamples {
		if speed < 5 || speed > 15 { // 允许一定误差
			t.Errorf("下载速度 %.2f KB/s 不在预期范围内(5-15 KB/s)", speed)
		}
	}
}
