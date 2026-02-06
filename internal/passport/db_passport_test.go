package passport

import (
	"context"
	"testing"

	"github.com/leeprovoost/go-rest-api-template/internal/passport/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPassportsByUser(t *testing.T) {
	srv := NewTestServer()
	passports, err := srv.passportStore.ListPassportsByUser(context.Background(), 0)
	require.NoError(t, err)
	assert.Len(t, passports, 1)
	assert.Equal(t, "012345678", passports[0].ID)
}

func TestListPassportsByUserNoResults(t *testing.T) {
	srv := NewTestServer()
	passports, err := srv.passportStore.ListPassportsByUser(context.Background(), 999)
	require.NoError(t, err)
	assert.Empty(t, passports)
}

func TestGetPassport(t *testing.T) {
	srv := NewTestServer()
	p, err := srv.passportStore.GetPassport(context.Background(), "012345678")
	require.NoError(t, err)
	assert.Equal(t, "HMPO", p.Authority)
	assert.Equal(t, 0, p.UserID)
}

func TestGetPassportFail(t *testing.T) {
	srv := NewTestServer()
	_, err := srv.passportStore.GetPassport(context.Background(), "nonexistent")
	assert.Error(t, err)
}

func TestAddPassport(t *testing.T) {
	srv := NewTestServer()
	p := models.Passport{
		ID:        "555666777",
		Authority: "IPS",
		UserID:    0,
	}
	created, err := srv.passportStore.AddPassport(context.Background(), p)
	require.NoError(t, err)
	assert.Equal(t, "555666777", created.ID)

	// Verify it's in the store
	fetched, err := srv.passportStore.GetPassport(context.Background(), "555666777")
	require.NoError(t, err)
	assert.Equal(t, "IPS", fetched.Authority)
}

func TestAddPassportDuplicate(t *testing.T) {
	srv := NewTestServer()
	p := models.Passport{
		ID:        "012345678", // already exists
		Authority: "HMPO",
		UserID:    0,
	}
	_, err := srv.passportStore.AddPassport(context.Background(), p)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestUpdatePassport(t *testing.T) {
	srv := NewTestServer()
	p := models.Passport{
		ID:        "012345678",
		Authority: "IPS",
		UserID:    0,
	}
	updated, err := srv.passportStore.UpdatePassport(context.Background(), p)
	require.NoError(t, err)
	assert.Equal(t, "IPS", updated.Authority)

	// Verify persisted
	fetched, err := srv.passportStore.GetPassport(context.Background(), "012345678")
	require.NoError(t, err)
	assert.Equal(t, "IPS", fetched.Authority)
}

func TestUpdatePassportFail(t *testing.T) {
	srv := NewTestServer()
	p := models.Passport{
		ID:        "nonexistent",
		Authority: "IPS",
	}
	_, err := srv.passportStore.UpdatePassport(context.Background(), p)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeletePassport(t *testing.T) {
	srv := NewTestServer()
	err := srv.passportStore.DeletePassport(context.Background(), "012345678")
	assert.NoError(t, err)

	// Verify deleted
	_, err = srv.passportStore.GetPassport(context.Background(), "012345678")
	assert.Error(t, err)
}

func TestDeletePassportFail(t *testing.T) {
	srv := NewTestServer()
	err := srv.passportStore.DeletePassport(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
