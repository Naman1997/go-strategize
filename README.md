# go-stratergize

[![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/Naman1997/go-strategize/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/Naman1997/go-strategize)](https://goreportcard.com/report/github.com/Naman1997/go-strategize)
[![Open Source Love svg1](https://badges.frapsoft.com/os/v1/open-source.svg?v=103)](https://github.com/ellerbrock/open-source-badges/)

## Overview

go-stratergize is a simple CLI program that provisions your infrastructure using terraform and configures it using ansible. To connect the two, it assumes that you're writing the ansible inventory to a file somewhere with terraform.

To learn how to do that trick read [local file example](https://github.com/Naman1997/proxmox-terraform-template-k8s/blob/main/main.tf#L97) and [hosts.tmpl](https://github.com/Naman1997/proxmox-terraform-template-k8s/blob/main/hosts.tmpl).

go-stratergize does the following in order:
- Clone your terraform and ansible playbooks
- [Optional] Copies your .tfvars file(terraform) over to the cloned terraform repo folder
- Makes sure some playbooks are available in the provided folder path for playbooks
- Executes `terraform init` and `terraform apply`
- Reads the (hopefully) generated ansible hosts file from the provided path
- Copies ssh public keys for passwordless auth using ssh-copy-id
- Executes `ansible-galaxy collection install`
- Executes all the playbooks present in a folder

It's specially useful for people who usually have very long running ansible playbooks and need them to be running in a closed network.
Imagine you need to provision some VMs and then install vanila k8s or whatever application on them. In a closed network you probably don't care about SSH key verification. Assuming the these two points are valid, you can utilize go-stratergize with flag '-strict=false' to skip hostname verification and go grab a cup of coffee while your infrastructure is being provisioned and configured!

## Options

- Template options
    - -proxmox-k8s=true : Uses 'https://github.com/Naman1997/proxmox-terraform-template-k8s' and 'https://github.com/Naman1997/cluster-management' for creating a k8s cluster

- Terraform Options:
    - -terraform=URL : URL for your terraform repo.
    - -var-file=path : Path to optional terraform.tfvars.

- Ansible Options:
    - -ansible=URL : URL for your ansible repo.
    - -inventory=path : Expected relative path of inventory file to repo folder. (default = /etc/ansible/hosts)
    - -ansible-req=path : Expected path to requirements.yaml file.
    - -ansible-play=path : Expected relative path of playbooks dir to repo folder.
    - -ansible-var=path : Expected path of your vars.json file.

- SSH Options:
    - -ssh-user : Username used for SSH and ansible (default = root)
    - -ssh-key : Private key for SSH. (default = ~/.ssh/id_rsa)
    - -strict=false : Do not ask for host verification. (default = true)

- Other Options:
    - version : Returns the version
    - help : Prints this help section
