<script setup>
import { useRouter } from 'vue-router'
import { computed, onMounted, ref, watch } from 'vue'
import { useLazyQuery, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import PageBar from '@/views/components/PageBar.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import moment from 'moment'

const router = useRouter()
const serverId = router.currentRoute.value.query.id
const serverName = router.currentRoute.value.query.name

const serverLogsContent = ref('')
const logId = ref(0)
const serverLogsContentRef = ref(null)

const {
  loading: serverLogsLoading,
  refetch: refetchServerLogs,
  result: serverLogsResultRaw,
  onError: onServerLogsError
} = useQuery(
  gql`
    query ($serverId: Uint!) {
      server(id: $serverId) {
        logs {
          id
          title
          updatedAt
        }
      }
    }
  `,
  {
    serverId
  }
)

onServerLogsError((error) => {
  toast.error(error.message)
})

const serverLogsResult = computed(() => {
  if (serverLogsResultRaw.value && serverLogsResultRaw.value.server) {
    return serverLogsResultRaw.value.server.logs
  }
  return []
})

watch(serverLogsResult, (logs) => {
  if (logs.length > 0) {
    if (logId.value === 0) selectServerLog(logs[0])
  }
})

const {
  load: fetchServerLogsContent,
  refetch: refetchServerLogsContent,
  loading: serverLogsContentLoading,
  onResult: onServerLogsContentResult,
  onError: onServerLogsContentError,
  variables: serverLogsContentVariables
} = useLazyQuery(gql`
  query GetServerLogsContent($logId: Uint!) {
    fetchServerLogContent(id: $logId)
  }
`)

onServerLogsContentResult((result) => {
  if (result.data && result.data.fetchServerLogContent) {
    serverLogsContent.value = result.data.fetchServerLogContent
    setTimeout(() => {
      if (serverLogsContentRef.value === null) {
        return
      }
      serverLogsContentRef.value.scrollTop = serverLogsContentRef.value.scrollHeight
    }, 1000)
  }
})

onServerLogsContentError((error) => {
  toast.error(error.message)
})

function loadServerLogsContent() {
  serverLogsContentVariables.value = { logId: logId.value }
  if (fetchServerLogsContent() === false) {
    refetchServerLogsContent()
  }
}

function selectServerLog(log) {
  logId.value = log.id
  loadServerLogsContent()
}

onMounted(() => {
  if (!serverId) {
    alert('Try to re-open the page from the server list')
  }
  if (serverName) {
    document.title = `${serverName} - Server Logs`
  }
})
</script>

<template>
  <PageBar>
    <template v-slot:title>{{ serverName }}</template>
    <template v-slot:subtitle>These logs are related to the actions performed on the server</template>
    <template v-slot:buttons>
      <FilledButton type="ghost" :click="refetchServerLogs">
        <font-awesome-icon
          icon="fa-solid fa-arrows-rotate"
          :class="{
            'animate-spin ': serverLogsLoading
          }" />&nbsp;&nbsp;Refresh List
      </FilledButton>
    </template>
  </PageBar>

  <div class="mt-8 flex w-full gap-4">
    <!--  Server logs list  -->
    <div
      class="scrollbox flex max-h-[80vh] w-[400px] flex-col gap-2 overflow-y-auto pr-2"
      v-if="serverLogsResult.length > 0">
      <div
        @click="() => selectServerLog(log)"
        :key="log.id"
        v-for="log in serverLogsResult"
        class="w-full cursor-pointer select-none rounded-lg border-2 border-secondary-200 p-3 hover:bg-secondary-200"
        :class="{
          'border-secondary-400 bg-secondary-200': log.id === logId
        }">
        <p class="font-medium">{{ log.title }}</p>
        <p>{{ moment(new Date(log.updatedAt)).format('Do MMMM YYYY - h:mm:ss a') }}</p>
      </div>
    </div>
    <div v-else class="max-h-[80vh] w-[400px]">
      <i v-if="serverLogsLoading">Loading logs. Please wait...</i>
      <i v-else>No logs found. Check later or <b>refresh list</b></i>
    </div>
    <!--  Server log result  -->
    <div
      class="relative max-h-[80vh] w-full overflow-y-hidden whitespace-pre-wrap rounded-lg border-2 border-secondary-200 bg-secondary-100 p-4">
      <FilledButton type="secondary" :click="loadServerLogsContent" class="absolute right-2 top-2">
        <font-awesome-icon
          icon="fa-solid fa-arrows-rotate"
          :class="{
            'animate-spin ': serverLogsContentLoading
          }" />&nbsp;&nbsp;Refresh Logs
      </FilledButton>
      <div v-if="logId === 0">
        <i>Select a record from left side to view content</i>
      </div>
      <div v-else-if="serverLogsContent.length === 0">
        <i>No Logs found. Check later or refresh logs</i>
      </div>
      <div v-else-if="serverLogsContentLoading">
        <i>Loading logs. Please wait...</i>
      </div>
      <div v-else class="scrollbox h-full overflow-y-auto scroll-smooth" ref="serverLogsContentRef">
        {{ serverLogsContent }}
      </div>
    </div>
  </div>
</template>

<style scoped></style>
