package specs

import (
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/dirrepo"
	"github.com/itaborai83/dsr/utils"
)

type SpecService struct {
	repo      common.SpecRepo
	validator common.SpecValidator
}

func RegisterServices() error {
	conf := common.GetConfig()
	err := utils.EnsureFolderExists(conf.DataDir, conf.SpecsDir)
	if err != nil {
		return fmt.Errorf("error ensuring data dir exists: %v", err)
	}
	dirrepo, err := dirrepo.NewDirRepo(conf.DataDir, conf.SpecsDir, conf.SpecEntry, nil)
	if err != nil {
		return fmt.Errorf("error creating spec repo: %v", err)
	}
	repo := NewSpecRepo(dirrepo)
	validator := NewSpecValidator(repo)
	service := NewSpecService(repo, validator)

	conf.SpecRepo = repo
	conf.SpecValidator = validator
	conf.SpecService = service
	return nil
}

func NewSpecService(repo common.SpecRepo, validator common.SpecValidator) *SpecService {
	return &SpecService{repo: repo, validator: validator}
}

func (s *SpecService) DoesSpecExist(specId string) (bool, error) {
	exists, err := s.repo.DoesSpecExist(specId)
	if err != nil {
		return false, fmt.Errorf("error checking if spec exists: %v", err)
	}
	return exists, nil
}

func (s *SpecService) GetSpec(specId string) (*common.TableSpec, error) {
	spec, err := s.repo.GetSpec(specId)
	if err != nil {
		return nil, fmt.Errorf("error getting spec: %v", err)
	}
	return spec, nil
}

func (s *SpecService) CreateSpec(spec *common.TableSpec) error {
	err := s.validator.ValidateCreation(spec)
	if err != nil {
		return err
	}
	err = s.repo.CreateSpec(spec.Id, spec)
	if err != nil {
		return fmt.Errorf("error saving spec: %v", err)
	}
	return nil
}

func (s *SpecService) UpdateSpec(specId string, spec *common.TableSpec) error {
	err := s.validator.ValidateUpdate(specId, spec)
	if err != nil {
		return err
	}
	err = s.repo.UpdateSpec(specId, spec)
	if err != nil {
		return fmt.Errorf("error updating spec: %v", err)
	}
	return nil
}

func (s *SpecService) DeleteSpec(specId string) error {
	err := s.validator.ValidateDeletion(specId)
	if err != nil {
		return err
	}
	err = s.repo.DeleteSpec(specId)
	if err != nil {
		return fmt.Errorf("error deleting spec: %v", err)
	}
	return nil
}

func (s *SpecService) ListSpecIds() ([]string, error) {
	ids, err := s.repo.ListSpecIds()
	if err != nil {
		return nil, fmt.Errorf("error listing spec ids: %v", err)
	}
	return ids, nil
}

func (s *SpecService) GetAllSpecs() ([]common.TableSpec, error) {
	ids, err := s.ListSpecIds()
	if err != nil {
		return nil, fmt.Errorf("error getting all specs: %v", err)
	}
	utils.GetLogger().Printf("ids: %v", ids)
	var specs []common.TableSpec
	for _, id := range ids {
		spec, err := s.GetSpec(id)
		if err != nil {
			return nil, fmt.Errorf("error getting spec: %v", err)
		}
		specs = append(specs, *spec)
	}
	return specs, nil
}
