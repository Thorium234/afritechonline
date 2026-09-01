'use client'

import { cn } from '@/lib/utils'

interface StatCardProps {
  label: string
  value: string | number
  hint?: string
  trend?: { value: number; positive?: boolean }
  icon?: React.ReactNode
  accent?: 'lime' | 'blue' | 'amber' | 'rose' | 'violet'
  className?: string
}

const ACCENTS: Record<NonNullable<StatCardProps['accent']>, string> = {
  lime: 'from-[rgba(210,255,42,0.18)] to-transparent text-[var(--accent)]',
  blue: 'from-[rgba(88,166,255,0.18)] to-transparent text-[#79b8ff]',
  amber: 'from-[rgba(210,153,34,0.18)] to-transparent text-[#f0c674]',
  rose: 'from-[rgba(248,81,73,0.18)] to-transparent text-[#ff8580]',
  violet: 'from-[rgba(163,113,247,0.18)] to-transparent text-[#c39bf3]',
}

export default function StatCard({
  label,
  value,
  hint,
  trend,
  icon,
  accent = 'lime',
  className,
}: StatCardProps) {
  return (
    <div className={cn('stat-card', className)}>
      <div className="flex items-start justify-between gap-3">
        <div>
          <div className="text-xs font-medium uppercase tracking-wider text-[var(--text-mute)]">
            {label}
          </div>
          <div className="mt-1.5 text-2xl font-semibold text-[var(--text)] tracking-tight">
            {value}
          </div>
          {hint && (
            <div className="mt-1 text-xs text-[var(--text-dim)]">{hint}</div>
          )}
        </div>
        {icon && (
          <div
            className={cn(
              'grid h-10 w-10 place-items-center rounded-xl bg-gradient-to-br',
              ACCENTS[accent]
            )}
          >
            {icon}
          </div>
        )}
      </div>
      {trend && (
        <div className="mt-3 flex items-center gap-1 text-xs">
          <span
            className={cn(
              'inline-flex items-center gap-0.5 font-medium',
              trend.positive ? 'text-[#6fdc8c]' : 'text-[#ff8580]'
            )}
          >
            <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}>
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d={trend.positive ? 'M5 15l7-7 7 7' : 'M19 9l-7 7-7-7'}
              />
            </svg>
            {Math.abs(trend.value)}%
          </span>
          <span className="text-[var(--text-mute)]">vs last period</span>
        </div>
      )}
    </div>
  )
}
