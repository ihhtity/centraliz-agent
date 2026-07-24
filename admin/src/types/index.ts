// 房间
export interface Room {
  id: number;
  merchsId: number;
  usersId: number;
  groupsId: number;
  rulesId: number;
  devicesId: number;
  ordersId: number;
  name: string;
  tag: string;
  status: string;
  boardNo: string;
  lockNo: string;
  image: string;
  price: number;
  freeTime: string;
  combo: string;
  createdAt: string;
  updatedAt: string;
}
// 设备
export interface Device {
  id: number;
  name: string;
  code: string;
  boardNo: string;
  merchsId: number;
  roomsId: number;
  groupsId: number;
  ordersId: number;
  status: string;
  type: string;
  lockCount: number;
  recharge: string;
  version: string;
  signal: string;
  heat: string;
  protectHeat: string;
  createdAt: string;
  updatedAt: string;
}
// 分组
export interface Group {
  id: number;
  name: string;
  merchsId: number;
  rulesId: number;
  devicesId: number;
  phone: string;
  count: number;
  roomCount: number;
  type: string;
  location: string;
  address: string;
  ruleName: string;
  bindNumber: string;
  consumePush: string;
  sort: number;
  tag: string;
  createdAt: string;
  updatedAt: string;
}
// 规则
export interface Rule {
  id: number;
  name: string;
  merchsId: number;
  type: string;
  mode: string;
  price: number;
  deposit: number;
  rate: number;
  durationUnit: string;
  autoEndTime: number;
  description: string;
  freeTime: number;
  autoRefund: boolean;
  manualRenew: boolean;
  timeOptions: string;
  sort: number;
  tag: string;
  duration: number;
  status: number;
  createdAt: string;
  updatedAt: string;
}
// 订单
export interface Order {
  id: number;
  orderNo: string;
  merchsId: number;
  usersId: number;
  roomsId: number;
  groupsId: number;
  rulesId: number;
  code: string;
  payCode: string;
  type: string;
  mode: string;
  status: string;
  tag: string;
  amount: number;
  duration: number;
  price: number;
  deposit: number;
  payPrice: number;
  payType: string;
  refundPrice: number;
  remark: string;
  userPhone: string;
  merchPhone: string;
  reqSeqId: string;
  reqDate: string;
  freeTime: string;
  startTime: string;
  endTime: string;
  payTime: string;
  refundTime: string;
  createdAt: string;
  updatedAt: string;
}
// 商户
export interface Merch {
  id: number;
  account: string;
  password: string;
  email: string;
  phone: string;
  role: string;
  status: string;
  logAt: string;
  createdAt: string;
}
// 设备日志
export interface DeviceLog {
  id: number;
  merchsId: number;
  usersId: number;
  devicesId: number;
  roomId: number;
  code: string;
  deviceName: string;
  roomName: string;
  type: string;
  control: string;
  status: string;
  occupant: string;
  phone: string;
  createdAt: string;
  updatedAt: string;
}
// 仪表盘统计
export interface DashboardStats {
  roomCount: number;
  deviceCount: number;
  groupCount: number;
  orderCount: number;
}
// 支付账号
export interface HuifuAccount {
  id: number;
  merchsId: number;
  code: string;
  sharing: string;
  account: string;
  phone: string;
  name: string;
  identity: string;
  card: string;
  encrypt: string;
  storename: string;
  area: string;
  picture: string;
  remarks: string;
  type: string;
  choose: string;
  share: string;
  rate: number;
  createdAt: string;
}
// 商户支付
export interface MerchPay {
  id: number;
  code: string;
  merchsId: number;
  name: string;
  reqDate: string;
  hfSeqId: string;
  originalPrice: number;
  price: number;
  locktotal: number;
  type: string;
  status: string;
  remarks: string;
  createdAt: string;
}
// 房间图片
export interface RoomImage {
  id: number;
  name: string;
  image: string;
  createdAt: string;
  updatedAt: string;
}
// 房间标签
export interface RoomTag {
  id: number;
  merchsId: number;
  name: string;
  createdAt: string;
  updatedAt: string;
}
// 子商户
export interface SubMerch {
  id: number;
  merchsId: number;
  account: string;
  password: string;
  email: string;
  phone: string;
  role: string;
  status: string;
  rule: string;
  logAt: string;
  createdAt: string;
}
// 用户
export interface User {
  id: number;
  merchsId: number;
  roomsId: number | null;
  ordersId: number | null;
  name: string;
  account: string;
  password: string;
  email: string | null;
  phone: string | null;
  privacy: string;
  status: string | null;
  avatarURL: string | null;
  unionId: string;
  openid: string;
  gopenid: string;
  updatedAt: string | null;
  createdAt: string | null;
}
// 微信用户
export interface WechatUser {
  id: number;
  openId: string;
  gopenId: string;
  unionId: string;
  sessionKey: string;
  accessToken: string;
  refreshToken: string;
  nickname: string;
  avatar: string;
  gender: number;
  country: string;
  province: string;
  city: string;
  language: string;
  platform: string;
  merchsId: number;
  usersId: number;
  status: string;
  expiredAt: string;
  createdAt: string;
  updatedAt: string;
}
// 工具调用
export interface ToolCall {
  tool_name: string;
  parameters: Record<string, any>;
}
// 工具执行结果
export interface ToolResult {
  tool_name: string;
  success: boolean;
  output: string;
  error?: string;
}
// 聊天消息
export interface ChatMessage {
  id: string;
  role: 'user' | 'assistant';
  content: string;
  thought?: string;
  tool_calls?: ToolCall[];
  tool_results?: ToolResult[];
  created_at: string;
}
// 聊天会话
export interface ChatSession {
  id: string;
  user_id: number;
  user_name: string;
  title: string;
  context?: string;
  messages: ChatMessage[];
  requires_confirm: boolean;
  affected_files: string[];
  created_at: string;
  updated_at: string;
}
// 聊天请求
export interface ChatRequest {
  session_id?: string;
  message: string;
}
// 确认请求
export interface ConfirmRequest {
  session_id: string;
  confirm: boolean;
  message_id?: string;
}