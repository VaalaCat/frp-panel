"use client"

import React from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listClient } from '@/api/client'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'

export interface ClientSelectorProps {
  clientID?: string
  setClientID: (clientID: string) => void
  onOpenChange?: () => void
}

export const ClientSelector: React.FC<ClientSelectorProps> = ({ 
  clientID, 
  setClientID, 
  onOpenChange 
}) => {
  const { t } = useTranslation()
  const handleClientChange = (value: string) => { setClientID(value) }
  const [keyword, setKeyword] = React.useState('')

  const { data: clientList, refetch: refetchClients } = useQuery({
    queryKey: ['listClient', keyword],
    queryFn: () => {
      return listClient({ page: 1, pageSize: 8, keyword: keyword })
    },
    placeholderData: keepPreviousData,
  })

  return (
    <Combobox
      placeholder={t('selector.client.placeholder')}
      dataList={clientList?.clients.map((client) => ({ 
        value: client.id || '', 
        label: client.id || '' 
      })) || []}
      setValue={handleClientChange}
      value={clientID}
      onKeyWordChange={setKeyword}
      onOpenChange={() => {
        onOpenChange && onOpenChange()
        refetchClients()
      }}
    />
  )
}
