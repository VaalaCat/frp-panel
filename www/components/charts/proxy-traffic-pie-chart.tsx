"use client"

import { Label, Pie, PieChart } from "recharts"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { formatBytes } from "@/lib/utils"
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"

const chartConfig = {
  trafficIn: {
    label: "入站",
  },
  trafficOut: {
    label: "出站",
  },
} satisfies ChartConfig

export function ProxyTrafficPieChart({ trafficIn, trafficOut, title, chartLabel }:
  { trafficIn: bigint,
    trafficOut: bigint,
    title: string,
    chartLabel: string,
  }) {
  const data = [
    { type: "trafficIn", data: Number(trafficIn), fill: "hsl(var(--chart-1))" },
    { type: "trafficOut", data: Number(trafficOut), fill: "hsl(var(--chart-2))" }]

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="font-mono">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square max-h-[250px]"
        >
          <PieChart>
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent hideLabel valueFormatter={(value) => formatBytes(Number(value))} />}
            />
            <Pie data={data}
              dataKey="data"
              nameKey="type"
              innerRadius={55} strokeWidth={10}>
              <Label
                content={({ viewBox }) => {
                  if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                    return (
                      <text x={viewBox.cx} y={viewBox.cy} textAnchor="middle" dominantBaseline="middle">
                        <tspan x={viewBox.cx} y={viewBox.cy} className="fill-foreground text-xl font-bold">
                          {formatBytes(Number(trafficIn) + Number(trafficOut))}
                        </tspan>
                        <tspan x={viewBox.cx} y={(viewBox.cy || 0) + 24} className="fill-muted-foreground" >
                          {chartLabel}
                        </tspan>
                      </text>
                    )
                  }
                }}
              />
            </Pie>
          </PieChart>
        </ChartContainer>
      </CardContent>
    </Card>
  )
}

