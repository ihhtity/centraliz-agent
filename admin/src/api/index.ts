import request from '@/utils/request';
import type { Room, Device, Group, Rule, Order, Merch, DeviceLog, DashboardStats, HuifuAccount, MerchPay, RoomImage, RoomTag, SubMerch, User, WechatUser } from '@/types';
// 获取看板统计信息
export const getDashboardStats = () => {
  return request.get<DashboardStats>('/admin/stats', { custom: { auth: true } });
};
// 获取趋势统计信息
export const getTrendStats = (params: { type: string; days?: number }) => {
  return request.get<{ data: { date: string; value: number }[] }>('/admin/stats/trend', { params, custom: { auth: true } });
};
// 登录
export const login = (data: { account: string; password: string; type: string }) => {
  return request.post<{ token: string }>('/merch/login', data);
};
// 获取房间列表
export const getRoomList = (params: {
  merchs_id?: number;
  groups_id?: number;
  name?: string;
  status?: string;
  board_no?: string;
  lock_no?: string;
  tag?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Room[]; total: number }>('/admin/room/list', { params, custom: { auth: true } });
};
// 获取房间详情
export const getRoomDetail = (id: number) => {
  return request.get<Room>(`/admin/room/${id}`, { custom: { auth: true } });
};
// 创建房间
export const createRoom = (data: Partial<Room>) => {
  return request.post<Room>('/admin/room', data, { custom: { auth: true } });
};
// 更新房间
export const updateRoom = (id: number, data: Partial<Room>) => {
  return request.put<Room>(`/admin/room/${id}`, data, { custom: { auth: true } });
};
// 删除房间
export const deleteRoom = (id: number) => {
  return request.delete<void>(`/admin/room/${id}`, { custom: { auth: true } });
};
// 批量删除房间
export const batchDeleteRoom = (data: { ids: string[] }) => {
  return request.post<void>('/admin/room/batch-delete', data, { custom: { auth: true } });
};
// 批量更新房间
export const batchUpdateRoom = (data: { ids: string[]; data: Partial<Room> }) => {
  return request.post<void>('/admin/room/batch-update', data, { custom: { auth: true } });
};
// 获取设备列表
export const getDeviceList = (params: {
  merchs_id?: number;
  groups_id?: number;
  name?: string;
  status?: string;
  type?: string;
  version?: string;
  signal?: string;
  heat?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Device[]; total: number }>('/admin/device/list', { params, custom: { auth: true } });
};
// 获取设备详情
export const getDeviceDetail = (id: number) => {
  return request.get<Device>(`/admin/device/${id}`, { custom: { auth: true } });
};
// 创建设备
export const createDevice = (data: Partial<Device>) => {
  return request.post<Device>('/admin/device', data, { custom: { auth: true } });
};
// 更新设备
export const updateDevice = (id: number, data: Partial<Device>) => {
  return request.put<Device>(`/admin/device/${id}`, data, { custom: { auth: true } });
};
// 删除设备
export const deleteDevice = (id: number) => {
  return request.delete<void>(`/admin/device/${id}`, { custom: { auth: true } });
};
// 批量删除设备
export const batchDeleteDevice = (data: { ids: string[] }) => {
  return request.post<void>('/admin/device/batch-delete', data, { custom: { auth: true } });
};
// 批量更新设备
export const batchUpdateDevice = (data: { ids: string[]; data: Partial<Device> }) => {
  return request.post<void>('/admin/device/batch-update', data, { custom: { auth: true } });
};
// 获取分组列表
export const getGroupList = (params: {
  merchs_id?: number;
  name?: string;
  type?: string;
  phone?: string;
  location?: string;
  address?: string;
  bind_number?: string;
  consume_push?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Group[]; total: number }>('/admin/group/list', { params, custom: { auth: true } });
};
// 获取分组详情
export const getGroupDetail = (id: number) => {
  return request.get<Group>(`/admin/group/${id}`, { custom: { auth: true } });
};
// 创建分组
export const createGroup = (data: Partial<Group>) => {
  return request.post<Group>('/admin/group', data, { custom: { auth: true } });
};
// 更新分组
export const updateGroup = (id: number, data: Partial<Group>) => {
  return request.put<Group>(`/admin/group/${id}`, data, { custom: { auth: true } });
};
// 删除分组
export const deleteGroup = (id: number) => {
  return request.delete<void>(`/admin/group/${id}`, { custom: { auth: true } });
};
// 批量删除分组
export const batchDeleteGroup = (data: { ids: string[] }) => {
  return request.post<void>('/admin/group/batch-delete', data, { custom: { auth: true } });
};
// 批量更新分组
export const batchUpdateGroup = (data: { ids: string[]; data: Partial<Group> }) => {
  return request.post<void>('/admin/group/batch-update', data, { custom: { auth: true } });
};
// 获取规则列表
export const getRuleList = (params: {
  merchs_id?: number;
  name?: string;
  type?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: Rule[]; total: number }>('/admin/rule/list', { params, custom: { auth: true } });
};
// 获取规则详情
export const getRuleDetail = (id: number) => {
  return request.get<Rule>(`/admin/rule/${id}`, { custom: { auth: true } });
};
// 创建规则
export const createRule = (data: Partial<Rule>) => {
  return request.post<Rule>('/admin/rule', data, { custom: { auth: true } });
};
// 更新规则
export const updateRule = (id: number, data: Partial<Rule>) => {
  return request.put<Rule>(`/admin/rule/${id}`, data, { custom: { auth: true } });
};
// 删除规则
export const deleteRule = (id: number) => {
  return request.delete<void>(`/admin/rule/${id}`, { custom: { auth: true } });
};
// 批量删除规则
export const batchDeleteRule = (data: { ids: string[] }) => {
  return request.post<void>('/admin/rule/batch-delete', data, { custom: { auth: true } });
};
// 批量更新规则
export const batchUpdateRule = (data: { ids: string[]; data: Partial<Rule> }) => {
  return request.post<void>('/admin/rule/batch-update', data, { custom: { auth: true } });
};
// 获取订单列表
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
// 获取订单详情
export const getOrderDetail = (id: number) => {
  return request.get<Order>(`/admin/order/${id}`, { custom: { auth: true } });
};
// 创建订单
export const createOrder = (data: Partial<Order>) => {
  return request.post<Order>('/admin/order', data, { custom: { auth: true } });
};
// 更新订单
export const updateOrder = (id: number, data: Partial<Order>) => {
  return request.put<Order>(`/admin/order/${id}`, data, { custom: { auth: true } });
};
// 删除订单
export const deleteOrder = (id: number) => {
  return request.delete<void>(`/admin/order/${id}`, { custom: { auth: true } });
};
// 批量删除订单
export const batchDeleteOrder = (data: { ids: string[] }) => {
  return request.post<void>('/admin/order/batch-delete', data, { custom: { auth: true } });
};
// 批量更新订单
export const batchUpdateOrder = (data: { ids: string[]; data: Partial<Order> }) => {
  return request.post<void>('/admin/order/batch-update', data, { custom: { auth: true } });
};
// 获取商户列表
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
// 获取商户详情
export const getMerchDetail = (id: number) => {
  return request.get<Merch>(`/admin/merch/${id}`, { custom: { auth: true } });
};
// 创建商户
export const createMerch = (data: Partial<Merch>) => {
  return request.post<Merch>('/admin/merch', data, { custom: { auth: true } });
};
// 更新商户
export const updateMerch = (id: number, data: Partial<Merch>) => {
  return request.put<Merch>(`/admin/merch/${id}`, data, { custom: { auth: true } });
};
// 删除商户
export const deleteMerch = (id: number) => {
  return request.delete<void>(`/admin/merch/${id}`, { custom: { auth: true } });
};
// 批量删除商户
export const batchDeleteMerch = (data: { ids: string[] }) => {
  return request.post<void>('/admin/merch/batch-delete', data, { custom: { auth: true } });
};
// 批量更新商户
export const batchUpdateMerch = (data: { ids: string[]; data: Partial<Merch> }) => {
  return request.post<void>('/admin/merch/batch-update', data, { custom: { auth: true } });
};
// 获取设备日志列表
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
// 获取设备日志详情
export const getDeviceLogDetail = (id: number) => {
  return request.get<DeviceLog>(`/admin/devicelog/${id}`, { custom: { auth: true } });
};
// 创建设备日志
export const createDeviceLog = (data: Partial<DeviceLog>) => {
  return request.post<DeviceLog>('/admin/devicelog', data, { custom: { auth: true } });
};
// 更新设备日志
export const updateDeviceLog = (id: number, data: Partial<DeviceLog>) => {
  return request.put<DeviceLog>(`/admin/devicelog/${id}`, data, { custom: { auth: true } });
};
// 批量更新设备日志
export const batchUpdateDeviceLog = (data: { ids: string[]; data: Partial<DeviceLog> }) => {
  return request.post<void>('/admin/devicelog/batch-update', data, { custom: { auth: true } });
};
// 批量删除设备日志
export const batchDeleteDeviceLog = (data: { ids: string[] }) => {
  return request.post<void>('/admin/devicelog/batch-delete', data, { custom: { auth: true } });
};
// 导入设备日志
export const importDeviceLog = (data: FormData) => {
  return request.post<void>('/admin/devicelog/import', data, { custom: { auth: true } });
};
// 获取汇付账号列表
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
// 获取汇付账号详情
export const getHuifuAccountDetail = (id: number) => {
  return request.get<HuifuAccount>(`/admin/huifu/${id}`, { custom: { auth: true } });
};
// 创建汇付账号
export const createHuifuAccount = (data: Partial<HuifuAccount>) => {
  return request.post<HuifuAccount>('/admin/huifu', data, { custom: { auth: true } });
};
// 更新汇付账号
export const updateHuifuAccount = (id: number, data: Partial<HuifuAccount>) => {
  return request.put<HuifuAccount>(`/admin/huifu/${id}`, data, { custom: { auth: true } });
};
// 删除汇付账号
export const deleteHuifuAccount = (id: number) => {
  return request.delete<void>(`/admin/huifu/${id}`, { custom: { auth: true } });
};
// 批量删除汇付账号
export const batchDeleteHuifuAccount = (data: { ids: string[] }) => {
  return request.post<void>('/admin/huifu/batch-delete', data, { custom: { auth: true } });
};
// 批量更新汇付账号
export const batchUpdateHuifuAccount = (data: { ids: string[]; data: Partial<HuifuAccount> }) => {
  return request.post<void>('/admin/huifu/batch-update', data, { custom: { auth: true } });
};
// 获取商户支付列表
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
// 获取商户支付详情
export const getMerchPayDetail = (id: number) => {
  return request.get<MerchPay>(`/admin/merchpay/${id}`, { custom: { auth: true } });
};
// 创建商户支付
export const createMerchPay = (data: Partial<MerchPay>) => {
  return request.post<MerchPay>('/admin/merchpay', data, { custom: { auth: true } });
};
// 更新商户支付
export const updateMerchPay = (id: number, data: Partial<MerchPay>) => {
  return request.put<MerchPay>(`/admin/merchpay/${id}`, data, { custom: { auth: true } });
};
// 批量更新商户支付
export const batchUpdateMerchPay = (data: { ids: string[]; data: Partial<MerchPay> }) => {
  return request.post<void>('/admin/merchpay/batch-update', data, { custom: { auth: true } });
};
// 批量删除商户支付
export const batchDeleteMerchPay = (data: { ids: string[] }) => {
  return request.post<void>('/admin/merchpay/batch-delete', data, { custom: { auth: true } });
};
// 导入商户支付
export const importMerchPay = (data: FormData) => {
  return request.post<void>('/admin/merchpay/import', data, { custom: { auth: true } });
};
// 获取房间图片列表
export const getRoomImageList = (params: {
  name?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: RoomImage[]; total: number }>('/admin/roomimg/list', { params, custom: { auth: true } });
};
// 获取房间图片详情
export const getRoomImageDetail = (id: number) => {
  return request.get<RoomImage>(`/admin/roomimg/${id}`, { custom: { auth: true } });
};
// 创建房间图片
export const createRoomImage = (data: Partial<RoomImage>) => {
  return request.post<RoomImage>('/admin/roomimg', data, { custom: { auth: true } });
};
// 更新房间图片
export const updateRoomImage = (id: number, data: Partial<RoomImage>) => {
  return request.put<RoomImage>(`/admin/roomimg/${id}`, data, { custom: { auth: true } });
};
// 删除房间图片
export const deleteRoomImage = (id: number) => {
  return request.delete<void>(`/admin/roomimg/${id}`, { custom: { auth: true } });
};
// 批量删除房间图片
export const batchDeleteRoomImage = (data: { ids: string[] }) => {
  return request.post<void>('/admin/roomimg/batch-delete', data, { custom: { auth: true } });
};
// 批量更新房间图片
export const batchUpdateRoomImage = (data: { ids: string[]; data: Partial<RoomImage> }) => {
  return request.post<void>('/admin/roomimg/batch-update', data, { custom: { auth: true } });
};
// 获取房间标签列表
export const getRoomTagList = (params: {
  merchs_id?: number;
  name?: string;
  page?: number;
  page_size?: number;
}) => {
  return request.get<{ data: RoomTag[]; total: number }>('/admin/roomtag/list', { params, custom: { auth: true } });
};
// 获取房间标签详情
export const getRoomTagDetail = (id: number) => {
  return request.get<RoomTag>(`/admin/roomtag/${id}`, { custom: { auth: true } });
};
// 创建房间标签
export const createRoomTag = (data: Partial<RoomTag>) => {
  return request.post<RoomTag>('/admin/roomtag', data, { custom: { auth: true } });
};
// 更新房间标签
export const updateRoomTag = (id: number, data: Partial<RoomTag>) => {
  return request.put<RoomTag>(`/admin/roomtag/${id}`, data, { custom: { auth: true } });
};
// 删除房间标签
export const deleteRoomTag = (id: number) => {
  return request.delete<void>(`/admin/roomtag/${id}`, { custom: { auth: true } });
};
// 批量删除房间标签
export const batchDeleteRoomTag = (data: { ids: string[] }) => {
  return request.post<void>('/admin/roomtag/batch-delete', data, { custom: { auth: true } });
};
// 批量更新房间标签
export const batchUpdateRoomTag = (data: { ids: string[]; data: Partial<RoomTag> }) => {
  return request.post<void>('/admin/roomtag/batch-update', data, { custom: { auth: true } });
};
// 获取子商户列表
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
// 获取子商户详情
export const getSubMerchDetail = (id: number) => {
  return request.get<SubMerch>(`/admin/submerch/${id}`, { custom: { auth: true } });
};
// 创建子商户
export const createSubMerch = (data: Partial<SubMerch>) => {
  return request.post<SubMerch>('/admin/submerch', data, { custom: { auth: true } });
};
// 更新子商户
export const updateSubMerch = (id: number, data: Partial<SubMerch>) => {
  return request.put<SubMerch>(`/admin/submerch/${id}`, data, { custom: { auth: true } });
};
// 删除子商户
export const deleteSubMerch = (id: number) => {
  return request.delete<void>(`/admin/submerch/${id}`, { custom: { auth: true } });
};
// 批量删除子商户
export const batchDeleteSubMerch = (data: { ids: string[] }) => {
  return request.post<void>('/admin/submerch/batch-delete', data, { custom: { auth: true } });
};
// 批量更新子商户
export const batchUpdateSubMerch = (data: { ids: string[]; data: Partial<SubMerch> }) => {
  return request.post<void>('/admin/submerch/batch-update', data, { custom: { auth: true } });
};
// 获取用户列表
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
// 获取用户详情
export const getUserDetail = (id: number) => {
  return request.get<User>(`/admin/user/${id}`, { custom: { auth: true } });
};
// 创建用户
export const createUser = (data: Partial<User>) => {
  return request.post<User>('/admin/user', data, { custom: { auth: true } });
};
// 更新用户
export const updateUser = (id: number, data: Partial<User>) => {
  return request.put<User>(`/admin/user/${id}`, data, { custom: { auth: true } });
};
// 删除用户
export const deleteUser = (id: number) => {
  return request.delete<void>(`/admin/user/${id}`, { custom: { auth: true } });
};
// 批量删除用户
export const batchDeleteUser = (data: { ids: string[] }) => {
  return request.post<void>('/admin/user/batch-delete', data, { custom: { auth: true } });
};
// 批量更新用户
export const batchUpdateUser = (data: { ids: string[]; data: Partial<User> }) => {
  return request.post<void>('/admin/user/batch-update', data, { custom: { auth: true } });
};
// 导入用户
export const importUser = (data: FormData) => {
  return request.post<void>('/admin/user/import', data, { custom: { auth: true } });
};
// 获取微信用户列表
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
// 获取微信用户详情
export const getWxUserDetail = (id: number) => {
  return request.get<WechatUser>(`/admin/wxuser/${id}`, { custom: { auth: true } });
};
// 创建微信用户
export const createWxUser = (data: Partial<WechatUser>) => {
  return request.post<WechatUser>('/admin/wxuser', data, { custom: { auth: true } });
};
// 批量更新微信用户
export const batchUpdateWxUser = (data: { ids: string[]; data: Partial<WechatUser> }) => {
  return request.post<void>('/admin/wxuser/batch-update', data, { custom: { auth: true } });
};
// 更新微信用户
export const updateWxUser = (id: number, data: Partial<WechatUser>) => {
  return request.put<WechatUser>(`/admin/wxuser/${id}`, data, { custom: { auth: true } });
};
// 批量删除微信用户
export const batchDeleteWxUser = (data: { ids: string[] }) => {
  return request.post<void>('/admin/wxuser/batch-delete', data, { custom: { auth: true } });
};
// 导入微信用户
export const importWxUser = (data: FormData) => {
  return request.post<void>('/admin/wxuser/import', data, { custom: { auth: true } });
};

// ==================== 智能开发助手 API ====================

import type { ChatSession, ChatRequest, ConfirmRequest } from '@/types';

export const chatWithAssistant = (data: ChatRequest) => {
  return request.post<{ session: ChatSession }>('/admin/assistant/chat', data, { custom: { auth: true } });
};

export const getSessions = () => {
  return request.get<{ sessions: ChatSession[] }>('/admin/assistant/sessions', { custom: { auth: true } });
};

export const createSession = (title?: string) => {
  return request.post<{ session: ChatSession }>('/admin/assistant/sessions', { title }, { custom: { auth: true } });
};

export const getSessionDetail = (id: string) => {
  return request.get<{ session: ChatSession }>(`/admin/assistant/sessions/${id}`, { custom: { auth: true } });
};

export const deleteSession = (id: string) => {
  return request.delete<void>(`/admin/assistant/sessions/${id}`, { custom: { auth: true } });
};

export const confirmAction = (data: ConfirmRequest) => {
  return request.post<{ session: ChatSession }>('/admin/assistant/confirm', data, { custom: { auth: true } });
};
