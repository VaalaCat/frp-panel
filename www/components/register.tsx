import { ZodEmailSchema, ZodStringSchema } from '@/lib/consts'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'
import * as z from 'zod'
import { Form, FormControl, FormField, FormItem, FormMessage, FormLabel } from '@/components/ui/form'
import { Input } from './ui/input'
import { register } from '@/api/auth'
import { Button } from './ui/button'

import { ExclamationTriangleIcon } from '@radix-ui/react-icons'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { useState } from 'react'
import { RespCode } from '@/lib/pb/common'
import { useRouter } from 'next/router'
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner'

export const RegisterSchema = z.object({
  username: ZodStringSchema,
  password: ZodStringSchema,
  email: ZodEmailSchema,
})

export function RegisterComponent() {
  const { t } = useTranslation();
  const form = useForm<z.infer<typeof RegisterSchema>>({
    resolver: zodResolver(RegisterSchema),
  })
  const router = useRouter()

  const [registerAlert, setRegisterAlert] = useState(false)
  const sleep = async (ms: number): Promise<void> => {
    return new Promise((resolve) => setTimeout(resolve, ms))
  }

  const onSubmit = async (values: z.infer<typeof RegisterSchema>) => {
    toast('auth.registering')
    try {
      const res = await register({ ...values })
      if (res.status?.code === RespCode.SUCCESS) {
        toast(t('auth.registerSuccess'))
        setRegisterAlert(false)
        await sleep(3000)
        router.push('/login')
      } else {
        toast(t('auth.registerFailed'), {
          description: res.status?.message
        })
        setRegisterAlert(true)
      }
    } catch (e) {
      toast(t('auth.registerFailed'), {
        description: (e as Error).message
      })
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
                  <Input type="text" placeholder={t('auth.usernamePlaceholder')} {...field} />
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
                  <Input type="email" placeholder={t('auth.emailPlaceholder')} {...field} />
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
                  <Input type="password" placeholder={t('auth.passwordPlaceholder')} {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          {registerAlert && (
            <Alert variant="destructive">
              <ExclamationTriangleIcon className="h-4 w-4" />
              <AlertTitle>{t('auth.error')}</AlertTitle>
              <AlertDescription>{t('auth.registerFailed')}</AlertDescription>
            </Alert>
          )}
          <Button type="submit">
            {t('common.register')}
          </Button>
        </form>
      </Form>
    </div>
  )
}
