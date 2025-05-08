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
