package cmd

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	SSL "github.com/swiftwave-org/swiftwave/pkg/ssl_manager"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/config/local_config"
	"github.com/swiftwave-org/swiftwave/swiftwave_service/db"
)

func init() {
	tlsCmd.AddCommand(tlsEnableCmd)
	tlsCmd.AddCommand(tlsDisableCmd)
	tlsCmd.AddCommand(generateCertificateCommand)
	tlsCmd.AddCommand(renewCertificateCommand)
}

var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "Manage TLS for swiftwave service",
	Long:  `Manage TLS for swiftwave service`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

var tlsEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable TLS for swiftwave service",
	Long:  `Enable TLS for swiftwave service`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.LocalConfig.ServiceConfig.UseTLS {
			printSuccess("TLS is already enabled")
			return
		}
		// Check if some certificate is already present
		if isFolderEmpty(config.LocalConfig.ServiceConfig.SSLCertDirectoryPath) {
			printError("No TLS certificate found")
			printInfo("Use `swiftwave tls generate` to generate a new certificate")
			return
		}
		config.LocalConfig.ServiceConfig.UseTLS = true
		config.LocalConfig.ServiceConfig.AutoRenewManagementNodeCert = true
		err := local_config.Update(config.LocalConfig)
		if err != nil {
			printError("Failed to update config")
			printError(err.Error())
			return
		}
		printSuccess("TLS has been enabled")
		restartSysctlService("swiftwave")
		// Restart local registry if required
		restartLocalRegistryIfRequired(cmd.Context())
	},
}

var tlsDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable TLS for swiftwave service",
	Long:  `Disable TLS for swiftwave service`,
	Run: func(cmd *cobra.Command, args []string) {
		lConfig := config.LocalConfig
		if !lConfig.ServiceConfig.UseTLS {
			printSuccess("TLS is already disabled")
			return
		}
		lConfig.ServiceConfig.UseTLS = false
		config.LocalConfig.ServiceConfig.AutoRenewManagementNodeCert = false
		err := local_config.Update(lConfig)
		if err != nil {
			printError("Failed to update config")
			printError(err.Error())
			return
		}
		printSuccess("TLS has been disabled")
		restartSysctlService("swiftwave")
		// Restart local registry if required
		restartLocalRegistryIfRequired(cmd.Context())
	},
}

var generateCertificateCommand = &cobra.Command{
	Use:   "generate",
	Short: "Generate TLS certificate for swiftwave endpoints",
	Long: `This command generates TLS certificate for swiftwave endpoints.
	It's not for generating certificates for domain of hosted applications`,
	Run: func(cmd *cobra.Command, args []string) {
		domain := config.LocalConfig.ServiceConfig.ManagementNodeAddress
		// Start http-01 challenge server
		echoServer := echo.New()
		echoServer.HideBanner = true
		echoServer.Pre(middleware.RemoveTrailingSlash())
		// Initiating database client
		dbClient, err := db.GetClient(config.LocalConfig, 10)
		if err != nil {
			printError("Failed to connect to database")
			return
		}
		// Initiating SSL Manager
		options := SSL.ManagerOptions{
			IsStaging:         config.SystemConfig.LetsEncryptConfig.Staging,
			Email:             config.SystemConfig.LetsEncryptConfig.EmailID,
			AccountPrivateKey: config.SystemConfig.LetsEncryptConfig.PrivateKey,
		}
		sslManager := SSL.Manager{}
		err = sslManager.Init(context.Background(), *dbClient, options)
		if err != nil {
			printError("Failed to initiate SSL Manager")
			return
		}
		// Check if there is already someone listening on port 80
		isServerStarted := false
		isPort80Blocked := checkIfPortIsInUse("80")
		isServicePortBlocked := checkIfPortIsInUse(strconv.Itoa(config.LocalConfig.ServiceConfig.BindPort))
		if isPort80Blocked {
			if isServicePortBlocked {
				printInfo("Already running swiftwave service will be used to solve http-01 challenge")
			} else {
				printError("Please stop the service running on port 80 temporarily")
				return
			}
		} else {
			isServerStarted = true
			// Start the server
			go func(sslManager *SSL.Manager) {
				sslManager.InitHttpHandlers(echoServer)
				err := echoServer.Start(":80")
				if err != nil {
					if errors.Is(err, http.ErrServerClosed) {
						printSuccess("http-01 challenge server has been stopped")
					} else {
						printError("Failed to start http-01 challenge server")
						os.Exit(1)
					}
				}
			}(&sslManager)
		}
		// Generate private key
		privateKey, err := generatePrivateKey()
		if err != nil {
			printError("Failed to generate private key")
			os.Exit(1)
			return
		}
		// Generate the certificate
		certificate, err := sslManager.ObtainCertificate(domain, privateKey)
		if err != nil {
			println(err.Error())
			printError("Failed to generate certificate")
			os.Exit(1)
			return
		}
		if isServerStarted {
			// Stop the http-01 challenge server
			err = echoServer.Server.Shutdown(context.Background())
			if err != nil {
				return
			}
		}
		// Store private key and certificate in the service.ssl_certificate_dir folder
		certDir := config.LocalConfig.ServiceConfig.SSLCertDirectoryPath
		// Store private key
		err = os.WriteFile(filepath.Join(certDir, "private.key"), []byte(privateKey), 0600)
		if err != nil {
			printError("Failed to store private key")
			os.Exit(1)
			return
		}
		// Store certificate
		err = os.WriteFile(filepath.Join(certDir, "certificate.crt"), []byte(certificate), 0600)
		if err != nil {
			printError("Failed to store certificate")
			os.Exit(1)
			return
		}
		// Print success message
		printSuccess("Successfully generated TLS certificate for " + domain)
		// Enable TLS for swiftwave service
		config.LocalConfig.ServiceConfig.UseTLS = true
		config.LocalConfig.ServiceConfig.AutoRenewManagementNodeCert = true
		err = local_config.Update(config.LocalConfig)
		if err != nil {
			printError("Failed to update config")
			os.Exit(1)
			return
		} else {
			printSuccess("TLS has been enabled")
			printInfo(fmt.Sprintf("Access dashboard at https://%s:%d", domain, config.LocalConfig.ServiceConfig.BindPort))
		}
		// Restart swiftwave service
		restartSysctlService("swiftwave")
		// Restart local registry if required
		restartLocalRegistryIfRequired(cmd.Context())
	},
}

