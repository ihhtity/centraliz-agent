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
				<view class="type-item" :class="{ active: resetType === 'phone' }" @click="resetType = 'phone'">
					{{ t('forgotPassword.byPhone') }}
				</view>
				<view class="type-item" :class="{ active: resetType === 'email' }" @click="resetType = 'email'">
					{{ t('forgotPassword.byEmail') }}
				</view>
			</view>

			<uv-form labelPosition="top" :model="form" :rules="resetRules" ref="formRef">
				<!-- 手机号找回 -->
				<template v-if="resetType === 'phone'">
					<uv-form-item :label="t('forgotPassword.phone')" borderBottom prop="phone">
						<uv-input v-model="form.phone" :placeholder="t('forgotPassword.phonePlaceholder')" border="none"
							shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('forgotPassword.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('forgotPassword.codePlaceholder')" border="none"
							shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode"
									:disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : t('forgotPassword.getCode') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>

				<!-- 邮箱找回 -->
				<template v-if="resetType === 'email'">
					<uv-form-item :label="t('forgotPassword.email')" borderBottom prop="email">
						<uv-input v-model="form.email" type="email" :placeholder="t('forgotPassword.emailPlaceholder')"
							border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('forgotPassword.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('forgotPassword.codePlaceholder')" border="none"
							shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode"
									:disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : t('forgotPassword.getCode') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>

				<uv-form-item :label="t('forgotPassword.newPassword')" borderBottom prop="newPassword">
					<uv-input v-model="form.newPassword" :type="showNewPassword ? 'text' : 'password'"
						:placeholder="t('forgotPassword.newPasswordPlaceholder')" border="none" shape="round">
						<template #suffix>
							<uv-icon :name="showNewPassword ? 'eye-fill' : 'eye-off'" size="20" color="#999"
								@click="showNewPassword = !showNewPassword"></uv-icon>
						</template>
					</uv-input>
				</uv-form-item>
				<uv-form-item :label="t('forgotPassword.confirmPassword')" borderBottom prop="confirmPassword">
					<uv-input v-model="form.confirmPassword" :type="showConfirmPassword ? 'text' : 'password'"
						:placeholder="t('forgotPassword.confirmPasswordPlaceholder')" border="none" shape="round">
						<template #suffix>
							<uv-icon :name="showConfirmPassword ? 'eye-fill' : 'eye-off'" size="20" color="#999"
								@click="showConfirmPassword = !showConfirmPassword"></uv-icon>
						</template>
					</uv-input>
				</uv-form-item>
			</uv-form>

			<uv-button type="primary" shape="round" @click="handleResetPassword" class="reset-btn">{{
				t('forgotPassword.resetButton') }}</uv-button>

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
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);
const form = reactive({
	phone: '17727293262',
	email: '2794159940@qq.com',
	account: '17727293262',
	password: '123456789',
	newPassword: '123456789',
	confirmPassword: '123456789',
	code: '',
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
		{ required: true, message: t('register.passwordRequired'), trigger: ['blur', 'change'] },
		{ min: 7, message: t('register.passwordLength'), trigger: ['blur', 'change'] },
		{ pattern: /^[^\u4e00-\u9fa5\uff01-\uff5e]+$/, message: t('register.passwordNoChinese'), trigger: ['blur', 'change'] }
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
	uni.showLoading({ title: t('forgotPassword.sendingCode') || '发送验证码...', duration: 10000 });
	let target = '';
	
	if (resetType.value === 'phone') {
		target = form.phone;
	} else if (resetType.value === 'email') {
		target = form.email;
	} else {
		return;
	}
	
	if (!target) {
		uni.showToast({ 
			title: resetType.value === 'phone' ? t('forgotPassword.phoneRequired') : t('forgotPassword.emailRequired'), 
			icon: 'none' 
		});
		return;
	}

	// 防止重复发送
	if (countdown.value > 0) return;
	
	try {
		const params = resetType.value === 'phone' 
			? { phone: target, type: 1 } 
			: { email: target, type: 2 };
				
		await proxy.$http.post('send-code', params)
		.then(res => {
			uni.hideLoading();
			if (res.code !== 200) {
				throw new Error(res.msg);
			}
			uni.showToast({ title: t('forgotPassword.codeSent'), icon: 'success' });
			
			// 开始倒计时
			countdown.value = 60;
			timer = setInterval(() => {
				countdown.value--;
				if (countdown.value <= 0) {
					clearInterval(timer);
					timer = null;
				}
			}, 1000);
		})
	} catch (error) {
		uni.hideLoading();
		uni.showToast({ title: error.message || t('forgotPassword.codeSendFailed') || '发送失败', icon: 'none', duration: 3000 });
	}
};

const handleResetPassword = async () => {
	try {
		await formRef.value?.validate();
		
		uni.showLoading({ title: t('forgotPassword.resetting') || '重置中...', duration: 10000 });
		
		// 构建重置密码数据
		const resetData = { type: resetType.value };
		
		if (resetType.value === 'phone') {
			resetData.phone = form.phone;
			resetData.code = form.code;
		} else if (resetType.value === 'email') {
			resetData.email = form.email;
			resetData.code = form.code;
		}
		
		resetData.newPassword = form.newPassword;
		resetData.confirmPassword = form.confirmPassword;

		// 根据用户角色确定请求路径
		let url = 'user/reset-password';
		if (proxy.$uv.http.getStorageSync('userRole') === 'merch') {
			url = 'merch/reset-password';
		}

		await proxy.$http.post(url, resetData).then(res => {
			uni.hideLoading();
			if (res.code === 200) {
				uni.showToast({ title: t('forgotPassword.resetSuccess') || '密码重置成功', icon: 'success' });
				setTimeout(() => {
					uni.redirectTo({ url: '/pages/login/login' });
				}, 1500);
			} else {
				throw new Error(res.msg || t('common.requestFailed') || '重置失败');
			}
		});
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
			uni.showToast({ icon: 'none', title: t('common.requestFailed') || '重置失败' });
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
			color: #3c9cff;
			border-bottom-color: #3c9cff;
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
}

.login-link {
	text-align: right;
	margin-top: 20rpx;
	font-size: 24rpx;
	color: #3c9cff;
	padding-right: 10rpx;
}
</style>