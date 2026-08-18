package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("load .env: %v", err)
	}

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	allowedEmails := os.Getenv("ALLOWED_EMAILS")
	if projectID == "" || allowedEmails == "" {
		log.Fatal("FIREBASE_PROJECT_ID and ALLOWED_EMAILS must be set")
	}
	allowed := make(map[string]struct{})
	for _, email := range strings.Split(allowedEmails, ",") {
		if normalized := strings.ToLower(strings.TrimSpace(email)); normalized != "" {
			allowed[normalized] = struct{}{}
		}
	}
	if len(allowed) == 0 {
		log.Fatal("ALLOWED_EMAILS must contain at least one email address")
	}

	ctx := context.Background()
	catalogDB, cleanupCatalog, err := openCatalogForServer(ctx)
	if err != nil {
		log.Fatalf("open catalog: %v", err)
	}
	defer catalogDB.Close()
	defer cleanupCatalog()
	firebaseApp, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
	if err != nil {
		log.Fatalf("initialize Firebase: %v", err)
	}
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		log.Fatalf("initialize Firebase Auth client: %v", err)
	}

	gcsSigner, err := newGCSMediaSigner(ctx)
	if err != nil {
		log.Fatalf("initialize GCS media signer: %v", err)
	}

	a := app{verifier: firebaseVerifier{client: client}, allowed: allowed, gcsMediaSigner: gcsSigner, db: catalogDB}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	address := ":" + port
	log.Printf("server listening on http://localhost%s", address)
	if err := http.ListenAndServe(address, newRouter(a)); err != nil {
		log.Fatal(err)
	}
}
