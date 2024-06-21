// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pinpoint_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	awstypes "github.com/aws/aws-sdk-go-v2/service/pinpoint/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfpinpoint "github.com/hashicorp/terraform-provider-aws/internal/service/pinpoint"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

/**
 Before running this test, the following ENV variable must be set:

 GCM_API_KEY - Google Cloud Messaging Api Key
 GCM_SERVICE_JSON - Google Cloud Messaging Token Json
**/

func TestAccPinpointGCMChannel_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var channel awstypes.GCMChannelResponse
	resourceName := "aws_pinpoint_gcm_channel.test_gcm_channel"

	if os.Getenv("GCM_API_KEY") == "" {
		t.Skipf("GCM_API_KEY env missing, skip test")
	}

	apiKey := os.Getenv("GCM_API_KEY")
	serviceJson := os.Getenv("GCM_SERVICE_JSON")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheckApp(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.PinpointServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckGCMChannelDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccGCMChannelConfigApiKey_basic(apiKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCMChannelExists(ctx, resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtFalse),
				),
			},
			{
				Config: testAccGCMChannelConfigServiceJson_basic(serviceJson),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCMChannelExists(ctx, resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtFalse),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key", "service_json"},
			},
			{
				Config: testAccGCMChannelConfigApiKey_basic(apiKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGCMChannelExists(ctx, resourceName, &channel),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtFalse),
				),
			},
		},
	})
}

func testAccCheckGCMChannelExists(ctx context.Context, n string, channel *awstypes.GCMChannelResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Pinpoint GCM Channel with that application ID exists")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointClient(ctx)

		output, err := tfpinpoint.FindGCMChannelByApplicationId(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*channel = *output

		return nil
	}
}

func testAccCheckGCMChannelDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).PinpointClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_pinpoint_gcm_channel" {
				continue
			}

			_, err := tfpinpoint.FindGCMChannelByApplicationId(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Pinpoint GCM Channel %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccGCMChannelConfigApiKey_basic(apiKey string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test_app" {}

resource "aws_pinpoint_gcm_channel" "test_gcm_channel_api_key" {
  application_id = aws_pinpoint_app.test_app.application_id
  enabled        = "false"
  api_key        = "%s"
}
`, apiKey)
}

func testAccGCMChannelConfigServiceJson_basic(serviceJson string) string {
	return fmt.Sprintf(`
resource "aws_pinpoint_app" "test_app" {}

resource "aws_pinpoint_gcm_channel" "test_gcm_channel_service_json" {
  application_id 								= aws_pinpoint_app.test_app.application_id
  enabled        								= "false"
	default_authentication_method = "TOKEN"
  service_json        					= "%s"
}
`, serviceJson)
}
