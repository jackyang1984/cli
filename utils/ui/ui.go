// Package ui will provide hooks into STDOUT, STDERR and STDIN. It will also
// handle translation as necessary.
//
// This package is explicitly designed for the CF CLI and is *not* to be used
// by any package outside of the commands package.
package ui

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"code.cloudfoundry.org/cli/utils/configv3"
	"github.com/fatih/color"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/vito/go-interact/interact"
)

const (
	red   color.Attribute = color.FgRed
	green                 = color.FgGreen
	// yellow                         = color.FgYellow
	// magenta                        = color.FgMagenta
	cyan = color.FgCyan
	// grey                           = color.FgWhite
	defaultFgColor = 38
)

//go:generate counterfeiter . Config

// Config is the UI configuration
type Config interface {
	// ColorEnabled enables or disabled color
	ColorEnabled() configv3.ColorSetting
	// Locale is the language to translate the output to
	Locale() string
}

//go:generate counterfeiter . TranslatableError

// TranslatableError it wraps the error interface adding a way to set the
// translation function on the error
type TranslatableError interface {
	// Returns back the untranslated error string
	Error() string
	Translate(func(string, ...interface{}) string) string
}

// UI is interface to interact with the user
type UI struct {
	// In is the input buffer
	In io.Reader
	// Out is the output buffer
	Out io.Writer
	// Err is the error buffer
	Err io.Writer

	colorEnabled configv3.ColorSetting
	translate    i18n.TranslateFunc
}

// NewUI will return a UI object where Out is set to STDOUT, In is set to STDIN,
// and Err is set to STDERR
func NewUI(c Config) (*UI, error) {
	translateFunc, err := GetTranslationFunc(c)
	if err != nil {
		return nil, err
	}

	return &UI{
		In:           os.Stdin,
		Out:          color.Output,
		Err:          os.Stderr,
		colorEnabled: c.ColorEnabled(),
		translate:    translateFunc,
	}, nil
}

// NewTestUI will return a UI object where Out, In, and Err are customizable, and
// colors are disabled
func NewTestUI(in io.Reader, out io.Writer, err io.Writer) *UI {
	return &UI{
		In:           in,
		Out:          out,
		Err:          err,
		colorEnabled: configv3.ColorDisabled,
		translate:    translationWrapper(i18n.IdentityTfunc()),
	}
}

// TranslateText returns the translated string with templateValues substituted into the
// template string.
func (ui *UI) TranslateText(stringTemplate string, templateValues ...map[string]interface{}) string {
	return ui.translate(stringTemplate, getFirstSet(templateValues))
}

// DisplayOK outputs a green translated "OK" message to UI.Out.
func (ui *UI) DisplayOK() {
	fmt.Fprintf(ui.Out, "%s\n", ui.colorize(ui.TranslateText("OK"), green, true))
}

// DisplayNewline outputs a newline to UI.Out.
func (ui *UI) DisplayNewline() {
	fmt.Fprintf(ui.Out, "\n")
}

// DisplayBoolPrompt outputs the prompt and waits for user input. It only
// allows for a boolean response. A default boolean response can be set with
// defaultResponse.
func (ui *UI) DisplayBoolPrompt(prompt string, defaultResponse bool) (bool, error) {
	response := defaultResponse
	interactivePrompt := interact.NewInteraction(fmt.Sprintf("%s%s", prompt, ui.colorize(">>", cyan, true)))
	interactivePrompt.Input = ui.In
	interactivePrompt.Output = ui.Out
	err := interactivePrompt.Resolve(&response)
	return response, err
}

// DisplayTable presents a two dimensional array of strings as a table to UI.Out
func (ui *UI) DisplayTable(prefix string, table [][]string, padding int) error {
	tw := tabwriter.NewWriter(ui.Out, 0, 1, padding, ' ', 0)

	for _, row := range table {
		fmt.Fprint(tw, prefix)
		fmt.Fprintln(tw, strings.Join(row, "\t"))
	}

	return tw.Flush()
}

