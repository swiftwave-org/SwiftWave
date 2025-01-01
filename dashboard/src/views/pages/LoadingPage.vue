<script setup>
import logo from '@/assets/images/logo-full.png'

defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

</script>

<template>
<div class="loading-page" v-show="show">
  <img alt="swiftwave brand logo" :src="logo" class="h-20" />
  <div class="spinner-box">
    <div class="pulse-container">
      <div class="pulse-bubble pulse-bubble-1"></div>
      <div class="pulse-bubble pulse-bubble-2"></div>
      <div class="pulse-bubble pulse-bubble-3"></div>
    </div>
  </div>
</div>
</template>

<style scoped>
.loading-page {
  @apply w-[100vw] h-[100vh] flex flex-col justify-center items-center absolute bg-white z-50
}

@keyframes pulse {
  from {
    opacity: 1;
    transform: scale(1);
  }
  to {
    opacity: .25;
    transform: scale(.75);
  }
}

.spinner-box {
  width: 300px;
  height: 160px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: transparent;
}

.pulse-container {
  width: 120px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pulse-bubble {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: #6B46C1;
}

.pulse-bubble-1 {
  animation: pulse .4s ease 0s infinite alternate;
}

.pulse-bubble-2 {
  animation: pulse .4s ease .2s infinite alternate;
}

.pulse-bubble-3 {
  animation: pulse .4s ease .4s infinite alternate;
}

</style>