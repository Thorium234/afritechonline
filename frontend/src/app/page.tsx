'use client'

export default function HomePage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-obsidian relative overflow-hidden">
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,rgba(210,255,42,0.03),transparent_70%)]" />
      <div className="max-w-md w-full card p-8 text-center relative animate-slide-up">
        <div className="mb-6">
          <div className="w-16 h-16 mx-auto mb-4 rounded-2xl bg-acid-lime/10 border border-acid-lime/20 flex items-center justify-center">
            <svg className="w-8 h-8 text-acid-lime" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h1 className="text-4xl font-bold text-pearl mb-2">Afritech Online</h1>
          <p className="text-slate text-sm">ISP Management Platform</p>
        </div>
        <a href="/login" className="btn-primary inline-block px-8 py-3 w-full">
          Login
        </a>
      </div>
    </div>
  )
}
