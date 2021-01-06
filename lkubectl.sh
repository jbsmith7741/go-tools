#!/bin/bash
function lkubectl {
  env=${e:=$env}
  # set the kube config file
  file="dev.kube"
  if [ "$env" != "" ]; then
      file="${env}.kube"
  fi
  if [ "$f" != "" ]; then
      file="$f"
  fi
  echo "file:$file"
  kubectl --kubeconfig="$file" $@
  reset
}

function reset {
  e=""
  env=""
}

function kubectx {
  env=${e:=$env}
  # set the kube config file
  file="dev.kube"
  if [ "$env" != "" ]; then
      file="${env}.kube"
  fi
  if [ "$f" != "" ]; then
      file="$f"
  fi

  # set cluster by regex match found in $1
  clusters=$(kubectl --kubeconfig=$file config get-contexts -o name)
  # for c in $clusters; do # bash
  for c in ${(f)clusters}; do # zsh fix to split by newlines
      if [[ "$1" == "list" ]]; then
          echo "$c"
      elif [[ $c =~ $1 ]]; then
          echo "cluster set to $c"
          kubectl --kubeconfig=$file config use-context $c
      fi
  done
  reset
}