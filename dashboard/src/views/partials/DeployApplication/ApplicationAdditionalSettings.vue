<script setup>
import { TabPanel } from '@headlessui/vue'
import Switch from '@/views/components/Switch.vue'
import { computed, reactive, ref } from 'vue'
import EnvironmentVariablesEditor from '@/views/partials/DeployApplication/EnvironmentVariablesEditor.vue'
import { v4 as uuidv4 } from 'uuid'
import PersistentVolumeBindingEditor from '@/views/partials/DeployApplication/PersistentVolumeBindingEditor.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import Disclosure from '@/views/components/Disclosure.vue'
import ConfigMountsEditor from '@/views/partials/DeployApplication/ConfigMountsEditor.vue'
import ConfigureDeploymentPreferredServers from '@/views/partials/ConfigureDeploymentPreferredServers.vue'
import DockerProxyPermissionChoose from '@/views/partials/DockerProxyPermissionChoose.vue'
import { preventSpaceInput } from '@/vendor/utils'

const props = defineProps({
  finalizeApplicationAdditionalSettingsAndDeploy: {
    type: Function,
    required: true
  },
  isDeployRequestSubmitting: {
    required: true
  }
})

const stateRef = reactive({
  replicas: 1,
  deploymentStrategy: 'replicated',
  hostname: '',
  // replicate -> false, global -> true,
  environmentVariablesKeys: [],
  environmentVariablesMap: {},
  persistentVolumeBindingKeys: [],
  persistentVolumeBindingsMap: {},
  configMountKeys: [],
  configMountsMap: {},
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

// Environment Variables Functions
const addEnvironmentVariable = () => {
  const key = uuidv4()
  stateRef.environmentVariablesKeys.push(key)
  stateRef.environmentVariablesMap[key] = {
    name: '',
    value: ''
  }
}

const deleteEnvironmentVariable = (key) => {
  let keys
  keys = stateRef.environmentVariablesKeys.filter((k) => k !== key)
  stateRef.environmentVariablesKeys = keys
  delete stateRef.environmentVariablesMap[key]
}

const changeDeploymentStrategy = (switchStatus) => {
  if (switchStatus) {
    stateRef.deploymentStrategy = 'global'
    stateRef.replicas = 0
  } else {
    stateRef.deploymentStrategy = 'replicated'
    stateRef.replicas = 1
  }
}

const onVariableNameChange = (key, value) => {
  stateRef.environmentVariablesMap[key].name = value
}

const onVariableValueChange = (key, value) => {
  stateRef.environmentVariablesMap[key].value = value
}

// Persistent Volume Binding Functions
const addPersistentVolumeBinding = () => {
  const key = uuidv4()
  stateRef.persistentVolumeBindingKeys.push(key)
  stateRef.persistentVolumeBindingsMap[key] = {
    persistentVolumeID: -1,
    mountingPath: ''
  }
}

const deletePersistentVolumeBinding = (key) => {
  let keys
  keys = stateRef.persistentVolumeBindingKeys.filter((k) => k !== key)
  stateRef.persistentVolumeBindingKeys = keys
  delete stateRef.persistentVolumeBindingsMap[key]
}

const onPersistentVolumeChange = (key, value) => {
  stateRef.persistentVolumeBindingsMap[key].persistentVolumeID = value
}

const onMountingPathChange = (key, value) => {
  stateRef.persistentVolumeBindingsMap[key].mountingPath = value
}

// Config Mount Functions
const addConfigMount = (details) => {
  const mountingPath = details.mountingPath
  // check if mounting path is already used
  for (const key in stateRef.configMountKeys) {
    if (stateRef.configMountsMap[key].mountingPath === mountingPath) {
      throw new Error('Mounting path already used')
    }
  }
  const key = uuidv4()
  stateRef.configMountKeys.push(key)
  stateRef.configMountsMap[key] = details
}

const deleteConfigMount = (key) => {
  let keys
  keys = stateRef.configMountKeys.filter((k) => k !== key)
  stateRef.configMountKeys = keys
  delete stateRef.configMountsMap[key]
}

const onConfigMountContentChange = (key, value) => {
  stateRef.configMountsMap[key].content = value
}

const submitDetails = () => {
  let environmentVariables = []
  for (let key in stateRef.environmentVariablesMap) {
    environmentVariables.push({
      key: stateRef.environmentVariablesMap[key].name,
      value: stateRef.environmentVariablesMap[key].value
    })
  }
  let details = {
    deploymentMode: stateRef.deploymentStrategy,
    replicas: stateRef.replicas,
    hostname: stateRef.hostname,
    environmentVariables: environmentVariables,
    persistentVolumeBindings: Object.values(stateRef.persistentVolumeBindingsMap),
    configMounts: Object.values(stateRef.configMountsMap),
    preferredServerHostnames: Object.values(stateRef.preferredServerHostnames),
    dockerProxyConfig: JSON.parse(JSON.stringify(stateRef.dockerProxyConfig)),
    customHealthCheck: JSON.parse(JSON.stringify(stateRef.customHealthCheck))
  }
  props.finalizeApplicationAdditionalSettingsAndDeploy(details)
}

const preferredServerHostnamesStr = computed(() => {
  return stateRef.preferredServerHostnames.join(', ')
})

const configureDeploymentPreferredServersRef = ref(null)
const openConfigureDeploymentPreferredServers = () => {
  if (configureDeploymentPreferredServersRef.value === null) {
    return
  }
  configureDeploymentPreferredServersRef.value.openModal()
}
</script>

<template>
  <TabPanel :key="3" class="flex w-full flex-col p-6">
    <!-- Environment Variables -->
    <p class="mt-6 text-base font-medium text-gray-900">Environment Variables</p>
    <EnvironmentVariablesEditor
      :add-environment-variable="addEnvironmentVariable"
      :delete-environment-variable="deleteEnvironmentVariable"
      :environment-variables-keys="stateRef.environmentVariablesKeys"
      :environment-variables-map="stateRef.environmentVariablesMap"
      :on-variable-name-change="onVariableNameChange"
      :on-variable-value-change="onVariableValueChange"
      class="mt-2" />
    <!-- Persistent Volumes -->
    <p class="mb-2 mt-6 text-base font-medium text-gray-900">Persistent Volumes</p>
    <PersistentVolumeBindingEditor
      :add-persistent-volume-binding="addPersistentVolumeBinding"
      :delete-persistent-volume-binding="deletePersistentVolumeBinding"
      :on-mounting-path-change="onMountingPathChange"
      :on-persistent-volume-change="onPersistentVolumeChange"
      :persistent-volume-binding-keys="stateRef.persistentVolumeBindingKeys"
      :persistent-volume-bindings-map="stateRef.persistentVolumeBindingsMap"
      class="mt-2" />
    <!-- Advanced Settings -->
    <Disclosure class="my-6">
      <template v-slot:title>Advanced Settings (Click to expand)</template>
      <template v-slot:body>
        <!-- Deployment Configuration -->
        <div class="mt-3 flex flex-row items-center">
          <p class="font-medium text-black">Deployment Strategy</p>
          <font-awesome-icon class="px-4" icon="fa-solid fa-arrow-right" />
          <div class="flex flex-row items-center gap-2">
            <p class="font-medium">Replicated</p>
            <input
              v-if="stateRef.deploymentStrategy === 'replicated'"
              class="block h-8 w-16 rounded-full border-gray-300 shadow-sm [appearance:textfield] focus:border-primary-500 focus:ring-primary-500 [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
              name="no_of_replicase"
              placeholder="No of Replicas"
              type="number"
              v-model="stateRef.replicas" />
            <p v-if="stateRef.deploymentStrategy === 'replicated'">replica(s)</p>
          </div>
          <Switch
            :enabled="stateRef.deploymentStrategy === 'global'"
            :onChange="changeDeploymentStrategy"
            class="mx-4" />
          <p class="font-medium">Global</p>
        </div>
        <!--    Deployment Preferred Server    -->
        <ConfigureDeploymentPreferredServers
          ref="configureDeploymentPreferredServersRef"
          :update-hostnames="(e) => (stateRef.preferredServerHostnames = e)"
          :hostnames="stateRef.preferredServerHostnames" />
        <!-- Preferred Servers  -->
        <div class="mt-3">
          <p class="font-medium text-black">Deployment Preferred Servers</p>
          <div class="mt-1">
            <input
              class="mt-1 block w-full cursor-pointer rounded-md border-gray-300 text-sm shadow-sm focus:border-primary-500 focus:ring-primary-500"
              placeholder="Click to add servers"
              type="text"
              @click="openConfigureDeploymentPreferredServers"
              v-model="preferredServerHostnamesStr"
              readonly />
          </div>
        </div>
        <!--    Config mounts    -->
        <div class="mt-3 w-full">
          <p class="mb-2 text-sm font-medium">Config Mounts</p>
          <ConfigMountsEditor
            :config-mounts-keys="stateRef.configMountKeys"
            :config-mounts-map="stateRef.configMountsMap"
            :add-config-mount="addConfigMount"
            :delete-config-mount="deleteConfigMount"
            :on-config-content-change="onConfigMountContentChange" />
        </div>
        <!--   Container Hostname    -->
        <div class="mt-3">
          <p class="font-medium text-black">Container Hostname</p>
          <div class="mt-1">
            <input
              class="mt-1 block w-full rounded-md border-gray-300 text-sm shadow-sm focus:border-primary-500 focus:ring-primary-500"
              placeholder="Provide container hostname"
              type="text"
              v-model="stateRef.hostname"
              @keydown="preventSpaceInput" />
          </div>
        </div>
        <!--   Healthcheck     -->
        <div class="mt-3 w-full">
          <div class="flex flex-row items-center gap-2">
            <p class="text-sm font-medium">Custom Healthcheck</p>
            <div class="multi-select">
              <div
                @click="() => (stateRef.customHealthCheck.enabled = true)"
                :class="{
                  active: stateRef.customHealthCheck.enabled
                }">
                Enabled
              </div>
              <div
                @click="() => (stateRef.customHealthCheck.enabled = false)"
                :class="{
                  active: !stateRef.customHealthCheck.enabled
                }">
                Disabled
              </div>
            </div>
          </div>
          <div class="mt-2" v-if="stateRef.customHealthCheck.enabled">
            <label class="block text-sm font-medium"
              >Healthcheck Test Command<span class="text-red-600"> *</span>
            </label>
            <div class="mt-1">
              <input
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                type="text"
                @input="(e) => (stateRef.customHealthCheck.test_command = e.target.value) || true"
                :value="stateRef.customHealthCheck.test_command" />
            </div>
          </div>
          <div class="mt-3 flex w-full flex-row gap-5" v-if="stateRef.customHealthCheck.enabled">
            <div class="w-1/5">
              <label class="block text-sm font-medium"
                >Check Interval (Seconds)<span class="text-red-600"> *</span>
              </label>
              <div class="mt-1">
                <input
                  autocomplete="off"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  type="number"
                  @change="(e) => (stateRef.customHealthCheck.interval_seconds = parseInt(e.target.value) || 0)"
                  :value="stateRef.customHealthCheck.interval_seconds" />
              </div>
            </div>
            <div class="w-1/5">
              <label class="block text-sm font-medium"
                >Check Timeout (Seconds)<span class="text-red-600"> *</span>
              </label>
              <div class="mt-1">
                <input
                  autocomplete="off"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  type="number"
                  @change="(e) => (stateRef.customHealthCheck.timeout_seconds = parseInt(e.target.value) || 0)"
                  :value="stateRef.customHealthCheck.timeout_seconds" />
              </div>
            </div>
            <div class="w-1/5">
              <label class="block text-sm font-medium"
                >Start Period (Seconds)<span class="text-red-600"> *</span>
              </label>
              <div class="mt-1">
                <input
                  autocomplete="off"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  type="number"
                  @change="(e) => (stateRef.customHealthCheck.start_period_seconds = parseInt(e.target.value) || 0)"
                  :value="stateRef.customHealthCheck.start_period_seconds" />
              </div>
            </div>
            <div class="w-1/5">
              <label class="block text-sm font-medium"
                >Start Interval (Seconds)<span class="text-red-600"> *</span>
              </label>
              <div class="mt-1">
                <input
                  autocomplete="off"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  type="number"
                  @change="(e) => (stateRef.customHealthCheck.start_interval_seconds = parseInt(e.target.value) || 0)"
                  :value="stateRef.customHealthCheck.start_interval_seconds" />
              </div>
            </div>
            <div class="w-1/5">
              <label class="block text-sm font-medium">Retries<span class="text-red-600"> *</span> </label>
              <div class="mt-1">
                <input
                  autocomplete="off"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  type="number"
                  @change="(e) => (stateRef.customHealthCheck.retries = parseInt(e.target.value) || 0)"
                  :value="stateRef.customHealthCheck.retries" />
              </div>
            </div>
          </div>
        </div>
        <!--   Docker Proxy     -->
        <div class="mt-5">
          <!-- Proxy Status   -->
          <div class="flex flex-row gap-2">
            <p class="font-medium">Docker Proxy Status</p>
            <div class="multi-select">
              <div
                @click="stateRef.dockerProxyConfig.enabled = true"
                :class="{
                  active: stateRef.dockerProxyConfig.enabled
                }">
                Enabled
              </div>
              <div
                @click="stateRef.dockerProxyConfig.enabled = false"
                :class="{
                  active: !stateRef.dockerProxyConfig.enabled
                }">
                Disabled
              </div>
            </div>
          </div>
          <p class="my-2 italic">
            <span class="text-red-600">* </span>Don't enable <b>Docker Proxy</b> if your application doesn't need access
            to the Docker socket and it increase attack surface of your serer if not configured and used properly.
          </p>
          <!--  Proxy Permission  -->
          <div
            class="mt-4 flex w-full flex-row gap-20 rounded-md border-2 border-secondary-300 p-2"
            v-if="stateRef.dockerProxyConfig.enabled">
            <div class="flex w-1/3 flex-col gap-2">
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Ping</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.ping"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.ping = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Version</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.version"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.version = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Info</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.info"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.info = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Events</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.events"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.events = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Auth</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.auth"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.auth = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Secrets</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.secrets"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.secrets = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Build</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.build"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.build = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Commit</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.commit"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.commit = value)" />
              </div>
            </div>
            <div class="flex w-1/3 flex-col gap-2">
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Configs</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.configs"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.configs = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Containers</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.containers"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.containers = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Distribution</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.distribution"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.distribution = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Exec</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.exec"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.exec = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Grpc</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.grpc"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.grpc = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Images</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.images"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.images = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Networks</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.networks"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.networks = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Nodes</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.nodes"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.nodes = value)" />
              </div>
            </div>
            <div class="flex w-1/3 flex-col gap-2">
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Plugins</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.plugins"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.plugins = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Services</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.services"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.services = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Session</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.session"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.session = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Swarm</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.swarm"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.swarm = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">System</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.system"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.system = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Tasks</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.tasks"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.tasks = value)" />
              </div>
              <div class="flex flex-row justify-between gap-2">
                <p class="w-1/3 font-normal text-gray-800">Volumes</p>
                <DockerProxyPermissionChoose
                  :value="stateRef.dockerProxyConfig.permission.volumes"
                  :on-change="(value) => (stateRef.dockerProxyConfig.permission.volumes = value)" />
              </div>
            </div>
          </div>
          <div class="mt-4" v-if="stateRef.dockerProxyConfig.enabled">
            <p class="font-medium">Docker Socket Proxy Usage Guide</p>
            <p class="mt-2">
              1. Use <span class="font-semibold" v-html="`{{DOCKER_PROXY_HOST}}`"></span> as value of environment
              variable. While deploying app, swiftwave will inject the docker proxy host.
            </p>
          </div>
        </div>
      </template>
    </Disclosure>
    <!-- Proceed to next -->
    <div class="mt-6 flex flex-row justify-end">
      <FilledButton type="primary" @click="submitDetails" :loading="isDeployRequestSubmitting">
        <font-awesome-icon icon="fa-solid fa-hammer" class="mr-2" />
        Confirm & Deploy Application
      </FilledButton>
    </div>
  </TabPanel>
</template>

<style scoped>
.multi-select {
  @apply flex h-fit w-min overflow-hidden rounded-md border border-secondary-400;

  div {
    @apply cursor-pointer border-r border-secondary-400 bg-secondary-100 px-2.5 py-0.5 text-sm transition-all;
  }

  div:last-child {
    @apply border-0;
  }

  .active {
    @apply bg-primary-600  text-white;
  }
}
</style>
