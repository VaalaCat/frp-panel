"use client"

import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useTranslation } from 'react-i18next'
import { Languages } from 'lucide-react'
import { useStore } from '@nanostores/react'
import { $language } from '@/store/user'
import { useEffect } from "react"

export function LanguageSwitcher() {
  const { i18n } = useTranslation()
  const currentLanguage = useStore($language)

  useEffect(() => {
    i18n.changeLanguage(currentLanguage)
  }, [currentLanguage, i18n])

  const toggleLanguage = (lang: string) => {
    $language.set(lang)
    i18n.changeLanguage(lang)
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon">
          <Languages className="h-5 w-5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem 
          onClick={() => toggleLanguage('zh')}
          className={currentLanguage === 'zh' ? 'bg-accent' : ''}
        >
          中文
        </DropdownMenuItem>
        <DropdownMenuItem 
          onClick={() => toggleLanguage('en')}
          className={currentLanguage === 'en' ? 'bg-accent' : ''}
        >
          English
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
