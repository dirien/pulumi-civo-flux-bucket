package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-civo/sdk/v2/go/civo"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		awsConfig := config.New(ctx, "aws")

		bucket, err := s3.NewBucket(ctx, "flux-source-bucket", &s3.BucketArgs{
			Acl: pulumi.String("private"),
			Tags: pulumi.StringMap{
				"Name": pulumi.String("flux-source-bucket"),
			},
			Bucket:       pulumi.StringPtr("flux-source-bucket"),
			ForceDestroy: pulumi.Bool(true),
		})

		firewall, err := civo.NewFirewall(ctx, "civo-firewall", &civo.FirewallArgs{
			Name:   pulumi.StringPtr("civo-firewall"),
			Region: pulumi.StringPtr("FRA1"),
		})
		if err != nil {
			return err
		}

		cluster, err := civo.NewKubernetesCluster(ctx, "civo-k3s-cluster", &civo.KubernetesClusterArgs{
			Name: pulumi.StringPtr("civo-k3s-cluster"),
			Pools: civo.KubernetesClusterPoolsArgs{
				Size:      pulumi.String("g4s.kube.medium"),
				NodeCount: pulumi.Int(3),
			},
			Region:     pulumi.StringPtr("FRA1"),
			FirewallId: firewall.ID(),
		})
		if err != nil {
			return err
		}

		ctx.Export("kubeconfig", pulumi.ToSecret(cluster.Kubeconfig))
		ctx.Export("accessKey", awsConfig.GetSecret("accessKey"))
		ctx.Export("secretKey", awsConfig.GetSecret("secretKey"))
		ctx.Export("bucket", bucket.Bucket)
		ctx.Export("bucket-region", bucket.Region)

		return nil
	})
}
