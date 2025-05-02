'use client'

import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { getUserInfo, updateUserInfo } from '@/api/user'
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Avatar } from '@/components/ui/avatar'
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card'
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogFooter as DialogFooterUI,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { UserAvatar } from '@/components/base/avatar'
import { toast } from 'sonner'
import { User } from '@/lib/pb/common'
import { useTranslation } from 'react-i18next'

const userSchema = z
  .object({
    UserID: z.string().optional(),
    TenantID: z.string().optional(),
    UserName: z.string().min(2, 'Name too short'),
    Email: z.string().email('Invalid email address'),
    Status: z.string().optional(),
    Role: z.string().optional(),
    NewPassword: z.string().optional(),
    ConfirmPassword: z.string().optional(),
  })
  .superRefine((data, ctx) => {
    const np = data.NewPassword
    const cp = data.ConfirmPassword
    if (np) {
      if (np.length < 6) {
        ctx.addIssue({
          path: ['NewPassword'],
          message: 'Password must be at least 6 characters',
          code: z.ZodIssueCode.custom,
        })
      }
      if (np !== cp) {
        ctx.addIssue({
          path: ['ConfirmPassword'],
          message: 'Passwords do not match',
          code: z.ZodIssueCode.custom,
        })
      }
    }
  })

type UserFormValues = z.infer<typeof userSchema>

export function UserProfileForm() {
  const { t } = useTranslation()
  const [loading, setLoading] = useState(true)
  const [initial, setInitial] = useState<UserFormValues | null>(null)
  const [open, setOpen] = useState(false)

  const form = useForm<UserFormValues>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      UserID: '',
      TenantID: '',
      UserName: '',
      Email: '',
      Role: '',
      NewPassword: '',
      ConfirmPassword: '',
    },
  })

  // Fetch on mount
  useEffect(() => {
    getUserInfo({})
      .then((res) => {
        const u = res.userInfo! as User
        form.reset({
          UserID: u.userID?.toString(),
          TenantID: u.tenantID?.toString(),
          UserName: u.userName || '',
          Email: u.email || '',
          Role: u.role || '',
          NewPassword: '',
          ConfirmPassword: '',
        })
        setInitial(form.getValues())
      })
      .finally(() => setLoading(false))
  }, [])

  const onSubmit = async (values: UserFormValues) => {
    try {
      const payload: any = {
        userID: values.UserID ? BigInt(values.UserID) : undefined,
        tenantID: values.TenantID ? BigInt(values.TenantID) : undefined,
        userName: values.UserName,
        email: values.Email,
      }
      if (values.NewPassword) {
        payload.rawPassword = values.NewPassword
      }
      await updateUserInfo({ userInfo: payload })
      toast.success(t('userInfo.profileUpdated'))
      form.reset({
        ...values,
        NewPassword: '',
        ConfirmPassword: '',
      })
      setInitial(form.getValues())
    } catch {
      toast(t('userInfo.updateFailed'))
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <p className="text-gray-500">{t('userInfo.loading')}</p>
      </div>
    )
  }

  return (
    <Card className="max-w-lg mx-auto">
      <CardHeader className="flex flex-row items-center space-x-4 border-b">
        <Avatar className="h-12 w-12 rounded-lg">
          <UserAvatar
            className="h-12 w-12"
            userInfo={{ userName: form.getValues().UserName, email: form.getValues().Email } as User}
          />
        </Avatar>
        <div>
          <CardTitle>{t('userInfo.yourProfile')}</CardTitle>
          <CardDescription>{t('userInfo.manageAccountDetails')}</CardDescription>
        </div>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 pt-2">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* UserID (read-only) */}
              <FormField
                control={form.control}
                name="UserID"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userInfo.userID')}</FormLabel>
                    <FormControl>
                      <Input {...field} disabled className="bg-gray-100 dark:bg-gray-700" />
                    </FormControl>
                  </FormItem>
                )}
              />

              {/* TenantID (read-only) */}
              <FormField
                control={form.control}
                name="TenantID"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userInfo.tenantID')}</FormLabel>
                    <FormControl>
                      <Input {...field} disabled className="bg-gray-100 dark:bg-gray-700" />
                    </FormControl>
                  </FormItem>
                )}
              />

              {/* UserName */}
              <FormField
                control={form.control}
                name="UserName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userInfo.name')}</FormLabel>
                    <FormControl>
                      <Input placeholder={t('userInfo.placeholderName')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Email */}
              <FormField
                control={form.control}
                name="Email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userInfo.email')}</FormLabel>
                    <FormControl>
                      <Input type="email" placeholder={t('userInfo.placeholderEmail')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Role (read-only badge style) */}
              <FormField
                control={form.control}
                name="Role"
                render={({ field }) => (
                  <FormItem className="md:col-span-2">
                    <FormLabel>{t('userInfo.role')}</FormLabel>
                    <FormDescription>{t('userInfo.roleDescription')}</FormDescription>
                    <FormControl>
                      <Input
                        {...field}
                        disabled
                        className="bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300"
                      />
                    </FormControl>
                  </FormItem>
                )}
              />

              {/* Change Password Section */}
              <div className="md:col-span-2">
                <h3 className="text-lg font-medium text-gray-700 dark:text-gray-300">{t('userInfo.changePassword')}</h3>
              </div>

              {/* New Password */}
              <FormField
                control={form.control}
                name="NewPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userInfo.newPassword')}</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder={t('userInfo.placeholderNewPassword')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Confirm Password */}
              <FormField
                control={form.control}
                name="ConfirmPassword"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{t('userInfo.confirmPassword')}</FormLabel>
                    <FormControl>
                      <Input type="password" placeholder={t('userInfo.placeholderConfirmNewPassword')} {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          </form>
        </Form>
      </CardContent>
      <CardFooter className="flex justify-end">
        <Dialog open={open} onOpenChange={setOpen}>
          <DialogTrigger asChild>
            <Button
              disabled={form.formState.isSubmitting || JSON.stringify(form.getValues()) === JSON.stringify(initial)}
              onClick={() => setOpen(true)}
            >
              {t('userInfo.saveChanges')}
            </Button>
          </DialogTrigger>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>{t('userInfo.confirmSaveTitle')}</DialogTitle>
              <DialogDescription>{t('userInfo.confirmSaveDescription')}</DialogDescription>
            </DialogHeader>
            <DialogFooterUI>
              <Button variant={'destructive'} onClick={() => setOpen(false)}>
                {t('userInfo.cancel')}
              </Button>
              <Button
                variant={'secondary'}
                onClick={() => {
                  form.handleSubmit(onSubmit)()
                  setOpen(false)
                }}
              >
                {t('userInfo.confirm')}
              </Button>
            </DialogFooterUI>
          </DialogContent>
        </Dialog>
      </CardFooter>
    </Card>
  )
}
