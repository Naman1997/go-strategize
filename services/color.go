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
		fmt.Println()
		os.Exit(0)
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
	if colorText != INPUT && colorText != ERROR {
		fmt.Println()
	}
}

func Help() {
	helpText := `
Usage: go-stratergize [options] [<arguments>]

The available commands for execution are listed below.
go-stratergize will attempt to do the following:
> Clone terraform & ansible repos
> Execute terraform apply
> Attempt to SSH in all VMs in ansible inventory
> Copy the ssh-key with ssh-copy-id
> Execute playbooks provided in the repo

Template Options:
  -proxmox-k8s=true     Uses the following repos for creating a k8s cluster:
  			> 'https://github.com/Naman1997/proxmox-terraform-template-k8s'
			> 'https://github.com/Naman1997/cluster-management'

Terraform Options:
  -terraform=URL        URL for your terraform repo. It's assumed that main.tf is
			in the root of this repo.
  -var-file=path        Path to terraform.tfvars that will be used with
			'terraform apply'.

Ansible Options:
  -ansible=URL          URL for your ansible repo.
  -inventory=path       Expected relative path of inventory file to repo folder.
  			This can be created after execution of terraform apply.
			(default = /etc/ansible/hosts)
  -ansible-req=path     Expected path to requirements.yaml file. This is mandatory
			if you're populating '-ansible' flag.
  -ansible-play=path    Expected relative path of playbooks dir to repo folder.
			This is mandatory if you're populating '-ansible' flag.
  -ansible-var=path     Expected path of your vars.json file. This is mandatory
			if you're populating '-ansible' flag.

SSH Options:
  -ssh-user             Name of user that will be used for SSH and ansible playbooks.
			(default = root)
  -ssh-key              Private key for SSH. (default = ~/.ssh/id_rsa)
  -strict=false         Do not ask for host verification. (default = true)

Other Options:
  version               Returns the version
  help                  Prints this help section
`
	fmt.Println(strings.TrimSpace(helpText))
}
