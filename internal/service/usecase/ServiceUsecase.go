package usecase

import (
	. "../../tools/response"
	"../repository"
	"net/http"
)

type ServiceUsecase struct {
	userRepository *repository.ServiceRepository
}

func NewServiceUsecase(serviceRepository *repository.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{
		userRepository: serviceRepository,
	}
}

func (serviceUsecase *ServiceUsecase) Clear() *Response {
	err := serviceUsecase.userRepository.Truncate()
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusOK, nil)
}

func (serviceUsecase *ServiceUsecase) Status() *Response {
	status, err := serviceUsecase.userRepository.Select()
	if err != nil {
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusOK, status)
}
