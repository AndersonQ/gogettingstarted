Task: read environment variables

To read environment variables we can use `os.LookupEnv` or `os.GetEnv`, check their docs to see the difference 
in-between them. It's the manual way, it'' good to know, however they aren't handy. There are libs which allow us to 
define a struct and annotate each struct field with the environment variabe name, a default value, mark it as required
and so on. I usually use [github.com/caarlos0/env](https://pkg.go.dev/github.com/caarlos0/env?tab=overview), check it out.

Open [config_test.go](config_test.go) and implement the tests

Task: add a logger

Go has the `log` package which provides basic log functionalities, whereas useful it's quite simple and lacks some
functionalities we need. Therefore, we'll use a log library. I've been using `zerolog` it's a 
_fast and simple logger dedicated to JSON output_, also with a pretty loggin on the console, not quite efficient as
the JSON form, but excellent for development. Check the [documentation](https://pkg.go.dev/github.com/rs/zerolog?tab=overview),
and on [config.go](solution/config.go) you will find the `func (Config) Logger` function showing how to 
initialise it.
