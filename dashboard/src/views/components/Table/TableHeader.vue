<script setup>
defineProps({
  align: {
    type: String,
    default: 'center',
    validator: (value) => {
      return ['left', 'center', 'right'].includes(value)
    }
  }
})
</script>

<template>
  <th
    :class="{
      'text-left': align === 'left',
      'text-center': align === 'center',
      'text-right': align === 'right'
    }"
    class="px-4 py-3.5 text-sm font-normal text-gray-700"
    scope="col">
    <span>
      <slot></slot>
    </span>
  </th>
</template>
