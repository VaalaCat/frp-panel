"use client"

import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useState } from "react"
import { useTranslation } from "react-i18next"

export interface IdInputProps {
  setKeyword: (keyword: string) => void
  keyword: string
  refetchTrigger?: (randStr: string) => void
}

const makeRandomKey = () => {
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID()
  }
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}

export const IdInput: React.FC<IdInputProps> = ({ setKeyword, keyword, refetchTrigger }) => {
  const { t } = useTranslation()
  const [input, setInput] = useState(keyword)

  return (
    <div className="flex flex-row gap-2 items-center">
      <Input
        className="text-sm"
        value={input}
        placeholder={t('input.keyword.placeholder')}
        onChange={(e) => setInput(e.target.value)}
      />
      <Button
        variant="outline"
        onClick={() => {
          setKeyword(input)
          refetchTrigger?.(makeRandomKey())
        }}
      >
        {t('input.search')}
      </Button>
    </div>
  )
}
