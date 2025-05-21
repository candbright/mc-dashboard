import request from '@/pkg/request/request.js'

/**
 * 获取服务器列表
 * @param {Object} params - 查询参数
 * @param {number} [params.page=0] - 页码，从0开始
 * @param {number} [params.size=0] - 每页数量
 * @param {string} [params.order='desc'] - 排序方向，可选值：'asc'（升序）或 'desc'（降序）
 * @param {string} [params.order_by='created_at'] - 排序字段，可选值：
 *   - 'id': 按服务器ID排序
 *   - 'name': 按服务器名称排序
 *   - 'version': 按服务器版本排序
 *   - 'created_at': 按创建时间排序
 * @returns {Promise} 返回服务器列表数据，包含以下字段：
 *   - total: 总记录数
 *   - items: 服务器列表，每个服务器对象包含：
 *     - id: 服务器ID
 *     - name: 服务器名称
 *     - description: 服务器描述
 *     - world_name: 世界名称
 *     - version: 服务器版本
 *     - active: 是否运行中
 *     - created_at: 创建时间
 */
export function requestServerInfos(params = {}) {
  return request.get('/servers', { params })
}

/**
 * 获取服务器信息
 * @param {string} id - 服务器ID
 * @returns {Promise} 返回服务器详细信息
 */
export function requestServerInfo(id) {
  return request.get(`/servers/${id}`)
}

/**
 * 创建服务器
 * @param {Object} params - 创建参数
 * @param {string} params.name - 服务器名称
 * @param {string} [params.description] - 服务器描述
 * @param {string} [params.world_name] - 世界名称
 * @param {string} [params.version] - 服务器版本
 * @returns {Promise} 返回创建结果
 */
export function createServer(params) {
  return request.post('/servers', params)
}

/**
 * 更新服务器信息
 * @param {string} id - 服务器ID
 * @param {Object} params - 更新参数
 * @param {string} params.name - 服务器名称
 * @param {string} [params.description] - 服务器描述
 * @returns {Promise} 返回更新结果
 */
export function updateServer(id, params) {
  return request.put(`/servers/${id}`, params)
}

/**
 * 启动服务器
 * @param {string} id - 服务器ID
 * @returns {Promise} 返回启动结果
 */
export function startServer(id) {
  return request.post(`/servers/${id}/start`)
}

/**
 * 停止服务器
 * @param {string} id - 服务器ID
 * @returns {Promise} 返回停止结果
 */
export function stopServer(id) {
  return request.post(`/servers/${id}/stop`)
}

/**
 * 删除服务器
 * @param {string} id - 服务器ID
 * @returns {Promise} 返回删除结果
 */
export function deleteServer(id) {
  return request.delete(`/servers/${id}`)
}

/**
 * 获取服务器属性配置
 * @param {string} id - 服务器ID
 * @returns {Promise} 返回服务器属性配置
 */
export function getServerProperties(id) {
  return request.get(`/servers/${id}/server_properties`)
}

/**
 * 更新服务器属性配置
 * @param {string} id - 服务器ID
 * @param {Object} properties - 服务器属性配置
 * @returns {Promise} 返回更新结果
 */
export function updateServerProperties(id, properties) {
  return request.put(`/servers/${id}/server_properties`, properties)
}

/**
 * 获取服务器白名单列表
 * @param {string} id - 服务器ID
 * @returns {Promise} 返回白名单列表
 */
export function getAllowList(id) {
  return request.get(`/servers/${id}/allowlist`)
}

/**
 * 添加白名单用户
 * @param {string} id - 服务器ID
 * @param {string} username - 用户名
 * @returns {Promise} 返回添加结果
 */
export function addAllowListUser(id, username) {
  return request.post(`/servers/${id}/allowlist/${username}`)
}

/**
 * 删除白名单用户
 * @param {string} id - 服务器ID
 * @param {string} username - 用户名
 * @returns {Promise} 返回删除结果
 */
export function deleteAllowListUser(id, username) {
  return request.delete(`/servers/${id}/allowlist/${username}`)
}

export function downloadLatestServer(id) {
  return request.post(`/server/${id}/download_start`)
}

/**
 * 获取存档文件列表
 * @param {Object} params - 查询参数
 * @param {number} [params.page=0] - 页码，从0开始
 * @param {number} [params.size=0] - 每页数量
 * @param {string} [params.order='desc'] - 排序方向，可选值：'asc'（升序）或 'desc'（降序）
 * @param {string} [params.order_by='created_at'] - 排序字段
 * @returns {Promise} 返回存档列表数据，包含以下字段：
 *   - total: 总记录数
 *   - items: 存档列表，每个存档对象包含：
 *     - id: 存档ID
 *     - name: 存档名称
 *     - size: 文件大小
 *     - last_modified: 最后修改时间
 *     - created_at: 创建时间
 */
export function requestSaveList(params = {}) {
  return request.get('/saves', { params })
}

/**
 * 删除存档文件
 * @param {number} id - 存档ID
 * @returns {Promise} 返回删除结果
 */
export function deleteSaveFile(id) {
  return request.delete(`/saves/${id}`)
}

/**
 * 应用存档到服务器
 * @param {number} saveId - 存档ID
 * @param {number} serverId - 服务器ID
 * @returns {Promise} 返回应用结果
 */
export function applySave(saveId, serverId) {
  return request.post('/saves/apply', {
    save_id: saveId,
    server_id: serverId
  })
}

/**
 * 上传存档文件
 * @param {File} file - 存档文件
 * @returns {Promise} 返回上传结果
 */
export function uploadServerFile(file) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/saves', formData, {
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
 * 获取服务器控制台日志
 * @param {string} id - 服务器ID
 * @param {number} [line=100] - 获取最后几行日志
 * @returns {Promise} 返回日志内容
 */
export function getConsoleLog(id, line = 100) {
  return request.get(`/servers/${id}/console_log`, {
    params: { line }
  })
}
