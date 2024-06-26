package main

import (
	"log"
	"net/http"
	"text/template"
)

func renderTemplate(w http.ResponseWriter, name string, html string, data []string) {
	if template.Must(template.New(name).Parse(html)).Execute(w, data) != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
//  connMAC, err := ScanMacAddress(r)
//	if err != nil {
//		log.Println("Error [ScanMacAddress]:", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}

//	var checkFlag bool
//	allowedMACs, err := ReadAllowedMacAddresses(ALLOWED_MAC_ADDRESSES_PATH)
//	if err != nil {
//		log.Println("Error [ReadAllowedMacAddresses]:", err)
//		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
//		return
//	}

//	for _, mac := range allowedMACs {
//		if mac.String() == connMAC {
//			checkFlag = true
//			break
//		}
//	}
//	if !checkFlag {
//		log.Println("Invalid user MAC")
//		http.Error(w, "Unauthorized", http.StatusUnauthorized)
//		return
//	}
  html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV</title>
		<style>
			body {
				font-family: 'Courier New', monospace;
				background-color: #000;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background-color: #1a1a1a;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0px 2px 10px rgba(255, 255, 255, 0.1);
				max-width: 400px;
				width: 90%;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #0f0;
				font-size: 2em;
			}
			ul {
				list-style: none;
				padding: 0;
				margin-top: 2rem;
			}
			li {
				margin-bottom: 1rem;
			}
			a {
				text-decoration: none;
				color: #0f0;
				font-weight: bold;
				font-size: 1.2em;
				margin-left: 1rem;
			}
			a:hover {
				color: #0a0;
			}
			hr {
				color: #0f0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>&lt; BrewTV &gt;</h1>
			<ul>
				<li><a href="/tv">[ TV ]</a></li>
				<br>
				<li><a href="/music">[ Music ]</a></li>
				<br>
				<li><a href="/ytpl">[ YouTube ]</a></li>
			</ul>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "index", html, make([]string, 0))
}

func MusicHandler(w http.ResponseWriter, r *http.Request) {
	music, err := listFilesByType(MUSIC_PATH, ACCEPTED_MUSIC_FORMAT)
	if err != nil {
		http.Error(w, "Failed to list music", http.StatusInternalServerError)
		return
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV Music</title>
		<style>
			body {
				font-family: 'Courier New', monospace;
				background-color: #000;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background-color: #1a1a1a;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0px 2px 10px rgba(255, 255, 255, 0.1);
				max-width: 400px;
				width: 90%;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #0f0;
				font-size: 2em;
			}
			ul {
				list-style: none;
				padding: 0;
				margin-top: 20px;
			}
			li {
				margin-bottom: 10px;
			}
			a {
				text-decoration: none;
				color: #0f0;
				font-weight: bold;
				font-size: 1.2em;
				margin-left: 10px;
			}
			a:hover {
				color: #0a0;
			}
			.back-button {
				margin-top: 20px;
			}
			.back-button a {
				text-decoration: none;
				color: #0f0;
				font-weight: bold;
				font-size: 1.2em;
			}
			.back-button a:hover {
				color: #0a0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>&lt; BrewTV Music &gt;</h1>
			<ul>
				{{range .}}
				<li><a href="/music/play?path={{.}}">[ {{.}} ]</a></li>
				{{end}}
			</ul>
			<div class="back-button">
				<a href="/">[ Back ]</a>
			</div>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "music", html, music)
}

func TVHandler(w http.ResponseWriter, r *http.Request) {
	tv, err := listFilesByType(TV_PATH, ACCEPTED_TV_FORMAT)
	if err != nil {
		http.Error(w, "Failed to list tv", http.StatusInternalServerError)
		return
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV TV</title>
		<style>
			body {
				font-family: 'Courier New', monospace;
				background-color: #000;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background-color: #1a1a1a;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0px 2px 10px rgba(255, 255, 255, 0.1);
				max-width: 400px;
				width: 90%;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #0f0;
				font-size: 2em;
			}
			ul {
				list-style: none;
				padding: 0;
				margin-top: 20px;
			}
			li {
				margin-bottom: 10px;
			}
			a {
				text-decoration: none;
				color: #0f0;
				font-weight: bold;
				font-size: 1.2em;
				margin-left: 10px;
			}
			a:hover {
				color: #0a0;
			}
			.back-button {
				margin-top: 20px;
			}
			.back-button a {
				text-decoration: none;
				color: #0f0;
				font-weight: bold;
				font-size: 1.2em;
			}
			.back-button a:hover {
				color: #0a0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>&lt; BrewTV TV &gt;</h1>
			<ul>
				{{range .}}
				<li><a href="/tv/play?path={{.}}">[ {{.}} ]</a></li>
				{{end}}
			</ul>
			<div class="back-button">
				<a href="/">[ Back ]</a>
			</div>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "tv", html, tv)
}

func YTPLVideoHandler(w http.ResponseWriter, request *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV YTPL</title>
		<style>
			body {
				font-family: 'Courier New', monospace;
				background-color: #000;
				margin: 0;
				padding: 0;
				display: flex;
				justify-content: center;
				align-items: center;
				min-height: 100vh;
			}
			.container {
				background-color: #1a1a1a;
				border-radius: 8px;
				padding: 20px;
				box-shadow: 0px 2px 10px rgba(255, 255, 255, 0.1);
				max-width: 400px;
				width: 90%;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #0f0;
				font-size: 2em;
			}
			form {
				margin-top: 20px;
			}
			label {
				display: block;
				margin-bottom: 10px;
				font-weight: bold;
				color: #0f0;
			}
			.center-input {
				display: flex;
				justify-content: center;
				align-items: center;
			}
			input[type="text"] {
				width: 100%;
				padding: 10px;
				border: none;
				border-radius: 4px;
				background-color: #333;
				color: #0f0;
			}
			button[type="submit"] {
				background-color: #0f0;
				color: #000;
				border: none;
				border-radius: 4px;
				padding: 12px 20px;
				cursor: pointer;
				margin-top: 10px;
			}
			button[type="submit"]:hover {
				background-color: #0a0;
			}
			.back-button {
				margin-top: 20px;
			}
			.back-button a {
				text-decoration: none;
				color: #0f0;
				font-weight: bold;
			}
			.back-button a:hover {
				color: #0a0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>&lt; BrewTV YTPL &gt;</h1>
			<form action="/ytpl/play" method="post">
				<label for="url">URL:</label>
				<div class="center-input">
					<input type="text" id="url" name="url" placeholder="https://www.youtube.com/watch?v=$VIDEO_ID$" required>
				</div>
				<button type="submit">[ Play ]</button>
			</form>
			<div class="back-button">
				<a href="/">[ Back ]</a>
			</div>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "ytpl", html, make([]string, 0))
}
