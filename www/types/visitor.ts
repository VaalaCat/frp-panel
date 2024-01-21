export interface VisitorTransport {
  useEncryption?: boolean
  useCompression?: boolean
}

export interface VisitorBaseConfig {
  name: string
  type: string
  transport?: VisitorTransport
  secretKey?: string
  serverUser?: string
  serverName?: string
  bindAddr?: string
  bindPort?: number
}

export type VisitorType = 'stcp' | 'xtcp' | 'sudp'

export type TypedVisitorConfig = STCPVisitorConfig | SUDPVisitorConfig | XTCPVisitorConfig

export interface STCPVisitorConfig extends VisitorBaseConfig {
  type: 'stcp'
}

export interface SUDPVisitorConfig extends VisitorBaseConfig {
  type: 'sudp'
}

export interface XTCPVisitorConfig extends VisitorBaseConfig {
  type: 'xtcp'
  protocol?: string
  keepTunnelOpen?: boolean
  maxRetriesAnHour?: number
  minRetryInterval?: number
  fallbackTo?: string
  fallbackTimeoutMs?: number
}

export interface ClientCommonConfig {
  user: string
}
