<template>
	<view class="login-container">
		<!-- 背景装饰 -->
		<view class="bg-decoration"></view>

		<view class="logo-area">
			<view class="logo-icon">
				<uv-icon name="account" size="60" color="#fff"></uv-icon>
			</view>
			<text class="title">{{ t('login.systemTitle') }}</text>
			<text class="subtitle">{{ t('login.systemSubtitle') }}</text>
		</view>

		<view class="form-card">
			<!-- 登录方式切换 -->
			<view class="login-type-switch">
				<view class="type-item" :class="{ active: loginType === 'account' }" @click="loginType = 'account'">
					{{ t('login.byAccount') }}
				</view>
				<view class="type-item" :class="{ active: loginType === 'phone' }" @click="loginType = 'phone'">
					{{ t('login.byPhone') }}
				</view>
				<view class="type-item" :class="{ active: loginType === 'email' }" @click="loginType = 'email'">
					{{ t('login.byEmail') }}
				</view>
			</view>

			<uv-form labelPosition="top" :model="form" :rules="loginRules" ref="formRef">
				<!-- 手机号登录 -->
				<template v-if="loginType === 'phone'">
					<uv-form-item :label="t('login.phone')" borderBottom prop="phone">
						<uv-input v-model="form.phone" :placeholder="t('login.phonePlaceholder')" border="none"
							shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('login.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('login.codePlaceholder')" border="none"
							shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode"
									:disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : (t('login.getCode') || '获取验证码') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>

				<!-- 邮箱登录 -->
				<template v-if="loginType === 'email'">
					<uv-form-item :label="t('login.email')" borderBottom prop="email">
						<uv-input v-model="form.email" type="email"
							:placeholder="t('login.emailPlaceholder') || '请输入邮箱'" border="none"
							shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('login.code')" borderBottom prop="code">
						<uv-input v-model="form.code" :placeholder="t('login.codePlaceholder')" border="none"
							shape="round">
							<template #suffix>
								<uv-button size="mini" shape="round" type="primary" plain @click="getCode"
									:disabled="countdown > 0">
									{{ countdown > 0 ? `${countdown}s` : (t('login.getCode') || '获取验证码') }}
								</uv-button>
							</template>
						</uv-input>
					</uv-form-item>
				</template>

				<!-- 账号密码登录 -->
				<template v-if="loginType === 'account'">
					<uv-form-item :label="t('login.account')" borderBottom prop="account">
						<uv-input v-model="form.account" :placeholder="t('login.accountPlaceholder') || '请输入账号'"
							border="none" shape="round"></uv-input>
					</uv-form-item>
					<uv-form-item :label="t('login.password')" borderBottom prop="password">
						<uv-input v-model="form.password" :type="showPassword ? 'text' : 'password'"
							:placeholder="t('login.passwordPlaceholder') || '请输入密码'" border="none" shape="round">
							<template #suffix>
								<uv-icon :name="showPassword ? 'eye-fill' : 'eye-off'" size="20" color="#999"
									@click="showPassword = !showPassword"></uv-icon>
							</template>
						</uv-input>
					</uv-form-item>
				</template>
			</uv-form>

			<uv-button type="primary" shape="round" @click="handleLogin" class="login-btn">{{ t('login.loginButton')
				}}</uv-button>

			<view class="forgot-password" @click="goToForgotPassword">
				<text>{{ t('login.forgotPassword') || '忘记密码?' }}</text>
			</view>
		</view>

		<view class="switch-role">
			<view class="role-tag" :class="{ active: userRole === 'merch' }" @click="userRole = 'merch'">{{
				t('login.merch') }}
			</view>
			<view class="role-tag" :class="{ active: userRole === 'user' }" @click="userRole = 'user'">{{
				t('login.user') }}
			</view>
		</view>

		<view class="register-link" @click="goToRegister">
			<text>{{ t('login.noAccount') }} {{ t('login.registerNow') }}</text>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, onUnmounted, getCurrentInstance } from 'vue';
// 修改: 引入 onLoad 以接收页面参数
import { onLoad } from '@dcloudio/uni-app';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
// 新增: 获取当前实例以访问全局属性
const { proxy } = getCurrentInstance();

const userRole = ref('merch'); // 'user', 'admin', 'merch'
const loginType = ref('account'); // 'phone', 'email', 'account'
const showPassword = ref(false);
const form = reactive({
	phone: '17727293262',
	email: '2794159940@qq.com',
	account: '17727293262',
	password: '12345678',
	code: ''
});

// 新增: 定义表单 ref
const formRef = ref(null);

const countdown = ref(0);
let timer = null;

// 新增: 接收页面参数
onLoad((options) => {
	if (options.userRole) {
		userRole.value = options.userRole;
	}
});

