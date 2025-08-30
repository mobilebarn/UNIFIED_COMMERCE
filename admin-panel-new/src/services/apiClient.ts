import axios from 'axios'

// Basic token storage; later replace with secure storage / refresh flow
let accessToken: string | null = null
export function setAccessToken(token: string){ accessToken = token }
export function getAccessToken(){ return accessToken }

export const api = axios.create({
  baseURL: 'http://localhost:4000', // GraphQL gateway or REST aggregation (adjust as needed)
  timeout: 8000,
})

api.interceptors.request.use(cfg => {
  if(accessToken) cfg.headers.Authorization = `Bearer ${accessToken}`
  return cfg
})

api.interceptors.response.use(r=>r, err => {
  // Simple token expiration handling placeholder
  if(err.response?.status === 401){
    // TODO: trigger refresh flow
  }
  return Promise.reject(err)
})
