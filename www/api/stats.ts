import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import { GetProxyStatsByClientIDRequest, GetProxyStatsByClientIDResponse } from '@/lib/pb/api_client'
import { GetProxyStatsByServerIDRequest, GetProxyStatsByServerIDResponse } from '@/lib/pb/api_server'
import { BaseResponse } from '@/types/api'

export const getProxyStatsByClientID = async (req: GetProxyStatsByClientIDRequest) => {
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
  const res = await http.post(API_PATH + '/proxy/get_by_cid', GetProxyStatsByClientIDRequest.toJson(req))
  return GetProxyStatsByClientIDResponse.fromJson((res.data as BaseResponse).body)
}

export const getProxyStatsByServerID = async (req: GetProxyStatsByServerIDRequest) => {
  const res = await http.post(API_PATH + '/proxy/get_by_sid', GetProxyStatsByServerIDRequest.toJson(req))
  return GetProxyStatsByServerIDResponse.fromJson((res.data as BaseResponse).body)
}