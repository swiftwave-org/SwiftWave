package worker

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"strings"

	haproxymanager "github.com/swiftwave-org/swiftwave/pkg/haproxy_manager"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/core"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/logger"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/manager"
	"gorm.io/gorm"
)

func (m Manager) SSLGenerate(request SSLGenerateRequest, ctx context.Context, _ context.CancelFunc) error {
	dbWithoutTx := m.ServiceManager.DbClient
	// fetch domain
	var domain core.Domain
	err := domain.FindById(ctx, dbWithoutTx, request.DomainId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// if not found, return nil as no queue is required
			return nil
		}
		return err
	}
	// If domain is IPv4, don't generate SSL
	if domain.IsIPv4() {
		return nil
	}
	// verify domain points to this server
	isDomainPointingToThisServer := m.ServiceManager.SslManager.VerifyDomain(domain.Name)
	if !isDomainPointingToThisServer {
		if domain.SSLStatus == core.DomainSSLStatusNone {
			// If SSL generation is invoked at the time of domain creation, don't mark it as failed if domain is not pointing to this server
			return nil
		}
		_ = domain.UpdateSSLStatus(ctx, dbWithoutTx, core.DomainSSLStatusFailed)
		logger.CronJobLoggerError.Println("Domain", domain.Name, "is not pointing to this server. Marking SSL Issue as failed")
		return nil
	}
	// generate private key [if not found]
	if domain.SSLPrivateKey == "" || strings.Compare(domain.SSLIssuer, "Let's Encrypt") != 0 {
		privateKey, err := generatePrivateKey()
		if err != nil {
			return err
		}
		domain.SSLPrivateKey = privateKey
		err = domain.Update(ctx, dbWithoutTx)
		if err != nil {
			logger.CronJobLoggerError.Println("Failed to update domain", domain.Name, "with private key.", err.Error(), "\nWill retry later")
			return err
		}
	}
	// obtain certificate
	fullChain, err := m.ServiceManager.SslManager.ObtainCertificate(domain.Name, domain.SSLPrivateKey)
	if err != nil {
		// don' requeue, if anything happen user can anytime re-request for certificate
		logger.CronJobLoggerError.Println("Failed to obtain certificate", err.Error())
		_ = domain.UpdateSSLStatus(ctx, dbWithoutTx, core.DomainSSLStatusFailed)
		return nil
	}
	// store certificate
	domain.SSLFullChain = fullChain
	// update status
	domain.SSLStatus = core.DomainSSLStatusIssued
	// enable auto renew
	domain.SslAutoRenew = true
	// update domain
	err = domain.Update(ctx, dbWithoutTx)
	if err != nil {
		logger.CronJobLoggerError.Println("Failed to update domain", domain.Name, "with SSL status", domain.SSLStatus, err.Error(), "\nWill retry later")
		return err
	}
	// fetch all proxy servers
	proxyServers, err := core.FetchProxyActiveServers(&m.ServiceManager.DbClient)
	if err != nil {
		logger.CronJobLoggerError.Println("Failed to fetch proxy servers while generating SSL certificate", err.Error(), "\nWill retry later")
		return err
	}
	// fetch all haproxy managers
	haproxyManagers, err := manager.HAProxyClients(context.Background(), proxyServers)
	if err != nil {
		logger.CronJobLoggerError.Println("Failed to fetch haproxy managers while generating SSL certificate", err.Error())
		return err
	}
	// map of server ip and transaction id
	transactionIdMap := make(map[*haproxymanager.Manager]string)
	isFailed := false

	for _, haproxyManager := range haproxyManagers {
		// generate a new transaction id for haproxy
		transactionId, err := haproxyManager.FetchNewTransactionId()
		if err != nil {
			logger.CronJobLoggerError.Println("Failed to fetch new transaction id while generating SSL certificate", err.Error(), "\nWill retry later")
			return err
		}
		// add to map
		transactionIdMap[haproxyManager] = transactionId
		// upload certificate to haproxy
		err = haproxyManager.UpdateSSL(transactionId, domain.Name, []byte(domain.SSLPrivateKey), []byte(domain.SSLFullChain))
		if err != nil {
			logger.CronJobLoggerError.Println("Failed to update SSL certificate while generating SSL certificate", err.Error(), "\nWill retry later")
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
			logger.CronJobLoggerError.Println("Failed to commit haproxy transaction while generating SSL certificate", err.Error())
			err := haproxyManager.DeleteTransaction(haproxyTransactionId)
			if err != nil {
				logger.CronJobLoggerError.Println("Failed to rollback haproxy transaction while generating SSL certificate", err.Error())
			}
		}
	}

	if isFailed {
		return domain.UpdateSSLStatus(ctx, dbWithoutTx, core.DomainSSLStatusFailed)
	}

	return nil
}

func (m Manager) SSLProxyUpdate(request SSLProxyUpdateRequest, ctx context.Context, _ context.CancelFunc) error {
	dbWithoutTx := m.ServiceManager.DbClient
	// fetch domain
	var domain core.Domain
	err := domain.FindById(ctx, dbWithoutTx, request.DomainId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// if not found, return nil as no queue is required
			return nil
		}
		return err
	}
	// If domain is IPv4, don't generate SSL
	if domain.IsIPv4() {
		return nil
	}
	// If domain doesn't have SSL, return
	if domain.SSLStatus == core.DomainSSLStatusNone || strings.Compare(domain.SSLFullChain, "") == 0 || strings.Compare(domain.SSLPrivateKey, "") == 0 {
		return nil
	}
	// fetch all proxy servers
	proxyServers, err := core.FetchProxyActiveServers(&m.ServiceManager.DbClient)
	if err != nil {
		return err
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
		// generate a new transaction id for haproxy
		transactionId, err := haproxyManager.FetchNewTransactionId()
		if err != nil {
			return err
		}
		// add to map
		transactionIdMap[haproxyManager] = transactionId
		// upload certificate to haproxy
		err = haproxyManager.UpdateSSL(transactionId, domain.Name, []byte(domain.SSLPrivateKey), []byte(domain.SSLFullChain))
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

	if isFailed {
		return domain.UpdateSSLStatus(ctx, dbWithoutTx, core.DomainSSLStatusFailed)
	} else {
		return domain.UpdateSSLStatus(ctx, dbWithoutTx, core.DomainSSLStatusIssued)
	}
}

// private functions
func generatePrivateKey() (string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", errors.New("unable to generate private key")
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemKey := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyBytes = pem.EncodeToMemory(&pemKey)
	return string(privateKeyBytes), nil
}
