import React from 'react'

interface AlertProps {
  type?: 'success' | 'error' | 'warning' | 'info'
  title?: string
  children: React.ReactNode
  onDismiss?: () => void
  icon?: React.ReactNode
}

export function Alert({ type = 'info', title, children, onDismiss, icon }: AlertProps) {
  const styles = {
    success: {
      bg: 'bg-[#0d3b16]/50',
      border: 'border-[#22c55e]/30',
      text: 'text-[#22c55e]',
      icon: (
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
        </svg>
      ),
    },
    error: {
      bg: 'bg-[#3b0d0d]/50',
      border: 'border-[#ef4444]/30',
      text: 'text-[#ef4444]',
      icon: (
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
        </svg>
      ),
    },
    warning: {
      bg: 'bg-[#3b2d0d]/50',
      border: 'border-[#f59e0b]/30',
      text: 'text-[#f59e0b]',
      icon: (
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
        </svg>
      ),
    },
    info: {
      bg: 'bg-[var(--bg-elev)]',
      border: 'border-[var(--border)]',
      text: 'text-[var(--text)]',
      icon: (
        <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
        </svg>
      ),
    },
  }

  const style = styles[type]

  return (
    <div className={`rounded-lg border p-4 ${style.bg} ${style.border} ${style.text}`}>
      <div className="flex gap-3">
        <div className="flex-shrink-0 mt-0.5">
          {icon || style.icon}
        </div>
        <div className="flex-1">
          {title && <h3 className="font-medium mb-1">{title}</h3>}
          <div className="text-sm opacity-90">{children}</div>
        </div>
        {onDismiss && (
          <button
            onClick={onDismiss}
            className="flex-shrink-0 p-1 hover:bg-white/10 rounded transition"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        )}
      </div>
    </div>
  )
}
