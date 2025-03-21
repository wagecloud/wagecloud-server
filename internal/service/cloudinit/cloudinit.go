package cloudinit

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kdomanski/iso9660"
	"github.com/wagecloud/wagecloud-server/internal/model"
	"github.com/wagecloud/wagecloud-server/internal/repository"
	"gopkg.in/yaml.v3"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCloudinit(
	cloudinitFile io.Writer,
	userdata model.Userdata,
	metadata model.Metadata,
) error {
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

	writer, err := iso9660.NewWriter()
	if err != nil {
		return fmt.Errorf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	if err = writer.AddFile(userdataReader, "user-data"); err != nil {
		return fmt.Errorf("failed to add user-data: %s", err)
	}

	if err = writer.AddFile(metadataReader, "meta-data"); err != nil {
		return fmt.Errorf("failed to add meta-data: %s", err)
	}

	if err = writer.WriteTo(cloudinitFile, "cidata"); err != nil {
		return fmt.Errorf("failed to write ISO image: %s", err)
	}

	return nil
}
