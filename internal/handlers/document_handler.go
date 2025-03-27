package handlers

import (
	"fmt"
	"hostel-management/internal/helpers"
	"hostel-management/internal/session"

	"github.com/gin-gonic/gin"
)

func DocumentsHandler(c *gin.Context) {
	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied")
		return
	}

	data := map[string]interface{}{
		"Page": "admin_documents",
		"Role": role,
	}

	c.HTML(200, "layout.html", data)
}

func CreateContractHandler(c *gin.Context) {
	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed")
		return
	}

	data := helpers.ContractData{
		FirstName:    c.PostForm("firstName"),
		LastName:     c.PostForm("lastName"),
		MiddleName:   c.PostForm("middleName"),
		CheckInDate:  c.PostForm("checkInDate"),
		CheckOutDate: c.PostForm("checkOutDate"),
		RoomNumber:   c.PostForm("roomNumber"),
		Amount:       c.PostForm("amount"),
	}

	fileBytes, err := helpers.GenerateContract(data)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// Отправляем файл в ответе
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	c.Header("Content-Disposition", "attachment; filename=Договор.docx")
	c.Header("Content-Length", fmt.Sprintf("%d", len(fileBytes)))
	c.Data(200, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", fileBytes)
}
