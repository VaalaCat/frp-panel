import { Inter } from 'next/font/google'
import { FRPCFormCard } from '@/components/frpc_card';
import { Providers } from '@/components/providers';
import { APITest } from '@/components/apitest';
import { Separator } from '@/components/ui/separator';
import { FRPSFormCard } from '@/components/frps_card';
import { RootLayout } from '@/components/layout';

const inter = Inter({ subsets: ['latin'] })

export default function Test() {
    return (
        <main
            className={`flex min-h-screen flex-col items-center justify-between p-2 ${inter.className}`}
        >
            <RootLayout>
                <Providers>
                    <div className='grid grid-cols-1 md:grid-cols-2 gap-8'>
                        <FRPCFormCard></FRPCFormCard>
                        <FRPSFormCard></FRPSFormCard>
                    </div>
                    <Separator className='my-2' />
                    <APITest />
                </Providers>
            </RootLayout>
        </main>
    )
}
