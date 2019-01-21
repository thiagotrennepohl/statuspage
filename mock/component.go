package mock

import (
	"github.com/involvestecnologia/statuspage/component"
	"github.com/involvestecnologia/statuspage/errors"
	"github.com/involvestecnologia/statuspage/models"
	"gopkg.in/mgo.v2/bson"
)

type MockComponentDAO struct {
	components []models.Component
}

func NewMockComponentDAO() component.Repository {
	return &MockComponentDAO{
		components: []models.Component{
			models.Component{
				Ref:     ZeroTimeHex,
				Name:    "first",
				Address: "",
			},
			models.Component{
				Ref:     bson.NewObjectIdWithTime(bson.Now()).Hex(),
				Name:    "first_comp_with_group",
				Address: "",
			},
			models.Component{
				Ref:     bson.NewObjectIdWithTime(bson.Now()).Hex(),
				Name:    "test",
				Address: "",
			},
			models.Component{
				Ref:     bson.NewObjectIdWithTime(bson.Now()).Hex(),
				Name:    "last_comp_with_group",
				Address: "",
			},
			models.Component{
				Ref:     bson.NewObjectIdWithTime(bson.Now()).Hex(),
				Name:    "last",
				Address: "",
			},
		},
	}
}

func (m *MockComponentDAO) List() ([]models.Component, error) {
	return m.components, nil
}
func (m *MockComponentDAO) Find(q map[string]interface{}) (models.Component, error) {
	var c models.Component
	if keyValue, hasKey := q["ref"]; hasKey {
		for _, c := range m.components {
			if c.Ref == keyValue {
				return c, nil
			}
		}
	} else {
		if keyValue, hasKey := q["name"]; hasKey {
			for _, c := range m.components {
				if c.Name == keyValue {
					return c, nil
				}
			}
		} else {
			return c, errors.E(errors.ErrInvalidQuery)
		}
	}

	return c, errors.E(errors.ErrNotFound)
}
func (m *MockComponentDAO) Insert(component models.Component) (string, error) {
	if component.Ref == "" {
		component.Ref = bson.NewObjectId().Hex()
	}
	m.components = append(m.components, component)
	return component.Ref, nil
}
func (m *MockComponentDAO) Update(ref string, component models.Component) error {
	for k, comp := range m.components {
		if comp.Ref == ref {
			m.components[k].Name = component.Name
			m.components[k].Address = component.Address
			return nil
		}
	}
	return errors.E(errors.ErrNotFound)
}
func (m *MockComponentDAO) Delete(ref string) error {
	for k, comp := range m.components {
		if comp.Ref == ref {
			m.components = append(m.components[:k], m.components[k+1:]...)
			return nil
		}
	}
	return errors.E(errors.ErrNotFound)
}
