package main

import (
	"encoding/json"
	"fmt"
	"github.com/carebdayrvis/go-ergast"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	db = 2

	pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}

			_, err = conn.Do("select", db)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}

	os.Exit(m.Run())

}

func TestLoadTeams(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/teams.json")
	assert.Nil(t, err)

	e := []Team{}

	assert.Nil(t, json.Unmarshal(b, &e))

	a, err := loadTeams("./testdata/teams.json")
	assert.Nil(t, err)
	assert.Equal(t, e, a)
}

func TestLoadSchedule(t *testing.T) {

	conn, err := redis.Dial("tcp", "localhost:6379")
	assert.Nil(t, err)

	_, err = conn.Do("select", db)
	assert.Nil(t, err)

	schedule := []ergast.Race{
		ergast.Race{
			Date:    ergast.ErgastDate{time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)},
			Circuit: ergast.Circuit{"Meow"},
		},
	}

	b, err := json.Marshal(schedule)
	assert.Nil(t, err)

	_, err = conn.Do("set", fmt.Sprintf("%v%v", SCHEDULE_CACHE_PREFIX, 2018), b)
	assert.Nil(t, err)

	a, err := loadSchedule(2018)
	assert.Nil(t, err)

	assert.Equal(t, schedule, a)
}

func TestLoadResult(t *testing.T) {

	conn, err := redis.Dial("tcp", "localhost:6379")
	assert.Nil(t, err)

	_, err = conn.Do("select", db)
	assert.Nil(t, err)

	race := ergast.Race{
		Date:    ergast.ErgastDate{time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC)},
		Circuit: ergast.Circuit{"Meow"},
	}

	b, err := json.Marshal(race)
	assert.Nil(t, err)

	_, err = conn.Do("set", fmt.Sprintf("%v%v:%v", RACE_CACHE_PREFIX, 2018, 1), b)
	assert.Nil(t, err)

	a, err := loadResult(2018, 1)
	assert.Nil(t, err)

	assert.Equal(t, race, a)
}
