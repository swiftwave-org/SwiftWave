<script setup>
import 'xterm/css/xterm.css'

import { useMutation, useQuery, useSubscription } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { useRouter } from 'vue-router'
import { computed, onMounted, ref, watch } from 'vue'
import Badge from '@/views/components/Badge.vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { toast } from 'vue-sonner'
import StatusPulse from '@/views/components/StatusPulse.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { camelCaseToSpacedCapitalized } from '@/vendor/utils.js'

const router = useRouter()
const deploymentId = router.currentRoute.value.params.deployment_id

// Fetch the deployment details
const { result: deploymentRaw, loading: deploymentLoading } = useQuery(
  gql`
    query ($id: String!) {
      deployment(id: $id) {
        id
        status
        upstreamType
        gitProvider
        commitMessage
        commitHash
        repositoryName
        repositoryOwner
        repositoryBranch
        dockerImage
        buildArgs {
          key
          value
        }
        createdAt
      }
    }
  `,
  {
    id: deploymentId
  },
  {
    pollInterval: 10000
  }
)

const deployment = computed(() => deploymentRaw.value?.deployment ?? {})
const buildArgs = computed(() => {
  const args = deploymentRaw.value?.deployment?.buildArgs ?? []
  return args
    .map((arg) => {
      return `${arg.key}=${arg.value}`
    })
    .join(' <b>|</b> ')
})

const deployedOn = computed(() => {
  const date = new Date(deploymentRaw.value?.deployment?.createdAt)
  return date.toLocaleString()
})

// Deployment logs
const showDeploymentLog = ref(false)
const terminal = new Terminal({
  convertEol: true,
  rows: 30,
  scrollback: 9999999
})
const fitAddon = new FitAddon()
terminal.loadAddon(fitAddon)

const { result: deploymentLogRaw, onError: onDeploymentLogError } = useSubscription(
  gql`
    subscription ($id: String!) {
      fetchDeploymentLog(id: $id) {
        content
      }
    }
  `,
  {
    id: deploymentId
  },
  {
    enabled: showDeploymentLog
  }
)

onDeploymentLogError((err) => {
  toast.error(err.message)
})

const deploymentLog = computed(() => deploymentLogRaw.value?.fetchDeploymentLog.content ?? '')
watch(deploymentLog, (value) => {
  if (value) {
    terminal.write(value)
  }
})

onMounted(() => {
  terminal.open(document.getElementById('terminal'))
  fitAddon.fit()
  showDeploymentLog.value = true
})

const isTerminalLoading = computed(() => {
  let status = deployment.value?.status ?? ''
  return status === 'pending' || status === 'deployPending'
})

// Cancel deployment
const {
  mutate: cancelDeployment,
  loading: cancelDeploymentLoading,
  onError: onCancelDeploymentError,
  onDone: onCancelDeploymentDone
} = useMutation(
  gql`
    mutation ($id: String!) {
      cancelDeployment(id: $id)
    }
  `,
  {
    fetchPolicy: 'no-cache',
    variables: {
      id: deploymentId
    }
  }
)

onCancelDeploymentDone((val) => {
  if (val.data.cancelDeployment) {
    toast.success('Deployment cancellation request sent.')
  } else {
    toast.error('Deployment cancellation request failed.')
  }
})

onCancelDeploymentError((err) => {
  toast.error(err.message)
})
</script>

