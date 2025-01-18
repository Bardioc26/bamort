import axios from 'axios'

const API = axios.create({
  baseURL: 'http://192.168.0.48:8180', // Replace with your backend URL
})

export default API
