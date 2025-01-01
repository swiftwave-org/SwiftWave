<script setup>
import { computed } from 'vue'

const props = defineProps({
  percentage: {
    type: Number,
    required: true
  },
  label: {
    type: String,
    required: true
  },
  hideLabel: {
    type: Boolean,
    default: false
  },
  hideHover: {
    type: Boolean,
    default: false
  },
  small: {
    type: Boolean,
    default: false
  },
  labelPosition: {
    type: String,
    default: 'bottom',
    validator: (value) => ['top', 'bottom', 'left', 'right'].includes(value)
  }
})

const percentage = computed(() => Math.min(props.percentage, 100))
const noOfGreenLines = computed(() => Math.floor(percentage.value / 10))
const noOfRedLines = computed(() => 10 - noOfGreenLines.value)
</script>

<template>
  <div class="group relative w-max">
    <div
      class="flex"
      :class="{
        'flex-col items-center': labelPosition === 'top' || labelPosition === 'bottom',
        'flex-row items-center': labelPosition === 'left' || labelPosition === 'right'
      }">
      <p v-if="!hideLabel && labelPosition === 'top'" class="mb-0.5 text-sm text-secondary-600">{{ label }}</p>
      <p v-if="!hideLabel && labelPosition === 'left'" class="mr-2 text-sm text-secondary-600">{{ label }}</p>
      <div class="flex flex-row gap-1">
        <div
          v-for="i in noOfGreenLines"
          :key="i"
          class="rounded-lg bg-green-500"
          :class="{
            'h-6 w-1.5': !small,
            'h-4 w-[0.25rem]': small
          }"></div>
        <div
          v-for="i in noOfRedLines"
          :key="i"
          class="rounded-lg bg-red-400"
          :class="{
            'h-6 w-1.5': !small,
            'h-4 w-[0.25rem]': small
          }"></div>
      </div>
      <p v-if="!hideLabel && labelPosition === 'bottom'" class="mt-0.5 text-sm text-secondary-600">{{ label }}</p>
      <p v-if="!hideLabel && labelPosition === 'right'" class="ml-2 text-sm text-secondary-600">{{ label }}</p>
    </div>
    <span
      v-if="!hideHover"
      class="pointer-events-none absolute top-0 w-max opacity-0 transition-opacity group-hover:opacity-100"
      :class="{
        'left-28 text-base': !small,
        'left-20 text-sm': small
      }">
      {{ percentage }}%
    </span>
  </div>
</template>

<style scoped></style>
