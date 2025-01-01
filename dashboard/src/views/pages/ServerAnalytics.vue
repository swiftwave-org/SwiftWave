<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import PageBar from '@/views/components/PageBar.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { useLazyQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import AreaChartTimeSeries from '@/views/components/AreaChartTimeSeries.vue'
import { humanizeDiskGB, humanizeMemoryGB, humanizeNetworkSpeed } from '@/vendor/utils.js'

const router = useRouter()
const serverInfo = reactive({
  id: router.currentRoute.value.query.id,
  hostname: '',
  ip: ''
})

const statsTimeframe = ref('last_1_hour')
const cpuUsageSeries = ref([])
const memoryUsageSeries = ref([])
const networkUsageSeries = ref([])
const diskUsageSeries = ref([])

// Fetch server info
const {
  load: loadServerInfo,
  loading: isServerInfoLoading,
  onResult: onServerInfoResult,
  onError: onServerInfoError
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      server(id: $id) {
        id
        hostname
        ip
      }
    }
  `,
  {
    id: serverInfo.id
  }
)

onServerInfoResult((result) => {
  serverInfo.id = result.data.server.id
  serverInfo.hostname = result.data.server.hostname
  serverInfo.ip = result.data.server.ip
  fetchAllAnalytics()
})

onServerInfoError((error) => {
  toast.error(error)
})

onMounted(() => {
  if (!serverInfo.id) {
    router.push({ name: 'Servers' })
  } else {
    loadServerInfo()
  }
})

// Fetch server analytics
const {
  load: loadServerAnalyticsRaw,
  refetch: refetchServerAnalyticsRaw,
  loading: isServerAnalyticsLoading,
  onResult: onServerAnalyticsResult,
  onError: onServerAnalyticsError
} = useLazyQuery(
  gql`
    query ($id: Uint!, $timeframe: ServerResourceAnalyticsTimeframe!) {
      serverResourceAnalytics(id: $id, timeframe: $timeframe) {
        cpu_usage_percent
        memory_total_gb
        memory_used_gb
        memory_cached_gb
        network_sent_kbps
        network_recv_kbps
        timestamp
      }
    }
  `,
  {
    id: serverInfo.id,
    timeframe: statsTimeframe
  }
)
onServerAnalyticsError((error) => {
  toast.error(error)
})

onServerAnalyticsResult((result) => {
  if (result.data.serverResourceAnalytics.length !== 0) {
    let cpuUsagePercentStat = []
    let memoryTotalGbStat = []
    let memoryUsedGbStat = []
    let memoryCachedGbStat = []
    let networkSentKbpsStat = []
    let networkRecvKbpsStat = []
    result.data.serverResourceAnalytics.forEach((d) => {
      cpuUsagePercentStat.push([new Date(d.timestamp).getTime(), d.cpu_usage_percent])
      memoryTotalGbStat.push([new Date(d.timestamp).getTime(), d.memory_total_gb])
      memoryUsedGbStat.push([new Date(d.timestamp).getTime(), d.memory_used_gb])
      memoryCachedGbStat.push([new Date(d.timestamp).getTime(), d.memory_cached_gb])
      networkSentKbpsStat.push([new Date(d.timestamp).getTime(), d.network_sent_kbps])
      networkRecvKbpsStat.push([new Date(d.timestamp).getTime(), d.network_recv_kbps])
    })
    cpuUsageSeries.value = [
      {
        name: 'CPU Usage',
        data: cpuUsagePercentStat
      }
    ]
    memoryUsageSeries.value = [
      {
        name: 'Total Memory',
        data: memoryTotalGbStat
      },
      {
        name: 'Used Memory',
        data: memoryUsedGbStat
      },
      {
        name: 'Cached Memory',
        data: memoryCachedGbStat
      }
    ]
    networkUsageSeries.value = [
      {
        name: 'Sent',
        data: networkSentKbpsStat
      },
      {
        name: 'Received',
        data: networkRecvKbpsStat
      }
    ]
  }
})

const loadServerAnalytics = () => {
  if (loadServerAnalyticsRaw() === false) {
    refetchServerAnalyticsRaw()
  }
}

// Fetch server disk usage
const {
  load: loadServerDiskUsageRaw,
  refetch: refetchServerDiskUsageRaw,
  loading: isServerDiskUsageLoading,
  onResult: onServerDiskUsageResult,
  onError: onServerDiskUsageError
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      serverDiskUsage(id: $id) {
        disks {
          path
          total_gb
          used_gb
          timestamp
        }
      }
    }
  `,
  {
    id: serverInfo.id
  }
)

