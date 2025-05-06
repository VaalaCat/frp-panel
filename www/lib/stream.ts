import { API_PATH } from './consts'

export const parseStreaming = async (
  controller: AbortController,
  id: string,
  pkgs: string[],
  onLog: (value: string) => void,
  onError?: (status: number) => void,
  onDone?: () => void,
) => {
  const decoder = new TextDecoder()
  let uint8Array = new Uint8Array()
  let chunks = ''
  let param: Record<string, string> = {
    id: id,
  }
  if (pkgs.length > 0) {
    param['pkgs'] = pkgs.join(',')
  }

  const response = await fetch(`${API_PATH}/log?${new URLSearchParams(param).toString()}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Accept: '*./*',
    },
    signal: controller.signal,
  })
  if (response.status !== 200) {
    onError?.(response.status)
    return
  } else {
    onError?.(200)
  }
  function decodeLog(chunk: string): string {
    const lines = chunk.split('\n')
    const newLines = lines.map((line) => {
      try {
        return JSON.parse(line)
      } catch (error) {
        return line
      }
    })
    const decodedLines = newLines.map((line) => {
      return Buffer.from(line, 'base64').toString('utf-8')
    })
    const splittedLines = decodedLines
      .map((line) => {
        return line.split('\n')
      })
      .flat()
    const trimmedLines = splittedLines.map((line) => {
      return line.trim().replaceAll('\r', '').replaceAll('\n', '').replaceAll('\t', '')
    })
    return trimmedLines.join('\n')
  }
  fetchStream(
    response,
    (chunk) => {
      //@ts-ignore
      uint8Array = new Uint8Array([...uint8Array, ...chunk])
      chunks = decoder.decode(uint8Array, { stream: true })
      onLog(decodeLog(chunks))
    },
    () => {
      onDone && onDone()
    },
  )
}

async function pump(
  reader: ReadableStreamDefaultReader<Uint8Array>,
  controller: ReadableStreamDefaultController,
  onChunk?: (chunk: Uint8Array) => void,
  onDone?: () => void,
): Promise<ReadableStreamReadResult<Uint8Array> | undefined> {
  const { done, value } = await reader.read()
  if (done) {
    onDone && onDone()
    controller.close()
    return
  }
  onChunk && onChunk(value)
  controller.enqueue(value)
  return pump(reader, controller, onChunk, onDone)
}
export const fetchStream = (
  response: Response,
  onChunk?: (chunk: Uint8Array) => void,
  onDone?: () => void,
): ReadableStream<string> => {
  const reader = response.body!.getReader()
  return new ReadableStream({
    start: (controller) => pump(reader, controller, onChunk, onDone),
  })
}
