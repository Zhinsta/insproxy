package main

var allowdHosts = map[string]bool{
	"distillery.s3.amazonaws.com":        true,
	"distilleryimage0.s3.amazonaws.com":  true,
	"distilleryimage1.s3.amazonaws.com":  true,
	"distilleryimage2.s3.amazonaws.com":  true,
	"distilleryimage3.s3.amazonaws.com":  true,
	"distilleryimage4.s3.amazonaws.com":  true,
	"distilleryimage5.s3.amazonaws.com":  true,
	"distilleryimage6.s3.amazonaws.com":  true,
	"distilleryimage7.s3.amazonaws.com":  true,
	"distilleryimage8.s3.amazonaws.com":  true,
	"distilleryimage9.s3.amazonaws.com":  true,
	"distilleryimage10.s3.amazonaws.com": true,
	"distilleryimage11.s3.amazonaws.com": true,
	"images.ak.instagram.com":            true,
	"origincache-ash.fbcdn.net":          true,
	"origincache-frc.fbcdn.net":          true,
	"origincache-prn.fbcdn.net":          true,
	"photos-a.ak.instagram.com":          true,
	"photos-b.ak.instagram.com":          true,
	"photos-c.ak.instagram.com":          true,
	"photos-d.ak.instagram.com":          true,
	"photos-e.ak.instagram.com":          true,
	"photos-f.ak.instagram.com":          true,
	"photos-g.ak.instagram.com":          true,
	"photos-h.ak.instagram.com":          true,
	"scontent-a.cdninstagram.com":        true,
	"scontent-b.cdninstagram.com":        true,
	"scontent-c.cdninstagram.com":        true,
	"scontent-d.cdninstagram.com":        true,
	"scontent-e.cdninstagram.com":        true,
	"scontent-f.cdninstagram.com":        true,
	"zhinsta.com:8080":                   true, // for debug
}

func isHostAllowed(host string) bool {
	_, ok := allowdHosts[host]
	return ok
}
