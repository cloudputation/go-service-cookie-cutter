# # !/bin/bash
#
# function auto_push {
#   API_VERSION=$(cat API_VERSION)
#
#
#   first_digit=$(echo ${API_VERSION} | cut -d"." -f1)
#   second_digit=$(echo ${API_VERSION} | cut -d"." -f2)
#   third_digit=$(echo ${API_VERSION} | cut -d"." -f3)
#
#   let "third_digit++"
#
#   new_API_VERSION="${first_digit}.${second_digit}.${third_digit}"
#
#   echo current api version is: $API_VERSION
#   echo $new_API_VERSION | tee  API_VERSION
#
#
#
#   git add .
#   git commit -m "auto commit"
#   git push
#
# }
#
# grep "production = true" GIT_CONTROLS/auto_push && auto_push || echo "Production deactivated."
#
# function sync_to_stage {
#   go mod tidy
#   rsync -a -P ./* devops@tool:~/dev/iterator-test/
#   rsync -a -P $SAVED_FILE devops@tool:~/dev/iterator-test/$SAVED_FILE
#   rsync -a -P ./go.mod devops@tool:~/dev/iterator-test/
#   rsync -a -P ./go.sum devops@tool:~/dev/iterator-test/
#   rsync -a -P ./.air.toml devops@tool:~/dev/iterator-test/
#   rsync -a -P ./.release/ devops@tool:~/dev/iterator-test/.release
#   ssh devops@tool "sudo cp /home/devops/dev/iterator-test/.release/defaults/test.config.hcl /etc/iterator/config.hcl &&
#                     cp -r /home/devops/dev/iterator-test/tests/env/terraform/* /opt/iterator/data/terraform-data"
#   # cp .release/defaults/test.config.hcl /etc/iterator/config.hcl
# }
#
# grep "staging = true" GIT_CONTROLS/auto_push && sync_to_stage || echo "Staging deactivated."
