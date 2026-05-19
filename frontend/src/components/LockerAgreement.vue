<template>
	<uv-popup ref="popup" mode="center" :round="20">
		<view class="agreement-container">
			<view class="agreement-header">
				<text class="title">{{ title || t('user.locker.agreementTitle') }}</text>
			</view>
			
			<scroll-view scroll-y class="agreement-content">
				<text class="content-text">{{ content || t('user.locker.agreementContent') }}</text>
			</scroll-view>
			
			<view class="agreement-footer">
				<button class="btn-disagree" @click="handleDisagree">{{ disagreeText || t('user.locker.disagree') }}</button>
				<button class="btn-agree" @click="handleAgree">{{ agreeText || t('user.locker.agreeAndContinue') }}</button>
			</view>
		</view>
	</uv-popup>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
	title: {
		type: String,
		default: ''
	},
	content: {
		type: String,
		default: ''
	},
	agreeText: {
		type: String,
		default: ''
	},
	disagreeText: {
		type: String,
		default: ''
	}
})

const emit = defineEmits(['agree', 'disagree'])

const popup = ref(null)

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
	width: 600rpx;
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	display: flex;
	flex-direction: column;
	max-height: 80vh;
}

.agreement-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
	
	.title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}
	
	.close-btn {
		padding: 10rpx;
	}
}

.agreement-content {
	flex: 1;
	min-height: 300rpx;
	max-height: 500rpx;
	margin-bottom: 30rpx;
	
	.content-text {
		font-size: 28rpx;
		color: #666;
		line-height: 1.6;
		white-space: pre-wrap; /* 保留换行符 */
	}
}

.agreement-footer {
	display: flex;
	justify-content: space-between;
	gap: 20rpx;
	
	button {
		flex: 1;
		height: 80rpx;
		line-height: 80rpx;
		border-radius: 40rpx;
		font-size: 28rpx;
		border: none;
		
		&::after {
			border: none;
		}
	}
	
	.btn-disagree {
		background-color: #f5f5f5;
		color: #666;
	}
	
	.btn-agree {
		background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
		color: #fff;
		box-shadow: 0 4rpx 10rpx rgba(79, 172, 254, 0.3);
	}
}
</style>