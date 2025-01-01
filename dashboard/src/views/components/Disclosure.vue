<script setup>
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/vue'
</script>

<template>
  <div class="w-full">
    <div class="mx-auto w-full rounded-2xl bg-white">
      <Disclosure v-slot="{ open }" as="div" class="mt-2">
        <DisclosureButton
          class="flex w-full justify-between rounded-lg bg-secondary-100 px-4 py-2 text-left text-sm font-medium text-secondary-900 hover:bg-secondary-200 focus:outline-none focus-visible:ring focus-visible:ring-secondary-500/75">
          <span><slot name="title"></slot></span>
          <font-awesome-icon
            :class="{
              'rotate-180 transform': open
            }"
            class="h-5 w-5 text-secondary-500"
            icon="fa-solid fa-chevron-down" />
        </DisclosureButton>
        <DisclosurePanel class="px-4 pb-2 pt-4 text-sm">
          <slot name="body"></slot>
        </DisclosurePanel>
      </Disclosure>
    </div>
  </div>
</template>
