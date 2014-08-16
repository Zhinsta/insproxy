package main

var allowdHosts = map[string]bool{
	"zhinsta.com:8080": true,
	"amazonaws.com":    true,
	"fbcdn.net":        true,
	"instagram.com":    true,
	"cdninstagram.com": true,
}

func isHostAllowed(host string) bool {
	_, ok := allowdHosts[host]
	return ok
}
