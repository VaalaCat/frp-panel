import { atom } from 'nanostores'

export const $clientTableRefetchTrigger = atom<number>(0)
export const $serverTableRefetchTrigger = atom<number>(0)
export const $proxyTableRefetchTrigger = atom<number>(0)