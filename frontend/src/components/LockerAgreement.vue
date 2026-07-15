<!-- 用户使用协议与隐私政策承诺书 -->
<template>
	<uv-popup ref="popup" mode="center" :round="20" custom-style="max-height: 80vh;min-height: 50rpx;">
		<view class="agreement-container">
			<view class="agreement-header">
				<text class="title">{{ currentTitle }}</text>
				<view class="close-btn" @click="close">
					<text class="close-icon">×</text>
				</view>
			</view>
			
			<scroll-view scroll-y class="agreement-content">
				<view class="agreement-section">
					<text class="section-title">{{ t('user.locker.privacyPolicy') }}</text>
					<text class="content-text">{{ privacyContent }}</text>
				</view>
				<view class="agreement-divider"></view>
				<view class="agreement-section">
					<text class="section-title">{{ t('user.locker.termsOfService') }}</text>
					<text class="content-text">{{ termsContent }}</text>
				</view>
			</scroll-view>
			
			<view class="agreement-footer">
				<button class="btn-disagree" @click="handleDisagree">{{ disagreeText || t('user.locker.disagree') }}</button>
				<button class="btn-agree" @click="handleAgree">{{ agreeText || t('user.locker.agreeAndContinue') }}</button>
			</view>
		</view>
	</uv-popup>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
	title: { // 自定义标题
		type: String,
		default: ''
	},
	content: { // 自定义内容
		type: String,
		default: ''
	},
	agreeText: { // 自定义同意按钮文本
		type: String,
		default: ''
	},
	disagreeText: { // 自定义不同意按钮文本
		type: String,
		default: ''
	},
	groupType: { // 分组类型（存柜/零售）
		type: String,
		default: '存柜'
	}
})

const emit = defineEmits(['agree', 'disagree'])

const popup = ref(null)

const isStorage = computed(() => props.groupType === '存柜')

const currentTitle = computed(() => {
	if (props.title) return props.title
	return isStorage.value ? t('user.locker.agreementTitle') : t('user.locker.retailAgreementTitle')
})

const privacyContent = computed(() => {
	if (props.content) return props.content
	return isStorage.value 
		? t('user.locker.storagePrivacy') 
		: t('user.locker.retailPrivacy')
})

const termsContent = computed(() => {
	return isStorage.value 
		? t('user.locker.storageTerms') 
		: t('user.locker.retailTerms')
})

const open = () => {
	popup.value.open()
}

const close = () => {
	popup.value.close()
}

const handleAgree = () => {
	emit('agree')
	close()
}

const handleDisagree = () => {
	emit('disagree')
	close()
}

defineExpose({
	open,
	close
})
</script>

<style lang="scss" scoped>
.agreement-container {
	background-color: #fff;
	border-radius: 24rpx;
	padding: 40rpx;
	display: flex;
	flex-direction: column;
	box-shadow: 0 10rpx 40rpx rgba(0, 0, 0, 0.1);
}

.agreement-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 30rpx;
	padding-bottom: 20rpx;
	border-bottom: 1rpx solid #f0f0f0;
	
	.title {
		font-size: 34rpx;
		font-weight: 600;
		color: #1a1a1a;
	}
	
	.close-btn {
		width: 60rpx;
		height: 60rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 50%;
		background-color: #f5f5f5;
		
		.close-icon {
			font-size: 40rpx;
			color: #999;
			line-height: 1;
		}
	}
}

.agreement-content {
	flex: 1;
	margin-bottom: 30rpx;
	padding-right: 10rpx;
	
	.agreement-section {
		padding: 20rpx 0;
	}
	
	.section-title {
		display: block;
		font-size: 30rpx;
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
		height: 10rpx;
		background-color: #f8f8f8;
		margin: 10rpx 0;
	}
}

.agreement-footer {
	display: flex;
	justify-content: space-between;
	gap: 24rpx;
	
	button {
		flex: 1;
		height: 88rpx;
		line-height: 88rpx;
		border-radius: 44rpx;
		font-size: 30rpx;
		font-weight: 500;
		border: none;
		
		&::after {
			border: none;
		}
	}
	
	.btn-disagree {
		background-color: #f5f5f5;
		color: #666;
		
		&:active {
			background-color: #eee;
		}
	}
	
	.btn-agree {
		background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
		color: #fff;
		box-shadow: 0 6rpx 16rpx rgba(79, 172, 254, 0.4);
		
		&:active {
			opacity: 0.9;
			transform: scale(0.98);
		}
	}
}
</style>