<template>
	<view class="container">
		<!-- 修改: 显式隐藏左侧图标 -->
		<uv-navbar title="自助服务" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack" />
		
		<view class="content">
			<!-- 新增: 售卖柜标题 -->
			<view class="section-title">{{ t('user.index.vendingMachine') }}</view>
			
			<!-- 新增: 商品网格 (一横2个) -->
			<view class="product-grid">
				<view 
					v-for="(item, index) in products" 
					:key="index" 
					class="product-item"
					@click="showProductDetail(item)"
				>
					<image :src="item.image" mode="aspectFit" class="product-image"></image>
					<view class="product-info">
						<text class="product-name">{{ item.name }}</text>
						<text class="product-price">¥{{ item.price.toFixed(2) }}</text>
					</view>
					
					<view class="buy-btn">
						<text>{{ t('user.index.buy') || '购买' }}</text>
					</view>
				</view>
			</view>
		</view>

		<!-- 新增: 底部 TabBar -->
		<uv-tabbar :value="0" :placeholder="true" @change="onTabChange">
			<uv-tabbar-item :text="t('tabBar.home')" icon="home" />
			<uv-tabbar-item :text="t('tabBar.profile')" icon="account" />
		</uv-tabbar>

		<!-- 商品详情弹窗 -->
		<view class="modal-mask" v-if="showModal" @click="closeModal">
			<view class="modal-content" @click.stop>
				<view class="modal-header">
					<text class="modal-title">{{ t('user.index.productDetail') || '商品详情' }}</text>
					<text class="modal-close" @click="closeModal">×</text>
				</view>
				<view class="modal-body" v-if="selectedProduct">
					<image :src="selectedProduct.image" mode="aspectFit" class="modal-image"></image>
					<view class="modal-info">
						<text class="modal-name">{{ selectedProduct.name }}</text>
						<text class="modal-desc">{{ selectedProduct.desc }}</text>
						<text class="modal-price">¥{{ selectedProduct.price.toFixed(2) }}</text>
					</view>
				</view>
				<view class="modal-footer">
					<uv-button 
						size="large" 
						class="modal-cancel-btn"
						@click="closeModal"
					>
						{{ t('common.cancel') || '取消' }}
					</uv-button>
					<uv-button 
						size="large" 
						type="error" 
						class="modal-confirm-btn"
						@click="confirmBuy"
					>
						{{ t('user.index.confirmPay') || '确认购买' }}
					</uv-button>
				</view>
			</view>
		</view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const products = ref([
	{ id: 1, name: '可口可乐', price: 3.5, desc: '经典碳酸饮料，清爽解渴，330ml', image: '//s1.chu0.com/pvimg/img/png/d0/d0ea20a9dd6e443a875cacde928233ac.png?imageMogr2/auto-orient/thumbnail/!239x320r/gravity/Center/crop/239x320/quality/85/%7CimageView2/2/w/239&e=2051020800&token=1srnZGLKZ0Aqlz6dk7yF4SkiYf4eP-YrEOdM1sob:hs2ceDgXyZ6pxA0LQzNEeVhAuzs=' },
	{ id: 2, name: '雪碧', price: 3.5, desc: '柠檬味碳酸饮料，清新爽口，330ml', image: '//s1.aigei.com/src/img/png/6f/6f36e30f8b2e40d9a01ce565a543a350.png?imageMogr2/auto-orient/thumbnail/!282x282r/gravity/Center/crop/282x282/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:fCdwrqUIJ6tze_fpN05PSrFUlKE=' },
	{ id: 3, name: '矿泉水', price: 2.0, desc: '天然矿泉水，纯净无污染，550ml', image: 'https://s1.aigei.com/src/img/png/0d/0d011e77ad9c44d9912b3033abc8a7f9.png?imageMogr2/auto-orient/thumbnail/!282x234r/gravity/Center/crop/282x234/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:8nOLNQcDUwI_Wqj9ovxNRcOJj-8=' },
	{ id: 4, name: '橙汁', price: 5.0, desc: '新鲜橙汁，富含维生素C，450ml', image: '//s1.aigei.com/src/img/png/d2/d268282c5074495a9e001a3a4890050d.png?imageMogr2/auto-orient/thumbnail/!239x320r/gravity/Center/crop/239x320/quality/85/%7CimageView2/2/w/239&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:DqwT9CkT65PH9vksjQRZiUrWv2M=' },
	{ id: 5, name: '薯片', price: 7.5, desc: '香脆薯片，经典原味，100g', image: 'https://s1.aigei.com/src/img/png/0e/0e7c923f0dee438f9ce8feed06c4e094.png?imageMogr2/auto-orient/thumbnail/!282x282r/gravity/Center/crop/282x282/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:Q5PwAkCssGc0zI3dys3SrUrDG9I=' },
	{ id: 6, name: '巧克力', price: 12.0, desc: '丝滑牛奶巧克力，入口即化，80g', image: '//s1.aigei.com/src/img/png/34/3459c99c233c4724b0e748f2a05b1d23.png?imageMogr2/auto-orient/thumbnail/!282x211r/gravity/Center/crop/282x211/quality/85/%7CimageView2/2/w/282&e=2051020800&token=P7S2Xpzfz11vAkASLTkfHN7Fw-oOZBecqeJaxypL:-0okH_kPtQqnKM2SnrzDsZ5NFwM=' },
])

const showModal = ref(false)
const selectedProduct = ref(null)

const showProductDetail = (item) => {
	selectedProduct.value = item
	showModal.value = true
}

const closeModal = () => {
	showModal.value = false
	selectedProduct.value = null
}

