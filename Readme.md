### LegalReferral BE

create a new db migration

    make new_migration name=otp_schema



https://github.com/sqlc-dev/sqlc/issues/1062

## Docker Image

    docker build -t legal-referral-api .

    docker run -p 8080:8080 legal-referral-api

AWS
 
    sudo yum install git -y
    git clone your-repo & git checkout your-branch
    sudo yum install -y docker
    sudo service docker start
    sudo usermod -a -G docker ec2-user

