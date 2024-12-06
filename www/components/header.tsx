import { Button } from './ui/button'
import { useStore } from '@nanostores/react'
import { useRouter } from 'next/router'
import { $platformInfo, $userInfo, $statusOnline, $token } from '@/store/user'
import { getUserInfo } from '@/api/user'
import { useQuery } from '@tanstack/react-query'
import { useEffect, useState } from 'react'
import { getPlatformInfo } from '@/api/platform'
import { useTranslation } from 'react-i18next'
import { LanguageSwitcher } from './language-switcher'

export const Header = ({ title }: { title?: string }) => {
  const router = useRouter()
  const isOnline = useStore($statusOnline)
  const token = useStore($token)
  const [isInitializing, setIsInitializing] = useState(true)
  const currentPath = router.pathname

  useEffect(() => {
    // 设置5秒延迟，等待状态初始化
    const timer = setTimeout(() => {
      setIsInitializing(false)
    }, 5000)

    return () => clearTimeout(timer)
  }, [isOnline, token])

  useEffect(() => {
    // 只有在初始化完成后才进行状态检查和跳转
    if (!isInitializing) {
      console.log('isInitializing', isOnline, token, currentPath)
      // 如果用户未登录且不在登录/注册页面，则跳转到登录页
      const isAuthPage = ['/login', '/register'].includes(currentPath)
      if ((!token || !isOnline) && !isAuthPage) {
        router.push('/login')
      }
    }
  }, [token, isOnline, router, isInitializing, currentPath])

  return (
    <div className="flex w-full justify-between items-center gap-2">
      {title && <p className='font-bold'>{title}</p>}
      {!title && <p></p>}
      <LanguageSwitcher />
    </div>
  )
}

export const RegisterAndLogin = () => {
  const router = useRouter()
  const userInfo = useStore($userInfo)
  const { t } = useTranslation()

  const platformInfo = useQuery({
    queryKey: ['platformInfo'],
    queryFn: getPlatformInfo,
  })

  useEffect(() => {
    $platformInfo.set(platformInfo.data)
  }, [platformInfo])

  const userInfoQuery = useQuery({
    queryKey: ['userInfo'],
    queryFn: getUserInfo,
  })

  useEffect(() => {
    $userInfo.set(userInfoQuery.data?.userInfo)
    $statusOnline.set(!!userInfoQuery.data?.userInfo)
  }, [userInfoQuery])

  return (
    <>
      {!userInfo && (
        <Button variant="ghost" size="sm" onClick={() => router.push('/login')}>
          {t('common.login')}
        </Button>
      )}
      {!userInfo && (
        <Button variant="ghost" size="sm" onClick={() => router.push('/register')}>
          {t('common.register')}
        </Button>
      )}
    </>
  )
}