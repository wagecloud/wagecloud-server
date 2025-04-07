package cloudinit

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/kdomanski/iso9660"
	"github.com/wagecloud/wagecloud-server/config"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"gopkg.in/yaml.v3"
)

type Service struct {
	repo *repository.RepositoryImpl
}

type ServiceInterface interface {
	CreateCloudinit(
		filename string,
		userdata model.Userdata,
		metadata model.Metadata,
		networkConfig model.NetworkConfig,
	) error
	CreateCloudinitByReader(
		filename string,
		userdata io.Reader,
		metadata io.Reader,
		networkConfig io.Reader,
	) error
	WriteCloudinit(
		cloudinitFile io.Writer,
		userdata io.Reader,
		metadata io.Reader,
		networkConfig io.Reader,
	) error
}

var _ ServiceInterface = (*Service)(nil)

func NewService(repo *repository.RepositoryImpl) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCloudinit(
	filename string,
	userdata model.Userdata,
	metadata model.Metadata,
	networkConfig model.NetworkConfig,
) error {
	cloudinitFile, err := os.Create(path.Join(
		config.GetConfig().App.CloudinitDir,
		filename,
	))
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %s", err)
	}

	userdataYaml, err := yaml.Marshal(userdata)
	if err != nil {
		return fmt.Errorf("failed to marshal userdata: %s", err)
	}
	userdataReader := bytes.NewReader(userdataYaml)

	metadataYaml, err := yaml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %s", err)
	}
	metadataReader := bytes.NewReader(metadataYaml)

	networkConfigYaml, err := yaml.Marshal(networkConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal network config: %s", err)
	}
	networkConfigReader := bytes.NewReader(networkConfigYaml)

	if err = s.WriteCloudinit(cloudinitFile, userdataReader, metadataReader, networkConfigReader); err != nil {
		return fmt.Errorf("failed to write cloudinit ISO: %s", err)
	}

	return nil
}

func (s *Service) CreateCloudinitByReader(
	filename string,
	userdata io.Reader,
	metadata io.Reader,
	networkConfig io.Reader,
) error {
	cloudinitFile, err := os.Create(path.Join(
		config.GetConfig().App.CloudinitDir,
		filename,
	))
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %s", err)
	}

	if err = s.WriteCloudinit(cloudinitFile, userdata, metadata, networkConfig); err != nil {
		return fmt.Errorf("failed to write cloudinit ISO: %s", err)
	}

	return nil
}

func (s *Service) WriteCloudinit(
	cloudinitFile io.Writer,
	userdata io.Reader,
	metadata io.Reader,
	networkConfig io.Reader,
) error {
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

	if err = writer.AddFile(networkConfig, "network-config"); err != nil {
		return fmt.Errorf("failed to add network-config: %s", err)
	}

	if err = writer.WriteTo(cloudinitFile, "cidata"); err != nil {
		return fmt.Errorf("failed to write ISO image: %s", err)
	}

	return nil
}
