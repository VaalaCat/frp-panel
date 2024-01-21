export interface QUICOptions {
  keepalivePeriod?: number
  maxIdleTimeout?: number
  maxIncomingStreams?: number
}

export interface WebServerConfig {
  addr?: string
  port?: number
  user?: string
  password?: string
  assetsDir?: string
  pprofEnable?: boolean
  tls?: TLSConfig
}

export interface TLSConfig {
  certFile?: string
  keyFile?: string
  trustedCaFile?: string
  serverName?: string
}

export interface LogConfig {
  to?: string
  level?: string
  maxDays: number
  disablePrintColor?: boolean
}

export interface HTTPPluginOptions {
  name: string
  addr: string
  path: string
  ops: string[]
  tlsVerify?: boolean
}

export interface HeaderOperations {
  set?: { [key: string]: string }
}

export type AuthMethod = 'token' | 'oidc'

export const AuthMethodToken: AuthMethod = 'token'
export const AuthMethodOIDC: AuthMethod = 'oidc'

export type AuthScope = 'HeartBeats' | 'NewWorkConns'

export const AuthScopeHeartBeats: AuthScope = 'HeartBeats'
export const AuthScopeNewWorkConns: AuthScope = 'NewWorkConns'

export interface PortsRange {
  start?: number
  end?: number
  single?: number
}

export type BandwidthUnit = 'MB' | 'KB'

export interface BandwidthQuantity {
  s: BandwidthUnit // MB or KB
  i: number // bytes
}
