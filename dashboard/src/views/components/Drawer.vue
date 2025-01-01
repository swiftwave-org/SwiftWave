<script setup>
import FilledButton from '@/views/components/FilledButton.vue'

defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  closeDrawer: {
    type: Function,
    default: () => {}
  },
  widthSize: {
    type: Number,
    default: 2,
    validator: (value) => {
      return value >= 1 && value <= 12
    }
  }
})
</script>

<template>
  <Transition>
    <div
      v-if="isOpen"
      :onclick="closeDrawer"
      class="absolute inset-0 flex h-screen max-h-screen w-screen overflow-hidden">
      <!--  Overlay background    -->
      <div class="background fixed inset-0 bg-black/25 transition-all ease-in-out"></div>
      <!-- Sidebar -->
      <div
        :onclick="(e) => e.stopPropagation()"
        class="sidebar absolute right-0 transform overflow-y-hidden bg-white p-4 transition-all ease-in-out"
        :class="{
          'w-1/12': widthSize === 1,
          'w-2/12': widthSize === 2,
          'w-3/12': widthSize === 3,
          'w-4/12': widthSize === 4,
          'w-5/12': widthSize === 5,
          'w-6/12': widthSize === 6,
          'w-7/12': widthSize === 7,
          'w-8/12': widthSize === 8,
          'w-9/12': widthSize === 9,
          'w-10/12': widthSize === 10,
          'w-11/12': widthSize === 11,
          'w-full': widthSize === 12
        }">
        <div class="flex w-full flex-col space-y-4 md:flex-row md:items-center md:justify-between md:space-y-0">
          <div>
            <h2 class="text-base font-semibold">
              <slot name="title"></slot>
            </h2>
            <p class="mt-1 text-sm text-gray-700">
              <slot name="subtitle"></slot>
            </p>
          </div>
          <FilledButton class="p-2" :click="closeDrawer">
            <font-awesome-icon icon="fa-solid fa-xmark" class="text-lg" />
          </FilledButton>
        </div>
        <div class="my-4 border border-gray-500"></div>
        <div class="body h-[90vh] overflow-y-auto">
          <slot name="body"></slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: all 0.3s;
}

.v-enter-from,
.v-leave-to {
  .background {
    opacity: 0;
  }

  .sidebar {
    opacity: 0;
    transform: translateX(100%);
  }
}

.body::-webkit-scrollbar {
  width: 12px;
}

.body::-webkit-scrollbar-thumb {
  @apply rounded-full shadow-[inset_0_0_10px_10px] shadow-primary-500;
  border: solid 3px transparent;
}
</style>
