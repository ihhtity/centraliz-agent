<template>
	<view class="container">
		<uv-navbar :title="t('admin.account.title')" :placeholder="true" rightIcon="plus"
			@leftClick="goBack" @rightClick="openEditModal" />

		<!-- 搜索栏 -->
		<view class="search-bar">
			<uv-search v-model="keyword" placeholder="搜索账号/姓名" shape="round" @search="handleSearch"
				@clear="handleSearch" bgColor="#f5f7fa"></uv-search>
		</view>

		<!-- 账号列表 -->
		<scroll-view scroll-y class="list-container">
			<view class="account-card" v-for="item in filteredList" :key="item.id">
				<view class="card-header-bg" @click="openEditModal(item)">
					<view class="card-header">
						<view class="user-info">
							<view class="avatar-placeholder">{{ item.name.charAt(0) }}</view>
							<view class="info-text">
								<text class="name">{{ item.name }}</text>
								<text class="username">@{{ item.username }}</text>
							</view>
						</view>
						<view class="actions">
							<uv-icon name="edit-pen" size="35" color="#3c9cff" @click="openEditModal(item)" />
							<uv-icon name="trash" size="35" color="#fa3534" @click="handleDelete(item)"
								style="margin-left: 20rpx;" />
						</view>
					</view>
					<view class="card-body">
						<view class="info-row">
							<text class="label">所属厂商</text>
							<uv-tags :text="getVendorName(item.vendor)" plain size="mini" type="warning"
								shape="circle"></uv-tags>
						</view>
						<view class="info-row">
							<text class="label">角色权限</text>
							<uv-tags :text="item.role === 'admin' ? '管理员' : '操作员'"
								:type="item.role === 'admin' ? 'primary' : 'success'" size="mini"
								shape="circle"></uv-tags>
						</view>
					</view>
				</view>
			</view>

			<view class="empty-state" v-if="filteredList.length === 0">
				<uv-empty text="暂无账号" mode="list" />
			</view>

			<!-- 底部占位，防止内容被遮挡 -->
			<view style="height: 50rpx;"></view>
		</scroll-view>

		<!-- 新增/编辑弹窗 -->
		<uv-popup ref="popup" mode="center" round="20" closeable>
			<view class="modal-content">
				<view class="modal-title">{{ isEdit ? '编辑账号' : '新增账号' }}</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> 姓名</text>
					<uv-input v-model="formData.name" border="surround" placeholder="请输入真实姓名" clearable></uv-input>
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> 登录账号</text>
					<uv-input v-model="formData.username" border="surround" placeholder="请输入登录账号" :disabled="isEdit"
						clearable></uv-input>
				</view>

				<view class="form-item" v-if="!isEdit">
					<text class="form-label"><text class="required">*</text> 初始密码</text>
					<uv-input v-model="formData.password" border="surround" type="password" placeholder="请设置初始密码"
						clearable></uv-input>
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> 所属厂商</text>
					<!-- 修改：使用 uv-input 模拟选择框，点击触发 picker -->
					<uv-input v-model="vendorDisplayText" border="surround" placeholder="请选择所属厂商" disabled
						suffixIcon="arrow-down" @click="openVendorPicker"></uv-input>
					<!-- 隐藏的 picker 组件 -->
					<uv-picker ref="vendorPicker" :columns="vendorColumns" keyName="label" @confirm="onVendorConfirm"
						:itemHeight="80" :visibleItemCount="5"></uv-picker>
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> 角色权限</text>
					<uv-radio-group v-model="formData.role" placement="row">
						<uv-radio name="admin" activeColor="#3c9cff" style="margin-right: 40rpx;">
							<template #icon>
								<view class="radio-custom-icon" :class="{ 'active': formData.role === 'admin' }"></view>
							</template>
							管理员
						</uv-radio>
						<uv-radio name="operator" activeColor="#3c9cff">
							<template #icon>
								<view class="radio-custom-icon" :class="{ 'active': formData.role === 'operator' }">
								</view>
							</template>
							操作员
						</uv-radio>
					</uv-radio-group>
				</view>

				<view class="modal-footer">
					<uv-button type="info" plain @click="closeModal" customStyle="flex: 1">取消</uv-button>
					<uv-button type="primary" @click="handleSubmit" customStyle="flex: 1">确定</uv-button>
				</view>
			</view>
		</uv-popup>
	</view>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue';

