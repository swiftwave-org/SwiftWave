package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"errors"
	"github.com/swiftwave-org/swiftwave/system_config"

	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
)

// BackupPersistentVolume is the resolver for the backupPersistentVolume field.
func (r *mutationResolver) BackupPersistentVolume(ctx context.Context, input model.PersistentVolumeBackupInput) (*model.PersistentVolumeBackup, error) {
	record := persistentVolumeBackupInputToDatabaseObject(&input)
	// check if s3 enabled
	if record.Type == core.S3Backup && !r.ServiceConfig.PersistentVolumeBackupConfig.S3Config.Enabled {
		return nil, errors.New("s3 backup is not enabled. Please enable it in the swiftwave configuration file")
	}
	if r.ServiceConfig.Mode == system_config.Cluster && record.Type == core.LocalBackup {
		return nil, errors.New("local backup is not supported in cluster mode, use s3 backup instead")
	}
	err := record.Create(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	// send to task queue
	err = r.WorkerManager.EnqueuePersistentVolumeBackupRequest(record.ID)
	if err != nil {
		return nil, err
	}
	return persistentVolumeBackupToGraphqlObject(record), nil
}

// DeletePersistentVolumeBackup is the resolver for the deletePersistentVolumeBackup field.
func (r *mutationResolver) DeletePersistentVolumeBackup(ctx context.Context, id uint) (bool, error) {
	record := core.PersistentVolumeBackup{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	err = record.Delete(ctx, r.ServiceManager.DbClient, r.ServiceConfig.ServiceConfig.DataDir, r.ServiceConfig.PersistentVolumeBackupConfig.S3Config)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeletePersistentVolumeBackupsByPersistentVolumeID is the resolver for the deletePersistentVolumeBackupsByPersistentVolumeId field.
func (r *mutationResolver) DeletePersistentVolumeBackupsByPersistentVolumeID(ctx context.Context, persistentVolumeID uint) (bool, error) {
	err := core.DeletePersistentVolumeBackupsByPersistentVolumeId(ctx, r.ServiceManager.DbClient, persistentVolumeID, r.ServiceConfig.ServiceConfig.DataDir, r.ServiceConfig.PersistentVolumeBackupConfig.S3Config)
	if err != nil {
		return false, err
	}
	return true, nil
}
