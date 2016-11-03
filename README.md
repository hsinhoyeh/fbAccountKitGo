## fbAccountKitGo

package `fbAccountKitGo` is a simple golang wrapper for accessing fbAccountKit. Before you use this package, please apply app id from [developer.facebook.com](https://developers.facebook.com) and enable `Allow SMS Login` tab under AccountKit tab.


###### INSTALL

```
go get github.com/hsinhoyeh/fbAccountKitGo

```


###### example
```
package main

const (
	facebookAppID = "<fill-your-appid>"
	appSecret = "<fill-appsecret-under-accountkit-tab>"
)

func main() {
	accountKit := NewAccountKit(
		facebookAppID,
		appSecret)

	// exchange auth code with access token
	c2t, _ := accountKit.VerifyByCode(code)

	// access user's profile from access token
	profile, _ := accountKit.VerifyByToken(
		*c2t.AccessToken,
	)
	fmt.Printf("profile:%v\n", profile)
}
```

For more details, please refer to testcase.