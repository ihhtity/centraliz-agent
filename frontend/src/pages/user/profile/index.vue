<template>
	<view class="container">
		<!-- 顶部导航栏，设置为透明或自定义颜色以配合背景 -->
		<uv-navbar :title="t('user.profile.title')" :placeholder="true" :leftIcon="''" />

		<view class="header-bg" @click="goToBasicInfo">
			<view class="header-content">
				<!-- 左侧头像 -->
				<view class="avatar-wrapper">
					<view class="avatar-default">
						<uv-icon name="account-fill" color="#fff" size="60" />
					</view>
				</view>
				
				<!-- 中间用户信息 -->
				<view class="info">
					<view class="info-row">
						<text class="info-text">昵称：</text>
						<text class="info-text">{{ userInfo.name || '隐藏' }}</text>
					</view>
					<view class="info-row">
						<text class="info-text">账号：</text>
						<text class="info-text">{{ userInfo.account || '-' }}</text>
					</view>
					<view class="info-row" v-if="userInfo.phone">
						<uv-icon name="phone" size="12" color="rgba(255,255,255,0.8)" />
						<text class="info-text">{{ userInfo.phone || '-' }}</text>
					</view>
					<view class="info-row" v-if="userInfo.email">
						<uv-icon name="email" size="12" color="rgba(255,255,255,0.8)" />
						<text class="info-text">{{ userInfo.email || '-'  }}</text>
					</view>
				</view>
				
				<!-- 右侧箭头 -->
				<view class="arrow-right">
					<uv-icon name="arrow-right" color="rgba(255,255,255,0.8)" size="20" />
				</view>
			</view>
		</view>

		<view class="menu-group">
			<uv-cell-group :border="false">
				<uv-cell :title="t('user.profile.wallet')" isLink icon="empty-coupon"
					:iconStyle="{ color: '#ff9900', marginRight: '10px' }" @click="goToWallet">
					<template #value>
						<text class="cell-value wallet-balance">¥{{ '0.00' }}</text>
					</template>
				</uv-cell>
				<uv-cell :title="t('user.profile.language')" isLink icon="setting"
					:iconStyle="{ color: '#909399', marginRight: '10px' }" @click="goToLanguage"></uv-cell>
				<uv-cell title="订单记录" isLink icon="order"
					:iconStyle="{ color: '#909399', marginRight: '10px' }" @click="goToOrderList"></uv-cell>
				<uv-cell :title="t('user.profile.contact')" isLink icon="phone"
					:iconStyle="{ color: '#3c9cff', marginRight: '10px' }" @click="goToContact"></uv-cell>
				<uv-cell :title="t('user.profile.logout')" isLink icon="share-square"
					:iconStyle="{ color: '#fa3534', marginRight: '10px' }" @click="logout"></uv-cell>
			</uv-cell-group>
		</view>

		<view class="footer-version">
			<text>v1.0.0</text>
		</view>

		<!-- 新增: 底部 TabBar -->
		<uv-tabbar :value="1" :placeholder="true" @change="onTabBarChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { onShow } from '@dcloudio/uni-app';

onShow(() => {
    userInfo.value = uni.getStorageSync('user')
})

const { t } = useI18n()
// 响应式用户信息对象 - 仅从本地存储获取
const userInfo = ref({})

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
		phoneNumber: userInfo.value.phone || '10086' // 示例号码
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

// 新增: TabBar 切换逻辑
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
	position: relative;
}

/* 顶部背景区域 */
.header-bg {
	width: 100%;
	height: 360rpx;
	background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
	border-bottom-left-radius: 40rpx;
	border-bottom-right-radius: 40rpx;
	position: relative;
	padding-top: var(--window-top);
	/* 适配刘海屏 */
	box-shadow: 0 10rpx 20rpx rgba(79, 172, 254, 0.2);
}

.header-content {
	display: flex;
	align-items: center;
	padding: 60rpx 30rpx;
	justify-content: space-between;
}

.avatar-wrapper {
	width: 140rpx;
	height: 140rpx;
	flex-shrink: 0;

	.avatar-img {
		width: 100%;
		height: 100%;
		border-radius: 50%;
		border: 4rpx solid rgba(255, 255, 255, 0.8);
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.1);
	}

	.avatar-default {
		width: 100%;
		height: 100%;
		border-radius: 50%;
		background-color: rgba(255, 255, 255, 0.2);
		display: flex;
		justify-content: center;
		align-items: center;
		border: 4rpx solid rgba(255, 255, 255, 0.8);
	}
}

.info {
	flex: 1;
	display: flex;
	flex-direction: column;
	justify-content: center;
	margin-left: 24rpx;
	color: #fff;

	.nickname {
		font-size: 36rpx;
		font-weight: 600;
		margin-bottom: 12rpx;
		text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.1);
	}

	.info-row {
		display: flex;
		align-items: center;
		margin-bottom: 8rpx;

		&:last-child {
			margin-bottom: 0;
		}

		.info-text {
			font-size: 24rpx;
			color: rgba(255, 255, 255, 0.9);
			margin-left: 8rpx;
		}
	}
}

.arrow-right {
	flex-shrink: 0;
	padding: 10rpx;
}

.menu-group {
	margin: -60rpx 20rpx 20rpx;
	/* 向上重叠背景 */
	position: relative;
	z-index: 10;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.05);
	background: #fff;

	::v-deep .uv-cell {
		padding: 26rpx 20rpx;

		.uv-cell__title {
			font-size: 28rpx;
			color: #333;
		}

		.uv-cell__value {
			font-size: 24rpx;
			color: #999;
		}
	}

	.cell-value {
		font-size: 24rpx;
		color: #999;

		&.verified {
			color: #19be6b;
		}

		&.wallet-balance {
			color: #ff9900;
			font-weight: 600;
		}
	}
}

.footer-version {
	text-align: center;
	padding: 40rpx 0;
	color: #ccc;
	font-size: 24rpx;
	/* 新增: 增加底部 padding 以避免被 TabBar 遮挡，虽然 placeholder=true 会处理，但这是双重保险 */
	padding-bottom: 120rpx;
}
</style>