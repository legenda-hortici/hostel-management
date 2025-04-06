package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/nguyenthenguyen/docx"
)

type ContractData struct {
	FirstName    string
	LastName     string
	MiddleName   string
	CheckInDate  string
	CheckOutDate string
	RoomNumber   string
	Amount       string
}

func GenerateContract(contractData ContractData) ([]byte, error) {
	// Генерация договора
	r, err := docx.ReadDocxFile("web/static/dosuments/Tipovoy_dogovor_nayma_pomeschenia_obschezhitia.docx")
	if err != nil {
		log.Fatalf("Ошибка при открытии шаблона: %s", err)
	}
	defer r.Close()

	// Заменяем плейсхолдеры
	doc := r.Editable()
	doc.Replace("LastName", contractData.LastName, -1)
	doc.Replace("FirstName", contractData.FirstName, -1)
	doc.Replace("MiddleName", contractData.MiddleName, -1)
	doc.Replace("Room", contractData.RoomNumber, -1)
	doc.Replace("ID", contractData.CheckInDate[8:10], -1)
	doc.Replace("IM", contractData.CheckInDate[5:7], -1)
	doc.Replace("IY", contractData.CheckInDate[:4], -1)
	doc.Replace("OutDay", contractData.CheckOutDate[8:10], -1)
	doc.Replace("OutMonth", contractData.CheckOutDate[5:7], -1)
	doc.Replace("OutYear", contractData.CheckOutDate[:4], -1)

	// Сохраняем измененный документ
	temp := "Tipovoy_dogovor_nayma_pomeschenia_obschezhitia.docx"
	err = doc.WriteToFile(temp)
	if err != nil {
		log.Fatalf("Ошибка при сохранении документа: %s", err)
	}

	defer os.Remove(temp) // Удаляем временный файл после использования

	// Читаем файл в байты
	fileBytes, err := os.ReadFile(temp)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	return fileBytes, nil
}
