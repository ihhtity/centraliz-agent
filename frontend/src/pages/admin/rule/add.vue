<template>
    <view class="container">
        <uv-navbar :title="'规则添加'" :placeholder="true" @leftClick="goBack" />
        
        <view class="form-wrapper">
            <!-- 规则名称 -->
            <view class="form-item">
                <text class="form-label">规则名称</text>
                <input class="form-input" v-model="form.name" placeholder="请输入规则名称" />
            </view>

            <!-- 规则类型 -->
            <view class="form-item">
                <text class="form-label">规则类型</text>
                <view class="radio-group">
                    <view 
                        class="radio-item" 
                        :class="{ active: form.type === 'free' }" 
                        @click="selectType('free')"
                    >
                        <text>免费模式</text>
                    </view>
                    <view 
                        class="radio-item" 
                        :class="{ active: form.type === 'charge' }" 
                        @click="selectType('charge')"
                    >
                        <text>收费模式</text>
                    </view>
                </view>
            </view>

            <!-- 模式类型 -->
            <view class="form-item">
                <text class="form-label">模式类型</text>
                <view class="radio-group">
                    <view 
                        v-for="mode in modeOptions" 
                        :key="mode.value"
                        class="radio-item" 
                        :class="{ active: form.mode === mode.value, disabled: mode.type !== form.type }" 
                        @click="selectMode(mode.value)"
                    >
                        <text>{{ mode.label }}</text>
                    </view>
                </view>
            </view>

            <!-- 价格（收费模式） -->
            <view v-if="form.type === 'charge' && form.mode !== 'pay_time'" class="form-item">
                <text class="form-label">单价（元）</text>
                <input class="form-input" v-model="form.price" placeholder="请输入单价" type="digit" />
            </view>

            <!-- 按时付费时间单位 -->
            <view v-if="form.mode == 'pay_hourly'" class="form-item">
                <text class="form-label">计费时间单位</text>
                <view class="unit-selector">
                    <view 
                        v-for="unit in timeUnits" 
                        :key="unit.value"
                        class="unit-item"
                        :class="{ active: form.durationUnit === unit.value }"
                        @click="form.durationUnit = unit.value"
                    >
                        <text>{{ unit.label }}</text>
                    </view>
                </view>
            </view>

            <!-- 押金 -->
            <view v-if="showDeposit" class="form-item">
                <text class="form-label">押金（元）</text>
                <input class="form-input" v-model="form.deposit" placeholder="请输入押金" type="digit" />
            </view>

            <!-- 免费时间（免费模式一存一取和收费模式） -->
            <view v-if="showFreeTime" class="form-item">
                <text class="form-label">免费时间（分钟）</text>
                <input class="form-input" v-model="form.freeTime" placeholder="0" type="number" />
                <text v-if="form.type === 'charge'" class="form-hint">在此时间内可以临时开锁和结束订单不收费</text>
                <text v-else class="form-hint">在此时间内可以临时开锁</text>
            </view>

            <!-- 自动结束时间（除了预付费模式） -->
            <view v-if="showAutoEnd" class="form-item">
                <text class="form-label">自动结束时间（分钟）</text>
                <input class="form-input" v-model="form.autoEndTime" placeholder="0" type="number" />
                <text class="form-hint">设置为0表示不自动结束</text>
            </view>

            <!-- 预付费功能开关 -->
            <view v-if="form.mode === 'pay_time'" class="form-item">
                <text class="form-label">功能开关</text>
                <view class="switch-group">
                    <view class="switch-item">
                        <text class="switch-label">自动退款</text>
                        <text class="switch-hint">提前结束订单按整小时退还剩余费用</text>
                        <switch :checked="form.autoRefund" @change="form.autoRefund = $event.detail.value" color="#1989fa" />
                    </view>
                    <view class="switch-item">
                        <text class="switch-label">手动续费</text>
                        <text class="switch-hint">订单有效期内可手动续费，续费金额=续费时长×单价</text>
                        <switch :checked="form.manualRenew" @change="form.manualRenew = $event.detail.value" color="#1989fa" />
                    </view>
                </view>
            </view>

            <!-- 预付费的时间选项 -->
            <view v-if="form.mode === 'pay_time'" class="form-item">
                <text class="form-label">时间选项（最多6个）</text>
                <view class="time-options">
                    <view v-for="(option, index) in form.timeOptions" :key="index" class="time-option-item">
                        <view class="option-header">
                            <text class="option-title">选项{{ index + 1 }}</text>
                            <view 
                                class="option-delete" 
                                v-if="form.timeOptions.length > 1" 
                                @click="removeTimeOption(index)"
                            >
                                <text>删除</text>
                            </view>
                        </view>
                        <view class="option-fields">
                            <view class="field-row">
                                <input class="option-input" v-model="option.title" placeholder="套餐标题" />
                            </view>
                            <view class="field-row">
                                <input class="option-input" v-model="option.duration" placeholder="时长" type="number" />
                                <view class="unit-selector-small">
                                    <view 
                                        v-for="unit in timeUnits" 
                                        :key="unit.value"
                                        class="unit-item-small"
                                        :class="{ active: option.durationUnit === unit.value }"
                                        @click="option.durationUnit = unit.value"
                                    >
                                        <text>{{ unit.label }}</text>
                                    </view>
                                </view>
                            </view>
                            <view class="field-row">
                                <input class="option-input" v-model="option.price" placeholder="单价（元）" type="digit" />
                            </view>
                            <view class="field-row">
                                <input class="option-input" v-model="option.discount" placeholder="优惠金额（元）" type="digit" />
                            </view>
                            <view class="field-row">
                                <text class="total-price">总价：{{ calculateTotal(option) }}元</text>
                            </view>
                        </view>
                    </view>
                </view>
                <view v-if="form.timeOptions.length < 6" class="add-option-btn" @click="addTimeOption">
                    <text>+ 添加选项</text>
                </view>
            </view>

            <!-- 规则描述 -->
            <view class="form-item">
                <text class="form-label">规则描述</text>
                <textarea 
                    class="form-textarea" 
                    v-model="form.description" 
                    placeholder="请输入规则描述"
                />
            </view>
        </view>

        <!-- 提交按钮 -->
        <view class="submit-btn-wrapper">
            <view class="submit-btn" @click="submitForm">
                <text>保存规则</text>
            </view>
        </view>
    </view>
