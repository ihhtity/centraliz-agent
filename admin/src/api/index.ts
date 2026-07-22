import request from '@/utils/request';
import type { Room, Device, Group, Rule, Order, Merch, DeviceLog, DashboardStats } from '@/types';

export const getDashboardStats = () => {
  return request.get<DashboardStats>('/admin/stats', { custom: { auth: true } });
};

export const login = (data: { account: string; password: string; type: string }) => {
  return request.post<{ token: string }>('/merch/login', data);
};

export const getRoomList = (params: {
  merchs_id?: number;
  groups_id?: number;
  name?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Room[]; total: number }>('/admin/room/list', { params, custom: { auth: true } });
};

export const getRoomDetail = (id: number) => {
  return request.get<Room>(`/admin/room/${id}`, { custom: { auth: true } });
};

export const createRoom = (data: Partial<Room>) => {
  return request.post<Room>('/admin/room', data, { custom: { auth: true } });
};

export const updateRoom = (id: number, data: Partial<Room>) => {
  return request.put<Room>(`/admin/room/${id}`, data, { custom: { auth: true } });
};

export const deleteRoom = (id: number) => {
  return request.delete<void>(`/admin/room/${id}`, { custom: { auth: true } });
};

export const batchDeleteRoom = (data: { ids: string[] }) => {
  return request.post<void>('/admin/room/batch-delete', data, { custom: { auth: true } });
};

export const getDeviceList = (params: {
  merchs_id?: number;
  groups_id?: number;
  name?: string;
  status?: string;
  type?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Device[]; total: number }>('/admin/device/list', { params, custom: { auth: true } });
};

export const getDeviceDetail = (id: number) => {
  return request.get<Device>(`/admin/device/${id}`, { custom: { auth: true } });
};

export const createDevice = (data: Partial<Device>) => {
  return request.post<Device>('/admin/device', data, { custom: { auth: true } });
};

export const updateDevice = (id: number, data: Partial<Device>) => {
  return request.put<Device>(`/admin/device/${id}`, data, { custom: { auth: true } });
};

export const deleteDevice = (id: number) => {
  return request.delete<void>(`/admin/device/${id}`, { custom: { auth: true } });
};

export const batchDeleteDevice = (data: { ids: string[] }) => {
  return request.post<void>('/admin/device/batch-delete', data, { custom: { auth: true } });
};

export const getGroupList = (params: {
  merchs_id?: number;
  name?: string;
  type?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Group[]; total: number }>('/admin/group/list', { params, custom: { auth: true } });
};

export const getGroupDetail = (id: number) => {
  return request.get<Group>(`/admin/group/${id}`, { custom: { auth: true } });
};

export const createGroup = (data: Partial<Group>) => {
  return request.post<Group>('/admin/group', data, { custom: { auth: true } });
};

export const updateGroup = (id: number, data: Partial<Group>) => {
  return request.put<Group>(`/admin/group/${id}`, data, { custom: { auth: true } });
};

export const deleteGroup = (id: number) => {
  return request.delete<void>(`/admin/group/${id}`, { custom: { auth: true } });
};

export const batchDeleteGroup = (data: { ids: string[] }) => {
  return request.post<void>('/admin/group/batch-delete', data, { custom: { auth: true } });
};

export const getRuleList = (params: {
  merchs_id?: number;
  name?: string;
  type?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Rule[]; total: number }>('/admin/rule/list', { params, custom: { auth: true } });
};

export const getRuleDetail = (id: number) => {
  return request.get<Rule>(`/admin/rule/${id}`, { custom: { auth: true } });
};

export const createRule = (data: Partial<Rule>) => {
  return request.post<Rule>('/admin/rule', data, { custom: { auth: true } });
};

export const updateRule = (id: number, data: Partial<Rule>) => {
  return request.put<Rule>(`/admin/rule/${id}`, data, { custom: { auth: true } });
};

export const deleteRule = (id: number) => {
  return request.delete<void>(`/admin/rule/${id}`, { custom: { auth: true } });
};

export const batchDeleteRule = (data: { ids: string[] }) => {
  return request.post<void>('/admin/rule/batch-delete', data, { custom: { auth: true } });
};

export const getOrderList = (params: {
  merchs_id?: number;
  users_id?: number;
  rooms_id?: number;
  status?: string;
  code?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Order[]; total: number }>('/admin/order/list', { params, custom: { auth: true } });
};

export const getOrderDetail = (id: number) => {
  return request.get<Order>(`/admin/order/${id}`, { custom: { auth: true } });
};

export const updateOrder = (id: number, data: Partial<Order>) => {
  return request.put<Order>(`/admin/order/${id}`, data, { custom: { auth: true } });
};

export const deleteOrder = (id: number) => {
  return request.delete<void>(`/admin/order/${id}`, { custom: { auth: true } });
};

export const batchDeleteOrder = (data: { ids: string[] }) => {
  return request.post<void>('/admin/order/batch-delete', data, { custom: { auth: true } });
};

export const getMerchList = (params: {
  account?: string;
  phone?: string;
  role?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Merch[]; total: number }>('/admin/merch/list', { params, custom: { auth: true } });
};

export const getMerchDetail = (id: number) => {
  return request.get<Merch>(`/admin/merch/${id}`, { custom: { auth: true } });
};

export const createMerch = (data: Partial<Merch>) => {
  return request.post<Merch>('/admin/merch', data, { custom: { auth: true } });
};

export const updateMerch = (id: number, data: Partial<Merch>) => {
  return request.put<Merch>(`/admin/merch/${id}`, data, { custom: { auth: true } });
};

export const deleteMerch = (id: number) => {
  return request.delete<void>(`/admin/merch/${id}`, { custom: { auth: true } });
};

export const batchDeleteMerch = (data: { ids: string[] }) => {
  return request.post<void>('/admin/merch/batch-delete', data, { custom: { auth: true } });
};

export const getDeviceLogList = (params: {
  merchs_id?: number;
  devices_id?: number;
  room_id?: number;
  type?: string;
  control?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: DeviceLog[]; total: number }>('/admin/devicelog/list', { params, custom: { auth: true } });
};

export const getDeviceLogDetail = (id: number) => {
  return request.get<DeviceLog>(`/admin/devicelog/${id}`, { custom: { auth: true } });
};

export const batchDeleteDeviceLog = (data: { ids: string[] }) => {
  return request.post<void>('/admin/devicelog/batch-delete', data, { custom: { auth: true } });
};
