<script setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'
import FilledButton from '@/views/components/FilledButton.vue'
import CreateIngressRuleModal from '@/views/partials/CreateIngressRuleModal.vue'
import IngressRuleList from '@/views/partials/IngressRuleList.vue'

const router = useRouter()

// Create new ingress rule
const newIngressRuleModalRef = ref(null)
const openNewIngressRuleModal = () => {
  if (!newIngressRuleModalRef.value?.openModal) return
  newIngressRuleModalRef.value.openModal(router.currentRoute.value.params.id)
}

// Ingress Rules List
const ingressRuleListRef = ref(null)
const refetchIngressRules = () => {
  if (!ingressRuleListRef.value?.refetchIngressRules) return
  ingressRuleListRef.value.refetchIngressRules()
}
</script>

<template>
  <!-- Create Ingress Rule Modal -->
  <CreateIngressRuleModal ref="newIngressRuleModalRef" :callback-on-create="refetchIngressRules" />
  <!-- Table -->
  <IngressRuleList ref="ingressRuleListRef" :application-id="router.currentRoute.value.params.id" />
  <!-- More actions  -->
  <div class="mt-4 flex items-center justify-center gap-3 text-sm">
    Expose your application to the world
    <FilledButton slim type="primary" :click="openNewIngressRuleModal">Add Ingress Rule</FilledButton>
  </div>
</template>

<style scoped></style>
