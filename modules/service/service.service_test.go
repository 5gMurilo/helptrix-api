package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	serviceinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/service"
	servicemodule "github.com/5gMurilo/helptrix-api/modules/service"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// mockServiceRepository implements IServiceRepository for tests.
type mockServiceRepository struct {
	CreateFn                      func(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	ExistsByNameAndUserFn         func(name string, userID uuid.UUID) (bool, error)
	ExistsByNameAndUserExcludingFn func(name string, userID uuid.UUID, excludeID uuid.UUID) (bool, error)
	UserHasCategoryFn             func(userID uuid.UUID, categoryID uint) (bool, error)
	ListFn                        func(userID uuid.UUID) ([]domain.ServiceResponseDTO, error)
	GetByIDFn                     func(serviceID uuid.UUID, userID uuid.UUID) (domain.ServiceResponseDTO, error)
	UpdateFn                      func(serviceID uuid.UUID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error)
	DeleteFn                      func(serviceID uuid.UUID, userID uuid.UUID) error
}

func (m *mockServiceRepository) Create(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	return m.CreateFn(userID, dto)
}

func (m *mockServiceRepository) ExistsByNameAndUser(name string, userID uuid.UUID) (bool, error) {
	return m.ExistsByNameAndUserFn(name, userID)
}

func (m *mockServiceRepository) ExistsByNameAndUserExcluding(name string, userID uuid.UUID, excludeID uuid.UUID) (bool, error) {
	return m.ExistsByNameAndUserExcludingFn(name, userID, excludeID)
}

func (m *mockServiceRepository) UserHasCategory(userID uuid.UUID, categoryID uint) (bool, error) {
	return m.UserHasCategoryFn(userID, categoryID)
}

func (m *mockServiceRepository) List(userID uuid.UUID) ([]domain.ServiceResponseDTO, error) {
	return m.ListFn(userID)
}

func (m *mockServiceRepository) GetByID(serviceID uuid.UUID, userID uuid.UUID) (domain.ServiceResponseDTO, error) {
	return m.GetByIDFn(serviceID, userID)
}

func (m *mockServiceRepository) Update(serviceID uuid.UUID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
	return m.UpdateFn(serviceID, userID, dto)
}

func (m *mockServiceRepository) Delete(serviceID uuid.UUID, userID uuid.UUID) error {
	return m.DeleteFn(serviceID, userID)
}

var _ serviceinterfaces.IServiceRepository = (*mockServiceRepository)(nil)

func validCreateDTO() domain.CreateServiceRequestDTO {
	return domain.CreateServiceRequestDTO{
		Name:          "Corte de Cabelo",
		Description:   "Corte profissional",
		ActuationDays: []string{"monday", "friday"},
		Value:         "49.90",
		StartTime:     "09:00",
		EndTime:       "18:00",
		OfferSince:    time.Now(),
		CategoryID:    1,
		Photos:        []string{},
	}
}

func defaultServiceResponse() domain.ServiceResponseDTO {
	return domain.ServiceResponseDTO{
		ID:            uuid.New(),
		Name:          "Corte de Cabelo",
		Description:   "Corte profissional",
		ActuationDays: []string{"monday", "friday"},
		Value:         decimal.RequireFromString("49.90"),
		StartTime:     "09:00",
		EndTime:       "18:00",
		OfferSince:    time.Now(),
		Category:      domain.ServiceCategoryDTO{ID: 1, Name: "Cabelereiro"},
	}
}

func defaultRepo() *mockServiceRepository {
	return &mockServiceRepository{
		UserHasCategoryFn: func(userID uuid.UUID, categoryID uint) (bool, error) {
			return true, nil
		},
		ExistsByNameAndUserFn: func(name string, userID uuid.UUID) (bool, error) {
			return false, nil
		},
		ExistsByNameAndUserExcludingFn: func(name string, userID uuid.UUID, excludeID uuid.UUID) (bool, error) {
			return false, nil
		},
		CreateFn: func(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{
				ID:            uuid.New(),
				Name:          dto.Name,
				Description:   dto.Description,
				ActuationDays: dto.ActuationDays,
				Value:         decimal.RequireFromString(dto.Value),
				StartTime:     dto.StartTime,
				EndTime:       dto.EndTime,
				OfferSince:    dto.OfferSince,
				Category:      domain.ServiceCategoryDTO{ID: dto.CategoryID, Name: "Cabelereiro"},
			}, nil
		},
		ListFn: func(userID uuid.UUID) ([]domain.ServiceResponseDTO, error) {
			return []domain.ServiceResponseDTO{defaultServiceResponse()}, nil
		},
		GetByIDFn: func(serviceID uuid.UUID, userID uuid.UUID) (domain.ServiceResponseDTO, error) {
			return defaultServiceResponse(), nil
		},
		UpdateFn: func(serviceID uuid.UUID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return defaultServiceResponse(), nil
		},
		DeleteFn: func(serviceID uuid.UUID, userID uuid.UUID) error {
			return nil
		},
	}
}

func TestServiceService_Create(t *testing.T) {
	t.Run("happy path: helper válido cria serviço com sucesso", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		userID := uuid.New()

		resp, err := svc.Create(userID, utils.UserTypeHelper, validCreateDTO())

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if resp.ID == (uuid.UUID{}) {
			t.Error("esperado ID não nulo na resposta")
		}
		if resp.Name != "Corte de Cabelo" {
			t.Errorf("esperado name 'Corte de Cabelo', obteve: %s", resp.Name)
		}
	})

	t.Run("usuário não-helper: retorna ErrHelperOnly", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.Create(uuid.New(), utils.UserTypeBusiness, validCreateDTO())

		if !errors.Is(err, utils.ErrHelperOnly) {
			t.Errorf("esperado ErrHelperOnly, obteve: %v", err)
		}
	})

	t.Run("value inválido: formato não parseável retorna erro", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		dto := validCreateDTO()
		dto.Value = "abc"

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrInvalidValueFormat) {
			t.Errorf("esperado ErrInvalidValueFormat, obteve: %v", err)
		}
	})

	t.Run("value <= 0: retorna erro", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		dto := validCreateDTO()
		dto.Value = "0"

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrValueNotPositive) {
			t.Errorf("esperado ErrValueNotPositive, obteve: %v", err)
		}
	})

	t.Run("start_time com formato inválido retorna erro", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		dto := validCreateDTO()
		dto.StartTime = "9:00"

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrInvalidStartTimeFormat) {
			t.Errorf("esperado ErrInvalidStartTimeFormat, obteve: %v", err)
		}
	})

	t.Run("end_time com formato inválido retorna erro", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		dto := validCreateDTO()
		dto.EndTime = "1800"

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrInvalidEndTimeFormat) {
			t.Errorf("esperado ErrInvalidEndTimeFormat, obteve: %v", err)
		}
	})

	t.Run("categoria não atribuída ao usuário: retorna ErrCategoryNotAssignedToUser", func(t *testing.T) {
		repo := defaultRepo()
		repo.UserHasCategoryFn = func(userID uuid.UUID, categoryID uint) (bool, error) {
			return false, nil
		}
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, validCreateDTO())

		if !errors.Is(err, utils.ErrCategoryNotAssignedToUser) {
			t.Errorf("esperado ErrCategoryNotAssignedToUser, obteve: %v", err)
		}
	})

	t.Run("nome duplicado: retorna ErrServiceNameNotUnique", func(t *testing.T) {
		repo := defaultRepo()
		repo.ExistsByNameAndUserFn = func(name string, userID uuid.UUID) (bool, error) {
			return true, nil
		}
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, validCreateDTO())

		if !errors.Is(err, utils.ErrServiceNameNotUnique) {
			t.Errorf("esperado ErrServiceNameNotUnique, obteve: %v", err)
		}
	})

	t.Run("erro do repo Create propagado", func(t *testing.T) {
		repoErr := errors.New("database connection failed")
		repo := defaultRepo()
		repo.CreateFn = func(userID uuid.UUID, dto domain.CreateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, repoErr
		}
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.Create(uuid.New(), utils.UserTypeHelper, validCreateDTO())

		if !errors.Is(err, repoErr) {
			t.Errorf("esperado erro do repo propagado, obteve: %v", err)
		}
	})
}

