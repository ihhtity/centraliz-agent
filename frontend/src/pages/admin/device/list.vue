<template>
	<view class="container">
		<uv-navbar :title="t('admin.device.title')" :placeholder="true" @leftClick="goBack">
			<template #right>
				<uv-icon name="plus" size="25" color="#3c9cff" @click="addDevice" />
			</template>
		</uv-navbar>

		<view class="list-wrapper">
			<uv-list>
				<uv-list-item 
					v-for="item in deviceList" 
					:key="item.id" 
					class="device-item"
					@click="editDevice(item)"
				>
					<template #header>
						<view class="device-info">
							<view class="device-icon">
								<uv-icon name="empty-favor" size="24" color="#3c9cff" />
							</view>
							<view class="device-text">
								<text class="name">{{ item.name }}</text>
								<text class="id">ID: {{ item.id }}</text>
							</view>
						</view>
					</template>
					
					<template #right>
						<view class="actions">
							<view :class="['status-dot', item.status === t('admin.device.online') ? 'online' : 'offline']"></view>
							<uv-button size="mini" type="error" shape="circle" icon="trash" @click.stop="deleteDevice(item.id)"></uv-button>
						</view>
					</template>
				</uv-list-item>
			</uv-list>
			
			<uv-empty v-if="deviceList.length === 0" mode="list" :text="t('admin.device.noDevice')"></uv-empty>
		</view>
	</view>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const deviceList = ref([
	{ id: 1, name: '集控柜A', status: t('admin.device.online') },
	{ id: 2, name: '集控柜B', status: t('admin.device.offline') }
]);

const goBack = () => {
	uni.navigateBack();
};

const addDevice = () => {
	uni.navigateTo({ url: '/pages/admin/device/edit?type=add' });
};

const editDevice = (item) => {
	uni.navigateTo({ url: `/pages/admin/device/edit?id=${item.id}` });
};

const deleteDevice = (id) => {
	uni.showModal({
		title: t('common.confirm'),
		content: t('admin.device.deleteConfirm'),
		success: function (res) {
			if (res.confirm) {
				// 调用删除接口
				console.log('删除设备', id);
			}
		}
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f7fa;
}

.list-wrapper {
	padding: 24rpx;
}

.device-item {
	background: #fff;
	margin-bottom: 24rpx;
	border-radius: 20rpx;
	padding: 30rpx;
	box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.03);
	
	.device-info {
		display: flex;
		align-items: center;
		
		.device-icon {
			width: 70rpx;
			height: 70rpx;
			background: rgba(60, 156, 255, 0.1);
			border-radius: 50%;
			display: flex;
			align-items: center;
			justify-content: center;
			margin-right: 20rpx;
		}
		
		.device-text {
			display: flex;
			flex-direction: column;
			
			.name {
				font-size: 30rpx;
				font-weight: bold;
				color: #333;
				margin-bottom: 6rpx;
			}
			
			.id {
				font-size: 22rpx;
				color: #999;
			}
		}
	}
	
	.actions {
		display: flex;
		align-items: center;
		gap: 20rpx;
		
		.status-dot {
			width: 16rpx;
			height: 16rpx;
			border-radius: 50%;
			
			&.online {
				background-color: #19be6b;
				box-shadow: 0 0 8rpx rgba(25, 190, 107, 0.4);
			}
			
			&.offline {
				background-color: #c0c4cc;
			}
		}
	}
}
</style>