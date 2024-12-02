import { User } from '@/lib/pb/common'
import Avatar, { genConfig } from 'react-nice-avatar'

export function UserAvatar({ userInfo, className }: { userInfo: User | undefined, className?: string }) {
	return <Avatar shape="rounded" className={className || "w-10 h-10"} {...genConfig(userInfo?.email)}></Avatar>
}