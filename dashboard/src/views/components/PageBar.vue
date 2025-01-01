<template>
  <div class="flex w-full flex-col space-y-2 md:flex-row md:items-center md:justify-between md:space-y-0">
    <div>
      <p class="text-base font-medium">
        <slot name="title"></slot>
      </p>
      <p class="text-sm text-gray-700">
        <slot name="subtitle"></slot>
      </p>
    </div>
    <div class="flex flex-row gap-2">
      <slot name="buttons"></slot>
    </div>
  </div>
</template>
