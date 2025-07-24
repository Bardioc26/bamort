import axios from 'axios'

const API = axios.create({
  //baseURL: 'http://192.168.0.48:8180', // Replace with your backend URL
  baseURL: 'http://localhost:8180', // Replace with your backend URL
})

// Request interceptor to add auth token
API.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

export default API
