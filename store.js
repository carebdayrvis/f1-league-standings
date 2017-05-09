const pug = require("pug")
const fs = require("fs")
const api = require("./api")
const teams = JSON.parse(fs.readFileSync("./teams.json").toString())





async function run() {
	try {
		let sched = await api.schedule(2017)
		let results = await api.results(sched)
		let fantasyStandings = standings(results, teams)

		let html = pug.renderFile("index.pug", {standings: fantasyStandings})

		fs.writeFileSync("./index.html", html)
	} catch(e) {
		console.error("error", e)
	}
}

function reverseResults(teams) {
	let numberOfDrivers = 20
	//let numberOfDrivers = teams.reduce((a, b) => {
	//	console.log(a, b)
	//	return {drivers: a.drivers.concat(b.drivers)}
	//}, {drivers: []}).drivers.length + 1


	let ret = []
	while (numberOfDrivers) {
		ret.push(numberOfDrivers)
		numberOfDrivers--
	}

	return ret


}

function pointsSort(a, b) {
	if (a.points > b.points) return 0
	if (b.points > a.points) return 1
	else return -1
}

function standings(results, teams) {
	let reversedResults = reverseResults(teams)

	return teams.map(t => {
		let pointsSoFar = 0
		t.results = []
		results.forEach(r => {
			let race = r.RaceName[0]
			let racePoints = 0
			let teamResult = []
			r.ResultsList[0].Result.forEach((result, i) => {
				let driver = result.Driver[0]
				let driverNum = parseInt(driver.PermanentNumber[0])
				if (t.drivers.indexOf(driverNum) > -1) {
					teamResult.push({driver: {name: driver["$"].code, num: driverNum}, position: i + 1, points: reversedResults[i]})
					pointsSoFar += reversedResults[i]
					racePoints += reversedResults[i]
				}
			})
			t.results.push({race: race, points: racePoints, result: teamResult.sort(pointsSort)})
		})

		t.results = t.results.reverse()
		
		t.points = pointsSoFar

		return t
	}).sort(pointsSort)

}


run()
