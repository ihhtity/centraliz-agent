<template>
	<view class="container">
		<!-- 顶部导航栏 -->
		<uv-navbar :title="t('user.profile.title')" :placeholder="true" :leftIcon="''" />

		<!-- 账号信息卡片 -->
		<view class="profile-card" @click="goToBasicInfo">
			<view class="profile-left">
				<view class="avatar-box">
					<uv-icon name="account-fill" size="48" color="#fff" />
				</view>
			</view>
			<view class="profile-right">
				<view class="profile-item">
					<text class="profile-label">{{ t('user.profile.nickname') }}</text>
					<text class="profile-value">{{ userInfo.name || t('user.profile.hidden') }}</text>
				</view>
				<view class="profile-item">
					<text class="profile-label">{{ t('user.profile.account') }}</text>
					<text class="profile-value">{{ userInfo.account || '-' }}</text>
				</view>
				<view class="profile-item" v-if="userInfo.phone">
					<uv-icon name="phone" size="14" color="#999" />
					<text class="profile-value">{{ userInfo.phone }}</text>
				</view>
				<view class="profile-item" v-if="userInfo.email">
					<uv-icon name="email" size="14" color="#999" />
					<text class="profile-value">{{ userInfo.email }}</text>
				</view>
			</view>
			<view class="arrow-icon">
				<uv-icon name="arrow-right" size="30" color="#999" />
			</view>
		</view>

		<!-- 功能菜单 -->
		<view class="section-title">{{ t('user.profile.quickActions') }}</view>
		<view class="grid-menu">
			<view class="menu-item" @click="goToWallet">
				<view class="icon-box orange">
					<uv-icon name="red-packet" size="32" color="#fff" />
				</view>
				<text>{{ t('user.profile.wallet') }}</text>
			</view>
			<view class="menu-item" @click="openDepositRefund">
				<view class="icon-box purple">
					<uv-icon name="empty-coupon" size="32" color="#fff" />
				</view>
				<text>{{ t('user.profile.depositRefund') }}</text>
			</view>
			<view class="menu-item" @click="goToOrderList">
				<view class="icon-box green">
					<uv-icon name="order" size="32" color="#fff" />
				</view>
				<text>{{ t('user.profile.orderRecord') }}</text>
			</view>
			<view class="menu-item" @click="goToContact">
				<view class="icon-box cyan">
					<uv-icon name="phone" size="32" color="#fff" />
				</view>
				<text>{{ t('user.profile.contact') }}</text>
			</view>
			<view class="menu-item" @click="goToLanguage">
				<view class="icon-box blue">
					<uv-icon name="setting" size="32" color="#fff" />
				</view>
				<text>{{ t('user.profile.language') }}</text>
			</view>
			<view class="menu-item" @click="logout">
				<view class="icon-box red">
					<uv-icon name="share-square" size="32" color="#fff" />
				</view>
				<text>{{ t('user.profile.logout') }}</text>
			</view>
		</view>

		<!-- 底部导航栏 -->
		<uv-tabbar :value="1" :placeholder="true" @change="onTabBarChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>

		<!-- 押金退款弹出层 -->
		<view v-if="depositPopupVisible" class="deposit-popup-mask" @click="closeDepositPopup">
			<view class="deposit-popup-content" @click.stop>
				<view class="deposit-popup-header">
					<text class="deposit-popup-title">{{ t('user.profile.depositInfo') }}</text>
					<view class="deposit-popup-close" @click="closeDepositPopup">
						<uv-icon name="close" size="24" color="#999" />
					</view>
				</view>
				<view class="deposit-popup-body">
					<view v-if="depositInfo.hasDeposit" class="deposit-info">
						<view class="deposit-info-item">
							<text class="deposit-info-label">{{ t('user.profile.depositAmount') }}</text>
							<text class="deposit-info-value">¥{{ depositInfo.deposit }}</text>
						</view>
						<view class="deposit-info-item">
							<text class="deposit-info-label">{{ t('user.profile.depositOrderId') }}</text>
							<text class="deposit-info-value">{{ depositInfo.code }}</text>
						</view>
						<view class="deposit-info-item">
							<text class="deposit-info-label">{{ t('user.profile.depositStatus') }}</text>
							<text class="deposit-info-value">{{ t('user.profile.depositPaid') }}</text>
						</view>
					</view>
					<view v-else class="deposit-empty">
						<uv-icon name="info-circle" size="60" color="#ccc" />
						<text class="deposit-empty-text">{{ t('user.profile.noDeposit') }}</text>
					</view>
				</view>
				<view class="deposit-popup-footer">
					<view v-if="depositInfo.hasDeposit" class="deposit-refund-btn" @click="handleDepositRefund">
						<text>{{ t('user.profile.applyRefund') }}</text>
					</view>
					<view v-else class="deposit-close-btn" @click="closeDepositPopup">
						<text>{{ t('common.close') }}</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { onShow } from '@dcloudio/uni-app';

