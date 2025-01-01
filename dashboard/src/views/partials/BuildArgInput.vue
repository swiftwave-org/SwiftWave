<script setup>
defineProps({
  argKey: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  value: {
    type: String,
    required: true
  },
  updateBuildArg: {
    type: Function,
    required: true
  }
})
</script>

<template>
  <div class="mt-4">
    <label :for="argKey" class="block text-sm font-medium text-gray-700"
      >{{ argKey }}

      <br v-if="description !== ''" />
      <span class="text-xs font-normal text-gray-800">{{ description }}</span>
    </label>
    <div class="mt-2">
      <input
        :id="argKey"
        :name="argKey"
        :value="value"
        autocomplete="off"
        class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
        placeholder="{{ description }}"
        type="text"
        @input="(e) => updateBuildArg(e.target.value)" />
    </div>
  </div>
</template>

<style scoped></style>
