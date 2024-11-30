import React from 'react'
import { Combobox } from './combobox'

export interface ProxySelectorProps {
    proxyName?: string
    setProxyname: (proxyName: string) => void
    proxyNames: string[]
}

export const ProxySelector: React.FC<ProxySelectorProps> = ({ proxyName, proxyNames ,setProxyname }) => {
    return <Combobox
        dataList={proxyNames.map((name) => ({ value: name, label: name }))}
        value={proxyName}
        setValue={setProxyname}
        notFoundText="未找到隧道"
        placeholder="隧道名称"
    />
}
