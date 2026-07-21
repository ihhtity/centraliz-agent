<template>
	<view class="container">
		<uv-navbar title="扫码" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack" />

		<view class="content">
			<!-- 微信登录区域 - 仅在微信环境下显示 -->
			<!-- #ifdef MP-WEIXIN -->
			<view class="login-card" @click="handleWXLogin">
				<view class="login-icon">
					<uv-icon name="weixin-fill" size="80" color="#fff" />
				</view>
				<text class="login-title">微信一键登录</text>
				<text class="login-desc">登录后可同步使用小程序和H5端</text>
			</view>
			<!-- #endif -->
			
			<!-- #ifdef H5 -->
			<view class="login-card" @click="handleWXLogin" v-if="isWeChat()">
				<view class="login-icon">
					<uv-icon name="weixin-fill" size="80" color="#fff" />
				</view>
				<text class="login-title">微信一键登录</text>
				<text class="login-desc">登录后可同步使用小程序和H5端</text>
			</view>
			<!-- #endif -->

			<!-- 扫码区域 -->
			<view class="scan-card" @click="handleScan">
				<view class="scan-icon">
					<uv-icon name="scan" size="80" color="#fff" />
				</view>
				<text class="scan-title">扫码</text>
				<text class="scan-desc">点击扫描二维码或条形码</text>
			</view>

			<!-- PC端手动输入 -->
			<view class="input-card">
				<view class="input-header">
					<uv-icon name="edit-pen" size="18" color="#64748b" />
					<text class="input-title">手动输入</text>
				</view>
				<view class="input-body">
					<uv-input
						v-model="manualInput"
						placeholder="请输入二维码内容"
						border="surround"
						:clearable="true"
						@confirm="handleManualSubmit"
					/>
					<view class="input-btn" @click="handleManualSubmit">
						<text class="btn-text">确认</text>
					</view>
				</view>
			</view>

			<!-- 扫描结果 -->
			<view class="result-card" v-if="scanResult">
				<view class="result-header">
					<text class="result-title">扫描结果</text>
					<view class="close-btn" @click="clearResult">
						<uv-icon name="close" size="20" color="#999" />
					</view>
				</view>
				<view class="result-body">
					<view class="result-row">
						<text class="result-label">类型</text>
						<text class="result-value">{{ getResultTypeText(scanResult.scanType) }}</text>
					</view>
					<view class="result-row">
						<text class="result-label">内容</text>
						<text class="result-value content">{{ scanResult.result }}</text>
					</view>
					<view class="result-row" v-if="scanResult.charSet">
						<text class="result-label">字符集</text>
						<text class="result-value">{{ scanResult.charSet }}</text>
					</view>
					<view class="result-row" v-if="scanResult.path">
						<text class="result-label">路径</text>
						<text class="result-value">{{ scanResult.path }}</text>
					</view>
				</view>
			</view>

			<!-- 使用说明 -->
			<view class="guide-card">
				<view class="guide-header">
					<uv-icon name="info-circle" size="18" color="#64748b" />
					<text class="guide-title">使用说明</text>
				</view>
				<view class="guide-steps">
					<view class="guide-step">
						<view class="step-num">1</view>
						<text class="step-text">移动端：点击扫码区域打开相机扫描</text>
					</view>
					<view class="guide-step">
						<view class="step-num">2</view>
						<text class="step-text">PC端：手动输入二维码内容</text>
					</view>
					<view class="guide-step">
						<view class="step-num">3</view>
						<text class="step-text">系统自动识别并显示结果</text>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app';

onLoad(() => {
    user.value = uni.getStorageSync('user')
    getPlatform()
    // console.log('页面加载',user.value)
})

// 用户数据
const user = ref({})
// 扫描结果
const scanResult = ref(null)
// 手动输入内容
const manualInput = ref('')

// 判断当前运行平台
const getPlatform = () => {
	const SystemInfo = uni.getStorageSync('SystemInfo')
	console.log(SystemInfo.platform, SystemInfo.uniPlatform)
}

