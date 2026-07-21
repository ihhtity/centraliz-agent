<!-- 首页/房间管理 -->
<template>
	<view class="container">
		<uv-navbar title="房间管理" :placeholder="true" :leftIcon="''" />

		<!-- 搜索栏 -->
		<view class="search-bar">
			<uv-icon name="search" color="#999" size="24" />
			<input v-model="searchQuery" type="text" placeholder="搜索房间..." />
		</view>

		<!-- 功能按钮栏 -->
		<view class="action-bar">
			<view class="action-btn group-btn" @click="handleGroup">
				<uv-icon name="grid" color="#fff" size="20" />
				<text class="btn-text">分组</text>
			</view>
			<view class="action-btn" @click="handleGroup">
				<text class="btn-text">锁定: {{ lockCount }}</text>
			</view>
			<view class="action-btn add-btn" @click="handleAdd">
				<uv-icon name="plus-circle" color="#fff" size="20" />
				<text class="btn-text">添加</text>
			</view>
			<view class="action-btn refresh-btn" @click="handleRefresh">
				<uv-icon name="reload" color="#fff" size="20" />
				<text class="btn-text">刷新</text>
			</view>
		</view>

		<!-- 房间分组标签栏 -->
		<view style="margin-bottom: 24rpx;">
			<uv-tabs :list="roomGroups" @click="selectGroup" />
		</view>

		<!-- 房间列表 -->
		<view class="room-grid" v-if="filteredRooms.length > 0">
			<view v-for="(room, index) in filteredRooms" :key="index" class="room-card" :class="{
				'empty': room.status === '空闲',
				'rented': room.status === '租用',
				'maintenance': room.status === '维修'
			}" @click="goToRoomDetail(room)">
				<!-- 房间信息 -->
				<view class="room-info">
					<text class="room-number">{{ room.name }}</text>
					<text class="room-number">{{ room.boardNo }}-{{ room.lockNo }}</text>
					<text class="room-tag">{{ room.tag }}</text>
				</view>
			</view>
		</view>
		<view style="margin-top: 280rpx;" v-else>
			<uv-empty mode="data" textSize="32" iconSize="150" />
		</view>

		<!-- 底部统计 -->
		<view class="stats-section-fixed">
			<view class="stat-row">
				<view class="stat-item" :class="{ active: !filterStatus }" @click="setFilter('')">
					<view class="color-indicator all"></view>
					<text class="stat-label">全部</text>
					<text class="stat-value">{{ totalRooms }}</text>
				</view>
				<view class="stat-item" :class="{ active: filterStatus === '空闲' }" @click="setFilter('空闲')">
					<view class="color-indicator blue"></view>
					<text class="stat-label">空闲</text>
					<text class="stat-value">{{ stats.empty }}</text>
				</view>
				<view class="stat-item" :class="{ active: filterStatus === '租用' }" @click="setFilter('租用')">
					<view class="color-indicator red"></view>
					<text class="stat-label">租用</text>
					<text class="stat-value">{{ stats.rented }}</text>
				</view>
				<view class="stat-item" :class="{ active: filterStatus === '维修' }" @click="setFilter('维修')">
					<view class="color-indicator orange"></view>
					<text class="stat-label">维修</text>
					<text class="stat-value">{{ stats.maintenance }}</text>
				</view>
			</view>
		</view>

		<!-- 底部导航栏 -->
		<uv-tabbar :value="tabbar" @change="editTabbar">
			<uv-tabbar-item text="房间" icon="home" />
			<uv-tabbar-item text="个人" icon="account" />
		</uv-tabbar>

		<!-- 前往分组页面 -->
		<uv-modal ref="groupModalRef" title="前往分组" :show-cancel-button="true" cancel-text="取消" confirm-text="确定"
			@confirm="handleGroupConfirm" @cancel="groupModalRef.close()">
			<view class="delete-modal-content">
				<text class="delete-tip">暂无房间分组，是否前往分组管理页面添加房间分组？</text>
			</view>
		</uv-modal>

		<!-- 添加房间弹窗 -->
		<uv-popup ref="popupRef" mode="bottom" closeable @close="closeAddRoomModal" :safeAreaInsetBottom="true">
			<view class="add-room-modal">
				<view class="modal-header">
					<text class="modal-title">添加房间</text>
				</view>
				<view class="modal-content">
					<view class="form-item">
						<text class="form-label">房间名称 <text class="required">*</text></text>
						<view class="input-box">
							<uv-input v-model="addRoomForm.name" placeholder="请输入房间名称（如：A01）" class="input-field"
								border="none" :focus="formFocus.name" />
						</view>
					</view>

					<view class="form-item">
						<text class="form-label">房间数量 <text class="required">*</text></text>
						<view class="input-box">
							<uv-input v-model="addRoomForm.count" placeholder="请输入房间数量（1-50）" class="input-field"
								type="number" border="none" :focus="formFocus.count" @input="handleCountChange" />
						</view>
					</view>

					<text class="form-label">设备选择 <text style="color: #ff4d4f;">*</text></text>
					<view class="form-item clickable" @click="goToDevicePicker">
						<text class="form-label">设备选择</text>
						<view class="arrow-wrap">
							<text class="value-text">{{ addRoomForm.deviceName || '未选择' }}</text>
							<uv-icon name="arrow-right" color="#ccc" size="20" />
						</view>
					</view>
					

					<view class="form-item">
						<text class="form-label">房间板号 <text class="required">*</text></text>
						<view class="input-box">
							<uv-input v-model="addRoomForm.boardNo" placeholder="01" class="input-field"
								type="text" border="none" :disabled="true" />
						</view>
					</view>

					<view class="form-item">
						<text class="form-label">房间锁号</text>
						<view class="input-box">
							<uv-input v-model="addRoomForm.lockNo" placeholder="留空自动分配" class="input-field"
								type="number" border="none" :focus="formFocus.lockNo"
								:disabled="addRoomForm.count > 1" />
						</view>
					</view>

					<view class="form-item">
						<text class="form-label">房间标签</text>
						<view class="input-box">
							<uv-input v-model="addRoomForm.tag" placeholder="请输入房间标签（如：普通柜）" class="input-field"
								border="none" :focus="formFocus.tag" />
						</view>
					</view>
				</view>
				<view class="modal-footer">
					<view class="button button-cancel" @click="closeAddRoomModal">取消</view>
					<view class="button button-confirm" @click="submitAddRoom">确认添加</view>
				</view>
			</view>
		</uv-popup>

		<!-- 设备选择器 -->
		<uv-picker ref="devicespicker" title="选择设备" round="10" keyName="name" :columns="deviceslist"
			@confirm="onDeviceConfirm" @cancel="devicespicker.close()" />
	</view>
