<script setup>
defineProps({
  type: {
    type: String,
    default: 'success',
    validator: (value) => {
      return ['primary', 'secondary', 'success', 'danger', 'warning'].includes(value)
    }
  },
  small: {
    type: Boolean,
    default: false
  }
})
</script>
<template>
  <div
    class="flex items-center p-0"
    :class="{
      'gap-3': !small,
      'gap-2': small,
      'text-sm': small
    }">
    <div
      class="rounded-full"
      :class="{
        'h-3 w-3': !small,
        'h-2.5 w-2.5': small,
        'bg-primary-500': type === 'primary',
        'bg-secondary-500': type === 'secondary',
        'bg-success-500': type === 'success',
        'bg-danger-500': type === 'danger',
        'bg-warning-500': type === 'warning'
      }"></div>
    <slot />
  </div>
</template>

<style scoped></style>
