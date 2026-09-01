import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL || process.env.BACKEND_URL || 'http://backend:8080'

async function forward(request: Request, path: string[]) {
  const url = new URL(path.join('/'), API_URL)
  const incoming = new URL(request.url)
  incoming.searchParams.forEach((v, k) => url.searchParams.set(k, v))

  const headers: Record<string, string> = {}
  const auth = request.headers.get('authorization')
  if (auth) headers['Authorization'] = auth
  const ct = request.headers.get('content-type')
  if (ct) headers['Content-Type'] = ct

  const init: RequestInit = {
    method: request.method,
    headers,
    cache: 'no-store',
  }
  if (request.method !== 'GET' && request.method !== 'HEAD') {
    init.body = await request.text()
  }
  const res = await fetch(url.toString(), init)
  const text = await res.text()
  return new NextResponse(text, {
    status: res.status,
    headers: { 'content-type': res.headers.get('content-type') || 'application/json' },
  })
}

export const GET = (req: Request, ctx: { params: { path: string[] } }) =>
  forward(req, ['/api/v1', ...ctx.params.path])

export const POST = (req: Request, ctx: { params: { path: string[] } }) =>
  forward(req, ['/api/v1', ...ctx.params.path])

export const PUT = (req: Request, ctx: { params: { path: string[] } }) =>
  forward(req, ['/api/v1', ...ctx.params.path])

export const PATCH = (req: Request, ctx: { params: { path: string[] } }) =>
  forward(req, ['/api/v1', ...ctx.params.path])

export const DELETE = (req: Request, ctx: { params: { path: string[] } }) =>
  forward(req, ['/api/v1', ...ctx.params.path])
