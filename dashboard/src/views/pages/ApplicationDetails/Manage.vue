<script setup>
import { useRouter } from 'vue-router'
import FilledButton from '@/views/components/FilledButton.vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import { getHttpBaseUrl } from '@/vendor/utils.js'

const router = useRouter()

// Restart Application
const {
  mutate: restartApplication,
  loading: restartApplicationLoading,
  onError: restartApplicationError,
  onDone: restartApplicationDone
} = useMutation(
  gql`
    mutation ($id: String!) {
      restartApplication(id: $id)
    }
  `,
  {
    fetchPolicy: 'no-cache',
    variables: {
      id: router.currentRoute.value.params.id
    }
  }
)

restartApplicationDone((result) => {
  if (result.data.restartApplication) {
    toast.success('Application restarted successfully !')
  } else {
    toast.error('Something went wrong !')
  }
})

restartApplicationError((error) => {
  toast.error(error.message)
})

const restartApplicationWithConfirmation = () => {
  const confirmation = confirm('Are you sure that you want to restart this application ?')
  if (confirmation) {
    restartApplication()
  }
}

// Rebuild Application
const {
  mutate: rebuildApplication,
  loading: rebuildApplicationLoading,
  onError: rebuildApplicationError,
  onDone: rebuildApplicationDone
} = useMutation(
  gql`
    mutation ($id: String!) {
      rebuildApplication(id: $id)
    }
  `,
  {
    fetchPolicy: 'no-cache',
    variables: {
      id: router.currentRoute.value.params.id
    }
  }
)

rebuildApplicationDone((result) => {
  if (result.data.rebuildApplication) {
    toast.success('Application rebuild request sent successfully !')
  } else {
    toast.error('Something went wrong !')
  }
  router.push({
    name: 'Application Details Deployments',
    params: {
      id: router.currentRoute.value.params.id
    }
  })
})

rebuildApplicationError((error) => {
  toast.error(error.message)
})

const rebuildApplicationWithConfirmation = () => {
  const confirmation = confirm('Are you sure that you want to rebuild this application ?')
  if (confirmation) {
    rebuildApplication()
  }
}

const openWebConsole = () => {
  const height = window.innerHeight * 0.7
  const width = window.innerWidth * 0.6
  const url = `${getHttpBaseUrl()}/console?application=${router.currentRoute.value.params.id}`
  window.open(url, '', `popup,height=${height},width=${width}`)
}
</script>

<template>
  <div class="flex flex-col items-start">
    <div class="flex w-full flex-row items-center justify-between rounded-md p-2">
      <div>
        <p class="inline-flex items-center gap-2 text-lg font-medium">SSH in Application</p>
        <p class="text-sm text-secondary-700">You can access the shell of the container running this application.</p>
      </div>
      <FilledButton type="primary" @click="openWebConsole">
        <font-awesome-icon icon="fa-solid fa-terminal" class="mr-2" />
        Open Console
      </FilledButton>
    </div>
    <div class="flex w-full flex-row items-center justify-between rounded-md p-2">
      <div>
        <p class="inline-flex items-center gap-2 text-lg font-medium">Restart Application</p>
        <p class="text-sm text-secondary-700">
          This will restart latest deployment of this app. <b>No configuration will be updated</b>
        </p>
      </div>
      <FilledButton type="primary" @click="restartApplicationWithConfirmation" :loading="restartApplicationLoading">
        <font-awesome-icon icon="fa-solid fa-rotate-right" class="mr-2" />
        Click to Restart
      </FilledButton>
    </div>

    <div class="flex w-full flex-row items-center justify-between rounded-md p-2">
      <div>
        <p class="inline-flex items-center gap-2 text-lg font-medium">Redeploy Application</p>
        <p class="text-sm text-secondary-700">This will trigger a new deployment with the latest source code.</p>
      </div>
      <FilledButton type="primary" @click="rebuildApplicationWithConfirmation" :loading="rebuildApplicationLoading">
        <font-awesome-icon icon="fa-solid fa-hammer" class="mr-2" />
        Click to Redeploy
      </FilledButton>
    </div>
  </div>
</template>

<style scoped></style>
