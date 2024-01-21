export interface APIMetadata {
  version: string
}

export interface BaseResponse {
  code: number
  msg: string
  body?: any
}
