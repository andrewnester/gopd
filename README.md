# gopd - Unofficial Golang SDK for PandaDoc API

See API reference here: https://developers.pandadoc.com

## Examples
### Authenticate
```
// Redirect to get access code
var url = auth.AuthenticateUrl("{redirect_url}", "read+write")
http.Redirect(w, r, url, 302)


creds, err := auth.CreateAccessToken(r.URL.Query().Get("code"), "read+write", "{redirect_url}")
if err != nil {
  log.Fatal(err)
} else {
// use access token here creds.AccessToken
}
```

### Load templates list
```
resp, err := gopd.GetTemplateList()
if err != nil {
  log.Fatal(err)
} else {
  for _, templ := range resp.Results {
    // process template data here
  }
}
```

### Load template details
```
template, err := gopd.GetTemplateDetails("{templateID}")
if err != nil {
  log.Fatal(err)
} else {
  // use template data here
}
```


### Load document details
```
doc, err := gopd.GetDocumentDetails("JMzSJBoeAP8TRN7w5HLzh3")
if err != nil {
  log.Fatal(err)
} else {
  // use  document data here
}
```

### Create new document
```
td := gopd.FromTemplateDocument{
  Name: "From API",
  TemplateUuid: "{template uuid here}",
  Recipients: []gopd.Recipient{
    {Email:"andrew.nester.dev@gmail.com", FirstName:"Andrew", LastName:"Nester"},
  },
}

resp, err := td.Create()
if err != nil {
  log.Fatal(err)
} else {
  // use response data here
}
```

### Send document
```
resp, err := gopd.SendDocument("{document_id}", messageText, isSilent)
if err != nil {
  log.Fatal(err)
} else {
  // use response data here.
}
```




