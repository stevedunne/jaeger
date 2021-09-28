
# to run this script use
#   .\build.ps1
#
# then to push to harbor
#   docker login harbor.dell.com
#   docker tag jaeger-agent-dell harbor.dell.com/distributed-tracing/jaeger-agent:1.25.2.dell2
#   docker push harbor.dell.com/distributed-tracing/jaeger-agent:1.25.2.dell2
#

$versionLabel = "github.com/jaegertracing/jaeger/pkg/version.latestVersion"
$versionValue = "1.25.2.dell6"
$versionDateLabel = "github.com/jaegertracing/jaeger/pkg/version.date" 
$versionDateValue = get-date -Format "yyyy-MM-ddTHH:mm:ssZ"

$ldflags = " -X $versionLabel=$versionValue -X $versionDateLabel=$versionDateValue"

# write-host "Building Agent $versionValue"
# go build -o $("jaeger-agent-$versionValue-win.exe") -ldflags $ldflags .\cmd\agent\main.go

write-host "Building Collector $versionValue"
go build -o $("jaeger-collector-$versionValue-win.exe") -ldflags $ldflags .\cmd\collector\main.go

# write-host "Building Ingester $versionValue"
# go build -o $("jaeger-ingester-$versionValue-win.exe") -ldflags $ldflags .\cmd\ingester\main.go

# #https://levelup.gitconnected.com/complete-guide-to-create-docker-container-for-your-golang-application-80f3fb59a15e

#docker build . --build-arg ldFlags=$ldflags --no-cache --progress plain -t jaeger-dell:$versionValue

#$containerId = docker run -d -p 1000:6831 -p 1001:14271 jaeger-dell

# #copy the file off the container to local
#docker cp "$($containerId):/jaeger-agent" "$($pwd)\jaeger-agent-$($versionValue)"
# docker cp "$($containerId):/jaeger-collector" "$($pwd)\jaeger-collector-$($versionValue)"
# docker cp "$($containerId):/jaeger-ingester" "$($pwd)\jaeger-ingester-$($versionValue)"

#docker container stop $containerId
#docker container rm $containerId
