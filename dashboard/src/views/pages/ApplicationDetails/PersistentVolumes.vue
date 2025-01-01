<script setup>

import PersistentVolumeBindingEditor from '@/views/partials/DeployApplication/PersistentVolumeBindingEditor.vue'
import { useRouter } from 'vue-router'
import newApplicationUpdater from '@/store/applicationUpdater.js'

const router = useRouter()
const applicationUpdater = newApplicationUpdater(router.currentRoute.value.params.id)()
</script>

<template>
  <PersistentVolumeBindingEditor
    :add-persistent-volume-binding="applicationUpdater.addPersistentVolumeBinding"
    :delete-persistent-volume-binding="applicationUpdater.deletePersistentVolumeBinding"
    :on-mounting-path-change="applicationUpdater.onPersistentVolumeMountingPathChange"
    :on-persistent-volume-change="applicationUpdater.onPersistentVolumeChange"
    :persistent-volume-binding-keys="applicationUpdater.persistentVolumeBindingsDetails.keys"
    :persistent-volume-bindings-map="applicationUpdater.persistentVolumeBindingsDetails.map"
  />
</template>

<style scoped>

</style>