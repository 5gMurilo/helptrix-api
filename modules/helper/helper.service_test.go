package helper_test

import (
	"errors"
	"testing"

	"github.com/5gMurilo/helptrix-api/core/domain"
	"github.com/5gMurilo/helptrix-api/core/utils"
	helpermodule "github.com/5gMurilo/helptrix-api/modules/helper"
	"github.com/google/uuid"
)

type mockHelperRepository struct {
	searchFn func(domain.HelperSearchParams) (domain.HelperListResponseDTO, error)
	called   bool
}

func (m *mockHelperRepository) Search(p domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
	m.called = true
	if m.searchFn != nil {
		return m.searchFn(p)
	}
	return domain.HelperListResponseDTO{}, nil
}

func TestSearch_BusinessUserSuccess(t *testing.T) {
	want := domain.HelperListResponseDTO{
		Data: []domain.HelperCardDTO{
			{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Name: "A"},
		},
		Total:    1,
		Page:     1,
		PageSize: 20,
	}
	mock := &mockHelperRepository{
		searchFn: func(p domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return want, nil
		},
	}
	svc := helpermodule.NewHelperService(mock)
	got, err := svc.Search(utils.UserTypeBusiness, domain.HelperSearchParams{Page: 1, PageSize: 20})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Data) != 1 || got.Data[0].Name != "A" {
		t.Errorf("got %+v, want one card named A", got)
	}
}

func TestSearch_HelperUserReturns403(t *testing.T) {
	mock := &mockHelperRepository{
		searchFn: func(domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			t.Fatal("repository must not be called for non-business user")
			return domain.HelperListResponseDTO{}, nil
		},
	}
	svc := helpermodule.NewHelperService(mock)
	_, err := svc.Search(utils.UserTypeHelper, domain.HelperSearchParams{Page: 1, PageSize: 20})
	if !errors.Is(err, utils.ErrBusinessOnly) {
		t.Fatalf("expected ErrBusinessOnly, got %v", err)
	}
}

func TestSearch_OtherUserTypeReturns403(t *testing.T) {
	mock := &mockHelperRepository{
		searchFn: func(domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			t.Fatal("repository must not be called")
			return domain.HelperListResponseDTO{}, nil
		},
	}
	svc := helpermodule.NewHelperService(mock)
	for _, ut := range []string{"", "unknown"} {
		t.Run(ut, func(t *testing.T) {
			_, err := svc.Search(ut, domain.HelperSearchParams{Page: 1, PageSize: 20})
			if !errors.Is(err, utils.ErrBusinessOnly) {
				t.Fatalf("expected ErrBusinessOnly for %q, got %v", ut, err)
			}
		})
	}
}

func TestSearch_NormalizesPageDefault(t *testing.T) {
	var gotParams domain.HelperSearchParams
	mock := &mockHelperRepository{
		searchFn: func(p domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			gotParams = p
			return domain.HelperListResponseDTO{Data: []domain.HelperCardDTO{}}, nil
		},
	}
	svc := helpermodule.NewHelperService(mock)
	_, err := svc.Search(utils.UserTypeBusiness, domain.HelperSearchParams{Page: 0, PageSize: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotParams.Page != 1 || gotParams.PageSize != 20 {
		t.Fatalf("expected Page=1 PageSize=20, got Page=%d PageSize=%d", gotParams.Page, gotParams.PageSize)
	}
}

func TestSearch_NegativePageFallsToDefault(t *testing.T) {
	var gotParams domain.HelperSearchParams
	mock := &mockHelperRepository{
		searchFn: func(p domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			gotParams = p
			return domain.HelperListResponseDTO{Data: []domain.HelperCardDTO{}}, nil
		},
	}
	svc := helpermodule.NewHelperService(mock)
	_, err := svc.Search(utils.UserTypeBusiness, domain.HelperSearchParams{Page: -5, PageSize: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotParams.Page != 1 {
		t.Fatalf("expected Page=1, got %d", gotParams.Page)
	}
}

func TestSearch_RepositoryErrorPropagates(t *testing.T) {
	dbErr := errors.New("db error")
	mock := &mockHelperRepository{
		searchFn: func(domain.HelperSearchParams) (domain.HelperListResponseDTO, error) {
			return domain.HelperListResponseDTO{}, dbErr
		},
	}
	svc := helpermodule.NewHelperService(mock)
	_, err := svc.Search(utils.UserTypeBusiness, domain.HelperSearchParams{Page: 1, PageSize: 20})
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected db error, got %v", err)
	}
}
