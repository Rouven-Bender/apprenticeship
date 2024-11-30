const redirectHomepage = event => {
	var sendbtn = document.getElementById("send")
	if (event.detail.target == sendbtn) {
		if (event.detail.successful) {
			window.location.replace("/home")
		}
	}
}

window.addEventListener("htmx:afterRequest", redirectHomepage)
