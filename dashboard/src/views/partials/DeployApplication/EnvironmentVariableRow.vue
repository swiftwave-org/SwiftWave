<script setup>
import FilledButton from '@/views/components/FilledButton.vue'
import TableRow from '@/views/components/Table/TableRow.vue'
import { ref } from 'vue'

defineProps({
  variableKey: {
    type: String,
    required: true
  },
  variableName: {
    type: String,
    required: true
  },
  variableValue: {
    type: String,
    required: true
  },
  onVariableNameChange: {
    type: Function,
    required: true
  },
  onVariableValueChange: {
    type: Function,
    required: true
  },
  deleteVariable: {
    type: Function,
    required: true
  }
})

const valueHidden = ref(true)
</script>

<template>
  <tr :key="variableKey">
    <TableRow class="!py-3">
      <input
        :key="`name-${variableKey}`"
        class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
        placeholder="Environment Variable Name"
        type="text"
        v-bind:value="variableName"
        @input="(event) => onVariableNameChange(variableKey, event.target.value)" />
    </TableRow>
    <TableRow class="!py-3">
      <div class="relative w-full">
        <input
          :key="`value-${variableKey}`"
          class="block w-full rounded-md border-gray-300 pe-14 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
          placeholder="Environment Variable Value"
          :type="valueHidden ? 'password' : 'text'"
          v-bind:value="variableValue"
          @input="(event) => onVariableValueChange(variableKey, event.target.value)" />
        <div
          class="absolute inset-y-0 right-0 flex h-full w-12 cursor-pointer items-center px-1 py-1"
          @click.stop="valueHidden = !valueHidden">
          <div
            class="flex h-full w-full items-center justify-center rounded-sm bg-gray-200 text-black transition-all duration-200 hover:bg-secondary-300">
            <font-awesome-icon v-if="valueHidden" icon="fa-solid fa-eye-slash" />
            <font-awesome-icon v-else icon="fa-solid fa-eye" />
          </div>
        </div>
      </div>
    </TableRow>
    <TableRow align="right" class="flex !py-3">
      <FilledButton
        :key="`delete-${variableKey}`"
        :click="() => deleteVariable(variableKey)"
        class="w-full"
        type="danger">
        <font-awesome-icon class="mr-2" icon="fa-solid fa-trash" />
        Delete
      </FilledButton>
    </TableRow>
  </tr>
</template>
