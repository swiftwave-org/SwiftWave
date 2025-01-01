<script setup>
// Modal related methods
import { toast } from 'vue-sonner'
import { ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import Divider from '@/views/components/Divider.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'

const props = defineProps({
  serverId: {
    type: Number,
    required: true
  },
  refetchServer: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isModalOpen = ref(false)
const proxyType = ref('active')

const openModal = () => {
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
}

defineExpose({
  openModal,
  closeModal
})

// Enable proxy
const {
  mutate: enableProxyRaw,
  loading: enableProxyLoading,
  onError: enableProxyError,
  onDone: enableProxyDone
} = useMutation(gql`
  mutation EnableProxy($serverId: Uint!, $type: ProxyType!) {
    enableProxyOnServer(id: $serverId, type: $type)
  }
`)

enableProxyError((error) => {
  toast.error(error.message)
})

enableProxyDone((val) => {
  if (val.data.enableProxyOnServer) {
    toast.success(
      'Request to enable proxy has been submitted\nKeep refreshing server list & check logs for update\nThis can take upto 5 minutes to reflect in the system'
    )
    closeModal()
    props.refetchServer()
  }
})

const enableProxy = () => {
  enableProxyRaw({
    serverId: props.serverId,
    type: proxyType.value
  })
}
</script>

<template>
  <teleport to="body">
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
      <template v-slot:header>Enable Ingress Proxy</template>
      <template v-slot:body>
        <div class="mt-6 text-base">
          <p class="font-semibold">Active proxy vs Backup proxy</p>
          <p class="mt-2">
            <b>Active Proxy : </b> As soon as, you add new domain, ingress rules, redirect rules, it will be added
            immediately to the active proxy.
          </p>
          <p class="mt-0.5">
            <b>Backup Proxy : </b> The configuration will be updated every 30 minutes. So you may not see changes
            instantly. This is useful when you want to keep a backup proxy server in case the active proxy server goes
            down.
          </p>
        </div>
        <Divider />
        <label class="mt-4 block text-sm font-medium text-gray-700" for="domain">Proxy Type</label>
        <div class="mt-1 space-x-2">
          <select
            v-model="proxyType"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
            <option value="active">Active</option>
            <option value="backup">Backup</option>
          </select>
        </div>
        <p class="mt-2 text-sm text-gray-700">
          If you are not sure or don't have any active proxy server running, select <b>Active</b> from the dropdown.
        </p>
        <FilledButton class="mt-4 w-full" type="primary" :loading="enableProxyLoading" :click="enableProxy"
          >Enable Proxy
        </FilledButton>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
