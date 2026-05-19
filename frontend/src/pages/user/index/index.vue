<template>
	<view class="container">
		<!-- 修改: 显式隐藏左侧图标 -->
		<uv-navbar :title="t('user.index.title')" :placeholder="true" :leftIcon="''" />
		
		<view class="content">
			<!-- 新增: 售卖柜标题 -->
			<view class="section-title">{{ t('user.index.vendingMachine') }}</view>
			
			<!-- 新增: 商品网格 (一横2个) -->
			<view class="product-grid">
				<view 
					v-for="(item, index) in products" 
					:key="index" 
					class="product-item"
				>
					<image :src="item.image" mode="aspectFit" class="product-image"></image>
					<view class="product-info">
						<text class="product-name">{{ item.name }}</text>
						<text class="product-price">¥{{ item.price.toFixed(2) }}</text>
					</view>
					
					<!-- 新增: 数量选择器 -->
					<view class="quantity-control">
						<uv-button 
							size="mini" 
							color="red"
							:disabled="item.count === 0"
							@click="decreaseCount(index)"
						>-</uv-button>
						<text class="count-text">{{ item.count }}</text>
						<uv-button 
							size="mini" 
							type="primary" 
							@click="increaseCount(index)"
						>+</uv-button>
					</view>
				</view>
			</view>
		</view>

		<!-- 新增: 底部结算栏 -->
		<view class="cart-bar" v-if="totalAmount > 0">
			<view class="cart-info">
				<text class="total-label">{{ t('user.index.total') || '合计:' }}</text>
				<text class="total-amount">¥{{ totalAmount.toFixed(2) }}</text>
			</view>
			<uv-button 
				type="error" 
				shape="circle" 
				@click="handleCheckout"
				class="checkout-btn"
			>
				{{ t('user.index.pay') || '支付' }}
			</uv-button>
		</view>

		<!-- 新增: 底部 TabBar -->
		<uv-tabbar :value="0" :placeholder="true" @change="onTabChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.order')" icon="order" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>
	</view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

// 引入 i18n
const { t } = useI18n()


// 新增: 模拟商品数据
const products = ref([
	{ id: 1, name: '可口可乐', price: 3.5, count: 0, image: '//s1.chu0.com/pvimg/img/png/d0/d0ea20a9dd6e443a875cacde928233ac.png?imageMogr2/auto-orient/thumbnail/!239x320r/gravity/Center/crop/239x320/quality/85/%7CimageView2/2/w/239&e=2051020800&token=1srnZGLKZ0Aqlz6dk7yF4SkiYf4eP-YrEOdM1sob:hs2ceDgXyZ6pxA0LQzNEeVhAuzs=' },
	{ id: 2, name: '雪碧', price: 3.5, count: 0, image: '//s1.aigei.com/src/img/png/6f/6f36e30f8b2e40d9a01ce565a543a350.png?imageMogr2/auto-orient/thumbnail/!282x282r/gravity/Center/crop/282x282/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:fCdwrqUIJ6tze_fpN05PSrFUlKE=' },
	{ id: 3, name: '矿泉水', price: 2.0, count: 0, image: 'https://s1.aigei.com/src/img/png/0d/0d011e77ad9c44d9912b3033abc8a7f9.png?imageMogr2/auto-orient/thumbnail/!282x234r/gravity/Center/crop/282x234/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:8nOLNQcDUwI_Wqj9ovxNRcOJj-8=' },
	{ id: 4, name: '橙汁', price: 5.0, count: 0, image: '//s1.aigei.com/src/img/png/d2/d268282c5074495a9e001a3a4890050d.png?imageMogr2/auto-orient/thumbnail/!239x320r/gravity/Center/crop/239x320/quality/85/%7CimageView2/2/w/239&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:DqwT9CkT65PH9vksjQRZiUrWv2M=' },
	{ id: 5, name: '薯片', price: 7.5, count: 0, image: 'https://s1.aigei.com/src/img/png/0e/0e7c923f0dee438f9ce8feed06c4e094.png?imageMogr2/auto-orient/thumbnail/!282x282r/gravity/Center/crop/282x282/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:Q5PwAkCssGc0zI3dys3SrUrDG9I=' },
	{ id: 6, name: '巧克力', price: 12.0, count: 0, image: '//s1.aigei.com/src/img/png/34/3459c99c233c4724b0e748f2a05b1d23.png?imageMogr2/auto-orient/thumbnail/!282x211r/gravity/Center/crop/282x211/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:-0okH_kPtQqnKM2SnrzDsZ5NFwM=' },
])

