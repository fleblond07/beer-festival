import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import Brasseries from '@/views/Brasseries.vue'
import { useBreweries } from '@/composables/useBreweries'

vi.mock('@/composables/useBreweries')

describe('Brasseries', () => {
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
  ]

  beforeEach(() => {
    vi.clearAllMocks()
  })

  const mountWithMocks = (
    breweries = mockBreweries,
    loading = false,
    error = null,
    fetchBreweries = vi.fn()
  ) => {
    vi.mocked(useBreweries).mockReturnValue({
      breweries: { value: breweries },
      loading: { value: loading },
      error: { value: error },
      fetchBreweries,
      fetchBreweriesByFestival: vi.fn(),
    })

    return mount(Brasseries, {
      global: {
        stubs: {
          BreweryList: {
            template: '<div data-testid="brewery-list"></div>',
            props: ['breweries', 'loading'],
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('should render the component', () => {
      const wrapper = mountWithMocks()
      expect(wrapper.exists()).toBe(true)
    })

    it('should render BreweryList component', () => {
      const wrapper = mountWithMocks()
      expect(wrapper.find('[data-testid="brewery-list"]').exists()).toBe(true)
    })

    it('should use useBreweries composable', () => {
      mountWithMocks()
      expect(useBreweries).toHaveBeenCalled()
    })
  })

  describe('lifecycle', () => {
    it('should fetch breweries on mount', async () => {
      const fetchBreweries = vi.fn()
      mountWithMocks(mockBreweries, false, null, fetchBreweries)

      // Wait for nextTick to ensure onMounted has run
      await new Promise(resolve => setTimeout(resolve, 0))

      expect(fetchBreweries).toHaveBeenCalledTimes(1)
    })

    it('should log error when error occurs during fetch', async () => {
      const consoleLogSpy = vi.spyOn(console, 'log')
      const errorMessage = 'Failed to fetch breweries'
      const fetchBreweries = vi.fn()

      vi.mocked(useBreweries).mockReturnValue({
        breweries: { value: [] },
        loading: { value: false },
        error: { value: errorMessage },
        fetchBreweries,
        fetchBreweriesByFestival: vi.fn(),
      })

      mount(Brasseries, {
        global: {
          stubs: {
            BreweryList: {
              template: '<div></div>',
            },
          },
        },
      })

      await new Promise(resolve => setTimeout(resolve, 0))

      expect(consoleLogSpy).toHaveBeenCalledWith(errorMessage)
      consoleLogSpy.mockRestore()
    })

    it('should not log when there is no error', async () => {
      const consoleLogSpy = vi.spyOn(console, 'log')
      const fetchBreweries = vi.fn()

      mountWithMocks(mockBreweries, false, null, fetchBreweries)

      await new Promise(resolve => setTimeout(resolve, 0))

      expect(consoleLogSpy).not.toHaveBeenCalled()
      consoleLogSpy.mockRestore()
    })
  })

  describe('styling', () => {
    it('should have min-h-screen class', () => {
      const wrapper = mountWithMocks()
      const root = wrapper.find('div')
      expect(root.classes()).toContain('min-h-screen')
    })

    it('should have bg-dark-bg class', () => {
      const wrapper = mountWithMocks()
      const root = wrapper.find('div')
      expect(root.classes()).toContain('bg-dark-bg')
    })

    it('should have text-white class', () => {
      const wrapper = mountWithMocks()
      const root = wrapper.find('div')
      expect(root.classes()).toContain('text-white')
    })
  })
})