// DisplayText combines the stringTemplate template with the key maps and then
// outputs it to the UI.Out file. Prior to outputting the stringTemplate, it
// is run through an internationalization function to translate it to a
// pre-configured language. Only the first map in templateValues is used.
func (ui *UI) DisplayText(stringTemplate string, templateValues ...map[string]interface{}) {
	fmt.Fprintf(ui.Out, "%s\n", ui.TranslateText(stringTemplate, templateValues...))
}

// DisplayTextWithKeyTranslations translates the templateValues listed in
// keysToTranslate, and then passes these values to DisplayText. Only the first
// map in templateValues is used.
func (ui *UI) DisplayTextWithKeyTranslations(stringTemplate string, keysToTranslate []string, templateValues ...map[string]interface{}) {
	firstTemplateValues := getFirstSet(templateValues)
	for _, key := range keysToTranslate {
		firstTemplateValues[key] = ui.TranslateText(fmt.Sprint(firstTemplateValues[key]))
	}
	fmt.Fprintf(ui.Out, "%s\n", ui.TranslateText(stringTemplate, firstTemplateValues))
}

// DisplayPair outputs the "attribute: stringTemplate" pair to UI.Out. templateValues
// are applied to the translation of stringTemplate, while attribute is
// translated directly.
func (ui *UI) DisplayPair(attribute string, stringTemplate string, templateValues ...map[string]interface{}) {
	fmt.Fprintf(ui.Out, "%s: %s\n", ui.TranslateText(attribute), ui.TranslateText(stringTemplate, templateValues...))
}

// DisplayHelpHeader translates and then bolds the help header. Sends output to
// UI.Out.
func (ui *UI) DisplayHelpHeader(text string) {
	fmt.Fprintf(ui.Out, "%s\n", ui.colorize(ui.TranslateText(text), defaultFgColor, true))
}

// DisplayHeaderFlavorText outputs the translated text, with cyan color templateValues,
// to UI.Out.
func (ui *UI) DisplayHeaderFlavorText(stringTemplate string, templateValues ...map[string]interface{}) {
	firstTemplateValues := getFirstSet(templateValues)
	for key, value := range firstTemplateValues {
		firstTemplateValues[key] = ui.colorize(fmt.Sprint(value), cyan, true)
	}
	fmt.Fprintf(ui.Out, "%s\n", ui.TranslateText(stringTemplate, firstTemplateValues))
}

// DisplayWarning applies translation to stringTemplate and displays the
// translated warning to UI.Err.
func (ui *UI) DisplayWarning(stringTemplate string, templateValues ...map[string]interface{}) {
	fmt.Fprintf(ui.Err, "%s\n", ui.TranslateText(stringTemplate, templateValues...))
}

// DisplayWarnings translates and displays the warnings.
func (ui *UI) DisplayWarnings(warnings []string) {
	for _, warning := range warnings {
		fmt.Fprintf(ui.Err, "%s\n", ui.TranslateText(warning))
	}
}

// DisplayError outputs the error to UI.Err and outputs a red translated
// "FAILED" to UI.Out.
func (ui *UI) DisplayError(err error) {
	errMsg := ""
	if translatableError, ok := err.(TranslatableError); ok {
		errMsg = translatableError.Translate(ui.translate)
	} else {
		errMsg = err.Error()
	}
	fmt.Fprintf(ui.Err, "%s\n", errMsg)
	fmt.Fprintf(ui.Out, "%s\n", ui.colorize(ui.TranslateText("FAILED"), red, true))
}

func (ui *UI) colorize(message string, textColor color.Attribute, bold bool) string {
	colorPrinter := color.New(textColor)
	switch ui.colorEnabled {
	case configv3.ColorEnabled:
		colorPrinter.EnableColor()
	case configv3.ColorDisabled:
		colorPrinter.DisableColor()
	}

	if bold {
		colorPrinter = colorPrinter.Add(color.Bold)
	}
	f := colorPrinter.SprintFunc()
	return f(message)
}

func getFirstSet(list []map[string]interface{}) map[string]interface{} {
	if list == nil || len(list) == 0 {
		return map[string]interface{}{}
	}
	return list[0]
}
