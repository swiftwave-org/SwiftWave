<script setup>
import Drawer from '@/views/components/Drawer.vue'
import { ref, watch } from 'vue'
import { useLazyQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import Badge from '@/views/components/Badge.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import DotLoader from '@/views/components/DotLoader.vue'

const props = defineProps({
  isDrawerOpen: Boolean,
  closeDrawer: Function,
  persistentVolumeId: Number,
  persistentVolumeName: String
})

const persistentVolumeNameRef = ref(props.persistentVolumeName)
const restores = ref([])

watch(
  () => props.persistentVolumeName,
  (newValue) => {
    persistentVolumeNameRef.value = newValue
  }
)
watch(
  () => props.isDrawerOpen,
  (newValue) => {
    if (newValue) {
      restores.value = []
      fetchPersistentVolumeRestores()
    }
  }
)

// Fetch Persistent Volume Restores
const {
  load: loadPersistentVolumeRestores,
  refetch: refetchPersistentVolumeRestores,
  loading: isPersistentVolumeRestoresLoading,
  onError: onPersistentVolumeRestoresError,
  variables: persistentVolumeRestoresVariables,
  onResult: onPersistentVolumeRestoresResult
} = useLazyQuery(
  gql`
    query ($id: Uint!) {
      persistentVolume(id: $id) {
        restores {
          id
          type
          status
          createdAt
          completedAt
        }
      }
    }
  `,
  null,
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

const fetchPersistentVolumeRestores = async () => {
  persistentVolumeRestoresVariables.value = { id: props.persistentVolumeId }
  if (!(await loadPersistentVolumeRestores())) {
    await refetchPersistentVolumeRestores()
  }
}

onPersistentVolumeRestoresResult((result) => {
  restores.value = result.data.persistentVolume.restores
})

onPersistentVolumeRestoresError((err) => {
  toast.error(err.message)
})
</script>

<template>
  <Drawer :close-drawer="closeDrawer" :is-open="isDrawerOpen" :width-size="3">
    <template v-slot:title>Restore History</template>
    <template v-slot:subtitle> Manage restore history for volume {{ persistentVolumeNameRef }}</template>
    <template v-slot:body>
      <FilledButton slim class="mb-2 w-full" type="secondary" :click="fetchPersistentVolumeRestores">
        <font-awesome-icon icon="fa-solid fa-rotate-right" class="mr-2" />
        Refresh List
      </FilledButton>
      <div class="mb-2 text-sm text-gray-500">{{ restores.length }} restores found for this volume</div>
      <div class="my-4 flex justify-center" v-if="isPersistentVolumeRestoresLoading">
        <DotLoader />
      </div>
      <div class="flex flex-col space-y-2" v-if="!isPersistentVolumeRestoresLoading">
        <div v-for="restore in restores" :key="restore.id" class="rounded-md border-2 border-gray-200 p-2 text-sm">
          <p>
            <Badge v-if="restore.status === 'pending'" type="warning">Pending</Badge>
            <Badge v-if="restore.status === 'failed'" type="danger">Failed</Badge>
            <Badge v-if="restore.status === 'success'" type="success">Success</Badge>
            <span class="ml-2 mr-1"><b>Initiated at</b> {{ new Date(restore.createdAt).toLocaleString() }}</span>
          </p>
          <div v-if="restore.status === 'success'" class="mt-1">
            Completed at {{ new Date(restore.completedAt).toLocaleString() }}
          </div>
        </div>
      </div>
    </template>
  </Drawer>
</template>

<style scoped></style>
