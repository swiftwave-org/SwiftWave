<script setup>
import { useLazyQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { reactive, ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { preventSpaceInput } from '@/vendor/utils.js'
import FilledButton from '@/views/components/FilledButton.vue'

const props = defineProps({
  imageRegistryCredentialId: {
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
  load: loadImageRegistryCredentialDetails,
  refetch: refetchImageRegistryCredentialDetails,
  loading: isImageRegistryCredentialDetailsLoading,
  onError: onImageRegistryCredentialDetailsError,
  onResult: onImageRegistryCredentialDetailsResult
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      imageRegistryCredential(id: $id) {
        id
        url
        username
        password
      }
    }
  `,
  {
    id: props.imageRegistryCredentialId
  }
)

const imageRegitryCredentialDetails = reactive({
  url: '',
  username: '',
  password: ''
})

onImageRegistryCredentialDetailsError((err) => {
  toast.error(err.message)
  closeModal()
})

onImageRegistryCredentialDetailsResult((result) => {
  if (result.data.imageRegistryCredential) {
    imageRegitryCredentialDetails.url = result.data.imageRegistryCredential.url
    imageRegitryCredentialDetails.username = result.data.imageRegistryCredential.username
    imageRegitryCredentialDetails.password = result.data.imageRegistryCredential.password
  }
})

const isModalOpen = ref(false)
const openModal = () => {
  imageRegitryCredentialDetails.url = ''
  imageRegitryCredentialDetails.username = ''
  imageRegitryCredentialDetails.password = ''
  isModalOpen.value = true
  if (loadImageRegistryCredentialDetails() === false) {
    refetchImageRegistryCredentialDetails()
  }
}
const closeModal = () => {
  isModalOpen.value = false
}

defineExpose({
  openModal,
  closeModal
})

// Update Image Registry Credential
const {
  mutate: updateImageRegistryCredential,
  loading: isImageRegistryCredentialUpdating,
  onError: onImageRegistryCredentialUpdateError,
  onDone: onImageRegistryCredentialUpdateSuccess
} = useMutation(
  gql`
    mutation ($id: Uint!, $input: ImageRegistryCredentialInput!) {
      updateImageRegistryCredential(id: $id, input: $input) {
        id
      }
    }
  `,
  {
    variables: {
      id: props.imageRegistryCredentialId,
      input: imageRegitryCredentialDetails
    }
  }
)

onImageRegistryCredentialUpdateError((err) => {
  toast.error(err.message)
})

onImageRegistryCredentialUpdateSuccess(() => {
  closeModal()
  toast.success('Image Registry Credential updated successfully')
  props.callbackOnPop()
})
</script>

<template>
  <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
    <template v-slot:header>Edit Image Registry Credential</template>
    <template v-slot:body>
      Update the required information of the credential.
      <div v-if="isImageRegistryCredentialDetailsLoading" class="my-4">Loading details ...</div>
      <form v-else @submit.prevent="">
        <!--  Name Field   -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="url"> URL (example: ghcr.io) </label>
          <div class="mt-1">
            <input
              id="url"
              v-model="imageRegitryCredentialDetails.url"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="url"
              placeholder="URL"
              type="text" />
          </div>
        </div>
        <!-- Username Field -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="username"> Image Registry Username</label>
          <div class="mt-1">
            <input
              id="username"
              v-model="imageRegitryCredentialDetails.username"
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
              v-model="imageRegitryCredentialDetails.password"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="password"
              placeholder="Image Registry Password"
              type="text" />
          </div>
        </div>
      </form>
    </template>
    <template v-slot:footer>
      <FilledButton
        class="w-full"
        :click="updateImageRegistryCredential"
        :loading="isImageRegistryCredentialUpdating"
        type="primary"
        :disabled="imageRegitryCredentialDetails.type === '' || imageRegitryCredentialDetails.name === ''">
        Update Credential
      </FilledButton>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
