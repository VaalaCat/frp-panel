import { LOCAL_STORAGE_TOKEN_KEY, SET_TOKEN_HEADER, X_CLIENT_REQUEST_ID } from '@/lib/consts'
import { $token } from '@/store/user'
import axios from 'axios'
import { v4 as uuidv4 } from 'uuid'

const instance = axios.create({})

instance.interceptors.request.use((request) => {
  let token = 'Bearer ' + localStorage.getItem(LOCAL_STORAGE_TOKEN_KEY)
  if (token) {
    request.headers.Authorization = token
    $token.set(token)
  }
  request.headers[X_CLIENT_REQUEST_ID] = uuidv4()
  return request
})

instance.interceptors.response.use((response) => {
  if (response.headers?.[SET_TOKEN_HEADER]) {
    localStorage.setItem(LOCAL_STORAGE_TOKEN_KEY, response.headers[SET_TOKEN_HEADER])
    $token.set(response.headers[SET_TOKEN_HEADER])
  }
  if (response.data.code != 200) {
    throw response.data.msg
  }
  return response
})

export default instance
