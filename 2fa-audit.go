package main

import (
  "flag"
  "fmt"
  "os"
  "github.com/octokit/go-octokit/octokit"
)

type Member struct {
  Login string `json:"login"`
}

var accessToken string
var orgName string

func init() {
  flag.StringVar(&accessToken, "token", "", "user access token")
  flag.Parse()
  orgName = flag.Arg(0)

  if accessToken == "" {
    accessToken = os.Getenv("GH_ACCESS_TOKEN")
  }

  fmt.Println("access token: ", accessToken)
  fmt.Println("org name: ", orgName)
}

func main() {
  if orgName == "" {
    fmt.Println("Usage: [-token <token>] <org>")
    fmt.Println("       Print list of organization members *without* 2fa enabled.")
    fmt.Println("       Looks for access token via -token or GH_ACCESS_TOKEN")
    fmt.Println("       Access token must be for member of target org")
    return
  }

  memberURL := octokit.Hyperlink("/orgs/{org}/members{?filter}") // type,page,per_page,sort
  url, err := memberURL.Expand(octokit.M{"org": orgName, "filter": "2fa_disabled"})

  if err != nil {
    fmt.Println("There was an error: ", err)
    return
  }

  client := octokit.NewClient(octokit.TokenAuth{AccessToken: accessToken})
  req, _ := client.NewRequest(url.String())
  var members []Member
  req.Get(&members)

  if len(members) == 0 { fmt.Println("No users have 2fa disabled.") }
  for i:=0; i<len(members); i++ {
    fmt.Println(members[i].Login)
  }
}
