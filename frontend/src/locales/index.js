import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN'
import zhTW from './zh-TW'
import enUS from './en-US'
import jaJP from './ja-JP'
import koKR from './ko-KR'
import frFR from './fr-FR'
import deDE from './de-DE'
import esES from './es-ES'
import ruRU from './ru-RU'
import arSA from './ar-SA'
import ptBR from './pt-BR'
import itIT from './it-IT'
import trTR from './tr-TR'
import thTH from './th-TH'

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
		'de-DE': deDE,  // 德语
		'es-ES': esES,  // 西班牙语
		'ru-RU': ruRU,  // 俄语
		'ar-SA': arSA,  // 阿拉伯语
		'pt-BR': ptBR,  // 葡萄牙语
		'it-IT': itIT,  // 意大利语
		'tr-TR': trTR,  // 土耳其语
		'th-TH': thTH   // 泰语
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
