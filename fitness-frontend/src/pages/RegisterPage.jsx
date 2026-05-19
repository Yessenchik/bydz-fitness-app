import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { authApi } from '../api/client'
import InputField from '../components/InputField'
import Spinner from '../components/Spinner'

const RegisterPage = () => {
    const navigate = useNavigate()

    const [form, setForm] = useState({
        email: '',
        password: '',
        confirmPassword: '',
        first_name: '',
        last_name: '',
        phone: '',
    })
    const [errors, setErrors] = useState({})
    const [apiError, setApiError] = useState('')
    const [success, setSuccess] = useState(false)
    const [loading, setLoading] = useState(false)

    const handleChange = (e) => {
        const { id, value } = e.target
        setForm((prev) => ({ ...prev, [id]: value }))
        if (errors[id]) setErrors((prev) => ({ ...prev, [id]: '' }))
    }

    const validate = () => {
        const errs = {}
        if (!form.first_name.trim()) errs.first_name = 'Введите имя'
        if (!form.last_name.trim())  errs.last_name  = 'Введите фамилию'
        if (!form.email) errs.email = 'Введите email'
        else if (!/\S+@\S+\.\S+/.test(form.email)) errs.email = 'Некорректный email'
        if (!form.phone.trim()) errs.phone = 'Введите телефон'
        if (!form.password) errs.password = 'Введите пароль'
        else if (form.password.length < 8) errs.password = 'Минимум 8 символов'
        else if (!/[A-Z]/.test(form.password)) errs.password = 'Нужна хотя бы одна заглавная буква'
        else if (!/[0-9]/.test(form.password)) errs.password = 'Нужна хотя бы одна цифра'
        if (form.password !== form.confirmPassword)
            errs.confirmPassword = 'Пароли не совпадают'
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
            await authApi.register({
                email:      form.email,
                password:   form.password,
                first_name: form.first_name,
                last_name:  form.last_name,
                phone:      form.phone,
            })
            setSuccess(true)
        } catch (err) {
            const message = err.response?.data?.message || 'Ошибка регистрации. Попробуйте снова.'
            setApiError(message)
        } finally {
            setLoading(false)
        }
    }

    // Экран успешной регистрации
    if (success) {
        return (
            <div className="min-h-screen flex items-center justify-center px-4">
                <div className="w-full max-w-md">
                    <div className="card text-center">
                        <div className="inline-flex items-center justify-center w-16 h-16
                            bg-brand-500/20 rounded-full mb-5 mx-auto
                            border border-brand-500/30">
                            <svg className="w-8 h-8 text-brand-400" fill="none" viewBox="0 0 24 24"
                                 stroke="currentColor" strokeWidth={2}>
                                <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
                            </svg>
                        </div>
                        <h2 className="text-2xl font-bold text-zinc-50 mb-2">
                            Регистрация успешна!
                        </h2>
                        <p className="text-zinc-400 mb-6 text-sm leading-relaxed">
                            На ваш email отправлено письмо для подтверждения.
                            После подтверждения вы сможете войти в аккаунт.
                        </p>
                        <Link to="/login" className="btn-primary inline-flex w-auto px-8">
                            Перейти ко входу
                        </Link>
                    </div>
                </div>
            </div>
        )
    }

    return (
        <div className="min-h-screen flex items-center justify-center px-4 py-12">
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute -top-40 -right-40 w-96 h-96 bg-brand-500/10
                        rounded-full blur-3xl" />
                <div className="absolute -bottom-40 -left-40 w-96 h-96 bg-brand-500/5
                        rounded-full blur-3xl" />
            </div>

            <div className="w-full max-w-md relative">
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
                    <p className="text-zinc-500 mt-1 text-sm">Создайте аккаунт и начните тренироваться</p>
                </div>

                <div className="card">
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

                    <form onSubmit={handleSubmit} className="space-y-4" noValidate>
                        {/* Имя и фамилия в строку */}
                        <div className="grid grid-cols-2 gap-3">
                            <InputField
                                label="Имя"
                                id="first_name"
                                placeholder="Иван"
                                value={form.first_name}
                                onChange={handleChange}
                                error={errors.first_name}
                                autoComplete="given-name"
                            />
                            <InputField
                                label="Фамилия"
                                id="last_name"
                                placeholder="Иванов"
                                value={form.last_name}
                                onChange={handleChange}
                                error={errors.last_name}
                                autoComplete="family-name"
                            />
                        </div>

                        <InputField
                            label="Email"
                            id="email"
                            type="email"
                            placeholder="you@example.com"
                            value={form.email}
                            onChange={handleChange}
                            error={errors.email}
                            autoComplete="email"
                        />

                        <InputField
                            label="Телефон"
                            id="phone"
                            type="tel"
                            placeholder="+7 999 000 00 00"
                            value={form.phone}
                            onChange={handleChange}
                            error={errors.phone}
                            autoComplete="tel"
                        />

                        <InputField
                            label="Пароль"
                            id="password"
                            type="password"
                            placeholder="Минимум 8 символов"
                            value={form.password}
                            onChange={handleChange}
                            error={errors.password}
                            autoComplete="new-password"
                        />

                        <InputField
                            label="Подтвердите пароль"
                            id="confirmPassword"
                            type="password"
                            placeholder="••••••••"
                            value={form.confirmPassword}
                            onChange={handleChange}
                            error={errors.confirmPassword}
                            autoComplete="new-password"
                        />

                        {/* Требования к паролю */}
                        <div className="text-xs text-zinc-500 space-y-1 px-1">
                            {[
                                { ok: form.password.length >= 8,        text: 'Минимум 8 символов' },
                                { ok: /[A-Z]/.test(form.password),      text: 'Заглавная буква' },
                                { ok: /[0-9]/.test(form.password),      text: 'Цифра' },
                            ].map(({ ok, text }) => (
                                <div key={text} className="flex items-center gap-1.5">
                                    <div className={`w-1.5 h-1.5 rounded-full transition-colors ${
                                        ok ? 'bg-brand-400' : 'bg-zinc-600'
                                    }`} />
                                    <span className={ok ? 'text-brand-400' : ''}>{text}</span>
                                </div>
                            ))}
                        </div>

                        <button type="submit" className="btn-primary mt-2" disabled={loading}>
                            {loading ? <Spinner /> : null}
                            {loading ? 'Регистрация...' : 'Создать аккаунт'}
                        </button>
                    </form>

                    <p className="text-center text-sm text-zinc-500 mt-6">
                        Уже есть аккаунт?{' '}
                        <Link to="/login" className="btn-ghost">
                            Войти
                        </Link>
                    </p>
                </div>
            </div>
        </div>
    )
}

export default RegisterPage