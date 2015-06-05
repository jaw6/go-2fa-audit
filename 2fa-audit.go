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

func makeURL(orgName string, page int) (urlString string, err error) {
  memberURL := octokit.Hyperlink("/orgs/{org}/members{?filter,page}") // type,page,per_page,sort
  url, err := memberURL.Expand(octokit.M{"org": orgName, "filter": "2fa_disabled", "page": page})

  if err != nil {
    return
  }

  return url.String(), err
}

func main() {
  if orgName == "" {
    fmt.Println("Usage: [-token <token>] <org>")
    fmt.Println("       Print list of organization members *without* 2fa enabled.")
    fmt.Println("       Looks for access token via -token or GH_ACCESS_TOKEN")
    fmt.Println("       Access token must be for member of target org")
    return
  }

  client := octokit.NewClient(octokit.TokenAuth{AccessToken: accessToken})
  var members []Member
  var pageMembers []Member
  page := 1

  for {
    url, err := makeURL(orgName, page)
    if err != nil { fmt.Println("There was an error: ", err); return }
    req, _ := client.NewRequest(url)
    req.Get(&pageMembers)
    if len(pageMembers) == 0 { break }
    members = append(members, pageMembers...)
    page += 1
  }

  if len(members) == 0 { fmt.Println("No users have 2fa disabled.") }
  for i:=0; i<len(members); i++ {
    fmt.Println(members[i].Login)
  }
}
