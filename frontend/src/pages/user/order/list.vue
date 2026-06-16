<template>
	<view class="container">
		<uv-navbar :title="t('user.order.title')" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

		<!-- 订单状态切换 Tab -->
		<view class="tabs-wrapper">
			<uv-tabs :list="tabs" lineColor="#3c9cff" lineWidth="40" lineHeight="6" @change="onTabChange"></uv-tabs>
		</view>
		<!-- 订单列表 -->
		<view class="section">
			<view class="order-list">
				<view class="order-item" v-for="order in orderList" :key="order.id" @click="viewDetail(order)">
					<view class="order-left">
						<view class="order-icon" :class="getStatusClass(order.status)">
							<text>{{ getStatusIcon(order.status) }}</text>
						</view>
					</view>
					<view class="order-center">
						<view class="order-name">{{ order.name || '订单 #' + order.code }}</view>
						<view class="order-time">{{ formatTime(order.createdAt) }}</view>
					</view>
					<view class="order-right">
						<text class="order-amount expense">¥{{ formatMoney(order.price) }}</text>
						<text class="order-status" :class="getStatusClass(order.status)">{{ order.status }}</text>
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
	</view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { onShow } from '@dcloudio/uni-app';

// 页面加载
onShow(() => {
	currentPage.value = 1;
	currentStatus.value = '全部';
	loadOrders();
});

// 国际化
const { t } = useI18n();
// 订单状态切换 Tab
const tabs = ref([
	{ name: '全部' },
	{ name: '未完成' },
	{ name: '进行中' },
	{ name: '已完成' },
	{ name: '申请退款' },
	{ name: '已退款' },
	{ name: '拒绝退款' },
]);
// 订单数据
const orderList = ref([]);
// 分页参数
const currentPage = ref(1);
// 每页数量
const pageSize = ref(20);
// 总订单数
const total = ref(0);
// 当前状态
const currentStatus = ref('全部');
// 跳转页码
const jumpPage = ref('');

// 加载订单列表
const loadOrders = async () => {
	try {
		const params = {
			page: currentPage.value,
			size: pageSize,
			status: currentStatus.value
		};

		const res = await uni.$uv.http.get('/user/order/list', {
			params,
			custom: { auth: true }
		});

		if (res.code === 200 && res.data) {
			orderList.value = res.data.list || [];
			total.value = res.data.total || 0;
		}
	} catch (e) {
		console.log('加载订单失败', e);
	}
};

// Tab切换
const onTabChange = (item) => {
	currentStatus.value = item.name;
	currentPage.value = 1;
	loadOrders();
};

// 查看详情
const viewDetail = (order) => {
	uni.navigateTo({
		url: `/pages/user/order/detail?id=${order.id}`
	});
};

// 上一页
const handlePrevPage = () => {
	if (currentPage.value > 1) {
		currentPage.value--;
		loadOrders();
	}
};

// 下一页
const handleNextPage = () => {
	if (currentPage.value < totalPages.value) {
		currentPage.value++;
		loadOrders();
	}
};

// 跳转页面
const handleJumpPage = () => {
	const page = parseInt(jumpPage.value);
	if (!isNaN(page) && page >= 1 && page <= totalPages.value) {
		currentPage.value = page;
		loadOrders();
	}
	jumpPage.value = '';
};

// 总页数
const totalPages = computed(() => {
	return Math.ceil(total.value / pageSize);
});

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
		case '申请退款':
			return 'warning';
		case '已退款':
			return 'danger';
		case '拒绝退款':
			return 'danger';
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
		default:
			return '○';
	}
};

// 返回上一页
const goBack = () => {
	uni.navigateBack({ delta: 1 });
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f8f9fa;
	padding-bottom: 160rpx;
}

.tabs-wrapper {
	background: #fff;
	padding-top: 10rpx;
}

.section {
	padding: 0 20rpx;
	margin-top: 16rpx;
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