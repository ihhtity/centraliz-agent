<!-- 商家资料页面 -->
<template>
	<view class="container">
		<uv-navbar :title="'个人中心'" :placeholder="true" :leftIcon="''" />

		<!-- 账号信息卡片 -->
		<view class="profile-card" @click="navTo('/pages/admin/profile/basic?merchs_id=' + merch.id)">
			<view class="profile-left">
				<view class="avatar-box">
					<uv-icon name="account" size="48" color="#fff" />
				</view>
			</view>
			<view class="profile-right">
				<view class="profile-item">
					<text class="profile-label">账号</text>
					<text class="profile-value">{{ profile.account || '-' }}</text>
				</view>
				<view class="profile-item">
					<text class="profile-label">邮箱</text>
					<text class="profile-value">{{ profile.email || '未绑定' }}</text>
				</view>
				<view class="profile-item">
					<text class="profile-label">手机</text>
					<text class="profile-value">{{ profile.phone || '未绑定' }}</text>
				</view>
			</view>
			<view>
				<uv-icon name="arrow-right" size="30" />
			</view>
		</view>

		<!-- 功能菜单 -->
		<view class="section-title">快捷操作</view>
		<view class="grid-menu">
			<view class="menu-item" @click="navTo('/pages/admin/profile/test')">
				<view class="icon-box blue">
					<uv-icon name="attach" size="32" color="#fff" />
				</view>
				<text>测试页面</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/profile/huifu')">
				<view class="icon-box blue">
					<uv-icon name="empty-coupon" size="32" color="#fff" />
				</view>
				<text>汇付测试</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/profile/basic?merchs_id=' + merch.id)">
				<view class="icon-box blue">
					<uv-icon name="account" size="32" color="#fff" />
				</view>
				<text>基本信息</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/device/list')">
				<view class="icon-box blue">
					<uv-icon name="setting" size="32" color="#fff" />
				</view>
				<text>设备管理</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/rule/manage')">
				<view class="icon-box orange">
					<uv-icon name="list" size="32" color="#fff" />
				</view>
				<text>规则管理</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/order/list')">
				<view class="icon-box green">
					<uv-icon name="order" size="32" color="#fff" />
				</view>
				<text>交易明细</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/account/manage?merchs_id=' + merch.id)">
				<view class="icon-box purple">
					<uv-icon name="empty-data" size="32" color="#fff" />
				</view>
				<text>子账号</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/huifu/list')">
				<view class="icon-box green">
					<uv-icon name="coupon" size="32" color="#fff" />
				</view>
				<text>收款账号</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/value-added/index')">
				<view class="icon-box purple">
					<uv-icon name="gift" size="32" color="#fff" />
				</view>
				<text>增值功能</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/expense/index')">
				<view class="icon-box orange">
					<uv-icon name="empty-coupon" size="32" color="#fff" />
				</view>
				<text>我的消费</text>
			</view>
			<view class="menu-item" @click="navTo('/pages/admin/profile/about')">
				<view class="icon-box cyan">
					<uv-icon name="info-circle" size="32" color="#fff" />
				</view>
				<text>关于我们</text>
			</view>
			<view class="menu-item" @click="logout">
				<view class="icon-box red">
					<uv-icon name="share-square" size="32" color="#fff" />
				</view>
				<text>退出登录</text>
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
// 商家资料页面脚本
import { onShow } from '@dcloudio/uni-app';
import { ref } from 'vue';

// 底部导航索引
const tabbar = ref(1);

// 商家信息
const profile = ref({
	account: '',  // 账号
	email: '',    // 邮箱
	phone: ''     // 手机号
});
const merch = ref({});

// 页面加载时获取商家信息
onShow(() => {
	merch.value = uni.getStorageSync('merch') || {};
	loadProfile();
});

// 获取商家信息
const loadProfile = async () => {
	const merchsId = merch.value.id || '';
	try {
		const res = await uni.$uv.http.get('/merch/profile', {
			params: { merchs_id: merchsId },
			custom: { auth: true }
		});
		if (res.code === 200 && res.data) {
			profile.value = {
				account: res.data.account || '',
				email: res.data.email || '',
				phone: res.data.phone || ''
			};
		}
	} catch (e) {
		console.error('加载失败', e);
	}
};

// 切换底部导航
const editTabbar = (e) => {
	tabbar.value = e;
	if (e === 0) {
		uni.reLaunch({
			url: '/pages/admin/index/index'
		});
	}
};

// 退出登录
const logout = () => {
	uni.showModal({
		title: '确认',
		content: '确定要退出登录吗？',
		cancelText: '取消',
		confirmText: '确定',
		success: (res) => {
			if (res.confirm) {
				uni.removeStorageSync('token');
				uni.removeStorageSync('userInfo');
				uni.reLaunch({ url: '/pages/login/login?isAdmin=true' });
			}
		}
	});
};

// 页面跳转
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
	background: linear-gradient(135deg, #3c9cff, #2b85e4);
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
	width: 80rpx;
}

.profile-value {
	font-size: 26rpx;
	color: #333;
	font-weight: 500;
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
</style>