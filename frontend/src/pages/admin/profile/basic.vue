<template>
	<view class="container">
		<uv-navbar title="基本信息" :placeholder="true" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container">
			<!-- 账号信息 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="account" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">账号</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value">{{ profile.account || '-' }}</text>
					</view>
				</view>
			</view>

			<!-- 密码 -->
			<view class="info-card">
				<view class="info-item clickable" @click="showChangePasswordModal = true">
					<view class="info-label-wrap">
						<uv-icon name="lock" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">密码</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value password">******</text>
						<view class="info-arrow">
							<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
						</view>
					</view>
				</view>
			</view>

			<!-- 邮箱 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="email" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">邮箱</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value" :class="{ 'unbind': !profile.email }">{{ profile.email || '未绑定'
							}}</text>
						<view v-if="profile.email" class="info-actions">
							<view class="action-btn" @click.stop="handleEmailUnbind">解绑</view>
							<view class="action-btn primary" @click.stop="showBindEmailModal = true">换绑</view>
						</view>
						<view v-else class="info-arrow" @click="showBindEmailModal = true">
							<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
						</view>
					</view>
				</view>
			</view>

			<!-- 手机号 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="phone" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">手机号</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value" :class="{ 'unbind': !profile.phone }">{{ profile.phone || '未绑定'
							}}</text>
						<view v-if="profile.phone" class="info-actions">
							<view class="action-btn" @click.stop="handlePhoneUnbind">解绑</view>
							<view class="action-btn primary" @click.stop="showBindPhoneModal = true">换绑</view>
						</view>
						<view v-else class="info-arrow" @click="showBindPhoneModal = true">
							<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
						</view>
					</view>
				</view>
			</view>

			<!-- 语言设置 -->
			<view class="info-card">
				<view class="info-item clickable" @click="goToLanguage">
					<view class="info-label-wrap">
						<uv-icon name="moments" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">语言</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value">{{ currentLanguageName }}</text>
						<view class="info-arrow">
							<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
						</view>
					</view>
				</view>
			</view>

			<!-- 创建时间 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="calendar" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">注册时间</text>
					</view>
					<text class="info-value">{{ formatTime(profile.createdAt) }}</text>
				</view>
			</view>
		</scroll-view>

		<!-- 修改密码弹窗 -->
		<view v-if="showChangePasswordModal" class="modal-overlay" @click="showChangePasswordModal = false">
			<view class="modal-content" @click.stop>
				<view class="modal-header">
					<text class="modal-title">修改密码</text>
					<view class="modal-close" @click="showChangePasswordModal = false">
						<uv-icon name="close"></uv-icon>
					</view>
				</view>
				<view class="modal-body">
					<view class="form-item">
						<text class="form-label">原密码</text>
						<view class="password-input-wrap">
							<input class="form-input" :type="showOldPassword ? 'text' : 'password'"
								v-model="passwordForm.oldPassword" placeholder="请输入原密码" />
							<uv-icon class="eye-icon" :name="showOldPassword ? 'eye-fill' : 'eye-off'" size="40rpx"
								color="#999" @click="showOldPassword = !showOldPassword"></uv-icon>
						</view>
					</view>
					<view class="form-item">
						<text class="form-label">新密码</text>
						<view class="password-input-wrap">
							<input class="form-input" :type="showNewPassword ? 'text' : 'password'"
								v-model="passwordForm.newPassword" placeholder="请输入新密码" />
							<uv-icon class="eye-icon" :name="showNewPassword ? 'eye-fill' : 'eye-off'" size="40rpx"
								color="#999" @click="showNewPassword = !showNewPassword"></uv-icon>
						</view>
					</view>
					<view class="form-item">
						<text class="form-label">确认密码</text>
						<view class="password-input-wrap">
							<input class="form-input" :type="showConfirmPassword ? 'text' : 'password'"
								v-model="passwordForm.confirmPassword" placeholder="请确认新密码" />
							<uv-icon class="eye-icon" :name="showConfirmPassword ? 'eye-fill' : 'eye-off'" size="40rpx"
								color="#999" @click="showConfirmPassword = !showConfirmPassword"></uv-icon>
						</view>
					</view>
				</view>
				<view class="modal-footer">
					<view class="btn-cancel" @click="showChangePasswordModal = false">取消</view>
					<view class="btn-confirm" @click="submitChangePassword">确认修改</view>
				</view>
			</view>
		</view>

		<!-- 绑定邮箱弹窗 -->
		<view v-if="showBindEmailModal" class="modal-overlay" @click="showBindEmailModal = false">
			<view class="modal-content" @click.stop>
				<view class="modal-header">
					<text class="modal-title">{{ profile.email ? '换绑邮箱' : '绑定邮箱' }}</text>
					<view class="modal-close" @click="showBindEmailModal = false">
						<uv-icon name="close"></uv-icon>
					</view>
				</view>
				<view class="modal-body">
					<view class="form-item">
						<text class="form-label">邮箱地址</text>
						<input class="form-input" v-model="emailForm.email" placeholder="请输入邮箱地址" />
					</view>
					<view class="form-item code-item">
						<input class="form-input code-input" v-model="emailForm.code" placeholder="请输入验证码" />
						<view class="send-code-btn" :class="{ disabled: emailForm.countdown > 0 }"
							@click="sendEmailCode">
							{{ emailForm.countdown > 0 ? emailForm.countdown + 's' : '发送验证码' }}
						</view>
					</view>
				</view>
				<view class="modal-footer">
					<view class="btn-cancel" @click="showBindEmailModal = false">取消</view>
					<view class="btn-confirm" @click="submitBindEmail">确认{{ profile.email ? '换绑' : '绑定' }}</view>
				</view>
			</view>
		</view>

		<!-- 绑定手机号弹窗 -->
		<view v-if="showBindPhoneModal" class="modal-overlay" @click="showBindPhoneModal = false">
			<view class="modal-content" @click.stop>
				<view class="modal-header">
					<text class="modal-title">{{ profile.phone ? '换绑手机号' : '绑定手机号' }}</text>
					<view class="modal-close" @click="showBindPhoneModal = false">
						<uv-icon name="close"></uv-icon>
					</view>
				</view>
				<view class="modal-body">
					<view class="form-item">
						<text class="form-label">手机号</text>
						<input class="form-input" v-model="phoneForm.phone" placeholder="请输入手机号" />
					</view>
					<view class="form-item code-item">
						<input class="form-input code-input" v-model="phoneForm.code" placeholder="请输入验证码" />
						<view class="send-code-btn" :class="{ disabled: phoneForm.countdown > 0 }"
							@click="sendPhoneCode">
							{{ phoneForm.countdown > 0 ? phoneForm.countdown + 's' : '发送验证码' }}
						</view>
					</view>
				</view>
				<view class="modal-footer">
					<view class="btn-cancel" @click="showBindPhoneModal = false">取消</view>
					<view class="btn-confirm" @click="submitBindPhone">确认{{ profile.phone ? '换绑' : '绑定' }}</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getLocale } from '@/locales/index'

