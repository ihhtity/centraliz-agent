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

export interface DashboardStats {
  roomCount: number;
  deviceCount: number;
  groupCount: number;
  orderCount: number;
}