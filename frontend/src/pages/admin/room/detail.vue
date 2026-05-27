<template>
	<view class="container">
		<uv-navbar :title="roomName" :placeholder="true" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container" enable-flex>
			<!-- 状态卡片 -->
			<view class="status-card">
				<view class="status-row">
					<view class="status-item">
						<view class="status-icon network">
							<uv-icon name="photo"></uv-icon>
						</view>
						<view class="status-info">
							<text class="status-label">网络</text>
							<text class="status-value online">{{ roomDetail.networkStatus || '在线' }}</text>
						</view>
					</view>
					<view class="status-item">
						<view class="status-icon signal">
							<uv-icon name="photo"></uv-icon>
						</view>
						<view class="status-info">
							<text class="status-label">信号</text>
							<text class="status-value">{{ roomDetail.signalStrength || 90 }}%</text>
						</view>
					</view>
				</view>
				<view class="status-row">
					<view class="status-item">
						<view class="status-icon power" :class="{ active: roomDetail.powerStatus === '开电' }">
							<uv-icon name="photo"></uv-icon>
						</view>
						<view class="status-info">
							<text class="status-label">电状态</text>
							<text class="status-value" :class="{ powerOn: roomDetail.powerStatus === '开电' }">{{ roomDetail.powerStatus || '关电' }}</text>
						</view>
					</view>
					<view class="status-item">
						<view class="status-icon lock" :class="{ unlocked: roomDetail.lockStatus === '开锁' }">
							<uv-icon name="photo"></uv-icon>
						</view>
						<view class="status-info">
							<text class="status-label">锁状态</text>
							<text class="status-value" :class="{ unlocked: roomDetail.lockStatus === '开锁' }">{{ roomDetail.lockStatus || '闭锁' }}</text>
						</view>
					</view>
				</view>
			</view>

			<!-- 房间基本信息 -->
			<view class="info-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="photo"></uv-icon>
					</view>
					<text class="card-title">房间信息</text>
				</view>
				<view class="info-grid">
					<view class="info-item">
						<text class="info-label">房间编号</text>
						<text class="info-value">{{ roomDetail.name || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">房间标签</text>
						<text class="info-value">{{ roomDetail.tag || '普通柜' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">所属分组</text>
						<text class="info-value">{{ roomDetail.groupName || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">房间状态</text>
						<view class="status-badge" :class="getStatusClass(roomDetail.status)">
							{{ getStatusText(roomDetail.status) }}
						</view>
					</view>
				</view>
			</view>

			<!-- 快捷操作区 -->
			<view class="action-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="photo"></uv-icon>
					</view>
					<text class="card-title">快捷操作</text>
				</view>
				<view class="action-grid">
					<view class="action-item" @click="handleOpenLock">
						<view class="action-icon purple">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">远程开锁</text>
					</view>
					<view class="action-item" @click="handlePowerOn">
						<view class="action-icon orange">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">开启电源</text>
					</view>
					<view class="action-item" @click="handlePowerAlways">
						<view class="action-icon yellow">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">常开电源</text>
					</view>
					<view class="action-item" @click="handlePowerOff">
						<view class="action-icon blue">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">关闭电源</text>
					</view>
					<view class="action-item" @click="handleVoice">
						<view class="action-icon blue-light">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">语音播放</text>
					</view>
					<view class="action-item" @click="handleStatusQuery">
						<view class="action-icon gray">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">状态查询</text>
					</view>
					<view class="action-item" @click="handleOperationLog">
						<view class="action-icon purple-light">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">操作记录</text>
					</view>
					<view class="action-item" @click="handleInterruptOrder">
						<view class="action-icon purple-dark">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">中断订单</text>
					</view>
					<view class="action-item" @click="handleReservation">
						<view class="action-icon orange-light">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">预约设置</text>
					</view>
					<view class="action-item" @click="handleEditName">
						<view class="action-icon green">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">修改名称</text>
					</view>
					<view class="action-item" @click="handleOneCode">
						<view class="action-icon orange-dark">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">一码一柜</text>
					</view>
					<view class="action-item" @click="handleImportQrcode">
						<view class="action-icon black">
							<uv-icon name="photo"></uv-icon>
						</view>
						<text class="action-text">导入二维码</text>
					</view>
				</view>
			</view>

			<!-- 设备列表 -->
			<view class="devices-card" v-if="devices.length > 0">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="photo"></uv-icon>
					</view>
					<text class="card-title">已绑定设备</text>
					<text class="device-count-text">{{ devices.length }}台</text>
				</view>
				<view class="devices-list">
					<view v-for="(device, index) in devices" :key="index" class="device-item">
						<view class="device-info">
							<text class="device-name">{{ device.name }}</text>
							<text class="device-code">{{ device.code }}</text>
						</view>
						<view class="device-status" :class="{ online: device.status === '在线' }">
							{{ device.status }}
						</view>
						<view class="device-action" @click="handleUnbindDevice(device.id)">
							<text class="unbind-text">解绑</text>
						</view>
					</view>
				</view>
			</view>

			<!-- 绑定设备入口 -->
			<view class="bind-device-card" v-if="(!roomDetail.deviceCount || roomDetail.deviceCount === 0)" @click="handleBindDevice">
				<view class="bind-icon">
					<uv-icon name="plus" color="#3c9cff" size="32"></uv-icon>
				</view>
				<text class="bind-text">绑定设备</text>
			</view>

			<!-- 解绑设备 -->
			<view class="danger-card" v-if="devices.length > 0" @click="handleUnbindAll">
				<view class="danger-icon">
					<uv-icon name="photo"></uv-icon>
				</view>
				<text class="danger-text">解绑所有设备</text>
			</view>
		</scroll-view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const roomId = ref('')
const roomName = ref('')
const roomDetail = ref({})

onLoad((options) => {
	if (options.id) {
		roomId.value = options.id
	}
	if (options.name) {
		roomName.value = decodeURIComponent(options.name)
	}
	loadRoomDetail()
})

const devices = ref([])

const loadRoomDetail = async () => {
	uni.showLoading({ title: '加载中...' })
	try {
		const token = uni.getStorageSync('token')
		const response = await uni.request({
			url: '/api/v1/room/' + roomId.value,
			method: 'GET',
			header: {
				'Authorization': 'Bearer ' + token,
				'Content-Type': 'application/json'
			}
		})
		
		if (response.data && response.data.code === 200) {
			roomDetail.value = response.data.data
			devices.value = response.data.data.devices || []
		}
	} catch (e) {
		console.error('加载房间详情失败', e)
		// 使用模拟数据
		roomDetail.value = {
			id: roomId.value,
			name: roomName.value,
			tag: '普通柜',
			status: '0',
			groupName: '默认分组',
			networkStatus: '在线',
			signalStrength: 90,
			powerStatus: '关电',
			lockStatus: '开锁',
			deviceCount: 0
		}
		devices.value = []
	} finally {
		uni.hideLoading()
	}
}

const goBack = () => {
	uni.navigateBack()
}

const getStatusText = (status) => {
	if (!status) return '未知'
	if (status === '空闲' || status === 0 || status === '0') return '空闲'
	if (status === '租用' || status === 1 || status === '1') return '租用中'
	if (status === '维修' || status === 2 || status === '2') return '维修中'
	return status
}

const getStatusClass = (status) => {
	if (!status) return 'status-unknown'
	if (status === '空闲' || status === 0 || status === '0') return 'status-free'
	if (status === '租用' || status === 1 || status === '1') return 'status-rent'
	if (status === '维修' || status === 2 || status === '2') return 'status-maint'
	return 'status-unknown'
}

const handleOpenLock = () => {
	uni.showModal({
		title: '确认开锁',
		content: '确定要远程开启此房间的锁吗？',
		success: (res) => {
			if (res.confirm) {
				uni.showToast({ title: '开锁指令已发送', icon: 'success' })
				roomDetail.value.lockStatus = '开锁'
			}
		}
	})
}

const handlePowerOn = () => {
	roomDetail.value.powerStatus = '开电'
	uni.showToast({ title: '电源已开启', icon: 'success' })
}

const handlePowerAlways = () => {
	uni.showToast({ title: '已设置为常开模式', icon: 'success' })
}

const handlePowerOff = () => {
	roomDetail.value.powerStatus = '关电'
	uni.showToast({ title: '电源已关闭', icon: 'success' })
}

const handleVoice = () => {
	uni.showToast({ title: '语音播放功能', icon: 'none' })
}

const handleStatusQuery = () => {
	loadRoomDetail()
}

const handleOperationLog = () => {
	uni.showToast({ title: '操作记录功能', icon: 'none' })
}

const handleInterruptOrder = () => {
	uni.showModal({
		title: '中断订单',
		content: '确定要中断当前订单吗？',
		success: (res) => {
			if (res.confirm) {
				uni.showToast({ title: '订单已中断', icon: 'success' })
			}
		}
	})
}

const handleReservation = () => {
	uni.showToast({ title: '预约设置功能', icon: 'none' })
}

const handleEditName = () => {
	uni.showModal({
		title: '修改名称',
		placeholderText: '请输入新的房间名称',
		success: (res) => {
			if (res.confirm && res.content) {
				roomDetail.value.name = res.content
				roomName.value = res.content
				uni.showToast({ title: '名称已修改', icon: 'success' })
			}
		}
	})
}

const handleOneCode = () => {
	uni.showToast({ title: '一码一柜功能', icon: 'none' })
}

const handleImportQrcode = () => {
	uni.showToast({ title: '导入二维码功能', icon: 'none' })
}

const handleBindDevice = () => {
	uni.showToast({ title: '绑定设备功能', icon: 'none' })
}

const handleUnbindDevice = (deviceId) => {
	uni.showModal({
		title: '解绑设备',
		content: '确定要解绑此设备吗？解绑后将无法控制该设备。',
		confirmColor: '#ff4d4f',
		success: async (res) => {
			if (res.confirm) {
				try {
					const token = uni.getStorageSync('token')
					const response = await uni.request({
						url: '/api/v1/room/unbind-device',
						method: 'POST',
						header: {
							'Authorization': 'Bearer ' + token,
							'Content-Type': 'application/json'
						},
						data: {
							roomId: roomId.value,
							deviceId: deviceId
						}
					})
					
					if (response.data && response.data.code === 200) {
						uni.showToast({ title: '设备已解绑', icon: 'success' })
						loadRoomDetail()
					} else {
						uni.showToast({ title: response.data.msg || '解绑失败', icon: 'none' })
					}
				} catch (e) {
					console.error('解绑设备失败', e)
					uni.showToast({ title: '解绑失败，请重试', icon: 'none' })
				}
			}
		}
	})
}

const handleUnbindAll = () => {
	uni.showModal({
		title: '解绑所有设备',
		content: `确定要解绑所有${devices.value.length}台设备吗？解绑后将无法控制这些设备。`,
		confirmColor: '#ff4d4f',
		success: async (res) => {
			if (res.confirm) {
				try {
					// 批量解绑所有设备
					for (const device of devices.value) {
						const token = uni.getStorageSync('token')
						await uni.request({
							url: '/api/v1/room/unbind-device',
							method: 'POST',
							header: {
								'Authorization': 'Bearer ' + token,
								'Content-Type': 'application/json'
							},
							data: {
								roomId: roomId.value,
								deviceId: device.id
							}
						})
					}
					uni.showToast({ title: '所有设备已解绑', icon: 'success' })
					loadRoomDetail()
				} catch (e) {
					console.error('解绑设备失败', e)
					uni.showToast({ title: '解绑失败，请重试', icon: 'none' })
				}
			}
		}
	})
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.scroll-container {
	height: calc(100vh - 88rpx);
	width: 93%;
	padding: 24rpx;
	padding-bottom: 40rpx;
}

.status-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 28rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	margin-bottom: 20rpx;

	.status-row {
		display: flex;
		justify-content: space-between;

		&:first-child {
			margin-bottom: 20rpx;
			padding-bottom: 20rpx;
			border-bottom: 1rpx solid #f5f5f5;
		}
	}

	.status-item {
		flex: 1;
		display: flex;
		align-items: center;

		&:first-child {
			margin-right: 20rpx;
		}
	}

	.status-icon {
		width: 64rpx;
		height: 64rpx;
		border-radius: 16rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 14rpx;

		&.network {
			background: rgba(24, 144, 255, 0.08);
		}

		&.signal {
			background: rgba(82, 196, 26, 0.08);
		}

		&.power {
			background: rgba(153, 153, 153, 0.08);

			&.active {
				background: rgba(255, 107, 107, 0.08);
			}
		}

		&.lock {
			background: rgba(153, 153, 153, 0.08);

			&.unlocked {
				background: rgba(255, 169, 64, 0.08);
			}
		}
	}

	.status-info {
		display: flex;
		flex-direction: column;

		.status-label {
			font-size: 24rpx;
			color: #999;
			margin-bottom: 4rpx;
		}

		.status-value {
			font-size: 28rpx;
			font-weight: 600;
			color: #333;

			&.online {
				color: #1890ff;
			}

			&.powerOn {
				color: #ff6b6b;
			}

			&.unlocked {
				color: #ffa940;
			}
		}
	}
}

.info-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 28rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	margin-bottom: 20rpx;

	.card-header {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;

		.header-icon {
			width: 44rpx;
			height: 44rpx;
			border-radius: 12rpx;
			background: rgba(41, 121, 255, 0.08);
			display: flex;
			align-items: center;
			justify-content: center;
		}

		.card-title {
			font-size: 30rpx;
			font-weight: 600;
			color: #1a1a1a;
			margin-left: 12rpx;
		}
	}

	.info-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 18rpx;
	}

	.info-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 12rpx 0;

		.info-label {
			font-size: 26rpx;
			color: #8c8c8c;
		}

		.info-value {
			font-size: 26rpx;
			color: #333;
			font-weight: 500;
		}

		.status-badge {
			font-size: 22rpx;
			padding: 6rpx 14rpx;
			border-radius: 16rpx;
			font-weight: 500;

			&.status-free {
				background: rgba(24, 144, 255, 0.1);
				color: #1890ff;
			}

			&.status-rent {
				background: rgba(255, 77, 79, 0.1);
				color: #ff4d4f;
			}

			&.status-maint {
				background: rgba(250, 173, 20, 0.1);
				color: #faad14;
			}

			&.status-unknown {
				background: rgba(153, 153, 153, 0.1);
				color: #999;
			}
		}
	}
}

.action-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 28rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	margin-bottom: 20rpx;

	.card-header {
		display: flex;
		align-items: center;
		margin-bottom: 20rpx;

		.header-icon {
			width: 44rpx;
			height: 44rpx;
			border-radius: 12rpx;
			background: rgba(250, 173, 20, 0.08);
			display: flex;
			align-items: center;
			justify-content: center;
		}

		.card-title {
			font-size: 30rpx;
			font-weight: 600;
			color: #1a1a1a;
			margin-left: 12rpx;
		}
	}

	.action-grid {
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 16rpx;
	}

	.action-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 16rpx 0;
		transition: all 0.2s ease;
		border-radius: 16rpx;

		&:active {
			background: #f5f5f5;
			transform: scale(0.96);
		}

		.action-icon {
			width: 80rpx;
			height: 80rpx;
			border-radius: 20rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-bottom: 10rpx;

			&.purple {
				background: linear-gradient(135deg, #9254de, #722ed1);
			}

			&.orange {
				background: linear-gradient(135deg, #ff7a45, #fa541c);
			}

			&.yellow {
				background: linear-gradient(135deg, #ffc53d, #fadb14);
			}

			&.blue {
				background: linear-gradient(135deg, #597ef7, #2f54eb);
			}

			&.blue-light {
				background: linear-gradient(135deg, #40a9ff, #1890ff);
			}

			&.gray {
				background: linear-gradient(135deg, #8c8c8c, #595959);
			}

			&.purple-light {
				background: linear-gradient(135deg, #b37feb, #9254de);
			}

			&.purple-dark {
				background: linear-gradient(135deg, #597ef7, #2f54eb);
			}

			&.orange-light {
				background: linear-gradient(135deg, #ff9c6e, #fa8c16);
			}

			&.green {
				background: linear-gradient(135deg, #73d13d, #52c41a);
			}

			&.orange-dark {
				background: linear-gradient(135deg, #ff7d00, #d4380d);
			}

			&.black {
				background: linear-gradient(135deg, #434343, #000000);
			}
		}

		.action-text {
			font-size: 22rpx;
			color: #595959;
			text-align: center;
		}
	}
}

.danger-card {
	background: linear-gradient(135deg, #ff4d4f, #d32f2f);
	border-radius: 20rpx;
	padding: 28rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	box-shadow: 0 6rpx 20rpx rgba(255, 77, 79, 0.25);
	transition: all 0.2s ease;

	&:active {
		transform: scale(0.98);
		opacity: 0.9;
	}

	.danger-icon {
		width: 56rpx;
		height: 56rpx;
		border-radius: 14rpx;
		background: rgba(255, 255, 255, 0.2);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-right: 14rpx;
	}

	.danger-text {
		font-size: 30rpx;
		font-weight: 600;
		color: #fff;
	}
}
</style>
