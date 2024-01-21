import { ZodEmailSchema, ZodStringSchema } from '@/lib/consts'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import * as z from 'zod'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from './ui/input'
import { register } from '@/api/auth'
import { Button } from './ui/button'

import { ExclamationTriangleIcon } from '@radix-ui/react-icons'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { useState } from 'react'
import { useToast } from './ui/use-toast'
import { RespCode } from '@/lib/pb/common'
import { useRouter } from 'next/router'
import { Toast } from './ui/toast'

export const RegisterSchema = z.object({
  username: ZodStringSchema,
  password: ZodStringSchema,
  email: ZodEmailSchema,
})

export const RegisterComponent = () => {
  const form = useForm<z.infer<typeof RegisterSchema>>({
    resolver: zodResolver(RegisterSchema),
  })
  const { toast } = useToast()
  const router = useRouter()

  const [registerAlert, setRegisterAlert] = useState(false)
  const sleep = async (ms: number): Promise<void> => {
    return new Promise((resolve) => setTimeout(resolve, ms))
  }

  const onSubmit = async (values: z.infer<typeof RegisterSchema>) => {
    toast({ title: '注册中，请稍候' })
    try {
      const res = await register({ ...values })
      if (res.status?.code === RespCode.SUCCESS) {
        toast({ title: '注册成功，正在跳转到登录' })
        setRegisterAlert(false)
        await sleep(3000)
        router.push('/login')
      } else {
        toast({ title: '注册失败' })
        setRegisterAlert(true)
      }
    } catch (e) {
      toast({ title: '注册失败' })
      console.log('register error', e)
      setRegisterAlert(true)
    }
  }

  return (
    <div className="w-full flex flex-col gap-6">
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="flex flex-col gap-4">
          <FormField
            control={form.control}
            name="username"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input type="text" placeholder="用户名" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="email"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input type="email" placeholder="邮箱地址" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="password"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input type="password" placeholder="密码" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          {registerAlert && (
            <Alert variant="destructive">
              <ExclamationTriangleIcon className="h-4 w-4" />
              <AlertTitle>错误</AlertTitle>
              <AlertDescription>注册失败，请重试</AlertDescription>
            </Alert>
          )}
          <Button type="submit">注册</Button>
        </form>
      </Form>
    </div>
  )
}
