package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-set"
	containermanger "github.com/swiftwave-org/swiftwave/container_manager"
	"gorm.io/gorm"
)

// This file contains the operations for the Application model.
// This functions will perform necessary validation before doing the actual database operation.

// Each function's argument format should be (ctx context.Context, db gorm.DB, ...)
// context used to pass some data to the function e.g. user id, auth info, etc.

func IsExistApplicationName(_ context.Context, db gorm.DB, dockerManager containermanger.Manager, name string) (bool, error) {
	// name cannot contain any special characters at the end
	if strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, ":") || strings.Contains(name, "*") || strings.Contains(name, "?") || strings.Contains(name, "\"") || strings.Contains(name, "<") || strings.Contains(name, ">") || strings.Contains(name, "|") || strings.Contains(name, "&") || strings.HasSuffix(name, "_") {
		return false, errors.New("application name cannot contain any special characters at the end")
	}
	// verify from database
	var count int64
	tx := db.Model(&Application{}).Where("name = ?", name).Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}
	if count > 0 {
		return true, nil
	}
	// verify from docker client
	_, err := dockerManager.GetService(name)
	if err == nil {
		return true, nil
	}
	return false, nil
}

func FindAllApplications(_ context.Context, db gorm.DB, includeGroupedApplications bool) ([]*Application, error) {
	var applications []*Application
	var tx *gorm.DB
	if includeGroupedApplications {
		tx = db.Find(&applications)
	} else {
		tx = db.Where("application_group_id IS NULL").Find(&applications)
	}
	return applications, tx.Error
}

type ApplicationDeploymentInfo struct {
	ApplicationID string
	DeploymentID  string
}

func FindApplicationsForForceUpdate(_ context.Context, db gorm.DB) ([]*ApplicationDeploymentInfo, error) {
	var deployments []*Deployment
	err := db.Model(&Deployment{}).Where("status = ?", DeploymentStatusDeployed).Scan(&deployments).Error
	if err != nil {
		return nil, err
	}
	var applicationDeploymentInfos []*ApplicationDeploymentInfo
	for _, deployment := range deployments {
		applicationDeploymentInfos = append(applicationDeploymentInfos, &ApplicationDeploymentInfo{
			ApplicationID: deployment.ApplicationID,
			DeploymentID:  deployment.ID,
		})
	}
	return applicationDeploymentInfos, nil
}

