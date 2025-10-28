export interface Festival {
  id: string
  name: string
  description: string
  startDate: string
  endDate: string
  city: string
  region: string
  location: {
    latitude: number
    longitude: number
  }
  image?: string
  website?: string
  breweryCount?: number
}
