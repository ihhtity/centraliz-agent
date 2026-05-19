<template>
	<view class="container">
		<!-- 修改: 添加返回按钮 -->
		<uv-navbar :title="deviceName" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack"></uv-navbar>
		
		<view class="status-card">
			<uv-icon name="checkmark-circle" color="#19be6b" size="60"></uv-icon>
			<text class="status-text">{{ t('user.device.running') }}</text>
			<text class="time-text">{{ t('user.device.duration') }}: {{ useTime }}</text>
		</view>
		
		<view class="control-panel">
			<uv-button type="error" size="large" @click="stopUse">{{ t('user.device.stopUse') }}</uv-button>
		</view>
		
		<view class="info-panel">
			<uv-cell-group :title="t('user.device.deviceInfo')">
				<uv-cell :title="t('admin.device.deviceCode')" :value="deviceId"></uv-cell>
				<uv-cell :title="t('user.device.currentLocation')" value="A区-01号位"></uv-cell>
			</uv-cell-group>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// 新增: 返回上一页方法
const goBack = () => {
	uni.navigateBack();
};

const deviceId = ref('');
const deviceName = ref('');
const useTime = ref('00:00:00');
let timer = null;
let seconds = 0;

onLoad((options) => {
	deviceId.value = options.id;
	deviceName.value = options.name;
	startTimer();
});

const startTimer = () => {
	timer = setInterval(() => {
		seconds++;
		const h = Math.floor(seconds / 3600).toString().padStart(2, '0');
		const m = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
		const s = (seconds % 60).toString().padStart(2, '0');
		useTime.value = `${h}:${m}:${s}`;
	}, 1000);
};

const stopUse = () => {
	uni.showModal({
		title: t('common.confirm'),
		content: t('common.confirm'), // 这里可能需要更具体的 content key，暂时复用 confirm 或添加新 key
		success: (res) => {
			if (res.confirm) {
				clearInterval(timer);
				// 调用停止接口
				uni.showToast({ title: t('common.operationSuccess'), icon: 'success' });
				setTimeout(() => {
					uni.navigateBack();
				}, 1500);
			}
		}
	});
};

onUnmounted(() => {
	if (timer) clearInterval(timer);
});
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding: 24rpx;
}
.status-card {
	background: #fff;
	border-radius: 24rpx;
	padding: 80rpx 0;
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-bottom: 40rpx;
	// 新增: 增加轻微阴影和边框
	box-shadow: 0 8rpx 24rpx rgba(25, 190, 107, 0.1);
	border: 1rpx solid rgba(25, 190, 107, 0.1);
	
	.status-text {
		font-size: 32rpx;
		font-weight: bold;
		margin-top: 24rpx;
		color: #333;
	}
	.time-text {
		font-size: 56rpx;
		font-weight: bold;
		background: linear-gradient(135deg, #19be6b 0%, #13ae5c 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		margin-top: 16rpx;
		font-family: 'DIN Alternate', sans-serif; // 假设有数字字体，否则回退
	}
	
	// 新增: 图标呼吸动画
	.uv-icon {
		animation: pulse 2s infinite;
	}
}

.control-panel {
	margin-bottom: 40rpx;
	
	// 优化按钮样式
	::v-deep .uv-button {
		border-radius: 44rpx !important;
		box-shadow: 0 8rpx 16rpx rgba(250, 53, 52, 0.2);
	}
}
.info-panel {
	background: #fff;
	border-radius: 20rpx;
	overflow: hidden;
	box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.03);
}

@keyframes pulse {
	0% { transform: scale(1); opacity: 1; }
	50% { transform: scale(1.1); opacity: 0.8; }
	100% { transform: scale(1); opacity: 1; }
}
</style>