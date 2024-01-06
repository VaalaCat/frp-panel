import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
    LoginRequest, LoginResponse,
    RegisterRequest, RegisterResponse
} from '@/lib/pb/api_auth'
import { BaseResponse } from '@/types/api'

export const login = async (req: LoginRequest) => {
    const res = await http.post(API_PATH + '/auth/login', LoginRequest.toJson(req))
    return LoginResponse.fromJson((res.data as BaseResponse).body)
}

export const register = async (req: RegisterRequest) => {
    const res = await http.post(API_PATH + '/auth/register', RegisterRequest.toJson(req))
    return RegisterResponse.fromJson((res.data as BaseResponse).body)
}