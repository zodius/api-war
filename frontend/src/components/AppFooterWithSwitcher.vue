<template>
  <v-footer height="40" app>
    <div class="pt-0">
      switch to
    </div>
    <v-btn v-for="icon in icons" :key="icon.label" :prepend-icon="icon.icon" variant="tonal" class="mx-4"
      @click="appStore.setMode(icon.label)">{{
      icon.label }}</v-btn>

    <div class="text-caption text-disabled" style="position: absolute; right: 16px;">
      Current mode: {{ currentType }}
      <v-spacer></v-spacer>
      &copy; 2024 <span class="d-none d-sm-inline-block">zodius</span>
    </div>
  </v-footer>
</template>

<script setup>
import { storeToRefs } from 'pinia'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

const { currentType } = storeToRefs(appStore)

const icons = [
  { icon: 'mdi-cog-outline', label: 'webservice' },
  { icon: 'mdi-book', label: 'restful' },
  { icon: 'mdi-database-search', label: 'graphql' },
]
</script>