<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import axios from 'axios'
import { useAuthStore } from '@/store/auth.js'
import { useLazyQuery, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import FilledButton from '@/views/components/FilledButton.vue'
import createTar from '@swiftwave/tartplus'
import DockerfileEditor from '@/views/partials/DeployApplication/DockerfileEditor.vue'
import BuildArgInput from '@/views/partials/BuildArgInput.vue'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import newApplicationUpdater from '@/store/applicationUpdater.js'
import { useRouter } from 'vue-router'
import CreateGitCredentialModal from '@/views/partials/CreateGitCredentialModal.vue'
import CreateImageRegistryCredentialModal from '@/views/partials/CreateImageRegistryCredentialModal.vue'
import ChooseOtherDockerConfigurationModal from '@/views/partials/ChooseOtherDockerConfigurationModal.vue'

const router = useRouter()
const authStore = useAuthStore()
const applicationUpdater = newApplicationUpdater(router.currentRoute.value.params.id)()

const applicationSourceType = computed(() => {
  if (applicationExistingDetailsResult.value) {
    return applicationExistingDetailsResult.value?.application?.latestDeployment?.upstreamType ?? ''
  } else {
    return null
  }
})
const sourceCodeCompressedFileFieldRef = ref(null)
const stateRef = reactive({
  sourceCodeCompressedFileName: '',
  gitCredentialID: 0,
  gitRepoUrl: '',
  gitBranch: '',
  codePath: '',
  command: '',
  imageRegistryCredentialID: 0,
  dockerImage: '',
  isUploadingSourceCode: false,
  detectedServiceName: '',
  dockerFile: '',
  dockerBuildArgs: [],
  buildArgs: {},
  isDockerFileEditorOpen: false,
  isDockerConfigurationGenerated: false
})
const availableGitBranches = ref([])

watch(
  stateRef,
  () => {
    applicationUpdater.updateApplicationSource(stateRef)
  },
  { deep: true }
)

watch(
  () => applicationUpdater.isConfigurationUpdated,
  (updateStatus) => {
    if (updateStatus === false) {
      prefillDetails()
    }
  }
)

const applicationExistingDetailsResult = computed(() => {
  return applicationUpdater.applicationExistingDetailsResult ?? {}
})

// Fetch git branches
const {
  load: fetchGitBranchesRaw,
  refetch: refetchGitBranchesRaw,
  loading: fetchingGitBranches,
  onError: onFetchGitBranchesError,
  onResult: onFetchGitBranchesResult,
  variables: fetchGitBranchesVariables
} = useLazyQuery(
  gql`
    query ($input: GitBranchesQueryInput!) {
      gitBranches(input: $input)
    }
  `,
  null,
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

const fetchGitBranches = () => {
  if (!stateRef.gitRepoUrl) {
    return
  }
  let gitRepoUrl = stateRef.gitRepoUrl.trim()
  fetchGitBranchesVariables.value = {
    input: {
      gitCredentialId: stateRef.gitCredentialID,
      repositoryUrl: gitRepoUrl
    }
  }
  if (fetchGitBranchesRaw() === false) {
    refetchGitBranchesRaw()
  }
}

onFetchGitBranchesResult((d) => {
  if (d.data && d.data.gitBranches) {
    availableGitBranches.value = d.data.gitBranches
    toast.success('Available branches fetched')
  }
})
onFetchGitBranchesError((err) => {
  toast.error(err.message)
  availableGitBranches.value = []
  stateRef.gitBranch = ''
})

function prefillDetails() {
  if (applicationExistingDetailsResult.value && applicationExistingDetailsResult.value.application) {
    stateRef.command = applicationExistingDetailsResult.value.application.command
    if (applicationExistingDetailsResult.value.application.latestDeployment.upstreamType === 'git') {
      stateRef.gitRepoUrl = applicationExistingDetailsResult.value.application.latestDeployment.repositoryUrl
      stateRef.gitBranch = applicationExistingDetailsResult.value.application.latestDeployment.repositoryBranch
      stateRef.gitCredentialID = applicationExistingDetailsResult.value.application.latestDeployment.gitCredentialID
      stateRef.codePath = applicationExistingDetailsResult.value.application.latestDeployment.codePath
    }
    stateRef.isDockerConfigurationGenerated = true
    stateRef.detectedServiceName = 'Taken from existing deployment'
    stateRef.dockerFile = applicationExistingDetailsResult.value.application.latestDeployment.dockerfile
    const buildArgs = applicationExistingDetailsResult.value.application.latestDeployment.buildArgs
    stateRef.buildArgs = {}
    stateRef.dockerBuildArgs = []
    for (const buildArg of buildArgs) {
      stateRef.buildArgs[buildArg.key] = buildArg.value
      stateRef.dockerBuildArgs.push({
        key: buildArg.key,
        description: '',
        value: buildArg.value
      })
    }
    stateRef.sourceCodeCompressedFileName =
      applicationExistingDetailsResult.value.application.latestDeployment.sourceCodeCompressedFileName
    stateRef.dockerImage = applicationExistingDetailsResult.value.application.latestDeployment.dockerImage
    stateRef.imageRegistryCredentialID =
      applicationExistingDetailsResult.value.application.latestDeployment.imageRegistryCredentialID
  }
}

watch(applicationExistingDetailsResult, () => {
  prefillDetails()
})

onMounted(() => {
  prefillDetails()
  fetchGitBranches()
})

const openDockerFileEditor = () => {
  stateRef.isDockerFileEditorOpen = true
}

const closeDockerFileEditor = () => {
  stateRef.isDockerFileEditorOpen = false
}

const enableGenerateConfigurationButton = computed(() => {
  if (applicationSourceType.value === 'git') {
    return stateRef.gitRepoUrl !== '' && stateRef.gitBranch !== ''
  } else if (applicationSourceType.value === 'sourceCode') {
    return stateRef.sourceCodeCompressedFileName !== ''
  } else if (applicationSourceType.value === 'image') {
    return stateRef.dockerImage !== ''
  } else {
    return false
  }
})

// List Image Registry Credentials query
const {
  result: imageRegistryCredentialList,
  onError: onImageRegistryCredentialListError,
  refetch: refetchImageRegistryCredentialList
} = useQuery(
  gql`
    query {
      imageRegistryCredentials {
        id
        url
        username
      }
    }
  `,
  null,
  {
    pollInterval: 10000
  }
)
const imageRegistryCredentials = computed(() => imageRegistryCredentialList.value?.imageRegistryCredentials ?? [])

onImageRegistryCredentialListError((err) => toast.error(err.message))
// Fetch git credentials
const {
  result: gitCredentialList,
  onError: onGitCredentialListError,
  refetch: refetchGitCredentialList
} = useQuery(
  gql`
    query {
      gitCredentials {
        id
        name
        type
      }
    }
  `,
  null,
  {
    pollInterval: 10000
  }
)
const gitCredentials = computed(() => gitCredentialList.value?.gitCredentials ?? [])

onGitCredentialListError((err) => toast.error(err.message))

const HTTP_BASE_URL = getHttpBaseUrl()

async function uploadTarFile(fileblob) {
  try {
    var data = new FormData()
    data.append('file', fileblob, 'file.tar')
    const res = await axios({
      method: 'post',
      url: `${HTTP_BASE_URL}/upload/code`,
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: authStore.FetchBearerToken()
      },
      data: data
    })
    return {
      success: true,
      message: res.data.message,
      file: res.data.file
    }
  } catch (error) {
    return {
      success: false,
      message: error.response.data.message,
      file: null
    }
  }
}

const uploadSourceCode = async () => {
  stateRef.isUploadingSourceCode = true
  try {
    const file = await createTar(sourceCodeCompressedFileFieldRef.value.files, ['.gitignore'])
    const res = await uploadTarFile(file)
    if (res.success) {
      stateRef.sourceCodeCompressedFileName = res.file
      toast.success(res.message)
    } else {
      toast.error(res.message)
    }
  } catch (e) {
    toast.error('failed to upload source code')
  }
  stateRef.isUploadingSourceCode = false
}

// Generate Configuration
const {
  load: generateConfigurationLoad,
  refetch: generateConfigurationRefetch,
  loading: dockerConfigGeneratorGenerating,
  onError: onGenerateConfigurationError,
  onResult: onGenerateConfigurationSuccess,
  variables: generateConfigurationVariables
} = useLazyQuery(
  gql`
    query ($input: DockerConfigGeneratorInput!) {
      dockerConfigGenerator(input: $input) {
        detectedServiceName
        dockerFile
        dockerBuildArgs {
          key
          description
          defaultValue
        }
      }
    }
  `,
  {
    input: {}
  },
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

onGenerateConfigurationError((err) => toast.error(err.message))

onGenerateConfigurationSuccess((res) => {
  if (res.data && res.data.dockerConfigGenerator) {
    updateDockerConfiguration(res.data.dockerConfigGenerator)
    closeDockerFileEditor()
  }
})

const updateDockerConfiguration = (dockerConfig) => {
  stateRef.detectedServiceName = dockerConfig.detectedServiceName
  stateRef.dockerFile = dockerConfig.dockerFile
  stateRef.dockerBuildArgs = dockerConfig.dockerBuildArgs
  // set default build args if not set
  for (const buildArg of stateRef.dockerBuildArgs) {
    stateRef.buildArgs[buildArg.key] = buildArg.defaultValue
  }
  // delete build args if not present in dockerBuildArgs
  for (const buildArgKey in stateRef.buildArgs) {
    if (!stateRef.dockerBuildArgs.some((buildArg) => buildArg.key === buildArgKey)) {
      delete stateRef.buildArgs[buildArgKey]
    }
  }
  stateRef.isDockerConfigurationGenerated = true
}

const updateBuildArg = (key, value) => {
  stateRef.buildArgs[key] = value
}

const generateConfiguration = () => {
  if (applicationSourceType.value === 'image') {
    stateRef.detectedServiceName = "😅 You don't need configuration for docker image"
    stateRef.isDockerConfigurationGenerated = true
  } else {
    let gitCredentialID = parseInt(stateRef.gitCredentialID.toString())
    generateConfigurationVariables.value.input = {
      sourceType: applicationSourceType.value,
      gitCredentialID: gitCredentialID === 0 ? null : gitCredentialID,
      repositoryBranch: stateRef.gitBranch === '' ? null : stateRef.gitBranch,
      repositoryUrl: stateRef.gitRepoUrl,
      codePath: stateRef.codePath,
      customDockerFile: '',
      sourceCodeCompressedFileName:
        stateRef.sourceCodeCompressedFileName === '' ? null : stateRef.sourceCodeCompressedFileName
    }
    if (generateConfigurationLoad() === false) {
      generateConfigurationRefetch()
    }
  }
}

const generateConfigurationForCustomDockerFile = (customDockerFile) => {
  generateConfigurationVariables.value.input = {
    sourceType: 'custom',
    gitCredentialID: null,
    repositoryBranch: null,
    repositoryUrl: null,
    codePath: null,
    customDockerFile: customDockerFile,
    sourceCodeCompressedFileName: null
  }
  if (generateConfigurationLoad() === false) {
    generateConfigurationRefetch()
  }
}

// Create Git Credential
const createGitCredentialModalRef = ref(null)
const openCreateGitCredentialModal = computed(() => createGitCredentialModalRef.value?.openModal ?? (() => {}))

// Create Image Registry Credential
const createImageRegistryCredentialModalRef = ref(null)
const openCreateImageRegistryCredentialModal = computed(
  () => createImageRegistryCredentialModalRef.value?.openModal ?? (() => {})
)

// Chose Other Docker Configuration
const chooseOtherDockerConfigurationModalRef = ref(null)
const openChooseOtherDockerConfigurationModal = computed(
  () => chooseOtherDockerConfigurationModalRef.value?.openModal ?? (() => {})
)
</script>

<template>
  <!--  Modals -->
  <CreateGitCredentialModal ref="createGitCredentialModalRef" :callback-on-create="refetchGitCredentialList" />
  <CreateImageRegistryCredentialModal
    ref="createImageRegistryCredentialModalRef"
    :callback-on-create="refetchImageRegistryCredentialList" />
  <ChooseOtherDockerConfigurationModal
    ref="chooseOtherDockerConfigurationModalRef"
    :on-apply-configuration="updateDockerConfiguration" />

  <div :key="2" class="mb-5 flex w-full flex-row justify-between p-6 pt-0">
    <div class="w-1/2 max-w-md">
      <!--  Git as Source  -->
      <div v-if="applicationSourceType === 'git'" class="w-full">
        <p class="text-xl font-medium">Git Repository Information</p>

        <!-- Git Credentials -->
        <div class="mt-6">
          <label class="block text-sm font-medium text-gray-700" for="git_credential"
            >Pick Git Credential (Optional)</label
          >
          <div class="mt-1">
            <select
              id="git_credential"
              v-model="stateRef.gitCredentialID"
              @change="fetchGitBranches"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              <option selected value="0">No Credential</option>
              <option v-for="credential in gitCredentials" :key="credential.id" :value="credential.id">
                {{ credential.name }} [{{ credential.type }}]
              </option>
            </select>
          </div>
          <p class="mt-2 flex items-center text-sm">
            Need to add credential for private repo ?
            <a @click="openCreateGitCredentialModal" class="ml-1.5 cursor-pointer font-bold text-primary-600"
              >Click Here</a
            >
          </p>
        </div>

        <!-- Git Repository URL -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="git_repo_url">Git Repository URL</label>
          <div class="mt-1">
            <input
              id="git_repo_url"
              v-model="stateRef.gitRepoUrl"
              v-debounce:1000ms="fetchGitBranches"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="name"
              placeholder="Enter Git Repository URL"
              type="text" />
          </div>
        </div>

        <!-- Git Branch -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="name"
            >Git Branch
            <span class="ml-2 italic" v-if="fetchingGitBranches"
              ><font-awesome-icon icon="fa-solid fa-spinner" class="animate-spin" />&nbsp;&nbsp;Fetching...</span
            ></label
          >
          <div class="mt-1">
            <select
              id="git_credential"
              v-model="stateRef.gitBranch"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              <option selected disabled value="">Select Branch</option>
              <option v-for="branch in availableGitBranches" :key="branch" :value="branch">
                {{ branch }}
              </option>
            </select>
          </div>
        </div>

        <!-- Code Path -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="name">Code Path</label>
          <div class="mt-1">
            <input
              id="name"
              v-model="stateRef.codePath"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="name"
              placeholder="Absolute path of code (optional)"
              type="text" />
            <p class="mt-1 text-xs text-gray-800">
              * You need to specify this if your code is not in root directory of git
            </p>
          </div>
        </div>
      </div>

      <!--  File upload source  -->
      <div v-else-if="applicationSourceType === 'sourceCode'" class="w-full">
        <p class="text-xl font-medium">Upload Source Code</p>
        <!--    Source Code -->
        <div class="mt-4">
          <label class="mb-2 block text-sm font-medium text-gray-900 dark:text-white" for="source_code"
            >Select Folder</label
          >
          <div class="mx-auto max-w-md space-y-8">
            <input
              ref="sourceCodeCompressedFileFieldRef"
              class="w-full cursor-pointer rounded-md bg-gray-100 text-sm text-black file:mr-4 file:cursor-pointer file:border-0 file:bg-gray-800 file:px-4 file:py-2 file:text-white file:hover:bg-gray-700 focus:outline-none"
              directory
              multiple
              type="file"
              webkitdirectory />
          </div>
        </div>

        <!-- Upload Code -->
        <FilledButton
          :loading="stateRef.isUploadingSourceCode"
          class="mt-4 w-full"
          type="secondary"
          @click="uploadSourceCode"
          >Upload Code
        </FilledButton>
      </div>
      <!--  Docker Source  -->
      <div v-else-if="applicationSourceType === 'image'" class="w-full">
        <!-- Docker Image URL-->
        <div class="mt-6">
          <label class="block text-sm font-medium text-gray-700" for="docker_image"
            >Docker Image <span class="text-red-600"> *</span>
          </label>
          <div class="mt-1">
            <input
              id="docker_image"
              v-model="stateRef.dockerImage"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="name"
              placeholder="Enter Docker Image URL"
              type="text" />
          </div>
        </div>
        <!-- Image Registry Credentials -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="image_registry_credential"
            >Pick Image Registry Credential (Optional)
          </label>
          <div class="mt-1">
            <select
              id="image_registry_credential"
              v-model="stateRef.imageRegistryCredentialID"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              <option selected value="0">No Credential</option>
              <option v-for="credential in imageRegistryCredentials" :key="credential.id" :value="credential.id">
                {{ credential.username }} - {{ credential.url }}
              </option>
            </select>
          </div>
          <p class="mt-2 flex items-center text-sm">
            Need to add credential for private registry ?
            <a @click="openCreateImageRegistryCredentialModal" class="ml-1.5 cursor-pointer font-bold text-primary-600"
              >Click Here</a
            >
          </p>
        </div>
      </div>

      <FilledButton
        :disabled="!enableGenerateConfigurationButton"
        :loading="dockerConfigGeneratorGenerating"
        class="mt-6 w-full"
        type="primary"
        @click="generateConfiguration"
        >Re-Generate Configuration
      </FilledButton>
    </div>

    <!-- just for padding purpose -->
    <div></div>

    <div v-if="stateRef.isDockerConfigurationGenerated" class="w-1/2 max-w-md">
      <p class="text-xl font-medium">Generated Configuration</p>
      <FilledButton
        class="mt-6 w-full"
        slim
        type="secondary"
        v-if="applicationSourceType !== 'image'"
        :click="openChooseOtherDockerConfigurationModal"
        >If detected service is incorrect, Click to change the configuration
      </FilledButton>
      <p class="mt-4 font-medium text-gray-700">
        🏂 Detected Service Name -
        <span class="font-normal text-primary-600">{{ stateRef.detectedServiceName }}</span>
      </p>
      <FilledButton v-if="applicationSourceType !== 'image'" class="mt-4 w-full" @click="openDockerFileEditor"
        >View / Modify Dockerfile
      </FilledButton>
      <!-- Docker Command-->
      <div class="mt-4">
        <label class="block text-sm font-medium text-gray-700" for="docker_command"
          >Docker Image Command (Optional)
        </label>
        <div class="mt-1">
          <input
            id="docker_command"
            v-model="stateRef.command"
            autocomplete="off"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            name="docker_command"
            placeholder="Enter Docker Command"
            type="text" />
          <p class="mt-1 text-xs text-gray-800">* It's just to override the default command of docker image</p>
        </div>
      </div>
      <div v-if="stateRef.dockerBuildArgs.length !== 0">
        <p class="mt-4 font-medium text-gray-700">🐳 Docker Build Args</p>
        <div class="w-full">
          <BuildArgInput
            v-for="buildArg in stateRef.dockerBuildArgs"
            :key="buildArg.key"
            :arg-key="buildArg.key"
            :description="buildArg.description"
            :update-build-arg="(val) => updateBuildArg(buildArg.key, val)"
            :value="stateRef.buildArgs[buildArg.key]" />
        </div>
      </div>
    </div>

    <!-- Dockerfile Editor -->
    <DockerfileEditor
      :close-modal="closeDockerFileEditor"
      :code="stateRef.dockerFile"
      :docker-configuration-generating="dockerConfigGeneratorGenerating"
      :is-open="stateRef.isDockerFileEditorOpen"
      :submit="generateConfigurationForCustomDockerFile" />
  </div>
</template>

<style scoped></style>
