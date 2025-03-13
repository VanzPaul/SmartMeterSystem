package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	// Serve static files (like HTMX)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Parse the HTML template
	tmpl := template.Must(template.ParseFiles("/home/xrrt/Documents/GoProjects/SmartMeterSystem/templates/login.html"))

	// Home route with login form
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Render the login form template
		tmpl.Execute(w, nil)
	})

	// Route to handle login form submission
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Get email and password from form
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Dummy validation (replace with real logic)
		if email != "user@example.com" || password != "password123" {
			// Return error message if credentials don't match
			fmt.Fprintf(w, `<p class="error">Invalid email or password. Please try again.</p>`)
			return
		}

		// If credentials are correct, return success message
		fmt.Fprintf(w, `<p style="color:green;">Login successful!</p>`)
	})

	// Start the server
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
