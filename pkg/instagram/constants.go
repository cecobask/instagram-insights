package instagram

const (
	OutputJson  Output = "json"
	OutputNone  Output = "none"
	OutputTable Output = "table"
	OutputYaml  Output = "yaml"
)

const (
	FlagOutput                 = "output"
	GoogleDriveHost            = "drive.google.com"
	GoogleDriveParsedUrlFormat = "https://drive.google.com/u/0/uc?id=%s&export=download&confirm=t"
	PathData                   = "instagram_data"
	PathDataArchive            = PathData + ".zip"
	PathDocs                   = "docs"
	PathFollowers              = PathData + "/followers_and_following/followers_*.json"
	PathFollowing              = PathData + "/followers_and_following/following.json"
	TableHeaderProfileUrl      = "PROFILE URL"
	TableHeaderTimestamp       = "TIMESTAMP"
	TableHeaderUsername        = "USERNAME"
)