</template>

<script setup>
import { reactive, computed } from 'vue';

// 商家数据
const merch = uni.getStorageSync('merch') || {};

// 表单数据
const form = reactive({
    name: '',
    type: 'free',
    mode: 'single',
    price: '1',
    deposit: '0',
    autoEndTime: '0',
    freeTime: '0',
    durationUnit: 'hour',
    autoRefund: false,
    manualRenew: false,
    description: '免费模式:您可以免费开锁，无需支付任何费用',
    timeOptions: [
        { title: '', duration: '', durationUnit: 'hour', price: '', discount: '0' }
    ]
});

// 模式选项
const modeOptions = [
    { value: 'single', label: '单次开锁', type: 'free' },
    { value: 'deposit', label: '一存一取', type: 'free' },
    { value: 'pay_single', label: '单次付费', type: 'charge' },
    { value: 'pay_deposit', label: '先存后取', type: 'charge' },
    { value: 'pay_hourly', label: '按时付费', type: 'charge' },
    { value: 'pay_time', label: '预付费', type: 'charge' }
];

// 时间单位选项
const timeUnits = [
    { value: 'minute', label: '分' },
    { value: 'hour', label: '时' },
    { value: 'day', label: '天' },
    { value: 'month', label: '月' }
];

// 判断模式是否可用
const isModeAvailable = (mode) => {
    const option = modeOptions.find(m => m.value === mode);
    if (!option) return false;
    return option.type === form.type;
};

// 选择规则类型
const selectType = (type) => {
    form.type = type;
    if (type === 'free') {
        form.mode = 'single';
        selectMode('single');
    } else if (type === 'charge') {
        form.mode = 'pay_single';
        selectMode('pay_single');
    }
};

