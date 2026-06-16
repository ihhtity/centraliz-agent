<!-- 电池指示器组件 -->
<template>
	<view class="battery-indicator" :class="batteryClass">
		<!-- 自定义电池图形 -->
		<view class="battery-body">
			<view class="battery-head"></view>
			<view class="battery-level" :style="{ width: percentage + '%' }"></view>
		</view>
		<!-- 电量百分比文本 -->
		<text v-if="showText" class="battery-text">{{ percentage }}%</text>
	</view>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
	percentage: {
		type: Number,
		default: 0,
		validator: (value) => value >= 0 && value <= 100
	},
	size: {
		type: [String, Number],
		default: 32
	},
	showText: {
		type: Boolean,
		default: true
	}
})

// 根据电量百分比确定颜色类
const batteryClass = computed(() => {
	if (props.percentage >= 80) return 'battery-high'
	if (props.percentage >= 30) return 'battery-medium'
	return 'battery-low'
})


</script>

<style lang="scss" scoped>
.battery-indicator {
	display: flex;
	align-items: center;
	gap: 8rpx;
	
	// 电池主体容器
	.battery-body {
		position: relative;
		width: 40rpx; // 默认宽度，可根据size prop进一步动态化，这里先固定比例美观为主
		height: 20rpx;
		border: 2rpx solid #999;
		border-radius: 4rpx;
		padding: 2rpx;
		box-sizing: border-box;
		display: flex;
		align-items: center;
		
		// 电池正极头
		&::before {
			content: '';
			position: absolute;
			right: -6rpx;
			top: 50%;
			transform: translateY(-50%);
			width: 4rpx;
			height: 10rpx;
			background-color: #999;
			border-top-right-radius: 2rpx;
			border-bottom-right-radius: 2rpx;
		}
		
		// 电量进度条
		.battery-level {
			height: 100%;
			border-radius: 2rpx;
			transition: width 0.3s ease, background-color 0.3s ease;
		}
	}
	
	.battery-text {
		font-size: 24rpx;
		font-weight: 500;
		min-width: 60rpx; // 防止数字变化导致抖动
		text-align: right;
	}
	
	// 高电量样式
	&.battery-high {
		.battery-body {
			border-color: #52c41a;
			&::before {
				background-color: #52c41a;
			}
		}
		.battery-level {
			background-color: #52c41a;
		}
		.battery-text {
			color: #52c41a;
		}
	}
	
	// 中电量样式
	&.battery-medium {
		.battery-body {
			border-color: #faad14;
			&::before {
				background-color: #faad14;
			}
		}
		.battery-level {
			background-color: #faad14;
		}
		.battery-text {
			color: #faad14;
		}
	}
	
	// 低电量样式
	&.battery-low {
		.battery-body {
			border-color: #ff4d4f;
			&::before {
				background-color: #ff4d4f;
			}
		}
		.battery-level {
			background-color: #ff4d4f;
		}
		.battery-text {
			color: #ff4d4f;
		}
	}
}
</style>