const confirmBuy = () => {
	if (!selectedProduct.value) return
	
	uni.showModal({
		title: t('user.index.confirmPay'),
		content: `${selectedProduct.value.name}\n${selectedProduct.value.desc}\n\n价格: ¥${selectedProduct.value.price.toFixed(2)}`,
		cancelText: t('common.cancel'),
		confirmText: t('common.confirm'),
		success: (res) => {
			if (res.confirm) {
				uni.showToast({
					title: t('user.index.paySuccess'),
					icon: 'success'
				})
				closeModal()
			}
		}
	})
}

const onTabChange = () => {
	uni.setStorageSync('userroute', "retail")
	uni.redirectTo({ url: '/pages/user/profile/index' })
}

const goBack = () => {
	uni.setStorageSync('userroute', "retail")
	uni.redirectTo({ url: '/pages/user/index/index'})
}
</script>

<style scoped lang="scss">
// 基础样式 - 小程序和移动端H5
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding-bottom: 100rpx;
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

.product-grid {
	display: flex;
	flex-wrap: wrap;
	justify-content: space-between;
}

.product-item {
	width: 48%;
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
	flex: 1;
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

.buy-btn {
	width: 100%;
	height: 60rpx;
	background-color: #ff4d4f;
	border-radius: 30rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	
	text {
		font-size: 26rpx;
		color: #fff;
		font-weight: bold;
	}
}

// ========== 商品详情弹窗样式 ==========
.modal-mask {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: rgba(0,0,0,0.5);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.modal-content {
	width: 80%;
	max-width: 600rpx;
	background-color: #fff;
	border-radius: 20rpx;
	overflow: hidden;
}

.modal-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 30rpx;
	border-bottom: 1rpx solid #f0f0f0;
}

.modal-title {
	font-size: 32rpx;
	font-weight: bold;
	color: #333;
}

.modal-close {
	font-size: 48rpx;
	color: #999;
	line-height: 1;
}

.modal-body {
	padding: 30rpx;
}

.modal-image {
	width: 100%;
	height: 300rpx;
	border-radius: 12rpx;
	margin-bottom: 20rpx;
}

.modal-info {
	text-align: center;
}

.modal-name {
	font-size: 36rpx;
	font-weight: bold;
	color: #333;
	display: block;
	margin-bottom: 15rpx;
}

.modal-desc {
	font-size: 28rpx;
	color: #666;
	display: block;
	margin-bottom: 20rpx;
	line-height: 1.5;
}

.modal-price {
	font-size: 44rpx;
	color: #ff4d4f;
	font-weight: bold;
	display: block;
}

.modal-footer {
	display: flex;
	border-top: 1rpx solid #f0f0f0;
	
	.modal-cancel-btn,
	.modal-confirm-btn {
		flex: 1;
		border-radius: 0;
	}
}

// ========== H5浏览器和PC端适配 ==========

@media screen and (min-width: 768px) {
	.container {
		padding-bottom: 120rpx;
	}
	
	.content {
		padding: 30rpx;
		max-width: 900px;
		margin: 0 auto;
	}
	
	.section-title {
		font-size: 36rpx;
		margin-bottom: 30rpx;
	}
	
	.product-grid {
		gap: 20rpx;
	}
	
	.product-item {
		width: calc(33.33% - 15rpx);
		padding: 25rpx;
		border-radius: 16rpx;
	}
	
	.product-image {
		height: 240rpx;
	}
	
	.product-name {
		font-size: 30rpx;
	}
	
	.product-price {
		font-size: 36rpx;
	}
	
	.buy-btn {
		height: 70rpx;
		
		text {
			font-size: 28rpx;
		}
	}
	
	.modal-content {
		max-width: 500px;
	}
	
	.modal-image {
		height: 280px;
	}
}

@media screen and (min-width: 1024px) {
	.container {
		padding-bottom: 140rpx;
		background-color: #e8e8e8;
	}
	
	.content {
		padding: 40rpx;
		max-width: 1200px;
		margin: 0 auto;
		background-color: #f5f7fa;
		border-radius: 20rpx;
		margin-top: 20rpx;
	}
	
	.section-title {
		font-size: 40rpx;
		margin-bottom: 40rpx;
	}
	
	.product-grid {
		gap: 30rpx;
		justify-content: flex-start;
	}
	
	.product-item {
		width: calc(25% - 25rpx);
		min-width: 280px;
		padding: 30rpx;
		border-radius: 20rpx;
		margin-bottom: 30rpx;
		transition: transform 0.3s ease, box-shadow 0.3s ease;
		
		&:hover {
			transform: translateY(-5rpx);
			box-shadow: 0 8rpx 20rpx rgba(0,0,0,0.1);
		}
	}
	
	.product-image {
		height: 280rpx;
	}
	
	.product-name {
		font-size: 32rpx;
		margin-bottom: 8rpx;
	}
	
	.product-price {
		font-size: 40rpx;
	}
	
	.buy-btn {
		height: 75rpx;
		
		text {
			font-size: 30rpx;
		}
	}
	
	.modal-content {
		max-width: 600px;
		border-radius: 24rpx;
	}
	
	.modal-header {
		padding: 36rpx;
	}
	
	.modal-title {
		font-size: 36rpx;
	}
	
	.modal-body {
		padding: 36rpx;
	}
	
	.modal-image {
		height: 320px;
		border-radius: 16rpx;
	}
	
	.modal-name {
		font-size: 40rpx;
	}
	
	.modal-desc {
		font-size: 32rpx;
	}
	
	.modal-price {
		font-size: 52rpx;
	}
}

@media screen and (min-width: 1440px) {
	.content {
		max-width: 1400px;
	}
	
	.product-item {
		width: calc(20% - 30rpx);
	}
}
</style>
