<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  subTitle: {
    type: String,
    default: ''
  },
  prefixText: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'success',
    validator: (value) => {
      return ['primary', 'secondary', 'success', 'danger', 'warning'].includes(value)
    }
  },
  showBody: {
    type: Boolean,
    default: true
  }
})
</script>

<template>
  <div class="my-3 flex items-center gap-4">
    <div
      class="flex h-12 w-12 items-center justify-center rounded-xl border text-xl font-bold text-white shadow-md"
      :class="{
        'border-primary-600 bg-primary-500': type === 'primary',
        'border-secondary-600 bg-secondary-500': type === 'secondary',
        'border-success-600 bg-success-500': type === 'success',
        'border-danger-600 bg-danger-500': type === 'danger',
        'border-warning-600 bg-warning-500': type === 'warning'
      }">
      {{ prefixText }}
    </div>
    <div class="flex flex-col">
      <div class="text-base font-medium text-gray-900">{{ title }}</div>
      <div class="text-sm text-gray-500">{{ subTitle }}</div>
    </div>
  </div>
  <div class="ml-16 mt-4" v-if="showBody">
    <slot></slot>
  </div>
</template>

<style scoped></style>
