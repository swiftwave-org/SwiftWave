package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

import (
	"context"
	"errors"

	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input *model.UserInput) (*model.User, error) {
	// Validate input
	if input.Username == "" {
		return nil, errors.New("username cannot be empty")
	}
	if input.Password == "" {
		return nil, errors.New("password cannot be empty")
	}
	// Find user by username
	_, err := core.FindUserByUsername(ctx, r.ServiceManager.DbClient, input.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	}
	user := core.User{
		Username: input.Username,
	}
	err = user.SetPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to set password")
	}
	user, err = core.CreateUser(ctx, r.ServiceManager.DbClient, user)
	if err != nil {
		return nil, err
	}
	return userToGraphqlObject(&user), nil
}

// DeleteUser is the resolver for the deleteUser field.
func (r *mutationResolver) DeleteUser(ctx context.Context, id uint) (bool, error) {
	// Fetch user
	user, err := core.FindUserByID(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		// Don't return error if record not found -- assume it's already deleted
		return true, nil
	}
	// Check if user is not current user
	currentUsername := ctx.Value("username").(string)
	if user.Username == currentUsername {
		return false, errors.New("cannot delete current user")
	}
	// Delete user
	err = core.DeleteUser(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return false, errors.New("failed to delete user")
	}
	return true, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, input *model.PasswordUpdateInput) (bool, error) {
	// Validate input
	if input.OldPassword == "" {
		return false, errors.New("old password cannot be empty")
	}
	if input.NewPassword == "" {
		return false, errors.New("new password cannot be empty")
	}
	// Change password
	username := ctx.Value("username").(string)
	err := core.ChangePassword(ctx, r.ServiceManager.DbClient, username, input.OldPassword, input.NewPassword)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	users, err := core.FindAllUsers(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	var result []*model.User
	for _, user := range users {
		result = append(result, userToGraphqlObject(&user))
	}
	return result, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id uint) (*model.User, error) {
	user, err := core.FindUserByID(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	return userToGraphqlObject(&user), nil
}

// CurrentUser is the resolver for the currentUser field.
func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	username := ctx.Value("username").(string)
	user, err := core.FindUserByUsername(ctx, r.ServiceManager.DbClient, username)
	if err != nil {
		return nil, err
	}
	return userToGraphqlObject(&user), nil
}
