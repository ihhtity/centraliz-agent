<!-- 关于我们页面 -->
<template>
	<view class="container">
		<uv-navbar :title="'关于我们'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

		<!-- 公司Logo区域 -->
		<view class="company-header">
			<view class="logo-box">
				<uv-icon name="home-fill" size="80" color="#3c9cff" />
			</view>
			<text class="company-name">宝士力得科技有限公司</text>
			<text class="company-slogan">智能科技 · 引领未来</text>
		</view>

		<!-- 公司介绍 -->
		<view class="section">
			<view class="section-title">公司简介</view>
			<view class="section-content">
				<text class="intro-text">
					宝士力得科技有限公司是一家专注于智能科技领域的高新技术企业。公司致力于为客户提供创新的技术解决方案，涵盖智能硬件、软件开发、系统集成等多个领域。
				</text>
				<text class="intro-text">
					我们秉承"科技创新、服务至上"的理念，不断追求卓越，为客户创造价值。凭借专业的技术团队和丰富的行业经验，我们已成为行业内的领先企业。
				</text>
			</view>
		</view>

		<!-- 联系信息 -->
		<view class="section">
			<view class="section-title">联系我们</view>
			<view class="contact-list">
				<view class="contact-item" @click="copyToClipboard('0755-88888888')">
					<view class="contact-icon">
						<uv-icon name="phone-fill" size="32" color="#3c9cff" />
					</view>
					<view class="contact-info">
						<text class="contact-label">客服电话</text>
						<text class="contact-value">0755-88888888</text>
					</view>
					<uv-icon name="arrow-right" size="24" color="#ccc" />
				</view>

				<view class="contact-item" @click="copyToClipboard('support@bsldtech.com')">
					<view class="contact-icon">
						<uv-icon name="email-fill" size="32" color="#10b981" />
					</view>
					<view class="contact-info">
						<text class="contact-label">电子邮箱</text>
						<text class="contact-value">support@bsldtech.com</text>
					</view>
					<uv-icon name="arrow-right" size="24" color="#ccc" />
				</view>

				<view class="contact-item" @click="copyToClipboard('www.bsldtech.cn')">
					<view class="contact-icon">
						<uv-icon name="attach" size="32" color="#a855f7" />
					</view>
					<view class="contact-info">
						<text class="contact-label">官方网站</text>
						<text class="contact-value">www.bsldtech.cn</text>
					</view>
					<uv-icon name="arrow-right" size="24" color="#ccc" />
				</view>

				<view class="contact-item">
					<view class="contact-icon">
						<uv-icon name="map-fill" size="32" color="#ff9900" />
					</view>
					<view class="contact-info">
						<text class="contact-label">公司地址</text>
						<text class="contact-value">深圳市南山区科技园南区高新南一道</text>
					</view>
				</view>

				<view class="contact-item">
					<view class="contact-icon">
						<uv-icon name="clock-fill" size="32" color="#fa3534" />
					</view>
					<view class="contact-info">
						<text class="contact-label">工作时间</text>
						<text class="contact-value">周一至周五 9:00-18:00</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 版本信息 -->
		<view class="version-section">
			<text class="version-text">版本号：v1.0.0</text>
			<text class="copyright">© 2026 宝士力得科技有限公司 版权所有</text>
		</view>
	</view>
</template>

<script setup>
// 关于我们页面脚本

// 返回上一页
const goBack = () => {
	uni.navigateBack();
};

