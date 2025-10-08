import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useFestivals } from '@/composables/useFestivals'
import dayjs from 'dayjs'

describe('useFestivals', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('initialization', () => {
    it('should initialize with mock festivals', () => {
      const { festivals } = useFestivals()
      expect(festivals.value).toBeDefined()
      expect(festivals.value.length).toBeGreaterThan(0)
    })

    it('should initialize with loading set to false', () => {
      const { loading } = useFestivals()
      expect(loading.value).toBe(false)
    })

    it('should initialize with no error', () => {
      const { error } = useFestivals()
      expect(error.value).toBeNull()
    })

    it('should have images for all festivals', () => {
      const { festivals } = useFestivals()
      festivals.value.forEach(festival => {
        expect(festival.image).toBeDefined()
        expect(festival.image).toBeTruthy()
        expect(festival.image).toMatch(/^https?:\/\//)
      })
    })

    it('should have unique images for each festival', () => {
      const { festivals } = useFestivals()
      const images = festivals.value.map(f => f.image)
      const uniqueImages = new Set(images)
      expect(uniqueImages.size).toBe(images.length)
    })
  })

  describe('nextFestival', () => {
    it('should return the next upcoming festival', () => {
      const { nextFestival } = useFestivals()
      const now = dayjs()

      if (nextFestival.value) {
        const festivalDate = dayjs(nextFestival.value.startDate)
        expect(festivalDate.isAfter(now)).toBe(true)
      }
    })

    it('should return null if there are no upcoming festivals', () => {
      const { festivals, nextFestival } = useFestivals()
      const pastDate = dayjs().subtract(1, 'year').format('YYYY-MM-DD')
      festivals.value = festivals.value.map(festival => ({
        ...festival,
        startDate: pastDate,
        endDate: pastDate,
      }))

      expect(nextFestival.value).toBeNull()
    })

    it('should return the closest upcoming festival when multiple exist', () => {
      const { festivals, nextFestival } = useFestivals()
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
    it('should filter out past festivals', () => {
      const { festivals, sortedFestivals } = useFestivals()
      const now = dayjs()
      const pastDate = now.subtract(1, 'month').format('YYYY-MM-DD')
      const futureDate = now.add(1, 'month').format('YYYY-MM-DD')
      const futureEndDate = now.add(2, 'months').format('YYYY-MM-DD')

      festivals.value = [
        { ...festivals.value[0], id: '1', startDate: pastDate, endDate: pastDate },
        { ...festivals.value[1], id: '2', startDate: futureDate, endDate: futureEndDate },
      ]

      const sorted = sortedFestivals.value
      expect(sorted.length).toBe(1)
      expect(sorted[0].id).toBe('2') // Only future festival
    })

    it('should sort upcoming festivals by date (earliest first)', () => {
      const { festivals, sortedFestivals } = useFestivals()
      const futureDate1 = dayjs().add(10, 'days').format('YYYY-MM-DD')
      const futureDate2 = dayjs().add(5, 'days').format('YYYY-MM-DD')
      const futureDate3 = dayjs().add(15, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...festivals.value[0], id: '1', startDate: futureDate1 },
        { ...festivals.value[1], id: '2', startDate: futureDate2 },
        { ...festivals.value[2], id: '3', startDate: futureDate3 },
      ]

      const sorted = sortedFestivals.value
      expect(sorted[0].id).toBe('2')
      expect(sorted[1].id).toBe('1')
      expect(sorted[2].id).toBe('3')
    })

    it('should not include festivals that have ended', () => {
      const { festivals, sortedFestivals } = useFestivals()
      const pastDate1 = dayjs().subtract(10, 'days').format('YYYY-MM-DD')
      const pastDate2 = dayjs().subtract(5, 'days').format('YYYY-MM-DD')
      const pastDate3 = dayjs().subtract(15, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...festivals.value[0], id: '1', startDate: pastDate1, endDate: pastDate1 },
        { ...festivals.value[1], id: '2', startDate: pastDate2, endDate: pastDate2 },
        { ...festivals.value[2], id: '3', startDate: pastDate3, endDate: pastDate3 },
      ]

      const sorted = sortedFestivals.value
      expect(sorted.length).toBe(0)
    })

    it('should include ongoing festivals (started but not ended)', () => {
      const { festivals, sortedFestivals } = useFestivals()
      const pastStartDate = dayjs().subtract(5, 'days').format('YYYY-MM-DD')
      const futureEndDate = dayjs().add(5, 'days').format('YYYY-MM-DD')
      const futureDate = dayjs().add(10, 'days').format('YYYY-MM-DD')
      const futureEndDate2 = dayjs().add(15, 'days').format('YYYY-MM-DD')

      festivals.value = [
        { ...festivals.value[0], id: '1', startDate: pastStartDate, endDate: futureEndDate },
        { ...festivals.value[1], id: '2', startDate: futureDate, endDate: futureEndDate2 },
      ]

      const sorted = sortedFestivals.value
      expect(sorted.length).toBe(2)
      expect(sorted[0].id).toBe('2') // Upcoming festival comes first
      expect(sorted[1].id).toBe('1') // Ongoing festival (already started, so not upcoming)
    })
  })

  describe('getFestivalsByRegion', () => {
    it('should return festivals for a specific region', () => {
      const { festivals, getFestivalsByRegion } = useFestivals()
      const testRegion = 'Île-de-France'

      festivals.value[0].region = testRegion
      const result = getFestivalsByRegion(testRegion)

      expect(result.length).toBeGreaterThan(0)
      result.forEach(festival => {
        expect(festival.region).toBe(testRegion)
      })
    })

    it('should return empty array for region with no festivals', () => {
      const { getFestivalsByRegion } = useFestivals()
      const result = getFestivalsByRegion('NonExistentRegion')
      expect(result).toEqual([])
    })

    it('should return all matching festivals for a region', () => {
      const { festivals, getFestivalsByRegion } = useFestivals()
      const testRegion = 'TestRegion'

      festivals.value[0].region = testRegion
      festivals.value[1].region = testRegion
      festivals.value[2].region = 'OtherRegion'

      const result = getFestivalsByRegion(testRegion)
      expect(result.length).toBe(2)
    })
  })

  describe('getFestivalById', () => {
    it('should return festival with matching id', () => {
      const { festivals, getFestivalById } = useFestivals()
      const testId = festivals.value[0].id
      const result = getFestivalById(testId)

      expect(result).toBeDefined()
      expect(result?.id).toBe(testId)
    })

    it('should return undefined for non-existent id', () => {
      const { getFestivalById } = useFestivals()
      const result = getFestivalById('non-existent-id')
      expect(result).toBeUndefined()
    })

    it('should return the correct festival when multiple exist', () => {
      const { festivals, getFestivalById } = useFestivals()
      const targetFestival = festivals.value[2]
      const result = getFestivalById(targetFestival.id)

      expect(result).toEqual(targetFestival)
    })
  })
})
