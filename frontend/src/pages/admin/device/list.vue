<!-- 设备列表页面 -->
<template>
	<view class="container">
		<uv-navbar title="设备列表" :placeholder="true" @leftClick="goBack">
			<template #right>
				<uv-icon name="plus" size="25" color="#3c9cff" @click="AddSheet" />
			</template>
		</uv-navbar>

		<!-- 搜索栏 -->
		<view class="search-section">
			<uv-search 
				v-model="searchName" 
				placeholder="搜索设备名称" 
				:showAction="false"
				@search="handleSearch"
				@clear="handleClearSearch"
			/>
		</view>

		<!-- 设备类型筛选 -->
		<view class="filter-section">
			<uv-tabs :list="filterTypes" v-model="filterType" @click="handleTypeChange" />
		</view>

		<view class="list-wrapper">
			<view class="device-list">
				<view 
					v-for="item in filteredDeviceList" 
					:key="item.id" 
					class="device-item"
				>
					<view class="device-left" @click.stop="showEditSheetFunc(item)">
						<view class="device-info">
							<view class="info-row">
								<text class="info-label">名称：</text>
								<text class="info-value">{{ item.name }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">状态：</text>
								<text :class="['info-value', 'status-text', getStatusClass(item.status)]">{{ item.status }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">类型：</text>
								<text class="info-value">{{ item.type }}</text>
							</view>
							<view class="info-row">
								<text class="info-label">创建时间：</text>
								<text class="info-value">{{ formatTime(item.createdAt) }}</text>
							</view>
						</view>
					</view>
					
					<view class="device-right">
						<view class="action-buttons">
							<text class="action-btn edit-btn" @click.stop="showEditSheetFunc(item)">编辑</text>
							<text class="action-btn change-btn" @click.stop="showActionSheet(item.id)">更换</text>
							<text class="action-btn delete-btn" @click.stop="deleteDevice(item.id)">删除</text>
						</view>
					</view>
				</view>
			</view>
			
			<uv-empty v-if="filteredDeviceList.length === 0" mode="list" text="暂无设备" />
		</view>

		<!-- 分页组件 -->
		<view class="pagination-wrapper">
			<view class="pagination">
				<view 
					class="page-btn prev-btn" 
					:class="{ disabled: currentPage === 1 }"
					@click="handlePrevPage"
				>
					<uv-icon name="arrow-left" size="16" color="#666" />
					<text>上一页</text>
				</view>
				
				<view class="page-info">
					<text class="current-page">{{ currentPage }}</text>
					<text class="separator">/</text>
					<text class="total-pages">{{ totalPages }}</text>
				</view>
				
				<view class="page-jump">
					<text>跳转</text>
					<input 
						v-model="jumpPage" 
						class="jump-input" 
						type="number"
						placeholder="页码"
						@blur="handleJumpPage"
						@confirm="handleJumpPage"
					/>
					<text>页</text>
				</view>
				
				<view 
					class="page-btn next-btn" 
					:class="{ disabled: currentPage === totalPages || totalPages === 0 }"
					@click="handleNextPage"
				>
					<text>下一页</text>
					<uv-icon name="arrow-right" size="16" color="#666" />
				</view>
			</view>
		</view>
		
		<!-- 更换设备弹出层 -->
		<view v-if="showSheet" class="action-sheet-mask" @click="hideActionSheet">
			<view class="action-sheet" @click.stop>
				<view class="sheet-title">选择更换设备操作</view>
				<view class="sheet-btns">
					<view class="sheet-btn" @click="handleScan">
						<uv-icon name="scan" size="32" color="#3c9cff" />
						<text>扫码</text>
					</view>
					<view class="sheet-btn" @click="handleEdit">
						<uv-icon name="edit-pen" size="32" color="#3c9cff" />
						<text>手动</text>
					</view>
				</view>
				<view class="sheet-cancel" @click="hideActionSheet">取消</view>
			</view>
		</view>
		
		<!-- 添加设备弹出层 -->
		<view v-if="showAddSheet" class="action-sheet-mask" @click="hideAddSheet">
			<view class="add-device-sheet" @click.stop>
				<view class="sheet-title">添加设备</view>
				
				<view class="type-section">
					<text class="section-label">选择设备类型</text>
					<uv-tabs :list="deviceTypes" @click="clickType" />
				</view>
				
				<view class="code-section">
					<text class="section-label">设备编码</text>
					<view class="code-input-wrapper">
						<input 
							v-model="deviceCode" 
							class="code-input" 
							placeholder="请输入设备编码"
							placeholder-class="input-placeholder"
						/>
						<view class="scan-btn" @click="scanDeviceCode">
							<uv-icon name="scan" size="28" color="#3c9cff" />
						</view>
					</view>
				</view>
				
				<view class="add-actions">
					<view class="add-btn cancel" @click="hideAddSheet">取消</view>
					<view class="add-btn confirm" @click="handleAddDevice">确认添加</view>
				</view>
			</view>
		</view>
		
		<!-- 编辑设备名称弹出层 -->
		<view v-if="showEditSheet" class="action-sheet-mask" @click="hideEditSheet">
			<view class="edit-device-sheet" @click.stop>
				<view class="sheet-title">修改设备名称</view>
				
				<view class="edit-form">
					<view class="form-item">
						<input 
							v-model="editDeviceName" 
							class="name-input" 
							placeholder="请输入设备名称"
							placeholder-class="input-placeholder"
						/>
					</view>
				</view>
				
				<view class="edit-actions">
					<view class="edit-btn cancel" @click="hideEditSheet">取消</view>
					<view class="edit-btn confirm" @click="handleUpdateDevice">确认</view>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref, computed } from 'vue';
import { onLoad, onPullDownRefresh } from '@dcloudio/uni-app';

// 设备列表数据
const deviceList = ref([]);
// 弹出层状态控制
const showSheet = ref(false);
const showAddSheet = ref(false);
const showEditSheet = ref(false);
// 当前操作的设备信息
const currentDeviceId = ref(null);
const currentDeviceName = ref('');
const editDeviceName = ref('');
const selectedType = ref('集控');
const deviceCode = ref('');

// 筛选相关
const filterType = ref(0);
const searchName = ref('');

// 分页相关
const currentPage = ref(1);
const pageSize = ref(50);
const total = ref(0);
const jumpPage = ref('');

// 筛选用设备类型选项（包含"全部"）
const filterTypes = [
	{ name: '全部' },
	{ name: '集控' },
	{ name: '摄像头' },
	{ name: '蓝牙锁' },
	{ name: '温湿度' },
	{ name: '门禁' },
	{ name: '照明' }
];

// 添加设备用设备类型选项
const deviceTypes = [
	{ name: '集控', icon: 'server' },
	{ name: '摄像头', icon: 'video' },
	{ name: '蓝牙锁', icon: 'lock' },
	{ name: '温湿度', icon: 'cloud' },
	{ name: '门禁', icon: 'key' },
	{ name: '照明', icon: 'bulb' }
];

// 计算总页数
const totalPages = computed(() => {
	if (total.value === 0) return 0;
	return Math.ceil(total.value / pageSize.value);
});

// 根据筛选条件过滤设备列表
const filteredDeviceList = computed(() => {
	let result = deviceList.value;
	
	if (filterType.value > 0) {
		const typeName = filterTypes[filterType.value].name;
		result = result.filter(item => item.type === typeName);
	}
	
	if (searchName.value.trim()) {
		const keyword = searchName.value.trim().toLowerCase();
		result = result.filter(item => 
			item.name.toLowerCase().includes(keyword)
		);
	}
	
	return result;
});

// 获取设备状态样式类
const getStatusClass = (status) => {
	switch (status) {
		case '在线':
			return 'online';
		case '离线':
			return 'offline';
		case '维修':
			return 'maintenance';
		default:
			return 'offline';
	}
};

// 返回上一页
const goBack = () => {
	uni.navigateBack();
};

// 处理类型筛选变化
const handleTypeChange = (item) => {
	filterType.value = item.index;
};

// 处理搜索
const handleSearch = () => {
};

// 处理清除搜索
const handleClearSearch = () => {
	searchName.value = '';
};

// 上一页
const handlePrevPage = () => {
	if (currentPage.value > 1) {
		currentPage.value--;
		fetchDeviceList();
	}
};

// 下一页
const handleNextPage = () => {
	if (currentPage.value < totalPages.value) {
		currentPage.value++;
		fetchDeviceList();
	}
};

// 跳转到指定页
const handleJumpPage = () => {
	const page = parseInt(jumpPage.value);
	if (isNaN(page) || page < 1 || page > totalPages.value) {
		uni.showToast({ title: `请输入1-${totalPages.value}之间的页码`, icon: 'none' });
		jumpPage.value = '';
		return;
	}
	
	if (page === currentPage.value) {
		jumpPage.value = '';
		return;
	}
	
	currentPage.value = page;
	jumpPage.value = '';
	fetchDeviceList();
};

// 添加设备的类型选择
const clickType = (type) => {
	selectedType.value = type.name;
};

// 显示编辑设备名称弹出层
const showEditSheetFunc = (item) => {
	currentDeviceId.value = item.id;
	currentDeviceName.value = item.name;
	editDeviceName.value = item.name;
	showEditSheet.value = true;
};

// 隐藏编辑弹出层并重置状态
const hideEditSheet = () => {
	showEditSheet.value = false;
	currentDeviceId.value = null;
	currentDeviceName.value = '';
	editDeviceName.value = '';
};

// 获取设备列表数据
const fetchDeviceList = () => {
	uni.showLoading({ title: '加载中' });
	
	const merch = uni.getStorageSync('merch') || {};
	const merchsId = merch.id || merch.merchsId || '';
	
	uni.$uv.http.get('/device/list', {
		params: { 
			merchs_id: merchsId,
			page: currentPage.value,
			page_size: pageSize.value
		},
		custom: { auth: true }
	}).then((res) => {
		deviceList.value = res.data.list || [];
		total.value = res.data.total || 0;
		uni.hideLoading();
		uni.stopPullDownRefresh();
	}).catch((err) => {
		deviceList.value = [];
		total.value = 0;
		uni.hideLoading();
		uni.stopPullDownRefresh();
		console.error('获取设备列表失败:', err);
	});
};

// 显示更换设备弹出层
const showActionSheet = (deviceId) => {
	currentDeviceId.value = deviceId;
	showSheet.value = true;
};

// 隐藏更换设备弹出层
const hideActionSheet = () => {
	showSheet.value = false;
	currentDeviceId.value = null;
};

// 显示添加设备弹出层
const AddSheet = () => {
	selectedType.value = '集控';
	deviceCode.value = '';
	showAddSheet.value = true;
};

// 隐藏添加设备弹出层并重置状态
const hideAddSheet = () => {
	showAddSheet.value = false;
	selectedType.value = '集控';
	deviceCode.value = '';
};

// 处理扫码更换设备
const handleScan = () => {
	hideActionSheet();
	uni.showToast({ title: '扫码功能开发中', icon: 'none' });
};

// 处理手动更换设备（跳转编辑页）
const handleEdit = () => {
	hideActionSheet();
	if (currentDeviceId.value) {
		uni.navigateTo({ url: `/pages/admin/device/edit?id=${currentDeviceId.value}` });
	}
};

// 扫码获取设备编码
const scanDeviceCode = () => {
	uni.scanCode({
		success: (res) => {
			deviceCode.value = res.result;
		},
		fail: (err) => {
			console.error('扫码失败:', err);
			uni.showToast({ title: '扫码失败', icon: 'none' });
		}
	});
};

// 处理添加设备
const handleAddDevice = () => {
	if (!deviceCode.value.trim()) {
		uni.showToast({ title: '请输入设备编码', icon: 'none' });
		return;
	}
	
	uni.showLoading({ title: '添加中' });
	
	const merch = uni.getStorageSync('merch') || {};
	const merchsId = merch.id || merch.merchsId || '';
	
	if (!merchsId) {
		uni.hideLoading();
		uni.showToast({ title: '请先登录', icon: 'none' });
		return;
	}
	
	uni.$uv.http.post('/device', {
		code: deviceCode.value.trim(),
		type: selectedType.value,
		merchsId: merchsId
	}, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		uni.showToast({ title: '添加成功', icon: 'success' });
		hideAddSheet();
		currentPage.value = 1;
		fetchDeviceList();
	}).catch((err) => {
		uni.hideLoading();
		console.error('添加设备失败:', err);
		uni.showToast({ title: '添加失败', icon: 'none' });
	});
};

