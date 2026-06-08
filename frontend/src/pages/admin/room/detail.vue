<!-- 房间详情页面 -->
<template>
	<view class="container">
		<uv-navbar :title="roomName || '房间详情'" :placeholder="true" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container" enable-flex>
			<!-- 房间基本信息 -->
			<view class="info-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="home-fill"  />
					</view>
					<text class="card-title">房间信息</text>
				</view>
				<view class="info-grid">
					<view class="info-item">
						<text class="info-label">房间标签</text>
						<text class="info-value">{{ roomDetail.tag || '普通柜' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">房间状态</text>
						<view class="status-badge" :class="getStatusClass(roomDetail.status)">
							{{ roomDetail.status }}
						</view>
					</view>
					<view class="info-item">
						<text class="info-label">所属分组</text>
						<text class="info-value">{{ roomDetail.groupName || '无' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">设备数量</text>
						<text class="info-value">{{ devices.length }} 台</text>
					</view>
				</view>
			</view>

			<!-- 设备列表 -->
			<view v-if="devices.length > 0" class="device-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="setting-fill"  />
					</view>
					<text class="card-title">设备列表</text>
				</view>
				<view v-for="device in devices" :key="device.id" class="device-item">
					<view class="device-info">
						<text class="device-name">{{ device.name || '设备' + device.id }}</text>
						<text class="device-status">{{ device.status || '正常' }}</text>
					</view>
					<view class="device-action" @click="handleUnbindDevice(device.id)">
						<text class="unbind-text">解绑</text>
					</view>
				</view>
			</view>

			<!-- 房间操作区 -->
			<view class="action-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="setting-fill"  />
					</view>
					<text class="card-title">房间操作</text>
				</view>
				<view class="action-grid">
					<view class="action-item" @click="handleOpenLock">
						<view class="action-icon">
							<uv-icon name="lock-opened-fill" size="40" color="primary" />
						</view>
						<text class="action-text">远程开锁</text>
					</view>
					<view class="action-item" @click="handleShowQrcode">
						<view class="action-icon">
							<uv-icon name="scan" size="40" color="#333" />
						</view>
						<text class="action-text">二维码</text>
					</view>
					<view class="action-item" @click="handleEditRoom">
						<view class="action-icon">
							<uv-icon name="edit-pen-fill" size="40" color="#ffb830" />
						</view>
						<text class="action-text">房间编辑</text>
					</view>
					<view class="action-item" @click="handleDeleteRoom">
						<view class="action-icon">
							<uv-icon name="trash-fill" size="40" color="red" />
						</view>
						<text class="action-text">房间删除</text>
					</view>
				</view>
			</view>
		</scroll-view>

		<!-- 编辑弹窗 -->
		<view v-if="showEditModal" class="modal-overlay" @click="showEditModal = false">
			<view class="modal-content" @click.stop>
				<view class="modal-header">
					<text class="modal-title">编辑房间</text>
					<view class="modal-close" @click="showEditModal = false">
						<uv-icon name="close"  />
					</view>
				</view>
				<view class="modal-body">
					<view class="form-item">
						<text class="form-label">房间名称</text>
						<input class="form-input" v-model="editForm.name" placeholder="请输入房间名称" />
					</view>
					<view class="form-item">
						<text class="form-label">房间标签</text>
						<input class="form-input" v-model="editForm.tag" placeholder="请输入房间标签" />
					</view>
					<view class="form-item">
						<text class="form-label">房间状态</text>
						<view class="status-options">
							<view v-for="option in statusOptions" :key="option.value" class="status-option"
								:class="{ active: editForm.status === option.value }"
								@click="editForm.status = option.value">
								{{ option.label }}
							</view>
						</view>
					</view>
				</view>
				<view class="modal-footer">
					<view class="btn-cancel" @click="showEditModal = false">取消</view>
					<view class="btn-confirm" @click="submitEdit">确认修改</view>
				</view>
			</view>
		</view>

		<!-- 删除确认弹窗 -->
		<view v-if="showDeleteModal" class="modal-overlay" @click="showDeleteModal = false">
			<view class="modal-content delete-modal" @click.stop>
				<view class="modal-icon">
					<uv-icon name="warning-fill" size="80rpx" color="#ff4d4f"  />
				</view>
				<text class="delete-title">确认删除房间</text>
				<text class="delete-desc">确定要删除此房间吗？此操作不可恢复。</text>
				<view class="modal-footer">
					<view class="btn-cancel" @click="showDeleteModal = false">取消</view>
					<view class="btn-confirm danger" @click="confirmDeleteRoom">确认删除</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const roomId = ref('')
const roomName = ref('')
const roomDetail = ref({})
const devices = ref([])

// 编辑弹窗相关
const showEditModal = ref(false)
const editForm = ref({
	name: '',
	tag: '',
	status: '0'
})

const statusOptions = [
	{ label: '空闲', value: '空闲' },
	{ label: '租用', value: '租用' },
	{ label: '维修', value: '维修' }
]

onLoad((options) => {
	if (options.id) roomId.value = options.id
	if (options.name) roomName.value = decodeURIComponent(options.name)
	loadRoomDetail()
})

// 加载房间详情
const loadRoomDetail = async () => {
	uni.showLoading({ title: '加载中...' })
	try {
		const res = await uni.$uv.http.get('/room/' + roomId.value, {
			custom: { auth: true }
		})
		if (res.code === 200 && res.data) {
			roomDetail.value = res.data
			devices.value = res.data.devices || []
		}
	} catch (e) {
		console.error('加载房间详情失败', e)
		uni.showToast({ title: '加载失败', icon: 'none' })
		roomDetail.value = {
			id: roomId.value,
			name: roomName.value,
			tag: '普通柜',
			status: '空闲',
		}
		devices.value = []
	} finally {
		uni.hideLoading()
	}
}

const goBack = () => uni.navigateBack()

const getStatusClass = (status) => {
	if (!status) return 'status-unknown'
	if (status === '空闲') return 'status-free'
	if (status === '租用') return 'status-rent'
	if (status === '维修') return 'status-maint'
	return 'status-unknown'
}

// 远程开锁
const handleOpenLock = async () => {
	try {
		uni.showLoading({ title: '开锁中...' })
		const result = await uni.$uv.http.post('/room/' + roomId.value + '/open-lock', {}, {
			custom: { auth: true }
		})
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '开锁成功', icon: 'success' })
		} else {
			uni.showToast({ title: result.msg || '开锁失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '开锁失败', icon: 'none' })
	}
}

// 显示二维码
const handleShowQrcode = async () => {
	try {
		uni.showLoading({ title: '获取中...' })
		const result = await uni.$uv.http.get('/room/' + roomId.value + '/qrcode', {
			custom: { auth: true }
		})
		uni.hideLoading()
		if (result.code === 200 && result.data) {
			const data = result.data
			uni.showActionSheet({
				itemList: ['H5二维码', '小程序二维码'],
				success: (res) => {
					if (res.tapIndex === 0) {
						uni.setClipboardData({
							data: data.h5QrcodeUrl,
							success: () => {
								uni.showToast({ title: 'H5链接已复制', icon: 'none' })
							}
						})
					} else {
						uni.setClipboardData({
							data: data.miniQrcodeUrl,
							success: () => {
								uni.showToast({ title: '小程序链接已复制', icon: 'none' })
							}
						})
					}
				}
			})
		} else {
			uni.showToast({ title: result.msg || '获取失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '获取失败', icon: 'none' })
	}
}

// 打开编辑弹窗
const handleEditRoom = () => {
	editForm.value = {
		name: roomDetail.value.name || '',
		tag: roomDetail.value.tag || '普通柜',
		status: roomDetail.value.status || '0'
	}
	showEditModal.value = true
}

// 提交编辑
const submitEdit = async () => {
	if (!editForm.value.name.trim()) {
		uni.showToast({ title: '请输入房间名称', icon: 'none' })
		return
	}
	if (!editForm.value.tag.trim()) {
		uni.showToast({ title: '请输入房间标签', icon: 'none' })
		return
	}
	try {
		uni.showLoading({ title: '保存中...' })
		const result = await uni.$uv.http.put('/room/' + roomId.value, {
			name: editForm.value.name.trim(),
			tag: editForm.value.tag.trim(),
			status: editForm.value.status
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '修改成功', icon: 'success' })
			roomName.value = editForm.value.name.trim()
			showEditModal.value = false
			loadRoomDetail()
		} else {
			uni.showToast({ title: result.msg || '修改失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '修改失败', icon: 'none' })
	}
}

// 确认删除弹窗
const showDeleteModal = ref(false)

// 打开删除确认弹窗
const handleDeleteRoom = () => {
	showDeleteModal.value = true
}

// 确认删除房间
const confirmDeleteRoom = async () => {
	showDeleteModal.value = false
	uni.showLoading({ title: '删除中...' })
	await uni.$uv.http.delete('/room/' + roomId.value, {}, { custom: { auth: true } })
		.then((res) => {
			uni.hideLoading()
			if (res.code === 200) {
				uni.showToast({ title: '删除成功', icon: 'success' })
				setTimeout(() => uni.navigateBack(), 1500)
			} else {
				uni.showToast({ title: res.msg || '删除失败', icon: 'none' })
			}
		})
		.catch((e) => {
			uni.hideLoading()
			uni.showToast({ title: '删除失败', icon: 'none' })
		})
}
// 解绑设备
const handleUnbindDevice = async (deviceId) => {
	try {
		uni.showLoading({ title: '解绑中...' })
		const result = await uni.$uv.http.post('/room/unbind-device', {
			roomId: roomId.value,
			deviceId: deviceId
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '解绑成功', icon: 'success' })
			loadRoomDetail()
		} else {
			uni.showToast({ title: result.msg || '解绑失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '解绑失败', icon: 'none' })
	}
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

.info-card,
.device-card,
.action-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 28rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
	margin-bottom: 20rpx;
}

.card-header {
	display: flex;
	align-items: center;
	margin-bottom: 20rpx;
}

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
}

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
}

.status-free {
	background: rgba(24, 144, 255, 0.1);
	color: #1890ff;
}

.status-rent {
	background: rgba(255, 77, 79, 0.1);
	color: #ff4d4f;
}

.status-maint {
	background: rgba(250, 173, 20, 0.1);
	color: #faad14;
}

.status-unknown {
	background: rgba(153, 153, 153, 0.1);
	color: #999;
}

.device-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f5f5f5;

	&:last-child {
		border-bottom: none;
	}
}

.device-info {
	display: flex;
	flex-direction: column;
}

.device-name {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;
}

.device-status {
	font-size: 24rpx;
	color: #999;
	margin-top: 4rpx;
}

.unbind-text {
	font-size: 26rpx;
	color: #ff4d4f;
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
	border-radius: 16rpx;

	&:active {
		background: #f5f5f5;
		transform: scale(0.96);
	}
}

.action-icon {
	width: 80rpx;
	height: 80rpx;
	border-radius: 20rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 10rpx;
}

.action-text {
	font-size: 22rpx;
	color: #595959;
	text-align: center;
}

/* 弹窗样式 */
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.modal-content {
	width: 600rpx;
	background: #fff;
	border-radius: 24rpx;
	overflow: hidden;
}

.modal-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 32rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.modal-title {
	font-size: 32rpx;
	font-weight: 600;
	color: #1a1a1a;
}

.modal-close {
	width: 48rpx;
	height: 48rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #999;
}

.modal-body {
	padding: 32rpx;
}

.form-item {
	margin-bottom: 28rpx;
}

.form-label {
	font-size: 26rpx;
	color: #666;
	margin-bottom: 12rpx;
	display: block;
}

.form-input {
	width: 100%;
	height: 80rpx;
	padding: 0 20rpx;
	border: 1rpx solid #e8e8e8;
	border-radius: 12rpx;
	font-size: 28rpx;
	box-sizing: border-box;
}

.status-options {
	display: flex;
	gap: 20rpx;
}

.status-option {
	flex: 1;
	height: 72rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	border: 1rpx solid #e8e8e8;
	border-radius: 12rpx;
	font-size: 26rpx;
	color: #666;
	transition: all 0.3s;

	&.active {
		background: #2979ff;
		border-color: #2979ff;
		color: #fff;
	}
}

.modal-footer {
	display: flex;
	border-top: 1rpx solid #f0f0f0;
}

.btn-cancel,
.btn-confirm {
	flex: 1;
	height: 88rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 30rpx;
}

.btn-cancel {
	color: #666;
	border-right: 1rpx solid #f0f0f0;
}

.btn-confirm {
	color: #2979ff;
	font-weight: 500;

	&.danger {
		color: #ff4d4f;
	}
}

/* 删除确认弹窗 */
.delete-modal {
	padding: 48rpx 32rpx;
	text-align: center;
}

.modal-icon {
	margin-bottom: 24rpx;
}

.delete-title {
	display: block;
	font-size: 34rpx;
	font-weight: 600;
	color: #1a1a1a;
	margin-bottom: 12rpx;
}

.delete-desc {
	display: block;
	font-size: 26rpx;
	color: #666;
	margin-bottom: 32rpx;
}
</style>
