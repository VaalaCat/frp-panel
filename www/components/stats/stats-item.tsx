'use client'

import { Card } from "@/components/ui/card"
import { ProxyInfo } from "@/lib/pb/common";
import { formatBytes } from "@/lib/utils";
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts'

const COLORS = ['#0088FE', '#00C49F'];

function preparePieData(inTraffic: bigint, outTraffic: bigint) {
  return [
    { name: 'In', value: Number(inTraffic) },
    { name: 'Out', value: Number(outTraffic) }
  ];
}

function calculateTotalTraffic(inTraffic: bigint, outTraffic: bigint): bigint {
  return inTraffic + outTraffic;
}

export default function TrafficStatsCard({ proxy }: { proxy: ProxyInfo }) {
  return <Card className="p-4 hover:bg-accent/50 transition-colors">
    <div className="items-center gap-4 grid grid-cols-4">
      {/* Server Info */}
      <div className="flex items-center gap-2 min-w-[160px]">
        <div className="w-2 h-2 rounded-full bg-green-500" />
        <span className="font-medium">{proxy.name}</span>
      </div>

      {/* Today's Traffic */}
      <div className="flex items-center gap-4 ">
        <div className="w-16 h-16">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={preparePieData(proxy.todayTrafficIn || BigInt(0), proxy.todayTrafficOut || BigInt(0))}
                cx="50%"
                cy="50%"
                innerRadius={15}
                outerRadius={30}
                paddingAngle={2}
                dataKey="value"
              >
                {preparePieData(proxy.todayTrafficIn || BigInt(0), proxy.todayTrafficOut || BigInt(0)).map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => formatBytes(value as number)} />
            </PieChart>
          </ResponsiveContainer>
        </div>
        <div className="flex flex-col gap-1">
          <div className="text-sm font-medium text-muted-foreground">Today&apos;s Traffic</div>
          <div className="flex flex-col gap-1">
            <div className="flex justify-between items-center">
              <span className="text-sm text-muted-foreground">In:</span>
              <span className="text-sm tabular-nums">{formatBytes(Number(proxy.todayTrafficIn || BigInt(0)))}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-muted-foreground">Out:</span>
              <span className="text-sm tabular-nums">{formatBytes(Number(proxy.todayTrafficOut || BigInt(0)))}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-muted-foreground">Total:</span>
              <span className="text-sm tabular-nums">{formatBytes(Number(calculateTotalTraffic(proxy.todayTrafficIn || BigInt(0), proxy.todayTrafficOut || BigInt(0))))}</span>
            </div>
          </div>
        </div>
      </div>

      {/* History Traffic */}
      <div className="flex items-center gap-4 min-w-[300px]">
        <div className="w-16 h-16">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={preparePieData(proxy.historyTrafficIn || BigInt(0), proxy.historyTrafficOut || BigInt(0))}
                cx="50%"
                cy="50%"
                innerRadius={15}
                outerRadius={30}
                paddingAngle={2}
                dataKey="value"
              >
                {preparePieData(proxy.historyTrafficIn || BigInt(0), proxy.historyTrafficOut || BigInt(0)).map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip formatter={(value) => formatBytes(value as number)} />
            </PieChart>
          </ResponsiveContainer>
        </div>
        <div className="flex flex-col gap-1">
          <div className="text-sm font-medium text-muted-foreground">History Traffic</div>
          <div className="flex flex-col gap-1">
            <div className="flex justify-between items-center">
              <span className="text-sm text-muted-foreground">In:</span>
              <span className="text-sm tabular-nums">{formatBytes(Number(proxy.historyTrafficIn || BigInt(0)))}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-muted-foreground">Out:</span>
              <span className="text-sm tabular-nums">{formatBytes(Number(proxy.historyTrafficOut || BigInt(0)))}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-muted-foreground">Total:</span>
              <span className="text-sm tabular-nums">{formatBytes(Number(calculateTotalTraffic(proxy.historyTrafficIn || BigInt(0), proxy.historyTrafficOut || BigInt(0))))}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Card>
}

