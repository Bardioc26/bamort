import axios from 'axios'

const API = axios.create({
  baseURL: 'http://localhost:8080', // Replace with your backend URL
})

export default API
