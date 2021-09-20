// Package grun wraps the gtk application startup process with flexible settings.
//
// New and Set applies Params to configure App.App and Win.
//
// Run launches Exec to work on lists of Actions.
//
// Actions can be any function with allowed arguments and returns, or lists of
// functions ([]func, map[string]func). Recursion is possible (lists in lists).
//
//   - Arguments:  none, App
//   - Returns:    none, gtk.Widgetter, error, Errors,
//                 (gtk.Widgetter, error), (gtk.Widgetter, Errors)
//
// So you can choose and edit your function parameters as you need.
//
//
// Hello World
//
// The basic example to start a gtk application
//   func main() { gapp.Run() }
//
//   var gapp = grun.App{
//     ID:     "com.github.gtkool4.hello.World",
//     Title:  "Hello world",
//     Width:  400,
//     Height: 400,
//     OnRun:  func() gtk.Widgetter { return gtk.NewLabel("Hello gotk4!") },
//   }
//
//
// Testable Proposal for the go gtk4 startup process level 2 API
//
// This package is a running example of a proposal to become the advised way to
// run GTK4 applications. It intends to have a fast and clean declaration for
// apps and flexible behavior for tests.
//
// Warning, this is an early discussion mode topic, so the code is subject to
// brutal change at any moment. So if you find this version perfect, make backups :) .
//
// Disclaimer, please forgive me for every spelling, grammar or other mistake in
// that doc as english is not my first language.
//
//
// Goals
//
// What this package tries to do:
//
//   - Gather clean settings declaration at the top of the file for apps.
//   - Easy to start with a one liner main that requires no testing.
//   - Hide the gtk.Application callback management and its basic needs:
//     Create app, connect callbacks, create window, set size, set title,
//     pack widget, show window, save pointers...
//   - Usage
//     - Can change between headless and with window.
//     - Can change between auto-close or not.
//     - Set those globally or locally.
//     - Do those changes easily (commenting preferred).
//     - No huge params list on functions.
//     - Short and readable (minimum boilerplate, and edits for tests).
//
// In summary, it should reduce at most the app startup process, especially for
// tests files that could have a lot of window creations.
//
// Bonus: It should be pretty easy to work with lists of interfaces.
// You can call Run with almost any kind of usable func (let us know what's
// missing). And Actions allow recursion to customize any kind of crazy config.
//
// There's much more to discover in the examples.
//
//
// Todo
//
//   - App.App open signal to open files from the command line or gui.
//   - ForceWindowInSingleTest Param: Need to detect if we're running a single
//     or package test to auto-toggle the show window.
//
//
// FEEDBACK - Evolution - Options - Need tests, comments and ideas
//
// Please test it and let us know if it was usable, or if you think some things
// could be improved, especially naming.
//
// This package, or an evolution of it, could remain mainstream as standalone or
// be integrated besides the libs, to help all gtk users start their programs so
// feedback will be really appreciated.
//
// Options (ideas possible to implement):
//   - Change OnStart callback to return an error. I feel this will mostly be used
//     to LoadConfig() and/or InitDB() so an error/exit management would be nice.
//   - MultiWindow flag: Allow any widget provided to open a window (at startup or later)
//   - Rename Run to Go ?
//   - Package name ideas:
//      -grun      Go/Gtk Run          My best candidate so far. Run is the package main call.
//      -gruntk    Run Gtk or reverse  Long for repetitive test typing.
//     - napp      New App             A nice option, but I think grun is better.
//     - appinfo   Application Info    Nice but a little long for repetitive calls
//     - gtg       Good To Go          I liked this idea a lot but it would be confusing with gtk
//                                     in test files.
//
//
// Vocabulary
//
// This documentation tries to always use the same term to talk about the same
// things, for clarity.
//
// List of terms defined for this documentation:
//
//   App            This package App object.
//   App.App        The *gtk.Application object pointer.
//   Action(s)      Any kind of usable function/closure/method, or list of.
//                  Usable on Run and after with Exec
//   Exec           Parse and calls Actions on Run and after.
//   GoExitCode     ExitCode returned from the go application.
//   GtkExitCode    App.App returned value. Used as App return value if > 0.
//   Param(s)       Setting(s) to apply before Run.
//   Run            App.App startup process with Exec.
//   Win            The *gtk.Window object pointer.
//
// Usage
//
// New creates the App object or it can be created manually.
//
// Run starts the App, blocking the main go loop until App is exited when
// the last connected window is closed or an exit was requested:
//
//   gapp := grun.NewSized(400, 200, Params...)
//   gapp.Run(Actions...)
//
//
// Paramaters functions
//
// Params and Actions are parameters functions which mean they are functions
// prepared to be called later. This has to be reminded as some things aren't
// always ready when the functions are created.
//
// So Params and Actions are a list of prepared calls that will be executed in
// the provided order.
//
//
// Actions
//
// List of types usable with Run or Exec:
//
// With widget for the window.
//   func() gtk.Widgetter               // Simple with widget.
//   func() (gtk.Widgetter, error)      // The same with errors.
//   func() (gtk.Widgetter, Errors)     // buildhelp (gtk.Builder) errors list.
//   func(*App) gtk.Widgetter           // To act on App or Win object.
//   func(*App) (gtk.Widgetter, error)  // ...
//
// Headless.
//   func()                   // Simple func or closure.
//   func() error             // With error testing.
//   func(*App)               // To act on App or Win object.
//   func(*App) error         // ...
//   func() func(*App)        // In case the Action is wrapped.
//
// Lists.
//   []interface{}            // Recursive list of any handled type.
//   map[string]interface{}:  // Warning, execution order from a map is random.
//                            // This is mostly for tests and serial queuing.
//
// String as label window (for tests)
//   string                   // Display a string.
//   func() string            // Or a returned string.
//   func(*AppInfo) string    // ...
//
//
// Callbacks
//
// With the advice to use application in GTK, callbacks are now our also our
// applications main entry point.
//
// They run in this order:
//
//   - OnInit           Optional (logger, config and DB init for example)
//   - OnRun            Where all the work is done, and/or in the Run arguments.
//     - Exec           Launch Actions.
//                      If an Action can create a widget:
//                        Create and configure the window.
//                        Create the widget.
//                        Pack the widget if it's not nil and show the window.
//                      If errors are returned, Stop.
//   - ..........       Application running........
//   - OnStop           Optional.
//
//
// Notes
//
//  - Actions set in OnRun are called before those provided in the Run call to
//    allow global actions before local actions
//  - Only one window will be created with the first valid widget found (so
//    there will be something to put inside).
//  - The returned exit code is the first positive between GtkExitCode and
//    GoExitCode (use the ExitCode method for GoExitCode).
//  - The returned exit code can be used with os.Exit but that prevents any
//    defer calls from running. Use at your own risks.
//
package grun

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// Format errors messages.
var (
	FmtErrExec        = "grun.Exec(%s): %w"               // Format: name, error
	FmtErrRun         = "grun.Run: %s"                    // Format: error
	FmtErrTypeUnknown = "grun exec func type unknown: %T" // Format: interface{}
	FmtErrLabel       = "errors:\n%s"                     // Format: error1\nerror2\nerror3...
	TxtErrNoWidget    = "grun.Errors.Widget called without widget"
)

