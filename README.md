# AWS-CLI-AUTH

aws-cli-auth is cli tool for authentication an AWS IAM User with MFA to assume a role when using AWS CLI on your machine.

## How to use
1. Configure AWS iam user(s) and role(s) according to aws best practices
1. clone this repo and create a config file, file type can be `json`, `yaml` or `toml`
1. The `config` file should look like this:
    #### **`config.yaml`**
    ```yaml 
    User:
      AccKeyId: "<IAM USER ACCESS KEY ID>"
      SecAccKey: "<IAM USER SECRET ACCESS KEY>"
    DefaultRegion: "<AWS REGION>"
    MFASerial: "<MFA SERIAL ARN>"
    RoleArn: "<ROLE-ARN>"
    SessionName: "<SESSION-NAME>"
    ```

    #### **`config.toml`**
    ```toml
    DefaultRegion = "<AWS REGION>"
    MFASerial = "<MFA SERIAL ARN>"
    RoleArn = "<AWS REGION>"
    SessionName = "<SESSION-NAME>"

    [User]
    AccKeyId = "<IAM USER ACCESS KEY ID>"
    SecAccKey = "<IAM USER SECRET ACCESS KEY>"
    ```

    #### **`config.json`**
    ```json
    {
        "DefaultRegion": "<AWS REGION>",
        "MFASerial": "<MFA SERIAL ARN>",
        "RoleArn": "<AWS REGION>",
        "SessionName": "<SESSION-NAME>",
        "User":{
            "AccKeyId": "<IAM USER ACCESS KEY ID>",
            "SecAccKey": "<IAM USER SECRET ACCESS KEY>"
        }
    }
    ```

1. run `go build .`
1. run `./aws-cli-auth -h` to see the help
1. run `./aws-cli-auth --config=config.[yaml|toml|json]` to request temporary AWS credentials


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


## Upcoming feature
- beautification