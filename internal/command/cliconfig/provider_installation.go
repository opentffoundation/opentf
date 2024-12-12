// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cliconfig

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl"
	hclast "github.com/hashicorp/hcl/hcl/ast"
	hcl2 "github.com/hashicorp/hcl/v2"
	hcl2syntax "github.com/hashicorp/hcl/v2/hclsyntax"
	svchost "github.com/hashicorp/terraform-svchost"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"

	"github.com/opentofu/opentofu/internal/addrs"
	"github.com/opentofu/opentofu/internal/getproviders"
	"github.com/opentofu/opentofu/internal/tfdiags"
)

// ProviderInstallation is the structure of the "provider_installation"
// nested block within the CLI configuration.
type ProviderInstallation struct {
	Methods []*ProviderInstallationMethod

	// DevOverrides allows overriding the normal selection process for
	// a particular subset of providers to force using a particular
	// local directory and disregard version numbering altogether.
	// This is here to allow provider developers to conveniently test
	// local builds of their plugins in a development environment, without
	// having to fuss with version constraints, dependency lock files, and
	// so forth.
	//
	// This is _not_ intended for "production" use because it bypasses the
	// usual version selection and checksum verification mechanisms for
	// the providers in question. To make that intent/effect clearer, some
	// OpenTofu commands emit warnings when overrides are present. Local
	// mirror directories are a better way to distribute "released"
	// providers, because they are still subject to version constraints and
	// checksum verification.
	DevOverrides map[addrs.Provider]getproviders.PackageLocalDir
}

