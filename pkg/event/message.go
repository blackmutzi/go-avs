package event

import (
	"fmt"
	"github.com/satori/go.uuid"
)

func NewMessageId() string {
	return fmt.Sprintf("%s", uuid.Must(uuid.NewV4()) )
}

func NewDialogReqId() string {
	return "dialog-" + NewMessageId()
}
