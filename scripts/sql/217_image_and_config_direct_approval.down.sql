/*
 * Copyright (c) 2024. Devtron Inc.
 */

UPDATE "public"."notification_templates" SET template_payload = '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "🛎️ Image approval requested | Application > {{appName}} | Environment > {{envName}}","html": "<table style=\"width: 600px; height: 485px;  border-collapse: collapse; padding: 20px;\"><tr style=\"background-color:#E5F2FF;\"><td colspan=\"2\" style=\"padding-left:16px;\"><h2 style=\"color:#000A14;\">Image approval request</h2><span>{{eventTime}}</span><br><span>by <strong style=\"color:#0066CC;\">{{triggeredBy}}</strong></span><br><br>{{#imageApprovalLink}}<a href=\"{{&imageApprovalLink}}\" style=\" height: 32px; padding: 7px 12px; line-height: 32px; font-size: 12px; font-weight: 600; border-radius: 4px; text-decoration: none; outline: none; min-width: 64px; text-transform: capitalize; text-align: center; background: #0066CC; color: #fff; border: 1px solid transparent; cursor: pointer;\">View request</a><br><br>{{/imageApprovalLink}}</td><td style=\"text-align: right;\"><img src=\"https://cdn.devtron.ai/images/img_build_notification.png\" style=\"height: 72px; width: 72px;\"></td></tr><tr><td colspan=\"3\"><hr><br><span>Application: <strong>{{appName}}</strong></span>&nbsp;&nbsp;|&nbsp;&nbsp;<span>Environment: <strong>{{envName}}</strong></span><br></span><br><br><hr><h3>Image Details</h3><span>Image tag <br><strong>{{imageTag}}</strong></span></td></tr><tr><td colspan=\"3\"><br><br><span style=\"display: {{commentDisplayStyle}}\">Comment<br><strong>{{comment}}</strong></span></td></tr><tr><td colspan=\"3\"><br><br><span style=\"display: {{tagDisplayStyle}}\">Tags<br><strong>{{tags}}</strong></span><br></td></tr></table>"}' WHERE event_type_id=4;
UPDATE "public"."notification_templates" SET template_payload = '{"from": "{{fromEmail}}", "to": "{{toEmail}}","subject": "🛎️ Config Change approval requested | Application > {{appName}} | Change Impacts > {{envName}}","html": "<table cellpadding=\"0\" style=\"font-family:Arial,Verdana,Helvetica;width:600px;height:485px;border-collapse:inherit;border-spacing:0;border:1px solid #d0d4d9;border-radius:8px;padding:16px 20px;margin:20px auto;box-shadow:0 0 8px 0 rgba(0,0,0,.1)\"><tr><td colspan=\"3\"><div style=\"height:28px;padding-bottom:16px;margin-bottom:20px;border-bottom:1px solid #edf1f5;max-width:600px\"><img style=\"height:100%\" src=\"https://devtron-public-asset.s3.us-east-2.amazonaws.com/images/devtron/devtron-logo.png\" alt=\"ci-triggered\"></div></td></tr><tr><td colspan=\"3\"><div style=\"background-color:#e5f2ff;border-top-left-radius:8px;border-top-right-radius:8px;padding:20px 20px 16px 20px;display:flex\"><div style=\"width:90%\"><div style=\"font-size:16px;line-height:24px;font-weight:600;margin-bottom:6px;color:#000a14\">Config Change approval request</div><span style=\"font-size:14px;line-height:20px;color:#000a14\">{{eventTime}}</span><br><div><span style=\"font-size:14px;line-height:20px;color:#000a14\">by</span><strong style=\"font-size:14px;line-height:20px;color:#06c;margin-left:4px\">{{triggeredBy}}</strong></div></div><div><img src=\"https://cdn.devtron.ai/images/img_build_notification.png\" style=\"height:72px;width:72px\"></div></div></td></tr><tr><td colspan=\"3\"><div style=\"background-color:#e5f2ff;border-bottom-left-radius:8px;border-bottom-right-radius:8px;padding:0 0 20px 20px\">{{#protectConfigLink}}<a href=\"{{&protectConfigLink}}\" style=\"height:32px;padding:7px 12px;line-height:32px;font-size:12px;font-weight:600;border-radius:4px;text-decoration:none;outline:0;min-width:64px;text-transform:capitalize;text-align:center;background:#06c;color:#fff;border:1px solid transparent;cursor:pointer\">View request</a>{{/protectConfigLink}}</div></td></tr><tr><tr><td><br></td></tr><td><div style=\"color:#3b444c;font-size:12px;padding-bottom:4px\">Application</div></td><td colspan=\"2\"><div style=\"color:#3b444c;font-size:12px;padding-bottom:4px\">Change Impacts</div></td><tr><td><div style=\"color:#000a14;font-size:13px\">{{appName}}</div></td><td><div style=\"color:#000a14;font-size:13px\">{{envName}}</div></td></tr><tr><td colspan=\"3\"><div style=\"font-weight:600;margin-top:20px;width:100%;border-top:1px solid #edf1f5;padding:16px 0;font-size:14px\">File Details</div></td></tr><tr><td><div style=\"color:#3b444c;font-size:12px;padding-bottom:4px\">File Type</div></td><td><div style=\"color:#3b444c;font-size:12px;padding-bottom:4px\">File Name</div></td></tr><tr></tr><tr><td><div style=\"color:#000a14;font-size:13px\">{{protectConfigFileType}}</div></td><td><div style=\"color:#000a14;font-size:13px\">{{protectConfigFileName}}</div></td><tr><td><br></td></tr><tr><td colspan=\"3\"><div style=\"display: {{commentDisplayStyle}};color:#3b444c;font-size:12px;padding-bottom:4px\">Comment</div></td></tr><tr><td colspan=\"3\"><div style=\"display: {{commentDisplayStyle}};color:#000a14;font-size:13px\">{{#protectConfigComment}}{{.}}<br>{{/protectConfigComment}}</div></tr><tr><td colspan=\"3\"><div style=\"border-top:1px solid #edf1f5;margin:20px 0 16px 0;height:1px\"></div></td></tr><tr><td colspan=\"2\" style=\"display:flex\"><span><a href=\"https://twitter.com/DevtronL\" target=\"_blank\" style=\"cursor:pointer;text-decoration:none;padding-right:12px;display:flex\"><div><img style=\"width:20px\" src=\"https://cdn.devtron.ai/images/twitter_social_dark.png\"></div></a></span><span><a href=\"https://www.linkedin.com/company/devtron-labs\" target=\"_blank\" style=\"cursor:pointer;text-decoration:none;padding-right:12px;display:flex\"><div><img style=\"width:20px\" src=\"https://cdn.devtron.ai/images/linkedin_social_dark.png\"></div></a></span><span><a href=\"https://devtron.ai/blog/\" target=\"_blank\" style=\"color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline;padding-right:12px\">Blog</a></span><span><a href=\"https://devtron.ai/\" target=\"_blank\" style=\"color:#000a14;font-size:13px;line-height:20px;cursor:pointer;text-decoration:underline\">Website</a></span></td><td colspan=\"2\" style=\"text-align:right\"><div style=\"color:#767d84;font-size:13px;line-height:20px\">&copy; Devtron Labs 2020</div></td></tr></table>"}' WHERE event_type_id=5;