import { useState, useCallback } from 'react'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

interface Toast {
  id: string
  message: string
  type: ToastType
  duration?: number
}

const toasts: Toast[] = []
const listeners: Set<(toasts: Toast[]) => void> = new Set()

function notifyListeners() {
  listeners.forEach((listener) => listener([...toasts]))
}

export function useToast() {
  const [, setToasts] = useState<Toast[]>([])

  const subscribe = useCallback((listener: (toasts: Toast[]) => void) => {
    listeners.add(listener)
    return () => listeners.delete(listener)
  }, [])

  const show = useCallback((message: string, type: ToastType = 'info', duration = 4000) => {
    const id = `toast-${Date.now()}-${Math.random()}`
    const toast: Toast = { id, message, type, duration }
    
    toasts.push(toast)
    notifyListeners()

    if (duration > 0) {
      setTimeout(() => {
        const index = toasts.findIndex((t) => t.id === id)
        if (index > -1) {
          toasts.splice(index, 1)
          notifyListeners()
        }
      }, duration)
    }

    return id
  }, [])

  const success = useCallback((message: string, duration?: number) => show(message, 'success', duration), [show])
  const error = useCallback((message: string, duration?: number) => show(message, 'error', duration), [show])
  const warning = useCallback((message: string, duration?: number) => show(message, 'warning', duration), [show])
  const info = useCallback((message: string, duration?: number) => show(message, 'info', duration), [show])

  const dismiss = useCallback((id: string) => {
    const index = toasts.findIndex((t) => t.id === id)
    if (index > -1) {
      toasts.splice(index, 1)
      notifyListeners()
    }
  }, [])

  return { show, success, error, warning, info, dismiss, subscribe }
}

export function ToastContainer() {
  const [toastList, setToastList] = useState<Toast[]>([])
  const { subscribe, dismiss } = useToast()

  useState(() => {
    return subscribe(setToastList)
  }, [subscribe])

  return (
    <div className="fixed bottom-4 right-4 z-50 space-y-2 max-w-md">
      {toastList.map((toast) => (
        <div
          key={toast.id}
          className={`rounded-lg border px-4 py-3 shadow-lg animate-fade-in flex items-start gap-3 ${
            toast.type === 'success'
              ? 'bg-[#0d3b16] border-[#22c55e]/30 text-[#22c55e]'
              : toast.type === 'error'
              ? 'bg-[#3b0d0d] border-[#ef4444]/30 text-[#ef4444]'
              : toast.type === 'warning'
              ? 'bg-[#3b2d0d] border-[#f59e0b]/30 text-[#f59e0b]'
              : 'bg-[var(--bg-elev)] border-[var(--border)] text-[var(--text)]'
          }`}
        >
          {toast.type === 'success' && (
            <svg className="w-5 h-5 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
            </svg>
          )}
          {toast.type === 'error' && (
            <svg className="w-5 h-5 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
          )}
          <div className="flex-1 min-w-0">
            <p className="text-sm font-medium">{toast.message}</p>
          </div>
          <button
            onClick={() => dismiss(toast.id)}
            className="flex-shrink-0 p-1 hover:bg-white/10 rounded transition"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      ))}
    </div>
  )
}
