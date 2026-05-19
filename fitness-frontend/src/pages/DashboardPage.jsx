import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { profileApi } from '../api/client'
import Spinner from '../components/Spinner'

// Маппинг ролей на русский
const roleLabels = {
    client:  { label: 'Клиент',    color: 'text-blue-400  bg-blue-400/10  border-blue-400/30'  },
    trainer: { label: 'Тренер',    color: 'text-purple-400 bg-purple-400/10 border-purple-400/30' },
    admin:   { label: 'Администратор', color: 'text-amber-400  bg-amber-400/10  border-amber-400/30'  },
}

// Иконки
const UserIcon = () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z" />
    </svg>
)

const MailIcon = () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M21.75 6.75v10.5a2.25 2.25 0 01-2.25 2.25h-15a2.25 2.25 0 01-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25m19.5 0v.243a2.25 2.25 0 01-1.07 1.916l-7.5 4.615a2.25 2.25 0 01-2.36 0L3.32 8.91a2.25 2.25 0 01-1.07-1.916V6.75" />
    </svg>
)

const PhoneIcon = () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 002.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 01-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 00-1.091-.852H4.5A2.25 2.25 0 002.25 4.5v2.25z" />
    </svg>
)

const ShieldIcon = () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
    </svg>
)

const LogoutIcon = () => (
    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
        <path strokeLinecap="round" strokeLinejoin="round"
              d="M15.75 9V5.25A2.25 2.25 0 0013.5 3h-6a2.25 2.25 0 00-2.25 2.25v13.5A2.25 2.25 0 007.5 21h6a2.25 2.25 0 002.25-2.25V15M12 9l-3 3m0 0l3 3m-3-3h12.75" />
    </svg>
)

