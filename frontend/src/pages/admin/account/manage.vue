<!-- 子账号管理页面 -->
<template>
	<view class="container">
		<uv-navbar :title="t('admin.account.title')" :placeholder="true" leftIcon="arrow-left" rightIcon="plus"
			@leftClick="goBack" @rightClick="openEditModal" />

		<!-- 搜索栏 -->
		<view class="search-bar">
			<uv-search v-model="keyword" :placeholder="t('admin.account.searchAccount')" shape="round"
				@search="handleSearch" @clear="handleSearch" bgColor="#f5f7fa"></uv-search>
		</view>

		<!-- 账号列表 -->
		<scroll-view scroll-y class="list-container">
			<view class="account-card" v-for="item in filteredList" :key="item.id">
				<view class="card-header-bg" @click="openEditModal(item)">
					<view class="card-header">
						<view class="user-info">
							<view class="avatar-placeholder">{{ item.account.charAt(0) }}</view>
							<view class="info-text">
								<text class="name">{{ item.account }}</text>
								<text class="role">{{ getRoleName(item.role) }}</text>
							</view>
						</view>
						<view class="actions">
							<uv-icon name="edit-pen" size="35" color="#3c9cff" @click="openEditModal(item)" />
							<uv-icon name="trash" size="35" color="#fa3534" @click.stop="handleDelete(item)"
								style="margin-left: 20rpx;" />
						</view>
					</view>
					<view class="card-body">
						<view class="info-row">
							<text class="label">{{ t('admin.account.email') }}</text>
							<text class="value">{{ item.email || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="label">{{ t('admin.account.phone') }}</text>
							<text class="value">{{ item.phone || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="label">{{ t('admin.account.status') }}</text>
							<uv-tags :text="getStatusName(item.status)"
								:type="item.status === '0' ? 'success' : 'warning'" size="mini"
								shape="circle"></uv-tags>
						</view>
						<view class="info-row">
							<text class="label">{{ t('admin.account.createTime') }}</text>
							<text class="value">{{ item.createdAt || '-' }}</text>
						</view>
					</view>
				</view>
			</view>

			<view class="empty-state" v-if="filteredList.length === 0 && !loading">
				<uv-empty :text="t('admin.account.noAccount')" mode="list" />
			</view>

			<!-- 底部占位，防止内容被遮挡 -->
			<view style="height: 50rpx;"></view>
		</scroll-view>

		<!-- 新增/编辑弹窗 -->
		<uv-popup ref="popup" mode="center" round="20" closeable>
			<view class="modal-content">
				<view class="modal-title">{{ isEdit ? t('admin.account.editAccount') : t('admin.account.addAccount') }}
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> {{ t('admin.account.account') }}</text>
					<uv-input v-model="formData.account" border="surround"
						:placeholder="t('admin.account.accountPlaceholder')" :disabled="isEdit" clearable></uv-input>
				</view>

				<view class="form-item" v-if="!isEdit">
					<text class="form-label"><text class="required">*</text> {{ t('admin.account.password') }}</text>
					<uv-input v-model="formData.password" border="surround" type="password"
						:placeholder="t('admin.account.passwordPlaceholder')" clearable></uv-input>
				</view>

				<view class="form-item" v-if="isEdit">
					<text class="form-label">{{ t('admin.account.newPassword') }}</text>
					<uv-input v-model="formData.password" border="surround" type="password"
						:placeholder="t('admin.account.newPasswordPlaceholder')" clearable></uv-input>
				</view>

				<view class="form-item">
					<text class="form-label">{{ t('admin.account.email') }}</text>
					<uv-input v-model="formData.email" border="surround"
						:placeholder="t('admin.account.emailPlaceholder')" clearable></uv-input>
				</view>

				<view class="form-item">
					<text class="form-label">{{ t('admin.account.phone') }}</text>
					<uv-input v-model="formData.phone" border="surround"
						:placeholder="t('admin.account.phonePlaceholder')" clearable></uv-input>
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> {{ t('admin.account.role') }}</text>
					<uv-picker ref="rolePicker" :columns="roleColumns" keyName="label" @confirm="onRoleConfirm"
						:itemHeight="80" :visibleItemCount="3">
						<view class="picker-trigger">
							<uv-input v-model="roleDisplayText" border="surround"
								:placeholder="t('admin.account.selectRole')" disabled
								suffixIcon="arrow-down"></uv-input>
						</view>
					</uv-picker>
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> {{ t('admin.account.status') }}</text>
					<uv-radio-group v-model="formData.status" placement="row">
						<uv-radio name="0" activeColor="#07c160" style="margin-right: 40rpx;">
							<template #icon>
								<view class="radio-custom-icon" :class="{ 'active': formData.status === '0' }"></view>
							</template>
							{{ t('admin.account.whitelist') }}
						</uv-radio>
						<uv-radio name="1" activeColor="#fa3534">
							<template #icon>
								<view class="radio-custom-icon" :class="{ 'active': formData.status === '1' }">
								</view>
							</template>
							{{ t('admin.account.blacklist') }}
						</uv-radio>
					</uv-radio-group>
				</view>

				<view class="form-item">
					<text class="form-label"><text class="required">*</text> {{ t('admin.account.permission') }}</text>
					<view class="permission-grid">
						<view v-for="perm in permissions" :key="perm.value"
							:class="['permission-item', { 'active': formData.rule.includes(perm.value) }]"
							@click="togglePermission(perm.value)">
							{{ perm.label }}
						</view>
					</view>
				</view>

				<view class="modal-footer">
					<uv-button type="info" plain @click="closeModal" customStyle="flex: 1">{{ t('common.cancel')
						}}</uv-button>
					<uv-button type="primary" @click="handleSubmit" customStyle="flex: 1">{{ t('common.confirm')
						}}</uv-button>
				</view>
			</view>
		</uv-popup>
	</view>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue';
import { onLoad } from '@dcloudio/uni-app';
import { useI18n } from 'vue-i18n';

const { t } = useI18n({ legacy: false });

// 角色选项配置
const roles = computed(() => [
	{ value: '0', label: t('admin.account.roleMerch') },
	{ value: '1', label: t('admin.account.roleManager') },
	{ value: '2', label: t('admin.account.roleAgent') },
]);

// 角色选择器列配置
const roleColumns = computed(() => [roles.value]);

// 权限选项
const permissions = computed(() => [
	{ value: 'device', label: t('admin.account.permDevice') },
	{ value: 'room', label: t('admin.account.permRoom') },
	{ value: 'order', label: t('admin.account.permOrder') },
	{ value: 'account', label: t('admin.account.permAccount') },
	{ value: 'huifu', label: t('admin.account.permHuifu') },
	{ value: 'rule', label: t('admin.account.permRule') },
]);

const keyword = ref('');
const list = ref([]);
const loading = ref(false);
const popup = ref(null);
const rolePicker = ref(null);
const isEdit = ref(false);
const currentId = ref(null);
const merchs_id = ref(''); // 商家ID参数

const formData = ref({
	account: '',
	password: '',
	email: '',
	phone: '',
	role: '1',
	status: '0',
	rule: ''
});

// 计算角色显示文本
const roleDisplayText = computed(() => {
	return getRoleName(formData.value.role);
});

// 过滤列表
const filteredList = computed(() => {
	if (!keyword.value) return list.value;
	const lowerKeyword = keyword.value.toLowerCase();
	return list.value.filter(item =>
		item.account.toLowerCase().includes(lowerKeyword) ||
		(item.email && item.email.toLowerCase().includes(lowerKeyword)) ||
		(item.phone && item.phone.includes(keyword.value))
	);
});

// 获取角色名称
const getRoleName = (value) => {
	if (!value) return '';
	const role = roles.value.find(item => item.value === value);
	return role ? role.label : value;
};

// 获取状态名称
const getStatusName = (value) => {
	if (!value) return '';
	return value === '0' ? t('admin.account.whitelist') : t('admin.account.blacklist');
};

// 页面加载时获取参数并加载数据
onLoad((options) => {
	merchs_id.value = options?.merchs_id || '';
	loadData();
});

// 加载子账号列表（根据商家ID参数筛选）
const loadData = async () => {
	loading.value = true;
	try {
		const params = {};
		if (merchs_id.value) {
			params.merchs_id = merchs_id.value;
		}
		const res = await uni.$uv.http.get('/submerch/list', {
			params: params,
			custom: { auth: true }
		});
		if (res.code === 200 && res.data) {
			list.value = res.data;
		}
	} catch (e) {
		console.error('加载子账号列表失败', e);
	} finally {
		loading.value = false;
	}
};

// 打开弹窗
const openEditModal = async (item = null) => {
	isEdit.value = !!item;
	if (item) {
		currentId.value = item.id;
		formData.value = {
			account: item.account || '',
			password: '',
			email: item.email || '',
			phone: item.phone || '',
			role: item.role || '1',
			status: item.status || '0',
			rule: item.rule || ''
		};
	} else {
		formData.value = {
			account: '',
			password: '',
			email: '',
			phone: '',
			role: '1',
			status: '0',
			rule: ''
		};
	}

	await nextTick();
	popup.value.open();
};

// 关闭弹窗
const closeModal = () => {
	popup.value.close();
};

// 角色选择确认
const onRoleConfirm = (e) => {
	if (e.value && e.value.length > 0) {
		formData.value.role = e.value[0].value;
	}
};

// 切换权限
const togglePermission = (value) => {
	const rule = formData.value.rule;
	if (rule.includes(value)) {
		formData.value.rule = rule.replace(value + ',', '').replace(',' + value, '').replace(value, '');
	} else {
		formData.value.rule = rule ? rule + ',' + value : value;
	}
};

// 提交表单
const handleSubmit = async () => {
	// 表单验证
	if (!formData.value.account.trim()) {
		uni.showToast({ title: t('admin.account.accountRequired'), icon: 'none' });
		return;
	}

	if (!isEdit.value && !formData.value.password.trim()) {
		uni.showToast({ title: t('admin.account.passwordRequired'), icon: 'none' });
		return;
	}

	if (!formData.value.role) {
		uni.showToast({ title: t('admin.account.roleRequired'), icon: 'none' });
		return;
	}

	if (!formData.value.status) {
		uni.showToast({ title: t('admin.account.statusRequired'), icon: 'none' });
		return;
	}

	if (!formData.value.rule) {
		uni.showToast({ title: t('admin.account.permissionRequired'), icon: 'none' });
		return;
	}

	try {
		uni.showLoading({ title: t('common.loading') });

		let result;
		if (isEdit.value) {
			// 更新子账号
			const updateData = {
				email: formData.value.email || null,
				phone: formData.value.phone || null,
				role: formData.value.role,
				status: formData.value.status,
				rule: formData.value.rule
			};
			if (formData.value.password.trim()) {
				updateData.password = formData.value.password;
			}

			result = await uni.$uv.http.put(`/submerch/${currentId.value}`, updateData, {
				custom: { auth: true }
			});
		} else {
			// 创建子账号
			result = await uni.$uv.http.post('/submerch', {
				account: formData.value.account,
				password: formData.value.password,
				email: formData.value.email || null,
				phone: formData.value.phone || null,
				role: formData.value.role,
				status: formData.value.status,
				rule: formData.value.rule
			}, {
				custom: { auth: true }
			});
		}

		if (result.code === 200) {
			uni.showToast({
				title: isEdit.value ? t('admin.account.editSuccess') : t('admin.account.addSuccess'),
				icon: 'success'
			});
			loadData();
			closeModal();
		} else {
			uni.showToast({ title: result.msg || t('common.error'), icon: 'none' });
		}
	} catch (e) {
		console.error('提交失败', e);
		uni.showToast({ title: t('common.error'), icon: 'none' });
	} finally {
		uni.hideLoading();
	}
};

// 删除子账号
const handleDelete = (item) => {
	uni.showModal({
		title: t('common.warning'),
		content: t('admin.account.deleteConfirm', { name: item.account }),
		confirmColor: '#fa3534',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: t('common.loading') });
					const result = await uni.$uv.http.delete(`/submerch/${item.id}`, {}, {
						custom: { auth: true }
					});

					if (result.code === 200) {
						uni.showToast({ title: t('admin.account.deleteSuccess'), icon: 'success' });
						loadData();
					} else {
						uni.showToast({ title: result.msg || t('common.error'), icon: 'none' });
					}
				} catch (e) {
					console.error('删除失败', e);
					uni.showToast({ title: t('common.error'), icon: 'none' });
				} finally {
					uni.hideLoading();
				}
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
		padding-bottom: 20rpx;
		border-bottom: 1rpx solid #f0f0f0;

		.user-info {
			display: flex;
			align-items: center;
			gap: 20rpx;

			.avatar-placeholder {
				width: 80rpx;
				height: 80rpx;
				border-radius: 50%;
				background: linear-gradient(135deg, #3c9cff 0%, #2b85e4 100%);
				display: flex;
				align-items: center;
				justify-content: center;
				color: #fff;
				font-size: 32rpx;
				font-weight: 600;
			}

			.info-text {
				display: flex;
				flex-direction: column;
				gap: 8rpx;

				.name {
					font-size: 32rpx;
					font-weight: 600;
					color: #333;
				}

				.role {
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
		padding-top: 20rpx;

		.info-row {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 12rpx 0;

			.label {
				font-size: 26rpx;
				color: #999;
			}

			.value {
				font-size: 26rpx;
				color: #333;
			}
		}
	}
}

.empty-state {
	padding: 100rpx 0;
}

.modal-content {
	width: 600rpx;
	padding: 40rpx;
	background: #fff;
	border-radius: 20rpx;

	.modal-title {
		font-size: 34rpx;
		font-weight: 600;
		text-align: center;
		margin-bottom: 40rpx;
		color: #333;
	}

	.form-item {
		margin-bottom: 30rpx;

		.form-label {
			font-size: 28rpx;
			color: #333;
			margin-bottom: 16rpx;
			display: block;

			.required {
				color: #fa3534;
				margin-right: 8rpx;
			}
		}
	}

	.permission-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;

		.permission-item {
			padding: 16rpx 24rpx;
			border-radius: 30rpx;
			font-size: 26rpx;
			color: #666;
			background: #f5f7fa;
			border: 1rpx solid #e8e8e8;
			transition: all 0.2s;

			&.active {
				background: #e6f4ff;
				color: #3c9cff;
				border-color: #3c9cff;
			}
		}
	}

	.modal-footer {
		display: flex;
		gap: 20rpx;
		margin-top: 40rpx;
	}
}

.radio-custom-icon {
	width: 36rpx;
	height: 36rpx;
	border-radius: 50%;
	border: 2rpx solid #d9d9d9;
	transition: all 0.2s;

	&.active {
		border-color: #3c9cff;
		background: #3c9cff;

		&::after {
			content: '';
			display: block;
			width: 16rpx;
			height: 16rpx;
			border-radius: 50%;
			background: #fff;
			margin: 7rpx auto;
		}
	}
}

.picker-trigger {
	width: 100%;
}
</style>