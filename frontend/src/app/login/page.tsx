'use client'

import { useState } from 'react'

export default function LoginPage() {
  const [identifier, setIdentifier] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ identifier, password }),
      })

      if (!res.ok) {
        const data = await res.json()
        setError(data.error?.message || 'Login failed')
        return
      }

      const data = await res.json()
      localStorage.setItem('access_token', data.data.tokens.access_token)
      localStorage.setItem('refresh_token', data.data.tokens.refresh_token)
      window.location.href = '/dashboard'
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-obsidian relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(210,255,42,0.03),transparent_70%)]" />
      <form onSubmit={handleSubmit} className="max-w-md w-full card p-8 relative animate-slide-up">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-pearl mb-2">Welcome back</h1>
          <p className="text-slate text-sm">Sign in to your account</p>
        </div>

        {error && (
          <div className="mb-6 p-3 rounded-lg bg-red-900/20 border border-red-500/30 text-red-400 text-sm animate-fade-in">
            {error}
          </div>
        )}

        <div className="space-y-5">
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Username or Email</label>
            <input
              type="text"
              value={identifier}
              onChange={(e) => setIdentifier(e.target.value)}
              className="input-field"
              placeholder="Enter your username or email"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="input-field"
              placeholder="Enter your password"
              required
            />
          </div>
          <button type="submit" className="btn-primary w-full" disabled={loading}>
            {loading ? 'Signing in...' : 'Sign in'}
          </button>
        </div>
      </form>
    </div>
  )
}
