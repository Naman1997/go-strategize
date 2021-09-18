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
	"sync"

	"github.com/Naman1997/go-stratergize/services"
	"github.com/relex/aini"
)

var (
	template bool   = false
	clone    bool   = true
	response string = "N"
	wg       sync.WaitGroup
)

const (
	base                   string = "https://github.com/Naman1997/"
	template_terraform     string = "proxmox-terraform-template-k8s.git"
	template_ansible       string = "cluster-management.git"
	inventory_default_path string = "/proxmox-terraform-template-k8s/ansible/hosts"
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
	inventory_flag := flag.String("inventory", "", "Expected file path to your ansible inventory")
	ssh_username_flag := flag.String("ssh-user", "root", "Username for SSH")
	ssh_key_flag := flag.String("ssh-key", "~/.ssh/id_rsa", "Private key for SSH")

	//Extract flag data
	flag.Parse()
	terraform_vars_file := *terraform_vars_file_flag
	terraform_repo := *terraform_repo_flag
	ansible_repo := *ansible_repo_flag
	inventory_file := *inventory_flag
	ssh_username := *ssh_username_flag
	ssh_key := *ssh_key_flag

	//Update clone flag if any of these flags are passed
	if len(terraform_repo) > 0 || len(ansible_repo) > 0 || len(inventory_file) > 0 {
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
	dir, err := os.Getwd()
	if template {
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		folder := services.FormatRepo(terraform_repo)
		services.Terraform_init(folder, dir)
		services.Terraform_apply(folder, dir)
	}

	//Validate ansible inventory exists
	if template {
		inventory_file = dir + inventory_default_path
	} else {
		inventory_file = services.HomeFix(inventory_file, homedir)
	}
	_, err = services.Exists(inventory_file, homedir)
	if err != nil {
		log.Fatalf("[ERROR] [Ansible inventory] %v", err)
	}

	// Parse the inventory
	file, err := os.Open(inventory_file)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	inventoryReader := bufio.NewReader(file)
	inventory, err := aini.Parse(inventoryReader)
	if err != nil {
		log.Fatalf("[ERROR] [Ansible inventory] %v", err)
	}

	//Attempt to SSH into all VMs
	for _, h := range inventory.Groups["all"].Hosts {
		wg.Add(1)
		services.ValidateConn(ssh_username, ssh_key, homedir, h.Vars["ansible_host"], h.Vars["ansible_port"], &wg)
	}
	wg.Wait()

	//Exit
	if template {
		fmt.Println("[INFO] Execution completed for template!")
	} else {
		fmt.Println("[ERROR] Sorry, non-template execution is not yet supported")
		os.Exit(1)
	}
}

func templateCheck() *bufio.Scanner {
	fmt.Print("[INPUT] Clone and execute default proxmox template?[y/N] ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input
}

func askRepoUrl(repotype string) string {
	fmt.Print("[INPUT] What's your " + repotype + " repo URL? ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	response := input.Text()
	if !services.IsURL(response) {
		log.Fatal("[ERROR] [Invalid value for", repotype, "flag] ")
	}
	return response
}
