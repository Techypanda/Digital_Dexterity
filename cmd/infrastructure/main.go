package main

import (
	"fmt"
	"techytechster/digitaldexterity/lib"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var DefaultPort = 8080

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		appLabels := pulumi.StringMap{
			"app": pulumi.String("digitaldexterityapi"),
		}
		apiSecretConfig, err := lib.LoadAPISecretsFromEnviron()
		if err != nil {
			return fmt.Errorf("failed to create load api secrets from env: %w", err)
		}
		err = lib.NewAPISecrets(ctx, appLabels, *apiSecretConfig)
		if err != nil {
			return fmt.Errorf("failed to create api secrets: %w", err)
		}
		githubSecretConfig, err := lib.LoadGithubSecretsFromEnviron()
		if err != nil {
			return fmt.Errorf("failed to load github secret config: %w", err)
		}
		err = lib.NewGithubSecret(ctx, appLabels, *githubSecretConfig)
		if err != nil {
			return fmt.Errorf("failed to create github secret: %w", err)
		}
		config, err := lib.LoadDeploymentConfigFromEnviron()
		if err != nil {
			return fmt.Errorf("failed to load deployment config from env: %w", err)
		}
		err = lib.NewAPIDeployment(ctx, appLabels, *config)
		if err != nil {
			return fmt.Errorf("failed to create api deployment: %w", err)
		}
		if err = lib.NewAPIService(ctx, DefaultPort); err != nil {
			return fmt.Errorf("failed to create api service: %w", err)
		}

		return nil
	})
}