// 处理更新设备名称
const handleUpdateDevice = () => {
	if (!editDeviceName.value.trim()) {
		uni.showToast({ title: '请输入设备名称', icon: 'none' });
		return;
	}
	
	if (editDeviceName.value === currentDeviceName.value) {
		uni.showToast({ title: '名称未改变', icon: 'none' });
		return;
	}
	
	uni.showLoading({ title: '更新中' });
	
	uni.$uv.http.put(`/device/${currentDeviceId.value}`, {
		name: editDeviceName.value.trim()
	}, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		uni.showToast({ title: '修改成功', icon: 'success' });
		hideEditSheet();
		fetchDeviceList();
	}).catch((err) => {
		uni.hideLoading();
		console.error('更新设备失败:', err);
		uni.showToast({ title: '修改失败', icon: 'none' });
	});
};

// 删除设备（带确认弹窗）
const deleteDevice = (id) => {
	uni.showModal({
		title: '确认删除',
		content: '确定要删除该设备吗？',
		success: function (res) {
			if (res.confirm) {
				uni.showLoading({ title: '删除中' });
				uni.$uv.http.delete(`/device/${id}`, {}, {
					custom: { auth: true }
				}).then(() => {
					uni.hideLoading();
					uni.showToast({ title: '删除成功', icon: 'success' });
					if (deviceList.value.length === 1 && currentPage.value > 1) {
						currentPage.value--;
					}
					fetchDeviceList();
				}).catch((err) => {
					uni.hideLoading();
					console.error('删除设备失败:', err);
					uni.showToast({ title: '删除失败', icon: 'none' });
				});
			}
		}
	});
};

