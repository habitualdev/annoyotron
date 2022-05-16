package main

import (
	"annoyotron/collection"
	"annoyotron/connect"
	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
	_ "embed"
)

//go:embed icon.png
var icon []byte

//go:embed template_config.yaml
var templateConfig []byte

func onReady() {
	systray.SetTitle("Annoyotron")
	systray.SetIcon(icon)
	systray.SetTooltip("Annoyotron")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
		os.Exit(0)
	}()
	systray.AddSeparator()
	viewConfig := systray.AddMenuItem("Edit Config", "Open the config in the default text editor.").ClickedCh
	go func() {
		for {
			<-viewConfig
			exec.Command("xdg-open", "config.yaml").Run()
		}
	}()
}

func onExit() {

}

func loopCheck(userData *collection.YamlConfig) string {
	panicString := ""
	for _, user := range userData.Users {
		for _, host := range user.Hosts {
			c := connect.SshClient{
				Hostname: host,
				Port:     22,
				Username: user.Username,
				Password: user.Password,
				KeyFile:  user.KeyFile,
			}
			panicString = c.SshDiskCheck()
		}
	}
	return panicString
}

func main() {

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		ioutil.WriteFile("config.yaml", templateConfig, 0644)
		println("Created config.yaml")
		println("Please edit it and restart the app.")
		exec.Command("xdg-open", "config.yaml").Run()
		os.Exit(0)
	}

	args := os.Args
	if len(args) == 1 {
		cwd, _ := os.Getwd()
		args := append(os.Args, "--detached")
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = cwd
		cmd.Start()
		println("Starting in background")
		cmd.Process.Release()
		os.Exit(0)
	}

	go systray.Run(onReady, onExit)

	yfile, err := ioutil.ReadFile("config.yaml")
	yData := collection.YamlConfig{}

	for {
		err = yaml.Unmarshal(yfile, &yData)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		if err != nil {
			log.Fatal(err)
		}
		panicString := loopCheck(&yData)
		if panicString != "" {
			dialog.Message("%s", panicString).Title("ANNOYOTRON ALERT").Info()
		}
		refreshDuration, _ := time.ParseDuration(strconv.Itoa(yData.RefreshTime) + "s")
		time.Sleep(refreshDuration)
	}
}
