package incident_test

import (
	"log"
	"os"
	"testing"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/involvestecnologia/statuspage/db"
	"github.com/involvestecnologia/statuspage/incident"
	"github.com/involvestecnologia/statuspage/mock"
	"github.com/involvestecnologia/statuspage/models"
	"github.com/stretchr/testify/assert"
)

var testSession *mgo.Session
var failureSession *mgo.Session
var dt, _ = time.Parse(time.RFC3339, time.Date(2018, time.March, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339))
var i = models.Incident{
	ComponentRef: mock.ZeroTimeHex,
	Status:       models.IncidentStatusOutage,
	Description:  "",
	Date:         dt,
}
var c = models.Component{
	Ref:     mock.ZeroTimeHex,
	Name:    "Component",
	Address: "",
}

func init() {
	testSession = db.InitMongo(os.Getenv("MONGO_URI"))
	err := testSession.DB("status").DropDatabase()
	if err != nil {
		log.Panicf("%s\n", err)
	}
	testSession.DB("status").C("Component").Insert(c)
}

func TestIncidentMongoDB_Repository_NewMongoRepository(t *testing.T) {
	repo := incident.NewMongoRepository(testSession)
	assert.Implements(t, (*incident.Repository)(nil), repo)
}

func TestIncidentMongoDB_Repository_Insert(t *testing.T) {

	repo := incident.NewMongoRepository(testSession)

	err := repo.Insert(i)
	assert.Nil(t, err)

	incidents, err := repo.Find(map[string]interface{}{"component_ref": i.ComponentRef})
	if assert.Nil(t, err) && assert.NotNil(t, incidents) {
		assert.Equal(t, []models.Incident{i}, incidents)
	}

	repo = incident.NewMongoRepository(failureSession)

	err = repo.Insert(i)
	assert.NotNil(t, err)

}

func TestIncidentMongoDB_Repository_Find(t *testing.T) {
	repo := incident.NewMongoRepository(testSession)

	incidents, err := repo.Find(map[string]interface{}{"component_ref": c.Ref})
	if assert.Nil(t, err) && assert.NotNil(t, incidents) {
		assert.Equal(t, []models.Incident{i}, incidents)
	}

	incidents, err = repo.Find(map[string]interface{}{"component_ref": "Invalid Ref"})
	assert.NotNil(t, err)
	assert.Nil(t, incidents)

	incidents, err = repo.Find(map[string]interface{}{"invalidQuery": "SomeValue"})
	assert.NotNil(t, err)
	assert.Nil(t, incidents)

	repo = incident.NewMongoRepository(failureSession)
	_, err = repo.Find(map[string]interface{}{"component_ref": c.Ref})
	assert.NotNil(t, err)

}

func TestIncidentMongoDB_Repository_FindOne(t *testing.T) {
	repo := incident.NewMongoRepository(testSession)

	inc, err := repo.FindOne(map[string]interface{}{"component_ref": c.Ref})
	if assert.Nil(t, err) && assert.NotNil(t, inc) {
		assert.Equal(t, i, inc)
	}

	_, err = repo.FindOne(map[string]interface{}{"component_ref": "Invalid Ref"})
	assert.NotNil(t, err)

	_, err = repo.FindOne(map[string]interface{}{"invalidQuery": "SomeValue"})
	assert.NotNil(t, err)

	repo = incident.NewMongoRepository(failureSession)
	_, err = repo.FindOne(map[string]interface{}{"component_ref": c.Ref})
	assert.NotNil(t, err)

}

func TestIncidentMongoDB_Repository_List(t *testing.T) {
	repo := incident.NewMongoRepository(testSession)

	startDt, err := time.Parse(time.RFC3339, time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339))
	assert.Nil(t, err)
	endDt, err := time.Parse(time.RFC3339, time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339))
	assert.Nil(t, err)
	unresolved := false

	incidents, err := repo.List(startDt, endDt, unresolved)
	if assert.Nil(t, err) && assert.NotNil(t, incidents) {
		assert.Equal(t, []models.Incident{i}, incidents)
	}

	incidents, err = repo.List(startDt, endDt, true)
	if assert.Nil(t, err) && assert.NotNil(t, incidents) {
		for _, i := range incidents {
			assert.False(t, i.Resolved)
		}
	}

	endDt, err = time.Parse(time.RFC3339, time.Date(2018, time.January, 2, 0, 0, 0, 0, time.UTC).Format(time.RFC3339))
	assert.Nil(t, err)
	incidents, err = repo.List(startDt, endDt, unresolved)
	if assert.Nil(t, err) && assert.Nil(t, incidents) {
		assert.IsType(t, []models.Incident{}, incidents)
	}

	endDt, err = time.Parse(time.RFC3339, time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339))
	assert.Nil(t, err)
	repo = incident.NewMongoRepository(failureSession)
	_, err = repo.List(startDt, endDt, unresolved)
	assert.NotNil(t, err)
}

func TestIncidentMongoDB_Repository_Update(t *testing.T) {

	repo := incident.NewMongoRepository(testSession)

	i.Status = models.IncidentStatusOK
	i.Resolved = true
	err := repo.Update(i)
	assert.Nil(t, err)

	incidents, err := repo.FindOne(map[string]interface{}{"component_ref": i.ComponentRef})
	if assert.Nil(t, err) && assert.NotNil(t, incidents) {
		assert.Equal(t, i, incidents)
	}
	i2 := i
	i2.ComponentRef = "Invalid Component Ref"
	err = repo.Update(i2)
	assert.NotNil(t, err)

	repo = incident.NewMongoRepository(failureSession)

	err = repo.Insert(i)
	assert.NotNil(t, err)

}
