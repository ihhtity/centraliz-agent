<template>
	<view class="forgot-container">
		<!-- 背景装饰 -->
		<view class="bg-decoration"></view>
		
		<view class="logo-area">
			<view class="logo-icon">
				<uv-icon name="lock" size="60" color="#fff"></uv-icon>
			</view>
			<text class="title">{{ t('forgotPassword.systemTitle') }}</text>
			<text class="subtitle">{{ t('forgotPassword.systemSubtitle') }}</text>
		</view>
		
		<view class="form-card">
			<view class="reset-type-switch">
				<view 
					class="type-item" 
					:class="{ active: resetType === 'phone' }" 
					@click="resetType = 'phone'"
				>
					{{ t('forgotPassword.byPhone') }}
				</view>
				<view 
					class="type-item" 
					:class="{ active: resetType === 'email' }" 
					@click="resetType = 'email'"
				>
					{{ t('forgotPassword.byEmail') }}
				</view>
			</view>
			
			<uv-form labelPosition="top" :model="form" :rules="resetRules" ref="formRef">
				<!-- 手机号找回 -->
				<template v-if="resetType === 'phone'">
					<uv-form-item :label="t('forgotPassword.phone')" borderBottom prop="phone">
						<uv-input v-model="form.phone" :placeholder="t('forgotPassword.phonePlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('forgotPassword.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('forgotPassword.codePlaceholder')" border="none" shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode" :disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : t('forgotPassword.getCode') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>
				
				<!-- 邮箱找回 -->
				<template v-if="resetType === 'email'">
					<uv-form-item :label="t('forgotPassword.email')" borderBottom prop="email">
						<uv-input v-model="form.email" type="email" :placeholder="t('forgotPassword.emailPlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('forgotPassword.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('forgotPassword.codePlaceholder')" border="none" shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode" :disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : t('forgotPassword.getCode') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>
				
				<uv-form-item :label="t('forgotPassword.newPassword')" borderBottom prop="newPassword">
					<uv-input v-model="form.newPassword" type="password" :placeholder="t('forgotPassword.newPasswordPlaceholder')" border="none" shape="round"></uv-input>
				</uv-form-item>
				<uv-form-item :label="t('forgotPassword.confirmPassword')" borderBottom prop="confirmPassword">
					<uv-input v-model="form.confirmPassword" type="password" :placeholder="t('forgotPassword.confirmPasswordPlaceholder')" border="none" shape="round"></uv-input>
				</uv-form-item>
			</uv-form>
			
			<uv-button type="primary" shape="round" @click="handleResetPassword" class="reset-btn">{{ t('forgotPassword.resetButton') }}</uv-button>
			
			<view class="login-link" @click="goToLogin">
				<text>{{ t('forgotPassword.backToLogin') }}</text>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, onUnmounted, getCurrentInstance } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { proxy } = getCurrentInstance();

const resetType = ref('phone'); // 'phone' or 'email'
const form = reactive({
	phone: '',
	email: '',
	code: '',
	newPassword: '',
	confirmPassword: ''
});

const formRef = ref(null);
const countdown = ref(0);
let timer = null;

// 表单验证规则
const resetRules = ref({
	phone: [
		{ required: true, message: t('forgotPassword.phoneRequired'), trigger: ['blur', 'change'] },
		{ pattern: /^1[3-9]\d{9}$/, message: t('forgotPassword.phoneInvalid'), trigger: ['blur', 'change'] }
	],
	email: [
		{ required: true, message: t('forgotPassword.emailRequired'), trigger: ['blur', 'change'] },
		{ type: 'email', message: t('forgotPassword.emailInvalid'), trigger: ['blur', 'change'] }
	],
	code: [
		{ required: true, message: t('forgotPassword.codeRequired'), trigger: ['blur', 'change'] }
	],
	newPassword: [
		{ required: true, message: t('forgotPassword.passwordRequired'), trigger: ['blur', 'change'] },
		{ pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,20}$/, message: t('forgotPassword.passwordStrength'), trigger: ['blur', 'change'] }
	],
	confirmPassword: [
		{ required: true, message: t('forgotPassword.confirmPasswordRequired'), trigger: ['blur', 'change'] },
		{ validator: validateConfirmPassword, trigger: ['blur', 'change'] }
	]
});

function validateConfirmPassword(rule, value, callback) {
	if (value !== form.newPassword) {
		callback(new Error(t('forgotPassword.passwordNotMatch')));
	} else {
		callback();
	}
}

