<template>
	<view class="container">
		<!-- 自定义导航栏 -->
		<uv-navbar :title="t('user.realname.title')" :placeholder="true" @leftClick="goBack"></uv-navbar>
		
		<!-- 状态横幅 -->
		<view class="status-banner">
			<uv-icon name="checkmark-circle" size="40" color="#19be6b"></uv-icon>
			<view class="status-text">
				<text class="main">{{t('user.realname.verified')}}</text>
				<text class="sub">{{t('user.realname.idCard')}}</text>
			</view>
		</view>
		
		<!-- 认证信息表单 -->
		<view class="form-card">
			<uv-form labelPosition="top" :model="form" ref="formRef">
				<uv-form-item :label="t('user.realname.realname')" prop="name">
					<uv-input v-model="form.name" :placeholder="t('user.realname.realnamePlaceholder')" border="surround" shape="round" disabled></uv-input>
				</uv-form-item>
				
				<uv-form-item :label="t('user.realname.idCard')" prop="idCard">
					<uv-input v-model="form.idCard" :placeholder="t('user.realname.idCardPlaceholder')" border="surround" shape="round" disabled></uv-input>
				</uv-form-item>
				
				<uv-form-item :label="t('user.realname.phone')" prop="phone">
					<uv-input v-model="form.phone" :placeholder="t('user.realname.phonePlaceholder')" border="surround" shape="round" disabled></uv-input>
				</uv-form-item>
			</uv-form>
			
			<view class="tips">
				<uv-icon name="info-circle" size="16" color="#909399"></uv-icon>
				<text>{{t('user.realname.tips')}}</text>
			</view>
		</view>
		
		<view class="footer-action" v-if="!isVerified">
			<uv-button type="primary" shape="round" @click="submitAuth">{{t('user.realname.submit')}}</uv-button>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const goBack = () => {
	uni.navigateBack();
};

const isVerified = ref(true); // 模拟已认证状态

const form = reactive({
	name: '张三',
	idCard: '110101199001011234',
	phone: '138****8888'
});

const submitAuth = () => {
	uni.showToast({ title: '提交成功', icon: 'success' });
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.status-banner {
	margin: 110rpx 30rpx 30rpx 30rpx;
	background: #fff;
	border-radius: 16rpx;
	padding: 30rpx;
	display: flex;
	align-items: center;
	gap: 20rpx;
	box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.03);
	
	.status-text {
		display: flex;
		flex-direction: column;
		
		.main {
			font-size: 32rpx;
			font-weight: bold;
			color: #333;
			margin-bottom: 6rpx;
		}
		
		.sub {
			font-size: 24rpx;
			color: #999;
		}
	}
}

.form-card {
	margin: 0 30rpx;
	background: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	box-shadow: 0 4rpx 12rpx rgba(0,0,0,0.03);
	
	.tips {
		display: flex;
		align-items: center;
		gap: 10rpx;
		margin-top: 20rpx;
		padding-top: 20rpx;
		border-top: 1rpx solid #f0f0f0;
		
		text {
			font-size: 22rpx;
			color: #909399;
		}
	}
}

::v-deep .uv-form-item__body__left__content__label {
	font-size: 28rpx;
	color: #333;
	font-weight: 500;
	white-space: nowrap;
}

.footer-action {
	padding: 40rpx 30rpx;
}
</style>