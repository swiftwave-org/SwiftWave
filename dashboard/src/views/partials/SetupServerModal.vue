<script setup>
import { toast } from 'vue-sonner'
import { computed, reactive, ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import Step from '@/views/components/Step.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import gql from 'graphql-tag'
import { useLazyQuery, useMutation, useQuery } from '@vue/apollo-composable'
import Code from '@/views/components/Code.vue'
import Table from '@/views/components/Table/Table.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import TableRow from '@/views/components/Table/TableRow.vue'
import Badge from '@/views/components/Badge.vue'
import { useRouter } from 'vue-router'

// Modal related methods
const props = defineProps({
  serverId: {
    type: Number,
    required: true
  },
  serverIp: {
    type: String,
    required: true
  },
  refetchServer: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const router = useRouter()
const isModalOpen = ref(false)
const verifyDependencyInterval = ref(null)

const openModal = () => {
  isModalOpen.value = true
}
const closeModal = () => {
  networkInterfacesOfServer.value = []
  info.step = 1
  info.sshVerified = -1
  info.dependenciesInstalled = -1
  info.dependencyReport = []
  info.isAutomatedDependenciesTriggered = false
  info.dockerUnixPath = '/var/run/docker.sock'
  info.swarmMode = 'manager'
  info.advertiseIP = ''
  isModalOpen.value = false
}

defineExpose({
  openModal,
  closeModal
})

// Other methods
// -1 = neutral, 0 = failed, 1 = success

const info = reactive({
  step: 1,
  sshVerified: -1,
  dependenciesInstalled: -1,
  dependencyReport: [],
  isAutomatedDependenciesTriggered: false,
  dockerUnixPath: '/var/run/docker.sock',
  swarmMode: 'manager',
  advertiseIP: ''
})

// Verify SSH connection
const {
  mutate: testSshAccessToServer,
  onError: onTestSshAccessToServerError,
  onDone: onTestSshAccessToServerSuccess,
  loading: isTestingSshAccess
} = useMutation(
  gql`
    mutation ($id: Uint!) {
      testSSHAccessToServer(id: $id)
    }
  `,
  {
    variables: {
      id: props.serverId
    }
  }
)

onTestSshAccessToServerError((err) => {
  toast.error(err.message)
  info.sshVerified = 0
})

onTestSshAccessToServerSuccess((result) => {
  if (result?.data?.testSSHAccessToServer ?? false) {
    info.sshVerified = 1
    info.step = 2
    verifyDependencyInterval.value = setInterval(() => {
      verifyDependencies()
    }, 5000)
  } else {
    info.sshVerified = 0
  }
})

// Verify dependencies
const {
  mutate: verifyDependencies,
  onError: onVerifyDependenciesError,
  onDone: onVerifyDependenciesSuccess,
  loading: isVerifyingDependencies
} = useMutation(
  gql`
    mutation ($id: Uint!) {
      checkDependenciesOnServer(id: $id) {
        name
        available
      }
    }
  `,
  {
    variables: {
      id: props.serverId
    }
  }
)

onVerifyDependenciesError((err) => {
  toast.error(err.message)
  info.dependenciesInstalled = -1
})

onVerifyDependenciesSuccess((result) => {
  if (result?.data?.checkDependenciesOnServer ?? []) {
    info.dependencyReport = result.data.checkDependenciesOnServer ?? []
    let isAvailable = true
    for (const dep of info.dependencyReport) {
      if (!dep.available) {
        isAvailable = false
        break
      }
    }
    info.dependenciesInstalled = isAvailable ? 1 : 0
    if (isAvailable) {
      info.step = 3
      if (verifyDependencyInterval.value !== null) {
        clearInterval(verifyDependencyInterval.value)
        verifyDependencyInterval.value = null
      }
      loadNetworkInterfacesOfServer()
    }
  }
})

// Install dependencies on server
const {
  mutate: installDependenciesMutate,
  onError: onInstallDependenciesOnServerError,
  onDone: onInstallDependenciesOnServerSuccess,
  loading: isInstallingDependencies
} = useMutation(gql`
  mutation ($id: Uint!) {
    installDependenciesOnServer(id: $id)
  }
`)

onInstallDependenciesOnServerError((err) => {
  toast.error(err.message)
})

onInstallDependenciesOnServerSuccess((val) => {
  if (val.data.installDependenciesOnServer) {
    toast.success('Automated dependencies installation started\nCheck logs for more info')
    info.isAutomatedDependenciesTriggered = true
  } else {
    toast.error('Failed to install dependencies')
    info.isAutomatedDependenciesTriggered = false
  }
})

const installDependenciesOnServer = () => {
  installDependenciesMutate({ id: props.serverId })
}

// Start Setup
const {
  mutate: setupServerRaw,
  loading: isSettingUpServer,
  onError: onSetupServerError,
  onDone: onSetupServerSuccess
} = useMutation(gql`
  mutation ($input: ServerSetupInput!) {
    setupServer(input: $input)
  }
`)

const setupServer = () => {
  setupServerRaw({
    input: {
      id: props.serverId,
      swarmMode: info.swarmMode,
      dockerUnixSocketPath: info.dockerUnixPath,
      advertiseIP: info.advertiseIP
    }
  })
}

onSetupServerError((err) => {
  toast.error(err.message)
})

onSetupServerSuccess(() => {
  toast.success('Server setup started successfully')
  closeModal()
  props.refetchServer()
})

// View logs
const viewLogs = () => {
  const url = router.resolve({
    name: 'Server Logs',
    query: {
      id: props.serverId,
      name: props.serverIp
    }
  }).href
  window.open(url, '_blank')
}

// Fetch public SSH key
const {
  result: fetchPublicSSHKeyResult,
  loading: isFetchingPublicSSHKey,
  onError: onFetchPublicSSHKeyError
} = useQuery(gql`
  query {
    publicSSHKey
  }
`)

onFetchPublicSSHKeyError((err) => {
  toast.error(err.message)
})

const publicSSHKey = computed(() => fetchPublicSSHKeyResult.value?.publicSSHKey ?? '')

// Fetch network interfaces of the server
const networkInterfacesOfServer = ref([])
const {
  load: loadNetworkInterfacesOfServerRaw,
  refetch: refetchNetworkInterfacesOfServerRaw,
  loading: isNetworkInterfacesOfServerLoading,
  onResult: onNetworkInterfacesOfServerResult,
  onError: onNetworkInterfacesOfServerError
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      networkInterfacesOnServer(id: $id) {
        name
        ip
      }
    }
  `,
  {
    id: props.serverId
  }
)

const loadNetworkInterfacesOfServer = () => {
  if (loadNetworkInterfacesOfServerRaw() === false) {
    refetchNetworkInterfacesOfServerRaw()
  }
}

onNetworkInterfacesOfServerError((err) => {
  toast.error(err.message)
})

onNetworkInterfacesOfServerResult((result) => {
  if (result.data.networkInterfacesOnServer) {
    networkInterfacesOfServer.value = result.data.networkInterfacesOnServer
    if (networkInterfacesOfServer.value.length > 0) {
      info.advertiseIP = networkInterfacesOfServer.value[0].ip
    }
  }
})
</script>

<template>
  <teleport to="body">
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" :key="serverId + '_setup_server_modal'" width="2xl">
      <template v-slot:header>Setup Server - {{ serverIp }}</template>
      <template v-slot:body>
        <div class="mt-6">
          <!--    Verify SSH Connection       -->
          <Step
            title="Verify SSH Connection"
            sub-title="Need SSH connectivity to prepare server"
            prefixText="1"
            :type="info.sshVerified === 1 ? 'success' : info.sshVerified === 0 ? 'danger' : 'secondary'"
            :show-body="info.sshVerified !== 1">
            <p class="my-2 text-sm text-gray-500" v-if="isFetchingPublicSSHKey">Fetching Instructions...</p>
            <div v-if="info.sshVerified === 0" class="mt-4">
              <p>Put this key in <b>/root/.ssh/authorized_keys</b></p>
              <Code>{{ publicSSHKey }}</Code>
              <p>Or, run this commands:</p>
              <Code
                >sudo mkdir -p /root/.ssh && echo "{{ publicSSHKey }}" | sudo tee -a /root/.ssh/authorized_keys</Code
              >
              <p class="font-medium italic">After adding the key, verify the SSH connection to proceed further</p>
            </div>
            <FilledButton class="mt-2" :click="testSshAccessToServer" :loading="isTestingSshAccess"
              >Click to Verify SSH Connection
            </FilledButton>
          </Step>

          <!--    Install Dependencies       -->
          <Step
            v-if="info.step >= 2"
            title="Install Required Packages"
            sub-title="To run swiftwave perfectly, it needs some additional packages"
            prefix-text="2"
            :type="
              info.dependenciesInstalled === 1 ? 'success' : info.dependenciesInstalled === 0 ? 'danger' : 'secondary'
            ">
            <FilledButton
              v-if="info.dependenciesInstalled === -1"
              class="mt-2"
              :click="verifyDependencies"
              :loading="isVerifyingDependencies"
              >Click to Verify Required Dependencies
            </FilledButton>
            <div v-else-if="info.dependenciesInstalled === 0">
              <!--    Show Report        -->
              <Table class="mt-4">
                <template v-slot:header>
                  <TableHeader align="left">Package Name</TableHeader>
                  <TableHeader align="center">Status</TableHeader>
                </template>
                <template v-slot:body>
                  <tr v-for="dep in info.dependencyReport" :key="dep.name">
                    <TableRow align="left">{{ dep.name }}</TableRow>
                    <TableRow align="center">
                      <Badge v-if="dep.available" type="success">Installed</Badge>
                      <Badge v-else type="danger">Not Installed</Badge>
                    </TableRow>
                  </tr>
                </template>
              </Table>
              <!--    Show utility options          -->
              <div class="mt-2 flex flex-row gap-2">
                <FilledButton class="w-full" :click="installDependenciesOnServer" :loading="isInstallingDependencies">
                  <font-awesome-icon icon="fa-solid fa-wrench" />&nbsp;&nbsp; Do Automated Installation
                </FilledButton>

                <FilledButton class="w-full" :click="verifyDependencies" :loading="isVerifyingDependencies">
                  <font-awesome-icon
                    icon="fa-solid fa-arrow-rotate-right"
                    v-show="!isVerifyingDependencies" />&nbsp;&nbsp; Verify Required Dependencies
                </FilledButton>
              </div>

              <div class="mt-2 rounded-md border border-secondary-200 p-2" v-if="info.isAutomatedDependenciesTriggered">
                <p class="italic">You can check logs of automated installation here</p>
                <FilledButton class="mt-2" type="primary" :click="viewLogs">
                  <font-awesome-icon icon="fa-solid fa-file-lines" />&nbsp;&nbsp; View Logs of Server
                </FilledButton>
              </div>
            </div>
          </Step>

          <!--    Setup Server      -->
          <Step
            v-if="info.step >= 3"
            title="Start setting up server"
            sub-title="Automated server setup will prepare the server for swiftwave"
            prefix-text="3"
            type="warning">
            <div>
              <label class="block text-sm font-medium text-gray-700" for="domain">Swarm Node Mode</label>
              <div class="mt-1 flex space-x-2">
                <select
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  v-model="info.swarmMode">
                  <option value="manager">Manager</option>
                  <option value="worker">Worker</option>
                </select>
              </div>
              <p class="mt-0.5 text-sm italic text-secondary-600">
                If you have no server registered already, must select <b>Manager</b>
              </p>
            </div>
            <div class="mt-3">
              <label class="block text-sm font-medium text-gray-700" for="domain"
                >Swarm Advertise IP<span class="ml-2 italic" v-if="isNetworkInterfacesOfServerLoading"
                  ><font-awesome-icon icon="fa-solid fa-spinner" class="animate-spin" />&nbsp;&nbsp;Fetching...</span
                ></label
              >
              <div class="mt-1 flex flex-col">
                <select
                  v-if="networkInterfacesOfServer.length > 1"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                  v-model="info.advertiseIP">
                  <option v-for="ip in networkInterfacesOfServer" :key="ip.ip" :value="ip.ip">
                    {{ ip.ip }} [{{ ip.name }}]
                  </option>
                </select>
                <p class="text-sm font-medium">{{ info.advertiseIP }}</p>
                <p v-if="networkInterfacesOfServer.length > 0" class="mt-0.5 text-sm italic text-secondary-600">
                  <span class="font-medium text-danger-500">Note:</span> Make sure that this IP address on your system
                  doesn't get changed on future.
                </p>
              </div>
            </div>
            <div class="mt-3">
              <label class="block text-sm font-medium text-gray-700" for="dockerUnixPath"
                >Docker UNIX Socket Path</label
              >
              <div class="mt-1">
                <input
                  id="dockerUnixPath"
                  v-model="info.dockerUnixPath"
                  class="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                  required
                  type="text" />
              </div>
              <p class="mt-0.5 text-sm italic text-secondary-600">Don't change it until you have special requirement</p>
            </div>
            <FilledButton class="mt-2 w-full" type="primary" :click="setupServer" :loading="isSettingUpServer">
              <font-awesome-icon icon="fa-solid fa-play" />&nbsp;&nbsp; Let's Start Setup
            </FilledButton>
          </Step>
        </div>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
