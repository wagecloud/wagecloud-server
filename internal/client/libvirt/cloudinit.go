package libvirt

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/kdomanski/iso9660"
	"gopkg.in/yaml.v3"
)

const cloudInitPrefix = "#cloud-config"

type CreateCloudinitParams struct {
	Filepath      string
	Userdata      Userdata
	Metadata      Metadata
	NetworkConfig NetworkConfig
}

func (s *ClientImpl) CreateCloudinit(ctx context.Context, params CreateCloudinitParams) error {
	cloudinitFile, err := os.Create(params.Filepath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %s", err)
	}

	// 1. Marshal userdata
	userdataYaml, err := yaml.Marshal(params.Userdata)
	if err != nil {
		return fmt.Errorf("failed to marshal userdata: %s", err)
	}
	// Add cloudinit prefix to userdata
	userdataReader := bytes.NewReader([]byte(cloudInitPrefix + "\n" + string(userdataYaml)))

	// 2. Marshal metadata
	metadataYaml, err := yaml.Marshal(params.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %s", err)
	}
	// Add cloudinit prefix to metadata
	metadataReader := bytes.NewReader(metadataYaml)

	// 3. Marshal network config
	networkConfigYaml, err := yaml.Marshal(params.NetworkConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal network config: %s", err)
	}
	networkConfigReader := bytes.NewReader(networkConfigYaml)

	if err = s.WriteCloudinit(ctx, userdataReader, metadataReader, networkConfigReader, cloudinitFile); err != nil {
		return fmt.Errorf("failed to write cloudinit ISO: %s", err)
	}

	return nil
}

type CreateCloudinitByReaderParams struct {
	Filepath      string
	Userdata      io.Reader
	Metadata      io.Reader
	NetworkConfig io.Reader
}

func (s *ClientImpl) CreateCloudinitByReader(ctx context.Context, params CreateCloudinitByReaderParams) error {
	cloudinitFile, err := os.Create(params.Filepath)
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %s", err)
	}

	if err = s.WriteCloudinit(ctx, params.Userdata, params.Metadata, params.NetworkConfig, cloudinitFile); err != nil {
		return fmt.Errorf("failed to write cloudinit ISO: %s", err)
	}

	return nil
}

func (s *ClientImpl) WriteCloudinit(ctx context.Context, userdata io.Reader, metadata io.Reader, networkConfig io.Reader, cloudinitFile io.Writer) error {
	writer, err := iso9660.NewWriter()
	if err != nil {
		return fmt.Errorf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	if err = writer.AddFile(userdata, "user-data"); err != nil {
		return fmt.Errorf("failed to add user-data: %s", err)
	}

	if err = writer.AddFile(metadata, "meta-data"); err != nil {
		return fmt.Errorf("failed to add meta-data: %s", err)
	}

	if err = writer.WriteTo(cloudinitFile, "cidata"); err != nil {
		return fmt.Errorf("failed to write ISO image: %s", err)
	}

	return nil
}
