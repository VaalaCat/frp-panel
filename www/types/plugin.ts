import { HeaderOperations } from './common'

export type ClientPluginType =
  | 'http_proxy'
  | 'http2https'
  | 'https2http'
  | 'https2https'
  | 'socks5'
  | 'static_file'
  | 'unix_domain_socket'

export type TypedClientPluginOptions =
| HTTP2HTTPSPluginOptions
| HTTPProxyPluginOptions
| HTTPS2HTTPPluginOptions
| HTTPS2HTTPSPluginOptions
| Socks5PluginOptions
| StaticFilePluginOptions
| UnixDomainSocketPluginOptions

export interface HTTP2HTTPSPluginOptions {
  type: 'http2https'
  localAddr?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
}

export interface HTTPProxyPluginOptions {
  type: 'http_proxy'
  httpUser?: string
  httpPassword?: string
}

export interface HTTPS2HTTPPluginOptions {
  type: 'https2http'
  localAddr?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
  crtPath?: string
  keyPath?: string
}

export interface HTTPS2HTTPSPluginOptions {
  type: 'https2https'
  localAddr?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
  crtPath?: string
  keyPath?: string
}

export interface Socks5PluginOptions {
  type: 'socks5'
  username?: string
  password?: string
}

export interface StaticFilePluginOptions {
  type: 'static_file'
  localPath?: string
  stripPrefix?: string
  httpUser?: string
  httpPassword?: string
}

export interface UnixDomainSocketPluginOptions {
  type: 'unix_domain_socket'
  unixPath?: string
}
