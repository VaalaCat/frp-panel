"use client"

import * as React from "react"

import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { cn } from "@/lib/utils"
import { useTranslation } from "react-i18next"

export interface BaseSelectorProps {
  value?: string
  setValue: (value: string) => void
  dataList: { value: string; label: string }[]
  placeholder?: string
  label?: string
  onOpenChange?: () => void
  className?: string
}

export function BaseSelector({ 
  value, 
  setValue, 
  dataList, 
  placeholder, 
  label, 
  onOpenChange, 
  className 
}: BaseSelectorProps) {
  const { t } = useTranslation()
  const defaultPlaceholder = t('selector.common.placeholder')

  return (
    <Select onValueChange={setValue} value={value} onOpenChange={onOpenChange}>
      <SelectTrigger className={cn("w-full", className)}>
        <SelectValue placeholder={placeholder || defaultPlaceholder} />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          {label && <SelectLabel>{label}</SelectLabel>}
          {dataList.map((item) => (
            <SelectItem key={item.value} value={item.value}>
              {item.label}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