const DashboardPage = () => {
    const navigate = useNavigate()
    const [profile, setProfile] = useState(null)
    const [loading, setLoading]  = useState(true)
    const [error, setError]      = useState('')

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const userId = localStorage.getItem('user_id')
                const { data } = await profileApi.getProfile(userId)
                setProfile(data)
            } catch {
                setError('Не удалось загрузить профиль')
            } finally {
                setLoading(false)
            }
        }
        fetchProfile()
    }, [])

    const handleLogout = () => {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('user_id')
        navigate('/login')
    }

    // ── Загрузка ─────────────────────────────
    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="text-center">
                    <Spinner size="lg" />
                    <p className="text-zinc-500 mt-4 text-sm">Загрузка профиля...</p>
                </div>
            </div>
        )
    }

    // ── Ошибка ───────────────────────────────
    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center px-4">
                <div className="card text-center max-w-sm">
                    <div className="w-12 h-12 bg-red-500/10 rounded-full flex items-center
                          justify-center mx-auto mb-4 border border-red-500/30">
                        <svg className="w-6 h-6 text-red-400" fill="none" viewBox="0 0 24 24"
                             stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round"
                                  d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
                        </svg>
                    </div>
                    <p className="text-red-400 mb-4">{error}</p>
                    <button onClick={handleLogout} className="btn-ghost text-sm">
                        Вернуться ко входу
                    </button>
                </div>
            </div>
        )
    }

    const role     = profile?.role || 'client'
    const roleInfo = roleLabels[role] || roleLabels.client
    const initials = [profile?.first_name?.[0], profile?.last_name?.[0]]
        .filter(Boolean).join('').toUpperCase() || '?'

    // Форматирование даты: ИСПРАВЛЕНО ТУТ (.seconds * 1000)
    const memberSince = profile?.created_at
        ? new Date(profile.created_at.seconds * 1000).toLocaleDateString('ru-RU', {
            day: 'numeric', month: 'long', year: 'numeric',
        })
        : '—'

    return (
        <div className="min-h-screen">
            {/* Фон */}
            <div className="fixed inset-0 overflow-hidden pointer-events-none">
                <div className="absolute -top-60 -right-60 w-[500px] h-[500px]
                        bg-brand-500/8 rounded-full blur-3xl" />
            </div>

            {/* Навбар */}
            <nav className="sticky top-0 z-10 bg-zinc-950/80 backdrop-blur-md
                      border-b border-zinc-800/50">
                <div className="max-w-5xl mx-auto px-4 h-16 flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        <div className="w-8 h-8 bg-brand-500/20 rounded-lg flex items-center
                            justify-center border border-brand-500/30">
                            <svg className="w-4 h-4 text-brand-400" fill="none" viewBox="0 0 24 24"
                                 stroke="currentColor" strokeWidth={1.5}>
                                <path strokeLinecap="round" strokeLinejoin="round"
                                      d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" />
                            </svg>
                        </div>
                        <span className="font-bold text-zinc-50">FitApp</span>
                    </div>

                    <button
                        onClick={handleLogout}
                        className="flex items-center gap-2 text-zinc-400 hover:text-zinc-100
                       transition-colors text-sm font-medium"
                    >
                        <LogoutIcon />
                        <span className="hidden sm:inline">Выйти</span>
                    </button>
                </div>
            </nav>

            {/* Контент */}
            <main className="max-w-5xl mx-auto px-4 py-8">
                {/* Hero-карточка профиля */}
                <div className="card mb-6">
                    <div className="flex flex-col sm:flex-row items-center sm:items-start gap-6">
                        {/* Аватар с инициалами */}
                        <div className="relative flex-shrink-0">
                            <div className="w-24 h-24 bg-gradient-to-br from-brand-500/30 to-brand-600/20
                              rounded-2xl flex items-center justify-center
                              border border-brand-500/30 text-3xl font-bold text-brand-400">
                                {initials}
                            </div>
                            {profile?.is_verified && (
                                <div className="absolute -bottom-1 -right-1 w-6 h-6 bg-brand-500
                                rounded-full flex items-center justify-center
                                border-2 border-zinc-900">
                                    <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd"
                                              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" />
                                    </svg>
                                </div>
                            )}
                        </div>

                        {/* Имя и мета */}
                        <div className="text-center sm:text-left flex-1">
                            <div className="flex flex-col sm:flex-row sm:items-center gap-3 mb-2">
                                <h2 className="text-2xl font-bold text-zinc-50">
                                    {profile?.first_name} {profile?.last_name}
                                </h2>
                                <span className={`inline-flex items-center px-3 py-1 rounded-full 
                                  text-xs font-medium border ${roleInfo.color}`}>
                  {roleInfo.label}
                </span>
                            </div>

                            <p className="text-zinc-500 text-sm">
                                Участник с {memberSince}
                            </p>

                            {profile?.is_verified ? (
                                <p className="text-brand-400 text-xs mt-1 flex items-center gap-1
                               justify-center sm:justify-start">
                                    <svg className="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd"
                                              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" />
                                    </svg>
                                    Email подтверждён
                                </p>
                            ) : (
                                <p className="text-amber-400 text-xs mt-1">
                                    ⚠ Email не подтверждён
                                </p>
                            )}
                        </div>
                    </div>
                </div>

                {/* Сетка информации */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">

                    {/* Контактная информация */}
                    <div className="card">
                        <h3 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-4">
                            Контактная информация
                        </h3>
                        <div className="space-y-4">
                            {[
                                { icon: <MailIcon />,  label: 'Email',   value: profile?.email || '—' },
                                { icon: <PhoneIcon />, label: 'Телефон', value: profile?.phone || '—' },
                            ].map(({ icon, label, value }) => (
                                <div key={label} className="flex items-center gap-3">
                                    <div className="w-9 h-9 bg-zinc-800 rounded-xl flex items-center
                                  justify-center text-zinc-400 flex-shrink-0">
                                        {icon}
                                    </div>
                                    <div>
                                        <p className="text-xs text-zinc-500">{label}</p>
                                        <p className="text-sm font-medium text-zinc-200">{value}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Аккаунт */}
                    <div className="card">
                        <h3 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-4">
                            Данные аккаунта
                        </h3>
                        <div className="space-y-4">
                            {[
                                { icon: <UserIcon />,   label: 'ID пользователя', value: profile?.user_id ? `#${profile.user_id.slice(0,8)}…` : '—' },
                                { icon: <ShieldIcon />, label: 'Роль',            value: roleInfo.label },
                            ].map(({ icon, label, value }) => (
                                <div key={label} className="flex items-center gap-3">
                                    <div className="w-9 h-9 bg-zinc-800 rounded-xl flex items-center
                                  justify-center text-zinc-400 flex-shrink-0">
                                        {icon}
                                    </div>
                                    <div>
                                        <p className="text-xs text-zinc-500">{label}</p>
                                        <p className="text-sm font-medium text-zinc-200">{value}</p>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                <div className="card">
                    <h3 className="text-sm font-semibold text-zinc-400 uppercase tracking-wider mb-4">
                        Активность
                    </h3>
                    <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
                        {[
                            { label: 'Тренировок',  value: '—', sub: 'всего'     },
                            { label: 'На этой неделе', value: '—', sub: 'тренировок' },
                            { label: 'Серия',       value: '—', sub: 'дней подряд' },
                            { label: 'Абонемент',   value: '—', sub: 'статус'    },
                        ].map(({ label, value, sub }) => (
                            <div key={label} className="stat-card text-center">
                                <p className="text-2xl font-bold text-zinc-100 mb-1">{value}</p>
                                <p className="text-xs font-medium text-zinc-300">{label}</p>
                                <p className="text-xs text-zinc-500 mt-0.5">{sub}</p>
                            </div>
                        ))}
                    </div>
                </div>
            </main>
        </div>
    )
}

export default DashboardPage