import { ProxyType, TypedProxyConfig } from "@/types/proxy"
import { useTranslation } from "react-i18next"
import { Button } from "../ui/button"
import { BaseDropdownMenu } from "../base/drop-down-menu"
import { deleteProxyConfig } from "@/api/proxy"
import { useMutation } from "@tanstack/react-query"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { ProxyConfigMutateForm } from "./mutate_proxy_config"
import { useEffect, useState } from "react"
import { Row } from "@tanstack/react-table"
import { ProxyConfigTableSchema } from "./proxy_config_item"
import { MoreHorizontal } from "lucide-react"
import { toast } from "sonner"
import { $proxyTableRefetchTrigger } from "@/store/refetch-trigger"

export interface ProxyConfigActionsProps {
  serverID: string
  clientID: string
  name: string
  row: Row<ProxyConfigTableSchema>
}

export function ProxyConfigActions({ serverID, clientID, name, row }: ProxyConfigActionsProps) {
  const { t } = useTranslation()
  const [proxyMutateFormOpen, setProxyMutateFormOpen] = useState(false)
  const [deleteWarnDialogOpen, setDeleteWarnDialogOpen] = useState(false)

  const deleteProxyConfigMutation = useMutation({
    mutationKey: ['deleteProxyConfig', serverID, clientID, name],
    mutationFn: () => deleteProxyConfig({
      serverId: serverID,
      clientId: clientID,
      name,
    }),
    onSuccess: () => {
      toast(t('proxy.action.delete_success'))
      $proxyTableRefetchTrigger.set(Math.random())
    },
    onError: (e) => {
      toast(t('proxy.action.delete_failed'), {
        description: JSON.stringify(e),
      })
      $proxyTableRefetchTrigger.set(Math.random())
    },
  })

  const menuActions = [[
    {
      name: t('proxy.action.edit'),
      onClick: () => { setProxyMutateFormOpen(true) },
    },
    {
      name: t('proxy.action.delete'),
      onClick: () => { setDeleteWarnDialogOpen(true) },
      className: 'text-destructive',
    },
  ]]
  return (<>
    <Dialog open={proxyMutateFormOpen} onOpenChange={setProxyMutateFormOpen}>
      <DialogContent className="max-h-screen overflow-auto">
        <ProxyConfigMutateForm
          disableChangeProxyName
          defaultProxyConfig={JSON.parse(row.original.config || '{}') as TypedProxyConfig}
          overwrite={true}
          defaultOriginalProxyConfig={row.original.originalProxyConfig}
        />
      </DialogContent>
    </Dialog>
    <Dialog open={deleteWarnDialogOpen} onOpenChange={setDeleteWarnDialogOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('proxy.action.delete_tunnel')}</DialogTitle>
          <DialogDescription>
            <p className="text-destructive">{t('proxy.action.delete_attention_title')}</p>
            <p className="text-gray-500 border-l-4 border-gray-500 pl-4 py-2">
              {t('proxy.action.delete_attention_description')}
            </p>
          </DialogDescription>
        </DialogHeader>
        <DialogFooter><DialogClose asChild><Button type="submit" onClick={() => {
          deleteProxyConfigMutation.mutate()
        }}>
          {t('proxy.action.delete_attention_confirm')}
        </Button></DialogClose></DialogFooter>
      </DialogContent>
    </Dialog>
    <BaseDropdownMenu
      menuGroup={menuActions}
      title={t('proxy.action.title')}
      trigger={<Button variant="ghost" className="h-8 w-8 p-0">
        <MoreHorizontal className="h-4 w-4" />
      </Button>} />
  </>)
}