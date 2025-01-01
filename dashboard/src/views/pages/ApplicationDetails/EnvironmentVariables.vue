<script setup>
import EnvironmentVariablesEditor from '@/views/partials/DeployApplication/EnvironmentVariablesEditor.vue'
import newApplicationUpdater from '@/store/applicationUpdater.js'
import { useRouter } from 'vue-router'

const router = useRouter()
const applicationUpdater = newApplicationUpdater(router.currentRoute.value.params.id)()

</script>

<template>
  <EnvironmentVariablesEditor
    :add-environment-variable="applicationUpdater.addEnvironmentVariable"
    :delete-environment-variable="applicationUpdater.deleteEnvironmentVariable"
    :environment-variables-keys="applicationUpdater.environmentVariableDetails.keys"
    :environment-variables-map="applicationUpdater.environmentVariableDetails.map"
    :on-variable-name-change="applicationUpdater.onEnvironmentVariableNameChange"
    :on-variable-value-change="applicationUpdater.onEnvironmentVariableValueChange"
  />

</template>

<style scoped>

</style>