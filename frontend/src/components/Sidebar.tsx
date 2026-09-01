'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useEffect } from 'react'
import { cn } from '@/lib/utils'

interface IconProps {
  d: string
}
function Icon({ d }: IconProps) {
  return (
    <svg className="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
      <path strokeLinecap="round" strokeLinejoin="round" d={d} />
    </svg>
  )
}

export interface NavItem {
  href: string
  label: string
  icon: string
  group?: string
}

const NAV: NavItem[] = [
  { href: '/dashboard', label: 'Overview', icon: 'M3 12l9-9 9 9M5 10v10a1 1 0 001 1h3v-6h6v6h3a1 1 0 001-1V10', group: 'Operations' },
  { href: '/customers', label: 'Customers', icon: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z', group: 'Operations' },
  { href: '/subscriptions', label: 'Subscriptions', icon: 'M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15', group: 'Operations' },
  { href: '/packages', label: 'Packages', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4', group: 'Operations' },
  { href: '/billing', label: 'Invoices', icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 7h6m-6 4h6', group: 'Billing' },
  { href: '/payments', label: 'Payments', icon: 'M3 10h18M5 6h14a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2zm3 6h2m4 0h2', group: 'Billing' },
  { href: '/reports', label: 'Reports', icon: 'M3 3v18h18M7 14l3-3 4 4 5-6', group: 'Insights' },
  { href: '/routers', label: 'Routers', icon: 'M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01M2 8.82a15 15 0 0120 0M5 12.859a10 10 0 0114 0', group: 'Network' },
  { href: '/settings', label: 'Settings', icon: 'M10.325 4.317a1.724 1.724 0 013.35 0 1.724 1.724 0 002.573 1.066 1.724 1.724 0 012.37 2.37 1.724 1.724 0 001.065 2.572 1.724 1.724 0 010 3.35 1.724 1.724 0 00-1.066 2.573 1.724 1.724 0 01-2.37 2.37 1.724 1.724 0 00-2.572 1.065 1.724 1.724 0 01-3.35 0 1.724 1.724 0 00-2.573-1.066 1.724 1.724 0 01-2.37-2.37 1.724 1.724 0 00-1.065-2.572 1.724 1.724 0 010-3.35 1.724 1.724 0 001.066-2.573 1.724 1.724 0 012.37-2.37 1.724 1.724 0 002.572-1.065zM12 15a3 3 0 100-6 3 3 0 000 6z', group: 'Account' },
]

interface SidebarProps {
  open: boolean
  onClose: () => void
}

export default function Sidebar({ open, onClose }: SidebarProps) {
  const pathname = usePathname()

  useEffect(() => {
    onClose()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pathname])

  const groups = NAV.reduce<Record<string, NavItem[]>>((acc, item) => {
    const g = item.group || 'Other'
    acc[g] = acc[g] || []
    acc[g].push(item)
    return acc
  }, {})

  return (
    <>
      {/* Mobile overlay */}
      <div
        className={cn(
          'fixed inset-0 z-40 bg-black/60 backdrop-blur-sm transition-opacity lg:hidden',
          open ? 'opacity-100' : 'pointer-events-none opacity-0'
        )}
        onClick={onClose}
      />
      <aside
        className={cn(
          'fixed lg:sticky top-0 left-0 z-50 h-screen w-64 shrink-0',
          'border-r border-[var(--border)] bg-[var(--bg)]',
          'transition-transform duration-200 ease-out lg:translate-x-0',
          open ? 'translate-x-0' : '-translate-x-full'
        )}
      >
        <div className="flex h-full flex-col">
          <div className="flex items-center gap-3 px-5 py-5">
            <div className="grid h-9 w-9 place-items-center rounded-xl bg-[var(--accent-soft)] border border-[var(--accent)]/30">
              <svg className="w-4 h-4 text-[var(--accent)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div className="leading-tight">
              <div className="text-[15px] font-semibold tracking-tight">Afritech</div>
              <div className="text-[11px] uppercase tracking-[0.18em] text-[var(--text-mute)]">Online</div>
            </div>
          </div>

          <nav className="flex-1 overflow-y-auto px-3 py-2 space-y-4">
            {Object.entries(groups).map(([group, items]) => (
              <div key={group}>
                <div className="px-3 pb-1.5 text-[10px] font-semibold uppercase tracking-[0.18em] text-[var(--text-mute)]">
                  {group}
                </div>
                <div className="space-y-0.5">
                  {items.map((item) => {
                    const active = pathname === item.href || pathname?.startsWith(item.href + '/')
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        className={cn('sidebar-link', active && 'active')}
                      >
                        <Icon d={item.icon} />
                        <span>{item.label}</span>
                      </Link>
                    )
                  })}
                </div>
              </div>
            ))}
          </nav>

          <div className="px-4 py-4 border-t border-[var(--border)]">
            <div className="rounded-lg bg-[var(--bg-elev)] border border-[var(--border)] p-3 text-xs text-[var(--text-dim)]">
              <div className="flex items-center gap-2 text-[var(--text)] font-medium">
                <span className="dot dot-success" />
                All systems operational
              </div>
              <div className="mt-1 text-[var(--text-mute)]">v0.1.0 · build local</div>
            </div>
          </div>
        </div>
      </aside>
    </>
  )
}
