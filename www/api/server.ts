import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
  DeleteServerRequest,
  DeleteServerResponse,
  GetServerRequest,
  GetServerResponse,
  InitServerRequest,
  InitServerResponse,
  ListServersRequest,
  ListServersResponse,
} from '@/lib/pb/api_server'
import { BaseResponse } from '@/types/api'

export const getServer = async (req: GetServerRequest) => {
  const res = await http.post(API_PATH + '/server/get', GetServerRequest.toJson(req))
  return GetServerResponse.fromJson((res.data as BaseResponse).body)
}

export const listServer = async (req: ListServersRequest) => {
  const res = await http.post(API_PATH + '/server/list', ListServersRequest.toJson(req))
  return ListServersResponse.fromJson((res.data as BaseResponse).body)
}

export const deleteServer = async (req: DeleteServerRequest) => {
  const res = await http.post(API_PATH + '/server/delete', DeleteServerRequest.toJson(req))
  return DeleteServerResponse.fromJson((res.data as BaseResponse).body)
}

export const initServer = async (req: InitServerRequest) => {
  const res = await http.post(API_PATH + '/server/init', InitServerRequest.toJson(req))
  return InitServerResponse.fromJson((res.data as BaseResponse).body)
}
