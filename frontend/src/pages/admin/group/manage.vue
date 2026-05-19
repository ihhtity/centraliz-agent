<template>
	<view class="container">
		<!-- 修改: 添加返回按钮 -->
		<uv-navbar :title="t('admin.group.title')" :placeholder="true" @leftClick="goBack">
			<template #right>
				<uv-icon name="plus" @click="addGroup"></uv-icon>
			</template>
		</uv-navbar>

		<uv-list>
			<uv-list-item 
				v-for="item in groupList" 
				:key="item.id" 
				:title="item.name" 
				:note="t('admin.group.groupList').replace(t('admin.group.groupList'), '') + t('admin.group.belongGroup') + ': ' + item.count"
				@click="editGroup(item)"
			>
				<template #right>
					<uv-button size="mini" type="error" @click.stop="deleteGroup(item.id)">{{ t('common.delete') }}</uv-button>
				</template>
			</uv-list-item>
		</uv-list>
	</view>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

// 新增: 返回上一页方法
const goBack = () => {
	uni.navigateBack();
};

const groupList = ref([
	{ id: 1, name: '默认分组', count: 5 },
	{ id: 2, name: 'VIP分组', count: 2 }
]);

const addGroup = () => {
	uni.showToast({ title: t('common.loading'), icon: 'none' });
};

const editGroup = (item) => {
	uni.showToast({ title: t('admin.group.title'), icon: 'none' });
};

const deleteGroup = (id) => {
	uni.showModal({
		title: t('common.confirm'),
		content: t('common.delete'), // 或者专门的 confirm delete key
		success: function (res) {
			if (res.confirm) {
				console.log('删除分组', id);
			}
		}
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f5f5;
}
</style>