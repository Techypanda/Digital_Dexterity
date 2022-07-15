package lib

import (
	"fmt"

	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewAPIService(ctx *pulumi.Context, port int) error {
	_, err := corev1.NewService(ctx, "digitaldexapi-service", &corev1.ServiceArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name: pulumi.String("digitaldexapi-service"),
		},
		Spec: &corev1.ServiceSpecArgs{
			Selector: pulumi.ToStringMap(map[string]string{
				"app": "digitaldexterityapi",
			}),
			Type: pulumi.String("NodePort"),
			Ports: &corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port: pulumi.Int(port),
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to create api service: %w", err)
	}

	return nil
}
