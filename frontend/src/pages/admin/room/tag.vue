<template>
	<view class="container">
		<uv-navbar title="房间标签" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack">
			<template #right>
				<uv-icon name="plus-circle" size="25" @click="showAddModal" />
			</template>
		</uv-navbar>
		
		<view class="content">
			<view class="search-bar">
				<uv-icon name="search" color="#999" size="24" />
				<input v-model="searchQuery" type="text" placeholder="搜索标签..." />
			</view>
			
			<view class="tag-list" v-if="tags.length > 0">
				<view v-for="(tag, index) in filteredTags" :key="index" class="tag-item">
					<view class="tag-info">
						<text class="tag-name">{{ tag.name }}</text>
						<text class="tag-time">{{ formatTime(tag.createdAt) }}</text>
					</view>
					<view class="tag-actions">
						<view class="action-icon edit" @click="editTag(tag)">
							<uv-icon name="edit-pen" color="#3c9cff" size="24" />
						</view>
						<view class="action-icon delete" @click="deleteTag(tag)">
							<uv-icon name="trash" color="#fa3534" size="24" />
						</view>
					</view>
				</view>
			</view>
			
			<view class="empty-state" v-else>
				<uv-empty mode="data" textSize="32" iconSize="150" />
			</view>
		</view>
		
		<uv-popup ref="popupRef" mode="bottom" closeable @close="closeModal" :safeAreaInsetBottom="true">
			<view class="modal-content">
				<view class="modal-header">
					<text class="modal-title">{{ isEdit ? '编辑标签' : '添加标签' }}</text>
				</view>
				<view class="modal-body">
					<view class="form-item">
						<text class="form-label">标签名称 <text class="required">*</text></text>
						<view class="input-box">
							<uv-input v-model="formData.name" placeholder="请输入标签名称" class="input-field" border="none" />
						</view>
					</view>
				</view>
				<view class="modal-footer">
					<view class="button button-cancel" @click="closeModal">取消</view>
					<view class="button button-confirm" @click="submitForm">{{ isEdit ? '保存修改' : '确认添加' }}</view>
				</view>
			</view>
		</uv-popup>
	</view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app';

const merch = uni.getStorageSync('merch')
const tags = ref([])
const searchQuery = ref('')
const popupRef = ref(null)
const isEdit = ref(false)
const formData = ref({ name: '' })
const editId = ref(null)

onLoad(() => {
	fetchTags()
})

// 获取房间标签列表
const fetchTags = () => {
	uni.showLoading({ title: '加载中...' })
	uni.$uv.http.get('/room/tag/list/'+merch.id, {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading()
		if (res.code === 200) {
			tags.value = res.data.list || []
		}
	}).catch((err) => {
		uni.hideLoading()
		uni.showToast({ title: err.message || '加载失败', icon: 'none' })
	})
}
// 删除房间标签
const deleteTag = (tag) => {
	uni.showModal({
		title: '确认删除',
		content: `确定要删除标签"${tag.name}"吗？`,
		cancelText: '取消',
		confirmText: '确定',
		success: (res) => {
			if (res.confirm) {
				uni.$uv.http.delete(`/room/tag/${tag.id}`, {
					custom: { auth: true }
				}).then((res) => {
					if (res.code === 200) {
						uni.showToast({ title: '删除成功', icon: 'success' })
						fetchTags()
					} else {
						uni.showToast({ title: res.msg || '删除失败', icon: 'none' })
					}
				}).catch((err) => {
					uni.showToast({ title: err.message || '删除失败', icon: 'none' })
				})
			}
		}
	})
}
// 提交表单
const submitForm = () => {
	if (!formData.value.name.trim()) {
		uni.showToast({ title: '请输入标签名称', icon: 'none' })
		return
	}
	
	const requestData = {
		merchsId: merch.id,
		name: formData.value.name.trim()
	}
	
	uni.showLoading({ title: '提交中...' })
	
	const request = isEdit.value 
		? uni.$uv.http.put(`/room/tag/${editId.value}`, requestData, { custom: { auth: true } })
		: uni.$uv.http.post('/room/tag', requestData, { custom: { auth: true } })
	
	request.then((res) => {
		uni.hideLoading()
		if (res.code === 200) {
			uni.showToast({ title: isEdit.value ? '修改成功' : '添加成功', icon: 'success' })
			closeModal()
			fetchTags()
		} else {
			uni.showToast({ title: res.msg || '操作失败', icon: 'none' })
		}
	}).catch((err) => {
		uni.hideLoading()
		uni.showToast({ title: err.message || '操作失败', icon: 'none' })
	})
}
// 过滤房间标签列表
const filteredTags = computed(() => {
	if (!searchQuery.value) return tags.value
	return tags.value.filter(tag => tag.name.includes(searchQuery.value))
})
// 显示添加标签弹窗
const showAddModal = () => {
	isEdit.value = false
	editId.value = null
	formData.value = { name: '' }
	popupRef.value.open()
}
// 显示编辑标签弹窗
const editTag = (tag) => {
	isEdit.value = true
	editId.value = tag.id
	formData.value = { name: tag.name }
	popupRef.value.open()
}
// 关闭弹窗
const closeModal = () => {
	popupRef.value.close()
}
// 格式化时间
const formatTime = (time) => {
	if (!time) return ''
	return time.replace('T', ' ').substring(0, 19)
}
// 返回上一页
const goBack = () => {
	uni.redirectTo({
		url: '/pages/admin/profile/index'
	})
}
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
	padding: 24rpx;
}

