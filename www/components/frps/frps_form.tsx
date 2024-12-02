import { ServerConfig } from '@/types/server'
import { useEffect, useState } from 'react'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import * as z from 'zod'
import { Button } from '@/components/ui/button'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { ZodIPSchema, ZodPortSchema, ZodStringSchema } from '@/lib/consts'
import { RespCode, Server } from '@/lib/pb/common'
import { updateFRPS } from '@/api/frp'
import { useMutation } from '@tanstack/react-query'
import { useToast } from '@/components/ui/use-toast'
import { Label } from '@radix-ui/react-label'

const ServerConfigSchema = z.object({
  bindAddr: ZodIPSchema.default('0.0.0.0').optional(),
  bindPort: ZodPortSchema.default(7000),
  proxyBindAddr: ZodIPSchema.optional(),
  vhostHTTPPort: ZodPortSchema.optional(),
  subDomainHost: ZodStringSchema.optional(),
})

export const ServerConfigZodSchema = ServerConfigSchema

export interface FRPSFormProps {
  serverID: string
  server: Server
}

const FRPSForm: React.FC<FRPSFormProps> = ({ serverID, server }) => {
  const [_, setFrpsConfig] = useState<ServerConfig | undefined>()
  const form = useForm<z.infer<typeof ServerConfigZodSchema>>({
    resolver: zodResolver(ServerConfigZodSchema),
  })
  const { toast } = useToast()

  const updateFrps = useMutation({ mutationFn: updateFRPS })

  useEffect(() => {
    setFrpsConfig(undefined)
    form.reset({})
  }, [])

  useEffect(() => {
    form.reset(JSON.parse(server?.config || '{}') as ServerConfig)
  }, [server])

  const onSubmit = async (values: z.infer<typeof ServerConfigZodSchema>) => {
    setFrpsConfig({ ...values })
    try {
      let resp = await updateFrps.mutateAsync({
        serverId: serverID,
        // @ts-ignore
        config: Buffer.from(
          JSON.stringify({
            ...values,
          } as ServerConfig),
        ),
      })
      toast({
        title: resp.status?.code === RespCode.SUCCESS ? '创建成功' : '创建失败',
        description: resp.status?.message,
      })
    } catch (error) {
      console.error(error)
      toast({ title: '创建服务端状态', description: '创建失败' })
    }
  }
  return (
    <div className="flex flex-col w-full pt-2">
      <Label className="text-sm font-medium">节点 {serverID} 的备注</Label>
      <p className="text-sm text-muted-foreground">可以到高级模式修改备注哦！</p>
      <p className="text-sm border rounded p-2 my-2">
        {server?.comment == undefined || server?.comment === '' ? '空空如也' : server?.comment}
      </p>
      {serverID && (
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="bindPort"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>FRPs 监听端口</FormLabel>
                  <FormControl>
                    <Input type="number" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
              defaultValue={7000}
            />
            <FormField
              control={form.control}
              name="bindAddr"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>FRPs 监听地址</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
              defaultValue="0.0.0.0"
            />
            <FormField
              control={form.control}
              name="proxyBindAddr"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>代理监听地址</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="vhostHTTPPort"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>HTTP 监听端口</FormLabel>
                  <FormControl>
                    <Input type="number" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="subDomainHost"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>域名后缀</FormLabel>
                  <FormControl>
                    <Input placeholder="example.com" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">提交</Button>
          </form>
        </Form>
      )}
    </div>
  )
}

export default FRPSForm
