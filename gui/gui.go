package gui

import (
	"log"
	"os"
	"unicode"

	"main/database"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	author     = "ТБД Богданович"
	windowName = "NAU"
)

type PasswordEntries struct {
	Username       *widget.Entry
	Password       *widget.Entry
	StrongPassword *widget.Check
}

type App struct {
	fyne.App
	db *database.DB
}

func New(db *database.DB) *App {
	return &App{
		App: app.New(),
		db:  db,
	}
}

func (a *App) PasswordEntries() PasswordEntries {
	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("Username")

	passwordEntry := widget.NewEntry()
	passwordEntry.SetPlaceHolder("Password")

	strongPasswordCheck := widget.NewCheck("Strong (If not checked set default rules for password)", func(_ bool) {})

	return PasswordEntries{
		Username:       usernameEntry,
		Password:       passwordEntry,
		StrongPassword: strongPasswordCheck,
	}
}

func (a *App) Build() error {
	myWindow := a.NewWindow(author)
	about := a.about()
	adduser := a.addUser()
	changePassword := a.changePassword()
	logIn := a.logIn()

	content := container.New(
		layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 50}),
		about, adduser, changePassword, logIn,
	)

	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.SetContent(content)
	myWindow.ShowAndRun()

	return nil
}

func (a *App) about() *widget.Button {
	about := widget.NewButton("About", func() {
		a.showWindowText("БІ-443Б Богданович Олексій")
	})

	return about
}

func (a *App) addUser() *widget.Button {
	addUserBtn := widget.NewButton("AddUser", func() {
		w := a.NewWindow(windowName)
		entries := a.PasswordEntries()

		submit := widget.NewButton("Enter", func() {

			username := entries.Username.Text
			password := entries.Password.Text

			if entries.StrongPassword.Checked && checkStrongPassword(password) && !a.db.HasUser(username) {
				err := a.db.CreateUser(username, password)
				if err != nil {
					log.Panic(err)
				}
				a.showWindowText("Registered with strong password")
			} else if a.db.HasUser(username) || entries.StrongPassword.Checked && !checkStrongPassword(password) {
				a.showWindowText("Rules dose not match or user with this name already exists")

			} else if !entries.StrongPassword.Checked && !a.db.HasUser(username) {
				err := a.db.CreateUser(username, password)
				if err != nil {
					log.Panic(err)
				}
				a.showWindowText("Registered with light password")
			}
		})

		content := container.NewVBox(
			widget.NewLabel("Enter user credentials:"),
			entries.Username,
			entries.Password,
			entries.StrongPassword,
			submit,
		)

		w.Resize(fyne.Size{Width: 500, Height: 500})
		w.SetContent(content)
		w.Show()
	})
	return addUserBtn
}

func (a *App) changePassword() *widget.Button {
	changePasswordBtn := widget.NewButton("ChangePass", func() {
		w := a.NewWindow(windowName)

		entries := a.PasswordEntries()

		submit := widget.NewButton("Enter", func() {

			username := entries.Username.Text
			password := entries.Password.Text

			if entries.StrongPassword.Checked && checkStrongPassword(password) && a.db.HasUser(username) {
				err := a.db.ChangePassword(password, username)
				if err != nil {
					log.Panic(err)
				}
				a.showWindowText("Updated")
			} else if !a.db.HasUser(username) || entries.StrongPassword.Checked && !checkStrongPassword(password) {
				a.showWindowText("Rules dose not match or user not exists")

			} else if !entries.StrongPassword.Checked && a.db.HasUser(username) {
				err := a.db.ChangePassword(password, username)
				if err != nil {
					log.Panic(err)
				}
				a.showWindowText("Updated")
			}
		})

		content := container.NewVBox(
			widget.NewLabel("Enter data to change password for existing user:"),
			entries.Username,
			entries.Password,
			entries.StrongPassword,
			submit,
		)

		w.Resize(fyne.Size{Width: 500, Height: 500})
		w.SetContent(content)
		w.Show()
	})
	return changePasswordBtn
}

func (a *App) logIn() *widget.Button {
	logInBtn := widget.NewButton("Log In", func() {
		w := a.NewWindow(windowName)

		entries := a.PasswordEntries()

		submit := widget.NewButton("Enter", func() {

			username := entries.Username.Text
			password := entries.Password.Text

			if a.db.LogIn(username, password) {
				a.authorization(username)
				return
			}
			a.showWindowText("Invalid credentials")
		})

		content := container.NewVBox(
			widget.NewLabel("Enter data to log in:"),
			entries.Username,
			entries.Password,
			submit,
		)

		w.Resize(fyne.Size{Width: 500, Height: 500})
		w.SetContent(content)
		w.Show()
	})
	return logInBtn
}

func (a *App) authorization(username string) {
	w := a.NewWindow(windowName)

	level, err := a.db.Level(username)
	if err != nil {
		log.Panic(err)
	}

	read := widget.NewButton("Read", func() {
		data, err := os.ReadFile("file.txt")
		if err != nil {
			log.Print(err)
			return
		}
		a.showWindowText(string(data))
	})
	write := widget.NewButton("Write", func() {
		data := []byte("Hello my Username is: " + username)
		err = os.WriteFile("file.txt", data, 0644)
		if err != nil {
			log.Print(err)
			return
		}
		a.showWindowText("Written")
	})
	deleteFile := widget.NewButton("Delete", func() {
		err = os.Remove("file.txt")
		if err != nil {
			log.Print(err)
			return
		}
		a.showWindowText("Deleted")
	})
	open := widget.NewButton("Open", func() {
		r, _ := fyne.LoadResourceFromURLString("https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png?20210701224649")
		img := canvas.NewImageFromResource(r)
		w := a.NewWindow(windowName)
		w.SetContent(img)
		w.Resize(fyne.Size{Width: 200, Height: 200})
		w.Show()
	})
	deleteExe := widget.NewButton("Delete Exe", func() {
		_ = os.Remove("main")
		if err != nil {
			log.Print(err)
			return
		}
		a.showWindowText(string("Exe file removed"))
	})

	content := container.NewVBox()
	switch level {
	case 0:
		content = container.NewVBox(read)
	case 1:
		content = container.NewVBox(read, write)
	case 2:
		content = container.NewVBox(read, write, deleteFile)
	case 3:
		content = container.NewVBox(read, write, deleteFile, open)
	case 4:
		content = container.NewVBox(read, write, deleteFile, open, deleteExe)

	default:
		a.showWindowText("User does not exists")
	}

	w.Resize(fyne.Size{Width: 500, Height: 500})
	w.SetContent(content)
	w.Show()
}

func (a *App) showWindowText(text string) {
	w := a.NewWindow(windowName)
	w.SetContent(widget.NewLabel(text))
	w.Resize(fyne.Size{Width: 200, Height: 200})
	w.Show()
}

func checkStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	containsLetter := false
	for _, char := range password {
		if unicode.IsLetter(char) {
			containsLetter = true
			break
		}
	}
	if !containsLetter {
		return false
	}

	var containsDigit bool
	for _, char := range password {
		if unicode.IsDigit(char) {
			containsDigit = true
			break
		}
	}
	if !containsDigit {
		return true
	}

	return true
}
