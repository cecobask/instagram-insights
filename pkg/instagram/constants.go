package instagram

const (
	FieldTimestamp = "timestamp"
	FieldUsername  = "username"
	Unlimited      = 0
	OrderAsc       = "asc"
	OrderDesc      = "desc"
	OutputJson     = "json"
	OutputNone     = "none"
	OutputTable    = "table"
	OutputYaml     = "yaml"
)

const (
	FlagLimit                  = "limit"
	FlagOrder                  = "order"
	FlagOutput                 = "output"
	FlagSortBy                 = "sort-by"
	GoogleDriveHost            = "drive.google.com"
	GoogleDriveParsedUrlFormat = "https://drive.google.com/u/0/uc?id=%s&export=download&confirm=t"
	PathData                   = "instagram_data"
	PathDataArchive            = PathData + ".zip"
	PathDocs                   = "docs"
	PathFollowers              = PathData + "/connections/followers_and_following/followers_*.json"
	PathFollowing              = PathData + "/connections/followers_and_following/following.json"
	TableHeaderProfileUrl      = "PROFILE URL"
	TableHeaderTimestamp       = "TIMESTAMP"
	TableHeaderUsername        = "USERNAME"
)
