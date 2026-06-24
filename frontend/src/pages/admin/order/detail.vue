<!-- 订单详情页面 -->
<template>
	<view class="container">
		<uv-navbar :title="'订单详情'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
		
		<!-- 订单状态 -->
		<view class="status-card" :class="getStatusClass(order.status)">
			<view class="status-icon">{{ getStatusIcon(order.status) }}</view>
			<view class="status-info">
				<text class="status-text">{{ order.status }}</text>
				<text class="status-desc">{{ getStatusDesc(order.status) }}</text>
			</view>
		</view>

		<!-- 订单信息 -->
		<view class="info-card">
			<view class="card-title">
				<uv-icon name="info-circle" size="20" color="#3c9cff" />
				<view>订单信息</view>
			</view>
			<view class="info-item">
				<text class="info-label">订单编号</text>
				<text class="info-value">{{ order.code || '暂无' }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">订单名称</text>
				<text class="info-value">{{ order.name || '暂无' }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">用户手机号</text>
				<text class="info-value">{{ order.userPhone || '暂无' }}</text>
			</view>
			<!-- <view class="info-item">
				<text class="info-label">商家手机号</text>
				<text class="info-value">{{ order.merchPhone || '暂无' }}</text>
			</view> -->
			<view class="info-item">
				<text class="info-label">使用时长</text>
				<text class="info-value">{{ formatDuration(order.duration) }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">开始时间</text>
				<text class="info-value">{{ formatTime(order.startTime) }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">结束时间</text>
				<text class="info-value">{{ formatTime(order.endTime) }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">支付时间</text>
				<text class="info-value">{{ formatTime(order.reqDate) }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">下单时间</text>
				<text class="info-value">{{ formatTime(order.createdAt) }}</text>
			</view>
		</view>

		<!-- 设备信息 -->
		<view class="info-card">
			<view class="card-title">
				<uv-icon name="setting" size="20" color="#3c9cff" />
				<view>设备信息</view>
			</view>
			<view class="info-item">
				<text class="info-label">设备名称</text>
				<text class="info-value">{{ order.deviceName || '暂无' }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">房间名称</text>
				<text class="info-value">{{ order.roomName || '暂无' }}</text>
			</view>
		</view>

		<!-- 金额信息 -->
		<view class="info-card">
			<view class="card-title">
				<uv-icon name="shopping-cart" size="20" color="#3c9cff" />
				<view>金额信息</view>
			</view>
			<view class="info-item">
				<text class="info-label">支付金额</text>
				<text class="info-value amount">¥{{ formatMoney(order.price) }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">押金</text>
				<text class="info-value amount">¥{{ formatMoney(order.deposit) }}</text>
			</view>
			<view class="info-item">
				<text class="info-label">数量</text>
				<text class="info-value">{{ order.amount || 0 }}</text>
			</view>
		</view>

		<!-- 操作按钮 -->
		<view class="action-bar">
			<uv-button text="联系用户" type="primary" @click="contactUser" />
			<uv-button text="确认退款" type="warning" @click="confirmRefund" v-if="order.status === '申请退款'" />
			<uv-button text="取消订单" type="error" @click="cancelOrder" v-if="order.status === '未完成'" />
			<uv-button text="拒绝退款" type="error" @click="rejectRefund" v-if="order.status === '申请退款'" />
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue';
import { onLoad } from '@dcloudio/uni-app';

// 订单数据
const order = ref({});

// 获取订单详情
onLoad(() => {
	loadOrderDetail();
});

const loadOrderDetail = async () => {
	try {
		const pages = getCurrentPages();
		const currentPage = pages[pages.length - 1];
		const options = currentPage.$page?.options || {};
		const orderId = options.id || 1;
		
		const res = await uni.$uv.http.get(`/order/${orderId}`, {
			custom: { auth: true }
		});
		
		if (res.code === 200 && res.data) {
			order.value = res.data;
		}
	} catch (e) {
		console.error('加载订单详情失败', e);
	}
};

// 格式化金额
const formatMoney = (amount) => {
	return parseFloat(amount || 0).toFixed(2);
};

// 格式化时间
const formatTime = (time) => {
	if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
};

// 格式化时长（分钟转小时/天）
const formatDuration = (minutes) => {
	if (!minutes || minutes <= 0) return '0分钟';
	if (minutes < 60) {
		return `${minutes} 分钟`;
	} else if (minutes < 1440) { // 24小时 = 1440分钟
		const hours = Math.floor(minutes / 60);
		const mins = minutes % 60;
		if (mins > 0) {
			return `${hours}小时${mins}分钟`;
		}
		return `${hours} 小时`;
	} else {
		const days = Math.floor(minutes / 1440);
		const hours = Math.floor((minutes % 1440) / 60);
		const mins = minutes % 60;
		if (hours > 0 && mins > 0) {
			return `${days}天${hours}小时${mins}分钟`;
		} else if (hours > 0) {
			return `${days}天${hours}小时`;
		} else if (mins > 0) {
			return `${days}天${mins}分钟`;
		}
		return `${days} 天`;
	}
};

// 获取状态样式类
const getStatusClass = (status) => {
	switch (status) {
		case '已完成':
			return 'success';
		case '进行中':
			return 'primary';
		case '申请退款':
			return 'warning';
		case '已退款':
			return 'danger';
		case '拒绝退款':
			return 'info';
		default:
			return 'default';
	}
};

// 获取状态图标
const getStatusIcon = (status) => {
	switch (status) {
		case '已完成':
			return '✓';
		case '进行中':
			return '○';
		case '申请退款':
			return '?';
		case '已退款':
			return '↩';
		case '拒绝退款':
			return '✗';
		default:
			return '-';
	}
};

// 获取状态描述
const getStatusDesc = (status) => {
	switch (status) {
		case '未完成':
			return '订单尚未完成';
		case '进行中':
			return '订单正在进行中';
		case '已完成':
			return '订单已完成';
		case '申请退款':
			return '用户已申请退款，请处理';
		case '已退款':
			return '退款已完成';
		case '拒绝退款':
			return '退款已拒绝';
		default:
			return '未知状态';
	}
};

// 取消订单
const cancelOrder = () => {
	uni.showModal({
		title: '确认取消',
		content: '确定要取消该订单吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					await uni.$uv.http.put(`/order/${order.value.id}`, {
						status: 'cancelled'
					}, {
						custom: { auth: true }
					});
					uni.showToast({
						title: '订单已取消',
						icon: 'success'
					});
					setTimeout(() => {
						uni.navigateBack();
					}, 1500);
				} catch (e) {
					uni.showToast({
						title: '取消失败',
						icon: 'error'
					});
				}
			}
		}
	});
};

// 确认退款
const confirmRefund = () => {
	uni.showModal({
		title: '确认退款',
		content: `确定要退款 ¥${formatMoney(order.value.price)} 给用户吗？`,
		success: async (res) => {
			if (res.confirm) {
				try {
					await uni.$uv.http.put(`/order/${order.value.id}`, {
						status: '4'
					}, {
						custom: { auth: true }
					});
					uni.showToast({
						title: '退款成功',
						icon: 'success'
					});
					setTimeout(() => {
						uni.navigateBack();
					}, 1500);
				} catch (e) {
					uni.showToast({
						title: '退款失败',
						icon: 'error'
					});
				}
			}
		}
	});
};

// 拒绝退款
const rejectRefund = () => {
	uni.showModal({
		title: '拒绝退款',
		content: '确定要拒绝用户的退款申请吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					await uni.$uv.http.put(`/order/${order.value.id}`, {
						status: '5'
					}, {
						custom: { auth: true }
					});
					uni.showToast({
						title: '已拒绝退款',
						icon: 'success'
					});
					setTimeout(() => {
						uni.navigateBack();
					}, 1500);
				} catch (e) {
					uni.showToast({
						title: '操作失败',
						icon: 'error'
					});
				}
			}
		}
	});
};

