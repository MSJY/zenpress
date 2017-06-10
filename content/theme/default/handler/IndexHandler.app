IndexHandler = fn(self) {
	//self.AddActionHook("makross_stop", IndexHandle)
	self.SetStore(map[string]var{
			"title": "#默认主题首页#",
			"oh":    "IndexHandler in default",
	})
	//self.DoActionHook("IndexHandler")
	self.AddFilterHook("index_template", indexFixTpl)
	//self.Stop()
	return self.Render("index")
}

IndexHandle = fn() {
	str = "<IndexHandle are Action in makross stop!!!!!!>"
	println(str)
	return str
}

indexFixTpl = fn(b) {
	return []byte(fmt.Sprintf("@### %s ###@", b))
}