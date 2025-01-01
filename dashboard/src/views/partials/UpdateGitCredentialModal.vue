<script setup>
import { useLazyQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { reactive, ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { preventSpaceInput } from '@/vendor/utils.js'
import FilledButton from '@/views/components/FilledButton.vue'

const props = defineProps({
  gitCredentialId: {
    type: Number,
    required: true
  },
  callbackOnPop: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const {
  load: loadGitCredentialDetails,
  refetch: refetchGitCredentialDetails,
  loading: isGitCredentialDetailsLoading,
  onError: onGitCredentialDetailsError,
  onResult: onGitCredentialDetailsResult
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      gitCredential(id: $id) {
        id
        name
        type
        username
      }
    }
  `,
  {
    id: props.gitCredentialId
  }
)

const gitCredentialDetails = reactive({
  name: '',
  type: '',
  username: '',
  password: '',
  sshPrivateKey: ''
})

onGitCredentialDetailsError((err) => {
  toast.error(err.message)
  closeModal()
})

onGitCredentialDetailsResult((result) => {
  if (result.data.gitCredential) {
    gitCredentialDetails.name = result.data.gitCredential.name
    gitCredentialDetails.type = result.data.gitCredential.type
    gitCredentialDetails.username = result.data.gitCredential.username
    gitCredentialDetails.password = ''
    gitCredentialDetails.sshPrivateKey = ''
  }
})

const isModalOpen = ref(false)
const openModal = () => {
  gitCredentialDetails.name = ''
  gitCredentialDetails.type = ''
  gitCredentialDetails.username = ''
  gitCredentialDetails.password = ''
  gitCredentialDetails.sshPrivateKey = ''
  isModalOpen.value = true
  if (loadGitCredentialDetails() === false) {
    refetchGitCredentialDetails()
  }
}
const closeModal = () => {
  isModalOpen.value = false
}

defineExpose({
  openModal,
  closeModal
})

// Update Git Credential
const {
  mutate: updateGitCredential,
  loading: isGitCredentialUpdating,
  onError: onGitCredentialUpdateError,
  onDone: onGitCredentialUpdateSuccess
} = useMutation(
  gql`
    mutation ($id: Uint!, $input: GitCredentialInput!) {
      updateGitCredential(id: $id, input: $input) {
        id
      }
    }
  `,
  {
    variables: {
      id: props.gitCredentialId,
      input: gitCredentialDetails
    }
  }
)

onGitCredentialUpdateError((err) => {
  toast.error(err.message)
})

onGitCredentialUpdateSuccess(() => {
  closeModal()
  toast.success('Git Credential updated successfully')
  props.callbackOnPop()
})
</script>

<template>
  <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
    <template v-slot:header>Edit Git Credential</template>
    <template v-slot:body>
      Enter the necessary information for configuring the new Git Credential.
      <div v-if="isGitCredentialDetailsLoading" class="my-4">Loading details ...</div>
      <form v-else @submit.prevent="">
        <!--  Name Field   -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="name">
            Name (Provide a name to identify the credential)
          </label>
          <div class="mt-1">
            <input
              id="name"
              v-model="gitCredentialDetails.name"
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
              v-model="gitCredentialDetails.type"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              <option selected value="">No Credential</option>
              <option value="http">HTTP</option>
              <option value="ssh">SSH</option>
            </select>
          </div>
          <p class="mt-1 text-sm text-gray-500">
            If some applications already using this credential, changing the authentication type can cause issues.
          </p>
        </div>
        <!-- Username Field -->
        <div class="mt-4" v-if="gitCredentialDetails.type === 'http'">
          <label class="block text-sm font-medium text-gray-700" for="username"> Git Username </label>
          <div class="mt-1">
            <input
              id="username"
              v-model="gitCredentialDetails.username"
              @keydown="preventSpaceInput"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="username"
              placeholder="Git Username"
              type="text" />
          </div>
        </div>
        <!-- Password Field -->
        <div class="mt-4" v-if="gitCredentialDetails.type === 'http'">
          <label class="block text-sm font-medium text-gray-700" for="password"> Git Password / Auth Token </label>
          <div class="mt-1">
            <input
              id="password"
              v-model="gitCredentialDetails.password"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="password"
              placeholder="Git Password"
              type="text" />
          </div>
        </div>
        <!-- Private Key Field -->
        <div class="mt-4" v-if="gitCredentialDetails.type === 'ssh'">
          <label class="block text-sm font-medium text-gray-700" for="ssh_private_key">
            Private Key (In OpenSSH Format)
          </label>
          <p class="mt-1 text-sm text-gray-500">Leave the input blank to keep the existing private key</p>
          <div class="mt-2">
            <textarea
              id="ssh_private_key"
              v-model="gitCredentialDetails.sshPrivateKey"
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
        :click="updateGitCredential"
        :loading="isGitCredentialUpdating"
        type="primary"
        :disabled="gitCredentialDetails.type === '' || gitCredentialDetails.name === ''">
        Update Credential
      </FilledButton>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
