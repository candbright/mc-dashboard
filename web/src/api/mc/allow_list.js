import request from '@/pkg/request/request.js'

export function addAllowList(id, username) {
  return request.post(`/server/${id}/allowlist/add`, { username: username })
}

export function deleteAllowList(id, username) {
  return request.post(`/server/${id}/allowlist/delete`, { username: username })
}
