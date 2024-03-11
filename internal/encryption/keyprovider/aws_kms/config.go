package aws_kms

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
	baselogging "github.com/hashicorp/aws-sdk-go-base/v2/logging"
	"github.com/opentofu/opentofu/internal/encryption/keyprovider"
	"github.com/opentofu/opentofu/internal/httpclient"
	"github.com/opentofu/opentofu/internal/logging"
	"github.com/opentofu/opentofu/version"
)

type Config struct {
	// KeyProvider Config
	KMSKeyID string `hcl:"kms_key_id"`
	KeySpec  string `hcl:"key_spec"`

	// Mirrored S3 Backend Config, mirror any changes
	AccessKey                      string                     `hcl:"access_key,optional"`
	Endpoints                      []ConfigEndpoints          `hcl:"endpoints,block"`
	MaxRetries                     int                        `hcl:"max_retries,optional"`
	Profile                        string                     `hcl:"profile,optional"`
	Region                         string                     `hcl:"region,optional"`
	SecretKey                      string                     `hcl:"secret_key,optional"`
	SkipCredsValidation            bool                       `hcl:"skip_credentials_validation,optional"`
	SkipRequestingAccountId        bool                       `hcl:"skip_requesting_account_id,optional"`
	STSRegion                      string                     `hcl:"sts_region,optional"`
	Token                          string                     `hcl:"token,optional"`
	HTTPProxy                      *string                    `hcl:"http_proxy,optional"`
	HTTPSProxy                     *string                    `hcl:"https_proxy,optional"`
	NoProxy                        string                     `hcl:"no_proxy,optional"`
	Insecure                       bool                       `hcl:"insecure,optional"`
	UseDualStackEndpoint           bool                       `hcl:"use_dualstack_endpoint,optional"`
	UseFIPSEndpoint                bool                       `hcl:"use_fips_endpoint,optional"`
	CustomCABundle                 string                     `hcl:"custom_ca_bundle,optional"`
	EC2MetadataServiceEndpoint     string                     `hcl:"ec2_metadata_service_endpoint,optional"`
	EC2MetadataServiceEndpointMode string                     `hcl:"ec2_metadata_service_endpoint_mode,optional"`
	SkipMetadataAPICheck           *bool                      `hcl:"skip_metadata_api_check,optional"`
	SharedCredentialsFiles         []string                   `hcl:"shared_credentials_files,optional"`
	SharedConfigFiles              []string                   `hcl:"shared_config_files,optional"`
	AssumeRole                     *AssumeRole                `hcl:"assume_role,optional"`
	AssumeRoleWithWebIdentity      *AssumeRoleWithWebIdentity `hcl:"assume_role_with_web_identity,optional"`
	AllowedAccountIds              []string                   `hcl:"allowed_account_ids,optional"`
	ForbiddenAccountIds            []string                   `hcl:"forbidden_account_ids,optional"`
	RetryMode                      string                     `hcl:"retry_mode,optional"`
}

func stringAttrEnvFallback(val string, env string) string {
	if val != "" {
		return val
	}
	return os.Getenv(env)
}

func stringArrayAttrEnvFallback(val []string, env string) []string {
	if len(val) != 0 {
		return val
	}
	envVal := os.Getenv(env)
	if envVal != "" {
		return []string{envVal}
	}
	return nil
}