<template>
  <div v-if="deploymentLoading">
    <p>Loading...</p>
  </div>
  <section v-else class="mx-auto w-full max-w-7xl text-sm">
    <div class="flex items-center gap-2">
      <p class="text-base font-bold">
        <font-awesome-icon icon="fa-solid fa-signal" />
        Status
      </p>
      <Badge v-if="deployment.status === 'deployed'" type="success"
        >{{ camelCaseToSpacedCapitalized(deployment.status) }}
      </Badge>
      <Badge v-else-if="deployment.status === 'pending'" type="warning">
        {{ camelCaseToSpacedCapitalized(deployment.status) }}
      </Badge>
      <Badge v-else-if="deployment.status === 'deployPending'" type="warning">
        {{ camelCaseToSpacedCapitalized(deployment.status) }}
      </Badge>
      <Badge v-else-if="deployment.status === 'deploying'" type="warning">
        {{ camelCaseToSpacedCapitalized(deployment.status) }}
      </Badge>
      <Badge v-else-if="deployment.status === 'failed'" type="danger"
        >{{ camelCaseToSpacedCapitalized(deployment.status) }}
      </Badge>
      <Badge v-else-if="deployment.status === 'stopped'" type="secondary">
        {{ camelCaseToSpacedCapitalized(deployment.status) }}
      </Badge>
    </div>
    <div class="mt-2 flex items-center gap-2 font-normal text-gray-800">
      <font-awesome-icon icon="fa-solid fa-fingerprint" />
      <p>{{ deployment.id }}</p>
    </div>
    <div class="mt-2 flex items-center gap-2 text-gray-800">
      <font-awesome-icon v-if="deployment.upstreamType === 'git'" icon="fa-solid fa-code-branch" />
      <font-awesome-icon v-if="deployment.upstreamType === 'image'" icon="fa-brands fa-docker" />
      <font-awesome-icon v-if="deployment.upstreamType === 'sourceCode'" icon="fa-solid fa-upload" />
      <p v-if="deployment.upstreamType === 'git'">
        {{ deployment.gitProvider }}@{{ deployment.repositoryOwner }}/{{ deployment.repositoryName }}:{{
          deployment.repositoryBranch
        }}
      </p>
      <p v-if="deployment.upstreamType === 'image'">{{ deployment.dockerImage }}</p>
      <p v-if="deployment.upstreamType === 'sourceCode'">Source-code uploaded manually</p>
    </div>
    <div class="mt-2 flex items-center gap-2 text-gray-800" v-if="deployment.upstreamType === 'git'">
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
        class="h-4 w-4">
        <circle cx="12" cy="12" r="3" />
        <line x1="3" x2="9" y1="12" y2="12" />
        <line x1="15" x2="21" y1="12" y2="12" />
      </svg>
      <p v-if="deployment.commitHash && deployment.commitMessage">
        {{ deployment.commitHash.slice(0, 7) }}
      </p>
      <p
        v-if="deployment.commitHash && deployment.commitMessage"
        class="nowrap max-w-[40vw] overflow-hidden text-ellipsis">
        {{ deployment.commitMessage }}
      </p>
      <p v-if="!(deployment.commitHash && deployment.commitMessage)" class="italic">not available</p>
    </div>
    <div class="mt-2 flex items-center gap-2 font-normal text-gray-800">
      <font-awesome-icon icon="fa-solid fa-calendar-days" />
      <p>{{ deployedOn }}</p>
    </div>
    <div class="mb-2 mt-2 flex items-center gap-2 font-normal text-gray-800" v-if="buildArgs.length !== 0">
      <font-awesome-icon icon="fa-solid fa-hammer" />
      <p><span class="font-medium">Build arguments :</span> <span v-html="buildArgs"></span></p>
    </div>
    <div
      v-if="deployment.status === 'pending'"
      class="mt-2 flex flex-row items-center justify-between rounded-md bg-red-100 px-3 py-2">
      <div>
        <p class="inline-flex items-center gap-2 text-lg font-medium">Cancel Deployment</p>
        <p class="text-sm text-secondary-700">
          If you are feeling deployment has been stuck for a long time, you can cancel the deployment.
        </p>
      </div>
      <FilledButton type="danger" @click="cancelDeployment" :loading="cancelDeploymentLoading"
        >Request Cancellation
      </FilledButton>
    </div>

    <hr class="mb-2 mt-2" />
    <p class="inline-flex items-center gap-2 text-base font-medium">
      Deployment Logs
      <StatusPulse v-if="isTerminalLoading" type="success" />
    </p>
    <p class="text-sm text-secondary-700">
      If you feel that deployment log is not automatically updating, please refresh the page.
    </p>
  </section>
  <div id="terminal" class="mt-3 w-full overflow-hidden rounded-md bg-black p-2"></div>
</template>

<style scoped></style>
