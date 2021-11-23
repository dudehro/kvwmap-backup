package structs

type HTMLTemplateData struct {
        Backup  *GdiBackup
        Vars    map[string]string
}

type Request struct {
	Success		bool
	Errors		[]string
	Messages	[]string
	Payload		interface{}
}
