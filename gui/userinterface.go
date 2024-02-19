package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"main/database"
	"os"
	"unicode"
)

func Ui(myApp fyne.App) {
	_ = godotenv.Load()

	db, err := database.Open(os.Getenv("DB_URL"))
	if err != nil {
		return
	}

	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS information(username varchar(50), user_password varchar(50))")
	if err != nil {
		return
	}

	myWindow := myApp.NewWindow("ТБД Богданович")

	about := widget.NewButton("About", func() {
		defaultWindow(myApp, "БІ-443Б Богданович Олексій")

	})

	adduser := widget.NewButton("AddUser", func() {
		w := myApp.NewWindow("NAU")

		usernameEntry := widget.NewEntry()
		usernameEntry.SetPlaceHolder("Username")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		strongPas := widget.NewCheck("Strong (If not checked set default rules for password)", func(check bool) {

		})

		submit := widget.NewButton("Enter", func() {

			username := usernameEntry.Text
			password := passwordEntry.Text

			if strongPas.Checked && CheckStrongPassword(password) {
				_ = db.CreateUser(username, password)
				defaultWindow(myApp, "Registered with strong password")
			} else if strongPas.Checked && !CheckStrongPassword(password) {
				defaultWindow(myApp, "Rules dose not match")

			} else if !strongPas.Checked {
				_ = db.CreateUser(username, password)
				defaultWindow(myApp, "Registered with light password")
			}

		})

		content := container.NewVBox(
			widget.NewLabel("Enter user credentials:"),
			usernameEntry,
			passwordEntry,
			strongPas,
			submit,
		)

		w.Resize(fyne.Size{Width: 500, Height: 500})
		w.SetContent(content)
		w.Show()
	})

	changePassword := widget.NewButton("ChangePass", func() {
		w := myApp.NewWindow("NAU")

		usernameEntry := widget.NewEntry()
		usernameEntry.SetPlaceHolder("ExistingUsername")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("NewPassword")

		strongPas := widget.NewCheck("Strong (If not checked set default rules for password)", func(check bool) {

		})

		submit := widget.NewButton("Enter", func() {
			username := usernameEntry.Text
			password := passwordEntry.Text

			if strongPas.Checked && CheckStrongPassword(password) {
				_ = db.CreateUser(username, password)
				defaultWindow(myApp, "Updated")
			} else if strongPas.Checked && !CheckStrongPassword(password) {
				defaultWindow(myApp, "Rules dose not match")

			} else if !strongPas.Checked {
				_ = db.ChangePassword(username, password)
				defaultWindow(myApp, "Updated")
			}

		})

		content := container.NewVBox(
			widget.NewLabel("Enter data to change password for existing user:"),
			usernameEntry,
			passwordEntry,
			submit,
			strongPas,
		)

		w.Resize(fyne.Size{Width: 500, Height: 500})
		w.SetContent(content)
		w.Show()
	})

	logIn := widget.NewButton("Log In", func() {
		w := myApp.NewWindow("NAU")

		usernameEntry := widget.NewEntry()
		usernameEntry.SetPlaceHolder("Username")

		passwordEntry := widget.NewPasswordEntry()
		passwordEntry.SetPlaceHolder("Password")

		submit := widget.NewButton("Enter", func() {
			username := usernameEntry.Text
			password := passwordEntry.Text
			if db.LogIn(username, password) {
				defaultWindow(myApp, "Sign in as: "+username)
				return
			}
			defaultWindow(myApp, "Invalid credentials")
		})

		content := container.NewVBox(
			widget.NewLabel("Enter data to log in:"),
			usernameEntry,
			passwordEntry,
			submit,
		)

		w.Resize(fyne.Size{Width: 500, Height: 500})
		w.SetContent(content)
		w.Show()
	})

	content := container.New(layout.NewGridWrapLayout(fyne.Size{Width: 100, Height: 50}), about, adduser, changePassword, logIn)

	myWindow.Resize(fyne.NewSize(500, 300))
	myWindow.SetContent(content)
	myWindow.ShowAndRun()

}

func CheckStrongPassword(password string) bool {
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

	containsDigit := false
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

func defaultWindow(myApp fyne.App, text string) {
	w1 := myApp.NewWindow("NAU")
	w1.SetContent(widget.NewLabel(text))
	w1.Resize(fyne.Size{Width: 200, Height: 200})
	w1.Show()
}
