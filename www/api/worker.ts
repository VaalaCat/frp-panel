import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
  CreateWorkerIngressRequest,
  CreateWorkerIngressResponse,
  CreateWorkerRequest,
  CreateWorkerResponse,
  GetWorkerIngressRequest,
  GetWorkerIngressResponse,
  GetWorkerRequest,
  GetWorkerResponse,
  GetWorkerStatusRequest,
  GetWorkerStatusResponse,
  InstallWorkerdRequest,
  InstallWorkerdResponse,
  ListWorkersRequest,
  ListWorkersResponse,
  RemoveWorkerRequest,
  RemoveWorkerResponse,
  UpdateWorkerRequest,
  UpdateWorkerResponse,
} from '@/lib/pb/api_client'
import { BaseResponse } from '@/types/api'
import { constants } from 'node:buffer'

export const getWorker = async (req: GetWorkerRequest) => {
  const res = await http.post(API_PATH + '/worker/get', GetWorkerRequest.toJson(req))
  return GetWorkerResponse.fromJson((res.data as BaseResponse).body)
}

export const createWorker = async (req: CreateWorkerRequest) => {
  const res = await http.post(API_PATH + '/worker/create', CreateWorkerRequest.toJson(req))
  return CreateWorkerResponse.fromJson((res.data as BaseResponse).body)
}

export const updateWorker = async (req: UpdateWorkerRequest) => {
  const res = await http.post(API_PATH + '/worker/update', UpdateWorkerRequest.toJson(req))
  return UpdateWorkerResponse.fromJson((res.data as BaseResponse).body)
}

export const removeWorker = async (req: RemoveWorkerRequest) => {
  const res = await http.post(API_PATH + '/worker/remove', RemoveWorkerRequest.toJson(req))
  return RemoveWorkerResponse.fromJson((res.data as BaseResponse).body)
}

export const listWorkers = async (req: ListWorkersRequest) => {
  const res = await http.post(API_PATH + '/worker/list', ListWorkersRequest.toJson(req))
  return ListWorkersResponse.fromJson((res.data as BaseResponse).body)
}

export const createWorkerIngress = async (req: CreateWorkerIngressRequest) => {
  const res = await http.post(API_PATH + '/worker/create_ingress', CreateWorkerIngressRequest.toJson(req))
  return CreateWorkerIngressResponse.fromJson((res.data as BaseResponse).body)
}

export const getWorkerIngress = async (req: GetWorkerIngressRequest) => {
  const res = await http.post(API_PATH + '/worker/get_ingress', GetWorkerIngressRequest.toJson(req))
  return GetWorkerIngressResponse.fromJson((res.data as BaseResponse).body)
}

export const getWorkerStatus = async (req: GetWorkerStatusRequest) => {
  const res = await http.post(API_PATH + '/worker/status', GetWorkerStatusRequest.toJson(req))
  return GetWorkerStatusResponse.fromJson((res.data as BaseResponse).body)
}

export const installWorkerd = async (req: InstallWorkerdRequest) => {
  const res = await http.post(API_PATH + '/client/install_workerd', InstallWorkerdRequest.toJson(req))
  return InstallWorkerdResponse.fromJson((res.data as BaseResponse).body)
}
