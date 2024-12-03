"use client"
import { FitAddon } from '@xterm/addon-fit'
import { useEffect, useState } from 'react'
import { useXTerm } from 'react-xtermjs'

const LogTerminalComponent = ({ logs, reset }: { logs: string, reset: number }) => {
  const { instance, ref } = useXTerm()
  const fitAddon = new FitAddon()
  const [prevLogs, setPrevLogs] = useState('')
  const [currentLine, setCurrentLine] = useState('')

  function removePrefix(str: string, prefix: string): string {
    if (str.startsWith(prefix)) {
      return str.slice(prefix.length);
    }
    return str;
  }

  useEffect(() => {
    const c = removePrefix(logs, prevLogs)
    c.trim()
    setPrevLogs(logs)
    if (c.length === 0) {
      return
    }
    const splittedLines = c.split('\n').map((line) => {
      return line.trim().replaceAll("\r", "").replaceAll("\n", "").replaceAll("\t", "");
    }).filter((line) => {
      return line.length > 0
    })

    setCurrentLine(splittedLines.join('\r\n'))
  }, [logs, prevLogs])

  useEffect(() => {
    if (instance) {
      instance.reset()
    }
  }, [reset, ref, instance])

  useEffect(() => {
    
    // Load the fit addon
    instance?.loadAddon(fitAddon)

    const handleResize = () => fitAddon.fit()
    fitAddon.fit()

    instance?.writeln(currentLine)
    // Handle resize event
    window.addEventListener('resize', handleResize)
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }, [ref, instance, currentLine])

  return <div ref={ref} style={{ height: '100%', width: '100%'}} />
}

export default LogTerminalComponent