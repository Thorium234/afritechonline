'use client'

import { cn } from '@/lib/utils'

export type StatusKind =
  | 'success'
  | 'warning'
  | 'danger'
  | 'info'
  | 'neutral'
  | 'inactive'

interface StatusBadgeProps {
  status: string
  className?: string
}

const RULES: Array<{ test: (s: string) => boolean; kind: StatusKind }> = [
  { test: (s) => /^(active|paid|completed|online|succeeded|success)$/i.test(s), kind: 'success' },
  { test: (s) => /^(pending|processing|partial|inprogress|in.progress)$/i.test(s), kind: 'warning' },
  { test: (s) => /^(expired|overdue|failed|cancelled|canceled|suspended|offline|error)$/i.test(s), kind: 'danger' },
  { test: (s) => /^(inactive|draft|unknown)$/i.test(s), kind: 'neutral' },
  { test: (s) => /^(new|info)$/i.test(s), kind: 'info' },
]

export default function StatusBadge({ status, className }: StatusBadgeProps) {
  const rule = RULES.find((r) => r.test(status))
  const kind: StatusKind = rule?.kind || 'neutral'
  const dotClass = {
    success: 'dot-success',
    warning: 'dot-warning',
    danger: 'dot-danger',
    info: 'dot-info',
    neutral: 'bg-[var(--text-mute)]',
    inactive: 'bg-[var(--text-mute)]',
  }[kind]
  return (
    <span className={cn('badge', `badge-${kind === 'inactive' ? 'neutral' : kind}`, className)}>
      <span className={cn('dot', dotClass)} />
      {status}
    </span>
  )
}
