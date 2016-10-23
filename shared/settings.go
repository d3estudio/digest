package shared

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	log "github.com/Sirupsen/logrus"
)

// Settings exposes all required settings loaded from environment variables
var Settings *settings

func init() {
	Settings = loadSettings()
	log.SetLevel(log.DebugLevel)
}

type settings struct {
	Token          string
	SilencerEmojis []string `default:"no_entry_sign"`
	RedisServer    string   `default:"redis"`
	RedisPort      int      `default:"6379"`
	MongoServer    string   `default:"mongo"`
	MongoDatabase  string   `default:"digest"`
	BaseQueueName  string   `default:"Digest"`
	TwitterKey     string   `default:""`
	TwitterSecret  string   `default:""`
}

func (c *settings) Validate() {
	if len(c.Token) == 0 {
		panic("You must define a slack access Token before using Digest. Please reffer to the documentation for further instructions.")
	}
	if len(c.SilencerEmojis) == 0 {
		log.Warning("Looks like you haven't defined at least one Silencer Emoji. Although it is not necessary for normal operation, it's a nice thing to have, you know?")
	}
}

func loadSettings() *settings {
	sets := settings{}
	st := reflect.TypeOf(sets)
	ps := reflect.ValueOf(&sets)
	si := ps.Elem()
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		var def *string
		if alias, ok := field.Tag.Lookup("default"); ok {
			if alias == "" {
				def = nil
			} else {
				def = &alias
			}
		}
		envName := fmt.Sprintf("DIGEST_%s", strings.ToUpper(snakeCase(field.Name)))
		envValue := os.Getenv(envName)
		if envValue == "" && def != nil {
			envValue = *def
		}

		kind := field.Type.Kind()
		if kind == reflect.Bool {
			si.Field(i).SetBool(boolFromString(envValue))
		} else if kind == reflect.Int {
			intValue, err := strconv.Atoi(envValue)
			if err == nil {
				si.Field(i).SetInt(int64(intValue))
			}
		} else if kind == reflect.Slice && reflect.TypeOf(field.Type.Elem()).Kind() == reflect.Ptr {
			// Assuming as String, this will have to be changed anytime soon if we support other slice types
			values := strings.Split(envValue, ",")
			si.Field(i).Set(reflect.ValueOf(values))
		} else {
			si.Field(i).SetString(envValue)
		}
	}
	return &sets
}

func boolFromString(value string) bool {
	value = strings.ToLower(value)
	return value == "yes" || value == "true" || value == "y"
}

func snakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) && runes[i-1] != '_' {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