// 页面加载时初始化数据
onLoad(() => {
	fetchDeviceList();
});

// 下拉刷新时重新获取数据
onPullDownRefresh(() => {
	currentPage.value = 1;
	fetchDeviceList();
});

// 格式化时间
const formatTime = (time) => {
	if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 140rpx;
}

.filter-section {
	background: #fff;
}

.search-section {
	padding: 16rpx 24rpx;
	background: #fff;
}

.list-wrapper {
	padding: 24rpx;
}

.device-list {
	display: flex;
	flex-direction: column;
}

.device-item {
	background: #fff;
	margin-bottom: 24rpx;
	border-radius: 20rpx;
	padding: 30rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
	display: flex;
	justify-content: space-between;
	align-items: center;
}

.device-left {
	flex: 1;
	
	.device-info {
		display: flex;
		flex-direction: column;
		gap: 12rpx;
	}
	
	.info-row {
		display: flex;
		align-items: center;
	}
	
	.info-label {
		font-size: 26rpx;
		color: #999;
		flex-shrink: 0;
	}
	
	.info-value {
			font-size: 28rpx;
			color: #333;
			font-weight: 500;
			
			&.status-text {
				&.online {
					color: #19be6b;
				}
				
				&.offline {
					color: #f56c6c;
				}
				
				&.maintenance {
					color: #f59e0b;
				}
			}
		}
	}

