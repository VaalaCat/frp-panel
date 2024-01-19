import { HTTPProxyConfig, TCPProxyConfig, TypedProxyConfig, UDPProxyConfig } from "@/types/proxy"
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
import { useToast } from "./ui/use-toast"
import { $clientProxyConfigs } from "@/store/proxy"
import { useStore } from "@nanostores/react"
import { YesIcon } from "./ui/icon"
export const TCPConfigSchema = z.object({
    remotePort: ZodPortSchema,
    localIP: ZodIPSchema.default("127.0.0.1"),
    localPort: ZodPortSchema,
})

export const UDPConfigSchema = z.object({
    remotePort: ZodPortSchema.optional(),
    localIP: ZodIPSchema.default("127.0.0.1"),
    localPort: ZodPortSchema,
})

export const HTTPConfigSchema = z.object({
    localPort: ZodPortSchema,
    localIP: ZodIPSchema.default("127.0.0.1"),
    subDomain: ZodStringSchema,
})

export interface ProxyFormProps {
    clientID: string
    serverID: string
    proxyName: string
    defaultProxyConfig?: TypedProxyConfig
}

export const TCPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName }) => {
    const [_, setTCPConfig] = useState<TCPProxyConfig | undefined>()

    const form = useForm<z.infer<typeof TCPConfigSchema>>({
        resolver: zodResolver(TCPConfigSchema),
    })

    useEffect(() => {
        setTCPConfig(undefined)
        form.reset({})
    }, [])

    const clientProxyConfigs = useStore($clientProxyConfigs)
    const onSubmit = async (values: z.infer<typeof TCPConfigSchema>) => {
        handleSave()
        setTCPConfig({ type: "tcp", ...values, name: proxyName })
        const newProxiyConfigs = clientProxyConfigs.map(proxyCfg => {
            if (proxyCfg.name === proxyName) {
                return { ...values, type: "tcp", name: proxyName } as TCPProxyConfig
            }
            return proxyCfg
        })
        $clientProxyConfigs.set(newProxiyConfigs)
    }

    const [isSaveDisabled, setSaveDisabled] = useState(false);

    const handleSave = () => {
        setSaveDisabled(true);
        setTimeout(() => {
            setSaveDisabled(false);
        }, 3000);
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                    defaultValue={defaultProxyConfig === undefined ? 1234 : defaultProxyConfig.localPort}
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
                    defaultValue={defaultProxyConfig === undefined ? "127.0.0.1" : defaultProxyConfig.localIP}
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
                    defaultValue={defaultProxyConfig === undefined ? 4321 : (defaultProxyConfig as TCPProxyConfig).remotePort}
                />
                <Button type="submit" disabled={isSaveDisabled} variant={"outline"}>
                    <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? "" : "hidden"}`}></YesIcon>
                    暂存修改</Button>
            </form>
        </Form>
    );
}

export const UDPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName }) => {
    const [_, setUDPConfig] = useState<UDPProxyConfig | undefined>()
    const form = useForm<z.infer<typeof UDPConfigSchema>>({
        resolver: zodResolver(UDPConfigSchema),
    })

    useEffect(() => {
        setUDPConfig(undefined)
        form.reset({})
    }, [])


    const clientProxyConfigs = useStore($clientProxyConfigs)
    const onSubmit = async (values: z.infer<typeof UDPConfigSchema>) => {
        handleSave()
        setUDPConfig({ type: "udp", ...values, name: proxyName })
        const newProxiyConfigs = clientProxyConfigs.map(proxyCfg => {
            if (proxyCfg.name === proxyName) {
                return { ...values, type: "udp", name: proxyName } as UDPProxyConfig
            }
            return proxyCfg
        })
        $clientProxyConfigs.set(newProxiyConfigs)
    }

    const [isSaveDisabled, setSaveDisabled] = useState(false);

    const handleSave = () => {
        setSaveDisabled(true);
        setTimeout(() => {
            setSaveDisabled(false);
        }, 3000);
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                    defaultValue={defaultProxyConfig === undefined ? 1234 : defaultProxyConfig.localPort}
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
                    defaultValue={defaultProxyConfig === undefined ? "127.0.0.1" : defaultProxyConfig.localIP}
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
                    defaultValue={defaultProxyConfig === undefined ? 4321 : (defaultProxyConfig as UDPProxyConfig).remotePort}
                />
                <Button type="submit" disabled={isSaveDisabled} variant={"outline"}>
                    <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? "" : "hidden"}`}></YesIcon>
                    暂存修改</Button>
            </form>
        </Form>
    );
}

export const HTTPProxyForm: React.FC<ProxyFormProps> = ({ serverID, clientID, defaultProxyConfig, proxyName }) => {
    const [_, setHTTPConfig] = useState<HTTPProxyConfig | undefined>()
    const form = useForm<z.infer<typeof HTTPConfigSchema>>({
        resolver: zodResolver(HTTPConfigSchema),
    })

    useEffect(() => {
        setHTTPConfig(undefined)
        form.reset({})
    }, [])

    const clientProxyConfigs = useStore($clientProxyConfigs)
    const onSubmit = async (values: z.infer<typeof HTTPConfigSchema>) => {
        handleSave()
        setHTTPConfig({ ...values, type: "http", name: proxyName })
        const newProxiyConfigs = clientProxyConfigs.map(proxyCfg => {
            if (proxyCfg.name === proxyName) {
                return { ...values, type: "http", name: proxyName } as HTTPProxyConfig
            }
            return proxyCfg
        })
        $clientProxyConfigs.set(newProxiyConfigs)
    }

    const [isSaveDisabled, setSaveDisabled] = useState(false);

    const handleSave = () => {
        setSaveDisabled(true);
        setTimeout(() => {
            setSaveDisabled(false);
        }, 3000);
    };

    return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                    defaultValue={defaultProxyConfig === undefined ? 1234 : defaultProxyConfig.localPort}
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
                    defaultValue={defaultProxyConfig === undefined ? "127.0.0.1" : defaultProxyConfig.localIP}
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
                    defaultValue={defaultProxyConfig === undefined ? "" : (defaultProxyConfig as HTTPProxyConfig).subdomain}
                />
                <Button type="submit" disabled={isSaveDisabled} variant={"outline"}>
                    <YesIcon className={`mr-2 h-4 w-4 ${isSaveDisabled ? "" : "hidden"}`}></YesIcon>
                    暂存修改</Button>
            </form>
        </Form>
    )
}