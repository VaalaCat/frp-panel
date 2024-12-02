import * as React from "react"
import {
  MessagesSquare,
  SquareTerminal,
  ServerCogIcon,
  ServerIcon,
  MonitorSmartphoneIcon,
  MonitorCogIcon,
  ChartNetworkIcon,
  icons,
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
import { $userInfo } from "@/store/user"
import { useStore } from "@nanostores/react"
import { TbBuildingTunnel } from "react-icons/tb"
import { RegisterAndLogin } from "./header"
import { useRouter } from "next/navigation"

const data = {
  teams: [
    {
      name: "Frp-Panel",
      logo: TbBuildingTunnel,
      plan: "Community Edition",
      url: "/",
    },
  ],
  navMain: [
    {
      title: "客户端",
      url: "/clients",
      icon: MonitorSmartphoneIcon,
      isActive: true,
    },
    {
      title: "服务端",
      url: "/servers",
      icon: ServerIcon,
    },
    {
      title: "编辑隧道",
      url: "/clientedit",
      icon: MonitorCogIcon,
    },
    {
      title: "编辑服务端",
      url: "/serveredit",
      icon: ServerCogIcon,
    },
    {
      title: "流量统计",
      url: "/clientstats",
      icon: ChartNetworkIcon,
    },
    {
      title: "实时日志",
      url: "/streamlog",
      icon: Scroll,
    }
  ]
}

export interface AppSidebarProps extends React.ComponentProps<typeof Sidebar> {
  chrildren?: React.ReactNode
  footer?: React.ReactNode
}

export function AppSidebar({ ...props }: AppSidebarProps) {
  const router = useRouter()
  const userInfo = useStore($userInfo)
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
              Frp-Panel
            </span>
            <span className="truncate text-xs font-mono">frp隧道面板</span>
          </div>
        </SidebarMenuButton>
        <NavMain items={data.navMain} />
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
