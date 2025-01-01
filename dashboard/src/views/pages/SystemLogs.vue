<script setup>
import { computed, ref, watch } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import PageBar from '@/views/components/PageBar.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import moment from 'moment'
import axios from 'axios'
import { getHttpBaseUrl } from '@/vendor/utils.js'

const systemLogsContent = ref('')
const logFileName = ref('')
const systemLogsContentRef = ref(null)
const systemLogsContentLoading = ref(false)

const {
  loading: systemLogsLoading,
  refetch: refetchSystemLogs,
  result: systemLogsResultRaw,
  onError: onSystemLogsError
} = useQuery(gql`
  query {
    fetchSystemLogRecords {
      name
      modTime
    }
  }
`)

onSystemLogsError((error) => {
  toast.error(error.message)
})

const systemLogsResult = computed(() => {
  if (systemLogsResultRaw.value && systemLogsResultRaw.value.fetchSystemLogRecords) {
    return systemLogsResultRaw.value.fetchSystemLogRecords
  }
  return []
})

watch(systemLogsResult, (logs) => {
  if (logs.length > 0) {
    if (logFileName.value === '') selectSystemLog(logs[0])
  }
})

async function loadServerLogsContent() {
  if (logFileName.value === '') return
  systemLogsContentLoading.value = true
  let config = {
    method: 'get',
    url: `${getHttpBaseUrl()}/log/${logFileName.value}`,
    headers: {
      Authorization: 'Bearer ' + (localStorage.getItem('token') || '')
    }
  }
  try {
    let res = await axios.request(config)
    systemLogsContentLoading.value = false
    systemLogsContent.value = res.data
    setTimeout(() => {
      if (systemLogsContentRef.value === null) {
        return
      }
      systemLogsContentRef.value.scrollTop = systemLogsContentRef.value.scrollHeight
    }, 500)
  } catch (e) {
    systemLogsContentLoading.value = false
    if (e.response) {
      toast.error(e.response.data.message || 'Unexpected error')
    } else {
      toast.error('Failed to send request')
    }
  }
}

function selectSystemLog(log) {
  logFileName.value = log.name
  loadServerLogsContent()
}
</script>

<template>
  <PageBar>
    <template v-slot:title>System Logs</template>
    <template v-slot:subtitle>It contains all the system logs and error logs</template>
    <template v-slot:buttons>
      <FilledButton type="ghost" :click="refetchSystemLogs">
        <font-awesome-icon
          icon="fa-solid fa-arrows-rotate"
          :class="{
            'animate-spin ': systemLogsLoading
          }" />&nbsp;&nbsp;Refresh List
      </FilledButton>
    </template>
  </PageBar>

  <div class="mt-8 flex w-full gap-4">
    <!--  System logs list  -->
    <div
      class="scrollbox flex max-h-[80vh] w-[400px] flex-col gap-2 overflow-y-auto pr-2"
      v-if="systemLogsResult.length > 0">
      <div
        @click="() => selectSystemLog(log)"
        :key="log.id"
        v-for="log in systemLogsResult"
        class="w-full cursor-pointer select-none rounded-lg border-2 border-secondary-200 p-3 hover:bg-secondary-200"
        :class="{
          'border-secondary-400 bg-secondary-200': log.name === logFileName
        }">
        <p class="font-medium">{{ log.name }}</p>
        <p>{{ moment(new Date(log.modTime)).format('Do MMMM YYYY - h:mm:ss a') }}</p>
      </div>
    </div>
    <div v-else class="max-h-[80vh] w-[400px]">
      <i v-if="systemLogsLoading">Loading logs. Please wait...</i>
      <i v-else>No logs found. Check later or <b>refresh list</b></i>
    </div>
    <!--  Server log result  -->
    <div
      class="relative max-h-[80vh] w-full overflow-y-hidden whitespace-pre-wrap rounded-lg border-2 border-secondary-200 bg-secondary-100 p-4">
      <FilledButton type="secondary" :click="loadServerLogsContent" class="absolute right-2 top-2">
        <font-awesome-icon
          icon="fa-solid fa-arrows-rotate"
          :class="{
            'animate-spin ': systemLogsContentLoading
          }" />&nbsp;&nbsp;Refresh Logs
      </FilledButton>
      <div v-if="logFileName === ''">
        <i>Select a record from left side to view content</i>
      </div>
      <div v-else-if="systemLogsContent.length === 0">
        <i>No Logs found.</i>
      </div>
      <div v-else-if="systemLogsContentLoading">
        <i>Loading logs. Please wait...</i>
      </div>
      <div v-else class="scrollbox h-full overflow-y-auto scroll-smooth" ref="systemLogsContentRef">
        {{ systemLogsContent }}
      </div>
    </div>
  </div>
</template>

<style scoped></style>
