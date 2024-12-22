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

export const ZodIPSchema = z.string({ required_error: 'validation.required' })
  .regex(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/, { message: 'validation.ipAddress' })
export const ZodStringSchema = z.string({ required_error: 'validation.required' })
  .min(1, { message: 'validation.required' })

export const ZodStringOptionalSchema = z.string().optional()
export const ZodEmailSchema = z.string({ required_error: 'validation.required' })
  .min(1, { message: 'validation.required' })
  .email({ message: 'auth.email.invalid' })

export const ConnectionProtocols = ["tcp", "kcp", "quic", "websocket", "wss"]

export const TypedProxyConfigValid = (typedProxyCfg: TypedProxyConfig | undefined): boolean => {
  return (typedProxyCfg?.localPort && typedProxyCfg.localIP && typedProxyCfg.name && typedProxyCfg.type) ? true : false
}

export const IsIDValid = (clientID: string | undefined): boolean => {
  if (clientID == undefined) {
    return false
  }
  const regex = /^[a-zA-Z0-9-_]+$/;
  return clientID.length > 0 && regex.test(clientID);
}

export const ClientConfigured = (client: Client | undefined): boolean => {
  if (client == undefined) {
    return false
  }
  return !((client.config == undefined || client.config == '') &&
    (client.clientIds == undefined || client.clientIds.length == 0))
}

// .refine((e) => e === "abcd@fg.com", "This email is not in our database")

export const ExecCommandStr = <T extends Client | Server>(
  type: "client" | "server",
  item: T,
  info: GetPlatformInfoResponse,
  fileName?: string,
) => {
  return `${fileName || 'frp-panel'} ${type} -s ${item.secret} -i ${item.id} -a ${info.globalSecret} -r ${info.masterRpcHost
    } -c ${info.masterRpcPort} -p ${info.masterApiPort} -e ${info.masterApiScheme}`
}

export const WindowsInstallCommand = <T extends Client | Server>(
  type: "client" | "server",
  item: T,
  info: GetPlatformInfoResponse,
) => {
  return `[Net.ServicePointManager]::SecurityProtocol = ` +
    `[Net.SecurityProtocolType]::Ssl3 -bor ` +
    `[Net.SecurityProtocolType]::Tls -bor ` +
    `[Net.SecurityProtocolType]::Tls11 -bor ` +
    `[Net.SecurityProtocolType]::Tls12;set-ExecutionPolicy RemoteSigned;` +
    `Invoke-WebRequest https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.ps1 ` +
    `-OutFile C:\install.ps1;powershell.exe C:\install.ps1 ${ExecCommandStr(type, item, info, ' ')}`
}

export const LinuxInstallCommand = <T extends Client | Server>(
  type: "client" | "server",
  item: T,
  info: GetPlatformInfoResponse,
) => {
  return `curl -sSL https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --${ExecCommandStr(type, item, info, ' ')}`
}

export const ClientEnvFile = <T extends Client | Server>(
  item: T,
  info: GetPlatformInfoResponse,
) => {
  return `CLIENT_ID=${item.id}
CLIENT_SECRET=${item.secret}
APP_SECRET=${info.globalSecret}
MASTER_RPC_HOST=${info.masterRpcHost}
MASTER_RPC_PORT=${info.masterRpcPort}
MASTER_API_HOST=${info.masterRpcHost}
MASTER_API_PORT=${info.masterApiPort}
MASTER_API_SCHEME=${info.masterApiScheme}`
}