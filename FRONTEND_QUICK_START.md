# Frontend Quick Start Guide - Building Forms & Pages

This guide covers how to build new forms and pages using the component library and patterns established in the project.

---

## 🏗️ Basic Page Structure

### List Page Template

```tsx
'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { Button } from '@/components/Button'
import { Input } from '@/components/Form'
import { Table, Pagination } from '@/components/Table'
import { EmptyState } from '@/components/Common'
import { useToast } from '@/components/Toast'
import { useApi } from '@/lib/useApi'
import { ApiError } from '@/lib/types'

export default function ListPage() {
  const toast = useToast()
  const [page, setPage] = useState(1)
  const [search, setSearch] = useState('')
  const [sortKey, setSortKey] = useState('created_at')
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('desc')

  const url = `/api/v1/items?page=${page}&search=${search}&sort=${sortKey}&dir=${sortDir}`
  const { data, loading, error } = useApi(url, {
    onError: (err) => toast.error(err.message),
  })

  const handleSort = (key: string) => {
    if (sortKey === key) {
      setSortDir(sortDir === 'asc' ? 'desc' : 'asc')
    } else {
      setSortKey(key)
      setSortDir('asc')
    }
  }

  return (
    <div>
      <PageHeader
        title="Items"
        subtitle="Manage your items"
        action={
          <Link href="/items/new">
            <Button>New Item</Button>
          </Link>
        }
      />

      {/* Search Bar */}
      <div className="mb-6">
        <Input
          type="search"
          placeholder="Search items..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value)
            setPage(1)
          }}
        />
      </div>

      {/* Error State */}
      {error && (
        <div className="mb-6">
          <Alert type="error" title="Error">
            {error.message}
          </Alert>
        </div>
      )}

      {/* Table */}
      <Table
        columns={[
          { key: 'name', label: 'Name', sortable: true },
          { key: 'status', label: 'Status' },
          { key: 'created_at', label: 'Created', sortable: true },
          {
            key: 'actions',
            label: 'Actions',
            render: (row) => (
              <div className="flex gap-2">
                <Link href={`/items/${row.id}/edit`}>
                  <Button variant="secondary" size="sm">Edit</Button>
                </Link>
                <Button
                  variant="danger"
                  size="sm"
                  onClick={() => handleDelete(row.id)}
                >
                  Delete
                </Button>
              </div>
            ),
          },
        ]}
        rows={data?.items || []}
        sortKey={sortKey}
        sortDir={sortDir}
        onSort={handleSort}
        loading={loading}
        empty={
          <EmptyState
            title="No items found"
            description="Create your first item to get started"
            action={{ label: 'New Item', href: '/items/new' }}
          />
        }
      />

      {/* Pagination */}
      {data && (
        <div className="mt-6">
          <Pagination
            currentPage={page}
            totalPages={Math.ceil(data.total / data.page_size)}
            onPageChange={setPage}
            isLoading={loading}
          />
        </div>
      )}
    </div>
  )
}
```

### Create/Edit Form Template

