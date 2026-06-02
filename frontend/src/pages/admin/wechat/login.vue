<!-- 微信授权页面 -->
<template>
	<view class="page">
		<uv-navbar title="微信授权" :placeholder="true" @leftClick="goBack" />
		<view class="content">
			<view class="card">
				<view class="logo-area">
					<view class="logo">
						<text class="logo-text">W</text>
					</view>
					<text class="title">微信授权</text>
					<text class="desc">获取您的微信身份信息</text>
				</view>
				
				<view class="button-area">
					<view 
						class="auth-btn wx-btn" 
						:class="{ disabled: !isMiniProgramEnv || isAuthLoading }"
						@click="handleMiniProgramLogin"
					>
						<text v-if="isAuthLoading && currentPlatform === 'miniprogram'" class="loading-icon">⏳</text>
						<text v-else class="btn-icon">💬</text>
						<text class="btn-text">{{ isAuthLoading && currentPlatform === 'miniprogram' ? '授权中...' : '小程序登录' }}</text>
					</view>
					
					<view 
						class="auth-btn mp-btn" 
						:class="{ disabled: !isH5Env || isAuthLoading }"
						@click="handleMPLogin"
					>
						<text v-if="isAuthLoading && currentPlatform === 'mp'" class="loading-icon">⏳</text>
						<text v-else class="btn-icon">🌐</text>
						<text class="btn-text">{{ isAuthLoading && currentPlatform === 'mp' ? '授权中...' : '公众号登录' }}</text>
					</view>
				</view>
				
				<view v-if="showResult" class="result-area">
					<view class="result-header">
						<text class="check-icon">✓</text>
					</view>
					<text class="result-title">授权成功</text>
					<view class="info-list">
						<view class="info-item">
							<text class="info-label">OpenID</text>
							<text class="info-value">{{ wechatInfo.openId }}</text>
						</view>
						<view class="info-item">
							<text class="info-label">UnionID</text>
							<text class="info-value">{{ wechatInfo.unionId || '未获取' }}</text>
						</view>
						<view class="info-item">
							<text class="info-label">平台</text>
							<text class="info-value">{{ platformText }}</text>
						</view>
					</view>
					<view class="action-area">
						<view class="action-btn" @click="copyOpenId">复制OpenID</view>
						<view class="action-btn primary" @click="reset">重新授权</view>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';

// 环境检测：区分小程序和H5环境
const isMiniProgramEnv = ref(false);
const isH5Env = ref(false);

// #ifdef MP-WEIXIN
isMiniProgramEnv.value = true;
// #endif

// #ifdef H5
isH5Env.value = true;
// #endif

// 是否正在授权中
const isAuthLoading = ref(false);
// 当前正在操作的平台
const currentPlatform = ref('');

// 返回上一页
const goBack = () => {
	uni.navigateBack()
}

// 微信用户信息
const wechatInfo = reactive({
	openId: '',
	unionId: '',
	platform: ''
});

// 是否显示授权结果
const showResult = computed(() => wechatInfo.openId !== '');

// 平台文本描述
const platformText = computed(() => {
	return wechatInfo.platform === 'miniprogram' ? '微信小程序' : '微信公众号';
});

// code有效期5分钟（微信官方规定）
const CODE_EXPIRE_TIME = 5 * 60 * 1000;

// 页面加载时检查URL中是否有公众号授权code
onMounted(() => {
	// #ifdef H5
	const code = new URLSearchParams(window.location.search).get('code');
	if (code && !wechatInfo.openId) {
		// 自动触发公众号授权处理
		setTimeout(() => {
			handleMPLogin();
		}, 500); // 延迟500ms，让页面完全渲染后再开始授权
	}
	// #endif
});