.device-right {
	flex-shrink: 0;
	margin-left: 30rpx;
	
	.action-buttons {
		display: flex;
		flex-direction: column;
		gap: 12rpx;
	}
	
	.action-btn {
		font-size: 26rpx;
		padding: 12rpx 24rpx;
		border-radius: 8rpx;
		text-align: center;
		white-space: nowrap;
		
		&.edit-btn {
			color: #3c9cff;
			background: rgba(60, 156, 255, 0.1);
		}
		
		&.change-btn {
			color: #67c23a;
			background: rgba(103, 194, 58, 0.1);
		}
		
		&.delete-btn {
			color: #f56c6c;
			background: rgba(245, 108, 108, 0.1);
		}
	}
}

.pagination-wrapper {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fff;
	box-shadow: 0 -4rpx 12rpx rgba(0, 0, 0, 0.05);
	padding: 20rpx 24rpx;
	padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
	z-index: 100;
}

.pagination {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 20rpx;
}

.page-btn {
	display: flex;
	align-items: center;
	gap: 8rpx;
	padding: 16rpx 24rpx;
	background: #f5f7fa;
	border-radius: 12rpx;
	font-size: 26rpx;
	color: #333;
	transition: all 0.2s;
	
	&:not(.disabled):active {
		background: #e8e8e8;
	}
	
	&.disabled {
		opacity: 0.4;
		color: #999;
	}
}

.page-info {
	display: flex;
	align-items: center;
	gap: 8rpx;
	font-size: 26rpx;
	color: #666;
	
	.current-page {
		color: #3c9cff;
		font-weight: 600;
		font-size: 28rpx;
	}
	
	.separator {
		color: #999;
	}
	
	.total-pages {
		color: #999;
	}
}

.page-jump {
	display: flex;
	align-items: center;
	gap: 12rpx;
	font-size: 26rpx;
	color: #666;
	
	.jump-input {
		width: 100rpx;
		height: 60rpx;
		padding: 0 12rpx;
		background: #f5f7fa;
		border-radius: 8rpx;
		font-size: 26rpx;
		color: #333;
		text-align: center;
	}
}

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

