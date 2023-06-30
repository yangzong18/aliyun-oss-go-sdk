package oss

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	arnDelimiter = ":"
	arnSections  = 6
	arnPrefix    = "arn:"

	// zero-indexed
	sectionPartition = 1
	sectionService   = 2
	sectionRegion    = 3
	sectionAccountID = 4
	sectionResource  = 5
)

// Arn represents an Amazon Resource Name.
type Arn struct {
	Partition   string
	Service     string
	Region      string
	AccountID   string
	Resource    string
	ArnResource *ArnResource
}

type OSSResource interface {
	getPartition() string
	getRegion() string
	getAccountId() string
	getType() string
	getParentOSSResource() OSSResource
}

type OSSResourceType string

const (
	ACCESS_POINT OSSResourceType = "accesspoint"
)

var (
	apNamePattern    = regexp.MustCompile("^[0-9a-zA-Z-]+$")
	accountIdPattern = regexp.MustCompile("^[0-9]+$")
)

type OSSArnConverter struct{}

// ParseArn parses an Arn Resource Name string into an Arn instance.
func ParseArn(strArn string) (*Arn, error) {
	if strArn == "" {
		return nil, fmt.Errorf("Malformed ARN: empty string")
	}

	parts := strings.Split(strArn, ":")

	if len(parts) < arnSections {
		return nil, fmt.Errorf("Malformed ARN: %s", strArn)
	}

	partition := parts[sectionPartition]
	service := parts[sectionService]
	region := parts[sectionRegion]
	accountID := parts[sectionAccountID]
	resource := parts[sectionResource]

	if len(parts) > arnSections {
		resource = strings.Join(parts[5:], arnDelimiter)
	}

	arnResource := ParseArnResource(resource)

	return &Arn{
		Partition:   partition,
		Service:     service,
		Region:      region,
		AccountID:   accountID,
		Resource:    resource,
		ArnResource: arnResource,
	}, nil
}

func HasArnPrefix(strArn string) bool {
	return strArn != "" && strings.HasPrefix(strArn, arnPrefix)
}
func (arn *Arn) GetArnAsString() string {
	sb := strings.Builder{}
	sb.WriteString("arn:")
	sb.WriteString(arn.Partition)
	sb.WriteString(":")
	sb.WriteString(arn.Service)
	sb.WriteString(":")
	sb.WriteString(arn.Region)
	sb.WriteString(":")
	sb.WriteString(arn.AccountID)
	sb.WriteString(":")
	sb.WriteString(arn.Resource)
	return sb.String()
}

func (arn *Arn) GetResourceAsString() string {
	return arn.Resource
}

// ArnResource represents the resource part of an Amazon Resource Name.
type ArnResource struct {
	ResourceType string
	Resource     string
	Qualifier    string
}

// ParseArnResource parses an Amazon Resource Name resource string into an ArnResource instance.
func ParseArnResource(strResource string) *ArnResource {
	var resourceTypeBoundary int
	var resourceIdBoundary int
	for i, ch := range strResource {
		if ch == ':' || ch == '/' {
			resourceTypeBoundary = i
			break
		}
	}
	if resourceTypeBoundary != 0 {
		for i := len(strResource) - 1; i > resourceTypeBoundary; i-- {
			ch := rune(strResource[i])

			if ch == ':' {
				resourceIdBoundary = i
				break
			}
		}
	}
	resourceType := ""
	resource := ""
	qualifier := ""
	if resourceTypeBoundary == 0 {
		// 'resource-id'
		resource = strResource
	} else if resourceIdBoundary == 0 {
		// 'resource-type:resource-id'
		resourceType = strResource[:resourceTypeBoundary]
		resource = strResource[resourceTypeBoundary+1:]
	} else {
		// 'resource-type:resource-id:qualifier'
		resourceType = strResource[:resourceTypeBoundary]
		resource = strResource[resourceTypeBoundary+1 : resourceIdBoundary]
		qualifier = strResource[resourceIdBoundary+1:]
	}
	return &ArnResource{
		ResourceType: resourceType,
		Resource:     resource,
		Qualifier:    qualifier,
	}
}

