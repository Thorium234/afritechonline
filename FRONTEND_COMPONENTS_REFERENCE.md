# Frontend Components Reference

## Overview

This guide covers all the frontend components available for building the Afritech Online ISP platform UI. All components follow the established design system using CSS variables and Tailwind CSS.

---

## Core Components

### Button

Versatile button component with multiple variants and sizes.

**Props:**
- `variant`: 'primary' | 'secondary' | 'danger' | 'ghost' | 'outline' (default: 'primary')
- `size`: 'xs' | 'sm' | 'md' | 'lg' (default: 'md')
- `loading`: boolean - Shows spinner and disables button
- `fullWidth`: boolean - Stretches button to full width
- `disabled`: boolean

**Usage:**
```tsx
import { Button } from '@/components/Button'

<Button variant="primary" size="md" onClick={handleClick}>
  Click me
</Button>

<Button loading={isLoading} fullWidth>
  Save Changes
</Button>

<Button variant="danger">Delete</Button>
```

**Variants:**
- **primary**: Blue accent background, main action button
- **secondary**: Outlined button, secondary actions
- **danger**: Red background, destructive actions
- **ghost**: Transparent, text-only button
- **outline**: Border only, subtle appearance

---

### Modal

Dialog/modal component with header, content area, and footer.

**Props:**
- `isOpen`: boolean - Controls modal visibility
- `onClose`: () => void - Called when modal closes
- `title`: string - Modal header title
- `description`: string (optional) - Subtitle in header
- `children`: ReactNode - Modal content
- `footer`: ReactNode (optional) - Footer actions
- `size`: 'sm' | 'md' | 'lg' | 'xl' (default: 'md')

**Usage:**
```tsx
import { Modal } from '@/components/Modal'

<Modal
  isOpen={isOpen}
  onClose={() => setIsOpen(false)}
  title="Confirm Action"
  description="This action cannot be undone"
  footer={
    <>
      <Button variant="secondary" onClick={() => setIsOpen(false)}>
        Cancel
      </Button>
      <Button variant="danger" onClick={handleConfirm}>
        Delete
      </Button>
    </>
  }
>
  Are you sure you want to delete this customer?
</Modal>
```

**Features:**
- Backdrop click to close
- Escape key support
- Scroll prevention on body
- Multiple size options

---

### Form Components

#### FormField
Wrapper component for form field with label, error display, and hint text.

```tsx
import { FormField, Input } from '@/components/Form'

<FormField 
  label="Email" 
  error={errors.email}
  required 
  hint="We'll never share your email"
>
  <Input 
    type="email" 
    placeholder="you@example.com"
    error={!!errors.email}
  />
</FormField>
```

#### Input
Standard text input field.

```tsx
<Input
  type="email"
  placeholder="Enter email"
  error={!!errors.email}
  value={email}
  onChange={(e) => setEmail(e.target.value)}
/>
```

#### Select
Dropdown select component.

```tsx
<Select
  value={status}
  onChange={(e) => setStatus(e.target.value)}
  options={[
    { value: 'ACTIVE', label: 'Active' },
    { value: 'INACTIVE', label: 'Inactive' },
    { value: 'SUSPENDED', label: 'Suspended' },
  ]}
/>
```

#### TextArea
Multi-line text input.

```tsx
<TextArea
  placeholder="Enter notes..."
  value={notes}
  onChange={(e) => setNotes(e.target.value)}
  rows={4}
/>
```

#### Checkbox
Checkbox input with optional label.

```tsx
<Checkbox
  label="I agree to the terms"
  checked={agreed}
  onChange={(e) => setAgreed(e.target.checked)}
/>
```

---

### Alert

Display important messages to users.

**Props:**
- `type`: 'success' | 'error' | 'warning' | 'info' (default: 'info')
- `title`: string (optional)
- `children`: ReactNode - Alert message content
- `onDismiss`: () => void (optional) - Called when close button clicked
- `icon`: ReactNode (optional) - Custom icon

