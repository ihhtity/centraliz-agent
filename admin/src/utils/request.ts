import axios, { InternalAxiosRequestConfig, AxiosResponse, AxiosRequestConfig } from 'axios';

interface ApiResponse<T = any> {
  code: number;
  msg: string;
  data: T;
}

interface CustomConfig {
  auth?: boolean;
  toast?: boolean;
}

type RequestConfig = AxiosRequestConfig & {
  custom?: CustomConfig;
};

const axiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'X-Requested-With': 'XMLHttpRequest',
    'Content-Type': 'application/json',
  },
});

axiosInstance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const reqConfig = config as RequestConfig;
    if (reqConfig?.custom?.auth) {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      } else {
        localStorage.removeItem('token');
        window.location.href = '/login';
        return Promise.reject(new Error('未登录'));
      }
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

axiosInstance.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const data = response.data;

    if (!data || typeof data !== 'object') {
      return Promise.reject({ msg: '响应格式错误' });
    }

    if (data.code !== 200) {
      if (data.code === 401) {
        localStorage.removeItem('token');
        window.location.href = '/login';
      }

      return Promise.reject(data);
    }

    return response;
  },
  (error) => {
    const response = error.response;

    if (response?.status === 401) {
      localStorage.removeItem('token');
      window.location.href = '/login';
      return Promise.reject({ msg: '登录已过期' });
    }

    let errorMsg = '网络请求失败';
    if (response) {
      switch (response.status) {
        case 429:
          errorMsg = (response.data as ApiResponse)?.msg || '请求过于频繁，请稍后再试';
          break;
        case 404:
          errorMsg = '请求接口不存在';
          break;
        case 500:
          errorMsg = (response.data as ApiResponse)?.msg || '服务器内部错误';
          break;
        default:
          errorMsg = `请求失败 [${response.status}]`;
      }
    } else if (error.message) {
      errorMsg = error.message;
    }

    return Promise.reject({ msg: errorMsg });
  }
);

const request = {
  get: async function <T = any>(url: string, config?: RequestConfig): Promise<ApiResponse<T>> {
    const res = await axiosInstance.get<ApiResponse<T>>(url, config);
    return res.data;
  },
  post: async function <T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    const res = await axiosInstance.post<ApiResponse<T>>(url, data, config);
    return res.data;
  },
  put: async function <T = any>(url: string, data?: any, config?: RequestConfig): Promise<ApiResponse<T>> {
    const res = await axiosInstance.put<ApiResponse<T>>(url, data, config);
    return res.data;
  },
  delete: async function <T = any>(url: string, config?: RequestConfig): Promise<ApiResponse<T>> {
    const res = await axiosInstance.delete<ApiResponse<T>>(url, config);
    return res.data;
  },
};

export default request;
