<!-- 操作记录页面 -->
<template>
	<view class="container">
		<uv-navbar title="操作记录" :placeholder="true" @leftClick="goBack" />

		<!-- 顶部统计卡片 -->
		<view class="stats-card">
			<view class="card-header">
				<text class="card-title">操作统计</text>
				<view class="date-picker" @click="showCalendar">
					<text>{{ dateRangeText }}</text>
					<uv-icon name="calendar" size="24" />
				</view>
			</view>
			<view class="stats-row">
				<view class="stat-item" :class="{ active: statusFilter === '成功' }" @click="filterByStatus('成功')">
					<text class="stat-value success">{{ successCount }}</text>
					<text class="stat-label">成功</text>
				</view>
				<view class="stat-divider"></view>
				<view class="stat-item" :class="{ active: statusFilter === '失败' }" @click="filterByStatus('失败')">
					<text class="stat-value fail">{{ failCount }}</text>
					<text class="stat-label">失败</text>
				</view>
				<view class="stat-divider"></view>
				<view class="stat-item" :class="{ active: statusFilter === '' }" @click="filterByStatus('')">
					<text class="stat-value total">{{ total }}</text>
					<text class="stat-label">总计</text>
				</view>
			</view>
		</view>

		<!-- 操作记录列表 -->
		<view class="list-wrapper">
			<view class="log-list">
				<view 
					v-for="log in filteredLogList" 
					:key="log.id" 
					class="log-item"
				>
					<view class="log-left">
						<view class="log-info">
							<view class="info-row">
								<text class="info-label">操作人：</text>
								<text class="info-value">{{ log.occupant || '用户' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">手机号：</text>
								<text class="info-value">{{ log.phone || '-' }}</text>
							</view>

							<view class="info-row">
								<text class="info-label">操作：</text>
								<text class="info-value">{{ log.control || '开锁' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">状态：</text>
								<text :class="['info-value', 'status-text', getStatusClass(log.status)]">{{ log.status || '失败' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">房间：</text>
								<text class="info-value">{{ log.roomName || '-' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">设备：</text>
								<text class="info-value">{{ log.deviceName || '-' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">编码：</text>
								<text class="info-value">{{ log.code || '-' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">使用设备：</text>
								<text class="info-value">{{ log.type || '手机' }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">操作时间：</text>
								<text class="info-value">{{ formatTime(log.createdAt) }}</text>
							</view>
						</view>
					</view>
				</view>
			</view>
			
			<uv-empty v-if="filteredLogList.length === 0" mode="list" text="暂无操作记录" />
		</view>

		<!-- 分页组件 -->
		<view class="pagination-wrapper">
			<view class="pagination">
				<view 
					class="page-btn prev-btn" 
					:class="{ disabled: currentPage === 1 }"
					@click="handlePrevPage"
				>
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
					<input 
						v-model="jumpPage" 
						class="jump-input" 
						type="number"
						placeholder="页码"
						@blur="handleJumpPage"
						@confirm="handleJumpPage"
					/>
					<text>页</text>
				</view>
				
				<view 
					class="page-btn next-btn" 
					:class="{ disabled: currentPage === totalPages || totalPages === 0 }"
					@click="handleNextPage"
				>
					<text>下一页</text>
					<uv-icon name="arrow-right" size="16" color="#666" />
				</view>
			</view>
		</view>

		<!-- 日期范围选择弹窗 -->
		<uv-calendars ref="calendar" :mode="'range'" :allowSameDay="true" :lunar="false" :endDate="newDate"
			closeOnClickOverlay @confirm="confirmCalendar" />
	</view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

onLoad((options) => {
	if (options.roomId) roomId.value = options.roomId
	if (options.roomName) roomName.value = decodeURIComponent(options.roomName)
	
	// 默认选择当天日期
	const today = getTodayStr();
	newDate.value = today;
	startDate.value = today;
	
	fetchLogList()
})

// 日期范围
const calendar = ref(null);
const startDate = ref('');
const endDate = ref('');
const newDate = ref('');

// 统计数据
const successCount = ref(0);
const failCount = ref(0);

// 列表数据
const roomId = ref('')
const roomName = ref('')
const logList = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const jumpPage = ref('')
const statusFilter = ref('')

// 日期范围显示文本
const dateRangeText = computed(() => {
	if (!startDate.value) return '选择日期'
	if (!endDate.value || startDate.value === endDate.value) return startDate.value
	return `${startDate.value} ~ ${endDate.value}`
})

// 获取今天的日期字符串
const getTodayStr = () => {
	const date = new Date();
	const year = date.getFullYear();
	const month = String(date.getMonth() + 1).padStart(2, '0');
	const day = String(date.getDate()).padStart(2, '0');
	return `${year}-${month}-${day}`;
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
	fetchLogList();
};

// 按状态筛选
const filterByStatus = (status) => {
	statusFilter.value = status;
	currentPage.value = 1;
	fetchLogList();
};

// 计算总页数
const totalPages = computed(() => {
	if (total.value === 0) return 0;
	return Math.ceil(total.value / pageSize.value);
});

// 过滤后的日志列表
const filteredLogList = computed(() => {
	return logList.value;
});

// 上一页
const handlePrevPage = () => {
	if (currentPage.value > 1) {
		currentPage.value--;
		fetchLogList();
	}
};

// 下一页
const handleNextPage = () => {
	if (currentPage.value < totalPages.value) {
		currentPage.value++;
		fetchLogList();
	}
};

// 跳转到指定页
const handleJumpPage = () => {
	const page = parseInt(jumpPage.value);
	if (isNaN(page) || page < 1 || page > totalPages.value) {
		uni.showToast({ title: `请输入1-${totalPages.value}之间的页码`, icon: 'none' });
		jumpPage.value = '';
		return;
	}
	
	if (page === currentPage.value) {
		jumpPage.value = '';
		return;
	}
	
	currentPage.value = page;
	jumpPage.value = '';
	fetchLogList();
};

// 获取操作记录列表
const fetchLogList = () => {
	uni.showLoading({ title: '加载中' });
	loading.value = true;
	
	const merch = uni.getStorageSync('merch') || {};
	const merchsId = merch.id || merch.merchsId || '';
	
	const params = {
		merchs_id: merchsId,
		roomId: roomId.value,
		page: currentPage.value,
		pageSize: pageSize.value,
		startDate: startDate.value,
		endDate: endDate.value,
		status: statusFilter.value
	};
	
	uni.$uv.http.get('/device/log/list', {
		params,
		custom: { auth: true }
	}).then((res) => {
		logList.value = res.data.list || [];
		total.value = res.data.total || 0;
		// 计算统计数据
		calculateStats();
		uni.hideLoading();
		loading.value = false;
	}).catch((err) => {
		logList.value = [];
		total.value = 0;
		uni.hideLoading();
		loading.value = false;
		console.error('获取操作记录失败:', err);
	});
};

// 计算统计数据
const calculateStats = () => {
	successCount.value = 0;
	failCount.value = 0;
	
	logList.value.forEach(log => {
		if (log.status === '成功') {
			successCount.value++;
		} else if (log.status === '失败') {
			failCount.value++;
		}
	});
};

const goBack = () => uni.navigateBack()

// 格式化时间
const formatTime = (time) => {
	if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
};

// 获取状态样式
const getStatusClass = (status) => {
	switch (status) {
		case '成功':
			return 'online';
		case '失败':
			return 'offline';
		default:
			return 'offline';
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 140rpx;
}

// 统计卡片
.stats-card {
	background: linear-gradient(135deg, #2979ff 0%, #1989fa 100%);
	margin: 24rpx;
	padding: 30rpx;
	border-radius: 20rpx;
	box-shadow: 0 8rpx 24rpx rgba(41, 121, 255, 0.3);
	
	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 30rpx;
	}
	
	.card-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #fff;
	}
	
	.date-picker {
		display: flex;
		align-items: center;
		gap: 12rpx;
		padding: 12rpx 24rpx;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 32rpx;
		
		text {
			font-size: 24rpx;
			color: #fff;
		}
	}
	
	.stats-row {
		display: flex;
		justify-content: space-around;
		align-items: center;
	}
	
	.stat-item {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 8rpx;
		flex: 1;
		padding: 16rpx 24rpx;
		border-radius: 12rpx;
		transition: all 0.2s;
		
		&.active {
			background: rgba(255, 255, 255, 0.3);
		}
	}
	
	.stat-value {
		font-size: 40rpx;
		font-weight: 700;
		
		&.success {
			color: #4ade80;
		}
		
		&.fail {
			color: #f87171;
		}
		
		&.total {
			color: #fff;
		}
	}
	
	.stat-label {
		font-size: 22rpx;
		color: rgba(255, 255, 255, 0.8);
	}
	
	.stat-divider {
		width: 1rpx;
		height: 60rpx;
		background: rgba(255, 255, 255, 0.3);
	}
}

.list-wrapper {
	padding: 24rpx;
}

.log-list {
	display: flex;
	flex-direction: column;
}

.log-item {
	background: #fff;
	margin-bottom: 24rpx;
	border-radius: 20rpx;
	padding: 30rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.log-left {
	flex: 1;
	
	.log-info {
		display: flex;
		flex-direction: column;
		gap: 12rpx;
	}
	
	.info-row {
		display: flex;
		align-items: center;
	}
	
	.info-label {
		font-size: 26rpx;
		color: #999;
		flex-shrink: 0;
	}
	
	.info-value {
		font-size: 28rpx;
		color: #333;
		font-weight: 500;
		
		&.status-text {
			&.online {
				color: #19be6b;
			}
			
			&.offline {
				color: #f56c6c;
			}
		}
	}
}

.pagination-wrapper {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fff;
	box-shadow: 0 -4rpx 12rpx rgba(0, 0, 0, 0.05);
	padding: 20rpx 24rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	z-index: 100;
}

.pagination {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 20rpx;
}

.page-btn {
	display: flex;
	align-items: center;
	gap: 8rpx;
	padding: 16rpx 24rpx;
	background: #f5f7fa;
	border-radius: 12rpx;
	font-size: 26rpx;
	color: #333;
	transition: all 0.2s;
	
	&:not(.disabled):active {
		background: #e8e8e8;
	}
	
	&.disabled {
		opacity: 0.4;
		color: #999;
	}
}

.page-info {
	display: flex;
	align-items: center;
	gap: 8rpx;
	font-size: 26rpx;
	color: #666;
	
	.current-page {
		color: #3c9cff;
		font-weight: 600;
		font-size: 28rpx;
	}
	
	.separator {
		color: #999;
	}
	
	.total-pages {
		color: #999;
	}
}

.page-jump {
	display: flex;
	align-items: center;
	gap: 12rpx;
	font-size: 26rpx;
	color: #666;
	
	.jump-input {
		width: 100rpx;
		height: 60rpx;
		padding: 0 12rpx;
		background: #f5f7fa;
		border-radius: 8rpx;
		font-size: 26rpx;
		color: #333;
		text-align: center;
	}
}
</style>
