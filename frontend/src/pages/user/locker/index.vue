<template>
	<view class="container">
		<uv-navbar :title="t('user.locker.title')" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack" />

		<view class="content">
			<!-- 新增: 场所地址 -->
			<view class="location-bar fade-in-up" style="animation-delay: 0.1s;">
				<uv-icon name="map" size="16" color="#4facfe" />
				<text class="location-text">{{ t('user.locker.locationPlaceholder') }}: 北京市朝阳区科技园区A座大厅</text>
			</view>

			<!-- 新增: 使用说明入口 -->
			<view class="guide-bar fade-in-up" style="animation-delay: 0.15s;" @click="showUsageGuide">
				<uv-icon name="info-circle" size="16" color="#faad14" />
				<text class="guide-text">{{ t('user.locker.usageGuideTitle') }}</text>
				<uv-icon name="arrow-right" size="12" color="#ccc" />
			</view>

			<!-- 柜子列表 -->
			<view class="locker-grid">
				<view v-for="(item, index) in lockerList" :key="item.id" 
					class="locker-item fade-in-up" 
					:style="{ animationDelay: (0.1 + index * 0.05) + 's' }"
					:class="{
						'is-occupied-other': item.status === 1 && !item.isMine,
						'is-mine': item.isMine,
						'is-free': item.status === 0,
						'is-maintenance': item.status === 2
					}" @click="handleLockerClick(item)">
					<!-- 状态角标 -->
					<view class="status-corner bounce-in">
						<text>{{ getStatusText(item) }}</text>
					</view>

					<view class="locker-main">
						<view class="locker-icon-wrapper" :class="{
							'pulse-free': item.status === 0,
							'pulse-occupied': item.status === 1 && !item.isMine,
							'pulse-mine': item.isMine
						}">
							<uv-icon name="empty-favor" :size="48" :color="getIconColor(item)" />
						</view>
						<view class="locker-info">
							<text class="locker-no">{{ item.no }}</text>
							<text class="locker-size">{{ getLockerSizeText(item.size) }}</text>
						</view>
					</view>

					<view class="locker-footer">
						<text class="locker-price">¥{{ item.price }}<text class="unit">{{ t('user.locker.pricePerHour') }}</text></text>
						<view class="action-hint" v-if="item.isMine || !item.status">
							<uv-icon name="arrow-right" size="12" color="#fff" />
						</view>
					</view>
				</view>
			</view>
		</view>

		<!-- 新增: 底部统计信息 (固定在底部) -->
		<view class="stats-fixed-wrapper">
			<view class="stats-container">
				<view class="stats-item">
					<view class="stats-icon-box free-bg">
						<uv-icon name="checkmark-circle" size="18" color="#1890ff" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.available') }}</text>
						<text class="stats-num num-free">{{ stats.free }}</text>
					</view>
				</view>
				<view class="divider"></view>
				<view class="stats-item">
					<view class="stats-icon-box used-bg">
						<uv-icon name="close-circle" size="18" color="#ff4d4f" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.occupied') }}</text>
						<text class="stats-num num-used">{{ stats.used }}</text>
					</view>
				</view>
				<view class="divider"></view>
				<view class="stats-item">
					<view class="stats-icon-box maint-bg">
						<uv-icon name="setting" size="18" color="#faad14" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.maintenance') }}</text>
						<text class="stats-num num-maintenance">{{ stats.maintenance }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 操作弹窗：仅当选中本人租用的柜子时显示 -->
		<uv-popup ref="actionPopup" mode="center" :round="20">
			<view class="popup-content" v-if="selectedLocker">
				<view class="popup-header">
					<text class="popup-title">{{ selectedLocker.no }} - {{ getLockerSizeText(selectedLocker.size)
						}}</text>
					<view class="close-btn" @click="closePopup">
						<uv-icon name="close" size="20" color="#999" />
					</view>
				</view>

				<view class="popup-actions">
					<view class="action-item" @click="handleTempUnlock">
						<view class="action-icon-bg warning-bg">
							<uv-icon name="lock-open" color="#fff" size="24" />
						</view>
						<text class="action-text">{{ t('user.locker.tempUnlock') }}</text>
					</view>

					<view class="action-item" @click="handleEndUsePreCheck">
						<view class="action-icon-bg error-bg">
							<uv-icon name="shopping-cart" color="#fff" size="24" />
						</view>
						<text class="action-text">{{ t('user.locker.endUse') }}</text>
					</view>
				</view>
			</view>
		</uv-popup>

		<!-- 新增: 用户同意使用柜子承诺书组件 -->
		<LockerAgreement 
			ref="agreementPopup" 
			:title="t('user.locker.agreementTitle')"
			:content="t('user.locker.agreementContent')"
			:agreeText="t('user.locker.agreeAndContinue')"
			:disagreeText="t('user.locker.disagree')"
			@agree="onAgreementConfirm"
			@disagree="onAgreementCancel"
		/>

		<!-- 新增: 使用说明弹窗 (复用 uv-popup 或简单实现，这里为了简洁使用 uv-popup 展示纯文本) -->
		<uv-popup ref="guidePopup" mode="center" :round="20">
			<view class="guide-popup-content">
				<view class="popup-header">
					<text class="popup-title">{{ t('user.locker.usageGuideTitle') }}</text>
				</view>
				<scroll-view scroll-y class="guide-scroll">
					<text class="guide-detail-text">{{ t('user.locker.usageGuideContent') }}</text>
				</scroll-view>
				<view class="guide-footer">
					<button class="btn-know" @click="guidePopup.close()">{{ t('common.ok') }}</button>
				</view>
			</view>
		</uv-popup>
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
// 新增: 引入承诺书组件
import LockerAgreement from '@/components/LockerAgreement.vue'

const { t } = useI18n()

// 模拟柜子数据
// status: 0-空闲, 1-占用, 2-维修
// isMine: 是否为当前用户租用
const lockerList = ref([
	{ id: 1, no: 'A01', size: 'large', price: 5, status: 0, isMine: false },
	{ id: 2, no: 'A02', size: 'medium', price: 3, status: 1, isMine: true }, // 本人租用
	{ id: 3, no: 'A03', size: 'small', price: 2, status: 1, isMine: false }, // 他人租用
	{ id: 4, no: 'A04', size: 'large', price: 5, status: 0, isMine: false },
	{ id: 5, no: 'B01', size: 'medium', price: 3, status: 0, isMine: false },
	{ id: 6, no: 'B02', size: 'small', price: 2, status: 1, isMine: true }, // 本人租用第二个柜子
	// 新增6个柜子
	{ id: 7, no: 'B03', size: 'large', price: 5, status: 0, isMine: false },
	{ id: 8, no: 'B04', size: 'medium', price: 3, status: 1, isMine: false },
	{ id: 9, no: 'C01', size: 'small', price: 2, status: 0, isMine: false },
	{ id: 10, no: 'C02', size: 'large', price: 5, status: 2, isMine: false }, // 维修中
	{ id: 11, no: 'C03', size: 'medium', price: 3, status: 0, isMine: false },
	{ id: 12, no: 'C04', size: 'small', price: 2, status: 1, isMine: true },
])

const selectedLocker = ref(null)
const actionPopup = ref(null)
// 新增: 承诺书弹窗引用
const agreementPopup = ref(null)
// 新增: 使用说明弹窗引用
const guidePopup = ref(null)
// 新增: 暂存待租用的柜子信息
const pendingRentLocker = ref(null)
// 新增: 用户是否已同意承诺书的状态
const hasAgreed = ref(false)

// 修改: 计算统计信息 (按状态统计)
const stats = computed(() => {
	const res = {
		free: 0,
		used: 0,
		maintenance: 0
	}

	lockerList.value.forEach(item => {
		if (item.status === 0) res.free++
		else if (item.status === 1) res.used++
		else if (item.status === 2) res.maintenance++
	})

	return res
})

const goBack = () => {
	uni.navigateBack()
}

const getLockerSizeText = (size) => {
	if (size === 'large') return t('user.locker.large')
	if (size === 'medium') return t('user.locker.medium')
	return t('user.locker.small')
}

const getStatusText = (item) => {
	if (item.status === 2) return t('user.locker.statusMaintenance')
	if (item.isMine) return t('user.locker.statusMine')
	if (item.status === 1) return t('user.locker.statusOccupied')
	return t('user.locker.statusFree')
}

const getIconColor = (item) => {
	if (item.status === 2) return '#bfbfbf' // 维修中灰色
	if (item.isMine) return '#ffccc7' // 本人租用图标白色
	if (item.status === 1) return '#ffe58f' // 他人占用图标白色
	return '#4facfe' // 空闲时蓝色
}

// 新增: 显示使用说明
const showUsageGuide = () => {
	guidePopup.value.open()
}

const handleLockerClick = (item) => {
	if (item.status === 1 && item.isMine) {
		selectedLocker.value = item
		actionPopup.value.open()
	} else if (item.status === 0) {
		// 空闲柜子，检查是否已同意承诺书
		if (!hasAgreed.value) {
			uni.showToast({ 
				title: t('user.locker.mustAgreeFirst'), 
				icon: 'none' 
			})
			// 可选：再次弹出承诺书引导用户同意
			agreementPopup.value.open()
			return
		}
		
		// 已同意，直接弹出确认租用对话框（或者根据需求直接租用，这里保留原有的确认逻辑）
		pendingRentLocker.value = item
		uni.showModal({
			title: t('common.confirm'),
			content: t('user.locker.confirmRent', { no: item.no, price: item.price }),
			success: (res) => {
				if (res.confirm) {
					// 模拟租用成功
					item.status = 1
					item.isMine = true
					uni.showToast({ title: t('user.locker.rentSuccess'), icon: 'success' })
				}
			}
		})
		pendingRentLocker.value = null
	} else if (item.status === 1 && !item.isMine) {
		// 他人占用
		uni.showToast({ title: t('user.locker.occupied'), icon: 'none' })
	} else if (item.status === 2) {
		// 维修中
		uni.showToast({ title: t('user.locker.underMaintenance'), icon: 'none' })
	}
}

// 新增: 承诺书同意回调
const onAgreementConfirm = () => {
	// 标记为已同意
	hasAgreed.value = true
	// 保存到本地存储，以便下次进入页面无需再次同意
	uni.setStorageSync('locker_agreement_agreed', true)
	
	// 如果之前有暂存的租用请求（例如从其他逻辑触发），可以在此处理
	// 但根据当前逻辑，同意后如果是首次进入，通常只是关闭弹窗让用户浏览
	// 如果用户是在点击柜子后被拦截并打开的承诺书，这里不需要自动触发租用，
	// 用户需要再次点击柜子进行租用，或者你可以选择在这里直接触发 pendingRentLocker 的逻辑
	
	if (pendingRentLocker.value) {
		const item = pendingRentLocker.value
		uni.showModal({
			title: t('common.confirm'),
			content: t('user.locker.confirmRent', { no: item.no, price: item.price }),
			success: (res) => {
				if (res.confirm) {
					item.status = 1
					item.isMine = true
					uni.showToast({ title: t('user.locker.rentSuccess'), icon: 'success' })
				}
			}
		})
		pendingRentLocker.value = null
	}
}

// 新增: 承诺书取消/不同意回调
const onAgreementCancel = () => {
	// 不同意，不改变 hasAgreed 状态，保持为 false
	// 用户可以浏览页面，但点击租用时会提示
	pendingRentLocker.value = null
	// 不再强制提示 mustAgree，因为用户选择了不同意，只需关闭弹窗
	// uni.showToast({ title: t('user.locker.mustAgree'), icon: 'none' }) 
}

const closePopup = () => {
	actionPopup.value.close()
	selectedLocker.value = null
}

const handleTempUnlock = () => {
	if (!selectedLocker.value) return

	// 关闭弹窗
	closePopup()

	// 模拟开锁API
	uni.showLoading({ title: t('common.loading') })
	setTimeout(() => {
		uni.hideLoading()
		uni.showToast({ title: t('user.locker.unlockSuccess'), icon: 'success' })
	}, 1000)
}

// 新增：结束使用预检查，引导至支付流程
const handleEndUsePreCheck = () => {
	if (!selectedLocker.value) return

	// 缓存当前选中的柜子信息，防止 closePopup 清空后无法访问
	const currentLocker = { ...selectedLocker.value }

	// 先关闭操作弹窗
	closePopup()

	// 计算费用（模拟）
	const cost = currentLocker.price * 2 // 假设使用了2小时

	uni.showModal({
		title: t('common.confirm'),
		content: t('user.locker.expectedCost', { cost: cost }),
		confirmText: t('user.locker.immediatePay'),
		success: (res) => {
			if (res.confirm) {
				handlePaymentAndEndUse(cost, currentLocker.id)
			}
		}
	})
}

// 新增：支付并结束使用
const handlePaymentAndEndUse = (amount, lockerId) => {
	uni.showLoading({ title: t('user.locker.paymentProcessing') })

	// 模拟支付请求
	setTimeout(() => {
		uni.hideLoading()

		// 模拟支付成功
		uni.showToast({ title: t('user.locker.paymentSuccess'), icon: 'success' })

		// 支付成功后执行结束使用逻辑
		setTimeout(() => {
			executeEndUse(lockerId)
		}, 1500)

	}, 1500)
}

const executeEndUse = (lockerId) => {
	// 模拟结束使用API
	uni.showLoading({ title: t('common.loading') })
	setTimeout(() => {
		uni.hideLoading()
		// 更新本地状态
		const index = lockerList.value.findIndex(l => l.id === lockerId)
		if (index !== -1) {
			lockerList.value[index].status = 0
			lockerList.value[index].isMine = false
		}
		selectedLocker.value = null
		uni.showToast({ title: t('user.locker.endSuccess'), icon: 'success' })
	}, 1000)
}

// 新增: 页面加载时检查同意状态
onMounted(() => {
	const agreed = uni.getStorageSync('locker_agreement_agreed')
	if (agreed) {
		hasAgreed.value = true
	} else {
		// 未同意，弹出承诺书
		// 使用 nextTick 或 setTimeout 确保 DOM 渲染完成后打开弹窗，避免某些 UI 库的警告
		setTimeout(() => {
			agreementPopup.value.open()
		}, 500)
	}
})

</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f8f9fc;
	/* 更柔和的背景色 */
	padding-bottom: 160rpx;
	box-sizing: border-box;
}

