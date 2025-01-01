<script setup>
// Toast
import { useMutation, useQuery } from '@vue/apollo-composable'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import gql from 'graphql-tag'
import { computed, reactive, ref } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import Table from '@/views/components/Table/Table.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import TableMessage from '@/views/components/Table/TableMessage.vue'
import ApplicationListRow from '@/views/partials/ApplicationListRow.vue'
import DeleteApplicationsModal from '@/views/partials/DeleteApplicationsModal.vue'
import RestartApplicationsModal from '@/views/partials/RestartApplicationsModal.vue'
import RebuildApplicationsModal from '@/views/partials/RebuildApplicationsModal.vue'
import EnvironmentVariablesEditor from '@/views/partials/DeployApplication/EnvironmentVariablesEditor.vue'
import { v4 as uuidv4 } from 'uuid'
import PersistentVolumeBindingEditor from '@/views/partials/DeployApplication/PersistentVolumeBindingEditor.vue'
import ConfigMountsEditor from '@/views/partials/DeployApplication/ConfigMountsEditor.vue'
import FilledButton from '@/views/components/FilledButton.vue'

// Get the application ID from the URL
const router = useRouter()
const applicationGroupId = router.currentRoute.value.params.id

// Fetch the application details
const {
  result: applicationGroupDetailsRaw,
  loading: applicationGroupDetailsLoading,
  refetch: refetchGroupApplicationDetails,
  onResult: onResultGroupApplicationDetails,
  onError: onErrorGroupApplicationDetails
} = useQuery(
  gql`
    query ($id: String!) {
      applicationGroup(id: $id) {
        id
        name
        logo
        applications {
          id
          name
          deploymentMode
          command
          replicas
          isSleeping
          resourceLimit {
            memoryMb
          }
          reservedResource {
            memoryMb
          }
          environmentVariables {
            key
            value
          }
          persistentVolumeBindings {
            id
            persistentVolumeID
            mountingPath
          }
          configMounts {
            content
            mountingPath
            uid
            gid
          }
          latestDeployment {
            upstreamType
            dockerfile
            buildArgs {
              key
              value
            }
            gitEndpoint
            gitProvider
            gitCredentialID
            repositoryUrl
            repositoryBranch
            codePath
            imageRegistryCredentialID
            dockerImage
            sourceCodeCompressedFileName
          }
          realtimeInfo {
            InfoFound
            DeploymentMode
            DesiredReplicas
            RunningReplicas
            HealthStatus
          }
          ingressRules {
            domain {
              name
            }
            protocol
            port
          }
          capabilities
          sysctls
          applicationGroupID
          customHealthCheck {
            enabled
            test_command
            interval_seconds
            timeout_seconds
            start_period_seconds
            start_interval_seconds
            retries
          }
          preferredServerHostnames
          dockerProxyConfig {
            enabled
            permission {
              ping
              version
              info
              events
              auth
              secrets
              build
              commit
              configs
              containers
              distribution
              exec
              grpc
              images
              networks
              nodes
              plugins
              services
              session
              swarm
              system
              tasks
              volumes
            }
          }
        }
      }
    }
  `,
  {
    id: applicationGroupId
  },
  {
    pollInterval: 30000
  }
)

onErrorGroupApplicationDetails(() => {
  toast.error('Failed to fetch application group details')
})

onResultGroupApplicationDetails(() => {
  for (const application of applications.value) {
    let envVariablesMap = {}
    application.environmentVariables.forEach((variable) => {
      const z = uuidv4()
      envVariablesMap[z] = {
        name: variable.key,
        value: variable.value
      }
    })
    environmentVariableDetails[application.id] = {
      keys: Object.keys(envVariablesMap),
      map: envVariablesMap
    }
    let persistentVolumeBindingsMap = {}
    application.persistentVolumeBindings.forEach((binding) => {
      const z = uuidv4()
      persistentVolumeBindingsMap[z] = {
        persistentVolumeID: binding.persistentVolumeID,
        mountingPath: binding.mountingPath
      }
    })
    persistentVolumeBindingDetails[application.id] = {
      keys: Object.keys(persistentVolumeBindingsMap),
      map: persistentVolumeBindingsMap
    }
    let configMountsMap = {}
    application.configMounts.forEach((configMount) => {
      const z = uuidv4()
      configMountsMap[z] = {
        content: configMount.content,
        mountingPath: configMount.mountingPath,
        uid: configMount.uid,
        gid: configMount.gid
      }
    })
    configMountDetails[application.id] = {
      keys: Object.keys(configMountsMap),
      map: configMountsMap
    }
    isAppInfoChanged[application.id] = false
  }
  if (applications.value.length > 0 && pageInfo.currentSelectedEnvironmentVariableApplicationId === '') {
    pageInfo.currentSelectedEnvironmentVariableApplicationId = applications.value[0].id
    pageInfo.currentSelectedPersistentVolumeApplicationId = applications.value[0].id
    pageInfo.currentSelectedConfigMountApplicationId = applications.value[0].id
  }
})

