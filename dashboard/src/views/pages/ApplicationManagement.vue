<script setup>
import PageBar from '@/views/components/PageBar.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { useRouter } from 'vue-router'
import TableHeader from '@/views/components/Table/TableHeader.vue'
import TableMessage from '@/views/components/Table/TableMessage.vue'
import Table from '@/views/components/Table/Table.vue'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed } from 'vue'
import { toast } from 'vue-sonner'
import ApplicationListRow from '@/views/partials/ApplicationListRow.vue'
import ProjectListRow from '@/views/partials/ProjectListRow.vue'

const router = useRouter()

const deployNewApplication = () => {
  router.push('/deploy/application')
}

const installApplicationFromAppStore = () => {
  router.push({ name: 'App Store' })
}

const {
  result: applicationsResult,
  refetch: refetchApplications,
  loading: isApplicationsLoading,
  onError: onApplicationsError
} = useQuery(
  gql`
    query {
      applications(includeGroupedApplications: false) {
        id
        name
        replicas
        isSleeping
        realtimeInfo {
          InfoFound
          DesiredReplicas
          RunningReplicas
          DeploymentMode
          HealthStatus
        }
        latestDeployment {
          status
          upstreamType
          createdAt
        }
      }
    }
  `,
  null,
  {
    pollInterval: 60000
  }
)

onApplicationsError((err) => {
  toast.error(err.message)
})

const {
  result: applicationGroupsResult,
  refetch: refetchApplicationGroups,
  loading: isApplicationGroupsLoading,
  onError: onApplicationGroupsError
} = useQuery(
  gql`
    query {
      applicationGroups {
        id
        name
        logo
        applications {
          realtimeInfo {
            InfoFound
            DeploymentMode
            DesiredReplicas
            RunningReplicas
            HealthStatus
          }
        }
      }
    }
  `,
  null,
  {
    pollInterval: 60000
  }
)

onApplicationGroupsError((err) => {
  toast.error(err.message)
})

const refreshData = () => {
  refetchApplications()
  refetchApplicationGroups()
}

const applications = computed(() => applicationsResult.value?.applications ?? [])
const applicationGroups = computed(() => applicationGroupsResult.value?.applicationGroups ?? [])
</script>

<template>
  <section class="mx-auto w-full max-w-7xl">
    <!-- Deploy Apps Page bar   -->
    <PageBar>
      <template v-slot:title>Deployed Services</template>
      <template v-slot:subtitle>Manage your deployed services</template>
      <template v-slot:buttons>
        <FilledButton :click="deployNewApplication" type="primary">
          <font-awesome-icon icon="fa-solid fa-hammer" class="mr-2" />
          Deploy App
        </FilledButton>
        <FilledButton :click="installApplicationFromAppStore" type="primary">
          <font-awesome-icon icon="fa-solid fa-store" class="mr-2" />
          App Store
        </FilledButton>
        <FilledButton type="ghost" :click="refreshData">
          <font-awesome-icon
            icon="fa-solid fa-arrows-rotate"
            :class="{
              'animate-spin ': isApplicationsLoading || isApplicationGroupsLoading
            }" />&nbsp;&nbsp; Refresh List
        </FilledButton>
      </template>
    </PageBar>

    <p class="mt-6 text-sm font-medium">
      <font-awesome-icon icon="fa-solid fa-hammer" class="me-1" />
      Deployed Applications
    </p>

    <!-- Applications Table -->
    <Table class="mt-2">
      <template v-slot:header>
        <TableHeader align="left">Application Name</TableHeader>
        <TableHeader align="center">Health Status</TableHeader>
        <TableHeader align="center">Replicas</TableHeader>
        <TableHeader align="center">Deploy Status</TableHeader>
        <TableHeader align="center">Last Deployment</TableHeader>
        <TableHeader align="right">View Details</TableHeader>
      </template>
      <template v-slot:message>
        <TableMessage v-if="applications.length === 0">
          No deployed applications found.<br />
          Click on the "Deploy App" button to deploy a new application.
        </TableMessage>
        <TableMessage v-if="isApplicationsLoading && applications.length === 0">
          Loading deployed applications...
        </TableMessage>
      </template>
      <template v-slot:body>
        <ApplicationListRow v-for="application in applications" :key="application.id" :application="application" />
      </template>
    </Table>

    <!--  Deployed Projects Bar   -->
    <p class="mt-6 text-sm font-medium">
      <font-awesome-icon icon="fa-solid fa-layer-group" class="me-1" />
      Deployed Projects
    </p>

    <!-- Projects Table -->
    <Table class="mt-2">
      <template v-slot:header>
        <TableHeader align="left">Project Name</TableHeader>
        <TableHeader align="center">Health Status</TableHeader>
        <TableHeader align="center">Total Services</TableHeader>
        <TableHeader align="center">Healthy Services</TableHeader>
        <TableHeader align="center">Unhealthy Services</TableHeader>
        <TableHeader align="right">View Details</TableHeader>
      </template>
      <template v-slot:message>
        <TableMessage v-if="applicationGroups.length === 0"> No projects found.</TableMessage>
        <TableMessage v-if="isApplicationGroupsLoading && applicationGroups.length === 0">
          Loading deployed projects...
        </TableMessage>
      </template>
      <template v-slot:body>
        <ProjectListRow v-for="group in applicationGroups" :key="group.id" :project="group" />
      </template>
    </Table>
  </section>
</template>

<style scoped></style>
