package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/amenzhinsky/go-polkit"
)

var (
	localeFlag        string
	checkAccessFlag   bool
	allowPasswordFlag bool
	verboseFlag       bool
)

var stdout = bufio.NewWriterSize(os.Stdout, 100)

func main() {
	flag.StringVar(&localeFlag, "locale", "", "output locale")
	flag.BoolVar(&checkAccessFlag, "check-access", false, "check access to an action")
	flag.BoolVar(&allowPasswordFlag, "allow-password", false, "ask user for password when needed")
	flag.BoolVar(&verboseFlag, "verbose", false, "enable verbose mode")
	flag.Parse()

	defer stdout.Flush()

	authority, err := polkit.NewAuthority()
	if err != nil {
		panic(err)
	}

	switch {
	case checkAccessFlag:
		if err = checkAccess(authority, flag.Args(), allowPasswordFlag); err != nil {
			panic(err)
		}
	default:
		if err = enumerateActions(authority, flag.Args(), localeFlag, verboseFlag); err != nil {
			panic(err)
		}
	}
}

func out(s string, v ...interface{}) {
	stdout.WriteString(fmt.Sprintf(s, v...))
}

func checkAccess(authority *polkit.Authority, actions []string, allowPassword bool) error {
	flags := polkit.CheckAuthorizationNone
	if allowPassword {
		flags = polkit.CheckAuthorizationAllowUserInteraction
	}

	var k, v string
	for _, action := range actions {
		result, err := authority.CheckAuthorization(action, nil, flags, "")
		if err != nil {
			return err
		}

		out("%s:\n", action)
		out("  Is authorized: %t\n", result.IsAuthorized)
		out("  Is challenge:  %t\n", result.IsChallenge)

		if len(result.Details) > 0 {
			out("  Details:\n")
			for k, v = range result.Details {
				out("    %s -> %s\n", k, v)
			}
		}

		out("\n")
	}

	return nil
}

func enumerateActions(authority *polkit.Authority, actionIDs []string, locale string, verbose bool) error {
	actions, err := authority.EnumerateActions(locale)
	if err != nil {
		return err
	}

	var k, v string
	for _, action := range actions {
		if skipAction(actionIDs, action.ActionID) {
			continue
		}

		out(action.ActionID)

		if verbose {
			out(":\n")
			out("  Descriprion:       %s\n", action.Description)
			out("  Message:           %s\n", action.Message)
			out("  Vendor name:       %s\n", action.VendorName)
			out("  Vendor URL:        %s\n", action.VendorURL)

			if action.IconName != "" {
				out("  Icon name:         %s\n", action.IconName)
			}

			out("  Implicit any:      %s\n", polkit.PKImplicitAuthorization(action.ImplicitAny))
			out("  Implicit inactive: %s\n", polkit.PKImplicitAuthorization(action.ImplicitInactive))
			out("  Implicit active:   %s\n", polkit.PKImplicitAuthorization(action.ImplicitActive))

			if len(action.Annotations) > 0 {
				out("  Annotations:\n")
				for k, v = range action.Annotations {
					out("    %s -> %s\n", k, v)
				}
			}
		}

		out("\n")
	}

	return nil
}

func skipAction(actions []string, action string) bool {
	if len(actions) == 0 {
		return false
	}

	for i := 0; i < len(actions); i++ {
		if actions[i] == action {
			return false
		}
	}

	return true
}
