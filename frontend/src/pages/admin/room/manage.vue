<template>
	<view class="container">
		<uv-navbar :title="t('admin.room.title')" :placeholder="true" rightIcon="plus" @leftClick="goBack"
			@rightClick="$refs.modal.open();addRoomType = false" />

		<!-- 新增: 搜索栏 -->
		<view class="search-bar fade-in-up">
			<uv-search v-model="searchKeyword" :placeholder="t('common.search') + t('admin.room.roomName')"
				shape="round" :bgColor="'#f0f2f5'" @search="handleSearch" @clear="handleSearchClear" />
		</view>

		<!-- 新增: 房间分组列表 -->
		<view class="room-groups">
			<uv-tabs :list="roomGroups" lineColor="#3c9cff" lineWidth="40" lineHeight="6" @click="onTabChange" />
		</view>

		<!-- 房间列表容器 -->
		<view class="content">
			<!-- 房间列表 -->
			<view class="locker-grid">
				<view v-for="(item, index) in filteredLockerList" :key="item.id" class="locker-item fade-in-up"
					:style="{ animationDelay: (index * 0.03) + 's' }" :class="{
						'is-occupied': item.status === 1,
						'is-free': item.status === 0,
						'is-maintenance': item.status === 2
					}" @click="handleLockerClick(item)">

					<!-- 状态角标 -->
					<view class="status-corner bounce-in">
						<text>{{ getStatusText(item) }}</text>
					</view>

					<!-- 修改: 房间主体信息显示 -->
					<view class="locker-main-info">
						<view class="info-row room-header">
							<text class="room-no">{{ item.no }}</text>
						</view>

						<!-- 锁信息区域 -->
						<view class="info-section lock-section">
							<view class="section-title">
								<uv-icon name="lock" size="14" color="#666" />
								<text>锁信息</text>
							</view>
							<view class="section-content">
								<text class="status-tag"
									:class="{ 'tag-error': !item.lockInfo?.isBound, 'tag-warning': item.lockInfo?.battery < 20 }">
									{{ item.lockInfo?.isBound ? (item.lockInfo.battery < 20 ? '低电' : '正常') : '未绑' }}
										</text>
										<!-- 修改: 使用 BatteryIndicator 组件显示电量 -->
										<BatteryIndicator v-if="item.lockInfo?.isBound"
											:percentage="item.lockInfo.battery" :size="24" :showText="false" />
										<text class="detail-text"
											:class="{ 'text-red': !item.lockInfo?.isGatewayBound }">
											网关: {{ item.lockInfo?.isGatewayBound ? '已绑' : '未绑' }}
										</text>
							</view>
						</view>

						<!-- 控电设备信息区域 -->
						<view class="info-section power-section">
							<view class="section-title">
								<uv-icon name="empty-favor" size="14" color="#666" />
								<text>控电</text>
							</view>
							<view class="section-content">
								<text class="status-tag" :class="{ 'tag-error': !item.powerControlInfo?.isBound }">
									{{ item.powerControlInfo?.isBound ? '已绑' : '未绑' }}
								</text>
								<text class="detail-text" v-if="item.powerControlInfo?.isBound">
									状态:
									<text
										:class="{ 'text-green': item.powerControlInfo.isOnline, 'text-red': !item.powerControlInfo.isOnline }">
										{{ item.powerControlInfo.isOnline ? '在线' : '离线' }}
									</text>
								</text>
								<text class="detail-text"
									v-if="item.powerControlInfo?.isBound && item.powerControlInfo.isOnline">
									开关: {{ item.powerControlInfo.switchStatus ? '开' : '关' }}
								</text>
							</view>
						</view>
					</view>
				</view>

				<!-- 新增: 无数据提示 -->
				<view v-if="filteredLockerList.length === 0" class="no-data-tip">
					<uv-empty mode="list" :text="t('common.noData')"></uv-empty>
				</view>

				<!-- 新增: 弹窗组件 -->
				<RoomDetailPopup v-if="isPopupVisible" :locker="selectedLocker" @close="closePopup" />
			</view>
		</view>

		<!-- 新增: 底部统计信息 (固定在底部) -->
		<view class="stats-fixed-wrapper">
			<view class="stats-container">
				<view class="stats-item">
					<view class="stats-icon-box free-bg">
						<uv-icon name="checkmark-circle" size="18" color="#1890ff" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.available') }}</text>
						<text class="stats-num num-free">{{ stats.free }}</text>
					</view>
				</view>
				<view class="divider"></view>
				<view class="stats-item">
					<view class="stats-icon-box used-bg">
						<uv-icon name="close-circle" size="18" color="#ff4d4f" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.occupied') }}</text>
						<text class="stats-num num-used">{{ stats.used }}</text>
					</view>
				</view>
				<view class="divider"></view>
				<view class="stats-item">
					<view class="stats-icon-box maint-bg">
						<uv-icon name="setting" size="18" color="#faad14" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.maintenance') }}</text>
						<text class="stats-num num-maintenance">{{ stats.maintenance }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 新增: 返回顶部按钮 -->
		<uv-back-top :scroll-top="scrollTop" :top="100" :bottom="120" :right="30" icon="arrow-upward" />

		<!-- 弹框 -->
		<view>
			<!-- 添加房间弹框 -->
			<uv-modal ref="modal" :title="t('admin.room.addRoom')" :showCancelButton="true" @confirm="addRoom" @cancel="$refs.modal.close()">
				<view class="add-room">
					<view class="room-groups">
						<uv-tabs :list="batchOperations" @click="addRoomType = !addRoomType" />
					</view>
					<view class="form-item" v-if="addRoomType">
						<text class="form-label">{{ t('admin.room.roomNum') }}：</text>
						<uv-input v-model="formData.sum" type="number" :placeholder="t('admin.room.roomNumPlaceholder')" />
					</view>
					<view class="form-item">
						<text class="form-label">{{ t('admin.room.roomName') }}：</text>
						<uv-input v-model="formData.name" :placeholder="t('admin.room.roomNamePlaceholder')" />
					</view>
				</view>
			</uv-modal>
		</view>
	</view>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { onPageScroll } from '@dcloudio/uni-app'
import BatteryIndicator from '@/components/BatteryIndicator.vue'

// 新增: 弹窗组件
import RoomDetailPopup from '@/components/RoomDetailPopup.vue'

const { t } = useI18n()

// status: 0-空闲, 1-租用, 2-维修
const lockerList = ref([
	{
		id: 1,
		no: 'A01',
		status: 0, // 空闲
		lockInfo: {
			isBound: true,
			battery: 85,
			isGatewayBound: true
		},
		powerControlInfo: {
			isBound: true,
			isOnline: true,
			switchStatus: false // false: 断电, true: 通电
		}
	},
	{
		id: 2,
		no: 'A02',
		status: 1, // 租用
		lockInfo: {
			isBound: true,
			battery: 42,
			isGatewayBound: true
		},
		powerControlInfo: {
			isBound: true,
			isOnline: true,
			switchStatus: true
		}
	},
	{
		id: 3,
		no: 'A03',
		status: 1, // 租用
		lockInfo: {
			isBound: true,
			battery: 12, // 低电量
			isGatewayBound: false // 未绑定网关
		},
		powerControlInfo: {
			isBound: true,
			isOnline: false, // 离线
			switchStatus: false
		}
	},
	{
		id: 4,
		no: 'A04',
		status: 0, // 空闲
		lockInfo: {
			isBound: false, // 未绑锁
			battery: 0,
			isGatewayBound: false
		},
		powerControlInfo: {
			isBound: false, // 未绑控电
			isOnline: false,
			switchStatus: false
		}
	},
	{
		id: 5,
		no: 'B01',
		status: 2, // 维修
		lockInfo: {
			isBound: true,
			battery: 90,
			isGatewayBound: true
		},
		powerControlInfo: {
			isBound: true,
			isOnline: true,
			switchStatus: false
		}
	},
	// 新增: 随机添加的10个房间数据
	{
		id: 6,
		no: 'B02',
		status: 0,
		lockInfo: { isBound: true, battery: 95, isGatewayBound: true },
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: false }
	},
	{
		id: 7,
		no: 'B03',
		status: 1,
		lockInfo: { isBound: true, battery: 60, isGatewayBound: true },
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: true }
	},
	{
		id: 8,
		no: 'C01',
		status: 0,
		lockInfo: { isBound: true, battery: 15, isGatewayBound: true }, // 低电
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: false }
	},
	{
		id: 9,
		no: 'C02',
		status: 1,
		lockInfo: { isBound: true, battery: 88, isGatewayBound: false }, // 未绑网关
		powerControlInfo: { isBound: true, isOnline: false, switchStatus: false } // 离线
	},
	{
		id: 10,
		no: 'C03',
		status: 2,
		lockInfo: { isBound: true, battery: 50, isGatewayBound: true },
		powerControlInfo: { isBound: false, isOnline: false, switchStatus: false } // 未绑控电
	},
	{
		id: 11,
		no: 'D01',
		status: 0,
		lockInfo: { isBound: false, battery: 0, isGatewayBound: false }, // 未绑锁
		powerControlInfo: { isBound: false, isOnline: false, switchStatus: false }
	},
	{
		id: 12,
		no: 'D02',
		status: 1,
		lockInfo: { isBound: true, battery: 30, isGatewayBound: true },
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: true }
	},
	{
		id: 13,
		no: 'D03',
		status: 0,
		lockInfo: { isBound: true, battery: 100, isGatewayBound: true },
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: false }
	},
	{
		id: 14,
		no: 'E01',
		status: 1,
		lockInfo: { isBound: true, battery: 5, isGatewayBound: true }, // 极低电
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: true }
	},
	{
		id: 15,
		no: 'E02',
		status: 2,
		lockInfo: { isBound: true, battery: 75, isGatewayBound: true },
		powerControlInfo: { isBound: true, isOnline: true, switchStatus: false }
	},
])
// 新增: 锁房分组列表数据
const roomGroups = ref([
	{ id: 'all', name: '全部' },
	{ id: 'floor1', name: '1楼' },
	{ id: 'floor2', name: '2楼' },
	{ id: 'floor3', name: '3楼' },
	{ id: 'floor4', name: '4楼' },
	{ id: 'floor5', name: '5楼' },
	{ id: 'floor6', name: '6楼' },
	{ id: 'floor7', name: '7楼' },
	{ id: 'floor8', name: '8楼' },
	{ id: 'floor9', name: '9楼' },
	{ id: 'floor10', name: '10楼' }
])
// 新增: 批量
const batchOperations = ref([
	{ name: '单个' },
	{ name: '批量' },
])
// 0-单个, 1-批量
const addRoomType = ref(false)
// 新增: 搜索关键词
const searchKeyword = ref('')
// 新增: 滚动距离
const scrollTop = ref(0)
// 新增: 表单数据
const formData = ref({
	sum: 1,
	name: ''
})

