"use client"

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { ProxyInfo } from "@/lib/pb/common"
import { formatBytes } from "@/lib/utils"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"

export function ProxyTrafficBarChart({ proxyInfo }:{ proxyInfo: ProxyInfo }) {
  const data = [
    {
      name: "入站",
      Today: Number(proxyInfo.todayTrafficIn),
      History: Number(proxyInfo.historyTrafficIn),
    },
    {
      name: "出站",
      Today: Number(proxyInfo.todayTrafficOut),
      History: Number(proxyInfo.historyTrafficOut),
    },
  ]

  return (
    <Card>
      <CardHeader>
        <CardTitle>流量详情</CardTitle>
      </CardHeader>
      <CardContent>
        <ChartContainer config={{}} className="h-[300px] w-full font-mono">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={data}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis tickLine={false} tickMargin={10} axisLine={false} dataKey="name" />
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent hideLabel valueFormatter={(value) => formatBytes(Number(value))} />}
              />
              <YAxis tickFormatter={(value) => formatBytes(Number(value))} />
              <Tooltip labelClassName="font-mono" wrapperClassName="font-mono" formatter={(value) => formatBytes(Number(value))} />
              <Legend />
              <Bar dataKey="Today" fill="hsl(var(--chart-1))" radius={4} />
              <Bar dataKey="History" fill="hsl(var(--chart-2))" radius={4} />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}

