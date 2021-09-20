///bin/true; exec /usr/bin/env go run "$0" "$@"   ## Shebang trick to directly run as script on unix like. Use once: chmod u+x file
package grun_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"github.com/gtkool4/grun"
)

// Start application. Move the main function to a dedicated package.
// This helps the main package stay clean and not requiring impossible tests.
//
func main() { App.Run() } // Can use argument: more functions to call.

// Define application information.
//
var App = grun.App{
	ID:     "com.github.gtkool4.grun.example", // GTK Application ID. Format: "org.gtk.example"
	Title:  "Basic Application",               // Window title
	Width:  400,                               // Window size
	Height: 200,                               //
	OnRun:  onRun,                             // Create widget function to fill the window.
}

// Create the minimal graphical interface for our window.
//
func onRun() gtk.Widgetter { return gtk.NewLabel("Hello, gotk4 !") }

// Create a simple Gtk4 Application in go.
//
func Example() {
	// Launch our first simple App.
	//
	fmt.Println("ExitCode  :", App.Run(grun.Exit(42))) // With autoclose for the test.

	// Another way to create and run the same application:
	//
	grun.NewSized(400, 200,
		grun.SetTitle("Basic Application"),
		grun.SetID("com.github.gtkool4.grun.example"),
	).Run(
		onRun,
		grun.Exit(0),
	)

	// The code below is not needed for our example. It shows other available fields.

	// Launch another App, with Actions provided on Run.
	//
	MoreAppFields.Run(
		grun.Println("--[ running ]--"),  // Basic Action to print.
		onRun,                            // Create the widget for the window..
		grun.ExitAfter(time.Second/4, 1), // Delayed autoclose App for tests.
	)

	// Define custom Param.
	grun.New(MoreParams).Run(grun.ExitAfter(time.Second/4, 1))

	// An error returned closes the application.
	grun.FmtErrRun = "error: %s"
	withErr := func() error { return errors.New("stop") }
	grun.NewMedium(grun.SetGuessName()).Run(withErr)

	// Headless mode.
	grun.New(
		grun.SetFmtTitleTest(), // Triggers GuessName to set ID and Title.
		grun.SetHeadless(),     // Prevents opening windows.
	).Run(
		onRun, // The returned widget won't be used.
		grun.Exit(0),
	)

	// Output:
	// ExitCode  : 42
	// --[ started ]--
	// --[ running ]--
	// --[ stopped ]--
	// App.ID    : com.github.gtkool4.grun.example3
	// Win.Title : Window Title
	// Win.Size  : 400 x 200
	// error: stop
}

//
//-------------------------------------------------[ MORE SETTINGS AS FIELDS ]--

// More App fields less used by basic applications.
var MoreAppFields = grun.App{
	ID:        "com.github.gtkool4.grun.example2", // GTK Application ID. Format: "org.gtk.example"
	Headless:  true,                               // Force without window. Usable until the first window is opened.
	GuessName: false,                              // Auto set ID and Title if empty.
	Flags:     gio.ApplicationNonUnique,           // See flags: https://pkg.go.dev/github.com/diamondburned/gotk4/pkg/gio/v2#ApplicationFlags
	Args:      []string{"--name=NAME"},            // GTK command line arguments: https://www.systutorials.com/docs/linux/man/7-gtk-options/
	OnInit:    OnInit,                             // Sets up the application when it first starts
	OnStop:    OnStop,                             // Called when the application is closing.
	FmtID:     "",                                 // With GuessName ID format: "%s.%s" as package repo
	FmtTitle:  "",                                 // With GuessName title format: "%s" as reformated ID (replace . with /)
}

func OnInit(*gtk.Application) { fmt.Println("--[ started ]--") }
func OnStop(*gtk.Application) { fmt.Println("--[ stopped ]--") }

//
//-------------------------------------------------[ MORE SETTINGS AS PARAMS ]--

func MoreParams(gapp *grun.App) {
	gapp.Set( // Other available options.
		grun.SetID("com.github.gtkool4.grun.example3"),
		grun.SetTitle("Window Title"),
		grun.SetSize(400, 200),
		grun.SetOnInit(func(*gtk.Application) {}),
		grun.SetOnRun(), // test empty.
		grun.SetOnRun(testUI),
		grun.SetOnStop(func(*gtk.Application) {}),
	)
	// or
	grun.Set(
		grun.SetFmtID(grun.FmtID),
		grun.SetFmtTitle("go test %s"),
		grun.SetGuessName(),
		grun.SetArgs("--class=APPCLASS"),
		grun.SetFlagNonUnique(),
	)(gapp)
}

//
//-------------------------------------------------[ TEST APP & WIN SETTINGS ]--

// Create and return a widget to fill the window. Prints applied settings for the test.
//
func testUI(app *grun.App) gtk.Widgetter {
	fmt.Printf("App.ID    : %s\nWin.Title : %s\nWin.Size  : %d x %d\n",
		app.App.ApplicationID(),
		app.Win.Title(),
		app.Win.ObjectProperty("default-width"),
		app.Win.ObjectProperty("default-height"),
	)
	return gtk.NewLabel("Hello, gotk4 !")
}
