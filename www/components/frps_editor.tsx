import { Label } from "@radix-ui/react-label"
import { Textarea } from "./ui/textarea"
import { FRPSFormProps } from "./frps_form"

export const FRPSEditor: React.FC<FRPSFormProps> = ({ server, serverID }) => {
	return (<div className="grid w-full gap-1.5">
		<Label className="text-sm font-medium">节点 {serverID} 配置文件`frps.json`内容</Label>
		<p className="text-sm text-muted-foreground">
			只需要配置端口和IP等字段，认证信息会由系统补全
		</p>
		<Textarea placeholder="配置文件内容" id="message" defaultValue={server?.config} />
	</div>)
}