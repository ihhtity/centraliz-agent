<template>
	<view class="container">
		<!-- 修改: 添加返回按钮 -->
		<uv-navbar :title="isEdit ? t('admin.device.editDevice') : t('admin.device.addDevice')" :placeholder="true" @leftClick="goBack" />
		
		<uv-form labelPosition="left" :model="form" ref="formRef" style="padding: 20rpx;">
			<uv-form-item :label="t('admin.device.deviceName')" prop="name">
				<uv-input v-model="form.name" :placeholder="t('admin.device.deviceName')" border="surround"></uv-input>
			</uv-form-item>
			<uv-form-item :label="t('admin.device.deviceCode')" prop="code">
				<uv-input v-model="form.code" :placeholder="t('admin.device.deviceCode')" border="surround"></uv-input>
			</uv-form-item>
			<uv-form-item :label="t('admin.group.belongGroup')" prop="groupId">
				<uv-picker :range="groupList" rangeKey="name" @confirm="onGroupConfirm"></uv-picker>
				<view class="picker-value">{{ t('admin.device.selectGroup') }}</view>
			</uv-form-item>
			<uv-form-item :label="t('admin.device.status')" prop="status">
				<uv-switch v-model="form.enabled" activeColor="#19be6b"></uv-switch>
			</uv-form-item>
		</uv-form>
		
		<view style="padding: 40rpx;">
			<uv-button type="primary" @click="submit">{{ t('common.save') }}</uv-button>
		</view>
	</view>
</template>

<script setup>
import { ref, reactive, computed } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// 新增: 返回上一页方法
const goBack = () => {
	uni.navigateBack();
};

const isEdit = ref(false);
const formRef = ref(null);
const form = reactive({
	name: '',
	code: '',
	groupId: '',
	enabled: true
});

const groupList = ref([
	{ id: 1, name: '默认分组' },
	{ id: 2, name: 'VIP分组' }
]);

const selectedGroupName = ref('');

onLoad((options) => {
	if (options.id) {
		isEdit.value = true;
		// 模拟获取设备详情
		form.name = '集控柜A';
		form.code = 'DEV001';
		form.groupId = 1;
		selectedGroupName.value = '默认分组';
	}
});

const onGroupConfirm = (e) => {
	const item = e.value[0];
	form.groupId = item.id;
	selectedGroupName.value = item.name;
};

const submit = () => {
	if (!form.name || !form.code) {
		return uni.showToast({ title: t('login.completeInfo'), icon: 'none' });
	}
	
	uni.showToast({ 
		title: isEdit.value ? t('common.operationSuccess') : t('common.operationSuccess'), 
		icon: 'success' 
	});
	
	setTimeout(() => {
		uni.navigateBack();
	}, 1500);
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f5f5;
}
.picker-value {
	font-size: 28rpx;
	color: #333;
	margin-left: 20rpx;
}
</style>