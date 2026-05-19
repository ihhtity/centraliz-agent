<template>
	<view class="container">
		<!-- 搜索栏 -->
		<view class="search-bar">
			<uv-icon name="search" color="#999" size="24" />
			<input
				v-model="searchQuery"
				type="text"
				placeholder="搜索房间..."
			/>
		</view>

		<!-- 功能按钮栏 -->
		<view class="action-bar">
			<view class="action-btn group-btn" @click="handleGroup">
				<uv-icon name="grid" color="#fff" size="20" />
				<text class="btn-text">分组</text>
			</view>
			<view class="action-btn add-btn" @click="handleAdd">
				<uv-icon name="edit-pen" color="#fff" size="20" />
				<text class="btn-text">添加</text>
			</view>
			<view class="action-btn refresh-btn" @click="handleRefresh">
				<uv-icon name="reload" color="#fff" size="20" />
				<text class="btn-text">刷新</text>
			</view>
		</view>

		<!-- 房间分组标签栏 -->
		<view style="margin-bottom: 24rpx;">
			<uv-tabs :list="roomGroups" @click="selectGroup" />
		</view>

		<!-- 房间列表 -->
		<view class="room-grid">
			<view
				v-for="(room, index) in filteredRooms"
				:key="index"
				class="room-card"
				:class="{
					'empty': room.status === '空闲',
					'rented': room.status === '租用',
					'maintenance': room.status === '维修'
				}"
			>
				<!-- 状态标签 -->
				<view class="status-badge" :class="getStatusClass(room.status)">
					{{ getStatusText(room.status) }}
				</view>
				
				<!-- 房间信息 -->
				<view class="room-info">
					<text class="room-number">{{ room.id }}</text>
					<text class="room-type">{{ room.type }}</text>
				</view>
			</view>
		</view>

		<!-- 底部统计 -->
		<view class="stats-section-fixed">
			<view class="stat-row">
				<view 
					class="stat-item" 
					:class="{ active: !filterStatus }"
					@click="setFilter('')"
				>
					<view class="color-indicator all"></view>
					<text class="stat-label">全部</text>
					<text class="stat-value">{{ totalRooms }}</text>
				</view>
				<view 
					class="stat-item" 
					:class="{ active: filterStatus === '空闲' }"
					@click="setFilter('空闲')"
				>
					<view class="color-indicator blue"></view>
					<text class="stat-label">空闲</text>
					<text class="stat-value">{{ stats.empty }}</text>
				</view>
				<view 
					class="stat-item" 
					:class="{ active: filterStatus === '租用' }"
					@click="setFilter('租用')"
				>
					<view class="color-indicator red"></view>
					<text class="stat-label">租用</text>
					<text class="stat-value">{{ stats.rented }}</text>
				</view>
				<view 
					class="stat-item" 
					:class="{ active: filterStatus === '维修' }"
					@click="setFilter('维修')"
				>
					<view class="color-indicator orange"></view>
					<text class="stat-label">维修</text>
					<text class="stat-value">{{ stats.maintenance }}</text>
				</view>
			</view>
		</view>

		<!-- 底部导航栏 -->
		<uv-tabbar :value="tabbar" @change="editTabbar">
			<uv-tabbar-item text="房间" icon="home" />
			<uv-tabbar-item text="个人" icon="account" />
		</uv-tabbar>
	</view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
// 模拟房间分组数据
const roomGroups = ref([
	{ id: '1', name: '关注' },
	{ id: '2', name: '推荐' },
	{ id: '3', name: '电影' },
	{ id: '4', name: '游戏' },
	{ id: '5', name: '体育' },
	{ id: '6', name: '音乐' },
	{ id: '7', name: '科技' },
	{ id: '8', name: '教育' }
]);
// 模拟房间数据
const rooms = ref([
	{ id: 'A01', type: '大柜', status: '空闲' },
	{ id: 'A02', type: '中柜', status: '租用' },
	{ id: 'A03', type: '小柜', status: '维修' },
	{ id: 'A04', type: '大柜', status: '空闲' },
	{ id: 'B01', type: '中柜', status: '空闲' },
	{ id: 'B02', type: '小柜', status: '租用' },
	{ id: 'B03', type: '大柜', status: '空闲' },
	{ id: 'B04', type: '中柜', status: '维修' },
	{ id: 'C01', type: '小柜', status: '维修' },
	{ id: 'C02', type: '大柜', status: '空闲' },
	{ id: 'C03', type: '中柜', status: '租用' },
	{ id: 'C04', type: '小柜', status: '维修' }
]);
// 筛选状态
const filterStatus = ref('');
// 全部房间总数计算
const totalRooms = computed(() => rooms.value.length);

// 搜索关键词
const searchQuery = ref('');

const handleGroup = () => {
	console.log('点击分组');
};

const handleAdd = () => {
	console.log('点击添加');
};

const handleRefresh = () => {
	// 模拟请求后端数据
	uni.showLoading({ title: '加载中...' });
	
	setTimeout(() => {
		// 这里可以替换为实际的API请求
		// rooms.value = newDataFromAPI;
		
		// 重置为全部数据
		filterStatus.value = '';
		uni.hideLoading();
		uni.showToast({ title: '刷新成功', icon: 'success' });
	}, 800);
};

// 过滤房间逻辑
const filteredRooms = computed(() => {
	let result = rooms.value;
	
	// 应用搜索过滤
	if (searchQuery.value) {
		result = result.filter(room =>
			room.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
			room.type.includes(searchQuery.value) ||
			room.status.includes(searchQuery.value)
		);
	}
	
	// 应用状态过滤
	if (filterStatus.value) {
		result = result.filter(room => room.status === filterStatus.value);
	}
	
	return result;
});

