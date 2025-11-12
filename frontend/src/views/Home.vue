<template>
  <div>
    <NextFestival :festival="nextFestival" :loading="loading" data-testid="next-festival-section" />

    <FestivalMap
      :festivals="sortedFestivals"
      data-testid="festival-map-section"
      @popup-click="handlePopupClick"
    />

    <FestivalList
      :festivals="sortedFestivals"
      :loading="loading"
      data-testid="festival-list-section"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NextFestival } from '@/components/NextFestival'
import { FestivalMap } from '@/components/FestivalMap'
import { FestivalList } from '@/components/FestivalList'
import { useFestivals } from '@/composables/useFestivals'
import type { Festival } from '@/types'

const { nextFestival, sortedFestivals, loading, fetchFestivals } = useFestivals()
const router = useRouter()

onMounted(() => {
  fetchFestivals()
})

const handlePopupClick = (festival: Festival) => {
  router.push(`/festival/${festival.id}`)
}
</script>

<style scoped>
:deep(.highlight-festival) {
  animation: highlight 2s ease-in-out;
}

@keyframes highlight {
  0% {
    box-shadow: 0 0 0 0 rgba(0, 100, 110, 0.7);
    transform: scale(1);
  }
  50% {
    box-shadow: 0 0 20px 10px rgba(0, 100, 110, 0.5);
    transform: scale(1.02);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(0, 100, 110, 0);
    transform: scale(1);
  }
}
</style>