// 选择模式类型
const selectMode = (mode) => {
    if (isModeAvailable(mode)) {
        form.mode = mode;
        switch (mode) {
            case 'single':
                form.description = '免费模式:您可以免费开锁，无需支付任何费用';
                break;
            case 'deposit':
                form.description = '免费模式:您可以免费存取，无需支付任何费用';
                break;
            case 'pay_single':
                form.description = '单次付费模式:您需要在开锁时支付费用，开锁后即可使用';
                break;
            case 'pay_deposit':
                form.description = '先存后取模式:您可以先存放物品，取走时需要支付费用才能打开。如果涉及押金可以手动申请退款';
                break;
            case 'pay_hourly':
                form.description = '按时付费模式:您可以先存放物品，取走时需要支付所有时长对应的费用才能打开。如果涉及押金可以手动申请退款';
                break;
            case 'pay_time':
                form.description = '预付费模式:您需要在使用时支付费用，使用后即可使用。如果涉及押金可以手动申请退款';
                break;
        }
    }
};

// 是否显示押金
const showDeposit = computed(() => {
    return form.type === 'charge';
});

// 是否显示自动结束时间
const showAutoEnd = computed(() => {
    return ['deposit', 'pay_deposit', 'pay_hourly'].includes(form.mode);
});

// 是否显示免费时间
const showFreeTime = computed(() => {
    return !['single', 'pay_single'].includes(form.mode);
});

// 计算总价
const calculateTotal = (option) => {
    const duration = parseInt(option.duration) || 0;
    const price = parseFloat(option.price) || 0;
    const discount = parseFloat(option.discount) || 0;
    return Math.max(0, price * duration - discount).toFixed(2);
};

// 添加时间选项
const addTimeOption = () => {
    if (form.timeOptions.length < 6) {
        form.timeOptions.push({ title: '', duration: '', durationUnit: 'hour', price: '', discount: '0' });
    }
};

// 删除时间选项
const removeTimeOption = (index) => {
    if (form.timeOptions.length > 1) {
        form.timeOptions.splice(index, 1);
    }
};

// 提交表单
const submitForm = async () => {
    if (!form.name) {
        uni.showToast({ title: '请输入规则名称', icon: 'none' });
        return;
    }
    // 检查是否有商家ID
    if (!merch.id) {
        uni.showToast({ title: '请先登录', icon: 'none' });
        return;
    }

    try {
        const data = {
            merchsId: merch.id,
            name: form.name,
            type: form.type,
            mode: form.mode
        };

        // 添加可选字段
        if (form.type === 'charge' && form.price && form.mode !== 'pay_time') {
            data.price = parseFloat(form.price);
        }
        if (form.deposit) {
            data.deposit = parseFloat(form.deposit);
        }
        if (form.mode === 'pay_hourly') {
            data.durationUnit = form.durationUnit || 'hour';
        }
        if (form.autoEndTime) {
            data.autoEndTime = parseInt(form.autoEndTime);
        }
        if (form.freeTime) {
            data.freeTime = parseInt(form.freeTime);
        }

        // 预付费的时间选项
        if (form.mode === 'pay_time') {
            const timeOptions = form.timeOptions.map(opt => ({
                title: opt.title,
                duration: parseInt(opt.duration) || 0,
                durationUnit: opt.durationUnit || 'hour',
                price: parseFloat(opt.price) || 0,
                discount: parseFloat(opt.discount) || 0
            }));
            data.timeOptions = JSON.stringify(timeOptions);
            data.autoRefund = form.autoRefund;
            data.manualRenew = form.manualRenew;
        }

        // 发送描述字段
        if (form.description.trim()) {
            data.description = form.description;
        }

        const res = await uni.$uv.http.post('/rule', data, {
            custom: { auth: true }
        });

        if (res.code === 200) {
            uni.showToast({ title: '创建成功', icon: 'success' });
            setTimeout(() => {
                uni.navigateBack();
            }, 1500);
        } else {
            uni.showToast({ title: res.msg || '创建失败', icon: 'none' });
        }
    } catch (error) {
        console.error('创建规则失败:', error);
        uni.showToast({ title: '创建失败', icon: 'none' });
    }
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
    padding-bottom: 140rpx;
}

