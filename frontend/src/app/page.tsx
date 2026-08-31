'use client'

export default function HomePage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-obsidian">
      <div className="max-w-md w-full card p-8 text-center">
        <h1 className="text-4xl font-bold text-acid-lime mb-4">Afritech Online</h1>
        <p className="text-slate mb-8">ISP Management Platform</p>
        <a href="/login" className="btn-primary inline-block px-8 py-3">
          Login
        </a>
      </div>
    </div>
  )
}
