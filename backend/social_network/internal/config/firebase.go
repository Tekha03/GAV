package config

type FirebaseConfig struct {
	Enabled         bool
	CredentialsFile string
}

func loadFirebase() FirebaseConfig {
	return FirebaseConfig{
		Enabled:         getEnv("FIREBASE_ENABLED", "false") == "true",
		CredentialsFile: getEnv("FIREBASE_CREDENTIALS_FILE", ""),
	}
}
