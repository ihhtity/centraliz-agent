import request from '@/utils/request';
import type { Room, Device, Group, Rule, Order, Merch, DeviceLog, DashboardStats, HuifuAccount, MerchPay, RoomImage, RoomTag, SubMerch, User, WechatUser } from '@/types';

export const getDashboardStats = () => {
  return request.get<DashboardStats>('/admin/stats', { custom: { auth: true } });
};

export const getTrendStats = (params: { type: string; days?: number }) => {
  return request.get<{ data: { date: string; value: number }[] }>('/admin/stats/trend', { params, custom: { auth: true } });
};

export const login = (data: { account: string; password: string; type: string }) => {
  return request.post<{ token: string }>('/merch/login', data);
};

export const getRoomList = (params: {
  merchs_id?: number;
  groups_id?: number;
  name?: string;
  status?: string;
  board_no?: string;
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

export const batchUpdateRoom = (data: { ids: string[]; data: Partial<Room> }) => {
  return request.post<void>('/admin/room/batch-update', data, { custom: { auth: true } });
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

export const batchUpdateDevice = (data: { ids: string[]; data: Partial<Device> }) => {
  return request.post<void>('/admin/device/batch-update', data, { custom: { auth: true } });
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

export const batchUpdateGroup = (data: { ids: string[]; data: Partial<Group> }) => {
  return request.post<void>('/admin/group/batch-update', data, { custom: { auth: true } });
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

export const batchUpdateRule = (data: { ids: string[]; data: Partial<Rule> }) => {
  return request.post<void>('/admin/rule/batch-update', data, { custom: { auth: true } });
};

export const getOrderList = (params: {
  merchs_id?: number;
  users_id?: number;
  rooms_id?: number;
  status?: string;
  code?: string;
  order_no?: string;
  user_phone?: string;
  pay_type?: string;
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

export const batchUpdateOrder = (data: { ids: string[]; data: Partial<Order> }) => {
  return request.post<void>('/admin/order/batch-update', data, { custom: { auth: true } });
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

export const batchUpdateMerch = (data: { ids: string[]; data: Partial<Merch> }) => {
  return request.post<void>('/admin/merch/batch-update', data, { custom: { auth: true } });
};

export const getDeviceLogList = (params: {
  merchs_id?: number;
  devices_id?: number;
  device_name?: string;
  room_id?: number;
  type?: string;
  control?: string;
  status?: string;
  code?: string;
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

export const importDeviceLog = (data: FormData) => {
  return request.post<void>('/admin/devicelog/import', data, { custom: { auth: true } });
};

export const getHuifuAccountList = (params: {
  merchs_id?: number;
  code?: string;
  account?: string;
  phone?: string;
  type?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: HuifuAccount[]; total: number }>('/admin/huifu/list', { params, custom: { auth: true } });
};

export const getHuifuAccountDetail = (id: number) => {
  return request.get<HuifuAccount>(`/admin/huifu/${id}`, { custom: { auth: true } });
};

export const createHuifuAccount = (data: Partial<HuifuAccount>) => {
  return request.post<HuifuAccount>('/admin/huifu', data, { custom: { auth: true } });
};

export const updateHuifuAccount = (id: number, data: Partial<HuifuAccount>) => {
  return request.put<HuifuAccount>(`/admin/huifu/${id}`, data, { custom: { auth: true } });
};

export const deleteHuifuAccount = (id: number) => {
  return request.delete<void>(`/admin/huifu/${id}`, { custom: { auth: true } });
};

export const batchDeleteHuifuAccount = (data: { ids: string[] }) => {
  return request.post<void>('/admin/huifu/batch-delete', data, { custom: { auth: true } });
};

export const batchUpdateHuifuAccount = (data: { ids: string[]; data: Partial<HuifuAccount> }) => {
  return request.post<void>('/admin/huifu/batch-update', data, { custom: { auth: true } });
};

export const getMerchPayList = (params: {
  merchs_id?: number;
  code?: string;
  name?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: MerchPay[]; total: number }>('/admin/merchpay/list', { params, custom: { auth: true } });
};

export const getMerchPayDetail = (id: number) => {
  return request.get<MerchPay>(`/admin/merchpay/${id}`, { custom: { auth: true } });
};

export const batchDeleteMerchPay = (data: { ids: string[] }) => {
  return request.post<void>('/admin/merchpay/batch-delete', data, { custom: { auth: true } });
};

export const importMerchPay = (data: FormData) => {
  return request.post<void>('/admin/merchpay/import', data, { custom: { auth: true } });
};

export const getRoomImageList = (params: {
  name?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: RoomImage[]; total: number }>('/admin/roomimg/list', { params, custom: { auth: true } });
};

export const getRoomImageDetail = (id: number) => {
  return request.get<RoomImage>(`/admin/roomimg/${id}`, { custom: { auth: true } });
};

export const createRoomImage = (data: Partial<RoomImage>) => {
  return request.post<RoomImage>('/admin/roomimg', data, { custom: { auth: true } });
};

export const updateRoomImage = (id: number, data: Partial<RoomImage>) => {
  return request.put<RoomImage>(`/admin/roomimg/${id}`, data, { custom: { auth: true } });
};

export const deleteRoomImage = (id: number) => {
  return request.delete<void>(`/admin/roomimg/${id}`, { custom: { auth: true } });
};

export const batchDeleteRoomImage = (data: { ids: string[] }) => {
  return request.post<void>('/admin/roomimg/batch-delete', data, { custom: { auth: true } });
};

export const batchUpdateRoomImage = (data: { ids: string[]; data: Partial<RoomImage> }) => {
  return request.post<void>('/admin/roomimg/batch-update', data, { custom: { auth: true } });
};

export const getRoomTagList = (params: {
  merchs_id?: number;
  name?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: RoomTag[]; total: number }>('/admin/roomtag/list', { params, custom: { auth: true } });
};

export const getRoomTagDetail = (id: number) => {
  return request.get<RoomTag>(`/admin/roomtag/${id}`, { custom: { auth: true } });
};

export const createRoomTag = (data: Partial<RoomTag>) => {
  return request.post<RoomTag>('/admin/roomtag', data, { custom: { auth: true } });
};

export const updateRoomTag = (id: number, data: Partial<RoomTag>) => {
  return request.put<RoomTag>(`/admin/roomtag/${id}`, data, { custom: { auth: true } });
};

export const deleteRoomTag = (id: number) => {
  return request.delete<void>(`/admin/roomtag/${id}`, { custom: { auth: true } });
};

export const batchDeleteRoomTag = (data: { ids: string[] }) => {
  return request.post<void>('/admin/roomtag/batch-delete', data, { custom: { auth: true } });
};

export const batchUpdateRoomTag = (data: { ids: string[]; data: Partial<RoomTag> }) => {
  return request.post<void>('/admin/roomtag/batch-update', data, { custom: { auth: true } });
};

export const getSubMerchList = (params: {
  merchs_id?: number;
  account?: string;
  phone?: string;
  role?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: SubMerch[]; total: number }>('/admin/submerch/list', { params, custom: { auth: true } });
};

export const getSubMerchDetail = (id: number) => {
  return request.get<SubMerch>(`/admin/submerch/${id}`, { custom: { auth: true } });
};

export const createSubMerch = (data: Partial<SubMerch>) => {
  return request.post<SubMerch>('/admin/submerch', data, { custom: { auth: true } });
};

export const updateSubMerch = (id: number, data: Partial<SubMerch>) => {
  return request.put<SubMerch>(`/admin/submerch/${id}`, data, { custom: { auth: true } });
};

export const deleteSubMerch = (id: number) => {
  return request.delete<void>(`/admin/submerch/${id}`, { custom: { auth: true } });
};

export const batchDeleteSubMerch = (data: { ids: string[] }) => {
  return request.post<void>('/admin/submerch/batch-delete', data, { custom: { auth: true } });
};

export const batchUpdateSubMerch = (data: { ids: string[]; data: Partial<SubMerch> }) => {
  return request.post<void>('/admin/submerch/batch-update', data, { custom: { auth: true } });
};

export const getUserList = (params: {
  merchs_id?: number;
  name?: string;
  account?: string;
  phone?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: User[]; total: number }>('/admin/user/list', { params, custom: { auth: true } });
};

export const getUserDetail = (id: number) => {
  return request.get<User>(`/admin/user/${id}`, { custom: { auth: true } });
};

export const updateUser = (id: number, data: Partial<User>) => {
  return request.put<User>(`/admin/user/${id}`, data, { custom: { auth: true } });
};

export const deleteUser = (id: number) => {
  return request.delete<void>(`/admin/user/${id}`, { custom: { auth: true } });
};

export const batchDeleteUser = (data: { ids: string[] }) => {
  return request.post<void>('/admin/user/batch-delete', data, { custom: { auth: true } });
};

export const batchUpdateUser = (data: { ids: string[]; data: Partial<User> }) => {
  return request.post<void>('/admin/user/batch-update', data, { custom: { auth: true } });
};

export const importUser = (data: FormData) => {
  return request.post<void>('/admin/user/import', data, { custom: { auth: true } });
};

export const getWxUserList = (params: {
  merchs_id?: number;
  nickname?: string;
  openid?: string;
  platform?: string;
  status?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: WechatUser[]; total: number }>('/admin/wxuser/list', { params, custom: { auth: true } });
};

export const getWxUserDetail = (id: number) => {
  return request.get<WechatUser>(`/admin/wxuser/${id}`, { custom: { auth: true } });
};

export const updateWxUser = (id: number, data: Partial<WechatUser>) => {
  return request.put<WechatUser>(`/admin/wxuser/${id}`, data, { custom: { auth: true } });
};

export const batchDeleteWxUser = (data: { ids: string[] }) => {
  return request.post<void>('/admin/wxuser/batch-delete', data, { custom: { auth: true } });
};

export const importWxUser = (data: FormData) => {
  return request.post<void>('/admin/wxuser/import', data, { custom: { auth: true } });
};
