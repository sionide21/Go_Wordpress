package wordpress

import (
	"xmlrpc"
	"os"
)

func (b *Blog) method(s string) xmlrpc.RemoteMethod {
	return xmlrpc.RemoteMethod{
		Method: s,
		Endpoint: b.XMLRPC,
		BaseParams: xmlrpc.Params(b.Id, b.Session.Username, b.Session.Password),
	}
}

func (b *Blog) GetTags() ([]Tag, os.Error) {
	tags, err := b.method("wp.getTags").CallArray()
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

func (b *Blog) GetCommentCount(p string) (CommentCount, os.Error) {
	count, err := b.method("wp.getCommentCount").CallStruct(p)
	if err != nil {
		return CommentCount{}, err
	}
	return CommentCount{
		Approved: getInt(count, "approved"), // SRSLY, one as string the rest as int?
		AwaitingModeration: count.GetInt("awaiting_moderation"),
		Spam: count.GetInt("spam"),
		Total: count.GetInt("total_comments"),
	},
		nil
}

func (b *Blog) GetPageTemplates() ([]PageTemplate, os.Error) {
	templates, err := b.method("wp.getPageTemplates").CallStruct()
	if err != nil {
		return nil, err
	}
	list := make([]PageTemplate, len(templates))
	i := 0
	for name := range templates {
		list[i] = PageTemplate{
			Name: name,
			Description: templates.GetString(name),
		}
		i++
	}
	return list, nil
}

func parseOptions(opts xmlrpc.StructValue, err os.Error) (map[string]Option, os.Error) {
	if err != nil {
		return nil, err
	}
	m := make(map[string]Option, len(opts))
	for name := range opts {
		opt := opts.GetStruct(name)
		m[name] = Option{
			Name: name,
			ReadOnly: opt.GetBoolean("readonly"),
			Desc: opt.GetString("desc"),
			Value: opt.GetRaw("value").String(), // Wordpress has a very bad xmlrpc impl. This is a necessary evil
		}
	}
	return m, nil
}

func (b *Blog) GetOptions(options ...) (map[string]Option, os.Error) {
	// By calling params manually we get an array of options as one param instaed of each option as a param
	return parseOptions(b.method("wp.getOptions").CallStruct(xmlrpc.Params(options)))
}

func (b *Blog) SetOptions(options map[string]string) (map[string]Option, os.Error) {
	return parseOptions(b.method("wp.setOptions").CallStruct(options))
}

func (b *Blog) SetOption(opt, val string) (map[string]Option, os.Error) {
	return b.SetOptions(map[string]string{opt: val})
}
