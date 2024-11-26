package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/microcosm-cc/bluemonday"
)

// Error handling for GetProjectRoot
func GetProjectRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}
	return cwd
}

// Refactored template rendering function
func renderTemplate(w http.ResponseWriter, tmplName string) {
	templatePath := filepath.Join(GetProjectRoot(), tmplName)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	// No need for RouteChecker middleware
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, // Directly using the mux without middleware
	}

	fmt.Println("server running @http://localhost:8080\n=====================================")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error running server: ", err)
	}
}

func RegisterRoutes(mux *http.ServeMux) {
	staticDir := GetProjectRoot()

	// Serve static assets (css, js, img, lib) from the correct folder
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(staticDir, "css")))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(staticDir, "js")))))
	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(filepath.Join(staticDir, "img")))))
	mux.Handle("/favicon/", http.StripPrefix("/favicon/", http.FileServer(http.Dir(filepath.Join(staticDir, "favicon")))))
	// mux.Handle("/scss/", http.StripPrefix("/scss/", http.FileServer(http.Dir(filepath.Join(staticDir, "scss")))))

	// Register the home page handler and email handler
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/resume", handleresume)
	mux.HandleFunc("/send-email", handleEmailSend)
}

// HomeHandler handles the home page request
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		renderTemplate(w, "index.html")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Contact page handler
func handleresume(w http.ResponseWriter, r *http.Request) {
	// Path to the PDF file
	filePath := "resume/SAMUELOKOTHOMULORESUME.pdf"

	// Set the appropriate headers for PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\"SAMUELOKOTHOMULORESUME.pdf\"")

	// Serve the PDF file
	http.ServeFile(w, r, filePath)
}

// // Contact page handler
// func handleContact(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "contact.html")
// }

// // Contact page handler
// func handleaward(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "award.html")
// }

// // About page handler
// func handleAbout(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "about.html")
// }

// // Blogs page handler
// func handleBlogs(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "blogs.html")
// }

// // Experience page handler
// func handleExperience(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "experience.html")
// }

// // My Work page handler
// func handleMyWork(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "mywork.html")
// }

// // Skills page handler
// func handleSkills(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "skills.html")
// }

// // Testimonial page handler
// func handleTestimonial(w http.ResponseWriter, r *http.Request) {
// 	renderTemplate(w, "testimonial.html")
// }

func handleEmailSend(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Initialize the bluemonday policy (sanitize HTML content)
	p := bluemonday.UGCPolicy()

	// Sanitize form fields
	name := p.Sanitize(r.FormValue("name"))
	email := p.Sanitize(r.FormValue("email"))
	subject := p.Sanitize(r.FormValue("subject"))
	message := p.Sanitize(r.FormValue("message"))

	// Validate Name (letters, spaces, hyphens, apostrophes, accents)
	nameRegex := `^[a-zA-Z\sÀ-ÿ'-]+$`
	if !regexp.MustCompile(nameRegex).MatchString(name) {
		http.Error(w, "Invalid name", http.StatusBadRequest)
		return
	}

	// Validate Email
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if !regexp.MustCompile(emailRegex).MatchString(email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Validate Subject (Alphanumeric and space, 1-100 characters)
	subjectRegex := `^[a-zA-Z0-9\s]+$`
	if !regexp.MustCompile(subjectRegex).MatchString(subject) {
		http.Error(w, "Invalid subject", http.StatusBadRequest)
		return
	}

	// Validate Message Length (1-500 characters)
	messageRegex := `^[a-zA-Z0-9\s.,!?-]+$`

	if !regexp.MustCompile(messageRegex).MatchString(message) {
		http.Error(w, "Invalid message", http.StatusBadRequest)
		return
	}

	// Check if all fields are provided
	if name == "" || email == "" || subject == "" || message == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Configure email settings
	from := "mcomulosammy37@gmail.com"
	password := "inss cfcv agtz njhn" // Use an app-specific password for Gmail
	to := "mcomulosammy37@gmail.com"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Compose the email with HTML content
	emailSubject := "New Message From Your Portfolio: " + subject
	emailBody := fmt.Sprintf(`
	<html>
	<head>
		<title>Someone has contucted you from your portfolio</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				color: #333;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				background: #f8f8f8;
				border-radius: 8px;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
			}
			.header {
				background-color: #2196f3;
				color: white;
				padding: 10px 20px;
				text-align: center;
				border-radius: 8px 8px 0 0;
			}
			.footer {
				background-color: #f1f1f1;
				padding: 10px;
				text-align: center;
				margin-top: 20px;
				font-size: 12px;
				color: #888;
				border-radius: 0 0 8px 8px;
			}
			.content {
				padding: 20px;
				background-color: white;
				border-radius: 8px;
			}
			.content p {
				font-size: 16px;
				line-height: 1.5;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="header">
				<h2>New Message From Your Portfolio</h2>
			</div>
			<div class="content">
				<p><strong>Name:</strong> %s</p>
				<p><strong>Email:</strong> %s</p>
				<p><strong>Subject:</strong> %s</p>
				<p><strong>Message:</strong><br>%s</p>
			</div>
			<div class="footer">
				<p>&copy; 2024 somuloportfolio. All rights reserved.</p>
			</div>
		</div>
	</body>
	</html>
	`, name, email, subject, message)

	// Compose the email message
	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, emailSubject, emailBody))

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Respond back with a success message in HTML format
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Send the confirmation HTML page
	fmt.Fprintf(w, `
	<html>
    <head>
        <title>Message Confirmation</title>
        <link rel="stylesheet" href="css/popup_submeet.css">
    </head>
    <body>
        <div class="confirmation custom-confirmation">
            <h2>Your Message Has Been Sent!</h2>
            <p>Thank you, %s, for your message. The following details have been received:</p>
            <ul>
                <li><strong>Name:</strong> %s</li>
                <li><strong>Email:</strong> %s</li>
                <li><strong>Subject:</strong> %s</li>
                <li><strong>Message:</strong> %s</li>
            </ul>
            <button class="back-button custom-back-button" onclick="javascript:history.back()">Go Back</button>
			<div class="footer">
			<p align="center" >&copy; 2024 samuel omulo. All rights reserved.</p>
	</div>
        </div>
		
    </body>
	

</html>
	`, name, name, email, subject, message)
}

