<template>
	<view class="container">
		<uv-navbar :title="t('admin.profile.title')" :placeholder="true" :leftIcon="''" />
		
		<!-- 功能菜单 -->
		<view class="section-title">{{ t('admin.index.quickActions') }}</view>
		<view class="grid-menu">
			<view class="menu-item" @click="navTo('/pages/admin/device/list')">
				<view class="icon-box blue">
					<uv-icon name="setting" size="32" color="#fff"></uv-icon>
				</view>
				<text>{{ t('admin.index.deviceManagement') }}</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/rule/manage')">
				<view class="icon-box orange">
					<uv-icon name="list" size="32" color="#fff"></uv-icon>
				</view>
				<text>{{ t('admin.index.ruleManagement') }}</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/account/manage')">
				<view class="icon-box purple">
					<uv-icon name="account" size="32" color="#fff"></uv-icon>
				</view>
				<text>{{ t('admin.account.title') }}</text>
			</view>
			<view class="menu-item" @click="logout">
				<view class="icon-box red">
					<uv-icon name="share-square" size="32" color="#fff"></uv-icon>
				</view>
				<text>{{ t('admin.index.logout') }}</text>
			</view>
		</view>

		<!-- 底部导航栏 -->
		<view>
			<uv-tabbar :value="tabbar" @change="editTabbar">
				<uv-tabbar-item text="房间" icon="home" />
				<uv-tabbar-item text="个人" icon="account" />
			</uv-tabbar>
		</view>
	</view>
</template>

<script setup>
import { useI18n } from 'vue-i18n';
import { ref } from 'vue';

const { t } = useI18n();
const tabbar = ref(1);

const editTabbar = (e) => {
	// 更新当前选中的 tab
	tabbar.value = e;
	
	// 根据点击的 tab 进行页面跳转
	if (e === 0) {
		uni.navigateTo({
			url: '/pages/admin/index/index'
		});
	}
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
				
				uni.reLaunch({ url: '/pages/login/login?isAdmin=true' });
			}
		}
	});
};

const navTo = (url) => {
	if (url === '/pages/login/login') {
		uni.reLaunch({ url });
	} else {
		uni.navigateTo({ url });
	}
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding: 24rpx;
}

.section-title {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 24rpx;
	padding-left: 10rpx;
	border-left: 8rpx solid #3c9cff;
	line-height: 1;
}

.grid-menu {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 24rpx;
	background: #fff;
	padding: 40rpx 20rpx;
	border-radius: 20rpx;
	box-shadow: 0 4rpx 16rpx rgba(0,0,0,0.04);
	
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
			box-shadow: 0 8rpx 16rpx rgba(0,0,0,0.1);
			
			&.blue { background: linear-gradient(135deg, #3c9cff, #2b85e4); }
			&.orange { background: linear-gradient(135deg, #ff9900, #f29100); }
			&.red { background: linear-gradient(135deg, #fa3534, #e63332); }
			&.purple { background: linear-gradient(135deg, #a855f7, #9333ea); }
		}
		
		text {
			font-size: 24rpx;
			color: #666;
		}
	}
}
</style>