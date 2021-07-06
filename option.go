package dlock

import "time"

// Options external option
type Options struct {
	// lock type: mysql/redis/etcd/zk
	Type string

	// common option
	// redis or ectd password
	Password string

	// mysql lock option
	IP   string
	Name string
	User string
	Port int64

	// cluster model ip address
	Cluster []string
	// connection timeout
	DialTimeout time.Duration

	// Tls Configs
	// cert file path
	CertFile string
	// cert key file path
	KeyFile string
	// ca file path
	CAFile string
	// skip https
	SkipSSL bool
}

const (
	// DefaultMysqlPort 3306
	DefaultMysqlPort = 3306
	// DefaultRedisPort 6379
	DefaultRedisPort = 6379
)

// WithDBOption setting database options
func WithDBOption(user, password, ip, database string, port int64) func(*Options) {
	if port <= 0 {
		port = DefaultMysqlPort
	}
	return func(opts *Options) {
		opts.User = user
		opts.Password = password
		opts.IP = ip
		opts.Name = database
		opts.Port = port
		opts.Type = MysqlLockType
	}
}

// WithRedisOption setting redis options
// any other options?
func WithRedisOption(password string, DialTimeout time.Duration, cluster ...string) func(*Options) {
	return func(opts *Options) {
		opts.Password = password
		opts.Cluster = cluster
		opts.DialTimeout = DialTimeout * time.Millisecond
		opts.Type = RedisLockType
	}
}

// WithEtcdOption setting etcd options
func WithEtcdOption(dialTimeout time.Duration, endpoints ...string) func(*Options) {
	return func(opts *Options) {
		opts.Cluster = endpoints
		opts.DialTimeout = dialTimeout * time.Millisecond
		opts.Type = EtcdLockType
	}
}

// WithEtcdAuthOption setting etcd options
func WithEtcdAuthOption(user, password, caFile, certFile, keyFile string, skipSSL bool) func(*Options) {
	return func(opts *Options) {
		opts.User = user
		opts.Password = password
		opts.Type = EtcdLockType
		opts.CAFile = caFile
		opts.CertFile = certFile
		opts.KeyFile = keyFile
		opts.SkipSSL = skipSSL
	}
}
