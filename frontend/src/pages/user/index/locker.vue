<template>
	<view class="container">
		<uv-navbar :title="t('user.locker.title')" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack" />
		<!-- 内容区域 -->
		<view class="content">
			<!-- 场所地址 -->
			<view class="location-bar fade-in-up" style="animation-delay: 0.1s;">
				<uv-icon name="map" size="16" color="#4facfe" />
				<text class="location-text">{{ t('user.locker.locationPlaceholder') }}: {{ selectedGroup ?
					selectedGroup.location : '-' }}</text>
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
			<!-- 搜索结果为空提示 -->
			<view v-else-if="hasData && filteredLockerGroups.length === 0" class="empty-container">
				<uv-empty mode="search" text="未找到匹配的柜子" />
			</view>
			<!-- 柜子列表（按分组显示） -->
			<view v-else class="locker-container">
				<view v-for="grouplist in filteredLockerGroups" :key="grouplist.groupId" class="group-section">
					<view class="locker-grid">
						<view v-for="(item, index) in grouplist.lockers" :key="item.id" class="locker-item fade-in-up"
							:style="{ animationDelay: (0.1 + index * 0.05) + 's' }" :class="{
								'is-occupied-other': item.status === '租用' && item.usersid !== user?.id,
								'is-mine': item.status === '租用' && item.usersid === user?.id,
								'is-free': item.status === '空闲',
								'is-maintenance': item.status === '维修'
							}" @click="handleLockerClick(item, grouplist)">

							<view class="locker-main">
								<view class="locker-icon-wrapper" :class="{
									'pulse-free': item.status === '空闲',
									'pulse-occupied': item.status === '租用' && item.usersid !== user?.id,
									'pulse-mine': item.status === '租用' && item.usersid === user?.id
								}">
									<uv-icon name="empty-favor" :size="48" :color="getIconColor(item)" />
								</view>
								<view class="locker-info">
									<text class="locker-no">{{ item.name }}</text>
									<text class="locker-size">{{ item.tag }}</text>
								</view>
							</view>

							<view class="locker-footer">
								<view v-if="grouplist.rules.type !== 'free'">
									<!-- 单次和存柜 -->
									<text class="locker-price"
										v-if="grouplist.rules.mode === 'pay_single' || grouplist.rules.mode === 'pay_deposit'">¥{{
											grouplist.rules.price || 0 }}</text>
									<!-- 按时租用 -->
									<text class="locker-price" v-if="grouplist.rules.mode === 'pay_hourly'">¥{{
										grouplist.rules.price || 0 }}<text class="unit">/{{
											getDurationUnitText(grouplist.rules.durationUnit) || '小时' }}</text></text>
									<!-- 按套餐租用 -->
									<text class="locker-price" v-if="grouplist.rules.mode === 'pay_time'">选择租用套餐</text>
								</view>
								<view v-else>
									<text class="locker-price">免费使用</text>
								</view>
								<view class="action-hint" v-if="item.usersid === user?.id || item.status === '空闲'">
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
		<!-- 操作弹窗 -->
		<view>
			<!-- 租用确认弹窗（根据规则类型显示不同内容） -->
			<uv-popup ref="rentPopup" mode="center" :round="20">
				<view class="rent-popup-content" v-if="selectedLocker && selectedRule">
					<view class="popup-header">
						<text class="popup-title">租用确认</text>
						<view class="close-btn" @click="rentPopup.close()">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>
					<!-- 租用信息 -->
					<view class="rent-info">
						<view class="info-row">
							<text class="info-label">柜子编号</text>
							<text class="info-value">{{ selectedLocker.name }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">柜子类型</text>
							<text class="info-value">{{ selectedLocker.tag }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">规则类型</text>
							<text class="info-value">{{ getRuleTypeText(selectedRule.type) }}</text>
						</view>
						<!-- 根据规则类型显示不同的价格信息 -->
						<view v-if="selectedRule.type === 'charge'" class="price-section">
							<view v-if="selectedRule.mode === 'pay_time'" class="time-options">
								<text class="section-title">选择套餐</text>
								<view class="options-grid">
									<view v-for="(option, index) in selectedRule.timeOptions" :key="index"
										class="time-option-item" :class="{ selected: selectedTimeOption === index }"
										@click="selectedTimeOption = index">
										<view class="option-top">
											<text class="option-title">{{ option.title }}</text>
											<text v-if="option.discount && option.discount > 0"
												class="option-discount">省 ¥{{ option.discount }}</text>
										</view>
										<text class="option-price">¥{{ calculateEndUseCost(option) }}</text>
									</view>
								</view>
							</view>
							<view v-else class="price-info">
								<view class="info-row">
									<text class="info-label">单价</text>
									<text class="info-value price-text" v-if="selectedRule.mode === 'pay_hourly'">¥{{
										selectedRule.price
									}}/{{
											getDurationUnitText(selectedRule.durationUnit) }}</text>
									<text class="info-value price-text" v-else>¥{{ selectedRule.price }}</text>
								</view>
								<!-- <view v-if="selectedRule.deposit > 0" class="info-row">
									<text class="info-label">押金</text>
									<text class="info-value price-text">¥{{ selectedRule.deposit }}</text>
								</view> -->
							</view>
						</view>
						<!-- 规则类型提示 -->
						<view v-if="selectedRule.type === 'free'">
							<!-- 免费模式提示 -->
							<view class="free-hint">
								<uv-icon name="info-circle" color="#1890ff" size="16" />
								<text class="hint-text">当前为免费模式,无需支付即可使用</text>
							</view>
							<!-- 免费时间提示 -->
							<view v-if="selectedRule.freeTime > 0 || selectedRule.autoEndTime > 0"
								class="free-time-hint">
								<uv-icon name="clock" color="#52c41a" size="16" />
								<view class="hint-content">
									<text class="hint-text">订单开始后,</text>
									<text class="hint-text" v-if="selectedRule.freeTime > 0">前{{
										selectedRule.freeTime }}分钟可以临时开锁.</text>
									<text class="hint-text" v-if="selectedRule.autoEndTime > 0">前{{
										selectedRule.autoEndTime }}分钟后自动结束,请您在指定时间内使用柜子并及时取走物品</text>
								</view>
							</view>
						</view>
						<view v-else>
							<!-- 免费时间提示 -->
							<view v-if="selectedRule.freeTime > 0 || selectedRule.autoEndTime > 0"
								class="free-time-hint">
								<uv-icon name="clock" color="#52c41a" size="16" />
								<view class="hint-content">
									<text class="hint-text">订单开始后,</text>
									<text class="hint-text"
										v-if="selectedRule.freeTime > 0 && selectedRule.mode != 'pay_time'">前{{
											selectedRule.freeTime }}分钟可以临时开锁或者免费结束订单.</text>
									<text class="hint-text" v-if="selectedRule.autoEndTime > 0">前{{
										selectedRule.autoEndTime }}分钟后自动结束,请您在指定时间内使用柜子并及时取走物品</text>
									<text class="hint-text"
										v-if="selectedRule.freeTime > 0 && selectedRule.mode == 'pay_time'">前{{
											selectedRule.freeTime }}分钟可以手动结束订单并全额退款.</text>
									<text class="hint-text" v-if="selectedRule.manualRenew">在订单使用中可以手动续费订单.</text>
									<text class="hint-text" v-if="selectedRule.autoRefund">提前结束订单按剩余时间部分退款</text>
								</view>
							</view>
						</view>
						<!-- 规则描述提示 -->
						<view class="rule-desc">
							<view class="desc-icon">
								<uv-icon name="info-circle" color="#1890ff" size="18" />
							</view>
							<text class="desc-text">{{ selectedRule.description }}</text>
						</view>
					</view>
					<!-- 租赁按钮组 -->
					<view class="rent-footer">
						<view class="btn-cancel" @click="rentPopup.close()">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" @click="confirmRent">{{ t('common.confirm') }}</view>
					</view>
				</view>
			</uv-popup>
			<!-- 本人租用的柜子时显示 -->
			<uv-popup ref="actionPopup" mode="center" :round="20">
				<view class="popup-content" v-if="selectedLocker">
					<!-- 弹窗标题 -->
					<view class="popup-header">
						<text class="popup-title">{{ selectedLocker.name }}</text>
						<view class="close-btn" @click="actionPopup.close()">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>
					<!-- 预付款使用信息 -->
					<view v-if="selectedRule.mode === 'pay_time'" class="combo-info-section">
						<view class="combo-info-header">
							<uv-icon name="clock" color="#faad14" size="20" />
							<text class="combo-info-title">使用时间</text>
						</view>
						<view class="combo-info-content">
							<view class="combo-info-row">
								<text class="combo-info-label">开始时间</text>
								<text class="combo-info-value">{{ selectedLocker?.combo?.startTime || '-' }}</text>
							</view>
							<view class="combo-info-row">
								<text class="combo-info-label">结束时间</text>
								<text class="combo-info-value">{{ selectedLocker?.combo?.endTime || '-' }}</text>
							</view>
						</view>
					</view>
					<!-- 操作按钮 -->
					<view class="popup-actions">
						<!-- 手动续费按钮 -->
						<view class="action-item" v-if="selectedRule.manualRenew && selectedRule.mode === 'pay_time'"
							@click="handleManualRenew">
							<view class="action-icon-bg success-bg">
								<uv-icon name="red-packet" color="#fff" size="24" />
							</view>
							<text class="action-text">{{ t('user.locker.manualRenew') }}</text>
						</view>
						<!-- 临时开锁按钮 -->
						<view class="action-item" v-if="checkingFreeTime" @click="handleTempUnlock">
							<view class="action-icon-bg warning-bg">
								<uv-icon name="lock-open" color="#fff" size="24" />
							</view>
							<text class="action-text">{{ t('user.locker.tempUnlock') }}</text>
						</view>
						<!-- 结束使用按钮 -->
						<view class="action-item" @click="handleEndUsePreCheck">
							<view class="action-icon-bg error-bg">
								<uv-icon name="shopping-cart" color="#fff" size="24" />
							</view>
							<text class="action-text">{{ t('user.locker.endUse') }}</text>
						</view>
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
							<text class="amount-value">¥{{ paymentOrder?.amount }}</text>
						</view>
						<view class="payment-detail">
							<view class="detail-row">
								<text class="detail-label">柜子编号</text>
								<text class="detail-value">{{ paymentOrder?.roomName }}</text>
							</view>
							<view class="detail-row">
								<text class="detail-label">柜子类型</text>
								<text class="detail-value">{{ paymentOrder?.roomTag }}</text>
							</view>
							<view class="detail-row">
								<text class="detail-label">金额</text>
								<text class="detail-value">¥{{ paymentOrder?.amount }}</text>
							</view>
						</view>
					</view>

					<view class="payment-methods" v-if="false">
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
							<text>立即支付 ¥{{ paymentOrder?.amount }}</text>
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
							<uv-icon name="empty-coupon" color="#faad14" size="48" />
						</view>
						<text class="deposit-desc">使用本柜子需先支付押金，押金可以手动申请退款，退款将在一周内审核通过后退还，若急需退款，请联系客服处理</text>
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
			<!-- 手动续费弹窗 -->
			<uv-popup ref="renewPopup" mode="center" :round="20">
				<view class="renew-popup-content">
					<view class="popup-header">
						<text class="popup-title">手动续费</text>
						<view class="close-btn" @click="closeRenewPopup">
							<uv-icon name="close" size="20" color="#999" />
						</view>
					</view>

					<view class="renew-info">
						<view class="renew-locker-info">
							<text class="info-label">柜子编号</text>
							<text class="info-value">{{ selectedLocker?.name || '-' }}</text>
						</view>
						<view class="renew-rule-info">
							<text class="info-label">计费规则</text>
							<text class="info-value">{{ selectedLocker?.combo?.price?.toFixed(2) || '0.00' }}元/{{
								getDurationUnitText(selectedLocker?.combo?.durationUnit || 'hour') }}</text>
						</view>
					</view>

					<view class="renew-input-section">
						<text class="input-label">续费时长（{{ getDurationUnitText(selectedLocker?.combo?.durationUnit ||
							'hour') }}）</text>
						<view class="input-wrapper">
							<input class="renew-input" v-model="renewDuration" type="number" placeholder="请输入续费时长"
								@input="calculateRenewAmount" :maxlength="4" />
							<text class="input-unit">{{ getDurationUnitText(selectedLocker?.combo?.durationUnit ||
								'hour') }}</text>
						</view>
					</view>

					<view class="renew-amount-section">
						<view class="amount-row">
							<text class="amount-label">续费金额</text>
							<text class="amount-value">¥{{ renewAmount.toFixed(2) }}</text>
						</view>
					</view>

					<view class="renew-footer">
						<view class="btn-cancel" @click="closeRenewPopup">{{ t('common.cancel') }}</view>
						<view class="btn-confirm" :class="{ disabled: !canSubmitRenew }" @click="confirmRenew">{{
							t('common.confirm') }}</view>
					</view>
				</view>
			</uv-popup>
			<!-- 使用说明弹窗 (复用 uv-popup 或简单实现，这里为了简洁使用 uv-popup 展示纯文本) -->
			<uv-popup ref="guidePopup" mode="center" :round="20">
				<view class="guide-popup-content">
					<view class="popup-header">
						<text class="popup-title">{{ t('user.locker.usageGuideTitle') }}</text>
					</view>
					<scroll-view scroll-y class="guide-scroll">
						<text class="guide-detail-text">{{ t('user.locker.usageGuideContent') }}</text>
						<text v-if="selectedRule?.description" class="guide-detail-text">.{{ selectedRule?.description }}</text>
					</scroll-view>
					<view class="guide-footer">
						<button class="btn-know" @click="guidePopup.close()">{{ t('common.ok') }}</button>
					</view>
				</view>
			</uv-popup>
			<!-- 用户同意使用柜子承诺书组件 -->
			<LockerAgreement ref="agreementPopup" groupType="存柜" :agreeText="t('user.locker.agreeAndContinue')"
				:disagreeText="t('user.locker.disagree')" @agree="onAgreementConfirm" @disagree="onAgreementCancel" />
		</view>
	</view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { generateLockCommand, recordOperationLog, formatDate, handleMiniProgramLogin, handleMPLogin } from '@/utils/utils'
// 引入承诺书组件
import LockerAgreement from '@/components/LockerAgreement.vue'

// 获取扫码进入页面的数据并打印到控制台
onLoad((options) => {
	// 从本地存储获取用户信息
	user.value = uni.getStorageSync('user') || null
	// 获取微信用户信息
	getWxData(user.value)
	// 存储扫码数据到本地存储
	storeScanData(options)
	// 通过scanData获取柜子数据
	fetchLockerData()
})
// 页面显示时刷新数据
onShow(() => {
	// 首次进入时onLoad已执行过fetchLockerData，这里不再重复执行
	if (!isFirstEnter) {
		fetchLockerData()
	}
	isFirstEnter = false
})
// 页面加载时检查同意状态
onMounted(() => {
	if (user.value && user.value.privacy === '0') {
		// 未同意，弹出承诺书
		agreementPopup.value.open()
	}
})

// 国际化翻译
const { t } = useI18n()
// 扫码数据
const scanData = ref(null)
// 用户信息
const user = ref(null)
// 首次进入标志
let isFirstEnter = true
// 加载状态
const loading = ref(false)
// 是否有数据
const hasData = ref(true)
// 柜子列表（按分组区分）
const lockerGroups = ref([])
// 选中的柜子
const selectedLocker = ref(null)
// 选中的分组
const selectedGroup = ref(null)
// 选中的规则
const selectedRule = ref(null)
// 检查柜子免费时间状态
const checkingFreeTime = ref(false)
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
// 搜索关键词
const searchKeyword = ref('')

// 过滤后的柜子分组
const filteredLockerGroups = computed(() => {
	if (!searchKeyword.value.trim()) {
		return lockerGroups.value
	}

	const keyword = searchKeyword.value.trim().toLowerCase()

	return lockerGroups.value.map(group => {
		const filteredLockers = group.lockers.filter(locker => {
			// 搜索柜子名称或编号
			const name = (locker.name || '').toLowerCase()
			const tag = (locker.tag || '').toLowerCase()
			return name.includes(keyword) || tag.includes(keyword)
		})

		return {
			...group,
			lockers: filteredLockers
		}
	}).filter(group => group.lockers.length > 0)
})
// 选中的时间套餐索引（预付费模式）
const selectedTimeOption = ref(0)
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
// 手动续费弹窗引用
const renewPopup = ref(null)
// 续费时长输入值
const renewDuration = ref('')
// 续费计算金额
const renewAmount = ref(0)

// 获取微信用户信息
const getWxData = (user) => {
	const ua = typeof navigator !== 'undefined' ? navigator.userAgent.toLowerCase() : ''
	const isInWechatBrowser = ua.includes('micromessenger')
	if (!isInWechatBrowser) {
		if (!user || !user.id) {
			uni.showModal({
				title: "登录提示",
				content: "请您先登录!才能使用柜子功能",
				cancelText: "取消",
				confirmText: "登录",
				success: async (res) => {
					if (res.confirm) {
						uni.removeStorageSync('user')
						uni.reLaunch({
							url: '/pages/login/login'
						})
					}
				}
			})
		}
		return;
	}

	// console.log('user', user)
	if (!user || !user.unionId) {
		// #ifdef MP-WEIXIN
		handleMiniProgramLogin(user)
		// #endif
		// #ifdef H5
		handleMPLogin(user)
		// #endif
		return
	}
	if (!user.openid) {
		// #ifdef MP-WEIXIN
		handleMiniProgramLogin(user)
		// #endif
		return
	}
	if (!user.gopenid) {
		// #ifdef H5
		handleMPLogin(user)
		return
		// #endif
	}
}
// 存储扫码数据到本地存储
const storeScanData = (options) => {
	if (options && options.id && options.type) {
		scanData.value = options
		uni.setStorageSync('scanData', options)
	} else {
		scanData.value = uni.getStorageSync('scanData') || null
	}
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

			if (res.code === 200 && res.data && res.data[0]?.rules) {
				lockerGroups.value = res.data
				for (let group of lockerGroups.value) {
					if (group.rules?.mode === 'pay_time' && group.rules?.timeOptions) {
						let rules = group.rules
						group.rules = {
							...rules,
							timeOptions: Array.isArray(rules.timeOptions)
								? rules.timeOptions
								: (typeof rules.timeOptions === 'string'
									? JSON.parse(rules.timeOptions)
									: [])
						}
					}
				}

				selectedGroup.value = lockerGroups.value[0].groups // 默认选中的分组
				selectedRule.value = lockerGroups.value[0]?.rules || {} // 默认选中的规则
				uni.setStorageSync('group', selectedGroup.value)
				// console.log(lockerGroups.value)
				// 检查是否有数据
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
const handleLockerClick = async (item, grouplist) => {
	if (!user.value || !user.value.id) {
		getWxData(user.value)
		return
	}

	selectedLocker.value = item // 选中的柜子
	selectedGroup.value = grouplist.groups // 选中的分组
	selectedRule.value = grouplist.rules // 选中的规则
	// console.log(grouplist);

	if (item.status === '租用' && item.usersid === user.value?.id) {
		// 检查柜子免费时间状态
		if (new Date(item.freeTime).getTime() > new Date().getTime() || selectedRule.value.mode === 'pay_time') {
			checkingFreeTime.value = true
		}
		// 预付费模式获取时间套餐
		if (selectedRule.value.mode === 'pay_time' && !item.combo.price) {
			try {
				selectedLocker.value.combo = JSON.parse(item.combo)
			} catch (e) {
				selectedLocker.value.combo = item.combo || {}
			}
			console.log(selectedLocker.value.combo)
		}
		// 弹出操作弹窗
		actionPopup.value.open()
	} else if (item.status === '空闲') {
		// 空闲柜子，检查是否已同意承诺书
		if (user.value.privacy === '0') {
			// 弹出承诺书
			agreementPopup.value.open()
			return
		}

		// 检查规则是否需要押金
		const deposit = selectedRule.value.deposit
		if (deposit && deposit > 0 && !hasDeposit.value) {
			// 需要押金，先检查押金状态
			checkingDeposit.value = true
			try {
				await checkDepositStatus(item.merchsId)
				if (!hasDeposit.value) {
					// 未支付押金，弹出押金支付弹窗
					depositAmount.value = deposit
					depositPopup.value.open()
					return
				}
			} catch (error) {
				uni.showToast({ title: error, icon: 'none', duration: 3000 })
				return
			} finally {
				checkingDeposit.value = false
			}
		}

		// 已同意且押金已支付（或不需要押金），弹出租用确认弹窗
		selectedTimeOption.value = 0 // 重置选中的套餐
		rentPopup.value.open()
	} else if (item.status === '租用' && item.usersid !== user.value?.id) {
		// 他人占用
		uni.showToast({ title: t('user.locker.occupied'), icon: 'none', duration: 3000 })
	} else if (item.status === '维修') {
		// 维修中
		uni.showToast({ title: t('user.locker.underMaintenance'), icon: 'none', duration: 3000 })
	}
}
// 确认租用
const confirmRent = async () => {
	if (!selectedLocker.value) return

	const locker = selectedLocker.value
	const rules = selectedRule.value

	try {
		// 关闭租用弹框
		rentPopup.value.close()
		// 免费模式直接租用并开锁
		if (rules.type === 'free' || rules.mode !== 'pay_time' && rules.mode !== 'pay_single') {
			// 创建订单
			const orderId = await createOrder()
			// 如果有订单ID，开锁
			if (orderId > 0) await handleUnlock('租用开锁', '', orderId)
			return
		}

		let amount = 0
		if (rules.mode === 'pay_time') {
			// 预付费模式
			const option = rules.timeOptions[selectedTimeOption.value]
			amount = parseFloat(calculateEndUseCost(option))
		} else if (rules.mode === 'pay_single') {
			// 单次付费
			amount = rules.price || 0
		}

		// 设置支付信息
		paymentOrder.value = {
			roomName: locker.name,
			roomTag: locker.tag,
			amount: amount
		}

		// 打开支付弹窗
		paymentPopup.value.open()
	} catch (error) {
		console.error('确认租用失败:', error)
		uni.showToast({ title: '确认租用失败', icon: 'none' })
	}
}
// 创建订单
const createOrder = async (amount = 0) => {
	try {
		const locker = selectedLocker.value
		const rules = selectedRule.value

		uni.showLoading({ title: '租用中...', mask: true })

		const orderData = {
			roomId: locker.id,
			merchsId: locker.merchsId,
			usersId: user.value?.id,
			groupsId: locker.groupsId,
			rulesId: rules.id,
			mode: rules.mode,
			type: rules.type,
			amount: amount,
			deposit: rules.deposit || 0,
			userPhone: user.value?.phone || '',
			merchPhone: selectedGroup.value?.phone || '',
		}
		// console.log(selectedGroup.value);return

		const res = await uni.$uv.http.post('/user/order/create', orderData, {
			custom: { auth: true }
		})

		if (res.code === 200) {
			return res.data.orderId
		} else {
			uni.showToast({ title: res.msg || '租用失败', icon: 'none' })
			return 0
		}
	} catch (error) {
		uni.showToast({ title: '租用失败', icon: 'none' })
		return 0
	}
}
// 结束使用,判断是否需要支付
const handleEndUsePreCheck = () => {
	if (!selectedLocker.value) return

	// 先关闭操作弹窗
	actionPopup.value.close()

	// 免费模式或者不需要再支付，直接结束使用
	if (checkingFreeTime.value || selectedRule.value.type === 'free' || selectedRule.value.mode === 'pay_time') {
		executeEndUse()
		return
	}

	// 根据规则模型计算费用（模拟）
	const price = calculateEndUseCost()

	// 设置支付信息
	paymentOrder.value = {
		roomName: selectedLocker.value.name,
		roomTag: selectedLocker.value.tag,
		amount: price,
		paytype: true
	}

	// 打开支付弹窗
	paymentPopup.value.open()
}
// 结束使用,执行结束使用逻辑
const executeEndUse = async () => {
	try {
		uni.showLoading({ title: "正在结束使用...", mask: true })

		const orderData = {
			roomId: selectedLocker.value.id,
			orderId: selectedLocker.value.ordersId,
			mode: selectedRule.value.mode,
		}
		// console.log(orderData);return

		const res = await uni.$uv.http.post('/user/order/complete', orderData, {
			custom: { auth: true }
		})

		uni.hideLoading()

		if (res.code === 200) {
			let result = await handleUnlock('结束开锁')
			if (!result.status) {
				setTimeout(() => {
					uni.navigateTo({
						url: `/pages/user/order/detail?id=${res.data.orderId}`
					})
				}, 2000)
			}

			// 更新柜子状态
			checkingFreeTime.value = false
			selectedLocker.value = null
			return res.data.orderId
		} else {
			uni.showToast({ title: res.msg || '结束使用失败', icon: 'none' })
			return null
		}
	} catch (error) {
		console.error('结束使用失败:', error)
		uni.hideLoading()
		uni.showToast({ title: '结束使用失败', icon: 'none' })
		return null
	}

}
// 支付并开始使用
const processPayment = async () => {
	// 校验支付信息
	if (!paymentOrder.value || !selectedLocker.value) return
	// 校验支付方式
	if (paymentOrder.value.paytype) {
		handlePaymentAndEndUse()
		return
	}

	try {
		uni.showLoading({ title: '支付中...', mask: true })

		const locker = selectedLocker.value
		const rules = selectedRule.value
		const amount = parseFloat(paymentOrder.value.amount)
		// console.log(selectedRule.value.timeOptions[selectedTimeOption.value]);return

		// 判断是否是续费
		if (paymentOrder.value.isRenew) {
			// 续费流程
			// 模拟支付请求
			const renewData = {
				orderId: locker.ordersId || 0,
				amount: amount,
				method: selectedPaymentMethod.value,
				combo: '',
			}

			if (rules.mode === 'pay_time') {
				let time = locker.combo.endTime
				locker.combo.endTime = formatDate(new Date(time).getTime() + paymentOrder.value.renewDuration * getDurationUnit(locker.combo.durationUnit) * 1000)
				renewData.combo = JSON.stringify(locker.combo)
				renewData.endTime = locker.combo.endTime
			}
			// console.log(locker.combo.endTime);return

			const renewRes = await uni.$uv.http.post('/user/order/renew', renewData, {
				custom: { auth: true }
			})

			if (renewRes.code === 200) {
				paymentPopup.value.close()
				uni.showToast({ title: '续费成功', icon: 'success', duration: 2000 })
				// 刷新柜子数据
				fetchLockerData()
				return renewRes.data
			} else {
				uni.showToast({ title: renewRes.msg || '续费失败', icon: 'none', duration: 2000 })
				return 0
			}
		} else {
			// 正常租用流程
			// 1. 先创建订单
			const orderId = await createOrder(amount)

			if (orderId < 1) {
				uni.showToast({ title: '支付失败', icon: 'none', duration: 2000 })
				return 0
			}

			// 2. 模拟支付请求
			const paymentData = {
				usersId: user.value?.id,
				orderId: orderId,
				amount: amount,
				method: selectedPaymentMethod.value,
				mode: rules.mode,
				combo: '',
			}

			if (rules.mode === 'pay_time') {
				let combo = rules.timeOptions[selectedTimeOption.value] || {}
				combo.startTime = formatDate(new Date())
				combo.endTime = formatDate(new Date().getTime() + combo.duration * getDurationUnit(combo.durationUnit) * 1000)
				paymentData.combo = JSON.stringify(combo)
				paymentData.endTime = combo.endTime
			}

			const payRes = await uni.$uv.http.post('/user/order/payment', paymentData, {
				custom: { auth: true }
			})

			if (payRes.code === 200) {
				paymentPopup.value.close()
				uni.showToast({ title: '支付成功', icon: 'success', duration: 2000 })
				await handleUnlock('支付开锁', '', orderId)
				return payRes.data
			} else {
				uni.showToast({ title: payRes.msg || '支付失败', icon: 'none', duration: 2000 })
				await deleteOrder(orderId)
				return 0
			}
		}
	}
	catch (error) {
		uni.showToast({ title: '支付失败', icon: 'none', duration: 2000 })
		return 0
	}
}
// 删除订单
const deleteOrder = async (orderId) => {
	if (orderId && orderId > 0) {
		try {
			await uni.$uv.http.delete(`/order/${orderId}`, {
				custom: { auth: true }
			})
		} catch (e) {
			console.log('删除订单失败', e)
		}
	}
}
// 支付并结束使用
const handlePaymentAndEndUse = async () => {
	uni.showLoading({ title: t('user.locker.paymentProcessing'), mask: true })

	if (paymentOrder.value.amount <= 0 || !selectedLocker.value) {
		uni.showToast({ title: '网络异常,请稍后重试', icon: 'none', duration: 2000 })
		return
	}

	try {
		uni.showLoading({ title: '支付中...', mask: true })

		let orderId = selectedLocker.value.ordersId
		if (orderId < 1) {
			uni.showToast({ title: '支付失败', icon: 'none', duration: 2000 })
			return 0
		}

		// 2. 模拟支付请求
		const paymentData = {
			orderId: orderId,
			amount: paymentOrder.value.amount,
		}

		const payRes = await uni.$uv.http.post('/user/order/end', paymentData, {
			custom: { auth: true }
		})

		if (payRes.code === 200) {
			paymentPopup.value.close()
			uni.showToast({ title: '支付成功', icon: 'success', duration: 2000 })
			await handleUnlock('结束开锁', '', orderId)
			return payRes.data
		} else {
			uni.showToast({ title: payRes.msg || '支付失败', icon: 'none', duration: 2000 })
			return 0
		}
	} catch (error) {
		uni.showToast({ title: '支付失败', icon: 'none', duration: 2000 })
		return 0
	}
}
// 检查押金状态
const checkDepositStatus = async () => {
	try {
		const res = await uni.$uv.http.post('/user/deposit/check', {
			merchsId: selectedLocker.value.merchsId,
			usersId: user.value?.id
		}, {
			custom: { auth: true }
		})

		if (res.code === 200) {
			hasDeposit.value = res.data.hasDeposit
			depositAmount.value = res.data.deposit || 0
			depositOrderId.value = res.data.orderId || 0
		} else {
			hasDeposit.value = false
			depositAmount.value = 0
			depositOrderId.value = 0
		}
	} catch (err) {
		throw "网络异常,请稍后重试"
	}
}
// 支付押金
const payDeposit = async () => {
	if (!selectedLocker.value) return

	const locker = selectedLocker.value
	const rules = selectedRule.value
	// console.log(locker)

	try {
		uni.showLoading({ title: '支付押金中...', mask: true })

		const depositData = {
			merchsId: locker.merchsId,
			usersId: user.value?.id,
			groupsId: locker.groupsId,
			rulesId: rules.id,
			amount: depositAmount.value,
			merchPhone: selectedGroup.value?.phone || '',
			userPhone: user.value?.phone || '',
		}

		const res = await uni.$uv.http.post('/user/deposit/pay', depositData, {
			custom: { auth: true }
		})

		if (res.code === 200) {
			hasDeposit.value = true
			depositOrderId.value = res.data.orderId
			depositPopup.value.close()
			uni.showToast({ title: '押金支付成功', icon: 'success' })

			// 押金支付成功后，自动弹出租用确认弹窗
			selectedTimeOption.value = 0
			rentPopup.value.open()
		} else {
			uni.showToast({ title: res.msg || '支付押金失败', icon: 'none' })
		}
	} catch (error) {
		console.error('支付押金失败:', error)
		uni.showToast({ title: '支付押金失败', icon: 'none' })
	}
}
// 开锁处理
const handleUnlock = async (control = '开锁', refresh = '', orderId = 0) => {
	let device = selectedLocker.value?.device
	// 记录操作日志（定义在 try 外部，确保 catch 块可以访问）
	let logData = {
		merchsId: selectedLocker.value?.merchsId || 0,
		devicesId: selectedLocker.value?.devicesId || 0,
		roomId: selectedLocker.value?.id || 0,
		roomName: selectedLocker.value?.name || '',
		code: device?.code || '',
		deviceName: device?.name || '',
		phone: user.value?.phone || '',
		control: control,
		status: "成功",
		occupant: "用户",
	}

	try {
		// 先刷新柜子数据
		if (refresh !== 'refresh') {
			await fetchLockerData()
		}

		uni.showLoading({ title: '开锁中...', mask: true })
		let hexData = "574B4C5909" + selectedLocker.value.boardNo + "82" + selectedLocker.value.lockNo
		// 生成锁命令
		const data = {
			code: device.code,
			command: generateLockCommand(hexData)
		}

		const result = await uni.$uv.http.post('/device/common', data, {
			custom: { auth: true }
		})

		if (result.code === 200) {
			uni.showToast({ title: '开锁成功', icon: 'success' })
			recordOperationLog(logData)
			return { status: true, message: '开锁成功' }
		} else {
			// 如果有订单ID，跳转到订单详情页
			if (orderId > 0) {
				setTimeout(() => {
					uni.navigateTo({
						url: `/pages/user/order/detail?id=${orderId}`
					})
				}, 2000)
			}

			uni.showToast({ title: result.msg || '开锁失败', icon: 'none' })
			logData.status = "失败"
			recordOperationLog(logData)
			return { status: false, message: result.msg || '开锁失败' }
		}
	} catch (e) {
		uni.showToast({ title: '开锁失败', icon: 'none' })
		logData.status = "失败"
		recordOperationLog(logData)
		return { status: false, message: '开锁失败' }
	}
}
// 临时开锁
const handleTempUnlock = async () => {
	// 关闭弹窗
	actionPopup.value.close()
	await handleUnlock('临时开锁', 'refresh')
}
// 打开手动续费弹窗
const handleManualRenew = () => {
	// 重置输入
	renewDuration.value = ''
	renewAmount.value = 0
	// 打开续费弹窗
	renewPopup.value.open()
	// 关闭操作弹窗
	actionPopup.value.close()
}
// 关闭续费弹窗
const closeRenewPopup = () => {
	renewDuration.value = ''
	renewAmount.value = 0
	renewPopup.value.close()
}
// 计算续费金额
const calculateRenewAmount = () => {
	const duration = parseInt(renewDuration.value) || 0
	if (duration > 0 && selectedLocker.value?.combo?.price) {
		const price = parseFloat(selectedLocker.value.combo.price) || 0
		renewAmount.value = duration * price
	} else {
		renewAmount.value = 0
	}
}
// 是否可以提交续费
const canSubmitRenew = computed(() => {
	const duration = parseInt(renewDuration.value) || 0
	return duration > 0 && renewAmount.value > 0 && selectedLocker.value?.combo
})
// 确认续费
const confirmRenew = () => {
	if (!canSubmitRenew.value) return

	// 设置支付订单信息
	paymentOrder.value = {
		roomName: selectedLocker.value?.name || '',
		roomTag: selectedLocker.value?.tag || '',
		amount: renewAmount.value,
		paytype: false,
		renewDuration: parseInt(renewDuration.value),
		isRenew: true
	}

	// 关闭续费弹窗
	closeRenewPopup()
	// 打开支付弹窗
	paymentPopup.value.open()
}
// 计算最终使用成本
const calculateEndUseCost = (option = {}) => {
	let rules = selectedRule.value || {}
	// 按小时计费模式
	if (rules.mode === 'pay_hourly') {
		let time = new Date().getTime() - new Date(selectedLocker.value.combo).getTime() // 使用时间（毫秒）
		let unittime = getDurationUnit(rules.durationUnit) // 时间单位
		let duration = Math.max(0, Math.ceil(time / (1000 * unittime))) // 计算使用时长
		// console.log(time,unittime,duration)
		return parseFloat(rules.price * duration)
	}
	// 预付费模式
	if (rules.mode === 'pay_time') {
		const duration = parseInt(option.duration) || 0
		const price = parseFloat(option.price) || 0
		const discount = parseFloat(option.discount) || 0
		return parseFloat(Math.max(0, price * duration - discount).toFixed(2))
	}
	// 存柜模式
	return parseFloat(rules.price)
}
// 显示使用说明
const showUsageGuide = () => {
	guidePopup.value.open()
}
// 承诺书同意回调
const onAgreementConfirm = async () => {
	try {
		if (!user.value || !user.value.id) {
			uni.showToast({ message: '用户不存在', icon: 'none', duration: 2000 })
			return
		}
		if (user.value.privacy === '1') {
			uni.showToast({ message: '您已同意承诺书', icon: 'none', duration: 2000 })
			return
		}

		const res = await uni.$uv.http.put('/user/profile', {
			id: user.value.id,
			privacy: '1'
		}, {
			custom: { auth: true }
		})

		if (res.code === 200) {
			user.value = res.data
			uni.setStorageSync('user', res.data)
			// 如果之前有暂存的租用请求（用户点击空闲柜子后被拦截），直接弹出租用确认弹窗
			if (selectedLocker.value) {
				selectedTimeOption.value = 0 // 重置选中的套餐
				rentPopup.value.open()
			}
		} else {
			uni.showToast({ message: '承诺书同意失败', icon: 'none', duration: 2000 })
			selectedLocker.value = null
		}
	} catch (e) {
		selectedLocker.value = null
	}
}
// 承诺书取消/不同意回调
const onAgreementCancel = () => {
	selectedLocker.value = null
}
// 扁平化的柜子列表（用于搜索和统计）
const allLockers = computed(() => {
	return lockerGroups.value.flatMap(item => item.lockers || [])
})
// 修改: 计算统计信息 (按状态统计，基于所有柜子数据)
const stats = computed(() => {
	return allLockers.value.reduce((acc, item) => {
		const statusMap = {
			'空闲': 'free',
			'租用': 'used',
			'维修': 'maintenance'
		}
		const key = statusMap[item.status]
		if (key) acc[key]++
		return acc
	}, { free: 0, used: 0, maintenance: 0 })
})
// 获取柜子图标颜色
const getIconColor = (item) => {
	if (item.status === '维修') return '#bfbfbf' // 维修中灰色
	if (item.usersid === user.value?.id) return '#ffccc7' // 本人租用图标白色
	if (item.status === '租用') return '#ffe58f' // 他人占用图标白色
	return '#4facfe' // 空闲时蓝色
}
// 获取规则类型文本
const getRuleTypeText = (type) => {
	const typeMap = {
		'free': '免费模式',
		'charge': '收费模式'
	}
	return typeMap[type] || '未知'
}
// 获取时间单位
const getDurationUnit = (unit) => {
	const unitMap = {
		'minute': 60,
		'hour': 60 * 60,
		'day': 60 * 60 * 24,
		'month': 60 * 60 * 24 * 30
	}
	return unitMap[unit] || 60
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
// 返回扫码
const goBack = () => {
	uni.setStorageSync('userroute', "locker")
	uni.redirectTo({ url: '/pages/user/index/index' })
}
// TabBar 切换逻辑
const onTabChange = () => {
	uni.setStorageSync('userroute', "locker")
	uni.redirectTo({ url: '/pages/user/profile/index' })
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
	max-height: 80vh;
	width: 650rpx;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.rent-info {
	flex: 1;
	overflow-y: auto;
	-webkit-overflow-scrolling: touch;

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
		margin-top: 16rpx;
	}

	.options-grid {
		display: flex;
		flex-wrap: wrap;
		gap: 16rpx;
	}

	.time-option-item {
		flex: 0 0 calc(40% - 8rpx);
		padding: 24rpx;
		background: #ffffff;
		border: 2rpx solid #e8e8e8;
		border-radius: 16rpx;
		transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
		box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.03);
		position: relative;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		gap: 12rpx;

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

		.option-top {
			display: flex;
			flex-direction: row;
			gap: 8rpx;
			align-items: center;
		}

		.option-title {
			font-size: 28rpx;
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

		.option-price {
			font-size: 36rpx;
			color: #ff4d4f;
			font-weight: 700;
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

	.rule-desc {
		display: flex;
		align-items: flex-start;
		padding: 20rpx 24rpx;
		background: linear-gradient(135deg, #e6f7ff 0%, #f0f5ff 100%);
		border: 1rpx solid #91d5ff;
		border-radius: 12rpx;
		margin-top: 20rpx;
		gap: 12rpx;

		.desc-icon {
			width: 36rpx;
			height: 36rpx;
			background: #fff;
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			flex-shrink: 0;
			box-shadow: 0 2rpx 8rpx rgba(24, 144, 255, 0.15);
		}

		.desc-text {
			font-size: 26rpx;
			color: #1890ff;
			line-height: 1.6;
			flex: 1;
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

		&.success-bg {
			background-color: #4facfe;
		}

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

// 预付款使用信息样式
.combo-info-section {
	background: linear-gradient(135deg, #fffbf0 0%, #fff7e6 100%);
	border-radius: 12rpx;
	padding: 24rpx;
	margin-bottom: 20rpx;
	border: 1rpx solid #ffe58f;

	.combo-info-header {
		display: flex;
		align-items: center;
		gap: 8rpx;
		margin-bottom: 16rpx;

		.combo-info-title {
			font-size: 28rpx;
			font-weight: 500;
			color: #d48806;
		}
	}

	.combo-info-content {
		background-color: rgba(255, 255, 255, 0.8);
		border-radius: 8rpx;
		padding: 16rpx;
	}

	.combo-info-row {
		display: flex;
		justify-content: space-between;
		align-items: center;

		&:not(:last-child) {
			margin-bottom: 12rpx;
		}

		.combo-info-label {
			font-size: 26rpx;
			color: #999;
		}

		.combo-info-value {
			font-size: 26rpx;
			color: #333;
			font-weight: 500;
		}
	}
}

// 手动续费弹窗样式
.renew-popup-content {
	background-color: #fff;
	border-radius: 20rpx;
	padding: 30rpx;
	width: 650rpx;
}

.renew-info {
	display: flex;
	flex-direction: column;
	gap: 20rpx;
	padding: 20rpx 0;
	border-bottom: 1rpx solid #f0f0f0;

	.renew-locker-info,
	.renew-rule-info {
		display: flex;
		justify-content: space-between;
		align-items: center;

		.info-label {
			font-size: 28rpx;
			color: #999;
		}

		.info-value {
			font-size: 28rpx;
			color: #333;
			font-weight: 500;
		}
	}
}

.renew-input-section {
	padding: 30rpx 0;

	.input-label {
		font-size: 28rpx;
		color: #333;
		display: block;
		margin-bottom: 20rpx;
		font-weight: 500;
	}

	.input-wrapper {
		display: flex;
		align-items: center;
		background-color: #f8f9fc;
		border-radius: 12rpx;
		padding: 0 20rpx;
		height: 80rpx;

		.renew-input {
			flex: 1;
			height: 100%;
			font-size: 32rpx;
			color: #333;
			text-align: left;
		}

		.input-unit {
			font-size: 28rpx;
			color: #999;
			margin-left: 10rpx;
		}
	}
}

.renew-amount-section {
	padding: 20rpx 0;
	border-top: 1rpx solid #f0f0f0;
	border-bottom: 1rpx solid #f0f0f0;

	.amount-row {
		display: flex;
		justify-content: space-between;
		align-items: center;

		.amount-label {
			font-size: 28rpx;
			color: #999;
		}

		.amount-value {
			font-size: 40rpx;
			color: #ff4d4f;
			font-weight: bold;
		}
	}
}

.renew-footer {
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

		&.disabled {
			background: #e8e8e8;
			color: #999;
		}
	}
}
</style>