.content {
	padding-bottom: 100rpx;
}

.search-bar {
	position: relative;
	width: 100%;
	height: 80rpx;
	background: #fff;
	border-radius: 40rpx;
	box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
	overflow: hidden;
	display: flex;
	align-items: center;
	padding: 0 20rpx;
	margin-bottom: 24rpx;
	
	input {
		flex: 1;
		font-size: 28rpx;
		color: #333;
		outline: none;
		background: transparent;
	}
	
	.uv-icon {
		margin-right: 16rpx;
	}
}

.tag-list {
	background: #fff;
	border-radius: 16rpx;
	overflow: hidden;
}

.tag-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 24rpx;
	border-bottom: 1rpx solid #f0f0f0;
	
	&:last-child {
		border-bottom: none;
	}
	
	&:active {
		background-color: #f8f9fa;
	}
}

.tag-info {
	flex: 1;
	display: flex;
	flex-direction: column;
	gap: 8rpx;
}

.tag-name {
	font-size: 30rpx;
	color: #333;
	font-weight: 500;
}

.tag-time {
	font-size: 24rpx;
	color: #999;
}

.tag-actions {
	display: flex;
	gap: 24rpx;
}

.action-icon {
	width: 60rpx;
	height: 60rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	background-color: #f5f7fa;
	
	&:active {
		background-color: #e8e8e8;
	}
}

.empty-state {
	margin-top: 100rpx;
	text-align: center;
}

.modal-content {
	padding: 0;
	background: #fafafa;
	border-radius: 24rpx 24rpx 0 0;
	height: auto;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.modal-header {
	background: #fff;
	padding: 32rpx;
	border-bottom: 1rpx solid #f0f0f0;
	text-align: center;
	
	.modal-title {
		font-size: 32rpx;
		font-weight: 600;
		color: #1a1a1a;
	}
}

.modal-body {
	padding: 32rpx;
	display: flex;
	flex-direction: column;
	gap: 32rpx;
}

.form-item {
	display: flex;
	flex-direction: column;
	gap: 16rpx;
	
	.form-label {
		font-size: 28rpx;
		color: #333;
		font-weight: 500;
		
		.required {
			color: #ff4d4f;
		}
	}
	
	.input-box {
		background: #fff;
		border-radius: 16rpx;
		padding: 0 24rpx;
		height: 88rpx;
		display: flex;
		align-items: center;
		border: 1rpx solid #e8e8e8;
		
		&:focus-within {
			border-color: #3c9cff;
			box-shadow: 0 0 0 4rpx rgba(60, 156, 255, 0.08);
		}
	}
	
	.input-field {
		flex: 1;
		font-size: 28rpx;
	}
}

.modal-footer {
	background: #fff;
	padding: 24rpx 32rpx;
	padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
	display: flex;
	gap: 16rpx;
	border-top: 1rpx solid #f0f0f0;
	
	.button {
		flex: 1;
		height: 88rpx;
		border-radius: 16rpx;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 30rpx;
		font-weight: 500;
		
		&:active {
			transform: scale(0.98);
			opacity: 0.9;
		}
	}
	
	.button-cancel {
		background: #f5f5f5;
		color: #666;
	}
	
	.button-confirm {
		background: #3c9cff;
		color: #fff;
	}
}
</style>