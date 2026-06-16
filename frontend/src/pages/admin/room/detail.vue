<!-- 房间详情页面 -->
<template>
	<view class="container">
		<uv-navbar title="房间详情" :placeholder="true" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container" enable-flex>
			<!-- 房间基本信息 -->
			<view class="info-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="home-fill" />
					</view>
					<text class="card-title">房间信息</text>
				</view>
				<view class="info-grid">
					<view class="info-item">
						<text class="info-label">房间名称</text>
						<text class="info-value">{{ roomDetail.name || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">房间状态</text>
						<view class="status-badge" :class="getStatusClass(roomDetail.status)">
							{{ roomDetail.status }}
						</view>
					</view>
					<view class="info-item">
						<text class="info-label">房间标签</text>
						<text class="info-value">{{ roomDetail.tag || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">所属分组</text>
						<text class="info-value">{{ roomDetail.groupName || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">设备名称</text>
						<text class="info-value">{{ device.name || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">板号-锁号</text>
						<text class="info-value">{{ roomDetail.boardNo || '-' }}-{{ roomDetail.lockNo || '-' }}</text>
					</view>
				</view>
			</view>

			<!-- 房间操作区 -->
			<view class="action-card">
				<view class="card-header">
					<view class="header-icon">
						<uv-icon name="setting-fill" />
					</view>
					<text class="card-title">房间操作</text>
				</view>
				<view class="action-grid">
					<view class="action-item" @click="showQRCodeModal">
						<view class="action-icon">
							<uv-icon name="scan" size="40" color="#333" />
						</view>
						<text class="action-text">二维码</text>
					</view>
					<view class="action-item" @click="handleOpenLock">
						<view class="action-icon">
							<uv-icon name="lock-opened-fill" size="40" color="primary" />
						</view>
						<text class="action-text">远程开锁</text>
					</view>
					<view class="action-item" @click="handleEditRoom">
						<view class="action-icon">
							<uv-icon name="edit-pen-fill" size="40" color="#ffb830" />
						</view>
						<text class="action-text">房间编辑</text>
					</view>
					<view class="action-item" @click="handleOperationLog">
						<view class="action-icon">
							<uv-icon name="order" size="40" color="#19be6b" />
						</view>
						<text class="action-text">操作记录</text>
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

		<!-- 弹框 -->
		<view>
			<!-- 二维码弹窗 -->
			<view v-if="showQRCodeSheet" class="action-sheet-mask" @click="closeQRCodeModal">
				<view class="qr-code-sheet" @click.stop>
					<view class="sheet-title">房间二维码</view>

					<view class="qr-code-content">
						<view class="qr-code-wrapper">
							<uv-qrcode ref="qrcodeRef" :value="qrCodeContent" :size="200" :isQueueLoadImage="true"
								:h5SaveIsDownload="true" :h5DownloadName="`${roomDetail.name}.png`" />
						</view>
						<view class="qr-code-desc">
							<text class="desc-text">扫码访问房间</text>
							<text class="desc-hint">{{ qrCodeContent }}</text>
						</view>
					</view>

					<view class="qr-actions">
						<view class="qr-btn cancel" @click="closeQRCodeModal">关闭</view>
						<view class="qr-btn confirm" @click="saveQRCode">保存图片</view>
					</view>
				</view>
			</view>
			<!-- 编辑弹窗 -->
			<view v-if="showEditModal" class="modal-overlay" @click="closeQRCodeModal">
				<view class="modal-content" @click.stop>
					<view class="modal-header">
						<text class="modal-title">编辑房间</text>
						<view class="modal-close" @click="showEditModal = false">
							<uv-icon name="close" />
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
							<text class="form-label">房间锁号</text>
							<input class="form-input" v-model="editForm.lockNo" placeholder="请输入房间锁号" />
						</view>
						<view class="form-item">
							<text class="form-label">房间状态</text>
							<view class="status-options">
								<view v-for="option in statusOptions" :key="option.value" class="status-option"
									:class="{ active: editForm.status === option.value }"
									@click="editForm.status = option.value">
									{{ option.value }}
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
						<uv-icon name="warning-fill" size="80rpx" color="#ff4d4f" />
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
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { generateQRCodeContent, generateLockCommand, recordOperationLog } from '@/utils/utils'

onLoad((options) => {
	if (options.id) roomId.value = options.id
	loadRoomDetail()
})

// 房间id
const roomId = ref('')
// 房间详情
const roomDetail = ref({})
// 设备详情
const device = ref({})
// 编辑弹窗相关
const showEditModal = ref(false)
// 编辑房间表单数据
const editForm = ref({
	name: '',
	tag: '',
	lockNo: '01',
	status: '0'
})
// 房间状态选项
const statusOptions = [{ value: '空闲' }, { value: '租用' }, { value: '维修' }]
// 二维码弹出层状态
const showQRCodeSheet = ref(false);
// 二维码内容
const qrCodeContent = ref('');
// 二维码组件ref
const qrcodeRef = ref(null);

// 远程开锁
const handleOpenLock = async () => {
	try {
		uni.showLoading({ title: '开锁中...' })
		let hexData = "574B4C5909" + roomDetail.value.boardNo + "82" + roomDetail.value.lockNo
		// 生成锁命令
		const data = {
			code: device.value.code,
			command: generateLockCommand(hexData)
		}
		// 记录操作日志
		const logData = {
			merchsId: roomDetail.value.merchsId,
			devicesId: roomDetail.value.devicesId,
			roomId: roomDetail.value.id,
			code: device.value.code,
			deviceName: device.value.name,
			roomName: roomDetail.value.name,
			control: "开锁",
			status: "成功",
			occupant: "商家",
		}

		const result = await uni.$uv.http.post('/device/common', data, {
			custom: { auth: true }
		})
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '开锁成功', icon: 'success' })
			recordOperationLog(logData)
		} else {
			uni.showToast({ title: result.msg || '开锁失败', icon: 'none' })
			logData.status = "失败"
			recordOperationLog(logData)
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '开锁失败', icon: 'none' })
		logData.status = "失败"
		recordOperationLog(logData)
	}
}

// 显示二维码弹窗
const showQRCodeModal = () => {
	if (!roomDetail.value.id) {
		uni.showToast({ title: '房间ID为空', icon: 'none' });
		return;
	}
	qrCodeContent.value = generateQRCodeContent('room', roomDetail.value.groupType, roomDetail.value.id);
	showQRCodeSheet.value = true;
};

// 关闭二维码弹窗
const closeQRCodeModal = () => {
	showQRCodeSheet.value = false;
};

// 保存二维码图片
const saveQRCode = async () => {
	uni.showLoading({ title: '保存中' });

	try {
		// #ifdef H5
		// H5端使用uv-qrcode组件的save方法
		if (qrcodeRef.value && qrcodeRef.value.save) {
			await qrcodeRef.value.save({
				content: qrCodeContent.value,
				success: (res) => {
					console.log('二维码保存成功:', res);
					uni.showToast({ title: '保存成功', icon: 'success' });
				},
				fail: (err) => {
					console.error('二维码保存失败:', err);
					uni.showToast({ title: '保存失败', icon: 'none' });
				}
			});

			closeQRCodeModal();
		} else {
			uni.hideLoading();
			uni.showToast({ title: '组件未就绪', icon: 'none' });
		}
		// #endif

		// #ifndef H5
		// 小程序端使用uv-qrcode组件的save方法
		if (qrcodeRef.value && qrcodeRef.value.save) {
			await qrcodeRef.value.save({
				content: qrCodeContent.value,
				success: (res) => {
					console.log('二维码保存成功:', res);
					uni.showToast({ title: '保存成功', icon: 'success' });
				},
				fail: (err) => {
					console.error('二维码保存失败:', err);
					uni.showToast({ title: '保存失败', icon: 'none' });
				}
			});

			closeQRCodeModal();
		} else {
			uni.showToast({ title: '组件未就绪', icon: 'none' });
		}
		// #endif
	} catch (err) {
		console.error('保存二维码失败:', err);
		uni.showToast({ title: '保存失败', icon: 'none' });
	}
};

// 加载房间详情
const loadRoomDetail = async () => {
	uni.showLoading({ title: '加载中...' })
	try {
		const res = await uni.$uv.http.get('/room/' + roomId.value, {
			custom: { auth: true }
		})
		if (res.code === 200 && res.data) {
			roomDetail.value = res.data
			device.value = res.data.device || {}
		}
	} catch (e) {
		console.error('加载房间详情失败', e)
		uni.showToast({ title: '加载失败', icon: 'none' })
		roomDetail.value = {
			id: roomId.value,
			name: roomDetail.value.name || '',
			tag: roomDetail.value.tag || '普通柜',
			lockNo: roomDetail.value.lockNo || '01',
			status: roomDetail.value.status || '空闲',
		}
	} finally {
		uni.hideLoading()
	}
}

// 格式化锁号为2-3位数
const formatLockNo = (lockNo) => {
	const num = parseInt(lockNo) || 0
	if (num < 0) return '00'
	if (num > 999) return '999'
	// 锁号：0-9 -> 2位(00-09), 10-99 -> 2位(10-99), 100-999 -> 3位(100-999)
	if (num < 10) {
		return num.toString().padStart(2, '0')
	}
	return num.toString()
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
	if (!editForm.value.lockNo.trim()) {
		uni.showToast({ title: '请输入房间锁号', icon: 'none' })
		return
	}

	let data = {
		name: editForm.value.name.trim(),
		tag: editForm.value.tag.trim(),
		status: editForm.value.status,
		lock_no: formatLockNo(editForm.value.lockNo.trim()),
	}

	try {
		uni.showLoading({ title: '保存中...' })
		const result = await uni.$uv.http.put('/room/' + roomId.value, data, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '修改成功', icon: 'success', duration: 1500 })
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

// 打开编辑弹窗
const handleEditRoom = () => {
	editForm.value = {
		name: roomDetail.value.name || '',
		tag: roomDetail.value.tag || '普通柜',
		lockNo: roomDetail.value.lockNo || '01',
		status: roomDetail.value.status || '空闲'
	}
	showEditModal.value = true
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

// 确认删除弹窗
const showDeleteModal = ref(false)

// 打开删除确认弹窗
const handleDeleteRoom = () => {
	showDeleteModal.value = true
}

// 获取房间状态类名
const getStatusClass = (status) => {
	if (!status) return 'status-unknown'
	if (status === '空闲') return 'status-free'
	if (status === '租用') return 'status-rent'
	if (status === '维修') return 'status-maint'
	return 'status-unknown'
}

// 查看操作记录
const handleOperationLog = () => {
	uni.navigateTo({
		url: `/pages/admin/device/log?roomId=${roomId.value}`
	})
}

const goBack = () => uni.navigateBack()
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

/* 弹窗遮罩样式 */
.action-sheet-mask {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: flex-end;
	z-index: 1000;
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

/* 二维码弹窗样式 */
.qr-code-sheet {
	width: 100%;
	background: #fff;
	border-radius: 32rpx 32rpx 0 0;
	padding-bottom: env(safe-area-inset-bottom);
	max-height: 80vh;
	overflow-y: auto;
	animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
	from {
		transform: translateY(100%);
		opacity: 0;
	}

	to {
		transform: translateY(0);
		opacity: 1;
	}
}

.sheet-title {
	text-align: center;
	padding: 40rpx;
	font-size: 32rpx;
	font-weight: 600;
	color: #333;
	border-bottom: 1rpx solid #f0f0f0;
	background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
}

.qr-code-content {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 32rpx;
}

.qr-code-wrapper {
	width: 360rpx;
	height: 360rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #ffffff;
	border: 2rpx solid #e8e8e8;
	border-radius: 20rpx;
	box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.12);
	overflow: hidden;
	flex-shrink: 0;
	transition: transform 0.3s ease, box-shadow 0.3s ease;

	&:active {
		transform: scale(0.98);
		box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
	}
}

.qr-code-wrapper canvas {
	width: 100% !important;
	height: 100% !important;
	display: block;
}

/* 二维码加载占位 */
.qr-code-loading {
	width: 60rpx;
	height: 60rpx;
	border: 4rpx solid #e8e8e8;
	border-top-color: #3c9cff;
	border-radius: 50%;
	animation: spin 1s linear infinite;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

/* 微信小程序特殊处理：禁止二维码区域滚动 */
/* #ifdef MP-WEIXIN */
.qr-code-content {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 32rpx;
	touch-action: none;
}

.qr-code-wrapper {
	pointer-events: none;
}

/* #endif */

.qr-code-desc {
	text-align: center;
	width: 100%;

	.desc-text {
		font-size: 32rpx;
		font-weight: 600;
		color: #333333;
		display: block;
		margin-bottom: 16rpx;
	}

	.desc-hint {
		font-size: 24rpx;
		color: #999999;
		display: block;
		word-break: break-all;
		max-width: 100%;
		line-height: 1.5;
		padding: 0 16rpx;
	}
}

.qr-actions {
	display: flex;
	padding: 0 32rpx 32rpx;
	gap: 24rpx;
	border-top: 1rpx solid #f5f5f5;
	padding-top: 28rpx;

	.qr-btn {
		flex: 1;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 44rpx;
		font-size: 30rpx;
		font-weight: 500;

		&.cancel {
			background: #f5f7fa;
			color: #666666;
			border: 1rpx solid #e8e8e8;
		}

		&.confirm {
			background: linear-gradient(135deg, #3c9cff 0%, #2979ff 100%);
			color: #ffffff;
			box-shadow: 0 4rpx 16rpx rgba(60, 156, 255, 0.3);
		}
	}
}

/* 遮罩动画 */
.action-sheet-mask {
	animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
	from {
		opacity: 0;
	}

	to {
		opacity: 1;
	}
}
</style>
