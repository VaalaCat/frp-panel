import { TbBuildingTunnel } from 'react-icons/tb'
import { Button } from './ui/button'
import { useStore } from '@nanostores/react'
import { useRouter } from 'next/router'
import { $platformInfo, $userInfo } from '@/store/user'
import { getUserInfo } from '@/api/user'
import { useQuery } from '@tanstack/react-query'
import { useEffect } from 'react'
import Gravatar from 'react-gravatar'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { LOCAL_STORAGE_TOKEN_KEY } from '@/lib/consts'
import { logout } from '@/api/auth'
import { getPlatformInfo } from '@/api/platform'

export const Header = () => {
  const router = useRouter()
  const userInfo = useStore($userInfo)

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
  }, [userInfoQuery])

  const redirToHome = () => {
    router.push('/')
  }

  return (
    <div className="flex flex-row h-10 items-center px-4 border-b">
      <TbBuildingTunnel />
      <p className="ml-2 font-mono" onClick={redirToHome}>
        frp-panel
      </p>
      {!userInfo && (
        <Button variant={'ghost'} className="ml-auto" size={'sm'} onClick={() => router.push('/login')}>
          登录
        </Button>
      )}
      {!userInfo && (
        <Button variant={'ghost'} className="ml-2" size={'sm'} onClick={() => router.push('/register')}>
          注册
        </Button>
      )}
      {userInfo && (
        <Button
          variant={'ghost'}
          className="ml-auto"
          size={'sm'}
          onClick={async () => {
            $userInfo.set(undefined)
            localStorage.removeItem(LOCAL_STORAGE_TOKEN_KEY)
            await logout()
            window.location.reload()
          }}
        >
          退出
        </Button>
      )}
      {userInfo && (
        <Avatar className="ml-2 w-7 h-7">
          <AvatarImage alt={'@' + userInfo.userName} asChild>
            <Gravatar email={userInfo.email} />
          </AvatarImage>
          <AvatarFallback>{userInfo.userName}</AvatarFallback>
        </Avatar>
      )}
    </div>
  )
}
