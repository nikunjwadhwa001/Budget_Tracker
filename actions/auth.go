package actions

import (
	"budget_tracker/models"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"golang.org/x/crypto/bcrypt"
)

// AuthNew renders the sign in page
func AuthNew(c buffalo.Context) error {
	c.Set("user", &models.User{})
	return c.Render(http.StatusOK, r.HTML("auth/new.plush.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	existingUser := &models.User{}
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(existingUser)

	// helper function to handle bad attempts
	bad := func() error {
		c.Flash().Add("danger", "Invalid email/password!")
		return c.Render(http.StatusUnauthorized, r.HTML("auth/new.plush.html"))
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return bad()
		}
		return err
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}

	// Check if user is verified
	if !existingUser.IsVerified {
		c.Flash().Add("warning", "Account not verified. Please verify your email.")
		c.Session().Set("pre_verification_user_id", existingUser.ID)
		return c.Redirect(http.StatusSeeOther, "/verify-otp")
	}

	c.Session().Set("current_user_id", existingUser.ID)
	c.Flash().Add("success", "Welcome back!")

	redirectPath := "/"
	if redirect, ok := c.Session().Get("redirectURL").(string); ok {
		redirectPath = redirect
		c.Session().Delete("redirectURL") // consume it
	}

	return c.Redirect(http.StatusSeeOther, redirectPath)
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	return c.Redirect(http.StatusSeeOther, "/")
}

// AuthVerifyOTP renders the OTP entry form
func AuthVerifyOTP(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/verify_otp.plush.html"))
}

// AuthVerifyOTPPost handles the OTP submission
func AuthVerifyOTPPost(c buffalo.Context) error {
	otp := c.Request().FormValue("OTP")
	userID := c.Session().Get("pre_verification_user_id")

	if userID == nil {
		c.Flash().Add("danger", "Session expired. Please sign up again.")
		return c.Redirect(http.StatusSeeOther, "/users/new")
	}

	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	if err := tx.Find(user, userID); err != nil {
		c.Flash().Add("danger", "User not found.")
		return c.Redirect(http.StatusSeeOther, "/users/new")
	}

	if user.OTPCode == otp {
		// Valid OTP (Add expiration check here if needed)
		user.IsVerified = true
		user.OTPCode = "" // Clear OTP
		if err := tx.Update(user); err != nil {
			return err
		}

		// Log them in
		c.Session().Set("current_user_id", user.ID)
		c.Session().Delete("pre_verification_user_id")
		c.Flash().Add("success", "Account verified! Welcome.")
		return c.Redirect(http.StatusSeeOther, "/")
	}

	c.Flash().Add("danger", "Invalid OTP. Please try again.")
	return c.Render(http.StatusUnprocessableEntity, r.HTML("auth/verify_otp.plush.html"))
}

// AuthForgotGet renders the forgot password page
func AuthForgotGet(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/forgot.plush.html"))
}

// AuthForgotPost handles the email submission to send OTP
func AuthForgotPost(c buffalo.Context) error {
	email := strings.TrimSpace(c.Request().FormValue("Email"))

	// WORKAROUND: Create fresh connection to ensure reliability
	db, err := pop.Connect("development")
	if err != nil {
		c.Logger().Error("Failed to connect to DB", err)
		return err
	}

	user := &models.User{}
	// Use ILIKE for case-insensitive match just in case
	if err := db.Where("email ILIKE ?", email).First(user); err != nil {
		c.Logger().Infof("DEBUG: User not found for email: '%s', Error: %v\n", email, err)
		c.Flash().Add("success", "If that email exists, we sent an OTP.")
		return c.Redirect(http.StatusSeeOther, "/reset-password")
	}

	c.Logger().Infof("DEBUG: User found: %s (ID: %s). Generating OTP...", user.Email, user.ID)

	// Generate OTP
	rand.Seed(time.Now().UnixNano())
	otp := strconv.Itoa(rand.Intn(900000) + 100000)
	expiry := time.Now().Add(10 * time.Minute)

	c.Logger().Infof("DEBUG: Updating OTP for Email: %s -> Code: %s", user.Email, otp)

	// Update by EMAIL to be sure
	query := "UPDATE users SET otp_code = '" + otp + "', otp_expires_at = '" + expiry.Format(time.RFC3339) + "' WHERE email = '" + user.Email + "'"
	if err := db.RawQuery(query).Exec(); err != nil {
		c.Logger().Error("Failed to update OTP", err)
		return err
	}

	// SIMULATE EMAIL SENDING
	msg := fmt.Sprintf("SENDING PASSWORD RESET OTP TO %s: %s", user.Email, otp)
	fmt.Println("\n\n========================================")
	fmt.Println(msg)
	fmt.Println("========================================\n")
	c.Logger().Info(msg)

	c.Session().Set("reset_user_id", user.ID)
	c.Flash().Add("success", "OTP sent to your email.")
	return c.Redirect(http.StatusSeeOther, "/reset-password")
}

// AuthResetGet renders the reset password page
func AuthResetGet(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/reset.plush.html"))
}

// AuthResetPost handles the password reset processing
func AuthResetPost(c buffalo.Context) error {
	otp := c.Request().FormValue("OTP")
	password := c.Request().FormValue("Password")
	conf := c.Request().FormValue("PasswordConfirmation")

	if password != conf {
		c.Flash().Add("danger", "Passwords do not match.")
		return c.Render(http.StatusUnprocessableEntity, r.HTML("auth/reset.plush.html"))
	}

	userID := c.Session().Get("reset_user_id")
	if userID == nil {
		c.Flash().Add("danger", "Session expired. Please try again.")
		return c.Redirect(http.StatusSeeOther, "/forgot-password")
	}

	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	if err := tx.Find(user, userID); err != nil {
		c.Flash().Add("danger", "User not found.")
		return c.Redirect(http.StatusSeeOther, "/forgot-password")
	}

	if user.OTPCode != otp {
		c.Flash().Add("danger", "Invalid OTP.")
		return c.Render(http.StatusUnprocessableEntity, r.HTML("auth/reset.plush.html"))
	}

	// Update Password (OTP will be cleared and hashing handled by model/Create but we are doing explicit Update)
	// IMPORTANT: Identify how password hashing is handled in User model.
	// The `Create` method handles hashing. But `ValidateAndUpdate` usually doesn't hashing automatically if we just set Password.
	// We need to hash it manually here or use a helper.

	// Let's re-hash manually to be safe since we are in `actions`.
	ph, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(ph)
	user.OTPCode = "" // Clear OTP

	if err := tx.Update(user); err != nil {
		return err
	}

	c.Session().Delete("reset_user_id")
	c.Flash().Add("success", "Password updated! Please log in.")
	return c.Redirect(http.StatusSeeOther, "/signin")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			if err := tx.Find(u, uid); err != nil {
				// invalid session
				c.Session().Clear()
				return next(c)
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(http.StatusSeeOther, "/signin")
		}
		return next(c)
	}
}
