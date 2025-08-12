package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	dbPath = "data/potion.db"
	tpl    = template.Must(template.ParseFiles("templates/index.html"))
	mu     sync.Mutex
)

func main() {
	db, err := sql.Open("sqlite", dsn(dbPath))
	if err != nil {
		log.Fatal("migrate/seed errorrr:", err)
	}

	defer db.Close()

	// Migrate
	if err := migrateAndSeed(db); err != nil {
		log.Fatal("Migrate seed error: ", err)
	}

	// here is the main endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		s, err := loadStore(db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		v := View{
			Store:        s,
			TotalPotions: s.Inventory["red"] + s.Inventory["blue"] + s.Inventory["green"],
			CurrentPrice: s.Prices[s.SelectedColor],
		}
		if err := tpl.ExecuteTemplate(w, "index.html", v); err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	// endpoint for the purchaces
	http.HandleFunc("/purchase", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", 405)
			return
		}
		color := r.FormValue("color")
		if color != "red" && color != "blue" && color != "green" {
			color = "red"
		}
		mu.Lock()
		defer mu.Unlock()
		if err := setSelectedColor(db, color); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if err := purchaseTxn(db); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func dsn(path string) string {
	return filepath.ToSlash(path) + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
}
