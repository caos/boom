set -e

function parse_yaml {
   local prefix=""
   local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034')
   sed -ne "s|^\($s\):|\1|" \
        -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
        -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p"  $1 |
   awk -F$fs '{
      indent = length($1)/2;
      vname[indent] = $2;
      for (i in vname) {if (i > indent) {delete vname[i]}}
      if (length($3) > 0) {
         vn=""; for (i=0; i<indent; i++) {vn=(vn)(vname[i])("_")}
         printf("%s%s%s=\"%s\"\n", "'$prefix'",vn, $2, $3);
      }
   }'
}

TOOLS_HOME="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

function doHelm {
  helm --home ${helmHome} $@
}

helmHome=${TOOLS_HOME}/helm
chartHome=${TOOLS_HOME}/charts
indexProtocol=https

doHelm init --client-only >& /dev/null

for file in toolsets/$1/*; do
   eval $(parse_yaml $file)

   if [ -z "$chart_index_name" ]; then
      indexName=stable
   else
      doHelm repo add ${chart_index_name} ${indexProtocol}://${chart_index_url} >& /dev/null
      indexName=$chart_index_name
   fi

   echo "fetching chart: $chart_name-$chart_version"
   doHelm repo update >& /dev/null
   doHelm fetch --untar --version=${chart_version} --untardir ${chartHome} ${indexName}/${chart_name} >& /dev/null
   echo "done fetching $chart_name-$chart_version"

   chart_index_name=""
   chart_index_url=""
done