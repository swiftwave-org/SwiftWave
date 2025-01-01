<script setup>
import { computed, nextTick, onMounted, ref, shallowRef, watch } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import OutlinedButton from '@/views/components/OutlinedButton.vue'
import DotLoader from '@/views/components/DotLoader.vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const storeEndpoints = shallowRef(['https://raw.githubusercontent.com/swiftwave-org/app-store/main/store.json'])
const apps = ref([])
const appsShown = ref([])
const searchText = ref('')
const isOptionsModalOpen = ref(false)
const selectedApp = ref({})
const selectedCategory = ref('')
const isLoading = ref(false)

watch(apps, () => {
  if (!apps.value.length) return
  searchApps()
})

onMounted(() => {
  fetchApps()
})

const closeModal = () => {
  isOptionsModalOpen.value = false
}

const openModal = () => {
  isOptionsModalOpen.value = true
}

function fetchApps() {
  isLoading.value = true
  // for each endpoint, fetch apps
  storeEndpoints.value.forEach((endpoint) => {
    fetch(endpoint)
      .then((response) => response.json())
      .then((data) => {
        apps.value = apps.value.concat(data)
        isLoading.value = false
      })
      .catch((error) => {
        console.log(error)
      })
  })
}

const categories = computed(() => {
  const appCategories = new Set()
  apps.value.forEach((app) => {
    appCategories.add(app.category)
  })
  return Array.from(appCategories).sort()
})

function searchApps() {
  // split search text by space and search for each word
  const searchWords = searchText.value.split(' ')
  appsShown.value = apps.value.filter((app) => {
    return (
      (selectedCategory.value === '' || app.category === selectedCategory.value) &&
      searchWords.every((word) => {
        return (
          app.title.toLowerCase().includes(word.toLowerCase()) ||
          app.description.toLowerCase().includes(word.toLowerCase())
        )
      })
    )
  })
}

const chooseApp = (app) => {
  if (!app) return
  if (app.stacks.length === 1) {
    openStackFileForInstall(app.stacks[0])
    return
  }
  selectedApp.value = app
  openModal()
}

const chooseCategory = (category) => {
  selectedCategory.value = category
  searchText.value = ''
  nextTick(() => {
    searchApps()
  })
}

const openStackFileForInstall = (stack) => {
  router.push({
    name: 'Install from App Store',
    query: {
      stack: stack.stack
    }
  })
}
</script>

<template>
  <!-- If loading   -->
  <div v-if="isLoading" class="flex h-full w-full items-center justify-center">
    <DotLoader />
  </div>
  <!-- Main -->
  <section v-else class="flex w-full flex-row items-start gap-2 overflow-hidden">
    <div class="navbar">
      <input
        class="block w-full rounded-md border-gray-300 text-sm shadow-sm focus:border-primary-500 focus:ring-primary-500"
        placeholder="Search Apps"
        v-model="searchText"
        @keydown.enter="searchApps"
        v-debounce:200ms="searchApps"
        type="text" />
      <div class="w-full select-none rounded-md px-2 py-2 text-sm font-medium text-black">Choose Category</div>
      <div
        class="nav-element"
        @click="chooseCategory('')"
        :class="{
          'nav-active': selectedCategory === ''
        }">
        All Apps
      </div>
      <div
        v-for="category in categories"
        :key="category"
        class="nav-element"
        @click="chooseCategory(category)"
        :class="{
          'nav-active': selectedCategory === category
        }">
        {{ category }}
      </div>
    </div>
    <!-- Apps List -->
    <div class="scrollbox h-full w-full overflow-auto">
      <!--    No app available -->
      <div v-if="appsShown.length === 0" class="flex h-full w-full flex-col items-center justify-center">
        <p class="text-5xl">🤔</p>
        <p class="ml-4 mt-10 text-xl font-medium">No apps found</p>
        <p class="mt-3">
          If you think the app should be here, Raise a request in
          <a href="https://github.com/swiftwave-org/app-store" target="_blank" class="font-semibold text-primary-600"
            >Swiftwave App Store</a
          >.
        </p>
      </div>
      <!--    Apps List (If available) -->
      <div v-else class="grid grid-cols-5 gap-2 pr-2">
        <!-- Component  -->
        <div
          @click="() => chooseApp(app)"
          v-for="app in appsShown"
          :key="app.id"
          class="flex h-[200px] cursor-pointer flex-col overflow-hidden rounded-xl border border-secondary-300 p-2 hover:border-primary-500 hover:shadow-sm">
          <!--    Header    -->
          <div class="flex flex-row gap-3 border-b pb-2">
            <div class="h-12 w-12 rounded-md p-1.5">
              <img :src="app.logo" class="h-full w-full" :alt="app.title" />
            </div>
            <div>
              <p class="text-base font-semibold text-gray-800">{{ app.title }}</p>
              <p class="text-sm">{{ app.category }}</p>
            </div>
          </div>
          <!--    Description Body    -->
          <div class="mt-2 h-full overflow-hidden text-ellipsis p-1 text-justify text-sm text-secondary-800">
            {{ app.description }}
          </div>
        </div>
      </div>
    </div>
  </section>
  <!-- Modal to show options -->
  <ModalDialog :close-modal="closeModal" :is-open="isOptionsModalOpen">
    <template v-slot:header>Install {{ selectedApp.title }}</template>
    <template v-slot:body>
      <p>Choose the preferred version -</p>
      <div class="mt-6 flex flex-col gap-2">
        <OutlinedButton
          :click="() => openStackFileForInstall(stack)"
          class="w-full"
          type="primary"
          v-for="stack in selectedApp.stacks"
          :key="stack.id">
          {{ stack.title }}
        </OutlinedButton>
      </div>
    </template>
  </ModalDialog>
</template>

<style scoped>
.scrollbox::-webkit-scrollbar {
  width: 6px;
}

.scrollbox::-webkit-scrollbar-track {
  @apply rounded-full bg-gray-200;
}

.scrollbox::-webkit-scrollbar-thumb {
  @apply rounded-full bg-primary-500;
}

.navbar {
  @apply flex h-min min-w-[200px] max-w-[200px] select-none flex-col flex-wrap gap-1 rounded-lg border border-secondary-300 p-1.5;
}

.nav-element {
  @apply min-w-max cursor-pointer rounded-md px-3 py-2 text-sm text-secondary-700 hover:bg-secondary-100;
}

.nav-active {
  @apply bg-secondary-100 font-medium text-black;
}
</style>
