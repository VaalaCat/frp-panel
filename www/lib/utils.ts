import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const BYTE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'] as const

export function formatBytes(bytes: number | bigint): string {
  if (typeof bytes === 'bigint') {
    if (bytes === BigInt(0)) return '0 B'
    const negative = bytes < BigInt(0)
    let value = negative ? -bytes : bytes
    const base = BigInt(1024)
    let unitIndex = 0
    let remainder = BigInt(0)
    while (value >= base && unitIndex < BYTE_UNITS.length - 1) {
      remainder = value % base
      value /= base
      unitIndex++
    }
    const fractional = unitIndex === 0 ? BigInt(0) : (remainder * BigInt(100)) / base
    const fractionStr = fractional > BigInt(0) ? `.${fractional.toString().padStart(2, '0').replace(/0+$/, '')}` : ''
    return `${negative ? '-' : ''}${value}${fractionStr} ${BYTE_UNITS[unitIndex]}`
  }
  if (bytes === 0) return '0 B'
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + BYTE_UNITS[i]
}

export function ObjToUint8Array(obj: any): Uint8Array {
  const buffer = Buffer.from(JSON.stringify(obj))
  return new Uint8Array(buffer.buffer, buffer.byteOffset, buffer.byteLength)
}

export function makeRandomTrigger() {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}