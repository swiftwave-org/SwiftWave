<script setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
import { useRouter } from 'vue-router'
import { ref, watch } from 'vue'

const props = defineProps({
  activeUrls: {
    type: Array,
    default: () => []
  }
})
const router = useRouter()
const disclosureBtnRef = ref(null)

watch(
  () => router.currentRoute.value,
  (currentRoute) => {
    if (props.activeUrls.includes(currentRoute.name)) {
      openMenu()
    } else {
      closeMenu()
    }
  }
)

const openMenu = () => {
  if (disclosureBtnRef.value) {
    if (disclosureBtnRef.value.classList.contains('is-active')) {
      return
    }
    const child = disclosureBtnRef.value?.firstChild
    if (child) {
      child.click()
    }
  }
}

const closeMenu = () => {
  if (disclosureBtnRef.value) {
    if (!disclosureBtnRef.value.classList.contains('is-active')) {
      return
    }
    const child = disclosureBtnRef.value?.firstChild
    if (child) {
      child.click()
    }
  }
}
</script>

<template>
  <Disclosure v-slot="{ open }">
    <div
      ref="disclosureBtnRef"
      class="mt-2.5"
      :class="{
        'is-active': open
      }">
      <DisclosureButton
        as="button"
        class="flex w-full items-center justify-between rounded-lg bg-gray-200 bg-opacity-20 bg-clip-padding px-3 py-1.5 text-gray-200 backdrop-blur-sm backdrop-filter">
        <div>
          <slot name="icon"></slot>
          <span class="mx-2 text-sm font-medium">
            <slot name="title"></slot>
          </span>
        </div>
        <font-awesome-icon
          icon="fa-solid fa-chevron-right"
          class="transform transition-all duration-200"
          :class="{
            'rotate-90': open
          }" />
      </DisclosureButton>
    </div>

    <div class="mt-2.5">
      <transition
        enter-active-class="transition duration-50 ease-in"
        enter-from-class="transform opacity-0"
        enter-to-class="transform opacity-100"
        leave-active-class="transition duration-50 ease-out"
        leave-from-class="transform opacity-100"
        leave-to-class="transform opacity-0">
        <DisclosurePanel class="mt-0">
          <div class="ml-4 mr-0 mt-0 transition-all">
            <slot name="content"></slot>
          </div>
        </DisclosurePanel>
      </transition>
    </div>
  </Disclosure>
</template>

<style scoped>
.router-link-exact-active {
  @apply bg-gray-100 text-gray-700;
}
</style>
