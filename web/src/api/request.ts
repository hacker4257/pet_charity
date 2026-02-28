import axios from 'axios'
import { message } from 'antd'



const request = axios.create({
    baseURL: '/api/v1',
    timeout: 10000,
});


//请求拦截器
request.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token');
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    },
    (error) => Promise.reject(error)
)

//response
request.interceptors.response.use(
    (response) => {
        const res = response.data
        if (res.code !== 0) {
            message.error(res.message || 'Request failed');
            return Promise.reject(new Error(res.message))
        }
        return res;
    },
    (error) => {
        if (error.response) {
            switch (error.response.status) {
                case 401:
                    //token过期
                    localStorage.removeItem('access_token')
                    window.location.href = '/login'
                    break
                case 403:
                    message.error('No permission')
                    break
                case 404:
                    message.error('Resource not found')
                    break
                default:
                    message.error(error.response.data?.message || 'Server error')
                
            }
        }else {
            message.error('Network error')
        }
        return Promise.reject(error)
    }
)

export default request;