var renewCertificateCommand = &cobra.Command{
	Use:   "renew",
	Short: "Renew TLS certificate for swiftwave endpoints",
	Long: `This command renews TLS certificates for swiftwave endpoints.
	It's not for renewing certificates for domain of hosted applications`,
	Run: func(cmd *cobra.Command, args []string) {
		sslCertificatePath := filepath.Join(config.LocalConfig.ServiceConfig.SSLCertDirectoryPath, "certificate.crt")
		if _, err := os.Stat(sslCertificatePath); os.IsNotExist(err) {
			printError("No TLS certificate found")
			printInfo("Use `swiftwave tls generate` to generate a new certificate")
			os.Exit(1)
		}
		isRenewalRequired, err := isRenewalImminent(sslCertificatePath)
		if err != nil {
			printError("Failed to check if renewal is required")
			printError(err.Error())
			os.Exit(1)
		}
		if isRenewalRequired {
			printSuccess("Renewal is required")
			generateCertificateCommand.Run(cmd, args)
		}
	},
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

func isFolderEmpty(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		return true
	}
	return len(files) == 0
}

func restartSysctlService(serviceName string) {
	// check if service is running
	// read the output of systemctl is-active <service_name>
	cmd := exec.Command("systemctl", "is-active", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return
	}
	if strings.TrimSpace(string(output)) == "active" {
		// restart the service
		cmd = exec.Command("systemctl", "restart", serviceName)
		err = cmd.Run()
		if err != nil {
			return
		}
		printSuccess(serviceName + " has been restarted")
	}
}

func daysUntilExpiration(certPath string) (int, error) {
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return 0, err
	}

	block, _ := pem.Decode(certBytes)
	if block == nil {
		return 0, fmt.Errorf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return 0, err
	}

	daysRemaining := int(time.Until(cert.NotAfter).Hours() / 24)
	return daysRemaining, nil
}

func isRenewalImminent(certPath string) (bool, error) {
	daysRemaining, err := daysUntilExpiration(certPath)
	if err != nil {
		return false, err
	}

	return daysRemaining <= 30, nil
}

func restartLocalRegistryIfRequired(ctx context.Context) {
	if config == nil || config.LocalConfig == nil || config.SystemConfig == nil {
		return
	}
	isRequired, err := isLocalRegistryRequired()
	if err != nil {
		printError("Failed to check if local registry is required")
		printError(err.Error())
		return
	}
	if isRequired {
		err := restartLocalRegistry(ctx)
		if err != nil {
			printError("Failed to restart local image registry")
			printError(err.Error())
		} else {
			printSuccess("Local registry restarted successfully")
		}
	} else {
		printInfo("[IGNORE] Local registry restart is not required, as not configured")
	}
}
