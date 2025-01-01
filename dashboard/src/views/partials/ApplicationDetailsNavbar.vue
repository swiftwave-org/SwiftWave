<script setup></script>

<template>
  <div class="navbar">
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Deployments',
        params: { id: $route.params.id }
      }">
      Deployments
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Resource Stats',
        params: { id: $route.params.id }
      }">
      Analytics
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Runtime Logs',
        params: { id: $route.params.id }
      }">
      Runtime Log
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Ingress Rules',
        params: { id: $route.params.id }
      }">
      Ingress Rules
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Persistent Volumes',
        params: { id: $route.params.id }
      }">
      Persistent Volume
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Environment Variables',
        params: { id: $route.params.id }
      }">
      Environment Variables
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Config Mounts',
        params: { id: $route.params.id }
      }">
      Static App Config
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Update Source',
        params: { id: $route.params.id }
      }">
      Update App Source
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Deployment Config',
        params: { id: $route.params.id }
      }">
      Deploy Config
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Webhook CI',
        params: { id: $route.params.id }
      }">
      Configure Webhook
    </RouterLink>
    <RouterLink
      class="nav-element"
      :to="{
        name: 'Application Details Manage',
        params: { id: $route.params.id }
      }">
      Actions
    </RouterLink>
    <RouterLink
      class="nav-element !text-red-500"
      :to="{
        name: 'Application Details Danger Zone',
        params: { id: $route.params.id }
      }">
      Destroy
    </RouterLink>
  </div>
</template>

<style scoped>
.navbar {
  @apply flex h-min select-none flex-col flex-wrap gap-1 rounded-lg border border-secondary-200 p-1.5;
}

.nav-element {
  @apply min-w-max rounded-md px-3 py-2 text-sm text-secondary-700 hover:bg-secondary-100;
}

.router-link-exact-active {
  @apply bg-secondary-100 font-medium text-black;
}
</style>
