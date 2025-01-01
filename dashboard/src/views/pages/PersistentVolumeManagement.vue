<script setup>
import PageBar from '@/views/components/PageBar.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { toast } from 'vue-sonner'
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import Table from '@/views/components/Table/Table.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import TableMessage from '@/views/components/Table/TableMessage.vue'
import { computed, reactive, ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import PersistentVolumeRow from '@/views/partials/PersistentVolumeRow.vue'
import PersistentVolumeBackups from '@/views/partials/PersistentVolumeBackups.vue'
import PersistentVolumeRestores from '@/views/partials/PersistentVolumeRestores.vue'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import axios from 'axios'
import { useAuthStore } from '@/store/auth.js'
import Badge from '@/views/components/Badge.vue'
import CreatePersistentVolumeModal from '@/views/partials/CreatePersistentVolumeModal.vue'
import SecuredText from '@/views/components/SecuredText.vue'

const authStore = useAuthStore()

// Create persistent volume
const createPersistentVolumeModalRef = ref(null)
const openCreatePersistentVolumeModal = computed(() => createPersistentVolumeModalRef.value?.openModal ?? (() => {}))

// Delete persistent volume
const {
  mutate: deletePersistentVolume,
  onDone: onDomainDeleteSuccess,
  onError: onDomainDeleteFail
} = useMutation(
  gql`
    mutation ($id: Uint!) {
      deletePersistentVolume(id: $id)
    }
  `,
  {
    variables: {
      id: ''
    }
  }
)

const deletePersistentVolumeWithConfirmation = (persistent_volume) => {
  if (
    confirm(
      'Are you sure you want to delete this persistent volume?\n[NOTE] If you have deleted volume bindings just a few seconds ago, please wait at-least a minute before deleting the persistent volume.'
    )
  ) {
    deletePersistentVolume({
      id: persistent_volume.id
    })
  }
}

onDomainDeleteSuccess(() => {
  toast.success(
    'Persistent volume deletion requested. If all conditions are met, the persistent volume will be deleted in the background.'
  )
  refetchPersistentVolumes()
})

onDomainDeleteFail((err) => {
  toast.error(err.message)
})

// Fetch persistent volumes
const {
  result: persistentVolumesRaw,
  refetch: refetchPersistentVolumes,
  loading: isPersistentVolumesLoading,
  onError: onPersistentVolumesError
} = useQuery(gql`
  query {
    persistentVolumes {
      id
      name
      type
      nfsConfig {
        host
        path
        version
      }
      cifsConfig {
        host
        share
        username
        password
        file_mode
        dir_mode
        uid
        gid
      }
    }
  }
`)

const persistentVolumes = computed(() => persistentVolumesRaw.value?.persistentVolumes ?? [])

onPersistentVolumesError((err) => {
  toast.error(err.message)
})

// Backup drawer
const selectedPersistentVolumeId = ref(-1)
const selectedPersistentVolumeName = ref('')
const isBackupDrawerOpen = ref(false)
const closeBackupDrawer = () => (isBackupDrawerOpen.value = false)
const openBackupDrawerForVolume = (id, name) => {
  selectedPersistentVolumeId.value = id
  selectedPersistentVolumeName.value = name
  isBackupDrawerOpen.value = true
}

// Restore drawer
const isRestoreDrawerOpen = ref(false)
const closeRestoreDrawer = () => (isRestoreDrawerOpen.value = false)
const openRestoreDrawerForVolume = (id, name) => {
  selectedPersistentVolumeId.value = id
  selectedPersistentVolumeName.value = name
  isRestoreDrawerOpen.value = true
}

// Restore now
const isRestoreNowModalOpen = ref(false)
const closeRestoreNowModal = () => (isRestoreNowModalOpen.value = false)
const openRestoreNowModal = (id, name) => {
  restoreFileFieldRef.value = null
  selectedPersistentVolumeId.value = id
  selectedPersistentVolumeName.value = name
  isRestoreNowModalOpen.value = true
}

const restoreFileFieldRef = ref(null)
const isRestoreNowButtonDisabled = computed(() => {
  return !restoreFileFieldRef.value?.files?.length
})
const isRestoreNowButtonLoading = ref(false)
const uploadPercentage = ref(0)
const uploadAndRestoreNow = () => {
  isRestoreNowButtonLoading.value = true
  const file = restoreFileFieldRef.value.files[0]
  const formData = new FormData()
  formData.append('file', file)
  const url = getHttpBaseUrl() + `/persistent-volume/${selectedPersistentVolumeId.value}/restore`
  axios
    .post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: authStore.FetchBearerToken()
      },
      onUploadProgress: (progressEvent) => {
        uploadPercentage.value = Math.round((progressEvent.loaded / progressEvent.total) * 100)
      }
    })
    .then((e) => {
      isRestoreNowButtonLoading.value = false
      toast.success('Restore initiated successfully')
      closeRestoreNowModal()
    })
    .catch((err) => {
      toast.error(err.message)
      isRestoreNowButtonLoading.value = false
    })
}

