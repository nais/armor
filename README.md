# armor

Simple API integration with Google compute API.

## Installation

```bash
make armor
```

`````bash
bin/armor
`````

## Endpoints

### Get

`/projects/{project}/policies/{policy}`  
`/projects/{project}/policies"`  
`/projects/{project}/policies/{policy}/rules/{priority}`  
`/projects/{project}/preConfiguredRules`  
`/projects/{project}/backendServices`  

### Post

NB requires policy or rule to be specified in the body.

`/projects/{project}/policies/{policy}/rules`  
`/projects/{project}/policies/{policy}`  
`/projects/{project}/policies/{policy}/backendServices/{backend}`  

### Delete

`/projects/{project}/policies/{policy}`  
`/projects/{project}/policies/{policy}/rules/{priority}`  

### Patch

NB requires policy or rule to be specified in the body.

`/projects/{project}/policies/{policy}`  
`/projects/{project}/policies/{policy}/rules/{priority}`  