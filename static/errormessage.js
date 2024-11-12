const changeErrorMessage = event => {
	document.getElementById("error").innerHTML = event.detail.xhr.responseText
};

const redirectHomepage = event => {
	var sendbtn = document.getElementById("send")
	console.log("send request")
	if (event.detail.target == sendbtn){
		if (event.detail.successful) {
			window.location.replace("/")
		}
	}
}

window.addEventListener("htmx:responseError", changeErrorMessage)
window.addEventListener("htmx:afterRequest", redirectHomepage)

function emptyErrorMessage() {
	document.getElementById("error").innerHTML = ""
}
