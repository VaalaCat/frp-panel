import { Label } from "@radix-ui/react-label"
import { Textarea } from "./ui/textarea"
import { FRPCFormProps } from "./frpc_form"
import { getClient } from "@/api/client";
import { useQuery } from "@tanstack/react-query";

export const FRPCEditor: React.FC<FRPCFormProps> = ({ clientID }) => {
	const { data: client, refetch: refetchClient } = useQuery({
		queryKey: ["getClient", clientID], queryFn: () => {
			return getClient({ clientId: clientID })
		}
	});

	return (<div className="grid w-full gap-1.5">
		<Label className="text-sm font-medium">客户端 {clientID} 配置文件`frpc.json`内容</Label>
		<p className="text-sm text-muted-foreground">
			只需要配置proxies和visitors字段，认证信息和服务器连接信息会由系统补全
		</p>
		<Textarea placeholder="配置文件内容" id="message" defaultValue={client?.client?.config} />
	</div>)
}