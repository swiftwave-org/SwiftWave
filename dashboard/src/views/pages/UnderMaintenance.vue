<script setup>
import logo from '@/assets/images/logo-full.png'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import { useRouter } from 'vue-router'
import { onMounted } from 'vue'

const router = useRouter()

const monitorIfServerIsUp = () => {
  const baseUrl = getHttpBaseUrl()
  const url = `${baseUrl}/healthcheck`
  const interval = setInterval(() => {
    fetch(url, {
      cache: 'no-cache'
    })
      .then((res) => {
        if (res.status === 200) {
          if (router.currentRoute.value.query.redirect) {
            router.push(router.currentRoute.value.query.redirect)
          } else {
            window.location.href = router.resolve({ name: 'Applications' }).href
          }
          clearInterval(interval)
        }
      })
      .catch(() => {
        console.log('Server is down')
      })
  }, 1000)
}

onMounted(() => {
  setTimeout(() => {
    monitorIfServerIsUp()
  }, 5000)
})
</script>

<template>
  <div class="flex h-full w-full max-w-7xl flex-col items-center justify-center gap-6 sm:px-0">
    <img alt="swiftwave brand logo" :src="logo" class="mb-8 h-10" />
    <h1 class="text-4xl font-bold text-gray-900">Under Maintenance</h1>
    <p class="text-center text-lg text-gray-600">
      Swiftwave will be online soon !<br />Will be auto redirected to the home page ASAP
    </p>
  </div>
</template>

<style scoped></style>
