package core

import (
	"fmt"

	"github.com/binarycraft007/wechat/core/utils"

	"github.com/kataras/go-events"

	"net/url"
)

type Core struct {
	Config utils.Config
	Events events.EventEmmiter
}

func New() (*Core, error) {
	core := Core{}

	core.Events = events.New()

	config, err := utils.NewConfig(utils.ConfigOption{})
	if err != nil {
		return nil, err
	}

	core.Config = *config

	core.Events.On("my_event", func(payload ...interface{}) {
		message := payload[0].(string)
		fmt.Println(message) // prints "this is my payload"
	})

	return &core, nil
}

func PreLogin() {
	params := url.Values{}
	params.Add("param1", "value1")
	params.Add("param2", "value2")
}
