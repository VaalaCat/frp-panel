'use client'

import React from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listClient } from '@/api/client'
import { listServer } from '@/api/server'
import { useTranslation } from 'react-i18next'
import { Client, Server } from '@/lib/pb/common'
import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { useDebouncedCallback } from 'use-debounce'
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from '@/components/ui/command'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { CaretSortIcon } from '@radix-ui/react-icons'

export type AgentType = 'client' | 'server'

export interface Agent {
	id: string
	label: string
	type: AgentType
	original: Client | Server
}

export interface AgentSelectorProps {
	value?: Agent
	onChange: (agent: Agent) => void
	placeholder?: string
	className?: string
}

export const AgentSelector: React.FC<AgentSelectorProps> = ({ value, onChange, placeholder, className }) => {
	const { t } = useTranslation()
	const [open, setOpen] = React.useState(false)
	const [keyword, setKeyword] = React.useState('')

	const debounced = useDebouncedCallback((v) => {
		setKeyword(v as string)
	}, 500)

	// 获取客户端列表
	const { data: clientList, refetch: refetchClients } = useQuery({
		queryKey: ['listClient', keyword],
		queryFn: () => {
			return listClient({ page: 1, pageSize: 20, keyword: keyword })
		},
		placeholderData: keepPreviousData,
	})

	// 获取服务器列表
	const { data: serverList, refetch: refetchServers } = useQuery({
		queryKey: ['listServer', keyword],
		queryFn: () => {
			return listServer({ page: 1, pageSize: 20, keyword: keyword })
		},
		placeholderData: keepPreviousData,
	})

	// 转换为统一的 Agent 格式
	const clientAgents: Agent[] = React.useMemo(() => {
		return (clientList?.clients || []).map((client) => ({
			id: `client-${client.id}`,
			label: client.id || '',
			type: 'client' as AgentType,
			original: client,
		}))
	}, [clientList])

	const serverAgents: Agent[] = React.useMemo(() => {
		return (serverList?.servers || []).map((server) => ({
			id: `server-${server.id}`,
			label: server.id || '',
			type: 'server' as AgentType,
			original: server,
		}))
	}, [serverList])

	const handleSelect = (agent: Agent) => {
		onChange(agent)
		setOpen(false)
	}

	const handleOpenChange = (open: boolean) => {
		setOpen(open)
		if (open) {
			refetchClients()
			refetchServers()
		}
	}

	const defaultPlaceholder = t('selector.common.placeholder')

	return (
		<Popover open={open} onOpenChange={handleOpenChange}>
			<PopoverTrigger asChild>
				<Button
					variant="outline"
					role="combobox"
					aria-expanded={open}
					className={cn('w-full justify-between font-normal px-3', className, !value && 'text-muted-foreground')}
				>
					{value ? value.label : placeholder || defaultPlaceholder}
					<CaretSortIcon className="h-4 w-4 opacity-50" />
				</Button>
			</PopoverTrigger>
			<PopoverContent className="w-[--radix-popover-trigger-width] p-0" align="start">
				<Command>
					<CommandInput onValueChange={(v) => debounced(v)} placeholder={placeholder || defaultPlaceholder} />
					<CommandList>
						<CommandEmpty>{t('selector.common.notFound')}</CommandEmpty>

						{clientAgents.length > 0 && (
							<CommandGroup heading={t('canvas.panel.clients')}>
								{clientAgents.map((agent) => (
									<CommandItem key={agent.id} value={agent.label} onSelect={() => handleSelect(agent)}>
										{agent.label}
										<Check className={cn('ml-auto', value?.id === agent.id ? 'opacity-100' : 'opacity-0')} />
									</CommandItem>
								))}
							</CommandGroup>
						)}

						{serverAgents.length > 0 && (
							<CommandGroup heading={t('canvas.panel.servers')}>
								{serverAgents.map((agent) => (
									<CommandItem key={agent.id} value={agent.label} onSelect={() => handleSelect(agent)}>
										{agent.label}
										<Check className={cn('ml-auto', value?.id === agent.id ? 'opacity-100' : 'opacity-0')} />
									</CommandItem>
								))}
							</CommandGroup>
						)}
					</CommandList>
				</Command>
			</PopoverContent>
		</Popover>
	)
}
