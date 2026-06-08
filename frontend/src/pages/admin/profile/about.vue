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
	background-color: #f5f7fa;
}

.company-header {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 60rpx 40rpx;
	background: linear-gradient(180deg, #fff 0%, #f5f7fa 100%);

	.logo-box {
		width: 160rpx;
		height: 160rpx;
		border-radius: 40rpx;
		background: linear-gradient(135deg, #3c9cff, #2b85e4);
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 30rpx;
		box-shadow: 0 8rpx 32rpx rgba(60, 156, 255, 0.3);
	}

	.company-name {
		font-size: 40rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 16rpx;
	}

	.company-slogan {
		font-size: 28rpx;
		color: #999;
	}
}

.section {
	background: #fff;
	margin: 24rpx;
	border-radius: 20rpx;
	padding: 32rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);

	.section-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
		margin-bottom: 24rpx;
		padding-left: 20rpx;
		border-left: 8rpx solid #3c9cff;
		line-height: 1;
	}

	.section-content {
		.intro-text {
			display: block;
			font-size: 28rpx;
			color: #666;
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
	.contact-item {
		display: flex;
		align-items: center;
		padding: 24rpx 0;
		border-bottom: 1rpx solid #f0f0f0;
		transition: all 0.2s;

		&:active {
			opacity: 0.7;
		}

		&:last-child {
			border-bottom: none;
		}

		.contact-icon {
			width: 64rpx;
			height: 64rpx;
			border-radius: 16rpx;
			background: #f5f7fa;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-right: 24rpx;
		}

		.contact-info {
			flex: 1;
			display: flex;
			flex-direction: column;

			.contact-label {
				font-size: 24rpx;
				color: #999;
				margin-bottom: 8rpx;
			}

			.contact-value {
				font-size: 28rpx;
				color: #333;
				font-weight: 500;
			}
		}
	}
}

.version-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 60rpx 40rpx;

	.version-text {
		font-size: 26rpx;
		color: #999;
		margin-bottom: 16rpx;
	}

	.copyright {
		font-size: 24rpx;
		color: #bbb;
	}
}
</style>
