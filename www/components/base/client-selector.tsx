import React from 'react'
import { useQuery } from '@tanstack/react-query'
import { listClient } from '@/api/client'
import { Combobox } from './combobox'

export interface ClientSelectorProps {
  clientID?: string
  setClientID: (clientID: string) => void
  onOpenChange?: () => void
}
export const ClientSelector: React.FC<ClientSelectorProps> = ({ clientID, setClientID, onOpenChange }) => {
  const handleClientChange = (value: string) => { setClientID(value) }

  const { data: clientList, refetch: refetchClients } = useQuery({
    queryKey: ['listClient'],
    queryFn: () => {
      return listClient({ page: 1, pageSize: 50, keyword: '' })
    },
  })

  return (
    <Combobox
      placeholder='客户端名称'
      dataList={clientList?.clients.map((client) => ({ value: client.id || '', label: client.id || '' })) || []}
      setValue={handleClientChange}
      value={clientID}
      onOpenChange={() => {
        onOpenChange && onOpenChange()
        refetchClients()
      }}
    />
  )
}
