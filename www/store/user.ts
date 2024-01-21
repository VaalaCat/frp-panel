import { GetPlatformInfoResponse } from '@/lib/pb/api_user'
import { User } from '@/lib/pb/common'
import { atom } from 'nanostores'

export const $userInfo = atom<User | undefined>()
export const $statusOnline = atom<boolean>(false)
export const $token = atom<string | undefined>()
export const $platformInfo = atom<GetPlatformInfoResponse | undefined>()
