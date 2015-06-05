package main

import (
  "fmt"
  "os"
  "github.com/octokit/go-octokit/octokit"
)

type Member struct {
  Login string `json:"login"`
}

func main() {
  if len(os.Args[1:]) < 1 {
    fmt.Println("Usage: given an org, print list of users *without* 2fa enabled")
    return
  }

  orgName := os.Args[1]

  memberURL := octokit.Hyperlink("/orgs/{org}/members{?filter}") // type,page,per_page,sort
  url, err := memberURL.Expand(octokit.M{"org": orgName, "filter": "2fa_disabled"})

  if err != nil {
    fmt.Println("There was an error: ", err)
    return
  }

  client := octokit.NewClient(nil)
  req, _ := client.NewRequest(url.String())
  var members []Member
  req.Get(&members)

  if len(members) == 0 { fmt.Println("No users have 2fa disabled.") }
  for i:=0; i<len(members); i++ {
    fmt.Println(members[i].Login)
  }
}