func (c Config) asAWSBase() (*awsbase.Config, error) {
	// Get endpoints to use
	endpoints, err := c.getEndpoints()
	if err != nil {
		return nil, err
	}

	// Get assume role
	assumeRole, err := c.AssumeRole.asAWSBase()
	if err != nil {
		return nil, err
	}

	// Get assume role with web identity
	assumeRoleWithWebIdentity, err := c.AssumeRoleWithWebIdentity.asAWSBase()
	if err != nil {
		return nil, err
	}

	// Validate region
	if c.Region == "" && os.Getenv("AWS_REGION") == "" && os.Getenv("AWS_DEFAULT_REGION") == "" {
		return nil, fmt.Errorf(`the "region" attribute or the "AWS_REGION" or "AWS_DEFAULT_REGION" environment variables must be set.`)
	}

	// Retry Mode
	if c.MaxRetries == 0 {
		c.MaxRetries = 5
	}
	var retryMode aws.RetryMode
	if len(c.RetryMode) != 0 {
		retryMode, err = aws.ParseRetryMode(c.RetryMode)
		if err != nil {
			return nil, fmt.Errorf("%w: expected %q or %q", err, aws.RetryModeStandard, aws.RetryModeAdaptive)
		}
	}

	// IDMS handling
	imdsEnabled := imds.ClientDefaultEnableState
	if c.SkipMetadataAPICheck != nil {
		if *c.SkipMetadataAPICheck {
			imdsEnabled = imds.ClientEnabled
		} else {
			imdsEnabled = imds.ClientDisabled
		}
	}

	// Validate account_ids
	if len(c.AllowedAccountIds) != 0 && len(c.ForbiddenAccountIds) != 0 {
		return nil, fmt.Errorf("conflicting config attributes: only allowed_account_ids or forbidden_account_ids can be specified, not both")
	}

	return &awsbase.Config{
		AccessKey:               c.AccessKey,
		CallerDocumentationURL:  "https://opentofu.org/docs/language/settings/backends/s3", // TODO
		CallerName:              "KMS Key Provider",
		IamEndpoint:             stringAttrEnvFallback(endpoints.IAM, "AWS_ENDPOINT_URL_IAM"),
		MaxRetries:              c.MaxRetries,
		RetryMode:               retryMode,
		Profile:                 c.Profile,
		Region:                  c.Region,
		SecretKey:               c.SecretKey,
		SkipCredsValidation:     c.SkipCredsValidation,
		SkipRequestingAccountId: c.SkipRequestingAccountId,
		StsEndpoint:             stringAttrEnvFallback(endpoints.STS, "AWS_ENDPOINT_URL_STS"),
		StsRegion:               c.STSRegion,
		Token:                   c.Token,

		// Note: we don't need to read env variables explicitly because they are read implicitly by aws-sdk-base-go:
		// see: https://github.com/hashicorp/aws-sdk-go-base/blob/v2.0.0-beta.41/internal/config/config.go#L133
		// which relies on: https://cs.opensource.google/go/x/net/+/refs/tags/v0.18.0:http/httpproxy/proxy.go;l=89-96
		HTTPProxy:            c.HTTPProxy,
		HTTPSProxy:           c.HTTPSProxy,
		NoProxy:              c.NoProxy,
		Insecure:             c.Insecure,
		UseDualStackEndpoint: c.UseDualStackEndpoint,
		UseFIPSEndpoint:      c.UseFIPSEndpoint,
		UserAgent: awsbase.UserAgentProducts{
			{Name: "APN", Version: "1.0"},
			{Name: httpclient.DefaultApplicationName, Version: version.String()},
		},
		CustomCABundle: stringAttrEnvFallback(c.CustomCABundle, "AWS_CA_BUNDLE"),

		EC2MetadataServiceEnableState:  imdsEnabled,
		EC2MetadataServiceEndpoint:     stringAttrEnvFallback(c.EC2MetadataServiceEndpoint, "AWS_EC2_METADATA_SERVICE_ENDPOINT"),
		EC2MetadataServiceEndpointMode: stringAttrEnvFallback(c.EC2MetadataServiceEndpointMode, "AWS_EC2_METADATA_SERVICE_ENDPOINT_MODE"),

		SharedCredentialsFiles:    stringArrayAttrEnvFallback(c.SharedCredentialsFiles, "AWS_SHARED_CREDENTIALS_FILE"),
		SharedConfigFiles:         stringArrayAttrEnvFallback(c.SharedConfigFiles, "AWS_SHARED_CONFIG_FILE"),
		AssumeRole:                assumeRole,
		AssumeRoleWithWebIdentity: assumeRoleWithWebIdentity,
		AllowedAccountIds:         c.AllowedAccountIds,
		ForbiddenAccountIds:       c.ForbiddenAccountIds,
	}, nil
}

func (c Config) Build() (keyprovider.KeyProvider, keyprovider.KeyMeta, error) {
	cfg, err := c.asAWSBase()
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	ctx, baselog := attachLoggerToContext(ctx)
	cfg.Logger = baselog

	_, awsConfig, awsDiags := awsbase.GetAwsConfig(ctx, cfg)

	if awsDiags.HasError() {
		out := "errors were encountered in aws kms configuration"
		for _, diag := range awsDiags.Errors() {
			out += "\n" + diag.Summary() + " : " + diag.Detail()
		}

		return nil, nil, fmt.Errorf(out)
	}

	return &keyProvider{
		Config: c,
		svc:    kms.NewFromConfig(awsConfig),
		ctx:    ctx,
	}, new(keyMeta), nil
}

// Mirrored from s3 backend config
func attachLoggerToContext(ctx context.Context) (context.Context, baselogging.HcLogger) {
	ctx, baselog := baselogging.NewHcLogger(ctx, logging.HCLogger().Named("backend-s3"))
	ctx = baselogging.RegisterLogger(ctx, baselog)
	return ctx, baselog
}
