package domain

type S3Event struct {
	Records []S3EventRecord `json:"Records"`
}

type S3EventRecord struct {
	S3 S3Entity `json:"s3"`
}

type S3UserIdentity struct {
	PrincipalID string `json:"principalId"`
}

type S3Entity struct {
	SchemaVersion   string   `json:"s3SchemaVersion"`
	ConfigurationID string   `json:"configurationId"`
	Bucket          S3Bucket `json:"bucket"`
	Object          S3Object `json:"object"`
}

type S3Bucket struct {
	Name          string         `json:"name"`
	OwnerIdentity S3UserIdentity `json:"ownerIdentity"`
	Arn           string         `json:"arn"` //nolint: stylecheck
}

type S3Object struct {
	Key           string `json:"key"`
	Size          int64  `json:"size,omitempty"`
	URLDecodedKey string `json:"urlDecodedKey"`
	VersionID     string `json:"versionId"`
	ETag          string `json:"eTag"`
	Sequencer     string `json:"sequencer"`
}
