'use client'

import React from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listNetworks } from '@/api/wg'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'

export interface NetworkSelectorProps {
  networkID?: number
  setNetworkID: (id?: number) => void
  onOpenChange?: () => void
}

export const NetworkSelector: React.FC<NetworkSelectorProps> = ({ networkID, setNetworkID, onOpenChange }) => {
  const { t } = useTranslation()
  const [keyword, setKeyword] = React.useState('')

  const { data, refetch, isFetching } = useQuery({
    queryKey: ['listNetworks', keyword],
    queryFn: () =>
      listNetworks({
        page: 1,
        pageSize: 20,
        keyword: keyword || undefined,
      }),
    placeholderData: keepPreviousData,
  })

  const items = (data?.networks ?? []).map((n) => ({ value: String(n.id), label: `${n.name} (${n.cidr})` }))

  return (
    <Combobox
      placeholder={t('wg.selector.network') as string}
      dataList={items}
      value={networkID ? String(networkID) : ''}
      setValue={(v) => setNetworkID(v ? Number(v) : undefined)}
      onKeyWordChange={setKeyword}
      onOpenChange={() => refetch()}
      isLoading={isFetching}
    />
  )
}


