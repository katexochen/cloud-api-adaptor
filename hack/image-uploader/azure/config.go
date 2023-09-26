package azure

type Config struct {
	SubscriptionID         string `yaml:"subscriptionID"`
	Location               string `yaml:"location"`
	ResourceGroup          string `yaml:"resourceGroup"`
	SharedImageGalleryName string `yaml:"sharedImageGallery"`
	ImageDefinitionName    string `yaml:"imageDefinition"`
	ImageVersionName       string `yaml:"imageVersion"`
	ImageOffer             string `yaml:"imageOffer"`
	ImageSKU               string `yaml:"imageSKU"`
	ImagePublisher         string `yaml:"imagePublisher"`
	DiskName               string `yaml:"diskName"`
	AttestationVariant     string `yaml:"attestationVariant"`
}
