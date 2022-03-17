# AWS-CLI-AUTH

aws-cli-auth is cli tool for authentication an AWS IAM User with MFA to assume a role when using AWS CLI on your machine.

## How to use
1. Create a `config.yaml` file in /config folder
1. Inside the `config.yaml` file, add the following:
    ```yaml
    DefaultRegion: "<AWS REGION>"
    MFASerial: "<MFA SERIAL ARN>"
    RoleArn: "<ROLE-ARN>"
    SessionName: "<SESSION-NAME>"
    ```

## AWS IAM user and role creation

AWS security best practices recommends enabling MFA for AWS account and using roles to grant limited access to resources for a limited amount of time.
In keeping with these security practices, I recommend:
- Create an AWS role and attach all the policies needed for that role
- Create an assume role policy for assuming that role
- Creating an AWS user with [MFA enabled](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html#enable-virt-mfa-for-iam-user) and attach an assumerole policy to that user
- Add the following trust relationship to the AWS IAM role, this trust relationship will only allow the user to assume role if the MFA code is provided
    ```json
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Principal": {
                    "AWS": "<User ARN>"
                },
                "Action": "sts:AssumeRole",
                "Condition": {
                    "Bool": {
                        "aws:MultiFactorAuthPresent": "true"
                    }
                }
            }
        ]
    }
    ```
