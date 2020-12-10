# API Usage

## Get an authentication token

This api returns a token. Default validity is 1 year.

 #### POST /tokens
 #### Payload
 ```
 {  
 	"username":"admin",
 	"password":"admin"
 }
 ```
 
 **Example**
 
`curl -XPOST http://localhost:8080/tokens -H "Content-Type: application/json" -d '{"username":"admin","password":"admin"}'`

Response:

`
{"code":200,"expire":"2020-06-18T20:58:56-07:00","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI1MzkxMzYsImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2MTAwMzEzNn0.FpPkK-5Tf-o_HAHjbf9FDc15w3Dq8jXPRco6ucR5vsk"}
`

## Create an event from Github Webhook
This api creates an event based on a request sent from Github Webhook.

#### POST /github/events
#### Header
```
X-Github-Event: pull_request
```
#### Payload
Partial view of the payload received, showing what is used by the service.
```
{  
	"action":"open",
	"pull_request": {
	  "head": {
	    "sha": "geloxl23l"
	    },
	  "id": 2132435,
	  "number" : 232,
	  "merged_commit_sha":"fdsdsf3r",
	  "merged": true
	},
	"repository": {
		"id":"",
		"name":""
	}
}
```
**Description**

See [Github documentation](https://developer.github.com/v3/activity/events/types/#pullrequestevent) 

## Create an event
This api creates an event and the resources the event may require (commit, incident).
It return the event created.

#### POST /events
#### Header
```
Authorization: Bearer <JWT obtained from POST /tokens api>
```
#### Payload
```
{  
	"type":"commit",
	"pipeline_id":"api",
	"environment":"dev",
	"commit_id":"1un1queid",
	"status":"success",   
	"timestamp":1561003507
}
```
**Description**
- type:
 	- "commit": creation of the commit
 	- "submit": submission of the commit for code review
 	- "approve": commit approved
 	- "deploy": commit deployed to `environment`
- pipeline_id: identifier of the source of events
- environment: identifier of the environment where the event occurred
- status: "success" or "failure" status of event
- timestamp: epoch time in seconds of the event

##### Example
`
curl -XPOST http://localhost:8080/events -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI1NDAwMzksImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2MTAwNDAzOX0.gKePomW78fxblYN9GIzNBgGbV8qFc0Kg71CBKqhzIxw" -d '{ "type":"commit", "pipeline_id":"api", "environment": "dev", "commit_id":"1un1queid", "status":"success", "timestamp":1561003507}'
`

Response:

```
{
	"id":"2b1914dc-8dce-4d1a-82e6-e612d7bbf69f",
	"pipeline_id":"api",
	"commit_id":"1un1queid",
	"incident_id":"",
	"timestamp":1561003507,
	"type":"commit",
	"status":"success",
	"environment":"dev"}
```

## Get Events

#### GET /events
#### Header
```
Authorization: Bearer <JWT obtained from POST /tokens api>
```
**Example**

`curl -XGET http://localhost:8080/events -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTI1NDAwMzksImlkIjoiYWRtaW4iLCJvcmlnX2lhdCI6MTU2MTAwNDAzOX0.gKePomW78fxblYN9GIzNBgGbV8qFc0Kg71CBKqhzIxw" `

**Response**
```
{ 
	"items":[
		{ "id":"05b7a1aa-8c69-464a-b325-5c430f7b6676",
		  "pipeline_id":"test",
		  "commit_id":"ab497a5a-c041-4e01-8b7f-c567170d62fc",
		  "incident_id":"",
		  "timestamp":1560964385,
		  "type":"deploy",
		  "status":"failure",
		  "environment":"dev"
		},
		.../....
		]
}
```