// 微信登录
const handleWXLogin = () => {
	uni.showLoading({ title: '登录中...' })
	
	// #ifdef MP-WEIXIN
	// 微信小程序登录
	uni.login({
		provider: 'weixin',
		success: (loginRes) => {
			// console.log('微信小程序登录成功:', loginRes)
			wxLoginRequest(loginRes.code, 'miniprogram')
		},
		fail: (err) => {
			console.error('微信登录失败:', err)
			uni.showToast({ title: '登录失败', icon: 'none' })
		}
	})
	// #endif
	
	// #ifdef H5
	// H5端微信公众号登录
	// 构建微信OAuth2授权URL
	const appId = 'wxfc42adf0c2f58bb6' // 公众号AppID
	// 网页授权域名（需在微信公众平台配置）
	const redirectUri = encodeURIComponent('https://centraliz.bsldtech.cn/pages/user/index/callback')
	const scope = 'snsapi_userinfo'
	const state = 'wx_login'
	
	const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${appId}&redirect_uri=${redirectUri}&response_type=code&scope=${scope}&state=${state}#wechat_redirect`
	
	// 在H5环境下跳转到微信授权页面
	window.location.href = authUrl
	// #endif
}

// 发起微信登录请求
const wxLoginRequest = async (code, platform) => {
    try {
        
        const requestData = {
            code: code,
            platform: platform
        }
        
        const res = await uni.$uv.http.post('/wechat/login', requestData)
        
        uni.hideLoading()
			
        if (res.code === 200) {
            // 登录成功，保存用户信息和token
            const userData = res.data.user
            const token = res.data.token
            
            // 保存到本地存储
            uni.setStorageSync('user', userData)
            uni.setStorageSync('token', token)
            
            uni.showToast({ title: '登录成功', icon: 'success', duration: 2000 })
            
            // 根据之前保存的路由跳转
            let userroute = uni.getStorageSync('userroute')
            if (userroute === 'locker') {
                uni.redirectTo({ url: '/pages/user/index/locker' })
            } else if (userroute === 'retail') {
                uni.redirectTo({ url: '/pages/user/index/retail' })
            } else {
                uni.redirectTo({ url: '/pages/user/index/locker' })
            }
        } else {
            uni.showToast({ title: res.msg || '登录失败', icon: 'none' })
        }
    } catch (e) {
        console.error('登录请求失败:', e)
        uni.showToast({ title: '网络错误，请重试', icon: 'none' })
    }
}

// 解析URL参数（适配小程序和H5）
const parseUrlParams = (url) => {
    const params = {
        type: '',
        groupType: '',
        id: 0
    }
    
    try {
        // 直接查找 ? 位置
        const questionIndex = url.indexOf('?')
        
        if (questionIndex === -1) {
            console.log('URL无参数')
            return params
        }
        
        // 提取 ? 后面的内容
        const queryString = url.substring(questionIndex + 1)
        console.log('查询字符串:', queryString)
        
        // 解析查询参数
        // #ifdef H5
        // H5环境：使用 URLSearchParams
        const urlParams = new URLSearchParams(queryString)
        params.type = urlParams.get('type') || ''
        params.groupType = urlParams.get('groupType') ? decodeURIComponent(urlParams.get('groupType')) : ''
        params.id = parseInt(urlParams.get('id'), 10) || 0
        // #endif
        
        // #ifndef H5
        // 小程序环境：手动解析参数
        const pairs = queryString.split('&')
        for (let i = 0; i < pairs.length; i++) {
            const pair = pairs[i].split('=')
            const key = decodeURIComponent(pair[0])
            const value = pair.length > 1 ? decodeURIComponent(pair[1]) : ''
            
            if (key === 'type') {
                params.type = value
            } else if (key === 'groupType') {
                params.groupType = value
            } else if (key === 'id') {
                params.id = parseInt(value, 10) || 0
            }
        }
        // #endif
        
        console.log('解析参数:', params)
    } catch (e) {
        console.error('URL参数解析失败:', e)
    }
    return params
}

