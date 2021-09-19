package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/Naman1997/go-stratergize/services"
	"github.com/relex/aini"
)

var (
	template bool   = false
	clone    bool   = true
	response string = "N"
	strict   bool   = true
)

const (
	base                          string = "https://github.com/Naman1997/"
	template_terraform            string = "proxmox-terraform-template-k8s"
	template_ansible              string = "cluster-management"
	inventory_template_path       string = "/proxmox-terraform-template-k8s/ansible/hosts"
	ansible_template_requirements string = "requirements.yaml"
	ansible_template_playbooks    string = "playbooks/"
	ansible_template_vars         string = "playbooks/vars.json"
	default_ansible_inventory     string = "/etc/ansible/hosts"
)

func main() {

	//Get home dir
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	homedir := usr.HomeDir

	//Get working dir
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	//Get flag inputs
	terraform_vars_file_flag := flag.String("var-file", "", "Path to .tfvars file")
	terraform_repo_flag := flag.String("terraform", "", "URL to your terraform repo")
	ansible_repo_flag := flag.String("ansible", "", "URL to your ansible repo")
	inventory_flag := flag.String("inventory", default_ansible_inventory, "Expected path of your ansible inventory file")
	ssh_username_flag := flag.String("ssh-user", "root", "Username for SSH")
	ssh_key_flag := flag.String("ssh-key", "~/.ssh/id_rsa", "Private key for SSH")
	ssh_strict_flag := flag.Bool("strict", true, "Private key for SSH")
	ansible_requirements_flag := flag.String("ansible-req", "", "Requirements file for ansible")
	ansible_playbooks_flag := flag.String("ansible-play", "", "Expected path of your ansible playbooks dir")
	ansible_vars_flag := flag.String("ansible-var", "", "Expected path of your ansible vars file")

	//go run . --var-file ~/proxmox-terraform-template-k8s/terraform.tfvars --ssh-user naman --strict=false --ansible-req cluster-management/requirements.yaml --ansible-play cluster-management/playbooks --ansible-var cluster-management/playbooks/vars.json

	//Extract flag data
	flag.Parse()
	terraform_vars_file := *terraform_vars_file_flag
	terraform_repo := *terraform_repo_flag
	ansible_repo := *ansible_repo_flag
	inventory_file := *inventory_flag
	ssh_username := *ssh_username_flag
	ssh_key := *ssh_key_flag
	strict := *ssh_strict_flag
	ansible_requirements := *ansible_requirements_flag
	ansible_playbooks := *ansible_playbooks_flag
	ansible_vars := *ansible_vars_flag

	//Update clone flag if any of these flags are passed
	if len(terraform_repo) > 0 || len(ansible_repo) > 0 || (len(inventory_file) > 0 && inventory_file != default_ansible_inventory) {
		clone = false
		fmt.Println("[WARN] Not using templates for current execution")
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
		ansible_requirements = ansible_template_requirements
		ansible_playbooks = ansible_template_playbooks
		ansible_vars = ansible_template_vars
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
	newfile := services.FormatRepo(terraform_repo) + "terraform.tfvars"
	if len(terraform_vars_file) > 0 {
		_, tfvars_exists := os.Stat(newfile)
		if tfvars_exists != nil {
			terraform_vars_file = services.Validate(terraform_vars_file, homedir)
			bytes, err := services.Copy(terraform_vars_file, newfile)
			if err != nil {
				log.Fatalf("[ERROR] %v", err)
			}
			m := fmt.Sprintf("[INFO] Copied %d bytes to "+newfile, bytes)
			fmt.Println(m)
		} else {
			fmt.Println("[SKIP] tfvars have already been copied!")
		}
	}

	//Validate ansible requirements file
	if len(ansible_requirements) > 0 {
		ansible_requirements = services.Validate(dir+"/"+template_ansible+"/"+ansible_requirements, homedir)
	}

	//Validate ansible vars file
	if len(ansible_vars) > 0 {
		ansible_vars = services.Validate(dir+"/"+template_ansible+"/"+ansible_vars, homedir)
	}

	//Validate ansible playbooks exists
	if len(ansible_playbooks) > 0 {
		ansible_playbooks = strings.TrimPrefix(ansible_playbooks, "/")
		if strings.Contains(ansible_playbooks, "~/") {
			ansible_playbooks = services.Validate(ansible_playbooks, homedir)
		} else {
			ansible_playbooks = dir + "/" + template_ansible + "/" + ansible_playbooks
			_ = services.Exists(ansible_playbooks, homedir)
		}
		ansible_playbooks = strings.TrimSuffix(ansible_playbooks, "/") + "/"
	} else {
		log.Fatalf("[ERROR] No such file or directory: %v", ansible_playbooks)
	}

	// Initialize and apply with terraform
	if template {
		folder := services.FormatRepo(terraform_repo)
		services.Terraform_init(folder, dir)
		services.Terraform_apply(folder, dir)
	}

	//Validate ansible inventory
	if template {
		inventory_file = services.Validate(dir+inventory_template_path, homedir)
	} else {
		inventory_file = services.Validate(inventory_file, homedir)
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
		services.ValidateConn(ssh_username, ssh_key, homedir, h.Vars["ansible_host"], h.Vars["ansible_port"], strict)
	}

	//Execute all ansible playbooks in the provided folder
	services.Ansible_galaxy(ansible_requirements)
	files, err := ioutil.ReadDir(ansible_playbooks)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".yaml") || strings.Contains(file.Name(), ".yml") {
			services.Ansible_playbook(ansible_playbooks+file.Name(), inventory_file, ansible_vars, ssh_username)
		}
	}
	fmt.Println("[INFO] Finished executing ansible playbook(s)")

	//Finish execution
	fmt.Println("[INFO] Execution completed for template!")
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
		log.Fatal("[ERROR] [Invalid value for", repotype, "repo] ")
	}
	return response
}
