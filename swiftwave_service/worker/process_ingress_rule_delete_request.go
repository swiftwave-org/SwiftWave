package worker

import (
	"context"
	"errors"
	"log"

	haproxymanager "github.com/swiftwave-org/swiftwave/haproxy_manager"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/manager"
	udpproxymanager "github.com/swiftwave-org/swiftwave/udp_proxy_manager"
	"gorm.io/gorm"
)

func (m Manager) IngressRuleDelete(request IngressRuleDeleteRequest, ctx context.Context, _ context.CancelFunc) error {
	dbWithoutTx := m.ServiceManager.DbClient
	// fetch ingress rule
	var ingressRule = &core.IngressRule{}
	err := ingressRule.FindById(ctx, dbWithoutTx, request.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	// check status should be deleting
	if ingressRule.Status != core.IngressRuleStatusDeleting {
		// dont requeue
		return nil
	}
	// fetch the domain
	domain := core.Domain{}
	if ingressRule.Protocol == core.HTTPProtocol || ingressRule.Protocol == core.HTTPSProtocol {
		if ingressRule.DomainID == nil {
			return errors.New("domain id is nil")
		}
		err = domain.FindById(ctx, dbWithoutTx, *ingressRule.DomainID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
	}

	// service name
	serviceName := ""

	if ingressRule.TargetType == core.ApplicationIngressRule {
		// fetch application
		application := &core.Application{}
		err = application.FindById(ctx, dbWithoutTx, *ingressRule.ApplicationID)
		if err != nil {
			return err
		}
		serviceName = application.Name
	} else if ingressRule.TargetType == core.ExternalServiceIngressRule {
		serviceName = ingressRule.ExternalService
	} else {
		return nil
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
	var haproxyManagers []*haproxymanager.Manager
	if isHAProxyAccessRequired(ingressRule) {
		haproxyManagers, err = manager.HAProxyClients(context.Background(), proxyServers)
		if err != nil {
			return err
		}
	}
	// fetch all udp proxy managers
	var udpProxyManagers []*udpproxymanager.Manager
	if isUDProxyAccessRequired(ingressRule) {
		udpProxyManagers, err = manager.UDPProxyClients(context.Background(), proxyServers)
		if err != nil {
			return err
		}
	}
	// map of server ip and transaction id
	transactionIdMap := make(map[*haproxymanager.Manager]string)
	isFailed := false

	// auth info
	authType := ingressRule.Authentication.AuthType
	authBasicUserlist := ""

	if authType == core.IngressRuleBasicAuthentication {
		record := &core.AppBasicAuthAccessControlList{}
		err = record.FindById(ctx, &dbWithoutTx, *ingressRule.Authentication.AppBasicAuthAccessControlListID)
		if err != nil {
			return err
		}
		authBasicUserlist = record.GeneratedName
	}

	for _, haproxyManager := range haproxyManagers {
		// check if ingress rules is not udp based
		if ingressRule.Protocol == core.UDPProtocol {
			continue
		}
		// backend protocol
		backendProtocol := ingressRuleProtocolToBackendProtocol(ingressRule.Protocol)
		// generate backend name
		backendName := haproxyManager.GenerateBackendName(backendProtocol, serviceName, int(ingressRule.TargetPort))
		// delete ingress rule from haproxy
		// create new haproxy transaction
		haproxyTransactionId, err := haproxyManager.FetchNewTransactionId()
		// store transaction id
		transactionIdMap[haproxyManager] = haproxyTransactionId
		if err != nil {
			continue
		}
		// delete ingress rule
		if ingressRule.Protocol == core.HTTPSProtocol {
			err = haproxyManager.DeleteHTTPSLink(haproxyTransactionId, backendName, domain.Name)
			if err != nil {
				// set status as failed and exit
				// because `DeleteHTTPSLink` can fail only if haproxy not working
				isFailed = true
				break
			}
			// remove https redirection if required
			if ingressRule.HttpsRedirect {
				err = haproxyManager.DisableHTTPSRedirection(haproxyTransactionId, domain.Name)
				if err != nil {
					isFailed = true
					break
				}
			}
		} else if ingressRule.Protocol == core.HTTPProtocol {
			if ingressRule.Port == 80 {
				err = haproxyManager.DeleteHTTPLink(haproxyTransactionId, backendName, domain.Name)
				if err != nil {
					// set status as failed and exit
					// because `DeleteHTTPLink` can fail only if haproxy not working
					isFailed = true
					break
				}
			} else {
				err = haproxyManager.DeleteTCPLink(haproxyTransactionId, backendName, int(ingressRule.Port), domain.Name, haproxymanager.HTTPMode)
				if err != nil {
					// set status as failed and exit
					// because `DeleteTCPLink` can fail only if haproxy not working
					isFailed = true
					break
				}
			}
		} else if ingressRule.Protocol == core.TCPProtocol {

			err = haproxyManager.DeleteTCPLink(haproxyTransactionId, backendName, int(ingressRule.Port), "", haproxymanager.TCPMode)
			if err != nil {
				// set status as failed and exit
				// because `DeleteTCPLink` can fail only if haproxy not working
				isFailed = true
				break
			}
		} else {
			// unknown protocol
			return nil
		}

		// disable basic auth if required
		if ingressRule.Protocol == core.HTTPProtocol || ingressRule.Protocol == core.HTTPSProtocol {
			if authType == core.IngressRuleBasicAuthentication {
				err = haproxyManager.RemoveBasicAuthentication(haproxyTransactionId, haproxymanager.HTTPMode, int(ingressRule.Port), domain.Name, authBasicUserlist)
				if err != nil {
					isFailed = true
					break
				}
			}
		}

		// delete backend
		backendUsedByOther := true
		var ingressRuleCheckForApplication core.IngressRule
		if ingressRule.TargetType == core.ApplicationIngressRule {
			err = m.ServiceManager.DbClient.Where("id != ? AND application_id = ? AND target_port = ?", ingressRule.ID, ingressRule.ApplicationID, ingressRule.TargetPort).First(&ingressRuleCheckForApplication).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					backendUsedByOther = false
				}
			}
		}
		if ingressRule.TargetType == core.ExternalServiceIngressRule {
			err = m.ServiceManager.DbClient.Where("id != ? AND external_service = ?", ingressRule.ID, ingressRule.ExternalService).First(&ingressRuleCheckForApplication).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					backendUsedByOther = false
				}
			}
		}
		if !backendUsedByOther {
			err = haproxyManager.DeleteBackend(haproxyTransactionId, backendName)
			if err != nil {
				// set status as failed and exit
				// because `DeleteBackend` can fail only if haproxy not working
				isFailed = true
				break
			}
		}
	}

	// delete ingress rule from udp proxy
	for _, udpProxyManager := range udpProxyManagers {
		if ingressRule.Protocol == core.UDPProtocol {
			err = udpProxyManager.Remove(udpproxymanager.Proxy{
				Port:       int(ingressRule.Port),
				TargetPort: int(ingressRule.TargetPort),
				Service:    serviceName,
			})
			if err != nil {
				// set status as failed and exit
				isFailed = true
				break
			}
		}
	}

	for haproxyManager, haproxyTransactionId := range transactionIdMap {
		if !isFailed {
			// commit the haproxy transaction
			err = haproxyManager.CommitTransaction(haproxyTransactionId)
			if err != nil {
				log.Println("committing haproxy transaction", haproxyTransactionId, err)
			}
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

	if isFailed {
		return ingressRule.UpdateStatus(ctx, dbWithoutTx, core.IngressRuleStatusFailed)
	} else {
		// delete ingress rule from database
		return ingressRule.Delete(ctx, dbWithoutTx, true)
	}
}
