<template>
	<view class="container">
		<uv-navbar :title="t('user.order.detail')" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
		<!-- 存在订单详情时 -->
		<view v-if="orderDetail">
			<!-- 订单详情 -->
			<view class="content">
				<!-- 订单状态 -->
				<view class="status-card" :class="getStatusClass(orderDetail.status)">
					<view class="status-icon">
						<text>{{ getStatusIcon(orderDetail.status) }}</text>
					</view>
					<view class="status-info">
						<text class="status-text">{{ orderDetail.status }}</text>
						<text class="order-code">订单编号：{{ orderDetail.code }}</text>
					</view>
				</view>
				<!-- 订单信息 -->
				<view class="info-card">
					<view class="card-title">
						<uv-icon name="info-circle" size="20" color="#3c9cff" />
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
							<text class="info-label">押金</text>
							<text class="info-value price">¥{{ formatMoney(orderDetail.deposit) }}</text>
						</view>
						<view class="info-item" v-if="orderDetail.status === '已退款'">
							<text class="info-label">退款金额</text>
							<text class="info-value price">¥{{ formatMoney(orderDetail.refundPrice) }}</text>
						</view>
						<view class="info-item">
							<text class="info-label">使用时长</text>
							<text class="info-value">{{ formatDuration(orderDetail.duration) }}</text>
						</view>
						<view class="info-item">
							<text class="info-label">订单备注</text>
							<text class="info-value">{{ orderDetail.remark || '-' }}</text>
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
						<!-- <view class="info-item">
						<text class="info-label">预定日期</text>
						<text class="info-value">{{ orderDetail.reqDate || '-' }}</text>
					</view> -->
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
						<!-- <view class="info-item">
						<text class="info-label">用户手机号</text>
						<text class="info-value">{{ orderDetail.userPhone || '-' }}</text>
					</view> -->
						<view class="info-item">
							<text class="info-label">商家手机号</text>
							<text class="info-value">{{ orderDetail.merchPhone || '-' }}</text>
						</view>
					</view>
				</view>
			</view>
			<!-- 操作按钮 -->
		<view class="action-bar">
			<uv-button text="联系客服" type="primary" @click="contactCustomerService" />
			<uv-button text="申请退款" type="warning" @click="showRefundModal"
				v-if="orderDetail.reqSeqId && (orderDetail.status === '已完成' || orderDetail.tag === '押金' && orderDetail.status === '进行中')" />
		</view>
		</view>
		<!-- 不存在订单详情时 -->
		<view v-else style="margin-top: 200px;">
			<uv-empty mode="data" text="暂无订单详情" />
		</view>
		<!-- 申请退款弹框 -->
		<uv-modal ref="refundModal" title="申请原因" :show-cancel-button="true" confirm-text="提交" cancel-text="取消"
			@confirm="submitRefund" @cancel="closeRefundModal">
				<textarea v-model="refundRemark" placeholder="请输入退款原因" :maxlength="100"  />
		</uv-modal>
	</view>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { onLoad } from '@dcloudio/uni-app';

onLoad((options) => {
	if (options && options.id) {
		orderId.value = options.id;
		loadOrderDetail();
	}
});

// 国际化
const { t } = useI18n();
// 订单详情
const orderDetail = ref(null);
// 订单ID
const orderId = ref('');
// 退款弹框
const refundModal = ref(null);
const refundRemark = ref('');

// 加载订单详情
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
// 显示退款弹框
const showRefundModal = () => {
	refundRemark.value = '';
	refundModal.value.open();
};

// 关闭退款弹框
const closeRefundModal = () => {
	refundModal.value.close();
};

// 提交退款申请
const submitRefund = async () => {
	if (!refundRemark.value.trim()) {
		uni.showToast({ title: '请输入退款原因', icon: 'none' });
		return;
	}
	if (refundRemark.value.length < 1 || refundRemark.value.length > 100) {
		uni.showToast({ title: '退款原因字数必须在1-100之间', icon: 'none' });
		return;
	}

	try {
		const res = await uni.$uv.http.put(`/user/order/${orderId.value}/refund`, {
			remark: refundRemark.value.trim()
		}, {
			custom: { auth: true }
		});

		if (res.code === 200) {
			uni.showToast({ title: '申请成功', icon: 'success' });
			closeRefundModal();
			setTimeout(() => {
				goBack();
			}, 1500);
		} else {
			uni.showToast({ title: res.msg || '申请失败', icon: 'none' });
		}
	} catch (e) {
		uni.showToast({ title: '申请失败', icon: 'none' });
	}
};
// 联系客服
const contactCustomerService = () => {
	if (!orderDetail.value.merchPhone) {
		uni.showToast({ title: '暂无客服手机号', icon: 'none', duration: 2000 });
		return;
	}

	uni.makePhoneCall({
		phoneNumber: orderDetail.value.merchPhone,
		success: () => {
			console.log('拨打电话成功');
		},
		fail: () => {
			uni.showToast({ title: '拨打电话失败', icon: 'none', duration: 2000 });
		}
	});
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
		case '申请退款':
			return '?';
		case '已退款':
			return '↩';
		case '拒绝退款':
			return '✗';
		case '进行中':
			return '○';
		default:
			return '-';
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
	if (!minutes || minutes <= 0) return '0分钟'
	if (minutes < 60) {
		return `${minutes} 分钟`
	} else if (minutes < 1440) { // 24小时 = 1440分钟
		const hours = Math.floor(minutes / 60)
		const mins = minutes % 60
		if (mins > 0) {
			return `${hours}小时${mins}分钟`
		}
		return `${hours} 小时`
	} else {
		const days = Math.floor(minutes / 1440)
		const hours = Math.floor((minutes % 1440) / 60)
		const mins = minutes % 60
		if (hours > 0 && mins > 0) {
			return `${days}天${hours}小时${mins}分钟`
		} else if (hours > 0) {
			return `${days}天${hours}小时`
		} else if (mins > 0) {
			return `${days}天${mins}分钟`
		}
		return `${days} 天`
	}
};
// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/user/order/list'
	});
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

	&.info {
		border-left-color: #c10781;

		.status-icon {
			background: #fff1f0;
			color: #c10781;
		}

		.status-text {
			color: #c10781;
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

.refund-modal-content {
	padding: 20rpx 0;
}

.remark-input-wrapper {
	position: relative;
}

.remark-input {
	width: 100%;
	height: 200rpx;
	padding: 20rpx;
	background: #f8f9fa;
	border-radius: 12rpx;
	font-size: 28rpx;
	line-height: 1.6;
	box-sizing: border-box;
}

.remark-count {
	position: absolute;
	right: 20rpx;
	bottom: 20rpx;
	font-size: 22rpx;
	color: #999;
}
</style>