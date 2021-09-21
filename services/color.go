package services

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

const (
	ERROR = "[ERROR] "
	WARN  = "[WARNING] "
	INPUT = "[INPUT] "
	INFO  = "[INFO] "
)

func ColorPrint(colorText string, text string, option ...interface{}) {
	if colorText == ERROR {
		fmt.Printf(color.RedString(colorText)+text, option...)
		os.Exit(1)
	} else if colorText == INFO {
		fmt.Printf(color.GreenString(colorText))
	} else if colorText == WARN {
		fmt.Printf(color.YellowString(colorText))
	} else if colorText == INPUT {
		fmt.Printf(color.BlueString(colorText))
	} else {
		fmt.Printf(colorText+text, option...)
		fmt.Println()
		return
	}

	fmt.Printf(text, option...)
	if colorText != INPUT {
		fmt.Println()
	}
}

/*
TODO: Add below options
Main commands:
  -template=true        Execute a proxmox template execution. Executes terraform apply
		        using 'https://github.com/Naman1997/proxmox-terraform-template-k8s'
		        and ansible using 'https://github.com/Naman1997/cluster-management'.
		        Will ignore repo related flags.
  -version          Show the current go-stratergize version
*/
func Help() {
	helpText := `
Usage: go-stratergize [command] [options] [<arguments>]

The available commands for execution are listed below.
go-stratergize will attempt to do the following:
> Clone terraform & ansible repos
> Execute terraform apply
> Attempt to SSH in all VMs in ansible inventory
> Copy the ssh-key with ssh-copy-id
> Execute playbooks provided in the repo

Terraform Options:
  -terraform=URL        URL for your terraform repo. It's assumed that main.tf is
			in the root of this repo.
  -var-file=path        Path to terraform.tfvars that will be used with
			'terraform apply'.

Ansible Options:
  -ansible=URL          URL for your ansible repo.
  -inventory=path       Expected path to ansible inventory. This can be created
  			after execution of terraform apply. (default = /etc/ansible/hosts)
  -ansible-req=path     Expected path to requirements.yaml file. This is mandatory
			if you're populating '-ansible' flag.
  -ansible-play=path    Expected path to playbooks dir. This is mandatory
			if you're populating '-ansible' flag.
  -ansible-var=path     Expected path of your vars.json file. This is mandatory
			if you're populating '-ansible' flag.

SSH Options:
  -ssh-user             Name of user that will be used for SSH and ansible playbooks.
			(default = root)
  -ssh-key              Private key for SSH. (default = ~/.ssh/id_rsa)
  -strict=false         Do not ask for host verification. (default = true)

`
	fmt.Println(strings.TrimSpace(helpText))
}
