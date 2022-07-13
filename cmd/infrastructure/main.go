package main

import (
	"techytechster/digitaldexterity/lib"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		appLabels := pulumi.StringMap{
			"app": pulumi.String("digitaldexterityapi"),
		}
		apiSecretConfig, err := lib.LoadAPISecretsFromEnviron()
		if err != nil {
			return err
		}
		err = lib.NewAPISecrets(ctx, appLabels, *apiSecretConfig)
		if err != nil {
			return err
		}
		githubSecretConfig, err := lib.LoadGithubSecretsFromEnviron()
		if err != nil {
			return err
		}
		err = lib.NewGithubSecret(ctx, appLabels, *githubSecretConfig)
		if err != nil {
			return err
		}
		config, err := lib.LoadDeploymentConfigFromEnviron()
		if err != nil {
			return err
		}
		err = lib.NewAPIDeployment(ctx, appLabels, *config)
		if err != nil {
			return err
		}
		err = lib.NewAPIService(ctx)
		return err
	})
}