```tsx
'use client'

import { useState, type FormEvent } from 'react'
import { useRouter, useParams } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { Button } from '@/components/Button'
import { FormField, Input, Select, TextArea } from '@/components/Form'
import { Alert } from '@/components/Alert'
import { useToast } from '@/components/Toast'
import { useApi, useApiMutation } from '@/lib/useApi'
import { ApiError } from '@/lib/types'

interface FormData {
  name: string
  email: string
  status: 'ACTIVE' | 'INACTIVE'
  notes: string
}

export default function FormPage() {
  const router = useRouter()
  const params = useParams()
  const toast = useToast()
  const isEdit = !!params.id

  // Fetch existing data if editing
  const { data: existingItem } = useApi(
    isEdit ? `/api/v1/items/${params.id}` : '',
    { skip: !isEdit }
  )

  const [form, setForm] = useState<FormData>({
    name: '',
    email: '',
    status: 'ACTIVE',
    notes: '',
  })

  const [errors, setErrors] = useState<Record<string, string>>({})
  const { execute, loading } = useApiMutation()

  // Populate form when data loads
  useEffect(() => {
    if (existingItem) {
      setForm(existingItem)
    }
  }, [existingItem])

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {}

    if (!form.name.trim()) {
      newErrors.name = 'Name is required'
    }

    if (!form.email.trim()) {
      newErrors.email = 'Email is required'
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
      newErrors.email = 'Invalid email format'
    }

    setErrors(newErrors)
    return Object.keys(newErrors).length === 0
  }

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()

    if (!validateForm()) {
      toast.error('Please fix the errors below')
      return
    }

    const method = isEdit ? 'PUT' : 'POST'
    const url = isEdit ? `/api/v1/items/${params.id}` : '/api/v1/items'

    const result = await execute(method, url, form)

    if (result) {
      toast.success(
        isEdit ? 'Item updated successfully' : 'Item created successfully'
      )
      router.push('/items')
    } else {
      toast.error('Failed to save item')
    }
  }

  return (
    <div>
      <PageHeader
        title={isEdit ? 'Edit Item' : 'New Item'}
        subtitle={isEdit ? 'Update item details' : 'Create a new item'}
      />

      <div className="max-w-2xl">
        <form onSubmit={onSubmit} className="card p-6 space-y-6">
          <FormField label="Name" error={errors.name} required>
            <Input
              type="text"
              placeholder="Item name"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              error={!!errors.name}
            />
          </FormField>

          <FormField label="Email" error={errors.email} required>
            <Input
              type="email"
              placeholder="email@example.com"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              error={!!errors.email}
            />
          </FormField>

          <FormField label="Status" required>
            <Select
              value={form.status}
              onChange={(e) => setForm({ ...form, status: e.target.value as any })}
              options={[
                { value: 'ACTIVE', label: 'Active' },
                { value: 'INACTIVE', label: 'Inactive' },
              ]}
            />
          </FormField>

          <FormField label="Notes">
            <TextArea
              placeholder="Additional notes..."
              value={form.notes}
              onChange={(e) => setForm({ ...form, notes: e.target.value })}
            />
          </FormField>

          {/* Form Actions */}
          <div className="flex gap-3 justify-end pt-4 border-t border-[var(--border)]">
            <Link href="/items">
              <Button variant="secondary">Cancel</Button>
            </Link>
            <Button type="submit" variant="primary" loading={loading}>
              {isEdit ? 'Save Changes' : 'Create'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}
```

---

## 🎯 Common Patterns

### Handling Validation Errors

```tsx
const onSubmit = async (e: FormEvent) => {
  e.preventDefault()
  setErrors({})

  const result = await execute('POST', '/api/v1/items', form)

  if (!result) {
    if (error instanceof ApiError) {
      // Handle validation errors
      if (error.statusCode === 422 && error.fields) {
        setErrors(error.fields)
        toast.error('Please fix the validation errors')
      } else {
        toast.error(error.message)
      }
    }
  }
}
```

### Confirmation Dialog for Delete

```tsx
const [showConfirm, setShowConfirm] = useState(false)
const [deleteId, setDeleteId] = useState<string | null>(null)

const handleDeleteClick = (id: string) => {
  setDeleteId(id)
  setShowConfirm(true)
}

const handleConfirmDelete = async () => {
  if (!deleteId) return

  const result = await execute('DELETE', `/api/v1/items/${deleteId}`)
  if (result) {
    toast.success('Item deleted successfully')
    // Refresh list
    refetch?.()
  }
  setShowConfirm(false)
}

// In JSX:
<ConfirmDialog
  isOpen={showConfirm}
  title="Delete Item"
  description="This action cannot be undone."
  confirmText="Delete"
  variant="danger"
  onConfirm={handleConfirmDelete}
  onCancel={() => setShowConfirm(false)}
/>
```

### Loading States

```tsx
// Loading skeleton while fetching
if (loading) {
  return <TableSkeleton />
}

// Loading button while submitting
<Button loading={isSubmitting}>
  Save
</Button>

// Show loading spinner overlay
{isSubmitting && (
  <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div className="animate-spin">
      <svg className="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        {/* spinner SVG */}
      </svg>
    </div>
  </div>
)}
```

### Date Input

```tsx
// For date fields in forms
<FormField label="Start Date" required>
  <Input
    type="date"
    value={form.startDate}
    onChange={(e) => setForm({ ...form, startDate: e.target.value })}
  />
</FormField>

// Alternative: Add date picker library
// npm install react-datepicker
import DatePicker from 'react-datepicker'
<DatePicker
  selected={form.startDate}
  onChange={(date) => setForm({ ...form, startDate: date })}
/>
```

### Dependent Fields

```tsx
// Show field based on another field's value
{form.paymentMethod === 'MPESA' && (
  <FormField label="M-Pesa Phone Number" required>
    <Input
      type="tel"
      value={form.mpesaPhone}
      onChange={(e) => setForm({ ...form, mpesaPhone: e.target.value })}
    />
  </FormField>
)}

{form.paymentMethod === 'MANUAL' && (
  <FormField label="Reference" required>
    <Input
      type="text"
      value={form.reference}
      onChange={(e) => setForm({ ...form, reference: e.target.value })}
    />
  </FormField>
)}
```

---

## 🔑 Key Files to Update Next

