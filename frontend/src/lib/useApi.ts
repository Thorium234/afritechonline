'use client'

import { useState, useCallback, useEffect } from 'react'
import { api } from './api'
import type { ApiError } from './types'

interface UseApiOptions {
  onSuccess?: (data: any) => void
  onError?: (error: ApiError) => void
}

export function useApi<T>(url: string, opts?: UseApiOptions) {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<ApiError | null>(null)

  const fetch = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await api.get<T>(url)
      setData(result)
      opts?.onSuccess?.(result)
    } catch (err) {
      const apiError = err instanceof ApiError ? err : new ApiError('Unknown error', 500)
      setError(apiError)
      opts?.onError?.(apiError)
    } finally {
      setLoading(false)
    }
  }, [url, opts])

  useEffect(() => {
    fetch()
  }, [fetch])

  return { data, loading, error, refetch: fetch }
}

export function useApiMutation<T>(opts?: UseApiOptions) {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<ApiError | null>(null)

  const execute = useCallback(async (
    method: 'POST' | 'PUT' | 'PATCH' | 'DELETE',
    url: string,
    body?: unknown
  ): Promise<T | null> => {
    try {
      setLoading(true)
      setError(null)
      const result = await api[method.toLowerCase() as keyof typeof api]<T>(url, { body })
      opts?.onSuccess?.(result)
      return result
    } catch (err) {
      const apiError = err instanceof ApiError ? err : new ApiError('Unknown error', 500)
      setError(apiError)
      opts?.onError?.(apiError)
      return null
    } finally {
      setLoading(false)
    }
  }, [opts])

  return { execute, loading, error }
}
