<!-- 编辑收款账号页面 -->
<template>
	<view class="container">
		<uv-navbar title="编辑收款账号" :placeholder="true" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container">
			<view class="form-section">
				<view class="section-title">账号类型</view>
				<view class="type-selector">
					<view class="type-item" :class="{ active: form.type === 'personal' }" @click="form.type = 'personal'">
						<uv-icon name="account" size="40rpx" :color="form.type === 'personal' ? '#10b981' : '#999'"></uv-icon>
						<text>个人账号</text>
					</view>
					<view class="type-item" :class="{ active: form.type === 'company' }" @click="form.type = 'company'">
						<uv-icon name="file-text" size="40rpx" :color="form.type === 'company' ? '#10b981' : '#999'"></uv-icon>
						<text>企业账号</text>
					</view>
				</view>
			</view>

			<view class="form-section">
				<view class="section-title">基本信息</view>
				<view class="form-item">
					<text class="form-label">账号</text>
					<input class="form-input" v-model="form.account" placeholder="请输入账号" />
				</view>
				<view class="form-item">
					<text class="form-label">姓名</text>
					<input class="form-input" v-model="form.name" placeholder="请输入姓名" />
				</view>
				<view class="form-item">
					<text class="form-label">手机号</text>
					<input class="form-input" v-model="form.phone" placeholder="请输入手机号" />
				</view>
				<view class="form-item">
					<text class="form-label">身份证</text>
					<input class="form-input" v-model="form.identity" placeholder="请输入身份证号" />
				</view>
				<view class="form-item">
					<text class="form-label">银行卡</text>
					<input class="form-input" v-model="form.card" placeholder="请输入银行卡号" />
				</view>
				<view class="form-item">
					<text class="form-label">汇付编码</text>
					<input class="form-input" :value="form.code || '系统生成'" disabled />
				</view>
				<view class="form-item">
					<text class="form-label">多方分账</text>
					<input class="form-input" :value="form.sharing || '系统生成'" disabled />
				</view>
			</view>

			<view class="form-section">
				<view class="section-title">分账设置</view>
				<view class="switch-item">
					<text class="switch-label">开启分账</text>
					<switch :checked="form.share === '1'" @change="onShareChange" color="#10b981" />
				</view>
				<view v-if="form.share === '1'" class="form-item">
					<text class="form-label">分账比例 (%)</text>
					<input class="form-input" type="digit" v-model="form.rate" placeholder="请输入分账比例" />
				</view>
			</view>

			<view v-if="form.type === 'company'" class="form-section">
				<view class="section-title">企业信息</view>
				<view class="form-item">
					<text class="form-label">店名</text>
					<input class="form-input" v-model="form.storename" placeholder="请输入店名" />
				</view>
				<view class="form-item">
					<text class="form-label">营业执照编码</text>
					<input class="form-input" v-model="form.encrypt" placeholder="请输入营业执照编码" />
				</view>
				<view class="form-item">
					<text class="form-label">经营地址</text>
					<textarea class="form-textarea" v-model="form.area" placeholder="请输入经营地址" />
				</view>
				<view class="form-item">
					<text class="form-label">使用场景描述</text>
					<textarea class="form-textarea" v-model="form.remarks" placeholder="请输入使用场景描述" />
				</view>
			</view>
		</scroll-view>

		<view class="bottom-btn" @click="submitForm">
			<text class="btn-text">保存</text>
		</view>
	</view>
</template>

<script setup>
import { reactive } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

// 收款账号表单数据
const form = reactive({
	type: 'personal',
	account: '',
	name: '',
	phone: '',
	identity: '',
	card: '',
	code: '',
	sharing: '',
	storename: '',
	encrypt: '',
	area: '',
	remarks: '',
	share: '0',
	rate: 0,
	id: ''
})

// 页面加载时获取收款账号详情
onLoad(async (options) => {
	form.id = options.id || ''
	if (form.id) {
		await loadAccountDetail()
	}
})

