import request from '@/pkg/request/request.js'

export function getAllowList(serverId) {
  return request.get(`/servers/${serverId}/allowlist`)
}

export function addAllowList(serverId, username) {
  return request.post(`/servers/${serverId}/allowlist/${username}`)
}

export function deleteAllowList(serverId, username) {
  return request.delete(`/servers/${serverId}/allowlist/${username}`)
}
