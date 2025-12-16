package main

import (
	"fmt"
	"os"

	"email-verifier/internal/service"
)

func main() {
	// --------------------------------------------------
	// ARGUMENT PARSING
	// --------------------------------------------------
	if len(os.Args) < 2 {
		fmt.Println("Usage: email-verifier <email>")
		os.Exit(1)
	}

	email := os.Args[1]

	// --------------------------------------------------
	// CALL CORE SERVICE
	// --------------------------------------------------
	result := service.VerifyEmail(email)

	// --------------------------------------------------
	// OUTPUT RESULTS
	// --------------------------------------------------
	fmt.Println("Email Verification Result")
	fmt.Println("-------------------------")
	fmt.Printf("Email   : %s\n", result.Email)
	fmt.Printf("Syntax  : %v\n", result.Syntax)
	fmt.Printf("Domain  : %v\n", result.Domain)
	fmt.Printf("MX      : %v\n", result.MX)
	fmt.Printf("SMTP    : %s\n", result.SMTP)

	// --------------------------------------------------
	// EXIT CODE (IMPORTANT FOR SCRIPTS)
	// --------------------------------------------------
	if !result.Syntax || !result.Domain || !result.MX {
		os.Exit(2)
	}

	os.Exit(0)
}
