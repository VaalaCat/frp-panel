'use client'

import React, { useEffect } from 'react'
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import { listWireGuards } from '@/api/wg'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'

export interface WireGuardSelectorProps {
  clientID?: string
  networkID?: number
  wireguardID?: number
  setWireguardID: (id?: number) => void
  onOpenChange?: () => void
}

export const WireGuardSelector: React.FC<WireGuardSelectorProps> = ({ clientID, networkID, wireguardID, setWireguardID, onOpenChange }) => {
  const { t } = useTranslation()
  const [keyword, setKeyword] = React.useState('')

  const [valueKey, setValueKey] = React.useState('')

  const { data, refetch, isFetching } = useQuery({
    queryKey: ['listWireGuards', clientID, networkID, keyword],
    queryFn: () =>
      listWireGuards({
        page: 1,
        pageSize: 50,
        clientId: clientID,
        networkId: networkID,
        keyword: keyword || undefined,
      }),
    placeholderData: keepPreviousData,
  })

  const items = (data?.wireguardConfigs ?? []).map((w) => ({
    value: `${w.id ?? ''}`,
    label: `${w.clientId ?? ''} ${w.localAddress ? `(${w.localAddress})` : ''}`.trim(),
  }))

  useEffect(() => {
    if (wireguardID && data?.wireguardConfigs) {
      const target = data.wireguardConfigs.find((w) => w.id === wireguardID)
      if (target) {
        setValueKey(String(target.id))
      }
    }
  }, [wireguardID, data?.wireguardConfigs])

  useEffect(() => {
    if (valueKey) {
      setWireguardID(Number(valueKey))
    } else {
      setWireguardID(undefined)
    }
  }, [valueKey, setWireguardID])

  return (
    <Combobox
      placeholder={t('wg.selector.clientWireguards') as string}
      dataList={items}
      value={valueKey}
      setValue={(v) => setValueKey(v)}
      onKeyWordChange={setKeyword}
      onOpenChange={() => refetch()}
      isLoading={isFetching}
    />
  )
}


