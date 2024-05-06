import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useState } from "react"

export interface IdInputProps {
  setKeyword: (keyword: string) => void
  keyword: string
  refetchTrigger?: (randStr: string) => void
}

export const IdInput: React.FC<IdInputProps> = ({ setKeyword, keyword, refetchTrigger }) => {
  const [input, setInput] = useState(keyword)

  return <div className="flex flex-1 flex-row gap-2">
    <Input className="max-w-40 h-auto" defaultValue={keyword} onChange={(e) => setInput(e.target.value)}></Input>
    <Button variant="outline" size={'sm'} onClick={() => {
      setKeyword(input)
      refetchTrigger && refetchTrigger(JSON.stringify(Math.random()))
    }}>搜索</Button>
  </div >
}
