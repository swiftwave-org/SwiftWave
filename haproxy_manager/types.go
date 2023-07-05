package haproxymanager

type HAProxySocket struct {
	Host string
	Port int
	isUnix bool // If true, Host is a path to a unix socket
	unixSocketPath string
	username string
	password string
}

type QueryParameter struct {
	key string
	value string
}

type QueryParameters []QueryParameter

type FrontendType string

const (
	HTTPFrontend FrontendType = "http"
	HTTPSFrontend FrontendType = "https"
	TCPFrontend FrontendType = "tcp"
)