// 根据参数跳转到对应页面
const navigateByParams = (params) => {
    params.id = params.id || 0
    params.type = params.type || 'group'
    params.groupType = params.groupType || '存柜'
    
    console.log('解析参数:', params)
    
    let url = '/pages/user/index/'
    
    if (params.groupType === '存柜') {
        uni.setStorageSync('userroute', "locker")
        if (params.type === 'group') {
            url += 'locker'
        } else if (params.type === 'room') {
            url += 'locker'
        }
    } else if (params.groupType === '零售') {
        uni.setStorageSync('userroute', "retail")
        if (params.type === 'group') {
            url += 'retail'
        } else if (params.type === 'room') {
            url += 'retail'
        }
    }
    
    // 将参数拼接到 URL 中（uni-app 的 redirectTo 不支持 query 参数）
    const queryStr = Object.keys(params)
        .map(key => `${key}=${encodeURIComponent(params[key])}`)
        .join('&')
    
    if (queryStr) {
        url += '?' + queryStr
    }
    
    console.log('跳转URL:', url)
    uni.redirectTo({ url })
}

// 扫码相关状态
const isScanning = ref(false)
const scanContainer = ref(null)

// 判断是否为移动端
const isMobile = () => {
    // #ifdef H5
    const userAgent = navigator.userAgent.toLowerCase()
    return /mobile|android|iphone|ipad|ipod|blackberry|windows phone/.test(userAgent)
    // #endif
    // #ifndef H5
    return true
    // #endif
}

// 判断是否为微信浏览器
const isWeChat = () => {
    // #ifdef H5
    const userAgent = navigator.userAgent.toLowerCase()
    return /micromessenger/.test(userAgent)
    // #endif
    // #ifndef H5
    return false
    // #endif
}

// H5端扫码实现
const h5Scan = () => {
    // #ifdef H5
    // 检查是否为移动端
    if (!isMobile()) {
        uni.showToast({ title: '请在移动端使用扫码功能', icon: 'none' })
        isScanning.value = false
        return
    }
    
    // 检查是否为微信浏览器
    if (isWeChat()) {
        // 使用微信JS-SDK进行扫码
        wechatScan()
        return
    }
    
    // 非微信浏览器，使用原生getUserMedia扫码
    nativeScan()
    // #endif
}

// 微信浏览器扫码（使用微信JS-SDK）
const wechatScan = () => {
    // #ifdef H5
    console.log('使用微信JS-SDK扫码')
    
    // 检查window.wx是否可用
    if (typeof window.wx === 'undefined') {
        // 初始化微信JS-SDK
        initWXConfig(() => {
            doWechatScan()
        })
    } else {
        doWechatScan()
    }
    // #endif
}

// 初始化微信JS-SDK配置
const initWXConfig = (callback) => {
    // #ifdef H5
    // 从后端获取签名配置
    uni.request({
        url: '/api/wechat/jssdk/config',
        method: 'GET',
        data: {
            url: window.location.href.split('#')[0]
        },
        success: (res) => {
            if (res.data.code === 200) {
                const config = res.data.data
                window.wx.config({
                    debug: false,
                    appId: config.appId,
                    timestamp: config.timestamp,
                    nonceStr: config.nonceStr,
                    signature: config.signature,
                    jsApiList: ['scanQRCode']
                })
                
                window.wx.ready(() => {
                    console.log('微信JS-SDK初始化成功')
                    callback()
                })
                
                window.wx.error((err) => {
                    console.error('微信JS-SDK初始化失败:', err)
                    uni.showToast({ title: '微信配置失败，请重试', icon: 'none' })
                    isScanning.value = false
                })
            } else {
                console.error('获取微信配置失败:', res.data.msg)
                uni.showToast({ title: '获取微信配置失败', icon: 'none' })
                isScanning.value = false
            }
        },
        fail: (err) => {
            console.error('请求微信配置失败:', err)
            // 降级为图片选择方式
            chooseImageScan()
        }
    })
    // #endif
}

