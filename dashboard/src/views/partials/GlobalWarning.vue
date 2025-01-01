<!-- This page is going to have global warning (which is going to take full screen to show)-->

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { computed, onMounted, watch } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { useLazyQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { useAuthStore } from '@/store/auth.js'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const {
  result: noOfPreparedServersResult,
  load: loadNoOfPreparedServersRaw,
  refetch: refetchNoOfPreparedServersRaw
} = useLazyQuery(gql`
  query {
    noOfPreparedServers
  }
`)

const showNoServerConfiguredWarning = computed(() => {
  if (!authStore.IsLoggedIn) return false
  if (['System Logs', 'Setup', 'Users', 'Servers', 'Server Logs'].includes(route.name)) return false
  if (!noOfPreparedServersResult.value) return false
  if (noOfPreparedServersResult.value.noOfPreparedServers === undefined) return false
  if (noOfPreparedServersResult.value.noOfPreparedServers === null) return false
  return noOfPreparedServersResult.value.noOfPreparedServers === 0
})

const loadNoOfPreparedServers = () => {
  if (!authStore.IsLoggedIn) return
  if (loadNoOfPreparedServersRaw() === false) {
    refetchNoOfPreparedServersRaw()
  }
}

watch(authStore, () => {
  if (!authStore.IsLoggedIn) return
  loadNoOfPreparedServers()
})

const setupServer = () => {
  router.push({ name: 'Servers' })
}

onMounted(() => {
  setInterval(loadNoOfPreparedServers, 10000)
})
</script>

<template>
  <teleport to="body">
    <ModalDialog :is-open="showNoServerConfiguredWarning" non-cancelable>
      <template #header>No Server Configured</template>
      <template #body>
        <p>You have not configured any server yet. Please configure a server before performing any other actions</p>
      </template>
      <template #footer>
        <FilledButton class="w-full" :click="setupServer">
          <font-awesome-icon icon="fa-solid fa-hammer" class="mr-2" />
          Setup Server
        </FilledButton>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
