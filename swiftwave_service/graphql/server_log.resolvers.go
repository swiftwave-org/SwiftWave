package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
)

// FetchServerLogContent is the resolver for the fetchServerLogContent field.
func (r *queryResolver) FetchServerLogContent(ctx context.Context, id uint) (string, error) {
	// fetch the server log content from the database
	serverLog, err := core.FetchServerLogContentByID(&r.ServiceManager.DbClient, id)
	if err != nil {
		return "", err
	}
	return serverLog, nil
}