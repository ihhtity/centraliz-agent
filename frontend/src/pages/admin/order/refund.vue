<!-- 退款管理页面 -->
<template>
	<view class="container">
		<uv-navbar :title="'退款管理'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
		
		<!-- 统计卡片 -->
		<view class="stats-card">
			<view class="stat-item">
				<text class="stat-value">{{ pendingCount }}</text>
				<text class="stat-label">待处理</text>
			</view>
			<view class="stat-divider"></view>
			<view class="stat-item">
				<text class="stat-value success">{{ approvedCount }}</text>
				<text class="stat-label">已同意</text>
			</view>
			<view class="stat-divider"></view>
			<view class="stat-item">
				<text class="stat-value danger">{{ rejectedCount }}</text>
				<text class="stat-label">已拒绝</text>
			</view>
		</view>

		<!-- 筛选栏 -->
		<view class="filter-bar">
			<view 
				class="filter-item" 
				:class="{ active: statusFilter === '' }"
				@click="setStatusFilter('')"
			>
				全部
			</view>
			<view 
				class="filter-item" 
				:class="{ active: statusFilter === '3' }"
				@click="setStatusFilter('3')"
			>
				待处理
			</view>
			<view 
				class="filter-item" 
				:class="{ active: statusFilter === '4' }"
				@click="setStatusFilter('4')"
			>
				已同意
			</view>
			<view 
				class="filter-item" 
				:class="{ active: statusFilter === '5' }"
				@click="setStatusFilter('5')"
			>
				已拒绝
			</view>
		</view>

		<!-- 退款列表 -->
		<view class="refund-list">
			<view class="refund-item" v-for="item in refundList" :key="item.id">
				<view class="refund-header">
					<text class="refund-code">订单号: {{ item.code }}</text>
					<text class="refund-status" :class="getStatusClass(item.status)">
						{{ getStatusText(item.status) }}
					</text>
				</view>
				<view class="refund-info">
					<view class="info-row">
						<text class="info-label">退款金额</text>
						<text class="info-value amount">¥{{ formatMoney(item.price) }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">申请时间</text>
						<text class="info-value">{{ formatTime(item.createdAt) }}</text>
					</view>
					<view class="info-row">
						<text class="info-label">用户手机号</text>
						<text class="info-value">{{ item.userPhone || '暂无' }}</text>
					</view>
				</view>
				<view class="refund-actions" v-if="item.status === '3'">
					<uv-button text="拒绝" type="default" size="small" @click="rejectRefund(item.id)" />
					<uv-button text="同意" type="primary" size="small" @click="approveRefund(item.id)" />
				</view>
				<view class="refund-actions" v-else>
					<uv-button text="查看详情" type="primary" size="small" @click="goToDetail(item.id)" />
				</view>
			</view>
		</view>

		<!-- 空状态 -->
		<uv-empty v-if="refundList.length === 0" mode="data" text="暂无退款记录" />
	</view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { onShow } from '@dcloudio/uni-app';

// 筛选条件
const statusFilter = ref('');
const currentPage = ref(1);
const pageSize = 20;

// 退款列表
const refundList = ref([]);
const total = ref(0);

// 统计数据
const pendingCount = ref(0);
const approvedCount = ref(0);
const rejectedCount = ref(0);

// 总页数
const totalPages = computed(() => {
	return Math.ceil(total.value / pageSize);
});

// 页面加载
onShow(() => {
	currentPage.value = 1;
	loadRefunds();
});

// 加载退款列表
const loadRefunds = async () => {
	try {
		const params = {
			page: currentPage.value,
			size: pageSize,
			status: statusFilter.value
		};
		
		const res = await uni.$uv.http.get('/order/refund/list', {
			params,
			custom: { auth: true }
		});
		
		if (res.code === 200 && res.data) {
			refundList.value = res.data.list || [];
			total.value = res.data.total || 0;
			calculateStats();
		}
	} catch (e) {
		console.error('加载退款列表失败', e);
	}
};

// 计算统计数据
const calculateStats = () => {
	let pending = 0;
	let approved = 0;
	let rejected = 0;
	
	refundList.value.forEach(item => {
		switch (item.status) {
			case '3':
				pending++;
				break;
			case '4':
				approved++;
				break;
			case '5':
				rejected++;
				break;
		}
	});
	
	pendingCount.value = pending;
	approvedCount.value = approved;
	rejectedCount.value = rejected;
};

// 设置状态筛选
const setStatusFilter = (status) => {
	statusFilter.value = status;
	currentPage.value = 1;
	loadRefunds();
};

// 格式化金额
const formatMoney = (amount) => {
	return parseFloat(amount || 0).toFixed(2);
};

// 格式化时间
const formatTime = (time) => {
	if (!time) return '-';
	try {
		const date = new Date(time);
		return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
	} catch (e) {
		return time;
	}
};

// 获取状态样式类
const getStatusClass = (status) => {
	switch (status) {
		case '3':
			return 'warning';
		case '4':
			return 'success';
		case '5':
			return 'danger';
		default:
			return 'default';
	}
};

// 获取状态文本
const getStatusText = (status) => {
	switch (status) {
		case '3':
			return '待处理';
		case '4':
			return '已同意';
		case '5':
			return '已拒绝';
		default:
			return '未知';
	}
};

// 同意退款
const approveRefund = (orderId) => {
	uni.showModal({
		title: '确认同意退款',
		content: '确定要同意该退款申请吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					await uni.$uv.http.put(`/order/${orderId}/refund/approve`, {}, {
						custom: { auth: true }
					});
					uni.showToast({
						title: '已同意退款',
						icon: 'success'
					});
					loadRefunds();
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

// 拒绝退款
const rejectRefund = (orderId) => {
	uni.showModal({
		title: '确认拒绝退款',
		content: '确定要拒绝该退款申请吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					await uni.$uv.http.put(`/order/${orderId}/refund/reject`, {}, {
						custom: { auth: true }
					});
					uni.showToast({
						title: '已拒绝退款',
						icon: 'success'
					});
					loadRefunds();
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
	background-color: #f5f7fa;
	padding-bottom: 180rpx;
}

.stats-card {
	display: flex;
	justify-content: space-around;
	align-items: center;
	background: #fff;
	margin: 24rpx;
	border-radius: 16rpx;
	padding: 32rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.stat-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	flex: 1;
}

.stat-value {
	font-size: 40rpx;
	font-weight: bold;
	color: #333;
	
	&.success {
		color: #4caf50;
	}
	
	&.danger {
		color: #f44336;
	}
}

.stat-label {
	font-size: 24rpx;
	color: #999;
	margin-top: 8rpx;
}

.stat-divider {
	width: 1rpx;
	height: 60rpx;
	background: #f0f0f0;
}

.filter-bar {
	display: flex;
	background: #fff;
	margin: 0 24rpx 24rpx;
	border-radius: 16rpx;
	padding: 8rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.filter-item {
	flex: 1;
	text-align: center;
	padding: 20rpx;
	font-size: 28rpx;
	color: #666;
	border-radius: 12rpx;
	
	&.active {
		background: #3c9cff;
		color: #fff;
	}
}

.refund-list {
	padding: 0 24rpx;
}

.refund-item {
	background: #fff;
	border-radius: 16rpx;
	padding: 28rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.refund-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
	padding-bottom: 16rpx;
	border-bottom: 1rpx solid #f5f5f5;
}

.refund-code {
	font-size: 26rpx;
	color: #666;
}

.refund-status {
	font-size: 24rpx;
	padding: 6rpx 16rpx;
	border-radius: 8rpx;
	
	&.warning {
		background: #fff3e0;
		color: #ff9800;
	}
	
	&.success {
		background: #e8f5e9;
		color: #4caf50;
	}
	
	&.danger {
		background: #ffebee;
		color: #f44336;
	}
}

.refund-info {
	margin-bottom: 20rpx;
}

.info-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 12rpx 0;
}

.info-label {
	font-size: 26rpx;
	color: #999;
}

.info-value {
	font-size: 26rpx;
	color: #333;
	
	&.amount {
		font-size: 30rpx;
		font-weight: bold;
		color: #f44336;
	}
}

.refund-actions {
	display: flex;
	gap: 20rpx;
	
	button {
		flex: 1;
		height: 72rpx;
		font-size: 26rpx;
	}
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 100rpx 0;
	background: #fff;
	margin: 0 24rpx;
	border-radius: 16rpx;
}

.empty-icon {
	margin-bottom: 20rpx;
}

.empty-text {
	font-size: 28rpx;
	color: #999;
}

.pagination-wrapper {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fff;
	padding: 20rpx 24rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.04);
}
</style>