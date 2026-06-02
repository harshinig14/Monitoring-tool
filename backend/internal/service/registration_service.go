package service

import (
	"MONITORING-TOOL/internal/models"
	"MONITORING-TOOL/internal/repository"
)

type RegistrationService struct {
	userRepo *repository.UserRepository
}

func NewRegistrationService(userRepo *repository.UserRepository) *RegistrationService {
	return &RegistrationService{userRepo: userRepo}
}

func (s *RegistrationService) RegisterDevice(req models.RegisterRequest) (int, error) {
	user, err := s.userRepo.GetUserByMachineName(req.MachineName)
	if err != nil {
		return 0, err
	}

	// Machine exists -> return existing UserID
	if user != nil {
		return user.UserID, nil
	}

	// Machine doesn't exist -> create user
	userID, err := s.userRepo.CreateUser(req.Username, req.MachineName, req.OSType)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
