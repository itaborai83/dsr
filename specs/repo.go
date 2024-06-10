package specs

import (
	"encoding/json"
	"fmt"

	"github.com/itaborai83/dsr/common"
	"github.com/itaborai83/dsr/utils"
)

type SpecRepo struct {
	repo common.Repo
}

func NewSpecRepo(repo common.Repo) *SpecRepo {
	return &SpecRepo{repo: repo}
}

func (s SpecRepo) DoesSpecExist(id string) (bool, error) {
	return s.repo.DoesEntryExist([]string{id})
}

func (s SpecRepo) GetSpec(id string) (*common.TableSpec, error) {
	data, err := s.repo.GetEntry([]string{id})
	if err != nil {
		return nil, err
	}
	spec := &common.TableSpec{}
	err = json.Unmarshal(data, spec)
	if err != nil {
		return nil, err
	}
	return spec, nil

}
func (s SpecRepo) CreateSpec(id string, spec *common.TableSpec) error {
	now := utils.GetNow()
	spec.CreatedAt = now
	spec.UpdatedAt = now

	data, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	err = s.repo.CreateEntry([]string{id}, data)
	if err != nil {
		return fmt.Errorf("error creating spec: %v", err)
	}
	return nil
}

func (s SpecRepo) UpdateSpec(id string, spec *common.TableSpec) error {
	now := utils.GetNow()
	spec.UpdatedAt = now
	data, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	err = s.repo.UpdateEntry([]string{id}, data)
	if err != nil {
		return fmt.Errorf("error updating spec: %v", err)
	}
	return nil
}

func (s SpecRepo) DeleteSpec(id string) error {
	err := s.repo.DeleteEntry([]string{id})
	if err != nil {
		return fmt.Errorf("error deleting spec: %v", err)
	}
	return nil
}

func (s SpecRepo) ListSpecIds() ([]string, error) {
	idStacks, err := s.repo.ListEntryIds([]string{})
	if err != nil {
		return nil, fmt.Errorf("error listing spec ids: %v", err)
	}
	ids := make([]string, len(idStacks))
	for i, idStack := range idStacks {
		lastIdx := len(idStack) - 1
		ids[i] = idStack[lastIdx]
	}
	return ids, nil
}
