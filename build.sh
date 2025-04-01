#!/usr/bin/env bash

Usage="Usage: build.sh -t [build mode: binary|docker] -o [binary|docker name] -v [binary|docker version] -m [dev|test|prod] [-e]"
envSetup="false"
envMode="prod"
buildMode=""
repoName=""
tagVer=""
dstName="dist"

while getopts "t:o:v:me" opt
do
  case $opt in
    t)
      buildMode=$OPTARG
      ;;
    o)
      repoName=$OPTARG
      ;;
    v)
      tagVer=$OPTARG
      ;;
    m)
      envMode=$OPTARG
      ;;
    e)
      envSetup="true"
      ;;
    ?)
      echo $Usage
      exit 1
      ;;
  esac
done

shift $(($OPTIND - 1))
if [[ -z ${buildMode} ]] || [[ -z ${repoName} ]] || [[ -z ${tagVer} ]]; then
  echo $Usage
  exit 1
fi

env_install(){
  if [ "$envSetup" = "true" ]; then
    sysUbuntu=`cat /etc/issue | grep "Ubuntu"`
    if [[ -z ${sysUbuntu} ]]; then
      # only centos7
      sysctos=`cat /etc/redhat-release | grep "CentOS Linux release 7"`
      if [[ -z ${sysCentOS} ]]; then
        echo "System unsupported, your must install docker and compose by yourself"
        exit 1
      else
        dockerPkg=$(rpm -qa | grep docker)
        if [[ -z ${dockerPkg} ]]; then
          curl -fsSL https://get.docker.com/ | sh
          systemctl enable docker
        fi
      fi
    else
      # ubuntu 20.04 or later
      dockerPkg=$(dpkg -l | grep docker.io)
      if [[ -z ${dockerPkg} ]]; then
        apt update
        apt upgrade
        apt install docker.io docker-compose -y
      fi
    fi
  fi
}

binary_build() {
  echo "building target: binary"
  rm -rf ${dstName}
  mkdir -p ${dstName}
  go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOARCH=amd64 \
    && go mod tidy \
    && go build -ldflags="-s -w" -o ${dstName}/${repoName}
  rm -f ${dstName}/VERSION
  touch ${dstName}/VERSION
  cp config.${envMode}.yaml ${dstName}/
  cp run.dist.sh ${dstName}/run.sh
  cp -R file ${dstName}/
  echo "$tagVer" >> ${dstName}/VERSION
}

docker_build() {
  echo "building target: docker"
  docker build -t ${repoName}:${tagVer} .
  noneRepo=`docker images | grep "<none>" | awk ' {print $3}'`
  if [[ -n ${noneRepo} ]]; then
    docker rmi ${noneRepo}
  fi
  rm -rf ${dstName}
  mkdir -p ${dstName}
  cd ${dstName}
  docker save ${repoName}:${tagVer} > ${repoName}_v${tagVer}.tar
  cd ..
  docker rmi `docker images | grep ${repoName} | awk ' {print $3}'`
}

env_install
if [ {$buildMode} = "docker" ]; then
  docker_build
elif [ ${buildMode} = "binary" ]; then
  binary_build
else
  echo "unknown build mode: ${buildMode}"
fi