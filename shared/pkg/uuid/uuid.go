package uuid

import "github.com/google/uuid"

func New() (string, error) {
    id, err := uuid.NewV7()
    if err != nil {
        return "", err
    }
    return id.String(), nil
}

func MustNew() string {
    id, err := uuid.NewV7()
    if err != nil {
        panic("failed to generate UUIDv7: " + err.Error())
    }
    return id.String()
}