<template>
	<view class="container">
		<uv-navbar :title="t('user.profile.basicInfo')" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container">
			<!-- 头像区域 -->
			<view class="avatar-card">
				<view class="avatar-section">
					<view class="avatar-wrapper">
						<view class="avatar-default">
							<uv-icon name="account-fill" color="#fff" size="80" />
						</view>
					</view>
				</view>
			</view>

			<!-- 用户ID -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="order" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">{{ t('user.profile.userId') }}</text>
					</view>
					<text class="info-value">{{ userInfo.id || '-' }}</text>
				</view>
			</view>

			<!-- 账号 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="order" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">账号</text>
					</view>
					<text class="info-value">{{ userInfo.account || '-' }}</text>
				</view>
			</view>

			<!-- 昵称 -->
			<view class="info-card">
				<view class="info-item clickable" @click="openEditNickname">
					<view class="info-label-wrap">
						<uv-icon name="account" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">{{ t('user.profile.nickname') }}</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value">{{ userInfo.name || '-' }}</text>
						<view class="info-arrow">
							<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
						</view>
					</view>
				</view>
			</view>

			<!-- 密码 -->
			<view class="info-card">
				<view class="info-item clickable" @click="openChangePasswordModal">
					<view class="info-label-wrap">
						<uv-icon name="lock" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">密码</text>
					</view>
					<view class="info-arrow">
						<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
					</view>
				</view>
			</view>

			<!-- 手机号 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="phone" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">{{ t('user.profile.phone') }}</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value" :class="{ 'unbind': !userInfo.phone }">{{ formattedPhone || '未绑定'
						}}</text>
						<view v-if="userInfo.phone" class="info-actions">
							<view class="action-btn" @click.stop="handlePhoneUnbind">解绑</view>
							<view class="action-btn primary" @click.stop="showBindPhoneModal = true">换绑</view>
						</view>
						<view v-else class="info-arrow" @click="showBindPhoneModal = true">
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
						<text class="info-label">{{ t('user.profile.email') }}</text>
					</view>
					<view class="info-value-wrap">
						<text class="info-value" :class="{ 'unbind': !userInfo.email }">{{ userInfo.email || '未绑定'
						}}</text>
						<view v-if="userInfo.email" class="info-actions">
							<view class="action-btn" @click.stop="handleEmailUnbind">解绑</view>
							<view class="action-btn primary" @click.stop="showBindEmailModal = true">换绑</view>
						</view>
						<view v-else class="info-arrow" @click="showBindEmailModal = true">
							<uv-icon name="arrow-right" size="28rpx" color="#ccc"></uv-icon>
						</view>
					</view>
				</view>
			</view>

			<!-- 创建时间 -->
			<view class="info-card">
				<view class="info-item">
					<view class="info-label-wrap">
						<uv-icon name="clock" size="28rpx" color="#999"></uv-icon>
						<text class="info-label">{{ t('user.profile.createTime') }}</text>
					</view>
					<text class="info-value">{{ formatTime(userInfo.createdAt) }}</text>
				</view>
			</view>
		</scroll-view>

		<!-- 弹框 -->
		<view>
			<!-- 修改昵称弹窗 -->
			<view v-if="showEditNicknameModal" class="modal-overlay" @click="showEditNicknameModal = false">
				<view class="modal-content" @click.stop>
					<view class="modal-header">
						<text class="modal-title">修改昵称</text>
						<view class="modal-close" @click="showEditNicknameModal = false">
							<uv-icon name="close"></uv-icon>
						</view>
					</view>
					<view class="modal-body">
						<view class="form-item">
							<text class="form-label">{{ t('user.profile.nickname') }}</text>
							<input class="form-input" v-model="nicknameForm.name" placeholder="请输入新昵称"
								:maxlength="20" />
						</view>
					</view>
					<view class="modal-footer">
						<view class="btn-cancel" @click="showEditNicknameModal = false">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" @click="submitEditNickname">保存</view>
					</view>
				</view>
			</view>
			<!-- 绑定手机号弹窗 -->
			<view v-if="showBindPhoneModal" class="modal-overlay" @click="showBindPhoneModal = false">
				<view class="modal-content" @click.stop>
					<view class="modal-header">
						<text class="modal-title">{{ userInfo.phone ? '换绑手机号' : '绑定手机号' }}</text>
						<view class="modal-close" @click="showBindPhoneModal = false">
							<uv-icon name="close"></uv-icon>
						</view>
					</view>
					<view class="modal-body">
						<view class="form-item">
							<text class="form-label">{{ t('user.profile.phone') }}</text>
							<input class="form-input" v-model="phoneForm.phone" placeholder="请输入手机号" />
						</view>
					</view>
					<view class="modal-footer">
						<view class="btn-cancel" @click="showBindPhoneModal = false">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" @click="submitBindPhone">{{ t('common.confirm') }}{{ userInfo.phone ?
							'换绑'
							: '绑定' }}</view>
					</view>
				</view>
			</view>
			<!-- 绑定邮箱弹窗 -->
			<view v-if="showBindEmailModal" class="modal-overlay" @click="showBindEmailModal = false">
				<view class="modal-content" @click.stop>
					<view class="modal-header">
						<text class="modal-title">{{ userInfo.email ? '换绑邮箱' : '绑定邮箱' }}</text>
						<view class="modal-close" @click="showBindEmailModal = false">
							<uv-icon name="close"></uv-icon>
						</view>
					</view>
					<view class="modal-body">
						<view class="form-item">
							<text class="form-label">{{ t('user.profile.email') }}</text>
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
						<view class="btn-cancel" @click="showBindEmailModal = false">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" @click="submitBindEmail">{{ t('common.confirm') }}{{ userInfo.email ?
							'换绑'
							: '绑定' }}</view>
					</view>
				</view>
			</view>
			<!-- 更改密码弹窗 -->
			<view v-if="showChangePasswordModal" class="modal-overlay" @click="showChangePasswordModal = false">
				<view class="modal-content" @click.stop>
					<view class="modal-header">
						<text class="modal-title">更改密码</text>
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
						<view class="btn-cancel" @click="showChangePasswordModal = false">{{ t('common.cancel') }}
						</view>
						<view class="btn-confirm" @click="submitChangePassword">{{ t('common.confirm') }}</view>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { onLoad } from '@dcloudio/uni-app';