// 设置筛选状态
const setFilter = (status) => {
	if (filterStatus.value === status) {
		filterStatus.value = ''; // 取消筛选
	} else {
		filterStatus.value = status;
	}
};

// 获取状态文本
const getStatusText = (status) => {
	switch(status) {
		case '空闲': return '空';
		case '租用': return '租';
		case '维修': return '维';
		default: return status;
	}
};

// 获取状态类名
const getStatusClass = (status) => {
	return {
		'empty': status === '空闲',
		'rented': status === '租用',
		'maintenance': status === '维修'
	};
};


// 统计状态
const stats = computed(() => {
	const empty = rooms.value.filter(r => r.status === '空闲').length;
	const rented = rooms.value.filter(r => r.status === '租用').length;
	const maintenance = rooms.value.filter(r => r.status === '维修').length;
	return { empty, rented, maintenance };
});

// 导航逻辑
const tabbar = ref(0);
// 点击 tab 切换页面
const editTabbar = (e) => {
	tabbar.value = e;
	if (e === 1) {
		uni.navigateTo({
			url: '/pages/admin/profile/index'
		});
	}
};

// 选择房间分组
const selectGroup = (e) => {
	console.log(e);
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding: 24rpx;
	padding-bottom: 80rpx;
}

.search-bar {
	position: relative;
	width: 92%;
	height: 80rpx;
	background: #fff;
	border-radius: 40rpx;
	box-shadow: 0 4rpx 16rpx rgba(0,0,0,0.04);
	overflow: hidden;
	display: flex;
	align-items: center;
	padding: 0 20rpx;
	margin-bottom: 24rpx;
	
	input {
		flex: 1;
		font-size: 28rpx;
		color: #333;
		outline: none;
		background: transparent;
	}
	
	.uv-icon {
		margin-right: 16rpx;
	}
}

.action-bar {
	display: flex;
	gap: 16rpx;
	
	.action-btn {
		flex: 1;
		display: flex;
    	flex-direction: row;
		justify-content: center;
		padding: 12rpx 16rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		color: #fff;
		text-align: center;
		cursor: pointer;
		
		.uv-icon {
			margin-bottom: 4rpx;
		}
		
		.btn-text {
			margin-left: 8rpx;
			font-size: 24rpx;
			font-weight: 500;
		}
	}
	
	.group-btn {
		background: linear-gradient(135deg, #4CAF50, #388E3C);
	}
	
	.count-btn {
		background: linear-gradient(135deg, #9E9E9E, #757575);
	}
	
	.add-btn {
		background: linear-gradient(135deg, #F44336, #D32F2F);
	}
	
	.refresh-btn {
		background: linear-gradient(135deg, #FF9800, #E67C00);
	}
}

.room-grid {
	display: grid;
	grid-template-columns: repeat(4, 1fr); // 改为4列
	gap: 20rpx; // 减小间距适应更多卡片
	margin-bottom: 30rpx;
}

.room-card {
	background: #ffffff;
	border-radius: 16rpx;
	padding: 20rpx 16rpx; // 调整内边距
	box-shadow: 0 6rpx 16rpx rgba(0, 0, 0, 0.08);
	position: relative;
	min-height: 140rpx;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	transition: all 0.2s ease;

	&:active {
		transform: translateY(4rpx);
		box-shadow: 0 8rpx 20rpx rgba(0, 0, 0, 0.12);
	}
	
	&.empty { background: #3c9cff }
	&.rented { background: #fa3534; }
	&.maintenance { background: #ff9900; }
	
	.status-badge {
		position: absolute;
		top: -10rpx;
		right: -10rpx;
		width: 44rpx;
		height: 44rpx;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24rpx;
		font-weight: bold;
		box-shadow: 0 4rpx 10rpx rgba(0,0,0,0.2);
		z-index: 2;
		
		&.empty { background: #3c9cff; color: white; }
		&.rented { background: #fa3534; color: white; }
		&.maintenance { background: #ff9900; color: white; }
	}
	
	.room-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8rpx;
		text-align: center;
	}
	
	.room-number {
		font-size: 28rpx;
		font-weight: 700; // 更粗的字体
		color: #fff; // 更深的颜色
		line-height: 1.2;
	}
	
	.room-type {
		font-size: 20rpx;
		color: #fff;
		padding: 2rpx 12rpx;
	}
}

.stats-section-fixed {
	position: fixed;
	bottom: 100rpx; // 在 tabbar 上方
	left: 0;
	right: 0;
	padding: 20rpx 0rpx;
	border-radius: 20rpx;
	background: #fff;
	box-shadow: 0 -2rpx 16rpx rgba(0,0,0,0.1);
	z-index: 99;
	width: 94%;
	margin: 0 auto;
	
	.stat-row {
		display: flex;
		justify-content: space-around;
		align-items: center;
	}
	
	.stat-item {
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: 8rpx 16rpx;
		border-radius: 8rpx;
		transition: all 0.2s;
		
		&:active {
			transform: scale(0.95);
		}
		
		&.active {
			background: rgba(0,0,0,0.05);
			box-shadow: 0 2rpx 8rpx rgba(0,0,0,0.1);
		}
		
		.color-indicator {
			width: 24rpx;
			height: 24rpx;
			border-radius: 4rpx;
			margin-right: 8rpx;
			
			&.all { background: #5ac725; }
			&.blue { background: linear-gradient(135deg, #3c9cff, #2b85e4); }
			&.red { background: linear-gradient(135deg, #fa3534, #e63332); }
			&.orange { background: linear-gradient(135deg, #ff9900, #f29100); }
		}
		
		.stat-label {
			font-size: 20rpx;
			color: #666;
			margin-right: 8rpx;
		}
		
		.stat-value {
			font-size: 28rpx;
			font-weight: bold;
			color: #333;
		}
	}
}
</style>