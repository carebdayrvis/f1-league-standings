package main

import (
	"fmt"
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
	FastestLapDuration string
	Result             ergast.Result
	QualifyingPosition int
	Score              int
}

type Race struct {
	Race             ergast.Race
	Score            int
	Results          []Result
	FastestLapPoint  bool
	LastInQualyPoint bool
	Modifier         string
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

		qualyByDriver := map[string]int{}

		fastestLapPoint := false
		lastInQualyPoint := false

		reversedScores := ReversedScores(len(race.Results))

		teamResults := []Result{}
		raceScore := 0

		for _, qresult := range race.QualifyingResults {
			_, hasDriverByNumber := driverMap[qresult.Driver.PermanentNumber]
			_, hasDriverByID := driverIDMap[qresult.Driver.DriverID]

			if !hasDriverByNumber && !hasDriverByID {
				continue
			}

			qualyByDriver[qresult.Driver.DriverID] = qresult.Position

			if qresult.Position == len(race.Results) {
				lastInQualyPoint = true
				raceScore += 1
			}
		}

		for _, result := range race.Results {

			_, hasDriverByNumber := driverMap[result.Driver.PermanentNumber]
			_, hasDriverByID := driverIDMap[result.Driver.DriverID]

			if !hasDriverByNumber && !hasDriverByID {
				continue
			}

			r := Result{
				Result:             result,
				Score:              reversedScores[result.Position-1],
				FastestLapDuration: FormatDuration(result.FastestLap.Time),
				QualifyingPosition: qualyByDriver[result.Driver.DriverID],
			}

			if r.Result.FastestLap.Rank == 1 {
				fastestLapPoint = true
				raceScore += 1
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
			FastestLapPoint:  fastestLapPoint,
			LastInQualyPoint: lastInQualyPoint,
		}

		teamRace.Modifier = FormatModifier(teamRace)

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

func FormatModifier(r Race) string {

	if !r.FastestLapPoint && !r.LastInQualyPoint {
		return ""
	}

	if !r.FastestLapPoint && r.LastInQualyPoint {
		return "(+1 for last in qualifying)"
	} else if r.FastestLapPoint && !r.LastInQualyPoint {
		return "(+1 for fastest lap)"
	} else if r.FastestLapPoint && r.LastInQualyPoint {
		return "(+1 for fastest lap, and last in qualifying)"
	}

	return ""
}

func FormatDuration(e ergast.ErgastDuration) string {

	d, _ := time.ParseDuration(e.String())

	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	d -= s * time.Second
	ms := d / time.Millisecond

	return fmt.Sprintf("%02d:%02d.%02d", m, s, ms)

}
