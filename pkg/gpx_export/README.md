

> strava oauth

### 0.Create your App

https://www.strava.com/settings/api

### 1.Requesting Access

https://www.strava.com/oauth/authorize

### params

#### client_id

The application’s ID, obtained during registration.

#### redirect_uri

URL to which the user will be redirected after authentication. Must be within the callback domain specified by the application. localhost and 127.0.0.1 are white-listed.

#### response_type

Must be code.

#### approval_prompt

force or auto, use force to always show the authorization prompt even if the user has already authorized the current application, default is auto.

#### scope

```
read: read public segments, public routes, public profile data, public posts, public events, club feeds, and leaderboards
read_all:read private routes, private segments, and private events for the user
profile:read_all: read all profile information even if the user has set their profile visibility to Followers or Only You
profile:write: update the user's weight and Functional Threshold Power (FTP), and access to star or unstar segments on their behalf
activity:read: read the user's activity data for activities that are visible to Everyone and Followers, excluding privacy zone data
activity:read_all: the same access as activity:read, plus privacy zone data and access to read the user's activities with visibility set to Only You
activity:write: access to create manual activities and uploads, and access to edit any activities that are visible to the app, based on activity read access level
```

例子：

```
# http://localhost:8000 can use  `python -m http.server` 

https://www.strava.com/oauth/authorize?client_id=xxx&redirect_uri=http://localhost:8000/token&response_type=code&scope=activity:read
```


### 2. oauth/token 

https://www.strava.com/oauth/token

```shell
# code from pre step get
curl --request POST \
  --url https://www.strava.com/oauth/token \
  --header 'content-type: multipart/form-data' \
  --form client_id=xxx \
  --form client_secret=xxxx \
  --form code=xxx \
  --form grant_type=authorization_code
```

- response 

```shell
{
	"token_type": "Bearer",
	"expires_at": 1703074024,
	"expires_in": 21600,
	"refresh_token": "xxxx",
	"access_token": "xxxx",
	"athlete": {
		"id": 126651800,
		"username": null,
		"resource_state": 2,
		"firstname": ".",
		"lastname": ".",
		"bio": null,
		"city": "Singapore",
		"state": null,
		"country": "Singapore",
		"sex": "M",
		"premium": false,
		"summit": false,
		"created_at": "2023-10-28T05:29:16Z",
		"updated_at": "2023-12-20T05:56:44Z",
		"badge_type_id": 0,
		"weight": null,
		"friend": null,
		"follower": null
	}
}
```

### Reference

https://developers.strava.com/

https://www.strava.com/settings/api

https://github.com/yihong0618/running_page

https://developers.strava.com/docs/getting-started/

https://developers.strava.com/docs/authentication/

https://stackoverflow.com/questions/52880434/problem-with-access-token-in-strava-api-v3-get-all-athlete-activities