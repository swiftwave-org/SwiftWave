<script setup>
import Switch from '@/views/components/Switch.vue'
import { computed, onMounted, reactive, ref } from 'vue'
import FilledButton from '@/views/components/FilledButton.vue'
import ModalDialog from '@/views/components/ModalDialog.vue'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'
import { useSystemConfigStore } from '@/store/systemConfig.js'
import { preventSpaceInput } from '@/vendor/utils.js'

const router = useRouter()
const systemConfigStore = useSystemConfigStore()
const isUpdateRequired = router.currentRoute.value.query.update === '1'
const isSimple = ref(true)
const isAdvanced = computed(() => !isSimple.value)
const toggleMode = () => (isSimple.value = !isSimple.value)
const formState = reactive({
  new_admin_credential: {
    username: '',
    password: ''
  },
  network_name: 'swiftwave_network',
  extra_restricted_ports: '',
  lets_encrypt: {
    email_address: '',
    staging_env: false
  },
  image_registry: {
    type: 'local',
    endpoint: '',
    namespace: '',
    username: '',
    password: ''
  },
  haproxy_config: {
    image: 'ghcr.io/swiftwave-org/haproxy:2.9'
  },
  udpproxy_config: {
    image: 'ghcr.io/swiftwave-org/udpproxy:latest'
  },
  pv_backup_config: {
    s3_config: {
      enabled: false,
      endpoint: '',
      region: 'us-east-1',
      bucket_name: '',
      access_key_id: '',
      secret_key: '',
      force_path_style: false
    }
  },
  pubsub_config: {
    type: 'local',
    buffer_length: 2000,
    redis_config: {
      host: 'localhost',
      port: 6379,
      password: '',
      database: 0
    }
  },
  task_queue_config: {
    type: 'local',
    remote_task_queue_type: 'none',
    max_outstanding_messages_per_queue: 1000,
    no_of_workers_per_queue: 1,
    amqp_config: {
      protocol: 'amqps',
      host: '',
      port: 5671,
      username: '',
      password: '',
      vhost: ''
    },
    redis_config: {
      host: '',
      port: 6379,
      password: '',
      database: 0
    }
  }
})
const timeCount = ref(5)
const setupSuccessful = ref(false)
const toggleStagingEnv = () => (formState.lets_encrypt.staging_env = !formState.lets_encrypt.staging_env)
const toggleImageRegistry = () =>
  (formState.image_registry.type = formState.image_registry.type === 'local' ? 'remote' : 'local')
const isRemoteImageRegistry = computed(() => formState.image_registry.type === 'remote')
const toggleS3Backup = () =>
  (formState.pv_backup_config.s3_config.enabled = !formState.pv_backup_config.s3_config.enabled)
const toggleS3ForcePathStyle = () =>
  (formState.pv_backup_config.s3_config.force_path_style = !formState.pv_backup_config.s3_config.force_path_style)
const isRemoteS3Backup = computed(() => formState.pv_backup_config.s3_config.enabled)

const isLocalTaskQueue = computed(() => formState.task_queue_config.type === 'local')
const isRedisTaskQueue = computed(
  () => formState.task_queue_config.remote_task_queue_type === 'redis' && formState.task_queue_config.type === 'remote'
)
const isAMQPTaskQueue = computed(
  () => formState.task_queue_config.remote_task_queue_type === 'amqp' && formState.task_queue_config.type === 'remote'
)

const switchTaskQueueType = (type) => {
  if (type === 'local') {
    formState.task_queue_config.type = 'local'
    formState.task_queue_config.remote_task_queue_type = 'none'
  } else if (type === 'redis') {
    formState.task_queue_config.type = 'remote'
    formState.task_queue_config.remote_task_queue_type = 'redis'
  } else if (type === 'amqp') {
    formState.task_queue_config.type = 'remote'
    formState.task_queue_config.remote_task_queue_type = 'amqp'
  }
}

const togglePubSub = () =>
  (formState.pubsub_config.type = formState.pubsub_config.type === 'local' ? 'remote' : 'local')
const isRemotePubSub = computed(() => formState.pubsub_config.type === 'remote')
const startCountDown = () => {
  const interval = setInterval(() => {
    timeCount.value--
    if (timeCount.value === 0) {
      clearInterval(interval)
      setupSuccessful.value = false
      router.push({ name: 'Maintenance' })
    }
  }, 1000)
}

