'use client'

import { createContext, useCallback, useContext, useEffect, useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import { api, clearTokens, getAccessToken, setTokens } from './api'
import type { AuthTokens, User } from './types'

interface AuthState {
  user: User | null
  loading: boolean
  login: (identifier: string, password: string) => Promise<void>
  logout: () => Promise<void>
  refresh: () => Promise<void>
  setUser: (u: User | null) => void
}

const AuthContext = createContext<AuthState | null>(null)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const router = useRouter()

  const fetchMe = useCallback(async (): Promise<User | null> => {
    if (!getAccessToken()) {
      return null
    }
    try {
      const res = await api.get<{ user: User }>('/api/v1/auth/me')
      return res.user
    } catch {
      return null
    }
  }, [])

  useEffect(() => {
    let mounted = true
    ;(async () => {
      const me = await fetchMe()
      if (mounted) {
        setUser(me)
        setLoading(false)
      }
    })()
    return () => {
      mounted = false
    }
  }, [fetchMe])

  const login = useCallback(
    async (identifier: string, password: string) => {
      const res = await api.post<{ user: User; tokens: AuthTokens }>(
        '/api/v1/auth/login',
        { identifier, password },
        { auth: false }
      )
      setTokens(res.tokens.access_token, res.tokens.refresh_token)
      setUser(res.user)
    },
    []
  )

  const logout = useCallback(async () => {
    try {
      await api.post('/api/v1/auth/logout', undefined, { auth: false }).catch(() => {})
    } finally {
      clearTokens()
      setUser(null)
      router.replace('/login')
    }
  }, [router])

  const refresh = useCallback(async () => {
    const me = await fetchMe()
    setUser(me)
  }, [fetchMe])

  const value = useMemo<AuthState>(
    () => ({ user, loading, login, logout, refresh, setUser }),
    [user, loading, login, logout, refresh]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

export function useRequireAuth() {
  const auth = useAuth()
  const router = useRouter()
  useEffect(() => {
    if (!auth.loading && !auth.user) {
      router.replace('/login')
    }
  }, [auth.loading, auth.user, router])
  return auth
}
