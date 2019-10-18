//Source: https://github.com/esimov/diagram/blob/master/version/version.go

package version

import "github.com/morgulbrut/listFiles/color"

// Version number.
const Version = "v1.0.1"

// DrawLogo draws diagram logo.
func DrawLogo() string {
	var logo string

	logo += "\n\n"
	logo += color.StringRandom("  ██╗    ██╗ █████╗██████╗██████╗██╗██╗    ██████╗ █████╗  \n")
	logo += color.StringRandom("  ██║    ██║██╔═══╝  ██╔═╝██╔═══╝██║██║    ██╔═══╝██╔═══╝  \n")
	logo += color.StringRandom("  ██║    ██║ █████╗  ██║  █████╗ ██║██║    █████╗  █████╗  \n")
	logo += color.StringRandom("  ██║    ██║ ╚═══██╗ ██║  ██╔══╝ ██║██║    ██╔══╝  ╚═══██╗ \n")
	logo += color.StringRandom("  ██████╗██║ █████╔╝ ██║  ██║    ██║██████╗██████╗ █████╔╝ \n")
	logo += color.StringRandom("  ╚═════╝╚═╝ ╚════╝  ╚═╝  ╚═╝    ╚═╝╚═════╝╚═════╝ ╚════╝    " + Version + "\n\n")

	return logo
}
