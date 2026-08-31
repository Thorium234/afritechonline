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
    <div className="flex min-h-screen bg-gray-50">
      <aside className="w-64 bg-white shadow">
        <div className="p-6">
          <h1 className="text-xl font-bold text-blue-600">Afritech Online</h1>
        </div>
        <nav className="mt-6">
          {navItems.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`block px-6 py-3 text-sm font-medium ${
                pathname === item.href
                  ? 'text-blue-600 bg-blue-50 border-r-4 border-blue-600'
                  : 'text-gray-600 hover:bg-gray-50'
              }`}
            >
              {item.label}
            </Link>
          ))}
        </nav>
        <div className="absolute bottom-0 w-64 p-6">
          <button
            onClick={() => {
              localStorage.removeItem('access_token')
              localStorage.removeItem('refresh_token')
              window.location.href = '/login'
            }}
            className="text-sm text-red-600 hover:text-red-700"
          >
            Logout
          </button>
        </div>
      </aside>
      <main className="flex-1">{children}</main>
    </div>
  )
}
