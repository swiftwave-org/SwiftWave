<script setup>

import { ref } from 'vue'

let isVisible = ref(false)

function toggleVisibility() {
  isVisible.value = !isVisible.value
}
</script>

<template>
  <div class="text-sm text-gray-900 min-w-[100px]">
    <span v-if="isVisible" class="flex items-center w-fit"><slot></slot><font-awesome-icon icon="fa-solid fa-eye" class="ml-2 cursor-pointer" :onclick="toggleVisibility" /></span>
    <span v-else class="flex items-center w-fit">********<font-awesome-icon icon="fa-solid fa-eye-slash" class="ml-2 cursor-pointer" :onclick="toggleVisibility" /></span>
  </div>
</template>

<style scoped>

</style>