package lib

import (
	"errors"
	"os"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type APIDeploymentConfiguration struct {
	Name         string
	Image        string
	SecretsStore string
	CorsList     string
	SecretKey    string
}

func NewAPIDeploymentConfiguration(name string, image string, secretStore string, corsList string, secretKey string) APIDeploymentConfiguration {
	return APIDeploymentConfiguration{
		Name:         name,
		Image:        image,
		SecretsStore: secretStore,
		CorsList:     corsList,
		SecretKey:    secretKey,
	}
}

func LoadDeploymentConfigFromEnviron() (*APIDeploymentConfiguration, error) {
	appName, found := os.LookupEnv("app_name")
	if !found {
		return nil, errors.New("app_name is undefined")
	}
	image, found := os.LookupEnv("image")
	if !found {
		return nil, errors.New("image is undefined")
	}
	secretsStore, found := os.LookupEnv("secrets_store")
	if !found {
		return nil, errors.New("secrets_store is undefined")
	}
	corsList, found := os.LookupEnv("cors_list")
	if !found {
		return nil, errors.New("cors_list is undefined")
	}
	secretKey, found := os.LookupEnv("secret_key")
	if !found {
		return nil, errors.New("secret_key is undefined")
	}
	config := NewAPIDeploymentConfiguration(appName, image, secretsStore, corsList, secretKey)
	return &config, nil
}

func NewAPIDeployment(ctx *pulumi.Context, appLabels pulumi.StringMap, config APIDeploymentConfiguration) error {
	_, err := appsv1.NewDeployment(ctx, "digitaldexapi-deployment", &appsv1.DeploymentArgs{
		Spec: appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: appLabels,
			},
			Replicas: pulumi.Int(1),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: appLabels,
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						corev1.ContainerArgs{
							Name:            pulumi.String(config.Name),
							Image:           pulumi.String(config.Image),
							ImagePullPolicy: pulumi.String("Always"),
							Ports: &corev1.ContainerPortArray{
								corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8080),
								},
							},
							Env: &corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name: pulumi.String("db_username"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										SecretKeyRef: &corev1.SecretKeySelectorArgs{
											Name: pulumi.String(config.SecretsStore),
											Key:  pulumi.String("db_username"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name: pulumi.String("db_password"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										SecretKeyRef: &corev1.SecretKeySelectorArgs{
											Name: pulumi.String(config.SecretsStore),
											Key:  pulumi.String("db_password"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name: pulumi.String("db_address"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										SecretKeyRef: &corev1.SecretKeySelectorArgs{
											Name: pulumi.String(config.SecretsStore),
											Key:  pulumi.String("db_address"),
										},
									},
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("cors_list"),
									Value: pulumi.String(config.CorsList),
								},
								&corev1.EnvVarArgs{ // TODO: make this a actual secret or remove it if its deprecated
									Name:  pulumi.String("secret_key"),
									Value: pulumi.String(config.SecretKey),
								},
							},
						},
					},
					ImagePullSecrets: &corev1.LocalObjectReferenceArray{
						&corev1.LocalObjectReferenceArgs{
							Name: pulumi.String("dockerconfig-digitaldex-github-com"),
						},
					},
				},
			},
		},
	})
	return err
}
