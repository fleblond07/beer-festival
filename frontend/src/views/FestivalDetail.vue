<template>
  <div class="min-h-screen bg-dark-bg text-white">
    <div v-if="loading" class="flex justify-center items-center min-h-screen">
      <div class="text-xl text-gray-300">Chargement...</div>
    </div>

    <div v-else-if="!festival" class="flex justify-center items-center min-h-screen">
      <div class="text-center">
        <h1 class="text-2xl font-bold mb-4">Festival non trouvé</h1>
        <router-link to="/" class="text-accent-cyan hover:text-accent-pink transition-colors">
          Retour à l'accueil
        </router-link>
      </div>
    </div>

    <div v-else class="container mx-auto px-4 py-8 max-w-6xl">
      <router-link
        to="/"
        class="inline-flex items-center text-gray-300 hover:text-accent-cyan transition-colors mb-6"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M15 19l-7-7 7-7"
          />
        </svg>
        Retour
      </router-link>

      <div class="bg-dark-card rounded-lg shadow-2xl overflow-hidden border border-dark-lighter">
        <div v-if="festival.image" class="h-96 overflow-hidden">
          <img
            :src="festival.image"
            :alt="festival.name"
            class="w-full h-full object-cover"
            data-testid="festival-detail-image"
          />
        </div>

        <div class="p-8">
          <div class="flex justify-between items-start mb-6">
            <div>
              <h1 class="text-4xl font-bold mb-2" data-testid="festival-detail-name">
                {{ festival.name }}
              </h1>
              <div class="flex items-center text-gray-300 mb-2">
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                  />
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
                <span class="text-lg" data-testid="festival-detail-location">
                  {{ festival.city }}, {{ festival.region }}
                </span>
              </div>
              <div class="flex items-center text-gray-300">
                <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
                <span class="text-lg" data-testid="festival-detail-dates">
                  {{ formattedDateRange }}
                </span>
              </div>
            </div>
            <span
              v-if="statusText"
              :class="statusClass"
              class="px-4 py-2 rounded-full text-sm font-semibold"
              data-testid="festival-detail-status"
            >
              {{ statusText }}
            </span>
          </div>

          <p
            class="text-gray-200 text-lg leading-relaxed mb-6"
            data-testid="festival-detail-description"
          >
            {{ festival.description }}
          </p>

          <div v-if="festival.website" class="mb-8">
            <a
              :href="festival.website"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center bg-accent-cyan text-dark-bg px-6 py-3 rounded-lg font-semibold hover:bg-accent-pink transition-colors"
              data-testid="festival-detail-website"
            >
              Visiter le site web
              <svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                />
              </svg>
            </a>
          </div>

          <div class="border-t border-dark-lighter pt-8">
            <h2 class="text-3xl font-bold mb-6">
              Brasseries
              <span v-if="festival.breweryCount" class="text-accent-cyan">
                ({{ festival.breweryCount }})
              </span>
            </h2>

            <div v-if="breweriesLoading" class="text-center py-12">
              <div class="text-lg text-gray-300">Chargement des brasseries...</div>
            </div>

            <div v-else-if="breweriesError" class="text-center py-12">
              <div class="text-lg text-red-400">{{ breweriesError }}</div>
            </div>

            <div v-else-if="breweries.length === 0" class="text-center py-12 text-gray-400 text-lg">
              Aucune brasserie pour ce festival
            </div>

            <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              <div
                v-for="brewery in breweries"
                :key="brewery.id"
                class="bg-dark-lighter rounded-lg p-6 hover:shadow-lg hover:shadow-accent-cyan/20 transition-all duration-300 border border-dark-lighter hover:border-accent-cyan"
                data-testid="brewery-card"
              >
                <div v-if="brewery.logo" class="mb-4 h-24 flex items-center justify-center">
                  <img
                    :src="brewery.logo"
                    :alt="brewery.name"
                    class="max-h-full max-w-full object-contain"
                    data-testid="brewery-logo"
                  />
                </div>
                <h3 class="text-xl font-bold mb-2 text-white" data-testid="brewery-name">
                  {{ brewery.name }}
                </h3>
                <p
                  v-if="brewery.description"
                  class="text-gray-300 text-sm mb-3 line-clamp-3"
                  data-testid="brewery-description"
                >
                  {{ brewery.description }}
                </p>
                <div v-if="brewery.city" class="text-gray-400 text-sm mb-2">
                  <svg
                    class="w-4 h-4 inline mr-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                    />
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                  </svg>
                  {{ brewery.city }}
                </div>
                <a
                  v-if="brewery.website"
                  :href="brewery.website"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-accent-cyan hover:text-accent-pink text-sm font-semibold transition-colors inline-flex items-center"
                  data-testid="brewery-website"
                >
                  Site web →
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useFestivals } from '@/composables/useFestivals'
import { useBreweries } from '@/composables/useBreweries'
import { formatDateRange, getDaysUntilText, isUpcoming } from '@/services/dateUtils'

const route = useRoute()
const festivalId = route.params.id as string

const { sortedFestivals, loading, fetchFestivals } = useFestivals()
const {
  breweries,
  loading: breweriesLoading,
  error: breweriesError,
  fetchBreweriesByFestival,
} = useBreweries()

const festival = computed(() => {
  return sortedFestivals.value.find(f => f.id === Number(festivalId))
})

const formattedDateRange = computed(() => {
  if (!festival.value) return ''
  return formatDateRange(festival.value.startDate, festival.value.endDate)
})

const statusText = computed(() => {
  if (!festival.value) return ''
  return getDaysUntilText(festival.value.startDate)
})

const statusClass = computed(() => {
  if (!festival.value) return ''
  if (isUpcoming(festival.value.startDate)) {
    return 'bg-accent-cyan/20 text-gray-200'
  }
  return 'bg-gray-700 text-gray-200'
})

onMounted(async () => {
  await fetchFestivals()
  if (festivalId) {
    await fetchBreweriesByFestival(festivalId)
  }
})
</script>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
