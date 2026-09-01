'use client'

import { useState } from 'react'
import { useAuth } from '@/lib/auth-context'
import { initials } from '@/lib/utils'

interface TopbarProps {
  onMenu: () => void
  title?: string
}

export default function Topbar({ onMenu, title }: TopbarProps) {
  const { user, logout } = useAuth()
  const [menuOpen, setMenuOpen] = useState(false)

  return (
    <header className="topbar">
      <div className="flex items-center gap-3 px-4 lg:px-8 py-3">
        <button
          onClick={onMenu}
          className="lg:hidden btn btn-ghost btn-sm"
          aria-label="Open menu"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        {title && (
          <div className="hidden md:block text-sm text-[var(--text-dim)]">
            {title}
          </div>
        )}
        <div className="ml-auto flex items-center gap-2">
          <div className="relative">
            <button
              onClick={() => setMenuOpen((v) => !v)}
              className="flex items-center gap-3 rounded-lg border border-[var(--border)] bg-[var(--bg-elev)] px-2.5 py-1.5 hover:border-[var(--border-strong)] transition"
            >
              <div className="grid h-7 w-7 place-items-center rounded-md bg-[var(--accent-soft)] border border-[var(--accent)]/30 text-[var(--accent)] text-[11px] font-bold">
                {initials(user?.username || user?.email)}
              </div>
              <div className="hidden sm:block text-left leading-tight">
                <div className="text-[13px] font-medium text-[var(--text)]">
                  {user?.username || 'User'}
                </div>
                <div className="text-[11px] text-[var(--text-mute)] capitalize">
                  {(user?.role || '').toLowerCase().replace('_', ' ')}
                </div>
              </div>
              <svg className="w-4 h-4 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 9l6 6 6-6" />
              </svg>
            </button>
            {menuOpen && (
              <>
                <div
                  className="fixed inset-0 z-30"
                  onClick={() => setMenuOpen(false)}
                />
                <div className="absolute right-0 z-40 mt-2 w-56 origin-top-right rounded-xl border border-[var(--border)] bg-[var(--bg-elev)] shadow-2xl py-1.5 animate-fade-in">
                  <div className="px-3 py-2 border-b border-[var(--border)]">
                    <div className="text-sm font-medium truncate">{user?.username}</div>
                    <div className="text-xs text-[var(--text-mute)] truncate">{user?.email}</div>
                  </div>
                  <button
                    onClick={() => {
                      setMenuOpen(false)
                      logout()
                    }}
                    className="w-full text-left px-3 py-2 text-sm hover:bg-white/5 flex items-center gap-2 text-[var(--danger)]"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
                      <path strokeLinecap="round" strokeLinejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    Sign out
                  </button>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
