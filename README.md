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

> Note: AWS functionality still exists but now we just sync data via the custom TCP connector and/or the fly ssh connection for images

## Uploading assets to Fly.io

> **Prerequisites**
> - Upload assets to your AWS bucket via the UI (to `app-data/` hardcoded in the aws cli)
> - Make sure templates are built and up to date (`task make-templates`)

To download to fly machine:
- rebuild fly app (`flyctl deploy`)
- ssh to fly machine (`flyctl ssh console`) visit website to start up VM
- run on fly machine `custom-tools sync-from-aws --bucket $BUCKET --data-path data/`

## How to sync files across system with custom tcp connection

The app can run a TCP replica on port `9000` and a client can push local recipe
markdown files to it. The connection is authenticated with `TCP_API_KEY`.

> **Current behaviour**
> The TCP sync writes received files into the app's markdown directory only. If
> the app is running with `--db`, this does not update `recipes.db`; regenerate
> the SQL database separately before expecting the synced recipes to appear in
> database-backed search/results.

1. Set the same API key on the app and the client:

```sh
export TCP_API_KEY="<shared-secret>"
```

2. Start the app with TCP file sync enabled. On Fly this is already configured in
   `fly.toml` with `--enable-filesync`; locally you can run:

```sh
bin/web --port 8000 --db recipes.db --static-path ./static/ --enable-filesync
```

3. Build the client CLI if needed:

```sh
task build-cli
```

4. From the machine that has the markdown files, run the TCP client and point
   `--directory` at the folder containing the `.md` recipe files:

```sh
bin/cli start-tcp --address recipe-book-go-htmx.fly.dev:9000 --directory ./static/recipe_mds/
```

For a local app, use the local TCP address instead:

```sh
bin/cli start-tcp --address 127.0.0.1:9000 --directory ./static/recipe_mds/
```

5. Optional: test authentication/connectivity without syncing files:

```sh
bin/cli start-tcp --address recipe-book-go-htmx.fly.dev:9000 --directory ./static/recipe_mds/ --ping-only
```
