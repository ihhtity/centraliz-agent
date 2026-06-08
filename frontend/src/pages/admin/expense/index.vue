<!-- 我的消费页面 -->
<template>
    <view class="container">
        <uv-navbar :title="'我的消费'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
        
        <!-- 消费统计 -->
        <view class="stats-card">
            <view class="stats-item">
                <text class="stats-value">¥{{ formatMoney(totalAmount) }}</text>
                <text class="stats-label">累计消费</text>
            </view>
            <view class="stats-divider"></view>
            <view class="stats-item">
                <text class="stats-value">{{ orderCount }}</text>
                <text class="stats-label">订单数量</text>
            </view>
            <view class="stats-divider"></view>
            <view class="stats-item">
                <text class="stats-value">{{ todayAmount }}</text>
                <text class="stats-label">今日消费</text>
            </view>
        </view>

        <!-- 订单列表 -->
        <view class="section">
            <view class="section-header">
                <text class="section-title">消费记录</text>
                <text class="section-count">共 {{ total }} 条</text>
            </view>
            
            <view class="order-list">
                <view class="order-item" v-for="order in orderList" :key="order.id" @click="goToDetail(order.id)">
                    <view class="order-icon" :class="getTypeIconClass(order.type)">
                        <text>{{ getTypeIcon(order.type) }}</text>
                    </view>
                    <view class="order-info">
                        <text class="order-name">{{ order.name || '未知商品' }}</text>
                        <text class="order-time">{{ formatTime(order.createdAt) }}</text>
                    </view>
                    <view class="order-right">
                        <text class="order-amount">-¥{{ formatMoney(order.price) }}</text>
                        <text class="order-status" :class="getStatusClass(order.status)">{{ getStatusText(order.status) }}</text>
                    </view>
                </view>
            </view>
            
            <uv-empty v-if="orderList.length === 0" mode="data" text="暂无消费记录" />
        </view>

        <!-- 分页 -->
        <view v-if="total > pageSize" class="pagination-wrapper">
            <view class="pagination">
                <view 
                    class="page-btn" 
                    :class="{ disabled: currentPage === 1 }"
                    @click="handlePrevPage"
                >
                    <uv-icon name="arrow-left" size="16" color="#666" />
                    <text>上一页</text>
                </view>
                <view class="page-info">
                    <text>{{ currentPage }} / {{ totalPages }}</text>
                </view>
                <view 
                    class="page-btn" 
                    :class="{ disabled: currentPage === totalPages }"
                    @click="handleNextPage"
                >
                    <text>下一页</text>
                    <uv-icon name="arrow-right" size="16" color="#666" />
                </view>
            </view>
        </view>
    </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

// 订单数据
const orderList = ref([]);
const currentPage = ref(1);
const pageSize = 10;
const total = ref(0);

// 统计数据
const totalAmount = ref(0);
const orderCount = ref(0);
const todayAmount = ref(0);

// 总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize));

// 页面加载
onMounted(() => {
    loadOrders();
});

// 加载订单列表
const loadOrders = async () => {
    try {
        const merch = uni.getStorageSync('merch') || {};
        const res = await uni.$uv.http.get('/merch-pay/list', {
            params: {
                merchs_id: merch.id,
                page: currentPage.value,
                size: pageSize
            },
            custom: { auth: true }
        });
        
        if (res.code === 200 && res.data) {
            orderList.value = res.data.list || [];
            total.value = res.data.total || 0;
            totalAmount.value = res.data.totalAmount || 0;
            orderCount.value = res.data.total || 0;
            todayAmount.value = res.data.todayAmount || 0;
        }
    } catch (error) {
        console.error('加载订单失败:', error);
    }
};

// 获取类型图标
const getTypeIcon = (type) => {
    const icons = { '0': '💬', '1': '📧', '2': '📢', '3': '⭐' };
    return icons[type] || '📦';
};

// 获取类型图标样式
const getTypeIconClass = (type) => {
    const classes = {
        '0': 'sms-icon',
        '1': 'email-icon',
        '2': 'ad-icon',
        '3': 'vip-icon'
    };
    return classes[type] || 'default-icon';
};

// 获取状态样式
const getStatusClass = (status) => {
    switch (status) {
        case '1': return 'success';
        case '0': return 'warning';
        case '2': return 'default';
        default: return 'default';
    }
};