// decodeProviderInstallationFromConfig uses the HCL AST API directly to
// decode "provider_installation" blocks from the given file.
//
// This uses the HCL AST directly, rather than HCL's decoder, because the
// intended configuration structure can't be represented using the HCL
// decoder's struct tags. This structure is intended as something that would
// be relatively easier to deal with in HCL 2 once we eventually migrate
// CLI config over to that, and so this function is stricter than HCL 1's
// decoder would be in terms of exactly what configuration shape it is
// expecting.
//
// Note that this function wants the top-level file object which might or
// might not contain provider_installation blocks, not a provider_installation
// block directly itself.
func decodeProviderInstallationFromConfig(hclFile *hclast.File) ([]*ProviderInstallation, tfdiags.Diagnostics) {
	var ret []*ProviderInstallation
	var diags tfdiags.Diagnostics

	root := hclFile.Node.(*hclast.ObjectList)

	// This is a rather odd hybrid: it's a HCL 2-like decode implemented using
	// the HCL 1 AST API. That makes it a bit awkward in places, but it allows
	// us to mimic the strictness of HCL 2 (making a later migration easier)
	// and to support a block structure that the HCL 1 decoder can't represent.
	for _, block := range root.Items {
		if block.Keys[0].Token.Value() != "provider_installation" {
			continue
		}
		// HCL only tracks whether the input was JSON or native syntax inside
		// individual tokens, so we'll use our block type token to decide
		// and assume that the rest of the block must be written in the same
		// syntax, because syntax is a whole-file idea.
		isJSON := block.Keys[0].Token.JSON
		if block.Assign.Line != 0 && !isJSON {
			// Seems to be an attribute rather than a block
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid provider_installation block",
				fmt.Sprintf("The provider_installation block at %s must not be introduced with an equals sign.", block.Pos()),
			))
			continue
		}
		if len(block.Keys) > 1 && !isJSON {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid provider_installation block",
				fmt.Sprintf("The provider_installation block at %s must not have any labels.", block.Pos()),
			))
		}

		pi := &ProviderInstallation{}
		devOverrides := make(map[addrs.Provider]getproviders.PackageLocalDir)

		body, ok := block.Val.(*hclast.ObjectType)
		if !ok {
			// We can't get in here with native HCL syntax because we
			// already checked above that we're using block syntax, but
			// if we're reading JSON then our value could potentially be
			// anything.
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid provider_installation block",
				fmt.Sprintf("The provider_installation block at %s must not be introduced with an equals sign.", block.Pos()),
			))
			continue
		}

		for _, methodBlock := range body.List.Items {
			if methodBlock.Assign.Line != 0 && !isJSON {
				// Seems to be an attribute rather than a block
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Invalid provider_installation method block",
					fmt.Sprintf("The items inside the provider_installation block at %s must all be blocks.", block.Pos()),
				))
				continue
			}
			if len(methodBlock.Keys) > 1 && !isJSON {
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Invalid provider_installation method block",
					fmt.Sprintf("The blocks inside the provider_installation block at %s may not have any labels.", block.Pos()),
				))
			}

			methodBody, ok := methodBlock.Val.(*hclast.ObjectType)
			if !ok {
				// We can't get in here with native HCL syntax because we
				// already checked above that we're using block syntax, but
				// if we're reading JSON then our value could potentially be
				// anything.
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Invalid provider_installation method block",
					fmt.Sprintf("The items inside the provider_installation block at %s must all be blocks.", block.Pos()),
				))
				continue
			}

			methodTypeStr := methodBlock.Keys[0].Token.Value().(string)
			var location ProviderInstallationLocation
			var include, exclude []string
			switch methodTypeStr {
			case "direct":
				type BodyContent struct {
					Include []string `hcl:"include"`
					Exclude []string `hcl:"exclude"`

					// A temporary extra setting available only for experimental builds (checked
					// in the validate step) which opts in to the not-yet-finalized alternative
					// provider registry service "oci-providers.v1", which allows installing
					// providers on a particular hostname directly from OCI distribution
					// repositories.
					OCIRegistryExperiment bool `hcl:"oci_registry_experiment"`
				}
				var bodyContent BodyContent
				err := hcl.DecodeObject(&bodyContent, methodBody)
				if err != nil {
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("Invalid %s block at %s: %s.", methodTypeStr, block.Pos(), err),
					))
					continue
				}
				location = ProviderInstallationDirect
				if bodyContent.OCIRegistryExperiment {
					location = ProviderInstallationDirectWithOCIExperiment
				}
				include = bodyContent.Include
				exclude = bodyContent.Exclude
			case "filesystem_mirror":
				type BodyContent struct {
					Path    string   `hcl:"path"`
					Include []string `hcl:"include"`
					Exclude []string `hcl:"exclude"`
				}
				var bodyContent BodyContent
				err := hcl.DecodeObject(&bodyContent, methodBody)
				if err != nil {
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("Invalid %s block at %s: %s.", methodTypeStr, block.Pos(), err),
					))
					continue
				}
				if bodyContent.Path == "" {
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("Invalid %s block at %s: \"path\" argument is required.", methodTypeStr, block.Pos()),
					))
					continue
				}
				location = ProviderInstallationFilesystemMirror(bodyContent.Path)
				include = bodyContent.Include
				exclude = bodyContent.Exclude
			case "network_mirror":
				type BodyContent struct {
					URL     string   `hcl:"url"`
					Include []string `hcl:"include"`
					Exclude []string `hcl:"exclude"`
				}
				var bodyContent BodyContent
				err := hcl.DecodeObject(&bodyContent, methodBody)
				if err != nil {
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("Invalid %s block at %s: %s.", methodTypeStr, block.Pos(), err),
					))
					continue
				}
				if bodyContent.URL == "" {
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("Invalid %s block at %s: \"url\" argument is required.", methodTypeStr, block.Pos()),
					))
					continue
				}
				location = ProviderInstallationNetworkMirror(bodyContent.URL)
				include = bodyContent.Include
				exclude = bodyContent.Exclude
			case "oci_mirror":
				var moreDiags tfdiags.Diagnostics
				location, include, exclude, moreDiags = decodeProviderInstallationOCIMirrorBlock(methodBody)
				diags = diags.Append(moreDiags)
				if moreDiags.HasErrors() {
					continue
				}
			case "dev_overrides":
				if len(pi.Methods) > 0 {
					// We require dev_overrides to appear first if it's present,
					// because dev_overrides effectively bypass the normal
					// selection process for a particular provider altogether,
					// and so they don't participate in the usual
					// include/exclude arguments and priority ordering.
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("The dev_overrides block at at %s must appear before all other installation methods, because development overrides always have the highest priority.", methodBlock.Pos()),
					))
					continue
				}

				// The content of a dev_overrides block is a mapping from
				// provider source addresses to local filesystem paths. To get
				// our decoding started, we'll use the normal HCL decoder to
				// populate a map of strings and then decode further from
				// that.
				var rawItems map[string]string
				err := hcl.DecodeObject(&rawItems, methodBody)
				if err != nil {
					diags = diags.Append(tfdiags.Sourceless(
						tfdiags.Error,
						"Invalid provider_installation method block",
						fmt.Sprintf("Invalid %s block at %s: %s.", methodTypeStr, block.Pos(), err),
					))
					continue
				}

				for rawAddr, rawPath := range rawItems {
					addr, moreDiags := addrs.ParseProviderSourceString(rawAddr)
					if moreDiags.HasErrors() {
						diags = diags.Append(tfdiags.Sourceless(
							tfdiags.Error,
							"Invalid provider installation dev overrides",
							fmt.Sprintf("The entry %q in %s is not a valid provider source string.\n\n%s", rawAddr, block.Pos(), moreDiags.Err().Error()),
						))
						continue
					}
					dirPath := filepath.Clean(rawPath)
					devOverrides[addr] = getproviders.PackageLocalDir(dirPath)
				}

				continue // We won't add anything to pi.MethodConfigs for this one

			default:
				diags = diags.Append(tfdiags.Sourceless(
					tfdiags.Error,
					"Invalid provider_installation method block",
					fmt.Sprintf("Unknown provider installation method %q at %s.", methodTypeStr, methodBlock.Pos()),
				))
				continue
			}

			pi.Methods = append(pi.Methods, &ProviderInstallationMethod{
				Location: location,
				Include:  include,
				Exclude:  exclude,
			})
		}

		if len(devOverrides) > 0 {
			pi.DevOverrides = devOverrides
		}

		ret = append(ret, pi)
	}

	return ret, diags
}

