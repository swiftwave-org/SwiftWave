<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth.js'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import FilledButton from '@/views/components/FilledButton.vue'

const router = useRouter()
const authStore = useAuthStore()
const pvBackupId = router.currentRoute.value.params.backup_id
const downloadUrl = getHttpBaseUrl() + `/persistent-volume/backup/${pvBackupId}/download`
const fileNameUrl = getHttpBaseUrl() + `/persistent-volume/backup/${pvBackupId}/filename`

const data = reactive({
  filenameFetched: false,
  filename: '',
  fileSizeMb: 0,
  progressPercentage: 0,
  mbps: 0,
  lastThreeMbps: [],
  fileDownloading: false,
  fileDownloaded: false,
  isErrored: false
})

const isStartDownloadButtonShown = computed(() => {
  return !data.fileDownloaded && !data.fileDownloading
})

const fileBlob = ref(null)

onMounted(function () {
  fetchFileName(fileNameUrl)
    .then((filename) => {
      data.filename = filename
      data.filenameFetched = true
    })
    .catch((error) => {
      console.error(error)
    })
})

const startDownload = () => {
  data.fileDownloading = true
  let xhr = new XMLHttpRequest()
  xhr.open('GET', downloadUrl, true)
  xhr.setRequestHeader('Authorization', authStore.FetchBearerToken())
  xhr.responseType = 'blob'
  let startTime = new Date().getTime()
  let startBytes = 0
  xhr.onprogress = function (event) {
    if (event.lengthComputable) {
      data.fileSizeMb = parseFloat((event.total / (1024 * 1024)).toFixed(2))
      let percentComplete = (event.loaded / event.total) * 100
      data.progressPercentage = parseFloat(percentComplete.toFixed(2))
      let currentTime = new Date().getTime()
      let timeElapsed = (currentTime - startTime) / 1000 // in seconds
      let loadedBytes = event.loaded - startBytes
      let speed = ((loadedBytes / timeElapsed) * 8) / 1000000 // Mbps
      startBytes = event.loaded
      startTime = currentTime
      data.lastThreeMbps.push(speed)
      if (data.lastThreeMbps.length > 5) {
        data.lastThreeMbps.shift()
      }
      data.mbps = parseFloat(
        parseFloat(data.lastThreeMbps.reduce((a, b) => a + b, 0) / data.lastThreeMbps.length).toFixed(1)
      )
    }
  }
  xhr.onload = function () {
    if (xhr.status === 200) {
      fileBlob.value = xhr.response
      data.fileDownloading = false
      data.fileDownloaded = true
      downloadBlob()
    }
  }
  xhr.onerror = function () {
    data.fileDownloading = false
    data.fileDownloaded = false
    data.isErrored = true
  }
  xhr.send()
}

const downloadBlob = () => {
  if (fileBlob.value) {
    const url = window.URL.createObjectURL(new Blob([fileBlob.value]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', data.filename)
    document.body.appendChild(link)
    link.style.display = 'none'
    link.click()
    link.remove()
  }
}

const fetchFileName = (url) => {
  return new Promise((resolve, reject) => {
    let filenameRequest = new XMLHttpRequest()
    filenameRequest.open('GET', url, true)
    filenameRequest.setRequestHeader('Authorization', authStore.FetchBearerToken())
    filenameRequest.responseType = 'text'
    filenameRequest.onload = function () {
      if (filenameRequest.status === 200) {
        resolve(filenameRequest.response)
      } else {
        reject('Error fetching filename')
      }
    }
    filenameRequest.send()
  })
}
</script>

<template>
  <section class="flex h-full w-full items-center justify-center">
    <section class="w-full max-w-screen-sm rounded-md border border-gray-400 p-2">
      <div class="font-medium">
        <font-awesome-icon icon="fa-solid fa-file" class="mr-2 text-blue-700" />
        {{ data.filename }}
        <span class="ml-4" v-if="data.fileSizeMb !== 0">({{ data.fileSizeMb }} MB)</span>
      </div>
      <FilledButton class="mt-4 w-full" :click="startDownload" v-if="isStartDownloadButtonShown" slim>
        <font-awesome-icon icon="fa-solid fa-circle-down" class="mr-2" />
        Click to Start Download
      </FilledButton>

      <div class="my-4 h-2.5 w-full rounded-full bg-gray-200 dark:bg-gray-700" v-if="data.fileDownloading">
        <div
          class="h-2.5 rounded-full bg-blue-600"
          :style="{
            width: `${data.progressPercentage}%`
          }"></div>
        <div class="flex w-full justify-between">
          <p>{{ data.progressPercentage }}%</p>
          <p>{{ data.mbps }} mbps</p>
        </div>
      </div>

      <FilledButton class="mt-4 w-full" :click="downloadBlob" v-if="data.fileDownloaded" slim>
        <font-awesome-icon icon="fa-solid fa-circle-down" class="mr-2" />
        Downloaded ! Click to Re-download File
      </FilledButton>
    </section>
  </section>
</template>

<style scoped></style>
