package proposal_test

import (
	"errors"
	"testing"
	"time"

	"github.com/5gMurilo/helptrix-api/core/domain"
	proposalinterfaces "github.com/5gMurilo/helptrix-api/core/interfaces/proposal"
	proposalmodule "github.com/5gMurilo/helptrix-api/modules/proposal"
	"github.com/5gMurilo/helptrix-api/core/utils"
	"github.com/google/uuid"
)

// mockProposalRepository implements IProposalRepository for unit tests.
type mockProposalRepository struct {
	CreateFn                      func(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.Proposal, error)
	FindByIDFn                    func(id uuid.UUID) (*domain.Proposal, error)
	UpdateStatusFn                func(id uuid.UUID, status string) (*domain.Proposal, error)
	ListByUserIDFn                func(userID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error)
	ListByHelperIDFn              func(helperID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error)
	HasBlockingProposalForHelperFn func(userID uuid.UUID, helperID uuid.UUID) (bool, error)
}

func (m *mockProposalRepository) Create(dto domain.CreateProposalRequestDTO, userID uuid.UUID) (domain.Proposal, error) {
	return m.CreateFn(dto, userID)
}

func (m *mockProposalRepository) FindByID(id uuid.UUID) (*domain.Proposal, error) {
	return m.FindByIDFn(id)
}

func (m *mockProposalRepository) UpdateStatus(id uuid.UUID, status string) (*domain.Proposal, error) {
	return m.UpdateStatusFn(id, status)
}

func (m *mockProposalRepository) ListByUserID(userID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
	return m.ListByUserIDFn(userID, statusFilter)
}

func (m *mockProposalRepository) ListByHelperID(helperID uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
	return m.ListByHelperIDFn(helperID, statusFilter)
}

func (m *mockProposalRepository) HasBlockingProposalForHelper(userID uuid.UUID, helperID uuid.UUID) (bool, error) {
	return m.HasBlockingProposalForHelperFn(userID, helperID)
}

var _ proposalinterfaces.IProposalRepository = (*mockProposalRepository)(nil)

// helper builders

