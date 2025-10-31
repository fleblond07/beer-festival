import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import NextFestival from '@/components/NextFestival/NextFestival.vue'
import type { Festival } from '@/types'
import dayjs from 'dayjs'

const createMockFestival = (overrides?: Partial<Festival>): Festival => ({
  id: '1',
  name: 'Test Festival',
  description: 'A great beer festival in the heart of France',
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

describe('NextFestival', () => {
  describe('with festival', () => {
    it('should render the section title', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="section-title"]').text()).toBe('Prochain Festival')
    })

    it('should render the section subtitle', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="section-subtitle"]').text()).toBe(
        'Ne manquez pas le prochain événement !'
      )
    })

    it('should render the festival name', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="festival-name"]').text()).toBe(festival.name)
    })

    it('should render the festival location', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      const location = wrapper.find('[data-testid="festival-location"]').text()
      expect(location).toContain(festival.city)
    })

    it('should render the festival dates', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      const dates = wrapper.find('[data-testid="festival-dates"]')
      expect(dates.exists()).toBe(true)
      expect(dates.text().length).toBeGreaterThan(0)
    })

    it('should render the festival description', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="festival-description"]').text()).toBe(festival.description)
    })

    it('should render the countdown badge', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      const countdown = wrapper.find('[data-testid="countdown-badge"]')
      expect(countdown.exists()).toBe(true)
      expect(countdown.text().length).toBeGreaterThan(0)
    })

    it('should render the festival image when provided', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      const image = wrapper.find('[data-testid="hero-image"]')
      expect(image.exists()).toBe(true)
      expect(image.attributes('src')).toBe(festival.image)
      expect(image.attributes('alt')).toBe(festival.name)
    })

    it('should not render the image container when image is not provided', () => {
      const festival = createMockFestival({ image: undefined })
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="hero-image-container"]').exists()).toBe(false)
    })

    it('should render brewery information when provided', () => {
      const festival = createMockFestival({ breweryCount: 75 })
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      const breweryInfo = wrapper.find('[data-testid="brewery-info"]')
      expect(breweryInfo.exists()).toBe(true)
      expect(breweryInfo.text()).toContain('75')
      expect(breweryInfo.text()).toContain('brasseries')
    })

    it('should not render brewery information when not provided', () => {
      const festival = createMockFestival({ breweryCount: undefined })
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="brewery-info"]').exists()).toBe(false)
    })

    it('should render website link when provided', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      const link = wrapper.find('[data-testid="website-link"]')
      expect(link.exists()).toBe(true)
      expect(link.attributes('href')).toBe(festival.website)
      expect(link.attributes('target')).toBe('_blank')
      expect(link.attributes('rel')).toBe('noopener noreferrer')
    })

    it('should not render website link when not provided', () => {
      const festival = createMockFestival({ website: undefined })
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="website-link"]').exists()).toBe(false)
    })

    it('should render the hero card', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="hero-card"]').exists()).toBe(true)
    })

    it('should not render no-festival message when festival exists', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="no-festival"]').exists()).toBe(false)
    })
  })

  describe('without festival', () => {
    it('should render no-festival message when festival is null', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: false },
      })

      const noFestival = wrapper.find('[data-testid="no-festival"]')
      expect(noFestival.exists()).toBe(true)
      expect(noFestival.text()).toContain('Aucun festival à venir')
    })

    it('should not render hero card when festival is null', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: false },
      })

      expect(wrapper.find('[data-testid="hero-card"]').exists()).toBe(false)
    })

    it('should render helpful message when no festival', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: false },
      })

      const noFestival = wrapper.find('[data-testid="no-festival"]')
      expect(noFestival.text()).toContain('Revenez bientôt')
    })

    it('should return empty string for formattedDateRange when festival is null', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: false },
      })

      expect(wrapper.find('[data-testid="festival-dates"]').exists()).toBe(false)
    })

    it('should return empty string for countdownText when festival is null', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: false },
      })

      expect(wrapper.find('[data-testid="countdown-badge"]').exists()).toBe(false)
    })
  })

  describe('loading state', () => {
    it('should render loading state when loading is true and no festival', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: true },
      })

      const loadingState = wrapper.find('[data-testid="loading-state"]')
      expect(loadingState.exists()).toBe(true)
      expect(loadingState.text()).toContain('Chargement...')
    })

    it('should not render no-festival message when loading', () => {
      const wrapper = mount(NextFestival, {
        props: { festival: null, loading: true },
      })

      expect(wrapper.find('[data-testid="no-festival"]').exists()).toBe(false)
    })

    it('should still render festival when loading and festival exists', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: true },
      })

      expect(wrapper.find('[data-testid="hero-card"]').exists()).toBe(true)
    })
  })

  describe('container', () => {
    it('should render the container', () => {
      const festival = createMockFestival()
      const wrapper = mount(NextFestival, {
        props: { festival, loading: false },
      })

      expect(wrapper.find('[data-testid="next-festival-container"]').exists()).toBe(true)
    })
  })
})
