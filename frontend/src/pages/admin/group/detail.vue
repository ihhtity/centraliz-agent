<template>
	<view class="container">
		<uv-navbar title="分组详情" :placeholder="true" @leftClick="goBack" />
		
		<scroll-view scroll-y class="content">
			<!-- 基本信息区域 -->
			<view class="section">
				<view class="section-title">
					<text class="title-text">基本信息</text>
				</view>
				
				<view class="form-item">
					<text class="label">分组名字:</text>
					<view class="input-wrap">
						<uv-input 
							v-model="formData.name" 
							class="form-input" 
							placeholder="请输入分组名称"
						/>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">联系电话:</text>
					<view class="input-wrap">
						<uv-input 
							v-model="formData.phone" 
							class="form-input" 
							type="number"
							placeholder="请输入联系电话"
						/>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">场地位置:</text>
					<view class="input-wrap">
						<uv-input 
							v-model="formData.location" 
							class="form-input" 
							placeholder="请输入场地位置"
						/>
					</view>
				</view>

				<view class="form-item">
					<text class="label">分组类型:</text>
					<view class="radio-group">
						<view 
							v-for="option in typeOptions" 
							:key="option.value"
							class="radio-option"
							:class="{ active: formData.type === option.value }"
							@click="formData.type = option.value"
						>
							<view class="radio-circle">
								<view v-if="formData.type === option.value" class="radio-dot"></view>
							</view>
							<text class="radio-label">{{ option.label }}</text>
						</view>
					</view>
				</view>
				
				<view class="form-item readonly">
					<text class="label">绑锁数量:</text>
					<view class="value-wrap">
						<text class="readonly-value">{{ groupDetail.count || 0 }}</text>
					</view>
				</view>
			</view>
			
			<!-- 设置区域 -->
			<view class="section">
				<view class="section-title">
					<text class="title-text">设置</text>
				</view>
				
				<view class="form-item">
					<text class="label">设备分配</text>
					<view class="radio-group">
						<view 
							v-for="option in deviceOptions" 
							:key="option.value"
							class="radio-option"
							:class="{ active: formData.deviceAssign === option.value }"
							@click="formData.deviceAssign = option.value"
						>
							<view class="radio-circle">
								<view v-if="formData.deviceAssign === option.value" class="radio-dot"></view>
							</view>
							<text class="radio-label">{{ option.label }}</text>
						</view>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">绑定号码</text>
					<view class="radio-group">
						<view 
							v-for="option in bindOptions" 
							:key="option.value"
							class="radio-option"
							:class="{ active: formData.bindNumber === option.value }"
							@click="formData.bindNumber = option.value"
						>
							<view class="radio-circle">
								<view v-if="formData.bindNumber === option.value" class="radio-dot"></view>
							</view>
							<text class="radio-label">{{ option.label }}</text>
						</view>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">消费推送</text>
					<view class="radio-group">
						<view 
							v-for="option in toggleOptions" 
							:key="option.value"
							class="radio-option"
							:class="{ active: formData.consumePush === option.value }"
							@click="formData.consumePush = option.value"
						>
							<view class="radio-circle">
								<view v-if="formData.consumePush === option.value" class="radio-dot"></view>
							</view>
							<text class="radio-label">{{ option.label }}</text>
						</view>
					</view>
				</view>
				
				<view class="form-item">
					<text class="label">控电设置</text>
					<view class="radio-group">
						<view 
							v-for="option in toggleOptions" 
							:key="option.value"
							class="radio-option"
							:class="{ active: formData.powerControl === option.value }"
							@click="formData.powerControl = option.value"
						>
							<view class="radio-circle">
								<view v-if="formData.powerControl === option.value" class="radio-dot"></view>
							</view>
							<text class="radio-label">{{ option.label }}</text>
						</view>
					</view>
				</view>
			</view>
			
			<!-- 功能入口区域 -->
			<view class="section">
				<view class="form-item clickable" @click="goToBillingRules">
					<text class="label">计费规则:</text>
					<view class="arrow-wrap">
						<text class="value-text">{{ billingRule }}</text>
						<uv-icon name="arrow-right" color="#ccc" size="28" />
					</view>
				</view>
				
				<view class="form-item clickable" @click="goToLockManagement">
					<text class="label">门锁管理:</text>
					<view class="arrow-wrap">
						<text class="value-text">门锁管理</text>
						<uv-icon name="arrow-right" color="#ccc" size="28" />
					</view>
				</view>
			</view>
		</scroll-view>
		
		<!-- 底部按钮 -->
		<view class="bottom-bar">
			<uv-button 
				type="primary" 
				shape="circle" 
				@click="submitForm"
			>
				确认提交
			</uv-button>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { onLoad } from '@dcloudio/uni-app';

const groupDetail = ref({});
const groupId = ref('');

// 表单数据
const formData = reactive({
	name: '',
	phone: '',
	location: '',
	type: '0',
	deviceAssign: 'manual',
	bindNumber: 'close',
	consumePush: 'open',
	powerControl: 'open'
});

