import * as React from "react"
import {
  SquareTerminal,
  ServerCogIcon,
  ServerIcon,
  MonitorSmartphoneIcon,
  MonitorCogIcon,
  ChartNetworkIcon,
  Scroll,
} from "lucide-react"

import { NavMain } from "@/components/nav-main"
import { NavUser } from "@/components/nav-user"
import { TeamSwitcher } from "@/components/team-switcher"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenuButton,
  SidebarRail,
} from "@/components/ui/sidebar"
import { $platformInfo, $userInfo } from "@/store/user"
import { useStore } from "@nanostores/react"
import { TbBuildingTunnel } from "react-icons/tb"
import { RegisterAndLogin } from "./header"
import { useRouter } from "next/navigation"
import { useQuery } from "@tanstack/react-query"
import { getPlatformInfo } from "@/api/platform"
import { teams, getNavItems } from '@/config/nav'
import { useTranslation } from 'react-i18next'

export interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  children?: React.ReactNode
  footer?: React.ReactNode
}

export function AppSidebar({ ...props }: AppSidebarProps) {
  const router = useRouter()
  const { t } = useTranslation()
  const userInfo = useStore($userInfo)
  const { data: platformInfo } = useQuery({
    queryKey: ['platformInfo'],
    queryFn: getPlatformInfo,
  })

  React.useEffect(() => {
    $platformInfo.set(platformInfo)
  }, [platformInfo])

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenuButton
          size="lg"
          className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          onClick={() => router.push("/")}
        >
          <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
            <TbBuildingTunnel className="size-4" />
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-semibold font-mono">
              {t('app.title')}
            </span>
            <span className="truncate text-xs font-mono">{t('app.subtitle')}</span>
          </div>
        </SidebarMenuButton>
        <NavMain items={getNavItems(t)} />
      </SidebarHeader>
      <SidebarContent>
        {props.children}
      </SidebarContent>
      <SidebarFooter>
        {props.footer}
        <div className="flex w-full flex-row group-data-[collapsible=icon]:flex-col-reverse gap-2 justify-between">
          {userInfo && <NavUser user={userInfo} />}
          {!userInfo && <RegisterAndLogin />}
        </div>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
