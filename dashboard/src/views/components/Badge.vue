<script setup>
defineProps({
  type: {
    type: String,
    default: 'success',
    validator: (value) => {
      return ['primary', 'secondary', 'success', 'danger', 'warning'].includes(value)
    }
  }
})
</script>

<template>
  <span
    :class="{
      'bg-primary-100 text-primary-800': type === 'primary',
      'bg-secondary-100 text-secondary-800': type === 'secondary',
      'bg-success-100 text-success-800': type === 'success',
      'bg-danger-100 text-danger-800': type === 'danger',
      'bg-warning-100 text-warning-800': type === 'warning'
    }"
    class="inline-flex rounded-full px-2 text-xs font-semibold leading-5">
    <slot></slot>
  </span>
</template>

<style scoped></style>