.content {
	padding: 24rpx;
}

// 新增: 场所地址样式
.location-bar {
	display: flex;
	align-items: center;
	background-color: #fff;
	padding: 20rpx 24rpx;
	border-radius: 12rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);

	.location-text {
		font-size: 26rpx;
		color: #333;
		margin-left: 12rpx;
		flex: 1;
	}
}

// 新增: 使用说明条目样式
.guide-bar {
	display: flex;
	align-items: center;
	background-color: #fff;
	padding: 20rpx 24rpx;
	border-radius: 12rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);
	
	.guide-text {
		font-size: 26rpx;
		color: #333;
		margin-left: 12rpx;
		flex: 1;
	}
}

.locker-grid {
	display: flex;
	flex-wrap: wrap;
	justify-content: space-between;
	gap: 20rpx;
	/* 使用 gap 替代 margin-bottom 以获得更均匀的间距 */
}

.locker-item {
	width: calc(50% - 10rpx);
	/* 精确计算宽度以适配 gap */
	height: 260rpx;
	background-color: #ffffff;
	border-radius: 20rpx;
	position: relative;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	padding: 24rpx;
	box-sizing: border-box;
	transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.2s ease;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	border: 2rpx solid transparent;
	overflow: hidden;

	/* 空闲状态 - 浅蓝色 */
	&.is-free {
		background: linear-gradient(145deg, #e6f7ff 0%, #ffffff 100%);
		border-color: #bae7ff;

		.locker-no,
		.locker-size {
			color: #0050b3;
		}

		.locker-price {
			color: #1890ff;
		}
	}

	/* 被他人占用状态 - 黄色/暖色 */
	&.is-occupied-other {
		background: linear-gradient(145deg, #fffbe6 0%, #ffffff 100%);
		border-color: #ffe58f;

		.locker-icon-wrapper {
			opacity: 0.9;
		}

		.locker-no,
		.locker-size {
			color: #ad6800;
		}

		.locker-price {
			color: #faad14;
		}
	}

	/* 我的柜子状态 - 红色/醒目 */
	&.is-mine {
		background: linear-gradient(145deg, #fff1f0 0%, #ffffff 100%);
		border: 2rpx solid #ffccc7;
		box-shadow: 0 8rpx 24rpx rgba(255, 77, 79, 0.15);
		transform: translateY(-4rpx);

		.locker-no,
		.locker-size {
			color: #cf1322;
		}

		.locker-price {
			color: #ff4d4f;
			font-weight: bold;
		}
	}

	/* 维修状态 - 灰色/不可用 */
	&.is-maintenance {
		background: linear-gradient(145deg, #f5f5f5 0%, #ffffff 100%);
		border: 2rpx dashed #d9d9d9;
		opacity: 0.8;

		.locker-no,
		.locker-size {
			color: #8c8c8c;
		}

		.locker-price {
			color: #bfbfbf;
		}
		
		.locker-icon-wrapper {
			filter: grayscale(100%);
		}
	}

	&:active {
		transform: scale(0.96);
	}
}

.status-corner {
	position: absolute;
	top: 0;
	right: 0;
	padding: 6rpx 16rpx;
	border-bottom-left-radius: 12rpx;
	z-index: 2;

	text {
		font-size: 25rpx;
		font-weight: bold;
		letter-spacing: 1rpx;
	}

	/* 空闲 - 蓝底白字或蓝字 */
	.is-free & {
		background-color: rgba(24, 144, 255, 0.1);

		text {
			color: #1890ff;
		}
	}

	/* 他人占用 - 黄底黄字 */
	.is-occupied-other & {
		background-color: rgba(250, 173, 20, 0.1);

		text {
			color: #faad14;
		}
	}

	/* 我的柜子 - 红底红字 */
	.is-mine & {
		background-color: rgba(255, 77, 79, 0.1);

		text {
			color: #ff4d4f;
		}
	}
	
	/* 维修 - 灰底灰字 */
	.is-maintenance & {
		background-color: rgba(0, 0, 0, 0.04);

		text {
			color: #8c8c8c;
		}
	}
}

.locker-main {
	display: flex;
	align-items: center;
	margin-top: 10rpx;
}

.locker-icon-wrapper {
	width: 80rpx;
	height: 80rpx;
	display: flex;
	justify-content: center;
	align-items: center;
	background-color: rgba(255, 255, 255, 0.9);
	border-radius: 50%;
	margin-right: 16rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
	transition: all 0.3s ease;

	/* 空闲状态呼吸灯 - 蓝色 */
	&.pulse-free {
		animation: pulse-blue 2s infinite ease-in-out;
	}

	/* 他人占用呼吸灯 - 黄色 */
	&.pulse-occupied {
		animation: pulse-yellow 2s infinite ease-in-out;
	}

	/* 本人租用呼吸灯 - 红色 */
	&.pulse-mine {
		animation: pulse-red 1.5s infinite ease-in-out;
		box-shadow: 0 0 0 0 rgba(255, 77, 79, 0.4);
	}
}

/* 蓝色呼吸动画 */
@keyframes pulse-blue {
	0% {
		box-shadow: 0 0 0 0 rgba(24, 144, 255, 0.4);
	}
	70% {
		box-shadow: 0 0 0 10rpx rgba(24, 144, 255, 0);
	}
	100% {
		box-shadow: 0 0 0 0 rgba(24, 144, 255, 0);
	}
}

/* 黄色呼吸动画 */
@keyframes pulse-yellow {
	0% {
		box-shadow: 0 0 0 0 rgba(250, 173, 20, 0.4);
	}
	70% {
		box-shadow: 0 0 0 10rpx rgba(250, 173, 20, 0);
	}
	100% {
		box-shadow: 0 0 0 0 rgba(250, 173, 20, 0);
	}
}

/* 红色呼吸动画 */
@keyframes pulse-red {
	0% {
		box-shadow: 0 0 0 0 rgba(255, 77, 79, 0.4);
	}
	70% {
		box-shadow: 0 0 0 12rpx rgba(255, 77, 79, 0);
	}
	100% {
		box-shadow: 0 0 0 0 rgba(255, 77, 79, 0);
	}
}

/* 列表项入场动画 */
@keyframes fade-in-up {
	from {
		opacity: 0;
		transform: translateY(20rpx);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.fade-in-up {
	animation: fade-in-up 0.5s ease-out forwards;
	opacity: 0; /* 初始隐藏，等待动画执行 */
}

/* 角标弹跳动画 */
@keyframes bounce-in {
	0% {
		opacity: 0;
		transform: scale(0.3);
	}
	50% {
		opacity: 1;
		transform: scale(1.05);
	}
	70% {
		transform: scale(0.9);
	}
	100% {
		transform: scale(1);
	}
}

.bounce-in {
	animation: bounce-in 0.6s cubic-bezier(0.215, 0.61, 0.355, 1) forwards;
}

.locker-info {
	display: flex;
	flex-direction: column;

	.locker-no {
		font-size: 36rpx;
		font-weight: 600;
		color: #333;
		line-height: 1.2;
	}

	.locker-size {
		font-size: 25rpx;
		margin-top: 4rpx;
		background-color: rgba(0, 0, 0, 0.03);
		padding: 2rpx 8rpx;
		border-radius: 6rpx;
		align-self: flex-start;
	}
}

.locker-footer {
	display: flex;
	justify-content: space-between;
	align-items: flex-end;
	margin-top: auto;

	.locker-price {
		font-size: 28rpx;
		font-weight: 500;
		color: #ff9900;

		.unit {
			font-size: 20rpx;
			font-weight: normal;
			opacity: 0.8;
			margin-left: 2rpx;
		}
	}

	.action-hint {
		width: 32rpx;
		height: 32rpx;
		background-color: #faad14;
		border-radius: 50%;
		display: flex;
		justify-content: center;
		align-items: center;
	}
}

// 新增: 底部固定统计样式
.stats-fixed-wrapper {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background-color: transparent; /* 改为透明，让阴影由容器承担 */
	z-index: 99;
	padding-bottom: env(safe-area-inset-bottom);
	pointer-events: none; /* 允许点击穿透空白区域，如果需要的话，或者保留none防止遮挡 */
	
	.stats-container {
		pointer-events: auto; /* 恢复统计栏的点击事件 */
		background-color: #fff;
		margin: 0 24rpx 24rpx; /* 增加外边距，悬浮感 */
		border-radius: 24rpx;
		padding: 24rpx 10rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 10rpx 30rpx rgba(0, 0, 0, 0.08);
		backdrop-filter: blur(10px); /* 毛玻璃效果 */

		.stats-item {
			display: flex;
			align-items: center;
			flex: 1;
			justify-content: center;
			
			.stats-icon-box {
				width: 48rpx;
				height: 48rpx;
				border-radius: 12rpx;
				display: flex;
				justify-content: center;
				align-items: center;
				margin-right: 12rpx;
				
				&.free-bg { background-color: rgba(24, 144, 255, 0.1); }
				&.used-bg { background-color: rgba(255, 77, 79, 0.1); }
				&.maint-bg { background-color: rgba(250, 173, 20, 0.1); }
			}
			
			.stats-info {
				display: flex;
				flex-direction: column;
				align-items: flex-start;
				
				.stats-label {
					font-size: 22rpx;
					color: #999;
					margin-bottom: 4rpx;
				}
				
				.stats-num {
					font-size: 34rpx;
					font-weight: 800;
					line-height: 1;
					
					&.num-free { color: #1890ff; }
					&.num-used { color: #ff4d4f; }
					&.num-maintenance { color: #faad14; }
				}
			}
		}
		
		.divider {
			width: 1rpx;
			height: 40rpx;
			background-color: #f0f0f0;
			margin: 0 10rpx;
		}
	}
}

// 新增: 使用说明弹窗样式
.guide-popup-content {
	width: 600rpx;
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	display: flex;
	flex-direction: column;
	max-height: 70vh;
}

.guide-scroll {
	flex: 1;
	min-height: 300rpx;
	max-height: 500rpx;
	margin-bottom: 30rpx;
	
	.guide-detail-text {
		font-size: 28rpx;
		color: #666;
		line-height: 1.6;
		white-space: pre-wrap;
	}
}

.guide-footer {
	.btn-know {
		background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
		color: #fff;
		border-radius: 40rpx;
		height: 80rpx;
		line-height: 80rpx;
		font-size: 28rpx;
		border: none;
		
		&::after {
			border: none;
		}
	}
}

.popup-content {
	background-color: #fff;
	border-radius: 20rpx;
	/* 修改为全圆角，适配中间弹出 */
	padding: 30rpx;
	min-height: 300rpx;
	width: 600rpx;
	/* 增加宽度限制，使中间弹窗更美观 */
}

.popup-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 40rpx;

	.popup-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}

	.close-btn {
		padding: 10rpx;
	}
}

.popup-actions {
	display: flex;
	justify-content: space-around;
	padding-bottom: 40rpx;
}

.action-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 200rpx;

	.action-icon-bg {
		width: 100rpx;
		height: 100rpx;
		border-radius: 50%;
		display: flex;
		justify-content: center;
		align-items: center;
		margin-bottom: 20rpx;

		&.warning-bg {
			background-color: #ff9900;
		}

		&.error-bg {
			background-color: #ff4d4f;
		}
	}

	.action-text {
		font-size: 28rpx;
		color: #333;
	}
}
</style>