<!-- 规则管理页面 -->
<template>
	<view class="container">
		<uv-navbar :title="t('admin.rule.title')" :placeholder="true" @leftClick="goBack">
			<template #right>
				<uv-icon name="plus" @click="addRule"></uv-icon>
			</template>
		</uv-navbar>

		<uv-list>
			<uv-list-item 
				v-for="item in ruleList" 
				:key="item.id" 
				:title="item.name" 
				:note="`${t('admin.rule.ruleType')}: ${item.type}`"
				@click="viewRule(item)"
			>
				<template #right>
					<uv-switch v-model="item.enabled" @change="toggleRule(item)"></uv-switch>
				</template>
			</uv-list-item>
		</uv-list>
	</view>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const goBack = () => {
	uni.navigateBack();
};

const ruleList = ref([
	{ id: 1, name: '超时自动断电', type: '时间控制', enabled: true },
	{ id: 2, name: '过载保护', type: '安全保护', enabled: true }
]);

const addRule = () => {
	uni.showToast({ title: t('common.loading'), icon: 'none' });
};

const toggleRule = (item) => {
	uni.showToast({ 
		title: item.enabled ? t('admin.rule.enable') : t('admin.rule.disable'), 
		icon: 'none' 
	});
};
</script>

<style lang="scss" scoped>
.container {
	min-height: 100vh;
	background-color: #f5f5f5;
}
</style>