</template>

<script setup>
import { ref, computed, toRaw } from 'vue';
import { onShow } from '@dcloudio/uni-app';
import { useI18n } from 'vue-i18n';

onShow(() => {
	merch.value = uni.getStorageSync('merch') || {};
	// 初始化房间分组数据
	fetchGroupData();
});

// 国际化语言
const { t } = useI18n();
// 商家数据
const merch = ref({});
// 房间分组数据
const roomGroups = ref([]);
// 房间数据
const rooms = ref([]);
// 设备数据
const deviceslist = ref([]);
// 筛选状态
const filterStatus = ref('');
// 分组筛选ID
const groupSelect = ref(null);
// 全部房间总数计算
const totalRooms = ref(0);
// 分组设备锁定数量总和
const lockCount = ref(0);
// 搜索关键词
const searchQuery = ref('');
// 导航逻辑
const tabbar = ref(0);
// 前往分组弹窗
const groupModalRef = ref(null);
// 添加房间弹窗
const popupRef = ref(null);
// 设备选择器
const devicespicker = ref(null);
// 添加房间表单数据
const addRoomForm = ref({
	name: '1',
	count: '1',
	boardNo: '01',
	lockNo: '01',
	tag: '普通柜',
	groupsId: '',
	devicesId: '',
	deviceName: '未绑定',
	lockCount: 0,
	deviceCount: 0,
});
// 添加房间表单焦点状态
const formFocus = ref({
	name: false,
	count: false,
	boardNo: false,
	lockNo: false,
	tag: false
});

