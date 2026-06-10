package utils

import (
	crand "crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID membuat string ULID unik yang aman untuk PublicID eksternal SaaS !
func GenerateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(crand.Reader, 0)
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {
		return ulid.Make().String()
	} // Fallback jika crypto/rand macet
	return id.String()
}
