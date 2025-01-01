<script setup>
import TableRow from '@/views/components/Table/TableRow.vue'
import TextButton from '@/views/components/TextButton.vue'
import Badge from '@/views/components/Badge.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { ref } from 'vue'
import UpdateGitCredentialModal from '@/views/partials/UpdateGitCredentialModal.vue'
import GitCredentialDetailsModal from '@/views/partials/GitCredentialDetailsModal.vue'

defineProps({
  gitCredential: {
    type: Object,
    required: true
  },
  deleteGitCredential: {
    type: Function,
    required: true
  },
  onUpdateGitCredential: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const updateGitCredentialModalRef = ref(null)
const gitCredentialDetailsModalRef = ref(null)

const openEditModal = () => {
  if (updateGitCredentialModalRef.value) updateGitCredentialModalRef.value.openModal()
}
const openDetailsModal = () => {
  if (gitCredentialDetailsModalRef.value) gitCredentialDetailsModalRef.value.openModal()
}
</script>

<template>
  <!-- Update Git Credential modal -->
  <UpdateGitCredentialModal
    :git-credential-id="gitCredential.id"
    ref="updateGitCredentialModalRef"
    :callback-on-pop="onUpdateGitCredential" />
  <!-- Git Credential details modal -->
  <GitCredentialDetailsModal :git-credential-id="gitCredential.id" ref="gitCredentialDetailsModalRef" />

  <tr>
    <TableRow align="left">
      <div class="text-sm font-medium text-gray-900">
        {{ gitCredential.name }}
      </div>
    </TableRow>
    <TableRow align="center">
      <Badge v-if="gitCredential.type === 'http'" type="success">HTTP Authentication</Badge>
      <Badge v-else-if="gitCredential.type === 'ssh'" type="warning">SSH Authentication</Badge>
      <Badge v-else type="danger">Unknown</Badge>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton slim :onclick="openDetailsModal">Show Details</FilledButton>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton type="secondary" slim :onclick="openEditModal">Edit Details</FilledButton>
    </TableRow>
    <TableRow align="right">
      <TextButton :click="() => deleteGitCredential(gitCredential)" type="danger"> Delete</TextButton>
    </TableRow>
  </tr>
</template>

<style scoped></style>