const persistentVolumeBindingDetails = reactive({})
const environmentVariableDetails = reactive({})
const configMountDetails = reactive({})

const applicationGroupDetails = computed(() => applicationGroupDetailsRaw.value?.applicationGroup ?? {})
const applications = computed(() => applicationGroupDetailsRaw.value?.applicationGroup?.applications ?? [])
const ingressRules = computed(() => {
  let records = []
  for (const application of applications.value) {
    records.push(...application.ingressRules)
  }
  return records
})

const totalServiceCount = computed(() => {
  if (applicationGroupDetails.value.applications.length === 0) {
    return 0
  }
  return applicationGroupDetails.value.applications.length
})

const healthyServiceCount = computed(() => {
  if (applicationGroupDetails.value.applications.length === 0) {
    return 0
  }
  return applicationGroupDetails.value.applications.filter((app) => app.realtimeInfo.HealthStatus === 'healthy').length
})

const unhealthyServiceCount = computed(() => {
  if (applicationGroupDetails.value.applications.length === 0) {
    return 0
  }
  return applicationGroupDetails.value.applications.filter((app) => app.realtimeInfo.HealthStatus === 'unhealthy')
    .length
})

const applicationIds = computed(() => {
  return applicationGroupDetails.value.applications.map((app) => app.id)
})

const deleteApplicationsModal = ref(null)
const restartApplicationsModal = ref(null)
const rebuildApplicationsModal = ref(null)

function deleteApplications() {
  if (deleteApplicationsModal.value) {
    deleteApplicationsModal.value.openModal()
  }
}

function restartApplications() {
  if (restartApplicationsModal.value) {
    restartApplicationsModal.value.openModal()
  }
}

function rebuildApplications() {
  if (rebuildApplicationsModal.value) {
    rebuildApplicationsModal.value.openModal()
  }
}

// page
const pageName = ref('deployed-apps')
const pageInfo = reactive({
  currentSelectedPersistentVolumeApplicationId: '',
  currentSelectedEnvironmentVariableApplicationId: '',
  currentSelectedConfigMountApplicationId: ''
})
const isAppInfoChanged = reactive({})
const isAnyAppInfoChanged = computed(() => {
  for (const app of applications.value) {
    if (isAppInfoChanged[app.id]) {
      return true
    }
  }
  return false
})

// Environment Variables Editor Related
const environmentVariableKeys = (app) => {
  console.log(environmentVariableDetails)
  return environmentVariableDetails[app.id].keys
}
const environmentVariableMap = (app) => {
  return environmentVariableDetails[app.id].map
}
const addEnvironmentVariable = (app) => {
  const key = uuidv4()
  environmentVariableDetails[app.id].map[key] = {
    name: '',
    value: ''
  }
  environmentVariableDetails[app.id].keys.push(key)
  isAppInfoChanged[app.id] = true
}
const deleteEnvironmentVariable = (app, key) => {
  delete environmentVariableDetails[app.id].map[key]
  environmentVariableDetails[app.id].keys = environmentVariableDetails[app.id].keys.filter((k) => k !== key)
  isAppInfoChanged[app.id] = true
}
const onEnvironmentVariableValueChange = (app, key, value) => {
  environmentVariableDetails[app.id].map[key].value = value
  isAppInfoChanged[app.id] = true
}
const onEnvironmentVariableNameChange = (app, key, name) => {
  environmentVariableDetails[app.id].map[key].name = name
  isAppInfoChanged[app.id] = true
}

