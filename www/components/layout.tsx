import { Toaster } from "@/components/ui/sonner"
import { AppSidebar } from "@/components/app-sidebar"
import { Separator } from "@/components/ui/separator"
import {
  SidebarInset,
  SidebarTrigger,
  useSidebar,
} from "@/components/ui/sidebar"
import { cn } from "@/lib/utils";

export function RootLayout({
  children,
  sidebarItems,
  sidebarFooter,
  mainHeader,
}: {
  children: React.ReactNode;
  sidebarItems?: React.ReactNode;
  sidebarFooter?: React.ReactNode;
  mainHeader?: React.ReactNode
}) {
  const { open, isMobile } = useSidebar()
  return (
    <>
      <AppSidebar footer={sidebarFooter} >{sidebarItems} </AppSidebar>
      <SidebarInset>
        <div className={cn("relative flex flex-1 flex-col overflow-hidden transition-all",
          isMobile && "w-[100dvw]",
          !isMobile && open && "w-[calc(100dvw-var(--sidebar-width))]",
          !isMobile && !open && "w-[calc(100dvw-var(--sidebar-width-icon))]"
        )}>
          <header className="flex flex-row h-12 items-center gap-2 w-full">
            <div className="flex flex-row items-center gap-2 px-4 w-full">
              <SidebarTrigger className="-ml-1" />
              <Separator orientation="vertical" className="mr-2 h-4" />
              {mainHeader}
            </div>
          </header>
          <div id="main-content"
            className="h-[calc(100dvh_-_48px)] overflow-auto w-full pb-4 px-4 pt-2">
            {children}
          </div>
        </div>
        <Toaster />
      </SidebarInset>
    </>
  )
}