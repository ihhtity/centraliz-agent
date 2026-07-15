<!-- 交易明细页面 -->
<template>
	<view class="container">
		<uv-navbar :title="'交易明细'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

		<!-- 顶部统计卡片 -->
		<view class="stats-card">
			<view class="card-header stats-line">
				<text class="card-title">交易统计</text>
				<view class="date-picker" @click="showCalendar">
					<text>{{ dateRangeText }}</text>
					<uv-icon name="calendar" size="24" />
				</view>
			</view>
			<view class="stats-row stats-line">
				<view class="stat-item">
					<text class="stat-value income">¥{{ formatMoney(totalIncome) }}</text>
					<text class="stat-label">收入</text>
				</view>
				<view class="stat-divider"></view>
				<view class="stat-item">
					<text class="stat-value expense">-¥{{ formatMoney(totalExpense) }}</text>
					<text class="stat-label">退款</text>
				</view>
			</view>
			<view class="stats-row">
				<view class="stat-item" :class="{ disabled: depositCount < 1 }" @click="getDeposit">
					<text class="stat-value deposit">¥{{ formatMoney(totalDeposit) }}</text>
					<text class="stat-label">押金</text>
				</view>
				<view class="stat-divider"></view>
				<view class="stat-item" :class="{ disabled: refundCount < 1 }" @click="getRefund">
					<text class="stat-value refund">{{ refundCount }}</text>
					<text class="stat-label">申请退款</text>
				</view>
			</view>
		</view>

		<!-- 订单列表 -->
		<view class="section">
			<view class="section-header">
				<text class="section-title">订单记录</text>
				<text class="section-count">共 {{ total }} 条</text>
			</view>

			<view class="order-list">
				<view class="order-item" v-for="order in orderList" :key="order.id" @click="goToDetail(order.id)">
					<view class="order-left">
						<view class="order-icon" :class="getStatusClass(order.status)">
							<text>{{ getStatusIcon(order.status) }}</text>
						</view>
					</view>
					<view class="order-center">
						<view class="order-name">{{ order.name || '订单 #' + order.id }}</view>
						<view class="order-time">{{ formatTime(order.createdAt) }}</view>
					</view>
					<view class="order-right">
						<text class="order-amount" :class="order.status === '已完成' ? 'income' : 'expense'">
							{{ order.status === '已完成' ? '+' : '-' }}¥{{ formatMoney(order.price ? order.price : order.deposit) }}
						</text>
						<text class="order-status" :class="getStatusClass(order.status)">{{ order.status
							}}</text>
					</view>
				</view>
			</view>

			<!-- 空状态 -->
			<uv-empty v-if="orderList.length === 0" mode="data" text="暂无订单记录" />
		</view>

		<!-- 分页组件 -->
		<view class="pagination-wrapper">
			<view class="pagination">
				<view class="page-btn prev-btn" :class="{ disabled: currentPage === 1 }" @click="handlePrevPage">
					<uv-icon name="arrow-left" size="16" color="#666" />
					<text>上一页</text>
				</view>

				<view class="page-info">
					<text class="current-page">{{ currentPage }}</text>
					<text class="separator">/</text>
					<text class="total-pages">{{ totalPages }}</text>
				</view>

				<view class="page-jump">
					<text>跳转</text>
					<input v-model="jumpPage" class="jump-input" type="number" placeholder="页码" @blur="handleJumpPage"
						@confirm="handleJumpPage" />
					<text>页</text>
				</view>

				<view class="page-btn next-btn" :class="{ disabled: currentPage === totalPages || totalPages === 0 }"
					@click="handleNextPage">
					<text>下一页</text>
					<uv-icon name="arrow-right" size="16" color="#666" />
				</view>
			</view>
		</view>

		<!-- 日期范围选择弹窗 -->
		<uv-calendars ref="calendar" :mode="'range'" :allowSameDay="true" :lunar="true" :endDate="newDate"
			closeOnClickOverlay @confirm="confirmCalendar" />
	</view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { onShow } from '@dcloudio/uni-app';

// 商户信息
const merch = ref({});

// 日期范围
const calendar = ref(null);
const startDate = ref('');
const endDate = ref('');
const newDate = ref('');

// 订单数据
const orderList = ref([]);
const currentPage = ref(1);
const pageSize = 20;
const total = ref(0);
const jumpPage = ref('');
const ordertype = ref(0); // 0: 所有订单, 1: 申请退款订单, 2: 押金订单

// 统计数据
const totalIncome = ref(0);
const totalExpense = ref(0);
const totalDeposit = ref(0);
const refundCount = ref(0);
const depositCount = ref(0);

// 页面加载
onShow(() => {
	merch.value = uni.getStorageSync('merch') || {};
	currentPage.value = 1;
	// 默认选择当天日期
	if (!startDate.value) {
		const today = getTodayStr();
		newDate.value = today;
		startDate.value = today;
	}
	if (ordertype.value === 0) {
		loadOrders();
	} else {
		getRefund();
	}
});