const isVolumeDetailsModalOpen = ref(false)
const selectedVolumeDetails = reactive({
  name: '',
  type: '',
  nfsConfig: {
    host: '',
    path: '',
    version: 0
  },
  cifsConfig: {
    host: '',
    share: '',
    username: '',
    password: '',
    file_mode: '',
    dir_mode: ''
  }
})
const closeVolumeDetailsModal = () => {
  isVolumeDetailsModalOpen.value = false
}

const showDetails = (volume) => {
  selectedVolumeDetails.name = volume.name
  selectedVolumeDetails.type = volume.type
  selectedVolumeDetails.nfsConfig.host = volume?.nfsConfig?.host ?? ''
  selectedVolumeDetails.nfsConfig.path = volume?.nfsConfig?.path ?? ''
  selectedVolumeDetails.nfsConfig.version = volume?.nfsConfig?.version ?? 0
  selectedVolumeDetails.cifsConfig.host = volume?.cifsConfig?.host ?? ''
  selectedVolumeDetails.cifsConfig.share = volume?.cifsConfig?.share ?? ''
  selectedVolumeDetails.cifsConfig.username = volume?.cifsConfig?.username ?? ''
  selectedVolumeDetails.cifsConfig.password = volume?.cifsConfig?.password ?? ''
  selectedVolumeDetails.cifsConfig.file_mode = volume?.cifsConfig?.file_mode ?? ''
  selectedVolumeDetails.cifsConfig.dir_mode = volume?.cifsConfig?.dir_mode ?? ''
  selectedVolumeDetails.cifsConfig.uid = volume?.cifsConfig?.uid ?? ''
  selectedVolumeDetails.cifsConfig.gid = volume?.cifsConfig?.gid ?? ''

  isVolumeDetailsModalOpen.value = true
}
</script>

