package contato

import (
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"arena-vip/app"

	"github.com/gin-gonic/gin"
)

var allowedModalidades = map[string]string{
	"bjj":        "Brazilian Jiu-Jitsu",
	"kickboxing": "Kickboxing",
	"judo":       "Judô",
	"kids":       "Turma Kids",
	"todos":      "Mais de uma modalidade",
}

func init() {
	app.GetInstance().Router.POST("/contato", handleContact)
}

func handleContact(c *gin.Context) {
	// Reject requests larger than 64 KB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64<<10)

	// Honeypot: bots fill this hidden field, humans don't
	if c.PostForm("website") != "" {
		c.JSON(http.StatusOK, gin.H{"ok": true}) // silent reject
		return
	}

	// Rate limit by client IP: 3 submissions per 10 minutes
	if isRateLimited(c.ClientIP()) {
		c.JSON(http.StatusTooManyRequests, gin.H{"ok": false, "error": "Muitas tentativas. Aguarde alguns minutos."})
		return
	}

	nome := strings.TrimSpace(c.PostForm("nome"))
	telefone := strings.TrimSpace(c.PostForm("telefone"))
	email := strings.TrimSpace(c.PostForm("email"))
	modalidade := strings.TrimSpace(c.PostForm("modalidade"))
	mensagem := strings.TrimSpace(c.PostForm("mensagem"))

	// Required field checks
	if nome == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Nome e email são obrigatórios."})
		return
	}

	// Basic email format check
	atIdx := strings.Index(email, "@")
	if atIdx < 1 || !strings.Contains(email[atIdx:], ".") {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Email inválido."})
		return
	}

	// Length limits
	if len(nome) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Nome demasiado longo."})
		return
	}
	if len(email) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Email demasiado longo."})
		return
	}
	if len(mensagem) > 5000 {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Mensagem demasiado longa (máx. 5000 caracteres)."})
		return
	}

	// Modalidade whitelist (empty is allowed — field is optional)
	modalidadeLabel := allowedModalidades[modalidade]
	if modalidade != "" && modalidadeLabel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "Modalidade inválida."})
		return
	}

	if err := sendContactEmail(nome, telefone, email, modalidadeLabel, mensagem); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": "Erro ao enviar mensagem. Tente novamente."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func sendContactEmail(nome, telefone, email, modalidade, mensagem string) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	contactTo := os.Getenv("CONTACT_TO")

	if smtpUser == "" || smtpPass == "" || contactTo == "" {
		return fmt.Errorf("SMTP not configured")
	}

	subject := fmt.Sprintf("Novo contacto Arena VIP: %s", nome)
	htmlBody := buildEmailHTML(nome, telefone, email, modalidade, mensagem)

	msg := []byte(strings.Join([]string{
		"From: Arena VIP <" + smtpUser + ">",
		"To: " + contactTo,
		"Reply-To: " + nome + " <" + email + ">",
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"",
		htmlBody,
	}, "\r\n"))

	auth := smtp.PlainAuth("", smtpUser, smtpPass, "smtp.gmail.com")
	return smtp.SendMail("smtp.gmail.com:587", auth, smtpUser, []string{contactTo}, msg)
}
