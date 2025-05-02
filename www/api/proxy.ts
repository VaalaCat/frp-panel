import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
  CreateProxyConfigRequest,
  CreateProxyConfigResponse,
  DeleteProxyConfigRequest,
  DeleteProxyConfigResponse,
  GetProxyConfigRequest,
  GetProxyConfigResponse,
  ListProxyConfigsRequest,
  ListProxyConfigsResponse,
  StartProxyRequest,
  StartProxyResponse,
  StopProxyRequest,
  StopProxyResponse,
  UpdateProxyConfigRequest,
  UpdateProxyConfigResponse,
} from '@/lib/pb/api_client'
import { BaseResponse } from '@/types/api'

export const createProxyConfig = async (req: CreateProxyConfigRequest) => {
  const res = await http.post(API_PATH + '/proxy/create_config', CreateProxyConfigRequest.toJson(req))
  return CreateProxyConfigResponse.fromJson((res.data as BaseResponse).body)
}

export const listProxyConfig = async (req: ListProxyConfigsRequest) => {
  const res = await http.post(API_PATH + '/proxy/list_configs', ListProxyConfigsRequest.toJson(req))
  return ListProxyConfigsResponse.fromJson((res.data as BaseResponse).body)
}

export const updateProxyConfig = async (req: UpdateProxyConfigRequest) => {
  const res = await http.post(API_PATH + '/proxy/update_config', UpdateProxyConfigRequest.toJson(req))
  return UpdateProxyConfigResponse.fromJson((res.data as BaseResponse).body)
}

export const deleteProxyConfig = async (req: DeleteProxyConfigRequest) => {
  const res = await http.post(API_PATH + '/proxy/delete_config', DeleteProxyConfigRequest.toJson(req))
  return DeleteProxyConfigResponse.fromJson((res.data as BaseResponse).body)
}

export const getProxyConfig = async (req: GetProxyConfigRequest) => {
  const res = await http.post(API_PATH + '/proxy/get_config', GetProxyConfigRequest.toJson(req))
  return GetProxyConfigResponse.fromJson((res.data as BaseResponse).body)
}

export const stopProxy = async (req: StopProxyRequest) => {
  const res = await http.post(API_PATH + '/proxy/stop_proxy', StopProxyRequest.toJson(req))
  return StopProxyResponse.fromJson((res.data as BaseResponse).body)
}

export const startProxy = async (req: StartProxyRequest) => {
  const res = await http.post(API_PATH + '/proxy/start_proxy', StartProxyRequest.toJson(req))
  return StartProxyResponse.fromJson((res.data as BaseResponse).body)
}
