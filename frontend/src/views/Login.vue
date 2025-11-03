<template>
  <div class="min-h-[calc(100vh-300px)] flex items-center justify-center p-4">
    <div class="bg-dark-card rounded-lg shadow-xl p-8 max-w-md w-full border border-dark-lighter">
      <h1 class="text-3xl font-bold text-white mb-6 text-center">Admin Login</h1>

      <form class="space-y-4" @submit.prevent="handleLogin">
        <div>
          <label for="email" class="block text-sm font-medium text-gray-300 mb-2"> Email </label>
          <input
            id="email"
            v-model="email"
            type="email"
            required
            class="w-full px-4 py-2 bg-dark-bg border border-dark-lighter rounded-md text-white focus:outline-none focus:ring-2 focus:ring-accent-cyan"
            placeholder="admin@example.com"
          />
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-300 mb-2">
            Password
          </label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            class="w-full px-4 py-2 bg-dark-bg border border-dark-lighter rounded-md text-white focus:outline-none focus:ring-2 focus:ring-accent-cyan"
            placeholder="••••••••"
          />
        </div>

        <div v-if="error" class="text-red-400 text-sm text-center" data-testid="login-error">
          {{ error }}
        </div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-gradient-to-r from-accent-cyan to-accent-purple text-white font-semibold py-2 px-4 rounded-md hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {{ loading ? 'Logging in...' : 'Login' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/services/auth'

const router = useRouter()
const { login } = useAuth()

const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const handleLogin = async () => {
  error.value = ''
  loading.value = true

  try {
    await login(email.value, password.value)
    router.push('/admin')
  } catch {
    error.value = 'Invalid email or password'
  } finally {
    loading.value = false
  }
}
</script>
