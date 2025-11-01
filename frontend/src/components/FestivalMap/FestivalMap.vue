<template>
  <div class="festival-map-container w-full py-16 px-4" data-testid="festival-map-container">
    <div class="max-w-6xl mx-auto">
      <h2 class="text-center text-4xl font-bold text-white mb-2" data-testid="section-title">
        Carte des Festivals
      </h2>
      <p class="text-center text-gray-400 mb-8" data-testid="section-subtitle">
        Explorez les festivals de bière partout en France
      </p>

      <div
        class="map-wrapper bg-dark-card rounded-xl shadow-2xl overflow-hidden border border-dark-lighter"
        style="height: 600px"
        data-testid="map-wrapper"
      >
        <div ref="mapContainer" class="w-full h-full" data-testid="map-container"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import type { Festival } from '@/types'

interface Props {
  festivals: Festival[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (_e: 'marker-click', _festival: Festival): void
  (_e: 'popup-click', _festival: Festival): void
}>()

const mapContainer = ref<HTMLElement | null>(null)
let map: L.Map | null = null
const markers: L.Marker[] = []

const initMap = () => {
  if (!mapContainer.value) return

  map = L.map(mapContainer.value).setView([46.603354, 1.888334], 6)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© OpenStreetMap contributors',
    maxZoom: 19,
  }).addTo(map)

  addMarkers()
}

const addMarkers = () => {
  if (!map) return

  markers.forEach(marker => marker.remove())
  markers.length = 0

  props.festivals.forEach(festival => {
    const marker = L.marker([festival.location.latitude, festival.location.longitude], {
      icon: createCustomIcon(),
    }).addTo(map!)

    marker.bindPopup(`
      <div class="festival-popup">
        <h3 style="font-weight: bold; margin-bottom: 8px;">${festival.name}</h3>
        <p style="margin-bottom: 4px;">${festival.city}</p>
        ${festival.breweryCount ? `<p style="color: #00d2d3; margin-bottom: 8px;">${festival.breweryCount} brasseries</p>` : ''}
        <button
          class="popup-scroll-btn"
          data-festival-id="${festival.id}"
          style="background: #00d2d3; color: #1a1a2e; padding: 8px 16px; border: none; border-radius: 6px; font-weight: 600; cursor: pointer; width: 100%; margin-top: 8px; transition: background 0.2s;"
          onmouseover="this.style.background='#00b8b9'"
          onmouseout="this.style.background='#00d2d3'"
        >
          Voir dans la liste ↓
        </button>
      </div>
    `)

    marker.on('click', () => {
      emit('marker-click', festival)
    })

    marker.on('popupopen', () => {
      const popupElement = marker.getPopup()?.getElement()
      const button = popupElement?.querySelector('.popup-scroll-btn')
      if (button) {
        button.addEventListener('click', () => {
          emit('popup-click', festival)
        })
      }
    })

    markers.push(marker)
  })

  if (markers.length > 0) {
    const group = L.featureGroup(markers)
    map?.fitBounds(group.getBounds().pad(0.1))
  }
}

const createCustomIcon = () => {
  return L.icon({
    iconUrl:
      'data:image/svg+xml;base64,' +
      btoa(`
      <svg width="32" height="32" viewBox="0 0 32 32" xmlns="http://www.w3.org/2000/svg">
        <circle cx="16" cy="16" r="14" fill="#00d2d3" stroke="#ffffff" stroke-width="3"/>
        <path d="M16 10 L16 22 M10 16 L22 16" stroke="#ffffff" stroke-width="2" stroke-linecap="round"/>
      </svg>
    `),
    iconSize: [32, 32],
    iconAnchor: [16, 16],
    popupAnchor: [0, -16],
  })
}

onMounted(() => {
  initMap()
})

onBeforeUnmount(() => {
  markers.forEach(marker => marker.remove())

  if (map) {
    map.remove()
    map = null
  }
})

watch(
  () => props.festivals,
  () => {
    if (map) {
      addMarkers()
    }
  }
)
</script>

<style scoped>
:deep(.leaflet-container) {
  font-family: inherit;
}

:deep(.leaflet-popup-content-wrapper) {
  border-radius: 8px;
}

:deep(.leaflet-popup-content) {
  margin: 12px;
}
</style>