// 小程序登录处理
const handleMiniProgramLogin = async () => {
	if (!isMiniProgramEnv.value || isAuthLoading.value) return;
	
	isAuthLoading.value = true;
	currentPlatform.value = 'miniprogram';
	
	try {
		const code = await getMiniProgramCode();
		if (code) {
			const result = await fetchWechatUserInfo(code, 'miniprogram');
			if (result.success) {
				uni.showToast({ title: '授权成功', icon: 'success', duration: 1500 });
			}
		} else {
			uni.showToast({ title: '获取code失败', icon: 'none', duration: 2000 });
		}
	} catch (error) {
		uni.showToast({ title: '授权失败', icon: 'none', duration: 2000 });
		console.error('小程序登录失败:', error);
	} finally {
		isAuthLoading.value = false;
		currentPlatform.value = '';
	}
};

// 获取小程序code，优先使用缓存
const getMiniProgramCode = async () => {
	const cachedCode = uni.getStorageSync('miniprogram_code');

	// 检查缓存是否有效且未过期
	if (cachedCode?.code && cachedCode?.timestamp) {
		const now = Date.now();
		if (now - cachedCode.timestamp < CODE_EXPIRE_TIME) {
			return cachedCode.code;
		}
	}

	// 重新获取code
	const res = await uni.login({ provider: 'weixin' });
	if (res.code) {
		uni.setStorageSync('miniprogram_code', {
			code: res.code,
			timestamp: Date.now()
		});
		return res.code;
	}
	return null;
};

