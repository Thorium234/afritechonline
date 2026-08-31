'use client'

import { useState } from 'react'

export default function LoginPage() {
  const [identifier, setIdentifier] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')

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
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-obsidian">
      <form onSubmit={handleSubmit} className="max-w-md w-full card p-8">
        <h1 className="text-3xl font-bold mb-2 text-pearl">Afritech Online</h1>
        <p className="text-slate mb-8 text-sm">ISP Management Platform</p>
        {error && (
          <div className="mb-4 p-3 rounded bg-red-900/20 border border-red-500/30 text-red-400 text-sm">
            {error}
          </div>
        )}
        <div className="mb-4">
          <label className="block text-sm font-medium mb-2 text-pearl">Username or Email</label>
          <input
            type="text"
            value={identifier}
            onChange={(e) => setIdentifier(e.target.value)}
            className="input-field"
            required
          />
        </div>
        <div className="mb-6">
          <label className="block text-sm font-medium mb-2 text-pearl">Password</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="input-field"
            required
          />
        </div>
        <button type="submit" className="btn-primary w-full">
          Login
        </button>
      </form>
    </div>
  )
}
