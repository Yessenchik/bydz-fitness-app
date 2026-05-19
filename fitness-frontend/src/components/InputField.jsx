const InputField = ({
                        label,
                        id,
                        type = 'text',
                        placeholder,
                        value,
                        onChange,
                        error,
                        autoComplete,
                        icon: Icon,
                    }) => {
    return (
        <div>
            <label htmlFor={id} className="label-base">
                {label}
            </label>
            <div className="relative">
                {Icon && (
                    <div className="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                        <Icon className="w-4 h-4 text-zinc-500" />
                    </div>
                )}
                <input
                    id={id}
                    type={type}
                    placeholder={placeholder}
                    value={value}
                    onChange={onChange}
                    autoComplete={autoComplete}
                    className={`input-base ${Icon ? 'pl-10' : ''} ${
                        error ? 'border-red-500 focus:ring-red-500' : ''
                    }`}
                />
            </div>
            {error && <p className="error-text">{error}</p>}
        </div>
    )
}

export default InputField