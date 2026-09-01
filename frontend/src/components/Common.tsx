'use client'

import React, { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/Button'
import { Modal } from '@/components/Modal'
import { Alert } from '@/components/Alert'
import { cn } from '@/lib/utils'

interface ConfirmDialogProps {
  isOpen: boolean
  title: string
  description: string
  confirmText?: string
  cancelText?: string
  variant?: 'danger' | 'warning' | 'info'
  isLoading?: boolean
  onConfirm: () => void | Promise<void>
  onCancel: () => void
}

export function ConfirmDialog({
  isOpen,
  title,
  description,
  confirmText = 'Confirm',
  cancelText = 'Cancel',
  variant = 'danger',
  isLoading = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  const [error, setError] = useState<string | null>(null)

  const handleConfirm = async () => {
    try {
      setError(null)
      await onConfirm()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onCancel}
      title={title}
      footer={
        <>
          <Button variant="secondary" onClick={onCancel} disabled={isLoading}>
            {cancelText}
          </Button>
          <Button
            variant={variant}
            onClick={handleConfirm}
            loading={isLoading}
          >
            {confirmText}
          </Button>
        </>
      }
    >
      <div className="space-y-4">
        <p className="text-[var(--text)]">{description}</p>
        {error && (
          <Alert type="error" onDismiss={() => setError(null)}>
            {error}
          </Alert>
        )}
      </div>
    </Modal>
  )
}

interface FormActionProps {
  label: string
  href?: string
  onClick?: () => void
  variant?: 'primary' | 'secondary' | 'danger'
  loading?: boolean
  disabled?: boolean
  icon?: React.ReactNode
}

export function FormActions({ actions, onCancel }: {
  actions: FormActionProps[]
  onCancel?: () => void
}) {
  return (
    <div className="flex gap-3 justify-end pt-4 border-t border-[var(--border)]">
      {onCancel && (
        <Link href="#">
          <Button variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        </Link>
      )}
      {actions.map((action, idx) => (
        <Button
          key={idx}
          variant={action.variant}
          onClick={action.onClick}
          disabled={action.disabled}
          loading={action.loading}
        >
          {action.icon}
          {action.label}
        </Button>
      ))}
    </div>
  )
}

interface DataRowProps {
  label: string
  value: React.ReactNode
  copyable?: boolean
}

export function DataRow({ label, value, copyable }: DataRowProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    if (typeof value === 'string') {
      navigator.clipboard.writeText(value)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  return (
    <div className="flex justify-between py-2 border-b border-[var(--border)]">
      <span className="text-[var(--text-dim)]">{label}</span>
      <div className="flex items-center gap-2">
        <span className="text-[var(--text)] font-medium">{value}</span>
        {copyable && typeof value === 'string' && (
          <button
            onClick={handleCopy}
            className="p-1 hover:bg-[var(--bg-elev)] rounded transition"
            title="Copy to clipboard"
          >
            {copied ? (
              <svg className="w-4 h-4 text-[var(--accent)]" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
              </svg>
            ) : (
              <svg className="w-4 h-4 text-[var(--text-dim)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            )}
          </button>
        )}
      </div>
    </div>
  )
}

interface StatProps {
  label: string
  value: string | number
  trend?: number // Positive for up, negative for down
  trendLabel?: string
}

export function Stat({ label, value, trend, trendLabel }: StatProps) {
  return (
    <div className="card p-4">
      <p className="text-sm text-[var(--text-dim)] mb-2">{label}</p>
      <div className="flex items-end justify-between">
        <p className="text-2xl font-bold text-[var(--text)]">{value}</p>
        {trend !== undefined && (
          <div className={cn(
            'text-sm font-medium flex items-center gap-1',
            trend > 0 ? 'text-[var(--accent)]' : trend < 0 ? 'text-[var(--danger)]' : 'text-[var(--text-dim)]'
          )}>
            {trend > 0 && <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M12 7a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0V8.414l-4.293 4.293a1 1 0 01-1.414-1.414L13.586 7H12z" clipRule="evenodd" />
            </svg>}
            {trend < 0 && <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M12 13a1 1 0 110 2H7a1 1 0 01-1-1V9a1 1 0 112 0v3.586l4.293-4.293a1 1 0 011.414 1.414L9.414 13H12z" clipRule="evenodd" />
            </svg>}
            {Math.abs(trend)}%{trendLabel && ` ${trendLabel}`}
          </div>
        )}
      </div>
    </div>
  )
}

interface EmptyStateProps {
  title: string
  description: string
  action?: {
    label: string
    href?: string
    onClick?: () => void
  }
  icon?: React.ReactNode
}

export function EmptyState({ title, description, action, icon }: EmptyStateProps) {
  return (
    <div className="card p-12 text-center">
      {icon && <div className="mb-4 flex justify-center">{icon}</div>}
      <h3 className="text-lg font-semibold text-[var(--text)] mb-2">{title}</h3>
      <p className="text-[var(--text-dim)] mb-6">{description}</p>
      {action && (
        action.href ? (
          <Link href={action.href}>
            <Button>{action.label}</Button>
          </Link>
        ) : (
          <Button onClick={action.onClick}>{action.label}</Button>
        )
      )}
    </div>
  )
}
