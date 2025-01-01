<script setup>
import { toast } from 'vue-sonner'
import { ref } from 'vue'

defineProps({
  showCopyButton: {
    type: Boolean,
    default: true
  }
})

const showCopyBorder = ref(false)
const textDivRef = ref(null)
const copyToClipboard = () => {
  if (textDivRef.value === null) {
    return
  }
  let isSuccess
  if ('clipboard' in navigator) {
    navigator.clipboard.writeText(textDivRef.value.innerText)
    isSuccess = true
  } else {
    const textArea = document.createElement('textarea')
    textArea.value = textDivRef.value.innerText
    textArea.style.opacity = '0'
    document.body.appendChild(textArea)
    textArea.focus()
    textArea.select()
    try {
      isSuccess = document.execCommand('copy')
    } catch (err) {
      isSuccess = false
    }
    document.body.removeChild(textArea)
  }
  if (isSuccess) {
    toast.success('Copied to clipboard')
    showCopyBorder.value = true
    setTimeout(() => {
      showCopyBorder.value = false
    }, 2000)
  } else {
    toast.error('Failed to copy to clipboard')
  }
}
</script>

<template>
  <div
    class="relative my-2 break-words rounded-lg border-2 border-secondary-300 bg-secondary-100 p-3 transition-all"
    :class="{
      'ring-2 ring-primary-300': showCopyBorder
    }">
    <div ref="textDivRef" class="select-text whitespace-pre-wrap break-words">
      <slot></slot>
    </div>
    <div
      v-if="showCopyButton"
      @click="copyToClipboard"
      class="absolute right-0 top-0 m-3 flex h-10 w-10 cursor-copy items-center justify-center rounded-lg border-2 border-primary-300 bg-white transition-all hover:bg-gray-200">
      <font-awesome-icon icon="fa-solid fa-clipboard" class="text-xl" />
    </div>
  </div>
</template>

<style scoped></style>
