url=235494794123.dkr.ecr.ap-northeast-2.amazonaws.com/prod-tmr686-backend-ecr

aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.ap-northeast-2.amazonaws.com"
docker pull $url:latest
docker stop tmr || true
docker rm tmr || true

# EC2에 temp_files 디렉토리가 없는 경우에만 생성
if [ ! -d "/home/ec2-user/tmr/temp_files" ]; then
    mkdir -p /home/ec2-user/tmr/temp_files
    chmod 777 /home/ec2-user/tmr/temp_files
fi

docker run -d \
  --name tmr \
  -p 80:8080 \
  -v /home/ec2-user/tmr/temp_files:/temp_files \
  $url:latest