// GuessName formats.
var (
	FmtID        = "com.github.gtkool4.default.%s.%s" // Format: repo, package
	FmtTitle     = "%s/%s"                            // Format: repo, package
	FmtTitleTest = "go test %s"                       // Format: repo/package
)

// Action defines an action usable on Run or after.
type Action interface{}

// Param defines a parameter usable before Run.
type Param func(*App)

// App defines application settings to run a GTK application.
type App struct {
	// App and Window settings.
	ID        string               // Format: "org.gtk.example"
	Title     string               // Window title
	Width     int                  // Window width
	Height    int                  // Window height
	Args      []string             // GTK command line arguments: https://www.systutorials.com/docs/linux/man/7-gtk-options/
	Flags     gio.ApplicationFlags // See flags: https://pkg.go.dev/github.com/diamondburned/gotk4/pkg/gio/v2#ApplicationFlags
	Headless  bool                 // Force without window
	GuessName bool                 // Auto set ID and Title if empty
	FmtID     string
	FmtTitle  string

	// Application callbacks (connected to application signals).
	OnInit func(*gtk.Application) // Sets up the application when it first starts
	OnRun  interface{}            // This corresponds to the application being launched by the desktop environment.
	OnStop func(*gtk.Application)

	// OnOpen        func(app *gtk.Application, files unsafe.Pointer, hint string, test string) // opens files and shows them in a new window. This corresponds to someone trying to open a document (or documents) using the application from the file browser, or similar.

	// Pointers.
	App *gtk.Application       // Set before OnNewApp
	Win *gtk.ApplicationWindow // Set before OnNewWin. Only set if OnNewWin is defined.

	// Private.
	exitCode int // Go exit code.
}

