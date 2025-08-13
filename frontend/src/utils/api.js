import axios from 'axios'

const API = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8180', // Use env variable with fallback
  //baseURL: 'http://192.168.0.48:8180', // Replace with your backend URL
  //baseURL: 'http://localhost:8180', // Replace with your backend URL
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

// Response interceptor to handle 401 errors
API.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response && error.response.status === 401) {
      // Token is invalid or expired
      console.warn('Authentication failed - token may be expired')
      localStorage.removeItem('token')
      // You might want to redirect to login here
      // window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default API
