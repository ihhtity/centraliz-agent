<template>
	<view class="register-container">
		<!-- 背景装饰 -->
		<view class="bg-decoration"></view>
		
		<view class="logo-area">
			<view class="logo-icon">
				<uv-icon name="account" size="60" color="#fff"></uv-icon>
			</view>
			<text class="title">{{ t('register.systemTitle') }}</text>
			<text class="subtitle">{{ t('register.systemSubtitle') }}</text>
		</view>
		
		<view class="form-card">
			<!-- 注册方式切换 -->
			<view class="register-type-switch">
				<view 
					class="type-item" 
					:class="{ active: registerType === 'phone' }" 
					@click="registerType = 'phone'"
				>
					{{ t('register.byPhone') }}
				</view>
				<view 
					class="type-item" 
					:class="{ active: registerType === 'email' }" 
					@click="registerType = 'email'"
				>
					{{ t('register.byEmail') }}
				</view>
				<view 
					class="type-item" 
					:class="{ active: registerType === 'account' }" 
					@click="registerType = 'account'"
				>
					{{ t('register.byAccount') }}
				</view>
			</view>
			
			<uv-form labelPosition="top" :model="form" :rules="registerRules" ref="formRef">
				<!-- 手机号注册 -->
				<template v-if="registerType === 'phone'">
					<uv-form-item :label="t('register.phone')" borderBottom prop="phone">
						<uv-input v-model="form.phone" :placeholder="t('register.phonePlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('register.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('register.codePlaceholder')" border="none" shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode" :disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : t('register.getCode') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>
				
				<!-- 邮箱注册 -->
				<template v-if="registerType === 'email'">
					<uv-form-item :label="t('register.email')" borderBottom prop="email">
						<uv-input v-model="form.email" type="email" :placeholder="t('register.emailPlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('register.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('register.codePlaceholder')" border="none" shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode" :disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : t('register.getCode') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>
				
				<!-- 账号密码注册 -->
				<template v-if="registerType === 'account'">
					<uv-form-item :label="t('register.account')" borderBottom prop="account">
						<uv-input v-model="form.account" :placeholder="t('register.accountPlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('register.password')" borderBottom prop="password">
						<uv-input v-model="form.password" type="password" :placeholder="t('register.passwordPlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('register.confirmPassword')" borderBottom prop="confirmPassword">
						<uv-input v-model="form.confirmPassword" type="password" :placeholder="t('register.confirmPasswordPlaceholder')" border="none" shape="round"></uv-input>
					</uv-form-item>
				</template>
			</uv-form>
			
			<uv-button type="primary" shape="round" @click="handleRegister" class="register-btn">{{ t('register.registerButton') }}</uv-button>
			
			<view class="login-link" @click="goToLogin">
				<text>{{ t('register.alreadyHaveAccount') }}</text>
			</view>
		</view>
		
		<view class="switch-role" @click="switchRole">
			<view class="role-tag" :class="{ active: isAdmin }">{{ t('login.admin') }}</view>
			<view class="role-tag" :class="{ active: !isAdmin }">{{ t('login.user') }}</view>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, onUnmounted, getCurrentInstance } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const { proxy } = getCurrentInstance();

const isAdmin = ref(false);
const registerType = ref('phone'); // 'phone', 'email', 'account'
const form = reactive({
	phone: '',
	email: '',
	account: '',
	password: '',
	confirmPassword: '',
	code: ''
});

const formRef = ref(null);
const countdown = ref(0);
let timer = null;

// 表单验证规则
const registerRules = ref({
	phone: [
		{ required: true, message: t('register.phoneRequired'), trigger: ['blur', 'change'] },
		{ pattern: /^1[3-9]\d{9}$/, message: t('register.phoneInvalid'), trigger: ['blur', 'change'] }
	],
	email: [
		{ required: true, message: t('register.emailRequired'), trigger: ['blur', 'change'] },
		{ type: 'email', message: t('register.emailInvalid'), trigger: ['blur', 'change'] }
	],
	account: [
		{ required: true, message: t('register.accountRequired'), trigger: ['blur', 'change'] },
		{ min: 4, max: 20, message: t('register.accountLength'), trigger: ['blur', 'change'] }
	],
	password: [
		{ required: true, message: t('register.passwordRequired'), trigger: ['blur', 'change'] },
		{ pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,20}$/, message: t('register.passwordStrength'), trigger: ['blur', 'change'] }
	],
	confirmPassword: [
		{ required: true, message: t('register.confirmPasswordRequired'), trigger: ['blur', 'change'] },
		{ validator: validateConfirmPassword, trigger: ['blur', 'change'] }
	],
	code: [
		{ required: true, message: t('register.codeRequired'), trigger: ['blur', 'change'] }
	]
});

