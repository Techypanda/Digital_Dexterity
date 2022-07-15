package lib

import (
	"fmt"
	"os"
	"strconv"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type APIDeploymentConfiguration struct {
	Name          string
	Image         string
	SecretsStore  string
	CorsList      string
	SecretKey     string
	ContainerPort int
}

func NewAPIDeploymentConfiguration(name string, image string, secretStore string, corsList string, secretKey string, containerPort int) APIDeploymentConfiguration {
	return APIDeploymentConfiguration{
		Name:          name,
		Image:         image,
		SecretsStore:  secretStore,
		CorsList:      corsList,
		SecretKey:     secretKey,
		ContainerPort: containerPort,
	}
}

var FallbackPort = 8080

func LoadDeploymentConfigFromEnviron() (*APIDeploymentConfiguration, error) {
	var appName, image, secretsStore, corsList, secretKey, containerPortEnv string

	var found bool

	containerPort := FallbackPort

	if appName, found = os.LookupEnv("app_name"); !found {
		return nil, EnvironmentMisconfiguredError("app_name is undefined")
	}

	if image, found = os.LookupEnv("image"); !found {
		return nil, EnvironmentMisconfiguredError("image is undefined")
	}

	if secretsStore, found = os.LookupEnv("secrets_store"); !found {
		return nil, EnvironmentMisconfiguredError("secrets_store is undefined")
	}

	if corsList, found = os.LookupEnv("cors_list"); !found {
		return nil, EnvironmentMisconfiguredError("cors_list is undefined")
	}

	if secretKey, found = os.LookupEnv("secret_key"); !found {
		return nil, EnvironmentMisconfiguredError("secret_key is undefined")
	}

	if containerPortEnv, found = os.LookupEnv("container_port"); found {
		if tmp, err := strconv.Atoi(containerPortEnv); err == nil {
			containerPort = tmp
		} else {
			return nil, ErrNotAValidPort
		}
	}

	config := NewAPIDeploymentConfiguration(appName, image, secretsStore, corsList, secretKey, containerPort)

	return &config, nil
}

func NewAPIDeployment(ctx *pulumi.Context, appLabels pulumi.StringMapInput, config APIDeploymentConfiguration) error {
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
									ContainerPort: pulumi.Int(config.ContainerPort),
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
									Name: pulumi.String("github_oauth"),
									ValueFrom: &corev1.EnvVarSourceArgs{
										SecretKeyRef: &corev1.SecretKeySelectorArgs{
											Name: pulumi.String(config.SecretsStore),
											Key:  pulumi.String("github_oauth"),
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

	if err != nil {
		return fmt.Errorf("failed to deploy api: %w", err)
	}

	return nil
}
