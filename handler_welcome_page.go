package main

import "net/http"

func handleWelcomePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	const page = `<html>
	<head></head>
	<body>
		<p> WELCOME PAGE from GO server. </p>
	</body>
	</html>
	`
	w.Write([]byte(page))
}
