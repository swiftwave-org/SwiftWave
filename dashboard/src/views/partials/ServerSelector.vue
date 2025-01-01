<script setup>
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, onMounted, ref, watch } from 'vue'

const props = defineProps({
  initialValue: {
    type: String,
    required: false,
    default: ''
  },
  onSelect: {
    type: Function,
    required: false,
    default: () => {}
  },
  fullWidth: {
    type: Boolean,
    required: false,
    default: false
  }
})

const selectedServerHostname = ref('')
const { result: serversRawResult } = useQuery(
  gql`
    query {
      servers {
        hostname
      }
    }
  `,
  null,
  {
    pollInterval: 60000,
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

const servers = computed(() => serversRawResult.value?.servers ?? [])

watch(selectedServerHostname, (newVal) => {
  if (!newVal) return
  props.onSelect(newVal)
})

onMounted(() => {
  if (props.initialValue) {
    selectedServerHostname.value = props.initialValue
  }
})
</script>

<template>
  <div
    :class="{
      'w-full': props.fullWidth
    }">
    <select
      class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
      v-model="selectedServerHostname">
      <option value="" disabled selected>Select Server</option>
      <option v-for="server in servers" :key="server.hostname" :value="server.hostname">
        {{ server.hostname }}
      </option>
    </select>
  </div>
</template>

<style scoped></style>
