import { AuthMethod, AuthScope, HTTPPluginOptions, LogConfig, PortsRange, QUICOptions, WebServerConfig } from './common'

export interface ServerConfig {
  auth?: AuthServerConfig
  bindAddr?: string
  bindPort?: number
  kcpBindPort?: number
  quicBindPort?: number
  proxyBindAddr?: string
  vhostHTTPPort?: number
  vhostHTTPTimeout?: number
  vhostHTTPSPort?: number
  tcpmuxHTTPConnectPort?: number
  tcpmuxPassthrough?: boolean
  subDomainHost?: string
  custom404Page?: string
  sshTunnelGateway?: SSHTunnelGateway
  webServer?: WebServerConfig
  enablePrometheus?: boolean
  log?: LogConfig
  transport?: ServerTransportConfig
  detailedErrorsToClient?: boolean
  maxPortsPerClient?: number
  userConnTimeout?: number
  udpPacketSize?: number
  natholeAnalysisDataReserveHours?: number
  allowPorts?: PortsRange[]
  httpPlugins?: HTTPPluginOptions[]
}

export interface AuthServerConfig {
  method?: AuthMethod
  additionalScopes?: AuthScope[]
  token?: string
  oidc?: AuthOIDCServerConfig
}

export interface AuthOIDCServerConfig {
  issuer?: string
  audience?: string
  skipExpiryCheck?: boolean
  skipIssuerCheck?: boolean
}

export interface ServerTransportConfig {
  tcpMux?: boolean
  tcpMuxKeepaliveInterval?: number
  tcpKeepalive?: number
  maxPoolCount?: number
  heartbeatTimeout?: number
  quic?: QUICOptions
  tls?: TLSServerConfig
}

export interface TLSServerConfig {
  force?: boolean
}

export interface SSHTunnelGateway {
  bindPort?: number
  privateKeyFile?: string
  autoGenPrivateKeyPath?: string
  authorizedKeysFile?: string
}
