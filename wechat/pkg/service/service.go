package service

type BootableService interface {
	Boot() error
}
