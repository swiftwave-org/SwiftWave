<script setup>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import TableRow from '@/views/components/Table/TableRow.vue'
import Badge from '@/views/components/Badge.vue'
import FilledButton from '@/views/components/FilledButton.vue'
import { ref } from 'vue'

defineProps({
  ingressRule: {
    type: Object,
    required: true
  },
  deleteIngressRule: {
    type: Function,
    required: true
  },
  enableHttpsRedirect: {
    type: Function,
    required: true
  },
  disableHttpsRedirect: {
    type: Function,
    required: true
  },
  recreateIngressRule: {
    type: Function,
    required: true
  },
  setupAuthentication: {
    type: Function,
    required: true
  },
  disableAuthentication: {
    type: Function,
    required: true
  },
  restrictTableWidth: {
    type: Boolean,
    default: false
  }
})

const actionsBtnRef = ref(null)
const actionsMenuRef = ref(null)

const onClickActions = () => {
  if (actionsBtnRef.value === null || actionsBtnRef.value.$el === null) {
    return
  }
  if (actionsMenuRef.value === null) {
    return
  }
  if (actionsMenuRef.value.style.display === 'block') {
    actionsMenuRef.value.style.display = 'none'
    return
  }
  const rect = actionsBtnRef.value.$el.getBoundingClientRect()
  const menuEl = actionsMenuRef.value
  menuEl.style.display = 'block'
  menuEl.style.minWidth = `${rect.width}px`
  menuEl.style.top = `${rect.top + rect.height + 8}px`
  menuEl.style.right = `${window.innerWidth - rect.left - rect.width}px`
}

const closeMenu = () => {
  if (!actionsMenuRef.value) {
    return
  }
  actionsMenuRef.value.style.display = 'none'
}

// on screen resize close the menu
window.addEventListener('resize', closeMenu)
// on click outside close the menu
window.addEventListener('click', (e) => {
  if (!actionsMenuRef.value || !actionsBtnRef.value.$el) {
    return
  }
  if (!actionsBtnRef.value.$el.contains(e.target)) {
    closeMenu()
  }
})
</script>

