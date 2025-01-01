<script setup>
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api'
import 'monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution.js'
import { computed, onMounted, reactive, ref, shallowRef } from 'vue'
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { toast } from 'vue-sonner'

const editor = ref()
const editorInstance = shallowRef()
const stateRef = reactive({
  stackName: '',
  stackConfig: '',
  verificationStatus: -1,
  message: '',
  error: '',
  deployedApplicationsResult: []
})

onMounted(() => {
  editorInstance.value = monaco.editor.create(editor.value, {
    language: 'yaml',
    minimap: {
      enabled: false
    },
    lineNumbersMinChars: 2,
    tabSize: 2
  })
})

const {
  mutate: cleanupStack,
  loading: cleanupStackLoading,
  onDone: onCleanupStackDone,
  onError: onCleanupStackError
} = useMutation(gql`
  mutation CleanupStack($input: StackInput!) {
    cleanupStack(input: $input)
  }
`)

onCleanupStackDone((res) => {
  let val = res?.data?.cleanupStack ?? ''
  if (val) {
    editorInstance.value.setValue(val)
    verifyStackConfigHelper()
  }
})

onCleanupStackError((err) => {
  stateRef.error = err.message
  stateRef.message = ''
})

const {
  mutate: verifyStack,
  loading: verifyStackLoading,
  onDone: onVerifyStackDone,
  onError: onVerifyStackError
} = useMutation(gql`
  mutation VerifyStack($input: StackInput!) {
    verifyStack(input: $input) {
      success
      message
      error
    }
  }
`)

onVerifyStackDone((res) => {
  if (!res?.data?.verifyStack) return
  stateRef.verificationStatus = res?.data?.verifyStack?.success ? 1 : 0
  stateRef.message = res?.data?.verifyStack?.message ?? ''
  stateRef.error = res?.data?.verifyStack?.error ?? ''
})

onVerifyStackError((err) => {
  stateRef.verificationStatus = 0
  stateRef.error = err.message
  stateRef.message = ''
})

const cleanupStackConfigHelper = () => {
  stateRef.stackConfig = editorInstance.value.getValue()
  cleanupStack({
    input: {
      content: stateRef.stackConfig,
      variables: [
        {
          name: 'STACK_NAME',
          value: stateRef.stackName
        }
      ]
    }
  })
}

const verifyStackConfigHelper = () => {
  stateRef.stackConfig = editorInstance.value.getValue()
  verifyStack({
    input: {
      content: stateRef.stackConfig,
      variables: [
        {
          name: 'STACK_NAME',
          value: stateRef.stackName
        }
      ]
    }
  })
}

const isCleanupAndVerifyStackConfigLoading = computed(() => {
  return cleanupStackLoading.value || verifyStackLoading.value
})

// Deploy Stack
const {
  mutate: deployStack,
  loading: deployStackLoading,
  onDone: onDeployStackDone,
  onError: onDeployStackError
} = useMutation(gql`
  mutation DeployStack($input: StackInput!) {
    deployStack(input: $input) {
      success
      message
      application {
        id
        name
      }
    }
  }
`)

onDeployStackDone((res) => {
  if (!res?.data?.deployStack) return
  stateRef.deployedApplicationsResult = res?.data?.deployStack ?? []
  isModalOpen.value = true
})

onDeployStackError((err) => {
  toast.error(err.message)
})

const deployStackHelper = () => {
  deployStack({
    input: {
      content: stateRef.stackConfig,
      variables: [
        {
          name: 'STACK_NAME',
          value: stateRef.stackName
        }
      ]
    }
  })
}

// Result modal
const isModalOpen = ref(false)
const openUrlInNewPage = (url) => {
  window.open(url)
}
</script>