// func handleAppointment(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	// Parse form data
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, "Error parsing form", http.StatusBadRequest)
// 		return
// 	}

// 	// Initialize the bluemonday policy (this will sanitize HTML content)
// 	p := bluemonday.UGCPolicy() // UGCPolicy allows a safe subset of HTML tags and attributes

// 	// Sanitize form fields
// 	name := p.Sanitize(r.FormValue("name"))
// 	email := p.Sanitize(r.FormValue("email"))
// 	mobile := p.Sanitize(r.FormValue("mobile"))
// 	service := p.Sanitize(r.FormValue("service"))
// 	date := p.Sanitize(r.FormValue("date"))
// 	time := p.Sanitize(r.FormValue("time"))
// 	message := p.Sanitize(r.FormValue("message"))

// 	// Define regex patterns for validation
// 	nameRegex := `^[a-zA-Z\s]+$`                                                                                // Only letters and spaces
// 	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`                                            // Email pattern
// 	mobileRegex := `^\d{10}$`                                                                                   // Only digits, length 10-15 digits
// 	dateRegex := `^(0[1-9]|1[0-2])\/(0[1-9]|[12][0-9]|3[01])\/\d{4}$`                                           // Date in dd/mm/yyyy format
// 	timeRegex := `^(?:([01]?[0-9]|2[0-3]):([0-5][0-9])(?:\s?(AM|PM|am|pm))?|([01]?[0-9]|2[0-3]):([0-5][0-9]))$` // Time in HH:MM format
// 	messageRegex := `^.{10,500}$`                                                                               // Non-empty, max 500 characters

// 	// Validate inputs using regex
// 	if !regexp.MustCompile(nameRegex).MatchString(name) {
// 		http.Error(w, "Invalid name. Only letters and spaces are allowed.", http.StatusBadRequest)
// 		return
// 	}

// 	if !regexp.MustCompile(emailRegex).MatchString(email) {
// 		http.Error(w, "Invalid email address.", http.StatusBadRequest)
// 		return
// 	}

// 	if !regexp.MustCompile(mobileRegex).MatchString(mobile) {
// 		http.Error(w, "Invalid mobile number. It must contain 10 digits.", http.StatusBadRequest)
// 		return
// 	}

// 	if !regexp.MustCompile(dateRegex).MatchString(date) {
// 		http.Error(w, "Invalid date. Please use the format dd/mm/yyyy.", http.StatusBadRequest)
// 		return
// 	}

