import {
  SquareTerminal,
  ServerCogIcon,
  ServerIcon,
  MonitorSmartphoneIcon,
  MonitorCogIcon,
  ChartNetworkIcon,
  Scroll,
  Cable,
} from "lucide-react"
import { TbBuildingTunnel } from "react-icons/tb"

export const teams = [
  {
    name: "Frp-Panel",
    logo: TbBuildingTunnel,
    plan: "Community Edition",
    url: "/",
  },
]

export const getNavItems = (t: any) => [
  {
    title: t('nav.clients'),
    url: "/clients",
    icon: MonitorSmartphoneIcon,
    isActive: true,
  },
  {
    title: t('nav.servers'),
    url: "/servers",
    icon: ServerIcon,
  },
  {
    title: t('nav.editTunnel'),
    url: "/proxies",
    icon: Cable,
  },
  {
    title: t('nav.editClient'),
    url: "/clientedit",
    icon: MonitorCogIcon,
  },
  {
    title: t('nav.editServer'),
    url: "/serveredit",
    icon: ServerCogIcon,
  },
  {
    title: t('nav.trafficStats'),
    url: "/clientstats",
    icon: ChartNetworkIcon,
  },
  {
    title: t('nav.realTimeLog'),
    url: "/streamlog",
    icon: Scroll,
  },
  {
    title: t('nav.console'),
    url: "/console",
    icon: SquareTerminal,
  },
]
