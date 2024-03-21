// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type Application struct {
	ID                       string                     `json:"id"`
	Name                     string                     `json:"name"`
	EnvironmentVariables     []*EnvironmentVariable     `json:"environmentVariables"`
	PersistentVolumeBindings []*PersistentVolumeBinding `json:"persistentVolumeBindings"`
	Capabilities             []string                   `json:"capabilities"`
	Sysctls                  []string                   `json:"sysctls"`
	RealtimeInfo             *RealtimeInfo              `json:"realtimeInfo"`
	LatestDeployment         *Deployment                `json:"latestDeployment"`
	Deployments              []*Deployment              `json:"deployments"`
	DeploymentMode           DeploymentMode             `json:"deploymentMode"`
	Replicas                 uint                       `json:"replicas"`
	IngressRules             []*IngressRule             `json:"ingressRules"`
	IsDeleted                bool                       `json:"isDeleted"`
	WebhookToken             string                     `json:"webhookToken"`
	IsSleeping               bool                       `json:"isSleeping"`
	Command                  string                     `json:"command"`
}

type ApplicationDeployResult struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	Application *Application `json:"application,omitempty"`
}

type ApplicationInput struct {
	Name                         string                          `json:"name"`
	EnvironmentVariables         []*EnvironmentVariableInput     `json:"environmentVariables"`
	PersistentVolumeBindings     []*PersistentVolumeBindingInput `json:"persistentVolumeBindings"`
	Capabilities                 []string                        `json:"capabilities"`
	Sysctls                      []string                        `json:"sysctls"`
	Dockerfile                   *string                         `json:"dockerfile,omitempty"`
	BuildArgs                    []*BuildArgInput                `json:"buildArgs"`
	DeploymentMode               DeploymentMode                  `json:"deploymentMode"`
	Replicas                     *uint                           `json:"replicas,omitempty"`
	UpstreamType                 UpstreamType                    `json:"upstreamType"`
	Command                      string                          `json:"command"`
	GitCredentialID              *uint                           `json:"gitCredentialID,omitempty"`
	GitProvider                  *GitProvider                    `json:"gitProvider,omitempty"`
	RepositoryOwner              *string                         `json:"repositoryOwner,omitempty"`
	RepositoryName               *string                         `json:"repositoryName,omitempty"`
	RepositoryBranch             *string                         `json:"repositoryBranch,omitempty"`
	CodePath                     *string                         `json:"codePath,omitempty"`
	SourceCodeCompressedFileName *string                         `json:"sourceCodeCompressedFileName,omitempty"`
	DockerImage                  *string                         `json:"dockerImage,omitempty"`
	ImageRegistryCredentialID    *uint                           `json:"imageRegistryCredentialID,omitempty"`
}

type ApplicationResourceAnalytics struct {
	CPUUsagePercent int       `json:"cpu_usage_percent"`
	MemoryUsedMb    uint64    `json:"memory_used_mb"`
	NetworkSentKbps uint64    `json:"network_sent_kbps"`
	NetworkRecvKbps uint64    `json:"network_recv_kbps"`
	Timestamp       time.Time `json:"timestamp"`
}

type BuildArg struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BuildArgInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CustomSSLInput struct {
	FullChain  string `json:"fullChain"`
	PrivateKey string `json:"privateKey"`
	SslIssuer  string `json:"sslIssuer"`
}

type Dependency struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

type Deployment struct {
	ID                           string                   `json:"id"`
	ApplicationID                string                   `json:"applicationID"`
	Application                  *Application             `json:"application"`
	UpstreamType                 UpstreamType             `json:"upstreamType"`
	GitCredentialID              uint                     `json:"gitCredentialID"`
	GitCredential                *GitCredential           `json:"gitCredential"`
	GitProvider                  GitProvider              `json:"gitProvider"`
	RepositoryOwner              string                   `json:"repositoryOwner"`
	RepositoryName               string                   `json:"repositoryName"`
	RepositoryBranch             string                   `json:"repositoryBranch"`
	CommitHash                   string                   `json:"commitHash"`
	CodePath                     string                   `json:"codePath"`
	SourceCodeCompressedFileName string                   `json:"sourceCodeCompressedFileName"`
	DockerImage                  string                   `json:"dockerImage"`
	ImageRegistryCredentialID    uint                     `json:"imageRegistryCredentialID"`
	ImageRegistryCredential      *ImageRegistryCredential `json:"imageRegistryCredential"`
	BuildArgs                    []*BuildArg              `json:"buildArgs"`
	Dockerfile                   string                   `json:"dockerfile"`
	Status                       DeploymentStatus         `json:"status"`
	CreatedAt                    time.Time                `json:"createdAt"`
}