// 	if !regexp.MustCompile(timeRegex).MatchString(time) {
// 		http.Error(w, "Invalid time. Please use the format HH:MM.", http.StatusBadRequest)
// 		return
// 	}

// 	if !regexp.MustCompile(messageRegex).MatchString(message) {
// 		http.Error(w, "Invalid message. It must be between 10 and 500 characters.", http.StatusBadRequest)
// 		return
// 	}

// 	// Check if all fields are provided (redundant check since regex handles the format)
// 	if name == "" || email == "" || mobile == "" || service == "" || date == "" || time == "" || message == "" {
// 		http.Error(w, "All fields are required", http.StatusBadRequest)
// 		return
// 	}

// 	// Configure email settings
// 	from := "mcomulosammy37@gmail.com"
// 	password := "inss cfcv agtz njhn" // Use an app-specific password
// 	to := "mcomulosammy37@gmail.com"
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	// Compose the HTML email
// 	emailSubject := "New Appointment Booking"
// 	emailBody := fmt.Sprintf(`
// 	<html>
// 	<head>
// 		<title>New Appointment Request</title>
// 		<style>
// 			body {
// 				font-family: Arial, sans-serif;
// 				color: #333;
// 			}
// 			.container {
// 				max-width: 600px;
// 				margin: 0 auto;
// 				padding: 20px;
// 				background: #f8f8f8;
// 				border-radius: 8px;
// 				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
// 			}
// 			.header {
// 				background-color: #2196f3;
// 				color: white;
// 				padding: 10px 20px;
// 				text-align: center;
// 				border-radius: 8px 8px 0 0;
// 			}
// 			.footer {
// 				background-color: #f1f1f1;
// 				padding: 10px;
// 				text-align: center;
// 				margin-top: 20px;
// 				font-size: 12px;
// 				color: #888;
// 				border-radius: 0 0 8px 8px;
// 			}
// 			.content {
// 				padding: 20px;
// 				background-color: white;
// 				border-radius: 8px;
// 			}
// 			.content p {
// 				font-size: 16px;
// 				line-height: 1.5;
// 			}
// 		</style>
// 	</head>
// 	<body>
// 		<div class="container">
// 			<div class="header">
// 				<h2>New Appointment Request</h2>
// 			</div>
// 			<div class="content">
// 				<p><strong>Name:</strong> %s</p>
// 				<p><strong>Email:</strong> %s</p>
// 				<p><strong>Mobile:</strong> %s</p>
// 				<p><strong>Service Requested:</strong> %s</p>
// 				<p><strong>Preferred Date:</strong> %s</p>
// 				<p><strong>Preferred Time:</strong> %s</p>
// 				<p><strong>Description:</strong><br>%s</p>
// 			</div>
// 			<div class="footer">
// 			    <p>&copy; 2024 Your Company Name. All rights reserved.</p>
// 		</div>
// 		</div>
// 	</body>
// 	</html>
// 	`, name, email, mobile, service, date, time, message)

// 	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", to, emailSubject, emailBody))

// 	// Authentication
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// Send email
// 	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
// 	if err != nil {
// 		log.Printf("Error sending email: %v", err)
// 		http.Error(w, "Failed to send email", http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond back with a success message in HTML format
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, `
// 	<html>
//     <head>
//         <title>Appointment Confirmation</title>
//         <link rel="stylesheet" href="css/popup_submeet.css">
//     </head>
//     <body>
//         <div class="confirmation custom-confirmation">
//             <h2>Appointment Successfully Booked!</h2>
//             <p>Thank you, %s, for booking an appointment. I have received the following details:</p>
//             <ul>
//                 <li><strong>Name:</strong> %s</li>
//                 <li><strong>Email:</strong> %s</li>
//                 <li><strong>Mobile:</strong> %s</li>
//                 <li><strong>Service:</strong> %s</li>
//                 <li><strong>Date:</strong> %s</li>
//                 <li><strong>Time:</strong> %s</li>
//                 <li><strong>Message:</strong> %s</li>
//             </ul>
//             <button class="back-button custom-back-button" onclick="javascript:history.back()">Go Back</button>
// 			<div class="footer">
// 			<p align="center" >&copy; 2024 samuel omulo. All rights reserved.</p>
// 		</div>
//         </div>
//         </div>
//     </body>
// </html>
// 	`, name, name, email, mobile, service, date, time, message)
// }
