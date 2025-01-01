<script setup>
import { useRouter } from 'vue-router'
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed } from 'vue'
import { toast } from 'vue-sonner'
import FilledButton from '@/views/components/FilledButton.vue'

const router = useRouter()
const applicationId = router.currentRoute.value.params.id

const {
  result: applicationDetailsRaw,
  loading: applicationDetailsLoading,
  refetch: refetchApplicationDetails
} = useQuery(
  gql`
    query ($id: String!) {
      application(id: $id) {
        id
        webhookToken
      }
    }
  `,
  {
    id: applicationId
  },
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

const webhookTriggerLink = computed(() => {
  if (applicationDetailsLoading.value) return 'Loading...'
  if (applicationDetailsRaw.value?.application?.webhookToken) {
    let token = applicationDetailsRaw.value?.application?.webhookToken ?? ''
    return location.origin + '/webhook/redeploy-app/' + applicationId + '/' + token
  } else {
    return 'Loading...'
  }
})

const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text)
  toast.success('Webhook link copied to clipboard !')
}

// Regenerate Webhook Token
const {
  mutate: regenerateWebhookToken,
  loading: regenerateWebhookTokenLoading,
  onError: regenerateWebhookTokenError,
  onDone: regenerateWebhookTokenDone
} = useMutation(
  gql`
    mutation ($id: String!) {
      regenerateWebhookToken(id: $id)
    }
  `,
  {
    fetchPolicy: 'no-cache',
    nextFetchPolicy: 'no-cache'
  }
)

regenerateWebhookTokenError((error) => {
  toast.error(error.message)
})

regenerateWebhookTokenDone((result) => {
  if (result.data.regenerateWebhookToken) {
    toast.success('Webhook token regenerated successfully !')
    refetchApplicationDetails()
  } else {
    toast.error('Something went wrong !')
  }
})

const regenerateWebhookTokenWithConfirmation = () => {
  if (
    confirm(
      'Are you sure you want to regenerate the webhook token ?\n\nThis will invalidate the previous token and you will have to update the webhook link in your git/docker repository.'
    )
  ) {
    regenerateWebhookToken({
      id: applicationId
    })
  }
}
</script>

<template>
  <!--  NOTE -->
  <div class="mb-8 rounded-md border-l-4 border-yellow-500 bg-yellow-100 p-3 text-yellow-700" role="alert">
    <p class="font-bold">NOTE:</p>
    <p>Webhook CI is only available for applications deployed using git/docker</p>
  </div>

  <p class="inline-flex items-center gap-2 text-lg font-medium">Webhook Based CI</p>
  <p class="text-sm text-secondary-700">
    You can configure your git/docker repository to trigger a new deployment on every push.
  </p>

  <!--  Link with a copy button -->
  <div class="mt-6">
    <div class="relative flex flex-row items-center gap-2">
      <input :value="webhookTriggerLink" class="w-full rounded-md border border-gray-300 p-2" readonly type="text" />
      <button
        class="absolute bottom-1 right-1 top-1 rounded-md bg-secondary-200 px-3 text-sm font-bold hover:bg-secondary-300"
        @click="copyToClipboard(webhookTriggerLink)">
        Copy
        <font-awesome-icon icon="fa-solid fa-copy" />
      </button>
    </div>
    <p class="mt-2 text-sm text-secondary-700">
      Copy the above link and paste it in your git/docker repository's webhook configuration.
    </p>
  </div>

  <!-- Regenerate Webhook tolen -->
  <div class="mt-6 flex w-full flex-row items-center justify-between rounded-md">
    <div>
      <p class="inline-flex items-center gap-2 text-lg font-medium">Regenerate Webhook Token</p>
      <p class="text-sm text-secondary-700">Regenerate the webhook token if you think it is compromised.</p>
    </div>
    <FilledButton
      type="primary"
      :loading="regenerateWebhookTokenLoading"
      @click="regenerateWebhookTokenWithConfirmation">
      Regenerate Token
    </FilledButton>
  </div>
</template>

<style scoped></style>