type DeploymentLog struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type DockerConfigBuildArg struct {
	Key          string `json:"key"`
	Description  string `json:"description"`
	Type         string `json:"type"`
	DefaultValue string `json:"defaultValue"`
}

type DockerConfigGeneratorInput struct {
	SourceType                   DockerConfigSourceType `json:"sourceType"`
	GitCredentialID              *uint                  `json:"gitCredentialID,omitempty"`
	GitProvider                  *GitProvider           `json:"gitProvider,omitempty"`
	RepositoryOwner              *string                `json:"repositoryOwner,omitempty"`
	RepositoryName               *string                `json:"repositoryName,omitempty"`
	RepositoryBranch             *string                `json:"repositoryBranch,omitempty"`
	CodePath                     *string                `json:"codePath,omitempty"`
	SourceCodeCompressedFileName *string                `json:"sourceCodeCompressedFileName,omitempty"`
	CustomDockerFile             *string                `json:"customDockerFile,omitempty"`
}

type DockerConfigGeneratorOutput struct {
	DetectedServiceName *string                 `json:"detectedServiceName,omitempty"`
	DockerFile          *string                 `json:"dockerFile,omitempty"`
	DockerBuildArgs     []*DockerConfigBuildArg `json:"dockerBuildArgs,omitempty"`
}

type Domain struct {
	ID            uint            `json:"id"`
	Name          string          `json:"name"`
	SslStatus     DomainSSLStatus `json:"sslStatus"`
	SslFullChain  string          `json:"sslFullChain"`
	SslPrivateKey string          `json:"sslPrivateKey"`
	SslIssuedAt   time.Time       `json:"sslIssuedAt"`
	SslIssuer     string          `json:"sslIssuer"`
	SslAutoRenew  bool            `json:"sslAutoRenew"`
	IngressRules  []*IngressRule  `json:"ingressRules"`
	RedirectRules []*RedirectRule `json:"redirectRules"`
}

type DomainInput struct {
	Name string `json:"name"`
}

type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EnvironmentVariableInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GitCredential struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Deployments []*Deployment `json:"deployments"`
}

type GitCredentialInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type GitCredentialRepositoryAccessInput struct {
	GitCredentialID  uint   `json:"gitCredentialId"`
	RepositoryURL    string `json:"repositoryUrl"`
	RepositoryBranch string `json:"repositoryBranch"`
}

type GitCredentialRepositoryAccessResult struct {
	GitCredentialID  uint           `json:"gitCredentialId"`
	GitCredential    *GitCredential `json:"gitCredential"`
	RepositoryURL    string         `json:"repositoryUrl"`
	RepositoryBranch string         `json:"repositoryBranch"`
	Success          bool           `json:"success"`
	Error            string         `json:"error"`
}