// 分组类型选项
const typeOptions = [
	{ value: '0', label: '存柜' },
	{ value: '1', label: '零售' }
];

// 选项配置
const deviceOptions = [
	{ value: 'auto', label: '自动' },
	{ value: 'manual', label: '手动' }
];

const bindOptions = [
	{ value: 'auto', label: '自动' },
	{ value: 'manual', label: '手动' },
	{ value: 'close', label: '关闭' }
];

const toggleOptions = [
	{ value: 'close', label: '关闭' },
	{ value: 'open', label: '开启' }
];

const billingRule = ref('预约');

onLoad((options) => {
	groupId.value = options.id || '';
	if (groupId.value) {
		fetchGroupDetail();
	}
});

const goBack = () => {
	uni.navigateBack();
};

const fetchGroupDetail = () => {
	uni.showLoading({ title: '加载中' });
	
	uni.$uv.http.get(`/group/${groupId.value}`, {
		custom: { auth: true }
	}).then((res) => {
		groupDetail.value = res.data;
		formData.name = res.data.name || '';
		formData.phone = res.data.phone || '';
		formData.type = res.data.type || '0';
		formData.location = res.data.location || '';
		uni.hideLoading();
	}).catch((err) => {
		uni.hideLoading();
		console.error('获取分组详情失败:', err);
	});
};

const goToBillingRules = () => {
	uni.showToast({ title: '功能开发中', icon: 'none' });
};

const goToLockManagement = () => {
	uni.showToast({ title: '功能开发中', icon: 'none' });
};

const submitForm = () => {
	if (!formData.name.trim()) {
		uni.showToast({ title: '请输入分组名称', icon: 'none' });
		return;
	}

	const data = {
		rulesId: groupDetail.value.rulesId,
		name: formData.name.trim(),
		phone: formData.phone.trim(),
		type: formData.type.trim(),
		location: formData.location.trim(),
		device_assign: formData.deviceAssign,
		bind_number: formData.bindNumber,
		consume_push: formData.consumePush,
		power_control: formData.powerControl
	};
	// console.log('formData',data);return
	
	uni.showLoading({ title: '保存中' });
	
	uni.$uv.http.put(`/group/${groupId.value}`, data, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		uni.showToast({ title: '保存成功', icon: 'success' });
		setTimeout(() => {
			uni.navigateBack();
		}, 1000);
	}).catch((err) => {
		uni.hideLoading();
		console.error('保存分组失败:', err);
		uni.showToast({ title: '保存失败', icon: 'none' });
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	display: flex;
	flex-direction: column;
}

.content {
	width: 93%;
	margin: 0 auto;
	flex: 1;
	padding: 20rpx;
	padding-bottom: 160rpx;
}

.section {
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);
}

.section-title {
	padding-bottom: 20rpx;
	margin-bottom: 8rpx;
	border-bottom: 1rpx solid #f0f0f0;
	
	.title-text {
		font-size: 30rpx;
		font-weight: 600;
		color: #333;
	}
}

.form-item {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f5f5f5;
	
	&:last-child {
		border-bottom: none;
	}
	
	&.readonly {
		.input-wrap {
			background: #f8f9fa;
		}
	}
	
	&.clickable {
		&:active {
			background: #f8f9fa;
		}
	}
	
	.label {
		font-size: 28rpx;
		color: #666;
		font-weight: 500;
		width: 200rpx;
		flex-shrink: 0;
	}
}

.input-wrap {
	flex: 1;
	text-align: right;
	
	.form-input {
		font-size: 28rpx;
		text-align: right;
	}
}

.value-wrap {
	flex: 1;
	text-align: right;
}

.readonly-value {
	font-size: 28rpx;
	color: #999;
	background: #f8f9fa;
	padding: 12rpx 20rpx;
	border-radius: 8rpx;
	display: inline-block;
}

.radio-group {
	display: flex;
	gap: 40rpx;
}

.radio-option {
	display: flex;
	align-items: center;
	gap: 12rpx;
	
	&.active {
		.radio-circle {
			border-color: #3c9cff;
			background: #fff;
			
			.radio-dot {
				background: #3c9cff;
				transform: scale(1);
			}
		}
		
		.radio-label {
			color: #3c9cff;
			font-weight: 500;
		}
	}
}

.radio-circle {
	width: 36rpx;
	height: 36rpx;
	border-radius: 50%;
	border: 3rpx solid #d9d9d9;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.2s;
}

.radio-dot {
	width: 20rpx;
	height: 20rpx;
	border-radius: 50%;
	background: #d9d9d9;
	transform: scale(0);
	transition: all 0.2s;
}

.radio-label {
	font-size: 28rpx;
	color: #333;
}

.arrow-wrap {
	display: flex;
	align-items: center;
	gap: 12rpx;
}

.value-text {
	font-size: 28rpx;
	color: #999;
}

.bottom-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	padding: 20rpx 30rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	background: #fff;
	box-shadow: 0 -4rpx 12rpx rgba(0, 0, 0, 0.08);
}
</style>