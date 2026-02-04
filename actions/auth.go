package actions

import (
	"budget_tracker/models"
	"database/sql"
	"net/http"
	"strings"

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
	c.Flash().Add("success", "You have been logged out!")
	return c.Redirect(http.StatusSeeOther, "/")
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
