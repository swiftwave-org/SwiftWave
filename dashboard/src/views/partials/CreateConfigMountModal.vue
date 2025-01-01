<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { computed, reactive, ref } from 'vue'
import { preventSpaceInput } from '@/vendor/utils.js'
import { toast } from 'vue-sonner'

const props = defineProps({
  createConfigMount: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isModalOpen = ref(false)
// Create persistent volume
const newConfigMountDetails = reactive({
  content: '',
  mountingPath: '',
  uid: 0,
  gid: 0
})

const openModal = () => {
  newConfigMountDetails.content = ''
  newConfigMountDetails.mountingPath = ''
  newConfigMountDetails.uid = 0
  newConfigMountDetails.gid = 0
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
}

const registerConfigMount = () => {
  try {
    props.createConfigMount(newConfigMountDetails)
    closeModal()
  } catch (e) {
    toast.error(e.message)
  }
}

const isAllFieldsFilled = computed(() => {
  return newConfigMountDetails.mountingPath !== '' && newConfigMountDetails.content !== ''
})

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
    <template v-slot:header>Add Config Mount</template>
    <template v-slot:body>
      Provide all the details for the new config mount.
      <form @submit.prevent="">
        <!--  Name Field   -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700"> Config Mount Path </label>
          <div class="mt-1">
            <input
              v-model="newConfigMountDetails.mountingPath"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="Provide path for the config mount"
              type="text"
              @keydown="preventSpaceInput" />
          </div>
        </div>

        <div class="mt-2 flex w-full flex-row gap-2">
          <!--   File UID     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">File UID</label>
            <div class="mt-1">
              <input
                v-model="newConfigMountDetails.uid"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="File UID"
                type="number" />
            </div>
          </div>
          <!--   File Gid     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">File GID</label>
            <div class="mt-1">
              <input
                v-model="newConfigMountDetails.gid"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="Mount GID"
                type="number" />
            </div>
          </div>
        </div>

        <div class="mt-2 w-full">
          <label class="block text-sm font-medium text-gray-700">Config Content</label>
          <div class="mt-1">
            <textarea
              rows="20"
              class="mt-2 w-full rounded-lg border-gray-300 align-top shadow-sm focus:border-primary-500 focus:ring-primary-500"
              v-model="newConfigMountDetails.content"></textarea>
          </div>
        </div>
      </form>
    </template>
    <template v-slot:footer>
      <FilledButton type="primary" :click="registerConfigMount" :disabled="!isAllFieldsFilled"
        >Confirm & Add
      </FilledButton>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
