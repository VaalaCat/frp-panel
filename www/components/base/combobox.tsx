"use client"

import * as React from "react"
import { Check, ChevronsUpDown } from "lucide-react"

import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
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
}

export function Combobox({ value, setValue, dataList, placeholder, notFoundText, onOpenChange, className }: ComboboxProps) {
  const [open, setOpen] = React.useState(false)

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
            ? dataList.find((item) => item.value === value)?.label
            : (placeholder||"请选择...")}
          <ChevronsUpDown className="opacity-50 h-[12px] w-[12px]" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[200px] p-0" align="start">
        <Command>
          <CommandInput placeholder={`${placeholder||"请选择..."}`} />
          <CommandList>
            <CommandEmpty>{notFoundText||"未找到结果"}</CommandEmpty>
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
