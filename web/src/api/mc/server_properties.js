import request from '@/pkg/request/request.js'
import { Utils } from '@/api/mc/util.js'

export function getServerProperties(id) {
  return request.post(`/server/${id}/server_properties/get`)
}

export function setServerProperties(id, serverProperties) {
  return request.post(`/server/${id}/server_properties/set`, Utils.convertNumbersToString(serverProperties))
}