// 获取验证码
const getCode = async () => {
	uni.showLoading({ title: t('login.sendingCode') || '发送验证码...', duration: 10000 });
	let target = '';

	if (loginType.value === 'phone') {
		target = form.phone;
	} else if (loginType.value === 'email') {
		target = form.email;
	} else {
		uni.showToast({ title: '请选择登录方式', icon: 'none' });
		return;
	}

	if (!target) {
		uni.showToast({
			title: loginType.value === 'phone' ? t('login.phoneRequired') || '请输入手机号' : t('login.emailRequired') || '请输入邮箱',
			icon: 'none'
		});
		return;
	}

	// 防止重复发送
	if (countdown.value > 0) return;

	try {
		const params = loginType.value === 'phone'
			? { phone: target, type: 1 }
			: { email: target, type: 2 };

		await proxy.$http.post('send-code', params)
			.then(res => {
				uni.hideLoading();
				if (res.code !== 200) {
					throw new Error(res.msg);
				}
				uni.showToast({ title: t('login.codeSent') || '验证码已发送', icon: 'success' });

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
		uni.showToast({ title: error.message || t('login.codeSendFailed') || '发送失败', icon: 'none', duration: 3000 });
	}
};

// 新增: 登录表单校验规则
const loginRules = ref({
	phone: [
		{ required: true, message: t('login.phonePlaceholder') || '请输入手机号', trigger: ['blur', 'change'] },
		{ pattern: /^1[3-9]\d{9}$/, message: '手机号格式错误', trigger: ['blur', 'change'] }
	],
	email: [
		{ required: true, message: t('login.emailPlaceholder') || '请输入邮箱', trigger: ['blur', 'change'] },
		{ type: 'email', message: '邮箱格式错误', trigger: ['blur', 'change'] }
	],
	account: [
		{ required: true, message: t('login.accountPlaceholder') || '请输入账号', trigger: ['blur', 'change'] }
	],
	password: [
		{ required: true, message: t('login.passwordPlaceholder') || '请输入密码', trigger: ['blur', 'change'] }
	],
	code: [
		{ required: true, message: t('login.codePlaceholder') || '请输入验证码', trigger: ['blur', 'change'] }
	]
});

// 修改: 登录方法，使用正确的API路径和参数
const handleLogin = async () => {
	try {
		// uv-ui validate 返回 Promise，校验失败会 reject
		await formRef.value?.validate();

		uni.showLoading({ title: t('login.logging') || '登录中...', duration: 10000 });

		let url = 'user/login';
		// 构建统一的登录数据，根据当前登录类型包含相应字段
		const loginData = {
			type: loginType.value,
			role: userRole.value
		};

		if (loginType.value === 'phone') {
			loginData.phone = form.phone;
			loginData.code = form.code;
		} else if (loginType.value === 'email') {
			loginData.email = form.email;
			loginData.code = form.code;
		} else if (loginType.value === 'account') {
			loginData.account = form.account;
			loginData.password = form.password;
		}

		if (userRole.value === 'merch') {
			url = 'merch/login';
		}

		await proxy.$http.post(url, loginData).then(res => {
			uni.hideLoading();
			if (res.code === 200) {
				uni.showToast({ title: t('login.loginSuccess'), icon: 'success' });
				// 保存token和用户信息
				if (res.data && res.data.token) {
					uni.setStorageSync('token', res.data.token);
					// 根据角色保存用户信息
					if (userRole.value === 'merch') {
						uni.setStorageSync('merch', res.data.merch || {});
					} else {
						uni.setStorageSync('user', res.data.user || {});
					}
					// 保存当前角色
					uni.setStorageSync('role', userRole.value);
				}

				// 根据角色跳转不同首页
				setTimeout(() => {
					let targetUrl = '';
					if (userRole.value === 'merch') {
						targetUrl = '/pages/admin/room/manage';
					} else {
						targetUrl = '/pages/user/index/locker';
					}
					uni.reLaunch({ url: targetUrl });
				}, 1500);
			} else {
				throw new Error(res.msg || t('common.requestFailed') || '登录失败');
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
			uni.showToast({ icon: 'none', title: t('common.requestFailed') || '登录失败' });
		}
	}
};

const goToForgotPassword = () => {
	uni.navigateTo({
		url: `/pages/login/forgot-password?userRole=${userRole.value}`
	});
};

const goToRegister = () => {
	uni.navigateTo({
		url: `/pages/login/register?userRole=${userRole.value}`
	});
};

// 新增: 组件卸载时清除定时器
onUnmounted(() => {
	if (timer) {
		clearInterval(timer);
	}
});

</script>

<style lang="scss" scoped>
.login-container {
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

.login-type-switch {
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

.login-btn {
	margin-top: 40rpx;
	font-size: 32rpx;
	font-weight: bold;
	box-shadow: 0 10rpx 20rpx rgba(60, 156, 255, 0.3);
}

.forgot-password {
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

.register-link {
	margin-top: 40rpx;
	text-align: center;
	font-size: 24rpx;
	color: #3c9cff;
}
</style>