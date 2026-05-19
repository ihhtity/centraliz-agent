<template>
	<view class="popup-container" @click="closePopup">
		<view class="popup-content" @click.stop>
			<view class="header">
				<view class="header-left">
					<uv-icon name="home-fill" size="20" color="#2979ff" />
					<text class="title">{{ locker.no }}</text>
				</view>
				<view class="close-btn" @click="closePopup">
					<uv-icon name="close" size="20" color="#999"></uv-icon>
				</view>
			</view>
			
			<scroll-view scroll-y class="content-scroll" enable-flex>
				<view class="content">
					<!-- 锁管理区域 -->
					<view class="info-section">
						<view class="section-header">
							<view class="header-icon-box">
								<uv-icon name="lock" size="14" color="#fff" />
							</view>
							<text class="section-title">锁管理</text>
						</view>
						<uv-grid :col="4" :border="false" :gap="10">
							<uv-grid-item>
								<uv-button type="primary" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">二维码</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="primary" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">导入码</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="success" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">获密码</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="error" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">远程开</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="primary" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">蓝牙开</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="warning" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">开锁记</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="error" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">中断单</uv-button>
							</uv-grid-item>
						</uv-grid>
					</view>

					<!-- 电控管理区域 -->
					<view class="info-section">
						<view class="section-header">
							<view class="header-icon-box bg-blue">
								<uv-icon name="setting" size="14" color="#fff" />
							</view>
							<text class="section-title">电控管理</text>
						</view>
						<view class="device-info-list">
							<view class="info-item">
								<text class="label">设备ID</text>
								<text class="value mono-font">{{ locker.deviceId || '-' }}</text>
							</view>
							<uv-line color="#f0f0f0" />
							<view class="info-item">
								<text class="label">信号强度</text>
								<text class="value">{{ locker.signal || '-' }}</text>
							</view>
							<uv-line color="#f0f0f0" />
							<view class="info-item">
								<text class="label">当前状态</text>
								<view class="status-tag" :class="getStatusClass(locker.status)">
									{{ locker.status || '-' }}
								</view>
							</view>
							<uv-line color="#f0f0f0" />
							<view class="info-item">
								<text class="label">实时电流</text>
								<text class="value">{{ locker.current || '-' }}</text>
							</view>
							<uv-line color="#f0f0f0" />
							<view class="info-item">
								<text class="label">保护状态</text>
								<text class="value">{{ locker.protection || '-' }}</text>
							</view>
						</view>
					</view>

					<!-- 其他功能区域 -->
					<view class="info-section">
						<view class="section-header">
							<view class="header-icon-box bg-gray">
								<uv-icon name="more" size="14" color="#fff" />
							</view>
							<text class="section-title">其他功能</text>
						</view>
						<uv-grid :col="4" :border="false" :gap="10">
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">网络高</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">开关控</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">状态查</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">改设置</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">预约设</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">控电记</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">会员充</uv-button>
							</uv-grid-item>
							<uv-grid-item>
								<uv-button type="info" shape="circle" size="mini" plain custom-style="font-size: 22rpx; padding: 0 10rpx;">解锁控</uv-button>
							</uv-grid-item>
						</uv-grid>
					</view>
				</view>
			</scroll-view>
		</view>
	</view>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
	locker: Object
});

const emit = defineEmits(['close']);

const closePopup = () => {
	emit('close');
};

// 根据状态返回对应的颜色类名
const getStatusClass = (status) => {
	// 修复: 确保 status 是字符串，防止 TypeError
	if (status === null || status === undefined) return 'status-default';
	
	const statusStr = String(status);
	
	// 如果原本是数字状态 (0:空闲, 1:租用, 2:维修)，可以根据需求映射颜色
	// 这里保持原有逻辑的兼容性，如果包含特定关键词则高亮
	if (statusStr.includes('正常') || statusStr.includes('在线') || statusStr === '0') return 'status-success'; // 假设0是空闲/正常
	if (statusStr.includes('故障') || statusStr.includes('离线') || statusStr === '2') return 'status-error';   // 假设2是维修/故障
	if (statusStr === '1') return 'status-warning'; // 假设1是租用/警告色
	
	return 'status-default';
};
</script>

<style lang="scss" scoped>
.popup-container {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: rgba(0, 0, 0, 0.6);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 999;
	backdrop-filter: blur(4rpx);

	.popup-content {
		background-color: #f5f6f8;
		border-radius: 24rpx;
		padding: 0;
		box-shadow: 0 12rpx 32rpx rgba(0, 0, 0, 0.1);
		width: 85%;
		max-width: 700rpx;
		max-height: 80vh;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		animation: popup-slide-up 0.3s ease-out;

		.header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 32rpx 32rpx 24rpx;
			background-color: #fff;
			border-bottom: 1rpx solid #f0f0f0;

			.header-left {
				display: flex;
				align-items: center;
				
				.title {
					margin-left: 12rpx;
					font-size: 34rpx;
					font-weight: 700;
					color: #303133;
					letter-spacing: 1rpx;
				}
			}

			.close-btn {
				padding: 8rpx;
				border-radius: 50%;
				background-color: #f5f7fa;
				display: flex;
				align-items: center;
				justify-content: center;
				transition: background-color 0.2s;
				
				&:active {
					background-color: #e4e7ed;
				}
			}
		}

		.content-scroll {
			flex: 1;
			overflow-y: auto;
			padding-bottom: 20rpx; /* 底部留白 */
		}

		.content {
			padding: 24rpx;

			.info-section {
				background-color: #fff;
				border-radius: 20rpx;
				padding: 24rpx;
				margin-bottom: 24rpx;
				box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
				transition: transform 0.2s;

				&:last-child {
					margin-bottom: 0;
				}

				.section-header {
					display: flex;
					align-items: center;
					margin-bottom: 24rpx;

					.header-icon-box {
						width: 40rpx;
						height: 40rpx;
						border-radius: 12rpx;
						background-color: #2979ff;
						display: flex;
						align-items: center;
						justify-content: center;
						margin-right: 12rpx;
						
						&.bg-blue {
							background-color: #2979ff;
						}
						
						&.bg-gray {
							background-color: #909399;
						}
					}

					.section-title {
						font-size: 30rpx;
						font-weight: 600;
						color: #303133;
					}
				}

				.device-info-list {
					.info-item {
						display: flex;
						justify-content: space-between;
						align-items: center;
						padding: 20rpx 0;
						
						.label {
							font-size: 28rpx;
							color: #606266;
							font-weight: 400;
						}
						
						.value {
							font-size: 28rpx;
							color: #303133;
							font-weight: 500;
							
							&.mono-font {
								font-family: monospace;
								color: #909399;
								font-size: 26rpx;
							}
						}

						.status-tag {
							padding: 4rpx 16rpx;
							border-radius: 8rpx;
							font-size: 24rpx;
							font-weight: 500;
							
							&.status-success {
								background-color: #f0f9eb;
								color: #67c23a;
								border: 1rpx solid #e1f3d8;
							}
							&.status-error {
								background-color: #fef0f0;
								color: #f56c6c;
								border: 1rpx solid #fde2e2;
							}
							&.status-warning {
								background-color: #fdf6ec;
								color: #e6a23c;
								border: 1rpx solid #faecd8;
							}
							&.status-default {
								background-color: #f4f4f5;
								color: #909399;
								border: 1rpx solid #e9e9eb;
							}
						}
					}
				}
			}
		}
	}
}

@keyframes popup-slide-up {
	from {
		opacity: 0;
		transform: translateY(20rpx);
	}
	to {
		opacity: 1;
		transform: translateY(0);
	}
}
</style>