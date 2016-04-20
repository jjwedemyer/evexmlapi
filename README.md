<p>This is a client library for connecting to Eve Online's XML API site using Golang.</p>


<h2>Usage</h2>

Start by creating a resource with:
	skillQueue := evexmlapi.NewCharSkillQueue()	
	<---The list of resource functions can be found in 'resources.go'.
	
Next create an httpRequest:
	httpRequest  := evexmlapi.XMLServerRequest()
	-	You can override the base URL with 
			httpRequest.overrideBaseURL = "https://api.testeveonline.com/"
	
	-	Set Params with 
			httpRequest.params = map[string][]string{
									"characterID": []string{},
									"keyID":       []string{},
									"vCode":       []string{},
								}
								
Then fetch your document:
	d, err := evexmlapi.Fetch(skillQueue, httpRequest)


<h3>Example</h3>
	package main 
	
	import (
		"github.com/jovon/eve-xmlapi-go"
		"fmt"
	)
	
	func main() {
		status := evexmlapi.NewServerStatus()
		hr := evexmlapi.NewRequest()
		d, err := evexmlapi.Fetch(status, hr)	
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", d)
	}