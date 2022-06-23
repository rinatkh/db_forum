package core

import (
	"fmt"
	"github.com/rinatkh/db_forum/internal/constants"

	"github.com/gofrs/uuid"
)

func GenUUID() (string, error) {
	v4, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("%w: %v", constants.ErrGenerateUUID, err)
	}
	return v4.String(), nil
}
