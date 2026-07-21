<!-- 收款账号列表页面 -->
<template>
	<view class="container">
		<uv-navbar title="收款账号" :placeholder="true" rightIcon="plus" @leftClick="goBack" @rightClick="goToCreate" />

		<scroll-view scroll-y class="scroll-container">
			<view v-if="accountList.length === 0" class="empty-box">
				<uv-empty mode="data" textSize="32" iconSize="150" />
			</view>

			<view v-else class="account-list">
				<view v-for="item in accountList" :key="item.id" class="account-card" @click="handleEdit(item)">
					<view class="card-header">
						<view class="type-tag" :class="item.type === 'personal' ? 'personal' : 'company'">
							{{ item.type === 'personal' ? '个人' : '企业' }}
						</view>
						<view class="status-tag" :class="item.choose === '1' ? 'active' : ''">
							{{ item.choose === '1' ? '已选择' : '未选择' }}
						</view>
					</view>
					<view class="card-body">
						<view class="info-row">
							<text class="info-label">账号</text>
							<text class="info-value">{{ item.account || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">姓名</text>
							<text class="info-value">{{ item.name || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">手机号</text>
							<text class="info-value">{{ item.phone || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">店名</text>
							<text class="info-value">{{ item.storename || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">汇付编码</text>
							<text class="info-value">{{ item.code || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">经营地址</text>
							<text class="info-value">{{ item.area || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">使用场景</text>
							<text class="info-value">{{ item.remarks || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">多方分账</text>
							<text class="info-value">{{ item.sharing || '-' }}</text>
						</view>
						<view class="info-row">
							<text class="info-label">分账状态</text>
							<text class="info-value">{{ item.share === '1' ? '开启' : '关闭' }}</text>
						</view>
						<view v-if="item.share === '1'" class="info-row">
							<text class="info-label">分账比例</text>
							<text class="info-value">{{ item.rate || 0 }}%</text>
						</view>
						<view class="info-row">
							<text class="info-label">创建时间</text>
							<text class="info-value">{{ formatDateTime(item.createdAt) }}</text>
						</view>
					</view>
					<view class="card-footer">
						<view class="action-btn" v-if="item.choose === '0'" @click.stop="handleChoose(item)">
							<uv-icon name="checkmark-circle" size="32rpx" color="#10b981" />
							<text class="active">选择</text>
						</view>
						<view class="action-btn" @click.stop="handleEdit(item)">
							<uv-icon name="edit-pen" size="32rpx" color="#2979ff" />
							<text class="primary">编辑</text>
						</view>
						<view class="action-btn" @click.stop="handleDelete(item)">
							<uv-icon name="trash" size="32rpx" color="#fa3534" />
							<text class="danger">删除</text>
						</view>
					</view>
				</view>
			</view>
		</scroll-view>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'

// 页面加载时获取账号列表
onShow(() => {
	loadAccountList()
})

// 收款账号列表数据
const accountList = ref([])
// 总条数
const total = ref(0)

// 获取收款账号列表
const loadAccountList = async () => {
	uni.showLoading({ title: '加载中...' })
	const merch = uni.getStorageSync('merch') || {};
	const merchsId = merch.id || merch.merchsId || '';
	try {
		const res = await uni.$uv.http.get('/huifu/list', {
			params: {
				merchs_id: merchsId,
			},
			custom: { auth: true }
		})
		if (res.code === 200 && res.data.total > 0) {
			accountList.value = res.data.list || []
			total.value = res.data.total || 0
		} else {
			accountList.value = []
			total.value = 0
			uni.showToast({ title: '暂无收款账号列表', icon: 'none' })
		}
	} catch (e) {
		console.error('加载失败', e)
	} finally {
		uni.hideLoading()
	}
}

// 选择默认收款账号
const handleChoose = async (item) => {
	uni.showModal({
		title: '确认选择',
		content: '确定选择该收款账号作为默认收款账号吗？',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '设置中...' })
					const result = await uni.$uv.http.put('/huifu/choose/' + item.id, {
						custom: { auth: true }
					})
					uni.hideLoading()
					if (result.code === 200) {
						uni.showToast({ title: '设置成功', icon: 'success' })
						loadAccountList()
					} else {
						uni.showToast({ title: result.msg || '设置失败', icon: 'none' })
					}
				} catch (e) {
					uni.hideLoading()
					uni.showToast({ title: '设置失败', icon: 'none' })
				}
			}
		}
	})
}

// 删除收款账号
const handleDelete = (item) => {
	uni.showModal({
		title: '确认删除',
		content: '确定删除该收款账号吗？',
		confirmColor: '#fa3534',
		success: async (res) => {
			if (res.confirm) {
				try {
					uni.showLoading({ title: '删除中...' })
					const result = await uni.$uv.http.delete('/huifu/' + item.id, {
						custom: { auth: true }
					})
					uni.hideLoading()
					if (result.code === 200) {
						uni.showToast({ title: '删除成功', icon: 'success' })
						loadAccountList()
					} else {
						uni.showToast({ title: result.msg || '删除失败', icon: 'none' })
					}
				} catch (e) {
					uni.hideLoading()
					uni.showToast({ title: '删除失败', icon: 'none' })
				}
			}
		}
	})
}

// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/profile/index'
	});
}

// 跳转到添加账号页面
const goToCreate = () => {
	uni.navigateTo({ url: '/pages/admin/huifu/create' })
}

// 跳转到编辑账号页面
const handleEdit = (item) => {
	uni.navigateTo({ url: '/pages/admin/huifu/edit?id=' + item.id })
}

// 格式化时间
const formatDateTime = (dateStr) => {
	if (!dateStr) return '-'
	try {
		const date = new Date(dateStr)
		if (isNaN(date.getTime())) return '-'
		const year = date.getFullYear()
		const month = String(date.getMonth() + 1).padStart(2, '0')
		const day = String(date.getDate()).padStart(2, '0')
		const hours = String(date.getHours()).padStart(2, '0')
		const minutes = String(date.getMinutes()).padStart(2, '0')
		const seconds = String(date.getSeconds()).padStart(2, '0')
		return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
	} catch (e) {
		return '-'
	}
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.scroll-container {
	height: calc(100vh - 88rpx - 100rpx);
	width: calc(100% - 48rpx);
	padding: 24rpx;
}

.empty-box {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	padding: 120rpx 0;
}

.empty-text {
	font-size: 28rpx;
	color: #999;
	margin-top: 24rpx;
}

.add-btn {
	margin-top: 32rpx;
	padding: 16rpx 48rpx;
	background: linear-gradient(135deg, #10b981, #059669);
	border-radius: 40rpx;
	font-size: 28rpx;
	color: #fff;
}

.account-list {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
}

.account-card {
	background: #fff;
	border-radius: 16rpx;
	padding: 24rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.card-header {
	display: flex;
	justify-content: space-between;
	margin-bottom: 16rpx;
}

.type-tag {
	padding: 6rpx 16rpx;
	border-radius: 8rpx;
	font-size: 22rpx;

	&.personal {
		background: rgba(41, 121, 255, 0.1);
		color: #2979ff;
	}

	&.company {
		background: rgba(168, 85, 247, 0.1);
		color: #a855f7;
	}
}

.status-tag {
	padding: 6rpx 16rpx;
	border-radius: 8rpx;
	font-size: 22rpx;
	background: #f5f5f5;
	color: #999;

	&.active {
		background: rgba(16, 185, 129, 0.1);
		color: #10b981;
	}
}

.card-body {
	padding: 16rpx 0;
	border-bottom: 1rpx solid #f0f0f0;
}

.info-row {
	display: flex;
	justify-content: space-between;
	padding: 8rpx 0;
}

.info-label {
	font-size: 26rpx;
	color: #999;
}

.info-value {
	font-size: 26rpx;
	color: #333;
	font-weight: 500;
}

.card-footer {
	display: flex;
	justify-content: flex-end;
	gap: 32rpx;
	padding-top: 16rpx;
}

.action-btn {
	display: flex;
	align-items: center;
	gap: 8rpx;
	font-size: 24rpx;
	color: #999;

	text.active {
		color: #10b981;
	}

	text.primary {
		color: #2979ff;
	}

	text.danger {
		color: #fa3534;
	}
}

.bottom-btn {
	position: fixed;
	left: 0;
	right: 0;
	bottom: 0;
	height: 100rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	gap: 16rpx;
	background: linear-gradient(135deg, #10b981, #059669);
	box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.1);
}

.btn-text {
	font-size: 30rpx;
	color: #fff;
	font-weight: 500;
}
</style>