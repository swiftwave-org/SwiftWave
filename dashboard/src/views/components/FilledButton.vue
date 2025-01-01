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
      return ['primary', 'success', 'warning', 'danger', 'secondary', 'info', 'ghost'].includes(value)
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
      'bg-primary-600 hover:bg-primary-600/80 focus-visible:outline-primary-600': type === 'primary',
      'bg-secondary-600 hover:bg-secondary-600/80 focus-visible:outline-secondary-600': type === 'secondary',
      'bg-success-600 hover:bg-success-600/80 focus-visible:outline-success-600': type === 'success',
      'bg-warning-600 hover:bg-warning-600/80 focus-visible:outline-warning-600': type === 'warning',
      'bg-danger-600 hover:bg-danger-600/80 focus-visible:outline-danger-600': type === 'danger',
      'bg-info-600 hover:bg-info-600/80 focus-visible:outline-info-600': type === 'info',
      'text-black-500 bg-secondary-100 font-medium hover:bg-secondary-200 focus-visible:outline-secondary-400':
        type === 'ghost',
      'font-semibold text-white': type !== 'ghost',
      'cursor-not-allowed opacity-50': disabled,
      'hover:bg-[type]-600/80': !disabled,
      'cursor-progress': loading,
      'px-3 py-2 text-sm': !slim,
      'px-2 py-1 text-xs': slim,
      'rounded-full': rounded,
      'rounded-md': !rounded
    }"
    :disabled="isDisabled"
    class="flex items-center justify-center shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
    type="button"
    @click.stop="click">
    <!--    spinner -->
    <svg
      v-if="loading"
      :class="{
        'h-5 w-5': !slim,
        'h-3 w-3': slim,
        'text-white': type !== 'ghost'
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
