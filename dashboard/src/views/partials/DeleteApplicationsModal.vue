<script setup>
import ModalDialog from '@/views/components/ModalDialog.vue'
import { ref } from 'vue'
import { useLazyQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import FilledButton from '@/views/components/FilledButton.vue'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'

const props = defineProps({
  applicationIds: {
    type: Array,
    required: true
  }
})

const router = useRouter()
const isOpen = ref(false)
const openModal = () => {
  isOpen.value = true
  fetchAllAppDetails()
}
const closeModal = () => {
  if (isDeleting.value) {
    toast.error('Wait until application deletion is completed')
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
        ingressRules {
          id
        }
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
const isDeleting = ref(false)
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

// Delete ingress rule
const { mutate: deleteIngressRule } = useMutation(gql`
  mutation ($id: Uint!) {
    deleteIngressRule(id: $id)
  }
`)

// Delete application
const { mutate: deleteApplication } = useMutation(gql`
  mutation ($id: String!) {
    deleteApplication(id: $id)
  }
`)

const deleteApplications = async () => {
  const confirmation = prompt('Type `delete` to confirm')
  if (confirmation !== 'delete') {
    return
  }

  isDeleting.value = true
  while (isDeleting.value) {
    const latestApplicationDetails = await fetchAllAppDetails()
    let isIngressRuleDeletionRequired = false
    // check if any ingress rules are using these applications
    for (const appId of props.applicationIds) {
      const app = latestApplicationDetails[appId]
      if (app?.ingressRules?.length > 0) {
        isIngressRuleDeletionRequired = true
        break
      }
    }
    if (isIngressRuleDeletionRequired) {
      for (const appId of props.applicationIds) {
        const app = latestApplicationDetails[appId]
        for (const ingressRule of app?.ingressRules ?? []) {
          if (ingressRule.status === 'deleting') continue
          try {
            await deleteIngressRule({
              id: ingressRule.id
            })
          } catch (error) {
            toast.error(error.message)
          }
        }
      }
      await new Promise((resolve) => setTimeout(resolve, 2000))
    } else {
      for (const appId of props.applicationIds) {
        try {
          await deleteApplication({
            id: appId
          })
        } catch (error) {
          toast.error(error.message)
        }
      }
      toast.success(
        props.applicationIds.length > 1
          ? 'Applications deletion initiated successfully ! It may take a few seconds to complete.'
          : 'Application deletion initiated successfully ! It may take a few seconds to complete.'
      )
      isDeleting.value = false
      closeModal()
      router.push('/applications')
    }
  }
}
</script>

<template>
  <ModalDialog :is-open="isOpen" :close-modal="closeModal">
    <template v-slot:header>
      <span>Delete Application<span v-if="props.applicationIds.length > 1">s</span></span>
    </template>
    <template v-slot:body>
      <div class="mb-2 mt-4 w-full rounded-md border border-warning-200 bg-warning-100 p-2 text-sm">
        Don't close this window until all the app<span v-if="props.applicationIds.length > 1">s</span> are deleted.
      </div>
      <p v-if="isLoadingApplicationDetailsFirstTime">Loading application details...</p>
      <div v-else class="mt-2">
        <div v-for="app in applicationDetails" :key="app.id" class="flex w-full flex-row items-center gap-2">
          <font-awesome-icon
            v-if="isDeleting"
            icon="fa-solid fa-circle-notch"
            class="animate-spin text-base text-warning-500" />
          <font-awesome-icon v-else icon="fa-regular fa-circle" class="text-base text-warning-500" />
          <span class="text-secondary-800"
            >{{ app.name }}&nbsp;<span v-if="(app?.ingressRules ?? []).length > 0"
              >({{ (app?.ingressRules ?? []).length }} Ingress Rules)</span
            ></span
          >
        </div>
      </div>
    </template>
    <template v-slot:footer>
      <FilledButton
        type="danger"
        class="w-full"
        :click="deleteApplications"
        :loading="isDeleting"
        :disabled="isDeleting || isLoadingApplicationDetailsFirstTime"
        >Delete Application<span v-if="props.applicationIds.length > 1">s</span></FilledButton
      >
    </template>
  </ModalDialog>
</template>

<style scoped></style>