//
//-------------------------------------------------------[ DEFAULTS APP INFO ]--

// New creates an App with Params.
func New(params ...Param) *App {
	return (&App{}).Set(params...)
}

// NewSized creates an App of given size.
func NewSized(w, h int, params ...Param) *App {
	return (&App{Width: w, Height: h}).Set(params...)
}

// NewTiny creates an App of size 400x200.
func NewTiny(params ...Param) *App { return NewSized(400, 200, params...) }

// NewSmall creates an App of size 600x400.
func NewSmall(params ...Param) *App { return NewSized(600, 400, params...) }

// NewMedium creates an App of size 800x600.
func NewMedium(params ...Param) *App { return NewSized(800, 600, params...) }

// NewLarge creates an App of size 1000x800.
func NewLarge(params ...Param) *App { return NewSized(1000, 800, params...) }

// Run starts the application and creates the window if needed.
// Actions are launched in the provided order.
// The window is created only for the first action that returned a valid widget.
//
// Locks the thread until application release when the last attached window is
// closed or an exit is requested.
//
// Returns an error code.
func (app *App) Run(calls ...interface{}) int {
	if app.OnRun != nil {
		calls = append([]interface{}{app.OnRun}, calls...)
	}
	var e error
	app.Init(func(_ *gtk.Application) { e = Exec(calls...)(app) })
	exitGtk := app.App.Run(app.Args)
	if e != nil {
		fmt.Printf(FmtErrRun+"\n", e)
		return 1
	}
	if app.ExitCode() != 0 {
		return app.ExitCode()
	}
	return exitGtk
}

//
//-----------------------------------------------------------[ INTERNAL WORK ]--

// Init creates the gtk.Application and connects its callbacks.
func (app *App) Init(call func(app *gtk.Application)) {
	if app.GuessName {
		repo, packag := packageName()
		if app.ID == "" {
			app.ID = fmt.Sprintf(firstNonEmpty(app.FmtID, FmtID), repo, packag) // "gtkelp.appinfo"
		}
		if app.Title == "" {
			app.FmtTitle = fmt.Sprintf(firstNonEmpty(app.FmtTitle, FmtTitle), repo, packag)
		}
	}
	app.App = gtk.NewApplication(app.ID, app.Flags)

	// Registered in their execution order to show how they are called.

	if app.OnInit != nil {
		app.App.Connect("startup", app.OnInit)
	}

	app.App.Connect("activate", call)

	if app.OnStop != nil {
		app.App.Connect("shutdown", app.OnStop)
	}
}

// NewWindow creates a new window and apply title and size settings.
func (app *App) NewWindow() *gtk.ApplicationWindow {
	win := gtk.NewApplicationWindow(app.App)
	if app.Title != "" {
		win.SetTitle(app.Title)
	}
	if app.Width > 0 && app.Height > 0 {
		win.SetDefaultSize(app.Width, app.Height)
	}
	return win
}

// Pack creates the widget and if it's usable, creates the window to pack it.
func (app *App) Pack(call func() gtk.Widgetter) {
	if app.Headless || app.Win != nil {
		call() // Drop widget. TODO: or append under the first widget or in its own window ?
		return
	}
	app.Win = app.NewWindow()
	w := call()
	if w == nil {
		// TODO: handle error: widget nil
		app.Win.Close()
		app.Win = nil
		return
	}
	app.Win.SetChild(w)
	app.Win.Show()
}

//
//----------------------------------------------------------[ LAUNCH ACTIONS ]--

