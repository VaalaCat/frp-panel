import { Providers } from './providers'
import { Toaster } from './ui/toaster'
import { Inter } from 'next/font/google'

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
  return (
    <main className={`${inter.className}`}>
      <div>
        <Providers>{header}</Providers>
      </div>
      <div className="flex">
        {sidebar}
        <div className="my-2 ml-0 mr-2 max-w-full w-full">{children}</div>
      </div>
      <Toaster />
    </main>
  )
}
