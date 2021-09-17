package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	services "github.com/Naman1997/go-stratergize/services"
)

var template bool = false

// var clone bool = true

func main() {
	dir, _ := os.Getwd()
	var terraform_vars_file string = *flag.String("var-file", "", "Path to .tfvars file")
	terraform_repo := flag.String("terraform", "", "Link to your terraform repo")
	ansible_repo := flag.String("ansible", "", "Link to your ansible repo")
	flag.Parse()
	// fixme: Do not ask for cloning template if terraform-repo and ansible-repo are defined
	// if len(terraform_repo) > 0 && len(ansible_repo) > 0 {
	// 	clone = false
	// }

	//Keep checking for [y/N]
	input := templateCheck()
	for !strings.EqualFold(input.Text(), "Y") && !strings.EqualFold(input.Text(), "N") {
		input = templateCheck()
	}

	//Clone terraform and ansible repos
	if strings.EqualFold(input.Text(), "Y") {
		template = true
		base := "https://github.com/Naman1997/"
		*terraform_repo = base + "proxmox-terraform-template-k8s.git"
		*ansible_repo = base + "cluster-management.git"
		services.CloneRepos(*terraform_repo, *ansible_repo)
	} else if strings.EqualFold(input.Text(), "N") && (len(*terraform_repo) > 0 || len(*ansible_repo) > 0) {
		services.CloneRepos(*terraform_repo, *ansible_repo)
	}

	//Copy over .tfvars file if specified
	if len(terraform_vars_file) > 0 {
		tfvars_exists, _ := services.Exists(terraform_vars_file)
		fmt.Println(tfvars_exists)
		if tfvars_exists {
			fmt.Println("Copying .tfvars file over to the required folder")
			cpCmd := exec.Command("cp", terraform_vars_file, dir+"/"+*terraform_repo+"/terraform.tfvars")
			_ = cpCmd.Run()
		} else {
			fmt.Println("Error: Provided path not found!")
			os.Exit(1)
		}
	}

	//Initialize and apply with terraform
	if template {
		folder := services.FormatRepo(*terraform_repo)
		services.Terraform_init(folder, dir)
		services.Terraform_apply(folder, dir)
	}

	// Attempt to SSH
	// fixme: Refactor this section
	// ssh_username := flag.String("ssh-user", "naman", "Username for SSH")
	// ssh_key := flag.String("ssh-key", "id_rsa", "Name of SSH private key")
	// flag.Parse()
	// conn, err := services.ConnectInsecure(*ssh_username, *ssh_key)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// output, err := conn.SendCommands("echo 'INFO: SSH successful with: '", "hostname")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(output))

	if template {
		fmt.Println("Execution completed for template!")
	} else {
		fmt.Println("Sorry, non-template execution is not yet supported")
		os.Exit(1)
	}
}

func templateCheck() *bufio.Scanner {
	fmt.Print("Clone and execute default proxmox template?[y/N]")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input
}
