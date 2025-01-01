<script setup>
import { useRouter } from 'vue-router'
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, ref } from 'vue'
import Badge from '@/views/components/Badge.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import ApplicationDetailsNavbar from '@/views/partials/ApplicationDetailsNavbar.vue'
import NewApplicationUpdaterStore from '@/store/applicationUpdater.js'
import FilledButton from '@/views/components/FilledButton.vue'
import { toast } from 'vue-sonner'
import { isNaN } from 'lodash'
import UptimeChart from '@/views/components/UptimeChart.vue'
import UpdateApplicationGroupModal from '@/views/partials/UpdateApplicationGroupModal.vue'
import { camelCaseToSpacedCapitalized } from '@/vendor/utils.js'

// Get the application ID from the URL
const router = useRouter()
const applicationId = router.currentRoute.value.params.id

// Fetch the application details
const {
  result: applicationDetailsRaw,
  loading: applicationDetailsLoading,
  refetch: refetchApplicationDetails
} = useQuery(
  gql`
    query ($id: String!) {
      application(id: $id) {
        id
        name
        isDeleted
        deploymentMode
        replicas
        isSleeping
        realtimeInfo {
          InfoFound
          DesiredReplicas
          RunningReplicas
          DeploymentMode
          HealthStatus
        }
        latestDeployment {
          id
          status
          upstreamType
          dockerImage
          gitProvider
          repositoryName
          repositoryOwner
          repositoryBranch
          codePath
          createdAt
        }
        ingressRules {
          domain {
            name
          }
          protocol
          port
        }
        applicationGroup {
          id
          name
        }
      }
    }
  `,
  {
    id: applicationId
  },
  {
    pollInterval: 10000
  }
)

const applicationDetails = computed(() => applicationDetailsRaw.value?.application ?? {})
const realtimeInfo = computed(() => applicationDetailsRaw.value?.application?.realtimeInfo ?? {})
const realtimeReplicaCountPercentage = computed(() => {
  try {
    return (realtimeInfo.value.RunningReplicas / applicationDetails.value.replicas) * 100
  } catch (e) {
    return 0
  }
})
const deploymentMode = computed(() => applicationDetails.value?.deploymentMode ?? '')

const isIngressRulesAvailable = computed(() => {
  return (applicationDetails.value?.ingressRules ?? []).length > 0
})

// Environment variables editor
const applicationUpdater = NewApplicationUpdaterStore(applicationId)()

// App Doze Mode
const {
  mutate: sleepApplication,
  onDone: onSleepApplicationDone,
  onError: onSleepApplicationError
} = useMutation(
  gql`
    mutation ($id: String!) {
      sleepApplication(id: $id)
    }
  `,
  {
    variables: {
      id: applicationId
    }
  }
)

onSleepApplicationDone(() => {
  toast.success('Application will be paused in a few seconds')
  refetchApplicationDetails()
})

onSleepApplicationError((error) => {
  toast.error(error.message)
})

const {
  mutate: wakeApplication,
  onDone: onWakeApplicationDone,
  onError: onWakeApplicationError
} = useMutation(
  gql`
    mutation ($id: String!) {
      wakeApplication(id: $id)
    }
  `,
  {
    variables: {
      id: applicationId
    }
  }
)

onWakeApplicationDone(() => {
  toast.success('Application will be resumed in a few seconds')
  refetchApplicationDetails()
})

onWakeApplicationError((error) => {
  toast.error(error.message)
})

// Restart Application
const {
  mutate: restartApplication,
  onError: restartApplicationError,
  onDone: restartApplicationDone
} = useMutation(
  gql`
    mutation ($id: String!) {
      restartApplication(id: $id)
    }
  `,
  {
    fetchPolicy: 'no-cache',
    variables: {
      id: router.currentRoute.value.params.id
    }
  }
)

restartApplicationDone((result) => {
  if (result.data.restartApplication) {
    toast.success('Application restarted successfully !')
  } else {
    toast.error('Something went wrong !')
  }
})

restartApplicationError((error) => {
  toast.error(error.message)
})

const restartApplicationWithConfirmation = () => {
  const confirmation = confirm('Are you sure that you want to restart this application ?')
  if (confirmation) {
    restartApplication()
  }
}

// Rebuild Application
const {
  mutate: rebuildApplication,
  onError: rebuildApplicationError,
  onDone: rebuildApplicationDone
} = useMutation(
  gql`
    mutation ($id: String!) {
      rebuildApplication(id: $id)
    }
  `,
  {
    fetchPolicy: 'no-cache',
    variables: {
      id: router.currentRoute.value.params.id
    }
  }
)

rebuildApplicationDone((result) => {
  if (result.data.rebuildApplication) {
    toast.success('Application rebuild request sent successfully !')
  } else {
    toast.error('Something went wrong !')
  }
  router.push({
    name: 'Application Details Deployments',
    params: {
      id: router.currentRoute.value.params.id
    }
  })
})

