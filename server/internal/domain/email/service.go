package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type Service interface {
	SendPasswordReset(email, token string) error
}

type SMTPService struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPService() *SMTPService {
	return &SMTPService{
		Host:     getEnv("SMTP_HOST", "localhost"),
		Port:     getEnv("SMTP_PORT", "1025"),
		Username: getEnv("SMTP_USERNAME", ""),
		Password: getEnv("SMTP_PASSWORD", ""),
		From:     getEnv("SMTP_FROM", "pedrohesm@gmail.com"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (s *SMTPService) SendPasswordReset(email, token string) error {
	baseURL := os.Getenv("FRONTEND_URL")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)
	
	subject := "Recuperação de Senha"
	body := s.getPasswordResetTemplate(resetURL)

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.From, email, subject, body)

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	if s.Username == "" && s.Password == "" {
		return smtp.SendMail(addr, nil, s.From, []string{email}, []byte(message))
	}

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)
	return smtp.SendMail(addr, auth, s.From, []string{email}, []byte(message))
}

func (s *SMTPService) getPasswordResetTemplate(resetURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Recuperação de Senha</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px 40px; text-align: center;">
                            <h1 style="margin: 0; color: #333333; font-size: 24px;">🔐 Recuperação de Senha</h1>
                        </td>
                    </tr>
                    
                    <!-- Content -->
                    <tr>
                        <td style="padding: 20px 40px;">
                            <p style="margin: 0 0 20px 0; color: #666666; font-size: 16px; line-height: 1.5;">
                                Olá,
                            </p>
                            <p style="margin: 0 0 20px 0; color: #666666; font-size: 16px; line-height: 1.5;">
                                Recebemos uma solicitação para redefinir a senha da sua conta. Se você não fez essa solicitação, pode ignorar este email.
                            </p>
                            <p style="margin: 0 0 30px 0; color: #666666; font-size: 16px; line-height: 1.5;">
                                Para redefinir sua senha, clique no botão abaixo:
                            </p>
                            
                            <!-- Button -->
                            <table role="presentation" style="width: 100%%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 0;">
                                        <a href="%s" style="display: inline-block; padding: 14px 40px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #ffffff; text-decoration: none; border-radius: 6px; font-size: 16px; font-weight: bold;">
                                            Redefinir Senha
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 30px 0 20px 0; color: #999999; font-size: 14px; line-height: 1.5;">
                                Ou copie e cole este link no seu navegador:
                            </p>
                            <p style="margin: 0; color: #667eea; font-size: 14px; word-break: break-all;">
                                %s
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f8f9fa; border-radius: 0 0 8px 8px;">
                            <p style="margin: 0 0 10px 0; color: #999999; font-size: 12px; line-height: 1.5;">
                                ⏱️ Este link expira em 1 hora por motivos de segurança.
                            </p>
                            <p style="margin: 0; color: #999999; font-size: 12px; line-height: 1.5;">
                                Se você não solicitou a redefinição de senha, entre em contato conosco imediatamente.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
	`, resetURL, resetURL)
}
