<script setup>

import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { onMounted, reactive } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'


const props = defineProps({
  isModalOpen: {
    type: Boolean,
    required: true
  },
  closeModal: {
    type: Function,
    required: true
  }
})

const passwordDetails = reactive({
  newPassword: '',
  oldPassword: ''
})
const resetPasswordDetails = () => {
  passwordDetails.newPassword = ''
  passwordDetails.oldPassword = ''
}

onMounted(() => resetPasswordDetails())

const {
  mutate: changePassword,
  loading: isChangingPassword,
  onError: onPasswordChangeFail,
  onDone: onPasswordChangeSuccess
} = useMutation(gql`
mutation ($input: PasswordUpdateInput) {
  changePassword(input: $input)
}`, {
  fetchPolicy: 'no-cache',
  variables: {
    input: passwordDetails
  }
})

onPasswordChangeSuccess((e) => {
  if (e.data.changePassword === true) {
    toast.success('Password changed successfully')
    props.closeModal()
  } else {
    toast.error('Password change failed')
  }
  resetPasswordDetails()
})

onPasswordChangeFail((error) => {
  toast.error(error.message)
  resetPasswordDetails()
})

</script>

<template>
  <Teleport to="body">
    <ModalDialog
      :close-modal="closeModal"
      :is-open="isModalOpen">
      <template v-slot:header>
        Change Password
      </template>
      <template v-slot:body>
        Enter details to change password.
        <form @submit.prevent="changePassword">
          <!-- Username Field -->
          <div class="mt-4">
            <label
              class="block text-sm font-medium text-gray-700"
              for="oldPassword">
              Old Password
            </label>
            <div class="mt-1">
              <input
                id="oldPassword"
                v-model="passwordDetails.oldPassword"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="Old Password"
                type="password" />
            </div>
          </div>
          <!-- Password Field -->
          <div class="mt-4">
            <label
              class="block text-sm font-medium text-gray-700"
              for="newPassword">
              New Password
            </label>
            <div class="mt-1">
              <input
                id="newPassword"
                v-model="passwordDetails.newPassword"
                autocomplete="off"
                class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 sm:text-sm"
                placeholder="New Password"
                type="password" />
            </div>
          </div>
        </form>
      </template>
      <template v-slot:footer>
        <FilledButton
          :click="changePassword"
          :loading="isChangingPassword"
          type="primary"
        >Change Password
        </FilledButton>
      </template>
    </ModalDialog>
  </Teleport>
</template>

<style scoped>

</style>