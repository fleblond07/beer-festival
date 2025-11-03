import { ref } from 'vue'

const API_URL = import.meta.env.VITE_API_URL

export interface User {
  id: string
  email: string
}

export interface LoginResponse {
  accessToken: string
  refreshToken: string
  user: User
}

const accessToken = ref<string | null>(localStorage.getItem('accessToken'))
const currentUser = ref<User | null>(null)

const storedUser = localStorage.getItem('user')
if (storedUser) {
  try {
    currentUser.value = JSON.parse(storedUser)
  } catch {
    localStorage.removeItem('user')
  }
}

export const useAuth = () => {
  const login = async (email: string, password: string): Promise<void> => {
    const response = await fetch(`${API_URL}/api/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    })

    if (!response.ok) {
      throw new Error('Invalid credentials')
    }

    const data: LoginResponse = await response.json()
    accessToken.value = data.accessToken
    currentUser.value = data.user

    localStorage.setItem('accessToken', data.accessToken)
    localStorage.setItem('refreshToken', data.refreshToken)
    localStorage.setItem('user', JSON.stringify(data.user))
  }

  const logout = () => {
    accessToken.value = null
    currentUser.value = null
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
  }

  const verify = async (): Promise<boolean> => {
    if (!accessToken.value) {
      return false
    }

    try {
      const response = await fetch(`${API_URL}/api/auth/verify`, {
        headers: {
          Authorization: `Bearer ${accessToken.value}`,
        },
      })

      const data = await response.json()
      if (data.valid && data.user) {
        currentUser.value = data.user
        return true
      }

      logout()
      return false
    } catch {
      logout()
      return false
    }
  }

  const isAuthenticated = () => {
    return !!accessToken.value && !!currentUser.value
  }

  return {
    login,
    logout,
    verify,
    isAuthenticated,
    user: currentUser,
    token: accessToken,
  }
}
