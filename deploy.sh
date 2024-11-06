url=235494794123.dkr.ecr.ap-northeast-2.amazonaws.com/prod-tmr686-backend-ecr

aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-2.amazonaws.com"
docker pull $url:latest
docker stop tmr || true
docker rm tmr || true
docker run -d --name tmr -p 80:8080 $url:latest
