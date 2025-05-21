import request from '@/pkg/request/request.js'

/**
 * 获取系统状态信息
 * @returns {Promise} 返回系统状态数据，包含以下字段：
 *   - cpu_usage: CPU使用率（百分比）
 *   - memory_total: 总内存（字节）
 *   - memory_used: 已用内存（字节）
 *   - disk_total: 总磁盘空间（字节）
 *   - disk_used: 已用磁盘空间（字节）
 *   - uptime: 系统运行时间（秒）
 *   - go_version: Go版本
 *   - os: 操作系统
 *   - arch: 系统架构
 */
export function requestSystemStatus() {
  return request.get('/system/status')
}
