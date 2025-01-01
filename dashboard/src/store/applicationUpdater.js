import { defineStore } from 'pinia'
import { useMutation, useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { computed, reactive, ref, watch } from 'vue'
import { v4 as uuidv4 } from 'uuid'
import { useRouter } from 'vue-router'

export default function newApplicationUpdater(applicationId) {
  const storeName = 'application_updater_' + applicationId
  return defineStore(storeName, () => {
    const router = useRouter()
    const isConfigurationUpdated = ref(false)
    const applyConfigurationChanges = () => {
      const appState = mergeChangesWithExistingApplicationDetails()
      appState.gitCredentialID = parseInt(appState.gitCredentialID)
      appState.imageRegistryCredentialID = parseInt(appState.imageRegistryCredentialID)
      appState.gitCredentialID = appState.gitCredentialID === 0 ? null : appState.gitCredentialID
      appState.imageRegistryCredentialID =
        appState.imageRegistryCredentialID === 0 ? null : appState.imageRegistryCredentialID
      deployApplication({
        input: appState,
        id: applicationId
      })
    }

    const cancelConfigurationChanges = () => {
      resetDetailsToApplicationDetails()
    }

    const {
      mutate: deployApplication,
      loading: isDeployRequestSubmitting,
      onDone: onDeployApplicationMutationDone,
      onError: onDeployApplicationMutationError
    } = useMutation(gql`
      mutation ($id: String!, $input: ApplicationInput!) {
        updateApplication(id: $id, input: $input) {
          id
          name
        }
      }
    `)

    onDeployApplicationMutationError((error) => {
      alert('Failed to update application\n' + error.message)
    })

    onDeployApplicationMutationDone(() => {
      refetchApplicationDetails()
      router.push({ name: 'Application Details Deployments', params: { id: applicationId } })
      alert('Application updated successfully\nNOTE: Wait for a few seconds to apply new changes')
    })

    const {
      result: applicationDetailsRaw,
      refetch: refetchApplicationDetails,
      loading: applicationDetailsLoading
    } = useQuery(
      gql`
        query ($id: String!) {
          application(id: $id) {
            deploymentMode
            replicas
            command
            hostname
            resourceLimit {
              memoryMb
            }
            reservedResource {
              memoryMb
            }
            environmentVariables {
              key
              value
            }
            persistentVolumeBindings {
              persistentVolumeID
              mountingPath
            }
            configMounts {
              content
              mountingPath
              uid
              gid
            }
            latestDeployment {
              upstreamType
              dockerfile
              buildArgs {
                key
                value
              }
              gitProvider
              gitEndpoint
              gitCredentialID
              repositoryUrl
              repositoryName
              repositoryOwner
              repositoryBranch
              codePath
              imageRegistryCredentialID
              dockerImage
              sourceCodeCompressedFileName
            }
            capabilities
            sysctls
            customHealthCheck {
              enabled
              test_command
              interval_seconds
              timeout_seconds
              start_period_seconds
              start_interval_seconds
              retries
            }
            dockerProxyHost
            preferredServerHostnames
            dockerProxyConfig {
              enabled
              permission {
                ping
                version
                info
                events
                auth
                secrets
                build
                commit
                configs
                containers
                distribution
                exec
                grpc
                images
                networks
                nodes
                plugins
                services
                session
                swarm
                system
                tasks
                volumes
              }
            }
          }
        }
      `,
      {
        id: applicationId
      }
    )

    watch(applicationDetailsRaw, () => {
      resetDetailsToApplicationDetails()
    })

    const dockerProxyHost = computed(() => {
      const applicationExistingDetails = applicationDetailsRaw.value?.application ?? {}
      if (applicationExistingDetails.length === 0) {
        return ''
      }
      return applicationExistingDetails.dockerProxyHost
    })

    const resetDetailsToApplicationDetails = () => {
      const environmentVariables = applicationDetailsRaw.value?.application?.environmentVariables ?? []
      let keys = []
      let map = {}
      environmentVariables.forEach((variable) => {
        const z = uuidv4()
        keys.push(z)
        map[z] = {
          name: variable.key,
          value: variable.value
        }
      })
      environmentVariableDetails.keys = keys
      environmentVariableDetails.map = map

      const persistentVolumeBindings = applicationDetailsRaw.value?.application?.persistentVolumeBindings ?? []
      keys = []
      map = {}
      persistentVolumeBindings.forEach((binding) => {
        const z = uuidv4()
        keys.push(z)
        map[z] = {
          persistentVolumeID: binding.persistentVolumeID,
          mountingPath: binding.mountingPath
        }
      })
      persistentVolumeBindingsDetails.keys = keys
      persistentVolumeBindingsDetails.map = map

      const configMounts = applicationDetailsRaw.value?.application?.configMounts ?? []
      keys = []
      map = {}
      configMounts.forEach((binding) => {
        const z = uuidv4()
        keys.push(z)
        map[z] = {
          content: binding.content,
          mountingPath: binding.mountingPath,
          uid: binding.uid,
          gid: binding.gid
        }
      })
      configMountDetails.keys = keys
      configMountDetails.map = map

      const applicationConfiguration = applicationDetailsRaw.value?.application ?? {}
      deploymentConfigurationDetails.deploymentMode = applicationConfiguration.deploymentMode
      deploymentConfigurationDetails.replicas = applicationConfiguration.replicas
      deploymentConfigurationDetails.hostname = applicationConfiguration.hostname
      deploymentConfigurationDetails.resourceLimit.memoryMb = applicationConfiguration.resourceLimit.memoryMb
      deploymentConfigurationDetails.reservedResource.memoryMb = applicationConfiguration.reservedResource.memoryMb
      deploymentConfigurationDetails.preferredServerHostnames = applicationConfiguration.preferredServerHostnames
      deploymentConfigurationDetails.dockerProxyConfig.enabled = applicationConfiguration.dockerProxyConfig.enabled
      deploymentConfigurationDetails.dockerProxyConfig.permission.ping =
        applicationConfiguration.dockerProxyConfig.permission.ping
      deploymentConfigurationDetails.dockerProxyConfig.permission.version =
        applicationConfiguration.dockerProxyConfig.permission.version
      deploymentConfigurationDetails.dockerProxyConfig.permission.info =
        applicationConfiguration.dockerProxyConfig.permission.info
      deploymentConfigurationDetails.dockerProxyConfig.permission.events =
        applicationConfiguration.dockerProxyConfig.permission.events
      deploymentConfigurationDetails.dockerProxyConfig.permission.auth =
        applicationConfiguration.dockerProxyConfig.permission.auth
      deploymentConfigurationDetails.dockerProxyConfig.permission.secrets =
        applicationConfiguration.dockerProxyConfig.permission.secrets
      deploymentConfigurationDetails.dockerProxyConfig.permission.build =
        applicationConfiguration.dockerProxyConfig.permission.build
      deploymentConfigurationDetails.dockerProxyConfig.permission.commit =
        applicationConfiguration.dockerProxyConfig.permission.commit
      deploymentConfigurationDetails.dockerProxyConfig.permission.configs =
        applicationConfiguration.dockerProxyConfig.permission.configs
      deploymentConfigurationDetails.dockerProxyConfig.permission.containers =
        applicationConfiguration.dockerProxyConfig.permission.containers
      deploymentConfigurationDetails.dockerProxyConfig.permission.distribution =
        applicationConfiguration.dockerProxyConfig.permission.distribution
      deploymentConfigurationDetails.dockerProxyConfig.permission.exec =
        applicationConfiguration.dockerProxyConfig.permission.exec
      deploymentConfigurationDetails.dockerProxyConfig.permission.grpc =
        applicationConfiguration.dockerProxyConfig.permission.grpc
      deploymentConfigurationDetails.dockerProxyConfig.permission.images =
        applicationConfiguration.dockerProxyConfig.permission.images
      deploymentConfigurationDetails.dockerProxyConfig.permission.networks =
        applicationConfiguration.dockerProxyConfig.permission.networks
      deploymentConfigurationDetails.dockerProxyConfig.permission.nodes =
        applicationConfiguration.dockerProxyConfig.permission.nodes
      deploymentConfigurationDetails.dockerProxyConfig.permission.plugins =
        applicationConfiguration.dockerProxyConfig.permission.plugins
      deploymentConfigurationDetails.dockerProxyConfig.permission.services =
        applicationConfiguration.dockerProxyConfig.permission.services
      deploymentConfigurationDetails.dockerProxyConfig.permission.session =
        applicationConfiguration.dockerProxyConfig.permission.session
      deploymentConfigurationDetails.dockerProxyConfig.permission.swarm =
        applicationConfiguration.dockerProxyConfig.permission.swarm
      deploymentConfigurationDetails.dockerProxyConfig.permission.system =
        applicationConfiguration.dockerProxyConfig.permission.system
      deploymentConfigurationDetails.dockerProxyConfig.permission.tasks =
        applicationConfiguration.dockerProxyConfig.permission.tasks
      deploymentConfigurationDetails.dockerProxyConfig.permission.volumes =
        applicationConfiguration.dockerProxyConfig.permission.volumes
      deploymentConfigurationDetails.customHealthCheck.enabled = applicationConfiguration.customHealthCheck.enabled
      deploymentConfigurationDetails.customHealthCheck.test_command =
        applicationConfiguration.customHealthCheck.test_command
      deploymentConfigurationDetails.customHealthCheck.interval_seconds =
        applicationConfiguration.customHealthCheck.interval_seconds
      deploymentConfigurationDetails.customHealthCheck.timeout_seconds =
        applicationConfiguration.customHealthCheck.timeout_seconds
      deploymentConfigurationDetails.customHealthCheck.start_period_seconds =
        applicationConfiguration.customHealthCheck.start_period_seconds
      deploymentConfigurationDetails.customHealthCheck.start_interval_seconds =
        applicationConfiguration.customHealthCheck.start_interval_seconds
      deploymentConfigurationDetails.customHealthCheck.retries = applicationConfiguration.customHealthCheck.retries
      sourceConfigurationRef.command = applicationConfiguration.command
      sourceConfigurationRef.gitCredentialID = applicationConfiguration.latestDeployment.gitCredentialID
      sourceConfigurationRef.gitProvider = applicationConfiguration.latestDeployment.gitProvider
      sourceConfigurationRef.gitEndpoint = applicationConfiguration.latestDeployment.gitEndpoint
      sourceConfigurationRef.repositoryUrl = applicationConfiguration.latestDeployment.repositoryUrl
      sourceConfigurationRef.repositoryName = applicationConfiguration.latestDeployment.repositoryName
      sourceConfigurationRef.repositoryOwner = applicationConfiguration.latestDeployment.repositoryOwner
      sourceConfigurationRef.repositoryBranch = applicationConfiguration.latestDeployment.repositoryBranch
      sourceConfigurationRef.codePath = applicationConfiguration.latestDeployment.codePath
      sourceConfigurationRef.imageRegistryCredentialID =
        applicationConfiguration.latestDeployment.imageRegistryCredentialID
      sourceConfigurationRef.dockerImage = applicationConfiguration.latestDeployment.dockerImage
      sourceConfigurationRef.sourceCodeCompressedFileName =
        applicationConfiguration.latestDeployment.sourceCodeCompressedFileName
      sourceConfigurationRef.dockerfile = applicationConfiguration.latestDeployment.dockerfile

      // reset isConfigurationUpdated
      isConfigurationUpdated.value = false
    }

    const environmentVariableDetails = reactive({
      keys: [],
      map: {}
    })

    const persistentVolumeBindingsDetails = reactive({
      keys: [],
      map: {}
    })

    const configMountDetails = reactive({
      keys: [],
      map: {}
    })

    const deploymentConfigurationDetails = reactive({
      deploymentMode: '',
      replicas: 0,
      hostname: '',
      resourceLimit: {
        memoryMb: 0
      },
      reservedResource: {
        memoryMb: 0
      },
      customHealthCheck: {
        enabled: false,
        test_command: '',
        interval_seconds: 0,
        timeout_seconds: 0,
        start_period_seconds: 0,
        start_interval_seconds: 0,
        retries: 0
      },
      preferredServerHostnames: [],
      dockerProxyConfig: {
        enabled: false,
        permission: {
          ping: 'none',
          version: 'none',
          info: 'none',
          events: 'none',
          auth: 'none',
          secrets: 'none',
          build: 'none',
          commit: 'none',
          configs: 'none',
          containers: 'none',
          distribution: 'none',
          exec: 'none',
          grpc: 'none',
          images: 'none',
          networks: 'none',
          nodes: 'none',
          plugins: 'none',
          services: 'none',
          session: 'none',
          swarm: 'none',
          system: 'none',
          tasks: 'none',
          volumes: 'none'
        }
      }
    })

    const updatePreferredServerHostnames = (hostnames) => {
      deploymentConfigurationDetails.preferredServerHostnames = [...hostnames]
      triggerUpdateHook()
    }

    const preferredServerHostnamesStr = computed(() => {
      return deploymentConfigurationDetails.preferredServerHostnames.join(', ')
    })

    const enableDockerProxy = () => {
      deploymentConfigurationDetails.dockerProxyConfig.enabled = true
      triggerUpdateHook()
    }

    const disableDockerProxy = () => {
      deploymentConfigurationDetails.dockerProxyConfig.enabled = false
      triggerUpdateHook()
    }

    const dockerProxyPermissionChanged = () => {
      // do the changes in ui directly and trigger this hook
      triggerUpdateHook()
    }

    const sourceConfigurationRef = reactive({
      command: '',
      gitCredentialID: 0,
      gitProvider: '',
      gitEndpoint: '',
      repositoryUrl: '',
      repositoryBranch: '',
      codePath: '',
      imageRegistryCredentialID: 0,
      dockerImage: '',
      sourceCodeCompressedFileName: '',
      dockerfile: '',
      buildArgs: {}
    })

    const addEnvironmentVariable = () => {
      const key = uuidv4()
      environmentVariableDetails.keys.push(key)
      environmentVariableDetails.map[key] = {
        name: '',
        value: ''
      }
      triggerUpdateHook()
    }

    const deleteEnvironmentVariable = (key) => {
      let keys
      keys = environmentVariableDetails.keys.filter((k) => k !== key)
      environmentVariableDetails.keys = keys
      delete environmentVariableDetails.map[key]
      triggerUpdateHook()
    }

    const onEnvironmentVariableNameChange = (key, value) => {
      environmentVariableDetails.map[key].name = value
      triggerUpdateHook()
    }

    const onEnvironmentVariableValueChange = (key, value) => {
      environmentVariableDetails.map[key].value = value
      triggerUpdateHook()
    }

    const addPersistentVolumeBinding = () => {
      const key = uuidv4()
      persistentVolumeBindingsDetails.keys.push(key)
      persistentVolumeBindingsDetails.map[key] = {
        persistentVolumeID: -1,
        mountingPath: ''
      }
      triggerUpdateHook()
    }

    const deletePersistentVolumeBinding = (key) => {
      let keys
      keys = persistentVolumeBindingsDetails.keys.filter((k) => k !== key)
      persistentVolumeBindingsDetails.keys = keys
      delete persistentVolumeBindingsDetails.map[key]
      triggerUpdateHook()
    }

    const onPersistentVolumeChange = (key, value) => {
      persistentVolumeBindingsDetails.map[key].persistentVolumeID = value
      triggerUpdateHook()
    }

    const onPersistentVolumeMountingPathChange = (key, value) => {
      persistentVolumeBindingsDetails.map[key].mountingPath = value
      triggerUpdateHook()
    }

    const addConfigMount = (details) => {
      const mountingPath = details.mountingPath
      // check if mounting path is already used
      for (const key in configMountDetails.map) {
        if (configMountDetails.map[key].mountingPath === mountingPath) {
          throw new Error('Mounting path already used')
        }
      }
      const key = uuidv4()
      configMountDetails.keys.push(key)
      configMountDetails.map[key] = details
      triggerUpdateHook()
    }

    const deleteConfigMount = (key) => {
      let keys
      keys = configMountDetails.keys.filter((k) => k !== key)
      configMountDetails.keys = keys
      delete configMountDetails.map[key]
      triggerUpdateHook()
    }

    const onConfigMountContentChange = (key, value) => {
      configMountDetails.map[key].content = value
      triggerUpdateHook()
    }

    // eslint-disable-next-line no-unused-vars
    const changeDeploymentStrategy = (switchStatus) => {
      alert(
        'Sorry, for change deployment strategy you need to delete and re-create the application\nIn future, we will support this feature'
      )

      // TODO: will be supported in future
      // if (switchStatus) {
      //   deploymentConfigurationDetails.deploymentMode = 'global'
      //   deploymentConfigurationDetails.replicas = 0
      // } else {
      //   deploymentConfigurationDetails.deploymentMode = 'replicated'
      //   deploymentConfigurationDetails.replicas = 1
      // }
      // triggerUpdateHook()
    }

    const onMemoryLimitChanged = (value) => {
      deploymentConfigurationDetails.resourceLimit.memoryMb = value
      triggerUpdateHook()
    }

    const onMemoryReservedChanged = (value) => {
      deploymentConfigurationDetails.reservedResource.memoryMb = value
      triggerUpdateHook()
    }

    const replicasCountChanged = () => {
      triggerUpdateHook()
    }

    // someInfoUpdated
    const triggerUpdateHook = () => {
      isConfigurationUpdated.value = checkIfApplicationDetailsAreChanged()
    }

    const { result: applicationExistingDetailsResult } = useQuery(
      gql`
        query ($id: String!) {
          application(id: $id) {
            name
            deploymentMode
            command
            hostname
            replicas
            resourceLimit {
              memoryMb
            }
            reservedResource {
              memoryMb
            }
            environmentVariables {
              key
              value
            }
            persistentVolumeBindings {
              id
              persistentVolumeID
              mountingPath
            }
            configMounts {
              content
              mountingPath
              uid
              gid
            }
            latestDeployment {
              upstreamType
              dockerfile
              buildArgs {
                key
                value
              }
              gitEndpoint
              gitProvider
              gitCredentialID
              repositoryUrl
              repositoryBranch
              codePath
              imageRegistryCredentialID
              dockerImage
              sourceCodeCompressedFileName
            }
            capabilities
            sysctls
            applicationGroupID
            customHealthCheck {
              enabled
              test_command
              interval_seconds
              timeout_seconds
              start_period_seconds
              start_interval_seconds
              retries
            }
            preferredServerHostnames
            dockerProxyConfig {
              enabled
              permission {
                ping
                version
                info
                events
                auth
                secrets
                build
                commit
                configs
                containers
                distribution
                exec
                grpc
                images
                networks
                nodes
                plugins
                services
                session
                swarm
                system
                tasks
                volumes
              }
            }
          }
        }
      `,
      {
        id: applicationId
      },
      {
        fetchPolicy: 'no-cache',
        nextFetchPolicy: 'no-cache'
      }
    )

    function checkIfApplicationDetailsAreChanged() {
      const applicationExistingDetails = applicationExistingDetailsResult.value?.application ?? {}

      // check if deployment mode is changed
      if (applicationExistingDetails.deploymentMode !== deploymentConfigurationDetails.deploymentMode) {
        return true
      }
      // check if replica count is changed
      if (applicationExistingDetails.replicas.toString() !== deploymentConfigurationDetails.replicas.toString()) {
        return true
      }
      // check if container hostname is changed
      if (applicationExistingDetails.hostname !== deploymentConfigurationDetails.hostname) {
        return true
      }
      // check if deploy proxy config changed
      if (
        applicationExistingDetails.dockerProxyConfig.enabled !==
          deploymentConfigurationDetails.dockerProxyConfig.enabled ||
        applicationExistingDetails.dockerProxyConfig.permission.ping !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.ping ||
        applicationExistingDetails.dockerProxyConfig.permission.version !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.version ||
        applicationExistingDetails.dockerProxyConfig.permission.info !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.info ||
        applicationExistingDetails.dockerProxyConfig.permission.events !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.events ||
        applicationExistingDetails.dockerProxyConfig.permission.auth !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.auth ||
        applicationExistingDetails.dockerProxyConfig.permission.secrets !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.secrets ||
        applicationExistingDetails.dockerProxyConfig.permission.build !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.build ||
        applicationExistingDetails.dockerProxyConfig.permission.commit !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.commit ||
        applicationExistingDetails.dockerProxyConfig.permission.configs !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.configs ||
        applicationExistingDetails.dockerProxyConfig.permission.containers !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.containers ||
        applicationExistingDetails.dockerProxyConfig.permission.distribution !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.distribution ||
        applicationExistingDetails.dockerProxyConfig.permission.exec !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.exec ||
        applicationExistingDetails.dockerProxyConfig.permission.grpc !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.grpc ||
        applicationExistingDetails.dockerProxyConfig.permission.images !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.images ||
        applicationExistingDetails.dockerProxyConfig.permission.networks !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.networks ||
        applicationExistingDetails.dockerProxyConfig.permission.nodes !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.nodes ||
        applicationExistingDetails.dockerProxyConfig.permission.plugins !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.plugins ||
        applicationExistingDetails.dockerProxyConfig.permission.services !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.services ||
        applicationExistingDetails.dockerProxyConfig.permission.session !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.session ||
        applicationExistingDetails.dockerProxyConfig.permission.swarm !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.swarm ||
        applicationExistingDetails.dockerProxyConfig.permission.system !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.system ||
        applicationExistingDetails.dockerProxyConfig.permission.tasks !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.tasks ||
        applicationExistingDetails.dockerProxyConfig.permission.volumes !==
          deploymentConfigurationDetails.dockerProxyConfig.permission.volumes
      ) {
        return true
      }

      // check if custom health check is changed
      if (
        applicationExistingDetails.customHealthCheck.enabled !==
          deploymentConfigurationDetails.customHealthCheck.enabled ||
        applicationExistingDetails.customHealthCheck.test_command !==
          deploymentConfigurationDetails.customHealthCheck.test_command ||
        applicationExistingDetails.customHealthCheck.interval_seconds !==
          deploymentConfigurationDetails.customHealthCheck.interval_seconds ||
        applicationExistingDetails.customHealthCheck.timeout_seconds !==
          deploymentConfigurationDetails.customHealthCheck.timeout_seconds ||
        applicationExistingDetails.customHealthCheck.start_period_seconds !==
          deploymentConfigurationDetails.customHealthCheck.start_period_seconds ||
        applicationExistingDetails.customHealthCheck.start_interval_seconds !==
          deploymentConfigurationDetails.customHealthCheck.start_interval_seconds ||
        applicationExistingDetails.customHealthCheck.retries !==
          deploymentConfigurationDetails.customHealthCheck.retries
      ) {
        return true
      }

      // check if preferred server hostnames are changed
      if (
        applicationExistingDetails.preferredServerHostnames.length !==
        deploymentConfigurationDetails.preferredServerHostnames.length
      ) {
        return true
      } else {
        if (
          JSON.stringify(applicationExistingDetails.preferredServerHostnames.sort()) !==
          JSON.stringify(deploymentConfigurationDetails.preferredServerHostnames.sort())
        ) {
          return true
        }
      }

      // check if resource limit is changed
      if (applicationExistingDetails.resourceLimit.memoryMb !== deploymentConfigurationDetails.resourceLimit.memoryMb) {
        return true
      }
      // check if reserved resource is changed
      if (
        applicationExistingDetails.reservedResource.memoryMb !==
        deploymentConfigurationDetails.reservedResource.memoryMb
      ) {
        return true
      }
      // check if environment variables are changed
      const existingEnvironmentVariables = applicationExistingDetails.environmentVariables ?? []
      const existingEnvironmentVariableKeys = existingEnvironmentVariables.map((variable) => variable.key)
      const existingEnvironmentVariableMap = existingEnvironmentVariables.reduce((map, variable) => {
        map[variable.key] = variable.value
        return map
      }, {})
      const newEnvironmentVariableKeys = environmentVariableDetails.keys
      const newEnvironmentVariableMap = environmentVariableDetails.keys.reduce((map, key) => {
        map[environmentVariableDetails.map[key].name] = environmentVariableDetails.map[key].value
        return map
      }, {})
      if (existingEnvironmentVariableKeys.length !== newEnvironmentVariableKeys.length) {
        return true
      }
      for (let i = 0; i < existingEnvironmentVariableKeys.length; i++) {
        const key = existingEnvironmentVariableKeys[i]
        if (existingEnvironmentVariableMap[key] !== newEnvironmentVariableMap[key]) {
          return true
        }
      }
      // check if config mounts are changed
      const existingConfigMounts = applicationExistingDetails.configMounts ?? []
      const existingConfigMountKeys = existingConfigMounts.map((mount) => mount.mountingPath)
      const existingConfigMountMap = existingConfigMounts.reduce((map, mount) => {
        map[mount.mountingPath] = mount
        return map
      }, {})
      const newConfigMountKeys = configMountDetails.keys
      const newConfigMountMap = configMountDetails.keys.reduce((map, key) => {
        map[configMountDetails.map[key].mountingPath] = configMountDetails.map[key]
        return map
      }, {})
      if (existingConfigMountKeys.length !== newConfigMountKeys.length) {
        return true
      }
      for (let i = 0; i < existingConfigMountKeys.length; i++) {
        const key = existingConfigMountKeys[i]
        if (
          existingConfigMountMap[key].content !== newConfigMountMap[key].content ||
          existingConfigMountMap[key].uid !== newConfigMountMap[key].uid ||
          existingConfigMountMap[key].gid !== newConfigMountMap[key].gid
        ) {
          return true
        }
      }
      // check if persistent volume bindings are changed
      const existingPersistentVolumeBindings = applicationExistingDetails.persistentVolumeBindings ?? []
      const existingPersistentVolumeBindingKeys = existingPersistentVolumeBindings.map(
        (binding) => binding.persistentVolumeID
      )
      const existingPersistentVolumeBindingMap = existingPersistentVolumeBindings.reduce((map, binding) => {
        map[binding.persistentVolumeID] = binding.mountingPath
        return map
      }, {})
      const newPersistentVolumeBindingKeys = persistentVolumeBindingsDetails.keys
      const newPersistentVolumeBindingMap = persistentVolumeBindingsDetails.keys.reduce((map, key) => {
        map[persistentVolumeBindingsDetails.map[key].persistentVolumeID] =
          persistentVolumeBindingsDetails.map[key].mountingPath
        return map
      }, {})
      if (existingPersistentVolumeBindingKeys.length !== newPersistentVolumeBindingKeys.length) {
        return true
      }
      for (let i = 0; i < existingPersistentVolumeBindingKeys.length; i++) {
        const key = existingPersistentVolumeBindingKeys[i]
        if (existingPersistentVolumeBindingMap[key] !== newPersistentVolumeBindingMap[key]) {
          return true
        }
      }

      // check if any source configuration is changed
      if (
        parseInt(sourceConfigurationRef.gitCredentialID) !== applicationExistingDetails.latestDeployment.gitCredentialID
      ) {
        return true
      }
      if (sourceConfigurationRef) {
        // check if source configuration is changed
        if (sourceConfigurationRef.repositoryUrl !== applicationExistingDetails.latestDeployment.repositoryUrl) {
          return true
        }
        if (sourceConfigurationRef.repositoryBranch !== applicationExistingDetails.latestDeployment.repositoryBranch) {
          return true
        }
        if (sourceConfigurationRef.codePath !== applicationExistingDetails.latestDeployment.codePath) {
          return true
        }
        if (
          sourceConfigurationRef.imageRegistryCredentialID !==
          applicationExistingDetails.latestDeployment.imageRegistryCredentialID
        ) {
          return true
        }
        if (sourceConfigurationRef.dockerImage !== applicationExistingDetails.latestDeployment.dockerImage) {
          return true
        }
        if (sourceConfigurationRef.command !== applicationExistingDetails.command) {
          return true
        }
        if (
          sourceConfigurationRef.sourceCodeCompressedFileName !==
          applicationExistingDetails.latestDeployment.sourceCodeCompressedFileName
        ) {
          return true
        }
        if (sourceConfigurationRef.dockerfile !== applicationExistingDetails.latestDeployment.dockerfile) {
          return true
        }
      }
      // check if build args are changed
      let existingBuildArgs = {}
      ;(applicationExistingDetails.latestDeployment.buildArgs ?? []).forEach((buildArg) => {
        existingBuildArgs[buildArg.key] = buildArg.value
      })
      if (Object.keys(existingBuildArgs).length !== Object.keys(sourceConfigurationRef.buildArgs).length) {
        return true
      }
      for (const key in existingBuildArgs) {
        if (existingBuildArgs[key] !== sourceConfigurationRef.buildArgs[key]) {
          return true
        }
      }
      return false
    }

    const mergeChangesWithExistingApplicationDetails = () => {
      const applicationExistingDetails = applicationExistingDetailsResult.value?.application ?? {}
      if (applicationExistingDetails.length === 0) {
        alert('Failed to fetch application details\nTry again after re-loading the page')
      }

      return {
        name: applicationExistingDetails.name,
        upstreamType: applicationExistingDetails.latestDeployment.upstreamType, // TODO Not allowed to change
        command: sourceConfigurationRef.command,
        deploymentMode: deploymentConfigurationDetails.deploymentMode,
        replicas: deploymentConfigurationDetails.replicas,
        hostname: deploymentConfigurationDetails.hostname,
        resourceLimit: {
          memoryMb: deploymentConfigurationDetails.resourceLimit.memoryMb
        },
        reservedResource: {
          memoryMb: deploymentConfigurationDetails.reservedResource.memoryMb
        },
        buildArgs: Object.entries(sourceConfigurationRef.buildArgs).map(([k, v]) => {
          return {
            key: k,
            value: v
          }
        }),
        environmentVariables: environmentVariableDetails.keys.map((key) => {
          return {
            key: environmentVariableDetails.map[key].name,
            value: environmentVariableDetails.map[key].value
          }
        }),
        configMounts: configMountDetails.keys.map((key) => configMountDetails.map[key]),
        persistentVolumeBindings: persistentVolumeBindingsDetails.keys.map((key) => {
          return {
            persistentVolumeID: persistentVolumeBindingsDetails.map[key].persistentVolumeID,
            mountingPath: persistentVolumeBindingsDetails.map[key].mountingPath
          }
        }),
        // update this part
        gitCredentialID: sourceConfigurationRef.gitCredentialID,
        repositoryUrl: sourceConfigurationRef.repositoryUrl,
        repositoryBranch: sourceConfigurationRef.repositoryBranch,
        codePath: sourceConfigurationRef.codePath,
        imageRegistryCredentialID: sourceConfigurationRef.imageRegistryCredentialID,
        dockerImage: sourceConfigurationRef.dockerImage,
        sourceCodeCompressedFileName: sourceConfigurationRef.sourceCodeCompressedFileName,
        dockerfile: sourceConfigurationRef.dockerfile,
        capabilities: applicationExistingDetails.capabilities,
        sysctls: applicationExistingDetails.sysctls,
        applicationGroupID: applicationExistingDetails.applicationGroupID,
        customHealthCheck: {
          enabled: deploymentConfigurationDetails.customHealthCheck.enabled,
          test_command: deploymentConfigurationDetails.customHealthCheck.test_command,
          interval_seconds: deploymentConfigurationDetails.customHealthCheck.interval_seconds,
          timeout_seconds: deploymentConfigurationDetails.customHealthCheck.timeout_seconds,
          start_period_seconds: deploymentConfigurationDetails.customHealthCheck.start_period_seconds,
          start_interval_seconds: deploymentConfigurationDetails.customHealthCheck.start_interval_seconds,
          retries: deploymentConfigurationDetails.customHealthCheck.retries
        },
        preferredServerHostnames: deploymentConfigurationDetails.preferredServerHostnames,
        dockerProxyConfig: {
          enabled: deploymentConfigurationDetails.dockerProxyConfig.enabled,
          permission: {
            ping: deploymentConfigurationDetails.dockerProxyConfig.permission.ping,
            version: deploymentConfigurationDetails.dockerProxyConfig.permission.version,
            info: deploymentConfigurationDetails.dockerProxyConfig.permission.info,
            events: deploymentConfigurationDetails.dockerProxyConfig.permission.events,
            auth: deploymentConfigurationDetails.dockerProxyConfig.permission.auth,
            secrets: deploymentConfigurationDetails.dockerProxyConfig.permission.secrets,
            build: deploymentConfigurationDetails.dockerProxyConfig.permission.build,
            commit: deploymentConfigurationDetails.dockerProxyConfig.permission.commit,
            configs: deploymentConfigurationDetails.dockerProxyConfig.permission.configs,
            containers: deploymentConfigurationDetails.dockerProxyConfig.permission.containers,
            distribution: deploymentConfigurationDetails.dockerProxyConfig.permission.distribution,
            exec: deploymentConfigurationDetails.dockerProxyConfig.permission.exec,
            grpc: deploymentConfigurationDetails.dockerProxyConfig.permission.grpc,
            images: deploymentConfigurationDetails.dockerProxyConfig.permission.images,
            networks: deploymentConfigurationDetails.dockerProxyConfig.permission.networks,
            nodes: deploymentConfigurationDetails.dockerProxyConfig.permission.nodes,
            plugins: deploymentConfigurationDetails.dockerProxyConfig.permission.plugins,
            services: deploymentConfigurationDetails.dockerProxyConfig.permission.services,
            session: deploymentConfigurationDetails.dockerProxyConfig.permission.session,
            swarm: deploymentConfigurationDetails.dockerProxyConfig.permission.swarm,
            system: deploymentConfigurationDetails.dockerProxyConfig.permission.system,
            tasks: deploymentConfigurationDetails.dockerProxyConfig.permission.tasks,
            volumes: deploymentConfigurationDetails.dockerProxyConfig.permission.volumes
          }
        }
      }
    }

    const gitRepoURL = computed(() => {
      const applicationExistingDetails = applicationExistingDetailsResult.value?.application ?? {}
      if (applicationExistingDetails.length === 0) {
        return ''
      }
      return applicationExistingDetails.latestDeployment.gitEndpoint
    })

    const updateApplicationSource = (source) => {
      sourceConfigurationRef.command = source.command
      sourceConfigurationRef.gitCredentialID = source.gitCredentialID
      sourceConfigurationRef.repositoryUrl = source.gitRepoUrl
      sourceConfigurationRef.repositoryBranch = source.gitBranch
      sourceConfigurationRef.codePath = source.codePath
      sourceConfigurationRef.imageRegistryCredentialID = source.imageRegistryCredentialID
      sourceConfigurationRef.dockerImage = source.dockerImage
      sourceConfigurationRef.sourceCodeCompressedFileName = source.sourceCodeCompressedFileName
      sourceConfigurationRef.dockerfile = source.dockerFile
      sourceConfigurationRef.buildArgs = source.buildArgs
      triggerUpdateHook()
    }

    return {
      isConfigurationUpdated,
      applyConfigurationChanges,
      cancelConfigurationChanges,
      applicationDetailsLoading,
      environmentVariableDetails,
      addEnvironmentVariable,
      deleteEnvironmentVariable,
      onEnvironmentVariableNameChange,
      onEnvironmentVariableValueChange,
      configMountDetails,
      addConfigMount,
      onConfigMountContentChange,
      deleteConfigMount,
      persistentVolumeBindingsDetails,
      addPersistentVolumeBinding,
      deletePersistentVolumeBinding,
      onPersistentVolumeChange,
      onPersistentVolumeMountingPathChange,
      onMemoryLimitChanged,
      onMemoryReservedChanged,
      deploymentConfigurationDetails,
      updatePreferredServerHostnames,
      preferredServerHostnamesStr,
      enableDockerProxy,
      disableDockerProxy,
      dockerProxyPermissionChanged,
      changeDeploymentStrategy,
      replicasCountChanged,
      isDeployRequestSubmitting,
      gitRepoURL,
      applicationExistingDetailsResult,
      updateApplicationSource,
      triggerUpdateHook,
      dockerProxyHost
    }
  })
}