onServerDiskUsageResult((result) => {
  if (!result.data.serverDiskUsage) {
    return
  }
  let diskUsageStatMap = {}
  result.data.serverDiskUsage.forEach((record) => {
    record.disks.forEach((d) => {
      if (!diskUsageStatMap[d.path]) {
        diskUsageStatMap[d.path] = {
          total_gb: [],
          used_gb: []
        }
      }
      diskUsageStatMap[d.path].total_gb.push([new Date(d.timestamp).getTime(), d.total_gb])
      diskUsageStatMap[d.path].used_gb.push([new Date(d.timestamp).getTime(), d.used_gb])
    })
  })
  let diskUsageStatTemp = {}
  for (const [key, value] of Object.entries(diskUsageStatMap)) {
    diskUsageStatTemp[key] = [
      {
        name: 'Total',
        data: value.total_gb
      },
      {
        name: 'Used',
        data: value.used_gb
      }
    ]
  }
  diskUsageSeries.value = diskUsageStatTemp
  diskUsageStatTemp = null
  diskUsageStatMap = null
})

onServerDiskUsageError((error) => {
  toast.error(error)
})

const loadServerDiskUsage = () => {
  if (loadServerDiskUsageRaw() === false) {
    refetchServerDiskUsageRaw()
  }
}

// Combined
const fetchAllAnalytics = () => {
  loadServerAnalytics()
  loadServerDiskUsage()
}
</script>

<template>
  <section class="mx-auto max-h-[100vh] w-full max-w-7xl overflow-y-hidden">
    <div v-if="isServerInfoLoading" class="w-full font-medium italic">Fetching server info...</div>
    <div v-else class="flex h-full w-full flex-col">
      <!-- Top Page bar   -->
      <PageBar>
        <template v-slot:title>Server {{ serverInfo.hostname }} ({{ serverInfo.ip }})</template>
        <template v-slot:subtitle>Monitor the resource analytics of your server</template>
        <template v-slot:buttons>
          <select
            class="block rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            v-model="statsTimeframe"
            @change="loadServerAnalytics">
            <option value="last_1_hour">Last 1 hour</option>
            <option value="last_3_hours">Last 3 hours</option>
            <option value="last_6_hours">Last 6 hours</option>
            <option value="last_12_hours">Last 12 hours</option>
            <option value="last_24_hours">Last 24 hours</option>
            <option value="last_7_days">Last 7 days</option>
            <option value="last_30_days">Last 30 days</option>
          </select>
          <FilledButton
            type="primary"
            :click="fetchAllAnalytics"
            :loading="isServerAnalyticsLoading || isServerDiskUsageLoading">
            Refresh Stats
          </FilledButton>
        </template>
      </PageBar>
      <div class="scrollbox mt-5 flex max-h-full w-full flex-col gap-5 overflow-y-auto overflow-x-clip pr-2">
        <!--  Cpu usage series  -->
        <AreaChartTimeSeries
          title="CPU Usage"
          :series="cpuUsageSeries"
          :y-axis-formatter="
            (val) => {
              return val + ' %'
            }
          "
          :y-axis-minimum-max="100" />
        <!--  Memory usage series  -->
        <AreaChartTimeSeries title="Memory Usage" :series="memoryUsageSeries" :y-axis-formatter="humanizeMemoryGB" />
        <!--  Network usage series  -->
        <AreaChartTimeSeries
          title="Network Usage"
          :series="networkUsageSeries"
          :y-axis-formatter="humanizeNetworkSpeed" />
        <!--    Disk usage    -->
        <div class="w-full">
          <p class="text-base font-medium">Disk Usage (Last 24 hours)</p>
          <div class="mt-2 grid w-full grid-cols-2 gap-4">
            <div v-for="(diskUsage, diskPath) in diskUsageSeries" :key="diskPath" class="w-full">
              <AreaChartTimeSeries
                :title="diskPath"
                :series="diskUsage"
                :toolbar="false"
                :y-axis-formatter="humanizeDiskGB" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
