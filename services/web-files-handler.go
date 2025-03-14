package services

import (
	"html/template"
	"net/http"

	"github.com/vanspaul/SmartMeterSystem/config"
	"github.com/vanspaul/SmartMeterSystem/utils"
)

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/home/xrrt/Documents/GoProjects/SmartMeterSystem/templates/home.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}

// HTMX
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/home/xrrt/Documents/GoProjects/SmartMeterSystem/templates/login.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("/home/xrrt/Documents/GoProjects/SmartMeterSystem/templates/register.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}

// Consumer Account and Balance pages
func ConsumerAccount(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/consumer/consumer_account.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}

func ConsumerBalance(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/consumer/consumer_balance.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}

// SubmitLogin submits the login form
func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Set Cache-Control header at the start, before any writes
	w.Header().Set("Cache-Control", "public, max-age=300") // 300 seconds = 5 minutes

	// Retrieve the session token cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get email from sessions map, and then login data from users map
	store := utils.GetStore(config.MaxWebUsers, config.MaxWebSessions, config.MaxMeterUsers, config.MaxMeterSessions)
	email := store.WebSessions[sessionCookie.Value]
	loginData, ok := store.WebUsers[email]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var tmpl *template.Template
	var tmplErr error

	// Dynamically load template based on the user role
	if loginData.Role == "consumer" {
		tmpl, tmplErr = template.ParseFiles("templates/consumer/consumer_dashboard.html")
	} else {
		tmpl, tmplErr = template.ParseFiles("templates/admin_dashboard.html")
	}
	if tmplErr != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	// Execute the template
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Execution error", http.StatusInternalServerError)
	}
}
