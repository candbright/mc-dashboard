package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemHandler 系统处理器
type SystemHandler struct {
}

// NewSystemHandler 创建系统处理器实例
func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// RegisterRoutes 注册路由
func (h *SystemHandler) RegisterRoutes(r *gin.RouterGroup) {
	system := r.Group("/system")
	{
		system.GET("/status", h.GetSystemStatus)
	}
}

// GetSystemStatus godoc
// @Summary 获取系统状态
// @Description 获取当前系统的运行状态，包括CPU、内存、磁盘使用情况等
// @Tags 系统
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} SystemStatus "系统状态信息"
// @Failure 500 {object} errors.AppError "服务器内部错误"
// @Router /system/status [get]
func (h *SystemHandler) GetSystemStatus(ctx *gin.Context) {
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		ctx.Error(err)
		return
	}

	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		ctx.Error(err)
		return
	}

	// 获取磁盘信息
	diskInfo, err := disk.Usage("/")
	if err != nil {
		ctx.Error(err)
		return
	}

	// 获取系统运行时间
	hostInfo, err := host.Info()
	if err != nil {
		ctx.Error(err)
		return
	}

	status := SystemStatus{
		CPUUsage:    cpuPercent[0],
		MemoryTotal: memInfo.Total,
		MemoryUsed:  memInfo.Used,
		DiskTotal:   diskInfo.Total,
		DiskUsed:    diskInfo.Used,
		Uptime:      hostInfo.Uptime,
		OS:          runtime.GOOS,
		Arch:        runtime.GOARCH,
	}

	ctx.JSON(http.StatusOK, status)
}

// SystemStatus 系统状态信息
type SystemStatus struct {
	CPUUsage    float64 `json:"cpu_usage"`    // CPU使用率
	MemoryTotal uint64  `json:"memory_total"` // 总内存（字节）
	MemoryUsed  uint64  `json:"memory_used"`  // 已用内存（字节）
	DiskTotal   uint64  `json:"disk_total"`   // 总磁盘空间（字节）
	DiskUsed    uint64  `json:"disk_used"`    // 已用磁盘空间（字节）
	Uptime      uint64  `json:"uptime"`       // 系统运行时间（秒）
	OS          string  `json:"os"`           // 操作系统
	Arch        string  `json:"arch"`         // 系统架构
}