// 厂商选项配置
const vendors = [
	{ value: 'baoshili', label: '宝士力得' },
	{ value: 'ttlock', label: '通通锁' },
	{ value: 'kuodao', label: '阔道' },
	{ value: 'xunming', label: '迅鸣' },
];
const vendorColumns = ref([vendors]);

const keyword = ref('');
const list = ref([]);
const popup = ref(null);
const vendorPicker = ref(null);
const isEdit = ref(false);
const currentId = ref(null);

const formData = ref({
	name: '',
	username: '',
	password: '',
	vendor: '',
	role: 'operator'
});

// 新增：计算厂商显示文本，用于输入框展示
const vendorDisplayText = computed(() => {
	return getVendorName(formData.value.vendor);
});

// 初始化加载数据
onMounted(() => {
	loadData();
});

// 从本地存储加载数据
const loadData = () => {
	const savedData = uni.getStorageSync('account_list');
	if (savedData && Array.isArray(savedData)) {
		list.value = savedData;
	} else {
		// 默认模拟数据
		list.value = [
			{ id: 1, name: '张三', username: 'zhangsan', vendor: 'baoshili', role: 'admin' },
			{ id: 2, name: '李四', username: 'lisi', vendor: 'ttlock', role: 'operator' },
			{ id: 3, name: '王五', username: 'wangwu', vendor: 'kuodao', role: 'operator' },
			{ id: 4, name: '赵六', username: 'zhaoliu', vendor: 'xunming', role: 'operator' },
		];
		saveData();
	}
};

// 保存数据到本地存储
const saveData = () => {
	uni.setStorageSync('account_list', list.value);
};

// 过滤列表
const filteredList = computed(() => {
	if (!keyword.value) return list.value;
	const lowerKeyword = keyword.value.toLowerCase();
	return list.value.filter(item =>
		item.name.toLowerCase().includes(lowerKeyword) ||
		item.username.toLowerCase().includes(lowerKeyword)
	);
});

// 获取厂商显示名称
const getVendorName = (value) => {
	if (!value) return '';
	const v = vendors.find(item => item.value === value);
	return v ? v.label : value;
};

// 打开弹窗
const openEditModal = async (item = null) => {
	isEdit.value = !!item;
	if (item) {
		currentId.value = item.id;
		// 编辑时不回填密码，避免明文显示或覆盖
		formData.value = {
			name: item.name,
			username: item.username,
			password: '',
			vendor: item.vendor,
			role: item.role
		};
	} else {
		formData.value = { name: '', username: '', password: '', vendor: '', role: 'operator' };
	}

	await nextTick();
	popup.value.open();
};

const closeModal = () => {
	popup.value.close();
};

// 新增：打开厂商选择器
const openVendorPicker = () => {
	vendorPicker.value.open();
};

// 厂商选择确认
const onVendorConfirm = (e) => {
	if (e.value && e.value.length > 0) {
		// e.value[0] 是选中的对象，取其中的 value 字段
		formData.value.vendor = e.value[0].value;
	}
};

// 提交表单
const handleSubmit = () => {
	// 表单验证
	if (!formData.value.name.trim()) {
		uni.showToast({ title: '请输入姓名', icon: 'none' });
		return;
	}
	if (!formData.value.username.trim()) {
		uni.showToast({ title: '请输入登录账号', icon: 'none' });
		return;
	}

	// 新增：用户名唯一性校验 (排除自身)
	const isUsernameExist = list.value.some(item =>
		item.username === formData.value.username && item.id !== currentId.value
	);
	if (isUsernameExist) {
		uni.showToast({ title: '该登录账号已存在', icon: 'none' });
		return;
	}

	if (!isEdit.value && !formData.value.password.trim()) {
		uni.showToast({ title: '请设置初始密码', icon: 'none' });
		return;
	}
	if (!formData.value.vendor) {
		uni.showToast({ title: '请选择所属厂商', icon: 'none' });
		return;
	}

	if (isEdit.value) {
		// 模拟更新
		const index = list.value.findIndex(item => item.id === currentId.value);
		if (index !== -1) {
			const updateData = {
				name: formData.value.name,
				username: formData.value.username,
				vendor: formData.value.vendor,
				role: formData.value.role
			};

			// 只有当密码不为空时才更新密码
			if (formData.value.password.trim()) {
				updateData.password = formData.value.password;
			} else {
				// 保持原密码不变
				updateData.password = list.value[index].password;
			}

			list.value[index] = {
				...list.value[index],
				...updateData
			};
			uni.showToast({ title: '修改成功', icon: 'success' });
		}
	} else {
		// 模拟新增
		const newItem = {
			id: Date.now(),
			name: formData.value.name,
			username: formData.value.username,
			password: formData.value.password,
			vendor: formData.value.vendor,
			role: formData.value.role
		};
		list.value.unshift(newItem);
		uni.showToast({ title: '添加成功', icon: 'success' });
	}

	saveData();
	closeModal();
};

