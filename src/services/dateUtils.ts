import dayjs from 'dayjs'
import 'dayjs/locale/fr'

dayjs.locale('fr')

export const formatDate = (date: string, format = 'DD MMMM YYYY'): string => {
  return dayjs(date).format(format)
}

export const formatDateRange = (startDate: string, endDate: string): string => {
  const start = dayjs(startDate)
  const end = dayjs(endDate)

  if (start.isSame(end, 'day')) {
    return formatDate(startDate)
  }

  if (start.isSame(end, 'month')) {
    return `${start.format('DD')} - ${end.format('DD MMMM YYYY')}`
  }

  return `${start.format('DD MMMM')} - ${end.format('DD MMMM YYYY')}`
}

export const isUpcoming = (date: string): boolean => {
  return dayjs(date).isAfter(dayjs())
}

export const isPast = (date: string): boolean => {
  return dayjs(date).isBefore(dayjs())
}

export const daysUntil = (date: string): number => {
  return dayjs(date).diff(dayjs(), 'day')
}

export const getDaysUntilText = (date: string): string => {
  const days = daysUntil(date)

  if (days < 0) {
    return 'Terminé'
  }

  if (days === 0) {
    return "Aujourd'hui"
  }

  if (days === 1) {
    return 'Demain'
  }

  if (days < 7) {
    return `Dans ${days} jours`
  }

  if (days < 30) {
    const weeks = Math.floor(days / 7)
    return `Dans ${weeks} semaine${weeks > 1 ? 's' : ''}`
  }

  const months = Math.floor(days / 30)
  return `Dans ${months} mois`
}
