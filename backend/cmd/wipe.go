package cmd

import (
	"fmt"
	"strings"

	"github.com/lalabuy948/linkopus/backend/di"
)

type WipeCommand struct {
	container *di.Container
}

func NewWipeCommand(c *di.Container) *WipeCommand {
	return &WipeCommand{c}
}

func (wc *WipeCommand) DeleteLinkMapAndView(date string, link string) {

	if link == "" && date == "" {
		fmt.Println("Please use -date or -link flag with appropriate string argument.")
	}

	if date != "" {
		dateSplit := strings.Split(date, "-")
		year := dateSplit[0]
		month := dateSplit[1]
		day := dateSplit[2]

		err := wc.container.CommandService.DeleteLinkMapByTime(year, month, day)
		if err != nil {
			fmt.Println(err)
		}

		err = wc.container.CommandService.DeleteLinkViewByTime(year, month, day)
		if err != nil {
			fmt.Println(err)
		}
	}

	if link != "" {
		err := wc.container.CommandService.DeleteLinkMapByLink(link)
		if err != nil {
			fmt.Println(err)
		}

		err = wc.container.CommandService.DeleteLinkViewByLink(link)
		if err != nil {
			fmt.Println(err)
		}
	}
}
