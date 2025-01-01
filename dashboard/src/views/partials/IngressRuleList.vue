<script setup>
import TableMessage from '@/views/components/Table/TableMessage.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import Table from '@/views/components/Table/Table.vue'
import { toast } from 'vue-sonner'
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, ref, watch } from 'vue'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import IngressRuleRow from '@/views/partials/IngressRuleRow.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import router from '@/router/index.js'

const props = defineProps({
  applicationId: {
    type: String,
    required: false,
    default: ''
  }
})

// Queries
const fetchAllIngressRulesQuery = gql`
  query {
    ingressRules {
      id
      status
      protocol
      domain {
        name
      }
      httpsRedirect
      port
      targetType
      externalService
      application {
        name
      }
      authenticationType
      basicAuthAccessControlListName
      targetPort
    }
  }
`
const applicationSpecificIngressRulesQuery = gql`
  query ($id: String!) {
    application(id: $id) {
      ingressRules {
        id
        status
        protocol
        domain {
          name
        }
        httpsRedirect
        port
        application {
          name
        }
        authenticationType
        basicAuthAccessControlListName
        targetType
        targetPort
      }
    }
  }
`
// Fetch ingress rules
const {
  result: ingressRulesRaw,
  refetch: refetchIngressRules,
  loading: isIngressRulesLoading,
  onError: onIngressRulesError
} = useQuery(
  props.applicationId ? applicationSpecificIngressRulesQuery : fetchAllIngressRulesQuery,
  {
    id: props.applicationId
  },
  {
    pollInterval: 10000
  }
)

const ingressRules = computed(() =>
  props.applicationId
    ? ingressRulesRaw.value?.application?.ingressRules ?? []
    : ingressRulesRaw.value?.ingressRules ?? []
)

onIngressRulesError((err) => {
  toast.error(err.message)
})

// Delete ingress rule
const {
  mutate: deleteIngressRule,
  onDone: onIngressDeleteSuccess,
  onError: onIngressRuleDeleteFail
} = useMutation(
  gql`
    mutation ($id: Uint!) {
      deleteIngressRule(id: $id)
    }
  `,
  {
    variables: {
      id: 0
    }
  }
)

const deleteIngressRulesWithConfirmation = (ingress_rule) => {
  if (confirm('Are you sure you want to delete this ingress rule ?')) {
    deleteIngressRule({
      id: ingress_rule.id
    })
  }
}

onIngressDeleteSuccess(() => {
  toast.success('Ingress Rule will be deleted shortly\nThis can take upto 5 minutes to reflect in the system')
  refetchIngressRules()
})

onIngressRuleDeleteFail((err) => {
  toast.error(err.message)
})

// Enable/Disable HTTPS Redirect
const {
  mutate: enableHttpsRedirectRaw,
  onDone: onEnableHttpsRedirectSuccess,
  onError: onEnableHttpsRedirectFail
} = useMutation(gql`
  mutation ($id: Uint!) {
    enableHttpsRedirectIngressRule(id: $id)
  }
`)

const enableHttpsRedirect = (ingress_rule) => {
  enableHttpsRedirectRaw({
    id: ingress_rule.id
  })
}

onEnableHttpsRedirectSuccess((res) => {
  if (res.data.enableHttpsRedirectIngressRule) {
    toast.success('Requested to enable HTTPS redirect. Refresh after few seconds')
  } else {
    toast.error('Failed to enable HTTPS redirect')
  }
})

onEnableHttpsRedirectFail((err) => {
  toast.error(err.message)
})

const {
  mutate: disableHttpsRedirectRaw,
  onDone: onDisableHttpsRedirectSuccess,
  onError: onDisableHttpsRedirectFail
} = useMutation(gql`
  mutation ($id: Uint!) {
    disableHttpsRedirectIngressRule(id: $id)
  }
`)

const disableHttpsRedirect = (ingress_rule) => {
  disableHttpsRedirectRaw({
    id: ingress_rule.id
  })
}

onDisableHttpsRedirectSuccess((res) => {
  if (res.data.disableHttpsRedirectIngressRule) {
    toast.success('Requested to disable HTTPS redirect. Refresh after few seconds')
  } else {
    toast.error('Failed to disable HTTPS redirect')
  }
})

onDisableHttpsRedirectFail((err) => {
  toast.error(err.message)
})

// Recreate ingress rule
const {
  mutate: recreateIngressRule,
  onDone: onRecreateIngressRuleSuccess,
  onError: onRecreateIngressRuleFail
} = useMutation(gql`
  mutation ($id: Uint!) {
    recreateIngressRule(id: $id)
  }
`)

