<script setup>
import { useRouter } from 'vue-router'
import FilledButton from '@/views/components/FilledButton.vue'
import DeleteApplicationsModal from '@/views/partials/DeleteApplicationsModal.vue'
import { ref } from 'vue'

const router = useRouter()
const deleteApplicationsModal = ref(null)

function deleteApplicationWithConfirmation() {
  if (deleteApplicationsModal.value) {
    deleteApplicationsModal.value.openModal()
  }
}
</script>

<template>
  <DeleteApplicationsModal ref="deleteApplicationsModal" :application-ids="[router.currentRoute.value.params.id]" />
  <div class="w-full rounded-md border border-warning-200 bg-warning-100 p-2">
    Use the below options with caution. These actions are non-reversible.
  </div>
  <div class="mt-3 flex flex-col items-start">
    <p class="font-bold text-danger-500">Do you like to delete this application ?</p>
    <p class="mt-2">This action will remove these stuffs -</p>
    <ul class="list-inside list-disc">
      <li>Application</li>
      <li>Ingress Rules</li>
      <li>Related Deployments</li>
      <li>Deployment Logs</li>
      <li>Environment Variables</li>
      <li>Persistent Volume Bindings</li>
      <li>Uploaded Source Code</li>
    </ul>

    <FilledButton class="mt-6" type="danger" :click="deleteApplicationWithConfirmation"
      >Delete Ingress Rules & Application
    </FilledButton>
  </div>
</template>

<style scoped></style>