.form-wrapper {
    padding: 30rpx;
}

.form-item {
    background-color: #fff;
    padding: 30rpx;
    margin-bottom: 20rpx;
    border-radius: 16rpx;
}

.form-label {
    font-size: 28rpx;
    color: #333;
    margin-bottom: 20rpx;
    display: block;
}

.form-input {
    width: 100%;
    height: 80rpx;
    border: 1rpx solid #e8e8e8;
    border-radius: 8rpx;
    padding: 0 20rpx;
    font-size: 28rpx;
    box-sizing: border-box;
}

.form-textarea {
    width: 100%;
    height: 200rpx;
    border: 1rpx solid #e8e8e8;
    border-radius: 8rpx;
    padding: 20rpx;
    font-size: 28rpx;
    box-sizing: border-box;
}

.form-hint {
    font-size: 24rpx;
    color: #999;
    margin-top: 10rpx;
    display: block;
}

.radio-group {
    display: flex;
    flex-wrap: wrap;
    gap: 20rpx;
}

.radio-item {
    padding: 20rpx 40rpx;
    background-color: #f5f5f5;
    border-radius: 8rpx;
    font-size: 28rpx;
    color: #666;
    transition: all 0.3s;

    &.active {
        background-color: #1989fa;
        color: #fff;
    }

    &.disabled {
        opacity: 0.4;
        pointer-events: none;
    }
}

.unit-selector {
    display: flex;
    gap: 20rpx;
    flex-wrap: wrap;
}

.unit-item {
    padding: 15rpx 30rpx;
    background-color: #f5f5f5;
    border-radius: 8rpx;
    font-size: 26rpx;
    color: #666;

    &.active {
        background-color: #1989fa;
        color: #fff;
    }
}

.unit-selector-small {
    display: flex;
    gap: 8rpx;
}

.unit-item-small {
    padding: 10rpx 20rpx;
    background-color: #f5f5f5;
    border-radius: 6rpx;
    font-size: 22rpx;
    color: #666;

    &.active {
        background-color: #1989fa;
        color: #fff;
    }
}

.time-options {
    margin-bottom: 20rpx;
}

.time-option-item {
    border: 1rpx solid #e8e8e8;
    border-radius: 8rpx;
    padding: 20rpx;
    margin-bottom: 15rpx;
}

.option-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15rpx;
}

.option-title {
    font-size: 26rpx;
    color: #333;
    font-weight: 500;
}

.option-delete {
    font-size: 24rpx;
    color: #f56c6c;
}

.option-fields {
    display: flex;
    flex-direction: column;
    gap: 15rpx;
}

.field-row {
    display: flex;
    align-items: center;
    gap: 15rpx;
}

.option-input {
    flex: 1;
    height: 70rpx;
    border: 1rpx solid #e8e8e8;
    border-radius: 6rpx;
    padding: 0 15rpx;
    font-size: 26rpx;
}

.total-price {
    font-size: 26rpx;
    color: #1989fa;
    font-weight: 500;
}

.add-option-btn {
    padding: 20rpx;
    border: 1rpx dashed #d9d9d9;
    border-radius: 8rpx;
    text-align: center;
    font-size: 28rpx;
    color: #1989fa;
}

.switch-group {
    display: flex;
    flex-direction: column;
    gap: 20rpx;
}

.switch-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15rpx 0;
}

.switch-label {
    font-size: 28rpx;
    color: #333;
    font-weight: 500;
}

.switch-hint {
    font-size: 24rpx;
    color: #999;
    margin-left: 15rpx;
    flex: 1;
}

.submit-btn-wrapper {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    padding: 20rpx 30rpx;
    padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
    background-color: #fff;
    box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.submit-btn {
    height: 88rpx;
    background: linear-gradient(135deg, #1989fa 0%, #40a9ff 100%);
    border-radius: 44rpx;
    display: flex;
    align-items: center;
    justify-content: center;

    text {
        font-size: 32rpx;
        color: #fff;
        font-weight: 500;
    }
}
</style>
