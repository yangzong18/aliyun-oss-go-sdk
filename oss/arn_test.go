package oss

import (
	. "gopkg.in/check.v1"
	"strings"
)

type ArnSuite struct{}

var _ = Suite(&ArnSuite{})

func (s *ArnSuite) TestArnBasic(c *C) {
	strArn := "arn:acs:oss:cn-hangzhou:12345:bucket/object"
	arn, err := ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "cn-hangzhou")
	c.Assert(arn.AccountID, Equals, "12345")

	c.Assert(arn.GetResourceAsString(), Equals, "bucket/object")
	c.Assert(arn.ArnResource.ResourceType, Equals, "bucket")
	c.Assert(arn.ArnResource.Resource, Equals, "object")
	c.Assert(arn.ArnResource.Qualifier, Equals, "")

	strArn = "arn:acs:oss:cn-hangzhou:12345:bucket/object:id"
	arn, err = ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "cn-hangzhou")
	c.Assert(arn.AccountID, Equals, "12345")
	c.Assert(arn.GetResourceAsString(), Equals, "bucket/object:id")
	c.Assert(arn.ArnResource.ResourceType, Equals, "bucket")
	c.Assert(arn.ArnResource.Resource, Equals, "object")
	c.Assert(arn.ArnResource.Qualifier, Equals, "id")

	strArn = "arn:acs:oss:cn-hangzhou:12345:id"
	arn, err = ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "cn-hangzhou")
	c.Assert(arn.AccountID, Equals, "12345")
	c.Assert(arn.GetResourceAsString(), Equals, "id")
	c.Assert(arn.ArnResource.ResourceType, Equals, "")
	c.Assert(arn.ArnResource.Resource, Equals, "id")
	c.Assert(arn.ArnResource.Qualifier, Equals, "")

	strArn = "arn:acs:oss:cn-hangzhou:12345:bucket:mybucket"
	arn, err = ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "cn-hangzhou")
	c.Assert(arn.AccountID, Equals, "12345")
	c.Assert(arn.GetResourceAsString(), Equals, "bucket:mybucket")
	c.Log(arn.ArnResource.ResourceType)
	c.Assert(arn.ArnResource.ResourceType, Equals, "bucket")
	c.Assert(arn.ArnResource.Resource, Equals, "mybucket")
	c.Assert(arn.ArnResource.Qualifier, Equals, "")

	strArn = "arn:acs:oss:cn-hangzhou:12345:bucket:mybucket:id"
	arn, err = ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "cn-hangzhou")
	c.Assert(arn.AccountID, Equals, "12345")
	c.Assert(arn.GetResourceAsString(), Equals, "bucket:mybucket:id")
	c.Assert(arn.ArnResource.ResourceType, Equals, "bucket")
	c.Assert(arn.ArnResource.Resource, Equals, "mybucket")
	c.Assert(arn.ArnResource.Qualifier, Equals, "id")

	strArn = "arn:acs:oss:::bucket"
	arn, err = ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "")
	c.Assert(arn.AccountID, Equals, "")
	c.Assert(arn.GetResourceAsString(), Equals, "bucket")
	c.Assert(arn.ArnResource.ResourceType, Equals, "")
	c.Assert(arn.ArnResource.Resource, Equals, "bucket")
	c.Assert(arn.ArnResource.Qualifier, Equals, "")
}

func (s *ArnSuite) TestArnWithSpecialChar(c *C) {
	strArn := "arn:acs#?:oss:cn-hangzhou:12345:bucket#@/object"
	arn, err := ParseArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(arn.Partition, Equals, "acs#?")
	c.Assert(arn.Service, Equals, "oss")
	c.Assert(arn.Region, Equals, "cn-hangzhou")
	c.Assert(arn.AccountID, Equals, "12345")
	c.Assert(arn.GetResourceAsString(), Equals, "bucket#@/object")
	c.Assert(arn.ArnResource.ResourceType, Equals, "bucket#@")
	c.Assert(arn.ArnResource.Resource, Equals, "object")
	c.Assert(arn.ArnResource.Qualifier, Equals, "")
}

func (s *ArnSuite) testArnInvalid(c *C) {
	_, err := ParseArn("")
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "Malformed ARN: empty string"), Equals, true)

	_, err = ParseArn("not arn")
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "Malformed ARN"), Equals, true)

	_, err = ParseArn("arn:::::")
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "Malformed ARN"), Equals, true)

	_, err = ParseArn("arn:acs::::")
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "Malformed ARN"), Equals, true)

	_, err = ParseArn("arn:acs:oss:::")
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "Malformed ARN"), Equals, true)
}

func (s *ArnSuite) testOssResourceTypeBasic(c *C) {
	ossResourceType, err := NewOSSResourceType("accesspoint")
	c.Assert(err, IsNil)
	c.Assert(ossResourceType, Equals, ACCESS_POINT)

	ossResourceType, err = NewOSSResourceType("abc")
	c.Assert(err, NotNil)
	c.Assert(strings.Contains(err.Error(), "invalid value for OSSResourceType"), Equals, true)
}

func (s *ArnSuite) TestOssArnConverter(c *C) {
	strArn := "arn:acs:oss:cn-hangzhou:12345:accesspoint/ap-test"
	res, err := NewOSSArnConverter().ConvertArn(strArn)
	testLogger.Print(res)
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
	apRes, ok := res.(*OSSAccessPointResource)
	c.Assert(ok, Equals, true)
	c.Assert(apRes.Partition, Equals, "acs")
	c.Assert(apRes.Region, Equals, "cn-hangzhou")
	c.Assert(apRes.AccountId, Equals, "12345")
	c.Assert(apRes.AccessPointName, Equals, "ap-test")
	c.Assert(apRes.parentOSSResource, IsNil)

	// AP name empty
	strArn = "arn:acs:oss:cn-hangzhou:12345:accesspoint/"
	_, err = NewOSSArnConverter().ConvertArn(strArn)
	c.Assert(err, NotNil)

	// No AP name
	strArn = "arn:acs:oss:cn-hangzhou:12345:accesspoint"
	_, err = NewOSSArnConverter().ConvertArn(strArn)
	c.Assert(err, NotNil)
	c.Assert(err.Error(), Equals, "Unknown ARN type")

	// Unsupported ARN
	strArn = "arn:acs:oss:cn-hangzhou:12345:cloudbox/cloudbox-id"
	_, err = NewOSSArnConverter().ConvertArn(strArn)
	c.Assert(err, IsNil)
	c.Assert(err.Error(), Equals, "Unknown ARN type")
}

func (s *ArnSuite) TestInvalidOssAccessPointArnArn(c *C) {
	strArn := "arn:acs:oss:cn-hangzhou:12345:accesspoint/adfd#-abc"
	_, err := NewOSSArnConverter().ConvertArn(strArn)
	c.Assert(err, NotNil)

	strArn = "arn:acs:oss:cn-hangzhou:12#ad345:accesspoint/ap-name"
	_, err = NewOSSArnConverter().ConvertArn(strArn)
	c.Assert(err, NotNil)
}
