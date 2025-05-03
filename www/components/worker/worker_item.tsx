import { MoreHorizontal, ArrowUpRight } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { removeWorker, getWorker, getWorkerIngress, getWorkerStatus } from '@/api/worker'
import { $workerTableRefetchTrigger } from '@/store/refetch-trigger'
import { useTranslation } from 'react-i18next'
import { useRouter } from 'next/router'
import { Worker } from '@/lib/pb/common'
import { toast } from 'sonner'
import { useMutation, useQuery } from '@tanstack/react-query'
import { WorkerStatus } from './worker_status'
import { ColumnDef, Row } from '@tanstack/react-table'

export type WorkerTableSchema = {
  workerId: string
  name: string
  userId: number
  tenantId: number
  socketAddress: string
  origin: Worker
}

export const columns: ColumnDef<WorkerTableSchema>[] = [
  {
    accessorKey: 'name',
    header: ({ column }: { column: ColumnDef<WorkerTableSchema> }) => {
      // eslint-disable-next-line react-hooks/rules-of-hooks
      const { t } = useTranslation()
      return t('worker.columns.name')
    },
    cell: ({ row }: { row: Row<WorkerTableSchema> }) => {
      const worker = row.original
      // eslint-disable-next-line react-hooks/rules-of-hooks
      const router = useRouter()

      return (
        <div className="flex items-center">
          <span className="mr-2 font-medium">{worker.name || worker.workerId}</span>
          <Button
            variant="ghost"
            size="icon"
            className="h-6 w-6 rounded-full"
            onClick={() =>
              router.push({
                pathname: '/worker-edit',
                query: { workerId: worker.workerId },
              })
            }
          >
            <ArrowUpRight className="h-3.5 w-3.5" />
          </Button>
        </div>
      )
    },
  },
  {
    id: 'status',
    header: ({ column }: { column: ColumnDef<WorkerTableSchema> }) => {
      // eslint-disable-next-line react-hooks/rules-of-hooks
      const { t } = useTranslation()
      return t('worker.columns.status')
    },
    cell: ({ row }: { row: Row<WorkerTableSchema> }) => {
      const workerId = row.getValue('workerId') as string

      // eslint-disable-next-line react-hooks/rules-of-hooks
      const { data: workerData } = useQuery({
        queryKey: ['getWorker', workerId],
        queryFn: () => getWorker({ workerId }),
        enabled: !!workerId,
      })

      return (
        <div className="flex justify-start">
          <WorkerStatus workerId={workerId} clients={workerData?.clients || []} />
        </div>
      )
    },
  },
  {
    accessorKey: 'workerId',
    header: ({ column }: { column: ColumnDef<WorkerTableSchema> }) => {
      // eslint-disable-next-line react-hooks/rules-of-hooks
      const { t } = useTranslation()
      return t('worker.columns.id')
    },
    cell: ({ row }: { row: Row<WorkerTableSchema> }) => {
      return <div className="font-mono text-sm text-nowarp whitespace-nowrap">{row.getValue('workerId')}</div>
    },
  },
  {
    id: 'actions',
    cell: ({ row }: { row: Row<WorkerTableSchema> }) => {
      const worker = row.original
      return <WorkerActions worker={worker} />
    },
  },
]

interface WorkerActionsProps {
  worker: WorkerTableSchema
}

export const WorkerActions: React.FC<WorkerActionsProps> = ({ worker }) => {
  const { t } = useTranslation()
  const router = useRouter()
  const del = useMutation({
    mutationFn: () => removeWorker({ workerId: worker.workerId }),
    onSuccess: () => {
      toast.success(t('worker.actions_menu.delete') + t('common.success'))
      $workerTableRefetchTrigger.set(Math.random())
    },
    onError: (err: any) => {
      toast(t('common.failed'), { description: err.message })
      $workerTableRefetchTrigger.set(Math.random())
    },
  })

  return (
    <Dialog>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="icon">
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end">
          <DropdownMenuLabel>{t('worker.actions_menu.title')}</DropdownMenuLabel>
          <DropdownMenuItem
            onClick={() =>
              router.push({
                pathname: '/worker-edit',
                query: { workerId: worker.workerId },
              })
            }
          >
            {t('worker.actions_menu.edit')}
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DialogTrigger asChild>
            <DropdownMenuItem className="text-destructive">{t('worker.actions_menu.delete')}</DropdownMenuItem>
          </DialogTrigger>
        </DropdownMenuContent>
      </DropdownMenu>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>{t('worker.delete.title')}</DialogTitle>
          <DialogDescription>{t('worker.delete.description', { name: worker.name })}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <DialogClose asChild>
            <Button variant="outline" className="mr-2">
              {t('common.cancel')}
            </Button>
          </DialogClose>
          <DialogClose asChild>
            <Button variant="destructive" onClick={() => del.mutate()}>
              {t('worker.delete.confirm')}
            </Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
