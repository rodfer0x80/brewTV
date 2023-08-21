package main

import (
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
	ScanMacAddress(r)
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV</title>
		<style>
			body {
				font-family: Arial, sans-serif;
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
				width: 400px;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #fff;
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
				color: #007bff;
				font-weight: bold;
				font-size: 18px; /* Adjust the font size */
				margin-left: 10px; /* Add margin for spacing */
			}
			a:hover {
				color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>BrewTV</h1>
			<ul>
				<li><a href="/library">Library</a></li>
				<li><a href="/ytpl">YouTube</a></li>
			</ul>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "index", html, make([]string, 0))
}

func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	library, err := listFilesByType(LIBRARY_PATH, ACCEPTED_VIDEO_FORMAT)
	if err != nil {
		http.Error(w, "Failed to list library", http.StatusInternalServerError)
		return
	}

	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV Library</title>
		<style>
			body {
				font-family: Arial, sans-serif;
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
				width: 400px;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #fff;
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
				color: #007bff;
				font-weight: bold;
				font-size: 18px; /* Adjust the font size */
				margin-left: 10px; /* Add margin for spacing */
			}
			a:hover {
				color: #0056b3;
			}
			.back-button {
				margin-top: 20px;
			}
			.back-button a {
				text-decoration: none;
				color: #007bff;
				font-weight: bold;
			}
			.back-button a:hover {
				color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>BrewTV Library</h1>
			<ul>
				{{range .}}
				<li><a href="/library/play?path={{.}}">{{.}}</a></li>
				{{end}}
			</ul>
			<div class="back-button">
				<a href="/">Back</a>
			</div>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "library", html, library)
}

func YTPLVideoHandler(w http.ResponseWriter, request *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>BrewTV YTPL</title>
		<style>
			body {
				font-family: Arial, sans-serif;
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
				width: 400px;
				text-align: center;
			}
			h1 {
				margin-top: 0;
				color: #fff;
			}
			form {
				margin-top: 20px;
			}
			label {
				display: block;
				margin-bottom: 10px;
				font-weight: bold;
				color: #fff;
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
				color: #fff;
			}
			button[type="submit"] {
				background-color: #007bff;
				color: #fff;
				border: none;
				border-radius: 4px;
				padding: 12px 20px;
				cursor: pointer;
				margin-top: 10px;
			}
			button[type="submit"]:hover {
				background-color: #0056b3;
			}
			.back-button {
				margin-top: 20px;
			}
			.back-button a {
				text-decoration: none;
				color: #007bff;
				font-weight: bold;
			}
			.back-button a:hover {
				color: #0056b3;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>BrewTV YTPL</h1>
			<form action="/ytpl/play" method="post">
				<label for="url">URL:</label>
				<div class="center-input">
					<input type="text" id="url" name="url" placeholder="https://www.youtube.com/watch?v=0123456789" required>
				</div>
				<button type="submit">Play</button>
			</form>
			<div class="back-button">
				<a href="/">Back</a>
			</div>
		</div>
	</body>
	</html>	
	`
	renderTemplate(w, "ytpl", html, make([]string, 0))
}
