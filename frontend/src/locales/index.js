import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import zhTW from './zh-TW'
import enUS from './en-US'
import jaJP from './ja-JP'
import koKR from './ko-KR'
import frFR from './fr-FR'
import deDE from './de-DE'

// 从本地存储获取语言设置，默认为简体中文
const getDefaultLocale = () => {
	const locale = uni.getStorageSync('locale')
	return locale || 'zh-CN'
}

const i18n = createI18n({
	legacy: false, // 使用 Composition API 模式
	locale: getDefaultLocale(), // 当前语言
	fallbackLocale: 'zh-CN', // 回退语言
	messages: {
		'zh-CN': zhCN,  // 简体中文
		'zh-TW': zhTW,  // 繁体中文
		'en-US': enUS,  // 英语
		'ja-JP': jaJP,  // 日语
		'ko-KR': koKR,  // 韩语
		'fr-FR': frFR,  // 法语
		'de-DE': deDE   // 德语
	},
	// 静默缺失翻译的警告
	silentTranslationWarn: true,
	silentFallbackWarn: true
})

// 切换语言的工具函数
export const setLocale = (locale) => {
	i18n.global.locale.value = locale
	uni.setStorageSync('locale', locale)
}

// 获取当前语言
export const getLocale = () => {
	return i18n.global.locale.value
}

export default i18n
