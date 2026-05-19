<template>
	<view class="container">
		<uv-navbar :title="t('user.device.selectDevice')" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
		
		<view class="device-list">
			<view v-for="item in devices" :key="item.id" class="item">
				<text>{{ item.name }}</text>
				<text :class="{ available: item.status === 1 }">
					{{ item.status === 1 ? t('user.device.available') : t('user.device.inUse') }}
				</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const devices = ref([])

// 新增: 返回上一页方法
const goBack = () => {
	uni.navigateBack();
};

</script>

<style scoped lang="scss">
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 20rpx;
}

.device-list {
	padding: 0 24rpx;
}

.item {
	background-color: #fff;
	margin-bottom: 20rpx;
	border-radius: 16rpx;
	box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.03);
	transition: all 0.2s ease;
	
	&:active {
		transform: scale(0.98);
		background-color: #fafafa;
	}
	
	text {
		padding: 24rpx 30rpx;
		font-size: 30rpx;
		font-weight: 600;
		color: #333;
	}
	
	text:last-child {
		font-size: 24rpx;
		color: #999;
		margin-top: 8rpx;
	}
	
	.available {
		color: #19be6b;
	}
}
</style>