function validateConfirmPassword(rule, value, callback) {
	if (value !== form.password) {
		callback(new Error(t('register.passwordNotMatch')));
	} else {
		callback();
	}
}

const getCode = async () => {
	let target = '';
	let type = '';
	
	if (registerType.value === 'phone') {
		target = form.phone;
		type = 'phone';
	} else if (registerType.value === 'email') {
		target = form.email;
		type = 'email';
	} else {
		return;
	}
	
	if (!target) {
		uni.showToast({ title: type === 'phone' ? t('register.phoneRequired') : t('register.emailRequired'), icon: 'none' });
		return;
	}

	// 防止重复发送
	if (countdown.value > 0) return;
	
	try {
		const params = registerType.value === 'phone' 
			? { phone: target, type: 1 } 
			: { email: target, type: 2 };
		
		await proxy.$http.post('/v1/user/send-code', params);
		uni.showToast({ title: t('register.codeSent'), icon: 'success' });
		
		// 开始倒计时
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

const handleRegister = async () => {
	try {
		await formRef.value?.validate();
		
		uni.showLoading({ title: t('register.registering') || '注册中...' });
		
		// 构建统一的注册数据，根据当前注册类型包含相应字段
		const registerData = {
			role: isAdmin.value ? 'admin' : 'user'
		};

		if (registerType.value === 'phone') {
			registerData.phone = form.phone;
			registerData.code = form.code;
		} else if (registerType.value === 'email') {
			registerData.email = form.email;
			registerData.code = form.code;
		} else if (registerType.value === 'account') {
			registerData.account = form.account;
			registerData.password = form.password;
		}
		
		await proxy.$http.post('/v1/user/register', registerData);
		
		uni.hideLoading();
		uni.showToast({ title: t('register.registerSuccess'), icon: 'success' });
		
		// 注册成功后自动登录或跳转到登录页
		setTimeout(() => {
			uni.redirectTo({
				url: `/pages/login/login${isAdmin.value ? '?isAdmin=true' : ''}`
			});
		}, 1500);
	} catch (error) {
		uni.hideLoading();
		// 处理表单验证错误 (通常是数组)
		if (Array.isArray(error) && error.length > 0 && error[0].message) {
			uni.showToast({ icon: 'none', title: error[0].message });
		} 
		// 处理 HTTP 请求错误或其他错误对象
		else if (error && error.message) {
			uni.showToast({ icon: 'none', title: error.message });
		} 
		// 兜底提示
		else {
			uni.showToast({ icon: 'none', title: t('common.requestFailed') || '注册失败' });
		}
	}
};

const goToLogin = () => {
	uni.redirectTo({ url: '/pages/login/login' });
};

const switchRole = () => {
	isAdmin.value = !isAdmin.value;
};

onUnmounted(() => {
	if (timer) {
		clearInterval(timer);
	}
});
</script>

<style lang="scss" scoped>
.register-container {
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
		background: linear-gradient(135deg, #3c9cff, #2b85e4);
		border-radius: 30rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		margin: 0 auto 30rpx;
		box-shadow: 0 10rpx 20rpx rgba(60, 156, 255, 0.3);
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

.register-type-switch {
	display: flex;
	justify-content: center;
	gap: 20rpx;
	margin-bottom: 30rpx;
	
	.type-item {
		padding: 10rpx 20rpx;
		border-radius: 20rpx;
		font-size: 24rpx;
		color: #999;
		background: rgba(255, 255, 255, 0.5);
		transition: all 0.3s;
		border: 2rpx solid transparent;
		
		&.active {
			background: #3c9cff;
			color: #fff;
			border-color: #3c9cff;
		}
	}
}
	
::v-deep .uv-form-item__body__left__content__label {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;
	white-space: nowrap;
}

.register-btn {
	margin-top: 40rpx;
	font-size: 32rpx;
	font-weight: bold;
	box-shadow: 0 10rpx 20rpx rgba(60, 156, 255, 0.3);
}

.login-link {
	text-align: right;
	margin-top: 20rpx;
	font-size: 24rpx;
	color: #3c9cff;
	padding-right: 10rpx;
}

.switch-role {
	margin-top: 60rpx;
	display: flex;
	gap: 20rpx;
	z-index: 1;
	
	.role-tag {
		padding: 10rpx 30rpx;
		border-radius: 30rpx;
		font-size: 24rpx;
		color: #999;
		background: rgba(255, 255, 255, 0.5);
		transition: all 0.3s;
		
		&.active {
			background: #3c9cff;
			color: #fff;
			box-shadow: 0 4rpx 10rpx rgba(60, 156, 255, 0.3);
		}
	}
}
</style>