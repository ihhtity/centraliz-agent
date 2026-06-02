<template>
	<view class="language-settings">
		<uv-navbar :title="t('user.language.title')" :placeholder="true" @leftClick="goBack" />
		
		<view class="language-list">
			<view 
				v-for="lang in languageList" 
				:key="lang.code"
				class="language-item"
				@click="selectLanguage(lang.code)"
			>
				<text class="language-name">{{ lang.name }}</text>
				<uv-icon 
					v-if="currentLocale === lang.code" 
					name="checkmark" 
					size="24" 
					color="#3c9cff"
				/>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { setLocale, getLocale } from '@/locales/index'

const { t } = useI18n()
const currentLocale = ref(getLocale())

const languageList = [
	{ code: 'zh-CN', name: '简体中文' },
	{ code: 'zh-TW', name: '繁體中文' },
	{ code: 'en-US', name: 'English' },
	{ code: 'ja-JP', name: '日本語' },
	{ code: 'ko-KR', name: '한국어' },
	{ code: 'fr-FR', name: 'Français' },
	{ code: 'de-DE', name: 'Deutsch' },
	{ code: 'es-ES', name: 'Español' },
	{ code: 'ru-RU', name: 'Русский' },
	{ code: 'ar-SA', name: 'العربية' },
	{ code: 'pt-BR', name: 'Português' },
	{ code: 'it-IT', name: 'Italiano' },
	{ code: 'tr-TR', name: 'Türkçe' },
	{ code: 'th-TH', name: 'ไทย' }
]

const goBack = () => {
	uni.navigateBack()
}

const selectLanguage = (code) => {
	if (currentLocale.value === code) {
		return
	}
	
	setLocale(code)
	currentLocale.value = code
	
	const pages = getCurrentPages()
	const prevPage = pages[pages.length - 2]
	if (prevPage && prevPage.$vm) {
		prevPage.$vm.currentLocale = code
	}
	
	uni.showToast({
		title: t('common.operationSuccess'),
		icon: 'success',
		duration: 1500
	})
	
	setTimeout(() => {
		uni.navigateBack()
	}, 1500)
}
</script>

<style scoped lang="scss">
.language-settings {
	min-height: 100vh;
	background-color: #f5f5f5;
	height: 100%;
	overflow-y: auto;
	
	.language-list {
		background-color: #fff;
		
		.language-item {
			display: flex;
			align-items: center;
			justify-content: space-between;
			padding: 32rpx 30rpx;
			border-bottom: 1rpx solid #f0f0f0;
			
			&:last-child {
				border-bottom: none;
			}
			
			&:active {
				background-color: #f9f9f9;
			}
			
			.language-name {
				font-size: 32rpx;
				color: #333;
			}
		}
	}
}
</style>