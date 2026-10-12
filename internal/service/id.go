package service

import (
    "crypto/rand"
    "fmt"
)

func newProductID() (string, error) {
    var id [16]byte

    if _, err := rand.Read(id[:]); err != nil {
        return "", fmt.Errorf("generate product ID: %w", err)
    }

    // UUID version 4.
    id[6] = (id[6] & 0x0f) | 0x40

    // UUID variant.
    id[8] = (id[8] & 0x3f) | 0x80

    return fmt.Sprintf(
        "%08x-%04x-%04x-%04x-%012x",
        id[0:4],
        id[4:6],
        id[6:8],
        id[8:10],
        id[10:16],
    ), nil
}