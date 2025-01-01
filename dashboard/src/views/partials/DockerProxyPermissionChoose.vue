<script setup>
defineProps({
  value: {
    type: String,
    default: ''
  },
  onChange: {
    type: Function,
    default: () => {}
  }
})
</script>

<template>
  <div class="multi-select">
    <div
      @click="onChange('none')"
      :class="{
        active: value === 'none'
      }">
      None
    </div>
    <div
      @click="onChange('read')"
      :class="{
        active: value === 'read'
      }">
      Read
    </div>
    <div
      @click="onChange('read_write')"
      :class="{
        active: value === 'read_write'
      }">
      Write
    </div>
  </div>
</template>

<style scoped>
.multi-select {
  @apply flex h-fit w-min overflow-hidden rounded-md border border-secondary-400;

  div {
    @apply cursor-pointer border-r border-secondary-400 bg-secondary-100 px-2.5 py-0.5 text-sm transition-all;
  }

  div:last-child {
    @apply border-0;
  }

  .active {
    @apply bg-primary-600  text-white;
  }
}
</style>
