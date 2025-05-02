import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { LucideIcon } from 'lucide-react'

export interface DropdownMenuProps {
  menuGroup: {
    name: string
    onClick: () => void
    icon?: LucideIcon
    className?: string
  }[][]
  title: string
  trigger: React.ReactNode
  extraButtons?: React.ReactNode
}

export function BaseDropdownMenu({ menuGroup, title, trigger, extraButtons }: DropdownMenuProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>{trigger || <Button variant="outline">Open</Button>}</DropdownMenuTrigger>
      <DropdownMenuContent className="w-fit">
        <DropdownMenuLabel>{title}</DropdownMenuLabel>
        <DropdownMenuSeparator />
        {menuGroup.map((items, id1) => (
          <DropdownMenuGroup key={id1}>
            {items.map((item, id2) => (
              <DropdownMenuItem onClick={item.onClick} key={id2} className={item.className}>
                {/* {<>{item.icon}</>} */}
                {item.name}
              </DropdownMenuItem>
            ))}
          </DropdownMenuGroup>
        ))}
        {extraButtons && <DropdownMenuGroup>{extraButtons}</DropdownMenuGroup>}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