func (application *Application) FindById(_ context.Context, db gorm.DB, id string) error {
	tx := db.Where("id = ?", id).First(&application)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (application *Application) FindByName(_ context.Context, db gorm.DB, name string) error {
	tx := db.Where("name = ?", name).First(&application)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (application *Application) Create(ctx context.Context, db gorm.DB, dockerManager containermanger.Manager, codeTarballDir string) error {
	// verify if there is no application with same name
	isExist, err := IsExistApplicationName(ctx, db, dockerManager, application.Name)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("application name not available")
	}
	// check resource limits and reserved resource
	if application.ResourceLimit.MemoryMB != 0 && application.ResourceLimit.MemoryMB < 6 {
		return errors.New("memory limit should be at least 6 MB or 0 for unlimited")
	}
	if application.ReservedResource.MemoryMB != 0 && application.ReservedResource.MemoryMB < 6 {
		return errors.New("reserved memory should be at least 6 MB or 0 for unlimited")
	}
	// Verify the PreferredServerHostnames
	if len(application.PreferredServerHostnames) > 0 {
		for _, preferredServerHostname := range application.PreferredServerHostnames {
			_, err := FetchServerIDByHostName(&db, preferredServerHostname)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return fmt.Errorf("invalid hostname %s provided for preferred server", preferredServerHostname)
				}
				return err
			}
		}
	}
	// State
	isGitCredentialExist := false
	isImageRegistryCredentialExist := false
	// For UpstreamType = Git, verify git record id
	if application.LatestDeployment.UpstreamType == UpstreamTypeGit {
		if application.LatestDeployment.GitCredentialID != nil {
			var gitCredential = &GitCredential{}
			err := gitCredential.FindById(ctx, db, *application.LatestDeployment.GitCredentialID)
			if err != nil {
				return err
			}
			isGitCredentialExist = true
		} else {
			isGitCredentialExist = false
		}
	}
	// For UpstreamType = Image, verify image registry credential id
	if application.LatestDeployment.UpstreamType == UpstreamTypeImage {
		if application.LatestDeployment.ImageRegistryCredentialID != nil {
			var imageRegistryCredential = &ImageRegistryCredential{}
			err := imageRegistryCredential.FindById(ctx, db, *application.LatestDeployment.ImageRegistryCredentialID)
			if err != nil {
				return err
			}
			isImageRegistryCredentialExist = true
		} else {
			isImageRegistryCredentialExist = false
		}
	}
	// For UpstreamType = SourceCode, verify source code compressed file exists
	if application.LatestDeployment.UpstreamType == UpstreamTypeSourceCode {
		tarballPath := filepath.Join(codeTarballDir, application.LatestDeployment.SourceCodeCompressedFileName)
		// Verify file exists
		if _, err := os.Stat(tarballPath); os.IsNotExist(err) {
			return errors.New("source code not found")
		}
	}
	// Validate DockerProxy configuration
	if application.DockerProxy.Enabled && len(application.PreferredServerHostnames) == 0 {
		return errors.New("you need to select exactly one preferred server for getting access to docker socket proxy")
	}
	// create application
	createdApplication := Application{
		ID:                       uuid.NewString(),
		Name:                     application.Name,
		DeploymentMode:           application.DeploymentMode,
		Replicas:                 application.Replicas,
		WebhookToken:             uuid.NewString(),
		Command:                  application.Command,
		Capabilities:             application.Capabilities,
		Sysctls:                  application.Sysctls,
		ResourceLimit:            application.ResourceLimit,
		ReservedResource:         application.ReservedResource,
		ApplicationGroupID:       application.ApplicationGroupID,
		DockerProxy:              application.DockerProxy,
		PreferredServerHostnames: application.PreferredServerHostnames,
		CustomHealthCheck:        application.CustomHealthCheck,
	}
	tx := db.Create(&createdApplication)
	if tx.Error != nil {
		return tx.Error
	}
	// create environment variables
	createdEnvironmentVariables := make([]EnvironmentVariable, 0)
	for _, environmentVariable := range application.EnvironmentVariables {
		createdEnvironmentVariable := EnvironmentVariable{
			ApplicationID: createdApplication.ID,
			Key:           environmentVariable.Key,
			Value:         environmentVariable.Value,
		}
		createdEnvironmentVariables = append(createdEnvironmentVariables, createdEnvironmentVariable)
	}
	if len(createdEnvironmentVariables) > 0 {
		tx = db.Create(&createdEnvironmentVariables)
		if tx.Error != nil {
			return tx.Error
		}
	}
	// create persistent volume bindings
	createdPersistentVolumeBindings := make([]PersistentVolumeBinding, 0)
	persistedVolumeBindingsMountingPathSet := set.From[string](make([]string, 0))
	for _, persistentVolumeBinding := range application.PersistentVolumeBindings {
		// check if mounting path is already used
		if persistedVolumeBindingsMountingPathSet.Contains(persistentVolumeBinding.MountingPath) {
			return errors.New("mounting path already used")
		} else {
			persistedVolumeBindingsMountingPathSet.Insert(persistentVolumeBinding.MountingPath)
		}
		// verify persistent volume exists
		var persistentVolume = &PersistentVolume{}
		err := persistentVolume.FindById(ctx, db, persistentVolumeBinding.PersistentVolumeID)
		if err != nil {
			return err
		}
		createdPersistentVolumeBinding := PersistentVolumeBinding{
			ApplicationID:      createdApplication.ID,
			PersistentVolumeID: persistentVolumeBinding.PersistentVolumeID,
			MountingPath:       persistentVolumeBinding.MountingPath,
		}
		createdPersistentVolumeBindings = append(createdPersistentVolumeBindings, createdPersistentVolumeBinding)
	}
	if len(createdPersistentVolumeBindings) > 0 {
		tx = db.Create(&createdPersistentVolumeBindings)
		if tx.Error != nil {
			return tx.Error
		}
	}
	// create config records
	configMountRecords := make([]ConfigMount, 0)
	configMountRecordsMountingPathSet := set.From[string](make([]string, 0))
	for _, configMount := range application.ConfigMounts {
		// check if mounting path is already used
		if configMountRecordsMountingPathSet.Contains(configMount.MountingPath) {
			return errors.New("mounting path already used")
		} else {
			configMountRecordsMountingPathSet.Insert(configMount.MountingPath)
		}
		configMount.ApplicationID = createdApplication.ID
		configMountRecords = append(configMountRecords, configMount)
	}
	if len(configMountRecords) > 0 {
		tx = db.Create(&configMountRecords)
		if tx.Error != nil {
			return tx.Error
		}
	}
	// handle other stuffs
	var gitCredentialID *uint = nil
	if isGitCredentialExist {
		gitCredentialID = application.LatestDeployment.GitCredentialID
	}
	var imageRegistryCredentialID *uint = nil
	if isImageRegistryCredentialExist {
		imageRegistryCredentialID = application.LatestDeployment.ImageRegistryCredentialID
	}
	// create deployment
	createdDeployment := Deployment{
		ApplicationID: createdApplication.ID,
		UpstreamType:  application.LatestDeployment.UpstreamType,
		// Fields for UpstreamType = Git
		GitCredentialID:  gitCredentialID,
		GitType:          application.LatestDeployment.GitType,
		GitProvider:      application.LatestDeployment.GitProvider,
		GitEndpoint:      application.LatestDeployment.GitEndpoint,
		GitSshUser:       application.LatestDeployment.GitSshUser,
		RepositoryOwner:  application.LatestDeployment.RepositoryOwner,
		RepositoryName:   application.LatestDeployment.RepositoryName,
		RepositoryBranch: application.LatestDeployment.RepositoryBranch,
		CommitHash:       application.LatestDeployment.CommitHash,
		CodePath:         application.LatestDeployment.CodePath,
		// Fields for UpstreamType = SourceCode
		SourceCodeCompressedFileName: application.LatestDeployment.SourceCodeCompressedFileName,
		// Fields for UpstreamType = Image
		DockerImage:               application.LatestDeployment.DockerImage,
		ImageRegistryCredentialID: imageRegistryCredentialID,
		// other fields
		Dockerfile: application.LatestDeployment.Dockerfile,
	}
	err = createdDeployment.Create(ctx, db)
	if err != nil {
		return err
	}
	// add build args to deployment
	createdBuildArgs := make([]BuildArg, 0)
	for _, buildArg := range application.LatestDeployment.BuildArgs {
		createdBuildArg := BuildArg{
			DeploymentID: createdDeployment.ID,
			Key:          buildArg.Key,
			Value:        buildArg.Value,
		}
		createdBuildArgs = append(createdBuildArgs, createdBuildArg)
	}
	if len(createdBuildArgs) > 0 {
		tx = db.Create(&createdBuildArgs)
		if tx.Error != nil {
			return tx.Error
		}
	}
	// update application details
	*application = createdApplication
	return nil
}

