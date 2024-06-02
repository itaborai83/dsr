package specs

import (
	"encoding/json"
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/config"
	"github.com/itaborai83/dsr/utils"
)

type SpecService struct {
	repo common.Repo
}

func RegisterServices() error {
	conf := config.GetConfig()
	err := utils.EnsureFolderExists(conf.DataDir, conf.SpecsDir)
	if err != nil {
		return fmt.Errorf("error ensuring spec directory exists: %v", err)
	}
	specsDir := conf.DataDir + "/" + conf.SpecsDir
	repo, err := common.NewFileRepo(specsDir, conf.SpecEntry)
	if err != nil {
		return fmt.Errorf("error creating spec repo: %v", err)
	}
	service := NewSpecService(repo)
	conf.SpecRepo = repo
	conf.SpecService = service
	return nil
}

func NewSpecService(repo common.Repo) *SpecService {
	return &SpecService{repo: repo}
}

func (s *SpecService) saveSpec(spec *TableSpec) error {
	now := utils.GetNow()
	if spec.CreatedAt == "" {
		spec.CreatedAt = now
		spec.UpdatedAt = now
		data, err := json.Marshal(spec)
		if err != nil {
			return fmt.Errorf("error marshalling spec: %v", err)
		}
		err = s.repo.CreateEntry(spec.Id, data)
		if err != nil {
			return fmt.Errorf("error creating spec: %v", err)
		}
	} else {
		spec.UpdatedAt = now
		data, err := json.Marshal(spec)
		if err != nil {
			return fmt.Errorf("error marshalling spec: %v", err)
		}
		err = s.repo.UpdateEntry(spec.Id, data)
		if err != nil {
			return fmt.Errorf("error updating spec: %v", err)
		}
	}
	return nil
}

func (s *SpecService) loadSpec(specId string) (*TableSpec, error) {
	data, err := s.repo.GetEntry(specId)
	if err != nil {
		return nil, fmt.Errorf("error getting spec: %v", err)
	}
	var spec TableSpec
	err = json.Unmarshal(data, &spec)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling spec: %v", err)
	}
	return &spec, nil
}

func (s *SpecService) DoesSpecExist(specId string) (bool, error) {
	return s.repo.DoesEntryExist(specId)
}

func (s *SpecService) GetSpec(specId string) (*TableSpec, error) {
	spec, err := s.loadSpec(specId)
	if err != nil {
		return nil, fmt.Errorf("error getting spec: %v", err)
	}
	return spec, nil
}

func (s *SpecService) validateSpecCreation(spec *TableSpec) error {
	err := ValidateTableSpec(spec)
	if err != nil {
		return fmt.Errorf("error validating spec: %v", err)
	}

	exists, err := s.DoesSpecExist(spec.Id)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if exists {
		return fmt.Errorf("spec already exists: %s", spec.Id)
	}

	return nil
}

func (s *SpecService) CreateSpec(spec *TableSpec) error {
	err := s.validateSpecCreation(spec)
	if err != nil {
		return err
	}
	err = s.saveSpec(spec)
	if err != nil {
		return fmt.Errorf("error saving spec: %v", err)
	}
	return nil
}

func (s *SpecService) validateSpecUpdate(specId string, spec *TableSpec) error {
	if specId != spec.Id {
		return fmt.Errorf("spec id mismatch")
	}

	err := ValidateTableSpec(spec)
	if err != nil {
		return fmt.Errorf("error validating spec: %v", err)
	}

	exists, err := s.DoesSpecExist(spec.Id)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", spec.Id)
	}

	return nil
}

func (s *SpecService) UpdateSpec(specId string, spec *TableSpec) error {
	err := s.validateSpecUpdate(specId, spec)
	if err != nil {
		return err
	}
	err = s.saveSpec(spec)
	if err != nil {
		return fmt.Errorf("error saving spec: %v", err)
	}
	return nil
}

func (s *SpecService) validateSpecDeletion(specId string) error {
	err := utils.ValidateId(specId)
	if err != nil {
		return fmt.Errorf("error validating spec id: %v", err)
	}
	exists, err := s.DoesSpecExist(specId)
	if err != nil {
		return fmt.Errorf("error checking if spec exists: %v", err)
	}
	if !exists {
		return fmt.Errorf("spec does not exist: %s", specId)
	}
	return nil
}

func (s *SpecService) DeleteSpec(specId string) error {
	err := s.validateSpecDeletion(specId)
	if err != nil {
		return err
	}

	err = s.repo.DeleteEntry(specId)
	if err != nil {
		return fmt.Errorf("error deleting spec: %v", err)
	}
	return nil
}

func (s *SpecService) ListSpecIds() ([]string, error) {
	ids, err := s.repo.ListEntryIds()
	if err != nil {
		return nil, fmt.Errorf("error listing spec ids: %v", err)
	}
	return ids, nil
}

func (s *SpecService) GetAllSpecs() ([]TableSpec, error) {
	ids, err := s.ListSpecIds()
	if err != nil {
		return nil, fmt.Errorf("error listing specs: %v", err)
	}
	specs := make([]TableSpec, 0)
	for _, id := range ids {
		spec, err := s.GetSpec(id)
		if err != nil {
			return nil, fmt.Errorf("error getting spec: %v", err)
		}
		specs = append(specs, *spec)
	}
	return specs, nil
}
