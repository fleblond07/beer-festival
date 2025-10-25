import { ref, computed } from 'vue'
import type { Festival } from '@/types'
import dayjs from 'dayjs'

const API_FESTIVAL_URL_PATH = `${import.meta.env.VITE_API_URL}/api/festivals`

export function useFestivals() {
  const festivals = ref<Festival[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const fetchFestivals = async () => {
    loading.value = true
    error.value = null
    try {
      const response = await fetch(API_FESTIVAL_URL_PATH)
      if (!response.ok) {
        throw new Error('Failed to fetch festivals')
      }
      const data = await response.json()
      festivals.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'An error occurred'
    } finally {
      loading.value = false
    }
  }

  const nextFestival = computed(() => {
    const now = dayjs()
    const upcoming = festivals.value
      .filter(festival => dayjs(festival.startDate).isAfter(now))
      .sort((a, b) => dayjs(a.startDate).diff(dayjs(b.startDate)))

    return upcoming[0] || null
  })

  const sortedFestivals = computed(() => {
    const now = dayjs()
    return [...festivals.value]
      .filter(festival => dayjs(festival.endDate).isAfter(now))
      .sort((a, b) => {
        const aDate = dayjs(a.startDate)
        const bDate = dayjs(b.startDate)
        const aIsUpcoming = aDate.isAfter(now)
        const bIsUpcoming = bDate.isAfter(now)

        if (aIsUpcoming && !bIsUpcoming) return -1
        if (!aIsUpcoming && bIsUpcoming) return 1

        return aDate.diff(bDate)
      })
  })

  const getFestivalsByRegion = (region: string): Festival[] => {
    return festivals.value.filter(festival => festival.region === region)
  }

  const getFestivalById = (id: string): Festival | undefined => {
    return festivals.value.find(festival => festival.id === id)
  }

  return {
    festivals,
    loading,
    error,
    fetchFestivals,
    nextFestival,
    sortedFestivals,
    getFestivalsByRegion,
    getFestivalById,
  }
}
