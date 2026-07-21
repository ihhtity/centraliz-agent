<!-- 添加收款账号页面 -->
<template>
	<view class="container">
		<uv-navbar title="添加收款账号" :placeholder="true" @leftClick="goBack" />

		<scroll-view scroll-y class="scroll-container">
			<view class="form-section">
				<view class="section-title">账号类型</view>
				<view class="type-selector">
					<view class="type-item" :class="{ active: form.type === 'personal' }"
						@click="form.type = 'personal'">
						<uv-icon name="account" size="40rpx"
							:color="form.type === 'personal' ? '#10b981' : '#999'"></uv-icon>
						<text>个人账号</text>
					</view>
					<view class="type-item" :class="{ active: form.type === 'company' }" @click="form.type = 'company'">
						<uv-icon name="file-text" size="40rpx"
							:color="form.type === 'company' ? '#10b981' : '#999'"></uv-icon>
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
			<text class="btn-text">提交</text>
		</view>
	</view>
</template>

<script setup>
// 添加收款账号脚本
import { reactive } from 'vue'

// 表单数据
const form = reactive({
	type: 'personal',    // 账号类型：personal-个人，company-企业
	account: '12345678',         // 账号
	name: '12345',            // 姓名
	phone: '13800000000',           // 手机号
	identity: '110101199003070011',        // 身份证（校验码正确）
	card: '4111111111111111',            // 银行卡（Luhn校验通过）
	storename: '12345',       // 店名（企业用）
	encrypt: '12345',         // 营业执照编码（企业用）
	area: '12345',            // 经营地址（企业用）
	remarks: '12345',         // 使用场景描述（企业用）
	share: '0',          // 是否分账：0-关闭，1-开启
	rate: 0              // 分账比例
})

// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/huifu/list'
	});
}

// 身份证号验证
const validateIdentity = (identity) => {
	const len = identity.length
	// 必须是15位或18位
	if (len !== 15 && len !== 18) {
		return false
	}
	// 15位身份证：全部为数字
	if (len === 15) {
		const numRegex = /^\d{15}$/
		return numRegex.test(identity)
	}
	// 18位身份证
	if (len === 18) {
		// 前17位必须是数字
		const prefix = identity.substring(0, 17)
		const numRegex = /^\d{17}$/
		if (!numRegex.test(prefix)) {
			return false
		}
		// 最后一位可以是数字或X
		const lastChar = identity[17].toUpperCase()
		if (!(/[\dX]/.test(lastChar))) {
			return false
		}
		// 校验码验证
		const coefficients = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
		const checkCodes = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2']
		let sum = 0
		for (let i = 0; i < 17; i++) {
			sum += parseInt(prefix[i]) * coefficients[i]
		}
		const remainder = sum % 11
		return checkCodes[remainder] === lastChar
	}
	return false
}

// 银行卡号验证（Luhn算法）
const validateCard = (card) => {
	const cardRegex = /^[0-9]{16,19}$/
	if (!cardRegex.test(card)) {
		return false
	}
	let sum = 0
	let alternate = false
	for (let i = card.length - 1; i >= 0; i--) {
		let digit = parseInt(card[i])
		if (alternate) {
			digit *= 2
			if (digit > 9) {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	return sum % 10 === 0
}

// 提交表单
const submitForm = async () => {
	// 表单验证
	if (!form.account) {
		uni.showToast({ title: '请输入账号', icon: 'none' })
		return
	}
	// 账号验证：6-20位，不能包含中文字符
	const accountRegex = /^[^\u4e00-\u9fa5]{6,20}$/
	if (!accountRegex.test(form.account)) {
		uni.showToast({ title: '账号长度必须在6-20位之间，且不能包含中文字符', icon: 'none' })
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
	// 手机号格式验证
	const phoneRegex = /^1[3-9]\d{9}$/
	if (!phoneRegex.test(form.phone)) {
		uni.showToast({ title: '手机号格式错误', icon: 'none' })
		return
	}
	if (!form.identity) {
		uni.showToast({ title: '请输入身份证号', icon: 'none' })
		return
	}
	// 身份证号格式验证
	if (!validateIdentity(form.identity)) {
		uni.showToast({ title: '身份证号格式错误', icon: 'none' })
		return
	}
	if (!form.card) {
		uni.showToast({ title: '请输入银行卡号', icon: 'none' })
		return
	}
	// 银行卡号格式验证
	if (!validateCard(form.card)) {
		uni.showToast({ title: '银行卡号格式错误', icon: 'none' })
		return
	}
	// 企业账号额外验证
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

	// 提交新增
	try {
		uni.showLoading({ title: '提交中...' })
		let merch = uni.getStorageSync('merch')

		const result = await uni.$uv.http.post('/huifu', {
			type: form.type,
			account: form.account,
			merchs_id: merch.id,
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
		}, {
			custom: { auth: true }
		})
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '添加成功', icon: 'success' })
			setTimeout(() => {
				goBack()
			}, 1500)
		} else {
			uni.showToast({ title: result.msg || '添加失败', icon: 'none' })
		}
	} catch (e) {
		uni.hideLoading()
		uni.showToast({ title: '添加失败', icon: 'none' })
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.scroll-container {
	height: calc(96vh - 88rpx - 140rpx);
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