// 加载订单列表
const loadOrders = async () => {
	if (!startDate.value) {
		uni.showToast({
			title: '请选择日期',
			icon: 'none',
			duration: 3000
		});
		return;
	}

	try {
		uni.showLoading({ title: '加载中...' });
		
		ordertype.value = 0;
		const params = {
			merch_id: merch.value.id,
			page: currentPage.value,
			size: pageSize,
			start_date: startDate.value,
			end_date: endDate.value
		};
		// console.log("params:", params);return;

		const res = await uni.$uv.http.get('/order/list', {
			params,
			custom: { auth: true }
		});
		uni.hideLoading();

		if (res.code === 200 && res.data) {
			orderList.value = res.data.list || [];
			total.value = res.data.total || 0;
			totalIncome.value = res.data.income || 0;
			totalExpense.value = res.data.expense || 0;
			totalDeposit.value = res.data.deposit || 0;
			refundCount.value = res.data.refund || 0;
			depositCount.value = total.value;
		}
	} catch (e) {
		uni.hideLoading();
		console.log('加载订单失败', e);
	}
};

// 获取申请退款订单列表
const getRefund = async () => {
	try {
		uni.showLoading({ title: '加载中...' });

		ordertype.value = 1;
		const params = {
			merch_id: merch.value.id,
			page: currentPage.value,
			size: pageSize,
			status: "申请退款"
		};

		const res = await uni.$uv.http.get('/order/refund/list', {
			params,
			custom: { auth: true }
		});
		uni.hideLoading();

		if (res.code === 200 && res.data) {
			orderList.value = res.data.list || [];
			total.value = res.data.total || 0;
			totalIncome.value = 0;
			totalExpense.value = 0;
			totalDeposit.value = 0;
			refundCount.value = res.data.refund || 0;
		}
	} catch (e) {
		uni.hideLoading();
		console.log('加载申请退款订单失败', e);
	}
};

// 获取押金订单列表
const getDeposit = async () => {
	try {
		uni.showLoading({ title: '加载中...' });

		ordertype.value = 2;
		const params = {
			merch_id: merch.value.id,
			page: currentPage.value,
			size: pageSize
		};

		const res = await uni.$uv.http.get('/order/deposit/list', {
			params,
			custom: { auth: true }
		});
		uni.hideLoading();

		if (res.code === 200 && res.data) {
			orderList.value = res.data.list || [];
			total.value = res.data.total || 0;
			totalIncome.value = 0;
			totalExpense.value = 0;
			totalDeposit.value = res.data.deposit || 0;
			depositCount.value = res.data.deposit || 0;
		}
	} catch (e) {
		uni.hideLoading();
		console.log('加载押金订单失败', e);
	}
};

// 显示日期选择器
const showCalendar = () => {
	calendar.value.open();
};

// 确认日期选择
const confirmCalendar = (date) => {
	startDate.value = date['range']['before'];
	endDate.value = date['range']['after'];
	calendar.value.close();
	currentPage.value = 1;
	loadOrders();
};

// 上一页
const handlePrevPage = () => {
	if (currentPage.value > 1) {
		currentPage.value--;
		if (ordertype.value === 0) {
			loadOrders();
		} else if (ordertype.value === 1) {
			getRefund();
		} else if (ordertype.value === 2) {
			getDeposit();
		}
	}
};

// 下一页
const handleNextPage = () => {
	if (currentPage.value < totalPages.value) {
		currentPage.value++;
		if (ordertype.value === 0) {
			loadOrders();
		} else if (ordertype.value === 1) {
			getRefund();
		} else if (ordertype.value === 2) {
			getDeposit();
		}
	}
};

// 跳转页面
const handleJumpPage = () => {
	const page = parseInt(jumpPage.value);
	if (!isNaN(page) && page >= 1 && page <= totalPages.value) {
		currentPage.value = page;
		if (ordertype.value === 0) {
			loadOrders();
		} else if (ordertype.value === 1) {
			getRefund();
		} else if (ordertype.value === 2) {
			getDeposit();
		}
	}
	jumpPage.value = '';
};

// 总页数
const totalPages = computed(() => {
	return Math.ceil(total.value / pageSize);
});

// 日期范围显示文本
const dateRangeText = computed(() => {
	if (!startDate.value) {
		return '选择日期';
	}
	if (!endDate.value) {
		return '全部日期';
	}
	if (startDate.value === endDate.value) {
		return startDate.value;
	}
	return `${startDate.value} ~ ${endDate.value}`;
});

