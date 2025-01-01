<script setup>
import ModalDialog from '@/views/components/ModalDialog.vue'
import { ref } from 'vue'
import { useLazyQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import FilledButton from '@/views/components/FilledButton.vue'
import { toast } from 'vue-sonner'

const props = defineProps({
  applicationIds: {
    type: Array,
    required: true
  },
  onDone: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isOpen = ref(false)
const openModal = () => {
  isOpen.value = true
  fetchAllAppDetails()
}
const closeModal = () => {
  if (isRebuilding.value) {
    toast.error('Wait until application rebuild initiation is completed')
    return
  }
  isLoadingApplicationDetailsFirstTime.value = true
  applicationDetails.value = {}
  isOpen.value = false
}
const { load: loadApplicationRaw, refetch: refetchApplicationRaw } = useLazyQuery(
  gql`
    query ($id: String!) {
      application(id: $id) {
        id
        name
      }
    }
  `,
  null,
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

const fetchLoadApplicationDetails = async (appId) => {
  let res = await loadApplicationRaw(
    null,
    {
      id: appId
    },
    null
  )
  if (!res) {
    res = await refetchApplicationRaw({
      id: appId
    })
  }
  return res?.application ?? res?.data?.application ?? {}
}

const applicationDetails = ref({})
const isLoadingApplicationDetailsFirstTime = ref(true)
const isRebuilding = ref(false)
const fetchAllAppDetails = async () => {
  let data = {}
  for (const appId of props.applicationIds) {
    data[appId] = await fetchLoadApplicationDetails(appId)
  }
  applicationDetails.value = data
  isLoadingApplicationDetailsFirstTime.value = false
  return data
}
defineExpose({
  openModal,
  closeModal
})

// Rebuild application
const { mutate: rebuildApplication } = useMutation(gql`
  mutation ($id: String!) {
    rebuildApplication(id: $id)
  }
`)

const rebuildApplications = async () => {
  isRebuilding.value = true
  for (const appId of props.applicationIds) {
    try {
      await rebuildApplication({
        id: appId
      })
    } catch (error) {
      toast.error(error.message)
    }
  }
  toast.success(
    props.applicationIds.length > 1
      ? 'Applications rebuild initiated successfully'
      : 'Application rebuild initiated successfully'
  )
  isRebuilding.value = false
  closeModal()
  if (props.onDone) {
    props.onDone()
  }
}
</script>

<template>
  <ModalDialog :is-open="isOpen" :close-modal="closeModal">
    <template v-slot:header>
      <span>Rebuild Application<span v-if="props.applicationIds.length > 1">s</span></span>
    </template>
    <template v-slot:body>
      Rebuild all the applications and redeploy them
      <p v-if="isLoadingApplicationDetailsFirstTime" class="mt-2">Loading application details...</p>
      <div v-else>
        <div v-for="app in applicationDetails" :key="app.id" class="mt-2 flex w-full flex-row items-center gap-2">
          <font-awesome-icon
            v-if="isRebuilding"
            icon="fa-solid fa-circle-notch"
            class="animate-spin text-base text-warning-500" />
          <font-awesome-icon v-else icon="fa-regular fa-circle" class="text-base text-warning-500" />
          <span class="text-secondary-800">{{ app.name }}</span>
        </div>
      </div>
    </template>
    <template v-slot:footer>
      <FilledButton
        type="primary"
        class="w-full"
        :click="rebuildApplications"
        :loading="isRebuilding"
        :disabled="isRebuilding || isLoadingApplicationDetailsFirstTime"
        >Rebuild Application<span v-if="props.applicationIds.length > 1">s</span></FilledButton
      >
    </template>
  </ModalDialog>
</template>

<style scoped></style>
