package database_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/l3lackShark/binance-ws-listener/database"
	"github.com/l3lackShark/binance-ws-listener/envvars"
	"github.com/stretchr/testify/assert"
)

func init() {
	envvars.LoadEnv()
}

func TestDatabase(t *testing.T) {
	db, err := database.New(os.Getenv("MONGO_CONN_URL"))
	ok := assert.NoError(t, err, fmt.Sprintf("database.New() error is non-nil: %e", err))
	if !ok {
		return
	}
	//connection established, try to drop the old test data and populate some new.
	err = db.RemoveAllDocumentsInCollection(database.TestDatabaseName, database.CollectionName)
	ok = assert.NoError(t, err, fmt.Sprintf("Failed to remove test collection: %e", err))
	if !ok {
		return
	}

	testCases := []database.Document{
		{
			Date:  "07.07.2020",
			Time:  "18:10",
			Price: "35326.24124",
		},
		{
			Date:  "08.07.2020",
			Time:  "18:12",
			Price: "35316.24124",
		},
		{
			Date:  "09.07.2020",
			Time:  "18:13",
			Price: "35396.24124",
		},
		{
			Date:  "07.07.2020",
			Time:  "18:15",
			Price: "37396.24124",
		},
	}

	for i, test := range testCases {
		t.Logf("Calling UpdateOrInsertOne(), case: %d", i)
		err := db.UpdateOrInsertOne(test, database.TestDatabaseName, database.CollectionName)
		ok := assert.NoError(t, err, fmt.Sprintf("Failed to pass test case suite on UpdateOrInsertOne(): %e", err))
		if !ok {
			return
		}
	}
	//test new doc

	// 	tests := []struct {
	// 	}{}

	// 	db.UpdateOrInsertOne(nil, database.TestDatabaseName, database.TestCollectionName)
	//
}
