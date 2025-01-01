<script setup>
import ModalDialog from '@/views/components/ModalDialog.vue'
import { reactive, ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import FilledButton from '@/views/components/FilledButton.vue'
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
const resetForm = () => {
  newGitCredential.name = ''
  newGitCredential.type = ''
  newGitCredential.username = ''
  newGitCredential.password = ''
  newGitCredential.sshPrivateKey = ''
}
const openModal = () => {
  resetForm()
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
  props.callbackOnPop()
  resetForm()
}

// New Git Credential form state
const newGitCredential = reactive({
  name: '',
  type: '',
  username: '',
  password: '',
  sshPrivateKey: ''
})

const {
  mutate: createGitCredential,
  loading: isGitCredentialCreating,
  onDone: onGitCredentialCreateSuccess,
  onError: onGitCredentialCreateFail
} = useMutation(
  gql`
    mutation ($input: GitCredentialInput!) {
      createGitCredential(input: $input) {
        id
        name
      }
    }
  `,
  {
    variables: {
      input: newGitCredential
    }
  }
)

onGitCredentialCreateSuccess(() => {
  closeModal()
  toast.success('Git Credential created successfully')
  props.callbackOnCreate()
})

onGitCredentialCreateFail((err) => {
  toast.error(err.message)
})

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <teleport to="body">
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
      <template v-slot:header>Add Git Credential</template>
      <template v-slot:body>
        Enter the necessary information for configuring the new Git Credential.
        <form @submit.prevent="">
          <!--  Name Field   -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="name">
              Name (Provide a name to identify the credential)
            </label>
            <div class="mt-1">
              <input
                id="name"
                v-model="newGitCredential.name"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                name="name"
                placeholder="Name"
                type="text" />
            </div>
          </div>
          <!-- Type Field -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="username"> Authentication Type </label>
            <div class="mt-1">
              <select
                id="git_credential"
                v-model="newGitCredential.type"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                <option selected value="">No Credential</option>
                <option value="http">HTTP</option>
                <option value="ssh">SSH</option>
              </select>
            </div>
          </div>
          <!-- Username Field -->
          <div class="mt-4" v-if="newGitCredential.type === 'http'">
            <label class="block text-sm font-medium text-gray-700" for="username"> Git Username </label>
            <div class="mt-1">
              <input
                id="username"
                v-model="newGitCredential.username"
                @keydown="preventSpaceInput"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                name="username"
                placeholder="Git Username"
                type="text" />
            </div>
          </div>
          <!-- Password Field -->
          <div class="mt-4" v-if="newGitCredential.type === 'http'">
            <label class="block text-sm font-medium text-gray-700" for="password"> Git Password / Auth Token </label>
            <div class="mt-1">
              <input
                id="password"
                v-model="newGitCredential.password"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                name="password"
                placeholder="Git Password"
                type="text" />
            </div>
          </div>
          <!-- Private Key Field -->
          <div class="mt-4" v-if="newGitCredential.type === 'ssh'">
            <label class="block text-sm font-medium text-gray-700" for="ssh_private_key">
              Private Key (In OpenSSH Format)
            </label>
            <p class="mt-1 text-sm text-gray-500">
              <b class="text-danger-500">NOTE:</b> Leave the input blank if you like to auto-generate the private key
            </p>
            <div class="mt-2">
              <textarea
                id="ssh_private_key"
                v-model="newGitCredential.sshPrivateKey"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="Private Key (In OpenSSH Format)"
                type="text"
                rows="5" />
            </div>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton
          :click="createGitCredential"
          :loading="isGitCredentialCreating"
          type="primary"
          :disabled="newGitCredential.type === '' || newGitCredential.name === ''">
          Add Now
        </FilledButton>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
