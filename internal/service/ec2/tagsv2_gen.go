// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ec2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// GetTag fetches an individual ec2 service tag for a resource.
// Returns whether the key value and any errors. A NotFoundError is used to signal that no value was found.
// This function will optimise the handling over listTags, if possible.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func GetTag(ctx context.Context, conn *ec2.Client, identifier, key string, optFns ...func(*ec2.Options)) (*string, error) {
	input := &ec2.DescribeTagsInput{
		Filters: []awstypes.Filter{
			{
				Name:   aws.String("resource-id"),
				Values: []string{identifier},
			},
			{
				Name:   aws.String(names.AttrKey),
				Values: []string{key},
			},
		},
	}

	output, err := conn.DescribeTags(ctx, input, optFns...)

	if err != nil {
		return nil, err
	}

	listTags := keyValueTagsV2(ctx, output.Tags)

	if !listTags.KeyExists(key) {
		return nil, tfresource.NewEmptyResultError(nil)
	}

	return listTags.KeyValue(key), nil
}

// []*SERVICE.Tag handling

// TagsV2 returns ec2 service tags.
func TagsV2(tags tftags.KeyValueTags) []awstypes.Tag {
	result := make([]awstypes.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := awstypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// keyValueTagsV2 creates tftags.KeyValueTags from ec2 service tags.
//
// Accepts the following types:
//   - []awstypes.Tag
//   - []awstypes.TagDescription
func keyValueTagsV2(ctx context.Context, tags any) tftags.KeyValueTags {
	switch tags := tags.(type) {
	case []awstypes.Tag:
		m := make(map[string]*string, len(tags))

		for _, tag := range tags {
			m[aws.ToString(tag.Key)] = tag.Value
		}

		return tftags.New(ctx, m)
	case []awstypes.TagDescription:
		m := make(map[string]*string, len(tags))

		for _, tag := range tags {
			m[aws.ToString(tag.Key)] = tag.Value
		}

		return tftags.New(ctx, m)
	default:
		return tftags.New(ctx, nil)
	}
}

// getTagsInV2 returns ec2 service tags from Context.
// nil is returned if there are no input tags.
func getTagsInV2(ctx context.Context) []awstypes.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := TagsV2(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOutV2 sets ec2 service tags in Context.
func setTagsOutV2(ctx context.Context, tags any) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(keyValueTagsV2(ctx, tags))
	}
}

// updateTagsV2 updates ec2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTagsV2(ctx context.Context, conn *ec2.Client, identifier string, oldTagsMap, newTagsMap any, optFns ...func(*ec2.Options)) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.EC2)
	if len(removedTags) > 0 {
		input := &ec2.DeleteTagsInput{
			Resources: []string{identifier},
			Tags:      TagsV2(removedTags),
		}

		_, err := conn.DeleteTags(ctx, input, optFns...)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.EC2)
	if len(updatedTags) > 0 {
		input := &ec2.CreateTagsInput{
			Resources: []string{identifier},
			Tags:      TagsV2(updatedTags),
		}

		_, err := conn.CreateTags(ctx, input, optFns...)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}
