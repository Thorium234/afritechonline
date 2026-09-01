'use client'

import { cn } from '@/lib/utils'

interface EmptyStateProps {
  icon?: React.ReactNode
  title: string
  description?: string
  action?: React.ReactNode
  className?: string
}

export default function EmptyState({ icon, title, description, action, className }: EmptyStateProps) {
  return (
    <div className={cn('card empty-state', className)}>
      <div className="empty-state-icon">{icon}</div>
      <div className="text-base font-semibold text-[var(--text)]">{title}</div>
      {description && <p className="mt-1 max-w-sm text-sm text-[var(--text-dim)]">{description}</p>}
      {action && <div className="mt-4">{action}</div>}
    </div>
  )
}
