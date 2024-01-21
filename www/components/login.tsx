import { ZodStringSchema } from '@/lib/consts'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import * as z from 'zod'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { Input } from './ui/input'
import { login } from '@/api/auth'
import { Button } from './ui/button'

import { ExclamationTriangleIcon } from '@radix-ui/react-icons'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { useState } from 'react'
import { useToast } from './ui/use-toast'
import { RespCode } from '@/lib/pb/common'
import { useRouter } from 'next/router'

export const LoginSchema = z.object({
  username: ZodStringSchema,
  password: ZodStringSchema,
})

export const LoginComponent = () => {
  const form = useForm<z.infer<typeof LoginSchema>>({
    resolver: zodResolver(LoginSchema),
  })
  const { toast } = useToast()
  const router = useRouter()

  const [loginAlert, setLoginAlert] = useState(false)

  const onSubmit = async (values: z.infer<typeof LoginSchema>) => {
    toast({ title: '登录中，请稍候' })
    try {
      const res = await login({ ...values })
      if (res.status?.code === RespCode.SUCCESS) {
        toast({ title: '登录成功，正在跳转到首页' })
        setTimeout(() => {
          router.push('/')
        }, 3000)
        setLoginAlert(false)
      } else {
        toast({ title: '登录失败' })
        setLoginAlert(true)
      }
    } catch (e) {
      toast({ title: '登录失败' })
      console.log('login error', e)
      setLoginAlert(true)
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
          {loginAlert && (
            <Alert variant="destructive">
              <ExclamationTriangleIcon className="h-4 w-4" />
              <AlertTitle>错误</AlertTitle>
              <AlertDescription>登录失败，请重试</AlertDescription>
            </Alert>
          )}
          <Button type="submit">登录</Button>
        </form>
      </Form>
    </div>
  )
}
