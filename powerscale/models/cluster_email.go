/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// ClusterEmail struct for ClusterEmail resource
type ClusterEmail struct {
	ID types.String `tfsdk:"id"`
	// V1ClusterEmailSettings
	Settings ClusterEmailSettings `tfsdk:"settings"`
}

// ClusterEmailSettings Cluster email notification settings.
type ClusterEmailSettings struct {
	// This setting determines how notifications will be batched together to be sent by email.  'none' means each notification will be sent separately.  'severity' means notifications of the same severity will be sent together.  'category' means notifications of the same category will be sent together.  'all' means all notifications will be batched together and sent in a single email.
	BatchMode types.String `tfsdk:"batch_mode"`
	// The address of the SMTP server to be used for relaying the notification messages.  An SMTP server is required in order to send notifications.  If this string is empty, no emails will be sent.
	MailRelay types.String `tfsdk:"mail_relay"`
	// The full email address that will appear as the sender of notification messages.
	MailSender types.String `tfsdk:"mail_sender"`
	// The subject line for notification messages from this cluster.
	MailSubject types.String `tfsdk:"mail_subject"`
	// Indicates if an SMTP authentication password is set.
	SMTPAuthPasswdSet types.Bool `tfsdk:"smtp_auth_passwd_set"`
	// Password to authenticate with if SMTP authentication is being used.
	SMTPAuthPasswd types.String `tfsdk:"smtp_auth_passwd"`
	// The type of secure communication protocol to use if SMTP is being used.  If 'none', plain text will be used, if 'starttls', the encrypted STARTTLS protocol will be used.
	SMTPAuthSecurity types.String `tfsdk:"smtp_auth_security"`
	// Username to authenticate with if SMTP authentication is being used.
	SMTPAuthUsername types.String `tfsdk:"smtp_auth_username"`
	// The port on the SMTP server to be used for relaying the notification messages.
	SMTPPort types.Int64 `tfsdk:"smtp_port"`
	// If true, this cluster will send SMTP authentication credentials to the SMTP relay server in order to send its notification emails.  If false, the cluster will attempt to send its notification emails without authentication.
	UseSMTPAuth types.Bool `tfsdk:"use_smtp_auth"`
	// Location of a custom template file that can be used to specify the layout of the notification emails.
	UserTemplate types.String `tfsdk:"user_template"`
}

// ClusterEmailDataSource struct for ClusterEmail data source
type ClusterEmailDataSource struct {
	ID types.String `tfsdk:"id"`
	// V1ClusterEmailSettings
	Settings V1ClusterEmailSettings `tfsdk:"settings"`
}

// V1ClusterEmailSettings Cluster email notification settings for data source.
type V1ClusterEmailSettings struct {
	// This setting determines how notifications will be batched together to be sent by email.  'none' means each notification will be sent separately.  'severity' means notifications of the same severity will be sent together.  'category' means notifications of the same category will be sent together.  'all' means all notifications will be batched together and sent in a single email.
	BatchMode types.String `tfsdk:"batch_mode"`
	// The address of the SMTP server to be used for relaying the notification messages.  An SMTP server is required in order to send notifications.  If this string is empty, no emails will be sent.
	MailRelay types.String `tfsdk:"mail_relay"`
	// The full email address that will appear as the sender of notification messages.
	MailSender types.String `tfsdk:"mail_sender"`
	// The subject line for notification messages from this cluster.
	MailSubject types.String `tfsdk:"mail_subject"`
	// Indicates if an SMTP authentication password is set.
	SMTPAuthPasswdSet types.Bool `tfsdk:"smtp_auth_passwd_set"`
	// The type of secure communication protocol to use if SMTP is being used.  If 'none', plain text will be used, if 'starttls', the encrypted STARTTLS protocol will be used.
	SMTPAuthSecurity types.String `tfsdk:"smtp_auth_security"`
	// Username to authenticate with if SMTP authentication is being used.
	SMTPAuthUsername types.String `tfsdk:"smtp_auth_username"`
	// The port on the SMTP server to be used for relaying the notification messages.
	SMTPPort types.Int64 `tfsdk:"smtp_port"`
	// If true, this cluster will send SMTP authentication credentials to the SMTP relay server in order to send its notification emails.  If false, the cluster will attempt to send its notification emails without authentication.
	UseSMTPAuth types.Bool `tfsdk:"use_smtp_auth"`
	// Location of a custom template file that can be used to specify the layout of the notification emails.
	UserTemplate types.String `tfsdk:"user_template"`
}
