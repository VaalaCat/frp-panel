import * as z from "zod"
import { Client, Server } from "./pb/common"

export const API_PATH = '/api/v1'
export const SET_TOKEN_HEADER = 'x-set-authorization'
export const X_CLIENT_REQUEST_ID = 'x-client-request-id'
export const LOCAL_STORAGE_TOKEN_KEY = 'token'
export const ZodPortSchema = z.coerce.number().min(1, {
	message: "端口号不能小于 1",
}).max(65535, { message: "端口号不能大于 65535" })
export const ZodIPSchema = z.string().regex(/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/,
	{ message: "请输入正确的IP地址" })
export const ZodStringSchema = z.string().min(1,
	{ message: "不能为空" })
export const ZodEmailSchema = z.string()
	.min(1, { message: "不能为空" })
	.email("是不是输错了邮箱地址呢?")
// .refine((e) => e === "abcd@fg.com", "This email is not in our database")

export const ExecCommandStr = <T extends Client | Server>(type: string, item: T) => {
	return `frp-panel ${type} -s ${item.secret} -i ${item.id}`
}