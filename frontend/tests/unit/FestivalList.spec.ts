import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import FestivalList from '@/components/FestivalList/FestivalList.vue'
import FestivalCard from '@/components/FestivalCard/FestivalCard.vue'
import type { Festival } from '@/types'
import dayjs from 'dayjs'

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

describe('FestivalList', () => {
  describe('rendering', () => {
    it('should render the section title', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      expect(wrapper.find('[data-testid="section-title"]').text()).toBe('Tous les Festivals')
    })

    it('should render the section subtitle', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      expect(wrapper.find('[data-testid="section-subtitle"]').text()).toBe(
        'Découvrez tous les festivals de bière en France'
      )
    })

    it('should render the container', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      expect(wrapper.find('[data-testid="festival-list-container"]').exists()).toBe(true)
    })
  })

  describe('with festivals', () => {
    it('should render the festivals grid', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      expect(wrapper.find('[data-testid="festivals-grid"]').exists()).toBe(true)
    })

    it('should render one festival card', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      expect(cards).toHaveLength(1)
    })

    it('should render multiple festival cards', () => {
      const festivals = [createMockFestival('1'), createMockFestival('2'), createMockFestival('3')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      expect(cards).toHaveLength(3)
    })

    it('should render correct number of festival cards', () => {
      const festivals = Array.from({ length: 10 }, (_, i) => createMockFestival(String(i + 1)))
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      expect(cards).toHaveLength(10)
    })

    it('should pass festival data to each card', () => {
      const festivals = [createMockFestival('1', { name: 'Unique Festival Name' })]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const card = wrapper.findComponent(FestivalCard)
      expect(card.props('festival')).toEqual(festivals[0])
    })

    it('should pass showStatus prop to festival cards', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const card = wrapper.findComponent(FestivalCard)
      expect(card.props('showStatus')).toBe(true)
    })

    it('should not render empty state when festivals exist', () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      expect(wrapper.find('[data-testid="empty-state"]').exists()).toBe(false)
    })

    it('should render festivals with unique keys', () => {
      const festivals = [createMockFestival('1'), createMockFestival('2'), createMockFestival('3')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      const keys = cards.map(card => card.attributes('data-testid'))
      expect(keys).toEqual(['festival-card-1', 'festival-card-2', 'festival-card-3'])
    })
  })

  describe('without festivals', () => {
    it('should render empty state when no festivals', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: false },
      })

      expect(wrapper.find('[data-testid="empty-state"]').exists()).toBe(true)
    })

    it('should not render festivals grid when no festivals', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: false },
      })

      expect(wrapper.find('[data-testid="festivals-grid"]').exists()).toBe(false)
    })

    it('should render empty state message', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: false },
      })

      const emptyState = wrapper.find('[data-testid="empty-state"]')
      expect(emptyState.text()).toContain('Aucun festival')
      expect(emptyState.text()).toContain("Il n'y a actuellement aucun festival à afficher")
    })

    it('should not render any festival cards when empty', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      expect(cards).toHaveLength(0)
    })
  })

  describe('loading state', () => {
    it('should render loading state when loading is true', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: true },
      })

      const loadingState = wrapper.find('[data-testid="loading-state"]')
      expect(loadingState.exists()).toBe(true)
      expect(loadingState.text()).toContain('Chargement...')
    })

    it('should not render festivals grid when loading', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: true },
      })

      expect(wrapper.find('[data-testid="festivals-grid"]').exists()).toBe(false)
    })

    it('should not render empty state when loading', () => {
      const wrapper = mount(FestivalList, {
        props: { festivals: [], loading: true },
      })

      expect(wrapper.find('[data-testid="empty-state"]').exists()).toBe(false)
    })
  })

  describe('interactions', () => {
    it('should emit festival-click event when card is clicked', async () => {
      const festivals = [createMockFestival('1')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const card = wrapper.findComponent(FestivalCard)
      await card.vm.$emit('click', festivals[0])

      expect(wrapper.emitted('festival-click')).toBeTruthy()
      expect(wrapper.emitted('festival-click')?.[0]).toEqual([festivals[0]])
    })

    it('should emit festival-click event for correct festival', async () => {
      const festivals = [createMockFestival('1'), createMockFestival('2'), createMockFestival('3')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      await cards[1].vm.$emit('click', festivals[1])

      expect(wrapper.emitted('festival-click')?.[0]).toEqual([festivals[1]])
    })

    it('should handle multiple festival clicks', async () => {
      const festivals = [createMockFestival('1'), createMockFestival('2')]
      const wrapper = mount(FestivalList, {
        props: { festivals, loading: false },
      })

      const cards = wrapper.findAllComponents(FestivalCard)
      await cards[0].vm.$emit('click', festivals[0])
      await cards[1].vm.$emit('click', festivals[1])

      expect(wrapper.emitted('festival-click')).toHaveLength(2)
      expect(wrapper.emitted('festival-click')?.[0]).toEqual([festivals[0]])
      expect(wrapper.emitted('festival-click')?.[1]).toEqual([festivals[1]])
    })
  })
})
