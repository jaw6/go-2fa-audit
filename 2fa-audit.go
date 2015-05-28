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
  fmt.Printf("Org Name is: %v", orgName)

  memberURL := octokit.Hyperlink("/orgs/{org}/members{?filter}") // type,page,per_page,sort
  url, err := memberURL.Expand(octokit.M{"org": orgName, "filter": "2fa_disabled"})

  fmt.Println("What is URL: " + url.String())

  if err != nil {
    fmt.Println("There was an error: ", err)
    return
  }
  fmt.Printf("URL is: %v\n", url)

  client := octokit.NewClient(nil)
  req, _ := client.NewRequest(url.String())
  var members []Member
  // result := client.get(url, &members)
  req.Get(&members)

  // if result.HasError() {
  //   fmt.Println("Error: ", err)
  //   return
  // }

  if len(members) == 0 { fmt.Println("No users have 2fa disabled.") }
  for i:=0; i<len(members); i++ {
    fmt.Println(members[i].Login)
  }
}
