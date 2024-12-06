"use client"

import React from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listServer } from '@/api/server'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'

export interface ServerSelectorProps {
  serverID?: string
  setServerID: (serverID: string) => void
  onOpenChange?: () => void
}

export const ServerSelector: React.FC<ServerSelectorProps> = ({ 
  serverID, 
  setServerID, 
  onOpenChange 
}) => {
  const { t } = useTranslation()
  const handleServerChange = (value: string) => { setServerID(value) }
  const [keyword, setKeyword] = React.useState('')

  const { data: serverList, refetch: refetchServers } = useQuery({
    queryKey: ['listServer', keyword],
    queryFn: () => {
      return listServer({ page: 1, pageSize: 8, keyword: keyword })
    },
    placeholderData: keepPreviousData,
  })

  return (
    <Combobox
      placeholder={t('selector.server.placeholder')}
      value={serverID}
      setValue={handleServerChange}
      dataList={serverList?.servers.map((server) => ({ 
        value: server.id || '', 
        label: server.id || '' 
      })) || []}
      onKeyWordChange={setKeyword}
      onOpenChange={() => {
        onOpenChange && onOpenChange()
        refetchServers()
      }}
    />
  )
}
