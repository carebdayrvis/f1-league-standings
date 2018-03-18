package main

import (
	"github.com/carebdayrvis/go-ergast"
	"sort"
	"time"
)

type TeamCollection []Team

func (t TeamCollection) Len() int           { return len(t) }
func (t TeamCollection) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t TeamCollection) Less(i, j int) bool { return t[i].Score < t[j].Score }

type RaceCollection []Race

func (t RaceCollection) Len() int      { return len(t) }
func (t RaceCollection) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t RaceCollection) Less(i, j int) bool {
	return t[i].Race.Date.Before(time.Date(t[j].Race.Date.Year(), t[j].Race.Date.Month(), t[j].Race.Date.Day(), 0, 0, 0, 0, time.UTC))
}

func Standings(teams []Team, results []ergast.Race) []Team {

	for i, _ := range teams {
		teams[i].CalculateResults(results)
	}

	sort.Sort(sort.Reverse(TeamCollection(teams)))

	for i, _ := range teams {
		teams[i].Standing = i + 1
	}

	return teams
}

type Team struct {
	Name              string
	Drivers           []int
	Races             []Race
	Score             int
	Standing          int
	DriversByDriverID map[string]int
}

type Result struct {
	Result ergast.Result
	Score  int
}

type Race struct {
	Race             ergast.Race
	Score            int
	Results          []Result
	FastestLap       ergast.Lap
	FastestLapDriver ergast.Driver
}

func (t *Team) CalculateResults(results []ergast.Race) {

	// Iterate through each race, finding the results for this team's drivers

	driverMap := map[int]bool{}
	driverIDMap := map[string]bool{}

	for _, driver := range t.Drivers {
		driverMap[driver] = true
	}

	for driver, _ := range t.DriversByDriverID {
		driverIDMap[driver] = true
	}

	teamRaces := []Race{}

	for _, race := range results {

		var fastestLap *ergast.Lap
		var fastestLapDriver ergast.Driver

		reversedScores := ReversedScores(len(results))
		teamResults := []Result{}
		raceScore := 0

		for _, result := range race.Results {

			_, hasDriverByNumber := driverMap[result.Driver.PermanentNumber]
			_, hasDriverByID := driverIDMap[result.Driver.DriverID]

			if !hasDriverByNumber && !hasDriverByID {
				continue
			}

			r := Result{
				Result: result,
				Score:  reversedScores[result.Position-1],
			}

			if r.Result.FastestLap.Rank == 1 {
				fastestLap = &r.Result.FastestLap
				fastestLapDriver = r.Result.Driver
			}

			raceScore += r.Score

			teamResults = append(teamResults, r)
		}

		t.Score += raceScore

		// zero out results
		race.Results = []ergast.Result{}
		teamRace := Race{
			Race:             race,
			Score:            raceScore,
			Results:          teamResults,
			FastestLapDriver: fastestLapDriver,
		}

		if fastestLap != nil {
			teamRace.FastestLap = *fastestLap
		}

		teamRaces = append(teamRaces, teamRace)
	}

	t.Races = teamRaces

	sort.Sort(sort.Reverse(RaceCollection(t.Races)))
}

func ReversedScores(numResults int) []int {
	reversedScore := numResults
	reversedScores := []int{}

	for {

		reversedScores = append(reversedScores, reversedScore)

		reversedScore--

		if reversedScore == 0 {
			break
		}
	}

	return reversedScores
}
