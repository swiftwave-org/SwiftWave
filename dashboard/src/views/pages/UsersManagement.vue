<script setup>
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, reactive, ref } from 'vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import PageBar from '@/views/components/PageBar.vue'
import { toast } from 'vue-sonner'
import VueQrcode from 'vue-qrcode'

import Table from '@/views/components/Table/Table.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import UserListRow from '@/views/partials/UserListRow.vue'
import TableMessage from '@/views/components/Table/TableMessage.vue'
import { preventSpaceInput } from '@/vendor/utils.js'
import Code from '@/views/components/Code.vue'
import Divider from '@/views/components/Divider.vue'

const isModalOpen = ref(false)
const openModal = () => {
  isModalOpen.value = true
}
const closeModal = () => {
  isModalOpen.value = false
}

// New user form state
const newUser = reactive({
  username: '',
  password: ''
})

const {
  mutate: createUser,
  loading: isUserCreating,
  onDone: onUserCreateSuccess,
  onError: onUserCreateFail
} = useMutation(
  gql`
    mutation ($input: UserInput!) {
      createUser(input: $input) {
        id
        username
        totpEnabled
      }
    }
  `,
  {
    variables: {
      input: newUser
    }
  }
)

onUserCreateSuccess(() => {
  closeModal()
  refetchUserList()
  newUser.username = ''
  newUser.password = ''
  toast.success('User created successfully')
})

onUserCreateFail((err) => {
  toast.error(err.message)
})

// Delete user mutation
const {
  mutate: deleteUser,
  onDone: onUserDeleteSuccess,
  onError: onUserDeleteFail
} = useMutation(gql`
  mutation ($id: Uint!) {
    deleteUser(id: $id)
  }
`)

const deleteUserWithConfirmation = (user) => {
  if (confirm(`Are you sure you want to delete user ${user.username}?`)) {
    deleteUser({ id: user.id })
  }
}

onUserDeleteSuccess(() => {
  refetchUserList()
  toast.success('User deleted successfully')
})

onUserDeleteFail((err) => {
  toast.error(err.message)
})

// User list query
const {
  result: userListResult,
  loading: isUserListLoading,
  refetch: refetchUserList,
  onError: onUserListFetchFailed
} = useQuery(
  gql`
    query {
      users {
        id
        username
        totpEnabled
      }
    }
  `,
  null,
  {
    pollInterval: 10000
  }
)
const users = computed(() => userListResult.value?.users)

onUserListFetchFailed((err) => {
  toast.error(err.message)
})

// TOTP Enable Related
const totpModalOpen = ref(false)
const enableTotpRequest = reactive({
  totpSecret: '',
  totpProvisioningUri: '',
  filledTotp: ''
})
const resetTotpRequest = () => {
  enableTotpRequest.totpSecret = ''
  enableTotpRequest.totpProvisioningUri = ''
  enableTotpRequest.filledTotp = ''
}
const closeTotpModal = () => {
  totpModalOpen.value = false
}

const {
  mutate: requestEnableTotp,
  loading: isRequestEnableTotpLoading,
  onDone: onRequestEnableTotpSuccess,
  onError: onRequestEnableTotpError
} = useMutation(gql`
  mutation {
    requestTotpEnable {
      totpSecret
      totpProvisioningUri
    }
  }
`)

onRequestEnableTotpSuccess((response) => {
  enableTotpRequest.totpSecret = response.data.requestTotpEnable.totpSecret
  enableTotpRequest.totpProvisioningUri = response.data.requestTotpEnable.totpProvisioningUri
  totpModalOpen.value = true
})

onRequestEnableTotpError((err) => {
  toast.error(err.message)
})

const requestEnableTotpWithConfirmation = () => {
  if (confirm(`Are you sure you want to enable TOTP 2FA for ?`)) {
    resetTotpRequest()
    requestEnableTotp()
  }
}

const {
  mutate: enableTotpRaw,
  loading: isEnableTotpLoading,
  onDone: onEnableTotpSuccess,
  onError: onEnableTotpError
} = useMutation(gql`
  mutation ($totp: String!) {
    enableTotp(totp: $totp)
  }
`)

const enableTotp = () => {
  enableTotpRaw({ totp: enableTotpRequest.filledTotp })
}

onEnableTotpSuccess(() => {
  closeTotpModal()
  resetTotpRequest()
  toast.success('TOTP 2FA enabled successfully')
  refetchUserList()
})

onEnableTotpError((err) => {
  toast.error(err.message)
})

// Disable TOTP Related
const {
  mutate: disableTotpRaw,
  loading: isDisableTotpLoading,
  onDone: onDisableTotpSuccess,
  onError: onDisableTotpError
} = useMutation(gql`
  mutation {
    disableTotp
  }
`)

