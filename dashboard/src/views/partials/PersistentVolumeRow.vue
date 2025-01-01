<script setup>
import TableRow from '@/views/components/Table/TableRow.vue'
import TextButton from '@/views/components/TextButton.vue'
import { useLazyQuery, useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, ref } from 'vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { round } from 'lodash'
import { toast } from 'vue-sonner'
import ModalDialog from '@/views/components/ModalDialog.vue'
import Badge from '@/views/components/Badge.vue'

const props = defineProps({
  volume: {
    type: Object,
    required: true
  },
  deletePersistentVolumeWithConfirmation: {
    type: Function,
    required: true
  },
  showBackups: {
    type: Function,
    required: true
  },
  showRestores: {
    type: Function,
    required: true
  },
  restoreNow: {
    type: Function,
    required: true
  },
  showDetails: {
    type: Function,
    required: true
  }
})

const sizeFetched = ref(false)

const {
  result: persistentVolumeSizeResult,
  load: loadPersistentVolumeSize,
  refetch: refetchPersistentVolumeSize,
  loading: isPersistentVolumeSizeLoading,
  onResult: onPersistentVolumeSizeResult,
  onError: onPersistentVolumeSizeError
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      persistentVolumeSizeMb(id: $id)
    }
  `,
  {
    id: props.volume.id
  },
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

onPersistentVolumeSizeResult(() => {
  sizeFetched.value = true
})

onPersistentVolumeSizeError((err) => {
  toast.error(err.message)
})

const fetchPersistentVolumeSize = async () => {
  if (!(await loadPersistentVolumeSize())) {
    await refetchPersistentVolumeSize()
  }
}
const persistentVolumeSize = computed(() => round(persistentVolumeSizeResult.value?.persistentVolumeSizeMb ?? 0, 2))

// Create backup
const backupType = ref('local')
const isBackupCreateModalOpen = ref(false)
const {
  mutate: createPersistentVolumeBackup,
  loading: isBackupCreating,
  onDone: onBackupCreateSuccess,
  onError: onBackupCreateFail
} = useMutation(gql`
  mutation ($input: PersistentVolumeBackupInput!) {
    backupPersistentVolume(input: $input) {
      id
    }
  }
`)

onBackupCreateSuccess(() => {
  toast.success('Backup will be created shortly ! Check the status in the backups list section.')
})

onBackupCreateFail((err) => {
  toast.error(err.message)
})

const closeBackupCreateModal = () => {
  isBackupCreateModalOpen.value = false
}

const openBackupCreateModal = () => {
  isBackupCreateModalOpen.value = true
}

const createBackup = () => {
  closeBackupCreateModal()

  createPersistentVolumeBackup({
    input: {
      persistentVolumeId: props.volume.id,
      type: backupType.value
    }
  })
}
</script>

<template>
  <Teleport to="body">
    <ModalDialog :close-modal="closeBackupCreateModal" :is-open="isBackupCreateModalOpen">
      <template v-slot:header>Backup Volume</template>
      <template v-slot:body>
        <form @submit.prevent="createPersistentVolumeBackup">
          <!-- Domains -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700">Choose Backup Type</label>
            <select
              class="mt-2 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
              v-model="backupType">
              <option value="local">Local</option>
              <option value="s3">S3</option>
            </select>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton :click="createBackup" :loading="isBackupCreating" type="primary">Create Backup</FilledButton>
      </template>
    </ModalDialog>
  </Teleport>
  <tr :key="volume.id">
    <TableRow align="left">
      <div class="text-sm font-medium text-gray-900">
        <Badge class="mr-1" type="success" v-if="volume.type === 'local'">Local</Badge>
        <Badge class="mr-1" type="warning" v-if="volume.type === 'nfs'">&nbsp;&nbsp;NFS&nbsp;&nbsp;</Badge>
        <Badge class="mr-1" type="warning" v-if="volume.type === 'cifs'">&nbsp;CIFS&nbsp;</Badge>
        {{ volume.name }}
      </div>
    </TableRow>
    <TableRow align="center">
      <FilledButton slim type="secondary" :click="showDetails"> Show Details</FilledButton>
    </TableRow>
    <TableRow align="center" flex>
      <p v-show="sizeFetched">{{ persistentVolumeSize }} MB</p>
      <FilledButton :click="fetchPersistentVolumeSize" slim type="secondary" :loading="isPersistentVolumeSizeLoading"
        >{{ sizeFetched ? 'Refresh' : 'Fetch Size' }}
      </FilledButton>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton slim type="primary" :loading="isBackupCreating" :click="openBackupCreateModal">
        <font-awesome-icon icon="fa-solid fa-plus" class="mr-2" />
        Create Backup
      </FilledButton>
      <FilledButton slim type="secondary" :click="showBackups">
        <font-awesome-icon icon="fa-solid fa-list-check" class="mr-2" />
        Show Backups
      </FilledButton>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton slim type="primary" :click="restoreNow">
        <font-awesome-icon icon="fa-solid fa-upload" class="mr-2" />
        Restore Now
      </FilledButton>
      <FilledButton slim type="secondary" :click="showRestores">
        <font-awesome-icon icon="fa-solid fa-list-check" class="mr-2" />
        Restore History
      </FilledButton>
    </TableRow>
    <TableRow align="right">
      <TextButton :click="() => deletePersistentVolumeWithConfirmation(volume)" type="danger">Delete</TextButton>
    </TableRow>
  </tr>
</template>

<style scoped></style>