// 联系用户
const contactUser = () => {
	if (!order.value.userPhone) {
		uni.showToast({
			title: '用户未绑定手机号',
			icon: 'none'
		});
		return;
	}
	uni.makePhoneCall({
		phoneNumber: order.value.userPhone,
		success: () => {
			console.log('拨打电话成功');
		},
		fail: () => {
			uni.showToast({
				title: '拨打电话失败',
				icon: 'error'
			});
		}
	});
};

// 返回上一页
const goBack = () => {
	uni.navigateBack();
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 160rpx;
}

.status-card {
	display: flex;
	align-items: center;
	padding: 32rpx;
	margin: 24rpx;
	border-radius: 16rpx;
	
	&.success {
		background: linear-gradient(135deg, #e8f5e9 0%, #c8e6c9 100%);
		
		.status-icon, .status-text {
			color: #4caf50;
		}
	}
	
	&.warning {
		background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
		
		.status-icon, .status-text {
			color: #ff9800;
		}
	}
	
	&.danger {
		background: linear-gradient(135deg, #ffebee 0%, #ffcdd2 100%);
		
		.status-icon, .status-text {
			color: #f44336;
		}
	}
	
	&.default {
		background: linear-gradient(135deg, #f5f5f5 0%, #e0e0e0 100%);
		
		.status-icon, .status-text {
			color: #999;
		}
	}

	&.primary {
		background: linear-gradient(135deg, #a2d0ff 0%, #80bfff 100%);
		
		.status-icon, .status-text {
			color: #007aff;
		}
	}
	
	&.info {
		background: linear-gradient(135deg, #fbf0ff 0%, #fff0 100%);
		
		.status-icon, .status-text {
			color: #c10781;
		}
	}
}

.status-icon {
	width: 80rpx;
	height: 80rpx;
	border-radius: 50%;
	background: rgba(255, 255, 255, 0.8);
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 40rpx;
	font-weight: bold;
	margin-right: 20rpx;
}

.status-info {
	display: flex;
	flex-direction: column;
}

.status-text {
	font-size: 32rpx;
	font-weight: bold;
}

.status-desc {
	font-size: 24rpx;
	color: #666;
	margin-top: 4rpx;
}

.info-card {
	background: #fff;
	border-radius: 16rpx;
	margin: 0 24rpx 24rpx;
	padding: 32rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.card-title {
	display: flex;
	align-items: center;
	gap: 10rpx;
	font-size: 28rpx;
	font-weight: 600;
	color: #333;
	margin-bottom: 20rpx;
	padding-bottom: 16rpx;
	border-bottom: 1rpx solid #f5f5f5;
}

.info-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 0;
	
	&:not(:last-child) {
		border-bottom: 1rpx solid #fafafa;
	}
}

.info-label {
	font-size: 28rpx;
	color: #999;
}

.info-value {
	font-size: 28rpx;
	color: #333;
	
	&.amount {
		font-size: 32rpx;
		font-weight: bold;
		color: red;
	}
}

.action-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	display: flex;
	gap: 20rpx;
	padding: 20rpx 24rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	background: #fff;
	box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.04);
	
	button {
		flex: 1;
		height: 88rpx;
	}
}
</style>