// ProviderInstallationMethod represents an installation method block inside
// a provider_installation block.
type ProviderInstallationMethod struct {
	Location ProviderInstallationLocation
	Include  []string `hcl:"include"`
	Exclude  []string `hcl:"exclude"`
}

// ProviderInstallationLocation is an interface type representing the
// different installation location types. The concrete implementations of
// this interface are:
//
//   - [ProviderInstallationDirect]:                 install from the provider's origin registry
//   - [ProviderInstallationFilesystemMirror] (dir): install from a local filesystem mirror
//   - [ProviderInstallationNetworkMirror] (host):   install from a network mirror
type ProviderInstallationLocation interface {
	providerInstallationLocation()
}

type providerInstallationDirect [0]byte

func (i providerInstallationDirect) providerInstallationLocation() {}

// ProviderInstallationDirect is a ProviderInstallationSourceLocation
// representing installation from a provider's origin registry.
var ProviderInstallationDirect ProviderInstallationLocation = providerInstallationDirect{}

func (i providerInstallationDirect) GoString() string {
	return "cliconfig.ProviderInstallationDirect"
}

type providerInstallationDirectWithOCIExperiment [0]byte

func (i providerInstallationDirectWithOCIExperiment) providerInstallationLocation() {}

// ProviderInstallationDirectWithOCIExperiment is a temporary variant
// of [ProviderInstallationDirect] which also supports the experimental
// "oci-providers.v1" service as an alternative to "providers.v1".
//
// If this experiment is successful and implemented as a non-experiment
// then this should be removed and its functionality incorporated into
// [ProviderInstallationDirect] directoy instead.
//
//nolint:gochecknoglobals // This part of the system was using global variables long before we had this linter
var ProviderInstallationDirectWithOCIExperiment ProviderInstallationLocation = providerInstallationDirectWithOCIExperiment{}

// ProviderInstallationFilesystemMirror is a ProviderInstallationSourceLocation
// representing installation from a particular local filesystem mirror. The
// string value is the filesystem path to the mirror directory.
type ProviderInstallationFilesystemMirror string

func (i ProviderInstallationFilesystemMirror) providerInstallationLocation() {}

func (i ProviderInstallationFilesystemMirror) GoString() string {
	return fmt.Sprintf("cliconfig.ProviderInstallationFilesystemMirror(%q)", i)
}

// ProviderInstallationNetworkMirror is a ProviderInstallationSourceLocation
// representing installation from a particular local network mirror. The
// string value is the HTTP base URL exactly as written in the configuration,
// without any normalization.
type ProviderInstallationNetworkMirror string

func (i ProviderInstallationNetworkMirror) providerInstallationLocation() {}

func (i ProviderInstallationNetworkMirror) GoString() string {
	return fmt.Sprintf("cliconfig.ProviderInstallationNetworkMirror(%q)", i)
}