onShow(() => {
    userInfo.value = uni.getStorageSync('user') || {}
})

const { t } = useI18n()
// 响应式用户信息对象
const userInfo = ref({})
// 押金弹出层显示状态
const depositPopupVisible = ref(false)
// 押金信息
const depositInfo = ref({
	hasDeposit: false,
	deposit: 0,
	orderId: 0,
	code: '',
})

// 打开押金退款弹出层
const openDepositRefund = async () => {
	if (!userInfo.value?.id) {
		uni.showToast({ title: t('common.loginFirst'), icon: 'none' })
		return
	}

	uni.showLoading({ title: t('common.loading') })
	try {
		// 检查押金状态
		const res = await uni.$uv.http.post('/user/deposit/check', {
			usersId: userInfo.value.id,
			merchsId: 0 // 商家ID，可根据实际情况传入
		}, { custom: { auth: true } })

		if (res.code === 200) {
			depositInfo.value = {
				hasDeposit: res.data.hasDeposit || false,
				deposit: res.data.deposit || 0,
				orderId: res.data.orderId || 0,
				code: res.data.code || '',
			}
			depositPopupVisible.value = true
		}
	} catch (error) {
		console.error('检查押金状态失败:', error)
		uni.showToast({ title: t('common.networkError'), icon: 'none' })
	} finally {
		uni.hideLoading()
	}
}
// 申请押金退款
const handleDepositRefund = () => {
	if (!depositInfo.value.hasDeposit || !depositInfo.value.orderId) {
		return
	}

	// 关闭弹退款弹出层
	closeDepositPopup()

	uni.showModal({
		title: t('common.confirm'),
		content: t('user.profile.depositRefundConfirm', { deposit: depositInfo.value.deposit }),
		cancelText: t('common.cancel'),
		confirmText: t('common.confirm'),
		success: async (res) => {
			if (res.confirm) {
				uni.showLoading({ title: t('common.loading') })
				try {
					const result = await uni.$uv.http.post('/user/deposit/refund', {
						usersId: userInfo.value.id,
						orderId: depositInfo.value.orderId
					}, { custom: { auth: true } })

					if (result.code === 200) {
						uni.showToast({ title: t('user.profile.depositRefundSuccess'), icon: 'none', duration: 2000 })
						depositInfo.value.hasDeposit = false
						depositInfo.value.deposit = 0
						depositInfo.value.orderId = 0
						depositPopupVisible.value = false
					} else {
						uni.showToast({ title: result.msg || t('common.operationFailed'), icon: 'none', duration: 2000 })
					}
				} catch (error) {
					console.error('押金退款失败:', error)
					uni.showToast({ title: t('common.networkError'), icon: 'none', duration: 2000 })
				} finally {
					uni.hideLoading()
				}
			}
		}
	})
}
// 关闭押金退款弹出层
const closeDepositPopup = () => {
	depositPopupVisible.value = false
}
// 跳转到基本资料页面
const goToBasicInfo = () => {
	uni.navigateTo({ url: '/pages/user/profile/basic' })
}
// 跳转到钱包页面
const goToWallet = () => {
	uni.navigateTo({ url: '/pages/user/wallet/index' });
}
// 跳转到联系页面
const goToContact = () => {
	uni.makePhoneCall({
		phoneNumber: userInfo.value.phone || '10086'
	})
}
// 跳转到语言设置页面
const goToLanguage = () => {
	uni.navigateTo({ url: '/pages/user/language/index' });
};
// 跳转到订单记录页面
const goToOrderList = () => {
	uni.navigateTo({ url: '/pages/user/order/list' });
};
// 退出登录
const logout = () => {
	uni.showModal({
		title: t('common.confirm'),
		content: t('user.profile.logoutConfirm'),
		cancelText: t('common.cancel'),
		confirmText: t('common.confirm'),
		success: (res) => {
			if (res.confirm) {
				// 清除本地存储
				uni.removeStorageSync('token')
				uni.removeStorageSync('user')

				uni.reLaunch({ url: '/pages/login/login' });
			}
		}
	});
};
// TabBar 切换逻辑
const onTabBarChange = () => {
	let userroute = uni.getStorageSync('userroute')
	if (userroute === 'locker') {
		uni.redirectTo({ url: '/pages/user/index/locker' })
	} else if (userroute === 'retail') {
		uni.redirectTo({ url: '/pages/user/index/retail' })
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding: 24rpx;
	padding-bottom: 140rpx;
}

.profile-card {
	display: flex;
	align-items: center;
	background: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.profile-left {
	margin-right: 60rpx;
}

.avatar-box {
	width: 100rpx;
	height: 100rpx;
	border-radius: 50%;
	background: linear-gradient(135deg, #4facfe, #00f2fe);
	display: flex;
	align-items: center;
	justify-content: center;
}

.profile-right {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: center;
}

.profile-item {
	display: flex;
	align-items: center;
	height: 56rpx;
}

.profile-label {
	font-size: 26rpx;
	color: #999;
	width: 100rpx;
}

.profile-value {
	font-size: 26rpx;
	color: #333;
	font-weight: 500;
}

.arrow-icon {
	padding: 10rpx;
}

.section-title {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 24rpx;
	padding-left: 10rpx;
	border-left: 8rpx solid #4facfe;
	line-height: 1;
}

.grid-menu {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 24rpx;
	background: #fff;
	padding: 40rpx 20rpx;
	border-radius: 20rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	.menu-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding: 10rpx 0;
		transition: all 0.2s;

		&:active {
			opacity: 0.7;
			transform: scale(0.95);
		}

		.icon-box {
			width: 80rpx;
			height: 80rpx;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-bottom: 16rpx;
			box-shadow: 0 8rpx 16rpx rgba(0, 0, 0, 0.1);

			&.blue {
				background: linear-gradient(135deg, #3c9cff, #2b85e4);
			}

			&.orange {
				background: linear-gradient(135deg, #ff9900, #f29100);
			}

			&.red {
				background: linear-gradient(135deg, #fa3534, #e63332);
			}

			&.purple {
				background: linear-gradient(135deg, #a855f7, #9333ea);
			}

			&.green {
				background: linear-gradient(135deg, #10b981, #059669);
			}

			&.cyan {
				background: linear-gradient(135deg, #06b6d4, #0891b2);
			}
		}

		text {
			font-size: 24rpx;
			color: #666;
		}
	}
}

/* 押金退款弹出层样式 */
.deposit-popup-mask {
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

.deposit-popup-content {
	width: 600rpx;
	background: #fff;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.15);
}

.deposit-popup-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 32rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.deposit-popup-title {
	font-size: 32rpx;
	font-weight: 600;
	color: #333;
}

.deposit-popup-close {
	padding: 10rpx;
}

.deposit-popup-body {
	padding: 40rpx 32rpx;
	min-height: 200rpx;
}

.deposit-info {
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.deposit-info-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	height: 60rpx;
}

.deposit-info-label {
	font-size: 28rpx;
	color: #999;
}

.deposit-info-value {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;
}

.deposit-empty {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 40rpx 0;
}

.deposit-empty-text {
	font-size: 28rpx;
	color: #999;
	margin-top: 16rpx;
}

.deposit-popup-footer {
	padding: 24rpx 32rpx 32rpx;
}

.deposit-refund-btn {
	width: 100%;
	height: 88rpx;
	background: linear-gradient(135deg, #4facfe, #00f2fe);
	border-radius: 44rpx;
	display: flex;
	align-items: center;
	justify-content: center;

	text {
		font-size: 32rpx;
		color: #fff;
		font-weight: 600;
	}
}

.deposit-close-btn {
	width: 100%;
	height: 88rpx;
	background: #f5f5f5;
	border-radius: 44rpx;
	display: flex;
	align-items: center;
	justify-content: center;

	text {
		font-size: 32rpx;
		color: #666;
	}
}
</style>
