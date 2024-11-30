import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import { GetProxyByCIDRequest, GetProxyByCIDResponse } from '@/lib/pb/api_client'
import { GetProxyBySIDRequest, GetProxyBySIDResponse } from '@/lib/pb/api_server'
import { BaseResponse } from '@/types/api'

export const getProxyStatsByClientID = async (req: GetProxyByCIDRequest) => {
  // return {
  //   proxyInfos: [
  //     {
  //       name: "test",
  //       historyTrafficIn: BigInt(1024 * 1024 * 1024),
  //       historyTrafficOut: BigInt(1024 * 1024 * 1024 * 2),
  //       todayTrafficIn: BigInt(1024 * 1024 * 1024 * 4),
  //       todayTrafficOut: BigInt(1024 * 1024 * 1024 * 5),
  //     },
  //   ],
  // } as GetProxyByCIDResponse
  const res = await http.post(API_PATH + '/proxy/get_by_cid', GetProxyByCIDRequest.toJson(req))
  return GetProxyByCIDResponse.fromJson((res.data as BaseResponse).body)
}

export const getProxyStatsByServerID = async (req: GetProxyBySIDRequest) => {
  const res = await http.post(API_PATH + '/proxy/get_by_sid', GetProxyBySIDRequest.toJson(req))
  return GetProxyBySIDResponse.fromJson((res.data as BaseResponse).body)
}