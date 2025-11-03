import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import App from '@/App.vue'

describe('Default App testing', () => {
  const mountWithRouter = () => {
    return mount(App, {
      global: {
        stubs: {
          RouterLink: {
            template: '<a><slot /></a>',
          },
          RouterView: {
            template: '<div></div>',
          },
        },
      },
    })
  }

  describe('rendering', () => {
    it('should render the app', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="app"]').exists()).toBe(true)
    })

    it('should render the site title', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="site-title"]').text()).toBe('Festivals de Bière')
    })

    it('should render the site tagline', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="site-tagline"]').text()).toContain(
        'Découvrez les meilleurs festivals de France'
      )
    })

    it('should render the logo icon', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="logo-icon"]').exists()).toBe(true)
    })

    it('should render footer text', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="footer-text"]').exists()).toBe(true)
    })

    it('should render lock icon', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="lock-icon"]').exists()).toBe(true)
    })

    it('should render admin link', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="admin-link"]').exists()).toBe(true)
    })
  })

  describe('structure', () => {
    it('should have header element', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('header').exists()).toBe(true)
    })

    it('should have main element', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('main').exists()).toBe(true)
    })

    it('should have footer element', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('footer').exists()).toBe(true)
    })
  })

  describe('layout', () => {
    it('should have min-h-screen class on root', () => {
      const wrapper = mountWithRouter()
      expect(wrapper.find('[data-testid="app"]').classes()).toContain('min-h-screen')
    })

    it('should have proper header styling', () => {
      const wrapper = mountWithRouter()
      const header = wrapper.find('header')
      expect(header.classes()).toContain('bg-gradient-to-r')
    })

    it('should have proper footer styling', () => {
      const wrapper = mountWithRouter()
      const footer = wrapper.find('footer')
      expect(footer.classes()).toContain('bg-dark-card')
    })
  })
})
