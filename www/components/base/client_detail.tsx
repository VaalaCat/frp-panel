import { Popover, PopoverTrigger } from "@radix-ui/react-popover"
import { Badge } from "../ui/badge"
import { ClientStatus } from "@/lib/pb/api_master"
import { PopoverContent } from "../ui/popover"

export const ClientDetail = ({ clientStatus }: { clientStatus: ClientStatus }) => {
  return (
    <Popover>
      <PopoverTrigger className='flex items-center'>
        <Badge variant={"secondary"} className='text-nowrap rounded-full h-6'>
          {clientStatus.version?.gitVersion}
        </Badge>
      </PopoverTrigger>
      <PopoverContent className="w-fit overflow-auto max-w-72 max-h-72 p-4 bg-white rounded-lg shadow-lg">
        <h3 className="text-lg font-semibold mb-4 text-center">客户端信息</h3>
        <div className="flex justify-between mb-2">
          <span className="font-medium">版本:</span>
          <span>{clientStatus.version?.gitVersion}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span className="font-medium">编译时间:</span>
          <span>{clientStatus.version?.buildDate}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span className="font-medium">Go版本:</span>
          <span>{clientStatus.version?.goVersion}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span className="font-medium">客户端平台:</span>
          <span>{clientStatus.version?.platform}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span className="font-medium">客户端地址:</span>
          <span>{clientStatus.addr}</span>
        </div>
        <div className="flex justify-between mb-2">
          <span className="font-medium">连接时间:</span>
          <span>{clientStatus.connectTime}</span>
        </div>
      </PopoverContent>
    </Popover>
  )
}