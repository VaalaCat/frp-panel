import * as z from 'zod'
import { Client, Server } from './pb/common'
import { GetPlatformInfoResponse } from './pb/api_user'

export const API_PATH = '/api/v1'
export const SET_TOKEN_HEADER = 'x-set-authorization'
export const X_CLIENT_REQUEST_ID = 'x-client-request-id'
export const LOCAL_STORAGE_TOKEN_KEY = 'token'
export const ZodPortSchema = z.coerce
  .number()
  .min(1, {
    message: '端口号不能小于 1',
  })
  .max(65535, { message: '端口号不能大于 65535' })
export const ZodIPSchema = z.string().regex(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/, { message: '请输入正确的IP地址' })
export const ZodStringSchema = z.string().min(1, { message: '不能为空' })
export const ZodEmailSchema = z.string().min(1, { message: '不能为空' }).email('是不是输错了邮箱地址呢?')
// .refine((e) => e === "abcd@fg.com", "This email is not in our database")

export const ExecCommandStr = <T extends Client | Server>(
  type: string,
  item: T,
  info: GetPlatformInfoResponse,
  fileName?: string,
) => {
  return `${fileName || 'frp-panel'} ${type} -s ${item.secret} -i ${item.id} -a ${info.globalSecret} -r ${
    info.masterRpcHost
  } -c ${info.masterRpcPort} -p ${info.masterApiPort} -e ${info.masterApiScheme}`
}

export const WindowsInstallCommand = <T extends Client | Server>(
  type: string,
  item: T,
  info: GetPlatformInfoResponse,
) => {
  return `Invoke-WebRequest -Uri 'https://github.com/your_repository/frp-panel/releases/latest/download/frp-panel-amd64.exe' -OutFile 'frp-panel.exe'
	Move-Item .\\frp-panel.exe C:\\Tools\\frp-panel.exe
	$command = "C:\\Tools\\${ExecCommandStr(type, item, info, 'frp-panel.exe')}"
	Set-ItemProperty -Path 'HKLM:\\SYSTEM\\CurrentControlSet\\Services\\FRPPanel' -Name 'ImagePath' -Value "\`"$command\`""
	New-Service -Name 'FRPPanel' -BinaryPathName 'C:\\Tools\\frp-panel.exe' -StartupType Automatic | Start-Service`
}

export const LinuxInstallCommand = <T extends Client | Server>(
  type: string,
  item: T,
  info: GetPlatformInfoResponse,
) => {
  return `curl -sSL https://raw.githubusercontent.com/VaalaCat/frp-panel/main/install.sh | bash -s --${ExecCommandStr(type, item, info, ' ')}`
}
