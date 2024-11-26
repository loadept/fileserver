import axios from 'axios'

const API_URL = import.meta.env.API_URL

const apiClient = axios.create({
  baseURL: API_URL
})

export default apiClient