**Usage:**
```tsx
import { Alert } from '@/components/Alert'

<Alert type="error" title="Error" onDismiss={handleDismiss}>
  Something went wrong. Please try again.
</Alert>

<Alert type="success">
  Customer created successfully!
</Alert>

<Alert type="warning">
  This action will suspend the account.
</Alert>
```

---

### Toast

Non-intrusive notifications that auto-dismiss.

**Usage:**
```tsx
import { useToast, ToastContainer } from '@/components/Toast'

function MyComponent() {
  const toast = useToast()

  return (
    <>
      <ToastContainer />
      <button onClick={() => toast.success('Saved!')}>
        Save
      </button>
    </>
  )
}
```

**Methods:**
- `toast.success(message, duration?)`
- `toast.error(message, duration?)`
- `toast.warning(message, duration?)`
- `toast.info(message, duration?)`
- `toast.show(message, type, duration?)`
- `toast.dismiss(id)`

Duration defaults to 4000ms (4 seconds), set to 0 for permanent.

---

### Table & Pagination

#### Table Component
Responsive data table with sorting and loading states.

```tsx
import { Table } from '@/components/Table'

<Table
  columns={[
    { key: 'name', label: 'Name', sortable: true },
    { key: 'email', label: 'Email', width: '200px' },
    { key: 'status', label: 'Status' },
  ]}
  rows={customers}
  sortKey={sortKey}
  sortDir={sortDir}
  onSort={(key) => setSortKey(key)}
  loading={isLoading}
/>
```

#### Pagination Component
Paginated navigation control.

```tsx
import { Pagination } from '@/components/Table'

<Pagination
  currentPage={page}
  totalPages={totalPages}
  onPageChange={setPage}
  isLoading={loading}
/>
```

---

### Skeleton

Loading placeholder components.

**Usage:**
```tsx
import { Skeleton, TableSkeleton, CardSkeleton } from '@/components/Skeleton'

// Generic skeleton
<div className="space-y-4">
  <Skeleton count={3} />
</div>

// Table skeleton
<TableSkeleton />

// Card skeleton
<CardSkeleton />

// List of cards
<ListSkeleton count={5} />
```

---

## Layout Components

### PageHeader
Displays page title with optional subtitle and actions.

**Usage:**
```tsx
import PageHeader from '@/components/PageHeader'

<PageHeader 
  title="Customers"
  subtitle="Manage subscriber accounts"
/>
```

### Sidebar
Main navigation sidebar (in app layout).

### Topbar
Header bar with user menu and context.

---

## Utility Hooks

### useApi
Fetch data from API endpoints.

```tsx
import { useApi } from '@/lib/useApi'

const { data, loading, error, refetch } = useApi<Customer[]>(
  '/api/v1/customers'
)

if (loading) return <div>Loading...</div>
if (error) return <div>Error: {error.message}</div>
if (!data) return <div>No data</div>

return (
  <ul>
    {data.map(customer => (
      <li key={customer.id}>{customer.full_name}</li>
    ))}
  </ul>
)
```

### useApiMutation
Execute POST/PUT/PATCH/DELETE operations.

```tsx
import { useApiMutation } from '@/lib/useApi'

const { execute, loading, error } = useApiMutation()

async function handleSubmit(formData) {
  const result = await execute('POST', '/api/v1/customers', formData)
  if (result) {
    console.log('Created:', result)
  }
}
```

### useToast
Manage toast notifications.

```tsx
import { useToast } from '@/components/Toast'

const toast = useToast()

// Show notifications
toast.success('Operation successful!')
toast.error('Operation failed!')
toast.warning('Are you sure?')
toast.info('FYI: Something happened')
```

### useAuth
Get current user and authentication methods.

```tsx
import { useAuth } from '@/lib/auth-context'

const { user, login, logout, loading } = useAuth()

if (!user) return <Redirect to="/login" />

return (
  <div>
    <p>Welcome, {user.username}!</p>
    <button onClick={logout}>Logout</button>
  </div>
)
```

---

## Design System

### CSS Variables
All components use CSS variables defined in `globals.css`:

