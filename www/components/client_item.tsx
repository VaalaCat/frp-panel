import { ColumnDef, Table } from "@tanstack/react-table"
import { MoreHorizontal } from "lucide-react"
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import { Button } from "@/components/ui/button"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useToast } from "./ui/use-toast"
import React, { useState } from "react"
import { useSearchParams } from "next/navigation"
import { ExecCommandStr } from "@/lib/consts"
import { useMutation, useQuery } from "@tanstack/react-query"
import { deleteClient, listClient } from "@/api/client"
import { useRouter } from "next/router"

export type ClientTableSchema = {
    id: string,
    status: "invalid" | "valid"
    secret: string
    config?: string
}

export const columns: ColumnDef<ClientTableSchema>[] = [
    {
        accessorKey: "id",
        header: "ID",
        cell: ({ row }) => {
            return <div className="font-mono">{row.original.id}</div>
        }
    },
    {
        accessorKey: "status",
        header: "状态",
        cell: ({ row }) => {
            const client = row.original
            return <div className={`font-medium ${client.status === "valid" ? "text-green-500" : "text-red-500"} min-w-12`}>{
                {
                    valid: "已配置",
                    invalid: "未配置",
                }[client.status]
            }</div>
        }
    },
    {
        accessorKey: "secret",
        header: "连接密钥",
        cell: ({ row }) => {
            const client = row.original
            return <ClientSecret client={client} />
        }
    },
    {
        id: "action",
        cell: ({ row, table }) => {
            const client = row.original
            return (<ClientActions client={client} table={table} />)
        },
    },
]

export const ClientSecret = ({ client }: { client: ClientTableSchema }) => {
    const [showSecrect, setShowSecrect] = useState<boolean>(false)
    const fakeSecret = Array.from({ length: client.secret.length }).map(() => '*').join('')
    const { toast } = useToast()
    return <div
        onMouseEnter={() => setShowSecrect(true)}
        onMouseLeave={() => setShowSecrect(false)}
        onClick={() => {
            navigator.clipboard.writeText(ExecCommandStr("client", client));
            toast({ description: "复制成功", });
        }}
        className="font-medium hover:rounded hover:bg-slate-100 p-2 font-mono">{
            showSecrect ? client.secret : fakeSecret
        }</div>
}

export interface ClientItemProps {
    client: ClientTableSchema
    table: Table<ClientTableSchema>
}

export const ClientActions: React.FC<ClientItemProps> = ({ client, table }) => {
    const { toast } = useToast()
    const router = useRouter();
    const fetchDataOptions = {
        pageIndex: table.getState().pagination.pageIndex,
        pageSize: table.getState().pagination.pageSize,
    }

    const dataQuery = useQuery({
        queryKey: ["listClient", fetchDataOptions],
        queryFn: async () => {
            return await listClient({
                page: fetchDataOptions.pageIndex + 1,
                pageSize: fetchDataOptions.pageSize
            })
        }
    })

    const removeClient = useMutation({
        mutationFn: deleteClient,
        onSuccess: () => {
            toast({ description: "删除成功" })
            dataQuery.refetch()
        },
        onError: () => {
            toast({ description: "删除失败" })
        }
    })

    return <Dialog>
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <Button variant="ghost" className="h-8 w-8 p-0">
                    <span className="sr-only">打开菜单</span>
                    <MoreHorizontal className="h-4 w-4" />
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                <DropdownMenuLabel>操作</DropdownMenuLabel>
                <DropdownMenuItem
                    onClick={() => {
                        navigator.clipboard.writeText(ExecCommandStr("client", client));
                        toast({ description: "复制成功", });
                    }}
                >
                    复制启动命令
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem onClick={() => {
                    router.push({ pathname: "/clientedit", query: { clientID: client.id } })
                }}>修改</DropdownMenuItem>
                <DialogTrigger asChild>
                    <DropdownMenuItem className="text-destructive">删除</DropdownMenuItem>
                </DialogTrigger>
            </DropdownMenuContent>
        </DropdownMenu>
        <DialogContent>
            <DialogHeader>
                <DialogTitle>确定删除该客户端?</DialogTitle>
                <DialogDescription>
                    <p className="text-destructive">
                        此操作无法撤消。您确定要永久从我们的服务器中删除该客户端?
                    </p>
                    <p className="text-gray-500 border-l-4 border-gray-500 pl-4 py-2">删除后运行中的客户端将无法通过现有参数再次连接，如果您需要删除客户端对外的连接，可以选择清空配置</p>
                </DialogDescription>
            </DialogHeader>
            <DialogFooter>
                <DialogClose asChild>
                    <Button type="submit" onClick={() => removeClient.mutate({ clientId: client.id })}>确定</Button>
                </DialogClose>
            </DialogFooter>
        </DialogContent>
    </Dialog>
}
