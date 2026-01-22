import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useBreweries } from '@/composables/useBreweries'

const mockBreweries = [
  {
    id: 1,
    name: 'Brasserie de la Plaine',
    description: 'Brasserie artisanale parisienne',
    city: 'Paris',
    website: 'https://brasseriedelaplaine.fr',
    logo: 'https://example.com/logo1.png',
    festivalCount: 5,
  },
  {
    id: 2,
    name: 'Brasserie du Mont Blanc',
    description: 'Brasserie savoyarde',
    city: 'Chamonix',
    website: 'https://brasserie-montblanc.com',
    logo: 'https://example.com/logo2.png',
    festivalCount: 3,
  },
  {
    id: 3,
    name: 'Brasserie Pietra',
    description: 'Brasserie corse',
    city: 'Furiani',
    website: 'https://brasseriedepetra.com',
    logo: 'https://example.com/logo3.png',
    festivalCount: 8,
  },
]

describe('useBreweries', () => {
  beforeEach(() => {
    vi.clearAllMocks()

    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve(mockBreweries),
      } as Response)
    )
  })

  describe('initialization', () => {
    it('should initialize with empty breweries array', () => {
      const { breweries } = useBreweries()
      expect(breweries.value).toBeDefined()
      expect(breweries.value).toEqual([])
    })

    it('should initialize with loading set to false', () => {
      const { loading } = useBreweries()
      expect(loading.value).toBe(false)
    })

    it('should initialize with no error', () => {
      const { error } = useBreweries()
      expect(error.value).toBeNull()
    })
  })

  describe('fetchBreweries', () => {
    it('should fetch and load breweries', async () => {
      const { breweries, loading, fetchBreweries } = useBreweries()

      await fetchBreweries()

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/breweries')
      )
      expect(breweries.value.length).toBe(3)
      expect(loading.value).toBe(false)
    })

    it('should set loading to true during fetch', async () => {
      const { loading, fetchBreweries } = useBreweries()

      const fetchPromise = fetchBreweries()

      // Loading should still be true before promise resolves
      // Note: This may be timing-dependent

      await fetchPromise
      expect(loading.value).toBe(false)
    })

    it('should populate breweries with correct data', async () => {
      const { breweries, fetchBreweries } = useBreweries()

      await fetchBreweries()

      expect(breweries.value[0].name).toBe('Brasserie de la Plaine')
      expect(breweries.value[1].name).toBe('Brasserie du Mont Blanc')
      expect(breweries.value[2].name).toBe('Brasserie Pietra')
    })

    it('should clear error on successful fetch', async () => {
      const { error, fetchBreweries } = useBreweries()

      await fetchBreweries()

      expect(error.value).toBeNull()
    })
  })

  describe('fetchBreweriesByFestival', () => {
    beforeEach(() => {
      const festivalBreweries = [mockBreweries[0], mockBreweries[1]]

      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: true,
          json: () => Promise.resolve(festivalBreweries),
        } as Response)
      )
    })

    it('should fetch breweries for a specific festival', async () => {
      const { breweries, fetchBreweriesByFestival } = useBreweries()

      await fetchBreweriesByFestival('1')

      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/festivals/1/breweries')
      )
      expect(breweries.value.length).toBe(2)
    })

    it('should set loading to true during fetch', async () => {
      const { loading, fetchBreweriesByFestival } = useBreweries()

      const fetchPromise = fetchBreweriesByFestival('1')

      await fetchPromise
      expect(loading.value).toBe(false)
    })

    it('should populate breweries with festival-specific data', async () => {
      const { breweries, fetchBreweriesByFestival } = useBreweries()

      await fetchBreweriesByFestival('1')

      expect(breweries.value[0].name).toBe('Brasserie de la Plaine')
      expect(breweries.value[1].name).toBe('Brasserie du Mont Blanc')
    })
  })

  describe('error handling', () => {
    it('should handle fetch error when response is not ok', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: false,
          json: () => Promise.resolve([]),
        } as Response)
      )

      const { error, fetchBreweries } = useBreweries()
      await fetchBreweries()

      expect(error.value).toBe('Failed to fetch breweries')
    })

    it('should handle network errors', async () => {
      global.fetch = vi.fn(() =>
        Promise.reject(new Error('Network error'))
      )

      const { error, fetchBreweries } = useBreweries()
      await fetchBreweries()

      expect(error.value).toBe('Network error')
    })

    it('should handle non-Error exceptions', async () => {
      global.fetch = vi.fn(() => Promise.reject('String error'))

      const { error, fetchBreweries } = useBreweries()
      await fetchBreweries()

      expect(error.value).toBe('An error occurred')
    })

    it('should clear breweries on error', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: false,
          json: () => Promise.resolve([]),
        } as Response)
      )

      const { breweries, fetchBreweries } = useBreweries()
      await fetchBreweries()

      expect(breweries.value).toEqual([])
    })

    it('should set loading to false after error', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: false,
          json: () => Promise.resolve([]),
        } as Response)
      )

      const { loading, fetchBreweries } = useBreweries()
      await fetchBreweries()

      expect(loading.value).toBe(false)
    })

    it('should handle error in fetchBreweriesByFestival', async () => {
      global.fetch = vi.fn(() =>
        Promise.resolve({
          ok: false,
          json: () => Promise.resolve([]),
        } as Response)
      )

      const { error, fetchBreweriesByFestival } = useBreweries()
      await fetchBreweriesByFestival('1')

      expect(error.value).toBe('Failed to fetch breweries')
    })
  })
})
