<script setup>
import { useAuthStore } from '@/store/auth.js'
import { RouterLink, useRouter } from 'vue-router'
import Logo from '@/assets/images/logo-full-inverse-subtitle.png'
import ChangePasswordModal from '@/views/partials/ChangePasswordModal.vue'
import { computed, onMounted, ref } from 'vue'
import SideBarOption from '@/views/partials/SideBarOption.vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import ModalDialog from '@/views/components/ModalDialog.vue'

const authStore = useAuthStore()
const router = useRouter()

const isChangePasswordModalOpen = ref(false)
const swVersion = ref('')
const openChangePasswordModal = () => {
  isChangePasswordModalOpen.value = true
}
const closeChangePasswordModal = () => {
  isChangePasswordModalOpen.value = false
}
const isShowSideBar = computed(() => {
  if (!authStore.IsLoggedIn) {
    return false
  } else {
    return !['Download Persistent Volume Backup', 'Maintenance', 'Setup'].includes(router.currentRoute.value.name)
  }
})

const logoutWithConfirmation = () => {
  if (confirm('Are you sure you want to logout?')) {
    authStore.Logout()
  }
}

const fetchSWVersion = () => {
  if (authStore.IsLoggedIn) {
    authStore.fetchSWVersion().then((v) => {
      swVersion.value = v
    })
  }
}

onMounted(() => {
  fetchSWVersion()
  const intervalId = setInterval(() => {
    if (authStore.IsLoggedIn) {
      if (swVersion.value === '') {
        fetchSWVersion()
      } else {
        clearInterval(intervalId)
      }
    }
  }, 2000)
})

// Restart system
const timeCount = ref(5)

const isSystemRestartModalOpen = ref(false)
const {
  mutate: restartSystem,
  onDone: onRestartSystemDone,
  onError: onRestartSystemError
} = useMutation(gql`
  mutation {
    restartSystem
  }
`)

onRestartSystemError((error) => {
  toast.error(error.message)
})

onRestartSystemDone((val) => {
  if (val.data.restartSystem) {
    toast.success('System restart requested')
    isSystemRestartModalOpen.value = true
    startCountDown()
  } else {
    toast.error('System restart failed')
  }
})

const systemRestart = () => {
  if (confirm('Are you sure you want to restart swiftwave ?\nYour deployed applications will not face any downtime.')) {
    restartSystem()
  }
}

const startCountDown = () => {
  const interval = setInterval(() => {
    timeCount.value--
    if (timeCount.value === 0) {
      clearInterval(interval)
      isSystemRestartModalOpen.value = false
      router.push({ name: 'Maintenance', query: { redirect: router.currentRoute.value.path } })
    }
  }, 1000)
}
</script>

