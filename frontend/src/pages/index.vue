<template>
  <v-container>
    <v-infinite-scroll side="both" @load="load" v-if="items.length > 0">
      <template v-for="row in (items.length - items.length % 12) / 12 " :key="row">
        <v-row>
          <v-col v-for="col in 12" :key="col"
            lg="1"
            md="2"
            sm="3"
          >
            <v-btn class="mx-0" @click="conquer(items[(row-1)*12+(col-1)].id, (row-1)*12+(col-1))"
              :color="items[(row - 1) * 12 + (col - 1)]?.[currentType] === username ? 'primary' : ''"
            >
              {{ items[(row - 1) * 12 + (col - 1)]?.[currentType] }}
            </v-btn>
          </v-col>
        </v-row>
      </template>
    </v-infinite-scroll>
  </v-container>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import axios from 'axios'

const batchSize = 1000
const maxItems = 1000000

const router = useRouter()
const appStore = useAppStore()
const { currentType, username } = storeToRefs(appStore)
const items = ref([])
const polling = ref(null)
const startPivot = ref(1)
const endPivot = ref(batchSize)

onMounted(async () => {
  await appStore.verifyToken()
  if(!appStore.isLoggedIn) {
    router.push('/login')
  }

  let data = await loadData(startPivot.value, endPivot.value)
  items.value = data

  startPolling()
})

onUnmounted(() => {
  clearTimeout(polling.value)
})

const startPolling = () => {
  polling.value = setTimeout(async () => {
    if (items.value.length == 0) {
      // no data, load data from startPivot to endPivot
      let data = await loadData(startPivot.value, endPivot.value)
      items.value = data
    } else {
      // reload data from currentStart to currentEnd
      let currentStart = parseInt(items.value[0].id, 10)
      let currentEnd = parseInt(items.value[items.value.length - 1].id, 10)
      let data = await loadData(currentStart, currentEnd)
      // replace data one by one to prevent flickering
      for (let i = 0; i < data.length; i++) {
        const idx = items.value.findIndex((v) => parseInt(data[i].id, 10) === parseInt(v.id, 10))
        if (idx !== -1) {
          items.value[idx] = data[i]
        }
      }
    }
    startPolling()
  }, 10000)
}

const load = ({ side, done }) => {
  setTimeout(async () => {
    if (side === 'start') {
      if (startPivot.value === 1) {
        // no more data
        done('empty')
        return
      }

      // move pivot
      if (startPivot.value - batchSize < 0) {
        startPivot.value = 1
        endPivot.value = batchSize
      } else {
        startPivot.value -= batchSize
        endPivot.value -= batchSize
      }

      // get currentStart
      const currentStart = Object.keys(items.value[0])[0]

      // check if currentStart is within the range of pivot, we only need to load data to fill the gap
      if (currentStart < endPivot.value) {
        let data = await loadData(startPivot.value, currentStart - 1)
        items.value = [...data, ...items.value]
      } else {
        // load data from startPivot to endPivot
        let data = await loadData(startPivot.value, endPivot.value)
        items.value = [...data, ...items.value]
      }

      // remove last data if it's more than 3 * batchSize
      items.value = items.value.slice(0, 3 * batchSize)
    } else if (side === 'end') {
      if (endPivot.value === maxItems) {
        // no more data
        done('empty')
        return
      }

      // move pivot
      if (endPivot.value + batchSize > maxItems) {
        startPivot.value = maxItems - batchSize
        endPivot.value = maxItems
      } else {
        startPivot.value += batchSize
        endPivot.value += batchSize
      }

      // get currentEnd
      const currentEnd = parseInt(Object.keys(items.value[items.value.length - 1])[0])

      // check if currentEnd is within the range of pivot, we only need to load data to fill the gap
      if (currentEnd < endPivot.value) {
        let data = await loadData(currentEnd + 1, endPivot.value)
        items.value = [...items.value, ...data]
      } else {
        // load data from startPivot to endPivot
        let data = await loadData(startPivot.value, endPivot.value)
        items.value = [...items.value, ...data]
      }

      // remove first data if it's more than 3 * batchSize
      items.value = items.value.slice(-3 * batchSize)
    }
    done('ok')
  }, 1000)
}

const loadData = async (start, end) => {
  // data: { id: { webservice: '', restful: '', graphql: '', grpc: ''}}

  let res = await axios.get('/map', {
    params: {
      start: start,
      end: end
    }
  })

  // mock data
  // const data = {
  //   data: Array.from({ length: end - start }, (v, index) => {
  //     index += start
  //     // each item has 10% chance to be random
  //     if (Math.random() > 0.1) {
  //       return {
  //         [index]: {
  //           webservice: "AAAAAA",
  //           restful: "BBBBBB",
  //           graphql: "CCCCCC",
  //           grpc: "DDDDDD"
  //         }
  //       }
  //     } else {
  //       let randomString = Math.random().toString(36).substring(7)
  //       return {
  //         [index]: {
  //           webservice: randomString,
  //           restful: randomString,
  //           graphql: randomString,
  //           grpc: randomString
  //         }
  //       }
  //     }
  //   })
  // }

  // convert data to array for easier manipulation
  let data = res.data
  let result = Object.keys(res.data).map((v) => {
    return {
      id: parseInt(v, 10),
      [currentType.value]: data[v][currentType.value]
    }
  })
  return result
}

const conquer = async (fieldid, itemIndex) => {
  await appStore.conquer(fieldid)
  // display the result immediately
  items.value[itemIndex][currentType.value] = username.value
}
</script>