// ProviderInstallationOCIMirror is a ProviderInstallationSourceLocation
// representing installation from an OCI registry that is being used in
// a similar way as a "network mirror" but using a non-OpenTofu-specific
// protocol.
type ProviderInstallationOCIMirror struct {
	// RepositoryAddrFunc represents the rule for translating a
	// provider source address into an OCI registry repository address.
	//
	// When loaded from a CLI configuration file, this wraps the evaluation
	// of an HCL template defined in the oci_mirror configuration block.
	// The HCL template expression is intentionally encapsulated here
	// so that the provider installation codepaths won't need to depend
	// on HCL directly to evaluate this.
	RepositoryAddrFunc func(addrs.Provider) (getproviders.OCIRepository, tfdiags.Diagnostics)
}

func (i ProviderInstallationOCIMirror) providerInstallationLocation() {}

func (i ProviderInstallationOCIMirror) GoString() string {
	return "cliconfig.ProviderInstallationOCIMirror(...)"
}

func decodeProviderInstallationOCIMirrorBlock(methodBody *hclast.ObjectType) (ProviderInstallationLocation, []string, []string, tfdiags.Diagnostics) {
	var diags tfdiags.Diagnostics
	type BodyContent struct {
		RepositoryTemplate string   `hcl:"repository_template"`
		Include            []string `hcl:"include"`
		Exclude            []string `hcl:"exclude"`
	}
	var bodyContent BodyContent
	err := hcl.DecodeObject(&bodyContent, methodBody)
	if err != nil {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Invalid provider_installation method block",
			fmt.Sprintf("Invalid oci_mirror block at %s: %s.", methodBody.Pos(), err),
		))
		return nil, nil, nil, diags
	}
	if bodyContent.RepositoryTemplate == "" {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Invalid provider_installation method block",
			fmt.Sprintf("Invalid oci_mirror block at %s: \"repository_template\" argument is required.", methodBody.Pos()),
		))
		return nil, nil, nil, diags
	}
	templateExpr, hclDiags := hcl2syntax.ParseTemplate([]byte(bodyContent.RepositoryTemplate), "<oci_mirror repository_template>", hcl2.InitialPos)
	diags = diags.Append(hclDiags)
	if hclDiags.HasErrors() {
		return nil, nil, nil, diags
	}
	location := ProviderInstallationOCIMirror{
		RepositoryAddrFunc: repositoryAddrFuncForHCLTemplate(templateExpr, methodBody),
	}
	include := bodyContent.Include
	exclude := bodyContent.Exclude

	diags = diags.Append(
		validateOCIMirrorTemplateExpr(templateExpr, include, methodBody),
	)

	return location, include, exclude, diags
}

func repositoryAddrFuncForHCLTemplate(templateExpr hcl2.Expression, methodBody *hclast.ObjectType) func(addrs.Provider) (getproviders.OCIRepository, tfdiags.Diagnostics) {
	pos := methodBody.Pos() // So that our closure won't prevent garbage collection of the whole methodBody

	// FIXME: Unfortunately HCLv1 doesn't retain source filename information as
	// part of its source positions and so our diagnostics in the function
	// below will only refer to the line and column. Since these errors
	// will ultimately be returned by the provider installer rather than the
	// CLI config loader we do at least mention in the messages that they
	// came from the CLI configuration, though if the user has multiple
	// CLI config files they'll need to figure out for themselves which
	// one we're talking about.
	return func(provider addrs.Provider) (getproviders.OCIRepository, tfdiags.Diagnostics) {
		var diags tfdiags.Diagnostics
		evalCtx := &hcl2.EvalContext{
			Variables: map[string]cty.Value{
				"hostname":  cty.StringVal(provider.Hostname.ForDisplay()),
				"namespace": cty.StringVal(provider.Namespace),
				"type":      cty.StringVal(provider.Type),
			},
		}
		v, hclDiags := templateExpr.Value(evalCtx)
		diags = diags.Append(hclDiags)
		if hclDiags.HasErrors() {
			// Since these diagnostics are coming from HCLv2 itself, they will
			// describe the source location as "<oci_mirror repository_template>"
			// rather than an actual file location. This is just another
			// unfortunate consequence of continuing to use legacy HCL
			// for the CLI configuration. :(
			return getproviders.OCIRepository{}, diags
		}

		v, err := convert.Convert(v, cty.String)
		if err != nil {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid oci_mirror repository template",
				fmt.Sprintf("Invalid oci_mirror repository template in CLI configuration at %s: %s.", pos, tfdiags.FormatError(err)),
			))
			return getproviders.OCIRepository{}, diags
		}
		if v.IsNull() {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid oci_mirror repository template",
				fmt.Sprintf("Invalid oci_mirror repository template in CLI configuration at %s: template result must not be null.", pos),
			))
			return getproviders.OCIRepository{}, diags
		}

		// We can assume that v is definitely known because the EvalContext didn't
		// include anything that could produce an unknown value, and HCL promises
		// not to invent its own unknown values if evaluation was successful.
		addrStr := v.AsString()
		slash := strings.IndexByte(addrStr, '/')
		if slash == -1 {
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid oci_mirror repository template",
				fmt.Sprintf("Invalid oci_mirror repository template in CLI configuration at %s: template result does not include an OCI registry hostname.", pos),
			))
			return getproviders.OCIRepository{}, diags
		}
		return getproviders.OCIRepository{
			Hostname: addrStr[:slash],
			Name:     addrStr[slash+1:],
		}, diags
	}
}

