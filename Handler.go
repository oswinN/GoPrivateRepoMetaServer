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
		// we have a regular request
		w.Write([]byte("<meta name=\"go-import\" content=\"" + importPrefix + " " + vcs + " " + repoRoot + id + "\">"))
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}
