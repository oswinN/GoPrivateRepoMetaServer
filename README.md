# GoPrivateRepoMetaServer
An example implementation for a basic http / https server that provides meta tags for private repositories. The server is written in golang.


# Why

Golang allows the use of private repositories, but requires the used domain to provide meta data information. If you want to set up a private repository using your own domain name, for example repo.mydomain.io/mymodule, in order to use it as a golang import, the golang tools (for example go get) will send a request to repo.mydomain.io/mymodule?go-get=1 and expect a meta data return in the following format.

        <meta name="go-import" content="import-prefix vcs repo-root">

The server has a very basic implementation for returning these meta tags using a single level path variable from the URL to point to a base repository URL with a single variable path at the end that uses the variable from the URL.


# Example Scenario

If your repository is in AWS CodeCommit with the URL:

        https://git-codecommit.us-west-2.amazonaws.com/v1/repos/mygreatgomodule

And assumed your own domain is repo.mydomain.io Then you would want to use an import in go code like:

        import "repo.mydomain.io/mygreatgomodule"

The go tools would make a request like the following:

        https://repo.mydomain.io/mygreatgomodule&go-get=1

The return has to point to your AWS CodeCommit repository like this:

        <meta name="go-import" content="repo.mydomain.io git https://git-codecommit.us-west-2.amazonaws.com/v1/repos/mygreatgomodule">

And this is what the server allows you to do when you run it at port 443 under the IP that your domain repo.mydomain.io points to. It uses the first path level passed in the URL to build the repo URL to return in the meta tag. So it turns

        https://repo.mydomain.io/${variable}&go-get=1

into

        <meta name="go-import" content="repo.mydomain.io git https://git-codecommit.us-west-2.amazonaws.com/v1/repos/${variable}">

You can use the server as is or customize the handler method to adopt your own rules of which repos to return for which import-prefix.


# How to use:

The server is a basic http/https server using bone and default libraries.

It uses a single configuration file (config.json) with content such as:

        {
            "ServerHost" : "repo.mydomain.io",
            "ServicePort" : 80,
            "VCSType" : "git",            
            "RepoBaseURL" : "https://git-codecommit.us-west-2.amazonaws.com/v1/repos/",
            "DebugOutput" : true,
            "CertFile" : "",
            "KeyFile" : ""
        }

Edit the config file with your host name, port and repository target base url. You can also set a certificate and key file name for enabling the https mode.

Setting the DebugOutput flag to true will add additional output on served requests, such as the complete request header & body.

# Customize

If you need a custom logic for how to return the meta tag and repository information you can just edit the **GoPrivateRepoMetaEndpointHandler** method inside **Handler.go**.