// 公众号登录处理
const handleMPLogin = async () => {
	if (!isH5Env.value || isAuthLoading.value) return;
	
	isAuthLoading.value = true;
	currentPlatform.value = 'mp';
	
	try {
		const appId = 'wxfc42adf0c2f58bb6';
		const currentUrl = window.location.href;
		const urlObj = new URL(currentUrl);
		const code = urlObj.searchParams.get('code');

		if (code) {
			// 已有code，调用后端获取用户信息
			const result = await fetchWechatUserInfo(code, 'mp');
			// 无论成功失败都清理URL中的code
			urlObj.searchParams.delete('code');
			urlObj.searchParams.delete('state');
			window.history.replaceState({}, document.title, urlObj.toString());
			// 显示错误信息
			if (!result.success) {
				uni.showToast({ title: result.msg || '获取失败', icon: 'none', duration: 2000 });
			}
		} else {
			// 无code，跳转微信授权页面
			const redirectUri = encodeURIComponent(window.location.href);
			const oauthUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${redirectUri}&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect`;
			window.location.href = oauthUrl;
		}
	} catch (error) {
		uni.showToast({ title: '授权失败', icon: 'none', duration: 2000 });
		console.error('公众号登录失败:', error);
	} finally {
		isAuthLoading.value = false;
		currentPlatform.value = '';
	}
};

// 调用后端获取微信用户信息，返回结果对象
const fetchWechatUserInfo = async (code, platform) => {
	try {
		const res = await uni.$uv.http.post('/wechat/login', { code, platform });

		if (res.code === 200 && res.data) {
			wechatInfo.openId = res.data.openId;
			wechatInfo.unionId = res.data.unionId || '';
			wechatInfo.platform = platform;
			uni.setStorageSync('wechatUserInfo', res.data);
			return { success: true, msg: '' };
		} else {
			// code无效或过期时清除缓存
			if (res.msg?.includes('40029')) {
				uni.removeStorageSync('miniprogram_code');
				return { success: false, msg: 'code已过期，请重试' };
			}
			return { success: false, msg: res.msg || '获取失败' };
		}
	} catch (error) {
		const errorMsg = error?.msg || error?.message || '';
		if (errorMsg.includes('40029') || errorMsg.includes('code无效')) {
			uni.removeStorageSync('miniprogram_code');
			return { success: false, msg: 'code已过期，请重试' };
		}
		return { success: false, msg: '网络请求失败' };
	}
};

// 复制OpenID到剪贴板
const copyOpenId = async () => {
	if (!wechatInfo.openId) return;
	try {
		await uni.setClipboardData({ data: wechatInfo.openId });
		uni.showToast({ title: '复制成功', icon: 'success', duration: 1500 });
	} catch (error) {
		uni.showToast({ title: '复制失败', icon: 'none', duration: 2000 });
	}
};

// 重置授权状态
const reset = () => {
	wechatInfo.openId = '';
	wechatInfo.unionId = '';
	wechatInfo.platform = '';
	uni.removeStorageSync('wechatUserInfo');
	uni.removeStorageSync('miniprogram_code');
};
</script>

<style lang="scss" scoped>
.page {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
	padding: 40rpx;
	box-sizing: border-box;
}

.content {
	width: 100%;
	max-width: 560rpx;
}

.card {
	background: #fff;
	border-radius: 32rpx;
	padding: 60rpx 40rpx;
	box-shadow: 0 20rpx 60rpx rgba(0, 0, 0, 0.15);
}

.logo-area {
	text-align: center;
	margin-bottom: 60rpx;
}

.logo {
	width: 100rpx;
	height: 100rpx;
	background: linear-gradient(135deg, #667eea, #764ba2);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 0 auto 30rpx;
	box-shadow: 0 8rpx 24rpx rgba(102, 126, 234, 0.4);
}

.logo-text {
	font-size: 48rpx;
	font-weight: bold;
	color: #fff;
}

.title {
	display: block;
	font-size: 40rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 16rpx;
}

.desc {
	display: block;
	font-size: 26rpx;
	color: #999;
}

.button-area {
	display: flex;
	flex-direction: column;
	gap: 24rpx;
}

.auth-btn {
	display: flex;
	align-items: center;
	justify-content: center;
	height: 96rpx;
	border-radius: 48rpx;
	transition: all 0.2s ease;
	
	&.wx-btn {
		background: linear-gradient(135deg, #07c160, #06ad56);
		
		.btn-icon, .btn-text, .loading-icon {
			color: #fff;
		}
	}
	
	&.mp-btn {
		background: #f5f7fa;
		border: 2rpx solid #e8ecf0;
		
		.btn-icon, .btn-text, .loading-icon {
			color: #576b95;
		}
	}
	
	&.disabled {
		opacity: 0.4;
		pointer-events: none;
	}
	
	&:active:not(.disabled) {
		transform: scale(0.98);
		opacity: 0.9;
	}
	
	.btn-icon, .loading-icon {
		font-size: 32rpx;
		margin-right: 12rpx;
	}
	
	.btn-text {
		font-size: 30rpx;
		font-weight: 500;
	}
}

.result-area {
	margin-top: 40rpx;
	padding-top: 40rpx;
	border-top: 1rpx solid #f0f2f5;
	text-align: center;
}

.result-header {
	width: 80rpx;
	height: 80rpx;
	background: linear-gradient(135deg, #07c160, #06ad56);
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	margin: 0 auto 24rpx;
}

.check-icon {
	font-size: 40rpx;
	color: #fff;
	font-weight: bold;
}

.result-title {
	display: block;
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
	margin-bottom: 32rpx;
}

.info-list {
	background: #f8f9fa;
	border-radius: 16rpx;
	padding: 8rpx 0;
	margin-bottom: 32rpx;
}

.info-item {
	display: flex;
	align-items: center;
	padding: 20rpx 24rpx;
	border-bottom: 1rpx solid #e8ecf0;
	
	&:last-child {
		border-bottom: none;
	}
	
	.info-label {
		font-size: 26rpx;
		color: #999;
		min-width: 100rpx;
		flex-shrink: 0;
	}
	
	.info-value {
		font-size: 26rpx;
		color: #333;
		font-family: -apple-system, BlinkMacSystemFont, monospace;
		word-break: break-all;
		flex: 1;
		text-align: right;
	}
}

.action-area {
	display: flex;
	gap: 20rpx;
	
	.action-btn {
		flex: 1;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 40rpx;
		font-size: 28rpx;
		background: #f5f7fa;
		color: #666;
		transition: all 0.2s;
		
		&:active {
			opacity: 0.7;
		}
		
		&.primary {
			background: linear-gradient(135deg, #667eea, #764ba2);
			color: #fff;
		}
	}
}
</style>
