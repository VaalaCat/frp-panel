import { BandwidthQuantity, HeaderOperations } from './common'
import { TypedClientPluginOptions } from './plugin'

export interface ProxyTransport {
  useEncryption?: boolean
  useCompression?: boolean
  bandwidthLimit?: BandwidthQuantity
  bandwidthLimitMode?: string
  proxyProtocolVersion?: string
}

export interface LoadBalancerConfig {
  group: string
  groupKey?: string
}

export interface ProxyBackend {
  localIP?: string
  localPort?: number
  plugin?: TypedClientPluginOptions
}

export interface HealthCheckConfig {
  type: string
  timeoutSeconds?: number
  maxFailed?: number
  intervalSeconds: number
  path?: string
}

export interface DomainConfig {
  customDomains?: string[]
  subdomain?: string
}

export interface ProxyBaseConfig {
  name: string
  type: string
  transport?: ProxyTransport
  metadatas?: { [key: string]: string }
  loadBalancer?: LoadBalancerConfig
  healthCheck?: HealthCheckConfig
  localIP?: string
  localPort?: number
  plugin?: TypedClientPluginOptions
}

export type TypedProxyConfig =
  | TCPProxyConfig
  | UDPProxyConfig
  | HTTPProxyConfig
  | HTTPSProxyConfig
  | TCPMuxProxyConfig
  | STCPProxyConfig
  | XTCPProxyConfig
  | SUDPProxyConfig

export type ProxyType = 'tcp' | 'udp' | 'tcpmux' | 'http' | 'https' | 'stcp' | 'xtcp' | 'sudp'

export interface TCPProxyConfig extends ProxyBaseConfig {
  type: 'tcp'
  remotePort?: number
}

export interface UDPProxyConfig extends ProxyBaseConfig {
  type: 'udp'
  remotePort?: number
}

export interface HTTPProxyConfig extends ProxyBaseConfig, DomainConfig {
  type: 'http'
  locations?: string[]
  httpUser?: string
  httpPassword?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
  routeByHTTPUser?: string
}

export interface HTTPSProxyConfig extends ProxyBaseConfig, DomainConfig {
  type: 'https'
}

export type TCPMultiplexerType = 'httpconnect'

export interface TCPMuxProxyConfig extends ProxyBaseConfig, DomainConfig {
  type: 'tcpmux'
  httpUser?: string
  httpPassword?: string
  routeByHTTPUser?: string
  multiplexer?: string
}

export interface STCPProxyConfig extends ProxyBaseConfig {
  type: 'stcp'
  secretKey?: string
  allowUsers?: string[]
}

export interface XTCPProxyConfig extends ProxyBaseConfig {
  type: 'xtcp'
  secretKey?: string
  allowUsers?: string[]
}

export interface SUDPProxyConfig extends ProxyBaseConfig {
  type: 'sudp'
  secretKey?: string
  allowUsers?: string[]
}
