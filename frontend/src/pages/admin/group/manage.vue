<!-- 分组管理页面 -->
<template>
	<view class="container">
		<uv-navbar :title="t('admin.group.title')" :placeholder="true" @leftClick="goBack">
			<template #right>
				<uv-icon name="scan" size="25" @click="openQRCodePopup" />
			</template>
		</uv-navbar>

		<view class="search-bar">
			<uv-icon name="search" color="#999" size="24" />
			<input v-model="searchQuery" type="text" :placeholder="t('common.searchPlaceholder')" />
			<uv-icon v-if="searchQuery" name="close-circle-fill" color="#999" size="24" @click="clearSearch" />
		</view>

		<view class="stats-section">
			<view class="stat-row">
				<view class="stat-item" :class="{ active: !filterType }" @click="setFilter('')">
					<view class="color-indicator all"></view>
					<text class="stat-label">全部</text>
					<text class="stat-value">{{ groupStats.total }}</text>
				</view>
				<view class="stat-item" :class="{ active: filterType === '存柜' }" @click="setFilter('存柜')">
					<view class="color-indicator blue"></view>
					<text class="stat-label">存柜</text>
					<text class="stat-value">{{ groupStats.storage }}</text>
				</view>
				<view class="stat-item" :class="{ active: filterType === '零售' }" @click="setFilter('零售')">
					<view class="color-indicator green"></view>
					<text class="stat-label">零售</text>
					<text class="stat-value">{{ groupStats.retail }}</text>
				</view>
			</view>
		</view>

		<scroll-view scroll-y @scrolltolower="loadMore" lower-threshold="100">
			<view v-if="groupList.length === 0" class="empty-state">
				<uv-empty mode="data" textSize="32" iconSize="150" />
			</view>

			<view v-else class="group-list">
				<view v-for="item in filteredGroups" :key="item.id" class="group-card" @click="viewGroupDetail(item)">
					<view class="card-header">
						<view class="group-icon" :class="getGroupIconClass(item.type)">
							<uv-icon :name="getGroupIcon(item.type)" color="#fff" size="32" />
						</view>
						<text class="group-name">{{ item.name }}</text>
					</view>
					<view class="card-body">
						<view class="info-row">
							<text class="info-label">手机号</text>
							<text class="info-value">{{ item.phone || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">房间数量</text>
							<text class="info-value">{{ item.count || 0 }}间</text>
						</view>
						<view class="info-row">
							<text class="info-label">房间类型</text>
							<view class="info-tag" :class="getGroupIconClass(item.type)">
								{{ item.type }}
							</view>
						</view>
						<view class="info-row" v-if="item.type === '存柜'">
							<text class="info-label">所选规则</text>
							<view class="info-value">
								{{ item.ruleName || '-' }}
							</view>
						</view>
						<view class="info-row">
							<text class="info-label">所在位置</text>
							<text class="info-value">{{ item.location || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">创建时间</text>
							<text class="info-value">{{ formatTime(item.createdAt) }}</text>
						</view>
					</view>
					<view class="card-footer">
						<view class="actions">
							<view class="action-btn edit" @click.stop="viewGroupDetail(item)">
								<uv-icon name="edit-pen" color="#666" size="28" />
							</view>
							<view class="action-btn delete" @click.stop="deleteGroup(item.id)">
								<uv-icon name="trash" color="#ff4d4f" size="28" />
							</view>
						</view>
					</view>
				</view>
			</view>
		</scroll-view>

		<view class="bottom-bar" @click="openAddModal">
			<uv-button type="primary" text="新增分组" size="large" shape="circle" />
		</view>

		<!-- 弹框 -->
		<view>
			<!-- 删除确认弹窗 -->
			<uv-modal ref="deleteModalRef" :show="showDeleteModal" title="删除分组" :show-cancel-button="true"
				cancel-text="取消" confirm-text="确定" @confirm="handleDeleteConfirm" @cancel="handleDeleteCancel">
				<view class="delete-modal-content">
					<text class="delete-tip">确定要删除该分组及其下的所有房间吗？此操作不可恢复。</text>
					<view class="password-section">
						<text class="password-label">请输入商家密码</text>
						<uv-input v-model="deletePassword" type="password" placeholder="请输入商家登录密码"
							class="password-input" />
					</view>
				</view>
			</uv-modal>
			<!-- 新增分组弹窗 -->
			<uv-popup ref="addPopupRef" mode="bottom" :closeable="false" :safeAreaInsetBottom="true">
				<view class="modal-container">
					<view class="modal-header">
						<text class="modal-title">{{ t('admin.group.create') }}</text>
						<view class="close-btn" @click="closeAddModal">
							<uv-icon name="close" color="#666" size="30" />
						</view>
					</view>

					<view class="modal-body">
						<view class="form-item">
							<text class="form-label">{{ t('admin.group.name') }}<text class="required">*</text></text>
							<view class="form-input-wrapper" :class="{ 'has-error': formError.name }">
								<uv-input v-model="addForm.name" :placeholder="t('admin.group.namePlaceholder')"
									class="form-input" />
							</view>
							<text v-if="formError.name" class="error-text">{{ formError.name }}</text>
						</view>

						<view class="form-item">
							<text class="form-label">{{ t('admin.group.type') }}</text>
							<view class="type-radio-group">
								<view v-for="type in groupTypes" :key="type.value" class="type-option"
									:class="{ active: addForm.type === type.value }" @click="addForm.type = type.value">
									<view class="type-icon" :class="'type-' + type.value">
										<uv-icon :name="type.icon" color="#3c9cff" size="24" />
									</view>
									<text class="type-label">{{ type.label }}</text>
								</view>
							</view>
						</view>

						<view class="form-item">
							<text class="form-label">{{ t('admin.group.phone') }}<text class="required">*</text></text>
							<view class="form-input-wrapper" :class="{ 'has-error': formError.phone }">
								<uv-input v-model="addForm.phone" type="number"
									:placeholder="t('admin.group.phonePlaceholder')" class="form-input" />
							</view>
							<text v-if="formError.phone" class="error-text">{{ formError.phone }}</text>
						</view>

						<view class="form-item">
							<text class="form-label">{{ t('admin.group.location') }}</text>
							<view class="form-input-wrapper">
								<uv-input v-model="addForm.location" :placeholder="t('admin.group.locationPlaceholder')"
									class="form-input" />
							</view>
						</view>
					</view>

					<view class="modal-footer">
						<view class="btn btn-cancel" @click="closeAddModal">
							<text>{{ t('common.cancel') }}</text>
						</view>
						<view class="btn btn-confirm" @click="submitForm">
							<text>{{ t('common.confirm') }}</text>
						</view>
					</view>
				</view>
			</uv-popup>
			<!-- 二维码类型弹窗 -->
			<uv-popup ref="qrCodePopupRef" mode="bottom" :closeable="false" :safeAreaInsetBottom="true">
				<view class="modal-container">
					<view class="modal-header">
						<text class="modal-title">{{ t('admin.group.createqrcode') }}</text>
						<view class="close-btn" @click="closeAddModal">
							<uv-icon name="close" color="#666" size="30" />
						</view>
					</view>

					<view class="modal-body">
						<view class="form-item">
							<text class="form-label">{{ t('admin.group.qrcodetype') }}</text>
							<view class="type-radio-group">
								<view v-for="type in groupTypes" :key="type.value" class="type-option"
									:class="{ active: addForm.type === type.value }" @click="addForm.type = type.value">
									<view class="type-icon" :class="'type-' + type.value">
										<uv-icon :name="type.icon" color="#3c9cff" size="24" />
									</view>
									<text class="type-label">{{ type.label }}</text>
								</view>
							</view>
						</view>
					</view>

					<view class="modal-footer">
						<view class="btn btn-cancel" @click="closeAddModal">
							<text>{{ t('common.cancel') }}</text>
						</view>
						<view class="btn btn-confirm" @click="showQRCodeModal">
							<text>{{ t('common.confirm') }}</text>
						</view>
					</view>
				</view>
			</uv-popup>
			<!-- 二维码弹窗 -->
			<view v-if="showQRCodeSheet" class="action-sheet-mask" @click="closeQRCodeModal">
				<view class="qr-code-sheet" @click.stop>
					<view class="sheet-title">总二维码</view>

					<view class="qr-code-content">
						<view class="qr-code-wrapper">
							<uv-qrcode ref="qrcodeRef" :value="qrCodeContent" :size="200" :isQueueLoadImage="true"
								:h5SaveIsDownload="true" :h5DownloadName="`${addForm.type}-总码.png`" />
						</view>
						<view class="qr-code-desc">
							<text class="desc-text">扫码访问分组</text>
							<text class="desc-hint">{{ qrCodeContent }}</text>
						</view>
					</view>

					<view class="qr-actions">
						<view class="qr-btn cancel" @click="closeQRCodeModal">关闭</view>
						<view class="qr-btn confirm" @click="saveQRCode">保存图片</view>
						<view class="qr-btn copy" @click="copyQRCode">复制地址</view>
					</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import { useI18n } from 'vue-i18n';
// 导入二维码工具函数
import { generateQRCodeContent } from '@/utils/utils';

const { t } = useI18n();
const groupList = ref([]);
const merch = uni.getStorageSync('merch') || {};
const searchQuery = ref('');
const addPopupRef = ref(null);
const formError = ref({});
const loadingMore = ref(false);
const filterType = ref('');
const deletePassword = ref('');
const deleteGroupId = ref(null);
const showDeleteModal = ref(false);
const deleteModalRef = ref(null);
// 二维码类型弹窗ref
const qrCodePopupRef = ref(null);
// 二维码弹出层状态
const showQRCodeSheet = ref(false);
// 二维码内容
const qrCodeContent = ref('');
// 二维码组件ref
const qrcodeRef = ref(null);
// 新增分组表单数据
const addForm = ref({
	name: '',
	type: '存柜',
	phone: '',
	location: ''
});
// 分组类型选项
const groupTypes = computed(() => [
	{ value: '存柜', label: t('admin.group.storage'), icon: 'empty-favor' },
	{ value: '零售', label: t('admin.group.retail'), icon: 'gift' }
]);

// 每次显示页面时刷新分组列表
onShow(() => {
	fetchGroupList();
});

// 获取分组列表数据
const fetchGroupList = () => {
	uni.showLoading({ title: t('common.loading') });

	const params = {};
	if (merch && merch.id) {
		params.merchs_id = merch.id;
	}

	uni.$uv.http.get('/group/list', {
		params: params,
		custom: { auth: true }
	}).then((res) => {
		groupList.value = res.data.list || [];
		uni.hideLoading();
		uni.stopPullDownRefresh();
	}).catch((err) => {
		uni.hideLoading();
		uni.stopPullDownRefresh();
	});
};
// 打开二维码类型弹窗
const openQRCodePopup = () => {
	qrCodePopupRef.value.open();
};
// 显示二维码弹窗
const showQRCodeModal = () => {
	if (!merch || !merch.id) {
		uni.showToast({ title: '商家ID为空', icon: 'none' });
		return;
	}

	qrCodeContent.value = generateQRCodeContent('total', addForm.value.type, merch.id);
	qrCodePopupRef.value.close();
	showQRCodeSheet.value = true;
};
// 关闭二维码弹窗
const closeQRCodeModal = () => {
	showQRCodeSheet.value = false;
};
// 保存二维码图片
const saveQRCode = async () => {
	uni.showLoading({ title: '保存中' });

	try {
		// #ifdef H5
		// H5端使用uv-qrcode组件的save方法
		if (qrcodeRef.value && qrcodeRef.value.save) {
			await qrcodeRef.value.save({
				content: qrCodeContent.value,
				success: (res) => {
					console.log('二维码保存成功:', res);
					uni.showToast({ title: '保存成功', icon: 'success' });
				},
				fail: (err) => {
					console.error('二维码保存失败:', err);
					uni.showToast({ title: '保存失败', icon: 'none' });
				}
			});

			closeQRCodeModal();
		} else {
			uni.hideLoading();
			uni.showToast({ title: '组件未就绪', icon: 'none' });
		}
		// #endif

		// #ifndef H5
		// 小程序端使用uv-qrcode组件的save方法
		if (qrcodeRef.value && qrcodeRef.value.save) {
			await qrcodeRef.value.save({
				content: qrCodeContent.value,
				success: (res) => {
					console.log('二维码保存成功:', res);
					uni.showToast({ title: '保存成功', icon: 'success' });
				},
				fail: (err) => {
					console.error('二维码保存失败:', err);
					uni.showToast({ title: '保存失败', icon: 'none' });
				}
			});

			closeQRCodeModal();
		} else {
			uni.showToast({ title: '组件未就绪', icon: 'none' });
		}
		// #endif
	} catch (err) {
		console.error('保存二维码失败:', err);
		uni.showToast({ title: '保存失败', icon: 'none' });
	}
};
// 复制二维码地址
const copyQRCode = () => {
	uni.setClipboardData({
		data: qrCodeContent.value,
		success: () => {
			uni.showToast({ title: '复制成功', icon: 'success' })
		},
		fail: () => {
			uni.showToast({ title: '复制失败', icon: 'none' })
		}
	})
}
// 根据搜索关键词筛选分组列表
const filteredGroups = computed(() => {
	let result = groupList.value;

	if (filterType.value) {
		result = result.filter(item => item.type === filterType.value);
	}

	if (!searchQuery.value) return result;
	const query = searchQuery.value.toLowerCase();
	return result.filter(item =>
		item.name?.toLowerCase().includes(query) ||
		item.location?.toLowerCase().includes(query)
	);
});
// 统计分组数量
const groupStats = computed(() => {
	const storage = groupList.value.filter(g => g.type === '存柜').length;
	const retail = groupList.value.filter(g => g.type === '零售').length;
	return {
		total: groupList.value.length,
		storage,
		retail
	};
});
// 设置筛选类型
const setFilter = (type) => {
	if (filterType.value === type) {
		filterType.value = '';
	} else {
		filterType.value = type;
	}
};
// 获取分组图标名称
const getGroupIcon = (type) => {
	return (type === '零售') ? 'gift' : 'empty-favor';
};
// 获取分组图标样式类名
const getGroupIconClass = (type) => {
	return (type === '零售') ? 'retail' : 'storage';
};
// 清空搜索关键词
const clearSearch = () => {
	searchQuery.value = '';
};
// 加载更多数据
const loadMore = () => {
	if (loadingMore.value) return;
	loadingMore.value = true;

	setTimeout(() => {
		loadingMore.value = false;
	}, 1000);
};
// 查看分组详情
const viewGroupDetail = (item) => {
	uni.navigateTo({
		url: `/pages/admin/group/detail?id=${item.id}&rulesId=${item.rulesId}`
	});
};
// 打开添加模态框
const openAddModal = () => {
	addForm.value = {
		name: '',
		type: '存柜',
		phone: '',
		location: ''
	};
	formError.value = {};
	addPopupRef.value.open();
};
// 关闭添加模态框
const closeAddModal = () => {
	addPopupRef.value.close();
	qrCodePopupRef.value.close();
		
	addForm.value = {
		name: '',
		type: '存柜',
		phone: '',
		location: ''
	};
	formError.value = {};
};
// 表单验证
const validateForm = () => {
	formError.value = {};

	if (!addForm.value.name || !addForm.value.name.trim()) {
		formError.value.name = t('admin.group.nameRequired');
		return false;
	}

	if (addForm.value.name.length > 100) {
		formError.value.name = '分组名称不能超过100个字符';
		return false;
	}

	// 手机号必填校验
	if (!addForm.value.phone || !addForm.value.phone.trim()) {
		formError.value.phone = '请输入联系电话';
		return false;
	}

	// 手机号格式验证
	const phoneRegex = /^1[3-9]\d{9}$/;
	if (!phoneRegex.test(addForm.value.phone.trim())) {
		formError.value.phone = '请输入有效的手机号码';
		return false;
	}

	return true;
};
// 提交表单（创建分组）
const submitForm = async () => {
	if (!validateForm()) {
		return;
	}

	if (!merch.id) {
		uni.showToast({ title: t('admin.group.merchantRequired'), icon: 'none' });
		return;
	}

	uni.showLoading({ title: t('common.creating') });

	const formData = {
		name: addForm.value.name,
		type: addForm.value.type,
		phone: addForm.value.phone,
		merchs_id: merch.id
	};

	// 可选字段
	if (addForm.value.location) {
		formData.location = addForm.value.location;
	}

	await uni.$uv.http.post('/group', formData, {
		custom: { auth: true }
	})
		.then((res) => {
			uni.hideLoading();
			if (res.code === 200) {
				uni.showToast({ title: t('common.createSuccess'), icon: 'success' });
				closeAddModal();
				fetchGroupList();
			} else {
				uni.showToast({ title: res.msg || t('common.createFailed'), icon: 'none' });
			}
		})
		.catch((err) => {
			uni.hideLoading();
		});
};
// 删除分组
const deleteGroup = (id) => {
	deleteGroupId.value = id;
	deletePassword.value = '';
	deleteModalRef.value.open();
};
// 处理删除确认
const handleDeleteConfirm = () => {
	if (!deletePassword.value.trim()) {
		uni.showToast({ title: '请输入商家密码', icon: 'none' });
		return;
	}

	showDeleteModal.value = false;
	confirmDeleteGroup(deleteGroupId.value, deletePassword.value);
};
// 处理删除取消
const handleDeleteCancel = () => {
	showDeleteModal.value = false;
	deletePassword.value = '';
};
// 确认删除分组
const confirmDeleteGroup = (id, password) => {
	uni.showLoading({ title: t('common.deleting') });

	uni.$uv.http.delete(`/group/${id}`, { password: password }, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		if (res.code === 200) {
			uni.showToast({ title: t('common.deleteSuccess'), icon: 'success' });
			fetchGroupList();
		} else {
			uni.showToast({ title: res.msg || t('common.deleteFailed'), icon: 'none' });
		}
	}).catch((err) => {
		uni.hideLoading();
		console.error('删除分组失败:', err);
		uni.showToast({ title: t('common.deleteFailed'), icon: 'none' });
	});
};
// 格式化时间
const formatTime = (time) => {
	if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
};
// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/room/manage'
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 140rpx;
}

.search-bar {
	display: flex;
	align-items: center;
	background: #fff;
	margin: 20rpx;
	padding: 0 24rpx;
	height: 80rpx;
	border-radius: 40rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);

	input {
		flex: 1;
		font-size: 28rpx;
		color: #333;
	}
}

.stats-section {
	padding: 0 20rpx;
	margin-bottom: 24rpx;

	.stat-row {
		display: flex;
		background: #fff;
		border-radius: 16rpx;
		padding: 20rpx;
		box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);
		justify-content: space-around;
	}

	.stat-item {
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: 8rpx 16rpx;
		border-radius: 8rpx;
		transition: all 0.2s;

		&:active {
			transform: scale(0.95);
		}

		&.active {
			background: rgba(0, 0, 0, 0.05);
			box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.1);
		}

		.color-indicator {
			width: 24rpx;
			height: 24rpx;
			border-radius: 4rpx;
			margin-right: 8rpx;

			&.all {
				background: #9E9E9E;
			}

			&.blue {
				background: linear-gradient(135deg, #3c9cff, #2b85e4);
			}

			&.green {
				background: linear-gradient(135deg, #4CAF50, #388E3C);
			}
		}

		.stat-label {
			font-size: 20rpx;
			color: #666;
			margin-right: 8rpx;
		}

		.stat-value {
			font-size: 28rpx;
			font-weight: bold;
			color: #333;
		}
	}
}

.loading-more {
	padding: 30rpx;
	display: flex;
	justify-content: center;
	align-items: center;
}

.empty-state {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 100rpx 40rpx;

	.empty-icon {
		width: 160rpx;
		height: 160rpx;
		background: #f5f5f5;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 30rpx;
	}

	.empty-text {
		font-size: 28rpx;
		color: #999;
		margin-bottom: 30rpx;
	}
}

.group-list {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.group-card {
	width: 86%;
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.06);
	transition: all 0.2s ease;

	&:active {
		transform: scale(0.99);
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.08);
	}
}

