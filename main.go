package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	services "github.com/Naman1997/go-strategize/services"
	"github.com/relex/aini"
)

var (
	templateAnsible   string = "cluster-management"
	templateTerraform string = "proxmox-terraform-template-k8s"
)

const (
	base                        string = "https://github.com/Naman1997/"
	inventoryTemplatePath       string = "ansible/hosts"
	ansibleTemplateRequirements string = "requirements.yaml"
	ansibleTemplatePlaybooks    string = "/playbooks/"
	ansibleTemplateVars         string = "playbooks/vars.json"
	defaultAnsibleInventory     string = "/etc/ansible/hosts"
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
	terraformVarsFileFlag := flag.String("var-file", "", "")
	terraformRepoFlag := flag.String("terraform", "", "")
	ansibleRepoFlag := flag.String("ansible", "", "")
	inventoryFlag := flag.String("inventory", defaultAnsibleInventory, "")
	sshUsernameFlag := flag.String("ssh-user", "root", "")
	sshKeyFlag := flag.String("ssh-key", "~/.ssh/id_rsa", "")
	sshStrictFlag := flag.Bool("strict", true, "")
	ansibleRequirementsFlag := flag.String("ansible-req", "", "")
	ansiblePlaybooksFlag := flag.String("ansible-play", "", "")
	ansibleVarsFlag := flag.String("ansible-var", "", "")
	proxmoxK8sFlag := flag.Bool("proxmox-k8s", false, "")

	//Extract flag data
	flag.Parse()
	terraformVarsFile := *terraformVarsFileFlag
	terraformRepo := *terraformRepoFlag
	ansibleRepo := *ansibleRepoFlag
	inventoryFile := *inventoryFlag
	sshUsername := *sshUsernameFlag
	sshKey := *sshKeyFlag
	strict := *sshStrictFlag
	ansibleRequirements := *ansibleRequirementsFlag
	ansiblePlaybooks := *ansiblePlaybooksFlag
	ansibleVars := *ansibleVarsFlag
	template := *proxmoxK8sFlag

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
	if template && (len(terraformRepo) > 0 || len(ansibleRepo) > 0) {
		services.ColorPrint(services.ERROR, "Ansible/Terraform repos cannot be provided when using a template")
	}

	//Clone terraform and ansible repos
	if template {
		template = true
		terraformRepo = base + templateTerraform
		ansibleRepo = base + templateAnsible
		ansibleRequirements = ansibleTemplateRequirements
		ansiblePlaybooks = ansibleTemplatePlaybooks
		ansibleVars = ansibleTemplateVars
		services.CloneRepos(terraformRepo, ansibleRepo, homedir)
	} else {

		//Make sure both repos are present
		if len(terraformRepo) == 0 {
			terraformRepo = askRepoUrl("terraform")
		} else if !services.IsURL(terraformRepo) {
			services.ColorPrint(services.ERROR, "Invalid URL for terraform repo")
		}
		if len(ansibleRepo) == 0 {
			ansibleRepo = askRepoUrl("ansible")
		} else if !services.IsURL(ansibleRepo) {
			services.ColorPrint(services.ERROR, "Invalid URL for ansible repo")
		}

		services.CloneRepos(terraformRepo, ansibleRepo, homedir)
	}

	//Update repo names
	if !template {
		templateAnsible = services.FormatRepo(ansibleRepo)
		templateTerraform = services.FormatRepo(terraformRepo)
	}

	//Copy over .tfvars file if specified
	newfile := filepath.Join(dir, templateTerraform, "terraform.tfvars")
	if len(terraformVarsFile) > 0 {
		_, tfvarsExists := os.Stat(newfile)
		if tfvarsExists != nil {
			terraformVarsFile = services.Validate(terraformVarsFile, homedir)
			bytes, err := services.Copy(terraformVarsFile, newfile)
			if err != nil {
				services.ColorPrint(services.ERROR, "%v", err)
			}
			services.ColorPrint(services.INFO, "Copied %d bytes to "+newfile, bytes)
		} else {
			services.ColorPrint(services.WARN, "tfvars have already been copied!")
		}
	}

	//Validate ansible requirements file
	if len(ansibleRequirements) > 0 {
		ansibleRequirements = services.Validate(filepath.Join(dir, templateAnsible, ansibleRequirements), homedir)
	} else {
		services.ColorPrint(services.WARN, "Ansible requirements file was not provided. Will not execute ansible galaxy collection install.")
	}

	//Validate ansible vars file
	if len(ansibleVars) > 0 {
		ansibleVars = services.Validate(filepath.Join(dir, templateAnsible, ansibleVars), homedir)
	} else {
		services.ColorPrint(services.WARN, "Ansible vars file was not provided.")
	}

	//Validate ansible inventory has been passed
	if len(inventoryFile) == 0 {
		services.ColorPrint(services.ERROR, "Inventory path cannot be empty")
	} else if inventoryFile == defaultAnsibleInventory && !template {
		services.ColorPrint(services.WARN, "Using /etc/ansible/hosts as the default inventory")
	}

	//Validate ansible playbooks exists
	if len(ansiblePlaybooks) == 0 {
		services.ColorPrint(services.ERROR, "Relative Folder path not provided for ansible playbooks")
	}
	ansiblePlaybooks = strings.TrimPrefix(ansiblePlaybooks, "/")
	if strings.Contains(ansiblePlaybooks, "~/") {
		ansiblePlaybooks = services.Validate(ansiblePlaybooks, homedir)
	} else {
		ansiblePlaybooks = filepath.Join(dir, templateAnsible, ansiblePlaybooks)
		_ = services.Exists(ansiblePlaybooks, homedir)
	}
	ansiblePlaybooks = strings.TrimSuffix(ansiblePlaybooks, "/") + "/"

	//Validate at least one yaml file is available in the dir
	files, err := ioutil.ReadDir(ansiblePlaybooks)
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
		services.ColorPrint(services.ERROR, "No yaml files found in path: "+ansiblePlaybooks)
	}

	// Initialize and apply with terraform
	terraformDir := filepath.Join(dir, templateTerraform)
	services.TerraformInit(terraformDir)
	services.TerraformApply(terraformDir)

	//Validate ansible inventory
	if template {
		inventoryFile = services.Validate(filepath.Join(dir, templateTerraform, inventoryTemplatePath), homedir)
	} else {
		inventoryFile = services.Validate(filepath.Join(dir, templateTerraform, inventoryFile), homedir)
	}

	// Parse the inventory
	file, err := os.Open(inventoryFile)
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
		services.ValidateConn(sshUsername, sshKey, homedir, h.Vars["ansible_host"], h.Vars["ansible_port"], strict)
	}

	//Execute all ansible galaxy collect if requirements file is present
	if len(ansibleRequirements) > 0 {
		services.AnsibleGalaxy(ansibleRequirements)
	}

	//Execute all ansible playbooks in the provided folder
	for _, file := range files {
		if strings.Contains(file.Name(), ".yaml") || strings.Contains(file.Name(), ".yml") {
			services.AnsiblePlaybook(ansiblePlaybooks+file.Name(), inventoryFile, ansibleVars, sshUsername)
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
