<script setup>
import ModalDialog from '@/views/components/ModalDialog.vue'
import DockerfileEditorOnly from '@/views/partials/DeployApplication/DockerfileEditorOnly.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { ref, watch } from 'vue'

const props = defineProps({
  code: {
    type: String,
    default: ''
  },
  submit: {
    type: Function,
    default: () => {}
  },
  isOpen: {
    type: Boolean,
    default: false
  },
  closeModal: {
    type: Function,
    default: () => {}
  },
  dockerConfigurationGenerating: {
    type: Boolean,
    default: false
  }
})

const codeCurrent = ref('')
const showCodeEditor = ref(false)

const codeChanged = (code) => {
  codeCurrent.value = code
}

const cancel = () => {
  codeCurrent.value = props.code
  props.closeModal()
  showCodeEditor.value = false
}

const generateDockerConfiguration = () => {
  props.submit(codeCurrent.value)
}

watch(
  () => props.isOpen,
  (isOpen) => {
    codeCurrent.value = props.code
    showCodeEditor.value = isOpen
  }
)
</script>

<template>
  <ModalDialog :is-open="isOpen" :non-cancelable="true" full-screen width="xl">
    <template v-slot:header>
      <p class="mb-2">Dockerfile Editor</p>
      <p class="mb-8 mr-8 text-sm font-normal text-gray-600">
        Edit the Dockerfile for your application. This will be used to build your application image.
      </p>
    </template>
    <template v-slot:body>
      <DockerfileEditorOnly v-if="showCodeEditor" :code="codeCurrent" :code-changed="codeChanged" :submit="submit" />
    </template>
    <template v-slot:footer>
      <div class="flex justify-end">
        <FilledButton class="ml-3" type="secondary" @click.prevent="cancel">Cancel</FilledButton>
        <FilledButton
          :loading="dockerConfigurationGenerating"
          class="ml-3"
          type="primary"
          @click.prevent="generateDockerConfiguration">
          Submit
        </FilledButton>
      </div>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
