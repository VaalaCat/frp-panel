import { TypedProxyConfig } from '@/types/proxy'
import { atom } from 'nanostores'

export const $clientProxyConfigs = atom<TypedProxyConfig[]>([])