const stats = computed(() => {
	const res = {
		free: 0,
		used: 0,
		maintenance: 0
	}

	lockerList.value.forEach(item => {
		if (item.status === 0) res.free++
		else if (item.status === 1) res.used++
		else if (item.status === 2) res.maintenance++
	})

	return res
})

const addRoom = () => {
	if (!formData.value.name) {
		uni.showToast({
			title: '请输入房间名称前缀',
			icon: 'none'
		})
		return
	}

	// 获取当前最大ID，用于生成新ID
	let currentMaxId = lockerList.value.length > 0
		? Math.max(...lockerList.value.map(item => item.id))
		: 0

	const roomsToAdd = []
	
	// 确定生成数量
	const count = addRoomType.value ? (parseInt(formData.value.sum) || 1) : 1

	if (addRoomType.value && (!formData.value.sum || formData.value.sum <= 0)) {
		uni.showToast({
			title: '请输入有效的房间数量',
			icon: 'none'
		})
		return
	}

	for (let i = 0; i < count; i++) {
		currentMaxId++
		
		// 生成房号：如果是批量，建议格式为 "名称+序号"，如果是单个，直接使用名称
		// 这里为了简单和通用，批量时追加序号，单个时直接使用输入的名称
		let roomNo = formData.value.name
		if (addRoomType.value) {
			// 假设用户输入的是前缀，如 "A"，生成 "A1", "A2"...
			// 或者如果用户输入的是完整名字，这里可以根据需求调整。
			// 此处采用：前缀 + 索引 (从1开始)
			roomNo = `${formData.value.name}${i + 1}`
		}

		const newRoom = {
			id: currentMaxId,
			no: roomNo,
			status: 0, // 空闲
			lockInfo: {
				isBound: false,
				battery: 0,
				isGatewayBound: false
			},
			powerControlInfo: {
				isBound: false,
				isOnline: false,
				switchStatus: false
			}
		}
		roomsToAdd.push(newRoom)
	}

	// 批量插入数组
	lockerList.value.push(...roomsToAdd)

	// 清空表单
	formData.value.name = ''
	if (addRoomType.value) {
		formData.value.sum = 1
	}

	uni.showToast({
		title: `成功添加 ${roomsToAdd.length} 个房间`,
		icon: 'success'
	})
}

