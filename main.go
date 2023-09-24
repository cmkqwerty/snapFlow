package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch email me at <a href=\"mailto:cemeke10@gmail.com\">cemeke10@gmail.com</a></p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<h1>FAQ Page</h1>
<ul>
	<li><b>Is there a free version?</b> Yes! We offer a free trial for 30 days on any paid plans.</li>
	<li><b>What happens after my free trial?</b> If you chose to keep your paid plan after the 30 day trial, you will need to enter your payment details and confirm.</li>
	<li><b>Can I upgrade or downgrade anytime?</b> Yes! SnapFlow is a pay-as-you-go service and you can upgrade, downgrade or cancel at any time.</li>
	<li><b>What are your support hours?</b> We have support staff answering emails 24/7, though response times may be a bit slower on weekends.</li>
	<li><b>How do I contact support?</b> Email us - <a href="mailto:support@snapflow.com">support@snapflow.com</a></li>
</ul>
`)
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		http.Error(w, "Page not found", http.StatusNotFound)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
