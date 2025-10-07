import { ref, computed } from 'vue'
import type { Festival } from '@/types'
import { mockFestivals } from '@/mocks/festivals'
import dayjs from 'dayjs'

export function useFestivals() {
  const festivals = ref<Festival[]>(mockFestivals)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const nextFestival = computed(() => {
    const now = dayjs()
    const upcoming = festivals.value
      .filter(festival => dayjs(festival.startDate).isAfter(now))
      .sort((a, b) => dayjs(a.startDate).diff(dayjs(b.startDate)))

    return upcoming[0] || null
  })

  const sortedFestivals = computed(() => {
    const now = dayjs()
    return [...festivals.value].sort((a, b) => {
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
    nextFestival,
    sortedFestivals,
    getFestivalsByRegion,
    getFestivalById,
  }
}
