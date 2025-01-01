<script setup>
import { useRouter } from 'vue-router'
import newApplicationUpdater from '@/store/applicationUpdater.js'
import ConfigMountsEditor from '@/views/partials/DeployApplication/ConfigMountsEditor.vue'

const router = useRouter()
const applicationUpdater = newApplicationUpdater(router.currentRoute.value.params.id)()
</script>

<template>
  <ConfigMountsEditor
    :config-mounts-keys="applicationUpdater.configMountDetails.keys"
    :config-mounts-map="applicationUpdater.configMountDetails.map"
    :add-config-mount="applicationUpdater.addConfigMount"
    :delete-config-mount="applicationUpdater.deleteConfigMount"
    :on-config-content-change="applicationUpdater.onConfigMountContentChange" />
</template>

<style scoped></style>
