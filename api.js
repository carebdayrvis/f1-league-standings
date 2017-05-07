const moment = require("moment")
const redis = require("redis")
const xml = require("xml2js")
const got = require("got")

const SCHEDULE_PREFIX = "k:schedule:" // k:schedule:YEAR
const RESULT_PREFIX = "k:result:" // k:result:YEAR:RACENUM


module.exports = class Api {
	static get host() {
		return "http://ergast.com/api/f1"
	}


	static async schedule(year) {
		let key = `k:schedule:${year}`
		try {
			return await Api.cached(key)
		} catch (e) {
			return new Promise((resolve, reject) => {
		      		let url = `${Api.host}/${year}`
		      		got(url).then(res => {
					xml.parseString(res.body, (err, val) => {
						if (err) return reject(err)
						let ret = val.MRData.RaceTable[0].Race
						Api.cache(key, JSON.stringify(ret))
						return resolve(ret)
					})
				})
				.catch(reject)
			})
		}
	}

	static async results(sched) {
		let racesSoFar = sched.filter(s => {
			let today = moment()
			let date = moment(s.Date[0], "YYYY-MM-DD")
			return date.isBefore(today)
		})

		return await Promise.all(racesSoFar.map(r => {
			return Api.result(2017, r["$"].round)
		}))
			
	}

	static async result(year, raceNum) {
		let key = RESULT_PREFIX + `${year}:${raceNum}`
		try {
			return await Api.cached(key)
		} catch(e) {
			return new Promise((resolve, reject) => {
				let url = `${Api.host}/${year}/${raceNum}/results`
				return got(url).then(res => xml.parseString(res.body, (err, val) => {
					if (err) reject(err)
					let race = val.MRData.RaceTable[0].Race[0]
			
					Api.cache(key, JSON.stringify(race))
					return resolve(val)
				}))
				.catch(e => reject(e))
			})
		}
	}

	static cache(key, val) {
		let client = redis.createClient()
		client.set(key, val, () => client.quit())
	}

	static cached(key) {
		return new Promise((resolve, reject) => {
			let client = redis.createClient()
			client.get(key, (err, val) => {
				client.quit()
				if (err) return reject(err)
				if (!val) return reject()
				console.log("getting cached", key)
				resolve(JSON.parse(val))
			})
		})
	}


}
