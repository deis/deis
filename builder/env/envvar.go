package env

import (
	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/log"
	"os"
)

// Get gets one or more environment variables and puts them into the context.
//
// Parameters passed in are of the form varname => defaultValue.
//
// 	r.Route("foo", "example").Does(envvar.Get).Using("HOME").WithDefault(".")
//
// As with all environment variables, the default value must be a string.
//
// WARNING: Since parameters are a map, order of processing is not
// guaranteed. If order is important, you'll need to call this command
// multiple times.
//
// For each parameter (`Using` clause), this command will look into the
// environment for a matching variable. If it finds one, it will add that
// variable to the context. If it does not find one, it will expand the
// default value (so you can set a default to something like "$HOST:$PORT")
// and also put the (unexpanded) default value back into the context in case
// any subsequent call to `os.Getenv` occurs.
func Get(c cookoo.Context, params *cookoo.Params) (interface{}, cookoo.Interrupt) {
	for name, def := range params.AsMap() {
		var val string
		if val = os.Getenv(name); len(val) == 0 {
			def := def.(string)
			val = os.ExpandEnv(def)
			// We want to make sure that any subsequent calls to Getenv
			// return the same default.
			os.Setenv(name, val)

		}
		c.Put(name, val)
		log.Debugf(c, "Name: %s, Val: %s", name, val)
	}
	return true, nil
}
