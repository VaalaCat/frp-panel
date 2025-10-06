import { Providers } from '@/components/providers'
import { RootLayout } from '@/components/layout'
import { Header } from '@/components/header'
import CanvasPanel from '@/components/canvas/CanvasPanel'

export default function CanvasPage() {
  return (
    <Providers>
      <RootLayout mainHeader={<Header />}>
        <div className="w-full">
          <CanvasPanel />
        </div>
      </RootLayout>
    </Providers>
  )
}
