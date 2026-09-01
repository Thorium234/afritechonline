'use client'

/**
 * Password validation following NIST SP 800-63B guidelines
 * https://pages.nist.gov/800-63-3/sp800-63b.html
 */

interface PasswordStrength {
  score: 0 | 1 | 2 | 3 | 4 // 0 = very weak, 4 = very strong
  feedback: string[]
  isValid: boolean
  percentage: number
}

const COMMON_PASSWORDS = new Set([
  'password', '123456', 'password123', '12345678', 'qwerty', 'abc123',
  '111111', '1234567', 'letmein', 'welcome', 'monkey', '1234567890',
  'dragon', 'master', 'sunshine', 'princess', 'football', 'soccer',
  'batman', 'superman', 'starwars', 'trustno1', 'mypassword', 'pass',
])

export function validatePassword(password: string): PasswordStrength {
  const feedback: string[] = []
  let score = 0
  const minLength = 14

  // Length check
  if (password.length < minLength) {
    feedback.push(`Password must be at least ${minLength} characters long`)
  } else if (password.length < 16) {
    score++
  } else if (password.length < 20) {
    score += 1.5
  } else {
    score += 2
  }

  // Character variety
  const hasLowercase = /[a-z]/.test(password)
  const hasUppercase = /[A-Z]/.test(password)
  const hasNumbers = /\d/.test(password)
  const hasSpecial = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password)

  if (!hasLowercase) {
    feedback.push('Add lowercase letters')
  } else {
    score += 0.5
  }

  if (!hasUppercase) {
    feedback.push('Add uppercase letters')
  } else {
    score += 0.5
  }

  if (!hasNumbers) {
    feedback.push('Add numbers')
  } else {
    score += 0.5
  }

  if (!hasSpecial) {
    feedback.push('Add special characters')
  } else {
    score += 0.5
  }

  // Common password check
  const passwordLower = password.toLowerCase()
  if (COMMON_PASSWORDS.has(passwordLower)) {
    feedback.push('This password is too common')
    score = Math.max(0, score - 2)
  }

  // Sequential/repeated characters
  if (/(.)\1{2,}/.test(password)) {
    feedback.push('Avoid repeating characters')
    score = Math.max(0, score - 0.5)
  }

  if (/(?:abc|bcd|cde|def|efg|fgh|ghi|hij|ijk|jkl|klm|lmn|mno|nop|opq|pqr|qrs|rst|stu|tuv|uvw|vwx|wxy|xyz|012|123|234|345|456|567|678|789|890|901)/i.test(password)) {
    feedback.push('Avoid sequential characters')
    score = Math.max(0, score - 0.5)
  }

  // Cap score at 4
  const finalScore = Math.min(4, Math.max(0, Math.round(score))) as 0 | 1 | 2 | 3 | 4
  const percentage = ((finalScore + 1) / 5) * 100

  return {
    score: finalScore,
    feedback,
    isValid: finalScore >= 2 && password.length >= minLength,
    percentage,
  }
}

export function getPasswordStrengthColor(score: 0 | 1 | 2 | 3 | 4): string {
  switch (score) {
    case 0:
      return 'bg-[var(--danger)]'
    case 1:
      return 'bg-[#f59e0b]' // orange
    case 2:
      return 'bg-[#eab308]' // yellow
    case 3:
      return 'bg-[#84cc16]' // lime
    case 4:
      return 'bg-[#22c55e]' // green
    default:
      return 'bg-[var(--border)]'
  }
}

export function getPasswordStrengthLabel(score: 0 | 1 | 2 | 3 | 4): string {
  switch (score) {
    case 0:
      return 'Very Weak'
    case 1:
      return 'Weak'
    case 2:
      return 'Fair'
    case 3:
      return 'Good'
    case 4:
      return 'Very Strong'
    default:
      return 'Unknown'
  }
}