const fetchDetails = async () => {
  if (!isUpdateRequired) return
  const { success, data } = await systemConfigStore.fetch()
  if (success) {
    Object.assign(formState, data)
  } else {
    toast.error('Failed to fetch the system configuration')
  }
}
onMounted(() => {
  fetchDetails()
})

const submitConfig = async () => {
  if (!isUpdateRequired) {
    if (!formState.new_admin_credential.username || !formState.new_admin_credential.password) {
      toast.error('Admin Username and Password are required')
      return
    }
  }
  let res = await systemConfigStore.submit(formState)
  if (res.success) {
    setupSuccessful.value = true
    toast.success(res.message)
    startCountDown()
  } else {
    toast.error(res.message)
  }
}
const updateConfig = async () => {
  let res = await systemConfigStore.update(formState)
  if (res.success) {
    setupSuccessful.value = true
    startCountDown()
  } else {
    toast.error(res.message)
  }
}
</script>

<template>
  <ModalDialog :is-open="setupSuccessful" non-cancelable>
    <template v-slot:header>
      <span v-show="!isUpdateRequired">🚀 Setup Successful</span>
      <span v-show="isUpdateRequired">🚀 Configuration Saved</span>
    </template>
    <template v-slot:body>
      <p class="mb-2">All the configurations are saved successfully. SwiftWave needs some time to be prepared.</p>
      <p>
        Redirecting to Maintenance Page in <b>{{ timeCount }}</b> seconds
      </p>
    </template>
  </ModalDialog>
  <div class="flex w-full max-w-7xl flex-col items-center gap-6 sm:px-0">
    <!-- Setup bar  -->
    <!-- logo | name | --spacing-- | -- switch [simple/advanced] -- -->
    <div class="flex w-full items-center justify-between py-2">
      <img src="@/assets/images/logo-full.png" alt="logo" class="max-h-10" />
      <div class="flex items-center gap-2 font-medium">
        Simple
        <Switch :enabled="!isSimple" :on-change="toggleMode" />
        Advanced
      </div>
    </div>
    <!--  New Admin info  -->
    <div class="info-section" v-if="!isUpdateRequired">
      <!--   label   -->
      <div class="label">
        <p>New Admin Credential</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Username</span>
          <span>Username for the new admin</span>
        </div>
        <input type="text" v-model="formState.new_admin_credential.username" @keydown="preventSpaceInput" />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Password</span>
          <span>Password for the new admin</span>
        </div>
        <input type="password" v-model="formState.new_admin_credential.password" />
      </div>
    </div>
    <!-- Networking section [advanced] -->
    <div class="info-section" v-if="isAdvanced">
      <!--   label   -->
      <div class="label">
        <p>Networking</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Overlay Network Name</span>
          <span>All the applications deployed through swiftwave will be connected to this overlay network</span>
        </div>
        <input type="text" v-model="formState.network_name" @keydown="preventSpaceInput" />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Extra Restricted Ports (comma seperated)</span>
          <span>Enter the ports you like to restrict in addition to the default restricted ports</span>
        </div>
        <input type="text" v-model="formState.extra_restricted_ports" @keydown="preventSpaceInput" />
      </div>
    </div>
    <!-- Lets Encrypt Info  -->
    <div class="info-section">
      <!--   label   -->
      <div class="label">
        <p>Let's Encrypt Account Info</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Email Address</span>
          <span>You will receive all the SSL expiry related notifications on this e-mail</span>
        </div>
        <input type="text" v-model="formState.lets_encrypt.email_address" @keydown="preventSpaceInput" />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Production Environment ?</span>
          <span>Do you need production SSL certificates which can be accepted by all browsers ?</span>
        </div>
        <div class="flex items-center gap-2 font-medium">
          <span>Yes (Recommended)</span>
          <Switch :on-change="toggleStagingEnv" :enabled="formState.lets_encrypt.staging_env" />
          <span>No (Staging Environment)</span>
        </div>
      </div>
    </div>
    <!-- Image Registry Config -->
    <div class="info-section">
      <!--   label   -->
      <div class="label">
        <p>Image Registry Info</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Registry Type</span>
          <span>It's always reliable to configure a remote image registry to store all the images</span>
        </div>
        <div class="flex items-center gap-2 font-medium">
          <span>Local</span>
          <Switch :on-change="toggleImageRegistry" :enabled="formState.image_registry.type === 'remote'" />
          <span>Remote</span>
        </div>
      </div>
      <div class="content" v-if="isRemoteImageRegistry">
        <div class="input-label">
          <span>Registry Endpoint</span>
          <span>For example, registry.digitalocean.com</span>
        </div>
        <input type="text" v-model="formState.image_registry.endpoint" />
      </div>
      <div class="content" v-if="isRemoteImageRegistry">
        <div class="input-label">
          <span>Registry Namespace</span>
          <span>You may leave this blank, if not required</span>
        </div>
        <input type="text" v-model="formState.image_registry.endpoint" />
      </div>
      <div class="content" v-if="isRemoteImageRegistry">
        <div class="input-label">
          <span>Registry Username</span>
          <span>Username to authenticate with the registry</span>
        </div>
        <input type="text" v-model="formState.image_registry.username" />
      </div>
      <div class="content" v-if="isRemoteImageRegistry">
        <div class="input-label">
          <span>Registry Password</span>
          <span>Password to authenticate with the registry</span>
        </div>
        <input type="password" v-model="formState.image_registry.password" />
      </div>
    </div>
    <!-- HAProxy config [advanced]  -->
    <div class="info-section" v-if="isAdvanced">
      <div class="label">
        <p>HAProxy Config</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Image</span>
          <span
            >Available Images -
            <a href="https://github.com/swiftwave-org/haproxy/pkgs/container/haproxy" target="_blank"
              >github.com/swiftwave-org/haproxy/pkgs/container/haproxy</a
            >
          </span>
        </div>
        <input type="text" v-model="formState.haproxy_config.image" />
      </div>
    </div>
    <!-- UDP Proxy Config [advanced] -->
    <div class="info-section" v-if="isAdvanced">
      <div class="label">
        <p>UDP Proxy Config</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Image</span>
          <span
            >Available Images -
            <a href="https://github.com/swiftwave-org/udpproxy/pkgs/container/udpproxy" target="_blank"
              >github.com/swiftwave-org/udpproxy/pkgs/container/udpproxy</a
            >
          </span>
        </div>
        <input type="text" v-model="formState.udpproxy_config.image" />
      </div>
    </div>
    <!-- PV Backup Config -->
    <div class="info-section">
      <div class="label">
        <p>Volume Backup Config</p>
        <span />
      </div>
      <div class="content">
        <div class="input-label">
          <span>S3 Enabled ?</span>
          <span>Do you want to store all the volume backups in S3 compatible storage ?</span>
        </div>
        <div class="flex items-center gap-2 font-medium">
          <span>Yes (Recommended)</span>
          <Switch :on-change="toggleS3Backup" :enabled="!formState.pv_backup_config.s3_config.enabled" />
          <span>No</span>
        </div>
      </div>
      <div class="content" v-if="isRemoteS3Backup">
        <div class="input-label">
          <span>S3 Endpoint</span>
          <span>For example, https://nyc3.digitaloceanspaces.com</span>
        </div>
        <input type="text" v-model="formState.pv_backup_config.s3_config.endpoint" />
      </div>
      <div class="content" v-if="isRemoteS3Backup">
        <div class="input-label">
          <span>S3 Region</span>
          <span>For example, us-east-1</span>
        </div>
        <input type="text" v-model="formState.pv_backup_config.s3_config.region" />
      </div>
      <div class="content" v-if="isRemoteS3Backup">
        <div class="input-label">
          <span>S3 Bucket Name</span>
          <span>For example, swiftwave-backups</span>
        </div>
        <input type="text" v-model="formState.pv_backup_config.s3_config.bucket_name" />
      </div>
      <div class="content" v-if="isRemoteS3Backup">
        <div class="input-label">
          <span>S3 Access Key ID</span>
          <span>Access key ID to authenticate with the S3 compatible storage</span>
        </div>
        <input type="text" v-model="formState.pv_backup_config.s3_config.access_key_id" />
      </div>
      <div class="content" v-if="isRemoteS3Backup">
        <div class="input-label">
          <span>S3 Secret Key</span>
          <span>Secret key to authenticate with the S3 compatible storage</span>
        </div>
        <input type="password" v-model="formState.pv_backup_config.s3_config.secret_key" />
      </div>
      <div class="content" v-if="isRemoteS3Backup">
        <div class="input-label">
          <span>Force Path Style ?</span>
          <span>Configures to use subdomain/virtual calling format</span>
        </div>
        <div class="flex items-center gap-2 font-medium">
          <span>Yes</span>
          <Switch
            :on-change="toggleS3ForcePathStyle"
            :enabled="!formState.pv_backup_config.s3_config.force_path_style" />
          <span>No</span>
        </div>
      </div>
    </div>
    <!-- PubSub config -->
    <div class="info-section">
      <div class="label">
        <p>PubSub Config</p>
        <span />
      </div>
      <div class="content" v-if="isAdvanced">
        <div class="input-label">
          <span>Buffer Length</span>
          <span>Maximum number of messages to store in buffer in the pubsub</span>
        </div>
        <input type="number" v-model="formState.pubsub_config.buffer_length" />
      </div>
      <div class="content">
        <div class="input-label">
          <span>PubSub Type</span>
          <span>It's always reliable to configure a remote pubsub to store all the messages</span>
        </div>
        <div class="flex items-center gap-2 font-medium">
          <span>Local</span>
          <Switch :on-change="togglePubSub" :enabled="formState.pubsub_config.type === 'remote'" />
          <span>Remote</span>
        </div>
      </div>
      <div class="rounded-md border border-success-600 bg-success-100 p-2 text-success-800" v-if="isRemotePubSub">
        You can configure Redis as a remote pubsub
      </div>
      <div class="content" v-if="isRemotePubSub">
        <div class="input-label">
          <span>Redis Host</span>
          <span>For example, localhost</span>
        </div>
        <input type="text" v-model="formState.pubsub_config.redis_config.host" />
      </div>
      <div class="content" v-if="isRemotePubSub">
        <div class="input-label">
          <span>Redis Port</span>
          <span>Port to connect with Redis server</span>
        </div>
        <input type="text" v-model="formState.pubsub_config.redis_config.port" />
      </div>
      <div class="content" v-if="isRemotePubSub">
        <div class="input-label">
          <span>Redis Password</span>
          <span>Password to authenticate with the pubsub</span>
        </div>
        <input type="password" v-model="formState.pubsub_config.redis_config.password" />
      </div>
      <div class="content" v-if="isRemotePubSub">
        <div class="input-label">
          <span>Redis Database ID</span>
          <span>Database to connect with the pubsub</span>
        </div>
        <input type="number" v-model="formState.pubsub_config.redis_config.database" />
      </div>
    </div>
    <!-- Task Queue Config -->
    <div class="info-section">
      <div class="label">
        <p>Task Queue Config</p>
        <span />
      </div>
      <div class="content" v-if="isAdvanced">
        <div class="input-label">
          <span>No of Workers Per Queue</span>
          <span>If you have a high load, you can increase the number of workers</span>
        </div>
        <input type="number" v-model="formState.task_queue_config.no_of_workers_per_queue" />
      </div>
      <div class="content">
        <div class="input-label">
          <span>Queue Type</span>
          <span>It's always reliable to configure a remote task queue to store all the messages</span>
        </div>
        <div class="multi-select">
          <div
            @click="switchTaskQueueType('local')"
            :class="{
              active: isLocalTaskQueue
            }">
            Local
          </div>
          <div
            @click="switchTaskQueueType('redis')"
            :class="{
              active: isRedisTaskQueue
            }">
            Redis
          </div>
          <div
            @click="switchTaskQueueType('amqp')"
            :class="{
              active: isAMQPTaskQueue
            }">
            AMQP
          </div>
        </div>
      </div>
      <div class="rounded-md border border-success-600 bg-success-100 p-2 text-success-800">
        <span v-if="isLocalTaskQueue"
          >Pre-configured PostgresSQL database will be used to store the tasks for reliability</span
        >
        <span v-else-if="isRedisTaskQueue">Configure Redis as task queue</span>
        <span v-else-if="isAMQPTaskQueue"
          >Configure AMQP compatible queue (RabbitMQ, LavinMQ, ActiveMQ etc.) as task queue</span
        >
      </div>

      <!--  Local Task Queue config    -->
      <div class="content" v-if="isAdvanced && isLocalTaskQueue">
        <div class="input-label">
          <span>Max Outstanding Messages Per Queue</span>
          <span>Maximum number of tasks</span>
        </div>
        <input type="number" v-model="formState.task_queue_config.max_outstanding_messages_per_queue" />
      </div>

      <!--   Redis Task Queue Config   -->
      <div class="content" v-if="isRedisTaskQueue">
        <div class="input-label">
          <span>Redis Host</span>
          <span>For example, localhost</span>
        </div>
        <input type="text" v-model="formState.task_queue_config.redis_config.host" />
      </div>
      <div class="content" v-if="isRedisTaskQueue">
        <div class="input-label">
          <span>Redis Port</span>
          <span>Port to connect with Redis server</span>
        </div>
        <input type="text" v-model="formState.task_queue_config.redis_config.port" />
      </div>
      <div class="content" v-if="isRedisTaskQueue">
        <div class="input-label">
          <span>Redis Password</span>
          <span>Password to authenticate with redis server</span>
        </div>
        <input type="password" v-model="formState.task_queue_config.redis_config.password" />
      </div>
      <div class="content" v-if="isRedisTaskQueue">
        <div class="input-label">
          <span>Redis Database ID</span>
          <span>Database to connect for task queue</span>
        </div>
        <input type="number" v-model="formState.task_queue_config.redis_config.database" />
      </div>

      <!--   AMQP Config   -->

      <div class="content" v-if="isAMQPTaskQueue">
        <div class="input-label">
          <span>Queue Protocol</span>
          <span>For example, amqp</span>
        </div>
        <select v-model="formState.task_queue_config.amqp_config.protocol">
          <option value="amqp">AMQP</option>
          <option value="amqps">AMQPS</option>
        </select>
      </div>
      <div class="content" v-if="isAMQPTaskQueue">
        <div class="input-label">
          <span>Queue Host</span>
          <span>For example, localhost</span>
        </div>
        <input type="text" v-model="formState.task_queue_config.amqp_config.host" />
      </div>
      <div class="content" v-if="isAMQPTaskQueue">
        <div class="input-label">
          <span>Queue Port</span>
          <span>Port to connect with AMQP server</span>
        </div>
        <input type="number" v-model="formState.task_queue_config.amqp_config.port" />
      </div>
      <div class="content" v-if="isAMQPTaskQueue">
        <div class="input-label">
          <span>Queue Username</span>
          <span>Username to authenticate with the queue</span>
        </div>
        <input type="text" v-model="formState.task_queue_config.amqp_config.username" />
      </div>
      <div class="content" v-if="isAMQPTaskQueue">
        <div class="input-label">
          <span>Queue Password</span>
          <span>Password to authenticate with the queue</span>
        </div>
        <input type="password" v-model="formState.task_queue_config.amqp_config.password" />
      </div>
      <div class="content" v-if="isAMQPTaskQueue">
        <div class="input-label">
          <span>Queue VHost</span>
          <span>Virtual host to connect with the queue</span>
        </div>
        <input type="text" v-model="formState.task_queue_config.amqp_config.vhost" />
      </div>
    </div>
    <!--  Confirm btn  -->
    <FilledButton
      class="mt-8"
      :click="submitConfig"
      :loading="systemConfigStore.isSubmitting"
      v-show="!isUpdateRequired">
      <font-awesome-icon icon="fa-solid fa-rocket" class="mr-3" />
      Save Configuration & Start Swiftwave
    </FilledButton>
    <FilledButton class="mt-8" :click="updateConfig" :loading="systemConfigStore.isUpdating" v-show="isUpdateRequired">
      <font-awesome-icon icon="fa-solid fa-rocket" class="mr-3" />
      Update Configuration & Re-Start Swiftwave
    </FilledButton>
  </div>
</template>

<style scoped>
.info-section {
  @apply flex w-full flex-col gap-3;

  .label {
    @apply flex w-full items-center gap-5 text-lg font-medium;

    p {
      @apply whitespace-nowrap;
    }

    span {
      @apply h-0.5 w-full bg-gray-300;
    }
  }

  .content {
    @apply grid grid-cols-2 gap-4;

    .input-label {
      @apply flex flex-col;

      span:first-child {
        @apply text-base font-medium text-gray-800;
      }

      span:last-child {
        @apply text-sm text-gray-500;
      }
    }

    select {
      @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500;
    }

    input {
      @apply w-full rounded-md border border-gray-300 p-2 shadow-sm focus:border-primary-500 focus:ring-primary-500;
    }
  }

  .multi-select {
    @apply flex h-fit w-min overflow-hidden rounded-md border border-secondary-400;

    div {
      @apply cursor-pointer border-r border-secondary-400 bg-secondary-100 px-5 py-2 transition-all hover:bg-secondary-400 hover:text-white;
    }

    div:last-child {
      @apply border-0;
    }

    .active {
      @apply bg-primary-600  text-white;
    }
  }
}
</style>
