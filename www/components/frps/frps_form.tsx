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
import { Label } from '@radix-ui/react-label'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'

const ServerConfigSchema = z.object({
  bindAddr: ZodIPSchema.default('0.0.0.0').optional(),
  bindPort: ZodPortSchema.default(7000),
  proxyBindAddr: ZodIPSchema.optional(),
  vhostHTTPPort: ZodPortSchema.optional(),
  subDomainHost: ZodStringSchema.optional(),
  publicHost: ZodStringSchema.optional(),
})

export const ServerConfigZodSchema = ServerConfigSchema

export interface FRPSFormProps {
  serverID: string
  server: Server
}

const FRPSForm: React.FC<FRPSFormProps> = ({ serverID, server }) => {
  const { t } = useTranslation()
  const form = useForm<z.infer<typeof ServerConfigZodSchema>>({
    resolver: zodResolver(ServerConfigZodSchema),
  })

  const updateFrps = useMutation({ mutationFn: updateFRPS })

  useEffect(() => {
    form.reset({})
  }, [])

  useEffect(() => {
    form.reset(JSON.parse(server?.config || '{}') as ServerConfig)
  }, [server])

  const onSubmit = async (values: z.infer<typeof ServerConfigZodSchema>) => {
    try {
      const {publicHost, ...rest} = values
      let resp = await updateFrps.mutateAsync({
        serverIp: publicHost,
        serverId: serverID,
        // @ts-ignore
        config: Buffer.from(
          JSON.stringify({
            ...rest,
          } as ServerConfig),
        ),
      })
      toast(resp.status?.code === RespCode.SUCCESS ? t('server.operation.update_success') : t('server.operation.update_failed'),{
        description: resp.status?.message,
      })
    } catch (error) {
      console.error(error)
      toast(t('server.operation.update_title'), {
        description: t('server.operation.update_failed')
      })
    }
  }

  return (
    <div className="flex flex-col w-full pt-2">
      <Label className="text-sm font-medium">{t('server.form.comment_title', { id: serverID })}</Label>
      <p className="text-sm text-muted-foreground">{t('server.form.comment_hint')}</p>
      <p className="text-sm border rounded p-2 my-2">
        {server?.comment == undefined || server?.comment === '' ? t('server.form.comment_empty') : server?.comment}
      </p>
      {serverID && (
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 px-0.5">
            <FormField
              control={form.control}
              name="publicHost"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('server.form.public_host')}</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
              defaultValue={server?.ip}
            />
            <FormField
              control={form.control}
              name="bindPort"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('server.form.bind_port')}</FormLabel>
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
                  <FormLabel>{t('server.form.bind_addr')}</FormLabel>
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
                  <FormLabel>{t('server.form.proxy_bind_addr')}</FormLabel>
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
                  <FormLabel>{t('server.form.vhost_http_port')}</FormLabel>
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
                  <FormLabel>{t('server.form.subdomain_host')}</FormLabel>
                  <FormControl>
                    <Input placeholder="example.com" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">{t('common.submit')}</Button>
          </form>
        </Form>
      )}
    </div>
  )
}

export default FRPSForm
