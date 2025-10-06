'use client'
import { ClientStatus } from '@/lib/pb/api_master'
import { terminalWebsocketUrl } from '@/lib/terminal'
import { FitAddon } from '@xterm/addon-fit'
import { CanvasAddon } from '@xterm/addon-canvas'
import { Unicode11Addon } from '@xterm/addon-unicode11'
import { useEffect } from 'react'
import { useXTerm } from 'react-xtermjs'

export interface TerminalComponentProps {
  isLoading: boolean
  clientStatus?: ClientStatus
  reset: number
  setStatus: (status: 'loading' | 'success' | 'error' | undefined) => void
}

const TerminalComponent = ({ isLoading, clientStatus, reset, setStatus }: TerminalComponentProps) => {
  const { instance: terminal, ref } = useXTerm()

  const fitAddon = new FitAddon()

  useEffect(() => {
    if (terminal) {
      terminal.reset()
    }
  }, [reset, ref, terminal])

  useEffect(() => {
    if (!terminal) {
      return
    }

    terminal.options.allowProposedApi = true
    terminal.options.cursorStyle = 'block'

    terminal.loadAddon(new CanvasAddon())
    terminal.loadAddon(fitAddon)
    terminal.loadAddon(new Unicode11Addon())
    terminal.unicode.activeVersion = '11'

    const handleResize = () => fitAddon.fit()

    fitAddon.fit()
    fitAddon.fit()

    // 监听窗口大小改变
    window.addEventListener('resize', handleResize)

    // 监听容器大小改变
    let resizeObserver: ResizeObserver | null = null
    if (ref.current) {
      resizeObserver = new ResizeObserver(() => {
        // 使用 requestAnimationFrame 来避免性能问题
        requestAnimationFrame(() => {
          fitAddon.fit()
        })
      })
      resizeObserver.observe(ref.current)
    }

    return () => {
      window.removeEventListener('resize', handleResize)
      if (resizeObserver) {
        resizeObserver.disconnect()
      }
    }
  }, [ref, terminal])

  useEffect(() => {
    if (!terminal) {
      return
    }

    // The terminal should be cleared on each reconnect
    // because all data is re-rendered from the backend.
    terminal.clear()

    // Focusing on connection allows users to reload the page and start
    // typing immediately.
    terminal.focus()

    // Disable input while we connect.
    terminal.options.disableStdin = true

    // Show a message if we failed to find the workspace or agent.
    if (isLoading) {
      return
    }

    if (!clientStatus) {
      terminal.writeln(`no client found with ID, is the program started?`)
      setStatus('error')
      return
    }

    // Hook up terminal events to the websocket.
    let websocket: WebSocket | null
    const disposers = [
      terminal.onData((data) => {
        websocket?.send(new TextEncoder().encode(JSON.stringify({ data: data })))
      }),
      terminal.onResize((event) => {
        try {
          websocket?.send(
            new TextEncoder().encode(
              JSON.stringify({
                height: event.rows,
                width: event.cols,
              }),
            ),
          )
        } catch (error) {
          console.error('Failed to send resize event:', error)
        }
      }),
    ]

    let disposed = false

    // Open the web socket and hook it up to the terminal.
    terminalWebsocketUrl(clientStatus.clientId, terminal.rows, terminal.cols)
      .then((url) => {
        if (disposed) {
          return // Unmounted while we waited for the async call.
        }
        websocket = new WebSocket(url)
        websocket.binaryType = 'arraybuffer'
        websocket.addEventListener('open', () => {
          // Now that we are connected, allow user input.
          terminal.options = {
            disableStdin: false,
            windowsMode: clientStatus.version?.platform.includes('windows'),
          }
          // Send the initial size.
          websocket?.send(
            new TextEncoder().encode(
              JSON.stringify({
                height: terminal.rows,
                width: terminal.cols,
              }),
            ),
          )
        })
        websocket.addEventListener('error', () => {
          terminal.options.disableStdin = true
          terminal.writeln(`socket errored`)
          setStatus('error')
        })
        websocket.addEventListener('close', () => {
          terminal.options.disableStdin = true
          setStatus(undefined)
        })
        websocket.addEventListener('message', (event) => {
          if (typeof event.data === 'string') {
            // This exclusively occurs when testing.
            // "jest-websocket-mock" doesn't support ArrayBuffer.
            terminal.write(event.data)
          } else {
            terminal.write(new Uint8Array(event.data))
          }
          setStatus('success')
        })
      })
      .catch((error) => {
        setStatus('error')
        if (disposed) {
          return // Unmounted while we waited for the async call.
        }
        terminal.writeln(error.message)
      })

    return () => {
      disposed = true // Could use AbortController instead?
      for (const d of disposers) {
        d.dispose()
      }
      websocket?.close(1000)
    }
  }, [terminal, isLoading, setStatus])

  return <div ref={ref} style={styles.terminal()} />
}

const styles = {
  terminal: () => ({
    width: '100%',
    height: '100%',
    overflow: 'hidden',
    flex: 1,
  }),
}

export default TerminalComponent
