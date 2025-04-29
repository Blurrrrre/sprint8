package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

var (
	// randSource source of pseudo random numbers.
	// To increase uniqueness, the current time in unix format (as a number) is used as seed
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange uses randSource to generate random numbers
	randRange = rand.New(randSource)
)

// getTestParcel returns a test parcel
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete checks adding, getting and deleting a parcel
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// get
	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, gotParcel.Client)
	require.Equal(t, parcel.Status, gotParcel.Status)
	require.Equal(t, parcel.Address, gotParcel.Address)

	// delete
	err = store.Delete(id)
	require.NoError(t, err)

	// try get again
	_, err = store.Get(id)
	require.Error(t, err)
}

// TestSetAddress last address update
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, gotParcel.Address)
}

// TestSetStatus checks for status update
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)

	// set status
	newStatus := ParcelStatusDelivered
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	gotParcel, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newStatus, gotParcel.Status)
}

// TestGetByClient checks the receipt of parcels by client ID
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// we assign the same client ID to all parcels
	client := randRange.Intn(10_000_000)
	for i := range parcels {
		parcels[i].Client = client
	}

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		expected, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, expected.Client, parcel.Client)
		require.Equal(t, expected.Status, parcel.Status)
		require.Equal(t, expected.Address, parcel.Address)
	}
}
