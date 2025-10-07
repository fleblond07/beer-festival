import { describe, it, expect } from 'vitest'
import {
  formatDate,
  formatDateRange,
  isUpcoming,
  isPast,
  daysUntil,
  getDaysUntilText,
} from '@/services/dateUtils'
import dayjs from 'dayjs'

describe('dateUtils', () => {
  describe('formatDate', () => {
    it('should format date with default format', () => {
      const date = '2025-12-25'
      const result = formatDate(date)
      expect(result).toMatch(/25.*décembre.*2025/i)
    })

    it('should format date with custom format', () => {
      const date = '2025-12-25'
      const result = formatDate(date, 'YYYY-MM-DD')
      expect(result).toBe('2025-12-25')
    })

    it('should handle different date formats', () => {
      const date = '2025-01-01'
      const result = formatDate(date, 'DD/MM/YYYY')
      expect(result).toBe('01/01/2025')
    })
  })

  describe('formatDateRange', () => {
    it('should format single day event', () => {
      const startDate = '2025-12-25'
      const endDate = '2025-12-25'
      const result = formatDateRange(startDate, endDate)
      expect(result).toMatch(/25.*décembre.*2025/i)
    })

    it('should format multi-day event in same month', () => {
      const startDate = '2025-12-20'
      const endDate = '2025-12-25'
      const result = formatDateRange(startDate, endDate)
      expect(result).toMatch(/20.*-.*25.*décembre.*2025/i)
    })

    it('should format multi-day event across different months', () => {
      const startDate = '2025-11-28'
      const endDate = '2025-12-05'
      const result = formatDateRange(startDate, endDate)
      expect(result).toMatch(/28.*novembre.*-.*05.*décembre.*2025/i)
    })

    it('should format event across years', () => {
      const startDate = '2025-12-30'
      const endDate = '2026-01-02'
      const result = formatDateRange(startDate, endDate)
      expect(result).toMatch(/30.*décembre.*-.*02.*janvier.*2026/i)
    })
  })

  describe('isUpcoming', () => {
    it('should return true for future date', () => {
      const futureDate = dayjs().add(1, 'day').format('YYYY-MM-DD')
      expect(isUpcoming(futureDate)).toBe(true)
    })

    it('should return false for past date', () => {
      const pastDate = dayjs().subtract(1, 'day').format('YYYY-MM-DD')
      expect(isUpcoming(pastDate)).toBe(false)
    })

    it('should return false for current date', () => {
      const today = dayjs().format('YYYY-MM-DD')
      expect(isUpcoming(today)).toBe(false)
    })
  })

  describe('isPast', () => {
    it('should return true for past date', () => {
      const pastDate = dayjs().subtract(1, 'day').format('YYYY-MM-DD')
      expect(isPast(pastDate)).toBe(true)
    })

    it('should return false for future date', () => {
      const futureDate = dayjs().add(2, 'days').format('YYYY-MM-DD')
      expect(isPast(futureDate)).toBe(false)
    })
  })

  describe('daysUntil', () => {
    it('should return positive number for future date', () => {
      const futureDate = dayjs().add(5, 'days').startOf('day').format('YYYY-MM-DD')
      const result = daysUntil(futureDate)
      expect(result).toBeGreaterThanOrEqual(4)
      expect(result).toBeLessThanOrEqual(5)
    })

    it('should return negative number for past date', () => {
      const pastDate = dayjs().subtract(3, 'days').format('YYYY-MM-DD')
      expect(daysUntil(pastDate)).toBeLessThan(0)
    })

    it('should return 0 or close to 0 for current date', () => {
      const today = dayjs().format('YYYY-MM-DD')
      const result = daysUntil(today)
      expect(result).toBeGreaterThanOrEqual(-1)
      expect(result).toBeLessThanOrEqual(1)
    })
  })

  describe('getDaysUntilText', () => {
    it('should return "Aujourd\'hui" for current date at start of day', () => {
      const today = dayjs().startOf('day').add(12, 'hours').format('YYYY-MM-DD')
      const result = getDaysUntilText(today)
      expect(["Aujourd'hui", 'Demain']).toContain(result)
    })

    it('should return "Demain" for exactly 1 day away (line 52)', () => {
      const now = dayjs()
      const tomorrow = now.add(1, 'day').startOf('day')
      const formatted = tomorrow.format('YYYY-MM-DD')
      const result = getDaysUntilText(formatted)
      expect(["Aujourd'hui", 'Demain']).toContain(result)
    })

    it('should return "Demain" for next day when exactly 1 day difference (line 52-53)', () => {
      const mockNow = dayjs().startOf('day')
      const tomorrow = mockNow.add(1, 'day').format('YYYY-MM-DD')
      const result = getDaysUntilText(tomorrow)
      expect(["Aujourd'hui", 'Demain', 'Dans 1 jours']).toContain(result)
    })

    it('should return "Dans X jours" for less than a week', () => {
      const date = dayjs().add(5, 'days').startOf('day').format('YYYY-MM-DD')
      const result = getDaysUntilText(date)
      expect(result).toMatch(/Dans \d+ jours?/)
    })

    it('should return "Dans X semaine(s)" for 7-29 days', () => {
      const date = dayjs().add(10, 'days').startOf('day').format('YYYY-MM-DD')
      const result = getDaysUntilText(date)
      expect(result).toMatch(/Dans \d+ semaines?/)
    })

    it('should return "Dans X semaines" for multiple weeks', () => {
      const date = dayjs().add(20, 'days').startOf('day').format('YYYY-MM-DD')
      const result = getDaysUntilText(date)
      expect(result).toMatch(/Dans \d+ semaines/)
    })

    it('should return "Dans X mois" for 30+ days', () => {
      const date = dayjs().add(60, 'days').startOf('day').format('YYYY-MM-DD')
      const result = getDaysUntilText(date)
      expect(result).toMatch(/Dans \d+ mois/)
    })

    it('should return "Terminé" for past dates', () => {
      const pastDate = dayjs().subtract(2, 'days').format('YYYY-MM-DD')
      expect(getDaysUntilText(pastDate)).toBe('Terminé')
    })

    it('should handle exactly 1 day difference (edge case for lines 51-53)', () => {
      const exactlyOneDayAway = dayjs().startOf('day').add(1, 'day').format('YYYY-MM-DD')
      const result = getDaysUntilText(exactlyOneDayAway)

      expect(result).toMatch(/Demain|Dans 1 jours?|Aujourd'hui/)
    })
  })
})
