import { cn } from '@/lib/utils'

interface TableProps {
  columns: Array<{
    key: string
    label: string
    width?: string
    sortable?: boolean
  }>
  rows: Array<Record<string, any>>
  onSort?: (key: string) => void
  sortKey?: string
  sortDir?: 'asc' | 'desc'
  loading?: boolean
  empty?: React.ReactNode
}

export function Table({ columns, rows, onSort, sortKey, sortDir, loading, empty }: TableProps) {
  if (loading) {
    return (
      <div className="card">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[var(--border)]">
              {columns.map((col) => (
                <th
                  key={col.key}
                  className="px-4 py-3 text-left text-xs font-semibold text-[var(--text-dim)] uppercase"
                  style={{ width: col.width }}
                >
                  {col.label}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {[...Array(5)].map((_, i) => (
              <tr key={i} className="border-b border-[var(--border)]">
                {columns.map((col) => (
                  <td key={col.key} className="px-4 py-3">
                    <div className="h-4 bg-[var(--bg-elev)] rounded animate-pulse" />
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    )
  }

  if (!rows.length) {
    return empty || (
      <div className="card text-center py-12">
        <p className="text-[var(--text-dim)]">No data found</p>
      </div>
    )
  }

  return (
    <div className="card overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-[var(--border)]">
            {columns.map((col) => (
              <th
                key={col.key}
                className={cn(
                  'px-4 py-3 text-left text-xs font-semibold text-[var(--text-dim)] uppercase',
                  col.sortable && 'cursor-pointer hover:bg-[var(--bg-elev)]'
                )}
                style={{ width: col.width }}
                onClick={() => col.sortable && onSort?.(col.key)}
              >
                <div className="flex items-center gap-2">
                  {col.label}
                  {col.sortable && sortKey === col.key && (
                    <svg className={cn('w-4 h-4', sortDir === 'desc' && 'rotate-180')} fill="currentColor" viewBox="0 0 20 20">
                      <path d="M3 3a1 1 0 000 2h11a1 1 0 100-2H3zM3 7a1 1 0 000 2h5a1 1 0 000-2H3zM3 11a1 1 0 100 2h4a1 1 0 100-2H3zM13 16a1 1 0 102 0v-5.586l1.293 1.293a1 1 0 001.414-1.414l-3-3a1 1 0 00-1.414 0l-3 3a1 1 0 101.414 1.414L13 10.414V16z" />
                    </svg>
                  )}
                </div>
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.map((row, idx) => (
            <tr key={idx} className="border-b border-[var(--border)] hover:bg-[var(--bg-elev)] transition">
              {columns.map((col) => (
                <td key={col.key} className="px-4 py-3 text-sm text-[var(--text)]">
                  {row[col.key]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  isLoading?: boolean
}

export function Pagination({ currentPage, totalPages, onPageChange, isLoading }: PaginationProps) {
  const pages = []
  const maxVisible = 7
  const halfVisible = Math.floor(maxVisible / 2)
  
  let start = Math.max(1, currentPage - halfVisible)
  let end = Math.min(totalPages, currentPage + halfVisible)
  
  if (end - start + 1 < maxVisible) {
    if (start === 1) {
      end = Math.min(totalPages, start + maxVisible - 1)
    } else {
      start = Math.max(1, end - maxVisible + 1)
    }
  }
  
  if (start > 1) {
    pages.push(1)
    if (start > 2) pages.push('...')
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  if (end < totalPages) {
    if (end < totalPages - 1) pages.push('...')
    pages.push(totalPages)
  }

  return (
    <div className="flex items-center justify-center gap-2">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1 || isLoading}
        className="p-2 hover:bg-[var(--bg-elev)] disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      
      {pages.map((page, idx) =>
        page === '...' ? (
          <span key={`${page}-${idx}`} className="px-2">
            ...
          </span>
        ) : (
          <button
            key={page}
            onClick={() => onPageChange(page as number)}
            disabled={isLoading}
            className={cn(
              'px-3 py-2 rounded-lg transition',
              currentPage === page
                ? 'bg-[var(--accent)] text-white'
                : 'hover:bg-[var(--bg-elev)]'
            )}
          >
            {page}
          </button>
        )
      )}
      
      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages || isLoading}
        className="p-2 hover:bg-[var(--bg-elev)] disabled:opacity-50 disabled:cursor-not-allowed rounded-lg transition"
      >
        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  )
}