// 获取收款账号详情
const loadAccountDetail = async () => {
	uni.showLoading({ title: '加载中...' })
	try {
		const res = await uni.$uv.http.get('/huifu/detail', {
			params: {
				id: form.id,
			},
			custom: { auth: true }
		})
		if (res.code === 200 && res.data) {
			Object.assign(form, {
				type: res.data.type || 'personal',
				account: res.data.account || '',
				name: res.data.name || '',
				phone: res.data.phone || '',
				identity: res.data.identity || '',
				card: res.data.card || '',
				code: res.data.code || '',
				sharing: res.data.sharing || '',
				storename: res.data.storename || '',
				encrypt: res.data.encrypt || '',
				area: res.data.area || '',
				remarks: res.data.remarks || '',
				share: res.data.share || '0',
				rate: res.data.rate || 0,
				id: res.data.id || form.id
			})
		}
	} catch (e) {
		console.error('加载失败', e)
	} finally {
		uni.hideLoading()
	}
}

// 提交表单
const submitForm = async () => {
	if (!form.account) {
		uni.showToast({ title: '请输入账号', icon: 'none' })
		return
	}
	if (!form.name) {
		uni.showToast({ title: '请输入姓名', icon: 'none' })
		return
	}
	if (!form.phone) {
		uni.showToast({ title: '请输入手机号', icon: 'none' })
		return
	}
	if (!form.identity) {
		uni.showToast({ title: '请输入身份证号', icon: 'none' })
		return
	}
	if (!form.card) {
		uni.showToast({ title: '请输入银行卡号', icon: 'none' })
		return
	}
	if (form.type === 'company') {
		if (!form.storename) {
			uni.showToast({ title: '请输入店名', icon: 'none' })
			return
		}
		if (!form.encrypt) {
			uni.showToast({ title: '请输入营业执照编码', icon: 'none' })
			return
		}
	}
	if (form.share === '1' && (!form.rate || form.rate <= 0)) {
		uni.showToast({ title: '请输入有效的分账比例', icon: 'none' })
		return
	}

	try {
		uni.showLoading({ title: '保存中...' })
		const result = await uni.$uv.http.put('/huifu', {params: {
			id: form.id,
			type: form.type,
			account: form.account,
			name: form.name,
			phone: form.phone,
			identity: form.identity,
			card: form.card,
			storename: form.storename,
			encrypt: form.encrypt,
			area: form.area,
			remarks: form.remarks,
			share: form.share,
			rate: parseFloat(form.rate) || 0
		}, custom: {
			auth: true
		}})
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '保存成功', icon: 'success' })
			setTimeout(() => {
				uni.navigateBack()
			}, 1500)
		} else {
			uni.showToast({ title: result.msg || '保存失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '保存失败', icon: 'none' })
	}
}

// 分账开关改变时更新表单数据
const onShareChange = (e) => {
	form.share = e.detail.value ? '1' : '0'
}

// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/huifu/list'
	});
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.scroll-container {
	height: calc(96vh - 88rpx - 100rpx);
	width: calc(100% - 48rpx);
	padding: 24rpx;
}

.form-section {
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin-bottom: 16rpx;
}

.section-title {
	font-size: 28rpx;
	color: #333;
	font-weight: 600;
	margin-bottom: 24rpx;
	padding-left: 12rpx;
	border-left: 6rpx solid #10b981;
}

.type-selector {
	display: flex;
	gap: 24rpx;
}

.type-item {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 24rpx;
	border-radius: 12rpx;
	background: #f5f5f5;
	transition: all 0.2s;
	&.active {
		background: rgba(16, 185, 129, 0.1);
		border: 2rpx solid #10b981;
	}
	text {
		font-size: 26rpx;
		color: #666;
		margin-top: 12rpx;
	}
}

.form-item {
	margin-bottom: 24rpx;
	&:last-child {
		margin-bottom: 0;
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
	background-color: #f5f5f5;
	&[disabled] {
		color: #999;
	}
}

.form-textarea {
	width: 100%;
	height: 160rpx;
	padding: 16rpx 20rpx;
	border: 1rpx solid #e8e8e8;
	border-radius: 12rpx;
	font-size: 28rpx;
	box-sizing: border-box;
}

.switch-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 16rpx 0;
}

.switch-label {
	font-size: 28rpx;
	color: #333;
}

.bottom-btn {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	height: 100rpx;
	padding-bottom: constant(safe-area-inset-bottom);
	padding-bottom: env(safe-area-inset-bottom);
	display: flex;
	align-items: center;
	justify-content: center;
	background: linear-gradient(135deg, #10b981, #059669);
	box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.1);
}

.btn-text {
	font-size: 30rpx;
	color: #fff;
	font-weight: 500;
}
</style>