const recreateIngressRuleWithConfirmation = (ingress_rule) => {
  if (confirm('Are you sure you want to recreate this ingress rule ?')) {
    recreateIngressRule({
      id: ingress_rule.id
    })
  }
}

onRecreateIngressRuleSuccess(() => {
  toast.success('Ingress Rule will be recreated shortly')
  refetchIngressRules()
})

onRecreateIngressRuleFail((err) => {
  toast.error(err.message)
})

// Setup Authentication
const { result: appBasicAuthAccessControlListsRaw, onError: onAppBasicAuthAccessControlListsError } = useQuery(
  gql`
    query {
      appBasicAuthAccessControlLists {
        id
        name
      }
    }
  `,
  null,
  {
    pollInterval: 30000
  }
)

const appBasicAuthAccessControlLists = computed(
  () => appBasicAuthAccessControlListsRaw.value?.appBasicAuthAccessControlLists ?? []
)

onAppBasicAuthAccessControlListsError((err) => {
  toast.error(err.message)
})

const isSetupAuthenticationModalOpen = ref(false)
const selectedIngressRuleForSetupAuthentication = ref(null)
const selectedAuthenticationType = ref('basic')

const selectedACLIdForSetupAuthentication = ref(0)

const openSetupAuthenticationModal = (ingress_rule) => {
  selectedIngressRuleForSetupAuthentication.value = ingress_rule
  isSetupAuthenticationModalOpen.value = true
}

const closeSetupAuthenticationModal = () => {
  isSetupAuthenticationModalOpen.value = false
  selectedIngressRuleForSetupAuthentication.value = null
}

watch(isSetupAuthenticationModalOpen, (isOpening) => {
  if (!isOpening) {
    selectedIngressRuleForSetupAuthentication.value = null
  }
  selectedACLIdForSetupAuthentication.value = 0
})

const isSetupAuthenticationButtonEnabled = computed(() => {
  if (!selectedAuthenticationType.value) return false
  if (selectedAuthenticationType.value === 'basic') {
    if (!selectedACLIdForSetupAuthentication.value) return false
    else return true
  }
  return false
})

const {
  mutate: setupAuthenticationRaw,
  loading: isSetupAuthenticationLoading,
  onError: onSetupAuthenticationError,
  onDone: onSetupAuthenticationDone
} = useMutation(gql`
  mutation ($id: Uint!, $appBasicAuthAccessControlListId: Uint!) {
    protectIngressRuleUsingBasicAuth(id: $id, appBasicAuthAccessControlListId: $appBasicAuthAccessControlListId)
  }
`)

const setupAuthentication = () => {
  if (
    !confirm(
      `This operation can take 5~6 seconds to apply.\nDon't leave this page until the request is completed.\n\nAre you sure you want to continue?`
    )
  ) {
    return
  }
  setupAuthenticationRaw({
    id: selectedIngressRuleForSetupAuthentication.value.id,
    appBasicAuthAccessControlListId: selectedACLIdForSetupAuthentication.value
  })
}

onSetupAuthenticationError((err) => {
  toast.error(err.message)
})

onSetupAuthenticationDone(() => {
  toast.success('Ingress Rule is now protected')
  refetchIngressRules()
  closeSetupAuthenticationModal()
})

const openCreateACLPage = () => {
  window.open(router.resolve({ name: 'Application Auth Basic ACL' }).href, '_blank')
}

// disable authentication
const isDisableAuthenticationModalOpen = ref(false)
const selectedIngressRuleForDisableAuthentication = ref(null)

const openDisableAuthenticationModal = (ingress_rule) => {
  selectedIngressRuleForDisableAuthentication.value = ingress_rule
  isDisableAuthenticationModalOpen.value = true
}

const closeDisableAuthenticationModal = () => {
  isDisableAuthenticationModalOpen.value = false
  selectedIngressRuleForDisableAuthentication.value = null
}

watch(isDisableAuthenticationModalOpen, (isOpening) => {
  if (!isOpening) {
    selectedIngressRuleForDisableAuthentication.value = null
  }
})

const {
  mutate: disableAuthenticationRaw,
  loading: isDisableAuthenticationLoading,
  onError: onDisableAuthenticationError,
  onDone: onDisableAuthenticationDone
} = useMutation(gql`
  mutation ($id: Uint!) {
    disableIngressRuleProtection(id: $id)
  }
`)

