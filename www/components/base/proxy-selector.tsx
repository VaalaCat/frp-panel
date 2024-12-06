"use client"

import React from 'react'
import { Combobox } from './combobox'
import { useTranslation } from 'react-i18next'

export interface ProxySelectorProps {
    proxyName?: string
    setProxyname: (proxyName: string) => void
    proxyNames: string[]
}

export const ProxySelector: React.FC<ProxySelectorProps> = ({ 
    proxyName, 
    proxyNames, 
    setProxyname 
}) => {
    const { t } = useTranslation()

    return (
        <Combobox
            dataList={proxyNames.map((name) => ({ 
                value: name, 
                label: name 
            }))}
            value={proxyName}
            setValue={setProxyname}
            notFoundText={t('selector.proxy.notFound')}
            placeholder={t('selector.proxy.placeholder')}
        />
    )
}