// 复制到剪贴板
const copyToClipboard = (text) => {
	uni.setClipboardData({
		data: text,
		success: () => {
			uni.showToast({
				title: '已复制',
				icon: 'success',
				duration: 1500
			});
		},
		fail: () => {
			uni.showToast({
				title: '复制失败',
				icon: 'none',
				duration: 1500
			});
		}
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background: linear-gradient(180deg, #f0f4f8 0%, #e8edf3 100%);
}

.company-header {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 80rpx 40rpx 60rpx;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	position: relative;
	overflow: hidden;

	&::before {
		content: '';
		position: absolute;
		top: -100rpx;
		right: -100rpx;
		width: 300rpx;
		height: 300rpx;
		background: rgba(255, 255, 255, 0.1);
		border-radius: 50%;
	}

	&::after {
		content: '';
		position: absolute;
		bottom: -80rpx;
		left: -80rpx;
		width: 200rpx;
		height: 200rpx;
		background: rgba(255, 255, 255, 0.08);
		border-radius: 50%;
	}

	.logo-box {
		width: 160rpx;
		height: 160rpx;
		border-radius: 50%;
		background: rgba(255, 255, 255, 0.2);
		backdrop-filter: blur(10px);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 30rpx;
		box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.15);
		position: relative;
		z-index: 1;
		transition: transform 0.3s ease;

		&:active {
			transform: scale(0.95);
		}
	}

	.company-name {
		font-size: 40rpx;
		font-weight: 700;
		color: #fff;
		margin-bottom: 16rpx;
		position: relative;
		z-index: 1;
		text-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
	}

	.company-slogan {
		font-size: 26rpx;
		color: rgba(255, 255, 255, 0.85);
		position: relative;
		z-index: 1;
	}
}

.section {
	background: #fff;
	margin: 24rpx;
	border-radius: 24rpx;
	padding: 36rpx;
	box-shadow: 0 4rpx 24rpx rgba(0, 0, 0, 0.04);
	transition: transform 0.2s ease, box-shadow 0.2s ease;

	&:active {
		transform: scale(0.99);
	}

	.section-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #1a1a2e;
		margin-bottom: 28rpx;
		display: flex;
		align-items: center;

		&::before {
			content: '';
			width: 6rpx;
			height: 32rpx;
			background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
			border-radius: 3rpx;
			margin-right: 16rpx;
		}
	}

	.section-content {
		.intro-text {
			display: block;
			font-size: 28rpx;
			color: #5a5a6e;
			line-height: 1.8;
			text-align: justify;
			margin-bottom: 20rpx;

			&:last-child {
				margin-bottom: 0;
			}
		}
	}
}

.contact-list {
	padding: 0 8rpx;

	.contact-item {
		display: flex;
		align-items: center;
		padding: 28rpx 0;
		border-bottom: 1rpx solid #f0f2f5;
		transition: all 0.2s ease;

		&:active {
			background: #f8f9fa;
			padding-left: 16rpx;
			border-radius: 16rpx;
		}

		&:last-child {
			border-bottom: none;
		}

		.contact-icon {
			width: 72rpx;
			height: 72rpx;
			border-radius: 20rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-right: 24rpx;
			background: linear-gradient(135deg, #f0f4f8 0%, #e8edf3 100%);
			transition: transform 0.2s ease;

			&:active {
				transform: scale(0.9);
			}
		}

		.contact-info {
			flex: 1;
			display: flex;
			flex-direction: column;

			.contact-label {
				font-size: 24rpx;
				color: #8a8a9a;
				margin-bottom: 6rpx;
			}

			.contact-value {
				font-size: 28rpx;
				color: #2d2d3d;
				font-weight: 500;
			}
		}
	}
}

.version-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 60rpx 40rpx 100rpx;

	.version-text {
		font-size: 26rpx;
		color: #9a9aaa;
		margin-bottom: 16rpx;
		display: flex;
		align-items: center;
		gap: 8rpx;
	}

	.copyright {
		font-size: 24rpx;
		color: #b8b8c8;
	}
}

// 渐变色图标背景
.icon-gradient-phone {
	background: linear-gradient(135deg, #3c9cff 0%, #2b85e4 100%);
}

.icon-gradient-email {
	background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.icon-gradient-web {
	background: linear-gradient(135deg, #a855f7 0%, #9333ea 100%);
}

.icon-gradient-map {
	background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.icon-gradient-clock {
	background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
}
</style>