func newProposal(userID, helperID uuid.UUID, status string) domain.Proposal {
	return domain.Proposal{
		ID:          uuid.New(),
		UserID:      userID,
		HelperID:    helperID,
		CategoryID:  1,
		Description: "Preciso de ajuda com encanamento",
		Value:       150.00,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func defaultRepo(userID, helperID uuid.UUID) *mockProposalRepository {
	proposal := newProposal(userID, helperID, utils.ProposalStatusPending)
	return &mockProposalRepository{
		HasBlockingProposalForHelperFn: func(uid uuid.UUID, hid uuid.UUID) (bool, error) {
			return false, nil
		},
		CreateFn: func(dto domain.CreateProposalRequestDTO, uid uuid.UUID) (domain.Proposal, error) {
			return proposal, nil
		},
		FindByIDFn: func(id uuid.UUID) (*domain.Proposal, error) {
			return &proposal, nil
		},
		UpdateStatusFn: func(id uuid.UUID, status string) (*domain.Proposal, error) {
			updated := proposal
			updated.Status = status
			return &updated, nil
		},
		ListByUserIDFn: func(uid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
			return []domain.ProposalResponseDTO{}, nil
		},
		ListByHelperIDFn: func(hid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
			return []domain.ProposalResponseDTO{}, nil
		},
	}
}

func validCreateDTO(helperID uuid.UUID) domain.CreateProposalRequestDTO {
	return domain.CreateProposalRequestDTO{
		HelperID:    helperID,
		CategoryID:  1,
		Description: "Preciso de ajuda com encanamento",
		Value:       150.00,
	}
}

// --- Create tests ---

func TestProposalService_Create_Success(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	resp, err := svc.Create(validCreateDTO(helperID), userID)

	if err != nil {
		t.Fatalf("esperado sem erro, obteve: %v", err)
	}
	if resp.UserID != userID {
		t.Errorf("esperado UserID %v, obteve %v", userID, resp.UserID)
	}
	if resp.Status != utils.ProposalStatusPending {
		t.Errorf("esperado status '%s', obteve '%s'", utils.ProposalStatusPending, resp.Status)
	}
}

func TestProposalService_Create_BlockedSameHelper(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	repo.HasBlockingProposalForHelperFn = func(uid uuid.UUID, hid uuid.UUID) (bool, error) {
		return true, nil
	}
	svc := proposalmodule.NewProposalService(repo)

	_, err := svc.Create(validCreateDTO(helperID), userID)

	if !errors.Is(err, utils.ErrProposalAlreadyActiveForHelper) {
		t.Errorf("esperado ErrProposalAlreadyActiveForHelper, obteve: %v", err)
	}
}

func TestProposalService_Create_HasBlockingProposalError(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	repoErr := errors.New("db error checking proposal for helper")
	repo.HasBlockingProposalForHelperFn = func(uid uuid.UUID, hid uuid.UUID) (bool, error) {
		return false, repoErr
	}
	svc := proposalmodule.NewProposalService(repo)

	_, err := svc.Create(validCreateDTO(helperID), userID)

	if !errors.Is(err, repoErr) {
		t.Errorf("esperado erro do repo propagado, obteve: %v", err)
	}
}

func TestProposalService_Create_AllowedWhenPreviousIsAccepted(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)

	var capturedHelperID uuid.UUID
	repo.HasBlockingProposalForHelperFn = func(uid uuid.UUID, hid uuid.UUID) (bool, error) {
		capturedHelperID = hid
		return false, nil
	}
	svc := proposalmodule.NewProposalService(repo)

	resp, err := svc.Create(validCreateDTO(helperID), userID)

	if err != nil {
		t.Fatalf("esperado sem erro quando proposta anterior esta accepted, obteve: %v", err)
	}
	if capturedHelperID != helperID {
		t.Errorf("esperado helperID %v passado ao repo, obteve %v", helperID, capturedHelperID)
	}
	if resp.Status != utils.ProposalStatusPending {
		t.Errorf("esperado status '%s', obteve '%s'", utils.ProposalStatusPending, resp.Status)
	}
}

func TestProposalService_Create_AllowedForDifferentHelper(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	otherHelperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	repo.HasBlockingProposalForHelperFn = func(uid uuid.UUID, hid uuid.UUID) (bool, error) {
		return hid == helperID, nil
	}
	svc := proposalmodule.NewProposalService(repo)

	resp, err := svc.Create(validCreateDTO(otherHelperID), userID)

	if err != nil {
		t.Fatalf("esperado sem erro para helper diferente, obteve: %v", err)
	}
	if resp.UserID != userID {
		t.Errorf("esperado UserID %v, obteve %v", userID, resp.UserID)
	}
}

func TestProposalService_Create_RepoError(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	repoErr := errors.New("db error on create")
	repo.CreateFn = func(dto domain.CreateProposalRequestDTO, uid uuid.UUID) (domain.Proposal, error) {
		return domain.Proposal{}, repoErr
	}
	svc := proposalmodule.NewProposalService(repo)

	_, err := svc.Create(validCreateDTO(helperID), userID)

	if !errors.Is(err, repoErr) {
		t.Errorf("esperado erro do repo propagado, obteve: %v", err)
	}
}

// --- GetByID tests ---

func TestProposalService_GetByID_FoundAsOwner(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	proposalID := uuid.New()
	p := newProposal(userID, helperID, utils.ProposalStatusPending)
	p.ID = proposalID
	repo.FindByIDFn = func(id uuid.UUID) (*domain.Proposal, error) {
		return &p, nil
	}
	svc := proposalmodule.NewProposalService(repo)

	resp, err := svc.GetByID(proposalID, userID)

	if err != nil {
		t.Fatalf("esperado sem erro, obteve: %v", err)
	}
	if resp.ID != proposalID {
		t.Errorf("esperado ID %v, obteve %v", proposalID, resp.ID)
	}
}

func TestProposalService_GetByID_FoundAsHelper(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	proposalID := uuid.New()
	p := newProposal(userID, helperID, utils.ProposalStatusPending)
	p.ID = proposalID
	repo.FindByIDFn = func(id uuid.UUID) (*domain.Proposal, error) {
		return &p, nil
	}
	svc := proposalmodule.NewProposalService(repo)

	resp, err := svc.GetByID(proposalID, helperID)

	if err != nil {
		t.Fatalf("esperado sem erro, obteve: %v", err)
	}
	if resp.HelperID != helperID {
		t.Errorf("esperado HelperID %v, obteve %v", helperID, resp.HelperID)
	}
}

func TestProposalService_GetByID_NotParticipant(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	stranger := uuid.New()
	_, err := svc.GetByID(uuid.New(), stranger)

	if !errors.Is(err, utils.ErrNotProposalParticipant) {
		t.Errorf("esperado ErrNotProposalParticipant, obteve: %v", err)
	}
}

func TestProposalService_GetByID_NotFound(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	repo.FindByIDFn = func(id uuid.UUID) (*domain.Proposal, error) {
		return nil, utils.ErrProposalNotFound
	}
	svc := proposalmodule.NewProposalService(repo)

	_, err := svc.GetByID(uuid.New(), userID)

	if !errors.Is(err, utils.ErrProposalNotFound) {
		t.Errorf("esperado ErrProposalNotFound, obteve: %v", err)
	}
}

// --- UpdateStatus tests ---

func TestProposalService_UpdateStatus_TerminalStatus(t *testing.T) {
	terminalStatuses := []string{
		utils.ProposalStatusRefused,
		utils.ProposalStatusCancelled,
		utils.ProposalStatusFinished,
	}

	for _, status := range terminalStatuses {
		status := status
		t.Run("terminal status: "+status, func(t *testing.T) {
			userID := uuid.New()
			helperID := uuid.New()
			repo := defaultRepo(userID, helperID)
			p := newProposal(userID, helperID, status)
			repo.FindByIDFn = func(id uuid.UUID) (*domain.Proposal, error) {
				return &p, nil
			}
			svc := proposalmodule.NewProposalService(repo)

			dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusAccepted}
			_, err := svc.UpdateStatus(p.ID, dto, helperID, utils.UserTypeHelper)

			if !errors.Is(err, utils.ErrProposalFinished) {
				t.Errorf("[%s] esperado ErrProposalFinished, obteve: %v", status, err)
			}
		})
	}
}