// 新增: 获取房间状态文字
const onTabChange = (tab) => {
	// console.log('Selected tab:', tab)
}

// 新增: 过滤后的房间列表
const filteredLockerList = computed(() => {
	if (!searchKeyword.value) {
		return lockerList.value
	}
	const keyword = searchKeyword.value.toLowerCase()
	return lockerList.value.filter(item => {
		return item.no.toLowerCase().includes(keyword) ||
			getStatusText(item).includes(keyword)
	})
})

const goBack = () => {
	uni.navigateBack()
}


const getStatusText = (item) => {
	if (item.status === 2) return t('user.locker.statusMaintenance')
	if (item.status === 1) return t('user.locker.statusRent')
	return t('user.locker.statusFree')
}

// 新增: 控制弹窗显示状态
const isPopupVisible = ref(false);
const selectedLocker = ref(null);

// 新增: 显示弹窗的方法
const showPopup = (item) => {
	selectedLocker.value = item;
	isPopupVisible.value = true;
};

// 新增: 关闭弹窗的方法
const closePopup = () => {
	isPopupVisible.value = false;
	selectedLocker.value = null;
};

// 修改: 点击房间事件处理
const handleLockerClick = (item) => {
	// 显示弹窗
	showPopup(item);
};

// 新增: 搜索处理
const handleSearch = () => {
}

