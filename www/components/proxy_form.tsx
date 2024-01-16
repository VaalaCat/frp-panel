import { HTTPProxyConfig, TCPProxyConfig, UDPProxyConfig } from "@/types/proxy"
import * as z from "zod"
import React from "react"
import { ZodIPSchema, ZodPortSchema, ZodStringSchema } from "@/lib/consts"
import { useEffect, useState } from "react"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import { Button } from "@/components/ui/button"
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form"
import { Input } from "@/components/ui/input"
import { useMutation } from '@tanstack/react-query'
import { updateFRPC } from "@/api/frp"
import { useToast } from "./ui/use-toast"
import { RespCode } from "@/lib/pb/common"
import { ClientConfig } from "@/types/client"
export const TCPConfigSchema = z.object({
    remotePort: ZodPortSchema,
    name: ZodStringSchema,
    localIP: ZodIPSchema.default("127.0.0.1"),
    localPort: ZodPortSchema,
})

export const UDPConfigSchema = z.object({
    remotePort: ZodPortSchema.optional(),
    name: ZodStringSchema,
    localIP: ZodIPSchema.default("127.0.0.1"),
    localPort: ZodPortSchema,
})

export const HTTPConfigSchema = z.object({
    name: ZodStringSchema,
    localPort: ZodPortSchema,
    localIP: ZodIPSchema.default("127.0.0.1"),
    subDomain: ZodStringSchema,
})

export interface ProxyFormProps {
    clientID: string
    serverID: string
}

export const TCPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID }) => {
    const [_, setTCPConfig] = useState<TCPProxyConfig | undefined>()
    const { toast } = useToast()
    const form = useForm<z.infer<typeof TCPConfigSchema>>({
        resolver: zodResolver(TCPConfigSchema),
    })
    const updateFrpc = useMutation({ mutationFn: updateFRPC, })

    useEffect(() => {
        setTCPConfig(undefined)
        form.reset({})
    }, [])

    const onSubmit = async (values: z.infer<typeof TCPConfigSchema>) => {
        setTCPConfig({ type: "tcp", ...values })
        try {
            const res = await updateFrpc.mutateAsync({
                config: Buffer.from(JSON.stringify({
                    proxies: [{ ...values, type: "tcp" }]
                } as ClientConfig)), serverId: serverID, clientId: clientID
            })
            toast({ title: "创建隧道状态", description: res.status?.code === RespCode.SUCCESS ? "创建成功" : "创建失败" })
        } catch (error) {
            console.error(error)
            toast({ title: "创建隧道状态", description: "创建失败" })
        }
    }

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 隧道名称 </FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue="tcptunnel"
                />
                <FormField
                    control={form.control}
                    name="localPort"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 本地端口 </FormLabel>
                            <FormControl>
                                <Input type="number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue={1234}
                />
                <FormField
                    control={form.control}
                    name="localIP"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 转发地址 </FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue="127.0.0.1"
                />
                <FormField
                    control={form.control}
                    name="remotePort"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 远端端口 </FormLabel>
                            <FormControl>
                                <Input type="number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue={4321}
                />
                <Button type="submit">提交</Button>
            </form>
        </Form>
    );
}

export const UDPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID }) => {
    const [_, setUDPConfig] = useState<UDPProxyConfig | undefined>()
    const { toast } = useToast()
    const form = useForm<z.infer<typeof UDPConfigSchema>>({
        resolver: zodResolver(UDPConfigSchema),
    })

    useEffect(() => {
        setUDPConfig(undefined)
        form.reset({})
    }, [])

    const updateFrpc = useMutation({ mutationFn: updateFRPC })

    const onSubmit = async (values: z.infer<typeof UDPConfigSchema>) => {
        setUDPConfig({ ...values, type: "udp" })
        try {
            const res = await updateFrpc.mutateAsync({
                config: Buffer.from(JSON.stringify({
                    proxies: [{ ...values }]
                } as ClientConfig)), serverId: serverID, clientId: clientID
            })
            toast({ title: "创建隧道状态", description: res.status?.code === RespCode.SUCCESS ? "创建成功" : "创建失败" })
        } catch (error) {
            toast({ title: "创建隧道状态", description: `创建失败: ${error}` })
        }
    }
    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 隧道名称 </FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue="udptunnel"
                />
                <FormField
                    control={form.control}
                    name="localPort"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 本地端口 </FormLabel>
                            <FormControl>
                                <Input type="number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue={1234}
                />
                <FormField
                    control={form.control}
                    name="localIP"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 转发地址 </FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue="127.0.0.1"
                />
                <FormField
                    control={form.control}
                    name="remotePort"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 远端端口 </FormLabel>
                            <FormControl>
                                <Input type="number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>)}
                    defaultValue={4321}
                />
                <Button type="submit">提交</Button>
            </form>
        </Form>
    );
}

export const HTTPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID }) => {
    const [_, setHTTPConfig] = useState<HTTPProxyConfig | undefined>()
    const { toast } = useToast()
    const form = useForm<z.infer<typeof HTTPConfigSchema>>({
        resolver: zodResolver(HTTPConfigSchema),
    })

    useEffect(() => {
        setHTTPConfig(undefined)
        form.reset({})
    }, [])

    const updateFrpc = useMutation({ mutationFn: updateFRPC })

    const onSubmit = async (values: z.infer<typeof HTTPConfigSchema>) => {
        setHTTPConfig({ ...values, type: "http" })
        const conf = {
            name: values.name,
            type: "http",
            subDomain: values.subDomain,
            localIP: values.localIP,
            localPort: values.localPort,
        } as HTTPProxyConfig
        try {
            const res = await updateFrpc.mutateAsync({
                config: Buffer.from(JSON.stringify({
                    proxies: [conf]
                } as ClientConfig)), serverId: serverID, clientId: clientID
            })
            toast({ title: "创建隧道状态", description: res.status?.code === RespCode.SUCCESS ? "创建成功" : "创建失败" })
        } catch (error) {
            toast({ title: "创建隧道状态", description: `创建失败: ${error}` })
        }
    }

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 隧道名称 </FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue="httptunnel"
                />
                <FormField
                    control={form.control}
                    name="localPort"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 本地端口 </FormLabel>
                            <FormControl>
                                <Input type="number" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue={1234}
                />
                <FormField
                    control={form.control}
                    name="localIP"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 转发地址 </FormLabel>
                            <FormControl>
                                <Input {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                    defaultValue="127.0.0.1"
                />
                <FormField
                    control={form.control}
                    name="subDomain"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel> 远端子域名 </FormLabel>
                            <FormControl>
                                <Input type="text" {...field} />
                            </FormControl>
                            <FormMessage />
                        </FormItem>)}
                    defaultValue={"sub"}
                />
                <Button type="submit">提交</Button>
            </form>
        </Form>
    )
}