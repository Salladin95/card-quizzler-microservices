package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/constants"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/entities"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/lib"
	"github.com/Salladin95/card-quizzler-microservices/mail-service/cmd/api/mail"
	"github.com/Salladin95/goErrorHandler"
	"github.com/Salladin95/rmqtools"
	"math/rand"
	"time"
)

type mailHandlers struct {
	gmailSender mail.GmailSender
	broker      rmqtools.MessageBroker
}

type MailHandlers interface {
	HandleEvent(_ string, payload []byte)
}

func NewMailHandlers(gmailSender mail.GmailSender, broker rmqtools.MessageBroker) MailHandlers {
	return &mailHandlers{gmailSender: gmailSender, broker: broker}
}

func (mh *mailHandlers) HandleEvent(_ string, payload []byte) {
	err := mh.sendEmailVerification(payload)
	if err != nil {
		fmt.Println(err)
	}

}

func (mh *mailHandlers) sendEmailVerification(payload []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var dto entities.SendEmailVerificationDto
	err := json.Unmarshal(payload, &dto)
	if err != nil {
		return goErrorHandler.OperationFailure("unmarshal dto", err)
	}

	secretCode := generateRandomSixDigitNumber()
	message := mail.EmailMessage{
		Subject: fmt.Sprintf("Confirm Your Email Address for - %s", mh.gmailSender.Name),
		Content: GenerateEmailVerificationRequestMessage(secretCode),
		To:      []string{dto.Email},
	}
	err = mh.gmailSender.SendEmail(message)

	if err != nil {
		return goErrorHandler.NewError(
			goErrorHandler.ErrInternalFailure,
			fmt.Errorf("failed to send email - %s", err.Error()),
		)
	}
	mh.broker.PushToQueue(
		ctx,
		constants.EmailVerificationCodeCommand,
		entities.EmailCode{Email: dto.Email, Code: secretCode},
	)
	lib.LogInfo(
		fmt.Sprintf(
			"verification code is sent to email - %s, code - %d", dto.Email, secretCode,
		),
	)
	return nil
}

func GenerateEmailVerificationRequestMessage(code int) string {
	return fmt.Sprintf(
		"<p>Hello there!</p>\n<p>Ready for your verification code? Here it is: <b>%d</b></p>\n<p>If you didn't light this up, feel free to ignore.</p>\n<p>Thanks for being part of our platform!</p>",
		code,
	)
}

func generateRandomSixDigitNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(900000) + 100000
}