// 执行微信扫码
const doWechatScan = () => {
    // #ifdef H5
    window.wx.scanQRCode({
        needResult: 1, // 1表示直接返回扫描结果
        scanType: ['qrCode', 'barCode'],
        success: (res) => {
            console.log('微信扫码成功:', res)
            const result = {
                result: res.resultStr,
                scanType: 'qrCode'
            }
            handleScanSuccess(result)
        },
        fail: (err) => {
            console.error('微信扫码失败:', err)
            // 如果用户取消或失败，提供图片选择方式
            if (err.errMsg && !err.errMsg.includes('cancel')) {
                uni.showModal({
                    title: '扫码失败',
                    content: '是否选择图片进行识别？',
                    success: (modalRes) => {
                        if (modalRes.confirm) {
                            chooseImageScan()
                        } else {
                            isScanning.value = false
                        }
                    }
                })
            } else {
                isScanning.value = false
            }
        }
    })
    // #endif
}

// 选择图片扫码（备选方案）
const chooseImageScan = () => {
    // #ifdef H5
    uni.chooseImage({
        count: 1,
        sizeType: ['compressed'],
        sourceType: ['album'],
        success: (res) => {
            console.log('选择图片成功:', res)
            // 将图片上传到后端进行识别
            uploadAndScan(res.tempFilePaths[0])
        },
        fail: (err) => {
            console.error('选择图片失败:', err)
            uni.showToast({ title: '请使用手动输入', icon: 'none' })
            isScanning.value = false
        }
    })
    // #endif
}

// 上传图片并识别
const uploadAndScan = (imagePath) => {
    // #ifdef H5
    uni.showLoading({ title: '识别中...' })
    
    uni.uploadFile({
        url: '/api/wechat/scan/image',
        filePath: imagePath,
        name: 'image',
        success: (res) => {
            uni.hideLoading()
            try {
                const data = JSON.parse(res.data)
                if (data.code === 200) {
                    const result = {
                        result: data.data.result,
                        scanType: 'qrCode'
                    }
                    handleScanSuccess(result)
                } else {
                    uni.showToast({ title: data.msg || '识别失败', icon: 'none' })
                    isScanning.value = false
                }
            } catch (e) {
                console.error('解析识别结果失败:', e)
                uni.showToast({ title: '识别失败', icon: 'none' })
                isScanning.value = false
            }
        },
        fail: (err) => {
            uni.hideLoading()
            console.error('上传图片失败:', err)
            uni.showToast({ title: '上传失败，请使用手动输入', icon: 'none' })
            isScanning.value = false
        }
    })
    // #endif
}

// 原生浏览器扫码（非微信浏览器）
const nativeScan = () => {
    // #ifdef H5
    console.log('使用原生扫码')
    
    // 检查摄像头权限
    navigator.permissions.query({ name: 'camera' }).then((permissionStatus) => {
        if (permissionStatus.state === 'denied') {
            uni.showModal({
                title: '权限不足',
                content: '请在系统设置中开启摄像头权限',
                showCancel: false,
                success: () => {
                    isScanning.value = false
                }
            })
            return
        }
        
        // 调用uni.scanCode，H5端会调用系统扫码能力
        uni.scanCode({
            scanType: ['qrCode', 'barCode'],
            success: (res) => {
                handleScanSuccess(res)
            },
            fail: (err) => {
                console.error('H5扫码失败:', err)
                handleScanFail(err)
            }
        })
    }).catch((err) => {
        console.error('检查权限失败:', err)
        // 直接尝试扫码
        uni.scanCode({
            scanType: ['qrCode', 'barCode'],
            success: (res) => {
                handleScanSuccess(res)
            },
            fail: (err) => {
                console.error('H5扫码失败:', err)
                handleScanFail(err)
            }
        })
    })
    // #endif
}

// 小程序端扫码
const miniProgramScan = () => {
    // #ifdef MP-WEIXIN
    uni.scanCode({
        scanType: ['qrCode', 'barCode'],
        success: (res) => {
            handleScanSuccess(res)
        },
        fail: (err) => {
            console.error('小程序扫码失败:', err)
            handleScanFail(err)
        }
    })
    // #endif
}

