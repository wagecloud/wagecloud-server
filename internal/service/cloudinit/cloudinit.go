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

type FinalUserData struct {
	Users []model.Userdata `yaml:"users"`
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

	var finalUserData FinalUserData = FinalUserData{
		Users: []model.Userdata{userdata},
	}

	finalUserDataReader, err := createFinalUserDataReader(finalUserData)
	if err != nil {
		return fmt.Errorf("failed to create cloudinit reader: %s", err)
	}

	metadataReader, err := createMetadataReader(metadata)
	if err != nil {
		return fmt.Errorf("failed to create metadata reader: %s", err)
	}

	networkConfigReader, err := createNetworkConfigReader(networkConfig)
	if err != nil {
		return fmt.Errorf("failed to create network config reader: %s", err)
	}

	if err = s.WriteCloudinit(cloudinitFile, finalUserDataReader, metadataReader, networkConfigReader); err != nil {
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

	if err = writer.WriteTo(cloudinitFile, "cidata"); err != nil {
		return fmt.Errorf("failed to write ISO image: %s", err)
	}

	return nil
}

func createFinalUserDataReader(final FinalUserData) (*bytes.Reader, error) {
	cloudInitYaml, err := yaml.Marshal(final)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal userdata: %s", err)
	}

	comment := "#cloud-config\n" // this is required for
	fullYaml := []byte(comment + string(cloudInitYaml))

	return bytes.NewReader(fullYaml), nil
}

func createUserDataReader(userdata model.Userdata) (*bytes.Reader, error) {
	cloudinitModel := FinalUserData{
		Users: []model.Userdata{userdata},
	}

	cloudInitYaml, err := yaml.Marshal(cloudinitModel)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal userdata: %s", err)
	}

	comment := "#cloud-config\n" // this is required for
	fullYaml := []byte(comment + string(cloudInitYaml))

	fmt.Println(string(fullYaml))

	return bytes.NewReader(fullYaml), nil
}

func createMetadataReader(metadata model.Metadata) (*bytes.Reader, error) {
	metadataYaml, err := yaml.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %s", err)
	}

	return bytes.NewReader(metadataYaml), nil
}

func createNetworkConfigReader(networkConfig model.NetworkConfig) (*bytes.Reader, error) {
	networkConfigYaml, err := yaml.Marshal(networkConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal network config: %s", err)
	}

	return bytes.NewReader(networkConfigYaml), nil
}
