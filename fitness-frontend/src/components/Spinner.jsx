const Spinner = ({ size = 'sm' }) => {
    const sizes = {
        sm: 'w-4 h-4',
        md: 'w-6 h-6',
        lg: 'w-10 h-10',
    }
    return (
        <div
            className={`${sizes[size]} border-2 border-white/30 border-t-white 
                  rounded-full animate-spin`}
        />
    )
}

export default Spinner