package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/confidential-containers/cloud-api-adaptor/hack/image-uploader/azure"
	"github.com/confidential-containers/cloud-api-adaptor/hack/image-uploader/uploader"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	configName      = "image-uploader.yml"
	configDir       = "image-uploader.targets"
	timestampFormat = "20060102150405"
)

func newCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "image-uploader",
		Short:            "image-uploader is a tool for uploading images to a cloud provider",
		PersistentPreRun: preRunRoot,
		RunE:             run,
		Args:             cobra.MatchAll(cobra.ExactArgs(2), isCSP(0)),
	}
	cmd.SetOut(os.Stdout)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	provider := args[0]
	imagePath := args[1]
	logger := log.New(cmd.OutOrStderr(), "", log.LstdFlags)

	configFile, err := os.Open(configName)
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer configFile.Close()
	var config uploader.Config
	if err := yaml.NewDecoder(configFile).Decode(&config); err != nil {
		return fmt.Errorf("decoding config file: %w", err)
	}

	var providerConfig, providerOverwrite io.ReadCloser
	providerConfig, err = os.Open(path.Join(configDir, args[0]+".yml"))
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer providerConfig.Close()
	if overwrite, err := os.Open(path.Join(configDir, args[0]+".overwrite.yml")); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("opening override config file: %w", err)
		}
	} else {
		providerOverwrite = overwrite
		defer providerOverwrite.Close()
	}

	var prepper Prepper
	var upload Uploader

	switch provider {
	case "azure":
		prepper = &azure.Prepper{}
		upload, err = azure.NewUploader(providerConfig, providerOverwrite, logger)
		if err != nil {
			return fmt.Errorf("creating azure uploader: %w", err)
		}
	}

	imagePath, err = prepper.Prepare(cmd.Context(), imagePath)
	if err != nil {
		return fmt.Errorf("preparing image: %w", err)
	}
	image, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("opening image: %w", err)
	}

	req := &uploader.Request{
		Config:    config,
		Image:     image,
		Timestamp: time.Now().UTC().Format(timestampFormat),
	}
	ref, err := upload.Upload(cmd.Context(), req)
	if err != nil {
		return fmt.Errorf("uploading image: %w", err)
	}

	fmt.Println(ref)
	return nil
}

func supportedCSPs() []string {
	return []string{"azure"}
}

func isCSP(position int) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, csp := range supportedCSPs() {
			if args[position] == csp {
				return nil
			}
		}
		return fmt.Errorf("unsupported cloud service provider: %s", args[position])
	}
}

type Prepper interface {
	Prepare(ctx context.Context, imagePath string) (string, error)
}

type Uploader interface {
	Upload(ctx context.Context, req *uploader.Request) (ref string, retErr error)
}
