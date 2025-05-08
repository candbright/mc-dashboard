import request from '@/pkg/request/request.js'

export function requestServerInfos(params = {}) {
  return request.post('/server/info/list', params)
}

export function requestServerInfo(id) {
  return request.post(`/server/${id}/info/get`)
}

export function downloadLatestServer(id) {
  return request.post(`/server/${id}/download_start`)
}

export function startServer(id) {
  return request.post(`/server/${id}/start`)
}

export function stopServer(id) {
  return request.post(`/server/${id}/stop`)
}

export function uploadServerFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/server/saves/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 300000, // 5分钟超时
    onUploadProgress: (progressEvent) => {
      const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      console.log(`上传进度: ${percentCompleted}%`)
    }
  })
}

/**
 * 获取存档文件列表
 * @param {Object} params - 查询参数
 * @param {number} params.page - 页码，从1开始
 * @param {number} params.size - 每页数量
 * @returns {Promise} 返回存档列表数据
 */
export function requestSaveList(params = {}) {
  return request.post('/server/saves/list', params)
}

/**
 * 删除存档文件
 * @param {string} filename - 存档文件名
 * @returns {Promise} 返回删除结果
 */
export function deleteSaveFile(filename) {
  return request.post('/server/saves/delete', { filename })
}

/**
 * 应用存档到服务器
 * @param {string} serverId - 服务器ID
 * @param {string} filename - 存档文件名
 * @returns {Promise} 返回应用结果
 */
export function applySave(serverId, filename) {
  return request.post('/server/saves/apply', {
    server_id: serverId,
    filename
  })
}