// Exec creates an Action that launch any kind of Actions.
func Exec(calls ...interface{}) func(*App) error {
	return func(app *App) error {
		var w gtk.Widgetter
		var e error
		for _, uncast := range calls {
			switch call := uncast.(type) {

			//
			// Widgets.

			case func() gtk.Widgetter:
				app.Pack(func() gtk.Widgetter { return call() })

			case func(app *App) gtk.Widgetter:
				app.Pack(func() gtk.Widgetter { return call(app) })

			case func() (gtk.Widgetter, error):
				app.Pack(func() gtk.Widgetter {
					w, e = call()
					if e != nil {
						return nil
					}
					return w
				})

			case func(app *App) (gtk.Widgetter, error):
				app.Pack(func() gtk.Widgetter {
					w, e = call(app)
					if e != nil {
						return nil
					}
					return w
				})

				// useful ???
			case chan gtk.Widgetter:
				app.Pack(func() gtk.Widgetter { w := <-call; close(call); return w }) // <3

				//
				// Errors: errors lists

			case func(app *App) (gtk.Widgetter, Errors):
				var errs Errors
				app.Pack(func() gtk.Widgetter {
					w, errs := call(app)
					if errs.IsError() {
						return nil
					}
					return w
				})
				if errs.IsError() {
					return errs.ToError()
				}

				//
				// Headless.

			case func():
				call()

			case func() error:
				return call()

			case Param:
				call(app)

			case func(app *App):
				call(app)

			case func(app *App) error:
				return call(app)

			case func() func(*App):
				call()(app)

				//
				// Lists.

			case []interface{}: // Recursion to allow any kind of crazy config.
				e = Exec(call...)(app)

			case map[string]interface{}: // Recursion... Déjà vu.
				var errs Errors
				for name, recall := range call { // Warning, from a map, the order is random
					e := Exec(recall)(app) // This is mostly for tests and serial queuing.
					if e != nil {
						errs.Append(fmt.Errorf(FmtErrExec, name, e))
					}
				}
				if errs.IsError() {
					e = errs.ToError()
				}

				//
				// String as label could be used for tests.

			case string:
				app.Pack(func() gtk.Widgetter { return gtk.NewLabel(call) })

			case func() string:
				app.Pack(func() gtk.Widgetter { return gtk.NewLabel(call()) })

			case func(app *App) string:
				app.Pack(func() gtk.Widgetter { return gtk.NewLabel(call(app)) })

			default:
				e = fmt.Errorf(FmtErrTypeUnknown, call)
			}

			if e != nil {
				return e
			}
		}
		return nil
	}
}

//
//--------------------------------------------------------------------[ EXIT ]--

// Exit closes the application and terminates Run. Stores the go exit code.
func (app *App) Exit(exitCode int) { app.exitCode = exitCode; app.App.Quit() }

// ExitCode returns the go exit code provided by any of the Exit method.
func (app *App) ExitCode() int { return app.exitCode }

//
//-----------------------------------------------------------------[ ACTIONS ]--

// Exit creates an Action that closes the application.
func Exit(exitCode int) Action {
	return func(app *App) { app.Exit(exitCode) }
}

// ExitAfter creates a Param that closes the application after duration.
// Usable at any moment.
func ExitAfter(d time.Duration, exitCode int) Param {
	return func(app *App) {
		go time.AfterFunc(d, func() { app.Exit(exitCode) })
	}
}

// Println creates an Action that prints data, usable in tests.
func Println(args ...interface{}) func() { return func() { fmt.Println(args...) } }

//
//--------------------------------------------------------------[ SET PARAMS ]--

// Set applies a list of Param on App.
func Set(calls ...Param) func(app *App) {
	return func(app *App) { app.Set(calls...) }
}

// Set applies a list of Param on App.
func (app *App) Set(calls ...Param) *App {
	for _, call := range calls {
		call(app)
	}
	return app
}

//
//------------[ PARAMS / ACTIONS - Usable before Run and until Win is opened ]--

// SetHeadless creates a Param that prevents windows to be opened.
// Usable until Win is opened.
func SetHeadless() Param {
	return func(app *App) { app.Headless = true }
}

// SetSize creates a Param that sets the window title.
// Usable until Win is opened.
func SetSize(w, h int) Param {
	return func(app *App) { app.Width = w; app.Height = h }
}

// SetTitle creates a Param that sets the window title.
// Usable until Win is opened.
func SetTitle(str string) Param {
	return func(app *App) { app.Title = str }
}

