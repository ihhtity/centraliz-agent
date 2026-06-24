export const Request = (app) => {
	// 初始化请求配置
	uni.$uv.http.setConfig((config) => {
		/* config 为默认全局配置*/
		// #ifdef H5
		// 开发环境通过vite代理访问，生产环境直接访问后端服务
		// 检测是否为开发环境（localhost或127.0.0.1）
		const isLocal = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
		if (isLocal) {
			config.baseURL = '/api/v1'; // 本地开发时代理到/v1前缀
		} else {
			// config.baseURL = 'https://bsldlock.bsldtech.cn/api/v1'; // 线上环境直接访问后端
			config.baseURL = '/api/v1';
		}
		// #endif

		// #ifdef MP-WEIXIN || APP-PLUS
		// config.baseURL = 'https://bsldlock.bsldtech.cn/api/v1';  // 小程序和APP直接访问后端
		config.baseURL = 'http://localhost:3300/api/v1'; // 本地开发时代理到/v1前缀
		// #endif

		//请求头
		config.header = {
			'X-Requested-With': 'XMLHttpRequest',
			'Content-Type': 'application/json'
		}

		return config
	})

	// 请求拦截
	uni.$uv.http.interceptors.request.use((config) => {
		// 初始化请求拦截器时，会执行此方法，此时data为undefined，赋予默认{}
		config.data = config.data || {}
		if (config?.custom?.auth || config?.data?.custom?.auth) {
			// 从localStorage获取token
			const token = uni.getStorageSync('token');
			
			if (token) {
				config.header.Authorization = `Bearer ${token}`;
			} else {
				// token不存在时跳转到登录页
				uni.removeStorageSync('token');
				uni.removeStorageSync('merch');
				uni.reLaunch({ url: '/pages/login/login' });
			}
		}
		return config
	}, config => { // 可使用async await 做异步操作
		return Promise.reject(config)
	})

	// 响应拦截
	uni.$uv.http.interceptors.response.use((response) => { 
		const data = response.data
		
		// 如果响应不是对象格式，尝试解析
		if (!data || typeof data !== 'object') {
			uni.showToast({ title: '响应格式错误', icon: 'none', duration: 3000 });
			return Promise.reject({ msg: '响应格式错误' });
		}

		// 处理业务错误（后端统一返回HTTP 200，通过code字段判断）
		// if (data.code !== 200 && data.code !== 0) { 
		// 	// 如果没有显式定义custom的toast参数为false的话，默认对报错进行toast弹出提示
		// 	const custom = response.config?.custom;
		// 	if (custom?.toast !== false) {
		// 		uni.showToast({ title: data.msg || '请求失败', icon: 'none', duration: 3000 });
		// 	}

		// 	// 始终返回reject，确保前端try-catch能正确捕获错误，finally能执行
		// 	return Promise.reject(data)
		// }
		return data === undefined ? {} : data
	}, (response) => { 
		const data = response.data
		// 对响应错误做点什么 （statusCode !== 200）
		// 注意：后端已调整，只有真正的服务器错误才会返回非200状态码
		
		// 401 未授权 - token过期或无效，直接跳转登录页
		if (response.statusCode === 401) {
			// 清除本地存储的用户信息
			uni.removeStorageSync('token');
			uni.removeStorageSync('user');
			uni.removeStorageSync('merch');
			
			// 直接跳转登录页，不显示toast避免重复
			uni.reLaunch({ url: '/pages/login/login' });
			return Promise.reject({ msg: '登录已过期' });
		}
		
		// 尝试解析响应数据
		// let responseData = response.data;
		// if (responseData && typeof responseData === 'string') {
		// 	try {
		// 		responseData = JSON.parse(responseData);
		// 	} catch (e) {
		// 		responseData = null;
		// 	}
		// }
		
		// 如果响应数据是对象格式且包含msg字段，显示错误信息
		// if (responseData && typeof responseData === 'object' && responseData.msg) {
		// 	uni.showToast({ title: responseData.msg, icon: 'none', duration: 3000 });
		// 	return Promise.reject(responseData);
		// }
		
		// 根据HTTP状态码设置错误消息
		// let message = '网络请求失败';
		// switch (response.statusCode) {
		// 	case 429:
		// 		message = responseData?.msg || '请求过于频繁，请稍后再试';
		// 		break;
		// 	case 404:
		// 		message = '请求接口不存在';
		// 		break;
		// 	case 500:
		// 		message = responseData?.msg || '服务器内部错误';
		// 		break;
		// 	default:
		// 		message = `请求失败 [${response.statusCode}]`;
		// }
		
		// uni.showToast({ title: message, icon: 'none', duration: 3000 });
		// return Promise.reject({ msg: message })
		return data === undefined ? {} : data
	})
	
	// 将 http 实例挂载到 app 的全局属性上
	app.config.globalProperties.$http = uni.$uv.http;
}