type ImageRegistryCredential struct {
	ID          uint          `json:"id"`
	URL         string        `json:"url"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Deployments []*Deployment `json:"deployments"`
}

type ImageRegistryCredentialInput struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type IngressRule struct {
	ID            uint              `json:"id"`
	DomainID      *uint             `json:"domainId,omitempty"`
	Domain        *Domain           `json:"domain,omitempty"`
	Protocol      ProtocolType      `json:"protocol"`
	Port          uint              `json:"port"`
	ApplicationID string            `json:"applicationId"`
	Application   *Application      `json:"application"`
	TargetPort    uint              `json:"targetPort"`
	Status        IngressRuleStatus `json:"status"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type IngressRuleInput struct {
	DomainID      *uint        `json:"domainId,omitempty"`
	ApplicationID string       `json:"applicationId"`
	Protocol      ProtocolType `json:"protocol"`
	Port          uint         `json:"port"`
	TargetPort    uint         `json:"targetPort"`
}

type Mutation struct {
}

type NFSConfig struct {
	Host    string `json:"host"`
	Path    string `json:"path"`
	Version int    `json:"version"`
}

type NFSConfigInput struct {
	Host    string `json:"host"`
	Path    string `json:"path"`
	Version int    `json:"version"`
}

type NewServerInput struct {
	IP   string `json:"ip"`
	User string `json:"user"`
}

type PasswordUpdateInput struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type PersistentVolume struct {
	ID                       uint                       `json:"id"`
	Name                     string                     `json:"name"`
	Type                     PersistentVolumeType       `json:"type"`
	NfsConfig                *NFSConfig                 `json:"nfsConfig"`
	PersistentVolumeBindings []*PersistentVolumeBinding `json:"persistentVolumeBindings"`
	Backups                  []*PersistentVolumeBackup  `json:"backups"`
	Restores                 []*PersistentVolumeRestore `json:"restores"`
}

type PersistentVolumeBackup struct {
	ID          uint                         `json:"id"`
	Type        PersistentVolumeBackupType   `json:"type"`
	Status      PersistentVolumeBackupStatus `json:"status"`
	SizeMb      float64                      `json:"sizeMb"`
	CreatedAt   time.Time                    `json:"createdAt"`
	CompletedAt time.Time                    `json:"completedAt"`
}

type PersistentVolumeBackupInput struct {
	PersistentVolumeID uint                       `json:"persistentVolumeId"`
	Type               PersistentVolumeBackupType `json:"type"`
}

type PersistentVolumeBinding struct {
	ID                 uint              `json:"id"`
	PersistentVolumeID uint              `json:"persistentVolumeID"`
	PersistentVolume   *PersistentVolume `json:"persistentVolume"`
	ApplicationID      string            `json:"applicationID"`
	Application        *Application      `json:"application"`
	MountingPath       string            `json:"mountingPath"`
}

type PersistentVolumeBindingInput struct {
	PersistentVolumeID uint   `json:"persistentVolumeID"`
	MountingPath       string `json:"mountingPath"`
}

type PersistentVolumeInput struct {
	Name      string               `json:"name"`
	Type      PersistentVolumeType `json:"type"`
	NfsConfig *NFSConfigInput      `json:"nfsConfig"`
}

type PersistentVolumeRestore struct {
	ID          uint                          `json:"id"`
	Type        PersistentVolumeRestoreType   `json:"type"`
	Status      PersistentVolumeRestoreStatus `json:"status"`
	CreatedAt   time.Time                     `json:"createdAt"`
	CompletedAt time.Time                     `json:"completedAt"`
}

type PersistentVolumeRestoreInput struct {
	PersistentVolumeID uint                        `json:"persistentVolumeId"`
	Type               PersistentVolumeRestoreType `json:"type"`
}

type Query struct {
}

type RealtimeInfo struct {
	InfoFound       bool           `json:"InfoFound"`
	DesiredReplicas int            `json:"DesiredReplicas"`
	RunningReplicas int            `json:"RunningReplicas"`
	DeploymentMode  DeploymentMode `json:"DeploymentMode"`
}

type RedirectRule struct {
	ID          uint               `json:"id"`
	DomainID    uint               `json:"domainId"`
	Domain      *Domain            `json:"domain"`
	Protocol    ProtocolType       `json:"protocol"`
	RedirectURL string             `json:"redirectURL"`
	Status      RedirectRuleStatus `json:"status"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type RedirectRuleInput struct {
	DomainID    uint         `json:"domainId"`
	Protocol    ProtocolType `json:"protocol"`
	RedirectURL string       `json:"redirectURL"`
}

type RuntimeLog struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type Server struct {
	ID                   uint         `json:"id"`
	IP                   string       `json:"ip"`
	Hostname             string       `json:"hostname"`
	User                 string       `json:"user"`
	SwarmMode            SwarmMode    `json:"swarmMode"`
	ScheduleDeployments  bool         `json:"scheduleDeployments"`
	DockerUnixSocketPath string       `json:"dockerUnixSocketPath"`
	ProxyEnabled         bool         `json:"proxyEnabled"`
	ProxyType            ProxyType    `json:"proxyType"`
	Status               ServerStatus `json:"status"`
	Logs                 []*ServerLog `json:"logs"`
}

type ServerDiskUsage struct {
	Path       string    `json:"path"`
	MountPoint string    `json:"mount_point"`
	TotalGb    float64   `json:"total_gb"`
	UsedGb     float64   `json:"used_gb"`
	Timestamp  time.Time `json:"timestamp"`
}

type ServerDisksUsage struct {
	Disks     []*ServerDiskUsage `json:"disks"`
	Timestamp time.Time          `json:"timestamp"`
}

type ServerLog struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ServerResourceAnalytics struct {
	CPUUsagePercent int       `json:"cpu_usage_percent"`
	MemoryTotalGb   float64   `json:"memory_total_gb"`
	MemoryUsedGb    float64   `json:"memory_used_gb"`
	MemoryCachedGb  float64   `json:"memory_cached_gb"`
	NetworkSentKbps uint64    `json:"network_sent_kbps"`
	NetworkRecvKbps uint64    `json:"network_recv_kbps"`
	Timestamp       time.Time `json:"timestamp"`
}

type ServerSetupInput struct {
	ID                   uint      `json:"id"`
	DockerUnixSocketPath string    `json:"dockerUnixSocketPath"`
	SwarmMode            SwarmMode `json:"swarmMode"`
}

type StackInput struct {
	Content   string               `json:"content"`
	Variables []*StackVariableType `json:"variables"`
}

type StackVariableType struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type StackVerifyResult struct {
	Success         bool     `json:"success"`
	Message         string   `json:"message"`
	Error           string   `json:"error"`
	ValidVolumes    []string `json:"validVolumes"`
	InvalidVolumes  []string `json:"invalidVolumes"`
	ValidServices   []string `json:"validServices"`
	InvalidServices []string `json:"invalidServices"`
}

type Subscription struct {
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApplicationResourceAnalyticsTimeframe string

const (
	ApplicationResourceAnalyticsTimeframeLast1Hour   ApplicationResourceAnalyticsTimeframe = "last_1_hour"
	ApplicationResourceAnalyticsTimeframeLast24Hours ApplicationResourceAnalyticsTimeframe = "last_24_hours"
	ApplicationResourceAnalyticsTimeframeLast7Days   ApplicationResourceAnalyticsTimeframe = "last_7_days"
	ApplicationResourceAnalyticsTimeframeLast30Days  ApplicationResourceAnalyticsTimeframe = "last_30_days"
)

var AllApplicationResourceAnalyticsTimeframe = []ApplicationResourceAnalyticsTimeframe{
	ApplicationResourceAnalyticsTimeframeLast1Hour,
	ApplicationResourceAnalyticsTimeframeLast24Hours,
	ApplicationResourceAnalyticsTimeframeLast7Days,
	ApplicationResourceAnalyticsTimeframeLast30Days,
}

func (e ApplicationResourceAnalyticsTimeframe) IsValid() bool {
	switch e {
	case ApplicationResourceAnalyticsTimeframeLast1Hour, ApplicationResourceAnalyticsTimeframeLast24Hours, ApplicationResourceAnalyticsTimeframeLast7Days, ApplicationResourceAnalyticsTimeframeLast30Days:
		return true
	}
	return false
}

func (e ApplicationResourceAnalyticsTimeframe) String() string {
	return string(e)
}

func (e *ApplicationResourceAnalyticsTimeframe) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ApplicationResourceAnalyticsTimeframe(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ApplicationResourceAnalyticsTimeframe", str)
	}
	return nil
}

func (e ApplicationResourceAnalyticsTimeframe) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DeploymentMode string

const (
	DeploymentModeReplicated DeploymentMode = "replicated"
	DeploymentModeGlobal     DeploymentMode = "global"
)

var AllDeploymentMode = []DeploymentMode{
	DeploymentModeReplicated,
	DeploymentModeGlobal,
}

func (e DeploymentMode) IsValid() bool {
	switch e {
	case DeploymentModeReplicated, DeploymentModeGlobal:
		return true
	}
	return false
}

func (e DeploymentMode) String() string {
	return string(e)
}

func (e *DeploymentMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DeploymentMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DeploymentMode", str)
	}
	return nil
}

func (e DeploymentMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DeploymentStatus string

const (
	DeploymentStatusPending       DeploymentStatus = "pending"
	DeploymentStatusDeployPending DeploymentStatus = "deployPending"
	DeploymentStatusDeploying     DeploymentStatus = "deploying"
	DeploymentStatusLive          DeploymentStatus = "live"
	DeploymentStatusStopped       DeploymentStatus = "stopped"
	DeploymentStatusFailed        DeploymentStatus = "failed"
	DeploymentStatusStalled       DeploymentStatus = "stalled"
)

var AllDeploymentStatus = []DeploymentStatus{
	DeploymentStatusPending,
	DeploymentStatusDeployPending,
	DeploymentStatusDeploying,
	DeploymentStatusLive,
	DeploymentStatusStopped,
	DeploymentStatusFailed,
	DeploymentStatusStalled,
}

func (e DeploymentStatus) IsValid() bool {
	switch e {
	case DeploymentStatusPending, DeploymentStatusDeployPending, DeploymentStatusDeploying, DeploymentStatusLive, DeploymentStatusStopped, DeploymentStatusFailed, DeploymentStatusStalled:
		return true
	}
	return false
}

func (e DeploymentStatus) String() string {
	return string(e)
}

func (e *DeploymentStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DeploymentStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DeploymentStatus", str)
	}
	return nil
}

func (e DeploymentStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DockerConfigSourceType string

const (
	DockerConfigSourceTypeGit        DockerConfigSourceType = "git"
	DockerConfigSourceTypeSourceCode DockerConfigSourceType = "sourceCode"
	DockerConfigSourceTypeCustom     DockerConfigSourceType = "custom"
)

var AllDockerConfigSourceType = []DockerConfigSourceType{
	DockerConfigSourceTypeGit,
	DockerConfigSourceTypeSourceCode,
	DockerConfigSourceTypeCustom,
}

func (e DockerConfigSourceType) IsValid() bool {
	switch e {
	case DockerConfigSourceTypeGit, DockerConfigSourceTypeSourceCode, DockerConfigSourceTypeCustom:
		return true
	}
	return false
}

func (e DockerConfigSourceType) String() string {
	return string(e)
}

func (e *DockerConfigSourceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DockerConfigSourceType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DockerConfigSourceType", str)
	}
	return nil
}

func (e DockerConfigSourceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DomainSSLStatus string

const (
	DomainSSLStatusNone    DomainSSLStatus = "none"
	DomainSSLStatusPending DomainSSLStatus = "pending"
	DomainSSLStatusIssued  DomainSSLStatus = "issued"
	DomainSSLStatusFailed  DomainSSLStatus = "failed"
)

var AllDomainSSLStatus = []DomainSSLStatus{
	DomainSSLStatusNone,
	DomainSSLStatusPending,
	DomainSSLStatusIssued,
	DomainSSLStatusFailed,
}

func (e DomainSSLStatus) IsValid() bool {
	switch e {
	case DomainSSLStatusNone, DomainSSLStatusPending, DomainSSLStatusIssued, DomainSSLStatusFailed:
		return true
	}
	return false
}

func (e DomainSSLStatus) String() string {
	return string(e)
}

func (e *DomainSSLStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DomainSSLStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DomainSSLStatus", str)
	}
	return nil
}

func (e DomainSSLStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type GitProvider string

const (
	GitProviderNone   GitProvider = "none"
	GitProviderGithub GitProvider = "github"
	GitProviderGitlab GitProvider = "gitlab"
)

var AllGitProvider = []GitProvider{
	GitProviderNone,
	GitProviderGithub,
	GitProviderGitlab,
}

func (e GitProvider) IsValid() bool {
	switch e {
	case GitProviderNone, GitProviderGithub, GitProviderGitlab:
		return true
	}
	return false
}

func (e GitProvider) String() string {
	return string(e)
}

func (e *GitProvider) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = GitProvider(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid GitProvider", str)
	}
	return nil
}

func (e GitProvider) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type IngressRuleStatus string

const (
	IngressRuleStatusPending  IngressRuleStatus = "pending"
	IngressRuleStatusApplied  IngressRuleStatus = "applied"
	IngressRuleStatusDeleting IngressRuleStatus = "deleting"
	IngressRuleStatusFailed   IngressRuleStatus = "failed"
)

var AllIngressRuleStatus = []IngressRuleStatus{
	IngressRuleStatusPending,
	IngressRuleStatusApplied,
	IngressRuleStatusDeleting,
	IngressRuleStatusFailed,
}

func (e IngressRuleStatus) IsValid() bool {
	switch e {
	case IngressRuleStatusPending, IngressRuleStatusApplied, IngressRuleStatusDeleting, IngressRuleStatusFailed:
		return true
	}
	return false
}

func (e IngressRuleStatus) String() string {
	return string(e)
}

func (e *IngressRuleStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = IngressRuleStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid IngressRuleStatus", str)
	}
	return nil
}

func (e IngressRuleStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PersistentVolumeBackupStatus string

const (
	PersistentVolumeBackupStatusPending PersistentVolumeBackupStatus = "pending"
	PersistentVolumeBackupStatusFailed  PersistentVolumeBackupStatus = "failed"
	PersistentVolumeBackupStatusSuccess PersistentVolumeBackupStatus = "success"
)

var AllPersistentVolumeBackupStatus = []PersistentVolumeBackupStatus{
	PersistentVolumeBackupStatusPending,
	PersistentVolumeBackupStatusFailed,
	PersistentVolumeBackupStatusSuccess,
}

func (e PersistentVolumeBackupStatus) IsValid() bool {
	switch e {
	case PersistentVolumeBackupStatusPending, PersistentVolumeBackupStatusFailed, PersistentVolumeBackupStatusSuccess:
		return true
	}
	return false
}

func (e PersistentVolumeBackupStatus) String() string {
	return string(e)
}

func (e *PersistentVolumeBackupStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PersistentVolumeBackupStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PersistentVolumeBackupStatus", str)
	}
	return nil
}

func (e PersistentVolumeBackupStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PersistentVolumeBackupType string

const (
	PersistentVolumeBackupTypeLocal PersistentVolumeBackupType = "local"
	PersistentVolumeBackupTypeS3    PersistentVolumeBackupType = "s3"
)

var AllPersistentVolumeBackupType = []PersistentVolumeBackupType{
	PersistentVolumeBackupTypeLocal,
	PersistentVolumeBackupTypeS3,
}

func (e PersistentVolumeBackupType) IsValid() bool {
	switch e {
	case PersistentVolumeBackupTypeLocal, PersistentVolumeBackupTypeS3:
		return true
	}
	return false
}

func (e PersistentVolumeBackupType) String() string {
	return string(e)
}

func (e *PersistentVolumeBackupType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PersistentVolumeBackupType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PersistentVolumeBackupType", str)
	}
	return nil
}

func (e PersistentVolumeBackupType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PersistentVolumeRestoreStatus string

const (
	PersistentVolumeRestoreStatusPending PersistentVolumeRestoreStatus = "pending"
	PersistentVolumeRestoreStatusFailed  PersistentVolumeRestoreStatus = "failed"
	PersistentVolumeRestoreStatusSuccess PersistentVolumeRestoreStatus = "success"
)

var AllPersistentVolumeRestoreStatus = []PersistentVolumeRestoreStatus{
	PersistentVolumeRestoreStatusPending,
	PersistentVolumeRestoreStatusFailed,
	PersistentVolumeRestoreStatusSuccess,
}

func (e PersistentVolumeRestoreStatus) IsValid() bool {
	switch e {
	case PersistentVolumeRestoreStatusPending, PersistentVolumeRestoreStatusFailed, PersistentVolumeRestoreStatusSuccess:
		return true
	}
	return false
}

func (e PersistentVolumeRestoreStatus) String() string {
	return string(e)
}

func (e *PersistentVolumeRestoreStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PersistentVolumeRestoreStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PersistentVolumeRestoreStatus", str)
	}
	return nil
}

func (e PersistentVolumeRestoreStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PersistentVolumeRestoreType string

const (
	PersistentVolumeRestoreTypeLocal PersistentVolumeRestoreType = "local"
)

var AllPersistentVolumeRestoreType = []PersistentVolumeRestoreType{
	PersistentVolumeRestoreTypeLocal,
}

func (e PersistentVolumeRestoreType) IsValid() bool {
	switch e {
	case PersistentVolumeRestoreTypeLocal:
		return true
	}
	return false
}

func (e PersistentVolumeRestoreType) String() string {
	return string(e)
}

func (e *PersistentVolumeRestoreType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PersistentVolumeRestoreType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PersistentVolumeRestoreType", str)
	}
	return nil
}

func (e PersistentVolumeRestoreType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type PersistentVolumeType string

const (
	PersistentVolumeTypeLocal PersistentVolumeType = "local"
	PersistentVolumeTypeNfs   PersistentVolumeType = "nfs"
)

var AllPersistentVolumeType = []PersistentVolumeType{
	PersistentVolumeTypeLocal,
	PersistentVolumeTypeNfs,
}

func (e PersistentVolumeType) IsValid() bool {
	switch e {
	case PersistentVolumeTypeLocal, PersistentVolumeTypeNfs:
		return true
	}
	return false
}

func (e PersistentVolumeType) String() string {
	return string(e)
}

func (e *PersistentVolumeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = PersistentVolumeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid PersistentVolumeType", str)
	}
	return nil
}

func (e PersistentVolumeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProtocolType string

const (
	ProtocolTypeHTTP  ProtocolType = "http"
	ProtocolTypeHTTPS ProtocolType = "https"
	ProtocolTypeTCP   ProtocolType = "tcp"
	ProtocolTypeUDP   ProtocolType = "udp"
)

var AllProtocolType = []ProtocolType{
	ProtocolTypeHTTP,
	ProtocolTypeHTTPS,
	ProtocolTypeTCP,
	ProtocolTypeUDP,
}

func (e ProtocolType) IsValid() bool {
	switch e {
	case ProtocolTypeHTTP, ProtocolTypeHTTPS, ProtocolTypeTCP, ProtocolTypeUDP:
		return true
	}
	return false
}

func (e ProtocolType) String() string {
	return string(e)
}

func (e *ProtocolType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProtocolType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProtocolType", str)
	}
	return nil
}

func (e ProtocolType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ProxyType string

const (
	ProxyTypeBackup ProxyType = "backup"
	ProxyTypeActive ProxyType = "active"
)

var AllProxyType = []ProxyType{
	ProxyTypeBackup,
	ProxyTypeActive,
}

func (e ProxyType) IsValid() bool {
	switch e {
	case ProxyTypeBackup, ProxyTypeActive:
		return true
	}
	return false
}

func (e ProxyType) String() string {
	return string(e)
}

func (e *ProxyType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ProxyType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ProxyType", str)
	}
	return nil
}

func (e ProxyType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RedirectRuleStatus string

const (
	RedirectRuleStatusPending  RedirectRuleStatus = "pending"
	RedirectRuleStatusApplied  RedirectRuleStatus = "applied"
	RedirectRuleStatusFailed   RedirectRuleStatus = "failed"
	RedirectRuleStatusDeleting RedirectRuleStatus = "deleting"
)

var AllRedirectRuleStatus = []RedirectRuleStatus{
	RedirectRuleStatusPending,
	RedirectRuleStatusApplied,
	RedirectRuleStatusFailed,
	RedirectRuleStatusDeleting,
}

func (e RedirectRuleStatus) IsValid() bool {
	switch e {
	case RedirectRuleStatusPending, RedirectRuleStatusApplied, RedirectRuleStatusFailed, RedirectRuleStatusDeleting:
		return true
	}
	return false
}

func (e RedirectRuleStatus) String() string {
	return string(e)
}

func (e *RedirectRuleStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RedirectRuleStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RedirectRuleStatus", str)
	}
	return nil
}

func (e RedirectRuleStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ServerResourceAnalyticsTimeframe string

const (
	ServerResourceAnalyticsTimeframeLast1Hour   ServerResourceAnalyticsTimeframe = "last_1_hour"
	ServerResourceAnalyticsTimeframeLast24Hours ServerResourceAnalyticsTimeframe = "last_24_hours"
	ServerResourceAnalyticsTimeframeLast7Days   ServerResourceAnalyticsTimeframe = "last_7_days"
	ServerResourceAnalyticsTimeframeLast30Days  ServerResourceAnalyticsTimeframe = "last_30_days"
)

var AllServerResourceAnalyticsTimeframe = []ServerResourceAnalyticsTimeframe{
	ServerResourceAnalyticsTimeframeLast1Hour,
	ServerResourceAnalyticsTimeframeLast24Hours,
	ServerResourceAnalyticsTimeframeLast7Days,
	ServerResourceAnalyticsTimeframeLast30Days,
}

func (e ServerResourceAnalyticsTimeframe) IsValid() bool {
	switch e {
	case ServerResourceAnalyticsTimeframeLast1Hour, ServerResourceAnalyticsTimeframeLast24Hours, ServerResourceAnalyticsTimeframeLast7Days, ServerResourceAnalyticsTimeframeLast30Days:
		return true
	}
	return false
}

func (e ServerResourceAnalyticsTimeframe) String() string {
	return string(e)
}

func (e *ServerResourceAnalyticsTimeframe) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ServerResourceAnalyticsTimeframe(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ServerResourceAnalyticsTimeframe", str)
	}
	return nil
}

func (e ServerResourceAnalyticsTimeframe) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ServerStatus string

const (
	ServerStatusNeedsSetup ServerStatus = "needs_setup"
	ServerStatusPreparing  ServerStatus = "preparing"
	ServerStatusOnline     ServerStatus = "online"
	ServerStatusOffline    ServerStatus = "offline"
)

var AllServerStatus = []ServerStatus{
	ServerStatusNeedsSetup,
	ServerStatusPreparing,
	ServerStatusOnline,
	ServerStatusOffline,
}

func (e ServerStatus) IsValid() bool {
	switch e {
	case ServerStatusNeedsSetup, ServerStatusPreparing, ServerStatusOnline, ServerStatusOffline:
		return true
	}
	return false
}

func (e ServerStatus) String() string {
	return string(e)
}

func (e *ServerStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ServerStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ServerStatus", str)
	}
	return nil
}

func (e ServerStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SwarmMode string

const (
	SwarmModeManager SwarmMode = "manager"
	SwarmModeWorker  SwarmMode = "worker"
)

var AllSwarmMode = []SwarmMode{
	SwarmModeManager,
	SwarmModeWorker,
}

func (e SwarmMode) IsValid() bool {
	switch e {
	case SwarmModeManager, SwarmModeWorker:
		return true
	}
	return false
}

func (e SwarmMode) String() string {
	return string(e)
}

func (e *SwarmMode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SwarmMode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SwarmMode", str)
	}
	return nil
}

func (e SwarmMode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type UpstreamType string

const (
	UpstreamTypeGit        UpstreamType = "git"
	UpstreamTypeSourceCode UpstreamType = "sourceCode"
	UpstreamTypeImage      UpstreamType = "image"
)

var AllUpstreamType = []UpstreamType{
	UpstreamTypeGit,
	UpstreamTypeSourceCode,
	UpstreamTypeImage,
}

func (e UpstreamType) IsValid() bool {
	switch e {
	case UpstreamTypeGit, UpstreamTypeSourceCode, UpstreamTypeImage:
		return true
	}
	return false
}

func (e UpstreamType) String() string {
	return string(e)
}

func (e *UpstreamType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = UpstreamType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid UpstreamType", str)
	}
	return nil
}

func (e UpstreamType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
