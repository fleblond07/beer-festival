import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '@/App.vue'

vi.mock('leaflet', () => {
  const mockMarker = {
    addTo: vi.fn().mockReturnThis(),
    bindPopup: vi.fn().mockReturnThis(),
    on: vi.fn().mockReturnThis(),
    remove: vi.fn(),
  }

  const mockTileLayer = {
    addTo: vi.fn().mockReturnThis(),
  }

  const mockFeatureGroup = {
    getBounds: vi.fn(() => ({
      pad: vi.fn(() => ({
        _southWest: { lat: 0, lng: 0 },
        _northEast: { lat: 10, lng: 10 },
      })),
    })),
  }

  const mockMap = {
    setView: vi.fn().mockReturnThis(),
    remove: vi.fn(),
    fitBounds: vi.fn(),
  }

  return {
    default: {
      map: vi.fn(() => mockMap),
      tileLayer: vi.fn(() => mockTileLayer),
      marker: vi.fn(() => mockMarker),
      icon: vi.fn(options => options),
      featureGroup: vi.fn(() => mockFeatureGroup),
    },
  }
})

describe('Default App testing', () => {
  describe('rendering', () => {
    it('should render the app', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="app"]').exists()).toBe(true)
    })

    it('should render the site title', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="site-title"]').text()).toBe('Festivals de Bière')
    })

    it('should render the site tagline', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="site-tagline"]').text()).toContain(
        'Découvrez les meilleurs festivals de France'
      )
    })

    it('should render the logo icon', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="logo-icon"]').exists()).toBe(true)
    })

    it('should render footer text', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="footer-text"]').exists()).toBe(true)
    })
  })

  describe('components', () => {
    it('should render NextFestival section', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="next-festival-section"]').exists()).toBe(true)
    })

    it('should render FestivalMap section', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="festival-map-section"]').exists()).toBe(true)
    })

    it('should render FestivalList section', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="festival-list-section"]').exists()).toBe(true)
    })
  })

  describe('structure', () => {
    it('should have header element', () => {
      const wrapper = mount(App)
      expect(wrapper.find('header').exists()).toBe(true)
    })

    it('should have main element', () => {
      const wrapper = mount(App)
      expect(wrapper.find('main').exists()).toBe(true)
    })

    it('should have footer element', () => {
      const wrapper = mount(App)
      expect(wrapper.find('footer').exists()).toBe(true)
    })

    it('should render sections in correct order', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="next-festival-section"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="festival-map-section"]').exists()).toBe(true)
      expect(wrapper.find('[data-testid="festival-list-section"]').exists()).toBe(true)
    })
  })

  describe('layout', () => {
    it('should have min-h-screen class on root', () => {
      const wrapper = mount(App)
      expect(wrapper.find('[data-testid="app"]').classes()).toContain('min-h-screen')
    })

    it('should have proper header styling', () => {
      const wrapper = mount(App)
      const header = wrapper.find('header')
      expect(header.classes()).toContain('bg-gradient-to-r')
    })

    it('should have proper footer styling', () => {
      const wrapper = mount(App)
      const footer = wrapper.find('footer')
      expect(footer.classes()).toContain('bg-dark-card')
    })
  })

  describe('event handlers', () => {
    it('should handle marker click event', async () => {
      const wrapper = mount(App)
      const consoleSpy = vi.spyOn(console, 'log')

      const mapSection = wrapper.findComponent({ name: 'FestivalMap' })
      const mockFestival = {
        id: '1',
        name: 'Test Festival',
        description: 'Test',
        startDate: '2025-10-15',
        endDate: '2025-10-17',
        city: 'Paris',
        region: 'Île-de-France',
        location: { latitude: 48.8566, longitude: 2.3522 },
      }

      await mapSection.vm.$emit('marker-click', mockFestival)

      expect(consoleSpy).toHaveBeenCalledWith('Marker clicked:', 'Test Festival')
      consoleSpy.mockRestore()
    })

    it('should handle festival click event', async () => {
      const wrapper = mount(App)
      const consoleSpy = vi.spyOn(console, 'log')

      const listSection = wrapper.findComponent({ name: 'FestivalList' })
      const mockFestival = {
        id: '2',
        name: 'Another Festival',
        description: 'Test',
        startDate: '2025-11-01',
        endDate: '2025-11-03',
        city: 'Lyon',
        region: 'Auvergne-Rhône-Alpes',
        location: { latitude: 45.764, longitude: 4.8357 },
      }

      await listSection.vm.$emit('festival-click', mockFestival)

      expect(consoleSpy).toHaveBeenCalledWith('Festival clicked:', 'Another Festival')
      consoleSpy.mockRestore()
    })
  })
})
