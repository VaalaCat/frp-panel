import { Inter } from 'next/font/google'
import { Providers } from '@/components/providers'
import { TbBuildingTunnel } from 'react-icons/tb'
import { LoginComponent } from '@/components/login'
import { useRouter } from 'next/router'
import { Toaster } from '@/components/ui/toaster'
import { useTranslation } from 'react-i18next'
import { LanguageSwitcher } from '@/components/language-switcher'
import Link from 'next/link'

const inter = Inter({ subsets: ['latin'] })

export default function LoginPage() {
  const router = useRouter()
  const { t } = useTranslation();
  
  return (
    <main className={`${inter.className} min-h-screen`}>
      <Providers>
        {/* Fixed Language Switcher */}
        <div className="fixed top-4 right-4 z-50">
          <LanguageSwitcher />
        </div>

        {/* Mobile Header */}
        <div className="fixed w-full flex items-center px-4 py-3 lg:hidden bg-white/80 backdrop-blur-sm border-b z-40">
          <div
            className="text-lg font-medium flex items-center"
            onClick={() => router.push('/')}
          >
            <div className="flex items-center rounded hover:bg-slate-100 p-2">
              <TbBuildingTunnel className="mr-2 h-6 w-6" />
              {t('app.title')}
            </div>
          </div>
        </div>

        <div className="container min-h-screen flex-col items-center justify-center grid lg:max-w-none lg:grid-cols-2 lg:px-0">
          {/* Left Panel */}
          <div className="relative hidden h-full flex-col bg-muted p-10 text-zinc-500 lg:flex dark:border-r">
            <div className="absolute inset-0 bg-zinc-900"></div>
            <div className="relative z-20">
              <div className="flex items-center text-lg font-medium" onClick={() => router.push('/')}>
                <div className="flex items-center rounded hover:bg-zinc-800 p-2 text-white">
                  <TbBuildingTunnel className="mr-2 h-8 w-8" />
                  {t('app.title')}
                </div>
              </div>
            </div>
            <div className="relative z-20 mt-auto">
              <blockquote className="space-y-2">
                <p className="text-lg leading-relaxed">
                  {t('app.description')}
                </p>
                <footer className="text-sm mt-4 opacity-80">
                  {t('app.github.navigate')} 
                  <a 
                    href="https://github.com/VaalaCat/frp-panel"
                    className="hover:text-white hover:underline ml-1"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    {t('app.github.repo')}
                  </a>
                </footer>
              </blockquote>
            </div>
          </div>

          {/* Right Panel - Login Form */}
          <div className="lg:p-8 flex items-center justify-center pt-20 lg:pt-0">
            <div className="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
              <div className="flex flex-col space-y-2 text-center">
                <h1 className="text-2xl font-semibold tracking-tight">
                  {t('auth.loginTitle')}
                </h1>
                <p className="text-sm text-muted-foreground">
                  {t('auth.inputCredentials')}
                </p>
              </div>
              <LoginComponent />
              <p className="px-8 text-center text-sm text-muted-foreground">
                {t('auth.noAccount')}{' '}
                <Link
                  className="underline underline-offset-4 hover:text-primary"
                  href="/register"
                >
                  {t('auth.register')}
                </Link>
              </p>
            </div>
          </div>
        </div>
        <Toaster />
      </Providers>
    </main>
  )
}
