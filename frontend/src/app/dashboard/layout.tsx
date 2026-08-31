'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

const navItems = [
  { href: '/dashboard', label: 'Dashboard' },
  { href: '/customers', label: 'Customers' },
  { href: '/packages', label: 'Packages' },
  { href: '/subscriptions', label: 'Subscriptions' },
  { href: '/billing', label: 'Billing' },
]

export default function DashboardLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()

  return (
    <div className="flex min-h-screen bg-obsidian">
      <aside className="w-64 bg-obsidian-light border-r border-white/10 flex flex-col">
        <div className="p-6">
          <h1 className="text-xl font-bold text-acid-lime">Afritech Online</h1>
          <p className="text-xs text-slate mt-1">ISP Management</p>
        </div>
        <nav className="flex-1 mt-6">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`nav-item block px-6 py-3 text-sm font-medium ${
                pathname === item.href ? 'active' : ''
              }`}
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="p-6">
          <button
            onClick={() => {
              localStorage.removeItem('access_token')
              localStorage.removeItem('refresh_token')
              window.location.href = '/login'
            }}
            className="text-sm text-red-400 hover:text-red-300 transition-colors"
          >
            Logout
          </button>
        </div>
      </aside>
      <main className="flex-1 bg-obsidian">{children}</main>
    </div>
  )
}
