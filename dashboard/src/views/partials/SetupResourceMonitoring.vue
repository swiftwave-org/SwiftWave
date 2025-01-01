<script setup>
// Modal related methods
import { toast } from 'vue-sonner'
import { ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import Step from '@/views/components/Step.vue'
import Code from '@/views/components/Code.vue'
import { getHttpBaseUrl } from '@/vendor/utils.js'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import FilledButton from '@/views/components/FilledButton.vue'

const props = defineProps({
  serverId: {
    type: Number,
    required: true
  },
  openWebConsole: {
    type: Function,
    required: false,
    default: () => {}
  }
})

const isModalOpen = ref(false)
const analyticsToken = ref('')

const openModal = () => {
  isModalOpen.value = true
  fetchAnalyticsServiceToken()
}
const closeModal = () => {
  isModalOpen.value = false
}

defineExpose({
  openModal,
  closeModal
})

// Analytics token related
const managementAddress = getHttpBaseUrl() + '/service/analytics'
const {
  mutate: fetchAnalyticsServiceTokenRaw,
  loading: isFetchingAnalyticsServiceToken,
  onDone: onFetchAnalyticsServiceTokenDone,
  onError: onFetchAnalyticsServiceTokenError
} = useMutation(gql`
  mutation ($serverId: Uint!, $rotate: Boolean!) {
    fetchAnalyticsServiceToken(id: $serverId, rotate: $rotate)
  }
`)

onFetchAnalyticsServiceTokenError((error) => {
  toast.error(error.message)
})

onFetchAnalyticsServiceTokenDone((val) => {
  analyticsToken.value = val?.data?.fetchAnalyticsServiceToken ?? ''
})

const fetchAnalyticsServiceToken = async (rotate = false) => {
  fetchAnalyticsServiceTokenRaw({
    serverId: props.serverId,
    rotate
  })
}
</script>

<template>
  <teleport to="body">
    <ModalDialog
      :close-modal="closeModal"
      :is-open="isModalOpen"
      :key="serverId + '_setup_resource_monitoring_modal'"
      width="2xl">
      <template v-slot:header>Setup Resource Monitoring</template>
      <template v-slot:body>
        <div class="mt-6">
          <Step
            title="Install swiftwave-stats-ninja package"
            sub-title="This package is responsible for collecting server and running container metrics"
            prefix-text="1"
            type="primary"
            show-body>
            <div>
              <p class="font-medium">For Debian/Ubuntu users,</p>
              <Code
                >sudo mkdir -p /etc/apt/keyrings<br />
                sudo mkdir -p /root/.gnupg<br />
                sudo gpg --no-default-keyring --keyring /etc/apt/keyrings/swiftwave.gpg --keyserver keyserver.ubuntu.com
                --recv-keys DD510C86CD3F6764<br />
                echo "deb [signed-by=/etc/apt/keyrings/swiftwave.gpg] http://deb.repo.swiftwave.org/
                swiftwave-stats-ninja stable" | sudo tee /etc/apt/sources.list.d/swiftwave-stats-ninja.list<br />
                sudo apt update<br />
                sudo apt install swiftwave-stats-ninja=2.0.1
              </Code>
              <div class="my-4"></div>
              <p class="font-medium">For Fedora/CentOS/Almalinux/RockyLinux users,</p>
              <Code
                >sudo dnf config-manager --add-repo http://rpm.repo.swiftwave.org/swiftwave.repo<br />
                sudo dnf install -y swiftwave-stats-ninja-2.0.1-1
              </Code>
            </div>
          </Step>
          <Step
            title="Enable swiftwave-stats-ninja"
            sub-title="This will start the swiftwave-stats-ninja service and enable it to start on boot"
            prefix-text="2"
            type="primary"
            show-body>
            <p v-if="isFetchingAnalyticsServiceToken" class="italic">Fetching stats-ninja authentication token ...</p>
            <Code v-else>sudo swiftwave-stats-ninja enable {{ managementAddress }} {{ analyticsToken }} </Code>
          </Step>
        </div>
      </template>
      <template v-slot:footer>
        <div class="flex w-full items-center justify-end gap-4 rounded-lg py-2">
          You can perform all the actions with
          <FilledButton type="secondary" :click="openWebConsole">
            <font-awesome-icon icon="fa-solid fa-terminal" />&nbsp;&nbsp;&nbsp;Web Console
          </FilledButton>
        </div>
      </template>
    </ModalDialog>
  </teleport>
</template>

<style scoped></style>
