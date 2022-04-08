package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

func main() {

	//inital vars
	jhost := ""
	host := ""
	command := ""
	username := ""
	pass := ""
	multihost := ""

	flag.StringVar(&host, "l", "", "Inital host we are jumping from.")
	flag.StringVar(&jhost, "j", "", "Host we are executing command against.")
	flag.StringVar(&command, "c", "", "Command we are executing.")
	flag.StringVar(&username, "u", "", "FORMAT: 'DOMAIN\\USER' we are authenticating with.")
	flag.StringVar(&pass, "p", "", "Password for user")
	flag.StringVar(&multihost, "m", "", "Add hosts comma seperated. FORMAT: 'host1,host2'")
	flag.Parse()

	os := runtime.GOOS
	switch os {
	case "windows":

		if host == "" || command == "" || username == "" || pass == "" {
			fmt.Println("[!] Must supply all flags and values.")
		} else if multihost != "" {

			//Check if user wants to run on multiple systems.
			fmt.Println("Are you sure about running on multiple systems? 'Yes or No' ")
			a1 := ""
			_, err := fmt.Scanln(&a1)
			if err != nil {
				log.Fatal(err)
			}

			if a1 == "Yes" || a1 == "yes" {
				// If yes will execute against read in hosts.
				mhexec(command, username, pass, host, multihost)
			} else {
				fmt.Println("Exiting")
			}
		} else if jhost == "" {

			fmt.Println("Must supply a jump host.")
		} else {
			fmt.Println("[*] Performing double hop.")
			psexecute(command, username, pass, host, jhost)
		}
	default:
		fmt.Println("Only supports windows")
	}

}

// Handles executing against single system.
func psexecute(command string, username string, pass string, host string, jhost string) {

	auth := `$username = "` + username + `";` + ` $password = "` + pass + `" ;$psHost = "` + host + `";$cred = New-Object System.Management.Automation.PSCredential -ArgumentList @($username,(ConvertTo-SecureString -String $password -AsPlainText -Force))`

	jump := auth + ";" + ` Invoke-Command -ComputerName $psHost -ScriptBlock { $username = "` + username + `";$password = "` + pass + `";$psHost = "` + jhost + `";$cred = New-Object System.Management.Automation.PSCredential -ArgumentList @($username,(ConvertTo-SecureString -String $password -AsPlainText -Force)); Invoke-Command -ComputerName $psHost -ScriptBlock { ` + command + ` } -credential $cred } -credential $cred`
	results, stdout, stderr := psrun(jump)

	if results != nil {
		fmt.Println(results)
	} else {
		fmt.Println("\n", strings.Replace(stdout, " ", "", -1), stderr)
	}
}

// Handles executing on multiple systems.
func mhexec(command string, username string, pass string, host string, multihost string) {
	// this function will handle execution for mutliple hosts.
	mh := []string{}

	for _, value := range strings.Split(multihost, ",") {
		// adds hosts from multihosts to slice.
		mh = append(mh, value)
	}

	auth := `$username = "` + username + `";` + ` $password = "` + pass + `" ;$psHost = "` + host + `";$cred = New-Object System.Management.Automation.PSCredential -ArgumentList @($username,(ConvertTo-SecureString -String $password -AsPlainText -Force))`

	for n := range mh {
		fmt.Printf("[*] Running against host: %s \n", mh[n])
		jump := auth + ";" + ` Invoke-Command -ComputerName $psHost -ScriptBlock { $username = "` + username + `";$password = "` + pass + `";$psHost = "` + mh[n] + `";$cred = New-Object System.Management.Automation.PSCredential -ArgumentList @($username,(ConvertTo-SecureString -String $password -AsPlainText -Force)); Invoke-Command -ComputerName $psHost -ScriptBlock { ` + command + ` } -credential $cred } -credential $cred`
		results, stdout, stderr := psrun(jump)

		if results != nil {
			fmt.Println(results)
		} else {
			fmt.Println("\n", strings.Replace(stdout, " ", "", -1), stderr)
		}
		fmt.Println("====================================================")
	}

}

// psrun executes powershell and returns output and errors.
func psrun(command string) (error, string, string) {

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("powershell.exe", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return err, stdout.String(), stderr.String()

}
