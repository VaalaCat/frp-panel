// 节点尺寸常量
export const NODE_WIDTH = 240
export const NODE_HEIGHT = 150

/**
 * 计算连接质量等级
 * 主要基于延迟判断，带宽作为参考
 */
export function calculateConnectionQuality(latency: number, bandwidth: number): 'excellent' | 'good' | 'fair' | 'poor' {
  // 优秀：延迟 0-50ms
  if (latency <= 50) return 'excellent'
  // 一般：延迟 51-200ms
  if (latency <= 200) return 'fair'
  // 过高：延迟 200ms 以上
  return 'poor'
}
