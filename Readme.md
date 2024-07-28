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

    aws secretsmanager get-secret-value --secret-id legalreferral-env --query SecretString --output text

Get Secret from AWS Secrets Manager and transform into app.env format

    aws secretsmanager get-secret-value --secret-id legalreferral --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env

Get Secret from AWS Secrets Manager and transform into json format

    aws secretsmanager get-secret-value --secret-id legalreferral-service-account-key --query SecretString --output text> service-account-key.json

Login to AWS ECR
https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login-password.html

    aws ecr get-login-password | docker login --username AWS --password-stdin 533267397749.dkr.ecr.ap-south-1.amazonaws.com

Docker Pull image

    docker pull 533267397749.dkr.ecr.ap-south-1.amazonaws.com/legal-referral:f4a3c4d6852f405eb0243c8f00808370ed986e6c

Run the image

    docker run -p 8080:8080 533267397749.dkr.ecr.ap-south-1.amazonaws.com/legal-referral:f4a3c4d6852f405eb0243c8f00808370ed986e6c

k8s Access Cluster

    aws eks update-kubeconfig --name legalreferral --region ap-south-1

Get Current User Identity

    aws sts get-caller-identity

AWS change profile

    export AWS_PROFILE=root

Get cluster info
    
    kubectl cluster-info

To apply the new configmap to the RBAC configuration of the cluster, run the following command

    kubectl apply -f eks/aws-auth.yaml

Get service

    kubectl get service

Get pods

    kubectl get pods

Apply Deployment

    kubectl apply -f eks/deployment.yaml

Apply Service
    
    kubectl apply -f eks/service.yaml

nslookup

    nslookup a05b7f12387364b1ab93c06f36486f89-204541799.ap-south-1.elb.amazonaws.com

https://repost.aws/knowledge-center/amazon-eks-cluster-access


Deploy to EC2 AWS
 
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