func TestServiceService_List(t *testing.T) {
	t.Run("helper válido lista serviços com sucesso", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		userID := uuid.New()

		result, err := svc.List(userID, utils.UserTypeHelper)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if len(result) == 0 {
			t.Error("esperado pelo menos um serviço na lista")
		}
	})

	t.Run("usuário não-helper: retorna ErrHelperOnly", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.List(uuid.New(), utils.UserTypeBusiness)

		if !errors.Is(err, utils.ErrHelperOnly) {
			t.Errorf("esperado ErrHelperOnly, obteve: %v", err)
		}
	})

	t.Run("erro do repo List propagado", func(t *testing.T) {
		repoErr := errors.New("db error on list")
		repo := defaultRepo()
		repo.ListFn = func(userID uuid.UUID) ([]domain.ServiceResponseDTO, error) {
			return nil, repoErr
		}
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.List(uuid.New(), utils.UserTypeHelper)

		if !errors.Is(err, repoErr) {
			t.Errorf("esperado erro do repo propagado, obteve: %v", err)
		}
	})
}

func TestServiceService_GetByID(t *testing.T) {
	t.Run("helper válido obtém serviço por ID com sucesso", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		result, err := svc.GetByID(uuid.New(), uuid.New(), utils.UserTypeHelper)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if result.ID == (uuid.UUID{}) {
			t.Error("esperado ID não nulo na resposta")
		}
	})

	t.Run("usuário não-helper: retorna ErrHelperOnly", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.GetByID(uuid.New(), uuid.New(), utils.UserTypeBusiness)

		if !errors.Is(err, utils.ErrHelperOnly) {
			t.Errorf("esperado ErrHelperOnly, obteve: %v", err)
		}
	})

	t.Run("serviço não encontrado: ErrServiceNotFound propagado", func(t *testing.T) {
		repo := defaultRepo()
		repo.GetByIDFn = func(serviceID uuid.UUID, userID uuid.UUID) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
		}
		svc := servicemodule.NewServiceService(repo)

		_, err := svc.GetByID(uuid.New(), uuid.New(), utils.UserTypeHelper)

		if !errors.Is(err, utils.ErrServiceNotFound) {
			t.Errorf("esperado ErrServiceNotFound, obteve: %v", err)
		}
	})
}

