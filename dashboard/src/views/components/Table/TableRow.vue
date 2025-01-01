<script setup>
defineProps({
  align: {
    type: String,
    required: false,
    default: 'left',
    validator: (value) => {
      return ['left', 'center', 'right'].includes(value)
    }
  },
  flex: {
    type: Boolean,
    required: false,
    default: false
  }
})
</script>

<template>
  <td class="px-4 py-4">
    <div
      :class="{
      'text-left': align === 'left',
      'text-center': align === 'center',
      'text-right': align === 'right',
      'justify-start': align === 'left',
      'justify-center': align === 'center',
      'justify-end': align === 'right',
      'flex': flex
    }"
      class="gap-2 whitespace-nowrap">
      <slot></slot>
    </div>
  </td>

</template>
