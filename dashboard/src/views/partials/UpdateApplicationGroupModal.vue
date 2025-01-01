<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { computed, ref } from 'vue'
import { useLazyQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import ComboBoxComponent from '@/views/components/ComboBoxComponent.vue'

const props = defineProps({
  applicationId: {
    type: String,
    required: true
  },
  currentGroupId: {
    type: Object || null,
    required: true
  },
  callbackOnUpdate: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isModalOpen = ref(false)
const selectedGroupId = ref(null)
const newGroupName = ref('')

const openModal = () => {
  selectedGroupId.value = props.currentGroupId
  fetchApplicationGroups()
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
}

const {
  result: applicationGroupsRaw,
  load: loadApplicationGroups,
  refetch: refetchApplicationGroups,
  onError: onApplicationGroupLoadFail
} = useLazyQuery(gql`
  query {
    applicationGroups {
      id
      name
    }
  }
`)

onApplicationGroupLoadFail((error) => {
  toast.error(error.message)
})

const applicationGroups = computed(() => {
  let l = applicationGroupsRaw?.value?.applicationGroups ?? []
  l.push({ id: null, name: 'No Project' })
  return l
})

function fetchApplicationGroups() {
  if (loadApplicationGroups() === false) {
    refetchApplicationGroups()
  }
}

const {
  mutate: updateApplicationGroupRaw,
  loading: isDomainRegistering,
  onDone: onApplicationGroupUpdateSuccess,
  onError: onApplicationGroupUpdateFail
} = useMutation(gql`
  mutation ($id: String!, $groupId: String) {
    updateApplicationGroup(id: $id, groupId: $groupId)
  }
`)

const { mutate: createApplicationGroupRaw, loading: isApplicationGroupCreating } = useMutation(gql`
  mutation ($name: String!) {
    createApplicationGroup(input: { name: $name }) {
      id
      name
    }
  }
`)

const updateApplicationGroup = async () => {
  if (selectedGroupId.value === 'new_group' && newGroupName.value) {
    try {
      const data = await createApplicationGroupRaw({ name: newGroupName.value })
      selectedGroupId.value = data.data?.createApplicationGroup.id
      await refetchApplicationGroups()
    } catch (e) {
      toast.error(e.message)
    }
  }
  await updateApplicationGroupRaw({
    id: props.applicationId,
    groupId: selectedGroupId.value
  })
}

onApplicationGroupUpdateSuccess((val) => {
  if (val.data?.updateApplicationGroup) {
    toast.success('Application group updated')
    closeModal()
    props.callbackOnUpdate()
  } else {
    toast.error('Application group update failed')
  }
})

onApplicationGroupUpdateFail((err) => {
  toast.error(err.message)
})

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
    <template v-slot:header>Update application project</template>
    <template v-slot:body>
      <p class="mb-4 mt-2 text-sm">Create a new project or pick an project group to assign your application to it.</p>
      <ComboBoxComponent
        :value="selectedGroupId"
        :options="applicationGroups"
        :title-from-option="(e) => e.name"
        :value-from-option="(e) => e.id"
        :new-option-data="(d) => ({ id: 'new_group', name: d })"
        :on-change="(e) => (selectedGroupId = e)"
        :on-new-option="(e) => (newGroupName = e)" />
    </template>
    <template v-slot:footer>
      <FilledButton
        class="w-full"
        :click="updateApplicationGroup"
        :loading="isDomainRegistering || isApplicationGroupCreating"
        type="primary"
        >Assign Application to Project
      </FilledButton>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