<template>
  <p class="text-xl font-semibold">
    <font-awesome-icon icon="fa-solid fa-cubes-stacked" class="mr-2 text-primary-600" />
    Deploy Stack
  </p>
  <section class="mx-auto mt-8 flex h-full w-full max-w-7xl space-x-8">
    <div class="h-full w-1/2">
      <!--  Stack Name  -->
      <div>
        <label class="block text-sm font-medium text-gray-700">Stack Name<span class="text-red-600"> *</span></label>
        <div class="mt-1">
          <input
            autocomplete="off"
            v-model="stateRef.stackName"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            placeholder="Enter stack name"
            type="text" />
        </div>
      </div>
      <!--   Stack Config (Yaml)   -->
      <div class="mt-3 h-[80%]">
        <label class="block text-sm font-medium text-gray-700"
          >Stack Config (YAML)<span class="text-red-600"> *</span></label
        >
        <div ref="editor" class="mt-1 h-full w-full overflow-hidden rounded-md border-2 border-primary-300" />
      </div>
    </div>
    <div class="w-1/2 pt-4">
      <p class="select-none text-base font-semibold">⚙️ Actions</p>
      <FilledButton
        class="mt-4 w-full"
        type="primary"
        :click="cleanupStackConfigHelper"
        :loading="isCleanupAndVerifyStackConfigLoading"
        >Cleanup & Verify Stack Config
      </FilledButton>
      <!--   Verification Status   -->
      <div class="mt-2 p-3" v-if="stateRef.verificationStatus !== -1">
        <font-awesome-icon
          v-if="stateRef.verificationStatus === 1"
          icon="fa-solid fa-circle-check"
          class="mr-1 text-xl text-success-500" />
        <font-awesome-icon
          v-else-if="stateRef.verificationStatus === 0"
          icon="fa-solid fa-circle-xmark"
          class="mr-1 text-xl text-danger-500" />
        {{ stateRef.verificationStatus === 1 ? 'Verified' : 'Verification Failed' }}
      </div>
      <!--    Verification Result  -->
      <div
        class="mt-2 whitespace-pre-line rounded-md border-2 border-danger-300 bg-danger-50 p-3"
        v-if="stateRef.error !== ''">
        {{ stateRef.error }}
      </div>
      <div
        class="mt-2 whitespace-pre-line rounded-md border-2 border-success-300 bg-success-50 p-3"
        v-if="stateRef.message !== ''">
        {{ stateRef.message }}
      </div>
      <!--   Deploy Button   -->
      <FilledButton
        v-if="stateRef.verificationStatus === 1"
        class="mt-6 w-full"
        type="primary"
        :loading="deployStackLoading"
        :click="deployStackHelper">
        <font-awesome-icon icon="fa-solid fa-hammer" class="mr-2" />
        Deploy Stack
      </FilledButton>
      <!--  Modal to show result    -->
      <ModalDialog :is-open="isModalOpen" non-cancelable>
        <template v-slot:header>🎉 Deployed Successfully</template>
        <template v-slot:body>
          <div class="flex flex-col space-y-3 pt-3">
            <div
              class="flex items-center space-x-2"
              v-for="result in stateRef.deployedApplicationsResult"
              :key="result.application.id">
              <font-awesome-icon
                v-if="result.success"
                icon="fa-solid fa-circle-check"
                class="text-base text-success-500" />
              <font-awesome-icon v-else icon="fa-solid fa-circle-xmark" class="text-base text-danger-500" />
              <p>
                {{ result.application?.name ?? 'N/A' }}
                <span v-if="result.message !== '' && !result.success"> - {{ result.message }}</span>
              </p>
              <FilledButton
                v-if="result.success"
                type="primary"
                slim
                :click="
                  () =>
                    openUrlInNewPage(
                      $router.resolve({
                        name: 'Application Details Deployments',
                        params: { id: result.application.id }
                      }).href
                    )
                ">
                View
              </FilledButton>
            </div>
            <div v-if="stateRef.deployedApplicationsResult.length === 0" class="text-center text-gray-500">
              No applications deployed
            </div>
          </div>
        </template>
        <template v-slot:footer>
          <FilledButton type="warning" class="w-full" :click="() => $router.replace('/applications')"
            >Go to Applications List
          </FilledButton>
        </template>
      </ModalDialog>
    </div>
  </section>
</template>

<style scoped></style>
