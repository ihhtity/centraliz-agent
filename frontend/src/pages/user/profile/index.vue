<template>
	<view class="container">
		<!-- 顶部导航栏，设置为透明或自定义颜色以配合背景 -->
		<uv-navbar :title="t('user.profile.title')" :placeholder="true" :leftIcon="''" />

		<view class="header-bg">
			<view class="header-content">
				<view class="avatar-wrapper" @click="handleEditAvatar">
					<!-- 如果有头像URL则显示图片，否则显示默认图标 -->
					<image v-if="userInfo.avatar" :src="userInfo.avatar" mode="aspectFill" class="avatar-img"></image>
					<view v-else class="avatar-default">
						<uv-icon name="account" color="#fff" size="60" />
					</view>
					<view class="edit-badge">
						<uv-icon name="camera-fill" color="#fff" size="16" />
					</view>
				</view>
				<view class="info">
					<text class="nickname">{{ userInfo.nickname || t('user.profile.username') }}</text>
					<view class="phone-tag" @click="handleCopyPhone">
						<uv-icon name="phone" size="14" color="#fff"></uv-icon>
						<text class="phone">{{ userInfo.phone || t('user.profile.phone') }}</text>
					</view>
				</view>
			</view>
		</view>

		<view class="menu-group">
			<uv-cell-group :border="false">
				<uv-cell :title="t('user.profile.wallet')" isLink icon="empty-coupon"
					:iconStyle="{ color: '#ff9900', marginRight: '10px' }" @click="goToWallet">
					<template #value>
						<text class="cell-value">{{ t('user.profile.balance') }}</text>
					</template>
				</uv-cell>
				<uv-cell :title="t('user.profile.realname')" isLink icon="checkmark-circle"
					:iconStyle="{ color: '#19be6b', marginRight: '10px' }" @click="goToRealname">
					<template #value>
						<text class="cell-value" :class="{ 'verified': true }">{{ t('user.profile.unverified') }}</text>
					</template>
				</uv-cell>
				<uv-cell :title="t('user.profile.language')" isLink icon="setting"
					:iconStyle="{ color: '#909399', marginRight: '10px' }" @click="goToLanguage"></uv-cell>
				<uv-cell :title="t('user.locker.title')" isLink icon="bag"
					:iconStyle="{ color: '#909399', marginRight: '10px' }" @click="goToLocker"></uv-cell>
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
		<uv-tabbar :value="2" :placeholder="true" @change="onTabBarChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.order')" icon="order" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>
	</view>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

// 模拟用户默认参数，实际项目中应从 pinia/vuex 或接口获取
const userInfo = reactive({
	nickname: '',
	phone: '',
	avatar: ''
})

// 初始化加载用户信息
const loadUserInfo = () => {
	// TODO: 替换为实际的获取用户信息逻辑
	// const store = useUserStore()
	// userInfo.nickname = store.nickname
	// userInfo.phone = store.phone
	// userInfo.avatar = store.avatar

	// 演示用默认值
	userInfo.nickname = 'User_8888'
	userInfo.phone = '138****8888'
}

loadUserInfo()

const handleEditAvatar = () => {
	uni.chooseImage({
		count: 1,
		success: (res) => {
			console.log('选择头像', res.tempFilePaths[0])
			// TODO: 上传头像逻辑
			uni.showToast({ title: '头像更换功能开发中', icon: 'none' })
		}
	})
}

const handleCopyPhone = () => {
	if (!userInfo.phone) return
	uni.setClipboardData({
		data: userInfo.phone,
		success: () => {
			uni.showToast({ title: t('common.copied'), icon: 'success' })
		}
	})
}

const goToWallet = () => {
	uni.navigateTo({ url: '/pages/user/wallet/index' });
}

const goToRealname = () => {
	uni.navigateTo({ url: '/pages/user/realname/index' });
}

const goToContact = () => {
	uni.makePhoneCall({
		phoneNumber: userInfo.phone || '10086' // 示例号码
	})
}

const goToLanguage = () => {
	uni.navigateTo({ url: '/pages/user/language/index' });
};

const goToLocker = () => {
	uni.navigateTo({ url: '/pages/user/locker/index' });
};

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
				uni.removeStorageSync('userInfo')

				uni.reLaunch({ url: '/pages/login/login' });
			}
		}
	});
};

// 新增: TabBar 切换逻辑
const onTabBarChange = (index) => {
	if (index === 0) {
		uni.redirectTo({ url: '/pages/user/index/index' })
	} else if (index === 1) {
		uni.redirectTo({ url: '/pages/user/order/list' })
	} else if (index === 2) {
		// 当前已是我的页，无需操作
		return
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
	padding: 40rpx 30rpx;
}

.avatar-wrapper {
	position: relative;
	width: 160rpx;
	height: 160rpx;
	margin-right: 30rpx;
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

	.edit-badge {
		position: absolute;
		bottom: 0;
		right: 0;
		width: 40rpx;
		height: 40rpx;
		background-color: #fff;
		border-radius: 50%;
		display: flex;
		justify-content: center;
		align-items: center;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
	}
}

.info {
	display: flex;
	flex-direction: column;
	justify-content: center;
	color: #fff;

	.nickname {
		font-size: 36rpx;
		font-weight: 600;
		margin-bottom: 12rpx;
		text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.1);
	}

	.phone-tag {
		display: inline-flex;
		align-items: center;
		background-color: rgba(255, 255, 255, 0.25);
		padding: 6rpx 16rpx;
		border-radius: 30rpx;
		width: fit-content;
		backdrop-filter: blur(4px);
		transition: all 0.3s;

		&:active {
			background-color: rgba(255, 255, 255, 0.4);
		}

		.phone {
			font-size: 24rpx;
			color: #fff;
			margin-left: 8rpx;
		}
	}
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