// 获取当天日期字符串
const getTodayStr = () => {
	const date = new Date();
	return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
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

// 跳转到订单详情
const goToDetail = (orderId) => {
	uni.navigateTo({
		url: `/pages/admin/order/detail?id=${orderId}`
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
	background-color: #f8f9fa;
	padding-bottom: 160rpx;
}

.stats-card {
	background: #fff;
	margin: 20rpx;
	border-radius: 12rpx;
	padding: 24rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.card-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.card-title {
	font-size: 28rpx;
	font-weight: 600;
	color: #333;
}

.date-picker {
	display: flex;
	align-items: center;
	gap: 8rpx;
	padding: 8rpx 16rpx;
	background: #f5f7fa;
	border-radius: 16rpx;
	font-size: 24rpx;
	color: #666;
}

.stats-row {
	display: flex;
	justify-content: space-around;
	align-items: center;
}

.stats-line {
	margin-bottom: 24rpx;
	padding-bottom: 16rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.stat-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	flex: 1;

	&.disabled {
		opacity: 0.4;
		pointer-events: none;
	}
}

.stat-value {
	font-size: 36rpx;
	font-weight: 600;

	&.income {
		color: #07c160;
	}

	&.expense {
		color: #ee0a24;
	}

	&.refund {
		color: #ff9500;
	}

	&.deposit {
		color: #3c9cff;
	}
}

.stat-label {
	font-size: 22rpx;
	color: #999;
	margin-top: 6rpx;
}

.stat-divider {
	width: 1rpx;
	height: 48rpx;
	background: #f0f0f0;
}

.section {
	padding: 0 20rpx;
}

.section-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16rpx;
}

.section-title {
	font-size: 26rpx;
	font-weight: 600;
	color: #333;
}

.section-count {
	font-size: 22rpx;
	color: #999;
}

.order-list {
	background: #fff;
	border-radius: 12rpx;
	overflow: hidden;
}

.order-item {
	display: flex;
	align-items: center;
	padding: 24rpx;
	border-bottom: 1rpx solid #f5f5f5;

	&:last-child {
		border-bottom: none;
	}

	&:active {
		background: #fafafa;
	}
}

.order-left {
	margin-right: 16rpx;
}

.order-icon {
	width: 48rpx;
	height: 48rpx;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 22rpx;
	background: #f5f7fa;
	color: #999;

	&.success {
		background: #e8fdf0;
		color: #07c160;
	}

	&.warning {
		background: #fff8f0;
		color: #ff9500;
	}

	&.danger {
		background: #fff1f0;
		color: #ee0a24;
	}

	&.primary {
		background: #a2d0ff;
		color: #007aff;
	}
		
	&.info {
		background: #fff1f0;
		color: #c10781;
	}
}

.order-center {
	flex: 1;
}

.order-name {
	font-size: 26rpx;
	color: #333;
	margin-bottom: 6rpx;
}

.order-time {
	font-size: 22rpx;
	color: #999;
}

.order-status {
	font-size: 20rpx;
	padding: 2rpx 8rpx;
	border-radius: 6rpx;

	&.success {
		background: #e8fdf0;
		color: #07c160;
	}

	&.warning {
		background: #fff8f0;
		color: #ff9500;
	}

	&.danger {
		background: #fff1f0;
		color: #ee0a24;
	}

	&.default {
		background: #f5f7fa;
		color: #999;
	}

	&.primary {
		background: #a2d0ff;
		color: #007aff;
	}
		
	&.info {
		background: #fbf0ff;
		color: #d800bd;
	}
}

.order-right {
	text-align: right;
	display: flex;
	flex-direction: column;
	align-items: flex-end;
	gap: 4rpx;
}

.order-amount {
	font-size: 28rpx;
	font-weight: 600;

	&.income {
		color: #07c160;
	}

	&.expense {
		color: #ee0a24;
	}
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 60rpx 0;
	background: #fff;
	border-radius: 12rpx;
}

.empty-icon {
	margin-bottom: 16rpx;
}

.empty-text {
	font-size: 24rpx;
	color: #999;
}

.pagination-wrapper {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fff;
	border-top: 1rpx solid #f0f0f0;
	padding: 16rpx 20rpx;
	padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
	z-index: 100;
}

.pagination {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 16rpx;
}

.page-btn {
	display: flex;
	align-items: center;
	gap: 6rpx;
	padding: 12rpx 20rpx;
	background: #f5f7fa;
	border-radius: 8rpx;
	font-size: 24rpx;
	color: #666;

	&:not(.disabled):active {
		background: #e8eaed;
	}

	&.disabled {
		opacity: 0.4;
		color: #ccc;
	}
}

.page-info {
	display: flex;
	align-items: center;
	gap: 6rpx;
	font-size: 24rpx;
	color: #666;

	.current-page {
		color: #333;
		font-weight: 600;
	}

	.separator {
		color: #ddd;
	}

	.total-pages {
		color: #999;
	}
}

.page-jump {
	display: flex;
	align-items: center;
	gap: 6rpx;
	font-size: 22rpx;
	color: #999;

	.jump-input {
		width: 80rpx;
		height: 48rpx;
		border: 1rpx solid #e8e8e8;
		border-radius: 6rpx;
		text-align: center;
		font-size: 22rpx;
	}
}
</style>