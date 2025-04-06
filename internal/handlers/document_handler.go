package handlers

import (
	"fmt"
	"hostel-management/pkg/helpers"
	"hostel-management/pkg/session"

	"github.com/gin-gonic/gin"
)

func DocumentsHandler(c *gin.Context) {

	const op = "handlers.DocumentsHandler.DocumentsHandler"

	role, exists := session.GetUserRole(c)
	if !exists || role != "admin" {
		c.String(403, "Access denied: %v: %v", role, op)
		return
	}

	data := map[string]interface{}{
		"Page": "admin_documents",
		"Role": role,
	}

	c.HTML(200, "layout.html", data)
}

func CreateContractHandler(c *gin.Context) {

	const op = "handlers.DocumentsHandler.CreateContractHandler"

	if c.Request.Method != "POST" {
		c.String(405, "Method not allowed: %v: %v", c.Request.Method, op)
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
	c.Data(201, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", fileBytes)
}
