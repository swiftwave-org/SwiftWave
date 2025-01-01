<script setup>
import { Switch } from '@headlessui/vue'

defineProps({
  enabled: {
    type: Boolean,
    required: false,
    default: false
  },
  onChange: {
    type: Function,
    required: false
  }
})
</script>

<template>
  <Switch
    :class="enabled ? 'bg-primary-600' : 'bg-primary-500'"
    :enabled="enabled"
    class="relative inline-flex h-[21px] w-[37px] shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75"
    @click="() => onChange(!enabled)">
    <span
      :class="enabled ? 'translate-x-4' : 'translate-x-0'"
      aria-hidden="true"
      class="pointer-events-none inline-block h-[17px] w-[17px] transform rounded-full bg-white shadow-lg ring-0 transition duration-200 ease-in-out" />
  </Switch>
</template>
