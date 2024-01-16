import { ProxyType } from "@/types/proxy"
import React from "react"
import { useState } from "react"
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select"
import { Label } from "@radix-ui/react-label"
import { HTTPProxyForm, TCPProxyForm, UDPProxyForm } from "./proxy_form"


export interface FRPCFormProps {
	clientID: string
	serverID: string
}

export const FRPCForm: React.FC<FRPCFormProps> = ({ clientID, serverID }) => {
	const [proxyType, setProxyType] = useState<ProxyType>('http')
	const handleTypeChange = (value: string) => {
		setProxyType(value as ProxyType)
	}

	return (
		<div className="flex flex-col w-full pt-2">
			<Select onValueChange={handleTypeChange} defaultValue="http">
				<Label className="text-sm font-medium">协议</Label>
				<SelectTrigger className="my-2">
					<SelectValue placeholder="类型" />
				</SelectTrigger>
				<SelectContent>
					<SelectItem value="http">http</SelectItem>
					<SelectItem value="tcp">tcp</SelectItem>
					<SelectItem value="udp">udp</SelectItem>
				</SelectContent>
			</Select>
			{
				proxyType === 'tcp' && serverID && clientID &&
				<TCPProxyForm
					serverID={serverID}
					clientID={clientID}
				/>
			}
			{
				proxyType === 'udp' && serverID && clientID &&
				<UDPProxyForm
					serverID={serverID}
					clientID={clientID}
				/>
			}
			{
				proxyType === 'http' && serverID && clientID &&
				<HTTPProxyForm
					serverID={serverID}
					clientID={clientID}
				/>
			}
		</div >
	)
}