const handleSearchClear = () => {
	searchKeyword.value = ''
}

// 新增: 使用 onPageScroll 钩子监听页面滚动
onPageScroll((e) => {
	scrollTop.value = e.scrollTop
})

onUnmounted(() => {
	// 移除监听
	// 由于使用了 onPageScroll 钩子，通常不需要手动移除 window 监听
	// 如果之前有手动绑定 window.scroll，则在此移除
})

</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	/* 优化: 使用更中性的浅灰背景，减少视觉疲劳 */
	padding-bottom: 180rpx;
	/* 增加底部留白，防止内容被底部统计栏遮挡 */
	box-sizing: border-box;
	/* 新增: 明确指定系统字体栈，避免浏览器尝试加载未定义的自定义字体，从而减少 Slow network 警告 */
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
}

// 新增: 搜索栏样式
.search-bar {
	padding: 20rpx 24rpx;
	background-color: #fff;
	margin-bottom: 16rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);
	/* 优化: 添加轻微阴影 */
}

.add-room {
	display: flex;
	flex-direction: column;
	align-items: center;
}
.form-item {
	display: flex;
	flex-direction: row;
	align-items: center;
}

.content {
	padding: 0 24rpx 24rpx;
}

.room-groups {
	padding: 0 24rpx;
	margin-bottom: 24rpx;

	/* 优化: 给 Tab 栏添加容器样式，使其更像一个独立的模块 */
	::v-deep .uv-tabs__wrapper {
		background-color: #fff;
		border-radius: 16rpx;
		padding: 4rpx;
		box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.03);
	}
}

