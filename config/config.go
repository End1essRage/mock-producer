package config

const (
	ENV_DEV  = "ENV_DEV"  // for local deploy, mocking services
	ENV_TEST = "ENV_TEST" // for deploy on test no secrets
	ENV_PROD = "ENV_PROD" // for deploy with vault
)

type Config struct {
	ENV   string `env:"ENV"`
	DEBUG string `env:"DEBUG"`

	//common
	//rmq
	RMQ_ADDRESS      string `env:"RMQ_ADDRESS"`
	RMQ_VIRTUAL_HOST string `env:"RMQ_VIRTUAL_HOST"`
	RMQ_PORT         string `env:"RMQ_PORT"`
	RMQ_EXCHANGE     string `env:"RMQ_EXCHANGE"`
	//for dev
	RMQ_LOGIN string `env:"RMQ_LOGIN"`
	RMQ_PWD   string `env:"RMQ_PWD"`

	SECRET_KEY_RMQ_LOGIN string `env:"SECRET_KEY_RMQ_LOGIN"`
	SECRET_KEY_RMQ_PWD   string `env:"SECRET_KEY_RMQ_PWD"`

	//vault
	VAULT_SERVER      string `env:"VAULT_SERVER"`
	VAULT_SECRET_PATH string `env:"VAULT_SECRET_PATH"`
	VAULT_MOUNT_POINT string `env:"VAULT_MOUNT_POINT"`
	VAULT_ROLE_NAME   string `env:"VAULT_ROLE_NAME"`
}
