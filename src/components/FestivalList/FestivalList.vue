<template>
  <div class="festival-list-container w-full py-16 px-4" data-testid="festival-list-container">
    <div class="max-w-6xl mx-auto">
      <h2 class="text-center text-4xl font-bold text-white mb-2" data-testid="section-title">
        Tous les Festivals
      </h2>
      <p class="text-center text-gray-400 mb-12" data-testid="section-subtitle">
        Découvrez tous les festivals de bière en France
      </p>

      <div v-if="festivals.length > 0" class="space-y-8">
        <div
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
          data-testid="festivals-grid"
        >
          <FestivalCard
            v-for="festival in festivals"
            :key="festival.id"
            :festival="festival"
            :show-status="true"
            @click="handleFestivalClick"
            :data-testid="`festival-card-${festival.id}`"
          />
        </div>
      </div>

      <div v-else class="text-center py-16" data-testid="empty-state">
        <svg
          class="w-16 h-16 mx-auto mb-4 text-gray-500"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
          />
        </svg>
        <h3 class="text-2xl font-bold text-white mb-2">Aucun festival</h3>
        <p class="text-gray-400">Il n'y a actuellement aucun festival à afficher.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Festival } from '@/types'
import { FestivalCard } from '@/components/FestivalCard'

interface Props {
  festivals: Festival[]
}

defineProps<Props>()

const emit = defineEmits<{
  (e: 'festival-click', festival: Festival): void
}>()

const handleFestivalClick = (festival: Festival) => {
  emit('festival-click', festival)
}
</script>
