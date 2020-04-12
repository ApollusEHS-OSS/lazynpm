package gui

import (
	"os"

	"github.com/jesseduffield/gocui"
)

// when a user runs lazynpm with the LAZYNPM_NEW_DIR_FILE env variable defined
// we will write the current directory to that file on exit so that their
// shell can then change to that directory. That means you don't get kicked
// back to the directory that you started with.
func (gui *Gui) recordCurrentDirectory() error {
	if os.Getenv("LAZYNPM_NEW_DIR_FILE") == "" {
		return nil
	}

	// determine current directory, set it in LAZYNPM_NEW_DIR_FILE
	dirName, err := os.Getwd()
	if err != nil {
		return err
	}

	return gui.OSCommand.CreateFileWithContent(os.Getenv("LAZYNPM_NEW_DIR_FILE"), dirName)
}

func (gui *Gui) handleQuitWithoutChangingDirectory(g *gocui.Gui, v *gocui.View) error {
	gui.State.RetainOriginalDir = true
	return gui.quit(v)
}

func (gui *Gui) handleQuit(g *gocui.Gui, v *gocui.View) error {
	gui.State.RetainOriginalDir = false
	return gui.quit(v)
}

func (gui *Gui) quit(v *gocui.View) error {
	if gui.State.Updating {
		return gui.createUpdateQuitConfirmation(gui.g, v)
	}

	if gui.Config.GetUserConfig().GetBool("confirmOnQuit") {
		return gui.createConfirmationPanel(createConfirmationPanelOpts{
			returnToView:       v,
			returnFocusOnClose: true,
			prompt:             gui.Tr.SLocalize("ConfirmQuit"),
			handleConfirm: func() error {
				return gocui.ErrQuit
			},
		})
	}

	return gocui.ErrQuit
}
