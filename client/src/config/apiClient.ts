import axios from 'axios'

const SECRET_API_URL = import.meta.env.SECRET_API_URL

const apiClient = axios.create({
  baseURL: SECRET_API_URL
})

export default apiClient
