import axios from "axios"
import { MessagePlugin } from 'tdesign-vue-next'
import { getToken, clearAuthData } from './auth'

interface HTTP_RESPONSE {
	state: string | number
	data: any
	msg: string
}

const instance = axios.create({
	baseURL: process.env.BASE_API_URL,
	timeout: 5 * 60 * 1000,
})

instance.interceptors.request.use(
	async function (config: any) {
		config.baseURL = process.env.BASE_API_URL
		// 使用新的认证工具获取 token
		const token = getToken()

		if (token) {
			config.headers.Authorization = `Bearer ${token}`
		}
		return config
	},
	function (error) {
		return Promise.reject(error)
	}
)

instance.interceptors.response.use(
	function (response) {

		switch (response.data.state) {
			// 成功状态 返回data数据源
			case 200:
				return response.data.data as HTTP_RESPONSE["data"]
			// 数据异常
			case 100102:
			case 100101:
			case 100009:
			case 100005:
				MessagePlugin.error(response.data.msg)
				return
			// 其他错误 返回整个服务端下发数据（包含状态码和消息）
			default:
				return response.data
		}
	},
	function (error) {
		if (error && !error.response) {
			clearAuthData()
			localStorage.clear()
		}
		if (error && error.response) {
			switch (error.response.status) {
				case 401:
					// 客户端环境
					if (window) {
						location.href = `${process.env.BASE}user/login`
					}
					break
				case 500:
				case 503:
          MessagePlugin.error(error.response.data.msg)
					break
				case 501:
					// 没有权限
					// window && (location.href = "/user/forbin")
					break
			}
		}
		return Promise.reject(error)
	}
)

export default instance

