"use client"

import * as React from "react"
import { Check, ChevronsUpDown } from "lucide-react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { useDebouncedCallback } from 'use-debounce'
import { useTranslation } from 'react-i18next'
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

export interface ComboboxProps {
  value?: string
  setValue: (value: string) => void
  dataList: { value: string; label: string }[]
  placeholder?: string
  notFoundText?: string
  onOpenChange?: () => void
  className?: string
  onKeyWordChange?: (keyword: string) => void
  keyword?: string
  isLoading?: boolean
}

export function Combobox({ 
  value, 
  setValue, 
  dataList, 
  placeholder, 
  notFoundText, 
  onOpenChange, 
  className, 
  keyword, 
  onKeyWordChange, 
  isLoading 
}: ComboboxProps) {
  const { t } = useTranslation()
  const [open, setOpen] = React.useState(false)
  const debounced = useDebouncedCallback(
    (v) => {
      onKeyWordChange && onKeyWordChange(v as string);
    },
    500,
  );

  const defaultPlaceholder = t('selector.common.placeholder')
  const defaultNotFoundText = t('selector.common.notFound')
  const loadingText = t('selector.common.loading')

  return (
    <Popover open={open} onOpenChange={(open) => {
      onOpenChange && onOpenChange()
      setOpen(open)}}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={cn("w-full justify-between font-normal", className)}
        >
          {value
            ? (dataList.find((item) => item.value === value)?.label || value)
            : (placeholder || defaultPlaceholder)}
          <ChevronsUpDown className="opacity-50 h-[12px] w-[12px]" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start">
        <Command>
          <CommandInput 
            value={keyword} 
            onValueChange={(v) => debounced(v)} 
            placeholder={placeholder || defaultPlaceholder} 
          />
          <CommandList>
            <CommandEmpty>{isLoading ? loadingText : (notFoundText || defaultNotFoundText)}</CommandEmpty>
            <CommandGroup>
              {dataList.map((item) => (
                <CommandItem
                  key={item.value}
                  value={item.value}
                  onSelect={(currentValue) => {
                    setValue(currentValue === value ? "" : currentValue)
                    setOpen(false)
                  }}
                >
                  {item.label}
                  <Check
                    className={cn(
                      "ml-auto",
                      value === item.value ? "opacity-100" : "opacity-0"
                    )}
                  />
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
