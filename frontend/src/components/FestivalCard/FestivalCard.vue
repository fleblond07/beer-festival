<template>
  <div
    class="festival-card bg-dark-card rounded-lg shadow-lg overflow-hidden hover:shadow-2xl hover:shadow-accent-cyan/20 transition-all duration-300 cursor-pointer border border-dark-lighter"
    data-testid="festival-card"
    @click="$emit('click', festival)"
  >
    <div v-if="festival.image" class="h-48 overflow-hidden">
      <img
        :src="festival.image"
        :alt="festival.name"
        class="w-full h-full object-cover"
        data-testid="festival-image"
      />
    </div>
    <div class="p-6">
      <div class="flex justify-between items-start mb-2">
        <h3 class="text-xl font-bold text-white" data-testid="festival-name">
          {{ festival.name }}
        </h3>
        <span
          v-if="showStatus"
          :class="statusClass"
          class="px-3 py-1 rounded-full text-xs font-semibold"
          data-testid="festival-status"
        >
          {{ statusText }}
        </span>
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
        <span class="text-sm" data-testid="festival-location"
          >{{ festival.city }}</span
        >
      </div>

      <div class="flex items-center text-gray-300 mb-3">
        <svg
          class="w-4 h-4 mr-2"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          data-testid="calendar-icon"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
          />
        </svg>
        <span class="text-sm" data-testid="festival-dates">{{ formattedDateRange }}</span>
      </div>

      <p class="text-gray-300 text-sm mb-4 line-clamp-2" data-testid="festival-description">
        {{ festival.description }}
      </p>

      <div class="flex items-center justify-between">
        <div v-if="festival.breweryCount" class="flex items-center text-gray-200">
          <svg
            class="w-5 h-5 mr-1"
            fill="currentColor"
            viewBox="0 0 20 20"
            data-testid="brewery-icon"
          >
            <path
              d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10v6a1 1 0 001 1h12a1 1 0 001-1v-6a1 1 0 00-1-1H4a1 1 0 00-1 1z"
            />
          </svg>
          <span class="text-sm font-semibold" data-testid="brewery-count">
            {{ festival.breweryCount }} brasseries
          </span>
        </div>
        <a
          v-if="festival.website"
          :href="festival.website"
          target="_blank"
          rel="noopener noreferrer"
          class="text-gray-200 hover:text-accent-pink text-sm font-semibold transition-colors"
          data-testid="festival-website"
          @click.stop
        >
          Site web →
        </a>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Festival } from '@/types'
import { formatDateRange, getDaysUntilText, isUpcoming } from '@/services/dateUtils'

interface Props {
  festival: Festival
  showStatus?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showStatus: true,
})

defineEmits<{
  (_e: 'click', _festival: Festival): void
}>()

const formattedDateRange = computed(() =>
  formatDateRange(props.festival.startDate, props.festival.endDate)
)

const statusText = computed(() => getDaysUntilText(props.festival.startDate))

const statusClass = computed(() => {
  if (isUpcoming(props.festival.startDate)) {
    return 'bg-accent-cyan/20 text-gray-200'
  }
  return 'bg-gray-700 text-gray-200'
})
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
