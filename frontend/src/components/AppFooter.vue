<template>
  <v-footer height="40" app>
    <template slot="content"></template>

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

</script>