"use client"

import React, { useEffect } from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listServer } from '@/api/server'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'
import { Server } from '@/lib/pb/common'

export interface ServerSelectorProps {
  serverID?: string
  setServerID: (serverID: string) => void
  onOpenChange?: () => void
  setServer?: (server: Server) => void
}

export const ServerSelector: React.FC<ServerSelectorProps> = ({
  serverID,
  setServerID,
  onOpenChange,
  setServer,
}) => {
  const { t } = useTranslation()
  const [keyword, setKeyword] = React.useState('')

  const { data: serverList, refetch: refetchServers } = useQuery({
    queryKey: ['listServer', keyword],
    queryFn: () => {
      return listServer({ page: 1, pageSize: 8, keyword: keyword })
    },
    placeholderData: keepPreviousData,
  })

  const handleServerChange = (value: string) => {
    setServerID(value)
  }

  useEffect(() => {
    if (serverID) {
      setServer && setServer(serverList?.servers.find((server) => server.id == serverID) || {})
    }
  }, [serverID])

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
