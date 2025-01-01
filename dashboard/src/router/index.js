import { createRouter, createWebHistory } from 'vue-router'
import UnderMaintenancePage from '@/views/pages/UnderMaintenance.vue'
import SetupSwiftwavePage from '@/views/pages/SetupSwiftwave.vue'
import LoginView from '@/views/pages/LoginView.vue'
import DeployApplicationPage from '@/views/pages/DeployApplication.vue'
import DeployStackPage from '@/views/pages/DeployStack.vue'
import AppStorePage from '@/views/pages/AppStore.vue'
import AppInstallPage from '@/views/pages/AppInstall.vue'
import ApplicationsPage from '@/views/pages/ApplicationManagement.vue'
import ApplicationDetailsPage from '@/views/pages/ApplicationDetails.vue'
import ApplicationDetailsDeploymentListPage from '@/views/pages/ApplicationDetails/DeploymentList.vue'
import ApplicationDetailsDeploymentDetailsPage from '@/views/pages/ApplicationDetails/DeploymentDetails.vue'
import ApplicationDetailsRuntimeLogsPage from '@/views/pages/ApplicationDetails/RuntimeLogs.vue'
import ApplicationDetailsIngressRulesPage from '@/views/pages/ApplicationDetails/IngressRules.vue'
import ApplicationDetailsUpdateSourcePage from '@/views/pages/ApplicationDetails/UpdateSourceCode.vue'
import ApplicationDetailsEnvironmentVariablesPage from '@/views/pages/ApplicationDetails/EnvironmentVariables.vue'
import ApplicationDetailsConfigMountsPage from '@/views/pages/ApplicationDetails/ConfigMounts.vue'
import ApplicationDetailsPersistentVolumesPage from '@/views/pages/ApplicationDetails/PersistentVolumes.vue'
import ApplicationDetailsDeploymentConfigPage from '@/views/pages/ApplicationDetails/DeploymentConfig.vue'
import ApplicationDetailsDangerZonePage from '@/views/pages/ApplicationDetails/Destroy.vue'
import ApplicationDetailsManagePage from '@/views/pages/ApplicationDetails/Manage.vue'
import ApplicationDetailsWebhookCIPage from '@/views/pages/ApplicationDetails/WebhookCI.vue'
import ApplicationDetailsResourceStatsPage from '@/views/pages/ApplicationDetails/ResourceStats.vue'
import PersistentVolumeManagementPage from '@/views/pages/PersistentVolumeManagement.vue'
import UserManagementPage from '@/views/pages/UsersManagement.vue'
import GitCredentialManagementPage from '@/views/pages/GitCredentialManagement.vue'
import ImageRegistryCredentialManagementPage from '@/views/pages/ImageRegistryCredentialManagement.vue'
import DomainManagementPage from '@/views/pages/DomainManagement.vue'
import RedirectRuleManagementPage from '@/views/pages/RedirectRuleManagement.vue'
import IngressRuleManagementPage from '@/views/pages/IngressRuleManagement.vue'
import DownloadPVBackupPage from '@/views/pages/DownloadPVBackup.vue'
import ServerManagementPage from '@/views/pages/ServerManagement.vue'
import ServerLogsPage from '@/views/pages/ServerLogs.vue'
import ServerAnalyticsPage from '@/views/pages/ServerAnalytics.vue'
import SystemLogsPage from '@/views/pages/SystemLogs.vue'
import AppAuthBasicAccessControlList from '@/views/pages/AppAuthBasicAccessControlList.vue'
import ApplicationGroupDetailsPage from '@/views/pages/ApplicationGroupDetails.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/setup',
      name: 'Setup',
      component: SetupSwiftwavePage
    },
    {
      path: '/maintenance',
      name: 'Maintenance',
      component: UnderMaintenancePage
    },
    {
      path: '/login',
      name: 'Login',
      component: LoginView
    },
    {
      path: '',
      redirect: '/applications'
    },

    {
      path: '/deploy',
      children: [
        {
          path: 'application',
          name: 'Deploy Application',
          component: DeployApplicationPage
        },
        {
          path: 'stack',
          name: 'Deploy Stack',
          component: DeployStackPage
        },
        {
          path: 'app-store',
          name: 'App Store',
          component: AppStorePage
        },
        {
          path: 'app-store/install',
          name: 'Install from App Store',
          component: AppInstallPage
        }
      ]
    },
    {
      path: '/applications',
      name: 'Applications',
      component: ApplicationsPage
    },
    {
      path: '/application/:id',
      name: 'Application Details',
      component: ApplicationDetailsPage,
      children: [
        {
          path: 'deployments',
          name: 'Application Details Deployments',
          component: ApplicationDetailsDeploymentListPage
        },
        {
          path: 'deployment/:deployment_id',
          name: 'Application Deployment Details',
          component: ApplicationDetailsDeploymentDetailsPage
        },
        {
          path: 'runtime_logs',
          name: 'Application Details Runtime Logs',
          component: ApplicationDetailsRuntimeLogsPage
        },
        {
          path: 'ingress_rules',
          name: 'Application Details Ingress Rules',
          component: ApplicationDetailsIngressRulesPage
        },
        {
          path: 'update_source',
          name: 'Application Details Update Source',
          component: ApplicationDetailsUpdateSourcePage
        },
        {
          path: 'environment_variables',
          name: 'Application Details Environment Variables',
          component: ApplicationDetailsEnvironmentVariablesPage
        },
        {
          path: 'persistent_volumes',
          name: 'Application Details Persistent Volumes',
          component: ApplicationDetailsPersistentVolumesPage
        },
        {
          path: 'config_mounts',
          name: 'Application Details Config Mounts',
          component: ApplicationDetailsConfigMountsPage
        },
        {
          path: 'deployment_config',
          name: 'Application Details Deployment Config',
          component: ApplicationDetailsDeploymentConfigPage
        },
        {
          path: 'danger_zone',
          name: 'Application Details Danger Zone',
          component: ApplicationDetailsDangerZonePage
        },
        {
          path: 'manage',
          name: 'Application Details Manage',
          component: ApplicationDetailsManagePage
        },
        {
          path: 'webhook_ci',
          name: 'Application Details Webhook CI',
          component: ApplicationDetailsWebhookCIPage
        },
        {
          path: 'resource_stats',
          name: 'Application Details Resource Stats',
          component: ApplicationDetailsResourceStatsPage
        }
      ]
    },
    {
      path: '/application_group/:id',
      name: 'Application Group Details',
      component: ApplicationGroupDetailsPage
    },
    {
      path: '/app_auth',
      children: [
        {
          path: 'basic_authentication',
          name: 'Application Auth Basic ACL',
          component: AppAuthBasicAccessControlList
        }
      ]
    },
    {
      path: '/persistent-volumes',
      name: 'Persistent Volumes',
      component: PersistentVolumeManagementPage
    },
    {
      path: '/users',
      name: 'Users',
      component: UserManagementPage
    },
    {
      path: '/git-credentials',
      name: 'Git Credentials',
      component: GitCredentialManagementPage
    },
    {
      path: '/image-registry-credentials',
      name: 'Image Registry Credentials',
      component: ImageRegistryCredentialManagementPage
    },
    {
      path: '/domains',
      name: 'Domains',
      component: DomainManagementPage
    },
    {
      path: '/redirect-rules',
      name: 'Redirect Rules',
      component: RedirectRuleManagementPage
    },
    {
      path: '/ingress-rules',
      name: 'Ingress Rules',
      component: IngressRuleManagementPage
    },
    {
      path: '/pv-backup-download/:backup_id',
      name: 'Download Persistent Volume Backup',
      component: DownloadPVBackupPage
    },
    {
      path: '/servers',
      name: 'Servers',
      component: ServerManagementPage
    },
    {
      path: '/server/logs',
      name: 'Server Logs',
      component: ServerLogsPage
    },
    {
      path: '/server/analytics',
      name: 'Server Analytics',
      component: ServerAnalyticsPage
    },
    {
      path: '/logs',
      name: 'System Logs',
      component: SystemLogsPage
    }
  ]
})

export default router
