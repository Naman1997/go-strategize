package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/Naman1997/go-stratergize/services"
	"github.com/relex/aini"
)

var (
	template           bool   = false
	strict             bool   = true
	template_ansible   string = "cluster-management"
	template_terraform string = "proxmox-terraform-template-k8s"
)

const (
	base                          string = "https://github.com/Naman1997/"
	inventory_template_path       string = "ansible/hosts"
	ansible_template_requirements string = "requirements.yaml"
	ansible_template_playbooks    string = "/playbooks/"
	ansible_template_vars         string = "playbooks/vars.json"
	default_ansible_inventory     string = "/etc/ansible/hosts"
)

func main() {

	//Get home dir
	usr, err := user.Current()
	if err != nil {
		services.ColorPrint(services.ERROR, "%v", err)
	}
	homedir := usr.HomeDir

	//Get working dir
	dir, err := os.Getwd()
	if err != nil {
		services.ColorPrint(services.ERROR, "%v", err)
	}

	flag.Usage = func() {
		services.Help()
		os.Exit(0)
	}

	//Get flag inputs
	terraform_vars_file_flag := flag.String("var-file", "", "")
	terraform_repo_flag := flag.String("terraform", "", "")
	ansible_repo_flag := flag.String("ansible", "", "")
	inventory_flag := flag.String("inventory", default_ansible_inventory, "")
	ssh_username_flag := flag.String("ssh-user", "root", "")
	ssh_key_flag := flag.String("ssh-key", "~/.ssh/id_rsa", "")
	ssh_strict_flag := flag.Bool("strict", true, "")
	ansible_requirements_flag := flag.String("ansible-req", "", "")
	ansible_playbooks_flag := flag.String("ansible-play", "", "")
	ansible_vars_flag := flag.String("ansible-var", "", "")
	proxmox_k8s_flag := flag.Bool("proxmox-k8s", false, "")

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
	template := *proxmox_k8s_flag

	args := flag.Args()
	for _, s := range args {
		if s == "help" {
			services.Help()
			os.Exit(0)
		} else if s == "version" {
			services.ColorPrint("", "go-stratergize v1.0")
			os.Exit(0)
		}
	}

	//Update clone flag if any of these flags are passed
	if template && (len(terraform_repo) > 0 || len(ansible_repo) > 0) {
		services.ColorPrint(services.ERROR, "Ansible/Terraform repos cannot be provided when using a template")
	}

	//Clone terraform and ansible repos
	if template {
		template = true
		terraform_repo = base + template_terraform
		ansible_repo = base + template_ansible
		ansible_requirements = ansible_template_requirements
		ansible_playbooks = ansible_template_playbooks
		ansible_vars = ansible_template_vars
		services.CloneRepos(terraform_repo, ansible_repo, homedir)
	} else {

		//Make sure both repos are present
		if len(terraform_repo) == 0 {
			terraform_repo = askRepoUrl("terraform")
		} else if !services.IsURL(terraform_repo) {
			services.ColorPrint(services.ERROR, "Invalid URL for terraform repo")
		}
		if len(ansible_repo) == 0 {
			ansible_repo = askRepoUrl("ansible")
		} else if !services.IsURL(ansible_repo) {
			services.ColorPrint(services.ERROR, "Invalid URL for ansible repo")
		}

		services.CloneRepos(terraform_repo, ansible_repo, homedir)
	}

	//Update repo names
	if !template {
		template_ansible = services.FormatRepo(ansible_repo)
		template_terraform = services.FormatRepo(terraform_repo)
	}

	//Copy over .tfvars file if specified
	newfile := filepath.Join(dir, template_terraform, "terraform.tfvars")
	if len(terraform_vars_file) > 0 {
		_, tfvars_exists := os.Stat(newfile)
		if tfvars_exists != nil {
			terraform_vars_file = services.Validate(terraform_vars_file, homedir)
			bytes, err := services.Copy(terraform_vars_file, newfile)
			if err != nil {
				services.ColorPrint(services.ERROR, "%v", err)
			}
			services.ColorPrint(services.INFO, "Copied %d bytes to "+newfile, bytes)
		} else {
			services.ColorPrint(services.WARN, "tfvars have already been copied!")
		}
	}

	//Validate ansible requirements file
	if len(ansible_requirements) > 0 {
		ansible_requirements = services.Validate(filepath.Join(dir, template_ansible, ansible_requirements), homedir)
	} else {
		services.ColorPrint(services.WARN, "Ansible requirements file was not provided. Will not execute ansible galaxy collection install.")
	}

	//Validate ansible vars file
	if len(ansible_vars) > 0 {
		ansible_vars = services.Validate(filepath.Join(dir, template_ansible, ansible_vars), homedir)
	} else {
		services.ColorPrint(services.WARN, "Ansible vars file was not provided.")
	}

	//Validate ansible inventory has been passed
	if len(inventory_file) == 0 {
		services.ColorPrint(services.ERROR, "Inventory path cannot be empty")
	} else if inventory_file == default_ansible_inventory && !template {
		services.ColorPrint(services.WARN, "Using /etc/ansible/hosts as the default inventory")
	}

	//Validate ansible playbooks exists
	if len(ansible_playbooks) == 0 {
		services.ColorPrint(services.ERROR, "Relative Folder path not provided for ansible playbooks")
	}
	ansible_playbooks = strings.TrimPrefix(ansible_playbooks, "/")
	if strings.Contains(ansible_playbooks, "~/") {
		ansible_playbooks = services.Validate(ansible_playbooks, homedir)
	} else {
		ansible_playbooks = filepath.Join(dir, template_ansible, ansible_playbooks)
		_ = services.Exists(ansible_playbooks, homedir)
	}
	ansible_playbooks = strings.TrimSuffix(ansible_playbooks, "/") + "/"

	//Validate at least one yaml file is available in the dir
	files, err := ioutil.ReadDir(ansible_playbooks)
	if err != nil {
		services.ColorPrint(services.ERROR, "%v", err)
	}
	yamlPresent := false
	for _, yaml := range files {
		if strings.Contains(yaml.Name(), ".yaml") || strings.Contains(yaml.Name(), ".yml") {
			yamlPresent = true
		}
	}
	if !yamlPresent {
		services.ColorPrint(services.ERROR, "No yaml files found in path: "+ansible_playbooks)
	}

	// Initialize and apply with terraform
	services.Terraform_init(template_terraform, dir)
	services.Terraform_apply(template_terraform, dir)

	//Validate ansible inventory
	if template {
		inventory_file = services.Validate(filepath.Join(dir, template_terraform, inventory_template_path), homedir)
	} else {
		inventory_file = services.Validate(filepath.Join(dir, template_terraform, inventory_file), homedir)
	}

	// Parse the inventory
	file, err := os.Open(inventory_file)
	if err != nil {
		services.ColorPrint(services.ERROR, "%v", err)
	}
	inventoryReader := bufio.NewReader(file)
	inventory, err := aini.Parse(inventoryReader)
	if err != nil {
		services.ColorPrint(services.ERROR, "%v", err)
	}

	//Attempt to SSH into all VMs
	for _, h := range inventory.Groups["all"].Hosts {
		services.ValidateConn(ssh_username, ssh_key, homedir, h.Vars["ansible_host"], h.Vars["ansible_port"], strict)
	}

	//Execute all ansible galaxy collect if requirements file is present
	if len(ansible_requirements) > 0 {
		services.Ansible_galaxy(ansible_requirements)
	}

	//Execute all ansible playbooks in the provided folder
	for _, file := range files {
		if strings.Contains(file.Name(), ".yaml") || strings.Contains(file.Name(), ".yml") {
			services.Ansible_playbook(ansible_playbooks+file.Name(), inventory_file, ansible_vars, ssh_username)
			services.ColorPrint(services.INFO, "Executing ansible playbook: %s", file.Name())
		}
	}
	services.ColorPrint(services.INFO, "Finished executing ansible playbook(s)")

	//Finish execution
	services.ColorPrint(services.INFO, "Execution completed for template!")
}

func askRepoUrl(repotype string) string {
	services.ColorPrint(services.INPUT, "What's your "+repotype+" repo URL? ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	response := input.Text()
	if !services.IsURL(response) {
		services.ColorPrint(services.ERROR, "Invalid URL for "+repotype+" repo")
	}
	return response
}

/*
go run . --var-file ~/proxmox-terraform-template-k8s/terraform.tfvars --ssh-user naman --strict=false --ansible=https://github.com/Naman1997/cluster-management --terraform=https://github.com/Naman1997/proxmox-terraform-template-k8s --ansible-play=playbooks/ --inventory=ansible/hosts --ansible-var=playbooks/vars.json
go run . --var-file ~/proxmox-terraform-template-k8s/terraform.tfvars --ssh-user naman --proxmox-k8s=true
*/