const profile = ref({
	account: '',
	email: '',
	phone: '',
	createdAt: ''
})

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

const currentLocale = ref(getLocale())

const currentLanguageName = computed(() => {
	const lang = languageList.find(l => l.code === currentLocale.value)
	return lang ? lang.name : '简体中文'
})

const showChangePasswordModal = ref(false)
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)
const passwordForm = reactive({
	oldPassword: '',
	newPassword: '',
	confirmPassword: ''
})

const showBindEmailModal = ref(false)
const emailForm = reactive({
	email: '',
	code: '',
	countdown: 0
})
let emailTimer = null

const showBindPhoneModal = ref(false)
const phoneForm = reactive({
	phone: '',
	code: '',
	countdown: 0
})
let phoneTimer = null
const merchsId = ref('')

onLoad((options) => {
	merchsId.value = options.merchs_id || ''
	if (!merchsId.value) {
		uni.showToast({ title: '请先登录', icon: 'none' })
		uni.navigateBack()
		return
	}
	loadProfile()
})

// 获取商户信息
const loadProfile = async () => {
	uni.showLoading({ title: '加载中...' })
	try {
		const res = await uni.$uv.http.get('/merch/profile', {
			params: {
				merchs_id: merchsId.value
			},
			custom: { auth: true }
		})
		if (res.code === 200 && res.data) {
			profile.value = res.data
		}
	} catch (e) {
		console.error('加载失败', e)
	} finally {
		uni.hideLoading()
	}
}

// 返回上一页
const goBack = () => uni.navigateBack()

// 切换语言
const goToLanguage = () => {
	uni.navigateTo({
		url: '/pages/user/language/index',
		events: {
			languageChanged: (data) => {
				currentLocale.value = data.code
			}
		}
	})
}

// 格式化时间
const formatTime = (time) => {
	if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
}

// 开始倒计时
const startCountdown = (form, timerRef) => {
	form.countdown = 60
	if (timerRef.value) {
		clearInterval(timerRef.value)
	}
	timerRef.value = setInterval(() => {
		form.countdown--
		if (form.countdown <= 0) {
			clearInterval(timerRef.value)
			timerRef.value = null
		}
	}, 1000)
}