// 处理扫码成功
const handleScanSuccess = (res) => {
    console.log('扫码结果:', res)
    isScanning.value = false
    scanResult.value = res
    
    // 解析URL参数
    const params = parseUrlParams(res.result)
    
    // 如果是有效的二维码内容，跳转到对应页面
    if (params.type || params.groupType || params.id) {
        navigateByParams(params)
    }
    
    uni.showToast({ title: '扫描成功', icon: 'success' })
}

// 处理扫码失败
const handleScanFail = (err) => {
    isScanning.value = false
    console.error('扫码失败:', err)
    
    const errMsg = err.errMsg || err.message || '扫描失败'
    
    // 处理用户取消的情况
    if (errMsg.includes('cancel') || errMsg.includes('取消')) {
        uni.showToast({ title: '已取消', icon: 'none' })
    } else if (errMsg.includes('permission') || errMsg.includes('权限')) {
        uni.showModal({
            title: '权限不足',
            content: '请在系统设置中开启摄像头权限',
            showCancel: false
        })
    } else {
        uni.showToast({ title: '请使用手动输入', icon: 'none' })
    }
}

// 处理扫码入口
const handleScan = () => {
    if (isScanning.value) {
        return
    }
    
    isScanning.value = true
    uni.showLoading({ title: '准备扫码...' })
    
    setTimeout(() => {
        uni.hideLoading()
        
        // #ifdef MP-WEIXIN
        miniProgramScan()
        // #endif
        
        // #ifdef H5
        h5Scan()
        // #endif
        
        // #ifndef H5
        // 其他非H5非小程序环境（如App）
        uni.scanCode({
            scanType: ['qrCode', 'barCode'],
            success: (res) => {
                handleScanSuccess(res)
            },
            fail: (err) => {
                console.error('扫码失败:', err)
                handleScanFail(err)
            }
        })
        // #endif
    }, 500)
}

// 处理手动输入提交
const handleManualSubmit = () => {
    if (!manualInput.value.trim()) {
        uni.showToast({ title: '请输入内容', icon: 'none' })
        return
    }
    
    const result = {
        result: manualInput.value.trim(),
        scanType: 'manual',
        charSet: 'UTF-8'
    }
    
    console.log('手动输入结果:', result)
    scanResult.value = result
    uni.showToast({ title: '已提交', icon: 'success' })
    
    // 清空输入框
    manualInput.value = ''
    
    // 解析URL参数并跳转
    const params = parseUrlParams(result.result)
    
    if (params.type || params.groupType || params.id) {
        navigateByParams(params)
    } else {
        // 如果不是有效URL，显示提示
        uni.showToast({ title: '非有效二维码', icon: 'none' })
    }
}

// 获取结果类型文本
const getResultTypeText = (type) => {
    if (type === 'qrCode') return '二维码'
    if (type === 'barCode') return '条形码'
    if (type === 'manual') return '手动输入'
    return type || '未知'
}

// 清除结果
const clearResult = () => {
    scanResult.value = null
}

