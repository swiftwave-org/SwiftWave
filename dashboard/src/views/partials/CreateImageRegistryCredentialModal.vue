<script setup>
import { preventSpaceInput } from '@/vendor/utils.js'
import { reactive, ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import ModalDialog from '@/views/components/ModalDialog.vue'
import FilledButton from '@/views/components/FilledButton.vue'

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
  newImageRegistryCredential.url = ''
  newImageRegistryCredential.username = ''
  newImageRegistryCredential.password = ''
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
  props.callbackOnPop()
}

// New Image Registry Credential form state
const newImageRegistryCredential = reactive({
  url: '',
  username: '',
  password: ''
})

const {
  mutate: createImageRegistryCredential,
  loading: isImageRegistryCredentialCreating,
  onDone: onImageRegistryCredentialCreateSuccess,
  onError: onImageRegistryCredentialCreateFail
} = useMutation(
  gql`
    mutation ($input: ImageRegistryCredentialInput!) {
      createImageRegistryCredential(input: $input) {
        id
        url
        username
        password
      }
    }
  `,
  {
    variables: {
      input: newImageRegistryCredential
    }
  }
)

onImageRegistryCredentialCreateSuccess(() => {
  closeModal()
  toast.success('Image Registry Credential created successfully')
  props.callbackOnCreate()
})

onImageRegistryCredentialCreateFail((err) => {
  toast.error(err.message)
})

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
    <template v-slot:header>Add Image Registry Credential</template>
    <template v-slot:body>
      Enter the necessary information for configuring the new Image Registry Credential.
      <form @submit.prevent="">
        <!--  Url Field   -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="url"> URL (example: ghcr.io) </label>
          <div class="mt-1">
            <input
              id="url"
              v-model="newImageRegistryCredential.url"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="url"
              placeholder="URL"
              type="text" />
          </div>
        </div>
        <!-- Username Field -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="username"> Image Registry Username </label>
          <div class="mt-1">
            <input
              id="username"
              v-model="newImageRegistryCredential.username"
              @keydown="preventSpaceInput"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="username"
              placeholder="Image Registry Username"
              type="text" />
          </div>
        </div>
        <!-- Password Field -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="password"> Image Registry Password </label>
          <div class="mt-1">
            <input
              id="password"
              v-model="newImageRegistryCredential.password"
              autocomplete="new-password"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="password"
              placeholder="Image Registry Password"
              type="text" />
          </div>
        </div>
      </form>
    </template>
    <template v-slot:footer>
      <FilledButton :click="createImageRegistryCredential" :loading="isImageRegistryCredentialCreating" type="primary"
        >Add Now
      </FilledButton>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
