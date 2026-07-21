<!-- 规则管理页面 -->
<template>
    <view class="container">
        <uv-navbar :title="'规则管理'" :placeholder="true" @leftClick="goBack" />
        
        <!-- 规则列表 -->
        <view class="rule-list">
            <view 
                v-for="item in ruleList" 
                :key="item.id" 
                class="rule-card"
            >
                <view class="rule-left" @click="editRule(item)">
                    <view class="rule-row">
                        <text class="label">名称:</text>
                        <text class="value">{{ item.name }}</text>
                    </view>
                    <view class="rule-row">
                        <text class="label">类型:</text>
                        <text class="value">{{ item.type === 'free' ? '免费' : '收费' }}</text>
                    </view>
                    <view class="rule-row">
                        <text class="label">模式:</text>
                        <text class="value">{{ getModeText(item.mode) }}</text>
                    </view>
                    <view v-if="item.price" class="rule-row">
                        <text class="label">单价:</text>
                        <text class="value">¥{{ item.price.toFixed(2) }}</text>
                    </view>
                    <view v-if="item.deposit" class="rule-row">
                        <text class="label">押金:</text>
                        <text class="value">¥{{ item.deposit.toFixed(2) }}</text>
                    </view>
                    <view class="rule-row">
                        <text class="label">创建时间:</text>
                        <text class="value">{{ formatDateTime(item.createdAt) }}</text>
                    </view>
                </view>
                <view class="rule-right">
                    <view class="action-btn edit" @click="editRule(item)">
                        <text>编辑</text>
                    </view>
                    <view class="action-btn delete" @click="deleteRule(item)">
                        <text>删除</text>
                    </view>
                </view>
            </view>
            
            <view v-if="ruleList.length === 0" class="empty-state">
                <text class="empty-icon">📋</text>
                <text class="empty-text">暂无规则，点击下方按钮添加</text>
            </view>
        </view>

        <!-- 底部添加按钮 -->
        <view class="add-btn-wrapper">
            <view class="add-btn" @click="addRule">
                <text>添加规则</text>
            </view>
        </view>
    </view>
</template>

<script setup>
import { ref } from 'vue';
import { onShow } from '@dcloudio/uni-app'

const ruleList = ref([]);

// 获取模式文本
const getModeText = (mode) => {
    const modes = {
        'single': '单次开锁',
        'deposit': '一存一取',
        'pay_single': '单次付费',
        'pay_deposit': '先存后取',
        'pay_hourly': '按时付费',
        'pay_time': '预付费'
    };
    return modes[mode] || mode;
};

// 格式化日期时间显示
const formatDateTime = (dateStr) => {
    if (!dateStr) return '';
    try {
        const date = new Date(dateStr);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        return `${year}-${month}-${day} ${hours}:${minutes}`;
    } catch {
        return dateStr;
    }
};

// 页面显示时加载规则列表
onShow(() => {
    loadRules();
});

// 加载规则列表
const loadRules = async () => {
    try {
        const merch = uni.getStorageSync('merch') || {};
        const res = await uni.$uv.http.get('/rule/list', {
            params: { merchs_id: merch.id },
            custom: { auth: true }
        });
        
        if (res.code === 200) {
            ruleList.value = res.data || [];
        }
    } catch (error) {
        console.error('加载规则失败:', error);
    }
};

// 返回上一页
const goBack = () => {
    uni.redirectTo({
		url: '/pages/admin/profile/index'
	});
};

// 编辑规则
const editRule = (item) => {
    uni.navigateTo({
        url: `/pages/admin/rule/edit?id=${item.id}`
    });
};

// 删除规则
const deleteRule = (item) => {
    uni.showModal({
        title: '确认删除',
        content: `确定要删除规则「${item.name}」吗？`,
        success: async (res) => {
            if (res.confirm) {
                try {
                    const result = await uni.$uv.http.delete(`/rule/${item.id}`, {
                        custom: { auth: true }
                    });
                    
                    if (result.code === 200) {
                        uni.showToast({ title: '删除成功', icon: 'success' });
                        loadRules();
                    } else {
                        uni.showToast({ title: result.msg || '删除失败', icon: 'none' });
                    }
                } catch (error) {
                    uni.showToast({ title: '删除失败', icon: 'none' });
                }
            }
        }
    });
};

// 添加规则
const addRule = () => {
    uni.navigateTo({
        url: '/pages/admin/rule/add'
    });
};
</script>

<style lang="scss" scoped>
.container {
    min-height: 100vh;
    background-color: #f5f7fa;
    padding-bottom: 140rpx;
}

.rule-list {
    padding: 20rpx;
}

.rule-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #fff;
    border-radius: 16rpx;
    padding: 24rpx;
    margin-bottom: 16rpx;
    box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.rule-left {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8rpx;
}

.rule-row {
    display: flex;
    align-items: center;
    
    .label {
        font-size: 24rpx;
        color: #999;
        margin-right: 8rpx;
        flex-shrink: 0;
    }
    
    .value {
        font-size: 24rpx;
        color: #333;
    }
}

.rule-right {
    display: flex;
    flex-direction: column;
    gap: 12rpx;
    flex-shrink: 0;
    margin-left: 20rpx;
}

.action-btn {
    padding: 12rpx 28rpx;
    border-radius: 8rpx;
    font-size: 24rpx;
    text-align: center;
    min-width: 80rpx;
    
    &.edit {
        background: #e6f7ff;
        color: #1890ff;
    }
    
    &.delete {
        background: #fff1f0;
        color: #ee0a24;
    }
    
    &:active {
        opacity: 0.7;
    }
}

.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 80rpx 40rpx;
    background: #fff;
    border-radius: 16rpx;
    
    .empty-icon {
        font-size: 80rpx;
        margin-bottom: 20rpx;
    }
    
    .empty-text {
        font-size: 26rpx;
        color: #999;
    }
}

.add-btn-wrapper {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 20rpx;
    padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
    background: #fff;
    border-top: 1rpx solid #f0f0f0;
}

.add-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10rpx;
    height: 88rpx;
    background: linear-gradient(135deg, #3c9cff 0%, #2b85e4 100%);
    border-radius: 44rpx;
    color: #fff;
    font-size: 30rpx;
    font-weight: 500;
    
    &:active {
        opacity: 0.8;
    }
}
</style>
