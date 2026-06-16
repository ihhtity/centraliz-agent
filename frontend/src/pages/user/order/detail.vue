<template>
	<view class="container">
		<uv-navbar :title="t('user.order.detail')" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

		<view class="content" v-if="orderDetail">
			<!-- 订单状态 -->
			<view class="status-card" :class="getStatusClass(orderDetail.status)">
				<view class="status-icon">
					<text>{{ getStatusIcon(orderDetail.status) }}</text>
				</view>
				<view class="status-info">
					<text class="status-text">{{ getStatusText(orderDetail.status) }}</text>
					<text class="order-code">订单编号：{{ orderDetail.code }}</text>
				</view>
			</view>

			<!-- 订单信息 -->
			<view class="info-card">
				<view class="card-title">
					<uv-icon name="info" size="20" color="#3c9cff" />
					<text>订单信息</text>
				</view>
				<view class="info-list">
					<view class="info-item">
						<text class="info-label">订单名称</text>
						<text class="info-value">{{ orderDetail.name || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">支付金额</text>
						<text class="info-value price">¥{{ formatMoney(orderDetail.price) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">商品数量</text>
						<text class="info-value">{{ orderDetail.amount }} 件</text>
					</view>
					<view class="info-item">
						<text class="info-label">使用时长</text>
						<text class="info-value">{{ orderDetail.duration }} 分钟</text>
					</view>
					<view class="info-item">
						<text class="info-label">押金</text>
						<text class="info-value">¥{{ formatMoney(orderDetail.deposit) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">下单时间</text>
						<text class="info-value">{{ formatTime(orderDetail.createdAt) }}</text>
					</view>
				</view>
			</view>

			<!-- 时间信息 -->
			<view class="info-card">
				<view class="card-title">
					<uv-icon name="clock" size="20" color="#3c9cff" />
					<text>时间信息</text>
				</view>
				<view class="info-list">
					<view class="info-item">
						<text class="info-label">预定日期</text>
						<text class="info-value">{{ orderDetail.reqDate || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">开始时间</text>
						<text class="info-value">{{ formatTime(orderDetail.startTime) }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">结束时间</text>
						<text class="info-value">{{ formatTime(orderDetail.endTime) }}</text>
					</view>
				</view>
			</view>

			<!-- 联系方式 -->
			<view class="info-card">
				<view class="card-title">
					<uv-icon name="phone" size="20" color="#3c9cff" />
					<text>联系方式</text>
				</view>
				<view class="info-list">
					<view class="info-item">
						<text class="info-label">用户手机号</text>
						<text class="info-value">{{ orderDetail.userPhone || '-' }}</text>
					</view>
					<view class="info-item">
						<text class="info-label">商家手机号</text>
						<text class="info-value">{{ orderDetail.merchPhone || '-' }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 底部操作按钮 -->
		<view class="bottom-btn" v-if="orderDetail && orderDetail.status === '申请退款'">
			<view class="btn-confirm" @click="handleCancelRefund">取消退款</view>
		</view>
		<view class="bottom-btn" v-else-if="orderDetail && orderDetail.status === '未完成'">
			<view class="btn-confirm" @click="handleApplyRefund">申请退款</view>
		</view>
		<view class="bottom-btn" v-else-if="orderDetail && orderDetail.status === '进行中'">
			<view class="btn-confirm" @click="handleComplete">完成订单</view>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { onLoad } from '@dcloudio/uni-app';

const { t } = useI18n();

const orderDetail = ref(null);
const orderId = ref('');

onLoad((options) => {
	if (options && options.id) {
		orderId.value = options.id;
		loadOrderDetail();
	}
});

const loadOrderDetail = async () => {
	try {
		const res = await uni.$uv.http.get(`/user/order/${orderId.value}`, {
			custom: { auth: true }
		});

		if (res.code === 200 && res.data) {
			orderDetail.value = res.data;
		}
	} catch (e) {
		console.log('加载订单详情失败', e);
	}
};

// 获取状态样式类
const getStatusClass = (status) => {
	switch (status) {
		case '已完成':
			return 'success';
		case '申请退款':
			return 'warning';
		case '已退款':
			return 'danger';
		case '拒绝退款':
			return 'danger';
		case '进行中':
			return 'primary';
		default:
			return 'default';
	}
};

// 获取状态图标
const getStatusIcon = (status) => {
	switch (status) {
		case '已完成':
			return '✓';
		case '申请退款':
			return '?';
		case '已退款':
			return '↩';
		case '拒绝退款':
			return '✗';
		case '进行中':
			return '○';
		default:
			return '○';
	}
};

// 获取状态文本
const getStatusText = (status) => {
	switch (status) {
		case '未完成':
			return '未完成';
		case '进行中':
			return '进行中';
		case '已完成':
			return '已完成';
		case '申请退款':
			return '申请退款';
		case '已退款':
			return '已退款';
		case '拒绝退款':
			return '拒绝退款';
		default:
			return '未知';
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

// 返回上一页
const goBack = () => {
	uni.navigateBack({ delta: 1 });
};

// 申请退款
const handleApplyRefund = async () => {
	try {
		const res = await uni.$uv.http.put(`/user/order/${orderId.value}/refund`, {}, {
			custom: { auth: true }
		});

		if (res.code === 200) {
			uni.showToast({ title: '申请成功', icon: 'success' });
			orderDetail.value.status = '申请退款';
		} else {
			uni.showToast({ title: res.msg || '申请失败', icon: 'none' });
		}
	} catch (e) {
		uni.showToast({ title: '申请失败', icon: 'none' });
	}
};

// 取消退款
const handleCancelRefund = async () => {
	try {
		const res = await uni.$uv.http.put(`/user/order/${orderId.value}/refund/cancel`, {}, {
			custom: { auth: true }
		});

		if (res.code === 200) {
			uni.showToast({ title: '取消成功', icon: 'success' });
			orderDetail.value.status = '未完成';
		} else {
			uni.showToast({ title: res.msg || '取消失败', icon: 'none' });
		}
	} catch (e) {
		uni.showToast({ title: '取消失败', icon: 'none' });
	}
};

// 完成订单
const handleComplete = async () => {
	try {
		const res = await uni.$uv.http.put(`/user/order/${orderId.value}/complete`, {}, {
			custom: { auth: true }
		});

		if (res.code === 200) {
			uni.showToast({ title: '完成成功', icon: 'success' });
			orderDetail.value.status = '已完成';
		} else {
			uni.showToast({ title: res.msg || '完成失败', icon: 'none' });
		}
	} catch (e) {
		uni.showToast({ title: '完成失败', icon: 'none' });
	}
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f8f9fa;
	padding-bottom: 160rpx;
}

.content {
	padding: 20rpx;
}

.status-card {
	display: flex;
	align-items: center;
	background: #fff;
	border-radius: 16rpx;
	padding: 32rpx;
	margin-bottom: 20rpx;
	border-left: 8rpx solid #e8e8e8;

	&.success {
		border-left-color: #07c160;
		.status-icon {
			background: #e8fdf0;
			color: #07c160;
		}
		.status-text {
			color: #07c160;
		}
	}

	&.primary {
		border-left-color: #3c9cff;
		.status-icon {
			background: #e8f0fe;
			color: #3c9cff;
		}
		.status-text {
			color: #3c9cff;
		}
	}

	&.warning {
		border-left-color: #ff9500;
		.status-icon {
			background: #fff8f0;
			color: #ff9500;
		}
		.status-text {
			color: #ff9500;
		}
	}

	&.danger {
		border-left-color: #ee0a24;
		.status-icon {
			background: #fff1f0;
			color: #ee0a24;
		}
		.status-text {
			color: #ee0a24;
		}
	}
}

.status-icon {
	width: 80rpx;
	height: 80rpx;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 36rpx;
	background: #f5f7fa;
	color: #999;
	margin-right: 20rpx;
}

.status-info {
	flex: 1;
}

.status-text {
	font-size: 32rpx;
	font-weight: 600;
	display: block;
	margin-bottom: 8rpx;
}

.order-code {
	font-size: 24rpx;
	color: #999;
}

.info-card {
	background: #fff;
	border-radius: 12rpx;
	padding: 24rpx;
	margin-bottom: 16rpx;
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

.info-list {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
}

.info-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.info-label {
	font-size: 26rpx;
	color: #999;
}

.info-value {
	font-size: 26rpx;
	color: #333;

	&.price {
		color: #ee0a24;
		font-weight: 600;
		font-size: 28rpx;
	}
}

.bottom-btn {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fff;
	border-top: 1rpx solid #f0f0f0;
	padding: 20rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	display: flex;
	gap: 20rpx;
}

.btn-confirm {
	flex: 1;
	height: 88rpx;
	background: linear-gradient(135deg, #3c9cff 0%, #36cbcb 100%);
	border-radius: 44rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #fff;
	font-size: 30rpx;
	font-weight: 500;

	&:active {
		opacity: 0.8;
	}
}
</style>