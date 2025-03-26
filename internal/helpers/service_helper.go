package helpers

import (
	"errors"
	"hostel-management/storage/models"
	"hostel-management/storage/repositories"
)

type ServiceHelper interface {
	ValidateServiceData(name, typeService, description, is_date, is_hostel, is_phone string, amount int) error
	ValidateStatementData(name, typeStatement, date, phone, status string, hostel, amount, users_id int) error
}

type serviceHelper struct {
	statementService repositories.StatementRepository
	serviceService   repositories.ServiceRepository
}

func TranslateTypeService(serviceType string) string {
	switch serviceType {
	case "payment":
		return "Платная"
	case "free":
		return "Бесплатная"
	default:
		return "Не указан"
	}
}

func TranslateService(service models.Service) models.Service {
	service.Type = TranslateTypeService(service.Type)
	return service
}

func TranslateServiceToggles(is_date, is_hostel, is_phone string) (bool, bool, bool) {
	return is_date == "on", is_hostel == "on", is_phone == "on"
}

func (sh *serviceHelper) ValidateServiceData(name, typeService, description, is_date, is_hostel, is_phone string, amount int) error {
	if name == "" || typeService == "" || description == "" {
		return errors.New("ValidateServiceData: заполните все обязательные поля")
	}

	// fmt.Println(name, typeService, description, is_date, is_hostel, is_phone, amount)

	is_dateT, is_hostelT, is_phoneT := TranslateServiceToggles(is_date, is_hostel, is_phone)

	err := sh.serviceService.CreateService(name, typeService, description, is_dateT, is_hostelT, is_phoneT, amount)
	if err != nil {
		return errors.New("ValidateServiceData: ошибка при создании услуги")
	}

	return nil
}

func TranslateStatus(status string) string {
	switch status {
	case "awaits":
		return "Ожидает"
	case "approved":
		return "Одобрено"
	case "denied":
		return "Отклонено"
	default:
		return "Не указан"
	}
}

func (sh *serviceHelper) ValidateStatementData(name, typeStatement, date, phone, status string, hostel, amount, users_id int) error {
	if typeStatement == "" || name == "" {
		return errors.New("ValidateStatementData: заполните все обязательные поля")
	}

	if date == "" {
		date = "Не указана"
	} else if phone == "" {
		phone = "Не указан"
	} else if hostel == 0 {
		hostel = 0
	} else if amount <= 0 {
		amount = 0
	}

	if typeStatement == "Платная" {
		typeStatement = "payment"
	} else {
		typeStatement = "free"
	}

	statement := models.Statement{
		Name:     name,
		Type:     typeStatement,
		Amount:   amount,
		Date:     date,
		Phone:    phone,
		Status:   status,
		Hostel:   hostel,
		Users_id: users_id,
	}

	err := sh.statementService.CreateStatementRequest(statement)
	if err != nil {
		return errors.New("ValidateStatementData: ошибка при создании заявки")
	}

	return nil
}
