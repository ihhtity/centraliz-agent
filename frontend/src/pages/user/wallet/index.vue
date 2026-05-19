<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<uv-navbar :title="t('user.wallet.title')" :placeholder="true" @leftClick="goBack"></uv-navbar>
		
		<!-- 顶部资产卡片 -->
		<view class="asset-card">
			<view class="card-bg"></view>
			<view class="card-content">
				<text class="label">{{t('user.wallet.balance')}}</text>
				<view class="amount-row">
					<text class="symbol">¥</text>
					<text class="amount">1,280.50</text>
				</view>
				<view class="actions">
					<view class="action-btn" @click="handleRecharge">
						<uv-icon name="plus-circle" size="24" color="#fff"></uv-icon>
						<text>{{t('user.wallet.recharge')}}</text>
					</view>
					<view class="action-btn" @click="handleWithdraw">
						<uv-icon name="share-square" size="24" color="#fff"></uv-icon>
						<text>{{t('user.wallet.withdraw')}}</text>
					</view>
				</view>
			</view>
		</view>
		
		<!-- 收支明细 -->
		<view class="section-title">
			<text>{{t('user.wallet.transactionHistory')}}</text>
		</view>
		
		<view class="record-list">
			<view class="record-item" v-for="(item, index) in records" :key="index">
				<view class="left">
					<view class="icon-box" :class="item.type">
						<uv-icon :name="item.icon" size="20" color="#fff"></uv-icon>
					</view>
					<view class="info">
						<text class="title">{{ item.title }}</text>
						<text class="time">{{ item.time }}</text>
					</view>
				</view>
				<view class="right">
					<text class="money" :class="{ 'income': item.amount > 0 }">
						{{ item.amount > 0 ? '+' : '' }}{{ item.amount }}
					</text>
				</view>
			</view>
			
			<uv-empty v-if="records.length === 0" mode="list" text="暂无数据"></uv-empty>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue';

// 新增: 返回上一页方法
const goBack = () => {
	uni.navigateBack();
};

// 模拟数据
const records = ref([
	{ title: '设备使用消费', time: '2023-10-27 14:30', amount: -5.00, type: 'expense', icon: 'shopping-cart' },
	{ title: '账户充值', time: '2023-10-26 09:15', amount: 100.00, type: 'income', icon: 'wallet' },
	{ title: '设备使用消费', time: '2023-10-25 18:20', amount: -3.50, type: 'expense', icon: 'shopping-cart' },
]);

const handleRecharge = () => {
	uni.showToast({ title: '充值功能开发中', icon: 'none' });
};

const handleWithdraw = () => {
	uni.showToast({ title: '提现功能开发中', icon: 'none' });
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.asset-card {
	position: relative;
	margin: 110rpx 30rpx 30rpx 30rpx;
	height: 320rpx;
	border-radius: 24rpx;
	overflow: hidden;
	box-shadow: 0 10rpx 30rpx rgba(60, 156, 255, 0.2);
	
	.card-bg {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: linear-gradient(135deg, #3c9cff 0%, #2b85e4 100%);
		z-index: 0;
	}
	
	.card-content {
		position: relative;
		z-index: 1;
		padding: 40rpx;
		color: #fff;
		display: flex;
		flex-direction: column;
		height: 100%;
		box-sizing: border-box;
		
		.label {
			font-size: 26rpx;
			opacity: 0.9;
			margin-bottom: 10rpx;
		}
		
		.amount-row {
			display: flex;
			align-items: baseline;
			margin-bottom: 40rpx;
			
			.symbol {
				font-size: 32rpx;
				margin-right: 8rpx;
			}
			
			.amount {
				font-size: 64rpx;
				font-weight: bold;
				font-family: 'DIN Alternate', sans-serif;
			}
		}
		
		.actions {
			display: flex;
			gap: 40rpx;
			margin-top: auto;
			
			.action-btn {
				display: flex;
				flex-direction: column;
				align-items: center;
				gap: 10rpx;
				
				text {
					font-size: 24rpx;
				}
			}
		}
	}
}

.section-title {
	padding: 30rpx 30rpx 20rpx;
	font-size: 30rpx;
	font-weight: bold;
	color: #333;
}

.record-list {
	padding: 0 30rpx 30rpx;
	
	.record-item {
		background: #fff;
		border-radius: 16rpx;
		padding: 24rpx;
		margin-bottom: 20rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 2rpx 10rpx rgba(0,0,0,0.02);
		
		.left {
			display: flex;
			align-items: center;
			gap: 20rpx;
			
			.icon-box {
				width: 70rpx;
				height: 70rpx;
				border-radius: 50%;
				display: flex;
				align-items: center;
				justify-content: center;
				
				&.expense {
					background: rgba(250, 53, 52, 0.1);
				}
				
				&.income {
					background: rgba(25, 190, 107, 0.1);
				}
			}
			
			.info {
				display: flex;
				flex-direction: column;
				
				.title {
					font-size: 28rpx;
					color: #333;
					margin-bottom: 6rpx;
				}
				
				.time {
					font-size: 22rpx;
					color: #999;
				}
			}
		}
		
		.right {
			.money {
				font-size: 32rpx;
				font-weight: bold;
				color: #333;
				font-family: 'DIN Alternate', sans-serif;
				
				&.income {
					color: #19be6b;
				}
			}
		}
	}
}
</style>