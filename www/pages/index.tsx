import { Inter } from 'next/font/google'
import { Providers } from '@/components/providers';
import { RootLayout } from '@/components/layout';
import { ClientItem } from '@/components/client_item';

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  return (
    <main
      className={`flex min-h-screen flex-col items-center justify-between p-2 ${inter.className}`}
    >
      <RootLayout>
        <Providers>
          <ClientItem Client={{
            id: "admin.test",
            config: "",
            secret: "admin123",
          }}></ClientItem>
        </Providers>
      </RootLayout>
    </main>
  )
}
