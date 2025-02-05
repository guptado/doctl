package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/digitalocean/doctl"
	"github.com/digitalocean/doctl/commands/displayers"
	"github.com/digitalocean/doctl/do"
	"github.com/digitalocean/godo"
)

// Network creates the partner commands
func Network() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:     "network",
			Short:   "Display commands that manage Network products",
			Long:    `The commands under ` + "`" + `doctl network` + "`" + ` are for managing Network products`,
			GroupID: manageResourcesGroup,
		},
	}

	cmd.PersistentFlags().String(doctl.ArgInterconnectAttachmentType, "partner", "Specify interconnect attachment type (e.g., partner)")
	viper.BindPFlag(strings.Join([]string{cmd.Use, doctl.ArgInterconnectAttachmentType}, "."), cmd.PersistentFlags().Lookup(doctl.ArgInterconnectAttachmentType))

	cmd.AddCommand(PartnerInterconnectAttachments())

	return cmd
}

// PartnerInterconnectAttachments creates the partner interconnect attachment command
func PartnerInterconnectAttachments() *Command {
	cmd := &Command{
		Command: &cobra.Command{
			Use:   "interconnect-attachment",
			Short: "Display commands that manage Partner Interconnect Attachments",
			Long: `The commands under ` + "`" + `doctl partner interconnect-attachment` + "`" + ` are for managing your Partner Interconnect Attachments.
With the Partner Interconnect Attachments commands, you can get or list, create, update, or delete Partner Interconnect Attachments, and manage their configuration details.`,
		},
	}

	interconnectAttachmentDetails := `
- The Partner Interconnect Attachment ID
- The Partner Interconnect Attachment Name
- The Partner Interconnect Attachment State
- The Partner Interconnect Attachment Connection Bandwidth in Mbps
- The Partner Interconnect Attachment Region
- The Partner Interconnect Attachment NaaS Provider
- The Partner Interconnect Attachment VPC network IDs
- The Partner Interconnect Attachment creation date, in ISO8601 combined date and time format
- The Partner Interconnect Attachment BGP Local ASN
- The Partner Interconnect Attachment BGP Local Router IP
- The Partner Interconnect Attachment BGP Peer ASN
- The Partner Interconnect Attachment BGP Peer Router IP
`

	cmdPartnerIAGet := CmdBuilder(cmd, RunPartnerInterconnectAttachmentGet, "get <interconnect-attachment-id>",
		"Retrieves a Partner Interconnect Attachment", "Retrieves information about a Partner Interconnect Attachment, including:"+interconnectAttachmentDetails, Writer,
		aliasOpt("g"), displayerType(&displayers.PartnerInterconnectAttachment{}))
	cmdPartnerIAGet.Example = `The following example retrieves information about a Partner Interconnect Attachment with the ID ` + "`" + `f81d4fae-7dec-11d0-a765-00a0c91e6bf6` + "`" +
		`: doctl network --type "partner" interconnect-attachment get f81d4fae-7dec-11d0-a765-00a0c91e6bf6`

	cmdPartnerIAList := CmdBuilder(cmd, RunPartnerInterconnectAttachmentList, "list", "List Network Interconnect Attachments", "Retrieves a list of the Network Interconnect Attachments on your account, including the following information for each:"+interconnectAttachmentDetails, Writer,
		aliasOpt("ls"), displayerType(&displayers.PartnerInterconnectAttachment{}))
	cmdPartnerIAList.Example = `The following example lists the Network Interconnect Attachments on your account :" + 
		" doctl network --type "partner" interconnect-attachment list --format Name,VPCIDs`

	cmdPartnerIADelete := CmdBuilder(cmd, RunPartnerInterconnectAttachmentDelete, "delete <interconnect-attachment-id>",
		"Deletes a Partner Interconnect Attachment", "Deletes information about a Partner Interconnect Attachment. This is irreversible ", Writer,
		aliasOpt("rm"), displayerType(&displayers.PartnerInterconnectAttachment{}))
	cmdPartnerIADelete.Example = `The following example deletes a Partner Interconnect Attachment with the ID ` + "`" + `f81d4fae-7dec-11d0-a765-00a0c91e6bf6` + "`" +
		`: doctl network --type "partner" interconnect-attachment delete f81d4fae-7dec-11d0-a765-00a0c91e6bf6`

	cmdPartnerIAUpdate := CmdBuilder(cmd, RunPartnerInterconnectAttachmentUpdate, "update <interconnect-attachment-id>",
		"Update a Partner Interconnect Attachment's name and configuration", `Use this command to update the name and and configuration of a Partner Interconnect Attachment`, Writer, aliasOpt("u"))
	AddStringFlag(cmdPartnerIAUpdate, doctl.ArgPartnerInterconnectAttachmentName, "", "",
		"The Partner Interconnect Attachment's name", requiredOpt())
	AddStringFlag(cmdPartnerIAUpdate, doctl.ArgPartnerInterconnectAttachmentVPCIDs, "", "",
		"The Partner Interconnect Attachment's vpc ids", requiredOpt())
	cmdPartnerIAUpdate.Example = `The following example updates the name of a Partner Interconnect Attachment with the ID ` +
		"`" + `f81d4fae-7dec-11d0-a765-00a0c91e6bf6` + "`" + ` to ` + "`" + `new-name` + "`" +
		`: doctl network --type "partner" interconnect-attachment update f81d4fae-7dec-11d0-a765-00a0c91e6bf6 --name new-name`

	return cmd
}

