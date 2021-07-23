To build run:

`go build -o go-rkvst`

By default it will use the profile option, to make use of it create a file called `.jitsuin` with the following example content:

```json
{ "profiles": [
    {
        "name": "default",
        "rkvst-url": "https://rkvst.poc.jitsuin.io/",
        "tenant-id": "XXXXXX",
        "client-id": "XXXXXX",
        "client-secret": "XXXXXX"
    },
    {
        "name": "mandy",
        "rkvst-url": "https://rkvst.poc.jitsuin.io/",
        "tenant-id": "XXXXXXX",
        "client-id": "XXXXXX",
        "client-secret": "XXXXXX"
    }]
} 
```

By default it will look for a `default` profile

Basic usage:

`./go-rkvst i`

Example help:

```bash
10:17:48 [wgodfrey] @ wg-dev:(~/workfolder/dev-work/jitsuin/personal/go-archivist) % 
-> ./go-rkvst                                                                                                                                                                                                                                                (main)
NAME:
   go-rkvst - Utility for Interacting with Jitsuin RKVST

USAGE:
   go-rkvst [global options] command [command options] [arguments...]

COMMANDS:
   init, i   
   asset, a  
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```


Init help:

```
NAME:
   go-rkvst init - 

USAGE:
   go-rkvst init [command options] [arguments...]

OPTIONS:
   --interactive, -i                    Interactive (default: false)
   --write-file, -w                     Writes the token to .auth_token in the current directory (default: false)
   --use-client-secret, -s              Chooses Client Secret Authentication (default: false)
   --base64, --b64                      Base64's the token (default: false)
   --read-file FILE, -r FILE            Select a FILE with Client Secret config (default: ".jitsuin")
   --profile PROFILE, -p PROFILE        Name a PROFILE to read in (default: "default")
   --inline, -l                         Requires RKVST_URL, RKVST_TENANT_ID, RKVST_CLIENT_ID and RKVST_CLIENT_SECRET to be set either in env or using flags (default: false)
   --rkvst-url value, --url value       RKVST Tenant ID [$RKVST_URL]
   --tenant-id value, --tid value       RKVST Tenant ID [$RKVST_TENANT_ID]
   --client-id value, --cid value       RKVST Client ID [$RKVST_CLIENT_ID]
   --client-secret value, --csec value  RKVST Cient Secret [$RKVST_CLIENT_SECRET]
   --help, -h                           show help (default: false)
```