// 提交修改密码
const submitChangePassword = async () => {
	if (!passwordForm.oldPassword) {
		uni.showToast({ title: '请输入原密码', icon: 'none' })
		return
	}
	if (!passwordForm.newPassword) {
		uni.showToast({ title: '请输入新密码', icon: 'none' })
		return
	}
	if (passwordForm.newPassword !== passwordForm.confirmPassword) {
		uni.showToast({ title: '两次密码不一致', icon: 'none' })
		return
	}
	try {
		uni.showLoading({ title: '修改中...' })
		const result = await uni.$uv.http.post('/merch/profile/password', {
			merchs_id: merchsId.value,
			oldPassword: passwordForm.oldPassword,
			newPassword: passwordForm.newPassword
		}, {
			custom: { auth: true }
		})
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '修改成功', icon: 'success' })
			showChangePasswordModal.value = false
			passwordForm.oldPassword = ''
			passwordForm.newPassword = ''
			passwordForm.confirmPassword = ''
		} else {
			uni.showToast({ title: result.msg || '修改失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '修改失败', icon: 'none' })
	}
}

// 发送邮箱验证码
const sendEmailCode = async () => {
	if (emailForm.countdown > 0) return
	if (!emailForm.email) {
		uni.showToast({ title: '请输入邮箱', icon: 'none' })
		return
	}
	try {
		const result = await uni.$uv.http.post('/send-code', {
			merchs_id: merchsId.value,
			email: emailForm.email,
			type: 2,
			purpose: profile.value.email ? 'rebind' : 'bind'
		})
		if (result.code === 200) {
			uni.showToast({ title: '验证码已发送', icon: 'success' })
			startCountdown(emailForm, { value: emailTimer })
		} else {
			uni.showToast({ title: result.msg || '发送失败', icon: 'none' })
		}
	} catch (e) {
		uni.showToast({ title: '发送失败', icon: 'none' })
	}
}

// 提交绑定邮箱
const submitBindEmail = async () => {
	if (!emailForm.email) {
		uni.showToast({ title: '请输入邮箱', icon: 'none' })
		return
	}
	if (!emailForm.code) {
		uni.showToast({ title: '请输入验证码', icon: 'none' })
		return
	}
	try {
		uni.showLoading({ title: '绑定中...' })
		const result = await uni.$uv.http.put('/merch/profile/email', {
			merchs_id: merchsId.value,
			email: emailForm.email,
			code: emailForm.code
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '绑定成功', icon: 'success' })
			showBindEmailModal.value = false
			emailForm.email = ''
			emailForm.code = ''
			emailForm.countdown = 0
			if (emailTimer) {
				clearInterval(emailTimer)
				emailTimer = null
			}
			loadProfile()
		} else {
			uni.showToast({ title: result.msg || '绑定失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '绑定失败', icon: 'none' })
	}
}

// 发送手机号验证码
const sendPhoneCode = async () => {
	if (phoneForm.countdown > 0) return
	if (!phoneForm.phone) {
		uni.showToast({ title: '请输入手机号', icon: 'none' })
		return
	}
	try {
		const result = await uni.$uv.http.post('/send-code', {
			merchs_id: merchsId.value,
			phone: phoneForm.phone,
			type: 1,
			purpose: profile.value.phone ? 'rebind' : 'bind'
		})
		if (result.code === 200) {
			uni.showToast({ title: '验证码已发送', icon: 'success' })
			startCountdown(phoneForm, { value: phoneTimer })
		} else {
			uni.showToast({ title: result.msg || '发送失败', icon: 'none' })
		}
	} catch (e) {
		uni.showToast({ title: '发送失败', icon: 'none' })
	}
}

// 提交绑定手机号
const submitBindPhone = async () => {
	if (!phoneForm.phone) {
		uni.showToast({ title: '请输入手机号', icon: 'none' })
		return
	}
	if (!phoneForm.code) {
		uni.showToast({ title: '请输入验证码', icon: 'none' })
		return
	}
	try {
		uni.showLoading({ title: '绑定中...' })
		const result = await uni.$uv.http.put('/merch/profile/phone', {
			merchs_id: merchsId.value,
			phone: phoneForm.phone,
			code: phoneForm.code
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '绑定成功', icon: 'success' })
			showBindPhoneModal.value = false
			phoneForm.phone = ''
			phoneForm.code = ''
			phoneForm.countdown = 0
			if (phoneTimer) {
				clearInterval(phoneTimer)
				phoneTimer = null
			}
			loadProfile()
		} else {
			uni.showToast({ title: result.msg || '绑定失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '绑定失败', icon: 'none' })
	}
}

// 解绑邮箱
const handleEmailUnbind = () => {
	uni.showModal({
		title: '确认解绑',
		content: '确定要解绑邮箱吗？',
		confirmColor: '#ff4d4f',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '解绑中...' })
					const result = await uni.$uv.http.delete('/merch/profile/email', {
						merchs_id: merchsId.value,
						custom: { auth: true }
					})
					uni.hideLoading()
					if (result.code === 200) {
						uni.showToast({ title: '解绑成功', icon: 'success' })
						loadProfile()
					} else {
						uni.showToast({ title: result.msg || '解绑失败', icon: 'none' })
					}
				} catch (e) {
					uni.hideLoading()
					uni.showToast({ title: '解绑失败', icon: 'none' })
				}
			}
		}
	})
}