func TestProposalService_UpdateStatus_InvalidTargetStatus(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	dto := domain.UpdateProposalStatusRequestDTO{Status: "invalid_status"}
	_, err := svc.UpdateStatus(uuid.New(), dto, helperID, utils.UserTypeHelper)

	if !errors.Is(err, utils.ErrProposalInvalidStatus) {
		t.Errorf("esperado ErrProposalInvalidStatus, obteve: %v", err)
	}
}

func TestProposalService_UpdateStatus_Unauthorized_NotHelper(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusAccepted}
	_, err := svc.UpdateStatus(uuid.New(), dto, userID, utils.UserTypeBusiness)

	if !errors.Is(err, utils.ErrProposalUnauthorized) {
		t.Errorf("esperado ErrProposalUnauthorized, obteve: %v", err)
	}
}

func TestProposalService_UpdateStatus_Unauthorized_WrongHelper(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	wrongHelper := uuid.New()
	dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusAccepted}
	_, err := svc.UpdateStatus(uuid.New(), dto, wrongHelper, utils.UserTypeHelper)

	if !errors.Is(err, utils.ErrProposalUnauthorized) {
		t.Errorf("esperado ErrProposalUnauthorized para helper errado, obteve: %v", err)
	}
}

func TestProposalService_UpdateStatus_CancelByOwner(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusCancelled}
	resp, err := svc.UpdateStatus(uuid.New(), dto, userID, utils.UserTypeBusiness)

	if err != nil {
		t.Fatalf("esperado sem erro ao cancelar como owner, obteve: %v", err)
	}
	if resp.Status != utils.ProposalStatusCancelled {
		t.Errorf("esperado status '%s', obteve '%s'", utils.ProposalStatusCancelled, resp.Status)
	}
}

func TestProposalService_UpdateStatus_CancelByHelper(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	svc := proposalmodule.NewProposalService(repo)

	dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusCancelled}
	resp, err := svc.UpdateStatus(uuid.New(), dto, helperID, utils.UserTypeHelper)

	if err != nil {
		t.Fatalf("esperado sem erro ao cancelar como helper, obteve: %v", err)
	}
	if resp.Status != utils.ProposalStatusCancelled {
		t.Errorf("esperado status '%s', obteve '%s'", utils.ProposalStatusCancelled, resp.Status)
	}
}

func TestProposalService_UpdateStatus_AllowedTransitions(t *testing.T) {
	type transition struct {
		from string
		to   string
	}

	allowedTransitions := []transition{
		{utils.ProposalStatusPending, utils.ProposalStatusAccepted},
		{utils.ProposalStatusPending, utils.ProposalStatusRefused},
		{utils.ProposalStatusPending, utils.ProposalStatusCancelled},
		{utils.ProposalStatusAccepted, utils.ProposalStatusInProgress},
		{utils.ProposalStatusAccepted, utils.ProposalStatusCancelled},
		{utils.ProposalStatusInProgress, utils.ProposalStatusFinished},
		{utils.ProposalStatusInProgress, utils.ProposalStatusCancelled},
	}

	for _, tr := range allowedTransitions {
		tr := tr
		t.Run(tr.from+"->"+tr.to, func(t *testing.T) {
			userID := uuid.New()
			helperID := uuid.New()
			repo := defaultRepo(userID, helperID)
			p := newProposal(userID, helperID, tr.from)
			repo.FindByIDFn = func(id uuid.UUID) (*domain.Proposal, error) {
				return &p, nil
			}
			repo.UpdateStatusFn = func(id uuid.UUID, status string) (*domain.Proposal, error) {
				updated := p
				updated.Status = status
				return &updated, nil
			}
			svc := proposalmodule.NewProposalService(repo)

			// determine requester: cancel can be by owner, others by helper
			requesterID := helperID
			requesterType := utils.UserTypeHelper
			if tr.to == utils.ProposalStatusCancelled {
				requesterID = userID
				requesterType = utils.UserTypeBusiness
			}

			dto := domain.UpdateProposalStatusRequestDTO{Status: tr.to}
			resp, err := svc.UpdateStatus(p.ID, dto, requesterID, requesterType)

			if err != nil {
				t.Errorf("[%s->%s] esperado sem erro, obteve: %v", tr.from, tr.to, err)
				return
			}
			if resp.Status != tr.to {
				t.Errorf("[%s->%s] esperado status '%s', obteve '%s'", tr.from, tr.to, tr.to, resp.Status)
			}
		})
	}
}