rebuildApplicationError((error) => {
  toast.error(error.message)
})

const rebuildApplicationWithConfirmation = () => {
  const confirmation = confirm('Are you sure that you want to rebuild this application ?')
  if (confirmation) {
    rebuildApplication()
  }
}

// Application group update
const applicationGroupUpdateModalRef = ref(null)
const openApplicationGroupUpdateModal = () => {
  if (applicationGroupUpdateModalRef.value) applicationGroupUpdateModalRef.value.openModal()
}
</script>

<template>
  <!-- Application group update modal -->
  <UpdateApplicationGroupModal
    ref="applicationGroupUpdateModalRef"
    :current-group-id="applicationDetails.applicationGroup?.id ?? null"
    :application-id="applicationDetails.id"
    :callback-on-update="refetchApplicationDetails" />

  <!-- Main -->
  <div v-if="applicationDetailsLoading">
    <p>Loading...</p>
  </div>
  <section v-else class="mx-auto w-full max-w-7xl">
    <!--  First line  -->
    <div class="flex w-full flex-row items-center justify-between">
      <!--   App name     -->
      <div class="flex items-center gap-2">
        <div class="flex flex-row items-center gap-2 overflow-hidden">
          <div
            @click="openApplicationGroupUpdateModal"
            class="flex cursor-pointer items-center justify-center rounded-full bg-secondary-600 px-3 py-1 text-sm font-medium text-white hover:bg-secondary-700">
            <span v-if="applicationDetails.applicationGroup">{{ applicationDetails.applicationGroup.name }}</span>
            <span v-else>no project</span>
            &nbsp;&nbsp;
            <font-awesome-icon icon="fa-solid fa-caret-down" />
          </div>
          <div class="flex items-center justify-center gap-2 font-medium">
            {{ applicationDetails.name }}
            <Badge v-if="applicationDetails.latestDeployment.status === 'deployed'" type="success">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
            <Badge v-else-if="applicationDetails.latestDeployment.status === 'pending'" type="warning">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
            <Badge v-else-if="applicationDetails.latestDeployment.status === 'deployPending'" type="warning">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
            <Badge v-else-if="applicationDetails.latestDeployment.status === 'deploying'" type="warning">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
            <Badge v-else-if="applicationDetails.latestDeployment.status === 'failed'" type="danger">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
            <Badge v-else-if="applicationDetails.latestDeployment.status === 'stopped'" type="secondary">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
            <Badge v-else-if="applicationDetails.latestDeployment.status === 'stalled'" type="secondary">
              {{ camelCaseToSpacedCapitalized(applicationDetails.latestDeployment.status) }}
            </Badge>
          </div>
        </div>
      </div>
      <!--     Status   -->
      <div class="text-center font-medium text-gray-800">
        <p v-if="applicationDetails.isSleeping" class="w-full pe-[5vw] text-center text-sm text-blue-600">
          <font-awesome-icon icon="fa-solid fa-bed" />
          Sleeping
        </p>
        <div v-else-if="realtimeInfo.InfoFound" class="flex flex-row items-center gap-5 px-3 text-center">
          <div
            v-if="applicationDetails.realtimeInfo.HealthStatus === 'healthy'"
            class="flex flex-row items-center text-sm text-gray-700">
            <font-awesome-icon icon="fa-solid fa-heart-circle-check" class="me-1 text-success-500" />
            Healthy
          </div>
          <div
            v-else-if="applicationDetails.realtimeInfo.HealthStatus === 'unhealthy'"
            class="flex flex-row items-center text-sm text-gray-700">
            <font-awesome-icon icon="fa-solid fa-heart-circle-exclamation" class="me-1 text-danger-500" />
            Unhealthy
          </div>
          <div
            v-else-if="applicationDetails.realtimeInfo.HealthStatus === 'unknown'"
            class="flex flex-row items-center text-sm text-gray-700">
            <font-awesome-icon icon="fa-solid fa-heart-circle-xmark" class="me-1 text-warning-600" />
            Unknown
          </div>
          <UptimeChart
            hide-hover
            small
            label-position="right"
            v-if="!isNaN(realtimeReplicaCountPercentage) && deploymentMode === 'replicated'"
            :percentage="realtimeReplicaCountPercentage"
            :label="`(${realtimeInfo.RunningReplicas ?? 0} / ${applicationDetails.replicas})`" />
          <p v-else-if="deploymentMode === 'global'" class="w-full text-center text-sm text-secondary-700">
            {{ realtimeInfo.RunningReplicas ?? 0 }} Instance{{ realtimeInfo.RunningReplicas > 1 ? 's' : '' }} running
          </p>
          <p v-else class="text-warning-600">Not Available</p>
        </div>
        <p v-else class="text-sm text-warning-600">Sorry, Health info not available currently</p>
      </div>
    </div>
    <!--  Second line  -->
    <div class="mt-3.5 flex w-full flex-row items-center justify-between">
      <!--   Deployment info   -->
      <div class="flex gap-2">
        <div class="flex items-center gap-2 text-gray-800">
          <div v-if="applicationDetails.latestDeployment.upstreamType === 'git'" class="flex gap-2">
            <div class="deployment-head">
              <font-awesome-icon icon="fa-brands fa-github" />
              {{ applicationDetails.latestDeployment.repositoryOwner }}/{{
                applicationDetails.latestDeployment.repositoryName
              }}
            </div>
            <div class="deployment-head">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="lucide lucide-git-branch h-4 w-4">
                <line x1="6" x2="6" y1="3" y2="15" />
                <circle cx="18" cy="6" r="3" />
                <circle cx="6" cy="18" r="3" />
                <path d="M18 9a9 9 0 0 1-9 9" />
              </svg>
              {{ applicationDetails.latestDeployment.repositoryBranch }}
            </div>
          </div>
          <p v-if="applicationDetails.latestDeployment.upstreamType === 'image'" class="deployment-head">
            <font-awesome-icon icon="fa-brands fa-docker" />
            {{ applicationDetails.latestDeployment.dockerImage }}
          </p>
          <p v-if="applicationDetails.latestDeployment.upstreamType === 'sourceCode'" class="deployment-head">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="lucide lucide-upload h-4 w-4">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
              <polyline points="17 8 12 3 7 8" />
              <line x1="12" x2="12" y1="3" y2="15" />
            </svg>
            Source-code uploaded manually
          </p>
        </div>
        <div class="flex items-center gap-2 text-gray-800">
          <div
            v-if="isIngressRulesAvailable"
            class="deployment-head max-w-[40vw]"
            :class="{
              '!pr-0.5': applicationDetails.ingressRules.length > 0
            }">
            <font-awesome-icon icon="fa-solid fa-globe" />
            <span v-for="(ingressRule, index) in applicationDetails.ingressRules" :key="index">
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
          <div v-else class="has-popover flex gap-2">
            <div class="deployment-head">
              <font-awesome-icon icon="fa-solid fa-globe" />
              <p class="text-warning-600">Not Exposed</p>
              <RouterLink
                :to="{
                  name: 'Application Details Ingress Rules',
                  params: { id: $route.params.id }
                }"
                class="font-semibold hover:cursor-pointer hover:text-primary-600">
                <font-awesome-icon icon="fa-solid fa-plus" />
              </RouterLink>
            </div>
            <div class="popover w-60">
              No Ingress Rules available. Click the <b>plus</b> button to add ingress rules if you want to expose the
              application to the internet.
            </div>
          </div>
        </div>
      </div>
      <!--    Quick Actions    -->
      <div class="quick-actions">
        <div class="button" v-if="applicationDetails.isSleeping" @click="wakeApplication">
          <font-awesome-icon icon="fa-solid fa-play" class="mr-1" />
          Resume
        </div>
        <div class="divider" v-if="applicationDetails.isSleeping"></div>
        <div class="button" v-if="!applicationDetails.isSleeping" @click="sleepApplication">
          <font-awesome-icon icon="fa-solid fa-pause" class="mr-1" />
          Pause
        </div>
        <div class="divider" v-if="!applicationDetails.isSleeping"></div>
        <div class="button" @click="rebuildApplicationWithConfirmation">
          <font-awesome-icon icon="fa-solid fa-hammer" class="mr-1" />
          Rebuild
        </div>
        <div class="divider"></div>
        <div class="button" @click="restartApplicationWithConfirmation">
          <font-awesome-icon icon="fa-solid fa-rotate-right" class="mr-1" />
          Restart
        </div>
      </div>
    </div>
    <div class="mt-8 flex w-full flex-row gap-5">
      <!--  Vertical navbar for links    -->
      <ApplicationDetailsNavbar />

      <div class="w-full">
        <!--  Nested Router View  -->
        <RouterView />
        <!--  Update Config Notify bar  -->
        <div
          v-if="applicationUpdater.isConfigurationUpdated"
          class="mt-4 flex flex-row items-center justify-end gap-2 rounded-md border border-gray-300 p-2">
          <span class="mr-4 font-medium">You have updated some of the configuration</span>
          <FilledButton
            :click="applicationUpdater.applyConfigurationChanges"
            :loading="applicationUpdater.isDeployRequestSubmitting"
            type="primary">
            Apply Changes
          </FilledButton>
          <FilledButton
            :click="applicationUpdater.cancelConfigurationChanges"
            :disabled="applicationUpdater.isDeployRequestSubmitting"
            type="secondary">
            Cancel
          </FilledButton>
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
</style>
