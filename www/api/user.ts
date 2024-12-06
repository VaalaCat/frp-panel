import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
  GetUserInfoRequest,
  GetUserInfoResponse,
  UpdateUserInfoRequest,
  UpdateUserInfoResponse,
} from '@/lib/pb/api_user'
import { $statusOnline, $userInfo } from '@/store/user'
import { BaseResponse } from '@/types/api'

export const getUserInfo = async (req: GetUserInfoRequest) => {
  const res = await http.post(API_PATH + '/user/get', GetUserInfoRequest.toJson(req))
  $userInfo.set(GetUserInfoResponse.fromJson((res.data as BaseResponse).body).userInfo)
  $statusOnline.set(!!GetUserInfoResponse.fromJson((res.data as BaseResponse).body).userInfo)
  return GetUserInfoResponse.fromJson((res.data as BaseResponse).body)
}

export const updateUserInfo = async (req: UpdateUserInfoRequest) => {
  const res = await http.post(API_PATH + '/user/update', UpdateUserInfoRequest.toJson(req))
  return UpdateUserInfoResponse.fromJson((res.data as BaseResponse).body)
}
