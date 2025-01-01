<script setup>
import TableRow from '@/views/components/Table/TableRow.vue'
import Badge from '@/views/components/Badge.vue'
import TextButton from '@/views/components/TextButton.vue'
import { useAuthStore } from '@/store/auth.js'
import FilledButton from '@/views/components/FilledButton.vue'

const currentUsername = useAuthStore().currentUsername
defineProps({
  user: {
    type: Object,
    required: true
  },
  deleteUser: {
    type: Function,
    required: true
  },
  enableTotpCurrentUser: {
    type: Function,
    required: false,
    default: () => {}
  },
  disableTotpCurrentUser: {
    type: Function,
    required: false,
    default: () => {}
  },
  isRequestRunningForTotp: {
    required: false,
    default: false
  }
})
</script>

<template>
  <tr>
    <TableRow align="left">
      <div class="text-sm font-medium text-gray-900">
        {{ user.username }}
      </div>
    </TableRow>
    <TableRow align="center">
      <Badge type="success">Active</Badge>
    </TableRow>
    <TableRow align="center">
      <span class="text-sm text-gray-700"> Administrator </span>
    </TableRow>
    <TableRow align="center" v-if="currentUsername === user.username" flex>
      <FilledButton
        type="primary"
        slim
        v-if="!user.totpEnabled"
        :loading="isRequestRunningForTotp"
        :click="enableTotpCurrentUser"
        >Enable TOTP
      </FilledButton>
      <FilledButton type="danger" slim v-else :loading="isRequestRunningForTotp" :click="disableTotpCurrentUser"
        >Disable TOTP
      </FilledButton>
    </TableRow>
    <TableRow align="center" v-else>
      <Badge type="success" v-if="user.totpEnabled">TOTP Enabled</Badge>
      <Badge type="danger" v-else>TOTP Disabled</Badge>
    </TableRow>
    <TableRow align="right">
      <TextButton :click="() => deleteUser(user)" type="danger" :disabled="currentUsername === user.username">
        Delete
      </TextButton>
    </TableRow>
  </tr>
</template>

<style scoped></style>
