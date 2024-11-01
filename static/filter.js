function filter() {
	filterOnlyActiv = document.getElementById("onlyActiv").checked+"" // to cast to string
	searchQuery = document.getElementById("search").value.toUpperCase()

	table = document.getElementById("data-table")
	tr = table.getElementsByTagName("tr")

	for (i = 1; i < tr.length; i++) { // 1 to skip the headers
		td = tr[i].getElementsByTagName("td")
		name = td[1].textContent.toUpperCase()
		activ = td[5].textContent.toLowerCase()
		if (searchQuery.length > 0 && (filterOnlyActiv == "true")){
			if (name.indexOf(searchQuery) > -1 && activ == "true") {
				tr[i].style.display="";
			} else {
				tr[i].style.display="none";
			}
		}
		if (searchQuery.length > 0 && !(filterOnlyActiv == "true")){
			if (name.indexOf(searchQuery) > -1) {
				tr[i].style.display="";
			} else {
				tr[i].style.display="none";
			}
		}
		if (searchQuery.length == 0 && (filterOnlyActiv == "true")){
			if (activ == "true") {
				tr[i].style.display="";
			} else {
				tr[i].style.display="none";
			}
		}
		if (searchQuery.length == 0 && !(filterOnlyActiv == "true")){
			tr[i].style.display="";
		}
	}
}