// Persistent Volume Binding Editor Related
const persistentVolumeBindingKeys = (app) => {
  return persistentVolumeBindingDetails[app.id].keys
}
const persistentVolumeBindingMap = (app) => {
  return persistentVolumeBindingDetails[app.id].map
}
const addPersistentVolumeBinding = (app) => {
  const key = uuidv4()
  persistentVolumeBindingDetails[app.id].map[key] = {
    persistentVolumeID: -1,
    mountingPath: ''
  }
  persistentVolumeBindingDetails[app.id].keys.push(key)
  isAppInfoChanged[app.id] = true
}
const deletePersistentVolumeBinding = (app, key) => {
  delete persistentVolumeBindingDetails[app.id].map[key]
  persistentVolumeBindingDetails[app.id].keys = persistentVolumeBindingDetails[app.id].keys.filter((k) => k !== key)
  isAppInfoChanged[app.id] = true
}
const onPersistentVolumeChange = (app, key, value) => {
  persistentVolumeBindingDetails[app.id].map[key].persistentVolumeID = value
  isAppInfoChanged[app.id] = true
}
const onPersistentVolumeMountingPathChange = (app, key, value) => {
  persistentVolumeBindingDetails[app.id].map[key].mountingPath = value
  isAppInfoChanged[app.id] = true
}

// Config Mount Editor Related
const configMountKeys = (app) => {
  return configMountDetails[app.id].keys
}
const configMountMap = (app) => {
  return configMountDetails[app.id].map
}
const addConfigMount = (app, details) => {
  const key = uuidv4()
  configMountDetails[app.id].map[key] = details
  configMountDetails[app.id].keys.push(key)
  isAppInfoChanged[app.id] = true
}
const deleteConfigMount = (app, key) => {
  delete configMountDetails[app.id].map[key]
  configMountDetails[app.id].keys = configMountDetails[app.id].keys.filter((k) => k !== key)
  isAppInfoChanged[app.id] = true
}
const onConfigMountContentChange = (app, key, content) => {
  configMountDetails[app.id].map[key].content = content
  isAppInfoChanged[app.id] = true
}

// Apply Changes

const isApplyingChanges = ref(false)
const { mutate: deployApplication } = useMutation(gql`
  mutation ($id: String!, $input: ApplicationInput!) {
    updateApplication(id: $id, input: $input) {
      id
      name
    }
  }
`)
const applyChanges = async () => {
  isApplyingChanges.value = true
  for (const appId in isAppInfoChanged) {
    if (!isAppInfoChanged[appId]) continue
    for (const application of applications.value) {
      if (application.id === appId) {
        let updatedPayload = {
          name: application.name,
          upstreamType: application.latestDeployment.upstreamType, // TODO Not allowed to change
          command: application.command,
          deploymentMode: application.deploymentMode,
          replicas: application.replicas,
          resourceLimit: {
            memoryMb: application.resourceLimit.memoryMb
          },
          reservedResource: {
            memoryMb: application.reservedResource.memoryMb
          },
          buildArgs: application.latestDeployment.buildArgs,
          environmentVariables: environmentVariableKeys(application).map((key) => {
            return {
              key: environmentVariableMap(application)[key].name,
              value: environmentVariableMap(application)[key].value
            }
          }),
          configMounts: configMountKeys(application).map((key) => configMountMap(application)[key]),
          persistentVolumeBindings: persistentVolumeBindingKeys(application).map(
            (key) => persistentVolumeBindingMap(application)[key]
          ),
          // update this part
          gitCredentialID:
            application.latestDeployment.gitCredentialID === 0 ? null : application.latestDeployment.gitCredentialID,
          repositoryUrl: application.latestDeployment.repositoryUrl,
          repositoryBranch: application.latestDeployment.repositoryBranch,
          codePath: application.latestDeployment.codePath,
          imageRegistryCredentialID:
            application.latestDeployment.imageRegistryCredentialID === 0
              ? null
              : application.latestDeployment.imageRegistryCredentialID,
          dockerImage: application.latestDeployment.dockerImage,
          sourceCodeCompressedFileName: application.latestDeployment.sourceCodeCompressedFileName,
          dockerfile: application.latestDeployment.dockerfile,
          capabilities: application.capabilities,
          sysctls: application.sysctls,
          applicationGroupID: application.applicationGroupID === 0 ? null : application.applicationGroupID,
          customHealthCheck: application.customHealthCheck,
          preferredServerHostnames: application.preferredServerHostnames,
          dockerProxyConfig: application.dockerProxyConfig
        }
        try {
          await deployApplication({
            input: updatedPayload,
            id: application.id
          })
        } catch (e) {
          toast.error(e.message)
        }
      }
    }
  }
  isApplyingChanges.value = false
  toast.success('Changes applied successfully')
  refetchGroupApplicationDetails()
}
</script>