.card-header {
	display: flex;
	align-items: center;
	margin-bottom: 20rpx;
	padding-bottom: 20rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.group-icon {
	width: 72rpx;
	height: 72rpx;
	border-radius: 18rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-right: 20rpx;

	&.storage {
		background: linear-gradient(135deg, #3c9cff, #2b85e4);
	}

	&.retail {
		background: linear-gradient(135deg, #4CAF50, #388E3C);
	}
}

.group-name {
	font-size: 32rpx;
	font-weight: 600;
	color: #333;
}

.card-body {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
	margin-bottom: 20rpx;
}

.info-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.info-label {
	font-size: 26rpx;
	color: #999;
}

.info-value {
	font-size: 26rpx;
	color: #333;
	font-weight: 500;
}

.info-tag {
	font-size: 22rpx;
	padding: 6rpx 16rpx;
	border-radius: 8rpx;

	&.storage {
		background: rgba(60, 156, 255, 0.15);
		color: #3c9cff;
	}

	&.retail {
		background: rgba(76, 175, 80, 0.15);
		color: #4CAF50;
	}
}

.card-footer {
	display: flex;
	justify-content: flex-end;
	padding-top: 16rpx;
	border-top: 1rpx solid #f0f0f0;
}

.actions {
	display: flex;
	gap: 24rpx;
}

.action-btn {
	width: 64rpx;
	height: 64rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	background: #f5f5f5;
	transition: all 0.2s ease;

	&:active {
		transform: scale(0.9);
	}

	&.delete {
		background: rgba(250, 53, 52, 0.1);
	}
}

.modal-container {
	background: #f8f9fa;
	border-radius: 24rpx 24rpx 0 0;
	padding: 0;
	display: flex;
	flex-direction: column;
	max-height: 90vh;
	position: relative;
	z-index: 1000;
}

.modal-header {
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 40rpx 30rpx;
	background: #fff;
	border-radius: 24rpx 24rpx 0 0;
	position: relative;

	.modal-title {
		font-size: 36rpx;
		font-weight: 600;
		color: #333;
	}

	.close-btn {
		position: absolute;
		right: 30rpx;
		top: 50%;
		transform: translateY(-50%);
		width: 56rpx;
		height: 56rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}
}

.modal-body {
	flex: 1;
	overflow-y: auto;
	padding: 30rpx;
	padding-bottom: 10rpx;
	background: #f8f9fa;
	flex-shrink: 1;
}

.form-item {
	margin-bottom: 32rpx;

	&:last-child {
		margin-bottom: 0;
	}

	.form-label {
		display: block;
		font-size: 28rpx;
		font-weight: 500;
		color: #333;
		margin-bottom: 16rpx;

		.required {
			color: #ff4d4f;
			margin-left: 6rpx;
		}
	}
}

.form-input-wrapper {
	display: flex;
	align-items: center;
	background: #fff;
	border-radius: 16rpx;
	padding: 0 24rpx;
	height: 92rpx;
	transition: all 0.3s ease;

	.form-input {
		flex: 1;
		font-size: 28rpx;
		color: #333;
	}

	&:focus-within {
		border-color: #3c9cff;
	}

	&.has-error {
		border-color: #ff4d4f;
	}
}

.error-text {
	font-size: 24rpx;
	color: #ff4d4f;
	margin-top: 12rpx;
	padding-left: 4rpx;
}

.type-radio-group {
	display: flex;
	gap: 20rpx;
}

.type-option {
	flex: 1;
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 28rpx 20rpx;
	background: #fff;
	border-radius: 16rpx;
	border: 2rpx solid #e8e8e8;
	transition: all 0.3s ease;

	&.active {
		border-color: #3c9cff;
		background: rgba(60, 156, 255, 0.05);

		.type-label {
			color: #3c9cff;
			font-weight: 600;
		}
	}

	&:active {
		transform: scale(0.96);
	}
}

.type-icon {
	width: 56rpx;
	height: 56rpx;
	border-radius: 14rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	margin-bottom: 12rpx;

	&.type-storage {
		background: linear-gradient(135deg, #3c9cff, #2b85e4);
	}

	&.type-retail {
		background: linear-gradient(135deg, #4CAF50, #388E3C);
	}
}

.type-label {
	font-size: 26rpx;
	color: #666;
}

.modal-footer {
	display: flex;
	gap: 20rpx;
	padding: 24rpx 30rpx;
	padding-bottom: calc(40rpx + env(safe-area-inset-bottom));
	background: #f8f9fa;
	flex-shrink: 0;

	.btn {
		flex: 1;
		height: 88rpx;
		border-radius: 16rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 32rpx;
		font-weight: 500;
		transition: all 0.3s ease;

		&:active {
			transform: scale(0.96);
		}
	}

	.btn-cancel {
		background: #d5d5d6;
		color: #666;

		&:active {
			background: #e8e8e8;
		}
	}

	.btn-confirm {
		background: #3c9cff;
		color: #fff;
	}
}

.bottom-bar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	padding: 20rpx 30rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	background: #fff;
	box-shadow: 0 -4rpx 12rpx rgba(0, 0, 0, 0.06);
	z-index: 666;
}

/* 删除弹窗样式 */
.delete-modal-content {
	padding: 20rpx 0;
}

.delete-tip {
	font-size: 28rpx;
	color: #666;
	line-height: 1.6;
	text-align: center;
	display: block;
	margin-bottom: 30rpx;
}

.password-section {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
}

.password-label {
	font-size: 26rpx;
	color: #333;
	font-weight: 500;
}

.password-input {
	background: #fff;
	border-radius: 12rpx;
	padding: 0 20rpx;
	height: 80rpx;
	font-size: 28rpx;
}

/* 添加设备弹出层样式 */
.action-sheet-mask {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: flex-end;
	z-index: 1000;
}

.sheet-title {
	text-align: center;
	padding: 40rpx;
	font-size: 32rpx;
	font-weight: 600;
	color: #333;
	border-bottom: 1rpx solid #f0f0f0;
}

/* 二维码弹窗样式 */
.qr-code-sheet {
	width: 100%;
	background: #fff;
	border-radius: 32rpx 32rpx 0 0;
	padding-bottom: env(safe-area-inset-bottom);
	max-height: 80vh;
	overflow-y: auto;
}

.qr-code-content {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 32rpx;
}

.qr-code-wrapper {
	width: 360rpx;
	height: 360rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	background: #ffffff;
	border: 2rpx solid #e8e8e8;
	border-radius: 20rpx;
	box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
	overflow: hidden;
	flex-shrink: 0;
}

.qr-code-wrapper canvas {
	width: 100% !important;
	height: 100% !important;
	display: block;
}

/* 微信小程序特殊处理：禁止二维码区域滚动 */
/* #ifdef MP-WEIXIN */
.qr-code-content {
	padding: 48rpx 32rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 32rpx;
	touch-action: none;
}

.qr-code-wrapper {
	pointer-events: none;
}
/* #endif */

.qr-code-desc {
	text-align: center;
	width: 100%;
	
	.desc-text {
		font-size: 32rpx;
		font-weight: 600;
		color: #333333;
		display: block;
		margin-bottom: 16rpx;
	}
	
	.desc-hint {
		font-size: 24rpx;
		color: #999999;
		display: block;
		word-break: break-all;
		max-width: 100%;
		line-height: 1.5;
		padding: 0 16rpx;
	}
}

.qr-actions {
	display: flex;
	padding: 0 32rpx 32rpx;
	gap: 24rpx;
	border-top: 1rpx solid #f5f5f5;
	padding-top: 28rpx;
	
	.qr-btn {
		flex: 1;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 44rpx;
		font-size: 30rpx;
		font-weight: 500;
		
		&.cancel {
			background: #f5f7fa;
			color: #666666;
			border: 1rpx solid #e8e8e8;
		}
		
		&.confirm {
			background: linear-gradient(135deg, #3c9cff 0%, #2979ff 100%);
			color: #ffffff;
			box-shadow: 0 4rpx 16rpx rgba(60, 156, 255, 0.3);
		}

		&.copy {
			background: linear-gradient(135deg, #ffc83c 0%, #ffab29 100%);
			color: #ffffff;
			box-shadow: 0 0.125rem 0.5rem rgb(255 60 60 / 30%);
		}
	}
}
</style>