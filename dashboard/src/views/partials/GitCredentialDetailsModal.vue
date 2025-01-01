<script setup>
import Badge from '@/views/components/Badge.vue'
import { useLazyQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { computed, ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import Code from '@/views/components/Code.vue'

const props = defineProps({
  gitCredentialId: {
    type: Number,
    required: true
  }
})

const {
  load: loadGitCredentialDetails,
  refetch: refetchGitCredentialDetails,
  result: gitCredentialDetailsRaw,
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
        sshPublicKey
      }
    }
  `,
  {
    id: props.gitCredentialId
  }
)

const gitCredentialDetails = computed(() => gitCredentialDetailsRaw.value?.gitCredential ?? {})
const fetchGitCredentialDetails = () => {
  if (loadGitCredentialDetails() === false) {
    refetchGitCredentialDetails()
  }
}

onGitCredentialDetailsError((err) => {
  toast.error(err.message)
})

onGitCredentialDetailsResult(() => {
  isModalOpen.value = true
})

const isModalOpen = ref(false)

const openModal = () => fetchGitCredentialDetails()

const closeModal = () => {
  isModalOpen.value = false
}

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <!-- Git Credential details modal -->
  <ModalDialog :is-open="isModalOpen" :close-modal="() => (isModalOpen = false)">
    <template v-slot:header>Git Credential Details</template>
    <template v-slot:body>
      <div v-if="isGitCredentialDetailsLoading">Loading details ...</div>
      <template v-else>
        <div class="mt-4 flex w-full flex-row gap-2">
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Name</label>
            <div class="mt-1">
              <p class="block w-full rounded-md focus:ring-primary-500 sm:text-sm">
                {{ gitCredentialDetails.name ?? '' }}
              </p>
            </div>
          </div>
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Authentication Mode</label>
            <div class="mt-1">
              <Badge v-if="gitCredentialDetails.type === 'http'" type="success">HTTP Based</Badge>
              <Badge v-else-if="gitCredentialDetails.type === 'ssh'" type="warning">SSH Based</Badge>
              <Badge v-else type="danger">Unknown</Badge>
            </div>
          </div>
        </div>
        <div class="mt-4 flex w-full flex-row gap-2" v-if="gitCredentialDetails.type === 'http'">
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Username</label>
            <div class="mt-1">
              <p class="block w-full rounded-md focus:ring-primary-500 sm:text-sm">
                {{ gitCredentialDetails.username ?? '' }}
              </p>
            </div>
          </div>
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Password</label>
            <div class="mt-1">
              <p class="block w-full rounded-md italic focus:ring-primary-500 sm:text-sm">Can't show (Confidential)</p>
            </div>
          </div>
        </div>
        <div class="mt-4 flex w-full flex-col gap-1" v-if="gitCredentialDetails.type === 'ssh'">
          <label class="block text-sm font-medium text-gray-700">SSH Public Key</label>
          <div>
            <Code>{{ gitCredentialDetails.sshPublicKey ?? '' }}</Code>
          </div>
        </div>
      </template>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
