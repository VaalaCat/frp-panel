import { useEffect } from 'react'
import { Button } from './ui/button'
import { useRouter } from 'next/router'

export interface SideBarItem {
  id: string
  label: string
  eventHandler: () => void
}

export interface SideBarProps {
  items?: SideBarItem[]
}

export const SideBar: React.FC<SideBarProps> = ({ items }) => {
  const router = useRouter()
  const defaultItems = [
    {
      id: 'clients',
      label: '客户端',
      eventHandler: () => {
        router.push('/clients')
      },
    },
    {
      id: 'servers',
      label: '服务端',
      eventHandler: () => {
        router.push('/servers')
      },
    },
    {
      id: 'clientedit',
      label: '编辑隧道',
      eventHandler: () => {
        router.push('/clientedit')
      },
    },
    {
      id: 'serveredit',
      label: '编辑端点',
      eventHandler: () => {
        router.push('/serveredit')
      },
    },
  ]
  return (
    <div className="w-48 h-full grid grid-cols-1 mt-1 min-w-24">
      {items &&
        items.map((item) => (
          <Button
            className={`mx-2 my-1 justify-start ${
              router.pathname.includes(item.id)
                ? 'bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground'
                : ''
            }`}
            variant={'ghost'}
            size={'sm'}
            key={item.id}
            onClick={item.eventHandler}
          >
            {item.label}
          </Button>
        ))}
      {!items &&
        defaultItems.map((item) => (
          <Button
            className={`mx-2 my-1 justify-start ${
              router.pathname.includes(item.id)
                ? 'bg-primary text-primary-foreground hover:bg-primary/90 hover:text-primary-foreground'
                : ''
            }`}
            variant={'ghost'}
            size={'sm'}
            key={item.id}
            onClick={item.eventHandler}
          >
            {item.label}
          </Button>
        ))}
    </div>
  )
}
