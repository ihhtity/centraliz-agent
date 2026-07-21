<!-- 增值功能页面 -->
<template>
    <view class="container">
        <uv-navbar :title="'增值功能'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
        
        <view class="section">
            <view class="section-title">服务套餐</view>
            <view class="product-list">
                <view class="product-card" v-for="product in products" :key="product.id" @click="handleBuy(product)">
                    <view class="product-icon" :class="product.iconClass">
                        <text>{{ product.icon }}</text>
                    </view>
                    <view class="product-info">
                        <text class="product-name">{{ product.name }}</text>
                        <text class="product-desc">{{ product.desc }}</text>
                    </view>
                    <view class="product-price">
                        <text class="price-symbol">¥</text>
                        <text class="price-value">{{ product.price }}</text>
                        <text class="price-unit">{{ product.unit }}</text>
                    </view>
                </view>
            </view>
        </view>

        <view class="section">
            <view class="section-title">我的服务</view>
            <view class="service-list">
                <view class="service-item" v-for="service in myServices" :key="service.type">
                    <view class="service-left">
                        <view class="service-icon" :class="getServiceIconClass(service.type)">
                            <text>{{ getServiceIcon(service.type) }}</text>
                        </view>
                        <view class="service-info">
                            <text class="service-name">{{ getServiceName(service.type) }}</text>
                            <text class="service-count">{{ service.count }} 次</text>
                        </view>
                    </view>
                    <view class="service-right">
                        <view class="buy-btn" @click.stop="handleBuy(getProductByType(service.type))">购买</view>
                    </view>
                </view>
            </view>
        </view>
    </view>
</template>

<script setup>
import { ref, reactive } from 'vue';

// 商品列表
const products = reactive([
    { id: 1, name: '短信套餐', desc: '发送短信通知', icon: '💬', iconClass: 'sms-icon', price: 10, unit: '100条', type: '0' },
    { id: 2, name: '邮箱服务', desc: '发送邮件通知', icon: '📧', iconClass: 'email-icon', price: 20, unit: '100封', type: '1' },
    { id: 3, name: '广告推广', desc: '平台广告展示', icon: '📢', iconClass: 'ad-icon', price: 100, unit: '月', type: '2' },
    { id: 4, name: '高级会员', desc: '尊享全部特权', icon: '⭐', iconClass: 'vip-icon', price: 200, unit: '年', type: '3' },
]);

// 我的服务
const myServices = ref([
    { type: '0', count: 150 },
    { type: '1', count: 50 },
    { type: '2', count: 0 },
    { type: '3', count: 0 },
]);

// 获取服务图标
const getServiceIcon = (type) => {
    const icons = { '0': '💬', '1': '📧', '2': '📢', '3': '⭐' };
    return icons[type] || '📦';
};

// 获取服务图标样式
const getServiceIconClass = (type) => {
    const classes = { '0': 'sms-icon', '1': 'email-icon', '2': 'ad-icon', '3': 'vip-icon' };
    return classes[type] || 'default-icon';
};

// 获取服务名称
const getServiceName = (type) => {
    const names = { '0': '短信套餐', '1': '邮箱服务', '2': '广告推广', '3': '高级会员' };
    return names[type] || '未知服务';
};

// 根据类型获取商品
const getProductByType = (type) => {
    return products.find(p => p.type === type);
};

// 购买服务
const handleBuy = (product) => {
    uni.showModal({
        title: '购买确认',
        content: `确认购买 ${product.name} (¥${product.price}/${product.unit})？`,
        success: (res) => {
            if (res.confirm) {
                // 调用购买接口
                buyProduct(product);
            }
        }
    });
};

// 调用购买接口
const buyProduct = async (product) => {
    try {
        const merch = uni.getStorageSync('merch') || {};
        const res = await uni.$uv.http.post('/merch-pay/create', {
            merchs_id: merch.id,
            type: product.type,
            name: product.name,
            price: product.price,
        }, { custom: { auth: true } });
        
        if (res.code === 200) {
            uni.showToast({ title: '购买成功', icon: 'success' });
            // 更新服务数量
            const service = myServices.value.find(s => s.type === product.type);
            if (service && product.type === '0') {
                service.count += 100;
            }
        } else {
            uni.showToast({ title: res.msg || '购买失败', icon: 'none' });
        }
    } catch (error) {
        uni.showToast({ title: '购买失败', icon: 'none' });
        console.error('购买失败:', error);
    }
};

// 返回上一页
const goBack = () => {
    uni.redirectTo({
		url: '/pages/admin/profile/index'
	});
};
</script>

<style lang="scss" scoped>
.container {
    min-height: 100vh;
    background-color: #f5f7fa;
    padding-bottom: env(safe-area-inset-bottom);
}

.section {
    padding: 24rpx;
}

.section-title {
    font-size: 28rpx;
    font-weight: 600;
    color: #333;
    margin-bottom: 20rpx;
}

.product-list {
    display: flex;
    flex-direction: column;
    gap: 16rpx;
}

.product-card {
    display: flex;
    align-items: center;
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    transition: background 0.2s;
    
    &:active {
        background: #fafafa;
    }
}

.product-icon {
    width: 80rpx;
    height: 80rpx;
    border-radius: 16rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36rpx;
    margin-right: 20rpx;
    
    &.sms-icon { background: #e8f5e9; }
    &.email-icon { background: #e3f2fd; }
    &.ad-icon { background: #fff3e0; }
    &.vip-icon { background: #f3e8ff; }
}

.product-info {
    flex: 1;
    
    .product-name {
        display: block;
        font-size: 28rpx;
        font-weight: 500;
        color: #333;
        margin-bottom: 6rpx;
    }
    
    .product-desc {
        display: block;
        font-size: 22rpx;
        color: #999;
    }
}

.product-price {
    display: flex;
    align-items: baseline;
    
    .price-symbol {
        font-size: 24rpx;
        color: #ee0a24;
        font-weight: 600;
    }
    
    .price-value {
        font-size: 36rpx;
        color: #ee0a24;
        font-weight: 600;
        margin-right: 4rpx;
    }
    
    .price-unit {
        font-size: 22rpx;
        color: #999;
    }
}

.service-list {
    background: #fff;
    border-radius: 16rpx;
    overflow: hidden;
}

.service-item {
    display: flex;
    align-items: center;
    padding: 24rpx;
    border-bottom: 1rpx solid #f5f5f5;
    
    &:last-child {
        border-bottom: none;
    }
}

.service-left {
    display: flex;
    align-items: center;
    flex: 1;
}

.service-icon {
    width: 64rpx;
    height: 64rpx;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28rpx;
    margin-right: 16rpx;
    
    &.sms-icon { background: #e8f5e9; }
    &.email-icon { background: #e3f2fd; }
    &.ad-icon { background: #fff3e0; }
    &.vip-icon { background: #f3e8ff; }
}

.service-info {
    .service-name {
        display: block;
        font-size: 26rpx;
        color: #333;
        margin-bottom: 4rpx;
    }
    
    .service-count {
        display: block;
        font-size: 22rpx;
        color: #999;
    }
}

.service-right {
    .buy-btn {
        padding: 12rpx 24rpx;
        background: #3c9cff;
        color: #fff;
        border-radius: 24rpx;
        font-size: 24rpx;
        
        &:active {
            opacity: 0.8;
        }
    }
}
</style>
