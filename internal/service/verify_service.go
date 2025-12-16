package service

import (
	"context"
	"net"
	"net/mail"
	"strings"
	"time"

	"email-verifier/internal/model"
)

const (
	dnsTimeout  = 5 * time.Second
	smtpTimeout = 5 * time.Second
)

// Minimal disposable domain list (extend later if needed)
var disposableDomains = map[string]bool{
	"mailinator.com":     true,
	"tempmail.com":       true,
	"10minutemail.com":   true,
	"guerrillamail.com":  true,
	"throwawaymail.com": true,
}

// VerifyEmail is the single source of truth for email verification.
// Used by API + CLI.
func VerifyEmail(email string) model.VerifyResponse {
	result := model.VerifyResponse{
		Email:  email,
		SMTP:  "unchecked",
		Status: "Invalid",
	}

	// --------------------------------------------------
	// 1. SYNTAX CHECK
	// --------------------------------------------------
	_, err := mail.ParseAddress(email)
	if err != nil {
		return result
	}
	result.Syntax = true

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return result
	}
	domain := strings.ToLower(parts[1])

	// --------------------------------------------------
	// 2. DISPOSABLE EMAIL CHECK
	// --------------------------------------------------
	if disposableDomains[domain] {
		result.Domain = true
		result.MX = true
		result.SMTP = "disposable email"
		result.Status = "Risky"
		return result
	}

	// --------------------------------------------------
	// 3. DOMAIN + MX CHECK (WITH TIMEOUT)
	// --------------------------------------------------
	ctx, cancel := context.WithTimeout(context.Background(), dnsTimeout)
	defer cancel()

	resolver := net.Resolver{}
	mxRecords, err := resolver.LookupMX(ctx, domain)
	if err != nil || len(mxRecords) == 0 {
		return result
	}

	result.Domain = true
	result.MX = true

	// --------------------------------------------------
	// 4. SMTP REACHABILITY CHECK (SAFE)
	// --------------------------------------------------
	mxHost := strings.TrimSuffix(mxRecords[0].Host, ".")
	address := mxHost + ":25"

	conn, err := net.DialTimeout("tcp", address, smtpTimeout)
	if err != nil {
		result.SMTP = "unreachable"
		return result
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(smtpTimeout))
	result.SMTP = "reachable"

	// --------------------------------------------------
	// 5. CATCH-ALL HEURISTIC (NON-INTRUSIVE)
	// --------------------------------------------------
	if isLikelyCatchAll(domain, mxHost) {
		result.SMTP = "reachable (catch-all domain)"
		result.Status = "Risky"
		return result
	}

	// --------------------------------------------------
	// 6. FINAL STATUS
	// --------------------------------------------------
	if result.Syntax && result.Domain && result.MX {
		result.Status = "Deliverable"
	}

	return result
}

// isLikelyCatchAll performs a SAFE heuristic check.
// It does NOT send MAIL FROM or RCPT TO.
func isLikelyCatchAll(domain, mxHost string) bool {
	address := mxHost + ":25"

	conn, err := net.DialTimeout("tcp", address, smtpTimeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	// If server always accepts connections without discrimination,
	// the domain is *likely* catch-all.
	return true
}
