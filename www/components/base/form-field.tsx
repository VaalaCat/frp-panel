import React from 'react'
import { Control } from 'react-hook-form'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useTranslation } from 'react-i18next'
import StringListInput from './list-input'

export const HostField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || '127.0.0.1'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}
export const PortField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: number
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || '1234'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}
export const SecretStringField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || "secret"} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}

export const StringField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <Input className='text-sm' placeholder={placeholder || '127.0.0.1'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}

export const StringArrayField = ({
  control,
  name,
  label,
  placeholder,
  defaultValue,
}: {
  control: Control<any>
  name: string
  label: string
  placeholder?: string
  defaultValue?: string[]
}) => {
  const { t } = useTranslation()
  return (
    <FormField
      name={name}
      control={control}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t(label)}</FormLabel>
          <FormControl>
            <StringListInput placeholder={placeholder || '/path'} {...field} />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
      defaultValue={defaultValue}
    />
  )
}
