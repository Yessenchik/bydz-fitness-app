import { Navigate, Outlet } from 'react-router-dom'

// Если токена нет — редиректим на /login
// Outlet рендерит дочерний маршрут
const ProtectedRoute = () => {
    const token = localStorage.getItem('access_token')
    return token ? <Outlet /> : <Navigate to="/login" replace />
}

export default ProtectedRoute