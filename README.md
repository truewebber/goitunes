## GoItunes
*this lib can help you to buy ios applications*

### simple start with
```go
itunes, err := NewGOiTunes(
    "your@apple.id",
    "auth_token_if_you_know_it_can_be_empty",
    "buying_cert_requared_for_buy_method",
    "dsid_if_you_know_it_can_be_empty",
    "geo_required",
    "machine_guid_required_for_buy",
    "simple_random_user_agent",
    "machine_name_required_for_buy",
)
```

yeah... maybe not so simple... okay, it's up to me to improve code and docs.

### Parameters  

 - *kbsync* - it's certificate which you can get only after authorize at `gsa.itunes.com` but i didn't realize how yet. 
 - *appleID* - it's your email and login at every apple services
 - *password* - you already know.
 - *geo* - to choose market, appstore is geo oriented, so as your account, if you registered as US, you are in US market with US apps
 - *user agent* - it's in every request to apple, don't know why, the best is win7 itunes, i think
 - *machine name* - required and can't be changed
 - *guid* - your machine id, required and can't be changed
 - *dsid* - some id apple gives you after login
 - *x-token* - some hash apple gives you after login

### Login

We can do simple login without access to payment details, but `x-token` and `DSID` are the same!
It means you don't need to think about them, so they are have short live time, but kbsync and machine info are permanent!  

So you just need to find out 3 params from one request!
