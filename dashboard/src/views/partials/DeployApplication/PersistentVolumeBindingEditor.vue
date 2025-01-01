<script setup>
import Table from '@/views/components/Table/Table.vue'
import TableMessage from '@/views/components/Table/TableMessage.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import TableRow from '@/views/components/Table/TableRow.vue'
import { computed, ref, toRef } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import CreatePersistentVolumeModal from '@/views/partials/CreatePersistentVolumeModal.vue'

const props = defineProps({
  persistentVolumeBindingKeys: {
    type: Array,
    required: true
  },
  persistentVolumeBindingsMap: {
    type: Object,
    required: true
  },
  addPersistentVolumeBinding: {
    type: Function,
    required: true
  },
  deletePersistentVolumeBinding: {
    type: Function,
    required: true
  },
  onPersistentVolumeChange: {
    type: Function,
    required: true
  },
  onMountingPathChange: {
    type: Function,
    required: true
  }
})

const persistentVolumeBindingKeys = toRef(props, 'persistentVolumeBindingKeys')

const {
  result: persistentVolumesResult,
  onError: onPersistentVolumesError,
  refetch: refetchPersistentVolumes
} = useQuery(
  gql`
    query {
      persistentVolumes {
        id
        name
      }
    }
  `,
  null,
  {
    pollInterval: 10000
  }
)

onPersistentVolumesError((err) => {
  toast.error(err.message)
})

const persistentVolumes = computed(() => persistentVolumesResult.value?.persistentVolumes ?? [])

// Create persistent volume
const createPersistentVolumeModalRef = ref(null)
const openPersistentVolumeModal = computed(() => createPersistentVolumeModalRef.value?.openModal ?? (() => {}))
</script>

<template>
  <CreatePersistentVolumeModal :callback-on-create="refetchPersistentVolumes" ref="createPersistentVolumeModalRef" />
  <Table>
    <template v-slot:header>
      <TableHeader align="center">Variable Name</TableHeader>
      <TableHeader align="center">Value</TableHeader>
      <TableHeader align="right" class="w-[80px]">Delete</TableHeader>
    </template>
    <template v-slot:message>
      <TableMessage v-if="persistentVolumeBindingKeys.length === 0" class="flex flex-col items-center">
        No Persistent Volume Bindings found.<br />
        If your application requires persistent volume bindings, you can add them here.<br />
        <FilledButton class="mt-3 max-w-fit" @click="addPersistentVolumeBinding"
          >Add Persistent Volume Binding
        </FilledButton>
      </TableMessage>
      <div v-else class="flex flex-col gap-2 px-6 py-2 text-sm text-gray-600">
        <p class="m-0 inline-flex items-center p-0">
          <FilledButton slim @click="addPersistentVolumeBinding" class="mr-2"
            >Add Persistent Volume Binding
          </FilledButton>
          Want to add more persistent volume bindings ?
        </p>
        <p class="inline-flex items-center">
          <FilledButton slim @click="openPersistentVolumeModal" class="mr-2">Create New Persistent Volume</FilledButton>
          Need a new persistent volume ?
        </p>
      </div>
    </template>
    <template v-slot:body>
      <tr v-for="key in persistentVolumeBindingKeys" :key="key">
        <TableRow>
          <select
            :key="`credential-${key}`"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            @change="(event) => onPersistentVolumeChange(key, event.target.value)">
            <option selected value="0">Select Persistent Volume</option>
            <option
              v-for="volume in persistentVolumes"
              :key="volume.id"
              :selected="volume.id === persistentVolumeBindingsMap[key]?.persistentVolumeID"
              :value="volume.id">
              {{ volume.name }}
            </option>
          </select>
        </TableRow>
        <TableRow>
          <input
            :key="`value-${key}`"
            class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
            placeholder="Mounted Path"
            type="text"
            v-bind:value="persistentVolumeBindingsMap[key]?.mountingPath ?? ''"
            @input="(event) => onMountingPathChange(key, event.target.value)" />
        </TableRow>
        <TableRow align="right" class="flex">
          <FilledButton
            :key="`delete-${key}`"
            :click="() => deletePersistentVolumeBinding(key)"
            class="w-full"
            type="danger">
            <font-awesome-icon class="mr-2" icon="fa-solid fa-trash" />
            Delete Volume
          </FilledButton>
        </TableRow>
      </tr>
    </template>
  </Table>
</template>

<style scoped></style>
