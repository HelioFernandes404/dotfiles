Update webhook-glpi deployment image tag on k3s cluster

Usage: `/update-webhook <ssh-machine> <image-tag>`

Example: `/update-webhook thinkpad-dev 1.0.0-feat_logging-improvements-31a4183`

Extract the SSH machine and image tag from the user's command, then run:

```bash
ssh <ssh-machine> "kubectl set image deployment/webhook-glpi webhook-glpi=730335278395.dkr.ecr.us-east-1.amazonaws.com/webhooks-alertmanager:<image-tag> -n monitoring"
```

Just run the command and report done. No verification, no status checks.
