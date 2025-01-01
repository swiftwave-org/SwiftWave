<script setup>
import Table from '@/views/components/Table/Table.vue'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import PageBar from '@/views/components/PageBar.vue'
import { computed, ref } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { toast } from 'vue-sonner'
import TableMessage from '@/views/components/Table/TableMessage.vue'
import ServerRow from '@/views/partials/ServerRow.vue'
import CreateServerModal from '@/views/partials/CreateServerModal.vue'

const createServerModal = ref(null)

const {
  result: serversResult,
  loading: isServersLoading,
  refetch: refetchServers,
  onError: onServersError
} = useQuery(
  gql`
    query {
      servers {
        id
        ip
        hostname
        user
        ssh_port
        swarmMode
        swarmNodeStatus
        maintenanceMode
        scheduleDeployments
        dockerUnixSocketPath
        proxyEnabled
        proxyType
        status
      }
    }
  `,
  null,
  {
    pollInterval: 10000
  }
)

onServersError((err) => {
  toast.error(err.message)
})

const servers = computed(() => serversResult.value?.servers ?? [])
const openCreateServerModal = () => {
  if (createServerModal.value) createServerModal.value.openModal()
}
</script>

<template>
  <!-- Modal to create server  -->
  <CreateServerModal :callback-on-create="refetchServers" ref="createServerModal" />
  <section class="mx-auto w-full max-w-7xl">
    <!-- Top Page bar   -->
    <PageBar>
      <template v-slot:title>Registered Servers</template>
      <template v-slot:subtitle>Take control of your servers</template>
      <template v-slot:buttons>
        <FilledButton type="primary" :click="openCreateServerModal">
          <font-awesome-icon icon="fa-solid fa-plus" />
          &nbsp;&nbsp; Add Server
        </FilledButton>
        <FilledButton type="ghost" :click="refetchServers">
          <font-awesome-icon
            icon="fa-solid fa-arrows-rotate"
            :class="{
              'animate-spin ': isServersLoading
            }" />&nbsp;&nbsp; Refresh List
        </FilledButton>
      </template>
    </PageBar>

    <!-- Table -->
    <Table class="mt-8">
      <template v-slot:header>
        <TableHeader align="left">Server</TableHeader>
        <TableHeader align="center">SSH</TableHeader>
        <TableHeader align="center">Node</TableHeader>
        <TableHeader align="center">Status</TableHeader>
        <TableHeader align="center">Maintenance</TableHeader>
        <TableHeader align="center">Swarm</TableHeader>
        <TableHeader align="center">Deployment</TableHeader>
        <TableHeader align="center">Proxy</TableHeader>
        <TableHeader align="center">Analytics</TableHeader>
        <TableHeader align="center">Logs</TableHeader>
        <TableHeader align="right">Actions</TableHeader>
      </template>
      <template v-slot:message>
        <TableMessage v-if="servers.length === 0">
          No servers found.<br />
          Click on the "Add Server" button to setup a new server.
        </TableMessage>
        <TableMessage v-if="isServersLoading"> Loading server list...</TableMessage>
      </template>
      <template v-slot:body>
        <ServerRow v-for="server in servers" :key="server.id" :server="server" :refetch-servers="refetchServers" />
      </template>
    </Table>
  </section>
</template>

<style scoped></style>
