export const Request = (vm) => {
	// 初始化请求配置
	uni.$uv.http.setConfig((config) => {
		/* config 为默认全局配置*/
		// #ifdef H5
		config.baseURL = '/api'; // 本地开发时代理到/api前缀
		// #endif

		// #ifdef MP-WEIXIN || APP-PLUS
		config.baseURL = 'https://bsldyyc.bsldtech.com/electronic';  // 线上
		// #endif

		//请求头
		config.header = {
			'X-Requested-With': 'XMLHttpRequest',
			'Content-Type': 'application/json'
		}

		return config
	})

	// 请求拦截
	uni.$uv.http.interceptors.request.use((config) => { // 可使用async await 做异步操作
		// 初始化请求拦截器时，会执行此方法，此时data为undefined，赋予默认{}
		config.data = config.data || {}
		// 根据custom参数中配置的是否需要token，添加对应的请求头
		if (config?.custom?.auth) {
			// 从localStorage获取token
			const token = uni.getStorageSync('token');
			if (token) {
				config.header.Authorization = `Bearer ${token}`;
			}
		}
		return config
	}, config => { // 可使用async await 做异步操作
		return Promise.reject(config)
	})

	// 响应拦截
	uni.$uv.http.interceptors.response.use((response) => { /* 对响应成功做点什么 可使用async await 做异步操作*/
		const data = response.data

		// 处理业务错误
		if (data.code !== 200 && data.code !== 0) { 
			// 如果没有显式定义custom的toast参数为false的话，默认对报错进行toast弹出提示
			const custom = response.config?.custom;
			if (custom?.toast !== false) {
				uni.$uv.toast(data.message || '请求失败')
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
		let message = '网络请求失败';
		if (response.statusCode === 401) {
			message = '登录已过期，请重新登录';
			// 清除本地存储的用户信息
			uni.removeStorageSync('token');
			uni.removeStorageSync('userInfo');
			uni.removeStorageSync('userRole');
			// 跳转到登录页
			setTimeout(() => {
				uni.reLaunch({ url: '/pages/login/login' });
			}, 1000);
		} else if (response.statusCode === 404) {
			message = '请求接口不存在';
		} else if (response.statusCode >= 500) {
			message = '服务器内部错误';
		}
		
		uni.$uv.toast(message);
		return Promise.reject(response)
	})
}