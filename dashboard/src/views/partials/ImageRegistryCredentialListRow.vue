<script setup>
import TableRow from '@/views/components/Table/TableRow.vue'
import TextButton from '@/views/components/TextButton.vue'
import SecuredText from '@/views/components/SecuredText.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import UpdateImageRegistryCredentialModal from '@/views/partials/UpdateImageRegistryCredentialModal.vue'
import { ref } from 'vue'

defineProps({
  imageRegistryCredential: {
    type: Object,
    required: true
  },
  deleteImageRegistryCredential: {
    type: Function,
    required: true
  },
  onUpdateImageRegistryCredential: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const updateImageRegistryCredential = ref(null)
const openEditModal = () => {
  if (updateImageRegistryCredential.value) updateImageRegistryCredential.value.openModal()
}
</script>

<template>
  <!-- Update Image Registry Credential modal -->
  <UpdateImageRegistryCredentialModal
    :image-registry-credential-id="imageRegistryCredential.id"
    ref="updateImageRegistryCredential"
    :callback-on-pop="onUpdateImageRegistryCredential" />

  <tr>
    <TableRow align="left">
      <div class="text-sm font-medium text-gray-900">
        {{ imageRegistryCredential.url }}
      </div>
    </TableRow>
    <TableRow align="center">
      <div class="text-sm text-gray-900">{{ imageRegistryCredential.username }}</div>
    </TableRow>
    <TableRow align="center" flex>
      <SecuredText>{{ imageRegistryCredential.password }}</SecuredText>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton type="secondary" slim :click="openEditModal">Edit Details</FilledButton>
    </TableRow>
    <TableRow align="right">
      <TextButton :click="() => deleteImageRegistryCredential(imageRegistryCredential)" type="danger">
        Delete
      </TextButton>
    </TableRow>
  </tr>
</template>

<style scoped></style>
