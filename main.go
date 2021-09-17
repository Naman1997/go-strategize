package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Naman1997/go-stratergize/services"
)

var (
	template bool   = false
	clone    bool   = true
	response string = "N"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	homedir := usr.HomeDir
	terraform_vars_file_flag := flag.String("var-file", "", "Path to .tfvars file")
	terraform_repo_flag := flag.String("terraform", "", "Link to your terraform repo")
	ansible_repo_flag := flag.String("ansible", "", "Link to your ansible repo")

	//Extract flag data
	flag.Parse()
	terraform_vars_file := *terraform_vars_file_flag
	terraform_repo := *terraform_repo_flag
	ansible_repo := *ansible_repo_flag

	if len(terraform_repo) > 0 || len(ansible_repo) > 0 {
		clone = false
	}

	//Check for template repos usage
	if clone {
		input := templateCheck()
		for !strings.EqualFold(input.Text(), "Y") && !strings.EqualFold(input.Text(), "N") {
			input = templateCheck()
		}
		response = input.Text()
	}

	//Clone terraform and ansible repos
	if strings.EqualFold(response, "Y") {
		template = true
		base := "https://github.com/Naman1997/"
		terraform_repo = base + "proxmox-terraform-template-k8s.git"
		ansible_repo = base + "cluster-management.git"
		services.CloneRepos(terraform_repo, ansible_repo)
	} else if strings.EqualFold(response, "N") {

		//Make sure both repos are present
		if len(terraform_repo) > 0 {
			terraform_repo = askRepoUrl("terraform", true).Text()
		}
		if len(ansible_repo) > 0 {
			ansible_repo = askRepoUrl("ansible", true).Text()
		}

		//Validate both repo URLs
		for !services.IsURL(terraform_repo) {
			terraform_repo = askRepoUrl("terraform", false).Text()
		}
		for !services.IsURL(ansible_repo) {
			ansible_repo = askRepoUrl("terraform", false).Text()
		}

		services.CloneRepos(terraform_repo, ansible_repo)
	}

	//Copy over .tfvars file if specified
	if len(terraform_vars_file) > 0 {
		if strings.Contains(terraform_vars_file, "~/") {
			terraform_vars_file = filepath.Join(homedir, terraform_vars_file[2:])
		}
		newfile := services.FormatRepo(terraform_repo) + "terraform.tfvars"
		bytes, err := services.Copy(terraform_vars_file, newfile)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		m := fmt.Sprintf("Copied %d bytes to "+newfile, bytes)
		fmt.Println(m)
	}

	// Initialize and apply with terraform
	if template {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		folder := services.FormatRepo(terraform_repo)
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
	fmt.Print("[INPUT] Clone and execute default proxmox template?[y/N]")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input
}

func askRepoUrl(repotype string, valid bool) *bufio.Scanner {
	if valid {
		fmt.Print("[INPUT] What's your " + repotype + " repo URL?")
	} else {
		fmt.Print("[INPUT] Please provide a valid URL for " + repotype + ":")
	}

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input
}
