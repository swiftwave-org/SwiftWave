<script setup>
import { computed } from 'vue'

const props = defineProps({
  loading: {
    default: false
  },
  disabled: {
    default: false
  },
  click: {
    type: Function,
    default: () => {}
  },
  type: {
    type: String,
    default: 'primary',
    validator(value) {
      return ['primary', 'success', 'warning', 'danger', 'secondary'].includes(value)
    }
  },
  slim: {
    type: Boolean,
    default: false
  },
  rounded: {
    type: Boolean,
    default: false
  }
})

const isDisabled = computed(() => {
  return props.disabled || props.loading
})
</script>

<template>
  <button
    :class="{
      'border-primary-600 hover:shadow-md focus-visible:outline-primary-700': type === 'primary',
      'border-secondary-600 hover:shadow-md focus-visible:outline-secondary-700': type === 'secondary',
      'border-success-600 hover:shadow-md focus-visible:outline-success-700': type === 'success',
      'border-warning-600 hover:shadow-md focus-visible:outline-warning-700': type === 'warning',
      'border-danger-600 hover:shadow-md focus-visible:outline-danger-700': type === 'danger',
      'text-primary-600 hover:shadow-md focus-visible:outline-primary-700': type === 'primary',
      'text-secondary-600 hover:shadow-md focus-visible:outline-secondary-700': type === 'secondary',
      'text-success-600 hover:shadow-md focus-visible:outline-success-700': type === 'success',
      'text-warning-600 hover:shadow-md focus-visible:outline-warning-700': type === 'warning',
      'text-danger-600 hover:shadow-md focus-visible:outline-danger-700': type === 'danger',
      'hover:border-[type]-600/80': !disabled,
      'hover:text-[type]-600/80': !disabled,
      'cursor-not-allowed opacity-50': disabled,
      'cursor-progress': loading,
      'py-2 text-sm': !slim,
      'px-2 py-1 text-xs': slim,
      'rounded-full': rounded,
      'rounded-md': !rounded
    }"
    :disabled="isDisabled"
    class="flex items-center justify-center border-2 px-3 font-semibold shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
    type="button"
    @click="click">
    <!--    spinner -->
    <svg
      v-if="loading"
      :class="{
        'h-5 w-5': !slim,
        'h-3 w-3': slim,
        'text-primary-700': type === 'primary',
        'text-secondary-700': type === 'secondary',
        'text-success-700': type === 'success',
        'text-warning-700': type === 'warning',
        'text-danger-700': type === 'danger',
        'text-info-700': type === 'info'
      }"
      class="-ml-1 mr-3 animate-spin"
      fill="none"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path
        class="opacity-75"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        fill="currentColor"></path>
    </svg>
    <!-- linear wave spinner -->

    <!-- text -->
    <slot></slot>
  </button>
</template>

<style scoped></style>
