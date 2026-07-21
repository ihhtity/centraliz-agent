<template>
	<view class="container">
		<uv-navbar :title="t('user.profile.agreement')" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
		
		<scroll-view scroll-y class="agreement-scroll">
			<view class="agreement-card">
				<view class="section-header">
					<text class="section-icon">📦</text>
					<text class="section-title">{{ t('user.locker.agreementTitle') }}</text>
				</view>
				
				<view class="agreement-section">
					<text class="subsection-title">{{ t('user.locker.privacyPolicy') }}</text>
					<text class="content-text">{{ t('user.locker.storagePrivacy') }}</text>
				</view>
				
				<view class="agreement-divider"></view>
				
				<view class="agreement-section">
					<text class="subsection-title">{{ t('user.locker.termsOfService') }}</text>
					<text class="content-text">{{ t('user.locker.storageTerms') }}</text>
				</view>
			</view>
			
			<view class="agreement-card">
				<view class="section-header">
					<text class="section-icon">🛒</text>
					<text class="section-title">{{ t('user.locker.retailAgreementTitle') }}</text>
				</view>
				
				<view class="agreement-section">
					<text class="subsection-title">{{ t('user.locker.privacyPolicy') }}</text>
					<text class="content-text">{{ t('user.locker.retailPrivacy') }}</text>
				</view>
				
				<view class="agreement-divider"></view>
				
				<view class="agreement-section">
					<text class="subsection-title">{{ t('user.locker.termsOfService') }}</text>
					<text class="content-text">{{ t('user.locker.retailTerms') }}</text>
				</view>
			</view>
		</scroll-view>
		
		<view class="footer">
			<view class="action-btn" :class="{ active: !hasAgreed, disabled: hasAgreed }" @click="handleAgree">
				<text>{{ hasAgreed ? t('user.profile.agreed') : t('common.agree') }}</text>
			</view>
			<view class="action-btn" :class="{ active: hasAgreed, disabled: !hasAgreed }" @click="handleDisagree">
				<text>{{ t('user.profile.disagree') }}</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { onLoad } from '@dcloudio/uni-app'

onLoad(() => {
	user.value = uni.getStorageSync('user') || {}
})

const { t } = useI18n()
const user = ref({})
const loading = ref(false)
const hasAgreed = computed(() => user.value.privacy === '1')
// 同意协议
const handleAgree = async () => {
	if (hasAgreed.value) return
	if (!user.value?.id) {
		uni.showToast({ title: t('common.loginFirst'), icon: 'none' })
		return
	}
	
	loading.value = true
	try {
		const res = await uni.$uv.http.put('/user/profile', {
			id: user.value.id,
			privacy: '1'
		}, { custom: { auth: true } })
		
		if (res.code === 200) {
			user.value.privacy = '1'
			uni.setStorageSync('user', user.value)
			uni.showToast({ title: t('user.profile.agreementSuccess'), icon: 'success', duration: 2000 })
		} else {
			uni.showToast({ title: res.msg || t('common.operationFailed'), icon: 'none', duration: 2000 })
		}
	} catch (error) {
		console.error('同意协议失败:', error)
		uni.showToast({ title: t('common.networkError'), icon: 'none', duration: 2000 })
	} finally {
		loading.value = false
	}
}
// 拒绝协议
const handleDisagree = async () => {
	if (!hasAgreed.value) return
	if (!user.value?.id) {
		uni.showToast({ title: t('common.loginFirst'), icon: 'none' })
		return
	}
	
	uni.showModal({
		title: t('common.confirm'),
		content: t('user.profile.disagreeConfirm'),
		cancelText: t('common.cancel'),
		confirmText: t('common.confirm'),
		success: async (res) => {
			if (res.confirm) {
				loading.value = true
				try {
					const res = await uni.$uv.http.put('/user/profile', {
						id: user.value.id,
						privacy: '0'
					}, { custom: { auth: true } })
					
					if (res.code === 200) {
						user.value.privacy = '0'
						uni.setStorageSync('user', user.value)
						uni.showToast({ title: t('user.profile.disagreeSuccess'), icon: 'none', duration: 2000 })
					} else {
						uni.showToast({ title: res.msg || t('common.operationFailed'), icon: 'none', duration: 2000 })
					}
				} catch (error) {
					console.error('拒绝协议失败:', error)
					uni.showToast({ title: t('common.networkError'), icon: 'none', duration: 2000 })
				} finally {
					loading.value = false
				}
			}
		}
	})
}
// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/user/profile/index'
	})
}
</script>

<style lang="scss" scoped>
.container {
	background-color: #f5f7fa;
	display: flex;
	flex-direction: column;
}

.agreement-scroll {
	flex: 1;
	padding-top: 24rpx;
	padding-bottom: 160rpx;
}

.agreement-card {
	background-color: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.section-header {
	display: flex;
	align-items: center;
	margin-bottom: 28rpx;
	
	.section-icon {
		font-size: 40rpx;
		margin-right: 16rpx;
	}
	
	.section-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #1a1a1a;
	}
}

.agreement-section {
	padding: 16rpx 0;
}

.subsection-title {
	display: block;
	font-size: 28rpx;
	font-weight: 600;
	color: #333;
	margin-bottom: 16rpx;
}

.content-text {
	font-size: 26rpx;
	color: #666;
	line-height: 1.8;
	white-space: pre-wrap;
}

.agreement-divider {
	height: 16rpx;
	background-color: #f8f8f8;
	margin: 20rpx 0;
	border-radius: 8rpx;
}

.footer {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fff;
	padding: 24rpx;
	padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
	box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.04);
	display: flex;
	gap: 24rpx;
}

.action-btn {
	flex: 1;
	height: 88rpx;
	line-height: 88rpx;
	border-radius: 44rpx;
	text-align: center;
	font-size: 30rpx;
	font-weight: 500;
	background-color: #f5f7fa;
	color: #999;
	border: 2rpx solid transparent;
	transition: all 0.3s ease;
	
	&.active {
		background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
		color: #fff;
		border-color: transparent;
		
		&:active {
			opacity: 0.85;
		}
	}
	
	&.disabled {
		background-color: #f0f0f0;
		color: #ccc;
		border-color: transparent;
	}
}
</style>