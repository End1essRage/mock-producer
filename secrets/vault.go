package secrets

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
	"gitlab.gitlab.bcs.ru/elma365/mock-producer/config"
)

const ServiceAccountTokenPath = "../var/run/secrets/kubernetes.io/serviceaccount/token"

type Credentials struct {
	RmqLogin string
	RmqPwd   string
}

type Vault struct {
	cfg *config.Config
}

func NewVault(cfg *config.Config) *Vault {
	return &Vault{cfg: cfg}
}

func (v *Vault) GetCredentials() (Credentials, error) {
	creds := Credentials{}

	//set config
	config := vault.DefaultConfig()
	config.Address = v.cfg.VAULT_SERVER

	//create client
	client, err := vault.NewClient(config)
	if err != nil {
		return creds, fmt.Errorf("ошибка инициализации клиента: %w", err)
	}

	//set auth
	k8sAuth, err := auth.NewKubernetesAuth(
		v.cfg.VAULT_ROLE_NAME,
		auth.WithServiceAccountTokenPath(ServiceAccountTokenPath),
	)
	if err != nil {
		return creds, fmt.Errorf("ошибка инициализации kubernetes аутентификации: %w", err)
	}

	authInfo, err := client.Auth().Login(context.TODO(), k8sAuth)
	if err != nil {
		return creds, fmt.Errorf("ошибка входа по keubernetes auth: %w", err)
	}
	if authInfo == nil {
		return creds, fmt.Errorf("ошибка, пустая authInfo")
	}
	//get secret
	secret, err := client.KVv2(v.cfg.VAULT_MOUNT_POINT).Get(context.Background(), v.cfg.VAULT_SECRET_PATH)
	if err != nil {
		return creds, fmt.Errorf("ошибка получения секрета: %w", err)
	}

	//extrude data
	return extrudeDataFromSecret(v.cfg, secret)
}

func extrudeDataFromSecret(cfg *config.Config, secret *vault.KVSecret) (Credentials, error) {
	creds := Credentials{}

	rlVal, ok := secret.Data[cfg.SECRET_KEY_RMQ_LOGIN].(string)
	if !ok {
		return creds, fmt.Errorf("value type assertion failed: %T %#v", secret.Data[cfg.SECRET_KEY_RMQ_LOGIN], secret.Data[cfg.SECRET_KEY_RMQ_LOGIN])
	}
	creds.RmqLogin = rlVal

	rpVal, ok := secret.Data[cfg.SECRET_KEY_RMQ_PWD].(string)
	if !ok {
		return creds, fmt.Errorf("value type assertion failed: %T %#v", secret.Data[cfg.SECRET_KEY_RMQ_PWD], secret.Data[cfg.SECRET_KEY_RMQ_PWD])
	}
	creds.RmqPwd = rpVal

	return creds, nil
}
