package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.48

import (
	"context"
	"errors"
	"strings"

	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/graphql/model"
)

// IngressRules is the resolver for the ingressRules field.
func (r *domainResolver) IngressRules(ctx context.Context, obj *model.Domain) ([]*model.IngressRule, error) {
	records, err := core.FindIngressRulesByDomainID(ctx, r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*model.IngressRule
	for _, record := range records {
		result = append(result, ingressRuleToGraphqlObject(record))
	}
	return result, nil
}

// RedirectRules is the resolver for the redirectRules field.
func (r *domainResolver) RedirectRules(ctx context.Context, obj *model.Domain) ([]*model.RedirectRule, error) {
	records, err := core.FindRedirectRulesByDomainID(ctx, r.ServiceManager.DbClient, obj.ID)
	if err != nil {
		return nil, err
	}
	var result []*model.RedirectRule
	for _, record := range records {
		result = append(result, redirectRuleToGraphqlObject(record))
	}
	return result, nil
}

// AddDomain is the resolver for the addDomain field.
func (r *mutationResolver) AddDomain(ctx context.Context, input model.DomainInput) (*model.Domain, error) {
	record := domainInputToDatabaseObject(&input)
	if record.Name == "" {
		return nil, errors.New("name is required")
	}
	err := record.Create(ctx, r.ServiceManager.DbClient)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, errors.New("domain with same name already exists")
		}
		return nil, err
	}
	// Just enqueue the domain for SSL generation
	_ = r.WorkerManager.EnqueueSSLGenerateRequest(record.ID)
	return domainToGraphqlObject(record), nil
}

// RemoveDomain is the resolver for the removeDomain field.
func (r *mutationResolver) RemoveDomain(ctx context.Context, id uint) (bool, error) {
	record := core.Domain{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return false, err
	}
	err = record.Delete(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return false, err
	}
	return true, nil
}

// IssueSsl is the resolver for the issueSSL field.
func (r *mutationResolver) IssueSsl(ctx context.Context, id uint) (*model.Domain, error) {
	// fetch record
	record := core.Domain{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	// verify domain configuration
	configured := r.ServiceManager.SslManager.VerifyDomain(record.Name)
	if !configured {
		return nil, errors.New("domain not configured")
	}
	// update record
	record.SSLStatus = core.DomainSSLStatusPending
	err = record.Update(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	// push task
	err = r.WorkerManager.EnqueueSSLGenerateRequest(record.ID)
	if err != nil {
		// rollback status
		record.SSLStatus = core.DomainSSLStatusNone
		err = record.Update(ctx, r.ServiceManager.DbClient)
		if err != nil {
			return nil, errors.New("failed to enqueue ssl generation request")
		}
		return nil, errors.New("failed to enqueue ssl generation request")
	}

	return domainToGraphqlObject(&record), nil
}

// AddCustomSsl is the resolver for the addCustomSSL field.
func (r *mutationResolver) AddCustomSsl(ctx context.Context, id uint, input model.CustomSSLInput) (*model.Domain, error) {
	// fetch record
	record := core.Domain{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	// validate certificate full chain
	err = ValidateSSLFullChainCertificate(input.FullChain)
	if err != nil {
		return nil, err
	}

	// validate certificate private key
	err = ValidateSSLPrivateKey(input.PrivateKey)
	if err != nil {
		return nil, err
	}
	// update record
	record.SSLPrivateKey = input.PrivateKey
	record.SSLFullChain = input.FullChain
	record.SSLStatus = core.DomainSSLStatusPending
	record.SslAutoRenew = false
	err = record.Update(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	// push task
	err = r.WorkerManager.EnqueueSSLProxyUpdateRequest(id)
	if err != nil {
		// rollback status
		record.SSLStatus = core.DomainSSLStatusNone
		_ = record.Update(ctx, r.ServiceManager.DbClient)
		return nil, errors.New("failed to enqueue ssl proxy update request")
	}
	return domainToGraphqlObject(&record), nil
}

// Domains is the resolver for the domains field.
func (r *queryResolver) Domains(ctx context.Context) ([]*model.Domain, error) {
	records, err := core.FindAllDomains(ctx, r.ServiceManager.DbClient)
	if err != nil {
		return nil, err
	}
	var result []*model.Domain
	for _, record := range records {
		result = append(result, domainToGraphqlObject(record))
	}
	return result, nil
}

// Domain is the resolver for the domain field.
func (r *queryResolver) Domain(ctx context.Context, id uint) (*model.Domain, error) {
	record := core.Domain{}
	err := record.FindById(ctx, r.ServiceManager.DbClient, id)
	if err != nil {
		return nil, err
	}
	return domainToGraphqlObject(&record), nil
}

// VerifyDomainConfiguration is the resolver for the verifyDomainConfiguration field.
func (r *queryResolver) VerifyDomainConfiguration(ctx context.Context, name string) (bool, error) {
	return r.ServiceManager.SslManager.VerifyDomain(name), nil
}

// Domain returns DomainResolver implementation.
func (r *Resolver) Domain() DomainResolver { return &domainResolver{r} }

type domainResolver struct{ *Resolver }
