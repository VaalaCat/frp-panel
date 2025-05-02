import * as z from 'zod'
import { Client, Server } from './pb/common'
import { GetPlatformInfoResponse } from './pb/api_user'
import { TypedProxyConfig } from '@/types/proxy'

export const API_PATH = '/api/v1'
export const SET_TOKEN_HEADER = 'x-set-authorization'
export const X_CLIENT_REQUEST_ID = 'x-client-request-id'
export const LOCAL_STORAGE_TOKEN_KEY = 'token'
export const ZodPortSchema = z.coerce
  .number({ required_error: 'validation.required' })
  .min(1, { message: 'validation.portRange.min' })
  .max(65535, { message: 'validation.portRange.max' })

export const ZodPortOptionalSchema = z.coerce
  .number({ required_error: 'validation.required' })
  .min(1, { message: 'validation.portRange.min' })
  .max(65535, { message: 'validation.portRange.max' })

export const ZodIPSchema = z
  .string({ required_error: 'validation.required' })
  .regex(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/, { message: 'validation.ipAddress' })
export const ZodStringSchema = z
  .string({ required_error: 'validation.required' })
  .min(1, { message: 'validation.required' })

export const ZodStringOptionalSchema = z.string().optional()
export const ZodEmailSchema = z
  .string({ required_error: 'validation.required' })
  .min(1, { message: 'validation.required' })
  .email({ message: 'auth.email.invalid' })

export const ConnectionProtocols = ['tcp', 'kcp', 'quic', 'websocket', 'wss']

export const TypedProxyConfigValid = (typedProxyCfg: TypedProxyConfig | undefined): boolean => {
  if (!typedProxyCfg) {
    return false
  }

  if (typedProxyCfg.plugin && typedProxyCfg.plugin.type) {
    if (typedProxyCfg.type === 'tcp' || typedProxyCfg.type === 'udp') {
      if (!typedProxyCfg.remotePort) {
        console.log('remotePort is undefined')
        return false
      }
    }
    return typedProxyCfg.name && typedProxyCfg.type ? true : false
  }

  if (typedProxyCfg.type === 'tcp' || typedProxyCfg.type === 'udp') {
    if (!typedProxyCfg.remotePort) {
      console.log('remotePort is undefined')
      return false
    }
  }

  return typedProxyCfg?.localPort && typedProxyCfg.localIP && typedProxyCfg.name && typedProxyCfg.type ? true : false
}

export const IsIDValid = (clientID: string | undefined): boolean => {
  if (clientID == undefined) {
    return false
  }
  const regex = /^[a-zA-Z0-9-_]+$/
  return clientID.length > 0 && regex.test(clientID)
}

export const ClientConfigured = (client: Client | undefined): boolean => {
  if (client == undefined) {
    return false
  }
  return !(
    (client.config == undefined || client.config == '') &&
    (client.clientIds == undefined || client.clientIds.length == 0)
  )
}

// .refine((e) => e === "abcd@fg.com", "This email is not in our database")

export const ExecCommandStr = <T extends Client | Server>(
  type: 'client' | 'server',
  item: T,
  info: GetPlatformInfoResponse,
  fileName?: string,
) => {
  return `${fileName || 'frp-panel'} ${type} -s ${item.secret} -i ${item.id} --api-url ${info.clientApiUrl} --rpc-url ${info.clientRpcUrl}`
}

export const JoinCommandStr = (info: GetPlatformInfoResponse, token: string, fileName?: string, clientID?: string) => {
  return `${fileName || 'frp-panel'} join${clientID ? ` -i ${clientID}` : ''} -j ${token} --api-url ${info.clientApiUrl} --rpc-url ${info.clientRpcUrl}`
}

export const WindowsInstallCommand = <T extends Client | Server>(
  type: 'client' | 'server',
  item: T,
  info: GetPlatformInfoResponse,
  github_proxy?: boolean,
) => {
  return (
    `[Net.ServicePointManager]::SecurityProtocol = ` +
    `[Net.SecurityProtocolType]::Ssl3 -bor ` +
    `[Net.SecurityProtocolType]::Tls -bor ` +
    `[Net.SecurityProtocolType]::Tls11 -bor ` +
    `[Net.SecurityProtocolType]::Tls12;set-ExecutionPolicy RemoteSigned;` +
    `Invoke-WebRequest ${github_proxy ? info.githubProxyUrl : ''}https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.ps1 ` +
    `-OutFile C:\install.ps1;powershell.exe C:\install.ps1 ${ExecCommandStr(type, item, info, ' ')}`
  )
}

export const LinuxInstallCommand = <T extends Client | Server>(
  type: 'client' | 'server',
  item: T,
  info: GetPlatformInfoResponse,
  github_proxy?: boolean,
) => {
  return `curl -fSL ${github_proxy ? info.githubProxyUrl : ''}https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --${ExecCommandStr(type, item, info, ' ')}`
}

export const ClientEnvFile = <T extends Client | Server>(item: T, info: GetPlatformInfoResponse) => {
  return `CLIENT_ID=${item.id}
CLIENT_SECRET=${item.secret}
CLIENT_API_URL=${info.clientApiUrl}
CLIENT_RPC_URL=${info.clientRpcUrl}`
}
