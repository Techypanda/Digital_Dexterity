package lib

import (
	b64 "encoding/base64"
	"fmt"
	"os"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type APISecretsConfig struct {
	DatabaseUsername   string
	DatabasePassword   string
	DatabaseAddress    string
	GithubOAuthSecrets string
	mappedSecrets      map[string]string
}

func NewAPISecretsConfig(dbUsername string, dbPassword string, dbAddress string, githubOauth string) APISecretsConfig {
	return APISecretsConfig{
		DatabaseUsername:   dbUsername,
		DatabasePassword:   dbPassword,
		DatabaseAddress:    dbAddress,
		GithubOAuthSecrets: githubOauth,
		mappedSecrets: map[string]string{
			"db_username":  dbUsername,
			"db_password":  dbPassword,
			"db_address":   dbAddress,
			"github_oauth": githubOauth,
		},
	}
}

func LoadAPISecretsFromEnviron() (*APISecretsConfig, error) {
	dbUsername, found := os.LookupEnv("db_username")
	if !found {
		return nil, EnvironmentMisconfiguredError("undefined db_username")
	}

	dbPassword, found := os.LookupEnv("db_password")
	if !found {
		return nil, EnvironmentMisconfiguredError("undefined db_password")
	}

	dbAddress, found := os.LookupEnv("db_address")
	if !found {
		return nil, EnvironmentMisconfiguredError("undefined db_address")
	}

	githubOAuth, found := os.LookupEnv("github_oauth")
	if !found {
		return nil, EnvironmentMisconfiguredError("undefined github_oauth")
	}

	config := NewAPISecretsConfig(dbUsername, dbPassword, dbAddress, githubOAuth)

	return &config, nil
}

func NewAPISecrets(ctx *pulumi.Context, appLabels pulumi.StringMapInput, secretsConfig APISecretsConfig) error {
	_, err := corev1.NewSecret(ctx, "digitaldexapi-secrets", &corev1.SecretArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:   pulumi.String("digitaldexapi-secrets"),
			Labels: appLabels,
		},
		StringData: pulumi.ToStringMap(secretsConfig.mappedSecrets),
	})

	if err != nil {
		return fmt.Errorf("failed to create api secrets: %w", err)
	}

	return nil
}

type GithubSecretConfig struct {
	GHToken      string
	GHUsername   string
	mappedSecret map[string]string
}

func NewGithubSecretConfig(ghToken string, ghUsername string) GithubSecretConfig {
	unencodedUP := fmt.Sprintf("%s:%s", ghUsername, ghToken)
	encodedUP := b64.StdEncoding.EncodeToString([]byte(unencodedUP))
	unencodedJSON := fmt.Sprintf(`
	{
		"auths":
		{
				"ghcr.io":
						{
								"auth":"%s"
						}
		}
	}
	`, encodedUP)
	encodedJSON := b64.StdEncoding.EncodeToString([]byte(unencodedJSON))

	return GithubSecretConfig{
		GHToken:    ghToken,
		GHUsername: ghUsername,
		mappedSecret: map[string]string{
			".dockerconfigjson": encodedJSON,
		},
	}
}

func NewGithubSecret(ctx *pulumi.Context, appLabels pulumi.StringMapInput, config GithubSecretConfig) error {
	_, err := corev1.NewSecret(ctx, "digitaldex-github-secret", &corev1.SecretArgs{
		Type: pulumi.String("kubernetes.io/dockerconfigjson"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:   pulumi.String("dockerconfig-digitaldex-github-com"),
			Labels: appLabels,
		},
		Data: pulumi.ToStringMap(config.mappedSecret),
	})

	return fmt.Errorf("failed to create github secret: %w", err)
}

func LoadGithubSecretsFromEnviron() (*GithubSecretConfig, error) {
	ghUsername, found := os.LookupEnv("gh_username")
	if !found {
		return nil, EnvironmentMisconfiguredError("gh_username is undefined")
	}

	ghToken, found := os.LookupEnv("gh_token")
	if !found {
		return nil, EnvironmentMisconfiguredError("gh_token is undefined")
	}

	config := NewGithubSecretConfig(ghToken, ghUsername)

	return &config, nil
}