const disableAuthentication = () => {
  if (
    !confirm(
      "This operation can take 5~6 seconds to apply.\nDon't leave this page until the request is completed.\n\nAre you sure you want to continue?"
    )
  ) {
    return
  }
  disableAuthenticationRaw({
    id: selectedIngressRuleForDisableAuthentication.value.id
  })
}

onDisableAuthenticationError((err) => {
  toast.error(err.message)
})

onDisableAuthenticationDone((res) => {
  if (res.data.disableIngressRuleProtection) {
    toast.success('Requested to disable authentication. Refresh after few seconds')
  } else {
    toast.error('Failed to disable authentication')
  }
  refetchIngressRules()
  closeDisableAuthenticationModal()
})

defineExpose({
  refetchIngressRules,
  isIngressRulesLoading
})
</script>

<template>
  <div class="w-full">
    <!-- Table -->
    <Table>
      <template v-slot:header>
        <TableHeader align="left">ID</TableHeader>
        <TableHeader align="center">Status</TableHeader>
        <TableHeader align="center">Ingress</TableHeader>
        <TableHeader align="center">
          <font-awesome-icon icon="fa-solid fa-arrow-right" />
        </TableHeader>
        <TableHeader align="center">Target</TableHeader>
        <TableHeader align="center">Authentication</TableHeader>
        <TableHeader align="center">HTTPS Redirect</TableHeader>
        <TableHeader align="right">Actions</TableHeader>
      </template>
      <template v-slot:message>
        <TableMessage v-if="ingressRules.length === 0">
          No Ingress Rules found.<br />
          Click on the "Add New" button to create a new Ingress rule.
        </TableMessage>
      </template>
      <template v-slot:body>
        <IngressRuleRow
          v-for="ingressRule in ingressRules"
          :key="ingressRule.id"
          :ingress-rule="ingressRule"
          :disable-https-redirect="() => disableHttpsRedirect(ingressRule)"
          :enable-https-redirect="() => enableHttpsRedirect(ingressRule)"
          :delete-ingress-rule="() => deleteIngressRulesWithConfirmation(ingressRule)"
          :recreate-ingress-rule="() => recreateIngressRuleWithConfirmation(ingressRule)"
          :setup-authentication="() => openSetupAuthenticationModal(ingressRule)"
          :disable-authentication="() => openDisableAuthenticationModal(ingressRule)"
          :restrict-table-width="props.applicationId.length > 0" />
      </template>
    </Table>
    <!-- Modal to protect ingress rule -->
    <ModalDialog :close-modal="closeSetupAuthenticationModal" :is-open="isSetupAuthenticationModalOpen">
      <template v-slot:header>Protect Ingress Rule</template>
      <template v-slot:body>
        Provide all the details to protect the ingress rule.
        <form @submit.prevent="">
          <!--   Auth Type     -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="name"> Authentication Type </label>
            <div class="mt-1">
              <select
                v-model="selectedAuthenticationType"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                <option value="basic">Basic Authentication</option>
              </select>
            </div>
          </div>
          <!--  User List Field   -->
          <div class="mt-4" v-if="selectedAuthenticationType === 'basic'">
            <label class="block text-sm font-medium text-gray-700" for="name"> Select User List </label>
            <div class="mt-1">
              <select
                v-model="selectedACLIdForSetupAuthentication"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm">
                <option :value="acl.id" :key="acl.id" v-for="acl in appBasicAuthAccessControlLists">
                  {{ acl.name }}
                </option>
              </select>
            </div>
            <p class="mt-2 flex items-center text-sm">
              Need to create a user list ?
              <a @click="openCreateACLPage" class="ml-1.5 cursor-pointer font-bold text-primary-600">Create ACL </a>
            </p>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton
          :click="setupAuthentication"
          :loading="isSetupAuthenticationLoading"
          :disabled="!isSetupAuthenticationButtonEnabled"
          type="primary"
          class="w-full"
          >Confirm & Protect
        </FilledButton>
      </template>
    </ModalDialog>
    <!-- Modal to disable authentication -->
    <ModalDialog :close-modal="closeDisableAuthenticationModal" :is-open="isDisableAuthenticationModalOpen">
      <template v-slot:header>Disable Authentication</template>
      <template v-slot:body> Are you sure you want to disable authentication for this ingress rule ?</template>
      <template v-slot:footer>
        <FilledButton
          :click="disableAuthentication"
          :loading="isDisableAuthenticationLoading"
          type="primary"
          class="w-full"
          >Confirm & Disable
        </FilledButton>
      </template>
    </ModalDialog>
  </div>
</template>

<style scoped></style>
