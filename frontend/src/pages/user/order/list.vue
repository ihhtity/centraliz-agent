<template>
	<view class="container">
		<uv-navbar :title="t('user.order.title')" :placeholder="true" :leftIcon="''" />
		
		<view class="tabs-wrapper">
			<uv-tabs :list="tabs" lineColor="#3c9cff" lineWidth="40" lineHeight="6" @change="onTabChange"></uv-tabs>
		</view>
		
		<view class="list-wrapper">
			<uv-list>
				<uv-list-item 
					v-for="item in orderList" 
					:key="item.id" 
					class="order-item"
					@click="viewDetail(item)"
				>
					<template #header>
						<view class="item-header">
							<text class="order-no">{{ t('user.order.orderNo') }}: {{ item.orderNo }}</text>
							<text :class="['status-badge', getStatusClass(item.status)]">{{ item.statusText }}</text>
						</view>
					</template>
					
					<template #default>
						<view class="item-body">
							<view class="info-row">
								<text class="label">{{ t('user.order.createTime') }}</text>
								<text class="value">{{ item.createTime }}</text>
							</view>
							<view class="info-row">
								<text class="label">{{ t('user.order.amount') }}</text>
								<text class="value price">¥{{ item.amount }}</text>
							</view>
						</view>
					</template>
				</uv-list-item>
			</uv-list>
			
			<uv-empty v-if="orderList.length === 0" mode="list" :text="t('user.order.noOrder')"></uv-empty>
		</view>

		<!-- 新增: 底部 TabBar -->
		<uv-tabbar :value="1" :placeholder="true" @change="onTabBarChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.order')" icon="order" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>
	</view>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const tabs = ref([
	{ name: t('user.order.all') || '全部' }, // 注意：如果 common.all 不存在，可能需要添加，这里暂时保留原样或添加 key
	{ name: t('user.order.pending') || '进行中' }, // 这里的逻辑可能需要调整，因为 tab 名字可能不完全对应订单状态
	{ name: t('user.order.completed') || '已完成' }
]);

// 修正 tabs 以匹配常见的订单状态翻译，或者直接使用硬编码如果翻译键不匹配
// 假设我们需要添加一个 'all' key 到 common，或者直接使用
// 为了简单，我们重新定义 tabs 使用具体的 key 如果存在，否则回退
const initTabs = () => {
    return [
        { name: '全部' }, // 暂时保留中文，除非添加 common.all
        { name: t('user.order.pending') }, // 待支付/进行中? 原文是进行中，pending是待支付。根据上下文，order list 通常是 全部/进行中/已完成
        { name: t('user.order.completed') }
    ]
}

const orderList = ref([
	{ id: 1, orderNo: 'ORD20231027001', createTime: '2023-10-27 10:00', amount: '5.00', status: 1, statusText: t('user.order.completed') },
	{ id: 2, orderNo: 'ORD20231027002', createTime: '2023-10-27 11:30', amount: '3.00', status: 0, statusText: t('user.order.pending') }, // 假设 0 是 pending/进行中
]);

const onTabChange = (item) => {
	console.log('切换标签', item);
};

const viewDetail = (item) => {
	// 查看详情逻辑
};

const getStatusClass = (status) => {
	return status === 1 ? 'status-success' : 'status-primary';
};

// 新增: TabBar 切换逻辑
const onTabBarChange = (index) => {
	if (index === 0) {
		uni.redirectTo({ url: '/pages/user/index/index' })
	} else if (index === 1) {
		// 当前已是订单页，无需操作
		return
	} else if (index === 2) {
		uni.redirectTo({ url: '/pages/user/profile/index' })
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.tabs-wrapper {
	background: #fff;
	padding-top: 10rpx;
}

.list-wrapper {
	padding: 24rpx;
}

.order-item {
	background: #fff;
	margin-bottom: 24rpx;
	border-radius: 20rpx;
	padding: 30rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
	
	.item-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20rpx;
		padding-bottom: 20rpx;
		border-bottom: 1rpx solid #f0f0f0;
		
		.order-no {
			font-size: 26rpx;
			color: #999;
		}
		
		.status-badge {
			font-size: 22rpx;
			padding: 4rpx 16rpx;
			border-radius: 20rpx;
			font-weight: 500;
		}
		
		.status-success {
			background-color: rgba(25, 190, 107, 0.1);
			color: #19be6b;
		}
		
		.status-primary {
			background-color: rgba(60, 156, 255, 0.1);
			color: #3c9cff;
		}
	}
	
	.item-body {
		.info-row {
			display: flex;
			justify-content: space-between;
			margin-bottom: 12rpx;
			
			&:last-child {
				margin-bottom: 0;
			}
			
			.label {
				font-size: 28rpx;
				color: #666;
			}
			
			.value {
				font-size: 28rpx;
				color: #333;
				
				&.price {
					color: #fa3534;
					font-weight: bold;
					font-size: 32rpx;
				}
			}
		}
	}
}
</style>