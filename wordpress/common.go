package wordpress

import (
	"xmlrpc"
	"os"
)

type session struct {
	Username string
	Password string
}

type Blog struct {
	Id      int
	Session session
	XMLRPC  string
	IsAdmin bool
	URL     string
	Name    string
}

type Tag struct {
	Id       int
	Name     string
	Count    int
	Slug     string
	HTML_URL string
	RSS_URL  string
}

type CommentCount struct {
	Approved           int
	AwaitingModeration int
	Spam               int
	Total              int
}

type PageTemplate struct {
	Name        string
	Description string
}

type Option struct {
	Name     string
	Desc     string
	ReadOnly bool
	Value    string
}

func GetUsersBlogs(url, username, password string) ([]Blog, os.Error) {
	getBlogs := xmlrpc.RemoteMethod{
		Endpoint: url,
		Method: "wp.getUsersBlogs",
	}
	res, err := getBlogs.CallArray(username, password)
	if err != nil {
		return nil, err
	}
	blogs := make([]Blog, len(res))
	for i, unCast := range res {
		v := unCast.(xmlrpc.StructValue)
		b := Blog{}
		b.Id = getInt(v, "blogid") // WTF wordpress, why is this a string?
		b.Session = session{username, password}
		b.XMLRPC = v.GetString("xmlrpc")
		b.IsAdmin = v.GetBoolean("isAdmin")
		b.URL = v.GetString("url")
		b.Name = v.GetString("blogName")
		blogs[i] = b
	}
	return blogs, nil
}
