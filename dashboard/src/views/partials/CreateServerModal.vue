<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { reactive, ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { preventSpaceInput } from '@/vendor/utils.js'

const props = defineProps({
  callbackOnCreate: {
    type: Function,
    required: false,
    default: () => {}
  },
  callbackOnPop: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isModalOpen = ref(false)

const openModal = () => {
  newServerDetails.ip = ''
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
  props.callbackOnPop()
}

// Register Server state
const newServerDetails = reactive({
  ip: '',
  user: 'root',
  ssh_port: 22
})

const {
  mutate: registerServer,
  loading: isServerRegistering,
  onDone: onServerRegisterSuccess,
  onError: onServerRegisterFail
} = useMutation(
  gql`
    mutation ($input: NewServerInput!) {
      createServer(input: $input) {
        ip
        ssh_port
        user
      }
    }
  `,
  {
    variables: {
      input: newServerDetails
    }
  }
)

onServerRegisterSuccess(() => {
  closeModal()
  newServerDetails.ip = ''
  newServerDetails.user = 'root'
  newServerDetails.ssh_port = 22
  toast.success('Server registered successfully\nProceed for server setup')
  props.callbackOnCreate()
})

onServerRegisterFail((err) => {
  toast.error(err.message)
})

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <teleport to="body">
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
      <template v-slot:header>Add New Server</template>
      <template v-slot:body>
        Enter all the details to register a new server
        <form @submit.prevent="">
          <!--  IP Field   -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="ip"> IP Address </label>
            <div class="mt-1">
              <input
                id="ip"
                v-model="newServerDetails.ip"
                @keydown="preventSpaceInput"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder=""
                type="text" />
            </div>
          </div>
          <!--  SSH Port Field   -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="ssh_port"> SSH Port </label>
            <div class="mt-1">
              <input
                id="ssh_port"
                v-model="newServerDetails.ssh_port"
                @keydown="preventSpaceInput"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder=""
                type="text" />
            </div>
          </div>
          <!--  User Field   -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="user"> User </label>
            <div class="mt-1">
              <input
                id="user"
                v-model="newServerDetails.user"
                @keydown="preventSpaceInput"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder=""
                type="text" />
              <p class="pt-2 text-xs text-gray-800">
                <span class="pr-1 font-semibold text-red-500">*</span>
                If you want to choose any other user than <b>root</b>, make sure that if you have setup permissions
                required for installation for that user.
              </p>
            </div>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton :click="registerServer" :loading="isServerRegistering" type="primary"
          >Register Server
        </FilledButton>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
