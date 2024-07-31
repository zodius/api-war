<template>
  <v-container>
    <v-table>
      <thead>
        <tr>
          <th>Rank</th>
          <th>Username</th>
          <th>Score</th>
          <th>Restful</th>
          <th>Graphql</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(score, index) in scorelist" :key="index">
          <td>{{ index + 1 }}</td>
          <td>{{ score?.username }}</td>
          <td>{{ score?.conquerFieldCount }}</td>
          <td>{{ score?.conquerHistoryCount?.restful }}</td>
          <td>{{ score?.conquerHistoryCount?.graphql }}</td>
        </tr>
      </tbody>
    </v-table>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'

const polling = ref(null)
const scorelist = ref([])

onMounted(() => {
  loadData()
  startPolling()
})

onUnmounted(() => {
  clearTimeout(polling.value)
})

const startPolling = () => {
  polling.value = setTimeout(async () => {
    loadData()
    startPolling()
  }, 5000)
}

const loadData = async () => {
  let res = await axios.get('/scoreboard')

  // display first 100 records
  let data = res.data
  data = Array.from(data.scoreList)
  data = data.sort((a, b) => b.conquerFieldCount - a.conquerFieldCount)
  data = data.slice(0, 100)
  scorelist.value = data
}

</script>