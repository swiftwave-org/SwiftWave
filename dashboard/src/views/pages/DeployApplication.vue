<script setup>
import { reactive, ref } from 'vue'
import { Tab, TabGroup, TabList, TabPanels } from '@headlessui/vue'
import ApplicationNameSelection from '@/views/partials/DeployApplication/ApplicationNameSelection.vue'
import ApplicationSourceSelection from '@/views/partials/DeployApplication/ApplicationSourceSelection.vue'
import ApplicationSourceConfiguration from '@/views/partials/DeployApplication/ApplicationSourceConfiguration.vue'
import ApplicationAdditionalSettings from '@/views/partials/DeployApplication/ApplicationAdditionalSettings.vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'
import ModalDialog from '@/views/components/ModalDialog.vue'

const router = useRouter()
const sectionNames = ['Application Name', 'Select Source', 'Application Source', 'Deploy Configuration']
const isApplicationDeployedSuccessfulModalOpen = ref(false)
const selectedTabIndex = ref(0)
const changeTab = (index) => {
  selectedTabIndex.value = index
}

// state
const newApplicationState = reactive({
  name: '',
  upstreamType: '',
  hostname: '',
  command: '',
  deploymentMode: '',
  replicas: 0,
  dockerfile: '',
  buildArgs: [],
  environmentVariables: [],
  persistentVolumeBindings: [],
  configMounts: [],
  gitCredentialID: 0,
  repositoryUrl: '',
  repositoryBranch: '',
  codePath: '',
  imageRegistryCredentialID: 0,
  dockerImage: '',
  sourceCodeCompressedFileName: '',
  capabilities: [],
  sysctls: [],
  resourceLimit: {
    memoryMb: 0
  },
  reservedResource: {
    memoryMb: 0
  },
  customHealthCheck: {
    enabled: false,
    test_command: '',
    interval_seconds: 0,
    timeout_seconds: 0,
    start_period_seconds: 0,
    start_interval_seconds: 0,
    retries: 0
  },
  preferredServerHostnames: [],
  dockerProxyConfig: {
    enabled: false,
    permission: {
      ping: 'none',
      version: 'none',
      info: 'none',
      events: 'none',
      auth: 'none',
      secrets: 'none',
      build: 'none',
      commit: 'none',
      configs: 'none',
      containers: 'none',
      distribution: 'none',
      exec: 'none',
      grpc: 'none',
      images: 'none',
      networks: 'none',
      nodes: 'none',
      plugins: 'none',
      services: 'none',
      session: 'none',
      swarm: 'none',
      system: 'none',
      tasks: 'none',
      volumes: 'none'
    }
  }
})

// Deploy application
const {
  mutate: deployApplication,
  loading: isDeployRequestSubmitting,
  onDone: onDeployApplicationMutationDone,
  onError: onDeployApplicationMutationError
} = useMutation(
  gql`
    mutation ($input: ApplicationInput!) {
      createApplication(input: $input) {
        id
        name
        latestDeployment {
          id
        }
      }
    }
  `,
  {
    variables: {
      input: newApplicationState
    }
  }
)

onDeployApplicationMutationDone((result) => {
  if (result.data.createApplication.latestDeployment === null) {
    toast.warning('Application is not deployed yet, please wait for a while and refresh the page')
    return
  }
  isApplicationDeployedSuccessfulModalOpen.value = true
  setTimeout(() => {
    router.push({
      name: 'Application Details Deployments',
      params: {
        id: result.data.createApplication.id,
        deployment_id: result.data.createApplication.latestDeployment.id
      }
    })
  }, 2000)
})

onDeployApplicationMutationError((err) => {
  toast.error(err.message)
})

// functions
const finalApplicationNameAndMoveToNextTab = (name) => {
  newApplicationState.name = name
  changeTab(1)
}

const finalizeApplicationSourceAndMoveToNextTab = (upstreamType) => {
  newApplicationState.upstreamType = upstreamType
  changeTab(2)
}

