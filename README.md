# ticken-pvtbc-connector
Module used to connect to the private blockhain

## How to install this package

1) Generate GITHUB_TOKEN [here](https://github.com/settings/tokens) 
2) `export GITHUB_TOKEN=xxx`
3) `git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/ticken-ts".insteadOf "https://github.com/ticken-ts"`
4) Set the following environment variable `export GOPRIVATE="github.com/ticken-ts/ticken-pvtbc-connector"`

***Important***:

* Set the expiration in 30 days 
* Do not select any other permissions to the access token furthermore than:r
  * `repo`: give full control of private repositories
  * `admin:org`:  give full control of orgs and teams, read and write org projects