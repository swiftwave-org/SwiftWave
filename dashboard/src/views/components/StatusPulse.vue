<script setup>
defineProps({
  type: {
    type: String,
    default: 'success',
    validator: (value) => {
      return ['primary', 'secondary', 'success', 'danger', 'warning'].includes(value)
    }
  },
  animated: {
    type: Boolean,
    default: true
  }
})
</script>

<template>
  <span class="relative flex h-3.5 w-3.5 items-center justify-center">
    <span
      class="absolute inline-flex h-full w-full rounded-full opacity-75"
      :class="{
        'bg-primary-400': type === 'primary',
        'bg-secondary-400': type === 'secondary',
        'bg-success-400': type === 'success',
        'bg-danger-400': type === 'danger',
        'bg-warning-400': type === 'warning',
        'animate-ping': animated
      }"></span>
    <span
      class="relative inline-flex h-3 w-3 rounded-full"
      :class="{
        'bg-primary-500': type === 'primary',
        'bg-secondary-500': type === 'secondary',
        'bg-success-500': type === 'success',
        'bg-danger-500': type === 'danger',
        'bg-warning-500': type === 'warning'
      }"></span>
  </span>
</template>

<style scoped></style>
