apps=$(curl -H "Authorization: token $DEIS_TOKEN" http://$DEIS_SERVER/v1/apps  | jq -r '.results | map(.id) | join(" ")')

for app in $apps; do
  echo "Resetting containers of $app"

  originalscale=$(curl -H "Authorization: token $DEIS_TOKEN" http://$DEIS_SERVER/v1/apps/$app/containers/ 2>/dev/null | jq -r '(.results) | [group_by(.type)[] | max_by(.num)] | [map(.type), map(.num)] | transpose | map([.[0], .[1] | tostring] | join("=")) | join(" ")')
  zeroscale=$(curl -H "Authorization: token $DEIS_TOKEN" http://$DEIS_SERVER/v1/apps/$app/containers/  2>/dev/null | jq -r '(.results) | unique_by(.type) | map([.type, "0"] | join("=")) | join(" ")')

  deis ps:scale $zeroscale -a $app
  deis ps:scale $originalscale -a $app
  echo
done
