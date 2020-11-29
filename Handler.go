package main

import (
	"net/http"

	"github.com/go-zoo/bone"
)

// TODO : Customize this handler to reflect your own requirements
func (n *GoPrivateRepoMetaEnpointServer) GoPrivateRepoMetaEndpointHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		goget := req.FormValue("go-get")
		if goget != "1" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		id := bone.GetValue(req, "id")
		importPrefix := n.config.ServerHost
		vcs := n.config.VCSType
		repoRoot := n.config.RepoBaseURL
		if id != "" {
			// we have a regular request
			w.Write([]byte("<meta name=\"go-import\" content=\"" + importPrefix + "/" + id + " " + vcs + " " + repoRoot + "/" + id + "\">"))
		} else {
			// this is a root request and we should return the parent meta and one meta tag for each library that is configured in "Modules"
			w.Write([]byte("<meta name=\"go-import\" content=\"" + importPrefix + " " + vcs + " " + repoRoot + "\">"))
			for _, cMod := range n.config.Modules {
				w.Write([]byte("<meta name=\"go-import\" content=\"" + importPrefix + "/" + id + " " + vcs + " " + repoRoot + "/" + cMod + "\">"))
			}
		}
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}