func (application *Application) Update(ctx context.Context, db gorm.DB, _ containermanger.Manager) (*ApplicationUpdateResult, error) {
	var err error
	// ensure that application is not deleted
	isDeleted, err := application.IsApplicationDeleted(ctx, db)
	if err != nil {
		return nil, err
	}
	if isDeleted {
		return nil, errors.New("application is deleted")
	}
	// check resource limits and reserved resource
	if application.ResourceLimit.MemoryMB != 0 && application.ResourceLimit.MemoryMB < 6 {
		return nil, errors.New("memory limit should be at least 6 MB or 0 for unlimited")
	}
	if application.ReservedResource.MemoryMB != 0 && application.ReservedResource.MemoryMB < 6 {
		return nil, errors.New("reserved memory should be at least 6 MB or 0 for unlimited")
	}
	// Verify the PreferredServerHostnames
	if len(application.PreferredServerHostnames) > 0 {
		for _, preferredServerHostname := range application.PreferredServerHostnames {
			_, err := FetchServerIDByHostName(&db, preferredServerHostname)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, fmt.Errorf("invalid hostname %s provided for preferred server", preferredServerHostname)
				}
				return nil, err
			}
		}
	}
	// check if docker proxy is enabled and preferred servers are not provided
	if application.DockerProxy.Enabled && len(application.PreferredServerHostnames) != 1 {
		return nil, errors.New("you must select preferred servers for deployment to get access to docker proxy")
	}
	// status
	isReloadRequired := false
	// fetch application with environment variables and persistent volume bindings
	var applicationExistingFull = &Application{}
	tx := db.Preload("EnvironmentVariables").Preload("ConfigMounts").Preload("PersistentVolumeBindings").Where("id = ?", application.ID).First(&applicationExistingFull)
	if tx.Error != nil {
		return nil, tx.Error
	}
	// check if DeploymentMode is changed
	if applicationExistingFull.DeploymentMode != application.DeploymentMode {
		// update deployment mode
		err = db.Model(&applicationExistingFull).Update("deployment_mode", application.DeploymentMode).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// check if Command is changed
	if applicationExistingFull.Command != application.Command {
		// update command
		err = db.Model(&applicationExistingFull).Update("command", application.Command).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// if replicated deployment, check if Replicas is changed
	if application.DeploymentMode == DeploymentModeReplicated && applicationExistingFull.Replicas != application.Replicas {
		// update replicas
		err = db.Model(&applicationExistingFull).Update("replicas", application.Replicas).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// check if Resource limits or reservations are changed
	if applicationExistingFull.ResourceLimit.MemoryMB != application.ResourceLimit.MemoryMB || applicationExistingFull.ReservedResource.MemoryMB != application.ReservedResource.MemoryMB {
		// update resource limits
		err = db.Model(&applicationExistingFull).Update("resource_limit_memory_mb", application.ResourceLimit.MemoryMB).Error
		if err != nil {
			return nil, err
		}
		// update reserved resource
		err = db.Model(&applicationExistingFull).Update("reserved_resource_memory_mb", application.ReservedResource.MemoryMB).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// create array of environment variables
	var newEnvironmentVariableMap = make(map[string]string)
	for _, environmentVariable := range application.EnvironmentVariables {
		newEnvironmentVariableMap[environmentVariable.Key] = environmentVariable.Value
	}
	// update environment variables -- if required
	if applicationExistingFull.EnvironmentVariables != nil {
		for _, environmentVariable := range applicationExistingFull.EnvironmentVariables {
			// check if environment variable is present in new environment variables
			if _, ok := newEnvironmentVariableMap[environmentVariable.Key]; ok {
				// check if value is changed
				if environmentVariable.Value != newEnvironmentVariableMap[environmentVariable.Key] {
					// update environment variable
					environmentVariable.Value = newEnvironmentVariableMap[environmentVariable.Key]
					err = environmentVariable.Update(ctx, db)
					if err != nil {
						return nil, err
					}
					// delete from newEnvironmentVariableMap
					delete(newEnvironmentVariableMap, environmentVariable.Key)
					// reload application
					isReloadRequired = true
				} else {
					// delete from newEnvironmentVariableMap
					delete(newEnvironmentVariableMap, environmentVariable.Key)
				}
			} else {
				// delete environment variable
				err = environmentVariable.Delete(ctx, db)
				if err != nil {
					return nil, err
				}
				// reload application
				isReloadRequired = true
			}
		}
	}
	// add new environment variables which are not present
	for key, value := range newEnvironmentVariableMap {
		environmentVariable := EnvironmentVariable{
			ApplicationID: application.ID,
			Key:           key,
			Value:         value,
		}
		err := environmentVariable.Create(ctx, db)
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// create map of config mounts
	var newConfigMountMap = make(map[string]ConfigMount)
	for _, configMount := range application.ConfigMounts {
		newConfigMountMap[configMount.MountingPath] = configMount
	}
	if applicationExistingFull.ConfigMounts != nil {
		for _, configMount := range applicationExistingFull.ConfigMounts {
			// check if config mount is present in new config mounts
			if _, ok := newConfigMountMap[configMount.MountingPath]; ok {
				// check if anything is changed
				if strings.Compare(configMount.Content, newConfigMountMap[configMount.MountingPath].Content) != 0 ||
					configMount.Uid != newConfigMountMap[configMount.MountingPath].Uid ||
					configMount.Gid != newConfigMountMap[configMount.MountingPath].Gid {
					// update config mount
					configMount.Content = newConfigMountMap[configMount.MountingPath].Content
					configMount.Uid = newConfigMountMap[configMount.MountingPath].Uid
					configMount.Gid = newConfigMountMap[configMount.MountingPath].Gid
					err = configMount.Update(ctx, db)
					if err != nil {
						return nil, err
					}
					// delete from newConfigMountMap
					delete(newConfigMountMap, configMount.MountingPath)
					// reload application
					isReloadRequired = true
				} else {
					// delete from newConfigMountMap
					delete(newConfigMountMap, configMount.MountingPath)
				}
			} else {
				err = configMount.Delete(ctx, db)
				if err != nil {
					return nil, err
				}
				// reload application
				isReloadRequired = true
			}
		}
	}
	// add new config mounts which are not present
	for mountingPath, record := range newConfigMountMap {
		configMount := ConfigMount{
			ApplicationID: application.ID,
			ConfigID:      "",
			Content:       record.Content,
			MountingPath:  mountingPath,
			Uid:           record.Uid,
			Gid:           record.Gid,
			FileMode:      444,
		}
		err := configMount.Create(ctx, db)
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// create array of persistent volume bindings
	var newPersistentVolumeBindingMap = make(map[string]uint)
	newPersistentVolumeBindingMountingPathSet := set.From[string](make([]string, 0))
	for _, persistentVolumeBinding := range application.PersistentVolumeBindings {
		// check if mounting path is already used
		if newPersistentVolumeBindingMountingPathSet.Contains(persistentVolumeBinding.MountingPath) {
			return nil, errors.New("duplicate mounting path found")
		} else {
			newPersistentVolumeBindingMountingPathSet.Insert(persistentVolumeBinding.MountingPath)
		}
		newPersistentVolumeBindingMap[persistentVolumeBinding.MountingPath] = persistentVolumeBinding.PersistentVolumeID
	}
	// update persistent volume bindings -- if required
	if applicationExistingFull.PersistentVolumeBindings != nil {
		for _, persistentVolumeBinding := range applicationExistingFull.PersistentVolumeBindings {
			// check if persistent volume binding is present in new persistent volume bindings
			if _, ok := newPersistentVolumeBindingMap[persistentVolumeBinding.MountingPath]; ok {
				// check if value is changed
				if persistentVolumeBinding.PersistentVolumeID != newPersistentVolumeBindingMap[persistentVolumeBinding.MountingPath] {
					// update persistent volume binding
					persistentVolumeBinding.PersistentVolumeID = newPersistentVolumeBindingMap[persistentVolumeBinding.MountingPath]
					err = persistentVolumeBinding.Update(ctx, db)
					if err != nil {
						return nil, err
					}
					// delete from newPersistentVolumeBindingMap
					delete(newPersistentVolumeBindingMap, persistentVolumeBinding.MountingPath)
					// reload application
					isReloadRequired = true
				} else {
					// delete from newPersistentVolumeBindingMap
					delete(newPersistentVolumeBindingMap, persistentVolumeBinding.MountingPath)
				}
			} else {
				// delete persistent volume binding
				err = persistentVolumeBinding.Delete(ctx, db)
				if err != nil {
					return nil, err
				}
				// reload application
				isReloadRequired = true
			}
		}
	}
	// add new persistent volume bindings which are not present
	for mountingPath, persistentVolumeID := range newPersistentVolumeBindingMap {
		persistentVolumeBinding := PersistentVolumeBinding{
			ApplicationID:      application.ID,
			PersistentVolumeID: persistentVolumeID,
			MountingPath:       mountingPath,
		}
		err := persistentVolumeBinding.Create(ctx, db)
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	isPreferredServerHostnameChanged := false
	// check if preferred servers a changed
	if len(application.PreferredServerHostnames) != len(applicationExistingFull.PreferredServerHostnames) {
		isPreferredServerHostnameChanged = true
	} else {
		// check if elements are changed
		for _, preferredServerHostname := range application.PreferredServerHostnames {
			isFound := false
			for _, preferredServerHostnameExisting := range applicationExistingFull.PreferredServerHostnames {
				if strings.Compare(preferredServerHostname, preferredServerHostnameExisting) == 0 {
					isFound = true
					break
				}
			}
			if !isFound {
				isPreferredServerHostnameChanged = true
			}
		}
	}
	if isPreferredServerHostnameChanged {
		// update preferred server hostnames
		err = db.Model(&applicationExistingFull).Update("preferred_server_hostnames", application.PreferredServerHostnames).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// check for changes in docker proxy configuration
	if !application.DockerProxy.Equal(&applicationExistingFull.DockerProxy) {
		// store docker proxy configuration
		err = db.Model(&applicationExistingFull).Select("docker_proxy_enabled",
			"docker_proxy_permission_ping", "docker_proxy_permission_version",
			"docker_proxy_permission_info", "docker_proxy_permission_events", "docker_proxy_permission_auth",
			"docker_proxy_permission_secrets", "docker_proxy_permission_build", "docker_proxy_permission_commit",
			"docker_proxy_permission_configs", "docker_proxy_permission_containers", "docker_proxy_permission_distribution",
			"docker_proxy_permission_exec", "docker_proxy_permission_grpc", "docker_proxy_permission_images",
			"docker_proxy_permission_networks", "docker_proxy_permission_nodes", "docker_proxy_permission_plugins",
			"docker_proxy_permission_services", "docker_proxy_permission_session", "docker_proxy_permission_swarm",
			"docker_proxy_permission_system", "docker_proxy_permission_tasks", "docker_proxy_permission_volumes").Updates(application).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// check for changes in custom health check
	if !application.CustomHealthCheck.Equal(&applicationExistingFull.CustomHealthCheck) {
		// store custom health check configuration
		err = db.Model(&applicationExistingFull).Select("custom_health_check_enabled",
			"custom_health_check_test_command", "custom_health_check_interval_seconds",
			"custom_health_check_timeout_seconds", "custom_health_check_start_period_seconds",
			"custom_health_check_start_interval_seconds", "custom_health_check_retries").Updates(application).Error
		if err != nil {
			return nil, err
		}
		// reload application
		isReloadRequired = true
	}
	// update deployment -- if required
	currentDeploymentID, err := FindCurrentDeployedDeploymentIDByApplicationId(ctx, db, application.ID)
	if err != nil {
		currentDeploymentID, err = FindLatestDeploymentIDByApplicationId(ctx, db, application.ID)
	}
	if err != nil {
		return nil, err
	}
	// set deployment id
	application.LatestDeployment.ID = currentDeploymentID
	// send call to update deployment
	updateDeploymentStatus, err := application.LatestDeployment.Update(ctx, db)
	if err != nil {
		return nil, err
	}
	return &ApplicationUpdateResult{
		ReloadRequired:  isReloadRequired,
		RebuildRequired: updateDeploymentStatus.RebuildRequired,
		DeploymentId:    updateDeploymentStatus.DeploymentId,
	}, nil
}

func (application *Application) SoftDelete(ctx context.Context, db gorm.DB, _ containermanger.Manager) error {
	// ensure that application is not deleted
	isDeleted, err := application.IsApplicationDeleted(ctx, db)
	if err != nil {
		return err
	}
	if isDeleted {
		return errors.New("application is deleted")
	}
	// ensure there is no ingress rule associated with this application
	ingressRules, err := FindIngressRulesByApplicationID(ctx, db, application.ID)
	if err != nil {
		return err
	}
	if len(ingressRules) > 0 {
		return errors.New("application has ingress rules associated with it")
	}
	// do soft delete
	tx := db.Model(&application).Update("is_deleted", true)
	return tx.Error
}

func (application *Application) HardDelete(ctx context.Context, db gorm.DB, _ containermanger.Manager) error {
	// ensure there is no ingress rule associated with this application
	ingressRules, err := FindIngressRulesByApplicationID(ctx, db, application.ID)
	if err != nil {
		return err
	}
	if len(ingressRules) > 0 {
		return errors.New("application has ingress rules associated with it")
	}
	// fetch application group id
	var applicationGroupID *string = nil
	err = db.Model(&application).Select("application_group_id").Where("id = ?", application.ID).Scan(&applicationGroupID).Error
	if err != nil {
		return err
	}
	// delete application
	err = db.Delete(&application).Error
	if err != nil {
		return err
	}
	if applicationGroupID != nil {
		// if application group is not associated with any other application, delete the group
		applicationGroup := &ApplicationGroup{
			ID: *applicationGroupID,
		}
		isAnyApplicationAssociatedWithGroup, err := applicationGroup.IsAnyApplicationAssociatedWithGroup(ctx, db)
		if err != nil {
			return err
		}
		if !isAnyApplicationAssociatedWithGroup {
			err = applicationGroup.Delete(ctx, db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (application *Application) IsApplicationDeleted(_ context.Context, db gorm.DB) (bool, error) {
	// verify from database
	var count int64
	tx := db.Model(&Application{}).Where("id = ? AND is_deleted = ?", application.ID, true).Count(&count)
	if tx.Error != nil {
		return false, tx.Error
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

func (application *Application) RebuildApplication(ctx context.Context, db gorm.DB) (deploymentId string, error error) {
	// fetch record
	err := application.FindById(ctx, db, application.ID)
	if err != nil {
		return "", err
	}
	// create a new deployment from latest deployment
	latestDeployment, err := FindCurrentDeployedDeploymentByApplicationId(ctx, db, application.ID)
	if err != nil {
		latestDeployment, err = FindLatestDeploymentByApplicationId(ctx, db, application.ID)
		if err != nil {
			return "", errors.New("failed to fetch latest deployment")
		}
	}

	// fetch build args
	buildArgs, err := FindBuildArgsByDeploymentId(ctx, db, latestDeployment.ID)
	if err != nil {
		return "", err
	}
	// add new deployment
	err = latestDeployment.Create(ctx, db)
	if err != nil {
		return "", err
	}
	// update build args
	for _, buildArg := range buildArgs {
		buildArg.ID = 0
		buildArg.DeploymentID = latestDeployment.ID
	}
	if len(buildArgs) > 0 {
		err = db.Create(&buildArgs).Error
		if err != nil {
			return "", err
		}
	}
	return latestDeployment.ID, nil
}

func (application *Application) RegenerateWebhookToken(ctx context.Context, db gorm.DB) error {
	// fetch record
	err := application.FindById(ctx, db, application.ID)
	if err != nil {
		return err
	}
	// update webhook token
	application.WebhookToken = uuid.NewString()
	tx := db.Model(&application).Update("webhook_token", application.WebhookToken)
	return tx.Error
}

func (application *Application) MarkAsSleeping(ctx context.Context, db gorm.DB) error {
	// fetch record
	err := application.FindById(ctx, db, application.ID)
	if err != nil {
		return err
	}
	if application.DeploymentMode == DeploymentModeGlobal {
		return errors.New("global deployment cannot be marked as sleeping")
	}
	// update is sleeping
	tx := db.Model(&application).Update("is_sleeping", true)
	return tx.Error
}

func (application *Application) MarkAsWake(ctx context.Context, db gorm.DB) error {
	// fetch record
	err := application.FindById(ctx, db, application.ID)
	if err != nil {
		return err
	}
	if application.DeploymentMode == DeploymentModeGlobal {
		return errors.New("global deployment cannot be marked as wake")
	}
	// update is sleeping
	tx := db.Model(&application).Update("is_sleeping", false)
	return tx.Error
}

func (application *Application) UpdateGroup(ctx context.Context, db gorm.DB, groupId *string) error {
	err := application.FindById(ctx, db, application.ID)
	if err != nil {
		return err
	}
	oldApplicationGroupID := ""
	if application.ApplicationGroupID != nil {
		oldApplicationGroupID = *application.ApplicationGroupID
	}
	if groupId != nil && strings.Compare(*groupId, "") == 0 {
		groupId = nil
	}
	err = db.Model(&application).Update("application_group_id", groupId).Error
	if err != nil {
		return err
	}
	if groupId == nil && strings.Compare(oldApplicationGroupID, "") != 0 {
		group := &ApplicationGroup{
			ID: oldApplicationGroupID,
		}
		isAnyApplicationAssociatedWithGroup, err := group.IsAnyApplicationAssociatedWithGroup(ctx, db)
		if err != nil {
			return err
		}
		if !isAnyApplicationAssociatedWithGroup {
			err = group.Delete(ctx, db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