const disableTotpWithConfirmation = () => {
  if (confirm(`Are you sure you want to disable TOTP 2FA for current user ?`)) {
    disableTotpRaw()
  }
}

onDisableTotpSuccess((response) => {
  if (response.data.disableTotp) {
    toast.success('TOTP 2FA disabled successfully')
    refetchUserList()
  } else {
    toast.error('Failed to disable TOTP 2FA')
  }
})

onDisableTotpError((err) => {
  toast.error(err.message)
})
</script>

<template>
  <!-- Modal for totp -->
  <ModalDialog :close-modal="closeTotpModal" :is-open="totpModalOpen">
    <template v-slot:header>TOTP 2FA</template>
    <template v-slot:body>
      <div class="mt-6 flex flex-col items-center">
        <p class="font-medium">Scan the QR code with Authenticator app.</p>
        <VueQrcode class="my-4" :value="enableTotpRequest.totpProvisioningUri" />
        <p class="font-medium">Or Paste the secret code in authenticator app.</p>
        <Code :show-copy-button="false">{{ enableTotpRequest.totpSecret }}</Code>
      </div>
      <Divider />
      <div class="flex w-full flex-col items-center gap-4">
        <p class="font-medium">Enter TOTP from app</p>
        <v-otp-input
          :num-inputs="6"
          input-classes="otp-input"
          :placeholder="['*', '*', '*', '*', '*', '*']"
          v-model:value="enableTotpRequest.filledTotp"
          @on-change="(v) => (enableTotpRequest.filledTotp = v)" />
      </div>
    </template>
    <template v-slot:footer>
      <FilledButton
        :click="enableTotp"
        :loading="isEnableTotpLoading"
        :disabled="enableTotpRequest.filledTotp.length !== 6"
        type="primary"
        class="w-full">
        Verify & Enable TOTP 2FA
      </FilledButton>
    </template>
  </ModalDialog>

  <section class="mx-auto w-full max-w-7xl">
    <!-- Modal for new user -->
    <ModalDialog :close-modal="closeModal" :is-open="isModalOpen">
      <template v-slot:header>Create new user</template>
      <template v-slot:body>
        Enter the username and password for the new user.
        <form @submit.prevent="createUser">
          <!-- Username Field -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="username"> Username </label>
            <div class="mt-1">
              <input
                id="username"
                v-model="newUser.username"
                @keydown="preventSpaceInput"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                name="username"
                placeholder="Username"
                type="text" />
            </div>
          </div>
          <!-- Password Field -->
          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700" for="password"> Password </label>
            <div class="mt-1">
              <input
                id="password"
                v-model="newUser.password"
                autocomplete="new-password"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                name="password"
                placeholder="Password"
                type="password" />
            </div>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton :click="createUser" :loading="isUserCreating" type="primary">Create</FilledButton>
      </template>
    </ModalDialog>

    <!-- Top Page bar   -->
    <PageBar>
      <template v-slot:title>Users</template>
      <template v-slot:subtitle>
        Registered users can access the SwiftWave dashboard, allowing them to perform all actions and access all
        features
      </template>
      <template v-slot:buttons>
        <FilledButton :click="openModal" type="primary">
          <font-awesome-icon icon="fa-solid fa-plus" class="mr-2" />
          Create User
        </FilledButton>
        <FilledButton type="ghost" :click="refetchUserList">
          <font-awesome-icon
            icon="fa-solid fa-arrows-rotate"
            :class="{
              'animate-spin ': isUserListLoading
            }" />&nbsp;&nbsp; Refresh List
        </FilledButton>
      </template>
    </PageBar>

    <!-- Tables -->
    <Table class="mt-8">
      <template v-slot:header>
        <TableHeader align="left">Username</TableHeader>
        <TableHeader align="center">Status</TableHeader>
        <TableHeader align="center">Role</TableHeader>
        <TableHeader align="center">2FA</TableHeader>
        <TableHeader align="right">Actions</TableHeader>
      </template>
      <template v-slot:message>
        <TableMessage v-if="!users"> Loading users...</TableMessage>
        <TableMessage v-else-if="users.length === 0">
          No users found.<br />
          Click on the "Create User" button to create a new user.
        </TableMessage>
      </template>
      <template v-slot:body>
        <UserListRow
          v-for="user in users"
          v-bind:key="user.id"
          :delete-user="deleteUserWithConfirmation"
          :user="user"
          :is-request-running-for-totp="isRequestEnableTotpLoading || isDisableTotpLoading"
          :enable-totp-current-user="requestEnableTotpWithConfirmation"
          :disable-totp-current-user="disableTotpWithConfirmation" />
      </template>
    </Table>
  </section>
</template>
