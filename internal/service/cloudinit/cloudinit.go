package cloudinit

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

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

type CreateCloudinitParams struct {
	Userdata model.Userdata
	Metadata model.Metadata
}

func (s *Service) CreateCloudinit(params CreateCloudinitParams) (io.Reader, error) {
	userdata, err := yaml.Marshal(params.Userdata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal userdata: %s", err)
	}
	userdataReader := bytes.NewReader(userdata)

	metadata, err := yaml.Marshal(params.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %s", err)
	}
	metadataReader := bytes.NewReader(metadata)

	iso, err := os.CreateTemp("tmp", "cloudinit_*"+".iso")
	if err != nil {
		return nil, fmt.Errorf("failed to create ISO image: %s", err)
	}

	writer, err := iso9660.NewWriter()
	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	if err = writer.AddFile(userdataReader, "user-data"); err != nil {
		return nil, fmt.Errorf("failed to add user-data: %s", err)
	}

	if err = writer.AddFile(metadataReader, "meta-data"); err != nil {
		return nil, fmt.Errorf("failed to add meta-data: %s", err)
	}

	if err = writer.WriteTo(iso, "cidata"); err != nil {
		return nil, fmt.Errorf("failed to write ISO image: %s", err)
	}

	log.Printf("ISO image created successfully")
	return iso, nil
}
