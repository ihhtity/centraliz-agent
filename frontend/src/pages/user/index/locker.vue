<template>
	<view class="container">
		<uv-navbar :title="t('user.locker.title')" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack" />

		<view class="content">
			<!-- 场所地址 -->
			<view class="location-bar fade-in-up" style="animation-delay: 0.1s;">
				<uv-icon name="map" size="16" color="#4facfe" />
				<text class="location-text">{{ t('user.locker.locationPlaceholder') }}: 北京市朝阳区科技园区A座大厅</text>
			</view>

			<!-- 使用说明入口 -->
			<view class="guide-bar fade-in-up" style="animation-delay: 0.15s;" @click="showUsageGuide">
				<uv-icon name="info-circle" size="16" color="#faad14" />
				<text class="guide-text">{{ t('user.locker.usageGuideTitle') }}</text>
				<uv-icon name="arrow-right" size="12" color="#ccc" />
			</view>

			<!-- 搜索栏 -->
			<view class="search-bar fade-in-up" style="animation-delay: 0.2s;">
				<uv-input v-model="searchKeyword" placeholder="搜索柜子编号，如 A01" prefixIcon="search" :clearable="true"
					shape="circle" class="search-input" />
			</view>

			<!-- 加载状态 -->
			<view v-if="loading" class="loading-container">
				<view class="loading-spinner"></view>
				<text class="loading-text">加载中...</text>
			</view>

			<!-- 空状态 -->
			<view v-else-if="!hasData" class="empty-container">
				<uv-empty mode="data" text="暂无柜子数据" />
			</view>

			<!-- 柜子列表（按分组显示） -->
			<view v-else class="locker-container">
				<view v-for="group in filteredLockerGroups" :key="group.groupId" class="group-section">
					<view class="locker-grid">
						<view v-for="(item, index) in group.lockers" :key="item.id" class="locker-item fade-in-up"
							:style="{ animationDelay: (0.1 + index * 0.05) + 's' }" :class="{
								'is-occupied-other': item.status === '租用' && item.usersid !== user.id,
								'is-mine': item.status === '租用' && item.usersid === user.id,
								'is-free': item.status === '空闲',
								'is-maintenance': item.status === '维修'
							}" @click="handleLockerClick(item)">

							<view class="locker-main">
								<view class="locker-icon-wrapper" :class="{
									'pulse-free': item.status === '空闲',
									'pulse-occupied': item.status === '租用' && item.usersid !== user.id,
									'pulse-mine': item.status === '租用' && item.usersid === user.id
								}">
									<uv-icon name="empty-favor" :size="48" :color="getIconColor(item)" />
								</view>
								<view class="locker-info">
									<text class="locker-no">{{ item.name }}</text>
									<text class="locker-size">{{ item.tag }}</text>
								</view>
							</view>

							<view class="locker-footer">
								<text class="locker-price">¥{{ item.price || 0 }}<text class="unit">{{
									t('user.locker.pricePerHour') }}</text></text>
								<view class="action-hint" v-if="item.usersid === user.id || item.status === '空闲'">
									<uv-icon name="arrow-right" size="12" color="#fff" />
								</view>
							</view>
						</view>
					</view>
				</view>
			</view>
		</view>

		<!-- 底部统计信息 (固定在底部) -->
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
						<uv-icon name="close-circle" size="18" color="#faad14" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.occupied') }}</text>
						<text class="stats-num num-used">{{ stats.used }}</text>
					</view>
				</view>
				<view class="divider"></view>
				<view class="stats-item">
					<view class="stats-icon-box maint-bg">
						<uv-icon name="setting" size="18" color="#999999" />
					</view>
					<view class="stats-info">
						<text class="stats-label">{{ t('user.locker.maintenance') }}</text>
						<text class="stats-num num-maintenance">{{ stats.maintenance }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 底部 TabBar -->
		<uv-tabbar :value="0" :placeholder="true" @change="onTabChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>

		<view>
			<!-- 租用确认弹窗（根据规则类型显示不同内容） -->
			<uv-popup ref="rentPopup" mode="center" :round="20">
				<view class="rent-popup-content" v-if="pendingRentLocker && pendingRentLocker.rules">
					<view class="popup-header">
						<text class="popup-title">租用确认</text>
						<view class="close-btn" @click="rentPopup.close()">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>

					<view class="rent-info">
						<view class="info-row">
							<text class="info-label">柜子编号</text>
							<text class="info-value">{{ pendingRentLocker.name }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">柜子类型</text>
							<text class="info-value">{{ pendingRentLocker.tag }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">规则类型</text>
							<text class="info-value">{{ getRuleTypeText(pendingRentLocker.rules.type) }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">计费模式</text>
							<text class="info-value">{{ getRuleModeText(pendingRentLocker.rules.mode) }}</text>
						</view>

						<!-- 根据规则类型显示不同的价格信息 -->
						<view v-if="pendingRentLocker.rules.type === 'charge'" class="price-section">
							<view v-if="pendingRentLocker.rules.mode === 'pay_time'" class="time-options">
								<text class="section-title">选择套餐</text>
								<view v-for="(option, index) in timeOptions" :key="index"
									class="time-option-item" :class="{ selected: selectedTimeOption === index }"
									@click="selectedTimeOption = index">
									<view class="option-left">
										<text class="option-title">{{ option.title }}</text>
										<text v-if="option.discount && option.discount > 0" class="option-discount">省 ¥{{ option.discount }}</text>
									</view>
									<view class="option-right">
										<text class="option-price">¥{{ calculateOptionTotal(option) }}<text class="option-price-unit"></text></text>
									</view>
								</view>
							</view>
							<view v-else class="price-info">
								<view class="info-row">
									<text class="info-label">单价</text>
									<text class="info-value price-text">¥{{ pendingRentLocker.rules.price }}/{{
										getDurationUnitText(pendingRentLocker.rules.durationUnit) }}</text>
								</view>
								<view v-if="pendingRentLocker.rules.deposit > 0" class="info-row">
									<text class="info-label">押金</text>
									<text class="info-value price-text">¥{{ pendingRentLocker.rules.deposit }}</text>
								</view>
							</view>
						</view>

						<!-- 免费模式提示 -->
						<view v-if="pendingRentLocker.rules.type === 'free'" class="free-hint">
							<uv-icon name="info-circle" color="#1890ff" size="16" />
							<text class="hint-text">当前为免费模式，无需支付即可使用</text>
						</view>

						<!-- 免费时间提示 -->
						<view v-if="pendingRentLocker.rules.freeTime > 0" class="free-time-hint">
							<uv-icon name="clock" color="#52c41a" size="16" />
							<text class="hint-text">前{{ pendingRentLocker.rules.freeTime }}分钟免费</text>
						</view>
					</view>

					<view class="rent-footer">
						<view class="btn-cancel" @click="rentPopup.close()">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" @click="confirmRent">{{ t('common.confirm') }}</view>
					</view>
				</view>
			</uv-popup>
			<!-- 支付弹窗（模拟支付） -->
			<uv-popup ref="paymentPopup" mode="center" :round="20">
				<view class="payment-popup-content">
					<view class="popup-header">
						<text class="popup-title">支付确认</text>
						<view class="close-btn" @click="paymentPopup.close()">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>

					<view class="payment-info">
						<view class="payment-amount">
							<text class="amount-label">支付金额</text>
							<text class="amount-value">¥{{ paymentAmount }}</text>
						</view>
						<view class="payment-detail">
							<view class="detail-row">
								<text class="detail-label">柜子编号</text>
								<text class="detail-value">{{ paymentOrder?.roomName }}</text>
							</view>
							<view class="detail-row">
								<text class="detail-label">计费方式</text>
								<text class="detail-value">{{ paymentOrder?.modeText }}</text>
							</view>
							<view v-if="paymentOrder?.deposit > 0" class="detail-row">
								<text class="detail-label">押金</text>
								<text class="detail-value">¥{{ paymentOrder?.deposit }}</text>
							</view>
						</view>
					</view>

					<view class="payment-methods">
						<text class="method-title">支付方式</text>
						<view class="method-list">
							<view class="method-item" :class="{ selected: selectedPaymentMethod === 'wechat' }"
								@click="selectedPaymentMethod = 'wechat'">
								<uv-icon name="weixin-fill" color="#09bb07" size="24" />
								<text class="method-name">微信支付</text>
							</view>
							<view class="method-item" :class="{ selected: selectedPaymentMethod === 'alipay' }"
								@click="selectedPaymentMethod = 'alipay'">
								<uv-icon name="zhifubao-circle-fill" color="#1677ff" size="24" />
								<text class="method-name">支付宝</text>
							</view>
						</view>
					</view>

					<view class="payment-footer">
						<view class="btn-pay" @click="processPayment">
							<text>立即支付 ¥{{ paymentAmount }}</text>
						</view>
					</view>
				</view>
			</uv-popup>
			<!-- 本人租用的柜子时显示 -->
			<uv-popup ref="actionPopup" mode="center" :round="20">
				<view class="popup-content" v-if="selectedLocker">
					<view class="popup-header">
						<text class="popup-title">{{ selectedLocker.name }}</text>
						<view class="close-btn" @click="actionPopup.close()">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>

					<view class="popup-actions">
						<view class="action-item" @click="handleTempUnlock">
							<view class="action-icon-bg warning-bg">
								<uv-icon name="lock-open" color="#fff" size="24" />
							</view>
							<text class="action-text">{{ t('user.locker.tempUnlock') }}</text>
						</view>

						<view class="action-item" @click="handleEndUsePreCheck">
							<view class="action-icon-bg error-bg">
								<uv-icon name="shopping-cart" color="#fff" size="24" />
							</view>
							<text class="action-text">{{ t('user.locker.endUse') }}</text>
						</view>
					</view>
				</view>
			</uv-popup>
			<!-- 押金支付弹窗 -->
			<uv-popup ref="depositPopup" mode="center" :round="20">
				<view class="deposit-popup-content">
					<view class="popup-header">
						<text class="popup-title">押金支付</text>
						<view class="close-btn" @click="depositPopup.close()">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>

					<view class="deposit-info">
						<view class="deposit-icon">
							<uv-icon name="wallet" color="#faad14" size="48" />
						</view>
						<text class="deposit-desc">使用本柜子需先支付押金，押金将在结束使用后退还</text>
						<view class="deposit-amount-row">
							<text class="deposit-label">押金金额</text>
							<text class="deposit-value">¥{{ depositAmount.toFixed(2) }}</text>
						</view>
					</view>

					<view class="deposit-footer">
						<view class="btn-cancel" @click="depositPopup.close()">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" @click="payDeposit">{{ t('common.confirm') }}</view>
					</view>
				</view>
			</uv-popup>

			<!-- 用户同意使用柜子承诺书组件 -->
			<LockerAgreement ref="agreementPopup" :title="t('user.locker.agreementTitle')"
				:content="t('user.locker.agreementContent')" :agreeText="t('user.locker.agreeAndContinue')"
				:disagreeText="t('user.locker.disagree')" @agree="onAgreementConfirm" @disagree="onAgreementCancel" />
			<!-- 使用说明弹窗 (复用 uv-popup 或简单实现，这里为了简洁使用 uv-popup 展示纯文本) -->
			<uv-popup ref="guidePopup" mode="center" :round="20">
				<view class="guide-popup-content">
					<view class="popup-header">
						<text class="popup-title">{{ t('user.locker.usageGuideTitle') }}</text>
					</view>
					<scroll-view scroll-y class="guide-scroll">
						<text class="guide-detail-text">{{ t('user.locker.usageGuideContent') }}</text>
					</scroll-view>
					<view class="guide-footer">
						<button class="btn-know" @click="guidePopup.close()">{{ t('common.ok') }}</button>
					</view>
				</view>
			</uv-popup>
		</view>
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { onLoad } from '@dcloudio/uni-app'
import { generateLockCommand, recordOperationLog } from '@/utils/utils'
// 引入承诺书组件
import LockerAgreement from '@/components/LockerAgreement.vue'

// 获取扫码进入页面的数据并打印到控制台
onLoad((options) => {
	// 从本地存储获取用户信息
	user.value = uni.getStorageSync('user') || null
	if (options && options.id && options.type) {
		storeScanData(options)
	} else {
		scanData.value = uni.getStorageSync('scanData') || null
	}

	// 通过scanData获取柜子数据
	fetchLockerData()
})

// 国际化翻译
const { t } = useI18n()
// 扫码数据
const scanData = ref(null)
// 用户信息
const user = ref(null)
// 加载状态
const loading = ref(false)
// 是否有数据
const hasData = ref(true)
// 柜子列表（按分组区分）
const lockerGroups = ref([])
// 选中的柜子
const selectedLocker = ref(null)
// 操作弹窗引用
const actionPopup = ref(null)
// 承诺书弹窗引用
const agreementPopup = ref(null)
// 使用说明弹窗引用
const guidePopup = ref(null)
// 租用确认弹窗引用
const rentPopup = ref(null)
// 支付弹窗引用
const paymentPopup = ref(null)
// 暂存待租用的柜子信息
const pendingRentLocker = ref(null)
// 用户是否已同意承诺书的状态
const hasAgreed = ref(false)
// 搜索关键词
const searchKeyword = ref('')
// 时间套餐列表
const timeOptions = ref([])
// 选中的时间套餐索引（预付费模式）
const selectedTimeOption = ref(0)
// 支付金额
const paymentAmount = ref('0')
// 支付订单信息
const paymentOrder = ref(null)
// 选中的支付方式
const selectedPaymentMethod = ref('wechat')
// 押金状态
const hasDeposit = ref(false)
// 押金金额
const depositAmount = ref(0)
// 押金订单ID
const depositOrderId = ref(0)
// 押金支付弹窗引用
const depositPopup = ref(null)
// 是否正在检查押金状态
const checkingDeposit = ref(false)

// 存储扫码数据到本地存储
const storeScanData = (options) => {
	scanData.value = options
	uni.setStorageSync('scanData', options)
	// console.log('URL 参数:', options)
}
// 获取柜子数据
const fetchLockerData = async () => {
	if (!scanData.value || !scanData.value.id || !scanData.value.type) {
		hasData.value = false
		return
	}

	loading.value = true
	try {
		const { id, type } = scanData.value
		let url = ''

		// 根据 type 类型获取不同数据
		if (type === 'room') {
			// 获取单个柜子数据
			url = `/user/room/${id}`
		} else if (type === 'group') {
			// 获取分组下的全部柜子数据
			url = `/user/room/group/${id}`
		} else if (type === 'total') {
			// 获取商家全部分组下的柜子数据
			url = `/user/room/merchant/${id}`
		}

		if (url) {
			const res = await uni.$uv.http.get(url, {
				custom: { auth: true }
			})

			if (res.code === 200 && res.data) {
				if (type === 'room') {
					// 单个柜子，包装成分组格式
					let groupRules = res.data.rules || null
					// 当规则的 mode 为 pay_time 时，将 timeOptions 转换为数组类型
					if (groupRules && groupRules.mode === 'pay_time' && groupRules.timeOptions) {
						groupRules = {
							...groupRules,
							timeOptions: Array.isArray(groupRules.timeOptions)
								? groupRules.timeOptions
								: (typeof groupRules.timeOptions === 'string'
									? JSON.parse(groupRules.timeOptions)
									: [])
						}
					}
					const lockerWithRules = { ...res.data, rules: groupRules }
					lockerGroups.value = [{
						groupId: res.data.groupsId || 0,
						groupName: res.data.groupName || '默认分组',
						groupType: res.data.groupType || '存柜',
						rules: groupRules,
						lockers: [lockerWithRules]
					}]
				} else {
					// 多个分组的数据，将分组的 rules 复制到每个柜子项中
					lockerGroups.value = res.data.map(group => {
						let groupRules = group.rules || null
						// 当规则的 mode 为 pay_time 时，将 timeOptions 转换为数组类型
						if (groupRules && groupRules.mode === 'pay_time' && groupRules.timeOptions) {
							groupRules = {
								...groupRules,
								timeOptions: Array.isArray(groupRules.timeOptions)
									? groupRules.timeOptions
									: (typeof groupRules.timeOptions === 'string'
										? JSON.parse(groupRules.timeOptions)
										: [])
							}
						}
						return {
							...group,
							rules: groupRules,
							lockers: group.lockers.map(locker => ({
								...locker,
								rules: groupRules
							}))
						}
					})
				}
				hasData.value = lockerGroups.value.length > 0 &&
					lockerGroups.value.some(g => g.lockers && g.lockers.length > 0)
			} else {
				hasData.value = false
			}
		}
	} catch (e) {
		console.error('获取柜子数据失败:', e)
		hasData.value = false
		uni.showToast({ title: '获取柜子数据失败', icon: 'none' })
	} finally {
		loading.value = false
	}
}
// 处理柜子点击事件
const handleLockerClick = async (item) => {
	if (item.status === '租用' && item.usersid === user.value?.id) {
		selectedLocker.value = item
		actionPopup.value.open()
	} else if (item.status === '空闲') {
		// 空闲柜子，检查是否已同意承诺书
		if (!hasAgreed.value) {
			// 暂存待租用的柜子信息
			pendingRentLocker.value = item
			// 弹出承诺书
			agreementPopup.value.open()
			return
		}

		// 暂存待租用的柜子信息
		pendingRentLocker.value = item
		selectedLocker.value = item

		// 检查规则是否需要押金
		const rules = item.rules
		if (rules && rules.deposit && rules.deposit > 0) {
			// 需要押金，先检查押金状态
			checkingDeposit.value = true
			try {
				await checkDepositStatus(item.merchsId)
				if (!hasDeposit.value) {
					// 未支付押金，弹出押金支付弹窗
					depositAmount.value = rules.deposit
					depositPopup.value.open()
					return
				}
			} catch (error) {
				console.error('检查押金状态失败:', error)
				uni.showToast({ title: '检查押金状态失败', icon: 'none' })
				return
			} finally {
				checkingDeposit.value = false
			}
		}

		// 已同意且押金已支付（或不需要押金），弹出租用确认弹窗
		// timeOptions 已在 fetchLockerData 中转换为数组类型
		timeOptions.value = pendingRentLocker.value.rules.timeOptions || []
		selectedTimeOption.value = 0 // 重置选中的套餐
		rentPopup.value.open()
	} else if (item.status === '租用' && item.usersid !== user.value?.id) {
		// 他人占用
		uni.showToast({ title: t('user.locker.occupied'), icon: 'none' })
	} else if (item.status === '维修') {
		// 维修中
		uni.showToast({ title: t('user.locker.underMaintenance'), icon: 'none' })
	}
}
// 检查押金状态
const checkDepositStatus = async (merchsId) => {
	try {
		const res = await uni.$uv.http.post('/user/deposit/check', {
			merchsId: merchsId
		}, {
			custom: { auth: true }
		})

		if (res.code === 200 && res.data) {
			hasDeposit.value = res.data.hasDeposit
			depositAmount.value = res.data.deposit || 0
			depositOrderId.value = res.data.orderId || 0
		} else {
			hasDeposit.value = false
			depositAmount.value = 0
			depositOrderId.value = 0
		}
	} catch (error) {
		console.error('检查押金状态失败:', error)
		throw error
	}
}

// 支付押金
const payDeposit = async () => {
	if (!pendingRentLocker.value) return

	const locker = pendingRentLocker.value
	const rules = locker.rules

	try {
		uni.showLoading({ title: '支付押金中...' })

		const depositData = {
			merchsId: locker.merchsId,
			groupsId: locker.groupsId,
			rulesId: rules.id,
			amount: depositAmount.value
		}

		const res = await uni.$uv.http.post('/user/deposit/pay', depositData, {
			custom: { auth: true }
		})

		uni.hideLoading()

		if (res.code === 200) {
			hasDeposit.value = true
			depositOrderId.value = res.data.orderId
			depositPopup.value.close()
			uni.showToast({ title: '押金支付成功', icon: 'success' })

			// 押金支付成功后，自动弹出租用确认弹窗
			timeOptions.value = rules.timeOptions || []
			selectedTimeOption.value = 0
			rentPopup.value.open()
		} else {
			uni.showToast({ title: res.msg || '支付押金失败', icon: 'none' })
		}
	} catch (error) {
		console.error('支付押金失败:', error)
		uni.hideLoading()
		uni.showToast({ title: '支付押金失败', icon: 'none' })
	}
}

// 临时开锁
const handleTempUnlock = async () => {
	let device = selectedLocker.value?.device
	// 记录操作日志（定义在 try 外部，确保 catch 块可以访问）
	let logData = {
		merchsId: selectedLocker.value?.merchsId || 0,
		devicesId: selectedLocker.value?.devicesId || 0,
		roomId: selectedLocker.value?.id || 0,
		code: device?.code || '',
		deviceName: device?.name || '',
		roomName: selectedLocker.value?.name || '',
		control: "临时开锁",
		status: "成功",
		occupant: "用户",
	}

	try {
		// console.log(selectedLocker.value.boardNo);return
		// 关闭弹窗
		actionPopup.value.close()

		uni.showLoading({ title: '开锁中...' })
		let hexData = "574B4C5909" + selectedLocker.value.boardNo + "82" + selectedLocker.value.lockNo
		// 生成锁命令
		const data = {
			code: device.code,
			command: generateLockCommand(hexData)
		}

		const result = await uni.$uv.http.post('/device/common', data, {
			custom: { auth: true }
		})
		uni.hideLoading()
		if (result.code === 200) {
			uni.showToast({ title: '开锁成功', icon: 'success' })
			recordOperationLog(logData)
		} else {
			uni.showToast({ title: result.msg || '开锁失败', icon: 'none' })
			logData.status = "失败"
			recordOperationLog(logData)
		}
	} catch (e) {
		console.error('开锁失败:', e)
		uni.hideLoading()
		uni.showToast({ title: '开锁失败', icon: 'none' })
		logData.status = "失败"
		recordOperationLog(logData)
	}
}
// 扁平化的柜子列表（用于搜索和统计）
const allLockers = computed(() => {
	return lockerGroups.value.flatMap(group => group.lockers || [])
})
// 过滤后的柜子列表（按分组）
const filteredLockerGroups = computed(() => {
	if (!searchKeyword.value.trim()) {
		return lockerGroups.value
	}
	const keyword = searchKeyword.value.trim().toUpperCase()
	return lockerGroups.value.map(group => ({
		...group,
		lockers: group.lockers.filter(item =>
			(item.no || item.name || '').toUpperCase().includes(keyword)
		)
	})).filter(group => group.lockers.length > 0)
})
// 修改: 计算统计信息 (按状态统计，基于所有柜子数据)
const stats = computed(() => {
	const res = {
		free: 0,
		used: 0,
		maintenance: 0
	}

	allLockers.value.forEach(item => {
		if (item.status === '空闲') res.free++
		else if (item.status === '租用') res.used++
		else if (item.status === '维修') res.maintenance++
	})

	return res
})
// 获取柜子图标颜色
const getIconColor = (item) => {
	if (item.status === '维修') return '#bfbfbf' // 维修中灰色
	if (item.usersid === user.value?.id) return '#ffccc7' // 本人租用图标白色
	if (item.status === '租用') return '#ffe58f' // 他人占用图标白色
	return '#4facfe' // 空闲时蓝色
}
// 显示使用说明
const showUsageGuide = () => {
	guidePopup.value.open()
}
// 承诺书同意回调
const onAgreementConfirm = () => {
	// 标记为已同意
	hasAgreed.value = true
	// 保存到本地存储，以便下次进入页面无需再次同意
	uni.setStorageSync('locker_agreement_agreed', true)

	// 如果之前有暂存的租用请求（用户点击空闲柜子后被拦截），直接弹出租用确认弹窗
	if (pendingRentLocker.value) {
		selectedLocker.value = pendingRentLocker.value
		selectedTimeOption.value = 0 // 重置选中的套餐
		rentPopup.value.open()
	}
}
// 承诺书取消/不同意回调
const onAgreementCancel = () => {
	// 不同意，不改变 hasAgreed 状态，保持为 false
	// 用户可以浏览页面，但点击租用时会提示
	pendingRentLocker.value = null
	// 不再强制提示 mustAgree，因为用户选择了不同意，只需关闭弹窗
	// uni.showToast({ title: t('user.locker.mustAgree'), icon: 'none' }) 
}
// 结束使用预检查，引导至支付流程
const handleEndUsePreCheck = () => {
	if (!selectedLocker.value) return

	// 缓存当前选中的柜子信息
	const currentLocker = { ...selectedLocker.value }
	// 先关闭操作弹窗
	actionPopup.value.close()
	// 计算费用（模拟）
	const cost = currentLocker.price * 2 // 假设使用了2小时

	uni.showModal({
		title: t('common.confirm'),
		content: t('user.locker.expectedCost', { cost: cost }),
		confirmText: t('user.locker.immediatePay'),
		success: (res) => {
			if (res.confirm) {
				handlePaymentAndEndUse(cost, currentLocker.id)
			}
		}
	})
}
// 支付并结束使用
const handlePaymentAndEndUse = (amount, lockerId) => {
	uni.showLoading({ title: t('user.locker.paymentProcessing') })

	// 模拟支付请求
	setTimeout(() => {
		uni.hideLoading()

		// 模拟支付成功
		uni.showToast({ title: t('user.locker.paymentSuccess'), icon: 'success' })

		// 支付成功后执行结束使用逻辑
		setTimeout(() => {
			executeEndUse(lockerId)
		}, 1500)

	}, 1500)
}
// 执行结束使用逻辑
const executeEndUse = (lockerId) => {
	// 模拟结束使用API
	uni.showLoading({ title: t('common.loading') })
	setTimeout(() => {
		uni.hideLoading()
		// 更新本地状态 - 遍历所有分组找到对应的柜子
		lockerGroups.value.forEach(group => {
			const locker = group.lockers.find(l => l.id === lockerId)
			if (locker) {
				locker.status = '空闲'
				locker.usersid = 0
				locker.ordersId = null
			}
		})
		selectedLocker.value = null
		uni.showToast({ title: t('user.locker.endSuccess'), icon: 'success' })
	}, 1000)
}
// TabBar 切换逻辑
const onTabChange = () => {
	uni.setStorageSync('userroute', "locker")
	uni.redirectTo({ url: '/pages/user/profile/index' })
}
// 页面加载时检查同意状态
onMounted(() => {
	const agreed = uni.getStorageSync('locker_agreement_agreed')
	if (agreed) {
		hasAgreed.value = true
	} else {
		// 未同意，弹出承诺书
		// 使用 nextTick 或 setTimeout 确保 DOM 渲染完成后打开弹窗，避免某些 UI 库的警告
		setTimeout(() => {
			agreementPopup.value.open()
		}, 500)
	}
})
// 返回扫码
const goBack = () => {
	uni.redirectTo({ url: '/pages/user/index/index' })
}

// 获取规则类型文本
const getRuleTypeText = (type) => {
	const typeMap = {
		'free': '免费模式',
		'charge': '收费模式'
	}
	return typeMap[type] || '未知'
}

// 获取规则模式文本
const getRuleModeText = (mode) => {
	const modeMap = {
		'single': '单次开锁',
		'deposit': '一存一取',
		'pay_single': '单次付费',
		'pay_deposit': '先存后取',
		'pay_hourly': '按时付费',
		'pay_time': '预付费'
	}
	return modeMap[mode] || '未知'
}

// 获取时长单位文本
const getDurationUnitText = (unit) => {
	const unitMap = {
		'minute': '分钟',
		'hour': '小时',
		'day': '天',
		'month': '月'
	}
	return unitMap[unit] || '小时'
}

// 计算套餐总价
const calculateOptionTotal = (option) => {
	const duration = parseInt(option.duration) || 0
	const price = parseFloat(option.price) || 0
	const discount = parseFloat(option.discount) || 0
	return Math.max(0, price * duration - discount).toFixed(2)
}

// 确认租用
const confirmRent = async () => {
	if (!pendingRentLocker.value) return

	const locker = pendingRentLocker.value
	const rules = locker.rules

	try {
		// 免费模式直接租用
		if (rules.type === 'free') {
			await createOrder(locker, rules, 0)
			rentPopup.value.close()
			return
		}

		// 收费模式需要支付
		let amount = 0
		let deposit = rules.deposit || 0
		let modeText = getRuleModeText(rules.mode)

		if (rules.mode === 'pay_time') {
			// 预付费模式
			const option = rules.timeOptions[selectedTimeOption.value]
			amount = parseFloat(calculateOptionTotal(option))
		} else if (rules.mode === 'pay_single') {
			// 单次付费
			amount = rules.price || 0
		} else if (rules.mode === 'pay_deposit') {
			// 先存后取
			amount = rules.price || 0
		} else if (rules.mode === 'pay_hourly') {
			// 按时付费，首次支付押金
			amount = deposit
			modeText = '押金'
		}

		// 设置支付信息
		paymentAmount.value = (amount + deposit).toFixed(2)
		paymentOrder.value = {
			roomName: locker.name,
			modeText: modeText,
			deposit: deposit,
			amount: amount
		}

		// 关闭租用弹窗，打开支付弹窗
		rentPopup.value.close()
		paymentPopup.value.open()
	} catch (error) {
		console.error('确认租用失败:', error)
		uni.showToast({ title: '确认租用失败', icon: 'none' })
	}
}

// 创建订单
const createOrder = async (locker, rules, amount) => {
	try {
		uni.showLoading({ title: '创建订单中...' })

		const orderData = {
			roomId: locker.id,
			merchsId: locker.merchsId,
			groupsId: locker.groupsId,
			rulesId: rules.id,
			mode: rules.mode,
			type: rules.type,
			amount: amount,
			deposit: rules.deposit || 0
		}

		const res = await uni.$uv.http.post('/user/order/create', {
			orderData: orderData,
			custom: { auth: true }
		})

		uni.hideLoading()

		if (res.code === 200) {
			// 更新柜子状态
			locker.status = '租用'
			locker.usersid = user.value.id
			locker.ordersId = res.data.orderId

			uni.showToast({ title: '租用成功', icon: 'success' })
			pendingRentLocker.value = null
			return res.data.orderId
		} else {
			uni.showToast({ title: res.msg || '创建订单失败', icon: 'none' })
			return null
		}
	} catch (error) {
		console.error('创建订单失败:', error)
		uni.hideLoading()
		uni.showToast({ title: '创建订单失败', icon: 'none' })
		return null
	}
}

// 处理支付
const processPayment = async () => {
	if (!paymentOrder.value || !pendingRentLocker.value) return

	try {
		uni.showLoading({ title: '支付中...' })

		// 1. 先创建订单
		const locker = pendingRentLocker.value
		const rules = locker.rules
		const amount = parseFloat(paymentAmount.value)
		
		const orderId = await createOrder(locker, rules, amount)
		
		if (!orderId) {
			uni.hideLoading()
			return
		}

		// 2. 模拟支付请求
		const paymentData = {
			orderId: orderId,
			amount: amount,
			method: selectedPaymentMethod.value
		}

		const payRes = await uni.$uv.http.post('/user/order/payment', paymentData, {
			custom: { auth: true }
		})

		uni.hideLoading()
		
		if (payRes.code === 200) {
			paymentPopup.value.close()
			uni.showToast({ title: '支付成功', icon: 'success' })
		} else {
			uni.showToast({ title: payRes.msg || '支付失败', icon: 'none' })
		}

	} catch (error) {
		console.error('支付失败:', error)
		uni.hideLoading()
		uni.showToast({ title: '支付失败', icon: 'none' })
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f8f9fc;
	/* 更柔和的背景色 */
	padding-bottom: 280rpx;
	/* 为底部统计栏和 TabBar 留出空间 */
	box-sizing: border-box;
}

.content {
	padding: 24rpx;
}

// 场所地址样式
.location-bar {
	display: flex;
	align-items: center;
	background-color: #fff;
	padding: 20rpx 24rpx;
	border-radius: 12rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);

	.location-text {
		font-size: 26rpx;
		color: #333;
		margin-left: 12rpx;
		flex: 1;
	}
}

// 使用说明条目样式
.guide-bar {
	display: flex;
	align-items: center;
	background-color: #fff;
	padding: 20rpx 24rpx;
	border-radius: 12rpx;
	margin-bottom: 20rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.02);

	.guide-text {
		font-size: 26rpx;
		color: #333;
		margin-left: 12rpx;
		flex: 1;
	}
}

// 搜索栏样式
.search-bar {
	margin-bottom: 20rpx;

	.search-input {
		background-color: #fff;
	}
}

.locker-grid {
	display: flex;
	flex-wrap: wrap;
	justify-content: space-between;
	padding-left: 0;
	margin: 0;
}

.locker-item {
	width: 48%;
	height: 260rpx;
	background-color: #ffffff;
	border-radius: 20rpx;
	position: relative;
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	padding: 24rpx;
	box-sizing: border-box;
	transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.2s ease;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	border: 2rpx solid transparent;
	overflow: hidden;

	/* 空闲状态 - 浅蓝色 */
	&.is-free {
		background: linear-gradient(145deg, #e6f7ff 0%, #ffffff 100%);
		border-color: #bae7ff;

		.locker-no,
		.locker-size {
			color: #0050b3;
		}

		.locker-price {
			color: #1890ff;
		}
	}

	/* 被他人占用状态 - 黄色/暖色 */
	&.is-occupied-other {
		background: linear-gradient(145deg, #fffbe6 0%, #ffffff 100%);
		border-color: #ffe58f;

		.locker-icon-wrapper {
			opacity: 0.9;
		}

		.locker-no,
		.locker-size {
			color: #ad6800;
		}

		.locker-price {
			color: #faad14;
		}
	}

	/* 我的柜子状态 - 红色/醒目 */
	&.is-mine {
		background: linear-gradient(145deg, #fff1f0 0%, #ffffff 100%);
		border: 2rpx solid #ffccc7;
		box-shadow: 0 8rpx 24rpx rgba(255, 77, 79, 0.15);
		transform: translateY(-4rpx);

		.locker-no,
		.locker-size {
			color: #cf1322;
		}

		.locker-price {
			color: #ff4d4f;
			font-weight: bold;
		}
	}

	/* 维修状态 - 灰色/不可用 */
	&.is-maintenance {
		background: linear-gradient(145deg, #f5f5f5 0%, #ffffff 100%);
		border: 2rpx dashed #d9d9d9;
		opacity: 0.8;

		.locker-no,
		.locker-size {
			color: #8c8c8c;
		}

		.locker-price {
			color: #bfbfbf;
		}

		.locker-icon-wrapper {
			filter: grayscale(100%);
		}
	}

	&:active {
		transform: scale(0.96);
	}
}

.status-corner {
	position: absolute;
	top: 0;
	right: 0;
	padding: 6rpx 16rpx;
	border-bottom-left-radius: 12rpx;
	z-index: 2;

	text {
		font-size: 25rpx;
		font-weight: bold;
		letter-spacing: 1rpx;
	}

	/* 空闲 - 蓝底白字或蓝字 */
	.is-free & {
		background-color: rgba(24, 144, 255, 0.1);

		text {
			color: #1890ff;
		}
	}

	/* 他人占用 - 黄底黄字 */
	.is-occupied-other & {
		background-color: rgba(250, 173, 20, 0.1);

		text {
			color: #faad14;
		}
	}

	/* 我的柜子 - 红底红字 */
	.is-mine & {
		background-color: rgba(255, 77, 79, 0.1);

		text {
			color: #ff4d4f;
		}
	}

	/* 维修 - 灰底灰字 */
	.is-maintenance & {
		background-color: rgba(0, 0, 0, 0.04);

		text {
			color: #8c8c8c;
		}
	}
}

.locker-main {
	display: flex;
	align-items: center;
	margin-top: 10rpx;
}

.locker-icon-wrapper {
	width: 80rpx;
	height: 80rpx;
	display: flex;
	justify-content: center;
	align-items: center;
	background-color: rgba(255, 255, 255, 0.9);
	border-radius: 50%;
	margin-right: 16rpx;
	box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.05);
	transition: all 0.3s ease;

	/* 空闲状态呼吸灯 - 蓝色 */
	&.pulse-free {
		animation: pulse-blue 2s infinite ease-in-out;
	}

	/* 他人占用呼吸灯 - 黄色 */
	&.pulse-occupied {
		animation: pulse-yellow 2s infinite ease-in-out;
	}

	/* 本人租用呼吸灯 - 红色 */
	&.pulse-mine {
		animation: pulse-red 1.5s infinite ease-in-out;
		box-shadow: 0 0 0 0 rgba(255, 77, 79, 0.4);
	}
}

/* 蓝色呼吸动画 */
@keyframes pulse-blue {
	0% {
		box-shadow: 0 0 0 0 rgba(24, 144, 255, 0.4);
	}

	70% {
		box-shadow: 0 0 0 10rpx rgba(24, 144, 255, 0);
	}

	100% {
		box-shadow: 0 0 0 0 rgba(24, 144, 255, 0);
	}
}

/* 黄色呼吸动画 */
@keyframes pulse-yellow {
	0% {
		box-shadow: 0 0 0 0 rgba(250, 173, 20, 0.4);
	}

	70% {
		box-shadow: 0 0 0 10rpx rgba(250, 173, 20, 0);
	}

	100% {
		box-shadow: 0 0 0 0 rgba(250, 173, 20, 0);
	}
}

/* 红色呼吸动画 */
@keyframes pulse-red {
	0% {
		box-shadow: 0 0 0 0 rgba(255, 77, 79, 0.4);
	}

	70% {
		box-shadow: 0 0 0 12rpx rgba(255, 77, 79, 0);
	}

	100% {
		box-shadow: 0 0 0 0 rgba(255, 77, 79, 0);
	}
}

/* 列表项入场动画 */
@keyframes fade-in-up {
	from {
		opacity: 0;
		transform: translateY(20rpx);
	}

	to {
		opacity: 1;
		transform: translateY(0);
	}
}

.fade-in-up {
	animation: fade-in-up 0.5s ease-out forwards;
	opacity: 0;
	/* 初始隐藏，等待动画执行 */
}

/* 角标弹跳动画 */
@keyframes bounce-in {
	0% {
		opacity: 0;
		transform: scale(0.3);
	}

	50% {
		opacity: 1;
		transform: scale(1.05);
	}

	70% {
		transform: scale(0.9);
	}

	100% {
		transform: scale(1);
	}
}

.bounce-in {
	animation: bounce-in 0.6s cubic-bezier(0.215, 0.61, 0.355, 1) forwards;
}

.locker-info {
	display: flex;
	flex-direction: column;

	.locker-no {
		font-size: 36rpx;
		font-weight: 600;
		color: #333;
		line-height: 1.2;
	}

	.locker-size {
		font-size: 25rpx;
		margin-top: 4rpx;
		background-color: rgba(0, 0, 0, 0.03);
		padding: 2rpx 8rpx;
		border-radius: 6rpx;
		align-self: flex-start;
	}
}

.locker-footer {
	display: flex;
	justify-content: space-between;
	align-items: flex-end;
	margin-top: auto;

	.locker-price {
		font-size: 28rpx;
		font-weight: 500;
		color: #ff9900;

		.unit {
			font-size: 20rpx;
			font-weight: normal;
			opacity: 0.8;
			margin-left: 2rpx;
		}
	}

	.action-hint {
		width: 32rpx;
		height: 32rpx;
		background-color: #faad14;
		border-radius: 50%;
		display: flex;
		justify-content: center;
		align-items: center;
	}
}

// 底部固定统计样式
.stats-fixed-wrapper {
	position: fixed;
	bottom: 100rpx;
	/* 留出 TabBar 高度 */
	left: 0;
	right: 0;
	background-color: transparent;
	/* 改为透明，让阴影由容器承担 */
	z-index: 99;
	padding-bottom: env(safe-area-inset-bottom);
	pointer-events: none;
	/* 允许点击穿透空白区域，如果需要的话，或者保留none防止遮挡 */

	.stats-container {
		pointer-events: auto;
		/* 恢复统计栏的点击事件 */
		background-color: #fff;
		margin: 0 24rpx 24rpx;
		/* 增加外边距，悬浮感 */
		border-radius: 24rpx;
		padding: 24rpx 10rpx;
		display: flex;
		justify-content: space-between;
		align-items: center;
		box-shadow: 0 10rpx 30rpx rgba(0, 0, 0, 0.08);
		backdrop-filter: blur(10px);
		/* 毛玻璃效果 */

		.stats-item {
			display: flex;
			align-items: center;
			flex: 1;
			justify-content: center;

			.stats-icon-box {
				width: 48rpx;
				height: 48rpx;
				border-radius: 12rpx;
				display: flex;
				justify-content: center;
				align-items: center;
				margin-right: 12rpx;

				&.free-bg {
					background-color: rgba(24, 144, 255, 0.1);
				}

				&.used-bg {
					background-color: #ffe58f;
				}

				&.maint-bg {
					background-color: rgba(153, 153, 153, 0.1);
				}
			}

			.stats-info {
				display: flex;
				flex-direction: column;
				align-items: flex-start;

				.stats-label {
					font-size: 22rpx;
					color: #999;
					margin-bottom: 4rpx;
				}

				.stats-num {
					font-size: 34rpx;
					font-weight: 800;
					line-height: 1;

					&.num-free {
						color: #1890ff;
					}

					&.num-used {
						color: #faad14; // 租用改为黄色
					}

					&.num-maintenance {
						color: #999999; // 维修改为灰色
					}
				}
			}
		}

		.divider {
			width: 1rpx;
			height: 40rpx;
			background-color: #f0f0f0;
			margin: 0 10rpx;
		}
	}
}

// 使用说明弹窗样式
.guide-popup-content {
	width: 600rpx;
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	display: flex;
	flex-direction: column;
	max-height: 70vh;
}

.guide-scroll {
	flex: 1;
	min-height: 300rpx;
	max-height: 500rpx;
	margin-bottom: 30rpx;

	.guide-detail-text {
		font-size: 28rpx;
		color: #666;
		line-height: 1.6;
		white-space: pre-wrap;
	}
}

// 押金支付弹窗样式
.deposit-popup-content {
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	min-height: 350rpx;
	width: 600rpx;
}

.deposit-info {
	display: flex;
	flex-direction: column;
	align-items: center;
	padding: 30rpx 0;

	.deposit-icon {
		width: 120rpx;
		height: 120rpx;
		background: linear-gradient(135deg, #fffbe6 0%, #fff7e6 100%);
		border-radius: 50%;
		display: flex;
		justify-content: center;
		align-items: center;
		margin-bottom: 24rpx;
	}

	.deposit-desc {
		font-size: 28rpx;
		color: #666;
		text-align: center;
		margin-bottom: 30rpx;
		line-height: 1.6;
	}

	.deposit-amount-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		width: 100%;
		padding: 24rpx 30rpx;
		background: linear-gradient(145deg, #fffbe6 0%, #fff7e6 100%);
		border-radius: 16rpx;
		border: 2rpx solid #ffe58f;

		.deposit-label {
			font-size: 28rpx;
			color: #999;
		}

		.deposit-value {
			font-size: 44rpx;
			color: #faad14;
			font-weight: 700;
		}
	}
}

.deposit-footer {
	display: flex;
	gap: 20rpx;
	margin-top: 30rpx;

	.btn-cancel,
	.btn-confirm {
		flex: 1;
		height: 80rpx;
		border-radius: 40rpx;
		display: flex;
		justify-content: center;
		align-items: center;
		font-size: 30rpx;
		font-weight: 500;
	}

	.btn-cancel {
		background-color: #f5f5f5;
		color: #666;
	}

	.btn-confirm {
		background: linear-gradient(135deg, #faad14 0%, #ffc53d 100%);
		color: #fff;
	}
}

// 租用确认弹窗样式
.rent-popup-content {
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	min-height: 400rpx;
	width: 650rpx;
}

.rent-info {
	.info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20rpx 0;
		border-bottom: 1rpx solid #f0f0f0;

		.info-label {
			font-size: 28rpx;
			color: #999;
		}

		.info-value {
			font-size: 28rpx;
			color: #333;
			font-weight: 500;

			&.price-text {
				color: #ff4d4f;
				font-weight: bold;
			}
		}
	}

	.price-section {
		margin-top: 24rpx;
		padding: 24rpx;
		background: linear-gradient(145deg, #fffbe6 0%, #fff7e6 100%);
		border-radius: 16rpx;
		border: 2rpx solid #ffe58f;
	}

	.section-title {
		font-size: 28rpx;
		color: #fa8c16;
		font-weight: 500;
		margin-bottom: 20rpx;
		display: block;
		position: relative;
		padding-left: 24rpx;

		&::before {
			content: '';
			position: absolute;
			left: 0;
			top: 50%;
			transform: translateY(-50%);
			width: 6rpx;
			height: 24rpx;
			background: linear-gradient(180deg, #fa8c16 0%, #ffa940 100%);
			border-radius: 3rpx;
		}
	}

	.time-options {
		display: flex;
		flex-direction: column;
		gap: 16rpx;
	}

	.time-option-item {
		padding: 28rpx 24rpx;
		background: #ffffff;
		border: 2rpx solid #e8e8e8;
		border-radius: 16rpx;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.03);
		position: relative;
		overflow: hidden;
		display: flex;
		justify-content: space-between;
		align-items: center;

		&::before {
			content: '';
			position: absolute;
			left: 0;
			top: 0;
			bottom: 0;
			width: 6rpx;
			background: transparent;
			transition: background 0.3s;
		}

		&.selected {
			border-color: #1890ff;
			background: linear-gradient(145deg, #ffffff 0%, #f0f7ff 100%);
			box-shadow: 0 4rpx 20rpx rgba(24, 144, 255, 0.15);
			transform: translateX(4rpx);

			&::before {
				background: linear-gradient(180deg, #1890ff 0%, #69c0ff 100%);
			}

			.option-title {
				color: #1890ff;
			}

			.option-price {
				color: #1890ff;
			}
		}

		&:active {
			transform: scale(0.98);
		}

		.option-left {
			display: flex;
			flex-direction: row;
			gap: 8rpx;
		}

		.option-title {
			font-size: 30rpx;
			color: #262626;
			font-weight: 600;
			display: block;
		}

		.option-discount {
			font-size: 22rpx;
			color: #52c41a;
			font-weight: 500;
			background: rgba(82, 196, 26, 0.1);
			padding: 4rpx 12rpx;
			border-radius: 20rpx;
			display: inline-block;
			align-self: flex-start;
		}

		.option-right {
			text-align: right;
		}

		.option-price {
			font-size: 38rpx;
			color: #ff4d4f;
			font-weight: 700;
			display: block;
			line-height: 1.2;
		}

		.option-price-unit {
			font-size: 24rpx;
			font-weight: 400;
			color: #999;
			margin-left: 4rpx;
		}
	}

	.price-info {
		display: flex;
		flex-direction: column;
		gap: 15rpx;
	}

	.free-hint,
	.free-time-hint {
		display: flex;
		align-items: center;
		padding: 15rpx;
		background-color: #f6ffed;
		border-radius: 8rpx;
		margin-top: 15rpx;

		.hint-text {
			font-size: 24rpx;
			color: #52c41a;
			margin-left: 10rpx;
		}
	}
}

.rent-footer {
	display: flex;
	gap: 20rpx;
	margin-top: 30rpx;

	.btn-cancel,
	.btn-confirm {
		flex: 1;
		height: 80rpx;
		border-radius: 40rpx;
		display: flex;
		justify-content: center;
		align-items: center;
		font-size: 30rpx;
	}

	.btn-cancel {
		background-color: #f5f5f5;
		color: #666;
	}

	.btn-confirm {
		background-color: #1890ff;
		color: #fff;
	}
}

// 支付弹窗样式
.payment-popup-content {
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	min-height: 500rpx;
	width: 650rpx;
}

.payment-info {
	.payment-amount {
		text-align: center;
		padding: 30rpx 0;
		border-bottom: 1rpx solid #f0f0f0;

		.amount-label {
			font-size: 26rpx;
			color: #999;
			display: block;
			margin-bottom: 10rpx;
		}

		.amount-value {
			font-size: 56rpx;
			color: #ff4d4f;
			font-weight: bold;
		}
	}

	.payment-detail {
		padding: 20rpx 0;

		.detail-row {
			display: flex;
			justify-content: space-between;
			align-items: center;
			padding: 15rpx 0;

			.detail-label {
				font-size: 28rpx;
				color: #999;
			}

			.detail-value {
				font-size: 28rpx;
				color: #333;
			}
		}
	}
}

.payment-methods {
	padding: 20rpx 0;
	border-top: 1rpx solid #f0f0f0;
	border-bottom: 1rpx solid #f0f0f0;

	.method-title {
		font-size: 28rpx;
		color: #333;
		font-weight: 500;
		display: block;
		margin-bottom: 20rpx;
	}

	.method-list {
		display: flex;
		gap: 20rpx;
	}

	.method-item {
		flex: 1;
		padding: 25rpx;
		background-color: #f9f9f9;
		border: 2rpx solid #e8e8e8;
		border-radius: 12rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 10rpx;
		transition: all 0.3s;

		&.selected {
			border-color: #1890ff;
			background-color: #e6f7ff;
		}

		.method-name {
			font-size: 26rpx;
			color: #333;
		}
	}
}

.payment-footer {
	margin-top: 30rpx;

	.btn-pay {
		width: 100%;
		height: 88rpx;
		background: linear-gradient(135deg, #ff6b6b 0%, #ff8e53 100%);
		border-radius: 44rpx;
		display: flex;
		justify-content: center;
		align-items: center;

		text {
			font-size: 32rpx;
			color: #fff;
			font-weight: 500;
		}
	}
}

.guide-footer {
	.btn-know {
		background: linear-gradient(90deg, #4facfe 0%, #00f2fe 100%);
		color: #fff;
		border-radius: 40rpx;
		height: 80rpx;
		line-height: 80rpx;
		font-size: 28rpx;
		border: none;

		&::after {
			border: none;
		}
	}
}

// 加载状态样式
.loading-container {
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	padding: 100rpx 0;

	.loading-spinner {
		width: 60rpx;
		height: 60rpx;
		border: 4rpx solid rgba(79, 172, 254, 0.2);
		border-top-color: #4facfe;
		border-radius: 50%;
		animation: spin 1s linear infinite;
	}

	.loading-text {
		font-size: 28rpx;
		color: #999;
		margin-top: 20rpx;
	}
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}

// 空状态容器
.empty-container {
	padding: 100rpx 0;
}

.popup-content {
	background-color: #fff;
	border-radius: 20rpx;
	/* 修改为全圆角，适配中间弹出 */
	padding: 30rpx;
	min-height: 300rpx;
	width: 600rpx;
	/* 增加宽度限制，使中间弹窗更美观 */
}

.popup-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 40rpx;

	.popup-title {
		font-size: 32rpx;
		font-weight: bold;
		color: #333;
	}

	.close-btn {
		padding: 10rpx;
	}
}

.popup-actions {
	display: flex;
	justify-content: space-around;
	padding-bottom: 40rpx;
}

.action-item {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 200rpx;

	.action-icon-bg {
		width: 100rpx;
		height: 100rpx;
		border-radius: 50%;
		display: flex;
		justify-content: center;
		align-items: center;
		margin-bottom: 20rpx;

		&.warning-bg {
			background-color: #ff9900;
		}

		&.error-bg {
			background-color: #ff4d4f;
		}
	}

	.action-text {
		font-size: 28rpx;
		color: #333;
	}
}

// ========== H5浏览器和PC端适配 ==========

// 平板设备 (768px - 1024px)
@media screen and (min-width: 768px) {
	.container {
		padding-bottom: 320rpx;
	}

	.content {
		padding: 30rpx;
		max-width: 900px;
		margin: 0 auto;
	}

	.location-bar {
		padding: 25rpx 30rpx;
		margin-bottom: 25rpx;

		.location-text {
			font-size: 28rpx;
		}
	}

	.guide-bar {
		padding: 25rpx 30rpx;
		margin-bottom: 25rpx;

		.guide-text {
			font-size: 28rpx;
		}
	}

	.search-bar {
		margin-bottom: 25rpx;
	}

	.locker-grid {
		justify-content: flex-start;
	}

	.locker-item {
		width: 31%;
		margin-right: 3%;
		margin-bottom: 25rpx;
		height: 280rpx;
		padding: 30rpx;
		border-radius: 24rpx;

		&:nth-child(3n) {
			margin-right: 0;
		}

		&:active {
			transform: scale(0.98);
		}
	}

	.status-corner {
		padding: 8rpx 20rpx;

		text {
			font-size: 26rpx;
		}
	}

	.locker-icon-wrapper {
		width: 90rpx;
		height: 90rpx;
		margin-right: 20rpx;
	}

	.locker-info {
		.locker-no {
			font-size: 40rpx;
		}

		.locker-size {
			font-size: 26rpx;
			padding: 4rpx 10rpx;
		}
	}

	.locker-footer {
		.locker-price {
			font-size: 30rpx;

			.unit {
				font-size: 22rpx;
			}
		}

		.action-hint {
			width: 36rpx;
			height: 36rpx;
		}
	}

	.stats-fixed-wrapper {
		bottom: 120rpx;

		.stats-container {
			margin: 0 30rpx 30rpx;
			padding: 30rpx 15rpx;
			border-radius: 28rpx;

			.stats-item {
				.stats-icon-box {
					width: 56rpx;
					height: 56rpx;
					margin-right: 15rpx;
				}

				.stats-info {
					.stats-label {
						font-size: 24rpx;
					}

					.stats-num {
						font-size: 38rpx;
					}
				}
			}

			.divider {
				height: 50rpx;
				margin: 0 15rpx;
			}
		}
	}

	.popup-content {
		width: 700rpx;
		padding: 40rpx;
	}

	.popup-header {
		margin-bottom: 50rpx;

		.popup-title {
			font-size: 36rpx;
		}
	}

	.popup-actions {
		padding-bottom: 50rpx;
	}

	.action-item {
		width: 240rpx;

		.action-icon-bg {
			width: 110rpx;
			height: 110rpx;
			margin-bottom: 25rpx;
		}

		.action-text {
			font-size: 30rpx;
		}
	}

	.guide-popup-content {
		width: 700rpx;
		padding: 40rpx;
	}

	.guide-scroll {
		.guide-detail-text {
			font-size: 30rpx;
		}
	}

	.guide-footer {
		.btn-know {
			height: 90rpx;
			line-height: 90rpx;
			font-size: 30rpx;
		}
	}
}

// PC端 (> 1024px)
@media screen and (min-width: 1024px) {
	.container {
		padding-bottom: 360rpx;
		background-color: #e8e8e8;
	}

	.content {
		padding: 40rpx;
		max-width: 1200px;
		margin: 0 auto;
		background-color: #f8f9fc;
		border-radius: 24rpx;
		margin-top: 20rpx;
	}

	.location-bar {
		padding: 30rpx 40rpx;
		margin-bottom: 30rpx;
		border-radius: 16rpx;

		.location-text {
			font-size: 30rpx;
		}
	}

	.guide-bar {
		padding: 30rpx 40rpx;
		margin-bottom: 30rpx;
		border-radius: 16rpx;

		.guide-text {
			font-size: 30rpx;
		}
	}

	.search-bar {
		margin-bottom: 30rpx;
	}

	.locker-grid {
		justify-content: flex-start;
	}

	.locker-item {
		width: 23%;
		margin-right: 2.66%;
		margin-bottom: 30rpx;
		min-width: 280px;
		height: 300rpx;
		padding: 35rpx;
		border-radius: 28rpx;
		transition: transform 0.3s ease, box-shadow 0.3s ease;

		&:nth-child(4n) {
			margin-right: 0;
		}

		&:hover {
			transform: translateY(-8rpx);
			box-shadow: 0 12rpx 30rpx rgba(0, 0, 0, 0.12);
		}

		&:active {
			transform: scale(0.99);
		}
	}

	.status-corner {
		padding: 10rpx 25rpx;

		text {
			font-size: 28rpx;
		}
	}

	.locker-icon-wrapper {
		width: 100rpx;
		height: 100rpx;
		margin-right: 25rpx;
	}

	.locker-info {
		.locker-no {
			font-size: 44rpx;
		}

		.locker-size {
			font-size: 28rpx;
			padding: 6rpx 12rpx;
		}
	}

	.locker-footer {
		.locker-price {
			font-size: 32rpx;

			.unit {
				font-size: 24rpx;
			}
		}

		.action-hint {
			width: 40rpx;
			height: 40rpx;
		}
	}

	.stats-fixed-wrapper {
		bottom: 140rpx;

		.stats-container {
			margin: 0 40rpx 40rpx;
			padding: 35rpx 20rpx;
			border-radius: 32rpx;
			max-width: 1200px;
			left: 50%;
			transform: translateX(-50%);
			position: relative;

			.stats-item {
				.stats-icon-box {
					width: 64rpx;
					height: 64rpx;
					margin-right: 20rpx;
				}

				.stats-info {
					.stats-label {
						font-size: 26rpx;
					}

					.stats-num {
						font-size: 42rpx;
					}
				}
			}

			.divider {
				height: 60rpx;
				margin: 0 20rpx;
			}
		}
	}

	.popup-content {
		width: 800rpx;
		padding: 50rpx;
		border-radius: 24rpx;
	}

	.popup-header {
		margin-bottom: 60rpx;

		.popup-title {
			font-size: 40rpx;
		}

		.close-btn {
			padding: 15rpx;
		}
	}

	.popup-actions {
		padding-bottom: 60rpx;
	}

	.action-item {
		width: 280rpx;

		.action-icon-bg {
			width: 120rpx;
			height: 120rpx;
			margin-bottom: 30rpx;
		}

		.action-text {
			font-size: 32rpx;
		}
	}

	.guide-popup-content {
		width: 800rpx;
		padding: 50rpx;
		border-radius: 24rpx;
	}

	.guide-scroll {
		.guide-detail-text {
			font-size: 32rpx;
		}
	}

	.guide-footer {
		.btn-know {
			height: 100rpx;
			line-height: 100rpx;
			font-size: 32rpx;
		}
	}
}

// 大屏幕PC (> 1440px)
@media screen and (min-width: 1440px) {
	.content {
		max-width: 1400px;
	}

	.locker-item {
		width: 18%;
		margin-right: 2.5%;
	}

	.locker-item:nth-child(4n) {
		margin-right: 2.5%;
	}

	.locker-item:nth-child(5n) {
		margin-right: 0;
	}

	.stats-fixed-wrapper {
		.stats-container {
			max-width: 1400px;
		}
	}
}
</style>