onLoad(() => {
	userInfo.value = uni.getStorageSync('user')
	loadUserInfo()
})

// 国际化翻译
const { t } = useI18n()
// 用户信息 - 仅从本地存储获取
const userInfo = ref({})
// 绑定手机号弹窗
const showBindPhoneModal = ref(false)
// 绑定手机号表单
const phoneForm = reactive({
	phone: ''
})
// 绑定邮箱弹窗
const showBindEmailModal = ref(false)
// 绑定邮箱表单
const emailForm = reactive({
	email: '',
	code: '',
	countdown: 0
})
// 绑定邮箱倒计时定时器
let emailTimer = null
// 修改昵称弹窗
const showEditNicknameModal = ref(false)
// 修改昵称表单
const nicknameForm = reactive({
	name: ''
})
// 更改密码弹窗
const showChangePasswordModal = ref(false)
// 更改密码表单
const passwordForm = reactive({
	oldPassword: '',
	newPassword: '',
	confirmPassword: ''
})
// 密码可见性切换
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

// 开始倒计时
const startCountdown = (form, timerRef) => {
	form.countdown = 60
	if (timerRef) {
		clearInterval(timerRef)
	}
	timerRef = setInterval(() => {
		form.countdown--
		if (form.countdown <= 0) {
			clearInterval(timerRef)
			timerRef = null
		}
	}, 1000)
}
// 提交绑定手机号
const submitBindPhone = async () => {
	const phone = phoneForm.phone.trim()
	if (!phone) {
		uni.showToast({ title: '请输入手机号', icon: 'none' })
		return
	}
	// 验证手机号格式
	if (!/^1[3-9]\d{9}$/.test(phone)) {
		uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
		return
	}
	try {
		uni.showLoading({ title: '绑定中...' })
		const result = await uni.$uv.http.put('/user/profile/' + userInfo.value.id, {
			id: userInfo.value.id,
			phone: phone
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '绑定成功', icon: 'success' })
			showBindPhoneModal.value = false
			phoneForm.phone = ''
			loadUserInfo()
		} else {
			uni.showToast({ title: result.msg || '绑定失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '绑定失败', icon: 'none' })
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
			id: userInfo.value.id,
			email: emailForm.email,
			type: 2,
			purpose: userInfo.value.email ? 'rebind' : 'bind',
			role: 'user'
		})
		if (result.code === 200) {
			uni.showToast({ title: '验证码已发送', icon: 'success' })
			startCountdown(emailForm, emailTimer)
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
		const result = await uni.$uv.http.put('/user/profile/email', {
			id: userInfo.value.id,
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
			loadUserInfo()
		} else {
			uni.showToast({ title: result.msg || '绑定失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '绑定失败', icon: 'none' })
	}
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
					const result = await uni.$uv.http.put('/user/profile/' + userInfo.value.id, {
						id: userInfo.value.id,
						phone: ''
					}, { custom: { auth: true } })
					uni.hideLoading()
					if (result.code === 200) {
						uni.showToast({ title: '解绑成功', icon: 'success' })
						loadUserInfo()
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
					const result = await uni.$uv.http.delete('/user/profile/email', {
						id: userInfo.value.id
					}, { custom: { auth: true } })
					uni.hideLoading()
					if (result.code === 200) {
						uni.showToast({ title: '解绑成功', icon: 'success' })
						loadUserInfo()
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
// 打开修改昵称弹窗
const openEditNickname = () => {
	nicknameForm.name = userInfo.value.name || ''
	showEditNicknameModal.value = true
}
// 提交修改昵称
const submitEditNickname = async () => {
	const newName = (nicknameForm.name || '').trim()
	if (!newName) {
		uni.showToast({ title: '请输入新昵称', icon: 'none' })
		return
	}
	if (newName.length > 20) {
		uni.showToast({ title: '昵称不能超过20个字符', icon: 'none' })
		return
	}
	if (newName === userInfo.value.name) {
		showEditNicknameModal.value = false
		return
	}
	try {
		uni.showLoading({ title: '保存中...' })
		const result = await uni.$uv.http.put('/user/profile', {
			id: userInfo.value.id,
			name: newName
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '修改成功', icon: 'success' })
			showEditNicknameModal.value = false
			loadUserInfo()
		} else {
			uni.showToast({ title: result.msg || '修改失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '修改失败', icon: 'none' })
	}
}
// 重新加载用户信息
const loadUserInfo = () => {
	// 从服务端获取最新信息
	uni.$uv.http.get('/user/profile/' + userInfo.value.id, { custom: { auth: true } })
		.then(res => {
			if (res.code === 200 && res.data) {
				userInfo.value = res.data
				uni.setStorageSync('user', res.data)
			}
		})
		.catch(() => { })
}
// 打开更改密码弹窗
const openChangePasswordModal = () => {
	passwordForm.oldPassword = ''
	passwordForm.newPassword = ''
	passwordForm.confirmPassword = ''
	showChangePasswordModal.value = true
}
// 提交更改密码
const submitChangePassword = async () => {
	const oldPassword = passwordForm.oldPassword.trim()
	const newPassword = passwordForm.newPassword.trim()
	const confirmPassword = passwordForm.confirmPassword.trim()

	if (!oldPassword) {
		uni.showToast({ title: '请输入原密码', icon: 'none' })
		return
	}
	if (!newPassword) {
		uni.showToast({ title: '请输入新密码', icon: 'none' })
		return
	}
	if (newPassword.length < 6 || newPassword.length > 20) {
		uni.showToast({ title: '密码长度必须在6-20位之间', icon: 'none' })
		return
	}
	// 检查是否包含中文
	if (/[\u4e00-\u9fa5]/.test(newPassword)) {
		uni.showToast({ title: '密码不能包含中文', icon: 'none' })
		return
	}
	if (newPassword !== confirmPassword) {
		uni.showToast({ title: '两次输入的密码不一致', icon: 'none' })
		return
	}

	try {
		uni.showLoading({ title: '保存中...' })
		const result = await uni.$uv.http.put(`/user/profile/${userInfo.value.id}`, {
			id: userInfo.value.id,
			oldPassword: oldPassword,
			newPassword: newPassword
		}, { custom: { auth: true } })
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '修改成功', icon: 'success' })
			showChangePasswordModal.value = false
		} else {
			uni.showToast({ title: result.msg || '修改失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '修改失败', icon: 'none' })
	}
}
// 格式化手机号
const formattedPhone = computed(() => {
	const phone = userInfo.value.phone || ''
	if (phone.length === 11) {
		return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
	}
	return phone || ''
})
// 格式化时间
const formatTime = (time) => {
	if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
}
// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/user/profile/index'
	});
};
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

.avatar-card {
	background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
	border-radius: 20rpx;
	overflow: hidden;
	margin-bottom: 16rpx;
}

.avatar-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 60rpx 0 40rpx;

	.avatar-wrapper {
		width: 160rpx;
		height: 160rpx;
		border-radius: 50%;
		overflow: hidden;
		border: 4rpx solid rgba(255, 255, 255, 0.8);
		box-shadow: 0 8rpx 20rpx rgba(0, 0, 0, 0.15);

		.avatar-img {
			width: 100%;
			height: 100%;
		}

		.avatar-default {
			width: 100%;
			height: 100%;
			background-color: rgba(255, 255, 255, 0.2);
			display: flex;
			justify-content: center;
			align-items: center;
		}
	}

	.avatar-tip {
		margin-top: 20rpx;
		font-size: 24rpx;
		color: rgba(255, 255, 255, 0.8);
	}
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

	&.unbind {
		color: #999;
	}
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