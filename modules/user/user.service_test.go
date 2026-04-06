package user_test

import (
	"errors"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	usermodule "github.com/5gMurilo/helptrix-api/modules/user"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

// mockUserRepo implementa IUserRepository para testes.
type mockUserRepo struct {
	GetProfileFn    func(userID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error)
	UpdateProfileFn func(userID uuid.UUID, dto domain.UpdateProfileRequestDTO) error
	DeleteProfileFn func(userID uuid.UUID) error
}

func (m *mockUserRepo) GetProfile(userID uuid.UUID, filters domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
	return m.GetProfileFn(userID, filters)
}

func (m *mockUserRepo) UpdateProfile(userID uuid.UUID, dto domain.UpdateProfileRequestDTO) error {
	return m.UpdateProfileFn(userID, dto)
}

func (m *mockUserRepo) DeleteProfile(userID uuid.UUID) error {
	return m.DeleteProfileFn(userID)
}

func (m *mockUserRepo) GetProfilePicture(userID uuid.UUID) (string, error) {
	return "", nil
}

func (m *mockUserRepo) UpdateProfilePicture(userID uuid.UUID, url string) error {
	return nil
}

func TestGetProfileService(t *testing.T) {
	targetID := uuid.New()
	requesterID := uuid.New()

	t.Run("business requester: filters passed through", func(t *testing.T) {
		catID := uint(5)
		filters := domain.ProfileFilters{
			CategoryID:    &catID,
			ActuationDays: []string{"monday", "friday"},
		}

		var capturedFilters domain.ProfileFilters
		var capturedID uuid.UUID

		repo := &mockUserRepo{
			GetProfileFn: func(userID uuid.UUID, f domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
				capturedID = userID
				capturedFilters = f
				return domain.GetProfileResponseDTO{ID: targetID, Reviews: []interface{}{}}, nil
			},
		}

		svc := usermodule.NewUserService(repo)
		resp, err := svc.GetProfile(requesterID, utils.UserTypeBusiness, targetID, filters)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if capturedID != targetID {
			t.Errorf("esperado targetID %v, obteve: %v", targetID, capturedID)
		}
		if capturedFilters.CategoryID == nil || *capturedFilters.CategoryID != catID {
			t.Error("esperado CategoryID preservado para requisitante business")
		}
		if len(capturedFilters.ActuationDays) != 2 {
			t.Error("esperado ActuationDays preservado para requisitante business")
		}
		if resp.ID != targetID {
			t.Errorf("esperado ID %v na resposta, obteve: %v", targetID, resp.ID)
		}
	})

	t.Run("non-business requester: filters stripped", func(t *testing.T) {
		catID := uint(5)
		filters := domain.ProfileFilters{
			CategoryID:    &catID,
			ActuationDays: []string{"monday"},
		}

		var capturedFilters domain.ProfileFilters

		repo := &mockUserRepo{
			GetProfileFn: func(userID uuid.UUID, f domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
				capturedFilters = f
				return domain.GetProfileResponseDTO{ID: targetID, Reviews: []interface{}{}}, nil
			},
		}

		svc := usermodule.NewUserService(repo)
		_, err := svc.GetProfile(requesterID, utils.UserTypeHelper, targetID, filters)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if capturedFilters.CategoryID != nil {
			t.Error("esperado CategoryID nil para requisitante não-business")
		}
		if len(capturedFilters.ActuationDays) != 0 {
			t.Error("esperado ActuationDays vazio para requisitante não-business")
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repoErr := errors.New("database connection failed")

		repo := &mockUserRepo{
			GetProfileFn: func(userID uuid.UUID, f domain.ProfileFilters) (domain.GetProfileResponseDTO, error) {
				return domain.GetProfileResponseDTO{}, repoErr
			},
		}

		svc := usermodule.NewUserService(repo)
		_, err := svc.GetProfile(requesterID, utils.UserTypeBusiness, targetID, domain.ProfileFilters{})

		if err == nil {
			t.Fatal("esperado erro propagado do repositório, obteve nil")
		}
		if !errors.Is(err, repoErr) {
			t.Errorf("esperado repoErr, obteve: %v", err)
		}
	})
}

func TestUpdateProfileService(t *testing.T) {
	ownerID := uuid.New()
	otherID := uuid.New()

	t.Run("happy path: requesterID == targetID, repo called", func(t *testing.T) {
		repoCalled := false
		dto := domain.UpdateProfileRequestDTO{Biography: "Meu perfil"}

		repo := &mockUserRepo{
			UpdateProfileFn: func(userID uuid.UUID, d domain.UpdateProfileRequestDTO) error {
				repoCalled = true
				if userID != ownerID {
					t.Errorf("esperado ownerID %v, obteve: %v", ownerID, userID)
				}
				return nil
			},
		}

		svc := usermodule.NewUserService(repo)
		err := svc.UpdateProfile(ownerID, ownerID, dto)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if !repoCalled {
			t.Error("esperado repositório chamado")
		}
	})

	t.Run("ownership error: requesterID != targetID returns ErrNotOwner", func(t *testing.T) {
		repoCalled := false

		repo := &mockUserRepo{
			UpdateProfileFn: func(userID uuid.UUID, d domain.UpdateProfileRequestDTO) error {
				repoCalled = true
				return nil
			},
		}

		svc := usermodule.NewUserService(repo)
		err := svc.UpdateProfile(otherID, ownerID, domain.UpdateProfileRequestDTO{})

		if err == nil {
			t.Fatal("esperado ErrNotOwner, obteve nil")
		}
		if !errors.Is(err, utils.ErrNotOwner) {
			t.Errorf("esperado ErrNotOwner, obteve: %v", err)
		}
		if repoCalled {
			t.Error("repositório não deve ser chamado quando ownership falha")
		}
	})

	t.Run("repo error propagated", func(t *testing.T) {
		repoErr := errors.New("update failed")

		repo := &mockUserRepo{
			UpdateProfileFn: func(userID uuid.UUID, d domain.UpdateProfileRequestDTO) error {
				return repoErr
			},
		}

		svc := usermodule.NewUserService(repo)
		err := svc.UpdateProfile(ownerID, ownerID, domain.UpdateProfileRequestDTO{})

		if err == nil {
			t.Fatal("esperado erro propagado do repositório, obteve nil")
		}
		if !errors.Is(err, repoErr) {
			t.Errorf("esperado repoErr, obteve: %v", err)
		}
	})
}

func TestDeleteProfileService(t *testing.T) {
	ownerID := uuid.New()
	otherID := uuid.New()

	t.Run("happy path: requesterID == targetID, repo called", func(t *testing.T) {
		repoCalled := false

		repo := &mockUserRepo{
			DeleteProfileFn: func(userID uuid.UUID) error {
				repoCalled = true
				if userID != ownerID {
					t.Errorf("esperado ownerID %v, obteve: %v", ownerID, userID)
				}
				return nil
			},
		}

		svc := usermodule.NewUserService(repo)
		err := svc.DeleteProfile(ownerID, ownerID)

		if err != nil {
			t.Fatalf("esperado sem erro, obteve: %v", err)
		}
		if !repoCalled {
			t.Error("esperado repositório chamado")
		}
	})

	t.Run("ownership error: requesterID != targetID returns ErrNotOwner", func(t *testing.T) {
		repoCalled := false

		repo := &mockUserRepo{
			DeleteProfileFn: func(userID uuid.UUID) error {
				repoCalled = true
				return nil
			},
		}

		svc := usermodule.NewUserService(repo)
		err := svc.DeleteProfile(otherID, ownerID)

		if err == nil {
			t.Fatal("esperado ErrNotOwner, obteve nil")
		}
		if !errors.Is(err, utils.ErrNotOwner) {
			t.Errorf("esperado ErrNotOwner, obteve: %v", err)
		}
		if repoCalled {
			t.Error("repositório não deve ser chamado quando ownership falha")
		}
	})

	t.Run("repo error propagated (ErrUserNotFound)", func(t *testing.T) {
		repo := &mockUserRepo{
			DeleteProfileFn: func(userID uuid.UUID) error {
				return utils.ErrUserNotFound
			},
		}

		svc := usermodule.NewUserService(repo)
		err := svc.DeleteProfile(ownerID, ownerID)

		if err == nil {
			t.Fatal("esperado ErrUserNotFound, obteve nil")
		}
		if !errors.Is(err, utils.ErrUserNotFound) {
			t.Errorf("esperado ErrUserNotFound, obteve: %v", err)
		}
	})
}
