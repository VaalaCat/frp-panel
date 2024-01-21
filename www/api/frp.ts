import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import { RemoveFRPCRequest, RemoveFRPCResponse, UpdateFRPCRequest, UpdateFRPCResponse } from '@/lib/pb/api_client'
import { RemoveFRPSRequest, RemoveFRPSResponse, UpdateFRPSRequest, UpdateFRPSResponse } from '@/lib/pb/api_server'
import { BaseResponse } from '@/types/api'

export const updateFRPS = async (req: UpdateFRPSRequest) => {
  const res = await http.post(API_PATH + '/frps/update', UpdateFRPSRequest.toJson(req))
  return UpdateFRPSResponse.fromJson((res.data as BaseResponse).body)
}

export const removeFRPS = async (req: RemoveFRPSRequest) => {
  const res = await http.post(API_PATH + '/frps/remove', RemoveFRPSRequest.toJson(req))
  return RemoveFRPSResponse.fromJson((res.data as BaseResponse).body)
}

export const updateFRPC = async (req: UpdateFRPCRequest) => {
  const res = await http.post(API_PATH + '/frpc/update', UpdateFRPCRequest.toJson(req))
  return UpdateFRPCResponse.fromJson((res.data as BaseResponse).body)
}

export const removeFRPC = async (req: RemoveFRPCRequest) => {
  const res = await http.post(API_PATH + '/frpc/remove', RemoveFRPCRequest.toJson(req))
  return RemoveFRPCResponse.fromJson((res.data as BaseResponse).body)
}
