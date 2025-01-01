<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { computed, reactive, ref } from 'vue'
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { preventSpaceInput } from '@/vendor/utils.js'

const props = defineProps({
  callbackOnCreate: {
    type: Function,
    required: false,
    default: () => {}
  },
  callbackOnPop: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isModalOpen = ref(false)
// Create persistent volume
const newPersistentVolumeDetails = reactive({
  name: '',
  type: 'local',
  nfsConfig: {
    host: '',
    path: '',
    version: 4
  },
  cifsConfig: {
    host: '',
    share: '',
    username: '',
    password: '',
    file_mode: '0777',
    dir_mode: '0777',
    uid: 0,
    gid: 0
  }
})

const openModal = () => {
  newPersistentVolumeDetails.name = ''
  newPersistentVolumeDetails.type = 'local'
  newPersistentVolumeDetails.nfsConfig.host = ''
  newPersistentVolumeDetails.nfsConfig.path = ''
  newPersistentVolumeDetails.nfsConfig.version = 4
  newPersistentVolumeDetails.cifsConfig.host = ''
  newPersistentVolumeDetails.cifsConfig.share = ''
  newPersistentVolumeDetails.cifsConfig.username = ''
  newPersistentVolumeDetails.cifsConfig.password = ''
  newPersistentVolumeDetails.cifsConfig.file_mode = '0777'
  newPersistentVolumeDetails.cifsConfig.dir_mode = '0777'
  newPersistentVolumeDetails.cifsConfig.uid = 0
  newPersistentVolumeDetails.cifsConfig.gid = 0
  isModalOpen.value = true
}

const closeModal = () => {
  isModalOpen.value = false
  props.callbackOnPop()
}

const {
  mutate: registerPersistentVolume,
  loading: isDomainRegistering,
  onDone: onDomainRegisterSuccess,
  onError: onDomainRegisterFail
} = useMutation(
  gql`
    mutation ($input: PersistentVolumeInput!) {
      createPersistentVolume(input: $input) {
        id
        name
      }
    }
  `,
  {
    variables: {
      input: newPersistentVolumeDetails
    }
  }
)

onDomainRegisterSuccess((result) => {
  closeModal()
  props.callbackOnCreate(result.data.createPersistentVolume)
})

onDomainRegisterFail((err) => {
  toast.error(err.message)
})

const { result: noOfServersResult } = useQuery(gql`
  query {
    noOfServers
  }
`)

const noOfServers = computed(() => noOfServersResult.value?.noOfServers ?? 0)

defineExpose({
  openModal,
  closeModal
})
</script>

<template>
  <ModalDialog :close-modal="closeModal" :is-open="isModalOpen" width="lg">
    <template v-slot:header>Add New Persistent Volume</template>
    <template v-slot:body>
      Enter a unique name for the persistent volume.
      <form @submit.prevent="">
        <!--  Name Field   -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700" for="name"> Persistent Volume </label>
          <div class="mt-1">
            <input
              id="name"
              v-model="newPersistentVolumeDetails.name"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              name="name"
              placeholder="Name of persistent volume"
              type="text"
              @keydown="preventSpaceInput" />
          </div>
        </div>
        <!--    Type Field      -->
        <div class="mt-2">
          <label class="block text-sm font-medium text-gray-700">Type</label>
          <select
            v-model="newPersistentVolumeDetails.type"
            class="mt-2 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
            <option value="local">Local</option>
            <option value="nfs">NFS</option>
            <option value="cifs">CIFS</option>
          </select>
          <div
            class="mb-5 mt-3 rounded border-s-4 border-danger-200 bg-danger-50 p-4"
            v-if="newPersistentVolumeDetails.type === 'local' && noOfServers > 1">
            <p class="block text-justify text-sm text-gray-900">
              You have <b>{{ noOfServers }} servers</b> configured for cluster mode.
              <br />
              Try to avoid create <b>Local</b> type persistent volume. Instead use <b>NFS</b> or <b>CIFS</b> persistent
              volume.
            </p>
          </div>
        </div>
        <!--   NFS Server Host    -->
        <div v-if="newPersistentVolumeDetails.type === 'nfs'" class="mt-2">
          <label class="block text-sm font-medium text-gray-700">NFS Server Host</label>
          <div class="mt-1">
            <input
              v-model="newPersistentVolumeDetails.nfsConfig.host"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="NFS Server Host"
              type="text" />
          </div>
          <p class="mt-1 text-sm text-gray-500">
            Example:
            <span class="text-gray-700"> nfs-server.example.com </span>
          </p>
        </div>
        <!--    NFS Share Path      -->
        <div v-if="newPersistentVolumeDetails.type === 'nfs'" class="mt-2">
          <label class="block text-sm font-medium text-gray-700">NFS Share Path</label>
          <div class="mt-1">
            <input
              v-model="newPersistentVolumeDetails.nfsConfig.path"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="NFS Share Path"
              type="text" />
          </div>
          <p class="mt-1 text-sm text-gray-500">
            Example:
            <span class="text-gray-700"> /mnt/nfs_share </span>
          </p>
        </div>
        <!--  Version -->
        <div v-if="newPersistentVolumeDetails.type === 'nfs'" class="mt-2">
          <label class="block text-sm font-medium text-gray-700">NFS Version</label>
          <select
            v-model="newPersistentVolumeDetails.nfsConfig.version"
            class="mt-2 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
            <option value="4">NFS v4</option>
            <option value="3">NFS v3</option>
            <option value="2">NFS v2</option>
          </select>
          <!-- RCPBind Alert  -->
          <div
            class="mb-5 mt-3 rounded border-s-4 border-danger-200 bg-danger-50 p-4"
            role="alert"
            v-if="
              newPersistentVolumeDetails.nfsConfig.version === '2' ||
              newPersistentVolumeDetails.nfsConfig.version === '3'
            ">
            <p class="block text-justify text-sm text-gray-900">
              During the installation of <b>swiftwave</b>, <b>rpcbind</b> service has been disabled to keep the servers
              secured. But, NFS v2, v3 has requirement of <b>rpcbind</b> service. Kindly enable <b>rpcbind</b> service
              on the servers before proceeding.

              <br /><br />
              You can checkout
              <a
                href="https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/9/html/securing_networks/securing-network-services_securing-networks"
                target="_blank"
                class="font-bold text-blue-500"
                >this blog from Red Hat</a
              >
              to learn best practices to secure <b>rpcbind</b> service.

              <br /><br />Ignore this warning if you have already enabled <b>rpcbind</b> service.
            </p>
          </div>
        </div>
        <!--   CIFS Host     -->
        <div v-if="newPersistentVolumeDetails.type === 'cifs'" class="mt-2">
          <label class="block text-sm font-medium text-gray-700">CIFS Host</label>
          <div class="mt-1">
            <input
              v-model="newPersistentVolumeDetails.cifsConfig.host"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="CIFS Host"
              type="text" />
          </div>
          <p class="mt-1 text-sm text-gray-500">
            Example:
            <span class="text-gray-700">uxxxxx.your-server.de</span>
          </p>
        </div>
        <!--   CIFS Share     -->
        <div v-if="newPersistentVolumeDetails.type === 'cifs'" class="mt-2">
          <label class="block text-sm font-medium text-gray-700">CIFS Share</label>
          <div class="mt-1">
            <input
              v-model="newPersistentVolumeDetails.cifsConfig.share"
              autocomplete="off"
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              placeholder="CIFS Share"
              type="text" />
          </div>
          <p class="mt-1 text-sm text-gray-500">
            Example:
            <span class="text-gray-700">//uxxxxx.your-server.de/backup</span>
          </p>
        </div>
        <div class="mt-2 flex w-full flex-row gap-2" v-if="newPersistentVolumeDetails.type === 'cifs'">
          <!--   CIFS Username     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">CIFS Username</label>
            <div class="mt-1">
              <input
                v-model="newPersistentVolumeDetails.cifsConfig.username"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="CIFS Username"
                type="text" />
            </div>
            <p class="mt-1 text-sm text-gray-500">
              Example:
              <span class="text-gray-700">uxxxxx</span>
            </p>
          </div>
          <!--   CIFS Password   -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">CIFS Password</label>
            <div class="mt-1">
              <input
                v-model="newPersistentVolumeDetails.cifsConfig.password"
                autocomplete="new-password"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="CIFS Password"
                type="password" />
            </div>
          </div>
        </div>

        <div class="mt-2 flex w-full flex-row gap-2" v-if="newPersistentVolumeDetails.type === 'cifs'">
          <!--   CIFS File Mode     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">CIFS File Mode</label>
            <div class="mt-1">
              <input
                v-model="newPersistentVolumeDetails.cifsConfig.file_mode"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="CIFS File Mode"
                type="text" />
            </div>
          </div>
          <!--   CIFS Dir Mode     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">CIFS Dir Mode</label>
            <div class="mt-1">
              <input
                v-model="newPersistentVolumeDetails.cifsConfig.dir_mode"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="CIFS Dir Mode"
                type="text" />
            </div>
          </div>
        </div>

        <div class="mt-2 flex w-full flex-row gap-2" v-if="newPersistentVolumeDetails.type === 'cifs'">
          <!--   CIFS UID     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Mount UID</label>
            <div class="mt-1">
              <input
                v-model="newPersistentVolumeDetails.cifsConfig.uid"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="CIFS UID"
                type="number" />
            </div>
          </div>
          <!--   CIFS Gid     -->
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Mount GID</label>
            <div class="mt-1">
              <input
                v-model="newPersistentVolumeDetails.cifsConfig.gid"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="CIFS GID"
                type="number" />
            </div>
          </div>
        </div>
      </form>
    </template>
    <template v-slot:footer>
      <FilledButton :click="registerPersistentVolume" :loading="isDomainRegistering" type="primary"
        >Register
      </FilledButton>
    </template>
  </ModalDialog>
</template>

<style scoped></style>
