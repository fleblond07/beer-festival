<template>
  <div id="app" class="min-h-screen" data-testid="app">
    <header class="bg-gradient-to-r from-accent-cyan to-accent-purple text-white shadow-2xl">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between">
          <router-link
            to="/"
            class="flex items-center space-x-3 hover:opacity-80 transition-opacity"
          >
            <svg class="w-10 h-10" fill="currentColor" viewBox="0 0 20 20" data-testid="logo-icon">
              <path
                d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10v6a1 1 0 001 1h12a1 1 0 001-1v-6a1 1 0 00-1-1H4a1 1 0 00-1 1z"
              />
            </svg>
            <div>
              <h1 class="text-3xl font-bold" data-testid="site-title">Festivals de Bière</h1>
              <p class="text-white/80 text-sm" data-testid="site-tagline">
                Découvrez les meilleurs festivals de France
              </p>
            </div>
          </router-link>
        </div>
      </div>
    </header>

    <main>
      <router-view />
    </main>

    <footer class="bg-dark-card text-white py-8 mt-16 border-t border-dark-lighter">
      <div class="max-w-7xl mx-auto px-4">
        <div class="relative flex justify-center items-center">
          <p class="text-gray-300 text-center" data-testid="footer-text">
            © 2025 Festivals de Bière. Fait avec ❤️ pour les amateurs de bière artisanale.
          </p>
          <router-link
            to="/login"
            class="absolute right-0 text-gray-400 hover:text-accent-cyan transition-colors"
            data-testid="admin-link"
          >
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              data-testid="lock-icon"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
              />
            </svg>
          </router-link>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts"></script>
