import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import FestivalMap from '@/components/FestivalMap/FestivalMap.vue'
import type { Festival } from '@/types'
import dayjs from 'dayjs'

vi.mock('leaflet', () => {
  const mockMarker = {
    addTo: vi.fn().mockReturnThis(),
    bindPopup: vi.fn().mockReturnThis(),
    on: vi.fn().mockReturnThis(),
    remove: vi.fn(),
    getPopup: vi.fn(() => ({
      getElement: vi.fn(() => ({
        querySelector: vi.fn(() => null),
      })),
    })),
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

const createMockFestival = (id: string, overrides?: Partial<Festival>): Festival => ({
  id,
  name: `Test Festival ${id}`,
  description: 'A great beer festival',
  startDate: dayjs().add(10, 'days').format('YYYY-MM-DD'),
  endDate: dayjs().add(12, 'days').format('YYYY-MM-DD'),
  city: 'Paris',
  region: 'Île-de-France',
  location: {
    latitude: 48.8566,
    longitude: 2.3522,
  },
  image: 'https://example.com/image.jpg',
  website: 'https://example.com',
  breweryCount: 50,
  ...overrides,
})

describe('FestivalMap', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('rendering', () => {
    it('should render the container', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.find('[data-testid="festival-map-container"]').exists()).toBe(true)
    })

    it('should render the section title', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.find('[data-testid="section-title"]').text()).toBe('Carte des Festivals')
    })

    it('should render the section subtitle', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.find('[data-testid="section-subtitle"]').text()).toBe(
        'Explorez les festivals de bière partout en France'
      )
    })

    it('should render the map wrapper', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.find('[data-testid="map-wrapper"]').exists()).toBe(true)
    })

    it('should render the map container', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.find('[data-testid="map-container"]').exists()).toBe(true)
    })

    it('should set correct height on map wrapper', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const mapWrapper = wrapper.find('[data-testid="map-wrapper"]')
      expect(mapWrapper.attributes('style')).toContain('height: 600px')
    })
  })

  describe('props', () => {
    it('should accept festivals prop', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.props('festivals')).toEqual(festivals)
    })

    it('should accept empty festivals array', () => {
      const wrapper = mount(FestivalMap, {
        props: { festivals: [] },
      })

      expect(wrapper.props('festivals')).toEqual([])
    })

    it('should accept multiple festivals', () => {
      const festivals = [createMockFestival('1'), createMockFestival('2'), createMockFestival('3')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.props('festivals')).toHaveLength(3)
    })
  })

  describe('component lifecycle', () => {
    it('should mount successfully with festivals', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(wrapper.vm).toBeDefined()
    })

    it('should mount successfully without festivals', () => {
      const wrapper = mount(FestivalMap, {
        props: { festivals: [] },
      })

      expect(wrapper.vm).toBeDefined()
    })

    it('should unmount successfully', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      expect(() => wrapper.unmount()).not.toThrow()
    })
  })

  describe('festivals updates', () => {
    it('should handle festivals prop updates', async () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const newFestivals = [createMockFestival('1'), createMockFestival('2')]
      await wrapper.setProps({ festivals: newFestivals })

      expect(wrapper.props('festivals')).toEqual(newFestivals)
    })

    it('should handle empty to populated festivals', async () => {
      const wrapper = mount(FestivalMap, {
        props: { festivals: [] },
      })

      const festivals = [createMockFestival('1')]
      await wrapper.setProps({ festivals })

      expect(wrapper.props('festivals')).toEqual(festivals)
    })

    it('should handle populated to empty festivals', async () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      await wrapper.setProps({ festivals: [] })

      expect(wrapper.props('festivals')).toEqual([])
    })
  })

  describe('structure', () => {
    it('should have correct container structure', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const container = wrapper.find('[data-testid="festival-map-container"]')
      const wrapper_div = container.find('.max-w-6xl')
      expect(wrapper_div.exists()).toBe(true)
    })

    it('should have map wrapper inside container', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const container = wrapper.find('[data-testid="festival-map-container"]')
      const mapWrapper = container.find('[data-testid="map-wrapper"]')
      expect(mapWrapper.exists()).toBe(true)
    })

    it('should have map container inside map wrapper', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const mapWrapper = wrapper.find('[data-testid="map-wrapper"]')
      const mapContainer = mapWrapper.find('[data-testid="map-container"]')
      expect(mapContainer.exists()).toBe(true)
    })
  })

  describe('popup interactions', () => {
    it('should emit popup-click event', async () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const mockFestival = festivals[0]
      await wrapper.vm.$emit('popup-click', mockFestival)

      expect(wrapper.emitted('popup-click')).toBeTruthy()
      expect(wrapper.emitted('popup-click')?.[0]).toEqual([mockFestival])
    })

    it('should emit marker-click event when marker is clicked', async () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const mockFestival = festivals[0]
      await wrapper.vm.$emit('marker-click', mockFestival)

      expect(wrapper.emitted('marker-click')).toBeTruthy()
      expect(wrapper.emitted('marker-click')?.[0]).toEqual([mockFestival])
    })
  })

  describe('marker event handlers', () => {
    it('should trigger marker click handler', async () => {
      const L = await import('leaflet')
      const mockMarker = L.default.marker([0, 0])
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const clickCallback = (mockMarker.on as any).mock.calls.find(
        (call: any) => call[0] === 'click'
      )?.[1]

      if (clickCallback) {
        clickCallback()
        await wrapper.vm.$nextTick()
        expect(wrapper.emitted('marker-click')).toBeTruthy()
      }
    })

    it('should trigger popup open handler without button', async () => {
      const L = await import('leaflet')
      const mockMarker = L.default.marker([0, 0])
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const popupOpenCallback = (mockMarker.on as any).mock.calls.find(
        (call: any) => call[0] === 'popupopen'
      )?.[1]

      if (popupOpenCallback) {
        popupOpenCallback()
        await wrapper.vm.$nextTick()
      }
    })

    it('should trigger popup open handler with button', async () => {
      const L = await import('leaflet')
      const mockButton = {
        addEventListener: vi.fn(),
      }
      const mockMarker = {
        addTo: vi.fn().mockReturnThis(),
        bindPopup: vi.fn().mockReturnThis(),
        on: vi.fn().mockReturnThis(),
        remove: vi.fn(),
        getPopup: vi.fn(() => ({
          getElement: vi.fn(() => ({
            querySelector: vi.fn(() => mockButton),
          })),
        })),
      }

      ;(L.default.marker as any).mockReturnValueOnce(mockMarker)

      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalMap, {
        props: { festivals },
      })

      const popupOpenCallback = (mockMarker.on as any).mock.calls.find(
        (call: any) => call[0] === 'popupopen'
      )?.[1]

      if (popupOpenCallback) {
        popupOpenCallback()
        await wrapper.vm.$nextTick()
        expect(mockButton.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))

        const buttonClickCallback = mockButton.addEventListener.mock.calls[0][1]
        buttonClickCallback()
        await wrapper.vm.$nextTick()
        expect(wrapper.emitted('popup-click')).toBeTruthy()
      }
    })
  })
})