func TestServiceService_Update(t *testing.T) {
	name := "Novo Nome"
	value := "99.99"
	start := "08:00"
	end := "17:00"
	catID := uint(2)

	t.Run("helper válido atualiza serviço com sucesso", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		dto := domain.UpdateServiceRequestDTO{Name: &name}
		result, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if result.ID == (uuid.UUID{}) {
			t.Error("esperado ID não nulo na resposta")
		}
	})

	t.Run("usuário não-helper: retorna ErrHelperOnly", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		dto := domain.UpdateServiceRequestDTO{Name: &name}
		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeBusiness, dto)

		if !errors.Is(err, utils.ErrHelperOnly) {
			t.Errorf("esperado ErrHelperOnly, obteve: %v", err)
		}
	})

	t.Run("serviço não encontrado: ErrServiceNotFound propagado", func(t *testing.T) {
		repo := defaultRepo()
		repo.UpdateFn = func(serviceID uuid.UUID, userID uuid.UUID, dto domain.UpdateServiceRequestDTO) (domain.ServiceResponseDTO, error) {
			return domain.ServiceResponseDTO{}, utils.ErrServiceNotFound
		}
		svc := servicemodule.NewServiceService(repo)

		dto := domain.UpdateServiceRequestDTO{Name: &name}
		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrServiceNotFound) {
			t.Errorf("esperado ErrServiceNotFound, obteve: %v", err)
		}
	})

	t.Run("value inválido: retorna ErrInvalidValueFormat", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		invalid := "abc"
		dto := domain.UpdateServiceRequestDTO{Value: &invalid}

		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrInvalidValueFormat) {
			t.Errorf("esperado ErrInvalidValueFormat, obteve: %v", err)
		}
	})

	t.Run("value não positivo: retorna ErrValueNotPositive", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		zero := "0"
		dto := domain.UpdateServiceRequestDTO{Value: &zero}

		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrValueNotPositive) {
			t.Errorf("esperado ErrValueNotPositive, obteve: %v", err)
		}
	})

	t.Run("start_time inválido: retorna ErrInvalidStartTimeFormat", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		invalid := "9:00"
		dto := domain.UpdateServiceRequestDTO{Value: &value, StartTime: &invalid}

		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrInvalidStartTimeFormat) {
			t.Errorf("esperado ErrInvalidStartTimeFormat, obteve: %v", err)
		}
	})

	t.Run("end_time inválido: retorna ErrInvalidEndTimeFormat", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)
		invalid := "1800"
		dto := domain.UpdateServiceRequestDTO{Value: &value, StartTime: &start, EndTime: &invalid}

		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrInvalidEndTimeFormat) {
			t.Errorf("esperado ErrInvalidEndTimeFormat, obteve: %v", err)
		}
	})

	t.Run("categoria não atribuída: retorna ErrCategoryNotAssignedToUser", func(t *testing.T) {
		repo := defaultRepo()
		repo.UserHasCategoryFn = func(userID uuid.UUID, categoryID uint) (bool, error) {
			return false, nil
		}
		svc := servicemodule.NewServiceService(repo)
		dto := domain.UpdateServiceRequestDTO{Value: &value, StartTime: &start, EndTime: &end, CategoryID: &catID}

		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrCategoryNotAssignedToUser) {
			t.Errorf("esperado ErrCategoryNotAssignedToUser, obteve: %v", err)
		}
	})

	t.Run("nome duplicado excluindo self: retorna ErrServiceNameNotUnique", func(t *testing.T) {
		repo := defaultRepo()
		repo.ExistsByNameAndUserExcludingFn = func(n string, userID uuid.UUID, excludeID uuid.UUID) (bool, error) {
			return true, nil
		}
		svc := servicemodule.NewServiceService(repo)
		dto := domain.UpdateServiceRequestDTO{Name: &name}

		_, err := svc.Update(uuid.New(), uuid.New(), utils.UserTypeHelper, dto)

		if !errors.Is(err, utils.ErrServiceNameNotUnique) {
			t.Errorf("esperado ErrServiceNameNotUnique, obteve: %v", err)
		}
	})
}

func TestServiceService_Delete(t *testing.T) {
	t.Run("helper válido deleta serviço com sucesso", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		err := svc.Delete(uuid.New(), uuid.New(), utils.UserTypeHelper)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
	})

	t.Run("usuário não-helper: retorna ErrHelperOnly", func(t *testing.T) {
		repo := defaultRepo()
		svc := servicemodule.NewServiceService(repo)

		err := svc.Delete(uuid.New(), uuid.New(), utils.UserTypeBusiness)

		if !errors.Is(err, utils.ErrHelperOnly) {
			t.Errorf("esperado ErrHelperOnly, obteve: %v", err)
		}
	})

	t.Run("serviço não encontrado: ErrServiceNotFound propagado", func(t *testing.T) {
		repo := defaultRepo()
		repo.DeleteFn = func(serviceID uuid.UUID, userID uuid.UUID) error {
			return utils.ErrServiceNotFound
		}
		svc := servicemodule.NewServiceService(repo)

		err := svc.Delete(uuid.New(), uuid.New(), utils.UserTypeHelper)

		if !errors.Is(err, utils.ErrServiceNotFound) {
			t.Errorf("esperado ErrServiceNotFound, obteve: %v", err)
		}
	})
}
