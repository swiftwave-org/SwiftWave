<script setup>
import { RouterView, useRouter } from 'vue-router'
import { computed, onBeforeMount, onMounted } from 'vue'
import { useAuthStore } from '@/store/auth.js'
import SideBar from '@/views/partials/SideBar.vue'
import LoadingPage from '@/views/pages/LoadingPage.vue'
import NotAvailableOnMobile from '@/views/pages/NotAvailableOnMobile.vue'
import GlobalWarning from '@/views/partials/GlobalWarning.vue'
import { Toaster } from 'vue-sonner'

const authStore = useAuthStore()
const router = useRouter()

onBeforeMount(() => {
  const token = localStorage.getItem('token')
  if (token !== null) {
    authStore.SetCredential(token)
  }
})

onMounted(() => {
  authStore.StartAuthChecker(() => {
    authStore.Logout()
  })
})

const isLoginPage = computed(() => router.currentRoute.value.name === 'Login')
</script>

<template>
  <Toaster position="top-center" theme="light" richColors />
  <LoadingPage :show="authStore.IsLoggingInProgress" />
  <div class="app">
    <SideBar class="w-80" />
    <div
      class="scrollbox flex max-h-[100vh] w-full flex-col items-center overflow-y-auto"
      :class="{
        'p-4': !isLoginPage
      }">
      <RouterView />
    </div>
    <GlobalWarning />
  </div>
  <NotAvailableOnMobile />
</template>

<style>
.app {
  display: flex;
  flex-direction: row;
  height: 100vh;
  width: 100vw;
}

.scrollbox::-webkit-scrollbar {
  width: 9px;
  height: 9px;
}

.scrollbox::-webkit-scrollbar-track {
  @apply rounded-full bg-gray-200;
}

.scrollbox::-webkit-scrollbar-thumb {
  @apply rounded-full bg-primary-500;
}

.terminal {
  margin-right: calc(0.5rem + 9px);
}

.xterm-viewport {
  right: calc(-0.5rem - 9px) !important;
  cursor: pointer !important;
  overflow-y: auto !important;
}

.xterm-viewport::-webkit-scrollbar {
  width: 9px !important;
  height: 9px !important;
}

.xterm-viewport::-webkit-scrollbar-track {
  @apply rounded-full bg-gray-200;
}

.xterm-viewport::-webkit-scrollbar-thumb {
  @apply rounded-full bg-primary-500;
}

.bg-color-1 {
  @apply bg-amber-300 !important;
}

.border-color-1 {
  @apply border-amber-400 !important;
}

.bg-color-2 {
  @apply bg-green-300 !important;
}

.border-color-2 {
  @apply border-green-400 !important;
}

.bg-color-3 {
  @apply bg-blue-300 !important;
}

.border-color-3 {
  @apply border-blue-400 !important;
}

.bg-color-4 {
  @apply bg-gray-300 !important;
}

.border-color-4 {
  @apply border-gray-400 !important;
}

/* v-otp-input custom class */
.otp-input-container {
  @apply flex flex-row justify-center;
}

.otp-input {
  @apply mx-1 h-10 w-10 !rounded-md !border !border-gray-300 !p-1 !text-center !text-base focus:!border-primary-500 focus:!ring-primary-500;
}

.otp-input::-webkit-inner-spin-button,
.otp-input::-webkit-outer-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

/* tooltip */
.tooltip {
  @apply invisible absolute -ms-12 -mt-9 rounded bg-black px-2 py-1 text-sm text-white shadow-sm;
}

.has-tooltip:hover .tooltip {
  @apply visible z-50;
}

/* popover */
.popover {
  @apply absolute left-1/2 top-[calc(100%+15px)] z-50 hidden max-w-[300px] -translate-x-1/2 flex-col items-center gap-1.5 rounded-lg bg-primary-500 p-2.5 text-sm text-secondary-100 shadow-md;
}

.has-popover {
  @apply relative;

  &:hover .popover {
    @apply block;
  }
}
</style>
