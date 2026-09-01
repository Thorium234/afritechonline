'use client'

import { useState, useCallback } from 'react'
import { cn } from '@/lib/utils'
import { validatePassword, getPasswordStrengthColor, getPasswordStrengthLabel } from '@/lib/passwordValidator'

interface PasswordInputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string
  showStrength?: boolean
  label?: string
}

export function PasswordInput({
  error,
  showStrength = false,
  label,
  value = '',
  onChange,
  ...props
}: PasswordInputProps) {
  const [showPassword, setShowPassword] = useState(false)
  const passwordStr = typeof value === 'string' ? value : ''
  const strength = showStrength ? validatePassword(passwordStr) : null

  return (
    <div className="space-y-2">
      {label && (
        <label className="block text-sm font-medium text-[var(--text)]">
          {label}
        </label>
      )}

      <div className="relative">
        <input
          type={showPassword ? 'text' : 'password'}
          value={value}
          onChange={onChange}
          className={cn(
            'w-full px-3 py-2 rounded-lg border bg-[var(--bg)] text-[var(--text)] placeholder:text-[var(--text-mute)]',
            error ? 'border-[var(--danger)]' : 'border-[var(--border)] focus:border-[var(--accent)]',
            'focus:outline-none focus:ring-2 focus:ring-[var(--accent)]/20 transition pr-10',
          )}
          {...props}
        />

        <button
          type="button"
          onClick={() => setShowPassword(!showPassword)}
          className="absolute right-3 top-1/2 -translate-y-1/2 p-1 hover:bg-white/5 rounded transition"
          tabIndex={-1}
        >
          {showPassword ? (
            <svg className="w-5 h-5 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
          ) : (
            <svg className="w-5 h-5 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-4.803m5.596-3.856a3.375 3.375 0 11-6.75 0 3.375 3.375 0 016.75 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 3l18 18" />
            </svg>
          )}
        </button>
      </div>

      {error && <p className="text-xs text-[var(--danger)]">{error}</p>}

      {showStrength && strength && (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <p className="text-xs text-[var(--text-dim)]">Strength</p>
            <p className={cn(
              'text-xs font-medium',
              strength.score === 0 && 'text-[var(--danger)]',
              strength.score === 1 && 'text-[#f59e0b]',
              strength.score === 2 && 'text-[#eab308]',
              strength.score === 3 && 'text-[#84cc16]',
              strength.score === 4 && 'text-[#22c55e]',
            )}>
              {getPasswordStrengthLabel(strength.score)}
            </p>
          </div>

          <div className="w-full h-2 bg-[var(--bg-elev)] rounded-full overflow-hidden">
            <div
              className={cn('h-full transition-all duration-300', getPasswordStrengthColor(strength.score))}
              style={{ width: `${strength.percentage}%` }}
            />
          </div>

          {strength.feedback.length > 0 && (
            <ul className="text-xs text-[var(--text-dim)] space-y-1">
              {strength.feedback.map((item, idx) => (
                <li key={idx} className="flex items-start gap-2">
                  <span className="text-[var(--danger)] mt-0.5">•</span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          )}

          {strength.isValid && (
            <p className="text-xs text-[#22c55e] flex items-center gap-1">
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              Password meets security requirements
            </p>
          )}
        </div>
      )}
    </div>
  )
}
