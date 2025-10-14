import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useFestivals } from '@/composables/useFestivals'
import { mockFestivals } from '@/mocks/festivals'
import dayjs from 'dayjs'

describe('useFestivals', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Mock fetch API
    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockFestivals),
      } as Response)
    )
  })

  describe('initialization', () => {
    it('should initialize with empty festivals array', () => {
      const { festivals } = useFestivals()
      expect(festivals.value).toBeDefined()
      expect(festivals.value).toEqual([])
    })

    it('should initialize with loading set to false', () => {
      const { loading } = useFestivals()
      expect(loading.value).toBe(false)
    })

    it('should initialize with no error', () => {
      const { error } = useFestivals()
      expect(error.value).toBeNull()
    })

    it('should fetch and load festivals on mount', async () => {
      const { festivals, loading, fetchFestivals } = useFestivals()
      await fetchFestivals()

      expect(global.fetch).toHaveBeenCalled()
      expect(festivals.value.length).toBeGreaterThan(0)
      expect(loading.value).toBe(false)
    })

    it('should have images for all festivals after fetch', async () => {
      const { festivals, fetchFestivals } = useFestivals()
      await fetchFestivals()

      festivals.value.forEach(festival => {
        expect(festival.image).toBeDefined()
        expect(festival.image).toBeTruthy()
        expect(festival.image).toMatch(/^https?:\/\//)
      })
    })

    it('should have unique images for each festival after fetch', async () => {
      const { festivals, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const images = festivals.value.map(f => f.image)
      const uniqueImages = new Set(images)
      expect(uniqueImages.size).toBe(images.length)
    })
  })

  describe('nextFestival', () => {
    it('should return the next upcoming festival', async () => {
      const { nextFestival, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const now = dayjs()

      if (nextFestival.value) {
        const festivalDate = dayjs(nextFestival.value.startDate)
        expect(festivalDate.isAfter(now)).toBe(true)
      }
    })

    it('should return null if there are no upcoming festivals', async () => {
      const { festivals, nextFestival, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const pastDate = dayjs().subtract(1, 'year').format('YYYY-MM-DD')
      festivals.value = festivals.value.map(festival => ({
        ...festival,
        startDate: pastDate,
        endDate: pastDate,
      }))

      expect(nextFestival.value).toBeNull()
    })

    it('should return the closest upcoming festival when multiple exist', async () => {
      const { festivals, nextFestival, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const futureDate1 = dayjs().add(10, 'days').format('YYYY-MM-DD')
      const futureDate2 = dayjs().add(5, 'days').format('YYYY-MM-DD')
      const futureDate3 = dayjs().add(15, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...festivals.value[0], startDate: futureDate1 },
        { ...festivals.value[1], startDate: futureDate2 },
        { ...festivals.value[2], startDate: futureDate3 },
      ]

      expect(nextFestival.value?.startDate).toBe(futureDate2)
    })
  })

  describe('sortedFestivals', () => {
    it('should filter out past festivals', async () => {
      const { festivals, sortedFestivals, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const now = dayjs()
      const pastDate = now.subtract(1, 'month').format('YYYY-MM-DD')
      const futureDate = now.add(1, 'month').format('YYYY-MM-DD')
      const futureEndDate = now.add(2, 'months').format('YYYY-MM-DD')

      festivals.value = [
        { ...mockFestivals[0], id: '1', startDate: pastDate, endDate: pastDate },
        { ...mockFestivals[1], id: '2', startDate: futureDate, endDate: futureEndDate },
      ]

      const sorted = sortedFestivals.value
      expect(sorted.length).toBe(1)
      expect(sorted[0].id).toBe('2') // Only future festival
    })

    it('should sort upcoming festivals by date (earliest first)', async () => {
      const { festivals, sortedFestivals, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const futureDate1 = dayjs().add(10, 'days').format('YYYY-MM-DD')
      const futureDate2 = dayjs().add(5, 'days').format('YYYY-MM-DD')
      const futureDate3 = dayjs().add(15, 'days').format('YYYY-MM-DD')
      const futureEndDate = dayjs().add(20, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...mockFestivals[0], id: '1', startDate: futureDate1, endDate: futureEndDate },
        { ...mockFestivals[1], id: '2', startDate: futureDate2, endDate: futureEndDate },
        { ...mockFestivals[2], id: '3', startDate: futureDate3, endDate: futureEndDate },
      ]

      const sorted = sortedFestivals.value
      expect(sorted[0].id).toBe('2')
      expect(sorted[1].id).toBe('1')
      expect(sorted[2].id).toBe('3')
    })

    it('should not include festivals that have ended', async () => {
      const { festivals, sortedFestivals, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const pastDate1 = dayjs().subtract(10, 'days').format('YYYY-MM-DD')
      const pastDate2 = dayjs().subtract(5, 'days').format('YYYY-MM-DD')
      const pastDate3 = dayjs().subtract(15, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...mockFestivals[0], id: '1', startDate: pastDate1, endDate: pastDate1 },
        { ...mockFestivals[1], id: '2', startDate: pastDate2, endDate: pastDate2 },
        { ...mockFestivals[2], id: '3', startDate: pastDate3, endDate: pastDate3 },
      ]

      const sorted = sortedFestivals.value
      expect(sorted.length).toBe(0)
    })

    it('should include ongoing festivals (started but not ended)', async () => {
      const { festivals, sortedFestivals, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const pastStartDate = dayjs().subtract(5, 'days').format('YYYY-MM-DD')
      const futureEndDate = dayjs().add(5, 'days').format('YYYY-MM-DD')
      const futureDate = dayjs().add(10, 'days').format('YYYY-MM-DD')
      const futureEndDate2 = dayjs().add(15, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...mockFestivals[0], id: '1', startDate: pastStartDate, endDate: futureEndDate },
        { ...mockFestivals[1], id: '2', startDate: futureDate, endDate: futureEndDate2 },
      ]

      const sorted = sortedFestivals.value
      expect(sorted.length).toBe(2)
      expect(sorted[0].id).toBe('2') // Upcoming festival comes first
      expect(sorted[1].id).toBe('1') // Ongoing festival (already started, so not upcoming)
    })
  })

  describe('getFestivalsByRegion', () => {
    it('should return festivals for a specific region', async () => {
      const { festivals, getFestivalsByRegion, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const testRegion = 'Île-de-France'

      festivals.value[0].region = testRegion
      const result = getFestivalsByRegion(testRegion)

      expect(result.length).toBeGreaterThan(0)
      result.forEach(festival => {
        expect(festival.region).toBe(testRegion)
      })
    })

    it('should return empty array for region with no festivals', async () => {
      const { getFestivalsByRegion, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const result = getFestivalsByRegion('NonExistentRegion')
      expect(result).toEqual([])
    })

    it('should return all matching festivals for a region', async () => {
      const { festivals, getFestivalsByRegion, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const testRegion = 'TestRegion'

      festivals.value[0].region = testRegion
      festivals.value[1].region = testRegion
      festivals.value[2].region = 'OtherRegion'

      const result = getFestivalsByRegion(testRegion)
      expect(result.length).toBe(2)
    })
  })

  describe('getFestivalById', () => {
    it('should return festival with matching id', async () => {
      const { festivals, getFestivalById, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const testId = festivals.value[0].id
      const result = getFestivalById(testId)

      expect(result).toBeDefined()
      expect(result?.id).toBe(testId)
    })

    it('should return undefined for non-existent id', async () => {
      const { getFestivalById, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const result = getFestivalById('non-existent-id')
      expect(result).toBeUndefined()
    })

    it('should return the correct festival when multiple exist', async () => {
      const { festivals, getFestivalById, fetchFestivals } = useFestivals()
      await fetchFestivals()

      const targetFestival = festivals.value[2]
      const result = getFestivalById(targetFestival.id)

      expect(result).toEqual(targetFestival)
    })
  })
})
