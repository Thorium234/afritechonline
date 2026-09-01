import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export async function GET(request: Request) {
  const token = request.headers.get('authorization')
  const res = await fetch(`${API_URL}/api/v1/packages`, {
    headers: { Authorization: token || '' },
  })
  const data = await res.json()
  return NextResponse.json(data, { status: res.status })
}