// 解绑手机号
const handlePhoneUnbind = () => {
	uni.showModal({
		title: '确认解绑',
		content: '确定要解绑手机号吗？',
		confirmColor: '#ff4d4f',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '解绑中...' })
					const result = await uni.$uv.http.delete('/merch/profile/phone', {
						merchs_id: merchsId.value,
						custom: { auth: true }
					})
					uni.hideLoading()
					if (result.code === 200) {
						uni.showToast({ title: '解绑成功', icon: 'success' })
						loadProfile()
					} else {
						uni.showToast({ title: result.msg || '解绑失败', icon: 'none' })
					}
				} catch (e) {
					uni.hideLoading()
					uni.showToast({ title: '解绑失败', icon: 'none' })
				}
			}
		}
	})
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.scroll-container {
	height: calc(100vh - 88rpx);
	width: calc(100% - 48rpx);
	padding: 24rpx;
}

.info-card {
	background: #fff;
	border-radius: 16rpx;
	padding: 0 24rpx;
	margin-bottom: 16rpx;
}

.info-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	height: 100rpx;
	border-bottom: 1rpx solid #f0f0f0;

	&:last-child {
		border-bottom: none;
	}

	&.clickable {
		&:active {
			background: #f9f9f9;
		}
	}
}

.info-label-wrap {
	display: flex;
	align-items: center;
	gap: 12rpx;
}

.info-label {
	font-size: 28rpx;
	color: #666;
}

.info-value-wrap {
	display: flex;
	align-items: center;
	gap: 12rpx;
}

.info-value {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;

	&.password {
		color: #999;
	}

	&.unbind {
		color: #999;
	}
}

.info-tip {
	font-size: 22rpx;
	color: #999;
	padding: 4rpx 12rpx;
	background: #f5f5f5;
	border-radius: 8rpx;
}

.info-arrow {
	padding: 8rpx;
}

.info-actions {
	display: flex;
	gap: 16rpx;
}

.action-btn {
	font-size: 24rpx;
	color: #999;
	padding: 6rpx 16rpx;
	border-radius: 8rpx;
	background: #f5f5f5;

	&.primary {
		color: #2979ff;
		background: rgba(41, 121, 255, 0.1);
	}
}

.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.modal-content {
	width: 600rpx;
	background: #fff;
	border-radius: 24rpx;
	overflow: hidden;
}

.modal-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 32rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.modal-title {
	font-size: 32rpx;
	font-weight: 600;
	color: #1a1a1a;
}

.modal-close {
	width: 48rpx;
	height: 48rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	color: #999;
}

.modal-body {
	padding: 32rpx;
}

.form-item {
	margin-bottom: 28rpx;

	&.code-item {
		display: flex;
		gap: 16rpx;
	}
}

.form-label {
	font-size: 26rpx;
	color: #666;
	margin-bottom: 12rpx;
	display: block;
}

.form-input {
	width: 100%;
	height: 80rpx;
	padding: 0 20rpx;
	border: 1rpx solid #e8e8e8;
	border-radius: 12rpx;
	font-size: 28rpx;
	box-sizing: border-box;

	&.code-input {
		flex: 1;
	}
}

.password-input-wrap {
	position: relative;
	display: flex;
	align-items: center;
}

.password-input-wrap .form-input {
	padding-right: 70rpx;
}

.eye-icon {
	position: absolute;
	right: 20rpx;
	z-index: 10;
}

.send-code-btn {
	width: 200rpx;
	height: 80rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 26rpx;
	color: #2979ff;
	border: 1rpx solid #2979ff;
	border-radius: 12rpx;

	&.disabled {
		color: #999;
		border-color: #e8e8e8;
	}
}

.modal-footer {
	display: flex;
	border-top: 1rpx solid #f0f0f0;
}

.btn-cancel,
.btn-confirm {
	flex: 1;
	height: 88rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	font-size: 30rpx;
}

.btn-cancel {
	color: #666;
	border-right: 1rpx solid #f0f0f0;
}

.btn-confirm {
	color: #2979ff;
	font-weight: 500;
}
</style>