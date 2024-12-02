import React from 'react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { TooltipProvider } from '@/components/ui/tooltip'
import {
  SidebarProvider,
} from "@/components/ui/sidebar"

const queryClient = new QueryClient()

export const Providers = ({ children }: { children: React.ReactNode }) => {
  return (
    <TooltipProvider>
      <SidebarProvider>
        <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
      </SidebarProvider>
    </TooltipProvider>
  )
}