### Priority 1: Core CRUD Forms
1. **customers/[id]/edit/page.tsx** - Edit customer form
2. **subscriptions/new/page.tsx** - Enhanced subscription creation
3. **packages/[id]/edit/page.tsx** - Edit package form
4. **payments/new/page.tsx** - Create payment form

### Priority 2: Delete Operations
1. Add delete buttons to all list pages
2. Implement ConfirmDialog for each
3. Handle delete API calls

### Priority 3: Advanced Features
1. **payments/mpesa/page.tsx** - M-Pesa payment initiation
2. **routers/new/page.tsx** - Create router form
3. **reports/page.tsx** - Enhance with charts

---

## 📱 Responsive Design Tips

```tsx
// Use Tailwind responsive classes
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* Cards adjust based on screen size */}
</div>

// Hide elements on mobile
<div className="hidden lg:block">
  {/* Only visible on large screens */}
</div>

// Stack on mobile
<div className="flex flex-col md:flex-row gap-4">
  {/* Column on mobile, row on desktop */}
</div>
```

---

## 🧪 Testing Forms

```tsx
import { render, screen, fireEvent } from '@testing-library/react'
import FormPage from './page'

describe('FormPage', () => {
  it('shows validation errors', async () => {
    render(<FormPage />)
    
    const submitButton = screen.getByRole('button', { name: /save/i })
    fireEvent.click(submitButton)
    
    expect(screen.getByText(/name is required/i)).toBeInTheDocument()
  })

  it('submits form on valid input', async () => {
    render(<FormPage />)
    
    fireEvent.change(screen.getByLabelText(/name/i), {
      target: { value: 'John' },
    })
    
    fireEvent.click(screen.getByRole('button', { name: /save/i }))
    
    // Wait for API call
    await waitFor(() => {
      expect(screen.getByText(/successfully/i)).toBeInTheDocument()
    })
  })
})
```

---

## 🚀 Component Usage Examples

### Button Variants

```tsx
<Button variant="primary">Primary</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="danger">Delete</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="outline">Outline</Button>
```

### Modal Usage

```tsx
import { Modal } from '@/components/Modal'

const [isOpen, setIsOpen] = useState(false)

return (
  <>
    <Button onClick={() => setIsOpen(true)}>Open Modal</Button>
    
    <Modal
      isOpen={isOpen}
      onClose={() => setIsOpen(false)}
      title="Confirm Action"
      footer={
        <>
          <Button variant="secondary" onClick={() => setIsOpen(false)}>
            Cancel
          </Button>
          <Button onClick={handleConfirm}>
            Confirm
          </Button>
        </>
      }
    >
      Are you sure?
    </Modal>
  </>
)
```

### Toast Notifications

```tsx
import { useToast } from '@/components/Toast'

const toast = useToast()

// In your component
toast.success('Operation successful!')
toast.error('Something went wrong')
toast.warning('Are you sure?')
toast.info('FYI: Something happened')
```

---

## 🎨 CSS Variables for Styling

All custom styles should use CSS variables:

```tsx
<div className="bg-[var(--bg)] text-[var(--text)]">
  <p className="text-[var(--text-dim)]">Dimmed text</p>
  <input className="border-[var(--border)]" />
</div>
```

---

## 🔄 API Integration Pattern

```tsx
// 1. Define form state
const [form, setForm] = useState(initialState)

// 2. Add validation
const validateForm = (): boolean => { ... }

// 3. Create mutation hook
const { execute, loading, error } = useApiMutation()

// 4. Handle submit
const onSubmit = async (e: FormEvent) => {
  e.preventDefault()
  if (!validateForm()) return
  
  const result = await execute('POST', '/api/v1/items', form)
  if (result) {
    toast.success('Success!')
    router.push('/items')
  }
}

// 5. Render form
return (
  <form onSubmit={onSubmit}>
    {/* Form fields */}
    <Button loading={loading}>Submit</Button>
  </form>
)
```

---

## ⚠️ Common Mistakes to Avoid

1. **Don't forget error handling** - Always show error messages to users
2. **Don't skip loading states** - Use loading prop on buttons/forms
3. **Don't make API calls directly** - Use useApi/useApiMutation hooks
4. **Don't hardcode styles** - Use CSS variables and Tailwind classes
5. **Don't forget validation** - Validate before submitting to API
6. **Don't forget accessibility** - Use proper labels and ARIA attributes
7. **Don't make untyped API responses** - Always define types in types.ts

---

## 📚 Resources

- [Component Reference](FRONTEND_COMPONENTS_REFERENCE.md)
- [Implementation Plan](FRONTEND_IMPLEMENTATION_PLAN.md)
- [TypeScript Types](frontend/src/lib/types.ts)
- [API Integration](frontend/src/lib/api.ts)
- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)

---

**This guide should help you quickly build forms and pages following established patterns!**
