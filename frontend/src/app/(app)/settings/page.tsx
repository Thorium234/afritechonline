'use client'

import PageHeader from '@/components/PageHeader'
import { useAuth } from '@/lib/auth-context'

export default function SettingsPage() {
  const { user } = useAuth()
  return (
    <div className="max-w-2xl">
      <PageHeader title="Settings" subtitle="Your account information" />

      <div className="card p-5 mb-4">
        <div className="text-sm font-semibold mb-3">Account</div>
        <div className="grid sm:grid-cols-2 gap-4 text-sm">
          <div>
            <div className="text-xs text-[var(--text-mute)]">Username</div>
            <div className="font-medium mt-0.5">{user?.username}</div>
          </div>
          <div>
            <div className="text-xs text-[var(--text-mute)]">Email</div>
            <div className="font-medium mt-0.5">{user?.email}</div>
          </div>
          <div>
            <div className="text-xs text-[var(--text-mute)]">Role</div>
            <div className="font-medium mt-0.5 capitalize">{(user?.role || '').toLowerCase().replace('_', ' ')}</div>
          </div>
          <div>
            <div className="text-xs text-[var(--text-mute)]">Status</div>
            <div className="font-medium mt-0.5">{user?.is_active ? 'Active' : 'Disabled'}</div>
          </div>
        </div>
      </div>

      <div className="card p-5">
        <div className="text-sm font-semibold mb-2">API</div>
        <p className="text-sm text-[var(--text-dim)]">
          The frontend talks to the Go backend at <code className="text-[var(--text)] bg-[var(--bg-elev-2)] px-1.5 py-0.5 rounded">/api/v1/*</code>.
          Set <code className="text-[var(--text)] bg-[var(--bg-elev-2)] px-1.5 py-0.5 rounded">NEXT_PUBLIC_API_URL</code> in <code className="text-[var(--text)] bg-[var(--bg-elev-2)] px-1.5 py-0.5 rounded">.env</code> to point at a different backend.
        </p>
      </div>
    </div>
  )
}
