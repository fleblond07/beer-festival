<template>
  <div class="next-festival-container w-full py-16 px-4" data-testid="next-festival-container">
    <div v-if="festival" class="max-w-6xl mx-auto">
      <h2 class="text-center text-4xl font-bold text-white mb-2" data-testid="section-title">
        Prochain Festival
      </h2>
      <p class="text-center text-gray-400 mb-8" data-testid="section-subtitle">
        Ne manquez pas le prochain événement !
      </p>

      <div
        class="next-festival-hero bg-gradient-to-r from-accent-purple to-accent-pink rounded-2xl shadow-2xl overflow-hidden transform transition-all duration-300 hover:scale-[1.02] border border-accent-purple/30 cursor-pointer"
        data-testid="hero-card"
        @click="navigateToDetail"
      >
        <div class="md:flex">
          <div
            v-if="festival.image"
            class="md:w-1/2 h-64 md:h-auto overflow-hidden"
            data-testid="hero-image-container"
          >
            <img
              :src="festival.image"
              :alt="festival.name"
              class="w-full h-full object-cover"
              data-testid="hero-image"
            />
          </div>
          <div class="md:w-1/2 p-8 md:p-12 text-white">
            <div class="flex items-center mb-4">
              <span
                class="bg-dark-card text-accent-cyan px-4 py-2 rounded-full text-sm font-bold border border-accent-cyan/50"
                data-testid="countdown-badge"
              >
                {{ countdownText }}
              </span>
            </div>

            <h3 class="text-3xl md:text-4xl font-bold mb-4" data-testid="festival-name">
              {{ festival.name }}
            </h3>

            <div class="flex items-center mb-3" data-testid="festival-location">
              <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z"
                  clip-rule="evenodd"
                />
              </svg>
              <span class="text-lg">{{ festival.city }}</span>
            </div>

            <div class="flex items-center mb-4" data-testid="festival-dates">
              <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                  clip-rule="evenodd"
                />
              </svg>
              <span class="text-lg">{{ formattedDateRange }}</span>
            </div>

            <p
              class="text-white/90 text-base mb-6 leading-relaxed"
              data-testid="festival-description"
            >
              {{ festival.description }}
            </p>

            <div class="flex items-center justify-between">
              <div
                v-if="festival.breweryCount"
                class="flex items-center"
                data-testid="brewery-info"
              >
                <svg class="w-6 h-6 mr-2" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10v6a1 1 0 001 1h12a1 1 0 001-1v-6a1 1 0 00-1-1H4a1 1 0 00-1 1z"
                  />
                </svg>
                <span class="text-lg font-semibold">{{ festival.breweryCount }} brasseries</span>
              </div>
              <a
                v-if="festival.website"
                :href="festival.website"
                target="_blank"
                rel="noopener noreferrer"
                class="bg-dark-card text-white px-6 py-3 rounded-lg font-bold hover:bg-accent-cyan transition-colors border border-white/30"
                data-testid="website-link"
                @click.stop
              >
                En savoir plus →
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="loading" class="max-w-6xl mx-auto text-center" data-testid="loading-state">
      <div class="bg-dark-card rounded-2xl p-12 border border-dark-lighter">
        <div class="flex justify-center mb-4">
          <svg class="w-16 h-16 text-accent-cyan animate-spin" fill="none" viewBox="0 0 24 24">
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
        </div>
        <h3 class="text-2xl font-bold text-white mb-2">Chargement...</h3>
        <p class="text-gray-400">Récupération des festivals en cours</p>
      </div>
    </div>

    <div v-else class="max-w-6xl mx-auto text-center" data-testid="no-festival">
      <div class="bg-dark-card rounded-2xl p-12 border border-dark-lighter">
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
            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
          />
        </svg>
        <h3 class="text-2xl font-bold text-white mb-2">Aucun festival à venir</h3>
        <p class="text-gray-400">
          Revenez bientôt pour découvrir les prochains festivals de bière !
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import type { Festival } from '@/types'
import { formatDateRange, getDaysUntilText } from '@/services/dateUtils'

interface Props {
  festival: Festival | null
  loading: boolean
}

const props = defineProps<Props>()
const router = useRouter()

const formattedDateRange = computed(() => {
  if (!props.festival) return ''
  return formatDateRange(props.festival.startDate, props.festival.endDate)
})

const countdownText = computed(() => {
  if (!props.festival) return ''
  return getDaysUntilText(props.festival.startDate)
})

const navigateToDetail = () => {
  if (props.festival) {
    router.push(`/festival/${props.festival.id}`)
  }
}
</script>
