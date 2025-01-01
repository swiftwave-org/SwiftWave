<script setup>
import Badge from '@/views/components/Badge.vue'
import TableRow from '@/views/components/Table/TableRow.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import TextButton from '@/views/components/TextButton.vue'
import { toast } from 'vue-sonner'
import { ref } from 'vue'

const props = defineProps({
  domain: {
    type: Object,
    required: true
  },
  deleteDomain: {
    type: Function,
    required: true
  },
  verifyDns: {
    type: Function,
    required: true
  },
  issueSsl: {
    type: Function,
    required: true
  },
  viewSsl: {
    type: Function,
    required: true
  }
})

const verifyingDns = ref(false)

const verifyDnsPointing = async () => {
  verifyingDns.value = true

  const verifyPromise = new Promise(function(resolve, reject) {
    props.verifyDns(props.domain.name)
    .then((isVerified) => {
      if (isVerified) {
        resolve()
      } else {
        reject()
      }
    })
  })

  toast.promise(verifyPromise, {
    loading: `Verifying DNS for ${props.domain.name}`,
    success: (_) => {
      verifyingDns.value = false
      return `DNS for ${props.domain.name} is pointed correctly`
    },
    error: (_) => {
      verifyingDns.value = false
      return `DNS for ${props.domain.name} is not pointed correctly`
    }
  });
}
</script>

<template>
  <tr>
    <TableRow align="left">
      <div class="text-sm font-medium text-gray-900">{{ domain.name }}</div>
    </TableRow>
    <TableRow align="center">
      <Badge v-if="domain.sslStatus === 'none'" type="secondary">N/A</Badge>
      <Badge v-else-if="domain.sslStatus === 'pending'" type="warning">Pending</Badge>
      <Badge v-else-if="domain.sslStatus === 'issued'" type="success">Issued</Badge>
      <Badge v-else-if="domain.sslStatus === 'failed'" type="danger">Failed</Badge>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton :click="() => viewSsl(domain.id)" :disabled="domain.sslStatus !== 'issued'" slim type="secondary"
        >View SSL
      </FilledButton>
    </TableRow>
    <TableRow align="center">
      <div class="text-sm text-gray-900">
        {{ domain.sslIssuer !== '' ? domain.sslIssuer : '---' }}
      </div>
    </TableRow>
    <TableRow align="center" flex>
      <FilledButton :click="() => issueSsl(domain)" slim type="secondary">Issue SSL</FilledButton>
    </TableRow>
    <TableRow align="center">
      <Badge v-if="domain.sslAutoRenew" type="success">Enabled</Badge>
      <Badge v-else type="danger">Disabled</Badge>
    </TableRow>

    <TableRow align="center" flex>
      <FilledButton :click="verifyDnsPointing" :loading="verifyingDns" slim type="secondary">Verify DNS</FilledButton>
    </TableRow>
    <TableRow align="right">
      <TextButton :click="() => deleteDomain(domain)" type="danger">Delete</TextButton>
    </TableRow>
  </tr>
</template>

<style scoped></style>
