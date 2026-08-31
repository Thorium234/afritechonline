'use client'

export default function HomePage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full bg-white p-8 rounded shadow text-center">
        <h1 className="text-3xl font-bold text-blue-600 mb-4">Afritech Online</h1>
        <p className="text-gray-600 mb-8">ISP Management Platform</p>
        <a href="/login" className="inline-block bg-blue-600 text-white px-6 py-3 rounded hover:bg-blue-700">
          Login
        </a>
      </div>
    </div>
  )
}
