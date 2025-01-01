<script setup>
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, ref, watch } from 'vue'

const props = defineProps({
  selector: {
    type: String,
    required: false,
    default: 'id',
    validator: (value) => ['id', 'name'].includes(value)
  },
  onSelect: {
    type: Function,
    required: false,
    default: () => {}
  },
  showCreateLink: {
    type: Boolean,
    required: false,
    default: false
  },
  fullWidth: {
    type: Boolean,
    required: false,
    default: false
  }
})

const selectedApplication = ref('')
const { result: applicationsResult, refetch: refetchApplications } = useQuery(
  gql`
    query {
      applications {
        id
        name
      }
    }
  `,
  null,
  {
    pollInterval: 0
  }
)

const applications = computed(() => applicationsResult.value?.applications ?? [])
const getValue = (application) => application[props.selector]

defineExpose({
  refetch: refetchApplications
})

watch(selectedApplication, () => {
  if (!selectedApplication.value) return
  props.onSelect(selectedApplication.value)
})

const createApplication = () => {
  alert("You can create application from 'Deploy Application', 'Deploy Stack' or 'App Store' options of sidebar")
}
</script>

<template>
  <div
    :class="{
      'w-full': fullWidth
    }">
    <select
      class="lock w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
      v-model="selectedApplication">
      <option value="" disabled selected>Select Application</option>
      <option v-for="application in applications" :key="application.id" :value="getValue(application)">
        {{ application.name }}
      </option>
    </select>
    <p class="ml-1 mt-2 flex items-center text-sm" v-if="showCreateLink">
      Need to create Application ?
      <a class="ml-1.5 cursor-pointer font-bold text-primary-600" @click="createApplication">Click here</a>
    </p>
  </div>
</template>

<style scoped></style>
