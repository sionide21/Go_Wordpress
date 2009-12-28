package wordpress

import (
	"xmlrpc"
	"os"
)

func (b *Blog) GetTags() ([]Tag, os.Error) {
	getTags := xmlrpc.RemoteMethod{
		Method: "wp.getTags",
		Endpoint: b.XMLRPC,
	}
	tags, err := getTags.CallArray(b.Id, b.Session.Username, b.Session.Password)
	if err != nil {
		return nil, err
	}
	list := make([]Tag, len(tags))
	for i, unCastV := range tags {
		v := unCastV.(xmlrpc.StructValue)
		t := Tag{}
		t.Id = getInt(v, "tag_id")
		t.Name = v.GetString("name")
		t.Count = getInt(v, "count")
		t.Slug = v.GetString("slug")
		t.HTML_URL = v.GetString("html_url")
		t.RSS_URL = v.GetString("rss_url")
		list[i] = t
	}
	return list, nil
}
