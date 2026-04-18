package utils

import (
	"errors"
	"regexp"
	"time"
)

const (
	UserTypeHelper   = "helper"
	UserTypeBusiness = "business"

	MaxCategoriesPerUserRegistration = 3

	ProposalStatusPending    = "pending"
	ProposalStatusAccepted   = "accepted"
	ProposalStatusRefused    = "refused"
	ProposalStatusInProgress = "in progress"
	ProposalStatusCancelled  = "cancelled"
	ProposalStatusFinished   = "finished"

	OTPStatusWaiting   = "waiting"
	OTPStatusConfirmed = "confirmed"
	OTPStatusExpired   = "expired"

	OTPExpirationDuration = 30 * time.Minute
)

var (
	CPFRegex      = regexp.MustCompile(`^\d{11}$`)
	CNPJRegex     = regexp.MustCompile(`^\d{14}$`)
	TimeHHMMRegex = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

	ErrUserAlreadyRegistered     = errors.New("user already registered with this email and role")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrUserNotFound              = errors.New("user not found")
	ErrNotOwner                  = errors.New("user is not the owner of this resource")
	ErrCategoryHasLinkedServices = errors.New("category has linked services and cannot be removed")
	ErrServiceNameNotUnique      = errors.New("a service with this name already exists for this user")
	ErrHelperOnly                = errors.New("this action is restricted to helper users")
	ErrBusinessOnly              = errors.New("this action is restricted to business users")
	ErrCategoryNotAssignedToUser = errors.New("category is not assigned to this user")

	ErrInvalidValueFormat     = errors.New("invalid value format")
	ErrValueNotPositive       = errors.New("value must be greater than zero")
	ErrInvalidStartTimeFormat = errors.New("invalid start_time format, use HH:MM")
	ErrInvalidEndTimeFormat   = errors.New("invalid end_time format, use HH:MM")

	ErrServiceNotFound = errors.New("service not found")

	ErrProposalNotFound               = errors.New("proposal not found")
	ErrProposalAlreadyActiveForHelper = errors.New("user already has an active proposal for this helper")
	ErrProposalFinished               = errors.New("proposal is already in a terminal status")
	ErrProposalCannotBeRefused        = errors.New("proposal cannot be refused in its current status")
	ErrProposalInvalidStatus          = errors.New("invalid proposal status or transition")
	ErrProposalUnauthorized           = errors.New("user is not authorized to perform this status change")
	ErrNotProposalParticipant         = errors.New("user is not a participant of this proposal")

	ErrOTPNotFound   = errors.New("otp not found")
	ErrOTPExpired    = errors.New("otp has expired")
	ErrOTPInvalid    = errors.New("invalid otp code")
	ErrOTPNotWaiting = errors.New("otp is not in waiting status")

	ErrInvalidImageType = errors.New("invalid image type")

	ErrReviewAlreadyExists    = errors.New("proposal has already been reviewed")
	ErrProposalNotFinished    = errors.New("proposal must be finished before creating review")
	ErrReviewProposalMismatch = errors.New("proposal does not match business and helper")
	ErrCannotReviewSelf       = errors.New("user cannot review themselves")
)