const getCode = async () => {
	let target = '';
	let type = '';
	
	if (resetType.value === 'phone') {
		target = form.phone;
		type = 'phone';
	} else if (resetType.value === 'email') {
		target = form.email;
		type = 'email';
	} else {
		return;
	}
	
	if (!target) {
		uni.showToast({ title: type === 'phone' ? t('forgotPassword.phoneRequired') : t('forgotPassword.emailRequired'), icon: 'none' });
		return;
	}
	
	try {
		const params = type === 'phone' ? { phone: target } : { email: target };
		await proxy.$http.post('/api/v1/auth/send-code', params);
		uni.showToast({ title: t('forgotPassword.codeSent'), icon: 'success' });
		
		countdown.value = 60;
		timer = setInterval(() => {
			countdown.value--;
			if (countdown.value <= 0) {
				clearInterval(timer);
				timer = null;
			}
		}, 1000);
	} catch (error) {
		console.error('发送验证码失败:', error);
	}
};

const handleResetPassword = async () => {
	try {
		await formRef.value?.validate();
		
		let resetData = {};
		if (resetType.value === 'phone') {
			resetData = { phone: form.phone, code: form.code, newPassword: form.newPassword };
		} else if (resetType.value === 'email') {
			resetData = { email: form.email, code: form.code, newPassword: form.newPassword };
		}
		
		uni.showLoading({ title: t('forgotPassword.resetting') });
		await proxy.$http.post('/v1/user/reset-password', resetData);
		uni.hideLoading();
		
		uni.showToast({ title: t('forgotPassword.resetSuccess'), icon: 'success' });
		setTimeout(() => {
			uni.redirectTo({ url: '/pages/login/login' });
		}, 1500);
	} catch (error) {
		uni.hideLoading();
		if (error && error.length && error[0].message) {
			uni.showToast({ icon: 'none', title: error[0].message });
		}
	}
};

const goToLogin = () => {
	uni.redirectTo({ url: '/pages/login/login' });
};

onUnmounted(() => {
	if (timer) {
		clearInterval(timer);
	}
});
</script>

<style lang="scss" scoped>
.forgot-container {
	min-height: 100vh;
	background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
	position: relative;
	overflow: hidden;
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 40rpx;
	box-sizing: border-box;
}

.bg-decoration {
	position: absolute;
	top: -100rpx;
	right: -100rpx;
	width: 400rpx;
	height: 400rpx;
	background: rgba(60, 156, 255, 0.1);
	border-radius: 50%;
	filter: blur(60rpx);
	z-index: 0;
}

.logo-area {
	margin-top: 100rpx;
	margin-bottom: 60rpx;
	text-align: center;
	z-index: 1;
	
	.logo-icon {
		width: 120rpx;
		height: 120rpx;
		background: linear-gradient(135deg, #ff6b6b, #ff4757);
		border-radius: 30rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto 30rpx;
		box-shadow: 0 10rpx 20rpx rgba(255, 107, 107, 0.3);
	}
	
	.title {
		font-size: 44rpx;
		font-weight: bold;
		color: #333;
		display: block;
		margin-bottom: 10rpx;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
	
	.subtitle {
		font-size: 24rpx;
		color: #666;
		letter-spacing: 2rpx;
	}
}

.form-card {
	width: 90%;
	background: #fff;
	border-radius: 30rpx;
	padding: 40rpx;
	box-shadow: 0 10rpx 30rpx rgba(0, 0, 0, 0.05);
	z-index: 1;
}

.reset-type-switch {
	display: flex;
	justify-content: center;
	gap: 40rpx;
	margin-bottom: 30rpx;
	
	.type-item {
		font-size: 28rpx;
		color: #999;
		padding-bottom: 10rpx;
		border-bottom: 4rpx solid transparent;
		transition: all 0.3s;
		
		&.active {
			color: #ff4757;
			border-bottom-color: #ff4757;
			font-weight: bold;
		}
	}
}
	
::v-deep .uv-form-item__body__left__content__label {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;
	white-space: nowrap;
}

.reset-btn {
	margin-top: 40rpx;
	font-size: 32rpx;
	font-weight: bold;
	background: linear-gradient(135deg, #ff6b6b, #ff4757);
	box-shadow: 0 10rpx 20rpx rgba(255, 107, 107, 0.3);
}

.login-link {
	text-align: right;
	margin-top: 20rpx;
	font-size: 24rpx;
	color: #ff4757;
	padding-right: 10rpx;
}
</style>