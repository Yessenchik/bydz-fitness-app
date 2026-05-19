import axios from 'axios'

const BASE_URL = 'http://localhost:8080'

// Создаём инстанс Axios с базовым URL
const api = axios.create({
    baseURL: BASE_URL,
    headers: { 'Content-Type': 'application/json' },
})

// Request интерцептор — автоматически добавляет Bearer токен
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('access_token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => Promise.reject(error)
)

// Response интерцептор — обрабатывает 401 (токен протух)
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            window.location.href = '/login'
        }
        return Promise.reject(error)
    }
)

// ── API методы ──────────────────────────────

export const authApi = {
    register: (data) =>
        api.post('/api/auth/register', data),

    login: (data) =>
        api.post('/api/auth/login', data),
}

export const profileApi = {
        getProfile: (userId) => api.get(`/api/profile?user_id=${userId}`),
}

export default api