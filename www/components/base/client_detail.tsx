"use client"

import { Popover, PopoverTrigger } from "@radix-ui/react-popover"
import { Badge } from "../ui/badge"
import { ClientStatus } from "@/lib/pb/api_master"
import { PopoverContent } from "../ui/popover"
import { useTranslation } from "react-i18next"
import { motion } from "framer-motion"
import { formatDistanceToNow } from 'date-fns'
import { zhCN, enUS } from 'date-fns/locale'

export const ClientDetail = ({ clientStatus }: { clientStatus: ClientStatus }) => {
  const { t, i18n } = useTranslation()

  const locale = i18n.language === 'zh' ? zhCN : enUS
  const connectTime = clientStatus.connectTime ? 
    formatDistanceToNow(new Date(parseInt(clientStatus.connectTime.toString())), { 
      addSuffix: true,
      locale 
    }) : '-'

  return (
    <Popover>
      <PopoverTrigger className='flex items-center'>
        <Badge 
          variant="secondary" 
          className='text-nowrap rounded-full h-6 hover:bg-secondary/80 transition-colors text-sm'
        >
          {clientStatus.version?.gitVersion || 'Unknown'}
        </Badge>
      </PopoverTrigger>
      <PopoverContent className="w-72 p-4 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-border">
        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.2 }}
        >
          <h3 className="text-base font-semibold mb-3 text-center text-foreground">
            {t('client.detail.title')}
          </h3>
          <div className="space-y-2">
            <div className="flex justify-between items-center py-1 border-b border-border">
              <span className="text-sm font-medium text-muted-foreground">{t('client.detail.version')}</span>
              <span className="text-sm text-foreground">{clientStatus.version?.gitVersion || '-'}</span>
            </div>
            <div className="flex justify-between items-center py-1 border-b border-border">
              <span className="text-sm font-medium text-muted-foreground">{t('client.detail.buildDate')}</span>
              <span className="text-sm text-foreground">{clientStatus.version?.buildDate || '-'}</span>
            </div>
            <div className="flex justify-between items-center py-1 border-b border-border">
              <span className="text-sm font-medium text-muted-foreground">{t('client.detail.goVersion')}</span>
              <span className="text-sm text-foreground">{clientStatus.version?.goVersion || '-'}</span>
            </div>
            <div className="flex justify-between items-center py-1 border-b border-border">
              <span className="text-sm font-medium text-muted-foreground">{t('client.detail.platform')}</span>
              <span className="text-sm text-foreground">{clientStatus.version?.platform || '-'}</span>
            </div>
            <div className="flex justify-between items-center py-1 border-b border-border">
              <span className="text-sm font-medium text-muted-foreground">{t('client.detail.address')}</span>
              <span className="text-sm text-foreground">{clientStatus.addr || '-'}</span>
            </div>
            <div className="flex justify-between items-center py-1 border-b border-border">
              <span className="text-sm font-medium text-muted-foreground">{t('client.detail.connectTime')}</span>
              <span className="text-sm text-foreground">{connectTime}</span>
            </div>
          </div>
        </motion.div>
      </PopoverContent>
    </Popover>
  )
}