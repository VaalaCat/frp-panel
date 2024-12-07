import { ZodStringSchema } from '@/lib/consts'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import * as z from 'zod'
import { Form, FormControl, FormField, FormItem, FormMessage, FormLabel } from '@/components/ui/form'
import { Input } from './ui/input'
import { login } from '@/api/auth'
import { Button } from './ui/button'

import { ExclamationTriangleIcon } from '@radix-ui/react-icons'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { useState } from 'react'
import { RespCode } from '@/lib/pb/common'
import { useRouter } from 'next/router'
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner'

export const LoginSchema = z.object({
  username: ZodStringSchema,
  password: ZodStringSchema,
})

export function LoginComponent() {
  const { t } = useTranslation();
  const form = useForm<z.infer<typeof LoginSchema>>({
    resolver: zodResolver(LoginSchema),
  })
  const router = useRouter()

  const [loginAlert, setLoginAlert] = useState(false)

  const onSubmit = async (values: z.infer<typeof LoginSchema>) => {
    toast(t('auth.loggingIn'))
    try {
      const res = await login({ ...values })
      if (res.status?.code === RespCode.SUCCESS) {
        toast(t('auth.loginSuccess'))
        setTimeout(() => {
          router.push('/')
        }, 3000)
        setLoginAlert(false)
      } else {
        toast(t('auth.loginFailed'), {
          description: res.status?.message
        })
        setLoginAlert(true)
      }
    } catch (e) {
      toast(t('auth.loginFailed'), {
        description: (e as Error).message
      })
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
                <FormLabel>{t('auth.usernamePlaceholder')}</FormLabel>
                <FormControl>
                  <Input type="text" placeholder={t('auth.usernamePlaceholder')} {...field} />
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
                <FormLabel>{t('auth.password')}</FormLabel>
                <FormControl>
                  <Input type="password" placeholder={t('auth.passwordPlaceholder')} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          {loginAlert && (
            <Alert variant="destructive">
              <ExclamationTriangleIcon className="h-4 w-4" />
              <AlertTitle>{t('auth.error')}</AlertTitle>
              <AlertDescription>{t('auth.loginFailed')}</AlertDescription>
            </Alert>
          )}
          <Button className="w-full" type="submit">
            {t('common.login')}
          </Button>
        </form>
      </Form>
    </div>
  )
}
