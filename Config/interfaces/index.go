package interfaces

import "Config/models"

type Config interface {
	AddApplication(user *models.Application) (*models.Application, error)
	UpdateConfig(user *models.Config1) (*models.Config, error)
	AddConfig(user *models.Config1) (error)
	GetConfigValue(user *models.Config1)([]*models.Config1,error)
}
