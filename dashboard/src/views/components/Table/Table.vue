<script setup>
defineProps({
  dividerBetweenRows: {
    type: Boolean,
    required: false,
    default: true
  }
})
</script>
<template>
  <div class="scrollbox max-h-[80vh] overflow-y-auto border border-gray-200 md:rounded-lg">
    <table class="min-w-full divide-y divide-gray-200">
      <thead class="bg-gray-50">
        <tr>
          <slot name="header"></slot>
        </tr>
      </thead>
      <tbody
        class="bg-white"
        :class="{
          'divide-y-0': dividerBetweenRows,
          'divide-y divide-gray-200': dividerBetweenRows
        }">
        <slot name="body"></slot>
      </tbody>
    </table>
    <slot name="message"></slot>
  </div>
</template>
