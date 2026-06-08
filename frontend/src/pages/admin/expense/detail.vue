<!-- 消费详情页面 -->
<template>
    <view class="container">
        <uv-navbar :title="'消费详情'" :placeholder="true" leftIcon="arrow-left" @leftClick="goBack" />
        
        <view v-if="order" class="detail-content">
            <!-- 订单状态 -->
            <view class="status-card" :class="getStatusBgClass(order.status)">
                <view class="status-icon">
                    <text>{{ getStatusIcon(order.status) }}</text>
                </view>
                <text class="status-text">{{ getStatusText(order.status) }}</text>
            </view>

            <!-- 订单信息 -->
            <view class="info-card">
                <view class="info-header">
                    <text class="info-title">订单信息</text>
                </view>
                <view class="info-item">
                    <text class="info-label">订单编号</text>
                    <text class="info-value">{{ order.code || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="info-label">商品名称</text>
                    <text class="info-value">{{ order.name || '-' }}</text>
                </view>
                <view class="info-item">
                    <text class="info-label">订单类型</text>
                    <text class="info-value">{{ getTypeName(order.type) }}</text>
                </view>
                <view class="info-item">
                    <text class="info-label">订单金额</text>
                    <text class="info-value price">¥{{ formatMoney(order.price) }}</text>
                </view>
                <view class="info-item">
                    <text class="info-label">创建时间</text>
                    <text class="info-value">{{ formatTime(order.createdAt) }}</text>
                </view>
                <view class="info-item" v-if="order.hfSeqId">
                    <text class="info-label">汇付订单号</text>
                    <text class="info-value">{{ order.hfSeqId }}</text>
                </view>
                <view class="info-item" v-if="order.remarks">
                    <text class="info-label">订单备注</text>
                    <text class="info-value">{{ order.remarks }}</text>
                </view>
            </view>

            <!-- 操作按钮 -->
            <view class="action-buttons" v-if="order.status === '0'">
                <view class="action-btn cancel" @click="handleCancel">取消订单</view>
                <view class="action-btn pay" @click="handlePay">立即支付</view>
            </view>
        </view>

        <uv-empty v-else mode="data" text="订单信息不存在" />
    </view>
</template>

<script setup>
import { ref, onMounted } from 'vue';

// 订单数据
const order = ref(null);

// 页面加载
onMounted(() => {
    const pages = getCurrentPages();
    const currentPage = pages[pages.length - 1];
    const orderId = currentPage.options?.id;
    
    if (orderId) {
        loadOrderDetail(orderId);
    }
});

// 加载订单详情
const loadOrderDetail = async (orderId) => {
    try {
        const res = await uni.$uv.http.get('/merch-pay/detail', {
            params: { id: orderId },
            custom: { auth: true }
        });
        
        if (res.code === 200 && res.data) {
            order.value = res.data;
        } else {
            uni.showToast({ title: res.msg || '获取订单失败', icon: 'none' });
        }
    } catch (error) {
        console.error('加载订单详情失败:', error);
        uni.showToast({ title: '获取订单失败', icon: 'none' });
    }
};

// 获取状态图标
const getStatusIcon = (status) => {
    switch (status) {
        case '0': return '⏳';
        case '1': return '✅';
        case '2': return '❌';
        default: return '❓';
    }
};

// 获取状态背景样式
const getStatusBgClass = (status) => {
    switch (status) {
        case '0': return 'status-pending';
        case '1': return 'status-success';
        case '2': return 'status-canceled';
        default: return 'status-default';
    }
};

// 获取状态文本
const getStatusText = (status) => {
    switch (status) {
        case '0': return '待支付';
        case '1': return '已支付';
        case '2': return '已关闭';
        default: return '未知状态';
    }
};

// 获取类型名称
const getTypeName = (type) => {
    const names = {
        '0': '短信套餐',
        '1': '邮箱服务',
        '2': '广告推广',
        '3': '高级会员'
    };
    return names[type] || '未知类型';
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

// 取消订单
const handleCancel = () => {
    uni.showModal({
        title: '确认取消',
        content: '确定要取消该订单吗？',
        success: async (res) => {
            if (res.confirm) {
                try {
                    const result = await uni.$uv.http.post('/merch-pay/cancel', {
                        id: order.value.id
                    }, { custom: { auth: true } });
                    
                    if (result.code === 200) {
                        uni.showToast({ title: '取消成功', icon: 'success' });
                        order.value.status = '2';
                    } else {
                        uni.showToast({ title: result.msg || '取消失败', icon: 'none' });
                    }
                } catch (error) {
                    uni.showToast({ title: '取消失败', icon: 'none' });
                }
            }
        }
    });
};

// 支付订单
const handlePay = () => {
    uni.showToast({ title: '支付功能开发中', icon: 'none' });
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
    padding-bottom: env(safe-area-inset-bottom);
}

.status-card {
    margin: 20rpx;
    padding: 32rpx;
    border-radius: 16rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 16rpx;
    
    &.status-pending {
        background: linear-gradient(135deg, #fff8f0 0%, #fff3e0 100%);
        .status-text { color: #ff9500; }
    }
    
    &.status-success {
        background: linear-gradient(135deg, #e8fdf0 0%, #d9f7be 100%);
        .status-text { color: #07c160; }
    }
    
    &.status-canceled {
        background: linear-gradient(135deg, #fef0f0 0%, #ffe4e4 100%);
        .status-text { color: #ee0a24; }
    }
    
    &.status-default {
        background: #f5f7fa;
        .status-text { color: #999; }
    }
}

.status-icon {
    font-size: 40rpx;
}

.status-text {
    font-size: 32rpx;
    font-weight: 600;
}

.info-card {
    margin: 20rpx;
    background: #fff;
    border-radius: 16rpx;
    overflow: hidden;
}

.info-header {
    padding: 24rpx;
    border-bottom: 1rpx solid #f5f5f5;
}

.info-title {
    font-size: 28rpx;
    font-weight: 600;
    color: #333;
}

.info-item {
    display: flex;
    align-items: center;
    padding: 20rpx 24rpx;
    border-bottom: 1rpx solid #fafafa;
    
    &:last-child {
        border-bottom: none;
    }
}

.info-label {
    font-size: 26rpx;
    color: #999;
    width: 160rpx;
}

.info-value {
    font-size: 26rpx;
    color: #333;
    flex: 1;
    
    &.price {
        color: #ee0a24;
        font-weight: 600;
    }
}

.action-buttons {
    display: flex;
    gap: 20rpx;
    padding: 20rpx;
}

.action-btn {
    flex: 1;
    padding: 24rpx;
    border-radius: 12rpx;
    text-align: center;
    font-size: 28rpx;
    font-weight: 500;
    transition: opacity 0.2s;
    
    &:active {
        opacity: 0.8;
    }
    
    &.cancel {
        background: #f5f7fa;
        color: #666;
    }
    
    &.pay {
        background: linear-gradient(135deg, #ee0a24 0%, #ff4d4f 100%);
        color: #fff;
    }
}
</style>