const finalizeApplicationSourceConfigurationAndMoveToNextTab = (configuration) => {
  // Store the configuration in the state
  // NOTE: Don't modify as configuration is a reference to the state of `ApplicationSourceConfiguration.vue`
  newApplicationState.dockerfile = configuration.dockerFile
  let buildArgs = []
  for (let key in configuration.buildArgs) {
    buildArgs.push({
      key: key,
      value: configuration.buildArgs[key]
    })
  }
  newApplicationState.command = configuration.command
  newApplicationState.buildArgs = buildArgs
  newApplicationState.gitCredentialID = parseInt(configuration.gitCredentialID)
  newApplicationState.imageRegistryCredentialID = parseInt(configuration.imageRegistryCredentialID)
  newApplicationState.gitCredentialID = configuration.gitCredentialID === 0 ? null : configuration.gitCredentialID
  newApplicationState.imageRegistryCredentialID =
    configuration.imageRegistryCredentialID === 0 ? null : configuration.imageRegistryCredentialID
  newApplicationState.repositoryUrl = configuration.gitRepoUrl
  newApplicationState.repositoryBranch = configuration.gitBranch
  newApplicationState.codePath = configuration.codePath
  newApplicationState.sourceCodeCompressedFileName = configuration.sourceCodeFile
  newApplicationState.dockerImage = configuration.dockerImage
  changeTab(3)
}

const finalizeApplicationAdditionalSettings = (additionalSettings) => {
  // Store the configuration in the state
  // NOTE: Don't modify as configuration is a reference to the state of `ApplicationAdditionalSettings.vue`
  newApplicationState.deploymentMode = additionalSettings.deploymentMode
  newApplicationState.replicas = additionalSettings.replicas
  newApplicationState.hostname = additionalSettings.hostname
  newApplicationState.environmentVariables = additionalSettings.environmentVariables
  newApplicationState.persistentVolumeBindings = additionalSettings.persistentVolumeBindings
  newApplicationState.configMounts = additionalSettings.configMounts
  newApplicationState.preferredServerHostnames = additionalSettings.preferredServerHostnames
  newApplicationState.dockerProxyConfig.enabled = additionalSettings.dockerProxyConfig.enabled
  newApplicationState.dockerProxyConfig.permission.ping = additionalSettings.dockerProxyConfig.permission.ping
  newApplicationState.dockerProxyConfig.permission.version = additionalSettings.dockerProxyConfig.permission.version
  newApplicationState.dockerProxyConfig.permission.info = additionalSettings.dockerProxyConfig.permission.info
  newApplicationState.dockerProxyConfig.permission.events = additionalSettings.dockerProxyConfig.permission.events
  newApplicationState.dockerProxyConfig.permission.auth = additionalSettings.dockerProxyConfig.permission.auth
  newApplicationState.dockerProxyConfig.permission.secrets = additionalSettings.dockerProxyConfig.permission.secrets
  newApplicationState.dockerProxyConfig.permission.build = additionalSettings.dockerProxyConfig.permission.build
  newApplicationState.dockerProxyConfig.permission.commit = additionalSettings.dockerProxyConfig.permission.commit
  newApplicationState.dockerProxyConfig.permission.configs = additionalSettings.dockerProxyConfig.permission.configs
  newApplicationState.dockerProxyConfig.permission.containers =
    additionalSettings.dockerProxyConfig.permission.containers
  newApplicationState.dockerProxyConfig.permission.distribution =
    additionalSettings.dockerProxyConfig.permission.distribution
  newApplicationState.dockerProxyConfig.permission.exec = additionalSettings.dockerProxyConfig.permission.exec
  newApplicationState.dockerProxyConfig.permission.grpc = additionalSettings.dockerProxyConfig.permission.grpc
  newApplicationState.dockerProxyConfig.permission.images = additionalSettings.dockerProxyConfig.permission.images
  newApplicationState.dockerProxyConfig.permission.networks = additionalSettings.dockerProxyConfig.permission.networks
  newApplicationState.dockerProxyConfig.permission.nodes = additionalSettings.dockerProxyConfig.permission.nodes
  newApplicationState.dockerProxyConfig.permission.plugins = additionalSettings.dockerProxyConfig.permission.plugins
  newApplicationState.dockerProxyConfig.permission.services = additionalSettings.dockerProxyConfig.permission.services
  newApplicationState.dockerProxyConfig.permission.session = additionalSettings.dockerProxyConfig.permission.session
  newApplicationState.dockerProxyConfig.permission.swarm = additionalSettings.dockerProxyConfig.permission.swarm
  newApplicationState.dockerProxyConfig.permission.system = additionalSettings.dockerProxyConfig.permission.system
  newApplicationState.dockerProxyConfig.permission.tasks = additionalSettings.dockerProxyConfig.permission.tasks
  newApplicationState.dockerProxyConfig.permission.volumes = additionalSettings.dockerProxyConfig.permission.volumes
  newApplicationState.customHealthCheck.enabled = additionalSettings.customHealthCheck.enabled
  newApplicationState.customHealthCheck.test_command = additionalSettings.customHealthCheck.test_command
  newApplicationState.customHealthCheck.interval_seconds = additionalSettings.customHealthCheck.interval_seconds
  newApplicationState.customHealthCheck.timeout_seconds = additionalSettings.customHealthCheck.timeout_seconds
  newApplicationState.customHealthCheck.start_period_seconds = additionalSettings.customHealthCheck.start_period_seconds
  newApplicationState.customHealthCheck.start_interval_seconds =
    additionalSettings.customHealthCheck.start_interval_seconds
  newApplicationState.customHealthCheck.retries = additionalSettings.customHealthCheck.retries
}