// SetFmtTitle creates a Param that sets the format title text.
// Activates GuessName.
// Only usable before Run.
func SetFmtTitle(str string) Param {
	return func(app *App) { app.FmtTitle = str; app.GuessName = true }
}

// SetFmtTitleTest creates a Param that sets the format title text
// to "go test .../package/repo".
// Activates GuessName.
// Only usable before Run.
func SetFmtTitleTest() Param {
	return func(app *App) { SetFmtTitle(FmtTitleTest)(app) }
}

// SetGuessName creates a Param that guesses missing application id and title.
// Only usable before Run.
func SetGuessName() Param {
	return func(app *App) { app.GuessName = true }
}

//
//-----------------------------------------[ PARAMS - Only usable before Run ]--

// SetOnInit creates a Param that sets the OnInit Action.
// Only usable before Run.
func SetOnInit(call func(*gtk.Application)) Param {
	return func(app *App) { app.OnInit = call }
}

// SetOnRun creates a Param that sets the OnRun Action.
// Only usable before Run.
func SetOnRun(call ...interface{}) Param {
	if len(call) == 0 {
		return func(app *App) {}
	}
	return func(app *App) { app.OnRun = call }
}

// SetOnStop creates a Param that sets the OnStop Action.
// Only usable before Run.
func SetOnStop(call func(*gtk.Application)) Param {
	return func(app *App) { app.OnStop = call }
}

// SetID creates a Param that sets the application ID.
// Only usable before Run.
func SetID(str string) Param {
	return func(app *App) { app.ID = str }
}

// SetArgs creates a Param that sets the gtk command lines args.
// Only usable before Run.
func SetArgs(args ...string) Param {
	return func(app *App) { app.Args = args }
}

// SetFmtID creates a Param that sets the format ID text.
// Activates GuessName.
// Only usable before Run.
func SetFmtID(str string) Param {
	return func(app *App) { app.FmtID = str; app.GuessName = true }
}

// SetFlagNonUnique creates a Param that activates the non unique application
// flag.
// Only usable before Run.
func SetFlagNonUnique() Param {
	return func(app *App) { app.Flags |= gio.ApplicationNonUnique }
}

//
//------------------------------------------------------------------[ ERRORS ]--

// Errors define an error list. They can be generated by the buildhelp package.
type Errors []error

// IsError returns true is the error list is not empty.
func (e Errors) IsError() bool {
	return len(e) > 0
}

// Append adds an error to the list.
func (e *Errors) Append(more ...error) {
	*e = append(*e, more...)
}

// ToError converts the error list to a single go error.
func (e Errors) ToError() error { return errors.New(e.Error()) }

// Error returns the list of errors as string. Acts as an error for fmt.
func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}
	list := make([]string, len(e))
	for i, err := range e {
		list[i] = err.Error()
	}
	return strings.Join(list, "\n")
}

// Widget returns either a new error label widget or the provided widget.
// If a widget is provided as optional parameter, it will be returned when no
// error is found to ensure a valid widget is returned.
func (e Errors) Widget(b ...gtk.Widgetter) gtk.Widgetter {
	switch {
	case e.IsError():
		return gtk.NewLabel(fmt.Sprintf(FmtErrLabel, e))

	case len(b) > 0 && b[0] != nil:
		return b[0]
	}
	return gtk.NewLabel(TxtErrNoWidget)
}

//
//-------------------------------------------------------------[ FORMAT NAME ]--

// packageName tries to find the repository and package name from the caller.
func packageName() (repo, packag string) {
	_, path, _, _ := runtime.Caller(3)                // path= .../gtk4/gtkelp/gtkest/gtkest.go
	repo, packag = filepath.Split(filepath.Dir(path)) // Drop filename and get package name

	if repo == "gtkelp" && packag == "gtkest" { // Coming from test package, retry one call further.
		_, path, _, _ := runtime.Caller(4)
		repo, packag = filepath.Split(filepath.Dir(path))
	}
	return filepath.Base(repo), packag // Trim all the path from the repo name
}

// firstNonEmpty returns the first non empty string found.
func firstNonEmpty(list ...string) string {
	for _, str := range list {
		if str != "" {
			return str
		}
	}
	return ""
}