<template>
  <section class="mx-auto w-full max-w-7xl">
    <!-- Drawer for persistent volume backups -->
    <PersistentVolumeBackups
      :is-drawer-open="isBackupDrawerOpen"
      :close-drawer="closeBackupDrawer"
      :persistent-volume-id="selectedPersistentVolumeId"
      :persistent-volume-name="selectedPersistentVolumeName" />

    <!-- Drawer for persistent volume restores -->
    <PersistentVolumeRestores
      :is-drawer-open="isRestoreDrawerOpen"
      :close-drawer="closeRestoreDrawer"
      :persistent-volume-id="selectedPersistentVolumeId"
      :persistent-volume-name="selectedPersistentVolumeName" />

    <!-- Modal for create persistent volumes -->
    <CreatePersistentVolumeModal ref="createPersistentVolumeModalRef" :callback-on-create="refetchPersistentVolumes" />

    <!--    Modal for create restore -->
    <ModalDialog :close-modal="closeRestoreNowModal" :is-open="isRestoreNowModalOpen" non-cancelable>
      <template v-slot:header>Restore `{{ selectedPersistentVolumeName }}` volume</template>
      <template v-slot:body>
        Choose the backup file (*.tar.gz) to restore the volume.
        <div class="">
          <label class="block text-sm font-medium text-gray-900 dark:text-white" for="source_code"
            >Select Restore File</label
          >
          <div class="mx-auto max-w-md space-y-8">
            <input
              @change="(e) => (restoreFileFieldRef = e.target)"
              class="w-full cursor-pointer rounded-md bg-gray-100 text-sm text-black file:mr-4 file:cursor-pointer file:border-0 file:bg-gray-800 file:px-4 file:py-2 file:text-white file:hover:bg-gray-700 focus:outline-none"
              accept=".tar.gz"
              type="file" />
          </div>
        </div>
      </template>
      <template v-slot:footer>
        <div class="flex w-full flex-col space-y-2">
          <FilledButton
            class="w-full"
            type="primary"
            :disabled="isRestoreNowButtonDisabled"
            :loading="isRestoreNowButtonLoading"
            :click="uploadAndRestoreNow"
            >Upload & Restore Now
            <span v-if="isRestoreNowButtonLoading" class="ml-2">{{ uploadPercentage }}%</span>
          </FilledButton>
          <FilledButton class="w-full" type="secondary" v-if="!isRestoreNowButtonLoading" :click="closeRestoreNowModal">
            Cancel
          </FilledButton>
        </div>
      </template>
    </ModalDialog>

    <!--  Show Volume Details  -->
    <ModalDialog :close-modal="closeVolumeDetailsModal" :is-open="isVolumeDetailsModalOpen">
      <template v-slot:header>Volume Details</template>
      <template v-slot:body>
        <div class="mt-4 flex w-full flex-row gap-2">
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Volume Name</label>
            <div class="mt-1">
              <p
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                {{ selectedVolumeDetails.name }}
              </p>
            </div>
          </div>
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Volume Type</label>
            <div class="mt-1">
              <Badge type="success" v-if="selectedVolumeDetails.type === 'local'">Local</Badge>
              <Badge type="warning" v-if="selectedVolumeDetails.type === 'nfs'"> &nbsp;&nbsp;NFS&nbsp;&nbsp;</Badge>
              <Badge type="warning" v-if="selectedVolumeDetails.type === 'cifs'">&nbsp;CIFS&nbsp;</Badge>
            </div>
          </div>
        </div>
        <div class="mt-4" v-if="selectedVolumeDetails.type === 'nfs'">
          <label class="block text-sm font-medium text-gray-700">NFS Config</label>
          <div class="mt-1">
            <p
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              {{ selectedVolumeDetails.nfsConfig.host }}:{{ selectedVolumeDetails.nfsConfig.path }}
            </p>
          </div>
        </div>
        <div class="mt-4" v-if="selectedVolumeDetails.type === 'nfs'">
          <label class="block text-sm font-medium text-gray-700">NFS Version</label>
          <div class="mt-1">
            <p
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              {{ selectedVolumeDetails.nfsConfig.version }}
            </p>
          </div>
        </div>
        <div class="mt-4" v-if="selectedVolumeDetails.type === 'cifs'">
          <label class="block text-sm font-medium text-gray-700">CIFS Host</label>
          <div class="mt-1">
            <p
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              {{ selectedVolumeDetails.cifsConfig.host }}
            </p>
          </div>
        </div>
        <div class="mt-4" v-if="selectedVolumeDetails.type === 'cifs'">
          <label class="block text-sm font-medium text-gray-700">CIFS Share</label>
          <div class="mt-1">
            <p
              class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
              {{ selectedVolumeDetails.cifsConfig.share }}
            </p>
          </div>
        </div>
        <div class="mt-4 flex w-full flex-row gap-2" v-if="selectedVolumeDetails.type === 'cifs'">
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Username</label>
            <div class="mt-1">
              <p
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                {{ selectedVolumeDetails.cifsConfig.username }}
              </p>
            </div>
          </div>
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Password</label>
            <div class="mt-1">
              <SecuredText>{{ selectedVolumeDetails.cifsConfig.password }}</SecuredText>
            </div>
          </div>
        </div>
        <div class="mt-4 flex w-full flex-row gap-2" v-if="selectedVolumeDetails.type === 'cifs'">
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">File Mode</label>
            <div class="mt-1">
              <p
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                {{ selectedVolumeDetails.cifsConfig.file_mode }}
              </p>
            </div>
          </div>
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">Dir Mode</label>
            <div class="mt-1">
              <p
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                {{ selectedVolumeDetails.cifsConfig.dir_mode }}
              </p>
            </div>
          </div>
        </div>
        <div class="mt-4 flex w-full flex-row gap-2" v-if="selectedVolumeDetails.type === 'cifs'">
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">UID</label>
            <div class="mt-1">
              <p
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                {{ selectedVolumeDetails.cifsConfig.uid }}
              </p>
            </div>
          </div>
          <div class="w-1/2">
            <label class="block text-sm font-medium text-gray-700">GID</label>
            <div class="mt-1">
              <p
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                {{ selectedVolumeDetails.cifsConfig.gid }}
              </p>
            </div>
          </div>
        </div>
      </template>
      <template v-slot:footer>
        <FilledButton :click="closeVolumeDetailsModal" type="primary">Close</FilledButton>
      </template>
    </ModalDialog>

    <!-- Top Page bar   -->
    <PageBar>
      <template v-slot:title>Persistent Volume</template>
      <template v-slot:subtitle>Manage Persistent Volume</template>
      <template v-slot:buttons>
        <FilledButton :click="openCreatePersistentVolumeModal" type="primary">
          <font-awesome-icon icon="fa-solid fa-plus" class="mr-2" />
          Add New
        </FilledButton>
        <FilledButton type="ghost" :click="refetchPersistentVolumes">
          <font-awesome-icon
            icon="fa-solid fa-arrows-rotate"
            :class="{
              'animate-spin ': isPersistentVolumesLoading
            }" />&nbsp;&nbsp; Refresh List
        </FilledButton>
      </template>
    </PageBar>

    <!-- Table -->
    <Table class="mt-8">
      <template v-slot:header>
        <TableHeader align="left">Volume Name</TableHeader>
        <TableHeader align="center">Details</TableHeader>
        <TableHeader align="center">Size</TableHeader>
        <TableHeader align="center">PV Backup</TableHeader>
        <TableHeader align="center">PV Restore</TableHeader>
        <TableHeader align="right">Actions</TableHeader>
      </template>
      <template v-slot:message>
        <TableMessage v-if="persistentVolumes.length === 0">
          No persistent volumes found.<br />
          Click on the "Add New" button to create a new persistent volume.
        </TableMessage>
      </template>
      <template v-slot:body>
        <PersistentVolumeRow
          :delete-persistent-volume-with-confirmation="deletePersistentVolumeWithConfirmation"
          v-for="volume in persistentVolumes"
          :key="volume.id"
          :show-backups="() => openBackupDrawerForVolume(volume.id, volume.name)"
          :show-restores="() => openRestoreDrawerForVolume(volume.id, volume.name)"
          :restore-now="() => openRestoreNowModal(volume.id, volume.name)"
          :show-details="() => showDetails(volume)"
          :volume="volume" />
      </template>
    </Table>
  </section>
</template>

<style scoped></style>
