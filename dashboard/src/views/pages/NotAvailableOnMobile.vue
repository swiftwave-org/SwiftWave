<script setup>
import logo from '@/assets/images/logo-full.png'
</script>

<template>
<teleport to="body">
  <div class="fixed top-0 left-0 w-full h-full bg-white flex flex-col justify-center items-center text-center z-50 md:hidden">
    <img alt="swiftwave brand logo" :src="logo" class="h-10 mb-5" />
    <h1 class="text-2xl font-bold mx-2">Not Available On Mobile</h1>
    <p class="text-gray-500 mx-2">Please use tablet or desktop browser to access this page.</p>
  </div>
</teleport>
</template>

<style scoped>

</style>