```css
--accent: RGB accent color
--accent-bright: Brighter accent
--accent-dark: Darker accent
--bg: Base background
--bg-elev: Elevated background (cards)
--bg-elev-hover: Elevated hover state
--text: Primary text color
--text-dim: Secondary text color
--text-mute: Muted text color
--border: Border color
--border-strong: Stronger border color
--danger: Danger/error color
```

### Spacing
Uses Tailwind's default spacing scale:
- `px-3` = 12px
- `py-2` = 8px
- `gap-3` = 12px spacing
- `mb-6` = 24px margin bottom

### Responsive Design
All components support Tailwind's responsive prefixes:
- `md:` - Medium screens and up (768px)
- `lg:` - Large screens and up (1024px)
- `xl:` - Extra large screens and up (1280px)

---

## Common Patterns

### Loading State
```tsx
import { Button } from '@/components/Button'
import { Skeleton } from '@/components/Skeleton'

if (loading) return <Skeleton count={5} />

return (
  <Button loading={isSubmitting}>
    Save Changes
  </Button>
)
```

### Error Handling
```tsx
import { Alert } from '@/components/Alert'
import { useToast } from '@/components/Toast'

const toast = useToast()

if (error) {
  return <Alert type="error" title="Error">{error.message}</Alert>
}

// Or with toast
toast.error('Operation failed!')
```

### Form Submission
```tsx
import { FormField, Input } from '@/components/Form'
import { Button } from '@/components/Button'

<form onSubmit={handleSubmit} className="space-y-6">
  <FormField label="Name" error={errors.name} required>
    <Input 
      value={form.name}
      onChange={(e) => setForm({...form, name: e.target.value})}
      error={!!errors.name}
    />
  </FormField>

  <FormField label="Email" error={errors.email} required>
    <Input 
      type="email"
      value={form.email}
      onChange={(e) => setForm({...form, email: e.target.value})}
      error={!!errors.email}
    />
  </FormField>

  <Button type="submit" fullWidth loading={loading}>
    Save
  </Button>
</form>
```

### Confirmation Modal
```tsx
import { Modal } from '@/components/Modal'
import { Button } from '@/components/Button'

<Modal
  isOpen={showConfirm}
  onClose={() => setShowConfirm(false)}
  title="Confirm Delete"
  footer={
    <>
      <Button variant="secondary" onClick={() => setShowConfirm(false)}>
        Cancel
      </Button>
      <Button variant="danger" onClick={handleDelete}>
        Delete
      </Button>
    </>
  }
>
  This action cannot be undone.
</Modal>
```

---

## Accessibility

All components follow WCAG 2.1 AA guidelines:
- Proper semantic HTML
- ARIA labels where needed
- Keyboard navigation support
- Sufficient color contrast
- Focus indicators

---

## Performance Tips

1. **Code Splitting**: Use Next.js dynamic imports for large components
2. **Memoization**: Use `React.memo` for frequently re-rendered components
3. **List Virtualization**: For large tables, consider virtualization
4. **Image Optimization**: Use Next.js `Image` component
5. **Bundle Analysis**: Use `@next/bundle-analyzer` to monitor size

---

## Testing Components

### Unit Testing Example
```tsx
import { render, screen } from '@testing-library/react'
import { Button } from '@/components/Button'

describe('Button', () => {
  it('renders with text', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByText('Click me')).toBeInTheDocument()
  })

  it('disables when loading', () => {
    render(<Button loading>Loading</Button>)
    expect(screen.getByRole('button')).toBeDisabled()
  })
})
```

---

## Component Checklist

When creating a new component, ensure:
- [ ] TypeScript types defined
- [ ] Props documented
- [ ] Usage examples provided
- [ ] Accessibility features added
- [ ] Mobile responsive
- [ ] Loading states handled
- [ ] Error states handled
- [ ] Theme-aware (uses CSS variables)
- [ ] No console errors in strict mode
- [ ] Unit tests written

---

## Contributing

When adding new components:
1. Create component file in `src/components/`
2. Export from `src/components/index.ts`
3. Add documentation here
4. Add example usage in Storybook (future)
5. Add unit tests
6. Ensure no breaking changes to existing components

