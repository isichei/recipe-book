# A Simple Personal Recipe Book App

A recipe app that takes Markdown files and serves them up in different webpages using HTMX and TEMPL.

Uses (task)[https://taskfile.dev/installation/] to run different commands to find which commands see:

```
task --list-all
```

To run:

```
task run-app 
```

> _**Note:**_ Requires images in `ui/static/img/` where each `<recipe-uid>.jpg` file is a thumbnail for the homepage. And the `<recipe-uid>` is the unique ID given to each recipe. Atm it is just a unique filename for the recipe.

## Installing templ on EC2 machine

> (Assuming you have already ran the setup script)

- `go mod download -json` gets you the installation path
- cd into that...
- `cd cmd/templ && go install` you should have the correct version runnable in `/home/ec2-user/go/bin/`

