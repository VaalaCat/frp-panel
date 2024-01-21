import { AuthMethod, AuthScope, LogConfig, QUICOptions, TLSConfig, WebServerConfig } from './common'
import { TypedProxyConfig } from './proxy'
import { TypedVisitorConfig } from './visitor'

export interface AuthOIDCClientConfig {
  clientID?: string
  clientSecret?: string
  audience?: string
  scope?: string
  tokenEndpointURL?: string
  additionalEndpointParams?: { [key: string]: string }
}

export interface AuthClientConfig {
  method?: AuthMethod
  additionalScopes?: AuthScope[]
  token?: string
  oidc?: AuthOIDCClientConfig
}

export interface ClientTransportConfig {
  protocol?: string
  dialServerTimeout?: number
  dialServerKeepAlive?: number
  connectServerLocalIP?: string
  proxyURL?: string
  poolCount?: number
  tcpMux?: boolean
  tcpMuxKeepaliveInterval?: number
  quic?: QUICOptions
  heartbeatInterval?: number
  heartbeatTimeout?: number
  tls?: TLSClientConfig
}

export interface TLSClientConfig {
  enable?: boolean
  disableCustomTLSFirstByte?: boolean
  tls?: TLSConfig
}

export interface CompleteTLSClientConfig extends TLSClientConfig {
  enable: boolean
  disableCustomTLSFirstByte: boolean
}

export interface AuthClientConfig {
  auth?: AuthClientConfig
  user?: string
  serverAddr?: string
  serverPort?: number
  natHoleStunServer?: string
  dnsServer?: string
  loginFailExit?: boolean
  start?: string[]
  log?: LogConfig
  webServer?: WebServerConfig
  transport?: ClientTransportConfig
  udpPacketSize?: number
  metadatas?: { [key: string]: string }
  includes?: string[]
}

export interface ClientConfig extends ClientCommonConfig {
  proxies?: TypedProxyConfig[]
  visitors?: TypedVisitorConfig[]
}

export interface ClientCommonConfig extends AuthClientConfig {
  auth?: AuthClientConfig
  user?: string
  serverAddr: string
  serverPort: number
  natHoleStunServer?: string
  dnsServer?: string
  loginFailExit?: boolean
  start?: string[]
  log?: LogConfig
  webServer?: WebServerConfig
  transport?: ClientTransportConfig
  udpPacketSize?: number
  metadatas?: { [key: string]: string }
  includes?: string[]
}
