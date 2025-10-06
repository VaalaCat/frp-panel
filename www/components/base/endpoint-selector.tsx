'use client'

import React from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listEndpoints } from '@/api/wg'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'

export interface EndpointSelectorProps {
  clientID: string
  endpointID?: number
  setEndpointID: (id?: number) => void
  onOpenChange?: () => void
}

export const EndpointSelector: React.FC<EndpointSelectorProps> = ({ clientID, endpointID, setEndpointID, onOpenChange }) => {
  const { t } = useTranslation()
  const [keyword, setKeyword] = React.useState('')

  const { data, refetch, isFetching } = useQuery({
    queryKey: ['listEndpoints', clientID, keyword],
    queryFn: () =>
      listEndpoints({
        page: 1,
        pageSize: 50,
        clientId: clientID,
        keyword: keyword || undefined,
      }),
    placeholderData: keepPreviousData,
    enabled: !!clientID,
  })

  const items = (data?.endpoints ?? []).map((e) => ({ value: String(e.id), label: `${e.host}:${e.port}` }))

  return (
    <Combobox
      placeholder={t('wg.selector.endpoint') as string}
      dataList={items}
      value={endpointID ? String(endpointID) : ''}
      setValue={(v) => setEndpointID(v ? Number(v) : undefined)}
      onKeyWordChange={setKeyword}
      onOpenChange={() => refetch()}
      isLoading={isFetching}
    />
  )
}