// 新增: 返回首页
const goBack = () => {
    let userroute = uni.getStorageSync('userroute')
    console.log(userroute)
	if (userroute === 'locker') {
		uni.redirectTo({ url: '/pages/user/index/locker' })
	} else if (userroute === 'retail') {
		uni.redirectTo({ url: '/pages/user/index/retail' })
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background: linear-gradient(180deg, #f8fafc 0%, #e2e8f0 100%);
	padding-bottom: 120rpx;
}

.content {
	padding: 40rpx 32rpx;
}

/* 微信登录卡片 */
.login-card {
	background: linear-gradient(135deg, #07c160 0%, #10b981 100%);
	border-radius: 24rpx;
	padding: 60rpx 40rpx;
	text-align: center;
	margin-bottom: 24rpx;
	transition: transform 0.2s ease;
	
	&:active {
		transform: scale(0.98);
	}
}

.login-icon {
	width: 140rpx;
	height: 140rpx;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 50%;
	display: flex;
	justify-content: center;
	align-items: center;
	margin: 0 auto 28rpx;
}

.login-title {
	font-size: 36rpx;
	font-weight: 600;
	color: #fff;
	display: block;
	margin-bottom: 12rpx;
}

.login-desc {
	font-size: 26rpx;
	color: rgba(255, 255, 255, 0.8);
	display: block;
}

/* 扫码卡片 */
.scan-card {
	background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
	border-radius: 24rpx;
	padding: 60rpx 40rpx;
	text-align: center;
	margin-bottom: 24rpx;
	transition: transform 0.2s ease;
	
	&:active {
		transform: scale(0.98);
	}
}

.scan-icon {
	width: 140rpx;
	height: 140rpx;
	background: rgba(255, 255, 255, 0.2);
	border-radius: 50%;
	display: flex;
	justify-content: center;
	align-items: center;
	margin: 0 auto 28rpx;
}

.scan-title {
	font-size: 36rpx;
	font-weight: 600;
	color: #fff;
	display: block;
	margin-bottom: 12rpx;
}

.scan-desc {
	font-size: 26rpx;
	color: rgba(255, 255, 255, 0.8);
	display: block;
}

/* 手动输入卡片 */
.input-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

.input-header {
	display: flex;
	align-items: center;
	margin-bottom: 24rpx;
}

.input-title {
	font-size: 28rpx;
	font-weight: 600;
	color: #1e293b;
	margin-left: 12rpx;
}

.input-body {
	display: flex;
	gap: 16rpx;
	align-items: center;
}

.input-btn {
	background: linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%);
	padding: 20rpx 40rpx;
	border-radius: 12rpx;
	flex-shrink: 0;
	transition: transform 0.2s ease;
	
	&:active {
		transform: scale(0.95);
	}
	
	.btn-text {
		font-size: 28rpx;
		font-weight: 500;
		color: #fff;
	}
}

/* 结果卡片 */
.result-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

.result-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 24rpx;
	padding-bottom: 20rpx;
	border-bottom: 1rpx solid #f1f5f9;
}

.result-title {
	font-size: 30rpx;
	font-weight: 600;
	color: #1e293b;
}

.close-btn {
	width: 48rpx;
	height: 48rpx;
	display: flex;
	justify-content: center;
	align-items: center;
	border-radius: 50%;
	background: #f8fafc;
	
	&:active {
		background: #e2e8f0;
	}
}

.result-body {
	.result-row {
		display: flex;
		flex-direction: column;
		margin-bottom: 20rpx;
		
		&:last-child {
			margin-bottom: 0;
		}
	}
	
	.result-label {
		font-size: 24rpx;
		color: #94a3b8;
		margin-bottom: 8rpx;
	}
	
	.result-value {
		font-size: 28rpx;
		color: #1e293b;
		font-weight: 500;
		
		&.content {
			background: #f8fafc;
			padding: 20rpx;
			border-radius: 12rpx;
			word-break: break-all;
			line-height: 1.6;
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', monospace;
		}
	}
}

/* 说明卡片 */
.guide-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 32rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

.guide-header {
	display: flex;
	align-items: center;
	margin-bottom: 28rpx;
}

.guide-title {
	font-size: 28rpx;
	font-weight: 600;
	color: #1e293b;
	margin-left: 12rpx;
}

.guide-steps {
	.guide-step {
		display: flex;
		align-items: flex-start;
		margin-bottom: 20rpx;
		
		&:last-child {
			margin-bottom: 0;
		}
	}
	
	.step-num {
		width: 40rpx;
		height: 40rpx;
		background: linear-gradient(135deg, #3b82f6 0%, #60a5fa 100%);
		border-radius: 50%;
		font-size: 24rpx;
		color: #fff;
		display: flex;
		justify-content: center;
		align-items: center;
		margin-right: 16rpx;
		flex-shrink: 0;
	}
	
	.step-text {
		font-size: 26rpx;
		color: #64748b;
		line-height: 1.5;
		padding-top: 6rpx;
	}
}
</style>
