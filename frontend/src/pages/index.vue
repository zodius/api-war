<template>
  <v-btn v-for="(item, id) in map" :key="id">
    {{ item.restful.substring(0, 6) }}
  </v-btn>
</template>

<script setup>
import axios from 'axios';
import { onMounted, onUnmounted, ref} from 'vue';

const map = ref(null);
const polling = ref(null);

onMounted(async () => {
  setPolling();
});

onUnmounted(() => {
  clearTimeout(polling.value);
});

const setPolling = () => {
  getMap();
  polling.value = setTimeout(async () => {
    setPolling();
  }, 10000);
}

const getMap = async () => {
  let res = await axios.get('http://localhost:8971/map');
  map.value = res.data;
  console.log(map.value);
}

</script>