func TestProposalService_UpdateStatus_InvalidTransition(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	// current status is pending, trying to go to finished (not allowed)
	svc := proposalmodule.NewProposalService(repo)

	dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusFinished}
	_, err := svc.UpdateStatus(uuid.New(), dto, helperID, utils.UserTypeHelper)

	if !errors.Is(err, utils.ErrProposalInvalidStatus) {
		t.Errorf("esperado ErrProposalInvalidStatus para transição inválida, obteve: %v", err)
	}
}

func TestProposalService_UpdateStatus_NotFound(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)
	repo.FindByIDFn = func(id uuid.UUID) (*domain.Proposal, error) {
		return nil, utils.ErrProposalNotFound
	}
	svc := proposalmodule.NewProposalService(repo)

	dto := domain.UpdateProposalStatusRequestDTO{Status: utils.ProposalStatusAccepted}
	_, err := svc.UpdateStatus(uuid.New(), dto, helperID, utils.UserTypeHelper)

	if !errors.Is(err, utils.ErrProposalNotFound) {
		t.Errorf("esperado ErrProposalNotFound, obteve: %v", err)
	}
}

// --- List tests ---

func TestProposalService_List_BusinessCallsListByUserID(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)

	listByUserIDCalled := false
	repo.ListByUserIDFn = func(uid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		listByUserIDCalled = true
		return []domain.ProposalResponseDTO{}, nil
	}
	listByHelperIDCalled := false
	repo.ListByHelperIDFn = func(hid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		listByHelperIDCalled = true
		return []domain.ProposalResponseDTO{}, nil
	}

	svc := proposalmodule.NewProposalService(repo)
	_, err := svc.List(userID, utils.UserTypeBusiness, "")

	if err != nil {
		t.Fatalf("esperado sem erro, obteve: %v", err)
	}
	if !listByUserIDCalled {
		t.Error("esperado que ListByUserID fosse chamado para business user")
	}
	if listByHelperIDCalled {
		t.Error("ListByHelperID não deveria ser chamado para business user")
	}
}

func TestProposalService_List_HelperCallsListByHelperID(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)

	listByUserIDCalled := false
	repo.ListByUserIDFn = func(uid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		listByUserIDCalled = true
		return []domain.ProposalResponseDTO{}, nil
	}
	listByHelperIDCalled := false
	repo.ListByHelperIDFn = func(hid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		listByHelperIDCalled = true
		return []domain.ProposalResponseDTO{}, nil
	}

	svc := proposalmodule.NewProposalService(repo)
	_, err := svc.List(helperID, utils.UserTypeHelper, "")

	if err != nil {
		t.Fatalf("esperado sem erro, obteve: %v", err)
	}
	if !listByHelperIDCalled {
		t.Error("esperado que ListByHelperID fosse chamado para helper user")
	}
	if listByUserIDCalled {
		t.Error("ListByUserID não deveria ser chamado para helper user")
	}
}

func TestProposalService_List_WithStatusFilter(t *testing.T) {
	userID := uuid.New()
	helperID := uuid.New()
	repo := defaultRepo(userID, helperID)

	capturedFilter := ""
	repo.ListByUserIDFn = func(uid uuid.UUID, statusFilter string) ([]domain.ProposalResponseDTO, error) {
		capturedFilter = statusFilter
		return []domain.ProposalResponseDTO{}, nil
	}

	svc := proposalmodule.NewProposalService(repo)
	svc.List(userID, utils.UserTypeBusiness, utils.ProposalStatusAccepted)

	if capturedFilter != utils.ProposalStatusAccepted {
		t.Errorf("esperado statusFilter '%s', obteve '%s'", utils.ProposalStatusAccepted, capturedFilter)
	}
}
