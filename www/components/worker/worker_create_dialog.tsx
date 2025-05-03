import React from 'react'
import { useMutation } from '@tanstack/react-query'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { createWorker } from '@/api/worker'
import { CreateWorkerRequest } from '@/lib/pb/api_client'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Form, FormField, FormItem, FormLabel, FormControl } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { ClientSelector } from '../base/client-selector'

const CreateWorkerSchema = z.object({
  clientId: z.string().min(1, 'worker.createClientRequired'),
  name: z.string().min(1, 'worker.createNameRequired'),
})

type CreateWorkerValues = z.infer<typeof CreateWorkerSchema>

export interface CreateWorkerDialogProps {
  refetchTrigger: React.Dispatch<React.SetStateAction<string>>
}

export const CreateWorkerDialog: React.FC<CreateWorkerDialogProps> = ({ refetchTrigger }) => {
  const { t } = useTranslation()
  const [open, setOpen] = React.useState(false)

  const form = useForm<CreateWorkerValues>({
    resolver: zodResolver(CreateWorkerSchema),
    defaultValues: {
      clientId: '',
      name: '',
    },
  })

  const { mutate, isPending } = useMutation({
    mutationFn: (values: CreateWorkerValues) => {
      const req: CreateWorkerRequest = {
        clientId: values.clientId,
        worker: { name: values.name },
      }
      return createWorker(req)
    },
    onSuccess: () => {
      toast.success(t('worker.create.success'))
      form.reset()
      setOpen(false)
      refetchTrigger(new Date().toISOString())
    },
    onError: (err: any) => {
      toast(err?.message || t('worker.create'))
    },
  })

  const onSubmit = (values: CreateWorkerValues) => {
    mutate(values)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button variant="outline">{t('worker.create.button')}</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t('worker.create.title')}</DialogTitle>
          <DialogDescription>{t('worker.create.description')}</DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="clientId"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('worker.create.clientIdLabel')}</FormLabel>
                  <FormControl>
                    <ClientSelector clientID={field.value} setClientID={field.onChange} />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>{t('worker.create.nameLabel')}</FormLabel>
                  <FormControl>
                    <Input placeholder={t('worker.create.namePlaceholder')} {...field} />
                  </FormControl>
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button type="submit" disabled={isPending}>
                {isPending ? t('worker.create.creating') : t('worker.create.submit')}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  )
}
