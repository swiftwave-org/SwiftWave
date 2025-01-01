<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { computed, ref } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'

const props = defineProps({
  hostnames: {
    type: Array,
    required: true
  },
  updateHostnames: {
    type: Function,
    required: true
  }
})

const isModalOpen = ref(false)
const { result: serversDetailsRaw, onError: onServersDetailsError } = useQuery(
  gql`
    query {
      servers {
        hostname
        swarmMode
        scheduleDeployments
      }
    }
  `,
  null,
  {
    fetchPolicy: 'network-only',
    pollInterval: 30000
  }
)

const serversDetails = computed(() => serversDetailsRaw.value?.servers ?? [])

onServersDetailsError((err) => {
  toast.error('Failed to fetch servers \n' + err.message)
})

const openModal = () => {
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const toggleHostnameEntry = (hostname) => {
  if (props.hostnames.includes(hostname)) {
    props.updateHostnames(props.hostnames.filter((h) => h !== hostname))
  } else {
    const newHostnames = [...props.hostnames, hostname]
    props.updateHostnames(newHostnames)
  }
}
defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <teleport to="body">
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
      <template v-slot:header>Configure Preferred Servers</template>
      <template v-slot:body>
        <p>Select the servers on which you wish to run your application.</p>
        <div class="mt-3 flex flex-col gap-1.5">
          <div
            v-for="server in serversDetails"
            :key="server.hostname"
            @click.stop="() => toggleHostnameEntry(server.hostname)">
            <input
              type="checkbox"
              :checked="hostnames.includes(server.hostname)"
              class="me-2 h-4 w-4 rounded border-gray-300 bg-gray-100 text-primary-500 focus:ring-1 focus:ring-primary-500" />
            {{ server.hostname }}&nbsp;&nbsp;[{{ (server.swarmMode ?? '').toUpperCase() }}]&nbsp;&nbsp;<span
              v-if="server.scheduleDeployments === false">
              class="font-medium text-danger-500" >Deployment Disabled</span
            >
          </div>
        </div>

        <div class="mb-5 mt-6 rounded border-s-4 border-danger-200 bg-warning-50 px-3 py-2" role="alert">
          <p class="block text-justify text-sm text-black">
            If you have plan to enable <b>Docker Socket Proxy</b>, choose only one server.
          </p>
        </div>
      </template>
      <template v-slot:footer>
        <FilledButton type="primary" :click="closeModal">Confirm and Close Modal</FilledButton>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
