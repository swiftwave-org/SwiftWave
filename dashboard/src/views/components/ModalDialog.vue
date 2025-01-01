<script setup>
import { Dialog, DialogDescription, DialogPanel, DialogTitle, TransitionChild, TransitionRoot } from '@headlessui/vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  closeModal: {
    type: Function,
    default: () => {}
  },
  nonCancelable: {
    type: Boolean,
    default: false
  },
  width: {
    type: String,
    default: 'md',
    validator: (value) => {
      return ['sm', 'md', 'lg', 'xl', '2xl', '4xl', '6xl'].includes(value)
    }
  }
})

const closeModalWithValidation = () => {
  if (props.nonCancelable) {
    return
  }
  props.closeModal()
}
</script>

<template>
  <TransitionRoot :show="isOpen" appear as="template">
    <Dialog as="div" class="relative z-10 select-none" @close="closeModalWithValidation">
      <TransitionChild
        as="template"
        enter="duration-300 ease-out"
        enter-from="opacity-0"
        enter-to="opacity-100"
        leave="duration-200 ease-in"
        leave-from="opacity-100"
        leave-to="opacity-0">
        <div class="fixed inset-0 bg-black/25" />
      </TransitionChild>
      <div class="fixed inset-0 overflow-y-auto">
        <div class="flex min-h-full items-center justify-center p-4 text-center">
          <TransitionChild
            as="template"
            enter="duration-300 ease-out"
            enter-from="opacity-0 scale-95"
            enter-to="opacity-100 scale-100"
            leave="duration-200 ease-in"
            leave-from="opacity-100 scale-100"
            leave-to="opacity-0 scale-95">
            <DialogPanel
              :class="{
                'max-w-sm': width === 'sm',
                'max-w-md': width === 'md',
                'max-w-lg': width === 'lg',
                'max-w-xl': width === 'xl',
                'max-w-2xl': width === '2xl',
                'max-w-4xl': width === '4xl',
                'max-w-6xl': width === '6xl'
              }"
              class="w-full transform overflow-hidden rounded-2xl bg-white p-4 text-left align-middle shadow-xl transition-all">
              <DialogTitle as="h3" class="text-lg font-semibold leading-6 text-gray-900">
                <slot name="header"></slot>
                <!-- Close button -->
                <button
                  v-show="!nonCancelable"
                  class="absolute right-4 top-4 rounded-md border-2 p-1 text-gray-400 transition-shadow duration-200 hover:text-gray-500 hover:ring-2 hover:ring-gray-400"
                  type="button"
                  @click="closeModalWithValidation">
                  <span class="sr-only">Close</span>
                  <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    <path
                      d="M6 18L18 6M6 6l12 12"
                      stroke="currentColor"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2" />
                  </svg>
                </button>
              </DialogTitle>

              <DialogDescription class="mt-2">
                <slot name="body"></slot>
              </DialogDescription>

              <div class="mt-4 flex flex-row justify-end gap-2">
                <slot name="footer"></slot>
              </div>
            </DialogPanel>
          </TransitionChild>
        </div>
      </div>
    </Dialog>
  </TransitionRoot>
</template>
