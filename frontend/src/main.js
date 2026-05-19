import { createSSRApp } from "vue";
import App from "./App.vue";
// 新增: 导入 uv-ui 组件库
import uvUI from '@climblee/uv-ui';
// 新增: 导入封装好的 http 请求工具
import i18n from './locales/index';
// 新增: 导入封装好的 http 请求工具
import { Request } from './utils/request';

export function createApp() {
	const app = createSSRApp(App);

	// 使用 uvUI 组件库
	app.use(uvUI);
	// 引入国际化插件
	app.use(i18n);
	// 引入请求封装
	Request(app)

	return {
		app,
	};
}