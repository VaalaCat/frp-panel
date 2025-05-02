import { GetPlatformInfoResponse } from '@/lib/pb/api_user'
import { User } from '@/lib/pb/common'
import { atom } from 'nanostores'
import { persistentAtom } from '@nanostores/persistent'
import { LOCAL_STORAGE_TOKEN_KEY } from '@/lib/consts'

export const $userInfo = atom<User | undefined>()
export const $statusOnline = atom<boolean>(false)
export const $token = persistentAtom<string | undefined>(LOCAL_STORAGE_TOKEN_KEY)
export const $platformInfo = atom<GetPlatformInfoResponse | undefined>()

// 创建持久化的语言设置
export const $language = persistentAtom<string>('user-language', 'zh', {
  encode: JSON.stringify,
  decode: JSON.parse,
})

export const $useServerGithubProxyUrl = persistentAtom<boolean>('use_server_github_proxy_url', false, {
  encode: JSON.stringify,
  decode: JSON.parse,
})
