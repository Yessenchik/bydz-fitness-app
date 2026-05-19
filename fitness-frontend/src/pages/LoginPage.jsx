import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { authApi } from '../api/client'
import InputField from '../components/InputField'
import Spinner from '../components/Spinner'

// Иконки (inline SVG как компоненты — без лишних зависимостей)
const MailIcon = ({ className }) => (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
    </svg>
)

const LockIcon = ({ className }) => (
    <svg className={className} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
    </svg>
)

const LoginPage = () => {
    const navigate = useNavigate()

    const [form, setForm] = useState({ email: '', password: '' })
    const [errors, setErrors] = useState({})
    const [apiError, setApiError] = useState('')
    const [loading, setLoading] = useState(false)

    const handleChange = (e) => {
        const { id, value } = e.target
        setForm((prev) => ({ ...prev, [id]: value }))
        // Сбрасываем ошибку поля при вводе
        if (errors[id]) setErrors((prev) => ({ ...prev, [id]: '' }))
    }

    const validate = () => {
        const errs = {}
        if (!form.email) errs.email = 'Введите email'
        else if (!/\S+@\S+\.\S+/.test(form.email)) errs.email = 'Некорректный email'
        if (!form.password) errs.password = 'Введите пароль'
        return errs
    }

    const handleSubmit = async (e) => {
        e.preventDefault()
        setApiError('')

        const validationErrors = validate()
        if (Object.keys(validationErrors).length > 0) {
            setErrors(validationErrors)
            return
        }

        setLoading(true)
        try {
            const { data } = await authApi.login({
                email: form.email,
                password: form.password,
            })

            // Сохраняем токены
            localStorage.setItem('access_token', data.access_token)
            localStorage.setItem('refresh_token', data.refresh_token)
            localStorage.setItem('user_id', data.user.user_id)

            navigate('/dashboard')
        } catch (err) {
            const message = err.response?.data?.message || 'Неверный email или пароль'
            setApiError(message)
        } finally {
            setLoading(false)
        }
    }

    return (
        <div className="min-h-screen flex items-center justify-center px-4">
            {/* Фоновые декоративные блобы */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute -top-40 -right-40 w-96 h-96 bg-brand-500/10
                        rounded-full blur-3xl" />
                <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-brand-500/5
                        rounded-full blur-3xl" />
            </div>

            <div className="w-full max-w-md relative">
                {/* Логотип / заголовок */}
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-16 h-16
                          bg-brand-500/20 rounded-2xl mb-4 border border-brand-500/30">
                        <svg className="w-8 h-8 text-brand-400" fill="none" viewBox="0 0 24 24"
                             stroke="currentColor" strokeWidth={1.5}>
                            <path strokeLinecap="round" strokeLinejoin="round"
                                  d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
                        </svg>
                    </div>
                    <h1 className="text-3xl font-bold text-zinc-50">FitApp</h1>
                    <p className="text-zinc-500 mt-1 text-sm">Войдите в свой аккаунт</p>
                </div>

                <div className="card">
                    {/* Глобальная ошибка API */}
                    {apiError && (
                        <div className="mb-5 px-4 py-3 bg-red-500/10 border border-red-500/30
                            rounded-xl text-red-400 text-sm flex items-center gap-2">
                            <svg className="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                                <path fillRule="evenodd"
                                      d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" />
                            </svg>
                            {apiError}
                        </div>
                    )}

                    <form onSubmit={handleSubmit} className="space-y-5" noValidate>
                        <InputField
                            label="Email"
                            id="email"
                            type="email"
                            placeholder="you@example.com"
                            value={form.email}
                            onChange={handleChange}
                            error={errors.email}
                            autoComplete="email"
                            icon={MailIcon}
                        />
                        <InputField
                            label="Пароль"
                            id="password"
                            type="password"
                            placeholder="••••••••"
                            value={form.password}
                            onChange={handleChange}
                            error={errors.password}
                            autoComplete="current-password"
                            icon={LockIcon}
                        />

                        <button type="submit" className="btn-primary mt-2" disabled={loading}>
                            {loading ? <Spinner /> : null}
                            {loading ? 'Вход...' : 'Войти'}
                        </button>
                    </form>

                    <p className="text-center text-sm text-zinc-500 mt-6">
                        Нет аккаунта?{' '}
                        <Link to="/register" className="btn-ghost">
                            Зарегистрироваться
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    )
}

export default LoginPage