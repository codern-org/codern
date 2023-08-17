package constant

import "os"

var IsDevelopment = os.Getenv("ENVIRONMENT") == "development"