// 获取状态文本
const getStatusText = (status) => {
    switch (status) {
        case '0': return '待支付';
        case '1': return '已支付';
        case '2': return '已关闭';
        default: return '未知';
    }
};

// 格式化金额
const formatMoney = (amount) => {
    return parseFloat(amount || 0).toFixed(2);
};

// 格式化时间
const formatTime = (time) => {
    if (!time) return '-'
	return time.replace('T', ' ').substring(0, 19)
};

// 上一页
const handlePrevPage = () => {
    if (currentPage.value > 1) {
        currentPage.value--;
        loadOrders();
    }
};

// 下一页
const handleNextPage = () => {
    if (currentPage.value < totalPages.value) {
        currentPage.value++;
        loadOrders();
    }
};

// 跳转到详情
const goToDetail = (orderId) => {
    uni.navigateTo({
        url: `/pages/admin/expense/detail?id=${orderId}`
    });
};

// 返回上一页
const goBack = () => {
    uni.navigateBack();
};
</script>

<style lang="scss" scoped>
.container {
    min-height: 100vh;
    background-color: #f5f7fa;
    padding-bottom: 160rpx;
}

.stats-card {
    display: flex;
    justify-content: space-around;
    align-items: center;
    background: linear-gradient(135deg, #3c9cff 0%, #6b8cff 100%);
    margin: 20rpx;
    border-radius: 16rpx;
    padding: 32rpx;
    color: #fff;
}

.stats-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    
    .stats-value {
        font-size: 36rpx;
        font-weight: 600;
    }
    
    .stats-label {
        font-size: 22rpx;
        opacity: 0.8;
        margin-top: 6rpx;
    }
}

.stats-divider {
    width: 1rpx;
    height: 60rpx;
    background: rgba(255, 255, 255, 0.3);
}

.section {
    padding: 0 20rpx;
}

.section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16rpx;
}

.section-title {
    font-size: 26rpx;
    font-weight: 600;
    color: #333;
}

.section-count {
    font-size: 22rpx;
    color: #999;
}

.order-list {
    background: #fff;
    border-radius: 12rpx;
    overflow: hidden;
}

.order-item {
    display: flex;
    align-items: center;
    padding: 24rpx;
    border-bottom: 1rpx solid #f5f5f5;
    
    &:last-child {
        border-bottom: none;
    }
    
    &:active {
        background: #fafafa;
    }
}

.order-icon {
    width: 48rpx;
    height: 48rpx;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 22rpx;
    margin-right: 16rpx;
    
    &.sms-icon { background: #e8f5e9; color: #07c160; }
    &.email-icon { background: #e3f2fd; color: #1890ff; }
    &.ad-icon { background: #fff3e0; color: #ff9500; }
    &.vip-icon { background: #f3e8ff; color: #9333ea; }
    &.default-icon { background: #f5f7fa; color: #999; }
}

.order-info {
    flex: 1;
    
    .order-name {
        display: block;
        font-size: 26rpx;
        color: #333;
        margin-bottom: 6rpx;
    }
    
    .order-time {
        display: block;
        font-size: 22rpx;
        color: #999;
    }
}

.order-right {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 4rpx;
    
    .order-amount {
        font-size: 28rpx;
        font-weight: 600;
        color: #ee0a24;
    }
    
    .order-status {
        font-size: 20rpx;
        padding: 2rpx 8rpx;
        border-radius: 6rpx;
        
        &.success {
            background: #e8fdf0;
            color: #07c160;
        }
        
        &.warning {
            background: #fff8f0;
            color: #ff9500;
        }
        
        &.default {
            background: #f5f7fa;
            color: #999;
        }
    }
}

.pagination-wrapper {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    background: #fff;
    border-top: 1rpx solid #f0f0f0;
    padding: 16rpx 20rpx;
    padding-bottom: calc(16rpx + env(safe-area-inset-bottom));
}

.pagination {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 40rpx;
}

.page-btn {
    display: flex;
    align-items: center;
    gap: 6rpx;
    padding: 12rpx 20rpx;
    background: #f5f7fa;
    border-radius: 8rpx;
    font-size: 24rpx;
    color: #666;
    
    &.disabled {
        opacity: 0.4;
        pointer-events: none;
    }
    
    &:active:not(.disabled) {
        background: #e8eaed;
    }
}

.page-info {
    font-size: 24rpx;
    color: #666;
}
</style>