<template>
  <aside
    v-if="isShowSideBar"
    class="scrollbox flex h-screen flex-col overflow-y-auto border-r bg-primary-600 px-2 pb-2 pt-6">
    <div class="px-3">
      <RouterLink to="/">
        <img :src="Logo" alt="logo" class="w-full max-w-40" />
      </RouterLink>
    </div>
    <div class="mt-6 flex flex-1 flex-col justify-between">
      <nav>
        <SideBarOption :active-urls="['Deploy Application', 'Deploy Stack', 'App Store', 'Install from App Store']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-hammer" />
          </template>
          <template #title> Deploy Application</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-50 hover:text-gray-700"
                to="/deploy/app-store">
                <font-awesome-icon icon="fa-solid fa-store" />
                <span class="mx-2 text-sm font-medium">App Store</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-50 hover:text-gray-700"
                to="/deploy/application">
                <font-awesome-icon icon="fa-solid fa-hammer" />
                <span class="mx-2 text-sm font-medium">Deploy App</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-50 hover:text-gray-700"
                to="/deploy/stack">
                <font-awesome-icon icon="fa-solid fa-cubes-stacked" />
                <span class="mx-2 text-sm font-medium">Deploy Stack</span>
              </RouterLink>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['Applications', 'Persistent Volumes']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-box" />
          </template>
          <template #title> Applications & Volumes</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/applications">
                <font-awesome-icon icon="fa-solid fa-box" />
                <span class="mx-2 text-sm font-medium">Applications</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/persistent-volumes">
                <font-awesome-icon icon="fa-solid fa-hard-drive" />
                <span class="mx-2 text-sm font-medium">Persistent Volumes</span>
              </RouterLink>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['Domains', 'Redirect Rules', 'Ingress Rules']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-route" />
          </template>
          <template #title>Manage Routing</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/domains">
                <font-awesome-icon icon="fa-solid fa-link" />
                <span class="mx-2 text-sm font-medium">Domains</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/ingress-rules">
                <font-awesome-icon icon="fa-solid fa-network-wired" />
                <span class="mx-2 text-sm font-medium">Ingress Rules</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/redirect-rules">
                <font-awesome-icon icon="fa-solid fa-location-arrow" />
                <span class="mx-2 text-sm font-medium">Redirect Rules</span>
              </RouterLink>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['Git Credentials', 'Image Registry Credentials']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-vault" />
          </template>
          <template #title>Manage Credentials</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/git-credentials">
                <font-awesome-icon icon="fa-solid fa-code-branch" />
                <span class="mx-2 text-sm font-medium">Git Credentials</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/image-registry-credentials">
                <font-awesome-icon icon="fa-solid fa-cloud" />
                <span class="mx-2 text-sm font-medium">Image Reg Credentials</span>
              </RouterLink>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['Application Auth Basic ACL']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-shield-halved" />
          </template>
          <template #title>Protect Application</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/app_auth/basic_authentication">
                <font-awesome-icon icon="fa-solid fa-user-shield" />
                <span class="mx-2 text-sm font-medium">Basic Authentication</span>
              </RouterLink>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['Servers']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-server" />
          </template>
          <template #title>Manage Servers</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/servers">
                <font-awesome-icon icon="fa-solid fa-list-ul" />
                <span class="mx-2 text-sm font-medium">Server list</span>
              </RouterLink>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['System Logs']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-gear" />
          </template>
          <template #title> Manage System</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/logs">
                <font-awesome-icon icon="fa-solid fa-file-waveform" />
                <span class="mx-2 text-sm font-medium">System Logs</span>
              </RouterLink>
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/setup?update=1">
                <font-awesome-icon icon="fa-solid fa-wrench" />
                <span class="mx-2 text-sm font-medium">System Configuration</span>
              </RouterLink>
              <div
                @click="systemRestart"
                class="flex transform cursor-pointer items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700">
                <font-awesome-icon icon="fa-solid fa-power-off" />
                <span class="mx-2 text-sm font-medium">System Restart</span>
              </div>
            </div>
          </template>
        </SideBarOption>

        <SideBarOption :active-urls="['Users']">
          <template #icon>
            <font-awesome-icon icon="fa-solid fa-user-tie" />
          </template>
          <template #title> Administration</template>
          <template #content>
            <div class="space-y-2">
              <RouterLink
                class="flex transform items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                to="/users">
                <font-awesome-icon icon="fa-solid fa-users" />
                <span class="mx-2 text-sm font-medium">Manage Users</span>
              </RouterLink>
              <div
                class="flex transform cursor-pointer items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                @click="openChangePasswordModal">
                <font-awesome-icon icon="fa-solid fa-key" />
                <span class="mx-2 text-sm font-medium">Change Password</span>
              </div>
              <a
                class="flex transform cursor-pointer items-center rounded-lg px-3 py-2 text-gray-200 transition-colors duration-300 hover:bg-gray-100 hover:text-gray-700"
                @click="logoutWithConfirmation">
                <font-awesome-icon icon="fa-solid fa-right-from-bracket" />
                <span class="mx-2 text-sm font-medium">Logout</span>
              </a>
            </div>
          </template>
        </SideBarOption>
      </nav>
    </div>
    <div class="flex justify-between px-2 text-sm font-medium text-white">
      <span>Auto-logout {{ authStore.sessionRelativeTimeoutStatus }}</span>
      <span> v{{ swVersion }}</span>
    </div>
    <ChangePasswordModal :is-modal-open="isChangePasswordModalOpen" :close-modal="closeChangePasswordModal" />
    <Teleport to="body">
      <!-- Modal for restart system -->
      <ModalDialog :is-open="isSystemRestartModalOpen" non-cancelable>
        <template v-slot:header>
          <span>🔌 Restarting System</span>
        </template>
        <template v-slot:body>
          <p class="mb-2">System restart has been requested.</p>
          <p>
            Redirecting to Maintenance Page in <b>{{ timeCount }}</b> seconds
          </p>
        </template>
      </ModalDialog>
    </Teleport>
  </aside>
</template>

<style scoped>
.router-link-exact-active {
  @apply bg-gray-100 text-gray-700;
}

.scrollbox::-webkit-scrollbar {
  width: 12px;
}

.scrollbox::-webkit-scrollbar-thumb {
  @apply rounded-full shadow-[inset_0_0_10px_10px_white];
  border: solid 3px transparent;
}
</style>