// 新增: 增加数量
const increaseCount = (index) => {
	products.value[index].count++
}

// 新增: 减少数量
const decreaseCount = (index) => {
	if (products.value[index].count > 0) {
		products.value[index].count--
	}
}

// 新增: 计算总金额
const totalAmount = computed(() => {
	return products.value.reduce((sum, item) => {
		return sum + (item.price * item.count)
	}, 0)
})

// 新增: 结算付款
const handleCheckout = () => {
	const selectedItems = products.value.filter(item => item.count > 0)
	if (selectedItems.length === 0) return
	
	uni.showModal({
		title: t('user.index.confirmPay'),
		content: `共 ${selectedItems.length} 种商品，合计 ¥${totalAmount.value.toFixed(2)}`,
		cancelText: t('common.cancel'),
		confirmText: t('common.confirm'),
		success: (res) => {
			if (res.confirm) {
				// 模拟支付成功
				uni.showToast({
					title: t('user.index.paySuccess'),
					icon: 'success'
				})
				
				// 清空购物车
				products.value.forEach(item => {
					item.count = 0
				})
				
				// 可选：跳转到订单页
				// setTimeout(() => {
				// 	uni.redirectTo({ url: '/pages/user/order/list' })
				// }, 1500)
			}
		}
	})
}

// 新增: TabBar 切换逻辑
const onTabChange = (index) => {
	if (index === 0) {
		// 当前已是首页，无需操作
		return
	} else if (index === 1) {
		uni.redirectTo({ url: '/pages/user/order/list' })
	} else if (index === 2) {
		uni.redirectTo({ url: '/pages/user/profile/index' })
	}
}
</script>

<style scoped lang="scss">
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 100rpx; /* 为底部结算栏和TabBar留出空间 */
}

.content {
	padding: 20rpx;
}

.section-title {
	font-size: 32rpx;
	font-weight: bold;
	margin-bottom: 20rpx;
	color: #333;
}

// 新增: 商品网格样式
.product-grid {
	display: flex;
	flex-wrap: wrap;
	justify-content: space-between;
}

.product-item {
	width: 48%; /* 一横2个，留少许间隙 */
	background-color: #fff;
	border-radius: 12rpx;
	padding: 20rpx;
	box-sizing: border-box;
	margin-bottom: 20rpx;
	display: flex;
	flex-direction: column;
	align-items: center;
	box-shadow: 0 2rpx 10rpx rgba(0,0,0,0.05);
}

.product-image {
	width: 100%;
	height: 200rpx;
	border-radius: 8rpx;
	margin-bottom: 10rpx;
}

.product-info {
	width: 100%;
	text-align: left;
	margin-bottom: 10rpx;
}

.product-name {
	font-size: 28rpx;
	color: #333;
	display: block;
	margin-bottom: 5rpx;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

.product-price {
	font-size: 32rpx;
	color: #ff4d4f;
	font-weight: bold;
}

.quantity-control {
	display: flex;
	align-items: center;
	justify-content: space-between;
	width: 100%;
	margin-top: 10rpx;
}

.count-text {
	font-size: 28rpx;
	font-weight: bold;
	color: #333;
	min-width: 40rpx;
	text-align: center;
}

// 新增: 底部结算栏样式
.cart-bar {
	position: fixed;
	bottom: 100rpx; /* TabBar 高度约为 100rpx，需调整以确保不被遮挡，具体视 TabBar 配置而定 */
	left: 0;
	right: 0;
	height: 100rpx;
	background-color: #fff;
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 0 30rpx;
	box-shadow: 0 -2rpx 10rpx rgba(0,0,0,0.05);
	z-index: 99;
}

.cart-info {
	display: flex;
	align-items: baseline;
}

.total-label {
	font-size: 28rpx;
	color: #666;
	margin-right: 10rpx;
}

.total-amount {
	font-size: 40rpx;
	color: #ff4d4f;
	font-weight: bold;
}

.checkout-btn {
	width: 200rpx;
	height: 70rpx;
	font-size: 28rpx;
}
</style>
