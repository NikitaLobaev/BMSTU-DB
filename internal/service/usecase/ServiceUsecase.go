package usecase

import (
	"github.com/NikitaLobaev/BMSTU-DB/internal/service/repository"
	. "github.com/NikitaLobaev/BMSTU-DB/internal/tools/response"
	"github.com/labstack/gommon/log"
	"net/http"
)

type ServiceUsecase struct {
	serviceRepository repository.ServiceRepository
}

func NewServiceUsecase(serviceRepository repository.ServiceRepository) *ServiceUsecase {
	return &ServiceUsecase{
		serviceRepository: serviceRepository,
	}
}

func (serviceUsecase *ServiceUsecase) Clear() *Response {
	err := serviceUsecase.serviceRepository.Truncate()
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusOK, nil)
}

func (serviceUsecase *ServiceUsecase) Status() *Response {
	status, err := serviceUsecase.serviceRepository.Select()
	if err != nil {
		log.Error(err)
		return NewResponse(http.StatusServiceUnavailable, nil)
	}

	return NewResponse(http.StatusOK, status)
}
