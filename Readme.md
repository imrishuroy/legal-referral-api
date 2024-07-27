### LegalReferral BE

create a new db migration

    make new_migration name=otp_schema

https://github.com/sqlc-dev/sqlc/issues/1062

## Docker Image

Build docker image
    docker build -t legalreferral:latest .

Remove docker images
    docker rmi image-name

Run docker image
    docker run --name legal-referral-api -p 8080:8080 legalreferral:latest

Run Docker image in detached mode
    docker run -d --name legal-referral-api -p 8080:8080 legalreferral:latest

Stop docker container
    docker stop legal-referral-api

Remove docker container
    docker rm legal-referral-api

Remove all docker containers
    docker rm $(docker ps -a -q)

Remove all docker images
    docker rmi $(docker images -q)

Run Docker compose 
    docker compose up

AWS

Get Secret from AWS Secrets Manager
    aws secretsmanager get-secret-value --secret-id legalreferral --query SecretString --output text

Get Secret from AWS Secrets Manager and transform into app.env format
aws secretsmanager get-secret-value --secret-id legalreferral --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env

Get Secret from AWS Secrets Manager and transform into json format
aws secretsmanager get-secret-value --secret-id legalreferral-service-account-key --query SecretString --output text> service-account-key.json

AWS
 
    sudo yum install git -y
    git clone your-repo & git checkout your-branch
    sudo yum install -y docker
    sudo service docker start
    sudo usermod -a -G docker ec2-user
    sudo df -h ( check disk space )

    zip files.zip *


https://www.youtube.com/watch?v=C_QzIpPsexs&t=372s


![img.png](img.png)

Deploy latest changes to elasticbeanstalk ( even if you have not committed the changes )

    eb deploy
eb deploy --staged

