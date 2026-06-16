<template>
	<view class="container">
		<view class="loading-content">
			<view class="loading-icon">
				<uv-icon name="loading" size="60" color="#3b82f6" />
			</view>
			<text class="loading-text">登录中...</text>
		</view>
	</view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app';

onLoad((options) => {
	console.log('回调页面参数:', options)
	
	// #ifdef H5
	// 尝试从URL参数获取token（后端重定向携带）
	const token = options?.token || getUrlParam('token')
	const userId = options?.user_id || getUrlParam('user_id')
	
	if (token && userId) {
		// 直接使用后端返回的token
		handleTokenLogin(token, userId)
	} else {
		// 尝试从URL参数获取code（前端直接处理）
		const code = options?.code || getUrlParam('code')
		const state = options?.state || getUrlParam('state')
		
		if (code && state === 'wx_login') {
			wxLoginRequest(code, 'mp')
		} else {
			uni.showToast({ title: '登录失败', icon: 'none' })
			setTimeout(() => {
				uni.redirectTo({ url: '/pages/user/index/index' })
			}, 1500)
		}
	}
	// #endif
	
	// #ifdef MP-WEIXIN
	// 小程序端直接跳转首页
	uni.redirectTo({ url: '/pages/user/index/index' })
	// #endif
})

// 获取URL参数
const getUrlParam = (name) => {
	const url = window.location.href
	const regex = new RegExp(`[?&]${name}=([^&#]*)`)
	const results = regex.exec(url)
	return results ? decodeURIComponent(results[1]) : ''
}

// 直接使用token登录
const handleTokenLogin = (token, userId) => {
	console.log('使用token登录:', token, userId)
	
	// 保存token到本地存储
	uni.setStorageSync('token', token)
	uni.setStorageSync('userId', userId)
	
	// 获取用户信息
	uni.request({
		url: '/api/user/profile',
		method: 'GET',
		header: {
			'Authorization': 'Bearer ' + token
		},
		success: (res) => {
			if (res.data.code === 200) {
				uni.setStorageSync('user', res.data.data)
				uni.showToast({ title: '登录成功', icon: 'success' })
				
				let userroute = uni.getStorageSync('userroute')
				setTimeout(() => {
					if (userroute === 'locker') {
						uni.redirectTo({ url: '/pages/user/index/locker' })
					} else if (userroute === 'retail') {
						uni.redirectTo({ url: '/pages/user/index/retail' })
					} else {
						uni.redirectTo({ url: '/pages/user/index/index' })
					}
				}, 1500)
			} else {
				uni.showToast({ title: '获取用户信息失败', icon: 'none' })
				setTimeout(() => {
					uni.redirectTo({ url: '/pages/user/index/index' })
				}, 1500)
			}
		},
		fail: (err) => {
			console.error('获取用户信息失败:', err)
			uni.showToast({ title: '网络错误', icon: 'none' })
			setTimeout(() => {
				uni.redirectTo({ url: '/pages/user/index/index' })
			}, 1500)
		}
	})
}

// 发起微信登录请求
const wxLoginRequest = (code, platform) => {
	uni.request({
		url: '/api/wechat/login',
		method: 'POST',
		data: {
			code: code,
			platform: platform
		},
		success: (res) => {
			console.log('登录接口返回:', res)
			
			if (res.data.code === 200) {
				const userData = res.data.data.user
				const token = res.data.data.token
				
				uni.setStorageSync('user', userData)
				uni.setStorageSync('token', token)
				
				uni.showToast({ title: '登录成功', icon: 'success' })
				
				let userroute = uni.getStorageSync('userroute')
				setTimeout(() => {
					if (userroute === 'locker') {
						uni.redirectTo({ url: '/pages/user/index/locker' })
					} else if (userroute === 'retail') {
						uni.redirectTo({ url: '/pages/user/index/retail' })
					} else {
						uni.redirectTo({ url: '/pages/user/index/index' })
					}
				}, 1500)
			} else {
				uni.showToast({ title: res.data.msg || '登录失败', icon: 'none' })
				setTimeout(() => {
					uni.redirectTo({ url: '/pages/user/index/index' })
				}, 1500)
			}
		},
		fail: (err) => {
			console.error('登录请求失败:', err)
			uni.showToast({ title: '网络错误', icon: 'none' })
			setTimeout(() => {
				uni.redirectTo({ url: '/pages/user/index/index' })
			}, 1500)
		}
	})
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background: linear-gradient(180deg, #f8fafc 0%, #e2e8f0 100%);
	display: flex;
	justify-content: center;
	align-items: center;
}

.loading-content {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.loading-icon {
	width: 120rpx;
	height: 120rpx;
	background: rgba(59, 130, 246, 0.1);
	border-radius: 50%;
	display: flex;
	justify-content: center;
	align-items: center;
	margin-bottom: 32rpx;
}

.loading-text {
	font-size: 30rpx;
	color: #64748b;
}
</style>
