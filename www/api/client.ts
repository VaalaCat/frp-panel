import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
  DeleteClientRequest,
  DeleteClientResponse,
  GetClientRequest,
  GetClientResponse,
  InitClientRequest,
  InitClientResponse,
  ListClientsRequest,
  ListClientsResponse,
} from '@/lib/pb/api_client'
import { BaseResponse } from '@/types/api'

export const getClient = async (req: GetClientRequest) => {
  const res = await http.post(API_PATH + '/client/get', GetClientRequest.toJson(req))
  return GetClientResponse.fromJson((res.data as BaseResponse).body)
}

export const listClient = async (req: ListClientsRequest) => {
  const res = await http.post(API_PATH + '/client/list', ListClientsRequest.toJson(req))
  return ListClientsResponse.fromJson((res.data as BaseResponse).body)
}

export const deleteClient = async (req: DeleteClientRequest) => {
  const res = await http.post(API_PATH + '/client/delete', DeleteClientRequest.toJson(req))
  return DeleteClientResponse.fromJson((res.data as BaseResponse).body)
}

export const initClient = async (req: InitClientRequest) => {
  console.log('attempting init client:', InitClientRequest.toJsonString(req))
  const res = await http.post(API_PATH + '/client/init', InitClientRequest.toJson(req))
  return InitClientResponse.fromJson((res.data as BaseResponse).body)
}