<template>
  <!-- Main -->
  <div v-if="applicationGroupDetailsLoading">
    <p>Loading...</p>
  </div>
  <section v-else class="mx-auto w-full max-w-7xl">
    <!--  Modals  -->
    <DeleteApplicationsModal ref="deleteApplicationsModal" :application-ids="applicationIds" />
    <RestartApplicationsModal
      ref="restartApplicationsModal"
      :application-ids="applicationIds"
      :on-done="refetchGroupApplicationDetails" />
    <RebuildApplicationsModal
      ref="rebuildApplicationsModal"
      :application-ids="applicationIds"
      :on-done="refetchGroupApplicationDetails" />
    <!--  First line  -->
    <div class="flex w-full flex-row items-center justify-between">
      <!--   App name     -->
      <div class="flex items-center gap-2">
        <div class="flex flex-row items-center gap-2 overflow-hidden">
          <div class="flex items-center justify-center gap-2 font-medium">
            <img
              v-if="applicationGroupDetails.logo"
              :src="applicationGroupDetails.logo"
              class="h-4 w-4 rounded-sm"
              alt="logo" />
            {{ applicationGroupDetails.name }}
          </div>
        </div>
      </div>
      <!--     Status   -->
      <div class="text-center font-medium text-gray-800">
        <div class="flex flex-row items-center gap-5 px-3 text-center">
          <div class="flex flex-row items-center text-sm text-gray-700">
            <font-awesome-icon icon="fa-solid fa-boxes-stacked" class="me-1 text-info-500" />
            {{ totalServiceCount }}&nbsp;Services
          </div>
          <div class="flex flex-row items-center text-sm text-gray-700">
            <font-awesome-icon icon="fa-solid fa-heart-circle-check" class="me-1 text-success-500" />
            {{ healthyServiceCount }}&nbsp;Healthy
          </div>
          <div class="flex flex-row items-center text-sm text-gray-700">
            <font-awesome-icon icon="fa-solid fa-heart-circle-exclamation" class="me-1 text-danger-500" />
            {{ unhealthyServiceCount }}&nbsp;Unhealthy
          </div>
        </div>
      </div>
    </div>
    <!--  Second line  -->
    <div class="mt-3.5 flex w-full flex-row items-center justify-between">
      <div class="flex gap-2">
        <div class="flex items-center gap-2 text-gray-800">
          <div
            v-if="ingressRules.length > 0"
            class="deployment-head max-w-[40vw]"
            :class="{
              '!pr-0.5': ingressRules.length > 0
            }">
            <font-awesome-icon icon="fa-solid fa-globe" />
            <span v-for="(ingressRule, index) in ingressRules" :key="index">
              <a
                :href="
                  ingressRule.protocol +
                  '://' +
                  ((ingressRule.domain?.name || null) ?? 'proxy_server_ip') +
                  ':' +
                  ingressRule.port.toString()
                "
                target="_blank"
                class="has-popover rounded-full bg-primary-500 px-2 py-1 text-secondary-100">
                <font-awesome-icon icon="fa-solid fa-link" class="mr-0.5 text-xs" />
                Link {{ index + 1 }}
                <div class="popover">
                  {{
                    ingressRule.protocol +
                    '://' +
                    ((ingressRule.domain?.name || null) ?? 'proxy_server_ip') +
                    ':' +
                    ingressRule.port.toString()
                  }}
                </div>
              </a>
            </span>
          </div>
          <div v-else class="has-popover flex cursor-pointer gap-2">
            <div class="deployment-head">
              <font-awesome-icon icon="fa-solid fa-globe" />
              <p class="text-warning-600">Not Exposed</p>
            </div>
            <div class="popover w-60">
              No Ingress Rules available. Please open the <b>application details</b> page and create ingress rule to
              expose your application to the internet.
            </div>
          </div>
        </div>
      </div>
      <!--    Quick Actions    -->
      <div class="quick-actions">
        <div class="divider"></div>
        <div class="button" @click="rebuildApplications">
          <font-awesome-icon icon="fa-solid fa-hammer" class="mr-1" />
          Rebuild & Deploy
        </div>
        <div class="divider"></div>
        <div class="button" @click="restartApplications">
          <font-awesome-icon icon="fa-solid fa-rotate-right" class="mr-1" />
          Restart All
        </div>
        <div class="divider"></div>
        <div class="button text-danger-500" @click="deleteApplications">
          <font-awesome-icon icon="fa-solid fa-trash" class="mr-1" />
          Delete All
        </div>
      </div>
    </div>
    <!--  main section  -->
    <div class="mt-8 flex w-full flex-row gap-5">
      <!--   navbar   -->
      <div class="navbar">
        <div
          class="nav-element"
          :class="{
            'router-link-exact-active': pageName === 'deployed-apps'
          }"
          @click="pageName = 'deployed-apps'">
          Deployed Apps
        </div>
        <div
          class="nav-element"
          :class="{ 'router-link-exact-active': pageName === 'persistent-volumes' }"
          @click="pageName = 'persistent-volumes'">
          Persistent Volume
        </div>
        <div
          class="nav-element"
          :class="{ 'router-link-exact-active': pageName === 'environment-variables' }"
          @click="pageName = 'environment-variables'">
          Environment Variables
        </div>
        <div
          class="nav-element"
          :class="{ 'router-link-exact-active': pageName === 'static-app-configs' }"
          @click="pageName = 'static-app-config'">
          Static App Config
        </div>
      </div>

      <div class="w-full">
        <!--    Deployed Apps  -->
        <div class="w-full" v-if="pageName === 'deployed-apps'">
          <Table>
            <template v-slot:header>
              <TableHeader align="left">Application Name</TableHeader>
              <TableHeader align="center">Health Status</TableHeader>
              <TableHeader align="center">Replicas</TableHeader>
              <TableHeader align="center">Deploy Status</TableHeader>
              <TableHeader align="center">Last Deployment</TableHeader>
              <TableHeader align="right">View Details</TableHeader>
            </template>
            <template v-slot:message>
              <TableMessage v-if="applications.length === 0">
                No applications found, in this project.<br />
                You can attach your app to new project by in application details page.
              </TableMessage>
            </template>
            <template v-slot:body>
              <ApplicationListRow
                v-for="application in applications"
                :key="application.id"
                :application="application" />
            </template>
          </Table>
        </div>
        <!--  Persistent Volume    -->
        <div v-else-if="pageName === 'persistent-volumes'" class="flex w-full flex-col gap-3">
          <div class="flex flex-row flex-wrap gap-2">
            <div class="w-min cursor-pointer rounded-md px-2 py-2 text-sm font-medium text-secondary-700">
              Applications
            </div>
            <div
              v-for="application in applications"
              v-bind:key="application.id"
              class="w-min cursor-pointer rounded-md border border-secondary-200 px-3 py-2 text-sm text-secondary-700 hover:bg-secondary-100"
              :class="{
                'border-secondary-400 bg-secondary-50':
                  pageInfo.currentSelectedPersistentVolumeApplicationId === application.id
              }"
              @click="pageInfo.currentSelectedPersistentVolumeApplicationId = application.id">
              {{ application.name }}
            </div>
          </div>
          <div
            class="w-full"
            v-for="application in applications"
            v-bind:key="application.id"
            v-show="pageInfo.currentSelectedPersistentVolumeApplicationId === application.id">
            <PersistentVolumeBindingEditor
              :on-mounting-path-change="(key, value) => onPersistentVolumeMountingPathChange(application, key, value)"
              :on-persistent-volume-change="(key, value) => onPersistentVolumeChange(application, key, value)"
              :delete-persistent-volume-binding="(key) => deletePersistentVolumeBinding(application, key)"
              :add-persistent-volume-binding="() => addPersistentVolumeBinding(application)"
              :persistent-volume-bindings-map="persistentVolumeBindingMap(application)"
              :persistent-volume-binding-keys="persistentVolumeBindingKeys(application)" />
          </div>
        </div>
        <!--  Environment Variables  -->
        <div v-else-if="pageName === 'environment-variables'" class="flex w-full flex-col gap-3">
          <div class="flex flex-row flex-wrap gap-2">
            <div class="w-min cursor-pointer rounded-md px-2 py-2 text-sm font-medium text-secondary-700">
              Applications
            </div>
            <div
              v-for="application in applications"
              v-bind:key="application.id"
              class="w-min cursor-pointer rounded-md border border-secondary-200 px-3 py-2 text-sm text-secondary-700 hover:bg-secondary-100"
              :class="{
                'border-secondary-400 bg-secondary-50':
                  pageInfo.currentSelectedEnvironmentVariableApplicationId === application.id
              }"
              @click="pageInfo.currentSelectedEnvironmentVariableApplicationId = application.id">
              {{ application.name }}
            </div>
          </div>
          <EnvironmentVariablesEditor
            v-for="application in applications"
            v-bind:key="application.id"
            v-show="pageInfo.currentSelectedEnvironmentVariableApplicationId === application.id"
            :on-variable-value-change="(key, value) => onEnvironmentVariableValueChange(application, key, value)"
            :on-variable-name-change="(key, name) => onEnvironmentVariableNameChange(application, key, name)"
            :delete-environment-variable="(key) => deleteEnvironmentVariable(application, key)"
            :add-environment-variable="() => addEnvironmentVariable(application)"
            :environment-variables-map="environmentVariableMap(application)"
            :environment-variables-keys="environmentVariableKeys(application)" />
        </div>
        <!--   Config Mounts   -->
        <div v-else-if="pageName === 'static-app-config'" class="flex w-full flex-col gap-3">
          <div class="flex flex-row flex-wrap gap-2">
            <div class="w-min cursor-pointer rounded-md px-2 py-2 text-sm font-medium text-secondary-700">
              Applications
            </div>
            <div
              v-for="application in applications"
              v-bind:key="application.id"
              class="w-min cursor-pointer rounded-md border border-secondary-200 px-3 py-2 text-sm text-secondary-700 hover:bg-secondary-100"
              :class="{
                'border-secondary-400 bg-secondary-50':
                  pageInfo.currentSelectedConfigMountApplicationId === application.id
              }"
              @click="pageInfo.currentSelectedConfigMountApplicationId = application.id">
              {{ application.name }}
            </div>
          </div>
          <div
            class="w-full"
            v-for="application in applications"
            v-bind:key="application.id"
            v-show="pageInfo.currentSelectedConfigMountApplicationId === application.id">
            <ConfigMountsEditor
              :on-config-content-change="(key, content) => onConfigMountContentChange(application, key, content)"
              :delete-config-mount="(key) => deleteConfigMount(application, key)"
              :add-config-mount="(details) => addConfigMount(application, details)"
              :config-mounts-map="configMountMap(application)"
              :config-mounts-keys="configMountKeys(application)" />
          </div>
        </div>
        <!--  Update Config Notify bar  -->
        <div
          v-if="isAnyAppInfoChanged"
          class="mt-4 flex flex-row items-center justify-end gap-2 rounded-md border border-gray-300 p-2">
          <span class="mr-4 font-medium">You have updated some of the configuration</span>
          <FilledButton type="primary" :click="applyChanges" :loading="isApplyingChanges"> Apply Changes</FilledButton>
          <FilledButton type="secondary" :click="refetchGroupApplicationDetails"> Cancel</FilledButton>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.deployment-head {
  @apply relative flex items-center justify-center gap-2.5  rounded-full border border-secondary-300 px-2 py-1 text-sm font-normal;
}

.quick-actions {
  @apply flex overflow-hidden rounded-full border border-secondary-300 text-sm  text-secondary-700;

  .button {
    @apply cursor-pointer px-2.5 py-1 hover:bg-secondary-200;
  }

  .divider {
    @apply h-auto w-[1px] bg-secondary-300;
  }
}

.navbar {
  @apply flex h-min select-none flex-col flex-wrap gap-1 rounded-lg border border-secondary-200 p-1.5;
}

.nav-element {
  @apply min-w-max cursor-pointer rounded-md px-3 py-2 text-sm text-secondary-700 hover:bg-secondary-100;
}

.router-link-exact-active {
  @apply bg-secondary-100 font-medium text-black;
}
</style>
