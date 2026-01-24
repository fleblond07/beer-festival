<template>
  <div class="min-h-screen bg-dark-bg text-white">
    <BreweryList :breweries="breweries" :loading="loading" />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useBreweries } from '@/composables/useBreweries'
import { BreweryList } from '@/components/BreweryList'

const { breweries, loading, error, fetchBreweries } = useBreweries()

onMounted(async () => {
  await fetchBreweries()
  if (error.value) {
    console.log(error.value)
  }
})
</script>
