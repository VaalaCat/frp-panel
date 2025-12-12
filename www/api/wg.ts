import http from '@/api/http'
import { API_PATH } from '@/lib/consts'
import {
	CreateNetworkRequest,
	CreateNetworkResponse,
	DeleteNetworkRequest,
	DeleteNetworkResponse,
	UpdateNetworkRequest,
	UpdateNetworkResponse,
	GetNetworkRequest,
	GetNetworkResponse,
	ListNetworksRequest,
	ListNetworksResponse,
	CreateEndpointRequest,
	CreateEndpointResponse,
	DeleteEndpointRequest,
	DeleteEndpointResponse,
	UpdateEndpointRequest,
	UpdateEndpointResponse,
	GetEndpointRequest,
	GetEndpointResponse,
	ListEndpointsRequest,
	ListEndpointsResponse,
	CreateWireGuardRequest,
	CreateWireGuardResponse,
	DeleteWireGuardRequest,
	DeleteWireGuardResponse,
	UpdateWireGuardRequest,
	UpdateWireGuardResponse,
	RestartWireGuardRequest,
	RestartWireGuardResponse,
	GetWireGuardRequest,
	GetWireGuardResponse,
	ListWireGuardsRequest,
	ListWireGuardsResponse,
	GetWireGuardRuntimeInfoRequest,
	GetWireGuardRuntimeInfoResponse,
	GetNetworkTopologyRequest,
	GetNetworkTopologyResponse,
	CreateWireGuardLinkRequest,
	CreateWireGuardLinkResponse,
	DeleteWireGuardLinkRequest,
	DeleteWireGuardLinkResponse,
	UpdateWireGuardLinkRequest,
	UpdateWireGuardLinkResponse,
	GetWireGuardLinkRequest,
	GetWireGuardLinkResponse,
	ListWireGuardLinksRequest,
	ListWireGuardLinksResponse,
} from '@/lib/pb/api_wg'
import { BaseResponse } from '@/types/api'

// Network
export const createNetwork = async (req: CreateNetworkRequest) => {
	const res = await http.post(API_PATH + '/wg/network/create', CreateNetworkRequest.toJson(req))
	return CreateNetworkResponse.fromJson((res.data as BaseResponse).body)
}
export const deleteNetwork = async (req: DeleteNetworkRequest) => {
	const res = await http.post(API_PATH + '/wg/network/delete', DeleteNetworkRequest.toJson(req))
	return DeleteNetworkResponse.fromJson((res.data as BaseResponse).body)
}
export const updateNetwork = async (req: UpdateNetworkRequest) => {
	const res = await http.post(API_PATH + '/wg/network/update', UpdateNetworkRequest.toJson(req))
	return UpdateNetworkResponse.fromJson((res.data as BaseResponse).body)
}
export const getNetwork = async (req: GetNetworkRequest) => {
	const res = await http.post(API_PATH + '/wg/network/get', GetNetworkRequest.toJson(req))
	return GetNetworkResponse.fromJson((res.data as BaseResponse).body)
}
export const listNetworks = async (req: ListNetworksRequest) => {
	const res = await http.post(API_PATH + '/wg/network/list', ListNetworksRequest.toJson(req))
	return ListNetworksResponse.fromJson((res.data as BaseResponse).body)
}
export const getNetworkTopology = async (req: GetNetworkTopologyRequest) => {
	const res = await http.post(API_PATH + '/wg/network/topology', GetNetworkTopologyRequest.toJson(req))
	return GetNetworkTopologyResponse.fromJson((res.data as BaseResponse).body)
}

