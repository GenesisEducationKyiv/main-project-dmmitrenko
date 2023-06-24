package controller

import (
	constants "CurrencyRateApp/domain"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type IEmailService interface {
	AddEmail(email string) error
	SendRateForSubscribedEmails(coin string, currency string) error
}

type EmailController struct {
	emailService IEmailService
}

func NewEmailController(emailService IEmailService) *EmailController {
	return &EmailController{
		emailService: emailService,
	}
}

// SubscribeEmail godoc
// @Summary Subscribe email to receive the current rate
// @Description Subscribe email
// @Tags subscription
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email address"
// @Success 200
// @Failure 400
// @Router /email [post]
func (r *EmailController) SubscribeEmail(c *gin.Context) {
	email := c.PostForm("email")

	if !isValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address."})
		return
	}

	err := r.emailService.AddEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email is successfully subscribed to the newsletter."})
}

// SubscribeEmail godoc
// @Summary Send an email with the current rate to all subscribed emails.
// @Description Send an emails
// @Tags subscription
// @Produce json
// @Success 200
// @Failure 500
// @Router /subscribe [post]
func (r *EmailController) SendEmails(c *gin.Context) {

	err := r.emailService.SendRateForSubscribedEmails(constants.BITCOIN, constants.UAH)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Letters sent successfully."})
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}