const finalizeApplicationAdditionalSettingsAndDeploy = (additionalSettings) => {
  finalizeApplicationAdditionalSettings(additionalSettings)
  deployApplication()
}

const onClickTab = (index) => {
  if (index < selectedTabIndex.value) {
    alert('If you want to change the previous configuration, you need to refresh the page and start over')
  }
}
</script>

<template>
  <ModalDialog :is-open="isApplicationDeployedSuccessfulModalOpen" non-cancelable>
    <template v-slot:header>
      <span>🚀 Application Deployment in Progress</span>
    </template>
    <template v-slot:body>
      <p class="mb-4">Application deployment has been started. Your application will be live shortly.</p>
      <p class="italic">Redirecting to Deployment Page in few seconds</p>
    </template>
  </ModalDialog>
  <div class="flex h-full w-full max-w-7xl flex-col items-center sm:px-0">
    <TabGroup :selected-index="selectedTabIndex">
      <TabList class="flex w-full max-w-4xl space-x-3 rounded-full bg-primary-600 p-1">
        <Tab
          v-for="(sectionName, index) in sectionNames"
          :key="sectionName"
          v-slot="{ selected }"
          as="template"
          @click="() => onClickTab(index)">
          <button :class="selected ? 'tab-button-selected' : 'tab-button-unselected'" class="tab-button">
            {{ sectionName }}
          </button>
        </Tab>
      </TabList>

      <TabPanels class="mt-6 flex h-full w-full flex-col items-center">
        <!-- Application Name Selection -->
        <ApplicationNameSelection :final-application-name-and-move-to-next-tab="finalApplicationNameAndMoveToNextTab" />
        <!-- Source Selection -->
        <ApplicationSourceSelection
          :finalize-application-source-and-move-to-next-tab="finalizeApplicationSourceAndMoveToNextTab" />
        <!--  Source Configuration -->
        <ApplicationSourceConfiguration
          :application-source-type="newApplicationState.upstreamType"
          :finalize-application-source-configuration-and-move-to-next-tab="
            finalizeApplicationSourceConfigurationAndMoveToNextTab
          " />
        <!-- Additional Settings  -->
        <ApplicationAdditionalSettings
          :finalize-application-additional-settings-and-deploy="finalizeApplicationAdditionalSettingsAndDeploy"
          :is-deploy-request-submitting="isDeployRequestSubmitting" />
      </TabPanels>
    </TabGroup>
  </div>
</template>

<style scoped>
.tab-button {
  @apply w-full rounded-full px-3 py-2 text-sm font-medium leading-5 focus:outline-none;
}

.tab-button-selected {
  @apply bg-gray-100 text-gray-900 shadow;
}

.tab-button-unselected {
  @apply text-gray-200 hover:bg-white/[0.1] hover:text-white;
}

.tab-panel {
  @apply mt-5;
}
</style>
