package env

import (
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/cookoo"
)

func TestGet(t *testing.T) {
	reg, router, cxt := cookoo.Cookoo()

	drink := "DEIS_DRINK_OF_CHOICE"
	cookies := "DEIS_FAVORITE_COOKIES"
	snack := "DEIS_SNACK_TIME"
	snackVal := fmt.Sprintf("$%s and $%s cookies", drink, cookies)

	// Set drink, but not cookies.
	os.Setenv(drink, "coffee")

	reg.Route("test", "Test route").
		Does(Get, "res").
		Using(drink).WithDefault("tea").
		Using(cookies).WithDefault("chocolate chip").
		Does(Get, "res2").
		Using(snack).WithDefault(snackVal)

	err := router.HandleRequest("test", cxt, true)
	if err != nil {
		t.Error(err)
	}

	// Drink should still be coffee.
	if coffee := cxt.Get(drink, "").(string); coffee != "coffee" {
		t.Errorf("A great sin has been committed. Expected coffee, but got '%s'", coffee)
	}
	// Env var should be untouched
	if coffee := os.Getenv(drink); coffee != "coffee" {
		t.Errorf("Environment was changed from 'coffee' to '%s'", coffee)
	}

	// Cookies should have been set to the default
	if cookies := cxt.Get(cookies, "").(string); cookies != "chocolate chip" {
		t.Errorf("Expected chocolate chip cookies, but instead, got '%s' :-(", cookies)
	}

	// In the environment, cookies should have been set.
	if cookies := os.Getenv(cookies); cookies != "chocolate chip" {
		t.Errorf("Expected environment to have chocolate chip cookies, but instead, got '%s'", cookies)
	}

	if both := cxt.Get(snack, "").(string); both != "coffee and chocolate chip cookies" {
		t.Errorf("Expected 'coffee and chocolate chip cookies'. Got '%s'", both)
	}

	if both := os.Getenv(snack); both != snackVal {
		t.Errorf("Expected %s to not be expanded. Got '%s'", snack, both)
	}
}
