import request from '@/pkg/request/request.js'
import { Utils } from '@/pkg/util/util.js'

export function getServerProperties(serverId) {
  return request.get(`/servers/${serverId}/server_properties`)
}

export function setServerProperties(serverId, serverProperties) {
  return request.put(`/servers/${serverId}/server_properties`, Utils.convertToString(serverProperties))
}
