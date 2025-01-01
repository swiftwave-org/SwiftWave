<script setup>
import { TabPanel } from '@headlessui/vue'
import { ref } from 'vue'
import FilledButton from '@/views/components/FilledButton.vue'

defineProps({
  finalizeApplicationSourceAndMoveToNextTab: {
    type: Function,
    required: true
  }
})

const sourceType = ref('')
const setSourceType = (type) => {
  sourceType.value = type
}
</script>

<template>
  <TabPanel :key="1" class="flex h-full flex-col items-center justify-center space-y-6">
    <p class="text-center text-lg font-medium text-gray-500">Select the source of your application</p>
    <div class="flex flex-row items-center justify-center space-x-6">
      <!--  Docker Image  -->
      <div
        :class="{
          optionButtonSelected: sourceType === 'image'
        }"
        class="optionButton"
        @click.prevent="() => setSourceType('image')">
        <font-awesome-icon class="optionIcon" icon="fa-brands fa-docker" />
        <p>Docker Image</p>
      </div>
      <!--  Git   -->
      <div
        :class="{
          optionButtonSelected: sourceType === 'git'
        }"
        class="optionButton"
        @click.prevent="() => setSourceType('git')">
        <font-awesome-icon class="optionIcon" icon="fa-brands fa-git" />
        <p>Git Repository</p>
      </div>
      <!--  Upload Code -->
      <div
        :class="{
          optionButtonSelected: sourceType === 'sourceCode'
        }"
        class="optionButton"
        @click.prevent="() => setSourceType('sourceCode')">
        <font-awesome-icon class="optionIcon" icon="fa-solid fa-upload" />
        <p>Upload Codes</p>
      </div>
    </div>

    <FilledButton
      :disabled="sourceType === ''"
      class="w-full"
      type="primary"
      @click="() => finalizeApplicationSourceAndMoveToNextTab(sourceType)">
      Proceed to Next Page
    </FilledButton>
  </TabPanel>
</template>

<style scoped>
.optionButton {
  @apply flex aspect-square h-fit cursor-pointer flex-col items-center justify-center rounded-md border-2 border-gray-400 p-8 hover:text-blue-700 hover:shadow-lg;
}

.optionButton p {
  @apply mt-4 text-center text-base font-medium text-gray-500;
}

.optionIcon {
  @apply text-5xl text-gray-600;
}

.optionButtonSelected {
  @apply border-primary-500;
}

.optionButtonSelected p {
  @apply text-primary-600;
}

.optionButtonSelected .optionIcon {
  @apply text-primary-600;
}
</style>