func validateOCIMirrorTemplateExpr(templateExpr hcl2.Expression, include []string, methodBody *hclast.ObjectType) tfdiags.Diagnostics {
	var diags tfdiags.Diagnostics
	var templateHasHostname, templateHasNamespace, templateHasType bool
	for _, traversal := range templateExpr.Variables() {
		switch name := traversal.RootName(); name {
		case "hostname":
			templateHasHostname = true
		case "namespace":
			templateHasNamespace = true
		case "type":
			templateHasType = true
		default:
			diags = diags.Append(tfdiags.Sourceless(
				tfdiags.Error,
				"Invalid oci_mirror repository template",
				fmt.Sprintf(
					"Invalid oci_mirror block at %s: the symbol %q is not available for an OCI mirror repository address template. Only \"hostname\", \"namespace\", and \"type\" are available.",
					methodBody.Pos(), name,
				),
			))
			// We continue anyway, because we might be able to collect other errors
			// if the template is invalid in multiple ways.
		}
	}

	// The template must include at least one reference to any source address
	// component that isn't isn't exactly matched by all of the "include" patterns.
	// It's okay to ignore any component that is exactly constrained by the
	// "include" patterns since they'd always have the same value anyway.
	includePatterns, err := getproviders.ParseMultiSourceMatchingPatterns(include)
	if err != nil {
		// Invalid patterns get caught later when we finally assemble the provider
		// sources, so we intentionally don't produce an error here to avoid
		// reporting the same problem twice. Instead, we just skip the
		// template checking altogether by returning early.
		return diags
	}

	hostnames := map[svchost.Hostname]struct{}{}
	namespaces := map[string]struct{}{}
	types := map[string]struct{}{}
	for _, pattern := range includePatterns {
		if pattern.Hostname != svchost.Hostname(getproviders.Wildcard) {
			hostnames[pattern.Hostname] = struct{}{}
		}
		if pattern.Namespace != getproviders.Wildcard {
			namespaces[pattern.Namespace] = struct{}{}
		}
		if pattern.Type != getproviders.Wildcard {
			types[pattern.Type] = struct{}{}
		}
	}
	if len(hostnames) != 1 && !templateHasHostname {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Invalid oci_mirror repository template",
			fmt.Sprintf("Invalid oci_mirror block at %s: template must refer to the \"hostname\" symbol unless the \"include\" argument selects exactly one registry hostname.", methodBody.Pos()),
		))
	}
	if len(namespaces) != 1 && !templateHasNamespace {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Invalid oci_mirror repository template",
			fmt.Sprintf("Invalid oci_mirror block at %s: template must refer to the \"namespace\" symbol unless the \"include\" argument selects exactly one provider namespace.", methodBody.Pos()),
		))
	}
	if len(types) != 1 && !templateHasType {
		diags = diags.Append(tfdiags.Sourceless(
			tfdiags.Error,
			"Invalid oci_mirror repository template",
			fmt.Sprintf("Invalid oci_mirror block at %s: template must refer to the \"type\" symbol unless the \"include\" argument selects exactly one provider.", methodBody.Pos()),
		))
	}

	return diags
}
