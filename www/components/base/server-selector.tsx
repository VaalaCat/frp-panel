import React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useQuery } from '@tanstack/react-query'
import { listServer } from '@/api/server'
import { Combobox } from './combobox'

export interface ServerSelectorProps {
  serverID?: string
  setServerID: (serverID: string) => void
  onOpenChange?: () => void
}
export const ServerSelector: React.FC<ServerSelectorProps> = ({ serverID, setServerID, onOpenChange }) => {
  const handleServerChange = (value: string) => { setServerID(value) }

  const { data: serverList, refetch: refetchServers } = useQuery({
    queryKey: ['listServer'],
    queryFn: () => {
      return listServer({ page: 1, pageSize: 50, keyword: '' })
    },
  })

  return (<Combobox
    placeholder='服务端名称'
    value={serverID}
    setValue={handleServerChange}
    dataList={serverList?.servers.map((server) => ({ value: server.id || '', label: server.id || '' })) || []}
    onOpenChange={() => {
      onOpenChange && onOpenChange()
      refetchServers()
    }}
  />)
}
