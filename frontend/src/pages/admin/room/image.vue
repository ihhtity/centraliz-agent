<template>
	<view class="container">
		<uv-navbar title="房间图片" :placeholder="true" :leftIcon="'arrow-left'" @leftClick="goBack">
			<template #right>
				<uv-icon name="plus-circle" size="25" @click="showAddModal" />
			</template>
		</uv-navbar>
		
		<view class="content">
			<view class="image-grid" v-if="images.length > 0">
				<view v-for="(image, index) in images" :key="index" class="image-item">
					<image :src="image.image" mode="aspectFit" class="image-preview" @click="previewImage(image.image)" />
					<view class="image-overlay">
						<view class="image-name">{{ image.name }}</view>
						<view class="image-actions">
							<view class="action-icon edit" @click="editImage(image)">
								<uv-icon name="edit-pen" color="#3c9cff" size="20" />
							</view>
							<view class="action-icon delete" @click="deleteImage(image)">
								<uv-icon name="trash" color="#fa3534" size="20" />
							</view>
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
					<text class="modal-title">{{ isEdit ? '编辑图片' : '添加图片' }}</text>
				</view>
				<view class="modal-body">
					<view class="form-item">
						<text class="form-label">图片名称 <text class="required">*</text></text>
						<view class="input-box">
							<uv-input v-model="formData.name" placeholder="请输入图片名称" class="input-field" border="none" />
						</view>
					</view>
					<view class="form-item" v-if="isEdit">
						<text class="form-label">图片URL</text>
						<view class="input-box">
							<uv-input v-model="formData.image" placeholder="请输入图片URL" class="input-field" :clearable="true" border="none" />
						</view>
					</view>
					<view class="form-item" v-else>
						<text class="form-label">图片URL <text class="required">*</text></text>
						<view class="input-box">
							<uv-input v-model="formData.image" placeholder="请输入图片URL" class="input-field" :clearable="true" border="none" />
						</view>
					</view>
					<view class="form-item">
						<text class="form-label">当前图片</text>
						<image :src="formData.image" mode="aspectFit" class="current-image" />
					</view>
				</view>
				<view class="modal-footer">
					<view class="button button-cancel" @click="closeModal">取消</view>
					<view class="button button-confirm" @click="submitForm">{{ isEdit ? '保存' : '确认' }}</view>
				</view>
			</view>
		</uv-popup>
	</view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app';

const images = ref([])
const popupRef = ref(null)
const isEdit = ref(false)
const editId = ref(null)
const formData = ref({
	name: '',
	image: ''
})

onLoad(() => {
	fetchImages()
})

// 获取图片列表
const fetchImages = () => {
	uni.showLoading({ title: '加载中...' })
	uni.$uv.http.get('/room/image/list', {
		custom: { auth: true }
	}).then((res) => {
		uni.hideLoading()
		if (res.code === 200) {
			images.value = res.data.list || []
		}
	}).catch((err) => {
		uni.hideLoading()
		uni.showToast({ title: err.message || '加载失败', icon: 'none' })
	})
}
// 提交表单
const submitForm = () => {
	if (!formData.value.name.trim()) {
		uni.showToast({ title: '请输入图片名称', icon: 'none' })
		return
	}
	
	const requestData = {
		name: formData.value.name.trim(),
		image: formData.value.image
	}
	
	if (!formData.value.image) {
		uni.showToast({ title: '请输入图片URL', icon: 'none' })
		return
	}
	// console.log(requestData);return
	
	uni.showLoading({ title: '提交中...' })
	
	const request = isEdit.value 
		? uni.$uv.http.put(`/room/image/${editId.value}`, requestData, { custom: { auth: true } })
		: uni.$uv.http.post('/room/image', requestData, { custom: { auth: true } })
	
	request.then((res) => {
		uni.hideLoading()
		if (res.code === 200) {
			uni.showToast({ title: isEdit.value ? '修改成功' : '添加成功', icon: 'success' })
			closeModal()
			fetchImages()
		} else {
			uni.showToast({ title: res.msg || '操作失败', icon: 'none' })
		}
	}).catch((err) => {
		uni.hideLoading()
		uni.showToast({ title: err.message || '操作失败', icon: 'none' })
	})
}
// 预览图片
const previewImage = (url) => {
	uni.previewImage({
		current: url,
		urls: images.value.map(img => img.image)
	})
}
// 显示添加图片弹窗
const showAddModal = () => {
	isEdit.value = false
	formData.value = {
		name: '',
		image: ''
	}
	popupRef.value.open()
}
// 编辑图片
const editImage = (image) => {
	isEdit.value = true
	editId.value = image.id
	formData.value = {
		name: image.name,
		image: image.image
	}
	popupRef.value.open()
}
// 删除图片
const deleteImage = (image) => {
	uni.showModal({
		title: '确认删除',
		content: `确定要删除图片"${image.name}"吗？`,
		cancelText: '取消',
		confirmText: '确定',
		success: (res) => {
			if (res.confirm) {
				uni.$uv.http.delete(`/room/image/${image.id}`, {
					custom: { auth: true }
				}).then((res) => {
					if (res.code === 200) {
						uni.showToast({ title: '删除成功', icon: 'success' })
						fetchImages()
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
// 关闭弹窗
const closeModal = () => {
	popupRef.value.close()
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

.action-bar {
	display: flex;
	gap: 16rpx;
	margin-bottom: 24rpx;
	
	.action-btn {
		flex: 1;
		display: flex;
		flex-direction: row;
		justify-content: center;
		padding: 16rpx;
		border-radius: 8rpx;
		font-size: 24rpx;
		color: #fff;
		text-align: center;
		cursor: pointer;
		
		.uv-icon {
			margin-right: 8rpx;
		}
		
		.btn-text {
			font-size: 26rpx;
			font-weight: 500;
		}
	}
	
	.add-btn {
		background: linear-gradient(135deg, #3c9cff, #2b85e4);
	}
}

.image-grid {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 16rpx;
}

.image-item {
	position: relative;
	background: #fff;
	border-radius: 12rpx;
	overflow: hidden;
	aspect-ratio: 1;
}

.image-preview {
	width: 100%;
	height: 100%;
}

.image-overlay {
	position: absolute;
	bottom: 0;
	left: 0;
	right: 0;
	background: linear-gradient(transparent, rgba(0,0,0,0.6));
	padding: 16rpx;
	display: flex;
	justify-content: space-between;
	align-items: flex-end;
}

.image-name {
	font-size: 22rpx;
	color: #fff;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
	flex: 1;
	margin-right: 12rpx;
}

.image-actions {
	display: flex;
	gap: 12rpx;
}

.action-icon {
	width: 48rpx;
	height: 48rpx;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	background-color: rgba(255,255,255,0.2);
	
	&:active {
		background-color: rgba(255,255,255,0.4);
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
	
	.upload-area {
		background: #fff;
		border-radius: 16rpx;
		padding: 40rpx;
		height: 200rpx;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		border: 2rpx dashed #e8e8e8;
		gap: 16rpx;
		
		&:active {
			border-color: #3c9cff;
			background: rgba(60, 156, 255, 0.05);
		}
	}
	
	.upload-text {
		font-size: 26rpx;
		color: #999;
	}
	
	.current-image {
		width: 100%;
		height: 200rpx;
		border-radius: 12rpx;
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