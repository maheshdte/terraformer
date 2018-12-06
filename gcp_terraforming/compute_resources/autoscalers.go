// AUTO-GENERATED CODE. DO NOT EDIT.
package computeTerrforming

import (
	"context"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"waze/terraform/gcp_terraforming/gcp_generator"
	"waze/terraform/terraform_utils"
)

var autoscalersIgnoreKey = map[string]bool{
	"^id$":                 true,
	"^self_link$":          true,
	"^fingerprint$":        true,
	"^label_fingerprint$":  true,
	"^creation_timestamp$": true,
}

var autoscalersAllowEmptyValues = map[string]bool{}

var autoscalersAdditionalFields = map[string]string{
	"project": os.Getenv("GOOGLE_CLOUD_PROJECT"),
}

type AutoscalersGenerator struct {
	gcp_generator.BasicGenerator
}

// Run on autoscalersList and create for each TerraformResource
func (AutoscalersGenerator) createResources(autoscalersList *compute.AutoscalersListCall, ctx context.Context, region, zone string) []terraform_utils.TerraformResource {
	resources := []terraform_utils.TerraformResource{}
	if err := autoscalersList.Pages(ctx, func(page *compute.AutoscalerList) error {
		for _, obj := range page.Items {
			resources = append(resources, terraform_utils.NewTerraformResource(
				zone+"/"+obj.Name,
				obj.Name,
				"google_compute_autoscaler",
				"google",
				nil,
				map[string]string{
					"name":    obj.Name,
					"project": os.Getenv("GOOGLE_CLOUD_PROJECT"),
					"region":  region,
					"zone":    zone,
				},
			))
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
	return resources
}

// Generate TerraformResources from GCP API,
// from each autoscalers create 1 TerraformResource
// Need autoscalers name as ID for terraform resource
func (g AutoscalersGenerator) Generate(zone string) ([]terraform_utils.TerraformResource, map[string]terraform_utils.ResourceMetaData, error) {
	region := strings.Join(strings.Split(zone, "-")[:len(strings.Split(zone, "-"))-1], "-")
	project := os.Getenv("GOOGLE_CLOUD_PROJECT")
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	autoscalersList := computeService.Autoscalers.List(project, zone)

	resources := g.createResources(autoscalersList, ctx, region, zone)
	metadata := terraform_utils.NewResourcesMetaData(resources, autoscalersIgnoreKey, autoscalersAllowEmptyValues, autoscalersAdditionalFields)
	return resources, metadata, nil

}