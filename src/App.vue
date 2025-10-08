<template>
  <div id="app" class="min-h-screen" data-testid="app">
    <header class="bg-gradient-to-r from-accent-cyan to-accent-purple text-white shadow-2xl">
      <div class="max-w-7xl mx-auto px-4 py-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
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
          </div>
        </div>
      </div>
    </header>

    <main>
      <NextFestival :festival="nextFestival" data-testid="next-festival-section" />

      <FestivalMap
        :festivals="sortedFestivals"
        data-testid="festival-map-section"
        @marker-click="handleMarkerClick"
      />

      <FestivalList
        :festivals="sortedFestivals"
        data-testid="festival-list-section"
        @festival-click="handleFestivalClick"
      />
    </main>

    <footer class="bg-dark-card text-white py-8 mt-16 border-t border-dark-lighter">
      <div class="max-w-7xl mx-auto px-4 text-center">
        <p class="text-gray-300" data-testid="footer-text">
          © 2025 Festivals de Bière. Fait avec ❤️ pour les amateurs de bière artisanale.
        </p>
        <p class="text-gray-400 text-sm mt-2" data-testid="footer-note">
          Les données sont fictives et à titre d'exemple uniquement.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { NextFestival } from '@/components/NextFestival'
import { FestivalMap } from '@/components/FestivalMap'
import { FestivalList } from '@/components/FestivalList'
import { useFestivals } from '@/composables/useFestivals'
import type { Festival } from '@/types'

const { nextFestival, sortedFestivals, fetchFestivals } = useFestivals()

onMounted(() => {
  fetchFestivals()
})

const handleMarkerClick = (festival: Festival) => {
  console.log('Marker clicked:', festival.name)
}

const handleFestivalClick = (festival: Festival) => {
  console.log('Festival clicked:', festival.name)
}
</script>