// 删除
const handleDelete = (item) => {
	// 新增：保护最后一个管理员不被删除
	if (item.role === 'admin') {
		const adminCount = list.value.filter(i => i.role === 'admin').length;
		if (adminCount <= 1) {
			uni.showToast({ title: '无法删除最后一个管理员账号', icon: 'none' });
			return;
		}
	}

	uni.showModal({
		title: '警告',
		content: `确定要删除账号 "${item.name}" (${item.username}) 吗？此操作不可恢复。`,
		confirmColor: '#fa3534',
		success: (res) => {
			if (res.confirm) {
				list.value = list.value.filter(i => i.id !== item.id);
				saveData();
				uni.showToast({ title: '删除成功', icon: 'success' });
			}
		}
	});
};

const handleSearch = () => {
	// 触发computed更新，无需额外操作
};

const goBack = () => {
	uni.navigateBack();
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	display: flex;
	flex-direction: column;
}

.search-bar {
	padding: 20rpx 30rpx;
	display: flex;
	align-items: center;
	gap: 20rpx;
	background: #fff;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);
	z-index: 10;

	.add-btn {
		background: linear-gradient(135deg, #3c9cff 0%, #2b85e4 100%);
		color: #fff;
		padding: 12rpx 24rpx;
		border-radius: 30rpx;
		font-size: 26rpx;
		display: flex;
		align-items: center;
		gap: 8rpx;
		flex-shrink: 0;
		box-shadow: 0 4rpx 10rpx rgba(60, 156, 255, 0.3);
		transition: all 0.2s;

		&:active {
			transform: scale(0.95);
			opacity: 0.9;
		}
	}
}

.list-container {
	flex: 1;
	padding: 20rpx 30rpx;
	box-sizing: border-box;
}

.account-card {
	background: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	margin-bottom: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	transition: transform 0.2s;

	&:active {
		transform: scale(0.98);
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24rpx;

		.user-info {
			display: flex;
			align-items: center;
			gap: 20rpx;

			.avatar-placeholder {
				width: 80rpx;
				height: 80rpx;
				background: #eef2f7;
				color: #3c9cff;
				border-radius: 50%;
				display: flex;
				align-items: center;
				justify-content: center;
				font-size: 32rpx;
				font-weight: bold;
			}

			.info-text {
				display: flex;
				flex-direction: column;

				.name {
					font-size: 32rpx;
					font-weight: bold;
					color: #333;
					margin-bottom: 4rpx;
				}

				.username {
					font-size: 24rpx;
					color: #999;
				}
			}
		}

		.actions {
			display: flex;
			align-items: center;
		}
	}

	.card-body {
		background: #f9fafc;
		border-radius: 12rpx;
		padding: 20rpx;

		.info-row {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 12rpx;
			font-size: 26rpx;

			&:last-child {
				margin-bottom: 0;
			}

			.label {
				color: #666;
				font-weight: 500;
			}
		}
	}
}

.modal-content {
	width: 640rpx;
	padding: 40rpx;
	background: #fff;
	border-radius: 24rpx;

	.modal-title {
		font-size: 34rpx;
		font-weight: bold;
		text-align: center;
		margin-bottom: 40rpx;
		color: #333;
	}

	.form-item {
		margin-bottom: 30rpx;

		.form-label {
			font-size: 28rpx;
			color: #333;
			margin-bottom: 12rpx;
			display: block;
			font-weight: 500;

			.required {
				color: #fa3534;
				margin-right: 4rpx;
			}
		}

		.radio-custom-icon {
			width: 32rpx;
			height: 32rpx;
			border: 2rpx solid #dcdfe6;
			border-radius: 50%;
			margin-right: 10rpx;
			position: relative;

			&.active {
				border-color: #3c9cff;
				background: #3c9cff;

				&::after {
					content: '';
					position: absolute;
					top: 50%;
					left: 50%;
					transform: translate(-50%, -50%);
					width: 12rpx;
					height: 12rpx;
					background: #fff;
					border-radius: 50%;
				}
			}
		}
	}

	.modal-footer {
		display: flex;
		gap: 20rpx;
		margin-top: 40rpx;
	}
}

.empty-state {
	margin-top: 100rpx;
}
</style>