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

export interface BaseSelectorProps {
  value?: string
  setValue: (value: string) => void
  dataList: { value: string; label: string }[]
  placeholder?: string
  label?: string
  onOpenChange?: () => void
  className?: string
}

export function BaseSelector({ value, setValue, dataList, placeholder, label, onOpenChange, className }: BaseSelectorProps) {
  return (
    <Select onValueChange={setValue} value={value} onOpenChange={onOpenChange}>
      <SelectTrigger className={cn("w-full", className)}>
        <SelectValue placeholder={placeholder || "请选择"} />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>{label}</SelectLabel>
          {
            dataList.map((item) => (
              <SelectItem key={item.value} value={item.value}>
                {item.label}
              </SelectItem>
            ))
          }
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}