<template>
  <tr :key="ingressRule.id">
    <TableRow align="left">
      <div class="text-sm font-medium text-gray-900">{{ ingressRule.id }}</div>
    </TableRow>
    <TableRow align="center">
      <Badge v-if="ingressRule.status === 'pending'" type="warning">Pending</Badge>
      <Badge v-else-if="ingressRule.status === 'applied'" type="success">Applied</Badge>
      <Badge v-else-if="ingressRule.status === 'failed'" type="danger">Failed</Badge>
      <Badge v-else-if="ingressRule.status === 'deleting'" type="danger">Deleting</Badge>
    </TableRow>
    <TableRow align="center">
      <div
        class="text-sm text-gray-900"
        :class="{
          'max-w-[200px] overflow-hidden text-ellipsis text-nowrap ': restrictTableWidth
        }">
        <a
          v-if="ingressRule.protocol === 'http' || ingressRule.protocol === 'https'"
          :href="ingressRule.protocol + '://' + ingressRule.domain.name + ':' + ingressRule.port.toString()"
          target="_blank"
          >{{ ingressRule.protocol }}://{{ ingressRule.domain.name }}:{{ ingressRule.port }}</a
        >
        <a v-else-if="ingressRule.protocol === 'tcp'" href="javascript:void(0);"
          >tcp://&lt;proxy-server-ip&gt;:{{ ingressRule.port }}</a
        >
        <a v-else-if="ingressRule.protocol === 'udp'" href="javascript:void(0);"
          >udp://&lt;proxy-server-ip&gt;:{{ ingressRule.port }}</a
        >
        <a v-else href="javascript:void(0);"><i>Unknown</i></a>
      </div>
    </TableRow>
    <TableRow align="center">
      <font-awesome-icon icon="fa-solid fa-arrow-right" />
    </TableRow>
    <TableRow align="center">
      <div class="text-sm text-gray-900">
        <Badge v-if="ingressRule.targetType === 'externalService'" type="warning">External Service</Badge>
        <Badge v-else-if="ingressRule.targetType === 'application'" type="success">Application</Badge>
        &nbsp;&nbsp;
        <a v-if="ingressRule.targetType === 'application'" href="javascript:void(0);"
          >{{ ingressRule.application.name }}:{{ ingressRule.targetPort }}</a
        >
        <a v-else href="javascript:void(0);">{{ ingressRule.externalService }}:{{ ingressRule.targetPort }}</a>
      </div>
    </TableRow>
    <TableRow align="center" flex>
      <div v-if="ingressRule.protocol === 'https' || ingressRule.protocol === 'http'" class="text-sm text-gray-900">
        <div v-if="ingressRule.authenticationType === 'none'" class="italic">No Authentication</div>
        <div
          v-else-if="ingressRule.authenticationType === 'basic'"
          class="flex items-center gap-2 text-sm text-gray-900">
          <font-awesome-icon icon="fa-solid fa-shield-halved" class="text-success-600" />
          <p><span class="font-medium">ACL</span> - {{ ingressRule.basicAuthAccessControlListName }}</p>
        </div>
      </div>
      <div v-else class="text-sm italic text-gray-900">N/A</div>
    </TableRow>
    <TableRow align="center" flex v-if="ingressRule.protocol === 'https'">
      <Badge v-if="ingressRule.httpsRedirect" type="success">Active</Badge>
      <Badge v-else type="secondary">Disabled</Badge>
    </TableRow>
    <TableRow align="center" v-else>
      <p class="text-sm font-medium italic text-gray-900">N/A</p>
    </TableRow>
    <TableRow align="right" flex>
      <FilledButton type="ghost" slim ref="actionsBtnRef" :click="onClickActions">
        <font-awesome-icon icon="fa-solid fa-ellipsis-vertical" />&nbsp;&nbsp;&nbsp;Show Actions
      </FilledButton>
    </TableRow>
  </tr>

  <div class="z-1 actions-menu" ref="actionsMenuRef" @click="closeMenu">
    <ul>
      <li v-if="ingressRule.httpsRedirect && ingressRule.protocol === 'https'" @click="disableHttpsRedirect">
        <font-awesome-icon icon="fa-solid fa-location-arrow" />&nbsp;&nbsp;&nbsp;Disable HTTPS Redirect
      </li>
      <li v-else-if="ingressRule.protocol === 'https'" @click="enableHttpsRedirect">
        <font-awesome-icon icon="fa-solid fa-location-arrow" />&nbsp;&nbsp;&nbsp;Enable HTTPS Redirect
      </li>
      <li @click="setupAuthentication" v-if="ingressRule.authenticationType === 'none'">
        <font-awesome-icon icon="fa-solid fa-shield-halved" />&nbsp;&nbsp;&nbsp;Setup Authentication
      </li>
      <li @click="disableAuthentication" v-if="ingressRule.authenticationType !== 'none'">
        <font-awesome-icon icon="fa-solid fa-shield-halved" />&nbsp;&nbsp;&nbsp;Disable Authentication
      </li>
      <li @click="recreateIngressRule">
        <font-awesome-icon icon="fa-solid fa-hammer" />&nbsp;&nbsp;&nbsp;Recreate & Fix
      </li>
      <li @click="deleteIngressRule">
        <p class="font-medium text-danger-500">
          <font-awesome-icon icon="fa-solid fa-trash" />&nbsp;&nbsp;&nbsp;Delete Ingress Rule
        </p>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.actions-menu {
  @apply absolute hidden rounded-md border border-gray-200 bg-white shadow-md;

  ul {
    li {
      @apply cursor-pointer px-4 py-2 text-sm text-gray-900 hover:bg-gray-100;
    }
  }
}
</style>
