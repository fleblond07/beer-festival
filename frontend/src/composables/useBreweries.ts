import { ref } from 'vue'
import type { Brewery } from '@/types'

const API_URL = import.meta.env.VITE_API_URL

export function useBreweries() {
  const breweries = ref<Brewery[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchBreweriesByFestival = async (festivalId: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_URL}/api/festivals/${festivalId}/breweries`)

      if (!response.ok) {
        throw new Error('Failed to fetch breweries')
      }

      const data = await response.json()
      breweries.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
      breweries.value = []
    } finally {
      loading.value = false
    }
  }

  return {
    breweries,
    loading,
    error,
    fetchBreweriesByFestival,
  }
}
