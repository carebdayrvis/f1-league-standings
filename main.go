package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/carebdayrvis/go-ergast"
	"github.com/garyburd/redigo/redis"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const SCHEDULE_CACHE_PREFIX string = "k:schedule:"
const RACE_CACHE_PREFIX string = "k:race:"

var pool *redis.Pool
var db int = 1

func main() {

	season := flag.Int("season", 2018, "season for which to generate results")
	teamsPath := flag.String("teams", "teams.json", "path to teams.json")
	flag.Parse()

	// Setup redis pool
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

	results, err := loadResults(*season)
	if err != nil {
		log.Fatal(err)
	}

	teams, err := loadTeams(*teamsPath)
	if err != nil {
		log.Fatal(err)
	}

	teams = Standings(teams, results)

	// Load template

	tplB, err := ioutil.ReadFile("views/standings.html")
	if err != nil {
		log.Fatal(err)
	}

	t, err := template.New("standings").Parse(string(tplB))
	if err != nil {
		log.Fatal(err)
	}

	page := struct {
		Teams TeamCollection
	}{
		Teams: teams,
	}

	log.Fatal(t.Execute(os.Stdout, page))

}

func loadTeams(path string) ([]Team, error) {
	teams := []Team{}

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return teams, err
	}

	return teams, json.Unmarshal(b, &teams)
}

func loadResults(season int) ([]ergast.Race, error) {

	// Load results for the schedule, caching them
	schedule, err := ergast.SeasonSchedule(season)
	if err != nil {
		return []ergast.Race{}, err
	}

	results := []ergast.Race{}

	// Get results for each race if it has happened
	for _, race := range schedule {

		if race.Date.After(time.Now()) {
			continue
		}

		result, err := loadResult(race.Season, race.Round)
		if err != nil {
			return []ergast.Race{}, err
		}

		results = append(results, result)
	}

	return results, nil

}

func loadSchedule(season int) ([]ergast.Race, error) {
	// Check cache for schedule
	conn := pool.Get()

	defer conn.Close()

	cacheKey := fmt.Sprintf("%v%v", SCHEDULE_CACHE_PREFIX, season)

	b, err := redis.Bytes(conn.Do("get", cacheKey))
	if err != nil && err != redis.ErrNil {
		return []ergast.Race{}, err
	} else if err == nil {
		schedule := []ergast.Race{}
		return schedule, json.Unmarshal(b, &schedule)
	}

	schedule, err := ergast.SeasonSchedule(season)
	if err != nil {
		return []ergast.Race{}, err
	}

	b, err = json.Marshal(schedule)
	if err != nil {
		return []ergast.Race{}, err
	}

	_, err = conn.Do("set", cacheKey, b)
	if err != nil {
		return []ergast.Race{}, err
	}

	return schedule, nil
}

func loadResult(season int, round int) (ergast.Race, error) {
	// Check cache for result
	conn := pool.Get()

	defer conn.Close()

	cacheKey := fmt.Sprintf("%v%v:%v", RACE_CACHE_PREFIX, season, round)

	b, err := redis.Bytes(conn.Do("get", cacheKey))
	if err != nil && err != redis.ErrNil {
		return ergast.Race{}, err
	} else if err == nil {
		race := ergast.Race{}
		return race, json.Unmarshal(b, &race)
	}

	race, err := ergast.SpecificResult(season, round)
	if err != nil {
		return ergast.Race{}, err
	}

	qualy, err := ergast.SpecificQualifying(season, round)
	if err != nil {
		return ergast.Race{}, err
	}

	race.QualifyingResults = qualy.QualifyingResults

	b, err = json.Marshal(race)
	if err != nil {
		return ergast.Race{}, err
	}

	_, err = conn.Do("set", cacheKey, b)
	if err != nil {
		return ergast.Race{}, err
	}

	return race, nil
}