.action-sheet {
	width: 100%;
	background: #fff;
	border-radius: 24rpx 24rpx 0 0;
	padding-bottom: calc(40rpx + env(safe-area-inset-bottom));
	animation: slideUp 0.3s ease;
}

.add-device-sheet {
	width: 100%;
	background: #fff;
	border-radius: 24rpx 24rpx 0 0;
	padding-bottom: calc(40rpx + env(safe-area-inset-bottom));
	animation: slideUp 0.3s ease;
	max-height: 80vh;
	overflow-y: auto;
}

.edit-device-sheet {
	width: 100%;
	background: #fff;
	border-radius: 24rpx 24rpx 0 0;
	padding-bottom: calc(40rpx + env(safe-area-inset-bottom));
	animation: slideUp 0.3s ease;
}

@keyframes slideUp {
	from {
		transform: translateY(100%);
	}
	to {
		transform: translateY(0);
	}
}

.sheet-title {
	text-align: center;
	padding: 40rpx;
	font-size: 32rpx;
	font-weight: 600;
	color: #333;
	border-bottom: 1rpx solid #f0f0f0;
}

.sheet-btns {
	display: flex;
	padding: 40rpx 60rpx;
	gap: 80rpx;
	
	.sheet-btn {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16rpx;
		
		text {
			font-size: 28rpx;
			color: #333;
		}
	}
}

.sheet-cancel {
	text-align: center;
	padding: 32rpx;
	font-size: 30rpx;
	color: #999;
	background: #f5f5f5;
	margin-top: 16rpx;
}

.edit-form {
	padding: 30rpx;
	
	.form-item {
		margin-bottom: 20rpx;
	}
	
	.name-input {
		width: 100%;
		height: 88rpx;
		padding: 0 24rpx;
		background: #f8f9fa;
		border-radius: 12rpx;
		font-size: 28rpx;
		color: #333;
		box-sizing: border-box;
	}
	
	.input-placeholder {
		color: #999;
	}
}

.edit-actions {
	display: flex;
	padding: 0 30rpx 30rpx;
	gap: 24rpx;
	
	.edit-btn {
		flex: 1;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 12rpx;
		font-size: 30rpx;
		font-weight: 500;
		
		&.cancel {
			background: #f5f5f5;
			color: #666;
		}
		
		&.confirm {
			background: #3c9cff;
			color: #fff;
		}
	}
}

.type-section {
	padding: 30rpx;
	
	.section-label {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 20rpx;
		display: block;
	}
	
	.type-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 20rpx;
	}
	
	.type-item {
		width: calc(25% - 15rpx);
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 12rpx;
		padding: 24rpx 12rpx;
		border-radius: 12rpx;
		background: #f8f9fa;
		transition: all 0.2s;
		
		&.active {
			background: rgba(60, 156, 255, 0.1);
			
			text {
				color: #3c9cff;
				font-weight: 500;
			}
		}
		
		text {
			font-size: 24rpx;
			color: #666;
		}
	}
}

.code-section {
	padding: 0 30rpx 30rpx;
	
	.section-label {
		font-size: 28rpx;
		color: #666;
		margin-bottom: 20rpx;
		display: block;
	}
	
	.code-input-wrapper {
		display: flex;
		align-items: center;
		background: #f8f9fa;
		border-radius: 12rpx;
		padding-right: 12rpx;
	}
	
	.code-input {
		flex: 1;
		height: 88rpx;
		padding: 0 24rpx;
		font-size: 28rpx;
		color: #333;
		background: transparent;
	}
	
	.scan-btn {
		width: 80rpx;
		height: 80rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(60, 156, 255, 0.1);
		border-radius: 10rpx;
	}
}

.add-actions {
	display: flex;
	padding: 0 30rpx 30rpx;
	gap: 24rpx;
	
	.add-btn {
		flex: 1;
		height: 88rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		border-radius: 12rpx;
		font-size: 30rpx;
		font-weight: 500;
		
		&.cancel {
			background: #f5f5f5;
			color: #666;
		}
		
		&.confirm {
			background: #3c9cff;
			color: #fff;
		}
	}
}

.method-btns {
	display: flex;
	padding: 20rpx 40rpx 30rpx;
	gap: 40rpx;
	
	.method-btn {
		flex: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 16rpx;
		padding: 30rpx;
		background: rgba(60, 156, 255, 0.05);
		border-radius: 16rpx;
		
		text {
			font-size: 28rpx;
			color: #3c9cff;
			font-weight: 500;
		}
	}
}
</style>