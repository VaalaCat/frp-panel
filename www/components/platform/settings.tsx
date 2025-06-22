import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import * as z from 'zod'
import { useStore } from '@nanostores/react'
import { $platformInfo } from '@/store/user'

import { $frontendPreference, FrontendPreference } from '@/store/user'

import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage, FormDescription } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Button } from '@/components/ui/button'
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter as DialogFooterUI,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from '@/components/ui/dialog'

// 表单校验 Schema
const platformSchema = z.object({
	useServerGithubProxyUrl: z.boolean().default(false),
	// githubProxyUrl 可选；若为空字符串则忽略，否则需为合法 URL
	githubProxyUrl: z.union([z.string().trim().url('Invalid URL'), z.literal('')]).optional(),
	clientApiUrl: z.union([z.string().trim().url('Invalid URL'), z.literal('')]).optional(),
	clientRpcUrl: z.union([z.string().trim().url('Invalid URL'), z.literal('')]).optional(),
})

type PlatformFormValues = z.infer<typeof platformSchema>

export function PlatformSettingsForm() {
	const { t } = useTranslation()
	const [loading, setLoading] = useState(true)
	const [initial, setInitial] = useState<PlatformFormValues | null>(null)
	const [open, setOpen] = useState(false)

	const platformInfo = useStore($platformInfo)

	const form = useForm<PlatformFormValues>({
		resolver: zodResolver(platformSchema),
		defaultValues: {
			useServerGithubProxyUrl: false,
			githubProxyUrl: '',
			clientApiUrl: '',
			clientRpcUrl: '',
		},
	})

	// 组件挂载时读取持久化设置
	useEffect(() => {
		const pref = ($frontendPreference.get() ?? {}) as FrontendPreference
		form.reset({
			useServerGithubProxyUrl: pref.useServerGithubProxyUrl ?? false,
			githubProxyUrl: pref.githubProxyUrl ?? '',
			clientApiUrl: pref.clientApiUrl ?? '',
			clientRpcUrl: pref.clientRpcUrl ?? '',
		})
		setInitial(form.getValues())
		setLoading(false)
	}, [])

	const onSubmit = (values: PlatformFormValues) => {
		const pref: FrontendPreference = {
			useServerGithubProxyUrl: values.useServerGithubProxyUrl,
			githubProxyUrl: values.githubProxyUrl?.trim() || undefined,
			clientApiUrl: values.clientApiUrl?.trim() || undefined,
			clientRpcUrl: values.clientRpcUrl?.trim() || undefined,
		}
		$frontendPreference.set(pref)
		toast.success(t('已更新平台设置'))
		// 重置 initial 状态 & 清空 dirty
		form.reset(values)
		setInitial(values)
	}

	if (loading) {
		return (
			<div className="flex items-center justify-center h-64">
				<p className="text-gray-500">{t('正在加载平台设置')}</p>
			</div>
		)
	}

	return (
		<Card className="max-w-lg mx-auto">
			<CardHeader className="border-b">
				<CardTitle>{t('平台设置')}</CardTitle>
				<CardDescription>{t('修改前端平台设置')}</CardDescription>
				<p className="text-xs text-muted-foreground mt-1 italic">{t('此配置仅保存在本地')}</p>
			</CardHeader>
			<CardContent>
				<Form {...form}>
					<form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 pt-2">
						{/* 使用服务器代理开关 */}
						<FormField
							control={form.control}
							name="useServerGithubProxyUrl"
							render={({ field }) => (
								<FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
									<div className="space-y-0.5">
										<FormLabel>{t('使用服务器 Github 代理')}</FormLabel>
										<FormDescription>{t('若开启，则使用后台配置的代理地址下载')}</FormDescription>
									</div>
									<FormControl>
										<Switch checked={field.value} onCheckedChange={field.onChange} />
									</FormControl>
								</FormItem>
							)}
						/>

						{/* 自定义 GitHub Proxy URL */}
						<FormField
							control={form.control}
							name="githubProxyUrl"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t('自定义 Github 代理地址')}</FormLabel>
									<FormControl>
										<Input placeholder={platformInfo?.githubProxyUrl || t('例如 https://ghproxy.com/')} {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						{/* 自定义 API URL */}
						<FormField
							control={form.control}
							name="clientApiUrl"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t('自定义 API URL')}</FormLabel>
									<FormControl>
										<Input placeholder={platformInfo?.clientApiUrl || t('例如 https://api.example.com/')} {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						{/* 自定义 RPC URL */}
						<FormField
							control={form.control}
							name="clientRpcUrl"
							render={({ field }) => (
								<FormItem>
									<FormLabel>{t('自定义 RPC URL')}</FormLabel>
									<FormControl>
										<Input placeholder={platformInfo?.clientRpcUrl || t('例如 https://rpc.example.com/')} {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>
					</form>
				</Form>
			</CardContent>
			<CardFooter className="flex justify-end">
				<Dialog open={open} onOpenChange={setOpen}>
					<DialogTrigger asChild>
						<Button
							disabled={
								form.formState.isSubmitting ||
								JSON.stringify(form.getValues()) === JSON.stringify(initial)
							}
							onClick={() => setOpen(true)}
						>
							{t('保存更改')}
						</Button>
					</DialogTrigger>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>{t('确认保存')}</DialogTitle>
							<DialogDescription>{t('确定保存修改？')}</DialogDescription>
						</DialogHeader>
						<DialogFooterUI>
							<Button variant={'destructive'} onClick={() => setOpen(false)}>
								{t('取消')}
							</Button>
							<Button
								variant={'secondary'}
								onClick={() => {
									form.handleSubmit(onSubmit)()
									setOpen(false)
								}}
							>
								{t('确认')}
							</Button>
						</DialogFooterUI>
					</DialogContent>
				</Dialog>
			</CardFooter>
		</Card>
	)
}
