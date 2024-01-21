import { HeaderOperations } from './common'

export interface ClientPluginOptions {}

export interface TypedClientPluginOptions {
  type: string
  clientPluginOptions?: ClientPluginOptions
}

export interface HTTP2HTTPSPluginOptions {
  type?: string
  localAddr?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
}

export interface HTTPProxyPluginOptions {
  type?: string
  httpUser?: string
  httpPassword?: string
}

export interface HTTPS2HTTPPluginOptions {
  type?: string
  localAddr?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
  crtPath?: string
  keyPath?: string
}

export interface HTTPS2HTTPSPluginOptions {
  type?: string
  localAddr?: string
  hostHeaderRewrite?: string
  requestHeaders?: HeaderOperations
  crtPath?: string
  keyPath?: string
}

export interface Socks5PluginOptions {
  type?: string
  username?: string
  password?: string
}

export interface StaticFilePluginOptions {
  type?: string
  localPath?: string
  stripPrefix?: string
  httpUser?: string
  httpPassword?: string
}

export interface UnixDomainSocketPluginOptions {
  type?: string
  unixPath?: string
}