// 获取分组数据
const fetchGroupData = () => {
	uni.showLoading({ title: '加载中...', duration: 13000 });
	const params = {};
	if (merch.value && merch.value.id) {
		params.merchs_id = merch.value.id;
	}

	uni.$uv.http.get('/group/list', {
		params: params,
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading();
		if (res.data && res.data.list && res.data.list.length > 0) {
			const groups = [];
			res.data.list.forEach(group => {
				groups.push({
					id: group.id.toString(),
					name: group.name || `分组${group.id}`
				});
			});
			roomGroups.value = groups;

			if (groups.length > 0) {
				if (groupSelect.value) {
					fetchRoomDataByGroup(groupSelect.value.id);
					return
				}

				groupSelect.value = groups[0];
				fetchRoomDataByGroup(groupSelect.value.id);
			} else {
				groupSelect.value = {};
				throw new Error('暂无分组数据');
			}
		} else {
			roomGroups.value = [];
			rooms.value = [];
			groupSelect.value = {};
			throw new Error('暂无分组数据');
		}
	}).catch((err) => {
		uni.hideLoading();
		uni.showToast({
			title: err.message || '网络异常，请检查网络连接',
			icon: 'none',
			duration: 3000,
		});
		groupSelect.value = {};
	});
};
// 根据分组ID获取房间数据
const fetchRoomDataByGroup = (groupId) => {
	uni.showLoading({ title: '加载中...', duration: 3000 });

	const params = {};
	// 只上传 groups_id
	if (groupId) {
		params.groups_id = groupId;
	}

	uni.$uv.http.get('/room/list', {
		params: params,
		custom: { auth: true }
	}).then((res) => {
		if (res.code !== 200) {
			throw new Error('暂无房间数据');
		}

		rooms.value = res.data.list;
		totalRooms.value = res.data.total;
		lockCount.value = res.data.lockCount;
		deviceslist.value = [toRaw(res.data.devices)] || [];
		uni.hideLoading();
	}).catch((err) => {
		uni.hideLoading();
		uni.showToast({
			title: err.message || '网络异常，请检查网络连接',
			icon: 'none',
			duration: 3000,
		});
	});
};
// 提交添加房间
const submitAddRoom = async () => {
	if (!addRoomForm.value.name.trim()) {
		uni.showToast({ title: '请输入房间名称', icon: 'none' });
		return;
	}

	// 验证必须选择设备
	if (!addRoomForm.value.devicesId) {
		uni.showToast({ title: '请选择设备', icon: 'none' });
		return;
	}

	const count = parseInt(addRoomForm.value.count);
	if (isNaN(count) || count < 1 || count > 50) {
		uni.showToast({ title: '房间数量必须在1-50之间', icon: 'none' });
		return;
	}

	const boardNo = parseInt(addRoomForm.value.boardNo);
	if (isNaN(boardNo) || boardNo < 1 || boardNo > 99) {
		uni.showToast({ title: '房间板号必须在1-99之间', icon: 'none' });
		return;
	}

	let lockNo = parseInt(addRoomForm.value.lockNo);
	if (isNaN(lockNo)) {
		lockNo = 0; // 0表示让后端自动分配
	}
	if (lockNo < 0 || lockNo > 999) {
		uni.showToast({ title: '房间锁号必须在0-999之间', icon: 'none' });
		return;
	}

	if (!merch.value.id) {
		uni.showToast({ title: '商家信息不存在', icon: 'none' });
		return;
	}

	// 验证设备是否已绑定房间数量已达上限
	if (addRoomForm.value.deviceCount + count >= addRoomForm.value.lockCount) {
		uni.showToast({ title: '该设备绑定房间数量已达上限', icon: 'none', duration: 3000 });
		return;
	}

	const requestData = {
		name: addRoomForm.value.name,
		count: count,
		board_no: boardNo,
		lock_no: formatLockNo(lockNo),
		merchs_id: merch.value.id,
		tag: addRoomForm.value.tag || '普通柜',
		device_id: parseInt(addRoomForm.value.devicesId),
	};

	if (addRoomForm.value.groupsId && addRoomForm.value.groupsId !== '0') {
		requestData.groups_id = parseInt(addRoomForm.value.groupsId);
	}

	uni.showLoading({ title: '添加中...' });

	try {
		const res = await uni.$uv.http.post('/room', requestData, {
			custom: { auth: true }
		});

		if (res.code === 200) {
			// 成功添加的房间数量
			const successCount = res.data?.count || 0;
			// 已存在的房间数量
			const skipCount = res.data?.skipCount || 0;
			let msg = `成功添加${successCount}个房间`;
			if (skipCount > 0) {
				msg += `，跳过${skipCount}个已存在的房间`;
			}
			uni.showToast({ title: msg, icon: 'success' });
			closeAddRoomModal();
			fetchRoomDataByGroup(groupSelect.value.id);
		} else {
			uni.showToast({ title: res.msg || '添加失败', icon: 'none' });
		}
	} catch (err) {
		uni.showToast({ title: err.msg || '添加失败', icon: 'none' });
	} finally {
		uni.hideLoading();
	}
};
// 添加房间
const handleAdd = () => {
	if (!groupSelect.value) {
		groupModalRef.value.open();
		return;
	}

	// 默认设置当前选择的分组
	addRoomForm.value.groupsId = groupSelect.value.id;
	addRoomForm.value.devicesId = 0
	addRoomForm.value.deviceName = '未绑定'
	popupRef.value.open()
};
// 打开设备选择器
const goToDevicePicker = () => {
	if (deviceslist.value.length === 0) {
		uni.showToast({ title: '暂无可用设备', icon: 'none' });
		return;
	}
	devicespicker.value.open();
};
// 设备选择确认
const onDeviceConfirm = (e) => {
	const selectedValue = e.value[0];
	addRoomForm.value.deviceName = selectedValue.name;
	addRoomForm.value.devicesId = selectedValue.id.toString();
	addRoomForm.value.lockCount = selectedValue.lockCount;
	addRoomForm.value.deviceCount = rooms.value.filter(r => r.devicesId === selectedValue.id.toString()).length;
	// 设置板号为设备的板号
	if (selectedValue.boardNo) {
		addRoomForm.value.boardNo = selectedValue.boardNo;
	} else {
		addRoomForm.value.boardNo = '01';
	}
	devicespicker.value.close();
};
// 格式化锁号为2-3位数
const formatLockNo = (lockNo) => {
  const num = parseInt(lockNo) || 0
  if (num < 0) return 0
  if (num > 999) return 999
  return num
};
// 关闭添加房间弹窗
const closeAddRoomModal = () => {
	popupRef.value.close();
};
// 选择房间分组
const selectGroup = (e) => {
	groupSelect.value = e;
	fetchRoomDataByGroup(e.id);
};
// 前往分组确认
const handleGroupConfirm = () => {
	uni.navigateTo({
		url: "/pages/admin/group/manage"
	});
	groupModalRef.value.close();
};
// 跳转到分组管理页面
const handleGroup = () => {
	uni.navigateTo({
		url: `/pages/admin/group/manage?id=${merch.value.id}`
	});
};
// 刷新房间列表
const handleRefresh = () => {
	filterStatus.value = '';
	fetchRoomDataByGroup(groupSelect.value.id);
};
// 房间数量变化处理
const handleCountChange = () => {
	const count = parseInt(addRoomForm.value.count) || 1;
	if (count > 1) {
		addRoomForm.value.lockNo = '0';
	}
};
// 过滤房间逻辑
const filteredRooms = computed(() => {
	let result = rooms.value;

	if (searchQuery.value) {
		result = result.filter(room =>
			room.id.toString().toLowerCase().includes(searchQuery.value.toLowerCase()) ||
			(room.tag && room.tag.includes(searchQuery.value)) ||
			room.status.includes(searchQuery.value) ||
			(room.name && room.name.toLowerCase().includes(searchQuery.value.toLowerCase()))
		);
	}

	if (filterStatus.value) {
		result = result.filter(room => room.status === filterStatus.value);
	}

	return result;
});
// 设置筛选状态
const setFilter = (status) => {
	if (filterStatus.value === status) {
		filterStatus.value = '';
	} else {
		filterStatus.value = status;
	}
};
// 统计状态
const stats = computed(() => {
	const empty = rooms.value.filter(r => r.status === '空闲').length;
	const rented = rooms.value.filter(r => r.status === '租用').length;
	const maintenance = rooms.value.filter(r => r.status === '维修').length;
	return { empty, rented, maintenance };
});
// 点击 tab 切换页面
const editTabbar = (e) => {
	tabbar.value = e;
	if (e === 1) {
		uni.reLaunch({
			url: '/pages/admin/profile/index'
		});
	}
};
// 跳转房间详情
const goToRoomDetail = (room) => {
	uni.navigateTo({
		url: `/pages/admin/room/detail?id=${room.id}&name=${encodeURIComponent(room.name)}`
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding: 24rpx;
	padding-bottom: 80rpx;
}

.search-bar {
	position: relative;
	width: 92%;
	height: 80rpx;
	background: #fff;
	border-radius: 40rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	overflow: hidden;
	display: flex;
	align-items: center;
	padding: 0 20rpx;
	margin-bottom: 24rpx;

	input {
		flex: 1;
		font-size: 28rpx;
		color: #333;
		outline: none;
		background: transparent;
	}

	.uv-icon {
		margin-right: 16rpx;
	}
}

.action-bar {
	display: flex;
	gap: 16rpx;

	.action-btn {
		flex: 1;
		display: flex;
		flex-direction: row;
		justify-content: center;
		padding: 12rpx 16rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		color: #fff;
		background: #5b626f;
		text-align: center;
		cursor: pointer;

		.uv-icon {
			margin-bottom: 4rpx;
		}

		.btn-text {
			margin-left: 8rpx;
			font-size: 24rpx;
			font-weight: 500;
		}
	}

	.group-btn {
		background: linear-gradient(135deg, #4CAF50, #388E3C);
	}

	.add-btn {
		background: linear-gradient(135deg, #F44336, #D32F2F);
	}

	.refresh-btn {
		background: linear-gradient(135deg, #FF9800, #E67C00);
	}
}

.room-grid {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 20rpx;
	margin-bottom: 30rpx;
}

.room-card {
	background: #ffffff;
	border-radius: 16rpx;
	padding: 20rpx 16rpx;
	box-shadow: 0 6rpx 16rpx rgba(0, 0, 0, 0.08);
	position: relative;
	min-height: 140rpx;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	transition: all 0.2s ease;

	&:active {
		transform: translateY(4rpx);
		box-shadow: 0 8rpx 20rpx rgba(0, 0, 0, 0.12);
	}

	&.empty {
		background: #3c9cff
	}

	&.rented {
		background: #fa3534;
	}

	&.maintenance {
		background: #ff9900;
	}

	.status-badge {
		position: absolute;
		top: -10rpx;
		right: -10rpx;
		width: 44rpx;
		height: 44rpx;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 24rpx;
		font-weight: bold;
		box-shadow: 0 4rpx 10rpx rgba(0, 0, 0, 0.2);
		z-index: 2;

		&.empty {
			background: #3c9cff;
			color: white;
		}

		&.rented {
			background: #fa3534;
			color: white;
		}

		&.maintenance {
			background: #ff9900;
			color: white;
		}
	}

	.room-info {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 8rpx;
		text-align: center;
	}

	.room-header {
		display: flex;
		align-items: center;
		gap: 8rpx;
	}

	.room-number {
		font-size: 28rpx;
		font-weight: 700;
		color: #fff;
		line-height: 1.2;
	}

	.device-badge {
		width: 32rpx;
		height: 32rpx;
		background: rgba(255, 255, 255, 0.2);
		border-radius: 6rpx;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.room-tag {
		font-size: 20rpx;
		color: #fff;
		padding: 2rpx 12rpx;
	}

	.device-count {
		font-size: 18rpx;
		color: rgba(255, 255, 255, 0.8);
	}
}

.stats-section-fixed {
	position: fixed;
	bottom: calc(100rpx + env(safe-area-inset-bottom));
	left: 0;
	right: 0;
	padding: 20rpx 0rpx;
	border-radius: 20rpx;
	background: #fff;
	box-shadow: 0 -2rpx 16rpx rgba(0, 0, 0, 0.1);
	z-index: 99;
	width: 94%;
	margin: 0 auto;

	.stat-row {
		display: flex;
		justify-content: space-around;
		align-items: center;
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
				background: #5ac725;
			}

			&.blue {
				background: linear-gradient(135deg, #3c9cff, #2b85e4);
			}

			&.red {
				background: linear-gradient(135deg, #fa3534, #e63332);
			}

			&.orange {
				background: linear-gradient(135deg, #ff9900, #f29100);
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

.add-room-modal {
	padding: 0;
	background: #fafafa;
	border-radius: 24rpx 24rpx 0 0;
	height: auto;
	display: flex;
	flex-direction: column;
	overflow: hidden;

	.modal-header {
		background: #fff;
		padding: 32rpx;
		border-bottom: 1rpx solid #f0f0f0;
		text-align: center;

		.modal-title {
			font-size: 32rpx;
			font-weight: 600;
			color: #1a1a1a;
		}
	}

	.modal-content {
		padding: 32rpx;
		display: flex;
		flex-direction: column;
		gap: 32rpx;
	}

	.form-item {
		display: flex;
		flex-direction: column;
		gap: 16rpx;

		.form-label {
			font-size: 28rpx;
			color: #333;
			font-weight: 500;

			.required {
				color: #ff4d4f;
			}
		}

		.input-box {
			background: #fff;
			border-radius: 16rpx;
			padding: 0 24rpx;
			height: 88rpx;
			display: flex;
			align-items: center;
			border: 1rpx solid #e8e8e8;
			transition: all 0.2s ease;

			&:focus-within {
				border-color: #3c9cff;
				box-shadow: 0 0 0 4rpx rgba(60, 156, 255, 0.08);
			}
		}

		.input-field {
			flex: 1;
			font-size: 28rpx;
		}

		.group-box {
			background: #f8f9fa;
			border: 1rpx solid #e8e8e8;

			.group-name {
				font-size: 28rpx;
				color: #3c9cff;
				font-weight: 500;
			}
		}

		&.clickable {
			flex-direction: row;
			align-items: center;
			background: #fff;
			border-radius: 16rpx;
			padding: 0 24rpx;
			height: 88rpx;
			border: 1rpx solid #e8e8e8;
			transition: all 0.2s ease;

			&:active {
				background: #f5f7fa;
				transform: scale(0.98);
			}

			.form-label {
				flex: 0 0 auto;
				margin-right: 16rpx;
			}

			.arrow-wrap {
				flex: 1;
				display: flex;
				flex-direction: row;
				justify-content: flex-end;
				align-items: center;

				.value-text {
					font-size: 28rpx;
					color: #666;

					&::placeholder {
						color: #bbb;
					}
				}

				.uv-icon {
					font-size: 24rpx;
					color: #ccc;
					transition: transform 0.2s ease;
				}
			}
		}
	}

	.modal-footer {
		background: #fff;
		padding: 24rpx 32rpx;
		padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
		display: flex;
		gap: 16rpx;
		border-top: 1rpx solid #f0f0f0;

		.button {
			flex: 1;
			height: 88rpx;
			border-radius: 16rpx;
			display: flex;
			align-items: center;
			justify-content: center;
			font-size: 30rpx;
			font-weight: 500;
			transition: all 0.2s ease;

			&:active {
				transform: scale(0.98);
				opacity: 0.9;
			}
		}

		.button-cancel {
			background: #f5f5f5;
			color: #666;
		}

		.button-confirm {
			background: #3c9cff;
			color: #fff;
		}
	}
}
</style>
