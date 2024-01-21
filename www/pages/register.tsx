import { Inter } from 'next/font/google'
import { Providers } from '@/components/providers'
import { TbBuildingTunnel } from 'react-icons/tb'
import { RegisterComponent } from '@/components/register'
import { useRouter } from 'next/router'
import { Toaster } from '@/components/ui/toaster'

const inter = Inter({ subsets: ['latin'] })

export default function Login() {
  const router = useRouter()
  return (
    <main className={`${inter.className}`}>
      <Providers>
        <div
          className="absolute text-lg font-medium left-1/2 transform -translate-x-1/2 mt-3 lg:hidden"
          onClick={() => router.push('/')}
        >
          <div className="flex rounded hover:bg-slate-100 p-2">
            <TbBuildingTunnel className="mr-2 h-8 w-8 pb-1" />
            FRP Panel
          </div>
        </div>
        <div className="container h-screen flex-col items-center justify-center grid lg:max-w-none lg:grid-cols-2 lg:px-0">
          <div className="relative hidden h-full flex-col bg-muted p-10 text-white lg:flex dark:border-r">
            <div className="absolute inset-0 bg-zinc-900"></div>
            <div className="relative flex items-center text-lg font-medium" onClick={() => router.push('/')}>
              <div className="flex rounded hover:bg-zinc-800 p-2">
                <TbBuildingTunnel className="mr-2 h-8 w-8 pb-1" />
                FRP Panel
              </div>
            </div>
            <div className="relative z-20 mt-auto">
              <blockquote className="space-y-2">
                <p className="text-lg">
                  A multi node frp webui and for <a href="https://github.com/fatedier/frp">[FRP]</a> server and client
                  management, which makes this project a [Cloudflare Tunnel] or [Tailscale Funnel] open source
                  alternative
                </p>
                <footer className="text-sm">
                  navigate to: <a href="https://github.com/VaalaCat/frp-panel">VaalaCat/frp-panel</a>
                </footer>
              </blockquote>
            </div>
          </div>
          <div className="lg:p-8 justify-center w-[300px]">
            <div className="flex flex-col justify-center space-y-6 w-[300px]">
              <div className="flex flex-col space-y-2 text-center">
                <h1 className="text-2xl font-semibold tracking-tight">注册</h1>
                <p className="text-sm text-muted-foreground">输入您的账号信息</p>
              </div>
              <div className="w-full justify-center">
                <div className="w-[300px]">
                  <RegisterComponent />
                </div>
              </div>
            </div>
          </div>
        </div>
        <Toaster />
      </Providers>
    </main>
  )
}
