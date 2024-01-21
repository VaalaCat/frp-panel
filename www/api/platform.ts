import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import { GetClientsStatusRequest, GetClientsStatusResponse } from '@/lib/pb/api_master'
import { GetPlatformInfoResponse } from '@/lib/pb/api_user'
import { BaseResponse } from '@/types/api'

export const getPlatformInfo = async () => {
  const res = await http.get(API_PATH + '/platform/baseinfo')
  return GetPlatformInfoResponse.fromJson((res.data as BaseResponse).body)
}

export const getClientsStatus = async (req: GetClientsStatusRequest) => {
  const res = await http.post(API_PATH + '/platform/clientsstatus', GetClientsStatusRequest.toJson(req))
  return GetClientsStatusResponse.fromJson((res.data as BaseResponse).body)
}
