const express = require("express")
const redis = require("redis")

const client = redis.createClient()
const app = express()

app.set("view engine", "pug")


app.use(express.static("public"))
app.use("/", (req, res, next) => {

	client.get("k:standings", (err, val) => {
		if (err) next(err)

		let standings = JSON.parse(val)
console.log(standings[0].results)
		res.render("index", {standings: standings})

	})

})

app.listen(4568)




