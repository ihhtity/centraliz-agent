<!-- 功能测试页面 -->
<template>
    <view class="container">
        <uv-navbar :title="'功能测试'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />

        <view class="section">
            <view class="section-title">微信授权</view>
            <view class="auth-card">
                <view 
                    class="auth-item" 
                    :class="{ disabled: !isMiniProgramEnv || isAuthLoading }"
                    @click="handleMiniProgramLogin"
                >
                    <view class="auth-icon mini-icon">💬</view>
                    <view class="auth-info">
                        <text class="auth-title">小程序登录</text>
                        <text class="auth-desc">微信小程序授权登录</text>
                    </view>
                    <view class="auth-arrow">
                        <text v-if="isAuthLoading && currentPlatform === 'miniprogram'" class="loading-text">授权中...</text>
                        <uv-icon v-else name="arrow-right" size="20" color="#ccc" />
                    </view>
                </view>
                
                <view 
                    class="auth-item" 
                    :class="{ disabled: !isH5Env || isAuthLoading }"
                    @click="handleMPLogin"
                >
                    <view class="auth-icon mp-icon">🌐</view>
                    <view class="auth-info">
                        <text class="auth-title">公众号登录</text>
                        <text class="auth-desc">微信公众号授权登录</text>
                    </view>
                    <view class="auth-arrow">
                        <text v-if="isAuthLoading && currentPlatform === 'mp'" class="loading-text">授权中...</text>
                        <uv-icon v-else name="arrow-right" size="20" color="#ccc" />
                    </view>
                </view>
            </view>
            
            <!-- 授权结果 -->
            <view v-if="showResult" class="result-card">
                <view class="result-header">
                    <view class="success-icon">✓</view>
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
                <view class="action-row">
                    <view class="action-btn" @click="copyOpenId">复制OpenID</view>
                    <view class="action-btn primary" @click="reset">重新授权</view>
                </view>
            </view>
        </view>

        <uv-cell-group>
            <uv-cell title="个人设置" :isLink="true" @click="goToSetting" />
            <uv-cell title="版本更新" :isLink="true" @click="goToUpdate" />
        </uv-cell-group>
    </view>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue';

// 环境检测
const isMiniProgramEnv = ref(false);
const isH5Env = ref(false);

// #ifdef MP-WEIXIN
isMiniProgramEnv.value = true;
// #endif

// #ifdef H5
isH5Env.value = true;
// #endif

// 授权状态
const isAuthLoading = ref(false);
const currentPlatform = ref('');

// 微信用户信息
const wechatInfo = reactive({
    openId: '',
    unionId: '',
    platform: ''
});

// 是否显示授权结果
const showResult = computed(() => wechatInfo.openId !== '');

// 平台文本
const platformText = computed(() => {
    return wechatInfo.platform === 'miniprogram' ? '微信小程序' : '微信公众号';
});

// code有效期
const CODE_EXPIRE_TIME = 5 * 60 * 1000;

// 页面加载时检查公众号授权code
onMounted(() => {
    // #ifdef H5
    const code = new URLSearchParams(window.location.search).get('code');
    if (code && !wechatInfo.openId) {
        setTimeout(() => {
            handleMPLogin();
        }, 500);
    }
    // #endif
});

// 小程序登录
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

// 获取小程序code
const getMiniProgramCode = async () => {
    const cachedCode = uni.getStorageSync('miniprogram_code');
    
    if (cachedCode?.code && cachedCode?.timestamp) {
        const now = Date.now();
        if (now - cachedCode.timestamp < CODE_EXPIRE_TIME) {
            return cachedCode.code;
        }
    }
    
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

// 公众号登录
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
            const result = await fetchWechatUserInfo(code, 'mp');
            urlObj.searchParams.delete('code');
            urlObj.searchParams.delete('state');
            window.history.replaceState({}, document.title, urlObj.toString());
            if (!result.success) {
                uni.showToast({ title: result.msg || '获取失败', icon: 'none', duration: 2000 });
            }
        } else {
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

// 调用后端获取微信用户信息
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

// 复制OpenID
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

// 跳转个人设置
const goToSetting = () => {
    uni.showToast({
        title: '个人设置',
        icon: 'none'
    });
};

// 跳转版本更新
const goToUpdate = () => {
    uni.showToast({
        title: '版本更新',
        icon: 'none'
    });
};

// 返回上一页
const goBack = () => {
    uni.navigateBack();
};
</script>

<style lang="scss" scoped>
.container {
    min-height: 100vh;
    background-color: #f5f7fa;
    padding-bottom: env(safe-area-inset-bottom);
}

.section {
    padding: 24rpx;
}

.section-title {
    font-size: 28rpx;
    font-weight: 600;
    color: #333;
    margin-bottom: 20rpx;
}

.auth-card {
    background: #fff;
    border-radius: 16rpx;
    overflow: hidden;
    margin-bottom: 24rpx;
}

.auth-item {
    display: flex;
    align-items: center;
    padding: 28rpx;
    border-bottom: 1rpx solid #f5f5f5;
    transition: background 0.2s;
    
    &:last-child {
        border-bottom: none;
    }
    
    &:active:not(.disabled) {
        background: #fafafa;
    }
    
    &.disabled {
        opacity: 0.4;
        pointer-events: none;
    }
}

.auth-icon {
    width: 80rpx;
    height: 80rpx;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36rpx;
    margin-right: 20rpx;
    
    &.mini-icon {
        background: linear-gradient(135deg, #e8fdf0, #d9f7e9);
    }
    
    &.mp-icon {
        background: linear-gradient(135deg, #fff8f0, #ffedd5);
    }
}

.auth-info {
    flex: 1;
    
    .auth-title {
        display: block;
        font-size: 28rpx;
        font-weight: 500;
        color: #333;
        margin-bottom: 6rpx;
    }
    
    .auth-desc {
        display: block;
        font-size: 22rpx;
        color: #999;
    }
}

.auth-arrow {
    .loading-text {
        font-size: 24rpx;
        color: #999;
    }
}

.result-card {
    background: #fff;
    border-radius: 16rpx;
    padding: 32rpx;
    text-align: center;
    margin-bottom: 24rpx;
}

.result-header {
    margin-bottom: 20rpx;
}

.success-icon {
    width: 80rpx;
    height: 80rpx;
    background: linear-gradient(135deg, #07c160, #06ad56);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto;
    font-size: 40rpx;
    color: #fff;
    font-weight: bold;
}

.result-title {
    display: block;
    font-size: 30rpx;
    font-weight: 600;
    color: #333;
    margin-bottom: 24rpx;
}

.info-list {
    background: #f8f9fa;
    border-radius: 12rpx;
    padding: 8rpx 0;
    margin-bottom: 24rpx;
}

.info-item {
    display: flex;
    align-items: center;
    padding: 16rpx 20rpx;
    border-bottom: 1rpx solid #e8ecf0;
    
    &:last-child {
        border-bottom: none;
    }
    
    .info-label {
        font-size: 24rpx;
        color: #999;
        min-width: 100rpx;
    }
    
    .info-value {
        flex: 1;
        font-size: 24rpx;
        color: #333;
        text-align: right;
        font-family: monospace;
        word-break: break-all;
    }
}

.action-row {
    display: flex;
    gap: 16rpx;
    
    .action-btn {
        flex: 1;
        height: 72rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 36rpx;
        font-size: 26rpx;
        background: #f5f7fa;
        color: #666;
        
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
