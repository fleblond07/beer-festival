<template>
  <div
    :id="`brewery-${brewery.id}`"
    class="festival-card bg-dark-card rounded-lg shadow-lg overflow-hidden hover:shadow-2xl hover:shadow-accent-cyan/20 transition-all duration-300 cursor-pointer border border-dark-lighter"
    data-testid="brewery-card"
  >
    <div v-if="brewery.logo" class="h-48 overflow-hidden">
      <img
        :src="brewery.logo"
        :alt="brewery.name"
        class="w-full h-full object-cover"
        data-testid="brewery-logo"
      />
    </div>
    <div class="p-6">
      <div class="flex justify-between items-start mb-2">
        <h3 class="text-xl font-bold text-white" data-testid="brewery-name">
          {{ brewery.name }}
        </h3>
      </div>

      <div class="flex items-center text-gray-300 mb-2">
        <svg
          class="w-4 h-4 mr-2"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          data-testid="location-icon"
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
        <span class="text-sm" data-testid="brewery-location">{{ brewery.city }}</span>
      </div>

      <p class="text-gray-300 text-sm mb-4 line-clamp-2" data-testid="brewery-description">
        {{ brewery.description }}
      </p>

      <div class="flex items-center justify-between">
        <div v-if="brewery.festivalCount" class="flex items-center text-gray-200">
          <svg
            class="w-5 h-5 mr-1"
            fill="currentColor"
            viewBox="0 0 20 20"
            data-testid="festival-icon"
          >
            <path
              d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10v6a1 1 0 001 1h12a1 1 0 001-1v-6a1 1 0 00-1-1H4a1 1 0 00-1 1z"
            />
          </svg>
          <span class="text-sm font-semibold" data-testid="festival-count">
            {{ brewery.festivalCount }} festivals
          </span>
        </div>
        <a
          v-if="brewery.website"
          :href="brewery.website"
          target="_blank"
          rel="noopener noreferrer"
          class="text-gray-200 hover:text-accent-pink text-sm font-semibold transition-colors"
          data-testid="brewery-website"
          @click.stop
        >
          Site web →
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// import { useRouter } from 'vue-router'
import type { Brewery } from '@/types'

interface Props {
  brewery: Brewery
  showStatus?: boolean
}

withDefaults(defineProps<Props>(), {
  showStatus: true,
})

// const router = useRouter()

// const navigateToDetail = () => {
//   router.push(`/brewery/${props.brewery.id}`)
// }
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
