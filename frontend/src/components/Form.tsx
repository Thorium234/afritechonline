import { cn } from '@/lib/utils'

interface FormFieldProps {
  label: string
  error?: string
  required?: boolean
  children: React.ReactNode
  hint?: string
}

export function FormField({ label, error, required, children, hint }: FormFieldProps) {
  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-[var(--text)]">
        {label}
        {required && <span className="text-[var(--danger)] ml-1">*</span>}
      </label>
      {children}
      {error && <p className="text-xs text-[var(--danger)]">{error}</p>}
      {hint && !error && <p className="text-xs text-[var(--text-mute)]">{hint}</p>}
    </div>
  )
}

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string
}

export function Input({ error, className, ...props }: InputProps) {
  return (
    <input
      className={cn(
        'w-full px-3 py-2 rounded-lg border bg-[var(--bg)] text-[var(--text)] placeholder:text-[var(--text-mute)]',
        error ? 'border-[var(--danger)]' : 'border-[var(--border)] focus:border-[var(--accent)]',
        'focus:outline-none focus:ring-2 focus:ring-[var(--accent)]/20 transition',
        className
      )}
      {...props}
    />
  )
}

interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  error?: string
  options: Array<{ value: string | number; label: string }>
}

export function Select({ error, options, className, ...props }: SelectProps) {
  return (
    <select
      className={cn(
        'w-full px-3 py-2 rounded-lg border bg-[var(--bg)] text-[var(--text)]',
        error ? 'border-[var(--danger)]' : 'border-[var(--border)] focus:border-[var(--accent)]',
        'focus:outline-none focus:ring-2 focus:ring-[var(--accent)]/20 transition',
        className
      )}
      {...props}
    >
      {options.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.label}
        </option>
      ))}
    </select>
  )
}

interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  error?: string
}

export function TextArea({ error, className, ...props }: TextAreaProps) {
  return (
    <textarea
      className={cn(
        'w-full px-3 py-2 rounded-lg border bg-[var(--bg)] text-[var(--text)] placeholder:text-[var(--text-mute)] resize-vertical',
        error ? 'border-[var(--danger)]' : 'border-[var(--border)] focus:border-[var(--accent)]',
        'focus:outline-none focus:ring-2 focus:ring-[var(--accent)]/20 transition',
        className
      )}
      {...props}
    />
  )
}

interface CheckboxProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string
  error?: string
}

export function Checkbox({ label, error, className, ...props }: CheckboxProps) {
  return (
    <div className="flex items-center gap-2">
      <input
        type="checkbox"
        className={cn(
          'w-5 h-5 rounded border-[var(--border)] accent-[var(--accent)]',
          error && 'border-[var(--danger)]',
          className
        )}
        {...props}
      />
      {label && <label className="text-sm text-[var(--text)]">{label}</label>}
    </div>
  )
}
