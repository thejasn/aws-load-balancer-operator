// Code generated by smithy-go-codegen DO NOT EDIT.

package wafregional

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// This is AWS WAF Classic Regional documentation. For more information, see AWS
// WAF Classic
// (https://docs.aws.amazon.com/waf/latest/developerguide/classic-waf-chapter.html)
// in the developer guide. For the latest version of AWS WAF, use the AWS WAFV2 API
// and see the AWS WAF Developer Guide
// (https://docs.aws.amazon.com/waf/latest/developerguide/waf-chapter.html). With
// the latest version, AWS WAF has a single set of endpoints for regional and
// global use. Associates a web ACL with a resource, either an application load
// balancer or Amazon API Gateway stage.
func (c *Client) AssociateWebACL(ctx context.Context, params *AssociateWebACLInput, optFns ...func(*Options)) (*AssociateWebACLOutput, error) {
	if params == nil {
		params = &AssociateWebACLInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "AssociateWebACL", params, optFns, c.addOperationAssociateWebACLMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*AssociateWebACLOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type AssociateWebACLInput struct {

	// The ARN (Amazon Resource Name) of the resource to be protected, either an
	// application load balancer or Amazon API Gateway stage. The ARN should be in one
	// of the following formats:
	//
	// * For an Application Load Balancer:
	// arn:aws:elasticloadbalancing:region:account-id:loadbalancer/app/load-balancer-name/load-balancer-id
	//
	// *
	// For an Amazon API Gateway stage:
	// arn:aws:apigateway:region::/restapis/api-id/stages/stage-name
	//
	// This member is required.
	ResourceArn *string

	// A unique identifier (ID) for the web ACL.
	//
	// This member is required.
	WebACLId *string

	noSmithyDocumentSerde
}

type AssociateWebACLOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationAssociateWebACLMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsAwsjson11_serializeOpAssociateWebACL{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsAwsjson11_deserializeOpAssociateWebACL{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpAssociateWebACLValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opAssociateWebACL(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opAssociateWebACL(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "waf-regional",
		OperationName: "AssociateWebACL",
	}
}
