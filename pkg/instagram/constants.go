package instagram

const (
	GoogleDriveHost            = "drive.google.com"
	GoogleDriveParsedUrlFormat = "https://drive.google.com/u/0/uc?id=%s&export=download&confirm=t"
	OutputNone                 = "none"
	OutputTable                = "table"
	PathData                   = "instagram_data"
	PathDataArchive            = PathData + ".zip"
	PathFollowers              = PathData + "/followers_and_following/followers_*.json"
	PathFollowing              = PathData + "/followers_and_following/following.json"
	TableHeaderProfileUrl      = "PROFILE URL"
	TableHeaderUsername        = "USERNAME"
)