func ensurePartnerAttachmentType(c *CmdConfig) error {
	attachmentType, err := c.Doit.GetString("network", doctl.ArgInterconnectAttachmentType)
	if err != nil {
		return err
	}
	if attachmentType != "partner" {
		return fmt.Errorf("unsupported attachment type: %s", attachmentType)
	}
	return nil
}

// RunPartnerInterconnectAttachmentGet retrieves an existing Partner Interconnect Attachment by its identifier.
func RunPartnerInterconnectAttachmentGet(c *CmdConfig) error {

	if err := ensurePartnerAttachmentType(c); err != nil {
		return err
	}

	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	iaID := c.Args[0]

	interconnectAttachment, err := c.VPCs().GetPartnerInterconnectAttachment(iaID)
	if err != nil {
		return err
	}

	item := &displayers.PartnerInterconnectAttachment{
		PartnerInterconnectAttachments: do.PartnerInterconnectAttachments{*interconnectAttachment},
	}
	return c.Display(item)
}

// RunPartnerInterconnectAttachmentList lists Partner Interconnect Attachment
func RunPartnerInterconnectAttachmentList(c *CmdConfig) error {

	if err := ensurePartnerAttachmentType(c); err != nil {
		return err
	}

	list, err := c.VPCs().ListPartnerInterconnectAttachments()
	if err != nil {
		return err
	}

	item := &displayers.PartnerInterconnectAttachment{PartnerInterconnectAttachments: list}
	return c.Display(item)
}

// RunPartnerInterconnectAttachmentUpdate updates an existing Partner Interconnect Attachment with new configuration.
func RunPartnerInterconnectAttachmentUpdate(c *CmdConfig) error {
	if err := ensurePartnerAttachmentType(c); err != nil {
		return err
	}

	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	peeringID := c.Args[0]

	r := new(godo.PartnerInterconnectAttachmentUpdateRequest)
	name, err := c.Doit.GetString(c.NS, doctl.ArgPartnerInterconnectAttachmentName)
	if err != nil {
		return err
	}
	r.Name = name

	vpcIDs, err := c.Doit.GetString(c.NS, doctl.ArgPartnerInterconnectAttachmentVPCIDs)
	if err != nil {
		return err
	}
	r.VPCIDs = strings.Split(vpcIDs, ",")

	interconnectAttachment, err := c.VPCs().UpdatePartnerInterconnectAttachment(peeringID, r)
	if err != nil {
		return err
	}

	item := &displayers.PartnerInterconnectAttachment{
		PartnerInterconnectAttachments: do.PartnerInterconnectAttachments{*interconnectAttachment},
	}
	return c.Display(item)
}

// RunPartnerInterconnectAttachmentDelete deletes an existing Partner Interconnect Attachment by its identifier.
func RunPartnerInterconnectAttachmentDelete(c *CmdConfig) error {

	if err := ensurePartnerAttachmentType(c); err != nil {
		return err
	}

	err := ensureOneArg(c)
	if err != nil {
		return err
	}
	iaID := c.Args[0]

	force, err := c.Doit.GetBool(c.NS, doctl.ArgForce)
	if err != nil {
		return err
	}

	if force || AskForConfirmDelete("Partner Interconnect Attachment", 1) == nil {

		vpcs := c.VPCs()
		err := vpcs.DeletePartnerInterconnectAttachment(iaID)
		if err != nil {
			return err
		}

		wait, err := c.Doit.GetBool(c.NS, doctl.ArgCommandWait)
		if err != nil {
			return err
		}

		if wait {
			notice("Partner Interconnect Attachment is in progress, waiting for Partner Interconnect Attachment to be deleted")

			err := waitForPIA(vpcs, iaID, "DELETED", true)
			if err != nil {
				return fmt.Errorf("Partner Interconnect Attachment couldn't be deleted : %v", err)
			}
			notice("Partner Interconnect Attachment is successfully deleted")
		} else {
			notice("Partner Interconnect Attachment deletion request accepted")
		}

	} else {
		return fmt.Errorf("operation aborted")
	}

	return nil
}

func waitForPIA(vpcService do.VPCsService, iaID string, wantStatus string, terminateOnNotFound bool) error {
	const maxAttempts = 360
	const errStatus = "ERROR"
	attempts := 0
	printNewLineSet := false

	for i := 0; i < maxAttempts; i++ {
		if attempts != 0 {
			fmt.Fprint(os.Stderr, ".")
			if !printNewLineSet {
				printNewLineSet = true
				defer fmt.Fprintln(os.Stderr)
			}
		}

		interconnectAttachment, err := vpcService.GetPartnerInterconnectAttachment(iaID)
		if err != nil {
			if terminateOnNotFound && strings.Contains(err.Error(), "not found") {
				return nil
			}
			return err
		}

		if interconnectAttachment.State == errStatus {
			return fmt.Errorf("Partner Interconnect Attachment (%s) entered status `%s`", iaID, errStatus)
		}

		if interconnectAttachment.State == wantStatus {
			return nil
		}

		attempts++
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("timeout waiting for Partner Interconnect Attachment (%s) to become %s", iaID, wantStatus)
}