// Endpoint
export const createEndpoint = async (req: CreateEndpointRequest) => {
	const res = await http.post(API_PATH + '/wg/endpoint/create', CreateEndpointRequest.toJson(req))
	return CreateEndpointResponse.fromJson((res.data as BaseResponse).body)
}
export const deleteEndpoint = async (req: DeleteEndpointRequest) => {
	const res = await http.post(API_PATH + '/wg/endpoint/delete', DeleteEndpointRequest.toJson(req))
	return DeleteEndpointResponse.fromJson((res.data as BaseResponse).body)
}
export const updateEndpoint = async (req: UpdateEndpointRequest) => {
	const res = await http.post(API_PATH + '/wg/endpoint/update', UpdateEndpointRequest.toJson(req))
	return UpdateEndpointResponse.fromJson((res.data as BaseResponse).body)
}
export const getEndpoint = async (req: GetEndpointRequest) => {
	const res = await http.post(API_PATH + '/wg/endpoint/get', GetEndpointRequest.toJson(req))
	return GetEndpointResponse.fromJson((res.data as BaseResponse).body)
}
export const listEndpoints = async (req: ListEndpointsRequest) => {
	const res = await http.post(API_PATH + '/wg/endpoint/list', ListEndpointsRequest.toJson(req))
	return ListEndpointsResponse.fromJson((res.data as BaseResponse).body)
}

// WireGuard
export const createWireGuard = async (req: CreateWireGuardRequest) => {
	const res = await http.post(API_PATH + '/wg/create', CreateWireGuardRequest.toJson(req))
	return CreateWireGuardResponse.fromJson((res.data as BaseResponse).body)
}
export const deleteWireGuard = async (req: DeleteWireGuardRequest) => {
	const res = await http.post(API_PATH + '/wg/delete', DeleteWireGuardRequest.toJson(req))
	return DeleteWireGuardResponse.fromJson((res.data as BaseResponse).body)
}
export const restartWireGuard = async (req: RestartWireGuardRequest) => {
	const res = await http.post(API_PATH + '/wg/restart', RestartWireGuardRequest.toJson(req))
	return RestartWireGuardResponse.fromJson((res.data as BaseResponse).body)
}
export const updateWireGuard = async (req: UpdateWireGuardRequest) => {
	const res = await http.post(API_PATH + '/wg/update', UpdateWireGuardRequest.toJson(req))
	return UpdateWireGuardResponse.fromJson((res.data as BaseResponse).body)
}
export const getWireGuard = async (req: GetWireGuardRequest) => {
	const res = await http.post(API_PATH + '/wg/get', GetWireGuardRequest.toJson(req))
	return GetWireGuardResponse.fromJson((res.data as BaseResponse).body)
}
export const listWireGuards = async (req: ListWireGuardsRequest) => {
	const res = await http.post(API_PATH + '/wg/list', ListWireGuardsRequest.toJson(req))
	return ListWireGuardsResponse.fromJson((res.data as BaseResponse).body)
}
export const getWireGuardRuntime = async (req: GetWireGuardRuntimeInfoRequest) => {
	const res = await http.post(API_PATH + '/wg/runtime/get', GetWireGuardRuntimeInfoRequest.toJson(req))
	return GetWireGuardRuntimeInfoResponse.fromJson((res.data as BaseResponse).body)
}

// WireGuard Link
export const createWireGuardLink = async (req: CreateWireGuardLinkRequest) => {
	const res = await http.post(API_PATH + '/wg/link/create', CreateWireGuardLinkRequest.toJson(req))
	return CreateWireGuardLinkResponse.fromJson((res.data as BaseResponse).body)
}
export const deleteWireGuardLink = async (req: DeleteWireGuardLinkRequest) => {
	const res = await http.post(API_PATH + '/wg/link/delete', DeleteWireGuardLinkRequest.toJson(req))
	return DeleteWireGuardLinkResponse.fromJson((res.data as BaseResponse).body)
}
export const updateWireGuardLink = async (req: UpdateWireGuardLinkRequest) => {
	const res = await http.post(API_PATH + '/wg/link/update', UpdateWireGuardLinkRequest.toJson(req))
	return UpdateWireGuardLinkResponse.fromJson((res.data as BaseResponse).body)
}
export const getWireGuardLink = async (req: GetWireGuardLinkRequest) => {
	const res = await http.post(API_PATH + '/wg/link/get', GetWireGuardLinkRequest.toJson(req))
	return GetWireGuardLinkResponse.fromJson((res.data as BaseResponse).body)
}
export const listWireGuardLinks = async (req: ListWireGuardLinksRequest) => {
	const res = await http.post(API_PATH + '/wg/link/list', ListWireGuardLinksRequest.toJson(req))
	return ListWireGuardLinksResponse.fromJson((res.data as BaseResponse).body)
}

