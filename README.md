# A Simple Personal Recipe Book App

A recipe app that takes Markdown files and serves them up in different webpages using HTMX and TEMPL.

Uses (task)[https://taskfile.dev/installation/] to run different commands to find which commands see:

```
task --list-all
```

To run locally:

```
task run-app 
```

## Uploading assets to Fly.io

> **Prerequisites**
> - Upload assets to your AWS bucket via the UI (to `app-data/` hardcoded in the aws cli)
> - Make sure templates are built and up to date (`task make-templates`)

To download to fly machine:
- rebuild fly app (`flyctl deploy`)
- ssh to fly machine (`flyctl ssh console`) visit website to start up VM
- run on fly machine `custom-tools sync-from-aws --bucket $BUCKET --data-path data/`

