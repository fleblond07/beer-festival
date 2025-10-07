import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FestivalCard from '@/components/FestivalCard/FestivalCard.vue'
import type { Festival } from '@/types'
import dayjs from 'dayjs'

const createMockFestival = (overrides?: Partial<Festival>): Festival => ({
  id: '1',
  name: 'Test Festival',
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

describe('FestivalCard', () => {
  describe('rendering', () => {
    it('should render the festival name', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="festival-name"]').text()).toBe(festival.name)
    })

    it('should render the festival location', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const location = wrapper.find('[data-testid="festival-location"]').text()
      expect(location).toContain(festival.city)
      expect(location).toContain(festival.region)
    })

    it('should render the festival description', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="festival-description"]').text()).toBe(festival.description)
    })

    it('should render the festival dates', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const dates = wrapper.find('[data-testid="festival-dates"]')
      expect(dates.exists()).toBe(true)
      expect(dates.text().length).toBeGreaterThan(0)
    })

    it('should render the festival image when provided', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const image = wrapper.find('[data-testid="festival-image"]')
      expect(image.exists()).toBe(true)
      expect(image.attributes('src')).toBe(festival.image)
      expect(image.attributes('alt')).toBe(festival.name)
    })

    it('should not render the image when not provided', () => {
      const festival = createMockFestival({ image: undefined })
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const image = wrapper.find('[data-testid="festival-image"]')
      expect(image.exists()).toBe(false)
    })

    it('should render the brewery count when provided', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const breweryCount = wrapper.find('[data-testid="brewery-count"]')
      expect(breweryCount.exists()).toBe(true)
      expect(breweryCount.text()).toContain(String(festival.breweryCount))
    })

    it('should not render brewery count when not provided', () => {
      const festival = createMockFestival({ breweryCount: undefined })
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const breweryCount = wrapper.find('[data-testid="brewery-count"]')
      expect(breweryCount.exists()).toBe(false)
    })

    it('should render the website link when provided', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const website = wrapper.find('[data-testid="festival-website"]')
      expect(website.exists()).toBe(true)
      expect(website.attributes('href')).toBe(festival.website)
      expect(website.attributes('target')).toBe('_blank')
      expect(website.attributes('rel')).toBe('noopener noreferrer')
    })

    it('should not render website link when not provided', () => {
      const festival = createMockFestival({ website: undefined })
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const website = wrapper.find('[data-testid="festival-website"]')
      expect(website.exists()).toBe(false)
    })
  })

  describe('status badge', () => {
    it('should show status badge by default', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="festival-status"]').exists()).toBe(true)
    })

    it('should hide status badge when showStatus is false', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival, showStatus: false },
      })

      expect(wrapper.find('[data-testid="festival-status"]').exists()).toBe(false)
    })

    it('should show green status for upcoming festivals', () => {
      const festival = createMockFestival({
        startDate: dayjs().add(10, 'days').format('YYYY-MM-DD'),
      })
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const status = wrapper.find('[data-testid="festival-status"]')
      expect(status.classes()).toContain('bg-green-100')
      expect(status.classes()).toContain('text-green-800')
    })

    it('should show gray status for past festivals', () => {
      const festival = createMockFestival({
        startDate: dayjs().subtract(10, 'days').format('YYYY-MM-DD'),
        endDate: dayjs().subtract(8, 'days').format('YYYY-MM-DD'),
      })
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const status = wrapper.find('[data-testid="festival-status"]')
      expect(status.classes()).toContain('bg-gray-100')
      expect(status.classes()).toContain('text-gray-800')
    })
  })

  describe('interactions', () => {
    it('should emit click event when card is clicked', async () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      await wrapper.find('[data-testid="festival-card"]').trigger('click')

      expect(wrapper.emitted('click')).toBeTruthy()
      expect(wrapper.emitted('click')?.[0]).toEqual([festival])
    })

    it('should not emit click event when website link is clicked', async () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      const website = wrapper.find('[data-testid="festival-website"]')
      await website.trigger('click')

      expect(wrapper.emitted('click')).toBeFalsy()
    })

    it('should have cursor-pointer class for interactivity', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="festival-card"]').classes()).toContain('cursor-pointer')
    })
  })

  describe('icons', () => {
    it('should render location icon', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="location-icon"]').exists()).toBe(true)
    })

    it('should render calendar icon', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="calendar-icon"]').exists()).toBe(true)
    })

    it('should render brewery icon when brewery count exists', () => {
      const festival = createMockFestival()
      const wrapper = mount(FestivalCard, {
        props: { festival },
      })

      expect(wrapper.find('[data-testid="brewery-icon"]').exists()).toBe(true)
    })
  })
})
