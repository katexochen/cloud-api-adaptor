package uploader

import (
	"io"
)

type CommonConfig struct{}

type Config struct {
	// CommonConfig `yaml:",inline"`
	ImageVersion string      `yaml:"imageVersion"`
	Name         string      `yaml:"name"`
	Azure        AzureConfig `yaml:"azure"`
}

type AzureConfig struct {
	SubscriptionID         string `yaml:"subscriptionID"`
	Location               string `yaml:"location"`
	ResourceGroup          string `yaml:"resourceGroup"`
	AttestationVariant     string `yaml:"attestationVariant"`
	SharedImageGalleryName string `yaml:"sharedImageGallery"`
	ImageDefinitionName    string `yaml:"imageDefinition"`
	ImageVersionName       string `yaml:"imageVersion"`
	ImageOffer             string `yaml:"imageOffer"`
	ImageSKU               string `yaml:"imageSKU"`
	ImagePublisher         string `yaml:"imagePublisher"`
	DiskName               string `yaml:"diskName"`
}

type Request struct {
	Image     io.ReadSeekCloser
	Timestamp string
	Size      int64
}
