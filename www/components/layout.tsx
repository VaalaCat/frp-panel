import { useStore } from '@nanostores/react'
import { Providers } from './providers'
import { Toaster } from './ui/toaster'
import { Inter } from 'next/font/google'
import { $language } from '@/lib/i18n'

const inter = Inter({ subsets: ['latin'] })

export const RootLayout = ({
  children,
  header,
  sidebar,
}: {
  children: React.ReactNode
  header: React.ReactNode
  sidebar?: React.ReactNode
}) => {
  const language = useStore($language)
  return (
    <main key={language} className={`${inter.className}`}>
      <div>
        <Providers>{header}</Providers>
      </div>
      <div className="flex">
        {sidebar}
        <div className="my-2 ml-0 mr-2 max-w-[calc(100vw-100px)] w-full">{children}</div>
      </div>
      <Toaster />
    </main>
  )
}
