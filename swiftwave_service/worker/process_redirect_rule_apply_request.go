package worker

import (
	"context"
	"errors"
	haproxymanager "github.com/swiftwave-org/swiftwave/pkg/haproxy_manager"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/manager"
	"gorm.io/gorm"
	"log"
)

func (m Manager) RedirectRuleApply(request RedirectRuleApplyRequest, ctx context.Context, _ context.CancelFunc) error {
	dbWithoutTx := m.ServiceManager.DbClient
	// fetch redirect rule
	redirectRule := &core.RedirectRule{}
	err := redirectRule.FindById(ctx, dbWithoutTx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	// ensure that the redirect rule is not being deleted
	if redirectRule.Status == core.RedirectRuleStatusDeleting {
		return nil
	}
	// fetch domain
	domain := &core.Domain{}
	err = domain.FindById(ctx, dbWithoutTx, redirectRule.DomainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	// fetch all proxy servers
	proxyServers, err := core.FetchProxyActiveServers(&m.ServiceManager.DbClient)
	if err != nil {
		return err
	}
	// don't attempt if no proxy servers are active
	if len(proxyServers) == 0 {
		return errors.New("no proxy servers are active")
	}
	// fetch all haproxy managers
	haproxyManagers, err := manager.HAProxyClients(context.Background(), proxyServers)
	if err != nil {
		return err
	}
	// map of server ip and transaction id
	transactionIdMap := make(map[*haproxymanager.Manager]string)
	isFailed := false

	for _, haproxyManager := range haproxyManagers {
		// fetch haproxy transaction
		haproxyTransactionId, err := haproxyManager.FetchNewTransactionId()
		if err != nil {
			return err
		}
		transactionIdMap[haproxyManager] = haproxyTransactionId
		// add redirect
		if redirectRule.Protocol == core.HTTPProtocol {
			err = haproxyManager.AddHTTPRedirectRule(haproxyTransactionId, domain.Name, redirectRule.RedirectURL)
		} else {
			err = haproxyManager.AddHTTPSRedirectRule(haproxyTransactionId, domain.Name, redirectRule.RedirectURL)
		}
		if err != nil {
			isFailed = true
			break
		}
	}

	for haproxyManager, haproxyTransactionId := range transactionIdMap {
		if !isFailed {
			// commit the haproxy transaction
			err = haproxyManager.CommitTransaction(haproxyTransactionId)
		}
		if isFailed || err != nil {
			isFailed = true
			log.Println("failed to commit haproxy transaction", err)
			err := haproxyManager.DeleteTransaction(haproxyTransactionId)
			if err != nil {
				log.Println("failed to rollback haproxy transaction", err)
			}
		}
	}

	// set status
	if isFailed {
		return redirectRule.UpdateStatus(ctx, dbWithoutTx, core.RedirectRuleStatusFailed)
	} else {
		return redirectRule.UpdateStatus(ctx, dbWithoutTx, core.RedirectRuleStatusApplied)
	}
}