func (arnRe *ArnResource) GetArnResourceAsString() string {
	sb := strings.Builder{}
	if arnRe.ResourceType != "" {
		sb.WriteString(arnRe.ResourceType)
		sb.WriteString(":")
	}
	sb.WriteString(arnRe.Resource)
	if arnRe.Qualifier != "" {
		sb.WriteString(":")
		sb.WriteString(arnRe.Qualifier)
	}
	return sb.String()
}

type OSSAccessPointResource struct {
	Partition         string
	Region            string
	AccountId         string
	AccessPointName   string
	parentOSSResource OSSResource
}

func (r *OSSAccessPointResource) GetParentOSSResource() OSSResource {
	panic("implement me")
}

type OSSAccessPointResourceBuilder struct {
	partition         string
	region            string
	accountId         string
	accessPointName   string
	parentOSSResource OSSResource
}

func NewOSSAccessPointResource(builder *OSSAccessPointResourceBuilder) *OSSAccessPointResource {
	var r *OSSAccessPointResource
	var accessPointName, partition, region, accountId string
	if builder.accessPointName != "" {
		accessPointName = builder.accessPointName
	}
	if builder.parentOSSResource == nil {
		if builder.partition != "" {
			partition = builder.partition
		}
		if builder.region != "" {
			region = builder.region
		}
		if builder.accountId != "" {
			accountId = builder.accountId
		}
	} else {
		parentOSSResource := validateParentOSSResource(builder.parentOSSResource)
		if builder.partition != "" {
			partition = parentOSSResource.getPartition()
		}
		if builder.region != "" {
			region = parentOSSResource.getRegion()
		}
		if builder.accountId != "" {
			accountId = parentOSSResource.getAccountId()
		}
	}
	r = &OSSAccessPointResource{
		AccessPointName: accessPointName,
		Partition:       partition,
		Region:          region,
		AccountId:       accountId,
	}
	return r
}

func validateParentOSSResource(parentOSSResource OSSResource) OSSResource {
	// TODO
	return nil
}

func (r *OSSAccessPointResource) getType() string {
	return "accesspoint"
}

func (r *OSSAccessPointResource) getPartition() string {
	return r.Partition
}

func (r *OSSAccessPointResource) getRegion() string {
	return r.Region
}

func (r *OSSAccessPointResource) getAccountId() string {
	return r.AccountId
}

func (r *OSSAccessPointResource) getAccessPointName() string {
	return r.AccessPointName
}

func (r *OSSAccessPointResource) getParentOSSResource() OSSResource {
	return r.parentOSSResource
}

func NewOSSResourceType(value string) (OSSResourceType, error) {
	switch value {
	case string(ACCESS_POINT):
		return ACCESS_POINT, nil
	default:
		return "", fmt.Errorf("invalid value for Oss Resource Type: %s", value)
	}
}

func NewOSSArnConverter() *OSSArnConverter {
	return &OSSArnConverter{}
}
func (c *OSSArnConverter) ConvertArn(strArn string) (OSSResource, error) {
	arn, err := ParseArn(strArn)
	if err != nil {
		return nil, err
	}

	resourceType := strings.ToLower(arn.ArnResource.ResourceType)
	switch resourceType {
	case string(ACCESS_POINT):
		return c.ParseOSSAccessPointArn(*arn)
	default:
		return nil, fmt.Errorf("unknown ARN type '%s'", arn.ArnResource.ResourceType)
	}
}

func (c *OSSArnConverter) ParseOSSAccessPointArn(arn Arn) (*OSSAccessPointResource, error) {
	//check ap name
	if !apNamePattern.MatchString(arn.ArnResource.Resource) {
		return nil, fmt.Errorf("Access Point Name arn is invalid")
	}

	//check account id
	if !accountIdPattern.MatchString(arn.AccountID) {
		return nil, fmt.Errorf("Account Id in arn is invalid")
	}

	return &OSSAccessPointResource{
		Partition:         arn.Partition,
		Region:            arn.Region,
		AccountId:         arn.AccountID,
		AccessPointName:   arn.ArnResource.Resource,
		parentOSSResource: nil,
	}, nil
}
