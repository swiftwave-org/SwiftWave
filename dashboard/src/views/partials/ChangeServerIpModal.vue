<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { preventSpaceInput } from '@/vendor/utils.js'
import { useRouter } from 'vue-router'

const props = defineProps({
  serverId: {
    type: Number,
    required: true
  },
  serverIp: {
    type: String,
    required: true
  }
})

const router = useRouter()
const isModalOpen = ref(false)
const ipChanged = ref(false)
const newServerIp = ref('')

const openModal = () => {
  newServerIp.value = props.serverIp
  ipChanged.value = false
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
}

const {
  mutate: changeServerIpRaw,
  loading: isRequestRunning,
  onDone: onChangeServerIpSuccess,
  onError: onChangeServerIpFail
} = useMutation(gql`
  mutation ($id: Uint!, $ip: String!) {
    changeServerIpAddress(id: $id, ip: $ip)
  }
`)

onChangeServerIpSuccess((val) => {
  if (val.data.changeServerIpAddress) {
    ipChanged.value = true
    startCountDown()
    closeModal()
  }
})

onChangeServerIpFail((err) => {
  toast.error(err.message)
})

const changeServerIp = () => {
  changeServerIpRaw({
    id: props.serverId,
    ip: newServerIp.value
  })
}

const timeCount = ref(5)

const startCountDown = () => {
  const interval = setInterval(() => {
    timeCount.value--
    if (timeCount.value === 0) {
      clearInterval(interval)
      ipChanged.value = false
      router.push({ name: 'Maintenance', query: { redirect: router.currentRoute.value.path } })
    }
  }, 1000)
}

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <teleport to="body">
    <ModalDialog :is-open="ipChanged" non-cancelable>
      <template v-slot:header>
        <span>🚀 Server IP Changed</span>
      </template>
      <template v-slot:body>
        <p class="mb-2">IP changed successfully! Swiftwave needs to restart.</p>
        <p>
          Redirecting to Maintenance Page in <b>{{ timeCount }}</b> seconds
        </p>
      </template>
    </ModalDialog>
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
      <template v-slot:header>Change Server IP</template>
      <template v-slot:body>
        Note: Changing the server IP address will restart the swiftwave service automatically.
        <form @submit.prevent="">
          <!--  IP Field   -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="ip"> IP Address </label>
            <div class="mt-1">
              <input
                id="ip"
                v-model="newServerIp"
                @keydown="preventSpaceInput"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                type="text" />
            </div>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton
          :click="changeServerIp"
          :loading="isRequestRunning"
          type="primary"
          :disabled="newServerIp === serverIp || newServerIp === ''"
          >Change Server IP
        </FilledButton>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