.locker-grid {
	display: flex;
	flex-wrap: wrap;
	justify-content: space-between;
	gap: 24rpx;
	/* 优化: 增加间距，使布局更呼吸感 */
}

.locker-item {
	width: calc(50% - 12rpx);
	/* 适配新的 gap */
	min-height: 340rpx;
	/* 稍稍微增加高度以容纳更舒适的排版 */
	background-color: #ffffff;
	border-radius: 24rpx;
	/* 优化: 更大的圆角 */
	position: relative;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	padding: 28rpx;
	box-sizing: border-box;
	transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.04);
	border: 1rpx solid rgba(0, 0, 0, 0.03);
	/* 优化: 极细的边框增加精致感 */
	overflow: hidden;
	will-change: transform, opacity;

	/* 优化: 增加点击态反馈 */
	&:active {
		transform: scale(0.98);
	}

	/* 空闲状态 - 清新蓝 */
	&.is-free {
		background: linear-gradient(135deg, #f0f9ff 0%, #ffffff 100%);
		border-color: #bae7ff;

		.room-no {
			color: #1890ff;
		}
	}

	/* 被占用状态 - 柔和红 */
	&.is-occupied {
		background: linear-gradient(135deg, #fff1f0 0%, #ffffff 100%);
		border-color: #ffccc7;

		.room-no {
			color: #ff4d4f;
		}
	}

	/* 维修状态 - 警示黄 */
	&.is-maintenance {
		background: linear-gradient(135deg, #fffbe6 0%, #ffffff 100%);
		border-color: #ffe58f;

		.room-no {
			color: #faad14;
		}
	}
}

// 新增: 房间主要信息布局
.locker-main-info {
	display: flex;
	flex-direction: column;
	flex: 1;
	margin-top: 16rpx;
}

.room-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20rpx;
	padding-bottom: 16rpx;
	border-bottom: 1rpx dashed rgba(0, 0, 0, 0.05);
	/* 优化: 添加分隔线 */

	.room-no {
		font-size: 40rpx;
		/* 优化: 稍大字体 */
		font-weight: 700;
		color: #333;
		letter-spacing: 1rpx;
	}
}

// 新增: 信息区块样式
.info-section {
	margin-bottom: 16rpx;

	.section-title {
		display: flex;
		align-items: center;
		font-size: 24rpx;
		color: #8c8c8c;
		margin-bottom: 10rpx;
		font-weight: 500;

		text {
			margin-left: 8rpx;
		}
	}

	.section-content {
		display: flex;
		flex-wrap: wrap;
		gap: 12rpx;
		align-items: center;
	}

	.status-tag {
		font-size: 22rpx;
		padding: 4rpx 10rpx;
		border-radius: 8rpx;
		background-color: #f5f5f5;
		color: #666;
		font-weight: 500;
		transition: all 0.2s;

		&.tag-error {
			background-color: #fff1f0;
			color: #ff4d4f;
			border: 1rpx solid #ffccc7;
		}

		&.tag-warning {
			background-color: #fffbe6;
			color: #faad14;
			border: 1rpx solid #ffe58f;
		}
	}

	.detail-text {
		font-size: 22rpx;
		color: #8c8c8c;
		display: flex;
		align-items: center;

		&.text-red {
			color: #ff4d4f;
			font-weight: 500;
		}

		&.text-green {
			color: #52c41a;
			font-weight: 500;
		}
	}
}

.status-corner {
	position: absolute;
	top: 0;
	right: 0;
	padding: 8rpx 20rpx;
	border-bottom-left-radius: 16rpx;
	z-index: 2;
	backdrop-filter: blur(4px);
	/* 优化: 毛玻璃效果 */

	text {
		font-size: 24rpx;
		font-weight: 600;
		letter-spacing: 1rpx;
	}

	/* 空闲 - 蓝底白字或蓝字 */
	.is-free & {
		background-color: rgba(24, 144, 255, 0.15);

		text {
			color: #1890ff;
		}
	}

	/* 租用 - 红底白字 */
	.is-occupied & {
		background-color: rgba(255, 77, 79, 0.15);

		text {
			color: #ff4d4f;
		}
	}

	/* 维修 - 黄底白字 */
	.is-maintenance & {
		background-color: rgba(250, 173, 20, 0.15);

		text {
			color: #faad14;
		}
	}
}

/* 列表项入场动画 */
@keyframes fade-in-up {
	from {
		opacity: 0;
		transform: translateY(30rpx);
	}

	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.fade-in-up {
	animation: fade-in-up 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards;
	/* 优化: 更自然的缓动曲线 */
	opacity: 0;
}

/* 角标弹跳动画 */
@keyframes bounce-in {
	0% {
		opacity: 0;
		transform: scale(0.5);
	}

	60% {
		opacity: 1;
		transform: scale(1.1);
	}

	100% {
		transform: scale(1);
	}
}

.bounce-in {
	animation: bounce-in 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
	/* 优化: 弹性效果 */
}

// 新增: 无数据提示样式
.no-data-tip {
	width: 100%;
	padding: 120rpx 0;
	display: flex;
	justify-content: center;
}

// 新增: 底部固定统计样式
.stats-fixed-wrapper {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	z-index: 99;
	padding-bottom: env(safe-area-inset-bottom);
	pointer-events: none;

	.stats-container {
		pointer-events: auto;
		background-color: rgba(255, 255, 255, 0.95);
		/* 优化: 半透明背景 */
		margin: 0 24rpx 30rpx;
		border-radius: 30rpx;
		/* 优化: 更大圆角 */
		padding: 20rpx 10rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05), 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
		/* 优化: 双向阴影增强悬浮感 */
		backdrop-filter: blur(20px);
		/* 优化: 更强毛玻璃 */
		border: 1rpx solid rgba(255, 255, 255, 0.5);

		.stats-item {
			display: flex;
			align-items: center;
			flex: 1;
			justify-content: center;
			padding: 0 10rpx;

			.stats-icon-box {
				width: 52rpx;
				/* 优化: 稍大图标容器 */
				height: 52rpx;
				border-radius: 14rpx;
				display: flex;
				justify-content: center;
				align-items: center;
				margin-right: 14rpx;
				box-shadow: 0 2rpx 6rpx rgba(0, 0, 0, 0.05);
				/* 优化: 图标阴影 */

				&.free-bg {
					background-color: #e6f7ff;
					color: #1890ff;
				}

				&.used-bg {
					background-color: #fff1f0;
					color: #ff4d4f;
				}

				&.maint-bg {
					background-color: #fffbe6;
					color: #faad14;
				}
			}

			.stats-info {
				display: flex;
				flex-direction: column;
				align-items: flex-start;

				.stats-label {
					font-size: 22rpx;
					color: #8c8c8c;
					margin-bottom: 2rpx;
				}

				.stats-num {
					font-size: 30rpx;
					/* 优化: 数字更大更醒目 */
					line-height: 1;
					font-family: 'DIN Alternate', 'Helvetica Neue', Arial, sans-serif;
					/* 优化: 使用等宽或数字友好字体 */

					&.num-free {
						color: #1890ff;
					}

					&.num-used {
						color: #ff4d4f;
					}

					&.num-maintenance {
						color: #faad14;
					}
				}
			}
		}

		.divider {
			width: 1rpx;
			height: 48rpx;
			background-color: #f0f0f0;
			margin: 0 8rpx;
		}
	}
}
</style>