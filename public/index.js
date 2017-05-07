function toggleResults(event) {
	var standing = event.target.parentNode.parentNode
	var results = standing.children[2]
	var currentDisplay = results.style.display
	console.log(currentDisplay)
	if (currentDisplay == "none" || currentDisplay == "") {
		results.style.display = "block"
	} else {
		results.style.display = "none"
	}
	
}
