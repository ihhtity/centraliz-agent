export const Request = (app) => {
	// 初始化请求配置
	uni.$uv.http.setConfig((config) => {
		/* config 为默认全局配置*/
		// #ifdef H5
		config.baseURL = '/api/v1'; // 本地开发时代理到/v1前缀
		// #endif

		// #ifdef MP-WEIXIN || APP-PLUS
		config.baseURL = 'https://centraliz.bsldtech.com/api/v1';  // 线上
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
		// 根据custom参数中配置的是否需要token，添加对应的请求头
		if (config?.custom?.auth) {
			// 从localStorage获取token
			const token = uni.getStorageSync('token');
			
			if (token) {
				config.header.Authorization = `Bearer ${token}`;
			} else {
				// token不存在时跳转到登录页
				uni.removeStorageSync('token');
				uni.removeStorageSync('userInfo');
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

		// 处理业务错误
		if (data.code !== 200 && data.code !== 0) { 
			// 如果没有显式定义custom的toast参数为false的话，默认对报错进行toast弹出提示
			const custom = response.config?.custom;
			if (custom?.toast !== false) {
				uni.$uv.toast(data.msg || '请求失败')
			}

			// 如果需要catch返回，则进行reject
			if (custom?.catch) {
				return Promise.reject(data)
			} else {
				// 否则返回一个pending中的promise，请求不会进入catch中
				return new Promise(() => { })
			}
		}
		return data === undefined ? {} : data
	}, (response) => { 
		// 对响应错误做点什么 （statusCode !== 200）
		
		// 401 未授权 - token过期或无效，直接跳转登录页
		if (response.statusCode === 401) {
			// 清除本地存储的用户信息
			uni.removeStorageSync('token');
			uni.removeStorageSync('userInfo');
			uni.removeStorageSync('merch');
			
			// 直接跳转登录页，不显示toast避免重复
			uni.reLaunch({ url: '/pages/login/login' });
			return Promise.reject(response);
		}
		
		// 429 限流 - 请求过多
		if (response.statusCode === 429) {
			const message = response.data?.msg || '请求过于频繁，请稍后再试';
			uni.showToast({ 
				title: message, 
				icon: 'none', 
				duration: 2000 
			});
			return Promise.reject(response);
		}
		
		// 其他错误处理
		let message = '网络请求失败';
		
		// 尝试从响应中获取后端返回的错误信息
		if (response.data && response.data.msg) {
			message = response.data.msg;
		} else if (response.statusCode === 404) {
			message = '请求接口不存在';
		} else if (response.statusCode === 400) {
			message = response.data?.msg || '请求参数错误';
		} else if (response.statusCode >= 500) {
			message = response.data?.msg || '服务器内部错误';
		}
		
		uni.showToast({ title: message, icon: 'none', duration: 3000 });
		return Promise.reject(response)
	})
	
	// 将 http 实例挂载到 app 的全局属性上
	app.config.globalProperties.$http = uni.$uv.http;
}
