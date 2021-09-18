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

const (
	base               string = "https://github.com/Naman1997/"
	template_terraform string = "proxmox-terraform-template-k8s.git"
	template_ansible   string = "cluster-management.git"
)

func main() {

	//Get home dir
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	homedir := usr.HomeDir

	//Get flag inputs
	terraform_vars_file_flag := flag.String("var-file", "", "Path to .tfvars file")
	terraform_repo_flag := flag.String("terraform", "", "URL to your terraform repo")
	ansible_repo_flag := flag.String("ansible", "", "URL to your ansible repo")
	inventory_flag := flag.String("inventory", "", "Path to your ansible inventory")

	//Extract flag data
	flag.Parse()
	terraform_vars_file := *terraform_vars_file_flag
	terraform_repo := *terraform_repo_flag
	ansible_repo := *ansible_repo_flag
	inventory := *inventory_flag

	//Update clone flag if any of these flags are passed
	if len(terraform_repo) > 0 || len(ansible_repo) > 0 || len(inventory) > 0 {
		clone = false
		fmt.Println("[INFO] Not using proxmox template for current execution")
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
		terraform_repo = base + template_terraform
		ansible_repo = base + template_ansible
		services.CloneRepos(terraform_repo, ansible_repo, homedir)
	} else if strings.EqualFold(response, "N") {

		//Make sure both repos are present
		if len(terraform_repo) == 0 {
			terraform_repo = askRepoUrl("terraform")
		}
		if len(ansible_repo) == 0 {
			ansible_repo = askRepoUrl("ansible")
		}

		//Validate ansible inventory
		_, err := services.Exists(inventory, homedir)
		if err != nil {
			log.Fatalf("[ERROR] [Invalid value for ansible flag] %v", err)
		}

		services.CloneRepos(terraform_repo, ansible_repo, homedir)
	}

	//Copy over .tfvars file if specified
	vars, err := services.Exists(terraform_vars_file, homedir)
	if err != nil {
		log.Fatalf("[ERROR] [Invalid value for var-file flag] %v", err)
	}
	if vars {
		if strings.Contains(terraform_vars_file, "~/") {
			terraform_vars_file = filepath.Join(homedir, terraform_vars_file[2:])
		}
		newfile := services.FormatRepo(terraform_repo) + "terraform.tfvars"
		bytes, err := services.Copy(terraform_vars_file, newfile)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		m := fmt.Sprintf("[INFO] Copied %d bytes to "+newfile, bytes)
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
	// conn, err := services.ConnectInsecure(*ssh_username, *ssh_key, homedir)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// output, err := conn.SendCommands("echo 'INFO: SSH successful with: '", "hostname")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(output))

	if template {
		fmt.Println("[INFO] Execution completed for template!")
	} else {
		fmt.Println("[ERROR] Sorry, non-template execution is not yet supported")
		os.Exit(1)
	}
}

func templateCheck() *bufio.Scanner {
	fmt.Print("[INPUT] Clone and execute default proxmox template?[y/N]")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input
}

func askRepoUrl(repotype string) string {
	fmt.Println("[INPUT] What's your " + repotype + " repo URL?")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	response := input.Text()
	if !services.IsURL(response) {
		log.Fatal("[ERROR] [Invalid value for", repotype, "flag] ")
	}
	return response
}
